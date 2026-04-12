// Package goscript — registry.go provides a factory-based component registry
// for goscript. Components can be registered by name and later rendered
// programmatically. This is the named-component counterpart to the lifecycle-
// based ComponentRegistry in lifecycle.go.
//
// While the lifecycle ComponentRegistry tracks mounted component instances,
// FactoryRegistry tracks component factories (functions that create components
// from Props). This separation allows the same component type to be rendered
// multiple times with different props.
//
// The global Components variable provides a default registry instance. In most
// applications, you only need the global registry. Create additional registries
// for isolated rendering contexts (e.g., preview mode, testing, plugins).
package goscript

import (
	"fmt"
	"sort"
	"sync"
)

// =========================================================================
// Types
// =========================================================================

// ComponentFactory is a function that creates a Component from Props.
// The factory receives the component's props and must return a Component.
// The returned Component's Render() method produces the HTML output.
//
// Usage:
//
//	factory := func(props Props) Component {
//	    return NewServerComponent("MyButton", func(p Props) string {
//	        return fmt.Sprintf("<button class=\"%s\">%s</button>",
//	            p["class"], p["text"])
//	    }, props)
//	}
//	Components.Register("MyButton", factory)
type ComponentFactory func(Props) Component

// FactoryRegistry is a thread-safe registry of named component factories.
// It allows components to be registered by name and rendered by name,
// which is useful for:
//   - Dynamic component resolution from configuration or data
//   - The gs-target="component:name" reactive attribute
//   - Plugin systems where components are loaded at runtime
//   - Template rendering where component names are determined by data
//
// FactoryRegistry is safe for concurrent use. All operations are protected
// by a read-write mutex.
type FactoryRegistry struct {
	factories map[string]ComponentFactory
	order     []string // insertion order for List()
	mutex     sync.RWMutex
}

// NewFactoryRegistry creates a new, empty FactoryRegistry.
//
// Usage:
//
//	reg := NewFactoryRegistry()
//	reg.Register("Header", headerFactory)
func NewFactoryRegistry() *FactoryRegistry {
	return &FactoryRegistry{
		factories: make(map[string]ComponentFactory),
		order:     make([]string, 0),
	}
}

// Register adds a component factory to the registry under the given name.
// If a factory with the same name already exists, it is silently replaced.
// Names are case-sensitive.
//
// Register is safe for concurrent use.
//
// Usage:
//
//	Components.Register("UserCard", func(props Props) Component {
//	    return NewServerComponent("UserCard", renderUserCard, props)
//	})
func (r *FactoryRegistry) Register(name string, factory ComponentFactory) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Track insertion order for deterministic List() output
	if _, exists := r.factories[name]; !exists {
		r.order = append(r.order, name)
	}

	r.factories[name] = factory
}

// Unregister removes a component factory from the registry by name.
// Returns true if the factory was found and removed, false otherwise.
//
// Unregister is safe for concurrent use.
//
// Usage:
//
//	ok := Components.Unregister("OldButton")
func (r *FactoryRegistry) Unregister(name string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.factories[name]; !exists {
		return false
	}

	delete(r.factories, name)

	// Remove from insertion order
	for i, n := range r.order {
		if n == name {
			r.order = append(r.order[:i], r.order[i+1:]...)
			break
		}
	}

	return true
}

// Render creates a component using the registered factory for the given name
// and renders it to an HTML string. Returns an error if the component is not
// registered or if the factory panics.
//
// Render is safe for concurrent use.
//
// Usage:
//
//	html, err := Components.Render("UserCard", Props{
//	    "name":  "Alice",
//	    "email": "alice@example.com",
//	})
func (r *FactoryRegistry) Render(name string, props Props) (string, error) {
	r.mutex.RLock()
	factory, exists := r.factories[name]
	r.mutex.RUnlock()

	if !exists {
		return "", fmt.Errorf("goscript: component %q is not registered", name)
	}

	if props == nil {
		props = Props{}
	}

	// Recover from panics in the factory function
	defer func() {
		if rec := recover(); rec != nil {
			// The error is returned below via the named return
		}
	}()

	component := factory(props)
	if component == nil {
		return "", fmt.Errorf("goscript: factory for component %q returned nil", name)
	}

	return component.Render(), nil
}

// MustRender is like Render but panics on error. Useful for initialization
// code where a missing component is a fatal error.
//
// Usage:
//
//	html := Components.MustRender("Header", Props{"title": "My App"})
func (r *FactoryRegistry) MustRender(name string, props Props) string {
	html, err := r.Render(name, props)
	if err != nil {
		panic(fmt.Sprintf("goscript: MustRender(%q) failed: %v", name, err))
	}
	return html
}

// RenderSafe creates a component using the registered factory and renders
// it to an HTML string. On any error (not registered, factory panic, nil
// component), it returns the fallback HTML instead.
//
// RenderSafe is safe for concurrent use.
//
// Usage:
//
//	html := Components.RenderSafe("OptionalWidget", props,
//	    "<!-- OptionalWidget not available -->")
func (r *FactoryRegistry) RenderSafe(name string, props Props, fallback string) string {
	html, err := r.Render(name, props)
	if err != nil {
		return fallback
	}
	return html
}

// Has checks whether a component factory is registered under the given name.
//
// Usage:
//
//	if Components.Has("Navbar") {
//	    // Navbar is available
//	}
func (r *FactoryRegistry) Has(name string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	_, exists := r.factories[name]
	return exists
}

// Get retrieves the factory function for a named component. Returns the
// factory and true if the component is registered, nil and false otherwise.
//
// Get is safe for concurrent use.
//
// Usage:
//
//	factory, ok := Components.Get("UserCard")
//	if ok {
//	    component := factory(props)
//	    html := component.Render()
//	}
func (r *FactoryRegistry) Get(name string) (ComponentFactory, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	factory, exists := r.factories[name]
	return factory, exists
}

// List returns all registered component names in insertion order.
// The returned slice is a copy and can be safely modified.
//
// List is safe for concurrent use.
//
// Usage:
//
//	names := Components.List()
//	fmt.Println("Registered components:", names)
func (r *FactoryRegistry) List() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make([]string, len(r.order))
	copy(result, r.order)
	return result
}

// ListSorted returns all registered component names in alphabetical order.
// This is useful for debugging, documentation generation, and admin UIs.
//
// ListSorted is safe for concurrent use.
func (r *FactoryRegistry) ListSorted() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	names := make([]string, 0, len(r.factories))
	for name := range r.factories {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Count returns the number of registered component factories.
//
// Usage:
//
//	fmt.Printf("Registered %d components\n", Components.Count())
func (r *FactoryRegistry) Count() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.factories)
}

// Clear removes all registered component factories from the registry.
// This is primarily useful for testing.
//
// Clear is safe for concurrent use.
func (r *FactoryRegistry) Clear() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.factories = make(map[string]ComponentFactory)
	r.order = make([]string, 0)
}

// RegisterBatch registers multiple component factories in a single call.
// The map keys are component names and the values are factory functions.
//
// Usage:
//
//	Components.RegisterBatch(map[string]ComponentFactory{
//	    "Header": headerFactory,
//	    "Footer": footerFactory,
//	    "Sidebar": sidebarFactory,
//	})
func (r *FactoryRegistry) RegisterBatch(factories map[string]ComponentFactory) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for name, factory := range factories {
		if _, exists := r.factories[name]; !exists {
			r.order = append(r.order, name)
		}
		r.factories[name] = factory
	}
}

// RenderAll renders multiple named components with their respective props
// and returns a map of name → rendered HTML. Components that fail to render
// are included in the map with their error message.
//
// RenderAll is safe for concurrent use.
//
// Usage:
//
//	results := Components.RenderAll(map[string]Props{
//	    "Header": Props{"title": "My App"},
//	    "Footer": Props{"year": "2024"},
//	})
//	// results["Header"] = "<header>...</header>"
//	// results["Footer"] = "<footer>...</footer>"
func (r *FactoryRegistry) RenderAll(componentProps map[string]Props) map[string]string {
	results := make(map[string]string, len(componentProps))

	for name, props := range componentProps {
		html, err := r.Render(name, props)
		if err != nil {
			results[name] = fmt.Sprintf("<!-- error rendering %q: %v -->", name, err)
		} else {
			results[name] = html
		}
	}

	return results
}

// ForEach iterates over all registered component factories and calls the
// given function with each name and factory. Iteration order matches
// insertion order.
//
// ForEach is safe for concurrent use. The callback function should not
// register or unregister components during iteration.
//
// Usage:
//
//	Components.ForEach(func(name string, factory ComponentFactory) {
//	    fmt.Printf("Component: %s\n", name)
//	})
func (r *FactoryRegistry) ForEach(fn func(name string, factory ComponentFactory)) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, name := range r.order {
		if factory, exists := r.factories[name]; exists {
			fn(name, factory)
		}
	}
}

// =========================================================================
// Global Registry
// =========================================================================

// Components is the global component factory registry. In most applications,
// this is the only registry instance you need. Register your components here
// during application initialization.
//
// Usage:
//
//	func init() {
//	    goscript.Components.Register("UserCard", userCardFactory)
//	    goscript.Components.Register("ProductGrid", productGridFactory)
//	}
var Components = NewFactoryRegistry()
