package components

import (
        "fmt"
        "strconv"
        "strings"

        "github.com/davidjeba/goscript/pkg/gouix"
)

// GoUIXCounterProps defines the props for the GoUIXCounter component
type GoUIXCounterProps struct {
        InitialCount int
        Title        string
        Theme        string
}

// GoUIXCounter is a hyper(reactive) counter component
type GoUIXCounter struct {
        gouix.HyperComponent
}

// NewGoUIXCounter creates a new counter component
func NewGoUIXCounter(id gouix.ComponentID, props gouix.Props) *GoUIXCounter {
        // Extract and validate props
        initialCount := 0
        if val, ok := props["initialCount"].(int); ok {
                initialCount = val
        }
        
        title := "Counter"
        if val, ok := props["title"].(string); ok {
                title = val
        }
        
        // Create initial state
        initialState := map[string]interface{}{
                "count": initialCount,
                "title": title,
        }
        
        // Create hyper(reactive) component
        base := gouix.NewHyperComponent(id, props, initialState)
        counter := &GoUIXCounter{
                HyperComponent: *base,
        }
        
        // Add event handlers
        counter.On("increment", counter.increment)
        counter.On("decrement", counter.decrement)
        counter.On("reset", counter.reset)
        
        return counter
}

// increment increases the counter
func (c *GoUIXCounter) increment(event gouix.Event) interface{} {
        count := c.GetState("count").(int)
        c.SetState("count", count+1)
        return nil
}

// decrement decreases the counter
func (c *GoUIXCounter) decrement(event gouix.Event) interface{} {
        count := c.GetState("count").(int)
        c.SetState("count", count-1)
        return nil
}

// reset resets the counter
func (c *GoUIXCounter) reset(event gouix.Event) interface{} {
        initialCount := 0
        if val, ok := c.GetProps()["initialCount"].(int); ok {
                initialCount = val
        }
        
        c.SetState("count", initialCount)
        return nil
}

// Render implements the Component interface
func (c *GoUIXCounter) Render() string {
        count := c.GetState("count").(int)
        title := c.GetState("title").(string)
        
        // Get theme from props
        theme := "light"
        if val, ok := c.GetProps()["theme"].(string); ok {
                theme = val
        }
        
        // Define theme styles
        themeStyles := map[string]map[string]string{
                "light": {
                        "bg":       "#ffffff",
                        "text":     "#333333",
                        "border":   "#dddddd",
                        "btnBg":    "#f0f0f0",
                        "btnText":  "#333333",
                        "btnHover": "#e0e0e0",
                },
                "dark": {
                        "bg":       "#333333",
                        "text":     "#ffffff",
                        "border":   "#555555",
                        "btnBg":    "#444444",
                        "btnText":  "#ffffff",
                        "btnHover": "#555555",
                },
                "blue": {
                        "bg":       "#e6f7ff",
                        "text":     "#0066cc",
                        "border":   "#99ccff",
                        "btnBg":    "#0066cc",
                        "btnText":  "#ffffff",
                        "btnHover": "#0052a3",
                },
        }
        
        // Use theme or default to light
        styles, ok := themeStyles[theme]
        if !ok {
                styles = themeStyles["light"]
        }
        
        // Create container style
        containerStyle := map[string]interface{}{
                "background-color": styles["bg"],
                "color":            styles["text"],
                "border":           "1px solid " + styles["border"],
                "border-radius":    "8px",
                "padding":          "16px",
                "text-align":       "center",
                "width":            "300px",
                "margin":           "20px auto",
                "box-shadow":       "0 4px 6px rgba(0, 0, 0, 0.1)",
        }
        
        // Create button style
        buttonStyle := map[string]interface{}{
                "background-color": styles["btnBg"],
                "color":            styles["btnText"],
                "border":           "none",
                "border-radius":    "4px",
                "padding":          "8px 16px",
                "margin":           "0 8px",
                "cursor":           "pointer",
                "font-size":        "14px",
                "transition":       "background-color 0.2s",
        }
        
        // Create count style
        countStyle := map[string]interface{}{
                "font-size":   "48px",
                "font-weight": "bold",
                "margin":      "16px 0",
        }
        
        // Create title style
        titleStyle := map[string]interface{}{
                "font-size":   "24px",
                "font-weight": "bold",
                "margin":      "0 0 16px 0",
        }
        
        // Create component ID for event handling
        componentID := string(c.GetID())
        
        return gouix.CreateElement("div", gouix.Props{
                "class": "counter",
                "style": containerStyle,
                "id":    componentID,
        },
                gouix.CreateElement("h2", gouix.Props{
                        "style": titleStyle,
                }, title),
                
                gouix.CreateElement("p", gouix.Props{
                        "style": countStyle,
                        "id":    componentID + "-count",
                }, strconv.Itoa(count)),
                
                gouix.CreateElement("div", gouix.Props{
                        "class": "counter-buttons",
                },
                        gouix.CreateElement("button", gouix.Props{
                                "style":   buttonStyle,
                                "onclick": fmt.Sprintf("_gouix.dispatchEvent('%s', 'decrement', {})", componentID),
                                "id":      componentID + "-decrement",
                        }, "−"),
                        
                        gouix.CreateElement("button", gouix.Props{
                                "style":   buttonStyle,
                                "onclick": fmt.Sprintf("_gouix.dispatchEvent('%s', 'reset', {})", componentID),
                                "id":      componentID + "-reset",
                        }, "Reset"),
                        
                        gouix.CreateElement("button", gouix.Props{
                                "style":   buttonStyle,
                                "onclick": fmt.Sprintf("_gouix.dispatchEvent('%s', 'increment', {})", componentID),
                                "id":      componentID + "-increment",
                        }, "+"),
                ),
        )
}

// DraggableGoUIXCounter is a counter that can be dragged
type DraggableGoUIXCounter struct {
        GoUIXCounter
}

// NewDraggableGoUIXCounter creates a new draggable counter
func NewDraggableGoUIXCounter(id gouix.ComponentID, props gouix.Props) *DraggableGoUIXCounter {
        counter := NewGoUIXCounter(id, props)
        
        // Create draggable counter
        draggable := &DraggableGoUIXCounter{
                GoUIXCounter: *counter,
        }
        
        // Enable drag
        draggable.EnableDrag(&gouix.DragConfig{
                Enabled: true,
                Axis:    "both",
        })
        
        return draggable
}

// Render implements the Component interface
func (c *DraggableGoUIXCounter) Render() string {
        // Get base rendering
        baseHTML := c.GoUIXCounter.Render()
        
        // Add draggable attributes
        return strings.Replace(baseHTML, "class=\"counter\"", "class=\"counter draggable\" draggable=\"true\"", 1)
}

// CounterWithHooks is a functional component that uses hooks
func CounterWithHooks(props gouix.Props, children ...interface{}) string {
        // Extract props
        initialCount := 0
        if val, ok := props["initialCount"].(int); ok {
                initialCount = val
        }
        
        title := "Counter with Hooks"
        if val, ok := props["title"].(string); ok {
                title = val
        }
        
        id := "counter-hooks"
        if val, ok := props["id"].(string); ok {
                id = val
        }
        
        // Use state hook
        count := gouix.UseState(initialCount)
        
        // Create container style
        containerStyle := map[string]interface{}{
                "background-color": "#f8f9fa",
                "color":            "#333333",
                "border":           "1px solid #dee2e6",
                "border-radius":    "8px",
                "padding":          "16px",
                "text-align":       "center",
                "width":            "300px",
                "margin":           "20px auto",
                "box-shadow":       "0 4px 6px rgba(0, 0, 0, 0.1)",
        }
        
        // Create button style
        buttonStyle := map[string]interface{}{
                "background-color": "#6c757d",
                "color":            "#ffffff",
                "border":           "none",
                "border-radius":    "4px",
                "padding":          "8px 16px",
                "margin":           "0 8px",
                "cursor":           "pointer",
                "font-size":        "14px",
                "transition":       "background-color 0.2s",
        }
        
        // Create counter element
        return gouix.CreateElement("div", gouix.Props{
                "class": "counter-hooks",
                "style": containerStyle,
                "id":    id,
        },
                // Title
                gouix.CreateElement("h3", nil, title),
                
                // Count display
                gouix.CreateElement("div", gouix.Props{
                        "style": map[string]interface{}{
                                "font-size":   "24px",
                                "font-weight": "bold",
                                "margin":      "16px 0",
                        },
                }, fmt.Sprintf("Count: %d", count.Get())),
                
                // Children
                gouix.CreateElement("div", nil, children...),
                
                // Buttons
                gouix.CreateElement("div", gouix.Props{
                        "style": map[string]interface{}{
                                "display":         "flex",
                                "justify-content": "center",
                                "margin-top":      "16px",
                        },
                },
                        // Decrement button
                        gouix.CreateElement("button", gouix.Props{
                                "style":   buttonStyle,
                                "onclick": fmt.Sprintf("_gouix.dispatchEvent('%s', 'decrement', {})", id),
                                "id":      id + "-decrement",
                        }, "-"),
                        
                        // Reset button
                        gouix.CreateElement("button", gouix.Props{
                                "style": map[string]interface{}{
                                        "background-color": "#dc3545",
                                        "color":            "#ffffff",
                                        "border":           "none",
                                        "border-radius":    "4px",
                                        "padding":          "8px 16px",
                                        "margin":           "0 8px",
                                        "cursor":           "pointer",
                                        "font-size":        "14px",
                                        "transition":       "background-color 0.2s",
                                },
                                "onclick": fmt.Sprintf("_gouix.dispatchEvent('%s', 'reset', {})", id),
                                "id":      id + "-reset",
                        }, "Reset"),
                        
                        // Increment button
                        gouix.CreateElement("button", gouix.Props{
                                "style":   buttonStyle,
                                "onclick": fmt.Sprintf("_gouix.dispatchEvent('%s', 'increment', {})", id),
                                "id":      id + "-increment",
                        }, "+"),
                ),
        )
}

// CanvasCounter is a counter rendered on a canvas
func CanvasCounter(canvas *gouix.Canvas, id gouix.ComponentID, x, y float64, props gouix.Props) {
        // Extract props
        initialCount := 0
        if val, ok := props["initialCount"].(int); ok {
                initialCount = val
        }
        
        title := "Canvas Counter"
        if val, ok := props["title"].(string); ok {
                title = val
        }
        
        // Create background rectangle
        bg := gouix.Rectangle(id+"-bg", x, y, 200, 150, gouix.Props{
                "fill":   "#ffffff",
                "stroke": "#dddddd",
        })
        canvas.AddElement(bg)
        
        // Create title text
        titleText := gouix.Text(id+"-title", x+10, y+20, title, gouix.Props{
                "fill": "#333333",
        })
        canvas.AddElement(titleText)
        
        // Create count text
        countText := gouix.Text(id+"-count", x+100, y+70, strconv.Itoa(initialCount), gouix.Props{
                "fill": "#333333",
        })
        canvas.AddElement(countText)
        
        // Create increment button
        incButton := gouix.Rectangle(id+"-inc", x+120, y+100, 60, 30, gouix.Props{
                "fill":   "#f0f0f0",
                "stroke": "#dddddd",
        })
        canvas.AddElement(incButton)
        
        // Create increment text
        incText := gouix.Text(id+"-inc-text", x+140, y+115, "+", gouix.Props{
                "fill": "#333333",
        })
        canvas.AddElement(incText)
        
        // Create decrement button
        decButton := gouix.Rectangle(id+"-dec", x+20, y+100, 60, 30, gouix.Props{
                "fill":   "#f0f0f0",
                "stroke": "#dddddd",
        })
        canvas.AddElement(decButton)
        
        // Create decrement text
        decText := gouix.Text(id+"-dec-text", x+40, y+115, "−", gouix.Props{
                "fill": "#333333",
        })
        canvas.AddElement(decText)
}

// HooksCounter is a functional component using hooks
func HooksCounter(props gouix.Props, children ...interface{}) string {
        // Extract props
        initialCount := 0
        if val, ok := props["initialCount"].(int); ok {
                initialCount = val
        }
        
        title := "Hooks Counter"
        if val, ok := props["title"].(string); ok {
                title = val
        }
        
        // Create a unique ID for this component instance
        id := "counter-hooks"
        if val, ok := props["id"].(string); ok {
                id = val
        }
        
        // In a real implementation, we would use hooks here
        // For now, we'll just render a static counter
        
        return gouix.CreateElement("div", gouix.Props{
                "class": "counter",
                "id":    id,
        },
                gouix.CreateElement("h2", nil, title),
                gouix.CreateElement("p", gouix.Props{
                        "id": id + "-count",
                }, strconv.Itoa(initialCount)),
                gouix.CreateElement("div", gouix.Props{
                        "class": "counter-buttons",
                },
                        gouix.CreateElement("button", gouix.Props{
                                "id": id + "-decrement",
                        }, "−"),
                        gouix.CreateElement("button", gouix.Props{
                                "id": id + "-reset",
                        }, "Reset"),
                        gouix.CreateElement("button", gouix.Props{
                                "id": id + "-increment",
                        }, "+"),
                ),
                gouix.Fragment(nil, children...),
        )
}