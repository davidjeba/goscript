package goscript

import (
	"sync"
)

// ComponentRegistry keeps track of mounted components
type ComponentRegistry struct {
	components map[string]LifecycleComponent
	mutex      sync.RWMutex
}

// NewComponentRegistry creates a new component registry
func NewComponentRegistry() *ComponentRegistry {
	return &ComponentRegistry{
		components: make(map[string]LifecycleComponent),
	}
}

// RegisterComponent adds a component to the registry
func (r *ComponentRegistry) RegisterComponent(id string, component LifecycleComponent) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.components[id] = component
	component.ComponentDidMount()
}

// UnregisterComponent removes a component from the registry
func (r *ComponentRegistry) UnregisterComponent(id string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if component, exists := r.components[id]; exists {
		component.ComponentWillUnmount()
		delete(r.components, id)
	}
}

// GetComponent retrieves a component from the registry
func (r *ComponentRegistry) GetComponent(id string) (LifecycleComponent, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	component, exists := r.components[id]
	return component, exists
}

// UpdateComponent updates a component's props if it exists
func (r *ComponentRegistry) UpdateComponent(id string, props Props) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if component, exists := r.components[id]; exists {
		if component.ShouldComponentUpdate(props) {
			// In a real implementation, we would update the props here
			// and trigger a re-render
			return true
		}
	}
	return false
}

// LifecycleComponentBase provides a base implementation of LifecycleComponent
type LifecycleComponentBase struct {
	BaseComponent
}

// ComponentDidMount is called when the component is mounted
func (l *LifecycleComponentBase) ComponentDidMount() {
	// Default implementation does nothing
}

// ComponentWillUnmount is called before the component is unmounted
func (l *LifecycleComponentBase) ComponentWillUnmount() {
	// Default implementation does nothing
}

// ShouldComponentUpdate determines if the component should update
func (l *LifecycleComponentBase) ShouldComponentUpdate(nextProps Props) bool {
	// Default implementation always updates
	return true
}

// Global component registry
var GlobalComponentRegistry = NewComponentRegistry()