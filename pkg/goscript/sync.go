package goscript

import (
	"sync"
	"time"
)

// SyncAction describes a queued offline-first update.
type SyncAction struct {
	ID        string                 `json:"id"`
	Scope     string                 `json:"scope"`
	Kind      string                 `json:"kind"`
	Payload   interface{}            `json:"payload"`
	CreatedAt time.Time              `json:"createdAt"`
	Meta      map[string]string      `json:"meta,omitempty"`
}

// SyncQueue stores actions until they can be flushed.
type SyncQueue struct {
	mu        sync.RWMutex
	items     []SyncAction
	listeners []func([]SyncAction)
}

// NewSyncQueue creates an empty queue.
func NewSyncQueue() *SyncQueue {
	return &SyncQueue{
		items: make([]SyncAction, 0),
	}
}

// Enqueue adds a new sync action.
func (q *SyncQueue) Enqueue(action SyncAction) {
	if action.CreatedAt.IsZero() {
		action.CreatedAt = time.Now().UTC()
	}

	q.mu.Lock()
	q.items = append(q.items, action)
	listeners := append([]func([]SyncAction){}, q.listeners...)
	snapshot := append([]SyncAction(nil), q.items...)
	q.mu.Unlock()

	for _, listener := range listeners {
		if listener != nil {
			go listener(snapshot)
		}
	}
}

// Drain returns all queued actions and clears the queue.
func (q *SyncQueue) Drain() []SyncAction {
	q.mu.Lock()
	defer q.mu.Unlock()

	out := append([]SyncAction(nil), q.items...)
	q.items = q.items[:0]
	return out
}

// Snapshot returns a copy of queued actions without clearing them.
func (q *SyncQueue) Snapshot() []SyncAction {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return append([]SyncAction(nil), q.items...)
}

// Subscribe registers a queue listener.
func (q *SyncQueue) Subscribe(listener func([]SyncAction)) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.listeners = append(q.listeners, listener)
}

