# GoScript 2.0 — Full-Stack Go Web Framework

**GoScript** is a Next.js-inspired full-stack web framework written entirely in Go. Version 2.0 brings 10 major new features including an App Router, Streaming SSR, Server/Client Components, and more — all while maintaining full backward compatibility with v1.

## Table of Contents

- [Overview](#overview)
- [Feature Comparison](#feature-comparison)
- [Quick Start](#quick-start)
- [Project Structure](#project-structure)
- [v2 Features](#v2-features)
- [v1 Sub-Modules](#v1-sub-modules)
- [Why GoScript over Next.js?](#why-goscript-over-nextjs)
- [License](#license)

---

## Overview

GoScript provides a complete framework for building modern web applications in Go, featuring:

- **Component System** — React-like component model with Props, State, Lifecycle hooks, and JSX-like templating
- **Server-Side Rendering** — Full SSR support with state hydration and streaming
- **App Router** — File-system based routing with layouts, dynamic segments, and route groups
- **API Routes** — Convention-based API endpoints with built-in rate limiting, auth, and CORS
- **SSG/ISR** — Static site generation with incremental revalidation
- **Hot Module Replacement** — Dev server with WebSocket-based live reload
- **Sub-Modules** — Gocsx, GoScale, GOPM, GoUIX, and Jetpack

---

## Feature Comparison

| Feature | GoScript v1 | GoScript v2 | Next.js 16 |
|---|---|---|---|
| Component Model | ✅ Props, State, Lifecycle | ✅ Props, State, Lifecycle + Server/Client Components | ✅ Server/Client Components |
| Routing | ✅ Basic Router with `:param` | ✅ App Router + Dynamic + Catch-all + Groups | ✅ App Router |
| Layouts | ❌ Manual | ✅ Nested Layout Conventions | ✅ Nested Layouts |
| SSR | ✅ Basic SSREngine | ✅ Streaming SSR with Suspense | ✅ Streaming SSR |
| SSG | ❌ | ✅ Static Site Generation | ✅ Static Export |
| ISR | ❌ | ✅ Incremental Static Regeneration | ✅ ISR |
| API Routes | ❌ Manual handlers | ✅ Convention-based with Rate Limit, Auth, CORS | ✅ Route Handlers |
| Middleware | ✅ Basic chaining | ✅ Composable Pipeline (Gzip, CORS, Security, Recovery, Rate Limit, Session) | ✅ Middleware |
| Error Boundaries | ❌ | ✅ ErrorBoundary with Fallback | ✅ error.tsx |
| Loading States | ❌ | ✅ LoadingBoundary + Skeleton Loaders | ✅ loading.tsx |
| Metadata/SEO | ❌ | ✅ Fluent Builder (OG, Twitter, JSON-LD, Robots) | ✅ generateMetadata |
| HMR | ❌ | ✅ WebSocket-based Dev Server | ✅ Fast Refresh |
| State Management | ✅ GlobalStore + Hooks | ✅ GlobalStore + Hooks | ✅ Zustand/Redux |
| JSX Templating | ✅ GSX Compiler | ✅ GSX Compiler | ✅ JSX/TSX |
| Canvas UI | ✅ GoUIX (2D Canvas) | ✅ GoUIX (2D Canvas) | ❌ |
| CSS Framework | ✅ Gocsx | ✅ Gocsx | ✅ CSS Modules/Tailwind |
| Package Manager | ✅ GOPM | ✅ GOPM | ✅ npm/pnpm |
| API System | ✅ GoScale API (GraphQL-like) | ✅ GoScale API + v2 API Router | ❌ |
| Database | ✅ GoScale DB | ✅ GoScale DB | ✅ Prisma |
| Edge Computing | ✅ GoScale Edge | ✅ GoScale Edge | ✅ Edge Runtime |
| Performance Monitor | ✅ Jetpack (Core/Frontend/Security) | ✅ Jetpack (Core/Frontend/Security) | ❌ |
| Language | Go 1.17 | **Go 1.22** | TypeScript |

---

## Quick Start

### Prerequisites

- Go 1.22 or later
- Git

### Installation

```bash
git clone https://github.com/davidjeba/goscript.git
cd goscript
go mod tidy
go run cmd/server/main.go
```

The server starts on `http://localhost:8080` by default. Set the `PORT` environment variable or use `-port` flag to customize.

### Your First Page

```go
package main

import (
    "fmt"
    "net/http"
    "github.com/davidjeba/goscript/pkg/goscript"
)

func main() {
    router := goscript.NewAppRouter("")

    router.RegisterRoute("/", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
        page := goscript.NewServerComponent("Home", func(props goscript.Props) string {
            return `<h1>Hello from GoScript v2!</h1>`
        }, nil)

        w.Header().Set("Content-Type", "text/html")
        fmt.Fprint(w, page.Render())
    }, []string{"GET"})

    http.ListenAndServe(":8080", router)
}
```

---

## Project Structure

```
goscript/
├── cmd/
│   ├── server/
│   │   └── main.go              # v2 demo server showcasing all features
│   ├── gocsx_demo/
│   │   └── main.go              # Gocsx CSS framework demo
│   ├── gocsx_2d_demo/
│   │   └── main.go              # Gocsx 2D canvas demo
│   ├── gocsx_3d_demo/
│   │   └── main.go              # Gocsx 3D WebGPU demo
│   ├── gopm/
│   │   └── main.go              # GOPM package manager CLI
│   ├── goscale_demo/
│   │   └── main.go              # GoScale API demo
│   └── gouix_demo/
│       └── main.go              # GoUIX canvas UI demo
├── pkg/
│   ├── goscript/                 # Core framework (v1 + v2)
│   │   ├── component.go         # v1: Component, BaseComponent, Props, CreateElement
│   │   ├── router.go            # v1: Router, RouteHandler, dynamic segments
│   │   ├── state.go             # v1: Store, GlobalStore, UseState
│   │   ├── ssr.go               # v1: SSREngine, RenderToString
│   │   ├── context.go           # v1: Context, ContextProvider, ContextConsumer
│   │   ├── assets.go            # v1: AssetManager
│   │   ├── lifecycle.go         # v1: LifecycleComponent, ComponentRegistry
│   │   ├── hooks.go             # v1: useState, useEffect, useContext, useMemo, useCallback, useRef
│   │   ├── jsx.go               # v1: JSXParser, GSXCompiler, TranspileGSX
│   │   ├── app-router.go        # v2: AppRouter, RouteGroup, layout nesting
│   │   ├── streaming-ssr.go     # v2: StreamSSREngine, SuspenseBoundary
│   │   ├── server-client-components.go  # v2: ServerComponent, ClientComponent
│   │   ├── api-routes.go        # v2: APIRouter, rate limiting, auth, CORS
│   │   ├── middleware-pipeline.go # v2: Pipeline, Gzip, CORS, Security, Recovery
│   │   ├── ssg-isr.go           # v2: SSGEngine, ISR revalidation
│   │   ├── error-loading-boundaries.go # v2: ErrorBoundary, LoadingBoundary, Skeleton
│   │   ├── metadata-seo.go      # v2: MetadataBuilder, OpenGraph, Twitter, JSON-LD
│   │   ├── hmr.go               # v2: DevServer, WebSocket HMR, file watching
│   │   └── types.go             # v2: Version, FeatureFlags, framework info
│   ├── components/
│   │   ├── home.go              # Sample home page component
│   │   ├── home.gsx             # GSX template for home page
│   │   ├── counter.go           # Counter component (Go)
│   │   ├── counter_gouix.go     # Counter component (GoUIX)
│   │   ├── counter_gouix_test.go
│   │   ├── counter_test.go
│   │   └── home_gouix.go        # Home component (GoUIX)
│   ├── gocsx/                   # CSS-in-Go framework
│   │   └── gocsx.go
│   ├── gopm/                    # Go Package Manager
│   │   └── package_manager.go
│   ├── goscale/
│   │   ├── api/
│   │   │   └── goscale.go       # GraphQL/gRPC hybrid API
│   │   ├── db/
│   │   │   └── goscaledb.go     # PostgreSQL-compatible database
│   │   └── edge/
│   │       └── edge.go          # Edge computing network
│   ├── gouix/                   # Interactive Canvas UI
│   │   ├── component.go         # Base component with drag & touch
│   │   ├── canvas.go            # SVG canvas rendering
│   │   ├── hyper.go             # Reactive signals and effects
│   │   ├── canvas_test.go
│   │   ├── component_test.go
│   │   └── reactive_test.go
│   └── jetpack/
│       ├── core/
│       │   └── jetpack.go       # Performance metrics engine
│       ├── frontend/
│       │   ├── lighthouse.go    # Lighthouse integration
│       │   └── performance_panel.go  # Floating performance panel
│       └── security/
│           └── security_monitor.go   # Vulnerability scanning
├── api/
│   └── index.go                 # API entry point
├── docs/
│   ├── component-system.md
│   └── gouix.md
├── static/
│   └── index.html               # Static HTML assets
├── github/workflows/
│   └── v0-deploy.yml
├── go.mod                       # Go 1.22
├── vercel.json
├── .gitignore
├── LICENSE                      # MIT
└── README.md
```

---

## v2 Features

### 1. App Router (`app-router.go`)

File-system based routing with nested layouts, dynamic segments, catch-all routes, and route groups.

```go
app := goscript.NewAppRouter("/app")

// Register routes with HTTP methods
app.RegisterRoute("/", homeHandler, []string{"GET"})
app.RegisterRoute("/users/:id", userHandler, []string{"GET"})
app.RegisterRoute("/docs/*path", docsHandler, []string{"GET"})

// Nested layouts
app.SetLayout("/", rootLayout)
app.SetLayout("/dashboard", dashboardLayout)

// Route groups
api := app.Group("/api/v1")
app.RegisterGroupRoute(api, "/users", usersHandler, []string{"GET"})
```

### 2. Streaming SSR (`streaming-ssr.go`)

Chunked HTML streaming with Suspense boundaries for progressive page loading.

```go
engine := goscript.NewStreamSSREngine(store)

boundaries := []goscript.SuspenseBoundary{
    {
        ID:       "comments",
        Fallback: goscript.NewSkeletonLoader(3),
        Loader: func(ctx interface{}) (goscript.Component, error) {
            // Async data loading
            return loadComments(), nil
        },
    },
}

engine.RenderStream(w, r, pageComponent, boundaries)
```

### 3. Server & Client Components (`server-client-components.go`)

Server components render pure HTML (zero JS). Client components include hydration markers and event handlers.

```go
// Server Component — no JavaScript shipped to the client
serverComp := goscript.NewServerComponent("Header", func(props goscript.Props) string {
    return `<header><h1>Pure Server Rendered</h1></header>`
}, nil)

// Client Component — includes hydration marker and event handlers
clientComp := goscript.NewClientComponent("Counter", goscript.Props{"count": 0})
clientComp.OnEvent("click", "incrementCounter()")
```

### 4. API Routes (`api-routes.go`)

Convention-based API router with rate limiting, authentication, CORS, and automatic JSON responses.

```go
api := goscript.NewAPIRouter()

api.GET("/api/users", listUsers,
    goscript.WithRateLimit(100, time.Minute),
    goscript.WithCORS("*"),
)

api.POST("/api/users", createUser,
    goscript.WithRateLimit(10, time.Minute),
    goscript.WithAuth([]string{"Bearer"}, validateToken),
)
```

### 5. Middleware Pipeline (`middleware-pipeline.go`)

Composable middleware with built-in implementations for Gzip, CORS, Security Headers, Rate Limiting, Recovery, Logging, Sessions, and more.

```go
pipeline := goscript.NewPipeline().
    Use(goscript.RecoveryMiddleware(log.Printf)).
    Use(goscript.LoggingMiddleware(log.Printf)).
    Use(goscript.RequestIDMiddleware()).
    Use(goscript.SecurityHeadersMiddleware()).
    Use(goscript.CORSMiddleware(goscript.DefaultCORSConfig())).
    Use(goscript.GzipMiddleware()).
    Use(goscript.RateLimitMiddleware(1000, 50))

pipeline.Execute(w, r, finalHandler)
```

### 6. SSG & ISR (`ssg-isr.go`)

Static site generation with incremental static regeneration for dynamic content.

```go
engine := goscript.NewSSGEngine("./out")

engine.AddPage(goscript.PageConfig{
    Path:       "/",
    Component:  homeComponent,
    RenderMode: goscript.RenderSSG,
})

engine.AddPage(goscript.PageConfig{
    Path:       "/blog",
    Component:  blogComponent,
    RenderMode: goscript.RenderISR,
    Revalidate: 60 * time.Second,
    Params:     []map[string]string{{"slug": "hello"}, {"slug": "world"}},
})

engine.Build(context.Background())
```

### 7. Error & Loading Boundaries (`error-loading-boundaries.go`)

Graceful error handling and skeleton loading states for component rendering.

```go
// Error Boundary — catches panics and renders fallback
errorBoundary := goscript.NewErrorBoundary(
    goscript.NewErrorComponent("Something went wrong"),
    riskyComponent,
    anotherComponent,
)

// Loading Boundary — shows skeleton while content loads
loadingBoundary := goscript.NewLoadingBoundary(
    goscript.NewSkeletonLoader(4),
    func(ctx context.Context) (goscript.Component, error) {
        return loadAsyncContent(), nil
    },
)
loadingBoundary.Load(ctx)
```

### 8. Metadata & SEO API (`metadata-seo.go`)

Fluent builder for generating complete `<head>` content with Open Graph, Twitter Cards, JSON-LD, and robot directives.

```go
metadata := goscript.NewMetadata().
    SetTitle("My Page").
    SetDescription("A page built with GoScript").
    SetCanonical("https://example.com/page").
    SetThemeColor("#ffffff").
    AddKeywords("goscript", "go", "framework").
    SetRobots(true, true).
    SetOpenGraph(goscript.OpenGraphMeta{
        Title: "My Page", Type: "article", Image: "https://example.com/og.png",
    }).
    SetTwitter(goscript.TwitterMeta{
        Card: "summary_large_image", Title: "My Page",
    }).
    AddJSONLD(map[string]interface{}{
        "@type": "Article", "headline": "My Page",
    }).
    Build()

headHTML := metadata.Render()
```

### 9. HMR Dev Server (`hmr.go`)

Hot Module Replacement development server with WebSocket-based live updates and file watching.

```go
devServer := goscript.NewDevServer(3001, router)
devServer.Watch("./pkg", "./cmd")
devServer.OnFileChange(func(path string) {
    log.Printf("File changed: %s", path)
})
devServer.Start(ctx)
```

### 10. Types & Feature Flags (`types.go`)

Central version info, feature flags, and framework metadata for diagnostics.

```go
fmt.Println(goscript.Version)           // "2.0.0"
fmt.Println(goscript.FrameworkName)     // "GoScript"
fmt.Println(goscript.IsFeatureEnabled("hmr")) // true
info := goscript.GetFrameworkInfo()
```

---

## v1 Sub-Modules

### Gocsx — CSS-in-Go Framework
A utility-first CSS framework for Go, providing Button, Card, and layout components with platform adapters. Compile-time CSS generation with zero runtime cost.

### GoScale API — GraphQL/gRPC Hybrid API
A high-performance API system combining gRPC-like performance with GraphQL-like flexibility. Supports schema definitions, resolvers, subscriptions, edge computing, and database integration.

### GoScale DB — PostgreSQL-Compatible Database
A high-performance database with caching, sharding, replication, time-series support, relationship management, and NoCode schema creation. Built on top of PostgreSQL with query caching and automatic optimization.

### GoScale Edge — Edge Computing Network
Distributed edge computing with worker pools, load balancing (round-robin, least-loaded, fastest), health checking, caching, and automatic synchronization with the parent API.

### GOPM — Go Package Manager
A comprehensive package manager with install, update, audit, publish, dependency resolution, caching, and project scaffolding commands. Also includes commands for CSS building, WebGPU shader compilation, 2D canvas, and UIX component management.

### GoUIX — Interactive Canvas UI Framework
A reactive UI framework with SVG-based canvas rendering, drag-and-drop, touch gestures, and hyper(reactive) signals with computed values and effects. Supports both 2D and 3D (WebGPU) rendering.

### Jetpack Core — Performance Monitoring
Real-time performance metrics engine supporting frontend (FPS, page load, Web Vitals), backend (API latency, throughput, error rates), database (query time, connection pools), and security metrics with alerting and export.

### Jetpack Frontend — Lighthouse & Performance Panel
Google Lighthouse integration for automated performance audits and a floating performance panel with real-time metrics display, charts, alerts, and configurable themes.

### Jetpack Security — Vulnerability Scanning
Automated security scanning for XSS, SQL injection, CSRF, outdated libraries, missing headers, TLS configuration, and anomaly detection with compliance checking and reporting.

---

## Why GoScript over Next.js?

| Concern | GoScript | Next.js |
|---|---|---|
| **Language** | Go — compiled, memory-safe, fast | JavaScript/TypeScript — interpreted/JIT |
| **Deployment** | Single binary, zero dependencies | Node.js runtime required |
| **Startup Time** | Milliseconds | Seconds |
| **Memory Usage** | ~20MB typical | ~200MB+ typical |
| **Concurrency** | Goroutines — millions of lightweight concurrent tasks | Event loop — single-threaded |
| **Type Safety** | Compile-time with generics (Go 1.22) | TypeScript (optional) |
| **Binary Size** | ~10MB static binary | ~100MB+ node_modules |
| **Learning Curve** | Go basics + framework concepts | React + Next.js + TypeScript + Node ecosystem |
| **Runtime Errors** | Virtually eliminated by compiler | Common (null refs, type mismatches) |
| **Cross-Compilation** | `GOOS=linux go build` — one command | Docker / platform-specific builds |
| **Standard Library** | Extensive built-in HTTP, crypto, encoding | Requires npm packages |
| **Database** | GoScale DB built-in | Prisma/Drizzle (external) |
| **Edge Computing** | GoScale Edge built-in | Vercel Edge only |
| **CSS Framework** | Gocsx built-in | Tailwind (external) |
| **Package Manager** | GOPM built-in | npm/pnpm (external) |

---

## License

GoScript is released under the **MIT License**. See [LICENSE](./LICENSE) for details.
