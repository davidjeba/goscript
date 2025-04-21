package gouix

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// ComponentID is a unique identifier for a component
type ComponentID string

// Props represents component properties
type Props map[string]interface{}

// EventHandler is a function that handles events
type EventHandler func(event Event) interface{}

// Event represents a UI event
type Event struct {
	Type    string
	Target  ComponentID
	Data    map[string]interface{}
	Bubbles bool
}

// Position represents x, y, z coordinates for canvas elements
type Position struct {
	X float64
	Y float64
	Z float64
}

// Size represents width and height
type Size struct {
	Width  float64
	Height float64
}

// DragConfig configures drag behavior
type DragConfig struct {
	Enabled       bool
	Axis          string // "x", "y", "both"
	Bounds        *Rect
	SnapToGrid    bool
	GridSize      float64
	OnDragStart   EventHandler
	OnDragMove    EventHandler
	OnDragEnd     EventHandler
	DragThreshold float64
}

// TouchConfig configures touch behavior
type TouchConfig struct {
	Enabled      bool
	Gestures     []string // "tap", "doubletap", "longpress", "swipe", "pinch", "rotate"
	OnTouchStart EventHandler
	OnTouchMove  EventHandler
	OnTouchEnd   EventHandler
}

// Rect represents a rectangle
type Rect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

// Component is the base interface for all components
type Component interface {
	// Core methods
	Render() string
	GetID() ComponentID
	GetProps() Props
	GetChildren() []interface{}
	
	// Event handling
	HandleEvent(event Event) interface{}
	
	// Lifecycle methods
	Mount()
	Unmount()
	Update(nextProps Props) bool
}

// BaseComponent provides a basic implementation of Component
type BaseComponent struct {
	id          ComponentID
	props       Props
	children    []interface{}
	state       map[string]interface{}
	events      map[string][]EventHandler
	position    *Position
	size        *Size
	dragConfig  *DragConfig
	touchConfig *TouchConfig
	mutex       sync.RWMutex
}

// NewBaseComponent creates a new BaseComponent
func NewBaseComponent(id ComponentID, props Props, children ...interface{}) *BaseComponent {
	if props == nil {
		props = Props{}
	}
	
	// Default drag configuration (disabled by default)
	dragConfig := &DragConfig{
		Enabled:       false,
		Axis:          "both",
		SnapToGrid:    false,
		GridSize:      10,
		DragThreshold: 5,
	}
	
	// Default touch configuration (enabled by default)
	touchConfig := &TouchConfig{
		Enabled:  true,
		Gestures: []string{"tap", "doubletap", "longpress", "swipe"},
	}
	
	return &BaseComponent{
		id:          id,
		props:       props,
		children:    children,
		state:       make(map[string]interface{}),
		events:      make(map[string][]EventHandler),
		dragConfig:  dragConfig,
		touchConfig: touchConfig,
	}
}

// GetID returns the component ID
func (b *BaseComponent) GetID() ComponentID {
	return b.id
}

// GetProps returns component props
func (b *BaseComponent) GetProps() Props {
	return b.props
}

// GetChildren returns component children
func (b *BaseComponent) GetChildren() []interface{} {
	return b.children
}

// SetState updates component state
func (b *BaseComponent) SetState(key string, value interface{}) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	
	oldValue := b.state[key]
	b.state[key] = value
	
	// Trigger re-render if value changed
	if !reflect.DeepEqual(oldValue, value) {
		// In a real implementation, this would trigger a re-render
	}
}

// GetState retrieves component state
func (b *BaseComponent) GetState(key string) interface{} {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.state[key]
}

// On registers an event handler
func (b *BaseComponent) On(eventType string, handler EventHandler) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	
	b.events[eventType] = append(b.events[eventType], handler)
}

// HandleEvent handles an event
func (b *BaseComponent) HandleEvent(event Event) interface{} {
	b.mutex.RLock()
	handlers, exists := b.events[event.Type]
	b.mutex.RUnlock()
	
	if !exists {
		return nil
	}
	
	var result interface{}
	for _, handler := range handlers {
		result = handler(event)
	}
	
	return result
}

// EnableDrag enables drag functionality
func (b *BaseComponent) EnableDrag(config *DragConfig) {
	if config == nil {
		config = &DragConfig{
			Enabled: true,
			Axis:    "both",
		}
	}
	
	b.mutex.Lock()
	defer b.mutex.Unlock()
	
	b.dragConfig = config
	b.dragConfig.Enabled = true
}

// DisableDrag disables drag functionality
func (b *BaseComponent) DisableDrag() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	
	b.dragConfig.Enabled = false
}

// SetPosition sets the component position
func (b *BaseComponent) SetPosition(x, y, z float64) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	
	if b.position == nil {
		b.position = &Position{}
	}
	
	b.position.X = x
	b.position.Y = y
	b.position.Z = z
}

// GetPosition gets the component position
func (b *BaseComponent) GetPosition() *Position {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	
	return b.position
}

// SetSize sets the component size
func (b *BaseComponent) SetSize(width, height float64) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	
	if b.size == nil {
		b.size = &Size{}
	}
	
	b.size.Width = width
	b.size.Height = height
}

// GetSize gets the component size
func (b *BaseComponent) GetSize() *Size {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	
	return b.size
}

// Render implements the Component interface
func (b *BaseComponent) Render() string {
	return ""
}

// Mount implements the Component interface
func (b *BaseComponent) Mount() {
	// Default implementation does nothing
}

// Unmount implements the Component interface
func (b *BaseComponent) Unmount() {
	// Default implementation does nothing
}

// Update implements the Component interface
func (b *BaseComponent) Update(nextProps Props) bool {
	// Default implementation always updates
	return true
}

// FunctionalComponent represents a function that renders a component
type FunctionalComponent func(props Props, children ...interface{}) string

// Render implements the Component interface for FunctionalComponent
func (f FunctionalComponent) Render() string {
	return f(nil)
}

// GetID implements the Component interface for FunctionalComponent
func (f FunctionalComponent) GetID() ComponentID {
	return ""
}

// GetProps implements the Component interface for FunctionalComponent
func (f FunctionalComponent) GetProps() Props {
	return nil
}

// GetChildren implements the Component interface for FunctionalComponent
func (f FunctionalComponent) GetChildren() []interface{} {
	return nil
}

// HandleEvent implements the Component interface for FunctionalComponent
func (f FunctionalComponent) HandleEvent(event Event) interface{} {
	return nil
}

// Mount implements the Component interface for FunctionalComponent
func (f FunctionalComponent) Mount() {
	// No-op for functional components
}

// Unmount implements the Component interface for FunctionalComponent
func (f FunctionalComponent) Unmount() {
	// No-op for functional components
}

// Update implements the Component interface for FunctionalComponent
func (f FunctionalComponent) Update(nextProps Props) bool {
	return true
}

// renderChild renders a single child
func renderChild(child interface{}) string {
	switch ch := child.(type) {
	case Component:
		return ch.Render()
	case FunctionalComponent:
		return ch(nil)
	case string:
		return ch
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		return fmt.Sprintf("%v", ch)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", ch)
	}
}

// CreateElement creates an HTML element
func CreateElement(component interface{}, props Props, children ...interface{}) string {
	var result strings.Builder

	switch c := component.(type) {
	case string:
		result.WriteString("<")
		result.WriteString(c)
		
		// Handle props/attributes
		if props != nil {
			for key, value := range props {
				// Special handling for event handlers
				if strings.HasPrefix(key, "on") && strings.HasPrefix(fmt.Sprintf("%T", value), "func(") {
					// In a real implementation, this would register event handlers
					result.WriteString(fmt.Sprintf(" %s=\"_gouix.handleEvent('%s', event)\"", key, key))
					continue
				}
				
				// Special handling for style
				if key == "style" && reflect.TypeOf(value).Kind() == reflect.Map {
					styleMap := value.(map[string]interface{})
					var styleStr strings.Builder
					for sk, sv := range styleMap {
						styleStr.WriteString(fmt.Sprintf("%s:%v;", sk, sv))
					}
					result.WriteString(fmt.Sprintf(" style=\"%s\"", styleStr.String()))
					continue
				}
				
				// Handle boolean attributes
				if val, ok := value.(bool); ok && val {
					result.WriteString(fmt.Sprintf(" %s", key))
				} else if val, ok := value.(bool); ok && !val {
					// Skip false boolean attributes
					continue
				} else if value != nil {
					// Handle regular attributes with values
					result.WriteString(fmt.Sprintf(" %s=\"%v\"", key, value))
				}
			}
			
			// Add drag and touch attributes if needed
			if dragEnabled, ok := props["draggable"].(bool); ok && dragEnabled {
				result.WriteString(" draggable=\"true\"")
				result.WriteString(" ondragstart=\"_gouix.dragStart(event)\"")
				result.WriteString(" ondrag=\"_gouix.drag(event)\"")
				result.WriteString(" ondragend=\"_gouix.dragEnd(event)\"")
			}
			
			// Add position and size if provided
			if position, ok := props["position"].(*Position); ok {
				result.WriteString(fmt.Sprintf(" style=\"position:absolute;left:%fpx;top:%fpx;z-index:%f\"", 
					position.X, position.Y, position.Z))
			}
		}
		
		// Handle self-closing tags
		if len(children) == 0 {
			// Self-closing tags in HTML5
			selfClosingTags := map[string]bool{
				"area": true, "base": true, "br": true, "col": true, "embed": true,
				"hr": true, "img": true, "input": true, "link": true, "meta": true,
				"param": true, "source": true, "track": true, "wbr": true,
			}
			
			if _, isSelfClosing := selfClosingTags[c]; isSelfClosing {
				result.WriteString("/>")
			} else {
				result.WriteString("></")
				result.WriteString(c)
				result.WriteString(">")
			}
		} else {
			result.WriteString(">")
			
			// Render children
			for _, child := range children {
				// Handle arrays/slices of children
				if reflect.TypeOf(child).Kind() == reflect.Slice {
					s := reflect.ValueOf(child)
					for i := 0; i < s.Len(); i++ {
						childItem := s.Index(i).Interface()
						result.WriteString(renderChild(childItem))
					}
				} else {
					result.WriteString(renderChild(child))
				}
			}
			
			result.WriteString("</")
			result.WriteString(c)
			result.WriteString(">")
		}
	case Component:
		result.WriteString(c.Render())
	case FunctionalComponent:
		result.WriteString(c(props, children...))
	case nil:
		// Handle nil component gracefully
		result.WriteString("<!-- nil component -->")
	default:
		// Handle unexpected component types
		result.WriteString(fmt.Sprintf("<!-- unknown component type: %T -->", c))
	}

	return result.String()
}

// Fragment is a special component that renders only its children
func Fragment(props Props, children ...interface{}) string {
	var result strings.Builder
	for _, child := range children {
		result.WriteString(renderChild(child))
	}
	return result.String()
}