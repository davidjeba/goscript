# GoScript: A Comprehensive Web Development Framework

GoScript is a full-stack web development framework written in Go, designed to provide a complete solution for building modern web, mobile, and AR/VR applications. It combines the performance of Go with the flexibility of modern web technologies, offering a unified approach to frontend, backend, and database development.

[![GitHub stars](https://img.shields.io/github/stars/davidjeba/goscript?style=social)](https://github.com/davidjeba/goscript/stargazers)
[![GitHub license](https://img.shields.io/github/license/davidjeba/goscript)](https://github.com/davidjeba/goscript/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/davidjeba/goscript)](https://goreportcard.com/report/github.com/davidjeba/goscript)
[![GoDoc](https://godoc.org/github.com/davidjeba/goscript?status.svg)](https://godoc.org/github.com/davidjeba/goscript)

GoScript aims to be the most comprehensive Go-based web development framework, combining the best features of modern web frameworks with the performance and type safety of Go. It provides a complete solution for building web applications, from frontend to backend, with built-in support for CSS, 3D rendering, API development, database integration, and performance monitoring.

## 🌟 Core Components

GoScript consists of several integrated components that work together to provide a complete development experience:

### 1. Full-Stack Web Framework

GoScript provides a familiar experience for TypeScript and React developers:

- **Enhanced Component System**
  - Class-based components with lifecycle methods
  - Functional components
  - Props validation
  - Component state management
  - Context API for state sharing
  - Hooks system (useState, useEffect, useContext, useMemo, useCallback, useRef)
  - JSX-like syntax parser
  - Fragment support
  - Component testing utilities

- **Core Framework Features**
  - Server-side rendering (SSR)
  - Client-side hydration
  - Global state management
  - Routing with middleware support
  - Static asset management
  - Hot-reloading for development
  - CLI for component generation

### 2. Gocsx: Utility-First CSS Framework

Gocsx is a Tailwind-inspired CSS framework that provides a utility-first approach to styling web applications. It includes:

- **Utility Classes**: A comprehensive set of utility classes for rapid UI development
- **Component System**: Pre-built components with variants and responsive design
- **Theme Support**: Customizable themes with dark/light mode support
- **Platform Adapters**: Support for web, mobile, and other platforms
- **CSS Generation**: Efficient CSS generation with minimal output size

[Learn more about Gocsx](./README_GOCSX.md)

### 3. GoEngine: 2D/3D Rendering Engine

GoEngine provides a unified approach to 2D and 3D rendering, with support for:

- **WebGPU Integration**: Modern 3D graphics capabilities with:
  - A Go wrapper around the WebGPU API
  - Shader compilation and management
  - Resource management for GPU buffers and textures
  - A render pipeline abstraction
- **Canvas2D API**: Simplified 2D drawing and animation
- **Scene Graph**: Hierarchical scene management for both 2D and 3D
- **Three.js-like API**: Familiar API for 3D scene management
- **Performance Optimization**: 
  - Automatic detection of interactive applications
  - Performance throttling based on device capabilities
  - Unified API for both 2D and 3D contexts

### 4. GoScale: API and Database System

GoScale provides a high-performance API and database system with:

- **GraphQL-like Flexibility**: Define schemas and queries with GraphQL-like syntax
- **gRPC-like Performance**: High-performance binary protocol
- **Edge Computing**: Distributed API processing
- **Database Integration**: Support for PostgreSQL, NoSQL, and time-series data
- **Schema Management**: Automatic schema generation and migration

[Learn more about GoScale](./README_GOSCALE.md)

### 5. GOPM: Go Package Manager

GOPM is a comprehensive package manager for Go projects, with special support for the GoScript ecosystem:

- **Package Management**: Install, update, and manage dependencies
- **Project Setup**: Initialize and configure projects
- **Build Tools**: Build, test, and deploy applications
- **Framework Integration**: Special commands for Gocsx, GoEngine, GoScale, and more

[Learn more about GOPM](./README_GOPM.md)

### 6. Jetpack: Performance Monitoring and Optimization

Jetpack provides comprehensive performance monitoring and optimization tools:

- **Real-time Metrics**: Monitor FPS, memory usage, API latency, and more
- **Google Lighthouse Integration**: Run Lighthouse audits and track Core Web Vitals
- **Performance Panel**: Floating translucent panel for real-time metrics visualization
- **Chrome DevTools Extension**: Advanced performance monitoring in Chrome DevTools
- **Security Monitoring**: Track vulnerabilities, suspicious activities, and security compliance

[Learn more about Jetpack](./README_JETPACK.md)

## 🚀 Getting Started

### Installation

```bash
# Install GOPM
go install github.com/davidjeba/goscript/cmd/gopm@latest

# Initialize a new project
gopm setup my-project
cd my-project

# Install dependencies
gopm get
```

### Quick Start: Web Application

```bash
# Initialize a new web application
gopm setup --template web my-app
cd my-app

# Start the development server
gopm run dev
```

### Creating Web Components

#### Class-based Component

```go
type MyComponent struct {
    goscript.LifecycleComponentBase
}

func NewMyComponent(props goscript.Props) *MyComponent {
    base := goscript.NewBaseComponent(props, nil)
    component := &MyComponent{}
    component.LifecycleComponentBase.BaseComponent = *base
    return component
}

func (c *MyComponent) Render() string {
    return goscript.CreateElement("div", nil, 
        goscript.CreateElement("h1", nil, "Hello, World!"),
    )
}
```

#### Functional Component

```go
func MyComponent(props goscript.Props) string {
    return goscript.CreateElement("div", nil, 
        goscript.CreateElement("h1", nil, "Hello, World!"),
    )
}
```

### Creating a 2D Canvas Application

```go
package main

import (
        "log"
        "net/http"

        "github.com/davidjeba/goscript/pkg/gocsx"
        "github.com/davidjeba/goscript/pkg/gocsx/engine"
)

func main() {
        // Create a new Gocsx instance
        g := gocsx.New()

        // Create a new engine with 2D context
        e := engine.NewEngine(&engine.EngineConfig{
                Context: engine.Context2D,
        })

        // Create a new Canvas2D
        canvas := engine.NewCanvas2D("main-canvas", 800, 600, e)

        // Set render callback
        canvas.SetRenderCallback(func(ctx *engine.Canvas2DContext, deltaTime float64) {
                // Clear the canvas
                ctx.ClearRect(0, 0, 800, 600)

                // Draw a rectangle
                ctx.FillStyle = "#ff0000"
                ctx.FillRect(100, 100, 200, 150)
        })

        // Start the engine
        e.Start()

        // Start the server
        log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Creating a 3D WebGPU Application

```go
package main

import (
        "log"
        "net/http"

        "github.com/davidjeba/goscript/pkg/gocsx"
        "github.com/davidjeba/goscript/pkg/gocsx/engine"
)

func main() {
        // Create a new Gocsx instance
        g := gocsx.New()

        // Create a new engine with 3D context
        e := engine.NewEngine(&engine.EngineConfig{
                Context: engine.Context3D,
        })

        // Create a new WebGPU instance
        webgpu := engine.NewWebGPU()

        // Create a new Three.js scene
        scene := engine.NewThreeJSScene(e, webgpu)

        // Create a camera
        scene.CreateCamera("main-camera", "Main Camera", [3]float64{0, 0, 5}, [3]float64{0, 0, 0})

        // Create a cube
        scene.CreateCube("cube1", "Cube 1", [3]float64{0, 0, 0}, 1.0, [3]float64{1, 0, 0})

        // Start the engine
        e.Start()

        // Start the server
        log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Using GoScale API and Database

```go
package main

import (
        "log"

        "github.com/davidjeba/goscript/pkg/goscale/api"
        "github.com/davidjeba/goscript/pkg/goscale/db"
)

func main() {
        // Initialize database
        database, err := db.NewGoScaleDB(&db.Config{
                ConnectionString: "postgres://user:password@localhost:5432/mydb",
                TimeSeriesEnabled: true,
        })
        if err != nil {
                log.Fatalf("Failed to initialize database: %v", err)
        }

        // Define schema
        schema := api.NewSchema()
        schema.AddType("User", map[string]string{
                "id":    "ID!",
                "name":  "String!",
                "email": "String!",
                "posts": "[Post]",
        })
        schema.AddType("Post", map[string]string{
                "id":      "ID!",
                "title":   "String!",
                "content": "String!",
                "author":  "User!",
        })

        // Initialize API
        apiServer := api.NewGoScaleAPI(&api.Config{
                Schema:  schema,
                DB:      database,
                Port:    8080,
                EdgeEnabled: true,
        })

        // Start API server
        log.Fatal(apiServer.Start())
}
```

### Using Jetpack Performance Monitoring

```go
package main

import (
        "log"
        "net/http"

        "github.com/davidjeba/goscript/pkg/jetpack/core"
        "github.com/davidjeba/goscript/pkg/jetpack/frontend"
)

func main() {
        // Initialize Jetpack
        jp := core.NewJetpack()
        jp.EnableDevMode()

        // Create performance panel
        panel := frontend.NewPerformancePanel(jp)
        panel.Show()

        // Register metrics
        fps := 60.0
        jp.RegisterMetric(
                core.MetricFPS,
                "fps",
                "Frames per second",
                "fps",
                &fps,
                []string{"performance"},
        )

        // Initialize Lighthouse monitor
        lighthouse := frontend.NewLighthouseMonitor(jp)
        
        // Run Lighthouse audit
        _, err := lighthouse.RunAudit("http://localhost:8080")
        if err != nil {
                log.Printf("Failed to run Lighthouse audit: %v", err)
        }

        // Start exporting metrics
        jp.ExportEnabled = true
        jp.ExportEndpoint = "http://metrics.example.com"
        jp.StartExporting()

        // Start HTTP server
        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                // Record FPS metric
                jp.RecordMetric("fps", 58.5)
                
                // Serve HTML with performance panel
                html := `<!DOCTYPE html><html><body><h1>Hello World</h1></body></html>`
                htmlWithPanel, _ := panel.InjectIntoHTML(html)
                w.Header().Set("Content-Type", "text/html")
                w.Write([]byte(htmlWithPanel))
        })

        log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## 📚 Documentation

### Component System
- [Component System](docs/component-system.md)

### Gocsx CSS Framework
- [Gocsx Documentation](./README_GOCSX.md)

### GoScale API and Database
- [GoScale Documentation](./README_GOSCALE.md)

### GOPM Package Manager
- [GOPM Documentation](./README_GOPM.md)

### Jetpack Performance Monitoring
- [Jetpack Documentation](./README_JETPACK.md)

## 🏗️ Architecture

GoScript follows a modular architecture that allows each component to be used independently or together as a complete framework:

```
GoScript
├── Gocsx (CSS Framework)
│   ├── Core
│   │   ├── Configuration
│   │   ├── CSS Generator
│   │   └── Component System
│   ├── Platforms
│   │   ├── Web
│   │   ├── Mobile
│   │   └── Desktop
│   └── Components
│       ├── Button
│       ├── Card
│       └── ...
├── GoEngine (2D/3D Rendering)
│   ├── Core
│   │   ├── Engine
│   │   └── Scene Graph
│   ├── WebGPU
│   │   ├── Renderer
│   │   └── Shaders
│   └── Canvas2D
│       ├── Renderer
│       └── Sprites
├── GoScale (API and Database)
│   ├── API
│   │   ├── Schema
│   │   ├── Resolvers
│   │   └── Edge Computing
│   └── Database
│       ├── PostgreSQL
│       ├── NoSQL
│       └── TimeSeries
├── GOPM (Package Manager)
│   ├── Core
│   │   ├── Package Management
│   │   └── Dependency Resolution
│   └── Commands
│       ├── CSS Commands
│       ├── WebGPU Commands
│       ├── API Commands
│       └── DB Commands
└── Jetpack (Performance Monitoring)
    ├── Core
    │   ├── Metrics
    │   └── Panel
    ├── Frontend
    │   ├── Lighthouse
    │   └── Web Vitals
    ├── Backend
    │   ├── API Monitoring
    │   └── System Metrics
    └── Security
        ├── Vulnerability Scanning
        └── Anomaly Detection
```

## 🔧 Configuration

GoScript uses a unified configuration approach across all components:

```json
{
  "gocsx": {
    "theme": "default",
    "breakpoints": {
      "sm": "640px",
      "md": "768px",
      "lg": "1024px",
      "xl": "1280px"
    }
  },
  "engine": {
    "webgpu": {
      "enabled": true,
      "shaders": "./shaders"
    },
    "canvas2d": {
      "enabled": true,
      "sprites": "./sprites"
    }
  },
  "goscale": {
    "api": {
      "port": 8080,
      "edge-enabled": true
    },
    "db": {
      "connection-string": "localhost:5432",
      "time-series-enabled": true
    }
  },
  "jetpack": {
    "monitoring": {
      "enabled": true,
      "metrics": ["fps", "memory_usage", "api_latency"]
    },
    "panel": {
      "enabled": true,
      "position": "bottom-right",
      "opacity": 0.8
    }
  }
}
```

## 📋 Feature Comparison

### GoScript vs React

- **Language**: GoScript uses Go, React uses JavaScript/TypeScript
- **Performance**: GoScript offers better performance due to Go's efficiency
- **Type Safety**: GoScript has stronger type safety through Go's type system
- **Learning Curve**: Familiar API for React developers, but requires Go knowledge
- **Ecosystem**: React has a larger ecosystem, but GoScript integrates with Go libraries
- **CSS Framework**: GoScript includes Gocsx, React requires external libraries
- **3D Rendering**: GoScript includes WebGPU integration, React requires external libraries
- **API System**: GoScript includes GoScale, React requires external libraries
- **Performance Monitoring**: GoScript includes Jetpack, React requires external libraries

### Gocsx vs Tailwind CSS

- **Language**: Gocsx uses Go, Tailwind uses JavaScript/CSS
- **Platforms**: Gocsx supports web, mobile, and AR/VR, Tailwind is web-only
- **Type Safety**: Gocsx has type safety, Tailwind does not
- **Components**: Gocsx has built-in components, Tailwind requires additional libraries
- **Customization**: Both have powerful customization options

### GoScript WebGPU vs Three.js

- **Language**: GoScript uses Go, Three.js uses JavaScript
- **Integration**: GoScript offers tighter integration with the application
- **Performance**: GoScript can achieve better performance through Go
- **Type Safety**: GoScript has stronger type safety
- **Features**: Three.js has more features currently, but GoScript is rapidly evolving

## 🔄 Roadmap

- **Mobile Platform Adapter**: Native mobile support for iOS and Android
- **AR/VR Platform Adapter**: Support for AR and VR applications
- **Advanced Component Library**: Expanded set of UI components
- **Testing Infrastructure**: Comprehensive testing tools
- **IDE Integration**: Integration with popular IDEs
- **Animation System**: Advanced animation and transition system
- **Machine Learning Integration**: Integration with ML frameworks
- **Serverless Deployment**: Support for serverless deployment
- **Multi-tenant Support**: Built-in multi-tenant capabilities
- **Internationalization**: Built-in i18n support

## 📦 Examples

Check out the examples in the repository:

### Web Components
- `pkg/components/counter.go`: Demonstrates class-based components with state
- `pkg/components/home.go`: Shows how to use the context API and functional components

### CSS Framework
- `cmd/gocsx_demo`: Basic CSS framework demo

### 2D and 3D Applications
- `cmd/gocsx_2d_demo`: 2D canvas application demo
- `cmd/gocsx_3d_demo`: 3D WebGPU application demo

### API and Database
- `cmd/goscale_demo`: API and database demo

### Performance Monitoring
- `cmd/jetpack_demo`: Performance monitoring demo

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

MIT License
