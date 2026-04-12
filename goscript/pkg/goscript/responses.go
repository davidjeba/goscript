// Package goscript — responses.go provides Go-side helpers for building HTTP
// responses that the goscript client runtime (__gs) interprets. These functions
// set GS-* response headers that the runtime's patched fetch handler reads to
// perform DOM swaps, fire events, and sync state — all without custom JavaScript.
//
// The primary use case is responding to reactive attribute requests (e.g. a
// gs-trigger="click /api/data" element sends a request, and the handler uses
// GoscriptResponse to return HTML fragments with swap instructions).
package goscript

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// =========================================================================
// Response Configuration
// =========================================================================

// responseConfig holds all options for a goscript response. It is built
// incrementally via ResponseOption functions.
type responseConfig struct {
	retarget  string                 // Override the target element for this swap
	reswap    string                 // Override the swap strategy
	triggers  []triggerEntry         // Client-side events to fire after swap
	pushURL   string                 // URL to push into browser history
	state     map[string]interface{} // State values to sync to the client
	retrigger bool                   // Whether to reprocess reactive elements
}

// triggerEntry represents a client-side event to fire.
type triggerEntry struct {
	Name   string      `json:"name"`   // Event name (e.g. "cart:updated")
	Detail interface{} `json:"detail"` // Event payload
}

// =========================================================================
// Response Options (Functional Options Pattern)
// =========================================================================

// ResponseOption is a function that modifies a responseConfig. Used with
// GoscriptResponse to customize response behavior.
type ResponseOption func(*responseConfig)

// WithRetarget overrides the target element for this swap. When set, the
// client runtime will swap the response into this element instead of the
// one specified in the originating element's gs-target attribute.
//
// Usage:
//
//	GoscriptResponse(w, "<div>Updated!</div>", WithRetarget("#sidebar"))
func WithRetarget(target string) ResponseOption {
	return func(c *responseConfig) {
		c.retarget = target
	}
}

// WithReswap overrides the swap strategy for this response. When set, the
// client runtime will use this strategy instead of the one specified in
// the originating element's gs-swap attribute.
//
// Usage:
//
//	GoscriptResponse(w, html, WithReswap(SwapOuterHTML))
func WithReswap(strategy string) ResponseOption {
	return func(c *responseConfig) {
		c.reswap = strategy
	}
}

// WithTrigger fires a client-side event after the swap is complete. The
// event is dispatched on the client's event bus. Multiple WithTrigger
// options can be combined.
//
// Usage:
//
//	GoscriptResponse(w, html, WithTrigger("item:added", itemID))
func WithTrigger(eventName string, detail interface{}) ResponseOption {
	return func(c *responseConfig) {
		c.triggers = append(c.triggers, triggerEntry{
			Name:   eventName,
			Detail: detail,
		})
	}
}

// WithPushURL pushes a URL into the browser's history after the swap.
// This enables SPA-like navigation without JavaScript.
//
// Usage:
//
//	GoscriptResponse(w, html, WithPushURL("/users/42"))
func WithPushURL(url string) ResponseOption {
	return func(c *responseConfig) {
		c.pushURL = url
	}
}

// WithState syncs a single server-side state value to the client runtime.
// The key-value pair is included in the GS-State response header.
//
// Usage:
//
//	GoscriptResponse(w, html, WithState("cartCount", 3))
func WithState(key string, value interface{}) ResponseOption {
	return func(c *responseConfig) {
		if c.state == nil {
			c.state = make(map[string]interface{})
		}
		c.state[key] = value
	}
}

// WithStateMap syncs multiple state values to the client runtime at once.
// The map is serialized into the GS-State response header.
//
// Usage:
//
//	GoscriptResponse(w, html, WithStateMap(map[string]interface{}{
//	    "cartCount": 3,
//	    "userLevel": "premium",
//	}))
func WithStateMap(state map[string]interface{}) ResponseOption {
	return func(c *responseConfig) {
		if c.state == nil {
			c.state = make(map[string]interface{})
		}
		for k, v := range state {
			c.state[k] = v
		}
	}
}

// WithRetrigger tells the client runtime to reprocess reactive elements
// in the swapped content. Set this to true when the response contains
// new elements with gs-trigger attributes.
//
// Usage:
//
//	GoscriptResponse(w, html, WithRetrigger(true))
func WithRetrigger(retrigger bool) ResponseOption {
	return func(c *responseConfig) {
		c.retrigger = retrigger
	}
}

// =========================================================================
// Primary Response Functions
// =========================================================================

// GoscriptResponse sends an HTML fragment response with goscript response
// headers. The client runtime reads these headers and performs the
// appropriate DOM operations (swap, trigger events, sync state, etc.).
//
// The html parameter is the raw HTML fragment to swap into the target
// element. Options can override the target, swap strategy, and more.
//
// This is the most commonly used response function for reactive handlers.
//
// Usage:
//
//	func handleUserList(w http.ResponseWriter, r *http.Request) {
//	    users := getUserList()
//	    html := renderUserList(users)
//	    GoscriptResponse(w, html, WithRetarget("#user-list"))
//	}
func GoscriptResponse(w http.ResponseWriter, html string, opts ...ResponseOption) {
	// Build config from options
	cfg := &responseConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	// Set content type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Set GS-Retarget if overridden
	if cfg.retarget != "" {
		w.Header().Set("GS-Retarget", cfg.retarget)
	}

	// Set GS-Reswap if overridden
	if cfg.reswap != "" {
		w.Header().Set("GS-Reswap", cfg.reswap)
	}

	// Set GS-Trigger if events are specified
	if len(cfg.triggers) > 0 {
		triggerJSON, err := json.Marshal(cfg.triggers)
		if err == nil {
			w.Header().Set("GS-Trigger", string(triggerJSON))
		}
	}

	// Set GS-Push-URL if specified
	if cfg.pushURL != "" {
		w.Header().Set("GS-Push-Url", cfg.pushURL)
	}

	// Set GS-State if state values are provided
	if len(cfg.state) > 0 {
		stateJSON, err := json.Marshal(cfg.state)
		if err == nil {
			w.Header().Set("GS-State", string(stateJSON))
		}
	}

	// Set GS-Retrigger header
	if cfg.retrigger {
		w.Header().Set("GS-Retrigger", "true")
	}

	// Write the HTML body
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, html)
}

// GoscriptRedirect sends a client-side redirect response. The client
// runtime will navigate to the specified URL without a full page reload.
//
// Usage:
//
//	GoscriptRedirect(w, "/login")
func GoscriptRedirect(w http.ResponseWriter, url string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("GS-Redirect", url)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "")
}

// GoscriptRefresh sends a response that triggers a full page refresh
// on the client side.
//
// Usage:
//
//	GoscriptRefresh(w)
func GoscriptRefresh(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("GS-Refresh", "true")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "")
}

// GoscriptTrigger sends a trigger-only response (no content swap).
// The client runtime will fire the specified event on the event bus
// but will not modify any DOM elements.
//
// Usage:
//
//	GoscriptTrigger(w, "notification", map[string]string{
//	    "message": "Data saved successfully",
//	    "type":    "success",
//	})
func GoscriptTrigger(w http.ResponseWriter, event string, detail interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Build trigger entry and marshal to JSON
	trigger := triggerEntry{
		Name:   event,
		Detail: detail,
	}
	triggerJSON, err := json.Marshal(trigger)
	if err == nil {
		w.Header().Set("GS-Trigger", string(triggerJSON))
	}

	w.WriteHeader(http.StatusNoContent)
}

// =========================================================================
// Batch Response Builder
// =========================================================================

// BatchResponse allows updating multiple DOM elements from a single server
// response. Use it when a single action needs to update several parts of
// the page simultaneously.
//
// The batch response uses a special JSON format in the response body that
// the client runtime parses and applies sequentially.
//
// Usage:
//
//	batch := NewBatchResponse(w)
//	batch.Swap("#user-list", renderUsers(users), SwapInnerHTML).
//	    Swap("#user-count", fmt.Sprintf("%d users", len(users)), SwapInnerHTML).
//	    Trigger("users:loaded", nil).
//	    Send()
type BatchResponse struct {
	writer   http.ResponseWriter
	swaps    []swapEntry
	triggers []triggerEntry
	state    map[string]interface{}
	pushURL  string
}

// swapEntry represents a single DOM swap operation in a batch.
type swapEntry struct {
	Target   string `json:"target"`   // CSS selector for the swap target
	Content  string `json:"content"`  // HTML content to swap
	Strategy string `json:"strategy"` // Swap strategy (innerHTML, outerHTML, etc.)
}

// batchResponseBody is the JSON structure sent to the client.
type batchResponseBody struct {
	Swaps   []swapEntry            `json:"swaps"`
	Triggers []triggerEntry        `json:"triggers,omitempty"`
	State    map[string]interface{} `json:"state,omitempty"`
	PushURL  string                `json:"pushUrl,omitempty"`
}

// NewBatchResponse creates a new batch response builder for the given
// HTTP response writer.
func NewBatchResponse(w http.ResponseWriter) *BatchResponse {
	return &BatchResponse{
		writer: w,
		swaps:  make([]swapEntry, 0),
		state:  make(map[string]interface{}),
	}
}

// Swap adds a DOM swap to the batch. The response HTML will be swapped
// into the target element using the specified strategy.
//
// Usage:
//
//	batch.Swap("#user-list", userHTML, SwapInnerHTML)
func (b *BatchResponse) Swap(target, content, strategy string) *BatchResponse {
	b.swaps = append(b.swaps, swapEntry{
		Target:   target,
		Content:  content,
		Strategy: strategy,
	})
	return b
}

// SwapComponent targets a component by name using the component:name
// selector format. This is a convenience wrapper around Swap.
//
// Usage:
//
//	batch.SwapComponent("navbar", navHTML, SwapInnerHTML)
func (b *BatchResponse) SwapComponent(name, content, strategy string) *BatchResponse {
	return b.Swap(fmt.Sprintf("[data-goscript-component=\"%s\"]", name), content, strategy)
}

// Trigger adds a client-side event trigger to the batch. The event will
// be dispatched on the client's event bus after all swaps are applied.
//
// Usage:
//
//	batch.Trigger("cart:updated", map[string]int{"count": newCount})
func (b *BatchResponse) Trigger(event string, detail interface{}) *BatchResponse {
	b.triggers = append(b.triggers, triggerEntry{
		Name:   event,
		Detail: detail,
	})
	return b
}

// State syncs a state key-value pair to the client. The client runtime
// will update its local state store with the provided values.
//
// Usage:
//
//	batch.State("cartCount", 5)
func (b *BatchResponse) State(key string, value interface{}) *BatchResponse {
	b.state[key] = value
	return b
}

// PushURL sets a URL to push into the browser history after the batch
// is applied.
//
// Usage:
//
//	batch.PushURL("/products/42")
func (b *BatchResponse) PushURL(url string) *BatchResponse {
	b.pushURL = url
	return b
}

// Send serializes all batch operations into a JSON response and writes it
// to the HTTP response writer. The client runtime parses this JSON and
// applies each operation sequentially.
func (b *BatchResponse) Send() {
	body := batchResponseBody{
		Swaps:   b.swaps,
		Triggers: b.triggers,
		State:    b.state,
		PushURL:  b.pushURL,
	}

	// Remove empty fields
	if len(body.Triggers) == 0 {
		body.Triggers = nil
	}
	if len(body.State) == 0 {
		body.State = nil
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		http.Error(b.writer, "failed to serialize batch response", http.StatusInternalServerError)
		return
	}

	b.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	b.writer.Header().Set("GS-Batch", "true")
	b.writer.WriteHeader(http.StatusOK)
	b.writer.Write(bodyJSON)
}

// =========================================================================
// Runtime Script Tag
// =========================================================================

// RenderGoScriptTag generates the <script> tag that includes the goscript
// client runtime. This tag should be placed in the <head> or before the
// closing </body> tag of the HTML document.
//
// The actual JavaScript file (runtime.js) is embedded in the Go binary
// via go:embed and served automatically by the framework at
// /__goscript/runtime.js.
//
// Usage:
//
//	fmt.Println(RenderGoScriptTag())
//	// Output: <script src="/__goscript/runtime.js"></script>
func RenderGoScriptTag() string {
	return `<script src="/__goscript/runtime.js"></script>`
}

// RenderGoScriptTagAsync generates an async/deferred script tag for the
// goscript runtime. Using async prevents the runtime from blocking HTML
// parsing, improving initial page load performance.
//
// Usage:
//
//	fmt.Println(RenderGoScriptTagAsync())
//	// Output: <script src="/__goscript/runtime.js" async></script>
func RenderGoScriptTagAsync() string {
	return `<script src="/__goscript/runtime.js" async></script>`
}

// RenderGoScriptTagDefer generates a deferred script tag for the goscript
// runtime. Using defer ensures the runtime loads after HTML parsing but
// before DOMContentLoaded fires.
//
// Usage:
//
//	fmt.Println(RenderGoScriptTagDefer())
//	// Output: <script src="/__goscript/runtime.js" defer></script>
func RenderGoScriptTagDefer() string {
	return `<script src="/__goscript/runtime.js" defer></script>`
}

// RenderHydrationScript generates the script tag that embeds server-side
// state for client hydration. The state parameter is serialized as JSON
// and assigned to window.__GOSCRIPT_STATE__.
//
// This should be placed before the runtime script tag so that state is
// available when the runtime initializes.
//
// Usage:
//
//	state := map[string]interface{}{"count": 0, "user": "alice"}
//	fmt.Println(RenderHydrationScript(state))
//	// Output: <script>window.__GOSCRIPT_STATE__={"count":0,"user":"alice"}</script>
func RenderHydrationScript(state map[string]interface{}) string {
	if state == nil {
		state = make(map[string]interface{})
	}
	stateJSON, err := json.Marshal(state)
	if err != nil {
		stateJSON = []byte("{}")
	}
	return fmt.Sprintf(`<script>window.__GOSCRIPT_STATE__=%s</script>`, string(stateJSON))
}

// =========================================================================
// Utility Functions
// =========================================================================

// SetGSHeader is a low-level helper that sets a GS-* response header.
// Use this for custom headers that don't have dedicated helper functions.
//
// Usage:
//
//	SetGSHeader(w, "GS-Custom", "some-value")
func SetGSHeader(w http.ResponseWriter, key, value string) {
	w.Header().Set(key, value)
}

// IsGSRequest checks whether an incoming request was made by the goscript
// client runtime. Requests from the runtime include the GS-Request header.
//
// Usage:
//
//	if IsGSRequest(r) {
//	    // This is a reactive request — return HTML fragment
//	} else {
//	    // This is a normal request — return full page
//	}
func IsGSRequest(r *http.Request) bool {
	return r.Header.Get("GS-Request") == "true"
}

// WantsHTML checks whether the client prefers an HTML response. This is
// useful for handlers that serve both JSON API and HTML reactive responses.
//
// Usage:
//
//	if WantsHTML(r) {
//	    GoscriptResponse(w, renderHTML(data))
//	} else {
//	    json.NewEncoder(w).Encode(data)
//	}
func WantsHTML(r *http.Request) bool {
	// Check GS-Request header (set by runtime)
	if IsGSRequest(r) {
		return true
	}
	// Check Accept header
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, "text/html")
}
