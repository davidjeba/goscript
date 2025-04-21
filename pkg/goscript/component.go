package goscript

import (
        "fmt"
        "reflect"
        "strings"
        "sync"
)

// Props represents component properties
type Props map[string]interface{}

// PropValidator validates a specific prop
type PropValidator func(value interface{}) error

// ReflectKind constants for prop types
const (
        ReflectKindInvalid reflect.Kind = reflect.Invalid
        ReflectKindBool    reflect.Kind = reflect.Bool
        ReflectKindInt     reflect.Kind = reflect.Int
        ReflectKindInt8    reflect.Kind = reflect.Int8
        ReflectKindInt16   reflect.Kind = reflect.Int16
        ReflectKindInt32   reflect.Kind = reflect.Int32
        ReflectKindInt64   reflect.Kind = reflect.Int64
        ReflectKindUint    reflect.Kind = reflect.Uint
        ReflectKindUint8   reflect.Kind = reflect.Uint8
        ReflectKindUint16  reflect.Kind = reflect.Uint16
        ReflectKindUint32  reflect.Kind = reflect.Uint32
        ReflectKindUint64  reflect.Kind = reflect.Uint64
        ReflectKindFloat32 reflect.Kind = reflect.Float32
        ReflectKindFloat64 reflect.Kind = reflect.Float64
        ReflectKindString  reflect.Kind = reflect.String
        ReflectKindSlice   reflect.Kind = reflect.Slice
        ReflectKindMap     reflect.Kind = reflect.Map
)

// PropType defines the expected type and validation for a prop
type PropType struct {
        Type      reflect.Kind
        Required  bool
        Validator PropValidator
        Default   interface{}
}

// PropTypes defines the expected props for a component
type PropTypes map[string]PropType

// Children represents component children
type Children []interface{}

// Component interface defines methods all components must implement
type Component interface {
        Render() string
        GetProps() Props
        GetChildren() Children
}

// LifecycleComponent extends Component with lifecycle methods
type LifecycleComponent interface {
        Component
        ComponentDidMount()
        ComponentWillUnmount()
        ShouldComponentUpdate(nextProps Props) bool
}

// BaseComponent provides a basic implementation of Component
type BaseComponent struct {
        props    Props
        children Children
        state    map[string]interface{}
        propTypes PropTypes
        mutex    sync.RWMutex
}

// NewBaseComponent creates a new BaseComponent
func NewBaseComponent(props Props, propTypes PropTypes, children ...interface{}) *BaseComponent {
        // Apply default values for missing props
        if props == nil {
                props = Props{}
        }
        
        // Apply default values for props
        for key, propType := range propTypes {
                if _, exists := props[key]; !exists && propType.Default != nil {
                        props[key] = propType.Default
                }
        }
        
        return &BaseComponent{
                props:     props,
                children:  children,
                state:     make(map[string]interface{}),
                propTypes: propTypes,
        }
}

// ValidateProps validates component props against propTypes
func (b *BaseComponent) ValidateProps() []error {
        var errors []error
        
        for key, propType := range b.propTypes {
                value, exists := b.props[key]
                
                // Check required props
                if propType.Required && !exists {
                        errors = append(errors, fmt.Errorf("required prop '%s' is missing", key))
                        continue
                }
                
                if exists {
                        // Check type
                        if value != nil && reflect.TypeOf(value).Kind() != propType.Type {
                                errors = append(errors, fmt.Errorf("prop '%s' should be of type %v, got %v", 
                                        key, propType.Type, reflect.TypeOf(value).Kind()))
                        }
                        
                        // Run custom validator if provided
                        if propType.Validator != nil {
                                if err := propType.Validator(value); err != nil {
                                        errors = append(errors, fmt.Errorf("prop '%s' validation failed: %v", key, err))
                                }
                        }
                }
        }
        
        return errors
}

// GetProps returns component props
func (b *BaseComponent) GetProps() Props {
        return b.props
}

// GetChildren returns component children
func (b *BaseComponent) GetChildren() Children {
        return b.children
}

// SetState updates component state
func (b *BaseComponent) SetState(key string, value interface{}) {
        b.mutex.Lock()
        defer b.mutex.Unlock()
        b.state[key] = value
}

// GetState retrieves component state
func (b *BaseComponent) GetState(key string) interface{} {
        b.mutex.RLock()
        defer b.mutex.RUnlock()
        return b.state[key]
}

// Render implements the Component interface
func (b *BaseComponent) Render() string {
        return ""
}

// FunctionalComponent represents a function that renders a component
type FunctionalComponent func(props Props, children ...interface{}) string

// Render implements the Component interface for FunctionalComponent
func (f FunctionalComponent) Render() string {
        return f(nil)
}

// GetProps implements the Component interface for FunctionalComponent
func (f FunctionalComponent) GetProps() Props {
        return nil
}

// GetChildren implements the Component interface for FunctionalComponent
func (f FunctionalComponent) GetChildren() Children {
        return nil
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

// CreateElement is the public version of our element creation function
func CreateElement(component interface{}, props Props, children ...interface{}) string {
        var result strings.Builder

        switch c := component.(type) {
        case string:
                result.WriteString("<")
                result.WriteString(c)
                
                // Handle props/attributes
                if props != nil {
                        for key, value := range props {
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
