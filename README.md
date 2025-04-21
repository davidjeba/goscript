# GoScript

GoScript is a comprehensive framework for building modern web, mobile, and AR/VR applications using Go. It combines the power of Go with modern web technologies to provide a seamless development experience.

## Core Components

### Full-Stack Web Framework

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

### Gocsx - CSS Framework

Gocsx (pronounced "gosix") is a utility-first CSS framework for Go, similar to Tailwind CSS but with extended capabilities for cross-platform development. It provides a unified styling approach for web, mobile, and AR/VR applications.

### WebGPU Integration

GoScript includes a powerful WebGPU integration that makes it easy to build 3D applications similar to Three.js directly in Go. The WebGPU integration provides:

- A Go wrapper around the WebGPU API
- Shader compilation and management
- Resource management for GPU buffers and textures
- A render pipeline abstraction

### Engine System

The engine system provides performance optimization for both 2D and 3D applications:

- Automatic detection of interactive applications
- Performance throttling based on device capabilities
- Unified API for both 2D and 3D contexts
- Scene graph system for both 2D and 3D

## Planned Features

- Multi-tenant and multi-language support
- Internationalization and localization utilities
- Database integration and ORM-like functionality
- Authentication and authorization system
- API documentation generation
- Mobile platform adapters
- AR/VR platform adapters

## Getting Started

### Installation

```bash
# Install Go (requires Go 1.19 or later)
go install github.com/davidjeba/goscript/cmd/goscript@latest

# Create a new project
goscript new my-project
cd my-project

# Start the development server
go run cmd/server/main.go
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

## Documentation

- [Component System](docs/component-system.md)

## Examples

Check out the examples in the repository:

### Web Components
- `pkg/components/counter.go`: Demonstrates class-based components with state
- `pkg/components/home.go`: Shows how to use the context API and functional components

### CSS Framework
- `cmd/gocsx_demo`: Basic CSS framework demo

### 2D and 3D Applications
- `cmd/gocsx_2d_demo`: 2D canvas application demo
- `cmd/gocsx_3d_demo`: 3D WebGPU application demo

## Features Comparison

### GoScript vs React

- **Language**: GoScript uses Go, React uses JavaScript/TypeScript
- **Performance**: GoScript offers better performance due to Go's efficiency
- **Type Safety**: GoScript has stronger type safety through Go's type system
- **Learning Curve**: Familiar API for React developers, but requires Go knowledge
- **Ecosystem**: React has a larger ecosystem, but GoScript integrates with Go libraries

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

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License
