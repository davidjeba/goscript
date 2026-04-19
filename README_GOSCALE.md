# GoScale - High-Performance API and Database System

GoScale is a comprehensive API and database system designed for high performance, flexibility, and edge computing capabilities. It combines the speed of gRPC with the flexibility of GraphQL, along with a custom database solution that integrates features from TimescaleDB, NoCode databases, and PostgreSQL.

## Core Components

### GoScaleAPI

GoScaleAPI provides a unified API system that combines:

- **gRPC-like Performance**: Optimized for high throughput and low latency
- **GraphQL-like Flexibility**: Schema-based API with queries, mutations, and subscriptions
- **Edge Computing Support**: Distributed API processing at the network edge
- **Real-time Subscriptions**: WebSocket-based real-time data updates
- **Middleware System**: Extensible request processing pipeline
- **Metrics and Monitoring**: Built-in performance tracking

### GoScaleDB

GoScaleDB is a high-performance database that combines:

- **TimescaleDB Features**: Time-series data management with retention policies
- **NoCode Database**: Schema-less data storage with validation
- **PostgreSQL Compatibility**: SQL query support and relational features
- **Relationship Management**: One-to-one, one-to-many, and many-to-many relationships
- **Query Caching**: Automatic caching of query results
- **Sharding and Replication**: Horizontal scaling with data distribution
- **Metrics and Monitoring**: Built-in performance tracking

### Edge Computing

The edge computing system provides:

- **Low-Latency Processing**: API processing at the network edge
- **Distributed Caching**: Cache data close to users
- **Load Balancing**: Distribute requests across edge nodes
- **Health Monitoring**: Automatic health checks and failover
- **Synchronization**: Keep edge nodes in sync with the central system
- **Metrics and Monitoring**: Track performance across the edge network

## Getting Started

### Installation

```bash
go get github.com/davidjeba/goscript/pkg/goscale
```

### Basic Usage

#### Creating an API

```go
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/davidjeba/goscript/pkg/goscale/api"
)

func main() {
	// Create a new GoScaleAPI instance
	apiConfig := &api.Config{
		DBConnectionString: "localhost:5432",
		EdgeEnabled:        true,
		EdgeNodes:          []string{"edge-1", "edge-2"},
		CompressionLevel:   5,
		BatchSize:          100,
		Timeout:            time.Second * 30,
		MaxConcurrent:      1000,
		EnableTimeSeries:   true,
		EnableRelationships: true,
		EnableNoCode:       true,
	}
	
	goscaleAPI := api.NewGoScaleAPI(apiConfig)
	
	// Create a schema
	schema := api.NewSchema()
	
	// Add a type
	userType := schema.AddType("User", "A user in the system")
	userType.AddField("id", "ID", "The user's ID")
	userType.AddField("name", "String", "The user's name")
	userType.AddField("email", "String", "The user's email")
	
	// Add a query
	getUserQuery := schema.AddQuery("getUser", "User", "Get a user by ID")
	getUserQuery.AddArg("id", "ID", nil, "The user's ID")
	getUserQuery.SetResolver(func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
		// Query logic here
		return map[string]interface{}{
			"id":    params["id"],
			"name":  "John Doe",
			"email": "john@example.com",
		}, nil
	})
	
	// Apply the schema
	err := goscaleAPI.ApplySchema(schema)
	if err != nil {
		log.Fatalf("Error applying schema: %v", err)
	}
	
	// Start the server
	http.Handle("/api", goscaleAPI)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

#### Using the Database

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/davidjeba/goscript/pkg/goscale/db"
)

func main() {
	// Create a new GoScaleDB instance
	dbConfig := &db.Config{
		ConnectionString:   "localhost:5432",
		MaxConnections:     100,
		QueryTimeout:       time.Second * 30,
		EnableTimeSeries:   true,
		EnableRelationships: true,
		EnableNoCode:       true,
	}
	
	database := db.NewGoScaleDB(dbConfig)
	
	// Connect to the database
	err := database.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer database.Close()
	
	// Create a schema
	schema, err := database.CreateSchema("app")
	if err != nil {
		log.Fatalf("Error creating schema: %v", err)
	}
	
	// Create a table
	columns := map[string]*db.Column{
		"id": {
			Name:     "id",
			Type:     "SERIAL",
			Nullable: false,
		},
		"name": {
			Name:     "name",
			Type:     "VARCHAR(255)",
			Nullable: false,
		},
		"email": {
			Name:     "email",
			Type:     "VARCHAR(255)",
			Nullable: false,
		},
		"created_at": {
			Name:     "created_at",
			Type:     "TIMESTAMP",
			Nullable: false,
			Default:  "NOW()",
		},
	}
	
	table, err := database.CreateTable("app", "users", columns, "id")
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	
	// Create an index
	_, err = database.CreateIndex("app", "users", "users_email_idx", []string{"email"}, true)
	if err != nil {
		log.Fatalf("Error creating index: %v", err)
	}
	
	// Insert data
	ctx := context.Background()
	data := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
	}
	
	id, err := database.Insert(ctx, "app", "users", data)
	if err != nil {
		log.Fatalf("Error inserting data: %v", err)
	}
	
	fmt.Printf("Inserted user with ID: %d\n", id)
	
	// Query data
	rows, err := database.Query(ctx, "SELECT * FROM app.users WHERE id = $1", id)
	if err != nil {
		log.Fatalf("Error querying data: %v", err)
	}
	
	fmt.Printf("Query result: %v\n", rows)
}
```

#### Setting Up Edge Computing

```go
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/davidjeba/goscript/pkg/goscale/api"
	"github.com/davidjeba/goscript/pkg/goscale/db"
	"github.com/davidjeba/goscript/pkg/goscale/edge"
)

func main() {
	// Create a new GoScaleAPI instance
	goscaleAPI := api.NewGoScaleAPI(nil)
	
	// Create an edge network
	edgeNetwork := edge.NewEdgeNetwork(goscaleAPI)
	
	// Create edge nodes
	edgeConfig1 := &edge.Config{
		ID:               "edge-1",
		Region:           "us-east",
		Capacity:         1000,
		CacheEnabled:     true,
		CacheTTL:         time.Minute * 5,
		DBConfig:         db.DefaultConfig(),
		SyncInterval:     time.Minute * 15,
		MaxConcurrent:    100,
		CompressionLevel: 5,
	}
	
	edgeNode1 := edge.NewEdgeNode(edgeConfig1, goscaleAPI)
	edgeNetwork.AddNode(edgeNode1)
	
	edgeConfig2 := &edge.Config{
		ID:               "edge-2",
		Region:           "us-west",
		Capacity:         1000,
		CacheEnabled:     true,
		CacheTTL:         time.Minute * 5,
		DBConfig:         db.DefaultConfig(),
		SyncInterval:     time.Minute * 15,
		MaxConcurrent:    100,
		CompressionLevel: 5,
	}
	
	edgeNode2 := edge.NewEdgeNode(edgeConfig2, goscaleAPI)
	edgeNetwork.AddNode(edgeNode2)
	
	// Register handlers for edge nodes
	for path, resolver := range goscaleAPI.GetResolvers() {
		edgeNode1.RegisterHandler(path, resolver)
		edgeNode2.RegisterHandler(path, resolver)
	}
	
	// Create HTTP handlers
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		goscaleAPI.ServeHTTP(w, r)
	})
	
	http.HandleFunc("/edge", func(w http.ResponseWriter, r *http.Request) {
		edgeNetwork.ServeHTTP(w, r)
	})
	
	// Start the server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Features

### API Features

- **Schema-based API**: Define your API using a schema with types, queries, mutations, and subscriptions
- **Resolver Functions**: Implement custom logic for each API operation
- **Middleware System**: Process requests through a pipeline of middleware functions
- **Real-time Subscriptions**: Subscribe to real-time data updates
- **Edge Computing**: Process API requests at the network edge for low latency
- **Metrics and Monitoring**: Track API performance and usage

### Database Features

- **Schema Management**: Create and manage database schemas and tables
- **Query Caching**: Automatically cache query results for improved performance
- **Time Series Data**: Store and query time series data with retention policies
- **NoCode Database**: Store schema-less data with validation
- **Relationship Management**: Define and query relationships between entities
- **Sharding and Replication**: Scale horizontally with data distribution
- **Metrics and Monitoring**: Track database performance and usage

### Edge Computing Features

- **Distributed Processing**: Process API requests at the network edge
- **Caching**: Cache data close to users for improved performance
- **Load Balancing**: Distribute requests across edge nodes
- **Health Monitoring**: Automatically check the health of edge nodes
- **Synchronization**: Keep edge nodes in sync with the central system
- **Metrics and Monitoring**: Track edge network performance and usage

## Performance

GoScale is designed for high performance:

- **Low Latency**: Process requests in microseconds
- **High Throughput**: Handle thousands of requests per second
- **Efficient Caching**: Reduce database load with intelligent caching
- **Edge Computing**: Process requests close to users for minimal latency
- **Optimized Data Storage**: Store and retrieve data efficiently
- **Horizontal Scaling**: Scale out with additional nodes

## Comparison with Other Technologies

### GoScale API vs GraphQL

- **Performance**: GoScale is optimized for performance, while GraphQL can be slower due to its flexibility
- **Flexibility**: Both offer schema-based APIs with queries, mutations, and subscriptions
- **Edge Computing**: GoScale includes built-in edge computing support
- **Language**: GoScale is written in Go, while GraphQL is language-agnostic
- **Ecosystem**: GraphQL has a larger ecosystem, but GoScale integrates better with Go applications

### GoScale API vs gRPC

- **Performance**: Both offer high performance, but GoScale includes additional optimizations
- **Flexibility**: GoScale offers GraphQL-like flexibility, while gRPC is more rigid
- **Edge Computing**: GoScale includes built-in edge computing support
- **Protocol**: GoScale uses HTTP/JSON, while gRPC uses HTTP/2 and Protocol Buffers
- **Code Generation**: gRPC requires code generation, while GoScale uses runtime reflection

### GoScale DB vs PostgreSQL

- **Performance**: GoScale includes additional optimizations and caching
- **Features**: GoScale includes time series and NoCode features
- **Flexibility**: GoScale offers more flexibility with its NoCode features
- **Ecosystem**: PostgreSQL has a larger ecosystem
- **Maturity**: PostgreSQL is more mature and battle-tested

### GoScale DB vs TimescaleDB

- **Performance**: Both offer high performance for time series data
- **Features**: GoScale includes NoCode and relationship features
- **Flexibility**: GoScale offers more flexibility with its NoCode features
- **Ecosystem**: TimescaleDB has a larger ecosystem for time series data
- **Maturity**: TimescaleDB is more mature for time series data

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License