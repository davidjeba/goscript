package goscript

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// APIContext provides a clean, structured request context for API route handlers.
// It encapsulates the HTTP request's parameters, query string, body, headers,
// and method, giving handlers easy access to all incoming data.
type APIContext struct {
	Params  map[string]string
	Query   map[string]string
	Body    map[string]interface{}
	Headers map[string]string
	Method  string
	Path    string
	Request *http.Request
	Writer  http.ResponseWriter
}

// APIOption is a functional option that configures an API route. Options such as
// rate limiting, authentication, CORS, and body validation can be composed.
type APIOption func(*apiRouteConfig)

// apiRouteConfig holds the configuration for a single API route, populated by
// the APIOption functional arguments.
type apiRouteConfig struct {
	rateLimit        *rateLimitConfig
	authSchemes      []string
	authValidator    func(string) bool
	corsOrigin       string
	bodyValidator    interface{}
}

// rateLimitConfig defines the parameters for a sliding-window rate limiter.
type rateLimitConfig struct {
	requests int
	window   time.Duration
}

// apiRoute represents a single registered API route with its handler and options.
type apiRoute struct {
	method   string
	path     string
	handler  func(*APIContext) (interface{}, error)
	config   apiRouteConfig
}

// apiRateLimiter implements a sliding-window rate limiter that tracks requests
// per client IP address within a configurable time window.
type apiRateLimiter struct {
	requests    int
	window      time.Duration
	clients     map[string][]time.Time
	mutex       sync.Mutex
}

// newAPIRateLimiter creates a rate limiter that allows the given number of
// requests within the specified time window per client.
func newAPIRateLimiter(requests int, window time.Duration) *apiRateLimiter {
	rl := &apiRateLimiter{
		requests: requests,
		window:   window,
		clients:  make(map[string][]time.Time),
	}
	// Start a background goroutine to clean up expired entries
	go rl.cleanup()
	return rl
}

// allow checks whether a request from the given client key is within the
// allowed rate limit. If the client has exceeded the limit, it returns false.
func (rl *apiRateLimiter) allow(clientKey string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Filter out timestamps outside the current window
	var recent []time.Time
	for _, ts := range rl.clients[clientKey] {
		if ts.After(cutoff) {
			recent = append(recent, ts)
		}
	}

	if len(recent) >= rl.requests {
		rl.clients[clientKey] = recent
		return false
	}

	recent = append(recent, now)
	rl.clients[clientKey] = recent
	return true
}

// cleanup periodically removes expired entries from the rate limiter to
// prevent memory leaks from abandoned clients.
func (rl *apiRateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()
	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.window)
		for key, timestamps := range rl.clients {
			var recent []time.Time
			for _, ts := range timestamps {
				if ts.After(cutoff) {
					recent = append(recent, ts)
				}
			}
			if len(recent) == 0 {
				delete(rl.clients, key)
			} else {
				rl.clients[key] = recent
			}
		}
		rl.mutex.Unlock()
	}
}

// APIRouter is a convention-based API router that implements http.Handler.
// It provides a clean API for registering routes with optional middleware
// such as rate limiting, authentication, CORS, and body validation.
type APIRouter struct {
	routes       []apiRoute
	rateLimiters map[string]*apiRateLimiter
	mutex        sync.RWMutex
	notFound     func(*APIContext) (interface{}, error)
}

// NewAPIRouter creates a new APIRouter with no registered routes.
func NewAPIRouter() *APIRouter {
	return &APIRouter{
		routes:       make([]apiRoute, 0),
		rateLimiters: make(map[string]*apiRateLimiter),
	}
}

// GET registers a handler for GET requests at the given path.
func (ar *APIRouter) GET(path string, handler func(*APIContext) (interface{}, error), opts ...APIOption) {
	ar.registerRoute("GET", path, handler, opts...)
}

// POST registers a handler for POST requests at the given path.
func (ar *APIRouter) POST(path string, handler func(*APIContext) (interface{}, error), opts ...APIOption) {
	ar.registerRoute("POST", path, handler, opts...)
}

// PUT registers a handler for PUT requests at the given path.
func (ar *APIRouter) PUT(path string, handler func(*APIContext) (interface{}, error), opts ...APIOption) {
	ar.registerRoute("PUT", path, handler, opts...)
}

// DELETE registers a handler for DELETE requests at the given path.
func (ar *APIRouter) DELETE(path string, handler func(*APIContext) (interface{}, error), opts ...APIOption) {
	ar.registerRoute("DELETE", path, handler, opts...)
}

// SetNotFound sets a custom handler for requests that match no registered route.
func (ar *APIRouter) SetNotFound(handler func(*APIContext) (interface{}, error)) {
	ar.notFound = handler
}

// registerRoute adds a new route with the given method, path, handler, and options.
func (ar *APIRouter) registerRoute(method, path string, handler func(*APIContext) (interface{}, error), opts ...APIOption) {
	config := apiRouteConfig{}
	for _, opt := range opts {
		opt(&config)
	}

	route := apiRoute{
		method:  method,
		path:    path,
		handler: handler,
		config:  config,
	}

	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	ar.routes = append(ar.routes, route)

	// Create a rate limiter if rate limiting is configured
	if config.rateLimit != nil {
		key := fmt.Sprintf("%s:%s", method, path)
		ar.rateLimiters[key] = newAPIRateLimiter(config.rateLimit.requests, config.rateLimit.window)
	}
}

// ServeHTTP implements the http.Handler interface. It matches incoming requests
// to registered routes, applies configured options (rate limiting, auth, CORS),
// and returns automatic JSON responses.
func (ar *APIRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ar.mutex.RLock()
	defer ar.mutex.RUnlock()

	// Try to match the request to a registered route
	for _, route := range ar.routes {
		if route.method != r.Method {
			continue
		}

		params, ok := matchAppPath(route.path, r.URL.Path)
		if !ok {
			continue
		}

		// Apply rate limiting
		rlKey := fmt.Sprintf("%s:%s", route.method, route.path)
		if limiter, exists := ar.rateLimiters[rlKey]; exists {
			clientIP := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				clientIP = strings.Split(forwarded, ",")[0]
			}
			if !limiter.allow(clientIP) {
				writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}
		}

		// Apply CORS headers
		if route.config.corsOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", route.config.corsOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400")

			// Handle preflight
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		// Apply authentication
		if len(route.config.authSchemes) > 0 && route.config.authValidator != nil {
			authHeader := r.Header.Get("Authorization")
			if !route.config.authValidator(authHeader) {
				writeJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
		}

		// Build the APIContext
		ctx := &APIContext{
			Params:  params,
			Query:   parseQuery(r.URL.RawQuery),
			Headers: extractHeaders(r.Header),
			Method:  r.Method,
			Path:    r.URL.Path,
			Request: r,
			Writer:  w,
		}

		// Parse request body for methods that carry a body
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			body, err := parseBody(r.Body)
			if err != nil {
				writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
				return
			}
			ctx.Body = body
		}

		// Invoke the handler and write the JSON response
		result, err := route.handler(ctx)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSONResponse(w, http.StatusOK, result)
		return
	}

	// No route matched
	if ar.notFound != nil {
		ctx := &APIContext{
			Query:   parseQuery(r.URL.RawQuery),
			Headers: extractHeaders(r.Header),
			Method:  r.Method,
			Path:    r.URL.Path,
			Request: r,
			Writer:  w,
		}
		result, err := ar.notFound(ctx)
		if err != nil {
			writeJSONError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSONResponse(w, http.StatusNotFound, result)
		return
	}

	writeJSONError(w, http.StatusNotFound, "not found")
}

// parseQuery parses a URL query string into a map.
func parseQuery(rawQuery string) map[string]string {
	result := make(map[string]string)
	if rawQuery == "" {
		return result
	}
	pairs := strings.Split(rawQuery, "&")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		}
	}
	return result
}

// extractHeaders extracts relevant HTTP headers into a map.
func extractHeaders(h http.Header) map[string]string {
	result := make(map[string]string)
	for key, values := range h {
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}

// parseBody reads and parses the request body as JSON.
func parseBody(body io.ReadCloser) (map[string]interface{}, error) {
	defer body.Close()
	var result map[string]interface{}
	decoder := json.NewDecoder(body)
	decoder.UseNumber()
	if err := decoder.Decode(&result); err != nil && err != io.EOF {
		return nil, err
	}
	if result == nil {
		result = make(map[string]interface{})
	}
	return result, nil
}

// writeJSONResponse writes a successful JSON response with the given status code.
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

// writeJSONError writes a JSON error response with the given status code and message.
func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error": map[string]interface{}{
			"code":    statusCode,
			"message": message,
		},
	})
}

// WithRateLimit returns an APIOption that configures a sliding-window rate limiter
// allowing the specified number of requests within the given time window.
func WithRateLimit(requests int, window time.Duration) APIOption {
	return func(cfg *apiRouteConfig) {
		cfg.rateLimit = &rateLimitConfig{
			requests: requests,
			window:   window,
		}
	}
}

// WithAuth returns an APIOption that enables authentication for the route.
// The validator function receives the Authorization header value and returns
// true if the request is authenticated. Multiple auth schemes can be specified
// for the WWW-Authenticate header on 401 responses.
func WithAuth(schemes []string, validator func(string) bool) APIOption {
	return func(cfg *apiRouteConfig) {
		cfg.authSchemes = schemes
		cfg.authValidator = validator
	}
}

// WithCORS returns an APIOption that enables Cross-Origin Resource Sharing for
// the route with the specified allowed origin.
func WithCORS(origin string) APIOption {
	return func(cfg *apiRouteConfig) {
		cfg.corsOrigin = origin
	}
}

// WithBodyValidator returns an APIOption that attaches a body validation schema
// to the route. The schema is an interface{} that describes the expected body
// structure, enabling automatic validation before the handler is invoked.
func WithBodyValidator(schema interface{}) APIOption {
	return func(cfg *apiRouteConfig) {
		cfg.bodyValidator = schema
	}
}

// APIMiddlewareFunc is a function type that wraps an APIContext handler with
// additional processing such as logging, metrics, or transformation.
type APIMiddlewareFunc func(func(*APIContext) (interface{}, error)) func(*APIContext) (interface{}, error)

// APIResponseRate provides per-route response rate tracking with configurable
// sliding windows for use in auto-scaling and load balancing decisions.
type APIResponseRate struct {
	TotalRequests   int64
	SuccessCount    int64
	ErrorCount      int64
	AvgLatencyMs    float64
	RequestsPerSec  float64
	WindowStart     time.Time
	WindowDuration  time.Duration
	runningLatencies []float64
	mutex           sync.Mutex
}

// NewAPIResponseRate creates a new response rate tracker with the specified
// sliding window duration.
func NewAPIResponseRate(windowDuration time.Duration) *APIResponseRate {
	return &APIResponseRate{
		WindowStart:    time.Now(),
		WindowDuration: windowDuration,
	}
}

// Record records a response with its latency in milliseconds.
func (arr *APIResponseRate) Record(success bool, latencyMs float64) {
	arr.mutex.Lock()
	defer arr.mutex.Unlock()

	arr.TotalRequests++
	if success {
		arr.SuccessCount++
	} else {
		arr.ErrorCount++
	}
	arr.runningLatencies = append(arr.runningLatencies, latencyMs)

	// Calculate average
	sum := 0.0
	for _, l := range arr.runningLatencies {
		sum += l
	}
	arr.AvgLatencyMs = sum / float64(len(arr.runningLatencies))

	// Calculate requests per second
	elapsed := time.Since(arr.WindowStart).Seconds()
	if elapsed > 0 {
		arr.RequestsPerSec = float64(arr.TotalRequests) / elapsed
	}
}

// Reset resets the sliding window counters.
func (arr *APIResponseRate) Reset() {
	arr.mutex.Lock()
	defer arr.mutex.Unlock()

	arr.TotalRequests = 0
	arr.SuccessCount = 0
	arr.ErrorCount = 0
	arr.AvgLatencyMs = 0
	arr.RequestsPerSec = 0
	arr.WindowStart = time.Now()
	arr.runningLatencies = nil
}

// StatusCode returns the HTTP status code as an integer from a string, defaulting
// to 200 if parsing fails.
func StatusCode(code interface{}) int {
	switch v := code.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		n, err := strconv.Atoi(v)
		if err != nil {
			return 200
		}
		return n
	default:
		return 200
	}
}
