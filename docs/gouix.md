# GoUIX: A Modern UI Framework for Go

GoUIX is a powerful, 100% Go-based UI framework designed to provide a modern, reactive, and flexible approach to building web interfaces. Drawing inspiration from the best parts of React, SvelteKit, and Flutter, GoUIX offers a unique Go-centric approach to UI development.

## Core Features

- **Component-Based Architecture**: Build UIs with reusable, composable components
- **Reactive State Management**: Automatic UI updates when state changes
- **Touch and Drag-and-Drop Support**: Built-in support for touch events and DnD
- **Canvas Rendering**: Create interactive canvas-based UIs with precise control
- **Server-Side Rendering**: Fast initial page loads with hydration
- **Type Safety**: Leverage Go's type system for safer UI code

## Component System

GoUIX provides multiple ways to create components:

### Class-Based Components

```go
// Define a component
type MyComponent struct {
    gouix.BaseComponent
}

// Create a new component
func NewMyComponent(id gouix.ComponentID, props gouix.Props) *MyComponent {
    base := gouix.NewBaseComponent(id, props)
    return &MyComponent{
        BaseComponent: *base,
    }
}

// Implement the Render method
func (c *MyComponent) Render() string {
    return gouix.CreateElement("div", nil, 
        gouix.CreateElement("h1", nil, "Hello, World!"),
    )
}
```

### Reactive Components

```go
// Define a reactive component
type Counter struct {
    gouix.ReactiveComponent
}

// Create a new reactive component
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
    
    return counter
}

// Implement the Render method
func (c *Counter) Render() string {
    count := c.GetState("count").(int)
    
    return gouix.CreateElement("div", nil, 
        gouix.CreateElement("p", nil, fmt.Sprintf("Count: %d", count)),
        gouix.CreateElement("button", gouix.Props{
            "onclick": fmt.Sprintf("_gouix.dispatchEvent('%s', 'increment', {})", c.GetID()),
        }, "Increment"),
    )
}
```

### Functional Components

```go
// Define a functional component
func Greeting(props gouix.Props, children ...interface{}) string {
    name := "World"
    if val, ok := props["name"].(string); ok {
        name = val
    }
    
    return gouix.CreateElement("div", nil, 
        gouix.CreateElement("h1", nil, fmt.Sprintf("Hello, %s!", name)),
        gouix.Fragment(nil, children...),
    )
}

// Use the functional component
Greeting(gouix.Props{"name": "John"}, 
    gouix.CreateElement("p", nil, "This is a child element"),
)
```

## Canvas Rendering

GoUIX provides a powerful canvas rendering system for creating interactive graphics:

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

// Add text
text := gouix.Text("text-1", 400, 300, "Hello, Canvas!", gouix.Props{
    "fill": "#000000",
})
canvas.AddElement(text)

// Render the canvas
canvas.Render()
```

## Drag and Drop

GoUIX components have built-in support for drag and drop:

```go
// Enable drag and drop on a component
component.EnableDrag(&gouix.DragConfig{
    Enabled: true,
    Axis: "both", // "x", "y", or "both"
    SnapToGrid: true,
    GridSize: 10,
    OnDragStart: func(event gouix.Event) interface{} {
        fmt.Println("Drag started")
        return nil
    },
    OnDragEnd: func(event gouix.Event) interface{} {
        fmt.Println("Drag ended")
        return nil
    },
})
```

## Touch Support

GoUIX components have built-in support for touch events:

```go
// Configure touch support
component.touchConfig = &gouix.TouchConfig{
    Enabled: true,
    Gestures: []string{"tap", "doubletap", "longpress", "swipe"},
    OnTouchStart: func(event gouix.Event) interface{} {
        fmt.Println("Touch started")
        return nil
    },
    OnTouchMove: func(event gouix.Event) interface{} {
        fmt.Println("Touch moved")
        return nil
    },
    OnTouchEnd: func(event gouix.Event) interface{} {
        fmt.Println("Touch ended")
        return nil
    },
}
```

## Reactive State Management

GoUIX provides a powerful reactive state management system:

```go
// Create a store
store := gouix.NewStore(map[string]interface{}{
    "count": 0,
    "user": map[string]interface{}{
        "name": "John",
        "age": 30,
    },
})

// Get a value
count := store.GetValue("count").(int)

// Set a value
store.Set("count", count+1)

// Subscribe to changes
unsubscribe := store.Subscribe("count", func(newValue, oldValue interface{}) {
    fmt.Printf("Count changed from %v to %v\n", oldValue, newValue)
})

// Create a computed value
store.AddComputed("doubleCount", func() interface{} {
    return store.GetValue("count").(int) * 2
}, "count")

// Get a computed value
doubleCount := store.GetComputed("doubleCount").Get().(int)
```

## Event Handling

GoUIX provides a flexible event handling system:

```go
// Register an event handler
component.On("click", func(event gouix.Event) interface{} {
    fmt.Println("Clicked!")
    return nil
})

// Dispatch an event
component.HandleEvent(gouix.Event{
    Type: "click",
    Target: component.GetID(),
    Data: map[string]interface{}{
        "x": 100,
        "y": 200,
    },
    Bubbles: true,
})
```

## Styling Components

GoUIX provides multiple ways to style components:

```go
// Inline styles
gouix.CreateElement("div", gouix.Props{
    "style": map[string]interface{}{
        "background-color": "#f0f0f0",
        "padding": "16px",
        "border-radius": "8px",
    },
}, "Styled div")

// CSS classes
gouix.CreateElement("div", gouix.Props{
    "class": "my-component primary large",
}, "Styled with classes")

// Conditional styling
isActive := true
gouix.CreateElement("div", gouix.Props{
    "class": fmt.Sprintf("button %s", map[bool]string{true: "active", false: ""}[isActive]),
}, "Conditional styling")
```

## Best Practices

1. **Component Organization**: Group related components in packages
2. **State Management**: Keep state as close as possible to where it's used
3. **Props Validation**: Validate props to catch errors early
4. **Event Handling**: Use event delegation for better performance
5. **Rendering Optimization**: Minimize DOM updates by using computed values
6. **Accessibility**: Ensure components are accessible with proper ARIA attributes
7. **Testing**: Write tests for components to ensure they render correctly

## Performance Considerations

1. **Minimize Renders**: Use computed values to avoid unnecessary renders
2. **Efficient DOM Updates**: Update only what has changed
3. **Lazy Loading**: Load components only when needed
4. **Memoization**: Cache expensive computations
5. **Event Delegation**: Use event delegation for better performance
6. **Virtualization**: Use virtualization for long lists

## Comparison with Other Frameworks

### GoUIX vs React

- **Language**: GoUIX uses 100% Go, React uses JavaScript/TypeScript
- **Rendering**: Both support SSR and client-side rendering
- **State Management**: GoUIX has built-in reactive state, React uses hooks/context
- **Performance**: GoUIX optimizes for server rendering, React for client rendering
- **Learning Curve**: GoUIX is simpler for Go developers

### GoUIX vs SvelteKit

- **Language**: GoUIX uses Go, SvelteKit uses JavaScript/TypeScript
- **Compilation**: Both compile components, but in different ways
- **Reactivity**: Both have built-in reactivity systems
- **Bundle Size**: GoUIX has no client-side runtime, SvelteKit has minimal runtime
- **Server Integration**: GoUIX has tighter integration with Go servers

### GoUIX vs Flutter

- **Language**: GoUIX uses Go, Flutter uses Dart
- **Platform**: GoUIX targets web, Flutter targets multiple platforms
- **Rendering**: GoUIX uses HTML/CSS/SVG, Flutter uses custom rendering
- **Component Model**: Both use a component-based architecture
- **State Management**: Both have reactive state management

## Conclusion

GoUIX provides a powerful, flexible, and Go-centric approach to building modern web UIs. By combining the best ideas from React, SvelteKit, and Flutter with Go's strengths, GoUIX offers a unique and productive development experience.