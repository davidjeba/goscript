package goscript

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// =========================================================================
// SSEWriter — Low-level Server-Sent Events writer
// =========================================================================

// SSEWriter wraps an http.ResponseWriter to send Server-Sent Events (SSE).
// All writes are serialized through a mutex so that concurrent goroutines
// can safely push events to the same client connection.
//
// Usage:
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//	    sse := NewSSEWriter(w)
//	    sse.WriteEvent("message", map[string]string{"text": "hello"})
//	    sse.Close()
//	}
type SSEWriter struct {
	w       http.ResponseWriter
	flusher http.Flusher
	mu      sync.Mutex
	closed  bool
}

// NewSSEWriter creates a new SSEWriter and sets the required headers for
// SSE (Content-Type: text/event-stream, Cache-Control: no-cache, etc.).
// It returns an error if the ResponseWriter does not implement http.Flusher,
// which is required for SSE streaming.
func NewSSEWriter(w http.ResponseWriter) (*SSEWriter, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("goscript: SSE requires a ResponseWriter that supports flushing")
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // Disable nginx buffering
	w.Header().Set("Access-Control-Allow-Origin", "*")

	return &SSEWriter{
		w:       w,
		flusher: flusher,
	}, nil
}

// WriteEvent sends a named event with data. The data is JSON-serialized.
// Each SSE event follows the format:
//
//	event: <name>
//	data: <json>
//
// Usage:
//
//	sse.WriteEvent("user:joined", map[string]interface{}{"name": "alice"})
func (s *SSEWriter) WriteEvent(event string, data interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return fmt.Errorf("goscript: SSE writer is closed")
	}

	var dataStr string
	switch v := data.(type) {
	case string:
		dataStr = v
	case []byte:
		dataStr = string(v)
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("goscript: SSE JSON marshal error: %w", err)
		}
		dataStr = string(jsonBytes)
	}

	// Write the SSE event in the standard format
	_, err := fmt.Fprintf(s.w, "event: %s\ndata: %s\n\n", event, dataStr)
	if err != nil {
		return fmt.Errorf("goscript: SSE write error: %w", err)
	}

	s.flusher.Flush()
	return nil
}

// Write sends a default event with data. This is equivalent to calling
// WriteEvent("message", data).
//
// Usage:
//
//	sse.Write("hello world")
func (s *SSEWriter) Write(data interface{}) error {
	return s.WriteEvent("message", data)
}

// WriteComment sends an SSE comment. Comments start with ':' and are ignored
// by the browser's EventSource API. They are useful as keep-alive heartbeats
// to prevent proxy timeouts.
//
// Usage:
//
//	sse.WriteComment("heartbeat")
func (s *SSEWriter) WriteComment(comment string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return fmt.Errorf("goscript: SSE writer is closed")
	}

	_, err := fmt.Fprintf(s.w, ": %s\n\n", comment)
	if err != nil {
		return fmt.Errorf("goscript: SSE comment write error: %w", err)
	}

	s.flusher.Flush()
	return nil
}

// WriteRaw writes raw SSE text without any processing. The caller is
// responsible for formatting the SSE protocol correctly (event:, data:,
// id:, etc.). This is useful for sending custom SSE fields like id: or retry:.
//
// Usage:
//
//	sse.WriteRaw("id: msg-123\nretry: 5000\ndata: hello\n\n")
func (s *SSEWriter) WriteRaw(raw string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return fmt.Errorf("goscript: SSE writer is closed")
	}

	_, err := fmt.Fprint(s.w, raw)
	if err != nil {
		return fmt.Errorf("goscript: SSE raw write error: %w", err)
	}

	s.flusher.Flush()
	return nil
}

// SetRetry instructs the client's EventSource to reconnect after the
// specified duration if the connection drops. The duration is in milliseconds.
//
// Usage:
//
//	sse.SetRetry(3000) // Reconnect after 3 seconds
func (s *SSEWriter) SetRetry(ms int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return fmt.Errorf("goscript: SSE writer is closed")
	}

	_, err := fmt.Fprintf(s.w, "retry: %d\n\n", ms)
	if err != nil {
		return fmt.Errorf("goscript: SSE retry write error: %w", err)
	}

	s.flusher.Flush()
	return nil
}

// Close sends a close event to signal the client that the stream is ending,
// then marks the writer as closed. The client runtime should handle the
// "__goscript_close__" event by cleaning up its EventSource connection.
func (s *SSEWriter) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true

	// Send a named close event so the client runtime can react
	_, err := fmt.Fprintf(s.w, "event: __goscript_close__\ndata: {}\n\n")
	if err != nil {
		return fmt.Errorf("goscript: SSE close error: %w", err)
	}

	s.flusher.Flush()
	return nil
}

// IsClosed returns whether the SSE writer has been closed.
func (s *SSEWriter) IsClosed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed
}

// =========================================================================
// EventHub — Pub/Sub for Server-Sent Events
// =========================================================================

// EventHub manages a publish/subscribe system for SSE topics. Multiple
// clients can subscribe to topics and receive events published to those
// topics. The hub is thread-safe and designed to be shared across goroutines.
//
// Usage:
//
//	hub := NewEventHub()
//
//	// Register the SSE handler on a route
//	router.GET("/events/:topic", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
//	    topic := params["topic"]
//	    hub.SSEHandler(topic)(w, r)
//	})
//
//	// Publish events from anywhere
//	hub.Publish("notifications", map[string]string{"msg": "new order"})
type EventHub struct {
	clients map[string]map[chan string]struct{}
	mu      sync.RWMutex
}

// NewEventHub creates a new EventHub with an initialized client map.
func NewEventHub() *EventHub {
	return &EventHub{
		clients: make(map[string]map[chan string]struct{}),
	}
}

// Subscribe adds a new client channel to a topic. The returned channel
// receives SSE-formatted strings whenever data is published to the topic.
// The caller should range over the channel to read events.
//
// Usage:
//
//	ch := hub.Subscribe("notifications")
//	for msg := range ch {
//	    sse.WriteRaw(msg)
//	}
func (h *EventHub) Subscribe(topic string) chan string {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[topic] == nil {
		h.clients[topic] = make(map[chan string]struct{})
	}

	// Use a buffered channel to avoid blocking publishers
	ch := make(chan string, 64)
	h.clients[topic][ch] = struct{}{}

	return ch
}

// Unsubscribe removes a client channel from a topic. If the channel is the
// last subscriber for the topic, the topic entry is removed from the hub.
// The channel is closed after removal.
func (h *EventHub) Unsubscribe(topic string, ch chan string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if subs, ok := h.clients[topic]; ok {
		delete(subs, ch)
		close(ch)
		if len(subs) == 0 {
			delete(h.clients, topic)
		}
	}
}

// Publish sends data to all subscribers of a specific topic. The data is
// JSON-serialized and formatted as an SSE event string. Returns the number
// of clients that received the event.
//
// Usage:
//
//	hub.Publish("notifications", map[string]string{"msg": "new order"})
func (h *EventHub) Publish(topic string, data interface{}) (int, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("goscript: EventHub JSON marshal error: %w", err)
	}

	return h.publishRaw(topic, fmt.Sprintf("data: %s\n\n", string(jsonBytes))), nil
}

// PublishEvent sends a named event with data to all subscribers of a topic.
// The output follows the SSE format with an explicit event: field.
//
// Usage:
//
//	hub.PublishEvent("notifications", "user:joined", map[string]string{"name": "alice"})
func (h *EventHub) PublishEvent(topic string, event string, data interface{}) (int, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("goscript: EventHub JSON marshal error: %w", err)
	}

	msg := fmt.Sprintf("event: %s\ndata: %s\n\n", event, string(jsonBytes))
	return h.publishRaw(topic, msg), nil
}

// publishRaw is the internal method that fans out a pre-formatted SSE message
// to all subscribers of a topic. Non-blocking sends are used to prevent slow
// clients from blocking the publisher.
func (h *EventHub) publishRaw(topic string, msg string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	subs, ok := h.clients[topic]
	if !ok {
		return 0
	}

	sent := 0
	for ch := range subs {
		select {
		case ch <- msg:
			sent++
		default:
			// Channel buffer full — skip this client to avoid blocking.
			// A monitoring system could track dropped messages.
		}
	}

	return sent
}

// Broadcast sends data to all subscribers across all topics. Returns the
// total number of clients that received the event.
//
// Usage:
//
//	hub.Broadcast(map[string]string{"alert": "system maintenance at midnight"})
func (h *EventHub) Broadcast(data interface{}) (int, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("goscript: EventHub JSON marshal error: %w", err)
	}

	msg := fmt.Sprintf("data: %s\n\n", string(jsonBytes))

	h.mu.RLock()
	defer h.mu.RUnlock()

	total := 0
	for _, subs := range h.clients {
		for ch := range subs {
			select {
			case ch <- msg:
				total++
			default:
				// Skip slow clients
			}
		}
	}

	return total, nil
}

// Topics returns a list of all active topic names.
func (h *EventHub) Topics() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	topics := make([]string, 0, len(h.clients))
	for topic := range h.clients {
		topics = append(topics, topic)
	}
	return topics
}

// SubscriberCount returns the number of subscribers for a specific topic.
func (h *EventHub) SubscriberCount(topic string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients[topic])
}

// TotalSubscribers returns the total number of subscribers across all topics.
func (h *EventHub) TotalSubscribers() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	count := 0
	for _, subs := range h.clients {
		count += len(subs)
	}
	return count
}

// SSEHandler returns an http.HandlerFunc that streams SSE events for a
// topic. The handler automatically handles client disconnection via
// request context cancellation, sends periodic heartbeat comments to
// keep the connection alive, and cleans up the subscription on disconnect.
//
// Usage:
//
//	router.GET("/events/:topic", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
//	    hub.SSEHandler(params["topic"])(w, r)
//	})
func (h *EventHub) SSEHandler(topic string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sse, err := NewSSEWriter(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Subscribe to the topic
		ch := h.Subscribe(topic)
		defer h.Unsubscribe(topic, ch)

		// Send initial retry instruction
		_ = sse.SetRetry(3000)

		// Heartbeat ticker to keep the connection alive
		heartbeat := time.NewTicker(30 * time.Second)
		defer heartbeat.Stop()

		// Context cancellation channel
		ctx := r.Context()
		ctxDone := ctx.Done()

		for {
			select {
			case <-ctxDone:
				// Client disconnected
				return

			case msg, ok := <-ch:
				if !ok {
					// Channel was closed (unsubscribed)
					return
				}
				if err := sse.WriteRaw(msg); err != nil {
					return
				}

			case <-heartbeat.C:
				// Send a keep-alive comment
				if err := sse.WriteComment("heartbeat"); err != nil {
					return
				}
			}
		}
	}
}

// SSETopicHandler returns a RouteHandler (compatible with the goscript
// Router) that streams SSE events for a topic. The topic is extracted from
// the "topic" route parameter.
//
// Usage:
//
//	router.GET("/events/:topic", hub.SSETopicHandler())
func (h *EventHub) SSETopicHandler() RouteHandler {
	return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		topic := params["topic"]
		if topic == "" {
			http.Error(w, "topic is required", http.StatusBadRequest)
			return
		}
		h.SSEHandler(topic)(w, r)
	}
}

// =========================================================================
// Helper functions for reactive SSE integration
// =========================================================================

// SSEStateChangeFormat formats a state change event that the goscript client
// runtime can interpret. The event name follows the convention "state:key"
// so the runtime can automatically dispatch it to the reactive state system.
//
// Usage:
//
//	hub.PublishEvent("state", SSEStateChangeFormat("user", userData))
func SSEStateChangeFormat(key string, value interface{}) string {
	return fmt.Sprintf("state:%s", key)
}

// SSEStateEventName returns the standard event name for state change
// notifications. The goscript client runtime listens for events matching
// this pattern to update its local state.
//
// Usage:
//
//	eventName := SSEStateEventName("cartCount")
//	hub.PublishEvent("state", eventName, 5)
func SSEStateEventName(key string) string {
	return fmt.Sprintf("state:%s", key)
}

// SSERenderEvent formats an SSE event string for a DOM swap that the client
// runtime can interpret. This is useful when the server wants to push
// HTML updates to specific elements without a client-initiated request.
//
// Usage:
//
//	event := fmt.Sprintf("event: swap\ndata: %s\n\n",
//	    SSERenderEvent("#user-list", "<li>New user</li>", "innerHTML"))
//	sse.WriteRaw(event)
func SSERenderEvent(target, html, swapStrategy string) string {
	swap := struct {
		Target   string `json:"target"`
		Content  string `json:"content"`
		Strategy string `json:"strategy"`
	}{
		Target:   target,
		Content:  html,
		Strategy: swapStrategy,
	}

	data, _ := json.Marshal(swap)
	return string(data)
}

// IsSSERequest checks whether a request is an SSE connection request by
// inspecting the Accept header for text/event-stream.
func IsSSERequest(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, "text/event-stream")
}

// SSEPreflightHandler returns an http.HandlerFunc that handles CORS preflight
// requests for SSE endpoints. Browsers send an OPTIONS request before
// establishing an SSE connection when the client and server are on different
// origins.
func SSEPreflightHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Cache-Control, Last-Event-ID")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusNoContent)
	}
}
