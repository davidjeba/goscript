# GoScript Component System

The GoScript component system is designed to provide a React-like experience for building web applications in Go. This document outlines the key features and usage patterns of the component system.

## Core Concepts

### Components

Components are the building blocks of a GoScript application. There are two types of components:

1. **Class-based Components**: These implement the `Component` interface and can have lifecycle methods, state, and props.
2. **Functional Components**: These are functions that take props and return HTML.

### Props

Props are properties passed to components. They can be validated using PropTypes.

### State

Components can maintain internal state that can change over time.

### Lifecycle Methods

Class-based components can implement lifecycle methods to run code at specific points in a component's lifecycle:

- `ComponentDidMount`: Called after a component is mounted
- `ComponentWillUnmount`: Called before a component is unmounted
- `ShouldComponentUpdate`: Called before rendering to determine if the component should update

### Context

Context provides a way to pass data through the component tree without having to pass props down manually at every level.

### Hooks

Hooks let you use state and other React features without writing a class.

## Usage Examples

### Creating a Class-based Component

```go
type MyComponent struct {
    goscript.LifecycleComponentBase
    // Component-specific fields
}

func NewMyComponent(props goscript.Props) *MyComponent {
    // Define prop types
    propTypes := goscript.PropTypes{
        "name": goscript.PropType{
            Type:     goscript.ReflectKindString,
            Required: true,
        },
    }
    
    // Create base component
    base := goscript.NewBaseComponent(props, propTypes)
    
    // Create component
    component := &MyComponent{}
    component.LifecycleComponentBase.BaseComponent = *base
    
    return component
}

func (c *MyComponent) Render() string {
    name := c.GetProps()["name"].(string)
    return goscript.CreateElement("div", nil, 
        goscript.CreateElement("h1", nil, "Hello, " + name),
    )
}

func (c *MyComponent) ComponentDidMount() {
    fmt.Println("Component mounted")
}

func (c *MyComponent) ComponentWillUnmount() {
    fmt.Println("Component will unmount")
}

func (c *MyComponent) ShouldComponentUpdate(nextProps goscript.Props) bool {
    return true
}
```

### Creating a Functional Component

```go
func MyFunctionalComponent(props goscript.Props) string {
    name, _ := props["name"].(string)
    return goscript.CreateElement("div", nil, 
        goscript.CreateElement("h1", nil, "Hello, " + name),
    )
}
```

### Using Context

```go
// Create a context
themeContext := goscript.WithContext(nil)
themeContext.Set("theme", "light")

// Create a provider
provider := goscript.CreateProvider(themeContext, "theme", "dark")

// Use the provider
html := provider(nil,
    goscript.CreateElement("div", nil,
        // Create a consumer
        goscript.CreateConsumer(themeContext, "theme", func(value interface{}) string {
            theme := value.(string)
            return goscript.CreateElement("p", nil, "Current theme: " + theme)
        })(nil),
    ),
)
```

### Using Hooks (Conceptual)

```go
func CounterWithHooks(props goscript.Props, componentID string) string {
    // Use state hook
    countValue, setCount := goscript.useState(componentID, 0)
    count := countValue.(int)
    
    // Use effect hook
    goscript.useEffect(componentID, func() func() {
        fmt.Println("Effect running")
        return func() {
            fmt.Println("Effect cleanup")
        }
    }, []interface{}{count})
    
    return goscript.CreateElement("div", nil,
        goscript.CreateElement("p", nil, fmt.Sprintf("Count: %d", count)),
        goscript.CreateElement("button", 
            goscript.Props{"onclick": "increment()"}, 
            "Increment"),
    )
}
```

## Prop Validation

GoScript supports prop validation similar to React's PropTypes:

```go
propTypes := goscript.PropTypes{
    "name": goscript.PropType{
        Type:     goscript.ReflectKindString,
        Required: true,
        Validator: func(value interface{}) error {
            name := value.(string)
            if len(name) < 3 {
                return fmt.Errorf("name must be at least 3 characters")
            }
            return nil
        },
    },
    "age": goscript.PropType{
        Type:     goscript.ReflectKindInt,
        Required: false,
        Default:  18,
    },
}
```

## JSX-like Syntax (Conceptual)

GoScript includes a JSX-like parser that can transform JSX-like syntax into Go code:

```jsx
// JSX-like syntax
<div className="container">
    <h1>Hello, {name}</h1>
    <p>This is a paragraph</p>
</div>
```

Gets transformed into:

```go
goscript.CreateElement("div", goscript.Props{"className": "container"},
    goscript.CreateElement("h1", nil, "Hello, ", name),
    goscript.CreateElement("p", nil, "This is a paragraph"),
)
```

## Testing Components

GoScript components can be easily tested:

```go
func TestMyComponent(t *testing.T) {
    component := NewMyComponent(goscript.Props{
        "name": "John",
    })
    
    html := component.Render()
    
    if !strings.Contains(html, "Hello, John") {
        t.Errorf("Expected 'Hello, John' in HTML, got: %s", html)
    }
}
```