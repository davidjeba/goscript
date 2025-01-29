package goscript

import (
	"net/http"
	"strings"
)

type RouteHandler func(w http.ResponseWriter, r *http.Request, params map[string]string)

type Route struct {
	Method  string
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

func (r *Router) Handle(method, path string, handler RouteHandler) {
	r.routes = append(r.routes, Route{Method: method, Path: path, Handler: handler})
}

func (r *Router) GET(path string, handler RouteHandler) {
	r.Handle("GET", path, handler)
}

func (r *Router) POST(path string, handler RouteHandler) {
	r.Handle("POST", path, handler)
}

func (r *Router) PUT(path string, handler RouteHandler) {
	r.Handle("PUT", path, handler)
}

func (r *Router) DELETE(path string, handler RouteHandler) {
	r.Handle("DELETE", path, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if route.Method == req.Method {
			params, ok := matchPath(route.Path, req.URL.Path)
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
	}

	http.NotFound(w, req)
}

func matchPath(routePath, requestPath string) (map[string]string, bool) {
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

