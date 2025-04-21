package api

import (
        "context"
        "encoding/json"
        "fmt"
        "net/http"
        "reflect"
        "sync"
        "time"

        "github.com/davidjeba/goscript/pkg/goscale/db"
)

// GoScaleAPI represents the main API system that combines gRPC-like performance
// with GraphQL-like flexibility
type GoScaleAPI struct {
        resolvers      map[string]Resolver
        middlewares    []Middleware
        subscriptions  map[string]*Subscription
        subMutex       sync.RWMutex
        dbConnection   *db.GoScaleDB
        edgeEnabled    bool
        edgeNodes      []string
        compressionLevel int
        batchSize      int
        timeout        time.Duration
        maxConcurrent  int
        metrics        *Metrics
}

// Resolver is a function that resolves a specific API request
type Resolver func(ctx context.Context, params map[string]interface{}) (interface{}, error)

// Middleware processes requests before they reach resolvers
type Middleware func(ctx context.Context, next Resolver) Resolver

// Subscription represents a real-time data subscription
type Subscription struct {
        topic     string
        clients   map[string]chan interface{}
        mutex     sync.RWMutex
}

// Metrics tracks API performance metrics
type Metrics struct {
        RequestCount      int64
        AvgResponseTime   float64
        ErrorRate         float64
        CacheHitRate      float64
        EdgeRequestCount  int64
        mutex             sync.RWMutex
        clients           map[string]chan interface{}
}

// NewGoScaleAPI creates a new instance of the GoScaleAPI
func NewGoScaleAPI(config *Config) *GoScaleAPI {
        if config == nil {
                config = DefaultConfig()
        }
        
        dbConfig := &db.Config{
                ConnectionString: config.DBConnectionString,
                MaxConnections: config.MaxDBConnections,
                QueryTimeout: config.Timeout,
                EnableTimeSeries: config.EnableTimeSeries,
                EnableRelationships: config.EnableRelationships,
                EnableNoCode: config.EnableNoCode,
        }
        
        return &GoScaleAPI{
                resolvers:      make(map[string]Resolver),
                middlewares:    []Middleware{},
                subscriptions:  make(map[string]*Subscription),
                dbConnection:   db.NewGoScaleDB(dbConfig),
                edgeEnabled:    config.EdgeEnabled,
                edgeNodes:      config.EdgeNodes,
                compressionLevel: config.CompressionLevel,
                batchSize:      config.BatchSize,
                timeout:        config.Timeout,
                maxConcurrent:  config.MaxConcurrent,
                metrics:        &Metrics{
                        clients: make(map[string]chan interface{}),
                },
        }
}

// Config contains configuration options for GoScaleAPI
type Config struct {
        DBConnectionString string
        MaxDBConnections   int
        EdgeEnabled        bool
        EdgeNodes          []string
        CompressionLevel   int
        BatchSize          int
        Timeout            time.Duration
        MaxConcurrent      int
        EnableTimeSeries   bool
        EnableRelationships bool
        EnableNoCode       bool
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
        return &Config{
                DBConnectionString: "localhost:5432",
                MaxDBConnections:   100,
                EdgeEnabled:        false,
                EdgeNodes:          []string{},
                CompressionLevel:   5,
                BatchSize:          100,
                Timeout:            time.Second * 30,
                MaxConcurrent:      1000,
                EnableTimeSeries:   true,
                EnableRelationships: true,
                EnableNoCode:       true,
        }
}

// RegisterResolver registers a new resolver function for a specific path
func (g *GoScaleAPI) RegisterResolver(path string, resolver Resolver) {
        g.resolvers[path] = resolver
}

// Use adds a middleware to the API
func (g *GoScaleAPI) Use(middleware Middleware) {
        g.middlewares = append(g.middlewares, middleware)
}

// CreateSubscription creates a new subscription topic
func (g *GoScaleAPI) CreateSubscription(topic string) *Subscription {
        g.subMutex.Lock()
        defer g.subMutex.Unlock()
        
        sub := &Subscription{
                topic:   topic,
                clients: make(map[string]chan interface{}),
        }
        
        g.subscriptions[topic] = sub
        return sub
}

// Subscribe adds a client to a subscription
func (s *Subscription) Subscribe(clientID string) chan interface{} {
        s.mutex.Lock()
        defer s.mutex.Unlock()
        
        ch := make(chan interface{}, 100)
        s.clients[clientID] = ch
        return ch
}

// Unsubscribe removes a client from a subscription
func (s *Subscription) Unsubscribe(clientID string) {
        s.mutex.Lock()
        defer s.mutex.Unlock()
        
        if ch, ok := s.clients[clientID]; ok {
                close(ch)
                delete(s.clients, clientID)
        }
}

// Publish sends data to all subscribers
func (s *Subscription) Publish(data interface{}) {
        s.mutex.RLock()
        defer s.mutex.RUnlock()
        
        for _, ch := range s.clients {
                select {
                case ch <- data:
                        // Data sent successfully
                default:
                        // Channel buffer is full, skip this message
                }
        }
}

// ServeHTTP implements the http.Handler interface
func (g *GoScaleAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        startTime := time.Now()
        
        // Parse the request
        var request struct {
                Query     string                 `json:"query"`
                Variables map[string]interface{} `json:"variables"`
                Operation string                 `json:"operation"`
        }
        
        if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
        }
        
        // Create context with timeout
        ctx, cancel := context.WithTimeout(r.Context(), g.timeout)
        defer cancel()
        
        // Apply middlewares
        var resolver Resolver
        if r, ok := g.resolvers[request.Operation]; ok {
                resolver = r
        } else {
                http.Error(w, "Unknown operation", http.StatusNotFound)
                return
        }
        
        for i := len(g.middlewares) - 1; i >= 0; i-- {
                resolver = g.middlewares[i](ctx, resolver)
        }
        
        // Execute the resolver
        result, err := resolver(ctx, request.Variables)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                g.updateMetrics(startTime, false)
                return
        }
        
        // Return the result
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
                "data": result,
        })
        
        g.updateMetrics(startTime, true)
}

// updateMetrics updates the API metrics
func (g *GoScaleAPI) updateMetrics(startTime time.Time, success bool) {
        duration := time.Since(startTime).Seconds()
        
        g.metrics.mutex.Lock()
        defer g.metrics.mutex.Unlock()
        
        g.metrics.RequestCount++
        g.metrics.AvgResponseTime = (g.metrics.AvgResponseTime*float64(g.metrics.RequestCount-1) + duration) / float64(g.metrics.RequestCount)
        
        if !success {
                g.metrics.ErrorRate = (g.metrics.ErrorRate*float64(g.metrics.RequestCount-1) + 1) / float64(g.metrics.RequestCount)
        } else {
                g.metrics.ErrorRate = (g.metrics.ErrorRate * float64(g.metrics.RequestCount-1)) / float64(g.metrics.RequestCount)
        }
}

// GetMetrics returns the current API metrics
func (g *GoScaleAPI) GetMetrics() *Metrics {
        g.metrics.mutex.RLock()
        defer g.metrics.mutex.RUnlock()
        
        return &Metrics{
                RequestCount:     g.metrics.RequestCount,
                AvgResponseTime:  g.metrics.AvgResponseTime,
                ErrorRate:        g.metrics.ErrorRate,
                CacheHitRate:     g.metrics.CacheHitRate,
                EdgeRequestCount: g.metrics.EdgeRequestCount,
        }
}

// Schema represents a GraphQL-like schema for the API
type Schema struct {
        Types       map[string]*Type
        Queries     map[string]*Field
        Mutations   map[string]*Field
        Subscriptions map[string]*Field
}

// Type represents a schema type
type Type struct {
        Name        string
        Fields      map[string]*Field
        Implements  []string
        Description string
}

// Field represents a field in a type
type Field struct {
        Name        string
        Type        string
        Args        map[string]*Argument
        Resolver    Resolver
        Description string
}

// Argument represents a field argument
type Argument struct {
        Name        string
        Type        string
        Default     interface{}
        Description string
}

// NewSchema creates a new schema
func NewSchema() *Schema {
        return &Schema{
                Types:       make(map[string]*Type),
                Queries:     make(map[string]*Field),
                Mutations:   make(map[string]*Field),
                Subscriptions: make(map[string]*Field),
        }
}

// AddType adds a type to the schema
func (s *Schema) AddType(name string, description string) *Type {
        t := &Type{
                Name:        name,
                Fields:      make(map[string]*Field),
                Description: description,
        }
        s.Types[name] = t
        return t
}

// AddField adds a field to a type
func (t *Type) AddField(name string, typeName string, description string) *Field {
        f := &Field{
                Name:        name,
                Type:        typeName,
                Args:        make(map[string]*Argument),
                Description: description,
        }
        t.Fields[name] = f
        return f
}

// AddArg adds an argument to a field
func (f *Field) AddArg(name string, typeName string, defaultValue interface{}, description string) *Argument {
        a := &Argument{
                Name:        name,
                Type:        typeName,
                Default:     defaultValue,
                Description: description,
        }
        f.Args[name] = a
        return a
}

// SetResolver sets the resolver for a field
func (f *Field) SetResolver(resolver Resolver) {
        f.Resolver = resolver
}

// AddQuery adds a query to the schema
func (s *Schema) AddQuery(name string, typeName string, description string) *Field {
        f := &Field{
                Name:        name,
                Type:        typeName,
                Args:        make(map[string]*Argument),
                Description: description,
        }
        s.Queries[name] = f
        return f
}

// AddMutation adds a mutation to the schema
func (s *Schema) AddMutation(name string, typeName string, description string) *Field {
        f := &Field{
                Name:        name,
                Type:        typeName,
                Args:        make(map[string]*Argument),
                Description: description,
        }
        s.Mutations[name] = f
        return f
}

// AddSubscription adds a subscription to the schema
func (s *Schema) AddSubscription(name string, typeName string, description string) *Field {
        f := &Field{
                Name:        name,
                Type:        typeName,
                Args:        make(map[string]*Argument),
                Description: description,
        }
        s.Subscriptions[name] = f
        return f
}

// ApplySchema applies a schema to a GoScaleAPI instance
func (g *GoScaleAPI) ApplySchema(schema *Schema) error {
        // Register query resolvers
        for name, field := range schema.Queries {
                if field.Resolver == nil {
                        return fmt.Errorf("query %s has no resolver", name)
                }
                g.RegisterResolver("query:"+name, field.Resolver)
        }
        
        // Register mutation resolvers
        for name, field := range schema.Mutations {
                if field.Resolver == nil {
                        return fmt.Errorf("mutation %s has no resolver", name)
                }
                g.RegisterResolver("mutation:"+name, field.Resolver)
        }
        
        // Register subscription resolvers
        for name, field := range schema.Subscriptions {
                if field.Resolver == nil {
                        return fmt.Errorf("subscription %s has no resolver", name)
                }
                g.RegisterResolver("subscription:"+name, field.Resolver)
                g.CreateSubscription(name)
        }
        
        return nil
}

// GenerateGRPCStubs generates gRPC-compatible stubs for the API
func (g *GoScaleAPI) GenerateGRPCStubs(schema *Schema) (string, error) {
        // This would generate actual gRPC stubs in a real implementation
        return "// Generated gRPC stubs", nil
}

// EnableEdgeComputing enables edge computing for the API
func (g *GoScaleAPI) EnableEdgeComputing(nodes []string) {
        g.edgeEnabled = true
        g.edgeNodes = nodes
}

// DisableEdgeComputing disables edge computing for the API
func (g *GoScaleAPI) DisableEdgeComputing() {
        g.edgeEnabled = false
        g.edgeNodes = []string{}
}

// GetDB returns the database connection
func (g *GoScaleAPI) GetDB() *db.GoScaleDB {
        return g.dbConnection
}

// GetResolvers returns the API resolvers
func (g *GoScaleAPI) GetResolvers() map[string]Resolver {
        return g.resolvers
}

// Close closes the API and all its resources
func (g *GoScaleAPI) Close() error {
        // Close all subscriptions
        g.subMutex.Lock()
        for _, sub := range g.subscriptions {
                sub.mutex.Lock()
                for clientID, ch := range sub.clients {
                        close(ch)
                        delete(sub.clients, clientID)
                }
                sub.mutex.Unlock()
        }
        g.subMutex.Unlock()
        
        // Close the database connection
        return g.dbConnection.Close()
}