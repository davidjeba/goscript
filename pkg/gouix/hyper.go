package gouix

import (
        "fmt"
        "reflect"
        "sync"
)

// Observer is a function that is called when a hyper(reactive) value changes
type Observer func(newValue interface{}, oldValue interface{})

// Signal represents a hyper(reactive) value
type Signal struct {
        value     interface{}
        observers []Observer
        mutex     sync.RWMutex
}

// NewSignal creates a new hyper(reactive) signal with an initial value
func NewSignal(initialValue interface{}) *Signal {
        return &Signal{
                value:     initialValue,
                observers: make([]Observer, 0),
        }
}

// Get returns the current value of the signal
func (s *Signal) Get() interface{} {
        s.mutex.RLock()
        defer s.mutex.RUnlock()
        return s.value
}

// Set updates the value of the signal and notifies observers
func (s *Signal) Set(newValue interface{}) {
        s.mutex.Lock()
        oldValue := s.value
        
        // Only update and notify if the value has changed
        if !reflect.DeepEqual(oldValue, newValue) {
                s.value = newValue
                observers := make([]Observer, len(s.observers))
                copy(observers, s.observers)
                s.mutex.Unlock()
                
                // Notify observers
                for _, observer := range observers {
                        observer(newValue, oldValue)
                }
        } else {
                s.mutex.Unlock()
        }
}

// Subscribe adds an observer to the signal
func (s *Signal) Subscribe(observer Observer) func() {
        s.mutex.Lock()
        defer s.mutex.Unlock()
        
        s.observers = append(s.observers, observer)
        
        // Return unsubscribe function
        return func() {
                s.mutex.Lock()
                defer s.mutex.Unlock()
                
                for i, o := range s.observers {
                        if reflect.ValueOf(o).Pointer() == reflect.ValueOf(observer).Pointer() {
                                s.observers = append(s.observers[:i], s.observers[i+1:]...)
                                break
                        }
                }
        }
}

// Computed represents a computed value that depends on other hyper(reactive) signals
type Computed struct {
        Signal
        compute     func() interface{}
        dependencies []*Signal
        unsubscribes []func()
}

// NewComputed creates a new computed value
func NewComputed(compute func() interface{}, dependencies ...*Signal) *Computed {
        c := &Computed{
                compute:      compute,
                dependencies: dependencies,
                unsubscribes: make([]func(), 0),
        }
        
        // Initialize with computed value
        c.Signal = *NewSignal(compute())
        
        // Subscribe to dependencies
        for _, dep := range dependencies {
                unsubscribe := dep.Subscribe(func(_, _  interface{}) {
                        c.recompute()
                })
                c.unsubscribes = append(c.unsubscribes, unsubscribe)
        }
        
        return c
}

// recompute updates the computed value
func (c *Computed) recompute() {
        newValue := c.compute()
        c.Signal.Set(newValue)
}

// Dispose cleans up subscriptions
func (c *Computed) Dispose() {
        for _, unsubscribe := range c.unsubscribes {
                unsubscribe()
        }
        c.unsubscribes = nil
}

// Effect runs a function when dependencies change
type Effect struct {
        run          func()
        dependencies []*Signal
        unsubscribes []func()
        disposed     bool
        mutex        sync.Mutex
}

// NewEffect creates a new effect
func NewEffect(run func(), dependencies ...*Signal) *Effect {
        e := &Effect{
                run:          run,
                dependencies: dependencies,
                unsubscribes: make([]func(), 0),
        }
        
        // Subscribe to dependencies
        for _, dep := range dependencies {
                unsubscribe := dep.Subscribe(func(_, _ interface{}) {
                        e.trigger()
                })
                e.unsubscribes = append(e.unsubscribes, unsubscribe)
        }
        
        // Run initially
        e.trigger()
        
        return e
}

// trigger runs the effect
func (e *Effect) trigger() {
        e.mutex.Lock()
        defer e.mutex.Unlock()
        
        if !e.disposed {
                e.run()
        }
}

// Dispose cleans up subscriptions
func (e *Effect) Dispose() {
        e.mutex.Lock()
        defer e.mutex.Unlock()
        
        e.disposed = true
        for _, unsubscribe := range e.unsubscribes {
                unsubscribe()
        }
        e.unsubscribes = nil
}

// Store represents a hyper(reactive) state store
type Store struct {
        state    map[string]*Signal
        computed map[string]*Computed
        mutex    sync.RWMutex
}

// NewStore creates a new store with initial state
func NewStore(initialState map[string]interface{}) *Store {
        store := &Store{
                state:    make(map[string]*Signal),
                computed: make(map[string]*Computed),
        }
        
        // Initialize state
        for key, value := range initialState {
                store.state[key] = NewSignal(value)
        }
        
        return store
}

// Get returns a signal for a state key
func (s *Store) Get(key string) *Signal {
        s.mutex.RLock()
        defer s.mutex.RUnlock()
        
        return s.state[key]
}

// Set updates a state value
func (s *Store) Set(key string, value interface{}) {
        s.mutex.Lock()
        
        // Create signal if it doesn't exist
        if _, exists := s.state[key]; !exists {
                s.state[key] = NewSignal(value)
                s.mutex.Unlock()
                return
        }
        
        signal := s.state[key]
        s.mutex.Unlock()
        
        // Update signal
        signal.Set(value)
}

// GetValue returns the current value for a state key
func (s *Store) GetValue(key string) interface{} {
        s.mutex.RLock()
        signal, exists := s.state[key]
        s.mutex.RUnlock()
        
        if !exists {
                return nil
        }
        
        return signal.Get()
}

// AddComputed adds a computed value to the store
func (s *Store) AddComputed(key string, compute func() interface{}, dependencies ...string) {
        s.mutex.Lock()
        defer s.mutex.Unlock()
        
        // Get dependency signals
        deps := make([]*Signal, len(dependencies))
        for i, depKey := range dependencies {
                deps[i] = s.state[depKey]
        }
        
        // Create computed
        s.computed[key] = NewComputed(compute, deps...)
}

// GetComputed returns a computed value
func (s *Store) GetComputed(key string) *Computed {
        s.mutex.RLock()
        defer s.mutex.RUnlock()
        
        return s.computed[key]
}

// Subscribe subscribes to changes in a state key
func (s *Store) Subscribe(key string, observer Observer) func() {
        s.mutex.RLock()
        signal, exists := s.state[key]
        s.mutex.RUnlock()
        
        if !exists {
                // Create signal if it doesn't exist
                s.mutex.Lock()
                signal = NewSignal(nil)
                s.state[key] = signal
                s.mutex.Unlock()
        }
        
        return signal.Subscribe(observer)
}

// Dispose cleans up the store
func (s *Store) Dispose() {
        s.mutex.Lock()
        defer s.mutex.Unlock()
        
        // Dispose computed values
        for _, computed := range s.computed {
                computed.Dispose()
        }
        
        s.state = make(map[string]*Signal)
        s.computed = make(map[string]*Computed)
}

// HyperComponent is a component that uses hyper(reactive) state
type HyperComponent struct {
        BaseComponent
        store *Store
}

// NewHyperComponent creates a new hyper(reactive) component
func NewHyperComponent(id ComponentID, props Props, initialState map[string]interface{}) *HyperComponent {
        base := NewBaseComponent(id, props)
        
        return &HyperComponent{
                BaseComponent: *base,
                store:         NewStore(initialState),
        }
}

// GetStore returns the component's store
func (h *HyperComponent) GetStore() *Store {
        return h.store
}

// GetState returns a state value
func (h *HyperComponent) GetState(key string) interface{} {
        return h.store.GetValue(key)
}

// SetState updates a state value
func (h *HyperComponent) SetState(key string, value interface{}) {
        h.store.Set(key, value)
}

// Watch subscribes to changes in a state key
func (h *HyperComponent) Watch(key string, observer Observer) func() {
        return h.store.Subscribe(key, observer)
}

// Unmount cleans up the component
func (h *HyperComponent) Unmount() {
        h.store.Dispose()
        h.BaseComponent.Unmount()
}

// currentHookComponent is the component currently being rendered
var currentHookComponent *HyperComponent

// SetCurrentHookComponent sets the current hook component
func SetCurrentHookComponent(component *HyperComponent) {
        currentHookComponent = component
}

// UseState creates a state hook
func UseState(initialValue interface{}) *Signal {
        if currentHookComponent == nil {
                // Create a temporary signal if no component is set
                return NewSignal(initialValue)
        }
        
        // Generate a unique key for this state
        key := fmt.Sprintf("state_%p", &initialValue)
        
        // Check if state already exists
        if signal := currentHookComponent.GetStore().Get(key); signal != nil {
                return signal
        }
        
        // Create new state
        currentHookComponent.GetStore().Set(key, initialValue)
        return currentHookComponent.GetStore().Get(key)
}