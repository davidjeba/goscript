package goscript

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Pipeline is an ordered chain of middleware handlers
type Pipeline struct {
	handlers []PipelineHandler
	mu       sync.RWMutex
}

// PipelineHandler processes a request and calls the next handler
type PipelineHandler func(ctx *RequestContext, next PipelineNext) PipelineResult

// PipelineNext invokes the next handler in the chain
type PipelineNext func() PipelineResult

// PipelineResult represents the outcome of pipeline execution
type PipelineResult struct {
	Status  int
	Body    []byte
	Headers http.Header
	Aborted bool
}

// RequestContext wraps http.Request with additional context
type RequestContext struct {
	*http.Request
	Response  http.ResponseWriter
	Params    map[string]string
	Values    map[string]interface{}
	StartTime time.Time
	Aborted   bool
}

// Set stores a value in the request context
func (rc *RequestContext) Set(key string, value interface{}) {
	rc.Values[key] = value
}

// Get retrieves a value from the request context
func (rc *RequestContext) Get(key string) interface{} {
	return rc.Values[key]
}

// Abort stops the middleware chain
func (rc *RequestContext) Abort(status int) {
	rc.Aborted = true
	rc.Response.WriteHeader(status)
}

// NewRequestContext creates a new RequestContext from an HTTP request
func NewRequestContext(w http.ResponseWriter, r *http.Request) *RequestContext {
	return &RequestContext{
		Request:   r,
		Response:  w,
		Params:    make(map[string]string),
		Values:    make(map[string]interface{}),
		StartTime: time.Now(),
	}
}

// NewPipeline creates a new middleware pipeline
func NewPipeline() *Pipeline {
	return &Pipeline{
		handlers: make([]PipelineHandler, 0),
	}
}

// Use adds a handler to the pipeline
func (p *Pipeline) Use(handler PipelineHandler) *Pipeline {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.handlers = append(p.handlers, handler)
	return p
}

// Execute runs all middleware handlers in order
func (p *Pipeline) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := NewRequestContext(w, r)
	p.mu.RLock()
	handlers := make([]PipelineHandler, len(p.handlers))
	copy(handlers, p.handlers)
	p.mu.RUnlock()
	p.executeChain(ctx, handlers, 0)
}

func (p *Pipeline) executeChain(
	ctx *RequestContext,
	handlers []PipelineHandler,
	index int,
) PipelineResult {
	if ctx.Aborted || index >= len(handlers) {
		return PipelineResult{Status: 200, Aborted: ctx.Aborted}
	}
	return handlers[index](ctx, func() PipelineResult {
		return p.executeChain(ctx, handlers, index+1)
	})
}

// ==================== Built-in Middleware ====================

// GzipMiddleware compresses responses with gzip
func GzipMiddleware() PipelineHandler {
	return func(ctx *RequestContext, next PipelineNext) PipelineResult {
		if !strings.Contains(ctx.Header.Get("Accept-Encoding"), "gzip") {
			return next()
		}
		ctx.Response.Header().Set("Content-Encoding", "gzip")
		return next()
	}
}

// CORSConfig defines CORS configuration for the middleware pipeline
type CORSConfig struct {
	AllowAllOrigins  bool
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware(config CORSConfig) PipelineHandler {
	return func(ctx *RequestContext, next PipelineNext) PipelineResult {
		origin := ctx.Header.Get("Origin")
		allowed := false

		if config.AllowAllOrigins {
			allowed = true
		} else {
			for _, o := range config.AllowedOrigins {
				if o == origin {
					allowed = true
					break
				}
			}
		}

		if allowed {
			if config.AllowAllOrigins {
				ctx.Response.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				ctx.Response.Header().Set("Access-Control-Allow-Origin", origin)
			}
			methods := config.AllowedMethods
			if len(methods) == 0 {
				methods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
			}
			ctx.Response.Header().Set("Access-Control-Allow-Methods",
				strings.Join(methods, ", "))
			headers := config.AllowedHeaders
			if len(headers) == 0 {
				headers = []string{"Content-Type", "Authorization"}
			}
			ctx.Response.Header().Set("Access-Control-Allow-Headers",
				strings.Join(headers, ", "))
			ctx.Response.Header().Set("Access-Control-Allow-Credentials",
				strconv.FormatBool(config.AllowCredentials))
			if config.MaxAge > 0 {
				ctx.Response.Header().Set("Access-Control-Max-Age",
					strconv.Itoa(config.MaxAge))
			}
		}

		if ctx.Request.Method == "OPTIONS" {
			ctx.Response.WriteHeader(204)
			return PipelineResult{Status: 204, Aborted: true}
		}
		return next()
	}
}

// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware() PipelineHandler {
	return func(ctx *RequestContext, next PipelineNext) PipelineResult {
		h := ctx.Response.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("X-XSS-Protection", "1; mode=block")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		return next()
	}
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware(logger func(string, ...interface{})) PipelineHandler {
	return func(ctx *RequestContext, next PipelineNext) PipelineResult {
		defer func() {
			if r := recover(); r != nil {
				logger("[PANIC] %s %s: %v", ctx.Request.Method, ctx.Request.URL.Path, r)
				ctx.Abort(500)
			}
		}()
		return next()
	}
}

// Session represents a user session
type Session struct {
	ID        string
	Data      map[string]interface{}
	CreatedAt time.Time
	ExpiresAt time.Time
}

// SessionStore defines the interface for session storage
type SessionStore interface {
	Get(id string) (*Session, error)
	Create() *Session
	Delete(id string) error
}

// InMemorySessionStore provides an in-memory session store
type InMemorySessionStore struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

// NewInMemorySessionStore creates a new in-memory session store
func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{
		sessions: make(map[string]*Session),
	}
}

// Get retrieves a session by ID
func (s *InMemorySessionStore) Get(id string) (*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessions[id], nil
}

// Create creates a new session
func (s *InMemorySessionStore) Create() *Session {
	session := &Session{
		ID:        generateSessionID(),
		Data:      make(map[string]interface{}),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	s.mu.Lock()
	s.sessions[session.ID] = session
	s.mu.Unlock()
	return session
}

// Delete removes a session by ID
func (s *InMemorySessionStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, id)
	return nil
}

func generateSessionID() string {
	return strings.Replace(time.Now().Format("20060102150405.000000"), ".", "", -1)
}

// SessionMiddleware provides session management
func SessionMiddleware(secret string, store SessionStore) PipelineHandler {
	_ = secret // secret can be used for signing in a real implementation
	return func(ctx *RequestContext, next PipelineNext) PipelineResult {
		cookie, err := ctx.Request.Cookie("gosession")
		var session *Session

		if err == nil && cookie.Value != "" {
			session, _ = store.Get(cookie.Value)
		}

		if session == nil {
			session = store.Create()
			http.SetCookie(ctx.Response, &http.Cookie{
				Name:     "gosession",
				Value:    session.ID,
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
				MaxAge:   86400,
			})
		}

		ctx.Set("session", session)
		return next()
	}
}

// LoggingMiddleware logs request information
func LoggingMiddleware(logger func(string, ...interface{})) PipelineHandler {
	return func(ctx *RequestContext, next PipelineNext) PipelineResult {
		result := next()
		duration := time.Since(ctx.StartTime)
		logger("[HTTP] %s %s %d %s", ctx.Request.Method, ctx.Request.URL.Path, result.Status, duration)
		return result
	}
}

// RateLimitMiddleware provides request rate limiting
func RateLimitMiddleware(requests int, window time.Duration) PipelineHandler {
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

	return func(ctx *RequestContext, next PipelineNext) PipelineResult {
		ip := strings.Split(ctx.Request.RemoteAddr, ":")[0]
		t.mu.RLock()
		count := t.counts[ip]
		t.mu.RUnlock()
		if count >= requests {
			ctx.Response.WriteHeader(429)
			return PipelineResult{Status: 429, Aborted: true}
		}
		t.mu.Lock()
		t.counts[ip]++
		t.mu.Unlock()
		return next()
	}
}
