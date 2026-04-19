# GoScript - The Go-native Language for the Web

GoScript is a Go-native language and runtime for building web apps without splitting the product across Go and JavaScript. It follows a familiar app structure for people coming from modern web frameworks, but its real purpose is different: GoScript gives Go developers a native path to the web.

## Quick Start

```bash
git clone https://github.com/davidjeba/goscript.git
cd goscript
go mod tidy
go run cmd/server/main.go        # Start showcase server
go run examples/counter/main.go  # Start counter example
```

The server starts at `http://localhost:8080`. Set `PORT` to customize.

---

## The .gs Language

`.gs` is GoScript's own language — Go-like syntax that compiles to browser JavaScript. Write components, state, and event handlers in `.gs`, compile with `gopm`, and ship zero-config JavaScript.

```go
// counter.gs — compiles to browser JavaScript
package main

import "goscript/dom"
import "goscript/state"

func Counter() dom.Element {
    count, setCount := state.Use(0)

    return dom.CreateElement("div", dom.Props{"class": "counter"},
        dom.CreateElement("h1", nil, count),
        dom.CreateElement("button", dom.Props{
            "onclick": func(e dom.Event) { setCount(count + 1) },
        }, "+ Increment"),
    )
}

func main() {
    dom.Mount("#app", Counter())
}
```

Compile with the CLI:

```bash
gopm compile counter.gs -o counter.js    # Single file
gopm build                                # All .gs files → ./build
gopm dev                                  # Dev server with HMR
```

The compilation pipeline: `Lexer.Tokenize()` → `Parser.Parse()` → `Generator.Generate()`

---

## Reactive Attributes

GoScript ships a client runtime that makes any HTML element reactive via data attributes. **Zero JavaScript required** — the runtime intercepts attributes and handles AJAX + DOM swapping automatically.

| Attribute | Description | Example |
|---|---|---|
| `gs-trigger` | Event to listen for | `gs-trigger="click"` |
| `gs-get` | URL to fetch (HTML fragment) | `gs-get="/api/counter?n=1"` |
| `gs-target` | CSS selector to update | `gs-target="#count"` |
| `gs-swap` | How to insert the response | `gs-swap="innerHTML"` |
| `gs-indicator` | CSS selector for loading state | `gs-indicator=".spinner"` |

```html
<div id="count">0</div>
<button gs-trigger="click"
        gs-get="/api/counter?n=1"
        gs-target="#count"
        gs-swap="innerHTML">
  + Increment
</button>
```

Include the runtime: `<script src="/__goscript/runtime.js"></script>`

---

## gopm CLI

The `gopm` command handles .gs compilation, dev server, builds, and project scaffolding.

```bash
gopm dev                        # Start dev server with HMR on :8080
gopm compile input.gs -o out.js # Compile single .gs → .js
gopm compile app.gs --watch     # Watch and recompile on change
gopm build                      # Compile all .gs files → ./build
gopm build --minify -o ./dist   # Minified production build
gopm init myapp                 # Scaffold new project
gopm version                    # Show version
gopm help                       # Show help
```

### Dev Server Flags

| Flag | Default | Description |
|---|---|---|
| `-p, --port` | 8080 | Server port |
| `--hmr` | true | Enable Hot Module Replacement |
| `--static` | ./static | Static files directory |
| `--open` | false | Open browser on start |

### Build Flags

| Flag | Default | Description |
|---|---|---|
| `-o, --output` | ./build | Output directory |
| `--minify` | false | Minify output JavaScript |

---

## Project Structure

```
goscript/
├── cmd/
│   ├── server/
│   │   └── main.go              # Showcase server (all features)
│   └── gopm/
│       └── main.go              # CLI: compile, dev, build, init
├── pkg/
│   ├── goscript/                # Core framework
│   │   ├── app-router.go        # App Router, layouts, route groups
│   │   ├── api-routes.go        # API router, rate limit, auth, CORS
│   │   ├── middleware-pipeline.go # Composable middleware
│   │   ├── metadata-seo.go      # Metadata/SEO builder
│   │   ├── streaming-ssr.go     # Streaming SSR with Suspense
│   │   ├── server-client-components.go # Server & Client Components
│   │   ├── ssg-isr.go           # Static generation + ISR
│   │   ├── hmr.go               # Dev server with WebSocket HMR
│   │   ├── error-loading-boundaries.go # Error & loading states
│   │   ├── component.go         # Component model, Props, State
│   │   ├── router.go            # v1 Router with dynamic segments
│   │   ├── state.go             # Store, GlobalStore, UseState
│   │   ├── ssr.go               # SSR engine
│   │   ├── hooks.go             # useState, useEffect, useContext, etc.
│   │   └── types.go             # Version, feature flags
│   ├── gscompiler/              # .gs → .js compiler
│   ├── gocsx/                   # CSS-in-Go framework
│   ├── goscale/                 # API, DB, Edge
│   ├── gouix/                   # Canvas UI framework
│   └── jetpack/                 # Performance monitoring
├── examples/
│   └── counter/
│       ├── main.go              # Minimal reactive counter (Go)
│       ├── counter.gs           # Client component (.gs → .js)
│       └── go.mod
└── README.md
```

---

## Framework API Reference

### App Router

```go
router := goscript.NewAppRouter("/")
router.RegisterRoute("/", homeHandler, []string{"GET"})
router.RegisterRoute("/users/:id", userHandler, []string{"GET"})
router.RegisterRoute("/docs/*path", docsHandler, []string{"GET"})

// Layouts wrap page content
router.SetLayout("/", rootLayout)
router.SetLayout("/dashboard", dashboardLayout)

// Route groups share a prefix
api := router.Group("/api/v1")
router.RegisterGroupRoute(api, "/users", usersHandler, []string{"GET"})
```

### API Routes

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

api.SetNotFound(func(ctx *goscript.APIContext) (interface{}, error) {
    return map[string]interface{}{"error": "not found"}, nil
})
```

### Middleware Pipeline

```go
pipeline := goscript.NewPipeline().
    Use(goscript.RequestIDMiddleware()).
    Use(goscript.GzipMiddleware()).
    Use(goscript.CORSMiddleware(goscript.DefaultCORSConfig())).
    Use(goscript.SecurityHeadersMiddleware()).
    Use(goscript.LoggingMiddleware(log.Printf)).
    Use(goscript.RecoveryMiddleware(log.Printf))

pipeline.Execute(w, r, finalHandler)
```

Built-in middleware: Gzip, CORS, Security Headers, Rate Limiting, Recovery, Logging, Request ID, Sessions.

### Metadata / SEO

```go
metadata := goscript.NewMetadata().
    SetTitle("My Page").
    SetDescription("A page built with GoScript").
    SetCanonical("https://example.com").
    SetThemeColor("#10b981").
    SetOpenGraph(goscript.OpenGraphMeta{
        Title: "My Page", Type: "website", Image: "/og.png",
    }).
    SetTwitter(goscript.TwitterMeta{
        Card: "summary_large_image", Title: "My Page",
    }).
    AddJSONLD(map[string]interface{}{"@type": "Article"}).
    Build()

headHTML := metadata.Render()
```

### Server Components

```go
page := goscript.NewServerComponent("Home", func(props goscript.Props) string {
    return `<h1>Hello from GoScript</h1>`
}, nil)
html := page.Render()
```

### Streaming SSR

```go
engine := goscript.NewStreamSSREngine(store)
boundaries := []goscript.SuspenseBoundary{{
    ID:       "comments",
    Fallback: goscript.NewSkeletonLoader(3),
    Loader: func(ctx interface{}) (goscript.Component, error) {
        return loadComments(), nil
    },
}}
engine.RenderStream(w, r, page, boundaries)
```

### SSG / ISR

```go
engine := goscript.NewSSGEngine("./out")
engine.AddPage(goscript.PageConfig{
    Path:       "/",
    Component:  homePage,
    RenderMode: goscript.RenderSSG,
})
engine.AddPage(goscript.PageConfig{
    Path:       "/blog",
    Component:  blogPage,
    RenderMode: goscript.RenderISR,
    Revalidate: 60 * time.Second,
})
engine.Build(context.Background())
```

---

## Counter Example

A complete, minimal GoScript app in ~80 lines of Go:

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    goscript "github.com/davidjeba/goscript/pkg/goscript"
)

func main() {
    port := "8080"
    if p := os.Getenv("PORT"); p != "" { port = p }

    router := goscript.NewAppRouter("/")
    router.RegisterRoute("/", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
        html := `<!DOCTYPE html><html><head><title>Counter</title>
        <script src="/__goscript/runtime.js"></script></head>
        <body><div id="count">0</div>
        <button gs-trigger="click" gs-get="/api/counter?n=1"
                gs-target="#count" gs-swap="innerHTML">+</button>
        </body></html>`
        w.Header().Set("Content-Type", "text/html")
        w.Write([]byte(html))
    }, []string{"GET"})

    count := 0
    http.HandleFunc("/api/counter", func(w http.ResponseWriter, r *http.Request) {
        if n := r.URL.Query().Get("n"); n != "" {
            fmt.Sscanf(n, "%d", &count)
        }
        w.Header().Set("Content-Type", "text/html")
        fmt.Fprintf(w, "%d", count)
    })

    log.Fatal(http.ListenAndServe(":"+port, router))
}
```

Run: `go run examples/counter/main.go`

---

## GoScript vs JavaScript

GoScript is not a Next.js replacement. The familiar app structure is there to reduce learning curve, but the actual comparison that matters is language ownership.

| Concern | GoScript | JavaScript |
|---|---|---|
| **Language role** | Go-native language and runtime for the web | Browser-native language and ecosystem |
| **Team workflow** | One Go mental model across product layers | Often split across Go plus JS/TS |
| **Deployment** | Go-shaped deliverable, simple ops story | Runtime plus package graph |
| **Startup** | Compiled and predictable | Dynamic and runtime-driven |
| **Memory** | Small, Go-shaped footprint | Usually larger due to runtime and tooling |
| **Concurrency** | Goroutines and Go primitives | Event loop and async patterns |
| **Type Safety** | Compile-time checks and generics | TypeScript layer, optional and separate |
| **Own Language** | GoScript language, designed for Go teams | JavaScript/TypeScript already native to the browser |
| **Reactivity** | GoScript runtime + reactive attributes | React state/hooks or other JS frameworks |
| **Best fit** | Go teams that want the web to feel native | Teams that want direct access to the JS ecosystem |

---

## License

GoScript is released under the **MIT License**. See [LICENSE](./LICENSE) for details.
