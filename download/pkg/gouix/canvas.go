package gouix

import (
        "fmt"
        "strings"
)

// CanvasElement represents an element on a canvas
type CanvasElement interface {
        Component
        GetPosition() *Position
        SetPosition(x, y, z float64)
        GetSize() *Size
        SetSize(width, height float64)
        Contains(x, y float64) bool
        Draw(ctx *CanvasContext)
}

// CanvasContext represents a canvas rendering context
type CanvasContext struct {
        Width  int
        Height int
        // In a real implementation, this would contain more rendering state
}

// BaseCanvasElement provides a basic implementation of CanvasElement
type BaseCanvasElement struct {
        BaseComponent
        shape     string // "rect", "circle", "path", etc.
        fillStyle string
        strokeStyle string
        lineWidth float64
}

// NewCanvasElement creates a new canvas element
func NewCanvasElement(id ComponentID, shape string, props Props) *BaseCanvasElement {
        base := NewBaseComponent(id, props)
        
        // Set default position and size
        base.SetPosition(0, 0, 0)
        base.SetSize(100, 100)
        
        // Enable drag by default for canvas elements
        base.EnableDrag(nil)
        
        return &BaseCanvasElement{
                BaseComponent: *base,
                shape:         shape,
                fillStyle:     "#000000",
                strokeStyle:   "#000000",
                lineWidth:     1.0,
        }
}

// Contains checks if a point is inside the element
func (c *BaseCanvasElement) Contains(x, y float64) bool {
        pos := c.GetPosition()
        size := c.GetSize()
        
        if c.shape == "rect" {
                return x >= pos.X && x <= pos.X+size.Width &&
                        y >= pos.Y && y <= pos.Y+size.Height
        } else if c.shape == "circle" {
                centerX := pos.X + size.Width/2
                centerY := pos.Y + size.Height/2
                radius := size.Width / 2
                
                dx := x - centerX
                dy := y - centerY
                
                return dx*dx + dy*dy <= radius*radius
        }
        
        // Default fallback
        return x >= pos.X && x <= pos.X+size.Width &&
                y >= pos.Y && y <= pos.Y+size.Height
}

// Draw draws the element on the canvas
func (c *BaseCanvasElement) Draw(ctx *CanvasContext) {
        // In a real implementation, this would use the canvas API
        // For now, we'll just generate SVG
}

// Render implements the Component interface
func (c *BaseCanvasElement) Render() string {
        pos := c.GetPosition()
        size := c.GetSize()
        
        // Generate SVG based on shape
        var svgContent string
        
        if c.shape == "rect" {
                svgContent = fmt.Sprintf("<rect x=\"%f\" y=\"%f\" width=\"%f\" height=\"%f\" fill=\"%s\" stroke=\"%s\" stroke-width=\"%f\" />",
                        pos.X, pos.Y, size.Width, size.Height, c.fillStyle, c.strokeStyle, c.lineWidth)
        } else if c.shape == "circle" {
                centerX := pos.X + size.Width/2
                centerY := pos.Y + size.Height/2
                radius := size.Width / 2
                
                svgContent = fmt.Sprintf("<circle cx=\"%f\" cy=\"%f\" r=\"%f\" fill=\"%s\" stroke=\"%s\" stroke-width=\"%f\" />",
                        centerX, centerY, radius, c.fillStyle, c.strokeStyle, c.lineWidth)
        } else if c.shape == "text" {
                text, _ := c.GetProps()["text"].(string)
                if text == "" {
                        text = "Text"
                }
                
                svgContent = fmt.Sprintf("<text x=\"%f\" y=\"%f\" fill=\"%s\">%s</text>",
                        pos.X, pos.Y+size.Height/2, c.fillStyle, text)
        }
        
        // Add event handlers
        id := string(c.GetID())
        svgContent = fmt.Sprintf("<g id=\"%s\" data-gouix-id=\"%s\" %s>%s</g>",
                id, id, c.generateEventAttributes(), svgContent)
        
        return svgContent
}

// generateEventAttributes generates event handler attributes
func (c *BaseCanvasElement) generateEventAttributes() string {
        var attrs []string
        
        // Add drag handlers if enabled
        if c.dragConfig.Enabled {
                attrs = append(attrs, "data-gouix-draggable=\"true\"")
        }
        
        // Add touch handlers if enabled
        if c.touchConfig.Enabled {
                attrs = append(attrs, "data-gouix-touch=\"true\"")
        }
        
        return strings.Join(attrs, " ")
}

// Canvas is a container for canvas elements
type Canvas struct {
        BaseComponent
        elements []CanvasElement
        width    int
        height   int
}

// NewCanvas creates a new canvas
func NewCanvas(id ComponentID, width, height int, props Props) *Canvas {
        if props == nil {
                props = Props{}
        }
        
        // Add width and height to props
        props["width"] = width
        props["height"] = height
        
        base := NewBaseComponent(id, props)
        
        return &Canvas{
                BaseComponent: *base,
                elements:      make([]CanvasElement, 0),
                width:         width,
                height:        height,
        }
}

// GetSize returns the size of the canvas
func (c *Canvas) GetSize() *Size {
        return &Size{
                Width:  float64(c.width),
                Height: float64(c.height),
        }
}

// AddElement adds an element to the canvas
func (c *Canvas) AddElement(element CanvasElement) {
        c.elements = append(c.elements, element)
}

// RemoveElement removes an element from the canvas
func (c *Canvas) RemoveElement(id ComponentID) {
        for i, element := range c.elements {
                if element.GetID() == id {
                        c.elements = append(c.elements[:i], c.elements[i+1:]...)
                        break
                }
        }
}

// GetElementAt gets the element at a specific position
func (c *Canvas) GetElementAt(x, y float64) CanvasElement {
        // Check elements in reverse order (top to bottom)
        for i := len(c.elements) - 1; i >= 0; i-- {
                element := c.elements[i]
                if element.Contains(x, y) {
                        return element
                }
        }
        
        return nil
}

// FindElementByID finds an element by ID
func (c *Canvas) FindElementByID(id ComponentID) CanvasElement {
        for _, element := range c.elements {
                if element.GetID() == id {
                        return element
                }
        }
        
        return nil
}

// Render implements the Component interface
func (c *Canvas) Render() string {
        var result strings.Builder
        
        // Start SVG
        result.WriteString(fmt.Sprintf("<svg id=\"%s\" width=\"%d\" height=\"%d\" xmlns=\"http://www.w3.org/2000/svg\">",
                c.GetID(), c.width, c.height))
        
        // Render elements
        for _, element := range c.elements {
                result.WriteString(element.Render())
        }
        
        // End SVG
        result.WriteString("</svg>")
        
        // Add canvas script
        result.WriteString("<script>")
        result.WriteString("if (typeof _gouix === 'undefined') { _gouix = {}; }")
        result.WriteString("_gouix.initCanvas = function(id) {")
        result.WriteString("  // Initialize canvas event handlers")
        result.WriteString("  var svg = document.getElementById(id);")
        result.WriteString("  if (!svg) return;")
        result.WriteString("  ")
        result.WriteString("  // Initialize draggable elements")
        result.WriteString("  var draggables = svg.querySelectorAll('[data-gouix-draggable=\"true\"]');")
        result.WriteString("  draggables.forEach(function(el) {")
        result.WriteString("    el.addEventListener('mousedown', _gouix.dragStart);")
        result.WriteString("  });")
        result.WriteString("  ")
        result.WriteString("  // Initialize touch elements")
        result.WriteString("  var touchables = svg.querySelectorAll('[data-gouix-touch=\"true\"]');")
        result.WriteString("  touchables.forEach(function(el) {")
        result.WriteString("    el.addEventListener('touchstart', _gouix.touchStart);")
        result.WriteString("    el.addEventListener('touchmove', _gouix.touchMove);")
        result.WriteString("    el.addEventListener('touchend', _gouix.touchEnd);")
        result.WriteString("  });")
        result.WriteString("};")
        result.WriteString("</script>")
        
        return result.String()
}

// Rectangle creates a rectangle canvas element
func Rectangle(id ComponentID, x, y, width, height float64, props Props) *BaseCanvasElement {
        if props == nil {
                props = Props{}
        }
        
        element := NewCanvasElement(id, "rect", props)
        element.SetPosition(x, y, 0)
        element.SetSize(width, height)
        
        if fill, ok := props["fill"].(string); ok {
                element.fillStyle = fill
        }
        
        if stroke, ok := props["stroke"].(string); ok {
                element.strokeStyle = stroke
        }
        
        if lineWidth, ok := props["lineWidth"].(float64); ok {
                element.lineWidth = lineWidth
        }
        
        return element
}

// Circle creates a circle canvas element
func Circle(id ComponentID, x, y, radius float64, props Props) *BaseCanvasElement {
        if props == nil {
                props = Props{}
        }
        
        element := NewCanvasElement(id, "circle", props)
        element.SetPosition(x-radius, y-radius, 0)
        element.SetSize(radius*2, radius*2)
        
        if fill, ok := props["fill"].(string); ok {
                element.fillStyle = fill
        }
        
        if stroke, ok := props["stroke"].(string); ok {
                element.strokeStyle = stroke
        }
        
        if lineWidth, ok := props["lineWidth"].(float64); ok {
                element.lineWidth = lineWidth
        }
        
        return element
}

// Text creates a text canvas element
func Text(id ComponentID, x, y float64, text string, props Props) *BaseCanvasElement {
        if props == nil {
                props = Props{}
        }
        
        props["text"] = text
        
        element := NewCanvasElement(id, "text", props)
        element.SetPosition(x, y, 0)
        
        // Default size for text
        width := float64(len(text) * 10)
        height := 20.0
        
        if w, ok := props["width"].(float64); ok {
                width = w
        }
        
        if h, ok := props["height"].(float64); ok {
                height = h
        }
        
        element.SetSize(width, height)
        
        if fill, ok := props["fill"].(string); ok {
                element.fillStyle = fill
        } else {
                element.fillStyle = "#000000"
        }
        
        return element
}