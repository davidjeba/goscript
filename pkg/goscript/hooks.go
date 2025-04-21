package goscript

import (
	"sync"
)

// HookContext represents the context for hooks
type HookContext struct {
	componentID string
	hooks       []interface{}
	index       int
	mutex       sync.RWMutex
}

// hookContexts stores hook contexts by component ID
var hookContexts = struct {
	contexts map[string]*HookContext
	mutex    sync.RWMutex
}{
	contexts: make(map[string]*HookContext),
}

// getCurrentHookContext gets the current hook context
func getCurrentHookContext(componentID string) *HookContext {
	hookContexts.mutex.RLock()
	context, exists := hookContexts.contexts[componentID]
	hookContexts.mutex.RUnlock()
	
	if !exists {
		hookContexts.mutex.Lock()
		context = &HookContext{
			componentID: componentID,
			hooks:       make([]interface{}, 0),
			index:       0,
		}
		hookContexts.contexts[componentID] = context
		hookContexts.mutex.Unlock()
	}
	
	return context
}

// resetHookContext resets the hook context for a component
func resetHookContext(componentID string) {
	hookContexts.mutex.Lock()
	defer hookContexts.mutex.Unlock()
	
	if context, exists := hookContexts.contexts[componentID]; exists {
		context.index = 0
	}
}

// cleanupHookContext removes the hook context for a component
func cleanupHookContext(componentID string) {
	hookContexts.mutex.Lock()
	defer hookContexts.mutex.Unlock()
	
	delete(hookContexts.contexts, componentID)
}

// useState is a hook for component state
func useState(componentID string, initialState interface{}) (interface{}, func(interface{})) {
	context := getCurrentHookContext(componentID)
	
	context.mutex.Lock()
	defer context.mutex.Unlock()
	
	// Initialize hook if needed
	if context.index >= len(context.hooks) {
		context.hooks = append(context.hooks, initialState)
	}
	
	// Get the current state
	stateIndex := context.index
	state := context.hooks[stateIndex]
	
	// Create the setState function
	setState := func(newState interface{}) {
		context.mutex.Lock()
		defer context.mutex.Unlock()
		
		context.hooks[stateIndex] = newState
		
		// In a real implementation, this would trigger a re-render
	}
	
	context.index++
	return state, setState
}

// useEffect is a hook for side effects
func useEffect(componentID string, effect func() func(), deps []interface{}) {
	context := getCurrentHookContext(componentID)
	
	context.mutex.Lock()
	defer context.mutex.Unlock()
	
	// Initialize hook if needed
	if context.index >= len(context.hooks) {
		context.hooks = append(context.hooks, struct {
			deps        []interface{}
			cleanup     func()
			initialized bool
		}{
			deps:        deps,
			cleanup:     nil,
			initialized: false,
		})
	}
	
	// Get the current effect state
	effectIndex := context.index
	effectState := context.hooks[effectIndex].(struct {
		deps        []interface{}
		cleanup     func()
		initialized bool
	})
	
	// Check if deps have changed
	depsChanged := !effectState.initialized
	if !depsChanged && deps != nil {
		oldDeps := effectState.deps
		if len(oldDeps) != len(deps) {
			depsChanged = true
		} else {
			for i := range deps {
				if deps[i] != oldDeps[i] {
					depsChanged = true
					break
				}
			}
		}
	}
	
	// Run effect if deps have changed
	if depsChanged {
		// Run cleanup if it exists
		if effectState.cleanup != nil {
			effectState.cleanup()
		}
		
		// Run effect
		cleanup := effect()
		
		// Update effect state
		context.hooks[effectIndex] = struct {
			deps        []interface{}
			cleanup     func()
			initialized bool
		}{
			deps:        deps,
			cleanup:     cleanup,
			initialized: true,
		}
	}
	
	context.index++
}

// useContext is a hook for consuming context
func useContext(componentID string, context *Context, key string) interface{} {
	hookContext := getCurrentHookContext(componentID)
	
	hookContext.mutex.Lock()
	defer hookContext.mutex.Unlock()
	
	// Initialize hook if needed
	if hookContext.index >= len(hookContext.hooks) {
		value, _ := context.Get(key)
		hookContext.hooks = append(hookContext.hooks, value)
	}
	
	// Get the current context value
	contextIndex := hookContext.index
	value, _ := context.Get(key)
	
	// Update the stored value
	hookContext.hooks[contextIndex] = value
	
	hookContext.index++
	return value
}

// useMemo is a hook for memoized values
func useMemo(componentID string, compute func() interface{}, deps []interface{}) interface{} {
	context := getCurrentHookContext(componentID)
	
	context.mutex.Lock()
	defer context.mutex.Unlock()
	
	// Initialize hook if needed
	if context.index >= len(context.hooks) {
		context.hooks = append(context.hooks, struct {
			value interface{}
			deps  []interface{}
		}{
			value: compute(),
			deps:  deps,
		})
	}
	
	// Get the current memo state
	memoIndex := context.index
	memoState := context.hooks[memoIndex].(struct {
		value interface{}
		deps  []interface{}
	})
	
	// Check if deps have changed
	depsChanged := false
	if deps != nil {
		oldDeps := memoState.deps
		if len(oldDeps) != len(deps) {
			depsChanged = true
		} else {
			for i := range deps {
				if deps[i] != oldDeps[i] {
					depsChanged = true
					break
				}
			}
		}
	}
	
	// Recompute if deps have changed
	if depsChanged {
		newValue := compute()
		context.hooks[memoIndex] = struct {
			value interface{}
			deps  []interface{}
		}{
			value: newValue,
			deps:  deps,
		}
	}
	
	context.index++
	return memoState.value
}

// useCallback is a hook for memoized callbacks
func useCallback(componentID string, callback func(), deps []interface{}) func() {
	return useMemo(componentID, func() interface{} {
		return callback
	}, deps).(func())
}

// useRef is a hook for mutable refs
func useRef(componentID string, initialValue interface{}) *struct{ Current interface{} } {
	context := getCurrentHookContext(componentID)
	
	context.mutex.Lock()
	defer context.mutex.Unlock()
	
	// Initialize hook if needed
	if context.index >= len(context.hooks) {
		context.hooks = append(context.hooks, &struct{ Current interface{} }{
			Current: initialValue,
		})
	}
	
	// Get the current ref
	refIndex := context.index
	ref := context.hooks[refIndex].(*struct{ Current interface{} })
	
	context.index++
	return ref
}