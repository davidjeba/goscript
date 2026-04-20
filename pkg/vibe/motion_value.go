package vibe

import (
	"reflect"
	"sync"
)

// ChangeEvent represents a motion value update.
type ChangeEvent struct {
	Previous interface{} `json:"previous"`
	Current  interface{} `json:"current"`
	Velocity float64     `json:"velocity"`
}

// ChangeHandler receives motion value updates.
type ChangeHandler func(ChangeEvent)

// MotionValue is a subscribable animatable value.
type MotionValue struct {
	mutex       sync.RWMutex
	value       interface{}
	velocity    float64
	subscribers map[int]ChangeHandler
	nextID      int
}

// NewMotionValue creates a new motion value.
func NewMotionValue(initial interface{}) *MotionValue {
	return &MotionValue{
		value:       initial,
		subscribers: make(map[int]ChangeHandler),
	}
}

// Get returns the current value.
func (mv *MotionValue) Get() interface{} {
	mv.mutex.RLock()
	defer mv.mutex.RUnlock()
	return mv.value
}

// Velocity returns the last measured velocity.
func (mv *MotionValue) Velocity() float64 {
	mv.mutex.RLock()
	defer mv.mutex.RUnlock()
	return mv.velocity
}

// Set updates the value and notifies subscribers.
func (mv *MotionValue) Set(next interface{}) {
	mv.mutex.Lock()
	previous := mv.value
	if reflect.DeepEqual(previous, next) {
		mv.mutex.Unlock()
		return
	}

	mv.value = next
	if prevFloat, ok := asFloat(previous); ok {
		if nextFloat, ok := asFloat(next); ok {
			mv.velocity = nextFloat - prevFloat
		}
	}

	event := ChangeEvent{
		Previous: previous,
		Current:  next,
		Velocity: mv.velocity,
	}

	handlers := make([]ChangeHandler, 0, len(mv.subscribers))
	for _, handler := range mv.subscribers {
		handlers = append(handlers, handler)
	}
	mv.mutex.Unlock()

	for _, handler := range handlers {
		handler(event)
	}
}

// Subscribe registers a change listener and returns an unsubscribe function.
func (mv *MotionValue) Subscribe(handler ChangeHandler) func() {
	mv.mutex.Lock()
	id := mv.nextID
	mv.nextID++
	mv.subscribers[id] = handler
	mv.mutex.Unlock()

	return func() {
		mv.mutex.Lock()
		delete(mv.subscribers, id)
		mv.mutex.Unlock()
	}
}

func asFloat(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}
