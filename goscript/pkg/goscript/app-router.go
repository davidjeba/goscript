package goscript

import (
        "context"
        "net/http"
        "os"
        "path/filepath"
        "regexp"
        "strings"
        "sync"
)

// RouteSegment represents a segment of a route path
type RouteSegment struct {
        Name       string
        IsDynamic  bool
        IsCatchAll bool
        ParamName  string
}

// RouteNode represents a node in the route tree
type RouteNode struct {
        Segment    RouteSegment
        Children   map[string]*RouteNode
        Handler    RouteHandler
        Layout     Component
        Loading    Component
        Error      Component
        Middleware []MiddlewareFunc
        ParamName  string
}

// MiddlewareFunc defines the middleware signature
type MiddlewareFunc func(http.Handler) http.Handler

// AppRouter is a file-system based router that mirrors Next.js App Router
type AppRouter struct {
        root       *RouteNode
        routes     []*RouteConfig
        middleware []MiddlewareFunc
        layouts    map[string]Component
        notFound   Component
        mu         sync.RWMutex
        basePath   string
}

// RouteConfig holds the configuration for a registered route
type RouteConfig struct {
        Path       string
        Handler    RouteHandler
        Layout     Component
        Loading    Component
        Error      Component
        Middleware []MiddlewareFunc
        Methods    []string
}

// NewAppRouter creates a new file-system based AppRouter
func NewAppRouter(basePath string) *AppRouter {
        return &AppRouter{
                root: &RouteNode{
                        Children: make(map[string]*RouteNode),
                },
                routes:   make([]*RouteConfig, 0),
                layouts:  make(map[string]Component),
                basePath: basePath,
        }
}

// AutoDiscover scans the given directory for route files and registers them automatically.
// It mirrors Next.js App Router conventions:
//   - page.go → route handler
//   - layout.go → nested layout
//   - loading.go → loading UI (Suspense boundary)
//   - error.go → error boundary
//   - (group)/ → route group (no URL segment)
//   - [id]/ → dynamic segment
//   - [...slug]/ → catch-all segment
//   - page.api.go → API route handler
func (ar *AppRouter) AutoDiscover(dir string) error {
        return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
                if err != nil {
                        return err
                }

                relPath, _ := filepath.Rel(dir, path)
                if relPath == "." {
                        return nil
                }

                if info.IsDir() {
                        return nil
                }

                base := filepath.Base(path)
                ext := filepath.Ext(path)
                if ext != ".go" {
                        return nil
                }

                // Determine route type from filename
                routePath := filepath.Dir(relPath)
                routePath = filepath.ToSlash(routePath)
                if routePath == "." {
                        routePath = "/"
                }

                // Handle dynamic segments: [param] and [...slug]
                segmentPattern := regexp.MustCompile(`\[(?:\.\.\.)?([^\]]+)\]`)
                routePath = segmentPattern.ReplaceAllStringFunc(routePath, func(match string) string {
                        if strings.HasPrefix(match, "[...") {
                                return "*"
                        }
                        return ":$1"
                })

                // Remove route group names from path
                groupPattern := regexp.MustCompile(`\([^)]+\)/`)
                routePath = groupPattern.ReplaceAllString(routePath, "")
                routePath = strings.TrimSuffix(routePath, "/")

                switch {
                case strings.HasSuffix(base, "page.go") && !strings.Contains(base, ".api."):
                        ar.RegisterRoute(routePath, nil, []string{"GET"})
                case strings.Contains(base, "page.api.go"):
                        ar.RegisterRoute(routePath, nil, []string{"GET", "POST", "PUT", "DELETE"})
                case strings.HasSuffix(base, "layout.go"):
                        // Layout would be loaded here
                case strings.HasSuffix(base, "loading.go"):
                        // Loading component would be loaded here
                case strings.HasSuffix(base, "error.go"):
                        // Error boundary would be loaded here
                case strings.HasSuffix(base, "middleware.go"):
                        // Route-level middleware would be loaded here
                }

                return nil
        })
}

// RegisterRoute adds a route to the router with optional layout and error handling
func (ar *AppRouter) RegisterRoute(path string, handler RouteHandler, methods []string) {
        if len(methods) == 0 {
                methods = []string{"GET"}
        }

        config := &RouteConfig{
                Path:    path,
                Handler: handler,
                Methods: methods,
        }

        ar.mu.Lock()
        ar.routes = append(ar.routes, config)
        ar.mu.Unlock()

        ar.buildRouteTree(path, config)
}

// Use adds global middleware applied to all routes
func (ar *AppRouter) Use(mw MiddlewareFunc) {
        ar.middleware = append(ar.middleware, mw)
}

// SetNotFound sets the 404 handler component
func (ar *AppRouter) SetNotFound(component Component) {
        ar.notFound = component
}

// SetLayout sets a layout for a specific route prefix
func (ar *AppRouter) SetLayout(prefix string, component Component) {
        ar.layouts[prefix] = component
}

// buildRouteTree constructs the route trie from a path
func (ar *AppRouter) buildRouteTree(path string, config *RouteConfig) {
        node := ar.root
        segments := strings.Split(strings.Trim(path, "/"), "/")

        if len(segments) == 1 && segments[0] == "" {
                node.Handler = config.Handler
                return
        }

        for _, seg := range segments {
                if seg == "" {
                        continue
                }

                var key string
                isDynamic := strings.HasPrefix(seg, ":")
                isCatchAll := strings.HasPrefix(seg, "*")

                if isCatchAll {
                        key = "*"
                } else if isDynamic {
                        key = ":"
                } else {
                        key = seg
                }

                if _, exists := node.Children[key]; !exists {
                        node.Children[key] = &RouteNode{
                                Children: make(map[string]*RouteNode),
                        }
                }

                node = node.Children[key]
                if isDynamic {
                        node.ParamName = strings.TrimPrefix(seg, ":")
                }
                if isCatchAll {
                        node.ParamName = strings.TrimPrefix(seg, "*")
                }
        }

        node.Handler = config.Handler
}

// Match finds the matching route and returns params
func (ar *AppRouter) Match(method, path string) (*RouteConfig, map[string]string) {
        ar.mu.RLock()
        defer ar.mu.RUnlock()

        node := ar.root
        segments := strings.Split(strings.Trim(path, "/"), "/")
        params := make(map[string]string)

        if len(segments) == 1 && segments[0] == "" {
                if node.Handler != nil {
                        return &RouteConfig{Path: "/", Handler: node.Handler, Methods: []string{"GET"}}, params
                }
                return nil, nil
        }

        for _, seg := range segments {
                if child, ok := node.Children[seg]; ok {
                        node = child
                        continue
                }
                if child, ok := node.Children[":"]; ok {
                        params[child.ParamName] = seg
                        node = child
                        continue
                }
                if child, ok := node.Children["*"]; ok {
                        params[child.ParamName] = seg
                        node = child
                        continue
                }
                return nil, nil
        }

        if node.Handler != nil {
                return &RouteConfig{
                        Path:    path,
                        Handler: node.Handler,
                        Methods: []string{"GET"},
                }, params
        }

        return nil, nil
}

// ServeHTTP implements http.Handler
func (ar *AppRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        config, params := ar.Match(r.Method, r.URL.Path)

        if config == nil || config.Handler == nil {
                if ar.notFound != nil {
                        w.Header().Set("Content-Type", "text/html")
                        w.Write([]byte(ar.notFound.Render()))
                } else {
                        http.NotFound(w, r)
                }
                return
        }

        ctx := context.WithValue(r.Context(), "params", params)
        r = r.WithContext(ctx)

        var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                config.Handler(w, r, params)
        })

        for i := len(ar.middleware) - 1; i >= 0; i-- {
                handler = ar.middleware[i](handler)
        }

        handler.ServeHTTP(w, r)
}

// GenerateRouteMap returns all registered routes (useful for debugging/sitemap)
func (ar *AppRouter) GenerateRouteMap() []string {
        var routes []string
        ar.mu.RLock()
        defer ar.mu.RUnlock()

        var walk func(node *RouteNode, path string)
        walk = func(node *RouteNode, path string) {
                if node.Handler != nil {
                        routes = append(routes, path)
                }
                for key, child := range node.Children {
                        var nextPath string
                        switch key {
                        case ":":
                                nextPath = path + "/:" + child.ParamName
                        case "*":
                                nextPath = path + "/..." + child.ParamName
                        default:
                                nextPath = path + "/" + key
                        }
                        walk(child, nextPath)
                }
        }
        walk(ar.root, "/")

        return routes
}

// RoutesCount returns the number of registered routes
func (ar *AppRouter) RoutesCount() int {
        ar.mu.RLock()
        defer ar.mu.RUnlock()
        return len(ar.routes)
}
