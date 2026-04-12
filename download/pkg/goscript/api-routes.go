package goscript

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// HTTPMethod represents allowed HTTP methods
type HTTPMethod string

const (
	MethodGet     HTTPMethod = "GET"
	MethodPost    HTTPMethod = "POST"
	MethodPut     HTTPMethod = "PUT"
	MethodPatch   HTTPMethod = "PATCH"
	MethodDelete  HTTPMethod = "DELETE"
	MethodOptions HTTPMethod = "OPTIONS"
)

// APIError represents an API error response
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return e.Message
}

// APIContext provides request context, params, and helpers
type APIContext struct {
	Request  *http.Request
	Response http.ResponseWriter
	Params   map[string]string
	Query    map[string]string
	Body     []byte
	Ctx      context.Context
}

// JSON writes a JSON response
func (ac *APIContext) JSON(status int, data interface{}) {
	ac.Response.Header().Set("Content-Type", "application/json")
	ac.Response.WriteHeader(status)
	json.NewEncoder(ac.Response).Encode(data)
}

// APIHandler is the function signature for API route handlers
type APIHandler func(*APIContext) (interface{}, error)

// APIRoute defines an API route with method, path, handler, and middleware
type APIRoute struct {
	Method      HTTPMethod
	Pattern     string
	Handler     APIHandler
	Middleware  []APIMiddlewareFunc
	RateLimit   *RateLimitConfig
	Auth        AuthConfig
	Description string
}

// APIMiddlewareFunc processes the request/response chain
type APIMiddlewareFunc func(*APIContext, func() (interface{}, error)) (interface{}, error)

// RateLimitConfig configures rate limiting per route
type RateLimitConfig struct {
	Requests int
	Window   time.Duration
}

// AuthConfig configures authentication for a route
type AuthConfig struct {
	Required bool
	Schemes  []string // "bearer", "basic", "cookie", "apikey"
	Roles    []string
}

// APIRouter manages API routes with convention-based patterns
type APIRouter struct {
	routes       []*APIRoute
	middleware   []APIMiddlewareFunc
	notFound     APIHandler
	errorHandler func(*APIContext, error)
}

// NewAPIRouter creates a new API router
func NewAPIRouter() *APIRouter {
	return &APIRouter{
		routes:       make([]*APIRoute, 0),
		middleware:   make([]APIMiddlewareFunc, 0),
		errorHandler: defaultAPIErrorHandler,
	}
}

// RouteOption configures an API route
type RouteOption func(*APIRoute)

// WithRateLimit adds rate limiting
func WithRateLimit(requests int, window time.Duration) RouteOption {
	return func(r *APIRoute) {
		r.RateLimit = &RateLimitConfig{Requests: requests, Window: window}
	}
}

// WithAuth adds authentication requirements
func WithAuth(schemes []string, roles []string) RouteOption {
	return func(r *APIRoute) {
		r.Auth = AuthConfig{Required: true, Schemes: schemes, Roles: roles}
	}
}

// Register registers a route with the given method, pattern, and handler
func (ar *APIRouter) Register(method HTTPMethod, pattern string, handler APIHandler, opts ...RouteOption) {
	route := &APIRoute{
		Method:  method,
		Pattern: pattern,
		Handler: handler,
	}
	for _, opt := range opts {
		opt(route)
	}
	ar.routes = append(ar.routes, route)
}

// GET registers a GET handler
func (ar *APIRouter) GET(pattern string, handler APIHandler, opts ...RouteOption) {
	ar.Register(MethodGet, pattern, handler, opts...)
}

// POST registers a POST handler
func (ar *APIRouter) POST(pattern string, handler APIHandler, opts ...RouteOption) {
	ar.Register(MethodPost, pattern, handler, opts...)
}

// DELETE registers a DELETE handler
func (ar *APIRouter) DELETE(pattern string, handler APIHandler, opts ...RouteOption) {
	ar.Register(MethodDelete, pattern, handler, opts...)
}

// ServeHTTP implements http.Handler
func (ar *APIRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := make(map[string]string)
	for k, v := range r.URL.Query() {
		if len(v) > 0 {
			query[k] = v[0]
		}
	}

	body, _ := io.ReadAll(r.Body)

	ctx := &APIContext{
		Request:  r,
		Response: w,
		Query:    query,
		Body:     body,
		Ctx:      r.Context(),
		Params:   make(map[string]string),
	}

	for _, route := range ar.routes {
		if !ar.methodMatches(route.Method, r.Method) {
			continue
		}
		params, ok := ar.matchPattern(route.Pattern, r.URL.Path)
		if !ok {
			continue
		}
		ctx.Params = params
		ar.executeHandler(route, ctx)
		return
	}

	if ar.notFound != nil {
		result, _ := ar.notFound(ctx)
		ar.writeJSON(w, 404, result)
	} else {
		ar.writeJSON(w, 404, map[string]string{"error": "not found"})
	}
}

// APICORSMiddleware handles CORS for API routes (renamed to avoid conflict with pipeline CORSMiddleware)
func APICORSMiddleware(allowedOrigins []string) APIMiddlewareFunc {
	originMap := make(map[string]bool)
	for _, o := range allowedOrigins {
		originMap[o] = true
	}

	return func(ctx *APIContext, next func() (interface{}, error)) (interface{}, error) {
		origin := ctx.Request.Header.Get("Origin")
		if originMap["*"] || originMap[origin] {
			ctx.Response.Header().Set("Access-Control-Allow-Origin", origin)
			ctx.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			ctx.Response.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			ctx.Response.Header().Set("Access-Control-Max-Age", "86400")
		}
		if ctx.Request.Method == "OPTIONS" {
			ctx.Response.WriteHeader(204)
			return nil, nil
		}
		return next()
	}
}

// APIRateLimitMiddleware tracks request counts per IP for API routes
func APIRateLimitMiddleware(requests int, window time.Duration) APIMiddlewareFunc {
	type tracker struct {
		counts map[string]int
		mu     sync.RWMutex
	}
	t := &tracker{counts: make(map[string]int)}

	go func() {
		for {
			time.Sleep(window)
			t.mu.Lock()
			t.counts = make(map[string]int)
			t.mu.Unlock()
		}
	}()

	return func(ctx *APIContext, next func() (interface{}, error)) (interface{}, error) {
		ip := strings.Split(ctx.Request.RemoteAddr, ":")[0]
		t.mu.RLock()
		count := t.counts[ip]
		t.mu.RUnlock()
		if count >= requests {
			return nil, &APIError{Code: 429, Message: "rate limit exceeded"}
		}
		t.mu.Lock()
		t.counts[ip]++
		t.mu.Unlock()
		return next()
	}
}

// RequestLogger middleware logs API requests
func RequestLogger(logger func(format string, args ...interface{})) APIMiddlewareFunc {
	return func(ctx *APIContext, next func() (interface{}, error)) (interface{}, error) {
		start := time.Now()
		result, err := next()
		duration := time.Since(start)
		logger("[API] %s %s -> %d (%s)", ctx.Request.Method, ctx.Request.URL.Path,
			ctx.Response.StatusCode(), duration)
		return result, err
	}
}

// defaultAPIErrorHandler handles API errors with a JSON response
func defaultAPIErrorHandler(ctx *APIContext, err error) {
	if apiErr, ok := err.(*APIError); ok {
		ctx.Response.Header().Set("Content-Type", "application/json")
		ctx.Response.WriteHeader(apiErr.Code)
		json.NewEncoder(ctx.Response).Encode(apiErr)
		return
	}
	ctx.Response.Header().Set("Content-Type", "application/json")
	ctx.Response.WriteHeader(500)
	json.NewEncoder(ctx.Response).Encode(map[string]string{
		"error": "internal server error",
	})
}

// methodMatches checks if the route method matches the request method
func (ar *APIRouter) methodMatches(route HTTPMethod, method string) bool {
	return string(route) == method
}

// matchPattern performs simple pattern matching with :param and *wildcard support
func (ar *APIRouter) matchPattern(pattern, path string) (map[string]string, bool) {
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	params := make(map[string]string)

	for i := 0; i < len(patternParts); i++ {
		if i >= len(pathParts) {
			return nil, false
		}

		pp := patternParts[i]
		pap := pathParts[i]

		if strings.HasPrefix(pp, ":") {
			params[pp[1:]] = pap
			continue
		}

		if strings.HasPrefix(pp, "*") {
			// wildcard matches rest
			params[pp[1:]] = strings.Join(pathParts[i:], "/")
			return params, true
		}

		if pp != pap {
			return nil, false
		}
	}

	if len(patternParts) != len(pathParts) {
		return nil, false
	}

	return params, true
}

// executeHandler runs the route handler with middleware
func (ar *APIRouter) executeHandler(route *APIRoute, ctx *APIContext) {
	handler := route.Handler

	// Wrap with route-specific middleware (in reverse order)
	for i := len(route.Middleware) - 1; i >= 0; i-- {
		mw := route.Middleware[i]
		next := handler
		handler = func(nextHandler APIHandler, middleware APIMiddlewareFunc) APIHandler {
			return func(ctx *APIContext) (interface{}, error) {
				return middleware(ctx, func() (interface{}, error) {
					return nextHandler(ctx)
				})
			}
		}(next, mw)
	}

	// Wrap with global middleware (in reverse order)
	for i := len(ar.middleware) - 1; i >= 0; i-- {
		mw := ar.middleware[i]
		next := handler
		handler = func(nextHandler APIHandler, middleware APIMiddlewareFunc) APIHandler {
			return func(ctx *APIContext) (interface{}, error) {
				return middleware(ctx, func() (interface{}, error) {
					return nextHandler(ctx)
				})
			}
		}(next, mw)
	}

	result, err := handler(ctx)
	if err != nil {
		ar.errorHandler(ctx, err)
		return
	}

	ar.writeJSON(ctx.Response, 200, result)
}

// writeJSON writes a JSON response
func (ar *APIRouter) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
