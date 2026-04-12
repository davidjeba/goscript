// Package main demonstrates the GoScript v2 framework with a full-featured
// development server that showcases all v2 improvements while maintaining
// backward compatibility with v1 components and APIs.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/davidjeba/goscript/pkg/goscript"
)

func main() {
	// Determine the server port from the PORT environment variable or -port flag.
	port := 8080
	if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}
	for i, arg := range os.Args[1:] {
		if arg == "-port" && i+1 < len(os.Args[1:]) {
			if p, err := strconv.Atoi(os.Args[i+2]); err == nil {
				port = p
			}
		}
	}

	// ---------------------------------------------------------------------------
	// 1. Create the v2 App Router with file-system conventions
	// ---------------------------------------------------------------------------
	appRouter := goscript.NewAppRouter("")

	// Define a root layout that wraps all pages
	rootLayout := goscript.NewServerComponent("RootLayout", func(props goscript.Props) string {
		childrenHTML, _ := props["children"].(string)
		return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>GoScript v2 Demo</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body { font-family: system-ui, -apple-system, sans-serif; background: #fafafa; color: #1a1a2e; }
    header { background: #16213e; color: #e94560; padding: 1rem 2rem; display: flex; align-items: center; gap: 1rem; }
    header h1 { font-size: 1.25rem; font-weight: 700; }
    header span { font-size: 0.85rem; color: #a8a8b3; }
    nav { background: #0f3460; padding: 0.5rem 2rem; display: flex; gap: 1.5rem; }
    nav a { color: #a8a8b3; text-decoration: none; font-size: 0.9rem; transition: color 0.2s; }
    nav a:hover { color: #e94560; }
    main { max-width: 960px; margin: 2rem auto; padding: 0 1rem; }
    .card { background: #fff; border-radius: 8px; padding: 1.5rem; margin-bottom: 1rem; box-shadow: 0 1px 3px rgba(0,0,0,0.08); }
    .card h2 { font-size: 1.1rem; margin-bottom: 0.5rem; color: #16213e; }
    .card p { color: #555; line-height: 1.6; font-size: 0.9rem; }
    .badge { display: inline-block; background: #e94560; color: #fff; font-size: 0.7rem; padding: 2px 8px; border-radius: 99px; margin-left: 8px; vertical-align: middle; }
    footer { text-align: center; padding: 2rem; color: #a8a8b3; font-size: 0.8rem; margin-top: 2rem; border-top: 1px solid #eee; }
  </style>
  %s
</head>
<body>
  <header>
    <h1>GoScript <span>v2.0.0</span></h1>
  </header>
  <nav>
    <a href="/">Home</a>
    <a href="/about">About</a>
    <a href="/api/status">API Status</a>
    <a href="/users/42">User Demo</a>
  </nav>
  <main>%s</main>
  <footer>GoScript v2 — Full-Stack Go Web Framework &middot; MIT License</footer>
</body>
</html>`, goscript.DefaultMetadata().SetTitle("GoScript v2 Demo").Build().Render(), childrenHTML)
	}, nil)
	appRouter.SetLayout("/", rootLayout)

	// ---------------------------------------------------------------------------
	// 2. Register pages using the v2 App Router
	// ---------------------------------------------------------------------------
	appRouter.RegisterRoute("/", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		homePage := goscript.NewServerComponent("HomePage", func(props goscript.Props) string {
			return fmt.Sprintf(`<div class="card">
<h2>Welcome to GoScript v2</h2>
<p>This demo showcases all 10 v2 improvements over the original v1 framework. Navigate the links above to explore each feature.</p>
</div>
<div class="card">
<h2>v2 Features</h2>
<p>App Router &middot; Streaming SSR &middot; Server/Client Components &middot; API Routes &middot; Middleware Pipeline &middot; SSG/ISR &middot; Error/Loading Boundaries &middot; Metadata/SEO &middot; HMR Dev Server &middot; Route Groups</p>
</div>
<div class="card">
<h2>Backward Compatible <span class="badge">v1</span></h2>
<p>All original v1 components, stores, routers, and sub-modules work unchanged. The v1 Router is accessible via appRouter.GetRootRouter().</p>
</div>`)
		}, nil)

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, appRouter.RenderWithLayout("/", homePage.Render()))
	}, []string{"GET"})

	appRouter.RegisterRoute("/about", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		aboutPage := goscript.NewServerComponent("AboutPage", func(props goscript.Props) string {
			return fmt.Sprintf(`<div class="card">
<h2>About GoScript <span class="badge">v2</span></h2>
<p>GoScript is a Next.js-inspired full-stack web framework written entirely in Go. Version 2 adds 10 major improvements including App Router conventions, streaming SSR with Suspense, server and client components, convention-based API routes, composable middleware, SSG/ISR, error and loading boundaries, a metadata API for SEO, HMR dev server, and route groups with dynamic segments.</p>
</div>
<div class="card">
<h2>Why Go?</h2>
<p>Go delivers fast compile times, zero-cost abstractions, single-binary deployments, excellent concurrency primitives, and memory safety without garbage collector pauses. GoScript brings the developer experience of modern JavaScript frameworks to the Go ecosystem.</p>
</div>`)
		}, nil)

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, appRouter.RenderWithLayout("/about", aboutPage.Render()))
	}, []string{"GET"})

	// Dynamic segment demo
	appRouter.RegisterRoute("/users/:id", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		userID := params["id"]
		userPage := goscript.NewServerComponent("UserPage", func(props goscript.Props) string {
			return fmt.Sprintf(`<div class="card">
<h2>User Profile</h2>
<p>User ID: <strong>%s</strong></p>
<p>This page demonstrates dynamic route segments. The <code>:id</code> parameter in the path is captured and rendered here via the App Router.</p>
</div>`, userID)
		}, nil)

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, appRouter.RenderWithLayout("/users/"+userID, userPage.Render()))
	}, []string{"GET"})

	// Register error and loading boundaries for demonstration
	appRouter.SetError("/", goscript.NewErrorComponent("An error occurred while loading this page."))
	appRouter.SetLoading("/", goscript.NewSkeletonLoader(4))

	// ---------------------------------------------------------------------------
	// 3. Create the v2 API Router with rate limiting and CORS
	// ---------------------------------------------------------------------------
	apiRouter := goscript.NewAPIRouter()

	apiRouter.GET("/api/status", func(ctx *goscript.APIContext) (interface{}, error) {
		return map[string]interface{}{
			"framework":    goscript.FrameworkName,
			"version":      goscript.Version,
			"features":     goscript.V2Features,
			"sub_modules":  goscript.V1SubModules,
			"uptime":       "running",
		}, nil
	})

	apiRouter.GET("/api/users/:id", func(ctx *goscript.APIContext) (interface{}, error) {
		userID := ctx.Params["id"]
		return map[string]interface{}{
			"id":    userID,
			"name":  "GoScript Developer",
			"email": "dev@choscript.dev",
		}, nil
	}, goscript.WithRateLimit(100, time.Minute), goscript.WithCORS("*"))

	apiRouter.POST("/api/users", func(ctx *goscript.APIContext) (interface{}, error) {
		return map[string]interface{}{
			"message": "User created successfully",
			"data":    ctx.Body,
		}, nil
	}, goscript.WithRateLimit(10, time.Minute), goscript.WithCORS("*"), goscript.WithAuth([]string{"Bearer"}, func(token string) bool {
		return len(token) > len("Bearer ")
	}))

	apiRouter.SetNotFound(func(ctx *goscript.APIContext) (interface{}, error) {
		return map[string]interface{}{
			"error":   "not_found",
			"message": "The requested API endpoint does not exist",
		}, nil
	})

	// ---------------------------------------------------------------------------
	// 4. Build the middleware pipeline (v2 composable middleware)
	// ---------------------------------------------------------------------------
	pipeline := goscript.NewPipeline()

	// Apply middleware in order: each wraps the next
	pipeline.Use(goscript.RecoveryMiddleware(func(format string, args ...interface{}) {
		log.Printf(format, args...)
	}))
	pipeline.Use(goscript.LoggingMiddleware(func(format string, args ...interface{}) {
		log.Printf(format, args...)
	}))
	pipeline.Use(goscript.RequestIDMiddleware())
	pipeline.Use(goscript.SecurityHeadersMiddleware())
	pipeline.Use(goscript.CORSMiddleware(goscript.DefaultCORSConfig()))
	pipeline.Use(goscript.GzipMiddleware())
	pipeline.Use(goscript.RateLimitMiddleware(1000, 50))

	// ---------------------------------------------------------------------------
	// 5. Create a streaming SSR engine for the streaming demo route
	// ---------------------------------------------------------------------------
	store := goscript.NewStore()
	store.SetState("siteName", "GoScript v2 Demo")
	streamEngine := goscript.NewStreamSSREngine(store)

	// ---------------------------------------------------------------------------
	// 6. Create SSG/ISR engine for static page generation
	// ---------------------------------------------------------------------------
	ssgEngine := goscript.NewSSGEngine("./out")
	ssgEngine.AddPage(goscript.PageConfig{
		Path:       "/",
		Component:  goscript.NewServerComponent("StaticHome", func(props goscript.Props) string {
			return `<div class="card"><h2>Static Home Page (SSG)</h2><p>This page was pre-generated at build time.</p></div>`
		}, nil),
		RenderMode: goscript.RenderSSG,
	})
	ssgEngine.AddPage(goscript.PageConfig{
		Path:       "/blog",
		Component:  goscript.NewServerComponent("BlogIndex", func(props goscript.Props) string {
			return `<div class="card"><h2>Blog (ISR)</h2><p>This page is revalidated every 60 seconds.</p></div>`
		}, nil),
		RenderMode: goscript.RenderISR,
		Revalidate: 60 * time.Second,
	})

	// Pre-build static pages
	if err := ssgEngine.Build(context.Background()); err != nil {
		log.Printf("[SSG] Build warning: %v", err)
	}

	// ---------------------------------------------------------------------------
	// 7. Create a route group for the /api/v2/ prefixed endpoints
	// ---------------------------------------------------------------------------
	apiV2Group := appRouter.Group("/api/v2")
	appRouter.RegisterGroupRoute(apiV2Group, "/health", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"healthy","version":"%s"}`, goscript.Version)
	}, []string{"GET"})

	appRouter.RegisterGroupRoute(apiV2Group, "/features", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		w.Header().Set("Content-Type", "application/json")
		features := make([]map[string]bool, 0, len(goscript.FeatureFlags))
		for k, v := range goscript.FeatureFlags {
			features = append(features, map[string]bool{k: v})
		}
		fmt.Fprintf(w, `{"features":%+v}`, features)
	}, []string{"GET"})

	// ---------------------------------------------------------------------------
	// 8. Combine everything into the final HTTP handler
	// ---------------------------------------------------------------------------
	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		// Streaming SSR demo endpoint
		case r.URL.Path == "/stream":
			component := goscript.NewServerComponent("StreamDemo", func(props goscript.Props) string {
				return `<div class="card"><h2>Streaming SSR Demo</h2><p>Watch this page load in chunks...</p></div>`
			}, nil)

			boundaries := []goscript.SuspenseBoundary{
				{
					ID:       "slow-data",
					Fallback: goscript.NewSkeletonLoader(3),
					Loader: func(ctx interface{}) (goscript.Component, error) {
						time.Sleep(2 * time.Second)
						return goscript.NewServerComponent("SlowData", func(props goscript.Props) string {
							return `<div class="card" style="border-left:3px solid #e94560;"><h2>Async Data Loaded!</h2><p>This content was streamed in after a 2-second delay while the skeleton was shown.</p></div>`
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
			appRouter.ServeHTTP(w, r)
		}
	})

	// ---------------------------------------------------------------------------
	// 9. Serve static assets from the ./static directory
	// ---------------------------------------------------------------------------
	assetManager := goscript.NewAssetManager("./static")
	staticHandler := assetManager.ServeAssets("/static")

	// Final mux combines static assets with the main handler
	finalMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) >= 8 && r.URL.Path[:8] == "/static/" {
			staticHandler.ServeHTTP(w, r)
			return
		}
		pipeline.Execute(w, r, mainHandler)
	})

	// ---------------------------------------------------------------------------
	// 10. Start the server
	// ---------------------------------------------------------------------------
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf(`
╔══════════════════════════════════════════╗
║         GoScript v2.0.0 Server          ║
║  Framework: %s            ║
║  Version:    %s                     ║
║  Port:       %d                       ║
╠══════════════════════════════════════════╣
║  Routes:                               ║
║    GET  /              (App Router)     ║
║    GET  /about         (App Router)     ║
║    GET  /users/:id     (Dynamic)        ║
║    GET  /api/status    (API + Rate)     ║
║    GET  /api/users/:id (API + CORS)     ║
║    POST /api/users     (API + Auth)     ║
║    GET  /api/v2/health (Group)          ║
║    GET  /api/v2/features(Group)          ║
║    GET  /stream        (Streaming SSR)  ║
║    GET  /static/...    (SSG Assets)     ║
║  Middleware:                           ║
║    Recovery → Logging → RequestID →    ║
║    Security → CORS → Gzip → RateLimit ║
╚══════════════════════════════════════════╝
`, goscript.FrameworkName, goscript.Version, port)

	log.Fatal(http.ListenAndServe(addr, finalMux))
}
