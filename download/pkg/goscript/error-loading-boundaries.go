package goscript

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

// ErrorBoundary catches errors from child components and renders a fallback
type ErrorBoundary struct {
	BaseComponent
	fallback  Component
	hasError  bool
	errorInfo error
	children  []Component
	resetKeys []interface{}
}

// NewErrorBoundary creates a new error boundary
func NewErrorBoundary(fallback Component, children ...Component) *ErrorBoundary {
	return &ErrorBoundary{
		fallback: fallback,
		children: children,
	}
}

// WithResetKeys allows the boundary to reset when keys change
func (eb *ErrorBoundary) WithResetKeys(keys ...interface{}) *ErrorBoundary {
	eb.resetKeys = keys
	return eb
}

// Render returns the fallback if an error occurred, otherwise renders children
func (eb *ErrorBoundary) Render() string {
	if eb.hasError {
		return eb.fallback.Render()
	}
	var result string
	for _, child := range eb.children {
		result += safeRender(child)
	}
	return result
}

// CatchError checks if a component rendered with an error
func (eb *ErrorBoundary) CatchError(err error) {
	if err != nil {
		eb.hasError = true
		eb.errorInfo = err
		log.Printf("[ErrorBoundary] Caught: %v", err)
	}
}

func safeRender(c Component) string {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ErrorBoundary] Panic recovered: %v", r)
		}
	}()
	return c.Render()
}

// LoadingBoundary wraps async components with loading states
type LoadingBoundary struct {
	ID            string
	Fallback      Component
	Children      []Component
	LoadingStates map[string]bool
	mu            sync.RWMutex
}

// IsLoading returns true if any children are loading
func (lb *LoadingBoundary) IsLoading() bool {
	lb.mu.RLock()
	defer lb.mu.RUnlock()
	for _, loading := range lb.LoadingStates {
		if loading {
			return true
		}
	}
	return false
}

// Render renders children or fallback based on loading state
func (lb *LoadingBoundary) Render() string {
	if lb.IsLoading() {
		fallbackHTML := ""
		if lb.Fallback != nil {
			fallbackHTML = lb.Fallback.Render()
		}
		return fmt.Sprintf(`<div id="%s" data-loading-boundary>%s</div>`, lb.ID, fallbackHTML)
	}

	var result string
	for _, child := range lb.Children {
		result += child.Render()
	}
	return fmt.Sprintf(`<div id="%s" data-loading-boundary="false">%s</div>`, lb.ID, result)
}

// MarkLoading sets the loading state for a specific child
func (lb *LoadingBoundary) MarkLoading(id string, loading bool) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	if lb.LoadingStates == nil {
		lb.LoadingStates = make(map[string]bool)
	}
	lb.LoadingStates[id] = loading
}

// TransitionBoundary manages animated transitions between states
type TransitionBoundary struct {
	ID       string
	Children []Component
	Key      interface{}
}

// Render renders the current children within a transition container
func (tb *TransitionBoundary) Render() string {
	var childrenHTML string
	for _, child := range tb.Children {
		childrenHTML += child.Render()
	}
	return fmt.Sprintf(`<div id="%s" data-transition-key="%v">%s</div>`, tb.ID, tb.Key, childrenHTML)
}

// ErrorInfo returns JSON-serialized error information
func (eb *ErrorBoundary) ErrorInfo() string {
	if eb.errorInfo == nil {
		return "{}"
	}
	info := map[string]string{
		"error": eb.errorInfo.Error(),
	}
	b, _ := json.Marshal(info)
	return string(b)
}
