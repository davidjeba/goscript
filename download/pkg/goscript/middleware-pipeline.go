package goscript

import (
	"compress/gzip"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// CORSConfig holds Cross-Origin Resource Sharing configuration. It specifies
// which origins, methods, and headers are permitted in cross-origin requests.
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig returns a sensible default CORS configuration that allows
// common methods and headers while being restrictive on origins.
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"X-Request-ID", "X-Response-Time"},
		AllowCredentials: false,
		MaxAge:           86400,
	}
}

// SessionStore defines the interface for session storage backends. Implementations
// must be safe for concurrent use.
type SessionStore interface {
	Get(sessionID string) (map[string]interface{}, error)
	Set(sessionID string, data map[string]interface{}) error
	Delete(sessionID string) error
}

// InMemorySessionStore is an in-memory implementation of SessionStore suitable for
// development and single-server deployments. Sessions expire after a configurable TTL.
type InMemorySessionStore struct {
	sessions map[string]*sessionEntry
	ttl      time.Duration
	mutex    sync.RWMutex
}

type sessionEntry struct {
	data      map[string]interface{}
	createdAt time.Time
}

// NewInMemorySessionStore creates a new InMemorySessionStore with the given TTL.
// Sessions that have not been accessed within the TTL are automatically cleaned up.
func NewInMemorySessionStore(ttl time.Duration) *InMemorySessionStore {
	store := &InMemorySessionStore{
		sessions: make(map[string]*sessionEntry),
		ttl:      ttl,
	}
	go store.cleanupExpired()
	return store
}

// Get retrieves the session data for the given session ID. Returns an error if
// the session does not exist or has expired.
func (s *InMemorySessionStore) Get(sessionID string) (map[string]interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	entry, exists := s.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}
	if time.Since(entry.createdAt) > s.ttl {
		return nil, fmt.Errorf("session expired")
	}

	// Return a copy to prevent external mutation
	result := make(map[string]interface{}, len(entry.data))
	for k, v := range entry.data {
		result[k] = v
	}
	return result, nil
}

// Set stores or updates the session data for the given session ID. If the session
// already exists, its TTL is refreshed.
func (s *InMemorySessionStore) Set(sessionID string, data map[string]interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	copy := make(map[string]interface{}, len(data))
	for k, v := range data {
		copy[k] = v
	}
	s.sessions[sessionID] = &sessionEntry{
		data:      copy,
		createdAt: time.Now(),
	}
	return nil
}

// Delete removes the session with the given ID.
func (s *InMemorySessionStore) Delete(sessionID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.sessions, sessionID)
	return nil
}

// cleanupExpired periodically removes expired sessions from memory.
func (s *InMemorySessionStore) cleanupExpired() {
	ticker := time.NewTicker(s.ttl / 2)
	defer ticker.Stop()
	for range ticker.C {
		s.mutex.Lock()
		now := time.Now()
		for id, entry := range s.sessions {
			if time.Since(entry.createdAt) > s.ttl {
				delete(s.sessions, id)
			}
		}
		s.mutex.Unlock()
	}
}

// Pipeline is a composable middleware pipeline that applies a chain of middleware
// functions around a final http.Handler. Middleware is applied in the order they
// are added, with the first middleware being the outermost wrapper.
type Pipeline struct {
	middlewares []func(http.Handler) http.Handler
}

// NewPipeline creates a new empty middleware Pipeline.
func NewPipeline() *Pipeline {
	return &Pipeline{
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

// Use adds a middleware function to the pipeline. The middleware wraps the next
// handler in the chain, enabling pre-processing and post-processing of requests.
func (p *Pipeline) Use(middleware func(http.Handler) http.Handler) *Pipeline {
	p.middlewares = append(p.middlewares, middleware)
	return p
}

// Execute runs the middleware pipeline against the final handler, writing the
// result to the provided http.ResponseWriter. If no middleware is registered,
// the handler is invoked directly.
func (p *Pipeline) Execute(w http.ResponseWriter, r *http.Request, handler http.Handler) {
	current := handler
	for i := len(p.middlewares) - 1; i >= 0; i-- {
		current = p.middlewares[i](current)
	}
	current.ServeHTTP(w, r)
}

// GzipMiddleware returns a middleware that compresses HTTP responses using gzip
// encoding. It only compresses responses for clients that indicate support via
// the Accept-Encoding header, and skips compression for responses smaller than
// 512 bytes where the overhead outweighs the benefit.
func GzipMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			gz := gzip.NewWriter(w)
			defer gz.Close()

			gzw := &gzipResponseWriter{
				Writer:         gz,
				ResponseWriter: w,
			}

			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Del("Content-Length")

			next.ServeHTTP(gzw, r)
		})
	}
}

// gzipResponseWriter wraps an http.ResponseWriter to intercept written content
// and pass it through a gzip writer.
type gzipResponseWriter struct {
	Writer         io.Writer
	ResponseWriter http.ResponseWriter
}

// Header returns the underlying ResponseWriter's header map.
func (gzw *gzipResponseWriter) Header() http.Header {
	return gzw.ResponseWriter.Header()
}

// Write writes the compressed data to the gzip writer.
func (gzw *gzipResponseWriter) Write(b []byte) (int, error) {
	return gzw.Writer.Write(b)
}

// WriteHeader sets the HTTP status code on the underlying ResponseWriter.
func (gzw *gzipResponseWriter) WriteHeader(statusCode int) {
	gzw.ResponseWriter.WriteHeader(statusCode)
}

// RequestIDMiddleware returns a middleware that generates a unique request ID
// for each incoming request and sets it in the X-Request-ID response header.
// The ID is a cryptographically random hex string of 16 bytes.
func RequestIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := generateRequestID()
			w.Header().Set("X-Request-ID", requestID)
			next.ServeHTTP(w, r)
		})
	}
}

// generateRequestID generates a unique request ID using crypto/rand.
func generateRequestID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// CORSMiddleware returns a middleware that handles Cross-Origin Resource Sharing
// based on the provided CORSConfig. It sets appropriate headers and handles
// preflight OPTIONS requests.
func CORSMiddleware(config CORSConfig) func(http.Handler) http.Handler {
	originPatterns := config.AllowOrigins
	methods := strings.Join(config.AllowMethods, ", ")
	headers := strings.Join(config.AllowHeaders, ", ")
	exposed := strings.Join(config.ExposeHeaders, ", ")
	maxAge := fmt.Sprintf("%d", config.MaxAge)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowedOrigin := ""
			for _, pattern := range originPatterns {
				if pattern == "*" || pattern == origin {
					allowedOrigin = pattern
					break
				}
			}

			if allowedOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
				w.Header().Set("Access-Control-Allow-Methods", methods)
				w.Header().Set("Access-Control-Allow-Headers", headers)
				if exposed != "" {
					w.Header().Set("Access-Control-Expose-Headers", exposed)
				}
				if config.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
				w.Header().Set("Access-Control-Max-Age", maxAge)
			}

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// SecurityHeadersMiddleware returns a middleware that adds security-related HTTP
// headers to every response, including X-Content-Type-Options, X-Frame-Options,
// X-XSS-Protection, Strict-Transport-Security, Referrer-Policy, and
// Content-Security-Policy.
func SecurityHeadersMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'")
			w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
			next.ServeHTTP(w, r)
		})
	}
}

// RecoveryMiddleware returns a middleware that recovers from panics in downstream
// handlers, logs the panic details using the provided logger function, and
// returns a 500 Internal Server Error response.
func RecoveryMiddleware(logger func(string, ...interface{})) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					if logger != nil {
						logger("[RECOVERY] panic recovered: %v, path: %s", err, r.URL.Path)
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, `{"success":false,"error":{"code":500,"message":"Internal Server Error"}}`)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// LoggingMiddleware returns a middleware that logs each incoming request with
// the method, path, status code, and duration. The logger function receives
// a format string and arguments.
func LoggingMiddleware(logger func(string, ...interface{})) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			lrw := &loggingResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(lrw, r)

			duration := time.Since(start)
			if logger != nil {
				logger("[HTTP] %s %s %d %s", r.Method, r.URL.Path, lrw.statusCode, duration)
			}
		})
	}
}

// loggingResponseWriter wraps http.ResponseWriter to capture the status code.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code.
func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

// RateLimitMiddleware returns a middleware that limits the number of requests
// per second per client IP, with a configurable burst allowance. The burst
// parameter allows a short burst of requests above the sustained rate.
func RateLimitMiddleware(requestsPerSecond int, burst int) func(http.Handler) http.Handler {
	limiter := newAPIRateLimiter(burst, time.Second)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				clientIP = strings.Split(forwarded, ",")[0]
			}
			if !limiter.allow(clientIP) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				fmt.Fprintf(w, `{"success":false,"error":{"code":429,"message":"Too Many Requests"}}`)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// SessionMiddleware returns a middleware that manages sessions using the provided
// secret for cookie signing and the SessionStore for persistence. It reads the
// session ID from a cookie named "goscript_session" and attaches session data
// to the request context.
func SessionMiddleware(secret string, store SessionStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("goscript_session")
			var sessionID string

			if err == nil && cookie.Value != "" {
				sessionID = cookie.Value
			} else {
				// Generate a new session ID
				sessionID = generateSessionID(secret)
				http.SetCookie(w, &http.Cookie{
					Name:     "goscript_session",
					Value:    sessionID,
					Path:     "/",
					HttpOnly: true,
					Secure:   true,
					SameSite: http.SameSiteLaxMode,
					MaxAge:   86400,
				})
			}

			// Store session ID in the request header for downstream access
			r.Header.Set("X-Session-ID", sessionID)
			next.ServeHTTP(w, r)
		})
	}
}

// generateSessionID generates a new session ID by combining a random nonce with
// a timestamp. The secret parameter adds entropy to the generation.
func generateSessionID(secret string) string {
	b := make([]byte, 24)
	rand.Read(b)
	return hex.EncodeToString(b) + hex.EncodeToString([]byte(secret))[:8]
}
