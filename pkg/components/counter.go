package components

import (
        "fmt"

        "github.com/davidjeba/goscript/pkg/goscript"
)

// CounterProps defines the props for the Counter component
type CounterProps struct {
        InitialCount int
        Title        string
}

// Counter is a stateful component that demonstrates the new component system
type Counter struct {
        goscript.LifecycleComponentBase
        count int
}

// NewCounter creates a new Counter component
func NewCounter(props goscript.Props) *Counter {
        // Define prop types for validation
        propTypes := goscript.PropTypes{
                "initialCount": goscript.PropType{
                        Type:     goscript.ReflectKindInt,
                        Required: false,
                        Default:  0,
                },
                "title": goscript.PropType{
                        Type:     goscript.ReflectKindString,
                        Required: false,
                        Default:  "Counter",
                },
        }

        // Create base component
        base := goscript.NewBaseComponent(props, propTypes)
        
        // Create counter component
        counter := &Counter{}
        counter.LifecycleComponentBase.BaseComponent = *base
        
        // Initialize count from props
        initialCount, _ := props["initialCount"].(int)
        counter.count = initialCount
        
        return counter
}

// Render implements the Component interface
func (c *Counter) Render() string {
        // Get props
        title, _ := c.GetProps()["title"].(string)
        
        // Create elements
        return goscript.CreateElement("div", goscript.Props{"class": "counter"},
                goscript.CreateElement("h2", nil, title),
                goscript.CreateElement("p", nil, fmt.Sprintf("Count: %d", c.count)),
                goscript.CreateElement("button", 
                        goscript.Props{
                                "onclick": "incrementCounter()",
                        }, 
                        "Increment"),
                goscript.CreateElement("button", 
                        goscript.Props{
                                "onclick": "decrementCounter()",
                        }, 
                        "Decrement"),
        )
}

// ComponentDidMount implements the LifecycleComponent interface
func (c *Counter) ComponentDidMount() {
        fmt.Println("Counter component mounted")
}

// ComponentWillUnmount implements the LifecycleComponent interface
func (c *Counter) ComponentWillUnmount() {
        fmt.Println("Counter component will unmount")
}

// ShouldComponentUpdate implements the LifecycleComponent interface
func (c *Counter) ShouldComponentUpdate(nextProps goscript.Props) bool {
        // Always update for now
        return true
}

// Increment increases the counter
func (c *Counter) Increment() {
        c.count++
}

// Decrement decreases the counter
func (c *Counter) Decrement() {
        c.count--
}

// FunctionalCounter is a functional component version of Counter
func FunctionalCounter(props goscript.Props) string {
        // Get props with defaults
        initialCount := 0
        if val, ok := props["initialCount"].(int); ok {
                initialCount = val
        }
        
        title := "Functional Counter"
        if val, ok := props["title"].(string); ok {
                title = val
        }
        
        // In a real implementation, we would use hooks for state
        // For now, we'll just use the initial count
        count := initialCount
        
        return goscript.CreateElement("div", goscript.Props{"class": "counter"},
                goscript.CreateElement("h2", nil, title),
                goscript.CreateElement("p", nil, fmt.Sprintf("Count: %d", count)),
                goscript.CreateElement("button", 
                        goscript.Props{
                                "onclick": "incrementCounter()",
                        }, 
                        "Increment"),
                goscript.CreateElement("button", 
                        goscript.Props{
                                "onclick": "decrementCounter()",
                        }, 
                        "Decrement"),
        )
}

// CounterWithHooks demonstrates how hooks would be used
// Note: This is a conceptual example, as the hooks system is not fully implemented
// This is commented out as the hooks system is not fully implemented yet
/*
func CounterWithHooks(props goscript.Props, componentID string) string {
        // Use state hook
        countValue, setCount := goscript.useState(componentID, 0)
        count := countValue.(int)
        
        // Use effect hook
        goscript.useEffect(componentID, func() func() {
                fmt.Println("Counter effect running")
                return func() {
                        fmt.Println("Counter effect cleanup")
                }
        }, []interface{}{count})
        
        // Create increment/decrement functions
        increment := func() {
                setCount(count + 1)
        }
        
        decrement := func() {
                setCount(count - 1)
        }
        
        // Get props with defaults
        title := "Hooks Counter"
        if val, ok := props["title"].(string); ok {
                title = val
        }
        
        return goscript.CreateElement("div", goscript.Props{"class": "counter"},
                goscript.CreateElement("h2", nil, title),
                goscript.CreateElement("p", nil, fmt.Sprintf("Count: %d", count)),
                goscript.CreateElement("button", 
                        goscript.Props{
                                "onclick": "increment()",
                                "id": "increment-btn",
                        }, 
                        "Increment"),
                goscript.CreateElement("button", 
                        goscript.Props{
                                "onclick": "decrement()",
                                "id": "decrement-btn",
                        }, 
                        "Decrement"),
        )
}
*/