package goscript

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// Props represents component properties as a key-value map
type Props map[string]interface{}

// Children is a slice of child Components
type Children []Component

// Component is the core interface that all GoScript components must implement
type Component interface {
	Render() string
}

// FunctionalComponent is a function that renders to HTML
type FunctionalComponent func(Props) string

// Render implements Component for FunctionalComponent
func (fc FunctionalComponent) Render() string {
	return fc(nil)
}

// RouteHandler is the function signature for page route handlers
type RouteHandler func(w http.ResponseWriter, r *http.Request, params map[string]string)

// BaseComponent provides a base implementation with props, children, and state
type BaseComponent struct {
	props   Props
	children Children
	state   map[string]interface{}
	mu      sync.RWMutex
}

// NewBaseComponent creates a new BaseComponent with the given props and children
func NewBaseComponent(props Props, children Children) *BaseComponent {
	if props == nil {
		props = make(Props)
	}
	return &BaseComponent{
		props:    props,
		children: children,
		state:    make(map[string]interface{}),
	}
}

// ValidateProps checks that required props are present
func (bc *BaseComponent) ValidateProps(required []string) error {
	for _, key := range required {
		if _, ok := bc.props[key]; !ok {
			return fmt.Errorf("missing required prop: %s", key)
		}
	}
	return nil
}

// GetProps returns the component's props
func (bc *BaseComponent) GetProps() Props {
	return bc.props
}

// GetChildren returns the component's children
func (bc *BaseComponent) GetChildren() Children {
	return bc.children
}

// SetState updates a value in the component's state
func (bc *BaseComponent) SetState(key string, value interface{}) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.state[key] = value
}

// GetState retrieves a value from the component's state
func (bc *BaseComponent) GetState(key string) interface{} {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.state[key]
}

// Render returns an empty string for the base component
func (bc *BaseComponent) Render() string {
	return ""
}

// LifecycleComponent extends Component with lifecycle hooks
type LifecycleComponent interface {
	Component
	Mount(ctx context.Context) error
	Update(prevProps Props)
	Unmount()
}

// LifecycleComponentBase provides default no-op implementations of lifecycle hooks
type LifecycleComponentBase struct {
	BaseComponent
}

// Mount is called when the component is first rendered
func (lcb *LifecycleComponentBase) Mount(ctx context.Context) error {
	return nil
}

// Update is called when the component's props change
func (lcb *LifecycleComponentBase) Update(prevProps Props) {}

// Unmount is called when the component is removed from the tree
func (lcb *LifecycleComponentBase) Unmount() {}

// Store provides a simple global state container with subscription support
type Store struct {
	state      map[string]interface{}
	subscribers map[string][]chan map[string]interface{}
	mu         sync.RWMutex
}

// NewStore creates a new Store instance
func NewStore() *Store {
	return &Store{
		state:       make(map[string]interface{}),
		subscribers: make(map[string][]chan map[string]interface{}),
	}
}

// GlobalStore is the default global store instance
var GlobalStore = NewStore()

// GetState retrieves a value from the store
func (s *Store) GetState(key string) interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state[key]
}

// SetState updates a value in the store and notifies subscribers
func (s *Store) SetState(key string, value interface{}) {
	s.mu.Lock()
	s.state[key] = value
	subs := s.subscribers[key]
	s.mu.Unlock()

	for _, ch := range subs {
		select {
		case ch <- map[string]interface{}{key: value}:
		default:
			// drop if channel is full
		}
	}
}

// Subscribe registers a channel to receive updates for a given key
func (s *Store) Subscribe(key string) chan map[string]interface{} {
	ch := make(chan map[string]interface{}, 10)
	s.mu.Lock()
	s.subscribers[key] = append(s.subscribers[key], ch)
	s.mu.Unlock()
	return ch
}

// CreateElement renders a component tree into an HTML string
func CreateElement(tag string, props Props, children ...Component) string {
	var sb strings.Builder
	sb.WriteString("<")
	sb.WriteString(tag)

	if props != nil {
		for key, val := range props {
			switch v := val.(type) {
			case string:
				sb.WriteString(fmt.Sprintf(` %s="%s"`, key, v))
			case bool:
				if v {
					sb.WriteString(fmt.Sprintf(` %s`, key))
				}
			case int, float64:
				sb.WriteString(fmt.Sprintf(` %s="%v"`, key, v))
			}
		}
	}

	if len(children) == 0 {
		sb.WriteString(" />")
		return sb.String()
	}

	sb.WriteString(">")
	for _, child := range children {
		sb.WriteString(renderChild(child))
	}
	sb.WriteString(fmt.Sprintf("</%s>", tag))
	return sb.String()
}

// Fragment renders children without a wrapping element
func Fragment(children ...Component) string {
	var sb strings.Builder
	for _, child := range children {
		sb.WriteString(renderChild(child))
	}
	return sb.String()
}

// renderChild renders a single child component to HTML
func renderChild(c Component) string {
	if c == nil {
		return ""
	}
	return c.Render()
}

// PropsToJSON marshals Props to a JSON string
func PropsToJSON(p Props) string {
	if p == nil {
		return "{}"
	}
	b, err := json.Marshal(p)
	if err != nil {
		return "{}"
	}
	return string(b)
}
