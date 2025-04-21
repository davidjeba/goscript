# GoUIX

GoUIX is a modern, 100% Go-based UI framework designed to provide a powerful, flexible, and reactive approach to building web interfaces. Drawing inspiration from React, SvelteKit, and Flutter, GoUIX offers a unique Go-centric approach to UI development.

![GoUIX Logo](https://via.placeholder.com/800x400?text=GoUIX)

## Features

- **Component-Based Architecture**
  - Class-based components with lifecycle methods
  - Functional components
  - Props validation
  - Fragment support
  - Component testing utilities

- **Reactive State Management**
  - Component-level state
  - Global stores
  - Computed values
  - Effects (side effects)
  - Context API for state sharing

- **Advanced UI Capabilities**
  - Touch events (tap, doubletap, longpress, swipe, pinch, rotate)
  - Drag and drop with precise control
  - Canvas-based rendering with SVG support
  - Position control with x, y, z coordinates
  - Responsive layouts

- **Developer Experience**
  - Server-side rendering (SSR)
  - Client-side hydration
  - Hot module replacement
  - Type safety with Go's type system
  - Familiar API for React/Flutter developers

## Installation

```bash
go get github.com/davidjeba/gouix
```

## Quick Start

### Create a Simple Component

```go
package main

import (
    "fmt"
    "github.com/davidjeba/gouix"
)

// Define a component
type Greeting struct {
    gouix.BaseComponent
}

// Create a new component
func NewGreeting(id gouix.ComponentID, props gouix.Props) *Greeting {
    base := gouix.NewBaseComponent(id, props)
    return &Greeting{
        BaseComponent: *base,
    }
}

// Implement the Render method
func (g *Greeting) Render() string {
    name := "World"
    if val, ok := g.GetProps()["name"].(string); ok {
        name = val
    }
    
    return gouix.CreateElement("div", nil, 
        gouix.CreateElement("h1", nil, fmt.Sprintf("Hello, %s!", name)),
    )
}

func main() {
    greeting := NewGreeting("greeting", gouix.Props{
        "name": "GoUIX",
    })
    
    html := greeting.Render()
    fmt.Println(html)
}
```

### Create a Counter with State

```go
// Define a reactive component
type Counter struct {
    gouix.ReactiveComponent
}

// Create a new counter
func NewCounter(id gouix.ComponentID, props gouix.Props) *Counter {
    // Create initial state
    initialState := map[string]interface{}{
        "count": 0,
    }
    
    // Create reactive component
    base := gouix.NewReactiveComponent(id, props, initialState)
    counter := &Counter{
        ReactiveComponent: *base,
    }
    
    // Add event handlers
    counter.On("increment", func(event gouix.Event) interface{} {
        count := counter.GetState("count").(int)
        counter.SetState("count", count+1)
        return nil
    })
    
    counter.On("decrement", func(event gouix.Event) interface{} {
        count := counter.GetState("count").(int)
        counter.SetState("count", count-1)
        return nil
    })
    
    return counter
}

// Implement the Render method
func (c *Counter) Render() string {
    count := c.GetState("count").(int)
    componentID := string(c.GetID())
    
    return gouix.CreateElement("div", nil, 
        gouix.CreateElement("p", nil, fmt.Sprintf("Count: %d", count)),
        gouix.CreateElement("button", gouix.Props{
            "onclick": fmt.Sprintf("_gouix.dispatchEvent('%s', 'decrement', {})", componentID),
        }, "-"),
        gouix.CreateElement("button", gouix.Props{
            "onclick": fmt.Sprintf("_gouix.dispatchEvent('%s', 'increment', {})", componentID),
        }, "+"),
    )
}
```

### Create a Canvas Application

```go
// Create a canvas
canvas := gouix.NewCanvas("my-canvas", 800, 600, nil)

// Add a rectangle
rect := gouix.Rectangle("rect-1", 50, 50, 200, 100, gouix.Props{
    "fill": "#ff0000",
    "stroke": "#000000",
})
canvas.AddElement(rect)

// Add a circle
circle := gouix.Circle("circle-1", 300, 200, 50, gouix.Props{
    "fill": "#0000ff",
})
canvas.AddElement(circle)

// Enable drag and drop
rect.EnableDrag(&gouix.DragConfig{
    Enabled: true,
    Axis: "both",
})

// Render the canvas
html := canvas.Render()
```

## Documentation

For detailed documentation, see the [GoUIX Documentation](docs/gouix.md).

## Examples

Check out the examples in the `pkg/components` directory:

- `counter_gouix.go`: Demonstrates reactive components with state
- `home_gouix.go`: Shows how to create a complete page with multiple components

## Why GoUIX?

### 100% Go

GoUIX is written entirely in Go, allowing you to build full-stack applications without context switching between languages.

### Modern Component Model

GoUIX provides a modern component model inspired by React and Flutter, but with Go's type safety and performance.

### Built for Touch and DnD

Unlike many frameworks where touch and drag-and-drop are afterthoughts, GoUIX is built from the ground up with these interactions in mind.

### Canvas-First Approach

GoUIX treats canvas rendering as a first-class citizen, allowing for pixel-perfect control and complex visualizations.

### Reactive by Default

State management is reactive by default, automatically updating the UI when state changes.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License