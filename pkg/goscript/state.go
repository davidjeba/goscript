package goscript

import (
	"reflect"
	"sync"
)

type Store struct {
	state     map[string]interface{}
	listeners map[string][]func(interface{})
	mu        sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		state:     make(map[string]interface{}),
		listeners: make(map[string][]func(interface{})),
	}
}

func (s *Store) GetState(key string) interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state[key]
}

func (s *Store) SetState(key string, value interface{}) {
	s.mu.Lock()
	oldValue := s.state[key]
	s.state[key] = value
	s.mu.Unlock()

	if !reflect.DeepEqual(oldValue, value) {
		s.notifyListeners(key, value)
	}
}

func (s *Store) Subscribe(key string, listener func(interface{})) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.listeners[key] = append(s.listeners[key], listener)
}

func (s *Store) notifyListeners(key string, value interface{}) {
	s.mu.RLock()
	listeners := s.listeners[key]
	s.mu.RUnlock()

	for _, listener := range listeners {
		go listener(value)
	}
}

// Global store instance
var GlobalStore = NewStore()

// UseState is a hook-like function for components to use state
func UseState(key string, initialValue interface{}) (interface{}, func(interface{})) {
	if GlobalStore.GetState(key) == nil {
		GlobalStore.SetState(key, initialValue)
	}

	setState := func(newValue interface{}) {
		GlobalStore.SetState(key, newValue)
	}

	return GlobalStore.GetState(key), setState
}

