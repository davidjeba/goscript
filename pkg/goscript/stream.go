package goscript

import (
	"sync"
	"time"
)

// StreamEvent represents a realtime message.
type StreamEvent struct {
	Topic     string
	Kind      string
	Payload   interface{}
	Source    string
	Timestamp time.Time
}

// RealtimeHub publishes events to topic subscribers.
type RealtimeHub struct {
	mu          sync.RWMutex
	subscribers  map[string]map[string]chan StreamEvent
	history     map[string][]StreamEvent
	bufferSize  int
}

// NewRealtimeHub creates a hub for realtime events.
func NewRealtimeHub(bufferSize int) *RealtimeHub {
	if bufferSize < 1 {
		bufferSize = 16
	}

	return &RealtimeHub{
		subscribers: make(map[string]map[string]chan StreamEvent),
		history:     make(map[string][]StreamEvent),
		bufferSize:  bufferSize,
	}
}

// Subscribe registers a subscriber for a topic.
func (h *RealtimeHub) Subscribe(topic, subscriberID string) <-chan StreamEvent {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.subscribers[topic] == nil {
		h.subscribers[topic] = make(map[string]chan StreamEvent)
	}

	ch := make(chan StreamEvent, h.bufferSize)
	h.subscribers[topic][subscriberID] = ch
	return ch
}

// Unsubscribe removes a subscriber from a topic.
func (h *RealtimeHub) Unsubscribe(topic, subscriberID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if topicSubs, ok := h.subscribers[topic]; ok {
		if ch, ok := topicSubs[subscriberID]; ok {
			close(ch)
			delete(topicSubs, subscriberID)
		}
	}
}

// Publish emits an event to a topic.
func (h *RealtimeHub) Publish(topic, kind string, payload interface{}, source string) StreamEvent {
	event := StreamEvent{
		Topic:     topic,
		Kind:      kind,
		Payload:   payload,
		Source:    source,
		Timestamp: time.Now().UTC(),
	}

	h.mu.Lock()
	h.history[topic] = append(h.history[topic], event)
	subs := h.subscribers[topic]
	h.mu.Unlock()

	for _, ch := range subs {
		select {
		case ch <- event:
		default:
		}
	}

	return event
}

// Ping publishes a ping event.
func (h *RealtimeHub) Ping(topic, source string) StreamEvent {
	return h.Publish(topic, "ping", map[string]interface{}{
		"alive": true,
	}, source)
}

// Poll returns a snapshot of recent events.
func (h *RealtimeHub) Poll(topic string, limit int) []StreamEvent {
	h.mu.RLock()
	defer h.mu.RUnlock()

	events := h.history[topic]
	if limit <= 0 || limit >= len(events) {
		out := make([]StreamEvent, len(events))
		copy(out, events)
		return out
	}

	out := make([]StreamEvent, limit)
	copy(out, events[len(events)-limit:])
	return out
}

