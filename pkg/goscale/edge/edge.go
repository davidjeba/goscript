package edge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/davidjeba/goscript/pkg/goscale/api"
	"github.com/davidjeba/goscript/pkg/goscale/db"
)

// EdgeNode represents an edge computing node that can process API requests
// with low latency at the network edge
type EdgeNode struct {
	ID              string
	Region          string
	Capacity        int
	Load            int
	APIHandlers     map[string]api.Resolver
	CacheEnabled    bool
	Cache           map[string]*CacheEntry
	CacheTTL        time.Duration
	CacheMutex      sync.RWMutex
	LocalDB         *db.GoScaleDB
	SyncInterval    time.Duration
	LastSyncTime    time.Time
	SyncMutex       sync.Mutex
	HealthStatus    string
	Metrics         *EdgeMetrics
	ParentAPI       *api.GoScaleAPI
	MaxConcurrent   int
	RequestQueue    chan *EdgeRequest
	WorkerPool      []*EdgeWorker
	CompressionLevel int
}

// EdgeRequest represents a request to be processed by the edge node
type EdgeRequest struct {
	Path       string
	Params     map[string]interface{}
	Context    context.Context
	ResultChan chan *EdgeResponse
}

// EdgeResponse represents a response from the edge node
type EdgeResponse struct {
	Result interface{}
	Error  error
}

// EdgeWorker represents a worker that processes edge requests
type EdgeWorker struct {
	ID         int
	RequestChan chan *EdgeRequest
	Node       *EdgeNode
	Active     bool
}

// CacheEntry represents a cached API response
type CacheEntry struct {
	Path       string
	Params     map[string]interface{}
	Result     interface{}
	Expiration time.Time
}

// EdgeMetrics tracks edge node performance metrics
type EdgeMetrics struct {
	RequestCount    int64
	AvgResponseTime float64
	ErrorRate       float64
	CacheHitRate    float64
	CPUUsage        float64
	MemoryUsage     float64
	NetworkIn       int64
	NetworkOut      int64
	mutex           sync.RWMutex
}

// Config contains configuration options for EdgeNode
type Config struct {
	ID               string
	Region           string
	Capacity         int
	CacheEnabled     bool
	CacheTTL         time.Duration
	DBConfig         *db.Config
	SyncInterval     time.Duration
	MaxConcurrent    int
	CompressionLevel int
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
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
}

// NewEdgeNode creates a new edge node
func NewEdgeNode(config *Config, parentAPI *api.GoScaleAPI) *EdgeNode {
	if config == nil {
		config = DefaultConfig()
	}
	
	node := &EdgeNode{
		ID:              config.ID,
		Region:          config.Region,
		Capacity:        config.Capacity,
		Load:            0,
		APIHandlers:     make(map[string]api.Resolver),
		CacheEnabled:    config.CacheEnabled,
		Cache:           make(map[string]*CacheEntry),
		CacheTTL:        config.CacheTTL,
		LocalDB:         db.NewGoScaleDB(config.DBConfig),
		SyncInterval:    config.SyncInterval,
		LastSyncTime:    time.Now(),
		HealthStatus:    "healthy",
		Metrics:         &EdgeMetrics{},
		ParentAPI:       parentAPI,
		MaxConcurrent:   config.MaxConcurrent,
		RequestQueue:    make(chan *EdgeRequest, config.MaxConcurrent*10),
		CompressionLevel: config.CompressionLevel,
	}
	
	// Initialize worker pool
	node.WorkerPool = make([]*EdgeWorker, config.MaxConcurrent)
	for i := 0; i < config.MaxConcurrent; i++ {
		worker := &EdgeWorker{
			ID:         i,
			RequestChan: make(chan *EdgeRequest, 10),
			Node:       node,
			Active:     true,
		}
		node.WorkerPool[i] = worker
		go worker.Start()
	}
	
	// Start the request dispatcher
	go node.startDispatcher()
	
	// Start the sync process
	go node.startSyncProcess()
	
	return node
}

// Start starts the edge worker
func (w *EdgeWorker) Start() {
	for w.Active {
		select {
		case req := <-w.RequestChan:
			startTime := time.Now()
			
			// Process the request
			var result interface{}
			var err error
			
			// Check cache first if enabled
			if w.Node.CacheEnabled {
				cacheKey := fmt.Sprintf("%s:%v", req.Path, req.Params)
				w.Node.CacheMutex.RLock()
				if entry, ok := w.Node.Cache[cacheKey]; ok && time.Now().Before(entry.Expiration) {
					result = entry.Result
					w.Node.CacheMutex.RUnlock()
					w.Node.updateMetrics(startTime, true, true)
					req.ResultChan <- &EdgeResponse{Result: result, Error: nil}
					continue
				}
				w.Node.CacheMutex.RUnlock()
			}
			
			// Get the handler
			handler, ok := w.Node.APIHandlers[req.Path]
			if !ok {
				err = fmt.Errorf("no handler found for path %s", req.Path)
				w.Node.updateMetrics(startTime, false, false)
				req.ResultChan <- &EdgeResponse{Result: nil, Error: err}
				continue
			}
			
			// Execute the handler
			result, err = handler(req.Context, req.Params)
			
			// Cache the result if successful and caching is enabled
			if err == nil && w.Node.CacheEnabled {
				cacheKey := fmt.Sprintf("%s:%v", req.Path, req.Params)
				w.Node.CacheMutex.Lock()
				w.Node.Cache[cacheKey] = &CacheEntry{
					Path:       req.Path,
					Params:     req.Params,
					Result:     result,
					Expiration: time.Now().Add(w.Node.CacheTTL),
				}
				w.Node.CacheMutex.Unlock()
			}
			
			w.Node.updateMetrics(startTime, err == nil, false)
			req.ResultChan <- &EdgeResponse{Result: result, Error: err}
		}
	}
}

// startDispatcher starts the request dispatcher
func (n *EdgeNode) startDispatcher() {
	for req := range n.RequestQueue {
		// Find an available worker
		workerFound := false
		for _, worker := range n.WorkerPool {
			select {
			case worker.RequestChan <- req:
				workerFound = true
				break
			default:
				// Worker is busy, try the next one
			}
			
			if workerFound {
				break
			}
		}
		
		// If no worker is available, process the request in the main goroutine
		if !workerFound {
			startTime := time.Now()
			
			// Get the handler
			handler, ok := n.APIHandlers[req.Path]
			if !ok {
				err := fmt.Errorf("no handler found for path %s", req.Path)
				n.updateMetrics(startTime, false, false)
				req.ResultChan <- &EdgeResponse{Result: nil, Error: err}
				continue
			}
			
			// Execute the handler
			result, err := handler(req.Context, req.Params)
			n.updateMetrics(startTime, err == nil, false)
			req.ResultChan <- &EdgeResponse{Result: result, Error: err}
		}
	}
}

// startSyncProcess starts the sync process
func (n *EdgeNode) startSyncProcess() {
	ticker := time.NewTicker(n.SyncInterval)
	defer ticker.Stop()
	
	for range ticker.C {
		n.SyncWithParent()
	}
}

// SyncWithParent synchronizes the edge node with the parent API
func (n *EdgeNode) SyncWithParent() error {
	n.SyncMutex.Lock()
	defer n.SyncMutex.Unlock()
	
	// In a real implementation, this would sync data with the parent API
	n.LastSyncTime = time.Now()
	
	return nil
}

// RegisterHandler registers a handler for a specific path
func (n *EdgeNode) RegisterHandler(path string, handler api.Resolver) {
	n.APIHandlers[path] = handler
}

// ProcessRequest processes an API request
func (n *EdgeNode) ProcessRequest(ctx context.Context, path string, params map[string]interface{}) (interface{}, error) {
	// Create a request
	resultChan := make(chan *EdgeResponse, 1)
	req := &EdgeRequest{
		Path:       path,
		Params:     params,
		Context:    ctx,
		ResultChan: resultChan,
	}
	
	// Add the request to the queue
	n.RequestQueue <- req
	
	// Wait for the response
	resp := <-resultChan
	return resp.Result, resp.Error
}

// ServeHTTP implements the http.Handler interface
func (n *EdgeNode) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	
	// Parse the request
	var request struct {
		Path       string                 `json:"path"`
		Params     map[string]interface{} `json:"params"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*30)
	defer cancel()
	
	// Process the request
	result, err := n.ProcessRequest(ctx, request.Path, request.Params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		n.updateMetrics(startTime, false, false)
		return
	}
	
	// Return the result
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": result,
	})
	
	n.updateMetrics(startTime, true, false)
}

// updateMetrics updates the edge node metrics
func (n *EdgeNode) updateMetrics(startTime time.Time, success, cacheHit bool) {
	duration := time.Since(startTime).Seconds()
	
	n.Metrics.mutex.Lock()
	defer n.Metrics.mutex.Unlock()
	
	n.Metrics.RequestCount++
	n.Metrics.AvgResponseTime = (n.Metrics.AvgResponseTime*float64(n.Metrics.RequestCount-1) + duration) / float64(n.Metrics.RequestCount)
	
	if !success {
		n.Metrics.ErrorRate = (n.Metrics.ErrorRate*float64(n.Metrics.RequestCount-1) + 1) / float64(n.Metrics.RequestCount)
	} else {
		n.Metrics.ErrorRate = (n.Metrics.ErrorRate * float64(n.Metrics.RequestCount-1)) / float64(n.Metrics.RequestCount)
	}
	
	if cacheHit {
		n.Metrics.CacheHitRate = (n.Metrics.CacheHitRate*float64(n.Metrics.RequestCount-1) + 1) / float64(n.Metrics.RequestCount)
	} else {
		n.Metrics.CacheHitRate = (n.Metrics.CacheHitRate * float64(n.Metrics.RequestCount-1)) / float64(n.Metrics.RequestCount)
	}
}

// GetMetrics returns the current edge node metrics
func (n *EdgeNode) GetMetrics() *EdgeMetrics {
	n.Metrics.mutex.RLock()
	defer n.Metrics.mutex.RUnlock()
	
	return &EdgeMetrics{
		RequestCount:    n.Metrics.RequestCount,
		AvgResponseTime: n.Metrics.AvgResponseTime,
		ErrorRate:       n.Metrics.ErrorRate,
		CacheHitRate:    n.Metrics.CacheHitRate,
		CPUUsage:        n.Metrics.CPUUsage,
		MemoryUsage:     n.Metrics.MemoryUsage,
		NetworkIn:       n.Metrics.NetworkIn,
		NetworkOut:      n.Metrics.NetworkOut,
	}
}

// ClearCache clears the edge node cache
func (n *EdgeNode) ClearCache() {
	n.CacheMutex.Lock()
	defer n.CacheMutex.Unlock()
	
	n.Cache = make(map[string]*CacheEntry)
}

// Close closes the edge node and all its resources
func (n *EdgeNode) Close() error {
	// Stop all workers
	for _, worker := range n.WorkerPool {
		worker.Active = false
	}
	
	// Close the request queue
	close(n.RequestQueue)
	
	// Close the local database
	return n.LocalDB.Close()
}

// EdgeNetwork represents a network of edge nodes
type EdgeNetwork struct {
	Nodes           map[string]*EdgeNode
	LoadBalancer    *LoadBalancer
	HealthChecker   *HealthChecker
	SyncManager     *SyncManager
	ParentAPI       *api.GoScaleAPI
	mutex           sync.RWMutex
}

// LoadBalancer distributes requests across edge nodes
type LoadBalancer struct {
	Strategy        string
	Network         *EdgeNetwork
	RequestCounter  int64
	mutex           sync.Mutex
}

// HealthChecker monitors the health of edge nodes
type HealthChecker struct {
	Network         *EdgeNetwork
	CheckInterval   time.Duration
	LastCheckTime   time.Time
	mutex           sync.Mutex
}

// SyncManager manages synchronization between edge nodes
type SyncManager struct {
	Network         *EdgeNetwork
	SyncInterval    time.Duration
	LastSyncTime    time.Time
	mutex           sync.Mutex
}

// NewEdgeNetwork creates a new edge network
func NewEdgeNetwork(parentAPI *api.GoScaleAPI) *EdgeNetwork {
	network := &EdgeNetwork{
		Nodes:     make(map[string]*EdgeNode),
		ParentAPI: parentAPI,
	}
	
	// Create the load balancer
	network.LoadBalancer = &LoadBalancer{
		Strategy: "round-robin",
		Network:  network,
	}
	
	// Create the health checker
	network.HealthChecker = &HealthChecker{
		Network:       network,
		CheckInterval: time.Minute,
		LastCheckTime: time.Now(),
	}
	
	// Create the sync manager
	network.SyncManager = &SyncManager{
		Network:      network,
		SyncInterval: time.Minute * 15,
		LastSyncTime: time.Now(),
	}
	
	// Start the health checker
	go network.HealthChecker.Start()
	
	// Start the sync manager
	go network.SyncManager.Start()
	
	return network
}

// AddNode adds a node to the network
func (n *EdgeNetwork) AddNode(node *EdgeNode) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	
	n.Nodes[node.ID] = node
}

// RemoveNode removes a node from the network
func (n *EdgeNetwork) RemoveNode(nodeID string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	
	if node, ok := n.Nodes[nodeID]; ok {
		node.Close()
		delete(n.Nodes, nodeID)
	}
}

// GetNode returns a node by ID
func (n *EdgeNetwork) GetNode(nodeID string) (*EdgeNode, error) {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	
	node, ok := n.Nodes[nodeID]
	if !ok {
		return nil, fmt.Errorf("node %s not found", nodeID)
	}
	
	return node, nil
}

// ProcessRequest processes a request through the edge network
func (n *EdgeNetwork) ProcessRequest(ctx context.Context, path string, params map[string]interface{}) (interface{}, error) {
	// Get the best node for this request
	node, err := n.LoadBalancer.GetBestNode(path, params)
	if err != nil {
		return nil, err
	}
	
	// Process the request on the selected node
	return node.ProcessRequest(ctx, path, params)
}

// Start starts the health checker
func (h *HealthChecker) Start() {
	ticker := time.NewTicker(h.CheckInterval)
	defer ticker.Stop()
	
	for range ticker.C {
		h.CheckHealth()
	}
}

// CheckHealth checks the health of all nodes
func (h *HealthChecker) CheckHealth() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	h.Network.mutex.RLock()
	nodes := make([]*EdgeNode, 0, len(h.Network.Nodes))
	for _, node := range h.Network.Nodes {
		nodes = append(nodes, node)
	}
	h.Network.mutex.RUnlock()
	
	for _, node := range nodes {
		// In a real implementation, this would check the node's health
		// For now, we'll just set a random health status
		if node.Metrics.ErrorRate > 0.5 {
			node.HealthStatus = "unhealthy"
		} else {
			node.HealthStatus = "healthy"
		}
	}
	
	h.LastCheckTime = time.Now()
}

// Start starts the sync manager
func (s *SyncManager) Start() {
	ticker := time.NewTicker(s.SyncInterval)
	defer ticker.Stop()
	
	for range ticker.C {
		s.SyncNodes()
	}
}

// SyncNodes synchronizes all nodes
func (s *SyncManager) SyncNodes() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	s.Network.mutex.RLock()
	nodes := make([]*EdgeNode, 0, len(s.Network.Nodes))
	for _, node := range s.Network.Nodes {
		nodes = append(nodes, node)
	}
	s.Network.mutex.RUnlock()
	
	for _, node := range nodes {
		node.SyncWithParent()
	}
	
	s.LastSyncTime = time.Now()
}

// GetBestNode returns the best node for a request
func (l *LoadBalancer) GetBestNode(path string, params map[string]interface{}) (*EdgeNode, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	
	l.Network.mutex.RLock()
	defer l.Network.mutex.RUnlock()
	
	if len(l.Network.Nodes) == 0 {
		return nil, errors.New("no nodes available")
	}
	
	// Get all healthy nodes
	healthyNodes := make([]*EdgeNode, 0)
	for _, node := range l.Network.Nodes {
		if node.HealthStatus == "healthy" {
			healthyNodes = append(healthyNodes, node)
		}
	}
	
	if len(healthyNodes) == 0 {
		return nil, errors.New("no healthy nodes available")
	}
	
	// Select a node based on the strategy
	var selectedNode *EdgeNode
	
	switch l.Strategy {
	case "round-robin":
		// Simple round-robin
		l.RequestCounter++
		selectedNode = healthyNodes[l.RequestCounter%int64(len(healthyNodes))]
	case "least-loaded":
		// Select the node with the lowest load
		minLoad := healthyNodes[0].Load
		selectedNode = healthyNodes[0]
		
		for _, node := range healthyNodes {
			if node.Load < minLoad {
				minLoad = node.Load
				selectedNode = node
			}
		}
	case "fastest":
		// Select the node with the lowest average response time
		minTime := healthyNodes[0].Metrics.AvgResponseTime
		selectedNode = healthyNodes[0]
		
		for _, node := range healthyNodes {
			if node.Metrics.AvgResponseTime < minTime {
				minTime = node.Metrics.AvgResponseTime
				selectedNode = node
			}
		}
	default:
		// Default to round-robin
		l.RequestCounter++
		selectedNode = healthyNodes[l.RequestCounter%int64(len(healthyNodes))]
	}
	
	return selectedNode, nil
}