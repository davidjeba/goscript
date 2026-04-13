// Package main is the GoScript 2.0 showcase server demonstrating all framework
// features: App Router, API Routes, reactive attributes, metadata/SEO,
// middleware pipeline, server-side rendering, and the .gs runtime.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	goscript "github.com/davidjeba/goscript/pkg/goscript"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ---------------------------------------------------------------------------
	// 1. App Router — page routes
	// ---------------------------------------------------------------------------
	router := goscript.NewAppRouter("/")

	// Root layout wraps all pages
	rootLayout := goscript.NewServerComponent("RootLayout", func(props goscript.Props) string {
		childrenHTML, _ := props["children"].(string)
		return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  %s
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body { font-family: 'Inter', system-ui, -apple-system, sans-serif; background: #fafafa; color: #1a1a2e; }
    header { background: #fff; border-bottom: 1px solid #e5e7eb; padding: 1rem 2rem; display: flex; align-items: center; gap: 1rem; }
    header .logo { color: #10b981; font-weight: 800; font-size: 1.25rem; }
    header .version { font-size: 0.75rem; background: #ecfdf5; color: #059669; padding: 2px 8px; border-radius: 99px; font-weight: 600; }
    nav { background: #fff; border-bottom: 1px solid #f3f4f6; padding: 0.5rem 2rem; display: flex; gap: 0; }
    nav a { color: #6b7280; text-decoration: none; font-size: 0.875rem; padding: 0.5rem 1rem; border-radius: 6px; transition: all 0.15s; }
    nav a:hover { color: #10b981; background: #ecfdf5; }
    main { max-width: 960px; margin: 2rem auto; padding: 0 1rem; }
    .card { background: #fff; border-radius: 12px; padding: 1.5rem; margin-bottom: 1rem; box-shadow: 0 1px 3px rgba(0,0,0,0.06); border: 1px solid #f3f4f6; }
    .card h2 { font-size: 1.125rem; margin-bottom: 0.5rem; color: #111827; }
    .card p { color: #6b7280; line-height: 1.6; font-size: 0.9rem; }
    .badge { display: inline-block; font-size: 0.7rem; padding: 2px 8px; border-radius: 99px; margin-left: 8px; vertical-align: middle; font-weight: 600; }
    .badge-new { background: #ecfdf5; color: #059669; }
    .badge-v1 { background: #f3f4f6; color: #6b7280; }
    pre { background: #1e293b; color: #e2e8f0; padding: 1rem; border-radius: 8px; overflow-x: auto; font-size: 0.85rem; line-height: 1.5; margin-top: 0.75rem; }
    pre code { font-family: 'JetBrains Mono', 'Fira Code', monospace; }
    .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 1rem; }
    footer { text-align: center; padding: 2rem; color: #9ca3af; font-size: 0.8rem; margin-top: 2rem; border-top: 1px solid #f3f4f6; }
    .counter-box { text-align: center; padding: 2rem; }
    .counter-value { font-size: 4rem; font-weight: 800; color: #10b981; margin: 1rem 0; }
    .btn { padding: 0.5rem 1.25rem; margin: 0.25rem; border: none; border-radius: 8px; font-size: 0.9rem; cursor: pointer; transition: all 0.15s; font-weight: 600; }
    .btn-primary { background: #10b981; color: #fff; }
    .btn-primary:hover { background: #059669; }
    .btn-danger { background: #ef4444; color: #fff; }
    .btn-danger:hover { background: #dc2626; }
    .btn-ghost { background: #f3f4f6; color: #374151; }
    .btn-ghost:hover { background: #e5e7eb; }
    .note { margin-top: 0.75rem; color: #9ca3af; font-size: 0.8rem; }
  </style>
  <script src="/__goscript/runtime.js"></script>
</head>
<body>
  <header>
    <span class="logo">GoScript</span>
    <span class="version">v2.0.0</span>
  </header>
  <nav>
    <a href="/">Home</a>
    <a href="/about">About</a>
    <a href="/api-demo">API Demo</a>
    <a href="/api/status" target="_blank">API Status</a>
  </nav>
  <main>%s</main>
  <footer>GoScript 2.0 — The Go-native web framework &middot; MIT License</footer>
</body>
</html>`, goscript.NewMetadata().
			SetTitle("GoScript 2.0 — Go-native Web Framework").
			SetDescription("A full-stack web framework written in Go with its own language (.gs), reactive runtime, and App Router").
			SetThemeColor("#10b981").
			Build().Render(), childrenHTML)
	}, nil)
	router.SetLayout("/", rootLayout)

	// ---- Home page ----
	router.RegisterRoute("/", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		home := goscript.NewServerComponent("HomePage", func(props goscript.Props) string {
			return fmt.Sprintf(`<div class="grid">
  <div class="card">
    <h2>Reactive Counter <span class="badge badge-new">goscript</span></h2>
    <p>Server-rendered HTML with zero JavaScript — goscript's reactive attributes handle interactivity.</p>
    <div class="counter-box">
      <div id="count-display" class="counter-value">0</div>
      <div>
        <button class="btn btn-danger"
          gs-trigger="click"
          gs-get="/api/counter?n=-1"
          gs-target="#count-display"
          gs-swap="innerHTML">−</button>
        <button class="btn btn-primary"
          gs-trigger="click"
          gs-get="/api/counter?n=1"
          gs-target="#count-display"
          gs-swap="innerHTML">+</button>
        <button class="btn btn-ghost"
          gs-trigger="click"
          gs-get="/api/counter?reset=1"
          gs-target="#count-display"
          gs-swap="innerHTML">Reset</button>
      </div>
      <p class="note">gs-trigger, gs-get, gs-target, gs-swap — zero JavaScript</p>
    </div>
  </div>

  <div class="card">
    <h2>Server-Side Rendered</h2>
    <p>This entire page is rendered by Go. No build step. No node_modules.</p>
    <pre><code>func homeHandler(w http.ResponseWriter, r *http.Request) {
    metadata := goscript.NewMetadata()
    html := renderPage(metadata, content)
    w.Write([]byte(html))
}</code></pre>
  </div>

  <div class="card">
    <h2>Why Go? <span class="badge badge-new">performance</span></h2>
    <p>Single binary deployment. Millisecond startup. 20MB memory. Millions of goroutines. Compile-time type safety.</p>
  </div>

  <div class="card">
    <h2>Backward Compatible <span class="badge badge-v1">v1</span></h2>
    <p>All v1 components, stores, routers, and sub-modules work unchanged.</p>
  </div>
</div>`)
		}, nil)

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, router.RenderWithLayout("/", home.Render()))
	}, []string{"GET"})

	// ---- About page ----
	router.RegisterRoute("/about", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		about := goscript.NewServerComponent("AboutPage", func(props goscript.Props) string {
			return fmt.Sprintf(`<div class="card">
  <h2>About GoScript 2.0 <span class="badge badge-new">v2</span></h2>
  <p>GoScript is a full-stack web framework written entirely in Go. It features its own language (.gs files) that compile to browser JavaScript, reactive HTML attributes for zero-JS interactivity, an App Router with layouts and dynamic segments, streaming SSR, SSG/ISR, and a comprehensive middleware pipeline.</p>
</div>

<div class="card">
  <h2>The .gs Language</h2>
  <p>.gs files are GoScript's own language — Go-like syntax that compiles to browser JavaScript. Write components, state, and event handlers in .gs, compile them with <code>gopm compile</code>, and ship zero-config JavaScript.</p>
  <pre><code>func Counter() dom.Element {
    count, setCount := state.Use(0)
    return dom.CreateElement("div", dom.Props{"class": "counter"},
        dom.CreateElement("h1", nil, count),
        dom.CreateElement("button", dom.Props{
            "onclick": func(e dom.Event) { setCount(count + 1) },
        }, "+"),
    )
}</code></pre>
</div>

<div class="card">
  <h2>Reactive Attributes</h2>
  <p>Add <code>gs-trigger</code>, <code>gs-get</code>, <code>gs-target</code>, and <code>gs-swap</code> to any HTML element for zero-JavaScript reactivity. The goscript runtime handles the rest.</p>
  <pre><code>&lt;button gs-trigger="click"
        gs-get="/api/counter?n=1"
        gs-target="#count"
        gs-swap="innerHTML"&gt;
  Increment
&lt;/button&gt;</code></pre>
</div>

<div class="card">
  <h2>Architecture</h2>
  <p>GoScript compiles to a single static binary. No Node.js runtime. No npm. No webpack. Just Go.</p>
  <pre><code>gopm build    # Compile .gs → .js, output to ./build
go build     # Single static binary
./server     # Deploy anywhere</code></pre>
</div>`)
		}, nil)

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, router.RenderWithLayout("/about", about.Render()))
	}, []string{"GET"})

	// ---- API Demo page ----
	router.RegisterRoute("/api-demo", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		apiDemo := goscript.NewServerComponent("APIDemoPage", func(props goscript.Props) string {
			return fmt.Sprintf(`<div class="card">
  <h2>API Routes</h2>
  <p>Convention-based API router with rate limiting, auth, and CORS — all built in.</p>
  <pre><code>api := goscript.NewAPIRouter()
api.GET("/api/health", healthHandler,
    goscript.WithRateLimit(100, time.Minute))
api.POST("/api/users", createUser,
    goscript.WithAuth([]string{"Bearer"}, validate))</code></pre>
</div>

<div class="card">
  <h2>Try the API</h2>
  <p>These endpoints are live on this server:</p>
  <div style="margin-top: 1rem; display: flex; flex-direction: column; gap: 0.5rem;">
    <a href="/api/status" target="_blank" style="padding: 0.5rem 1rem; background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 6px; text-decoration: none; color: #166534; font-size: 0.875rem; font-family: monospace;">GET /api/status → framework info</a>
    <a href="/api/users" target="_blank" style="padding: 0.5rem 1rem; background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 6px; text-decoration: none; color: #166534; font-size: 0.875rem; font-family: monospace;">GET /api/users → user list</a>
    <a href="/api/counter?n=5" target="_blank" style="padding: 0.5rem 1rem; background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 6px; text-decoration: none; color: #166534; font-size: 0.875rem; font-family: monospace;">GET /api/counter?n=5 → HTML fragment</a>
  </div>
</div>

<div class="card">
  <h2>Middleware Pipeline</h2>
  <p>Composable middleware with built-in implementations:</p>
  <pre><code>pipeline := goscript.NewPipeline().
    Use(goscript.RequestIDMiddleware()).
    Use(goscript.GzipMiddleware()).
    Use(goscript.CORSMiddleware(goscript.DefaultCORSConfig())).
    Use(goscript.SecurityHeadersMiddleware()).
    Use(goscript.LoggingMiddleware(log.Printf)).
    Use(goscript.RecoveryMiddleware(log.Printf))</code></pre>
</div>`)
		}, nil)

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, router.RenderWithLayout("/api-demo", apiDemo.Render()))
	}, []string{"GET"})

	// Set error and loading boundaries
	router.SetError("/", goscript.NewErrorComponent("Something went wrong."))
	router.SetLoading("/", goscript.NewSkeletonLoader(4))

	// ---------------------------------------------------------------------------
	// 2. API Router
	// ---------------------------------------------------------------------------
	apiRouter := goscript.NewAPIRouter()

	apiRouter.GET("/api/status", func(ctx *goscript.APIContext) (interface{}, error) {
		return map[string]interface{}{
			"status":    "healthy",
			"version":   goscript.Version,
			"framework": goscript.FrameworkName,
			"features":  goscript.V2Features,
			"runtime":   "goscript",
		}, nil
	})

	apiRouter.GET("/api/users", func(ctx *goscript.APIContext) (interface{}, error) {
		return map[string]interface{}{
			"users": []map[string]interface{}{
				{"id": 1, "name": "Alice", "role": "admin", "email": "alice@go.dev"},
				{"id": 2, "name": "Bob", "role": "user", "email": "bob@go.dev"},
				{"id": 3, "name": "Charlie", "role": "user", "email": "charlie@go.dev"},
			},
			"total": 3,
		}, nil
	})

	apiRouter.GET("/api/users/:id", func(ctx *goscript.APIContext) (interface{}, error) {
		userID := ctx.Params["id"]
		return map[string]interface{}{
			"id":    userID,
			"name":  "GoScript Developer",
			"email": "dev@choscript.dev",
		}, nil
	}, goscript.WithRateLimit(100, time.Minute))

	apiRouter.POST("/api/users", func(ctx *goscript.APIContext) (interface{}, error) {
		return map[string]interface{}{
			"message": "User created successfully",
			"data":    ctx.Body,
		}, nil
	}, goscript.WithRateLimit(10, time.Minute), goscript.WithCORS("*"))

	apiRouter.SetNotFound(func(ctx *goscript.APIContext) (interface{}, error) {
		return map[string]interface{}{
			"error":   "not_found",
			"message": "The requested API endpoint does not exist",
		}, nil
	})

	// ---------------------------------------------------------------------------
	// 3. Counter endpoint — returns HTML fragment for goscript reactive swap
	// ---------------------------------------------------------------------------
	counter := 0

	// ---------------------------------------------------------------------------
	// 4. Middleware Pipeline
	// ---------------------------------------------------------------------------
	pipeline := goscript.NewPipeline().
		Use(goscript.RequestIDMiddleware()).
		Use(goscript.GzipMiddleware()).
		Use(goscript.CORSMiddleware(goscript.DefaultCORSConfig())).
		Use(goscript.SecurityHeadersMiddleware()).
		Use(goscript.LoggingMiddleware(func(format string, args ...interface{}) {
			log.Printf(format, args...)
		})).
		Use(goscript.RecoveryMiddleware(func(format string, args ...interface{}) {
			log.Printf(format, args...)
		}))

	// ---------------------------------------------------------------------------
	// 5. Streaming SSR Engine
	// ---------------------------------------------------------------------------
	store := goscript.NewStore()
	store.SetState("siteName", "GoScript 2.0 Server")
	streamEngine := goscript.NewStreamSSREngine(store)

	// ---------------------------------------------------------------------------
	// 6. SSG/ISR Engine
	// ---------------------------------------------------------------------------
	ssgEngine := goscript.NewSSGEngine("./out")
	ssgEngine.AddPage(goscript.PageConfig{
		Path:       "/",
		Component:  goscript.NewServerComponent("StaticHome", func(props goscript.Props) string {
			return `<div class="card"><h2>Static Home (SSG)</h2><p>Pre-generated at build time.</p></div>`
		}, nil),
		RenderMode: goscript.RenderSSG,
	})
	ssgEngine.AddPage(goscript.PageConfig{
		Path:       "/blog",
		Component:  goscript.NewServerComponent("BlogISR", func(props goscript.Props) string {
			return `<div class="card"><h2>Blog (ISR)</h2><p>Revalidated every 60s.</p></div>`
		}, nil),
		RenderMode: goscript.RenderISR,
		Revalidate: 60 * time.Second,
	})

	if err := ssgEngine.Build(context.Background()); err != nil {
		log.Printf("[SSG] Build warning: %v", err)
	}

	// ---------------------------------------------------------------------------
	// 7. Route groups
	// ---------------------------------------------------------------------------
	apiV2Group := router.Group("/api/v2")
	router.RegisterGroupRoute(apiV2Group, "/health", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"healthy","version":"%s"}`, goscript.Version)
	}, []string{"GET"})
	router.RegisterGroupRoute(apiV2Group, "/features", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"features":%+v}`, goscript.FeatureFlags)
	}, []string{"GET"})

	// ---------------------------------------------------------------------------
	// 8. Combine into final handler
	// ---------------------------------------------------------------------------
	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		// Counter API — returns HTML fragment for goscript reactive swap
		case r.URL.Path == "/api/counter":
			n := r.URL.Query().Get("n")
			_, reset := r.URL.Query()["reset"]
			if reset {
				counter = 0
			} else if n != "" {
				delta := 0
				fmt.Sscanf(n, "%d", &delta)
				counter += delta
			}
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, "%d", counter)
			return

		// Streaming SSR demo
		case r.URL.Path == "/stream":
			component := goscript.NewServerComponent("StreamDemo", func(props goscript.Props) string {
				return `<div class="card"><h2>Streaming SSR</h2><p>Loading in chunks...</p></div>`
			}, nil)
			boundaries := []goscript.SuspenseBoundary{
				{
					ID:       "slow-data",
					Fallback: goscript.NewSkeletonLoader(3),
					Loader: func(ctx interface{}) (goscript.Component, error) {
						time.Sleep(2 * time.Second)
						return goscript.NewServerComponent("SlowData", func(props goscript.Props) string {
							return `<div class="card" style="border-left:3px solid #10b981;"><h2>Streamed In!</h2><p>This content arrived after 2s while the skeleton was shown.</p></div>`
						}, nil), nil
					},
				},
			}
			streamEngine.RenderStream(w, r, component, boundaries)
			return

		// API router
		case len(r.URL.Path) >= 4 && r.URL.Path[:4] == "/api":
			apiRouter.ServeHTTP(w, r)
			return

		// Static SSG pages
		case r.URL.Path == "/static/" || r.URL.Path == "/static":
			ssgEngine.ServeSSG(w, r)
			return

		// Default: AppRouter
		default:
			router.ServeHTTP(w, r)
		}
	})

	// Static file serving
	assetManager := goscript.NewAssetManager("./static")
	staticHandler := assetManager.ServeAssets("/static")

	finalMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) >= 8 && r.URL.Path[:8] == "/static/" {
			staticHandler.ServeHTTP(w, r)
			return
		}
		pipeline.Execute(w, r, mainHandler)
	})

	// ---------------------------------------------------------------------------
	// 9. Start
	// ---------------------------------------------------------------------------
	fmt.Printf(`
╔═══════════════════════════════════════════════╗
║          GoScript 2.0 Server                 ║
║          The Go-native web framework         ║
╠═══════════════════════════════════════════════╣
║  Server:    http://localhost:%-15s ║
║  API:       http://localhost:%-15s ║
║  Stream:    http://localhost:%-15s ║
║  API V2:    http://localhost:%-15s ║
║  Framework: %-33s ║
║  Version:   %-33s ║
╠═══════════════════════════════════════════════╣
║  Routes:                                     ║
║    GET  /              (App Router + Layout)  ║
║    GET  /about         (App Router + Layout)  ║
║    GET  /api-demo      (App Router + Layout)  ║
║    GET  /api/status    (API)                  ║
║    GET  /api/users     (API)                  ║
║    GET  /api/users/:id (API + Rate Limit)     ║
║    POST /api/users     (API + Auth + CORS)    ║
║    GET  /api/counter   (HTML Fragment)        ║
║    GET  /api/v2/health (Route Group)          ║
║    GET  /api/v2/features(Route Group)         ║
║    GET  /stream        (Streaming SSR)        ║
║    GET  /static/...    (SSG Assets)           ║
║  Middleware:                                 ║
║    RequestID → Gzip → CORS → Security →     ║
║    Logging → Recovery                        ║
╚═══════════════════════════════════════════════╝
`, port, "api/status", "stream", "api/v2/health", goscript.FrameworkName, goscript.Version)

	if err := http.ListenAndServe(":"+port, finalMux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
