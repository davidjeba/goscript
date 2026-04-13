package goscript

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/davidjeba/goscript/pkg/gslib"
)

// runtimeJSPath is the URL path at which the runtime is served.
const runtimeJSPath = "/__goscript/runtime.js"

// RuntimeJS returns the embedded runtime.js content as a string.
// The actual embedding happens in pkg/gslib/runtime_embed.go via go:embed.
// This function provides convenient access from the goscript package.
//
// Usage:
//
//	js := RuntimeJS()
//	fmt.Println(len(js)) // ~40KB uncompressed
func RuntimeJS() string {
	return gslib.RuntimeJS
}

// RuntimeHandler returns an http.HandlerFunc that serves the embedded
// runtime.js with appropriate caching headers. The handler sets a long
// max-age cache header (1 year) keyed on the response's ETag so that
// browsers cache the runtime aggressively while still allowing cache
// busting when the binary changes.
//
// Usage:
//
//	mux := http.NewServeMux()
//	mux.Handle("/__goscript/runtime.js", RuntimeHandler())
func RuntimeHandler() http.HandlerFunc {
	data := RuntimeJS()
	etag := fmt.Sprintf(`"%x"`, len(data))

	return func(w http.ResponseWriter, r *http.Request) {
		// Only serve GET/HEAD requests
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Normalize path: strip trailing slash
		reqPath := strings.TrimRight(r.URL.Path, "/")
		if reqPath != runtimeJSPath {
			http.NotFound(w, r)
			return
		}

		// ETag-based conditional request support
		if match := r.Header.Get("If-None-Match"); match != "" {
			if match == etag {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		// Set caching headers — the runtime changes only when the binary
		// is rebuilt, so we use a very long max-age with ETag validation.
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		w.Header().Set("ETag", etag)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))

		if r.Method == http.MethodHead {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, data)
	}
}

// RuntimeRouteHandler returns a RouteHandler (compatible with the goscript
// Router) that serves the embedded runtime.js. Use this when registering
// the route via the goscript Router directly.
//
// Usage:
//
//	router.GET("/__goscript/runtime.js", RuntimeRouteHandler())
func RuntimeRouteHandler() RouteHandler {
	return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		RuntimeHandler()(w, r)
	}
}

// RegisterRuntimeRoute adds the /__goscript/runtime.js route to a goscript
// Router instance. This is the simplest way to wire up the embedded runtime.
//
// Usage:
//
//	router := goscript.NewRouter()
//	goscript.RegisterRuntimeRoute(router)
func RegisterRuntimeRoute(router *Router) {
	router.GET(runtimeJSPath, RuntimeRouteHandler())
}

// RegisterRuntimeRouteStd adds the /__goscript/runtime.js route to a
// standard library http.ServeMux. This supports projects that use
// http.NewServeMux or http.DefaultServeMux instead of the goscript Router.
//
// Usage:
//
//	mux := http.NewServeMux()
//	goscript.RegisterRuntimeRouteStd(mux)
func RegisterRuntimeRouteStd(mux *http.ServeMux) {
	mux.HandleFunc(runtimeJSPath, RuntimeHandler())
}
