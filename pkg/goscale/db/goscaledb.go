package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

// GoScaleDB is a high-performance database that combines features of
// PostgreSQL, TimescaleDB, and NoCode databases with a focus on
// simplicity, robustness, and scalability.
type GoScaleDB struct {
	conn            *sql.DB
	queryCache      map[string]*CachedQuery
	cacheMutex      sync.RWMutex
	schemas         map[string]*Schema
	schemaMutex     sync.RWMutex
	timeSeries      *TimeSeriesManager
	relationships   *RelationshipManager
	noCode          *NoCodeManager
	metrics         *Metrics
	config          *Config
	shards          []*Shard
	shardingEnabled bool
	replicaNodes    []string
	replicationMode string
	migrationLock   sync.Mutex
}

// Config contains configuration options for GoScaleDB
type Config struct {
	ConnectionString   string
	MaxConnections     int
	QueryTimeout       time.Duration
	EnableTimeSeries   bool
	EnableRelationships bool
	EnableNoCode       bool
	ShardingEnabled    bool
	ShardCount         int
	ReplicationMode    string
	ReplicaNodes       []string
	AutoMigrate        bool
	CacheSize          int
	CacheTTL           time.Duration
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		ConnectionString:   "localhost:5432",
		MaxConnections:     100,
		QueryTimeout:       time.Second * 30,
		EnableTimeSeries:   true,
		EnableRelationships: true,
		EnableNoCode:       true,
		ShardingEnabled:    false,
		ShardCount:         1,
		ReplicationMode:    "async",
		ReplicaNodes:       []string{},
		AutoMigrate:        true,
		CacheSize:          1000,
		CacheTTL:           time.Minute * 5,
	}
}

// CachedQuery represents a cached query result
type CachedQuery struct {
	Query      string
	Result     interface{}
	Expiration time.Time
}

// Schema represents a database schema
type Schema struct {
	Name    string
	Tables  map[string]*Table
	Version int
}

// Table represents a database table
type Table struct {
	Name       string
	Columns    map[string]*Column
	Indexes    map[string]*Index
	PrimaryKey string
	TimeColumn string
}

// Column represents a table column
type Column struct {
	Name     string
	Type     string
	Nullable bool
	Default  interface{}
}

// Index represents a table index
type Index struct {
	Name    string
	Columns []string
	Unique  bool
}

// Metrics tracks database performance metrics
type Metrics struct {
	QueryCount      int64
	AvgQueryTime    float64
	CacheHitRate    float64
	WriteCount      int64
	AvgWriteTime    float64
	ErrorRate       float64
	mutex           sync.RWMutex
}

// TimeSeriesManager manages time-series data
type TimeSeriesManager struct {
	db             *GoScaleDB
	enabledTables  map[string]bool
	retentionPolicies map[string]time.Duration
	aggregationFuncs map[string]string
	mutex          sync.RWMutex
}

// RelationshipManager manages relationships between entities
type RelationshipManager struct {
	db             *GoScaleDB
	relationships  map[string]map[string]*Relationship
	mutex          sync.RWMutex
}

// Relationship represents a relationship between two entities
type Relationship struct {
	Name        string
	SourceTable string
	TargetTable string
	Type        string // OneToOne, OneToMany, ManyToMany
	JoinTable   string
	SourceKey   string
	TargetKey   string
}

// NoCodeManager manages NoCode database features
type NoCodeManager struct {
	db          *GoScaleDB
	schemas     map[string]*NoCodeSchema
	mutex       sync.RWMutex
}

// NoCodeSchema represents a NoCode schema
type NoCodeSchema struct {
	Name    string
	Fields  map[string]*NoCodeField
	Version int
}

// NoCodeField represents a field in a NoCode schema
type NoCodeField struct {
	Name     string
	Type     string
	Required bool
	Default  interface{}
	Validators []string
}

// Shard represents a database shard
type Shard struct {
	ID       int
	Conn     *sql.DB
	KeyRange [2]int64
	Tables   []string
}

// NewGoScaleDB creates a new instance of GoScaleDB
func NewGoScaleDB(config *Config) *GoScaleDB {
	if config == nil {
		config = DefaultConfig()
	}
	
	db := &GoScaleDB{
		queryCache:      make(map[string]*CachedQuery),
		schemas:         make(map[string]*Schema),
		config:          config,
		shardingEnabled: config.ShardingEnabled,
		replicaNodes:    config.ReplicaNodes,
		replicationMode: config.ReplicationMode,
		metrics:         &Metrics{},
	}
	
	// Initialize time series manager if enabled
	if config.EnableTimeSeries {
		db.timeSeries = &TimeSeriesManager{
			db:             db,
			enabledTables:  make(map[string]bool),
			retentionPolicies: make(map[string]time.Duration),
			aggregationFuncs: make(map[string]string),
		}
	}
	
	// Initialize relationship manager if enabled
	if config.EnableRelationships {
		db.relationships = &RelationshipManager{
			db:            db,
			relationships: make(map[string]map[string]*Relationship),
		}
	}
	
	// Initialize NoCode manager if enabled
	if config.EnableNoCode {
		db.noCode = &NoCodeManager{
			db:      db,
			schemas: make(map[string]*NoCodeSchema),
		}
	}
	
	// Initialize shards if sharding is enabled
	if config.ShardingEnabled && config.ShardCount > 0 {
		db.shards = make([]*Shard, config.ShardCount)
		for i := 0; i < config.ShardCount; i++ {
			// In a real implementation, we would connect to different shards
			db.shards[i] = &Shard{
				ID:       i,
				KeyRange: [2]int64{int64(i) * (1<<63) / int64(config.ShardCount), (int64(i)+1) * (1<<63) / int64(config.ShardCount)},
				Tables:   []string{},
			}
		}
	}
	
	return db
}

// Connect establishes a connection to the database
func (db *GoScaleDB) Connect() error {
	var err error
	db.conn, err = sql.Open("postgres", db.config.ConnectionString)
	if err != nil {
		return err
	}
	
	db.conn.SetMaxOpenConns(db.config.MaxConnections)
	db.conn.SetMaxIdleConns(db.config.MaxConnections / 2)
	
	// Connect to shards if sharding is enabled
	if db.shardingEnabled {
		for i, shard := range db.shards {
			// In a real implementation, we would connect to different shard databases
			shard.Conn, err = sql.Open("postgres", fmt.Sprintf("%s_shard_%d", db.config.ConnectionString, i))
			if err != nil {
				return err
			}
		}
	}
	
	// Initialize time series features if enabled
	if db.config.EnableTimeSeries {
		err = db.initializeTimeSeries()
		if err != nil {
			return err
		}
	}
	
	return nil
}

// Close closes the database connection
func (db *GoScaleDB) Close() error {
	// Close shard connections
	if db.shardingEnabled {
		for _, shard := range db.shards {
			if shard.Conn != nil {
				shard.Conn.Close()
			}
		}
	}
	
	// Close main connection
	if db.conn != nil {
		return db.conn.Close()
	}
	
	return nil
}

// Query executes a query and returns the results
func (db *GoScaleDB) Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	startTime := time.Now()
	
	// Check cache first
	cacheKey := fmt.Sprintf("%s:%v", query, args)
	db.cacheMutex.RLock()
	if cached, ok := db.queryCache[cacheKey]; ok && time.Now().Before(cached.Expiration) {
		db.cacheMutex.RUnlock()
		db.updateMetrics(startTime, true, true, false)
		return cached.Result.([]map[string]interface{}), nil
	}
	db.cacheMutex.RUnlock()
	
	// Execute the query
	rows, err := db.conn.QueryContext(ctx, query, args...)
	if err != nil {
		db.updateMetrics(startTime, false, false, false)
		return nil, err
	}
	defer rows.Close()
	
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		db.updateMetrics(startTime, false, false, false)
		return nil, err
	}
	
	// Prepare result
	var result []map[string]interface{}
	
	// Scan rows
	for rows.Next() {
		// Create a slice of interface{} to hold the values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		
		// Initialize the pointers
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		
		// Scan the row into the slice
		if err := rows.Scan(valuePtrs...); err != nil {
			db.updateMetrics(startTime, false, false, false)
			return nil, err
		}
		
		// Create a map for this row
		row := make(map[string]interface{})
		
		// Convert the values to their appropriate types
		for i, col := range columns {
			val := values[i]
			
			// Handle null values
			if val == nil {
				row[col] = nil
				continue
			}
			
			// Handle different types
			switch v := val.(type) {
			case []byte:
				// Try to unmarshal as JSON first
				var jsonVal interface{}
				if err := json.Unmarshal(v, &jsonVal); err == nil {
					row[col] = jsonVal
				} else {
					// If not JSON, use as string
					row[col] = string(v)
				}
			default:
				row[col] = v
			}
		}
		
		result = append(result, row)
	}
	
	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		db.updateMetrics(startTime, false, false, false)
		return nil, err
	}
	
	// Cache the result
	db.cacheMutex.Lock()
	db.queryCache[cacheKey] = &CachedQuery{
		Query:      query,
		Result:     result,
		Expiration: time.Now().Add(db.config.CacheTTL),
	}
	db.cacheMutex.Unlock()
	
	db.updateMetrics(startTime, true, false, false)
	return result, nil
}

// Execute executes a non-query SQL statement
func (db *GoScaleDB) Execute(ctx context.Context, query string, args ...interface{}) (int64, error) {
	startTime := time.Now()
	
	// Execute the statement
	result, err := db.conn.ExecContext(ctx, query, args...)
	if err != nil {
		db.updateMetrics(startTime, false, false, true)
		return 0, err
	}
	
	// Get the number of affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.updateMetrics(startTime, false, false, true)
		return 0, err
	}
	
	// Invalidate cache for write operations
	db.cacheMutex.Lock()
	db.queryCache = make(map[string]*CachedQuery)
	db.cacheMutex.Unlock()
	
	// Replicate to replica nodes if replication is enabled
	if len(db.replicaNodes) > 0 {
		go db.replicateQuery(query, args...)
	}
	
	db.updateMetrics(startTime, true, false, true)
	return rowsAffected, nil
}

// Transaction executes a function within a transaction
func (db *GoScaleDB) Transaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	
	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	
	return tx.Commit()
}

// CreateSchema creates a new database schema
func (db *GoScaleDB) CreateSchema(name string) (*Schema, error) {
	db.schemaMutex.Lock()
	defer db.schemaMutex.Unlock()
	
	if _, ok := db.schemas[name]; ok {
		return nil, fmt.Errorf("schema %s already exists", name)
	}
	
	schema := &Schema{
		Name:    name,
		Tables:  make(map[string]*Table),
		Version: 1,
	}
	
	db.schemas[name] = schema
	
	// Create the schema in the database
	_, err := db.Execute(context.Background(), fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", name))
	if err != nil {
		delete(db.schemas, name)
		return nil, err
	}
	
	return schema, nil
}

// GetSchema returns a schema by name
func (db *GoScaleDB) GetSchema(name string) (*Schema, error) {
	db.schemaMutex.RLock()
	defer db.schemaMutex.RUnlock()
	
	schema, ok := db.schemas[name]
	if !ok {
		return nil, fmt.Errorf("schema %s not found", name)
	}
	
	return schema, nil
}

// CreateTable creates a new table in a schema
func (db *GoScaleDB) CreateTable(schemaName, tableName string, columns map[string]*Column, primaryKey string) (*Table, error) {
	db.schemaMutex.Lock()
	defer db.schemaMutex.Unlock()
	
	schema, ok := db.schemas[schemaName]
	if !ok {
		return nil, fmt.Errorf("schema %s not found", schemaName)
	}
	
	if _, ok := schema.Tables[tableName]; ok {
		return nil, fmt.Errorf("table %s already exists in schema %s", tableName, schemaName)
	}
	
	table := &Table{
		Name:       tableName,
		Columns:    columns,
		Indexes:    make(map[string]*Index),
		PrimaryKey: primaryKey,
	}
	
	schema.Tables[tableName] = table
	
	// Build the CREATE TABLE statement
	query := fmt.Sprintf("CREATE TABLE %s.%s (", schemaName, tableName)
	
	i := 0
	for name, column := range columns {
		if i > 0 {
			query += ", "
		}
		
		query += fmt.Sprintf("%s %s", name, column.Type)
		
		if !column.Nullable {
			query += " NOT NULL"
		}
		
		if column.Default != nil {
			query += fmt.Sprintf(" DEFAULT %v", column.Default)
		}
		
		if name == primaryKey {
			query += " PRIMARY KEY"
		}
		
		i++
	}
	
	query += ")"
	
	// Create the table in the database
	_, err := db.Execute(context.Background(), query)
	if err != nil {
		delete(schema.Tables, tableName)
		return nil, err
	}
	
	// Enable time series for this table if it has a time column and time series is enabled
	if db.config.EnableTimeSeries && table.TimeColumn != "" {
		err = db.timeSeries.EnableTimeSeriesForTable(schemaName, tableName, table.TimeColumn)
		if err != nil {
			return nil, err
		}
	}
	
	return table, nil
}

// GetTable returns a table by name
func (db *GoScaleDB) GetTable(schemaName, tableName string) (*Table, error) {
	db.schemaMutex.RLock()
	defer db.schemaMutex.RUnlock()
	
	schema, ok := db.schemas[schemaName]
	if !ok {
		return nil, fmt.Errorf("schema %s not found", schemaName)
	}
	
	table, ok := schema.Tables[tableName]
	if !ok {
		return nil, fmt.Errorf("table %s not found in schema %s", tableName, schemaName)
	}
	
	return table, nil
}

// CreateIndex creates a new index on a table
func (db *GoScaleDB) CreateIndex(schemaName, tableName, indexName string, columns []string, unique bool) (*Index, error) {
	db.schemaMutex.Lock()
	defer db.schemaMutex.Unlock()
	
	schema, ok := db.schemas[schemaName]
	if !ok {
		return nil, fmt.Errorf("schema %s not found", schemaName)
	}
	
	table, ok := schema.Tables[tableName]
	if !ok {
		return nil, fmt.Errorf("table %s not found in schema %s", tableName, schemaName)
	}
	
	if _, ok := table.Indexes[indexName]; ok {
		return nil, fmt.Errorf("index %s already exists on table %s.%s", indexName, schemaName, tableName)
	}
	
	index := &Index{
		Name:    indexName,
		Columns: columns,
		Unique:  unique,
	}
	
	table.Indexes[indexName] = index
	
	// Build the CREATE INDEX statement
	uniqueStr := ""
	if unique {
		uniqueStr = "UNIQUE "
	}
	
	columnStr := ""
	for i, col := range columns {
		if i > 0 {
			columnStr += ", "
		}
		columnStr += col
	}
	
	query := fmt.Sprintf("CREATE %sINDEX %s ON %s.%s (%s)", uniqueStr, indexName, schemaName, tableName, columnStr)
	
	// Create the index in the database
	_, err := db.Execute(context.Background(), query)
	if err != nil {
		delete(table.Indexes, indexName)
		return nil, err
	}
	
	return index, nil
}

// Insert inserts a new row into a table
func (db *GoScaleDB) Insert(ctx context.Context, schemaName, tableName string, data map[string]interface{}) (int64, error) {
	// Get the table
	table, err := db.GetTable(schemaName, tableName)
	if err != nil {
		return 0, err
	}
	
	// Build the INSERT statement
	columns := ""
	values := ""
	args := []interface{}{}
	
	i := 0
	for col, val := range data {
		if _, ok := table.Columns[col]; !ok {
			return 0, fmt.Errorf("column %s not found in table %s.%s", col, schemaName, tableName)
		}
		
		if i > 0 {
			columns += ", "
			values += ", "
		}
		
		columns += col
		values += fmt.Sprintf("$%d", i+1)
		args = append(args, val)
		
		i++
	}
	
	query := fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES (%s) RETURNING %s", schemaName, tableName, columns, values, table.PrimaryKey)
	
	// Execute the query
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	
	if len(rows) == 0 {
		return 0, errors.New("no rows returned after insert")
	}
	
	// Get the primary key value
	pkVal, ok := rows[0][table.PrimaryKey]
	if !ok {
		return 0, fmt.Errorf("primary key %s not found in returned row", table.PrimaryKey)
	}
	
	// Convert to int64
	var id int64
	switch v := pkVal.(type) {
	case int64:
		id = v
	case int:
		id = int64(v)
	case float64:
		id = int64(v)
	default:
		return 0, fmt.Errorf("primary key %s is not a number", table.PrimaryKey)
	}
	
	return id, nil
}

// Update updates rows in a table
func (db *GoScaleDB) Update(ctx context.Context, schemaName, tableName string, data map[string]interface{}, where string, args ...interface{}) (int64, error) {
	// Get the table
	table, err := db.GetTable(schemaName, tableName)
	if err != nil {
		return 0, err
	}
	
	// Build the UPDATE statement
	set := ""
	setArgs := []interface{}{}
	
	i := 0
	for col, val := range data {
		if _, ok := table.Columns[col]; !ok {
			return 0, fmt.Errorf("column %s not found in table %s.%s", col, schemaName, tableName)
		}
		
		if i > 0 {
			set += ", "
		}
		
		set += fmt.Sprintf("%s = $%d", col, i+1)
		setArgs = append(setArgs, val)
		
		i++
	}
	
	// Adjust the placeholder indices in the WHERE clause
	for j := range args {
		where = fmt.Sprintf(where, fmt.Sprintf("$%d", i+j+1))
	}
	
	query := fmt.Sprintf("UPDATE %s.%s SET %s WHERE %s", schemaName, tableName, set, where)
	
	// Combine the arguments
	allArgs := append(setArgs, args...)
	
	// Execute the query
	return db.Execute(ctx, query, allArgs...)
}

// Delete deletes rows from a table
func (db *GoScaleDB) Delete(ctx context.Context, schemaName, tableName string, where string, args ...interface{}) (int64, error) {
	// Get the table
	_, err := db.GetTable(schemaName, tableName)
	if err != nil {
		return 0, err
	}
	
	// Build the DELETE statement
	query := fmt.Sprintf("DELETE FROM %s.%s WHERE %s", schemaName, tableName, where)
	
	// Execute the query
	return db.Execute(ctx, query, args...)
}

// initializeTimeSeries initializes time series features
func (db *GoScaleDB) initializeTimeSeries() error {
	// In a real implementation, this would initialize TimescaleDB extensions
	return nil
}

// EnableTimeSeriesForTable enables time series features for a table
func (ts *TimeSeriesManager) EnableTimeSeriesForTable(schemaName, tableName, timeColumn string) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	
	key := fmt.Sprintf("%s.%s", schemaName, tableName)
	ts.enabledTables[key] = true
	
	// In a real implementation, this would create a hypertable
	query := fmt.Sprintf("SELECT create_hypertable('%s.%s', '%s')", schemaName, tableName, timeColumn)
	_, err := ts.db.Execute(context.Background(), query)
	
	return err
}

// SetRetentionPolicy sets a retention policy for a time series table
func (ts *TimeSeriesManager) SetRetentionPolicy(schemaName, tableName string, retention time.Duration) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	
	key := fmt.Sprintf("%s.%s", schemaName, tableName)
	if !ts.enabledTables[key] {
		return fmt.Errorf("table %s.%s is not a time series table", schemaName, tableName)
	}
	
	ts.retentionPolicies[key] = retention
	
	// In a real implementation, this would set a retention policy
	query := fmt.Sprintf("SELECT add_retention_policy('%s.%s', INTERVAL '%d seconds')", schemaName, tableName, int(retention.Seconds()))
	_, err := ts.db.Execute(context.Background(), query)
	
	return err
}

// CreateRelationship creates a relationship between two tables
func (rm *RelationshipManager) CreateRelationship(name, sourceTable, targetTable, relType, sourceKey, targetKey string) (*Relationship, error) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	
	if _, ok := rm.relationships[sourceTable]; !ok {
		rm.relationships[sourceTable] = make(map[string]*Relationship)
	}
	
	if _, ok := rm.relationships[sourceTable][name]; ok {
		return nil, fmt.Errorf("relationship %s already exists for table %s", name, sourceTable)
	}
	
	rel := &Relationship{
		Name:        name,
		SourceTable: sourceTable,
		TargetTable: targetTable,
		Type:        relType,
		SourceKey:   sourceKey,
		TargetKey:   targetKey,
	}
	
	// For many-to-many relationships, create a join table
	if relType == "ManyToMany" {
		joinTable := fmt.Sprintf("%s_%s", sourceTable, targetTable)
		rel.JoinTable = joinTable
		
		// Create the join table
		// This is simplified; in a real implementation, we would create the actual table
	}
	
	rm.relationships[sourceTable][name] = rel
	
	return rel, nil
}

// GetRelatedEntities gets entities related to a source entity
func (rm *RelationshipManager) GetRelatedEntities(ctx context.Context, sourceTable, relationshipName string, sourceID interface{}) ([]map[string]interface{}, error) {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()
	
	if _, ok := rm.relationships[sourceTable]; !ok {
		return nil, fmt.Errorf("no relationships found for table %s", sourceTable)
	}
	
	rel, ok := rm.relationships[sourceTable][relationshipName]
	if !ok {
		return nil, fmt.Errorf("relationship %s not found for table %s", relationshipName, sourceTable)
	}
	
	var query string
	var args []interface{}
	
	switch rel.Type {
	case "OneToOne", "OneToMany":
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", rel.TargetTable, rel.TargetKey)
		args = []interface{}{sourceID}
	case "ManyToMany":
		query = fmt.Sprintf("SELECT t.* FROM %s t JOIN %s j ON t.%s = j.%s WHERE j.%s = $1",
			rel.TargetTable, rel.JoinTable, rel.TargetKey, rel.TargetKey, rel.SourceKey)
		args = []interface{}{sourceID}
	default:
		return nil, fmt.Errorf("unknown relationship type: %s", rel.Type)
	}
	
	return rm.db.Query(ctx, query, args...)
}

// CreateNoCodeSchema creates a new NoCode schema
func (nc *NoCodeManager) CreateNoCodeSchema(name string) (*NoCodeSchema, error) {
	nc.mutex.Lock()
	defer nc.mutex.Unlock()
	
	if _, ok := nc.schemas[name]; ok {
		return nil, fmt.Errorf("NoCode schema %s already exists", name)
	}
	
	schema := &NoCodeSchema{
		Name:    name,
		Fields:  make(map[string]*NoCodeField),
		Version: 1,
	}
	
	nc.schemas[name] = schema
	
	return schema, nil
}

// AddNoCodeField adds a field to a NoCode schema
func (nc *NoCodeManager) AddNoCodeField(schemaName, fieldName, fieldType string, required bool, defaultValue interface{}, validators []string) (*NoCodeField, error) {
	nc.mutex.Lock()
	defer nc.mutex.Unlock()
	
	schema, ok := nc.schemas[schemaName]
	if !ok {
		return nil, fmt.Errorf("NoCode schema %s not found", schemaName)
	}
	
	if _, ok := schema.Fields[fieldName]; ok {
		return nil, fmt.Errorf("field %s already exists in NoCode schema %s", fieldName, schemaName)
	}
	
	field := &NoCodeField{
		Name:       fieldName,
		Type:       fieldType,
		Required:   required,
		Default:    defaultValue,
		Validators: validators,
	}
	
	schema.Fields[fieldName] = field
	
	return field, nil
}

// CreateNoCodeEntity creates a new entity in a NoCode schema
func (nc *NoCodeManager) CreateNoCodeEntity(ctx context.Context, schemaName string, data map[string]interface{}) (int64, error) {
	nc.mutex.RLock()
	schema, ok := nc.schemas[schemaName]
	if !ok {
		nc.mutex.RUnlock()
		return 0, fmt.Errorf("NoCode schema %s not found", schemaName)
	}
	nc.mutex.RUnlock()
	
	// Validate the data against the schema
	for name, field := range schema.Fields {
		if field.Required {
			if _, ok := data[name]; !ok {
				return 0, fmt.Errorf("required field %s is missing", name)
			}
		}
		
		if val, ok := data[name]; ok {
			// Validate the type
			switch field.Type {
			case "string":
				if _, ok := val.(string); !ok {
					return 0, fmt.Errorf("field %s must be a string", name)
				}
			case "number":
				switch val.(type) {
				case int, int64, float64:
					// Valid number types
				default:
					return 0, fmt.Errorf("field %s must be a number", name)
				}
			case "boolean":
				if _, ok := val.(bool); !ok {
					return 0, fmt.Errorf("field %s must be a boolean", name)
				}
			case "object":
				if _, ok := val.(map[string]interface{}); !ok {
					return 0, fmt.Errorf("field %s must be an object", name)
				}
			case "array":
				if _, ok := val.([]interface{}); !ok {
					return 0, fmt.Errorf("field %s must be an array", name)
				}
			}
			
			// Apply validators
			for _, validator := range field.Validators {
				// In a real implementation, we would apply the validators
			}
		} else if field.Default != nil {
			// Use default value
			data[name] = field.Default
		}
	}
	
	// Create the entity
	return nc.db.Insert(ctx, "nocode", schemaName, data)
}

// GetNoCodeEntity gets an entity from a NoCode schema
func (nc *NoCodeManager) GetNoCodeEntity(ctx context.Context, schemaName string, id int64) (map[string]interface{}, error) {
	nc.mutex.RLock()
	_, ok := nc.schemas[schemaName]
	if !ok {
		nc.mutex.RUnlock()
		return nil, fmt.Errorf("NoCode schema %s not found", schemaName)
	}
	nc.mutex.RUnlock()
	
	// Get the entity
	rows, err := nc.db.Query(ctx, fmt.Sprintf("SELECT * FROM nocode.%s WHERE id = $1", schemaName), id)
	if err != nil {
		return nil, err
	}
	
	if len(rows) == 0 {
		return nil, fmt.Errorf("entity with ID %d not found in schema %s", id, schemaName)
	}
	
	return rows[0], nil
}

// updateMetrics updates the database metrics
func (db *GoScaleDB) updateMetrics(startTime time.Time, success, cacheHit, isWrite bool) {
	duration := time.Since(startTime).Seconds()
	
	db.metrics.mutex.Lock()
	defer db.metrics.mutex.Unlock()
	
	if isWrite {
		db.metrics.WriteCount++
		db.metrics.AvgWriteTime = (db.metrics.AvgWriteTime*float64(db.metrics.WriteCount-1) + duration) / float64(db.metrics.WriteCount)
	} else {
		db.metrics.QueryCount++
		db.metrics.AvgQueryTime = (db.metrics.AvgQueryTime*float64(db.metrics.QueryCount-1) + duration) / float64(db.metrics.QueryCount)
	}
	
	if cacheHit {
		db.metrics.CacheHitRate = (db.metrics.CacheHitRate*float64(db.metrics.QueryCount-1) + 1) / float64(db.metrics.QueryCount)
	} else {
		db.metrics.CacheHitRate = (db.metrics.CacheHitRate * float64(db.metrics.QueryCount-1)) / float64(db.metrics.QueryCount)
	}
	
	if !success {
		db.metrics.ErrorRate = (db.metrics.ErrorRate*float64(db.metrics.QueryCount+db.metrics.WriteCount-1) + 1) / float64(db.metrics.QueryCount+db.metrics.WriteCount)
	} else {
		db.metrics.ErrorRate = (db.metrics.ErrorRate * float64(db.metrics.QueryCount+db.metrics.WriteCount-1)) / float64(db.metrics.QueryCount+db.metrics.WriteCount)
	}
}

// GetMetrics returns the current database metrics
func (db *GoScaleDB) GetMetrics() *Metrics {
	db.metrics.mutex.RLock()
	defer db.metrics.mutex.RUnlock()
	
	return &Metrics{
		QueryCount:   db.metrics.QueryCount,
		AvgQueryTime: db.metrics.AvgQueryTime,
		CacheHitRate: db.metrics.CacheHitRate,
		WriteCount:   db.metrics.WriteCount,
		AvgWriteTime: db.metrics.AvgWriteTime,
		ErrorRate:    db.metrics.ErrorRate,
	}
}

// replicateQuery replicates a query to replica nodes
func (db *GoScaleDB) replicateQuery(query string, args ...interface{}) {
	// In a real implementation, this would send the query to replica nodes
}

// GetShardForKey returns the shard for a given key
func (db *GoScaleDB) GetShardForKey(key int64) (*Shard, error) {
	if !db.shardingEnabled {
		return nil, errors.New("sharding is not enabled")
	}
	
	for _, shard := range db.shards {
		if key >= shard.KeyRange[0] && key < shard.KeyRange[1] {
			return shard, nil
		}
	}
	
	return nil, fmt.Errorf("no shard found for key %d", key)
}