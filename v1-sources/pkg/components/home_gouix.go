package components

import (
        "fmt"
        
        "github.com/davidjeba/goscript/pkg/gouix"
)

// GoUIXHomePage is the home page component
type GoUIXHomePage struct {
        gouix.BaseComponent
        counters      []*GoUIXCounter
        dragCounters  []*DraggableGoUIXCounter
        canvas        *gouix.Canvas
}

// NewGoUIXHomePage creates a new home page
func NewGoUIXHomePage(id gouix.ComponentID, props gouix.Props) *GoUIXHomePage {
        base := gouix.NewBaseComponent(id, props)
        
        // Create canvas
        canvas := gouix.NewCanvas("home-canvas", 800, 400, nil)
        
        // Create home page
        home := &GoUIXHomePage{
                BaseComponent: *base,
                counters:      make([]*GoUIXCounter, 0),
                dragCounters:  make([]*DraggableGoUIXCounter, 0),
                canvas:        canvas,
        }
        
        // Create counters
        counter1 := NewGoUIXCounter("counter-1", gouix.Props{
                "initialCount": 0,
                "title":        "Counter 1",
                "theme":        "light",
        })
        
        counter2 := NewGoUIXCounter("counter-2", gouix.Props{
                "initialCount": 10,
                "title":        "Counter 2",
                "theme":        "dark",
        })
        
        counter3 := NewDraggableGoUIXCounter("counter-3", gouix.Props{
                "initialCount": 5,
                "title":        "Draggable Counter",
                "theme":        "blue",
        })
        
        home.counters = append(home.counters, counter1, counter2)
        home.dragCounters = append(home.dragCounters, counter3)
        
        // Add canvas counter
        CanvasCounter(canvas, "canvas-counter", 50, 50, gouix.Props{
                "initialCount": 20,
                "title":        "Canvas Counter",
        })
        
        // Add event handlers
        home.On("addCounter", home.addCounter)
        
        return home
}

// addCounter adds a new counter
func (h *GoUIXHomePage) addCounter(event gouix.Event) interface{} {
        // Create a new counter with a unique ID
        id := gouix.ComponentID(fmt.Sprintf("counter-%d", len(h.counters)+1))
        
        counter := NewGoUIXCounter(id, gouix.Props{
                "initialCount": 0,
                "title":        fmt.Sprintf("Counter %d", len(h.counters)+1),
                "theme":        "light",
        })
        
        h.counters = append(h.counters, counter)
        
        return nil
}

// Render implements the Component interface
func (h *GoUIXHomePage) Render() string {
        // Create header style
        headerStyle := map[string]interface{}{
                "background-color": "#4a90e2",
                "color":            "#ffffff",
                "padding":          "20px",
                "text-align":       "center",
                "margin-bottom":    "20px",
                "box-shadow":       "0 2px 4px rgba(0, 0, 0, 0.1)",
        }
        
        // Create container style
        containerStyle := map[string]interface{}{
                "max-width":  "1200px",
                "margin":     "0 auto",
                "padding":    "0 20px",
                "font-family": "system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif",
        }
        
        // Create button style
        buttonStyle := map[string]interface{}{
                "background-color": "#4a90e2",
                "color":            "#ffffff",
                "border":           "none",
                "border-radius":    "4px",
                "padding":          "10px 20px",
                "margin":           "20px 0",
                "cursor":           "pointer",
                "font-size":        "16px",
                "transition":       "background-color 0.2s",
        }
        
        // Create section style
        sectionStyle := map[string]interface{}{
                "margin-bottom": "40px",
        }
        
        // Create grid style
        gridStyle := map[string]interface{}{
                "display":               "grid",
                "grid-template-columns": "repeat(auto-fill, minmax(300px, 1fr))",
                "gap":                   "20px",
                "margin-bottom":         "40px",
        }
        
        // Create component ID for event handling
        componentID := string(h.GetID())
        
        // Render counters
        var counterElements []interface{}
        for _, counter := range h.counters {
                counterElements = append(counterElements, counter.Render())
        }
        
        return gouix.CreateElement("div", gouix.Props{
                "class": "home-page",
                "style": containerStyle,
                "id":    componentID,
        },
                // Header
                gouix.CreateElement("header", gouix.Props{
                        "style": headerStyle,
                },
                        gouix.CreateElement("h1", nil, "GoUIX Demo"),
                        gouix.CreateElement("p", nil, "A powerful UI framework for Go"),
                ),
                
                // Main content
                gouix.CreateElement("main", nil,
                        // Regular counters section
                        gouix.CreateElement("section", gouix.Props{
                                "style": sectionStyle,
                        },
                                gouix.CreateElement("h2", nil, "Regular Counters"),
                                gouix.CreateElement("div", gouix.Props{
                                        "class": "counters-grid",
                                        "style": gridStyle,
                                }, counterElements...),
                                gouix.CreateElement("button", gouix.Props{
                                        "style":   buttonStyle,
                                        "onclick": fmt.Sprintf("_gouix.dispatchEvent('%s', 'addCounter', {})", componentID),
                                        "id":      componentID + "-add-counter",
                                }, "Add Counter"),
                        ),
                        
                        // Draggable counter section
                        gouix.CreateElement("section", gouix.Props{
                                "style": sectionStyle,
                        },
                                gouix.CreateElement("h2", nil, "Draggable Counter"),
                                gouix.CreateElement("p", nil, "You can drag this counter around the page:"),
                                gouix.CreateElement("div", gouix.Props{
                                        "class": "draggable-container",
                                        "style": map[string]interface{}{
                                                "position": "relative",
                                                "height":   "200px",
                                                "border":   "1px dashed #ccc",
                                                "margin":   "20px 0",
                                        },
                                }, h.dragCounters[0].Render()),
                        ),
                        
                        // Canvas section
                        gouix.CreateElement("section", gouix.Props{
                                "style": sectionStyle,
                        },
                                gouix.CreateElement("h2", nil, "Canvas Rendering"),
                                gouix.CreateElement("p", nil, "Counters can also be rendered on a canvas:"),
                                h.canvas.Render(),
                        ),
                        
                        // Hooks section
                        gouix.CreateElement("section", gouix.Props{
                                "style": sectionStyle,
                        },
                                gouix.CreateElement("h2", nil, "Functional Components with Hooks"),
                                gouix.CreateElement("p", nil, "Example of a counter using hooks:"),
                                HooksCounter(gouix.Props{
                                        "initialCount": 15,
                                        "title":        "Hooks Demo",
                                        "id":           "hooks-counter",
                                },
                                        gouix.CreateElement("p", gouix.Props{
                                                "style": map[string]interface{}{
                                                        "font-size":  "12px",
                                                        "font-style": "italic",
                                                        "color":      "#666",
                                                },
                                        }, "This is a child element passed to the hooks component"),
                                ),
                        ),
                ),
                
                // Footer
                gouix.CreateElement("footer", gouix.Props{
                        "style": map[string]interface{}{
                                "text-align":    "center",
                                "margin-top":    "40px",
                                "padding":       "20px",
                                "border-top":    "1px solid #eee",
                                "color":         "#666",
                        },
                },
                        gouix.CreateElement("p", nil, "GoUIX - A 100% Go-based UI framework"),
                ),
                
                // Client-side script
                gouix.CreateElement("script", nil, `
                        // Initialize GoUIX
                        if (typeof _gouix === 'undefined') {
                                _gouix = {
                                        components: {},
                                        
                                        // Event handling
                                        dispatchEvent: function(componentId, eventType, data) {
                                                const event = {
                                                        type: eventType,
                                                        target: componentId,
                                                        data: data || {},
                                                        bubbles: true
                                                };
                                                
                                                // Send event to server
                                                console.log('Event:', event);
                                                
                                                // For demo purposes, update the UI directly
                                                if (eventType === 'increment') {
                                                        const countEl = document.getElementById(componentId + '-count');
                                                        if (countEl) {
                                                                countEl.textContent = (parseInt(countEl.textContent) + 1).toString();
                                                        }
                                                } else if (eventType === 'decrement') {
                                                        const countEl = document.getElementById(componentId + '-count');
                                                        if (countEl) {
                                                                countEl.textContent = (parseInt(countEl.textContent) - 1).toString();
                                                        }
                                                } else if (eventType === 'reset') {
                                                        const countEl = document.getElementById(componentId + '-count');
                                                        if (countEl) {
                                                                countEl.textContent = '0';
                                                        }
                                                }
                                        },
                                        
                                        // Drag handling
                                        dragStart: function(event) {
                                                const el = event.target;
                                                el.style.opacity = '0.8';
                                                
                                                // Store initial position
                                                el._startX = event.clientX;
                                                el._startY = event.clientY;
                                                el._initialLeft = parseInt(el.style.left || '0');
                                                el._initialTop = parseInt(el.style.top || '0');
                                                
                                                event.dataTransfer.setData('text/plain', el.id);
                                                event.dataTransfer.effectAllowed = 'move';
                                        },
                                        
                                        drag: function(event) {
                                                // Handled by the browser
                                        },
                                        
                                        dragEnd: function(event) {
                                                const el = event.target;
                                                el.style.opacity = '1';
                                                
                                                // Calculate new position
                                                const dx = event.clientX - el._startX;
                                                const dy = event.clientY - el._startY;
                                                
                                                el.style.left = (el._initialLeft + dx) + 'px';
                                                el.style.top = (el._initialTop + dy) + 'px';
                                        },
                                        
                                        // Touch handling
                                        touchStart: function(event) {
                                                const el = event.target;
                                                const touch = event.touches[0];
                                                
                                                // Store initial position
                                                el._startX = touch.clientX;
                                                el._startY = touch.clientY;
                                                el._initialLeft = parseInt(el.style.left || '0');
                                                el._initialTop = parseInt(el.style.top || '0');
                                        },
                                        
                                        touchMove: function(event) {
                                                const el = event.target;
                                                const touch = event.touches[0];
                                                
                                                // Calculate new position
                                                const dx = touch.clientX - el._startX;
                                                const dy = touch.clientY - el._startY;
                                                
                                                el.style.left = (el._initialLeft + dx) + 'px';
                                                el.style.top = (el._initialTop + dy) + 'px';
                                                
                                                event.preventDefault();
                                        },
                                        
                                        touchEnd: function(event) {
                                                // Touch ended
                                        }
                                };
                        }
                        
                        // Initialize canvas
                        if (document.getElementById('home-canvas')) {
                                _gouix.initCanvas('home-canvas');
                        }
                        
                        // Make draggable elements draggable
                        document.querySelectorAll('.draggable').forEach(function(el) {
                                el.addEventListener('dragstart', _gouix.dragStart);
                                el.addEventListener('drag', _gouix.drag);
                                el.addEventListener('dragend', _gouix.dragEnd);
                                
                                // Also add touch events
                                el.addEventListener('touchstart', _gouix.touchStart);
                                el.addEventListener('touchmove', _gouix.touchMove);
                                el.addEventListener('touchend', _gouix.touchEnd);
                        });
                `),
        )
}