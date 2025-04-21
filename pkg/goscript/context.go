package goscript

import (
        "strings"
        "sync"
)

// Context represents a component context
type Context struct {
        values map[string]interface{}
        parent *Context
        mutex  sync.RWMutex
}

// NewContext creates a new context
func NewContext(parent *Context) *Context {
        return &Context{
                values: make(map[string]interface{}),
                parent: parent,
        }
}

// Set sets a value in the context
func (c *Context) Set(key string, value interface{}) {
        c.mutex.Lock()
        defer c.mutex.Unlock()
        c.values[key] = value
}

// Get gets a value from the context
func (c *Context) Get(key string) (interface{}, bool) {
        c.mutex.RLock()
        defer c.mutex.RUnlock()
        
        // Check current context
        if value, exists := c.values[key]; exists {
                return value, true
        }
        
        // Check parent context
        if c.parent != nil {
                return c.parent.Get(key)
        }
        
        return nil, false
}

// CreateProvider creates a context provider component
func CreateProvider(context *Context, key string, value interface{}) FunctionalComponent {
        return func(props Props, children ...interface{}) string {
                // Set the value in the context
                context.Set(key, value)
                
                // Render children
                var result strings.Builder
                for _, child := range children {
                        result.WriteString(renderChild(child))
                }
                return result.String()
        }
}

// CreateConsumer creates a context consumer component
func CreateConsumer(context *Context, key string, render func(value interface{}) string) FunctionalComponent {
        return func(props Props, children ...interface{}) string {
                // Get the value from the context
                value, _ := context.Get(key)
                
                // Render using the provided function
                return render(value)
        }
}

// GlobalContext is the root context
var GlobalContext = NewContext(nil)

// WithContext creates a new context with the given parent
func WithContext(parent *Context) *Context {
        if parent == nil {
                parent = GlobalContext
        }
        return NewContext(parent)
}

// ContextProvider is a component that provides context values
type ContextProvider struct {
        LifecycleComponentBase
        context *Context
        key     string
        value   interface{}
}

// NewContextProvider creates a new context provider
func NewContextProvider(context *Context, key string, value interface{}, props Props, children ...interface{}) *ContextProvider {
        base := NewBaseComponent(props, nil, children...)
        provider := &ContextProvider{
                context: context,
                key:     key,
                value:   value,
        }
        provider.LifecycleComponentBase.BaseComponent = *base
        return provider
}

// Render implements the Component interface
func (p *ContextProvider) Render() string {
        // Set the value in the context
        p.context.Set(p.key, p.value)
        
        // Render children
        var result strings.Builder
        for _, child := range p.GetChildren() {
                result.WriteString(renderChild(child))
        }
        return result.String()
}

// ComponentDidMount implements the LifecycleComponent interface
func (p *ContextProvider) ComponentDidMount() {
        // Set the value in the context when mounted
        p.context.Set(p.key, p.value)
}

// ComponentWillUnmount implements the LifecycleComponent interface
func (p *ContextProvider) ComponentWillUnmount() {
        // No need to remove from context, as contexts are hierarchical
}

// ContextConsumer is a component that consumes context values
type ContextConsumer struct {
        LifecycleComponentBase
        context *Context
        key     string
        render  func(value interface{}) string
}

// NewContextConsumer creates a new context consumer
func NewContextConsumer(context *Context, key string, render func(value interface{}) string, props Props) *ContextConsumer {
        base := NewBaseComponent(props, nil)
        consumer := &ContextConsumer{
                context: context,
                key:     key,
                render:  render,
        }
        consumer.LifecycleComponentBase.BaseComponent = *base
        return consumer
}

// Render implements the Component interface
func (c *ContextConsumer) Render() string {
        // Get the value from the context
        value, _ := c.context.Get(c.key)
        
        // Render using the provided function
        return c.render(value)
}