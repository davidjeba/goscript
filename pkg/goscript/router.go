package goscript

import (
	"net/http"
	"strings"
)

type RouteHandler func(w http.ResponseWriter, r *http.Request, params map[string]string)

type Route struct {
	Path    string
	Handler RouteHandler
}

type Router struct {
	routes     []Route
	middleware []func(http.HandlerFunc) http.HandlerFunc
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Use(middleware func(http.HandlerFunc) http.HandlerFunc) {
	r.middleware = append(r.middleware, middleware)
}

func (r *Router) AddRoute(path string, handler RouteHandler) {
	r.routes = append(r.routes, Route{Path: path, Handler: handler})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		params, ok := matchRoute(route.Path, req.URL.Path)
		if ok {
			handler := func(w http.ResponseWriter, r *http.Request) {
				route.Handler(w, r, params)
			}

			// Apply middleware
			for i := len(r.middleware) - 1; i >= 0; i-- {
				handler = r.middleware[i](handler)
			}

			handler(w, req)
			return
		}
	}

	http.NotFound(w, req)
}

func matchRoute(routePath, requestPath string) (map[string]string, bool) {
	routeParts := strings.Split(strings.Trim(routePath, "/"), "/")
	requestParts := strings.Split(strings.Trim(requestPath, "/"), "/")

	if len(routeParts) != len(requestParts) {
		return nil, false
	}

	params := make(map[string]string)

	for i, routePart := range routeParts {
		if strings.HasPrefix(routePart, ":") {
			params[routePart[1:]] = requestParts[i]
		} else if routePart != requestParts[i] {
			return nil, false
		}
	}

	return params, true
}

// Middleware for multi-tenant support
func TenantMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenant := r.Header.Get("X-Tenant")
		ctx := SetTenantContext(r.Context(), tenant)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// Middleware for multi-language support
func LanguageMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := r.Header.Get("Accept-Language")
		ctx := SetLanguageContext(r.Context(), lang)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

