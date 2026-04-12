# GoScript 2.0 — Full-Stack Go Web Framework

> A Next.js-inspired full-stack web framework written in Go, designed to challenge React/Next.js with compiled performance, zero cold starts, and built-in everything.

[![Go Report Card](https://goreportcard.com/badge/github.com/davidjeba/goscript)](https://goreportcard.com/report/github.com/davidjeba/goscript)

## What's New in 2.0

GoScript 2.0 is a **complete rewrite** that brings every major Next.js feature to Go:

| Feature | Original goscript | GoScript 2.0 | Next.js 16 |
|---------|------------------|---------------|------------|
| File-System Routing | Manual registration | Auto-discovery (page.go, layout.go, loading.go, error.go) | Yes |
| Streaming SSR | None | Chunked HTML + Suspense boundaries | Yes |
| Server & Client Components | None | ServerComponent (zero JS) + ClientComponent (hydrated) | Yes |
| API Routes | Manual handlers | Convention-based (.api.go) with CORS, rate limiting, auth | Yes |
| Middleware Pipeline | Basic chain | Composable pipeline with 8+ built-in handlers | Yes |
| SSG / ISR | None | Static generation + incremental regeneration | Yes |
| Error Boundaries | None | ErrorBoundary + LoadingBoundary + SkeletonLoader | Yes |
| Metadata & SEO API | None | Fluent builder (OG, Twitter, JSON-LD, robots) | Yes |
| HMR Dev Server | None | WebSocket-based Hot Module Replacement | Yes |
| CLI Tool | Placeholder | Full gopm CLI (init, dev, build, generate) | Yes |
| SSR Performance | Template-based | Compiled Go binary (no runtime) | V8 runtime |
| Cold Start | Fast | Instant (compiled binary) | Slow (Node.js boot) |
| Bundle Size | N/A | ~5KB hydration runtime | Varies |
| Built-in CSS | Gocsx | Gocsx (retained) | Tailwind CSS |

## Quick Start

```bash
# Clone the repository
git clone https://github.com/davidjeba/goscript.git
cd goscript

# Run the server
go run cmd/server/main.go

# Or use the CLI
go run cmd/gopm/main.go init my-app
cd my-app
go run cmd/gopm/main.go dev
```

## Project Structure

```
goscript/
├── cmd/
│   ├── server/main.go          # HTTP server entrypoint
│   └── gopm/main.go            # CLI tool (init, dev, build, generate)
├── pkg/goscript/
│   ├── types.go                # Core types: Component, Props, Store, CreateElement
│   ├── app-router.go           # File-system based routing (App Router)
│   ├── streaming-ssr.go        # Streaming SSR with Suspense boundaries
│   ├── server-client-components.go  # Server & Client component system
│   ├── api-routes.go           # Convention-based API routes
│   ├── middleware-pipeline.go  # Composable middleware chain
│   ├── ssg-isr.go              # Static site generation & incremental regeneration
│   ├── error-loading-boundaries.go  # Error & loading boundary system
│   ├── metadata-seo.go         # Metadata & SEO API
│   └── hmr.go                  # Hot Module Replacement dev server
├── app/
│   ├── layout.go               # Root layout component
│   └── page.go                 # Home page component
├── api/
│   └── hello.go                # Example API route
└── go.mod
```

## Core Packages

### 1. App Router (`app-router.go`)

File-system based routing that mirrors Next.js App Router conventions:

```go
router := goscript.NewAppRouter("/")
router.RegisterRoute("/", homeHandler, []string{"GET"})
router.RegisterRoute("/posts/:id", postHandler, []string{"GET"})
router.RegisterRoute("/docs/*path", docsHandler, []string{"GET"})
```

Supports:
- Dynamic segments: `[id]` → `:id`
- Catch-all segments: `[...slug]` → `*slug`
- Route groups: `(marketing)/` — no URL segment
- Convention files: `page.go`, `layout.go`, `loading.go`, `error.go`, `middleware.go`

### 2. Streaming SSR (`streaming-ssr.go`)

Chunked HTML streaming with Suspense boundaries:

```go
ssr := goscript.NewStreamSSREngine(store)
ssr.RenderStream(w, r, mainComponent, []goscript.SuspenseBoundary{
    {
        ID:       "user-profile",
        Fallback: skeletonComponent,
        Loader:   func(ctx context.Context) (Component, error) {
            data := fetchUserProfile(ctx)
            return ProfileComponent(data), nil
        },
    },
})
```

### 3. Server & Client Components (`server-client-components.go`)

```go
// Server component — renders on the server, ships ZERO JS
profile := goscript.NewServerComponent("Profile", func(props Props) string {
    return fmt.Sprintf("<div>%s</div>", props["name"])
}, userData)

// Client component — hydrates on the client with minimal JS
counter := goscript.NewClientComponent("Counter", props)
counter.OnEvent("click", "increment()")
```

### 4. API Routes (`api-routes.go`)

```go
api := goscript.NewAPIRouter()

api.GET("/api/users/:id", func(ctx *goscript.APIContext) (interface{}, error) {
    return map[string]string{"id": ctx.Params["id"]}, nil
}, goscript.WithRateLimit(100, time.Minute), goscript.WithAuth([]string{"bearer"}, nil))

api.POST("/api/users", createUserHandler, goscript.WithBodyValidator(UserSchema{}))
```

### 5. Middleware Pipeline (`middleware-pipeline.go`)

```go
pipeline := goscript.NewPipeline().
    Use(goscript.GzipMiddleware()).
    Use(goscript.RequestIDMiddleware()).
    Use(goscript.CORSMiddleware(goscript.DefaultCORSConfig())).
    Use(goscript.SecurityHeadersMiddleware()).
    Use(goscript.RecoveryMiddleware(log.Printf)).
    Use(goscript.LoggingMiddleware(log.Printf)).
    Use(goscript.SessionMiddleware("secret", sessionStore))

pipeline.Execute(w, r)
```

Built-in middleware: Gzip, RequestID, CORS, Security Headers, Recovery, Logger, Rate Limit, Session.

### 6. SSG & ISR (`ssg-isr.go`)

```go
engine := goscript.NewSSGEngine("./dist")

engine.AddPage(goscript.PageConfig{
    Path:       "/about",
    Component:  aboutComponent,
    RenderMode: goscript.RenderSSG,
})

engine.AddPage(goscript.PageConfig{
    Path:       "/blog/:slug",
    Component:  blogComponent,
    RenderMode: goscript.RenderISR,
    Revalidate: 5 * time.Minute,
    Params:     []map[string]string{{"slug": "hello"}, {"slug": "world"}},
})

engine.Build(context.Background())
```

### 7. Error & Loading Boundaries (`error-loading-boundaries.go`)

```go
boundary := goscript.NewErrorBoundary(
    goscript.NewSkeletonLoader(5),
    userProfileComponent,
    commentsComponent,
)
```

### 8. Metadata & SEO (`metadata-seo.go`)

```go
metadata := goscript.NewMetadata().
    SetTitle("My Page").
    SetDescription("A description of my page").
    SetCanonical("https://example.com/page").
    SetThemeColor("#10b981").
    AddKeywords("go", "web", "framework").
    SetOpenGraph(goscript.OpenGraphMeta{
        Title: "My Page", Type: "article",
        Image: "https://example.com/og.png",
    }).
    SetTwitter(goscript.TwitterMeta{
        Card: "summary_large_image",
    }).
    AddJSONLD(map[string]interface{}{
        "@type": "WebPage",
        "name":  "My Page",
    }).
    Build()

headHTML := metadata.Render() // Full <head> content
```

### 9. HMR Dev Server (`hmr.go`)

```go
dev := goscript.NewDevServer(8080, router)
dev.Watch("./app", "./api", "./pkg")
dev.OnFileChange(func(path string) {
    fmt.Println("Changed:", path)
})
dev.Start(context.Background())
```

### 10. CLI (`cmd/gopm/main.go`)

```
gopm init [name]       Create a new GoScript project
gopm dev               Start dev server with HMR
gopm build             Build for production
gopm start             Start production server
gopm generate page     Generate a new page
gopm generate component Generate a new component
gopm generate api      Generate a new API route
gopm generate middleware Generate new middleware
gopm generate layout   Generate a new layout
gopm lint              Run linter
gopm test              Run tests
gopm version           Show version
```

## Why GoScript over Next.js?

| Advantage | Detail |
|-----------|--------|
| **Compiled Performance** | Go compiles to a native binary — no V8 runtime, no JIT warmup |
| **Instant Cold Starts** | Zero ms cold start — perfect for serverless and edge |
| **Minimal JS Bundle** | ~5KB hydration runtime vs 50-200KB+ React bundles |
| **Type Safety** | Go's type system catches errors at compile time |
| **Single Binary** | Deploy one executable — no node_modules, no runtime |
| **Built-in Everything** | Router, SSR, API, middleware, CLI, CSS, 3D all included |
| **Lower Memory** | Go uses 10-50x less memory than Node.js |

## License

MIT License
