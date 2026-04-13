// Package state provides reactive state management functions for goscript .gs files.
// These functions are compiled to JavaScript calls to __gs.state by the GS compiler.
//
// The state system mirrors a subset of React's hooks API, providing reactive state
// slots that notify subscribers on every change. Each state slot is identified by a
// unique numeric ID managed internally by the runtime.
//
// This package serves as API documentation and type definitions for .gs developers.
// The functions listed here map to JavaScript operations in the goscript client runtime
// (pkg/gslib/runtime.js). You do not need to import this package in Go server code —
// it exists solely for .gs files.
//
// # State Model
//
// goscript state is built around two primitives:
//   - UseState: creates a reactive [getter, setter] pair
//   - UseReducer: creates state driven by a reducer function
//
// Hooks like UseEffect, UseRef, UseMemo, and UseCallback provide lifecycle
// and optimization patterns familiar to React developers.
//
// # Usage in .gs code
//
//	import "goscript/state"
//
//	func Counter() {
//	    count, setCount := state.UseState(0)
//
//	    incBtn := dom.CreateElement("button")
//	    dom.SetTextContent(incBtn, "Increment")
//	    dom.AddEventListener(incBtn, "click", func() {
//	        cur := count()
//	        setCount(cur + 1)
//	    })
//
//	    label := dom.CreateElement("span")
//	    dom.SetTextContent(label, fmt.Sprintf("Count: %d", count()))
//
//	    root := dom.CreateElement("div")
//	    dom.AppendChild(root, label)
//	    dom.AppendChild(root, incBtn)
//	    return root
//	}
package state

// Store represents a global reactive state store.
// In .gs code, the GS compiler manages store instances via __gs.state.
// Stores contain named state slots that can be read and updated reactively.
type Store struct {
	ID    int
	Value interface{}
}

// Ref represents a mutable reference object, similar to React's useRef.
// A Ref has a single field, Current, that can be read and written freely
// without triggering re-renders.
type Ref struct {
	Current interface{}
}

// Getter is a function that returns the current value of a state slot.
// In compiled .gs code, this maps to the getter function returned by __gs.useState.
type Getter func() interface{}

// Setter is a function that updates a state slot. It accepts either a direct
// value or an updater function that receives the previous value and returns
// the new value.
// In compiled .gs code, this maps to the setter function returned by __gs.useState.
type Setter func(value interface{})

// Dispatcher is a function that sends an action to a reducer-based state slot.
// The reducer receives (prevState, action) and returns the new state.
type Dispatcher func(action interface{})

// Reducer is a function that takes the current state and an action, and returns
// the new state. Used with UseReducer.
type Reducer func(state interface{}, action interface{}) interface{}

// Unsubscriber is a function that cancels a previously registered subscription.
// Call it to stop receiving state change notifications.
type Unsubscriber func()

// ---------------------------------------------------------------------------
// Core State Primitives
// ---------------------------------------------------------------------------

// CreateStore creates a new named store with the given initial value.
// The store is registered in the global state registry and can be accessed
// by its ID. In .gs code this maps to internal __gs.state management.
//
// Parameters:
//   - initial: the initial value for the store
//
// Returns a Store instance.
//
// Example (.gs):
//
//	store := state.CreateStore(map[string]interface{}{
//	    "count": 0,
//	    "name":  "goscript",
//	})
func CreateStore(initial interface{}) Store { return Store{} }

// UseState creates a reactive state slot and returns a [getter, setter] pair.
//
// The getter is a zero-argument function that returns the current value.
// The setter accepts either a direct value or an updater function(prev) => next.
// If the new value differs from the old value (strict equality), all subscribers
// are notified with (newVal, oldVal).
//
// In .gs code this compiles to __gs.useState(initial).
//
// Parameters:
//   - initial: the initial value, or a function () => initialValue (thunk)
//
// Returns [getter, setter].
//
// Example (.gs):
//
//	count, setCount := state.UseState(0)
//
//	// Read current value
//	cur := count()
//
//	// Update with direct value
//	setCount(5)
//
//	// Update with updater function
//	setCount(func(prev interface{}) interface{} {
//	    return prev.(int) + 1
//	})
func UseState(initial interface{}) []interface{} { return nil }

// UseReducer creates a state slot driven by a reducer function.
// Returns [stateGetter, dispatch].
//
// The reducer function receives (currentState, action) and must return the new state.
// The dispatch function sends an action to the reducer, which computes the new state.
//
// In .gs code this compiles to __gs.useReducer(reducer, initial).
//
// Parameters:
//   - reducer: a function(state, action) => newState
//   - initial: the initial state value
//
// Returns [getter, dispatch].
//
// Example (.gs):
//
//	reducer := func(s interface{}, a interface{}) interface{} {
//	    st := s.(map[string]interface{})
//	    act := a.(string)
//	    if act == "increment" {
//	        st["count"] = st["count"].(int) + 1
//	    }
//	    return st
//	}
//
//	stateGetter, dispatch := state.UseReducer(reducer, map[string]interface{}{
//	    "count": 0,
//	})
//
//	// Dispatch an action
//	dispatch("increment")
func UseReducer(reducer Reducer, initial interface{}) []interface{} { return nil }

// ---------------------------------------------------------------------------
// Lifecycle Hooks
// ---------------------------------------------------------------------------

// UseEffect runs a side-effect function. In the current runtime, the effect
// runs immediately when called. If the function returns a cleanup function,
// it will be available for future dependency-based re-running.
//
// In .gs code this compiles to __gs.useEffect(fn, deps).
//
// Parameters:
//   - fn: the effect function. May optionally return a cleanup function.
//   - deps: dependency array (reserved for future use; pass nil for now)
//
// Returns a cleanup function (if fn returned one), or nil.
//
// Example (.gs):
//
//	state.UseEffect(func() interface{} {
//	    // Setup: subscribe to an event
//	    unsub := realtime.On("tick", func(detail interface{}) {
//	        fmt.Println("tick:", detail)
//	    })
//	    // Return cleanup function
//	    return func() { unsub() }
//	}, nil)
func UseEffect(fn func() interface{}, deps []interface{}) interface{} { return nil }

// UseRef creates a mutable reference object with an initial value.
// Unlike UseState, reading or writing Ref.Current does not trigger re-renders.
// Useful for holding DOM references, timers, or any mutable value that
// persists across renders without causing updates.
//
// In .gs code this compiles to __gs.useRef(initial).
//
// Parameters:
//   - initial: the initial value for ref.Current
//
// Returns a Ref with the given initial value.
//
// Example (.gs):
//
//	inputRef := state.UseRef(nil)
//
//	// Later, assign a DOM element to the ref
//	inputRef.Current = dom.GetElementByID("myInput")
//
//	// Read the ref
//	el := inputRef.Current
//	dom.Focus(el)
func UseRef(initial interface{}) Ref { return Ref{} }

// UseMemo computes a memoized value. In the current runtime, the function
// is called immediately on every invocation. Dependency-based caching will
// be added in a future version.
//
// Use useMemo to avoid expensive recalculation on every render.
//
// In .gs code this compiles to __gs.useMemo(fn, deps).
//
// Parameters:
//   - fn: a computation function that returns the memoized value
//   - deps: dependency array (reserved for future use; pass nil for now)
//
// Returns the computed value.
//
// Example (.gs):
//
//	expensive := state.UseMemo(func() interface{} {
//	    return computeExpensiveThing(items)
//	}, nil)
func UseMemo(fn func() interface{}, deps []interface{}) interface{} { return nil }

// UseCallback returns a memoized callback function. In the current runtime,
// the function is returned unchanged. Dependency-based caching will be
// added in a future version.
//
// Use useCallback to prevent unnecessary re-renders when passing callbacks
// to child components.
//
// In .gs code this compiles to __gs.useCallback(fn, deps).
//
// Parameters:
//   - fn: the callback function to memoize
//   - deps: dependency array (reserved for future use; pass nil for now)
//
// Returns the (memoized) function.
//
// Example (.gs):
//
//	handleClick := state.UseCallback(func(args ...interface{}) interface{} {
//	    fmt.Println("clicked")
//	    return nil
//	}, nil)
//	dom.AddEventListener(btn, "click", handleClick)
func UseCallback(fn Func, deps []interface{}) Func { return nil }

// Func is a generic function type used for callbacks and event handlers.
type Func func(...interface{}) interface{}

// ---------------------------------------------------------------------------
// State Subscription & Serialization
// ---------------------------------------------------------------------------

// Subscribe registers a callback to be notified when a specific state slot changes.
// The callback receives (newVal, oldVal) whenever the state is updated.
//
// In .gs code this compiles to __gs.subscribe(stateId, callback).
//
// Parameters:
//   - stateId: the numeric ID of the state slot to watch
//   - callback: a function(newVal, oldVal) called on every change
//
// Returns an Unsubscriber function. Call it to stop receiving notifications.
//
// Example (.gs):
//
//	unsub := state.Subscribe(1, func(newVal interface{}, oldVal interface{}) {
//	    fmt.Println("State changed:", oldVal, "->", newVal)
//	})
//
//	// Later, stop listening
//	unsub()
func Subscribe(stateId int, callback func(newVal interface{}, oldVal interface{})) Unsubscriber {
	return func() {}
}

// GetState returns a snapshot of all state as a map.
// This is useful for debugging, serialization, or sending state to the server.
//
// In .gs code this compiles to __gs.getState().
//
// Returns a map[string]interface{} where keys are stringified state IDs
// and values are the current state values.
//
// Example (.gs):
//
//	snapshot := state.GetState()
//	fmt.Println("All state:", snapshot)
func GetState() map[string]interface{} { return nil }

// Hydrate populates the client state from server-provided data.
// This is called automatically during initialization from
// window.__GOSCRIPT_STATE__. You rarely need to call it manually
// unless implementing custom hydration logic.
//
// In .gs code this compiles to __gs.hydrate(serverState).
//
// Parameters:
//   - serverState: a map of state ID to value, typically from the server
//
// Example (.gs):
//
//	state.Hydrate(map[string]interface{}{
//	    "1": 42,
//	    "2": "hello",
//	})
func Hydrate(serverState map[string]interface{}) {}
