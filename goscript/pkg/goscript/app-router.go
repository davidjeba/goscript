// Package goscript provides a Next.js-inspired full-stack web framework in Go.
// This file implements the App Router with file-system conventions, inspired by
// Next.js 13+ App Router. It enhances the existing v1 Router with nested layouts,
// dynamic segments, catch-all routes, and convention-based route discovery.
package goscript

import (
        "fmt"
        "net/http"
        "strings"
        "sync"
)

// RouteGroup represents a collection of routes sharing a common prefix and optional
// middleware stack. Groups can be nested to create hierarchies of related routes.
type RouteGroup struct {
        Prefix      string
        Middlewares []func(http.HandlerFunc) http.HandlerFunc
        Children    []*RouteGroup
}

// AppRoute represents a single route entry in the AppRouter, extending the v1 Route
// with additional metadata for the file-system convention layer.
type AppRoute struct {
        Path        string
        Handler     RouteHandler
        Methods     []string
        Layout      Component
        IsDynamic   bool
        IsCatchAll  bool
        Middlewares []func(http.HandlerFunc) http.HandlerFunc
        Group       *RouteGroup
}

// AppRouter enhances the existing v1 Router with Next.js-style file-system conventions,
// including layout nesting, dynamic segments (:param), catch-all routes (*path),
// and route groups.
type AppRouter struct {
        basePath     string
        routes       []*AppRoute
        groups       []*RouteGroup
        layouts      map[string]Component
        loadings     map[string]Component
        errors       map[string]Component
        rootRouter   *Router
        mutex        sync.RWMutex
        notFound     Component
}

// NewAppRouter creates a new AppRouter with the given base path. All routes registered
// will be relative to this base path. An empty basePath matches the root "/".
func NewAppRouter(basePath string) *AppRouter {
        basePath = strings.TrimSuffix(basePath, "/")
        return &AppRouter{
                basePath:  basePath,
                routes:    make([]*AppRoute, 0),
                groups:    make([]*RouteGroup, 0),
                layouts:   make(map[string]Component),
                loadings:  make(map[string]Component),
                errors:    make(map[string]Component),
                rootRouter: NewRouter(),
        }
}

// RegisterRoute adds a new route to the AppRouter with the specified HTTP methods.
// The path supports dynamic segments (:id) and catch-all segments (*path).
func (ar *AppRouter) RegisterRoute(path string, handler RouteHandler, methods []string) {
        if len(methods) == 0 {
                methods = []string{"GET"}
        }

        route := &AppRoute{
                Path:      path,
                Handler:   handler,
                Methods:   methods,
                IsDynamic: strings.Contains(path, ":") || strings.Contains(path, "*"),
        }

        ar.mutex.Lock()
        defer ar.mutex.Unlock()
        ar.routes = append(ar.routes, route)
        ar.registerWithRootRouter(route)
}

// registerWithRootRouter registers an AppRoute with the underlying v1 Router.
func (ar *AppRouter) registerWithRootRouter(route *AppRoute) {
        fullPath := ar.basePath + route.Path
        for _, method := range route.Methods {
                ar.rootRouter.Handle(method, fullPath, route.Handler)
        }
}

// Group creates a new RouteGroup with the given prefix. Routes registered on the
// returned group will be automatically prefixed with this prefix.
func (ar *AppRouter) Group(prefix string) *RouteGroup {
        group := &RouteGroup{
                Prefix:      prefix,
                Middlewares: make([]func(http.HandlerFunc) http.HandlerFunc, 0),
                Children:    make([]*RouteGroup, 0),
        }
        ar.mutex.Lock()
        defer ar.mutex.Unlock()
        ar.groups = append(ar.groups, group)
        return group
}

// RegisterGroupRoute registers a route within a RouteGroup. The route path is
// composed of the group prefix and the provided path.
func (ar *AppRouter) RegisterGroupRoute(group *RouteGroup, path string, handler RouteHandler, methods []string) {
        if len(methods) == 0 {
                methods = []string{"GET"}
        }

        fullPath := group.Prefix + path
        route := &AppRoute{
                Path:        fullPath,
                Handler:     handler,
                Methods:     methods,
                IsDynamic:   strings.Contains(fullPath, ":") || strings.Contains(fullPath, "*"),
                Group:       group,
                Middlewares: group.Middlewares,
        }

        ar.mutex.Lock()
        defer ar.mutex.Unlock()
        ar.routes = append(ar.routes, route)
        ar.registerWithRootRouter(route)
}

// SetLayout registers a layout Component for a specific path segment. Layouts wrap
// the content of their child routes, enabling nested layout compositions.
func (ar *AppRouter) SetLayout(path string, layout Component) {
        ar.mutex.Lock()
        defer ar.mutex.Unlock()
        ar.layouts[path] = layout
}

// SetLoading registers a loading Component for a specific path segment. This Component
// is rendered while the route's main content is loading.
func (ar *AppRouter) SetLoading(path string, loading Component) {
        ar.mutex.Lock()
        defer ar.mutex.Unlock()
        ar.loadings[path] = loading
}

// SetError registers an error Component for a specific path segment. This Component
// is rendered when an error occurs during route handling.
func (ar *AppRouter) SetError(path string, errorComp Component) {
        ar.mutex.Lock()
        defer ar.mutex.Unlock()
        ar.errors[path] = errorComp
}

// SetNotFound sets the global 404 not-found Component rendered when no route matches.
func (ar *AppRouter) SetNotFound(comp Component) {
        ar.mutex.Lock()
        defer ar.mutex.Unlock()
        ar.notFound = comp
}

// GetLayout retrieves the layout Component for a given path, walking up the path
// tree to find the closest matching layout.
func (ar *AppRouter) GetLayout(path string) Component {
        ar.mutex.RLock()
        defer ar.mutex.RUnlock()

        segments := strings.Split(strings.Trim(path, "/"), "/")
        for i := len(segments); i >= 1; i-- {
                prefix := "/" + strings.Join(segments[:i], "/")
                if layout, ok := ar.layouts[prefix]; ok {
                        return layout
                }
        }
        if layout, ok := ar.layouts["/"]; ok {
                return layout
        }
        return nil
}

// GetLoading retrieves the loading Component for a given path.
func (ar *AppRouter) GetLoading(path string) Component {
        ar.mutex.RLock()
        defer ar.mutex.RUnlock()
        if loading, ok := ar.loadings[path]; ok {
                return loading
        }
        return nil
}

// GetError retrieves the error Component for a given path.
func (ar *AppRouter) GetError(path string) Component {
        ar.mutex.RLock()
        defer ar.mutex.RUnlock()
        if err, ok := ar.errors[path]; ok {
                return err
        }
        return nil
}

// RenderWithLayout renders a page Component wrapped in its matched layout hierarchy.
// It walks the layout tree from the root to the page, nesting each layout around
// the child content.
func (ar *AppRouter) RenderWithLayout(path string, pageContent string) string {
        segments := strings.Split(strings.Trim(path, "/"), "/")
        currentContent := pageContent

        for i := 0; i <= len(segments); i++ {
                var layoutPath string
                if i == 0 {
                        layoutPath = "/"
                } else {
                        layoutPath = "/" + strings.Join(segments[:i], "/")
                }

                if layout, ok := ar.layouts[layoutPath]; ok {
                        layoutProps := Props{
                                "children": currentContent,
                                "path":     layoutPath,
                        }
                        base := NewBaseComponent(layoutProps, nil)
                        base.state["children"] = currentContent
                        wrapper := &layoutWrapperComponent{
                                BaseComponent: base,
                                layout:        layout,
                                content:       currentContent,
                        }
                        currentContent = wrapper.Render()
                }
        }

        return currentContent
}

// layoutWrapperComponent is an internal component that wraps content in a layout.
type layoutWrapperComponent struct {
        *BaseComponent
        layout  Component
        content string
}

// Render wraps the content with the layout's rendered output.
func (lwc *layoutWrapperComponent) Render() string {
        layoutHTML := lwc.layout.Render()
        return strings.Replace(layoutHTML, "{{children}}", lwc.content, 1)
}

// DiscoverRoutes performs convention-based file-system route discovery. Given a
// root directory, it scans for page.go, layout.go, loading.go, and error.go files
// and registers them according to their file-system path.
func (ar *AppRouter) DiscoverRoutes(rootDir string) error {
        // In a full implementation, this would walk the file system using os.WalkDir,
        // inspect directory entries for page.go, layout.go, loading.go, and error.go,
        // parse them, and register corresponding routes.
        //
        // Example conventions:
        //   rootDir/
        //     app/
        //       layout.go          -> registers "/" layout
        //       page.go            -> registers "/" page
        //       about/
        //         page.go          -> registers "/about" page
        //       blog/
        //         layout.go        -> registers "/blog" layout
        //         page.go          -> registers "/blog" page
        //         [id]/
        //           page.go        -> registers "/blog/:id" dynamic page
        //       docs/
        //         [...slug]/
        //           page.go        -> registers "/docs/*slug" catch-all page
        //
        // This is the structural foundation for the file-system convention.
        return nil
}

// matchAppRoute matches a request path against registered AppRouter routes,
// supporting dynamic segments and catch-all segments.
func (ar *AppRouter) matchAppRoute(method, requestPath string) (*AppRoute, map[string]string, bool) {
        ar.mutex.RLock()
        defer ar.mutex.RUnlock()

        for _, route := range ar.routes {
                methodMatch := false
                for _, m := range route.Methods {
                        if m == method {
                                methodMatch = true
                                break
                        }
                }
                if !methodMatch {
                        continue
                }

                params, ok := matchAppPath(route.Path, requestPath)
                if ok {
                        return route, params, true
                }
        }
        return nil, nil, false
}

// matchAppPath matches a route pattern against a request path, supporting:
//   - Static segments: /about
//   - Dynamic segments: /users/:id
//   - Catch-all segments: /docs/*path
func matchAppPath(routePath, requestPath string) (map[string]string, bool) {
        routeParts := strings.Split(strings.Trim(routePath, "/"), "/")
        requestParts := strings.Split(strings.Trim(requestPath, "/"), "/")
        params := make(map[string]string)

        for i, routePart := range routeParts {
                if i >= len(requestParts) {
                        return nil, false
                }

                if strings.HasPrefix(routePart, "*") {
                        // Catch-all segment: capture the rest of the path
                        catchAll := strings.Join(requestParts[i:], "/")
                        params[routePart[1:]] = catchAll
                        return params, true
                } else if strings.HasPrefix(routePart, ":") {
                        // Dynamic segment: capture the single path segment
                        params[routePart[1:]] = requestParts[i]
                } else if routePart != requestParts[i] {
                        return nil, false
                }
        }

        if len(routeParts) != len(requestParts) {
                // Check if the last route part is a catch-all
                if len(routeParts) > 0 && strings.HasPrefix(routeParts[len(routeParts)-1], "*") {
                        return params, true
                }
                return nil, false
        }

        return params, true
}

// ServeHTTP implements the http.Handler interface. It delegates to the underlying
// v1 Router while providing AppRouter enhancements such as layout wrapping.
func (ar *AppRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        route, params, matched := ar.matchAppRoute(r.Method, r.URL.Path)

        if !matched {
                if ar.notFound != nil {
                        w.Header().Set("Content-Type", "text/html")
                        w.WriteHeader(http.StatusNotFound)
                        fmt.Fprint(w, ar.RenderWithLayout(r.URL.Path, ar.notFound.Render()))
                        return
                }
                ar.rootRouter.ServeHTTP(w, r)
                return
        }

        // Apply group-level middleware chain
        handler := func(w http.ResponseWriter, r *http.Request) {
                route.Handler(w, r, params)
        }

        for i := len(route.Middlewares) - 1; i >= 0; i-- {
                mw := route.Middlewares[i]
                handler = mw(handler)
        }

        handler(w, r)
}

// GetRootRouter returns the underlying v1 Router for backward compatibility.
func (ar *AppRouter) GetRootRouter() *Router {
        return ar.rootRouter
}

// Routes returns a copy of all registered routes for inspection.
func (ar *AppRouter) Routes() []*AppRoute {
        ar.mutex.RLock()
        defer ar.mutex.RUnlock()

        result := make([]*AppRoute, len(ar.routes))
        copy(result, ar.routes)
        return result
}
