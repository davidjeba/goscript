package goscript

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// HydrationData represents data to be sent to the client for hydration. It
// holds three independent namespaces — State (global application state),
// Props (per-component properties), and Config (app configuration) — that
// are serialized into a single window.__GOSCRIPT_HYDRATION__ object.
//
// The struct is safe for concurrent use; all map mutations are protected by
// a sync.RWMutex.
type HydrationData struct {
	State  map[string]interface{} `json:"state,omitempty"`
	Props  map[string]interface{} `json:"props,omitempty"`
	Config map[string]interface{} `json:"config,omitempty"`
	mu     sync.RWMutex
}

// NewHydrationData creates a new HydrationData with all maps initialized
// and ready for use.
func NewHydrationData() *HydrationData {
	return &HydrationData{
		State:  make(map[string]interface{}),
		Props:  make(map[string]interface{}),
		Config: make(map[string]interface{}),
	}
}

// Set adds a key-value pair to the state namespace. State values are
// typically global application state that the client runtime should
// initialize with (e.g. user session, feature flags, counters).
// The method returns the receiver for chaining.
//
// Usage:
//
//	hydration.Set("user", "alice").Set("theme", "dark")
func (h *HydrationData) Set(key string, value interface{}) *HydrationData {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.State[key] = value
	return h
}

// SetProps adds component props for hydration. The componentName is used as
// the key in the props namespace so that the client runtime can look up
// props per-component.
//
// Usage:
//
//	hydration.SetProps("Navbar", map[string]interface{}{"isLoggedIn": true})
func (h *HydrationData) SetProps(componentName string, props map[string]interface{}) *HydrationData {
	h.mu.Lock()
	defer h.mu.Unlock()
	if props != nil {
		h.Props[componentName] = props
	}
	return h
}

// SetConfig sets a single app configuration value. Config is distinct from
// state — it represents application-wide settings that don't change during
// the session (e.g. API base URL, environment, feature flags).
//
// Usage:
//
//	hydration.SetConfig("apiBase", "https://api.example.com")
func (h *HydrationData) SetConfig(key string, value interface{}) *HydrationData {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Config[key] = value
	return h
}

// GetState retrieves a value from the state namespace.
func (h *HydrationData) GetState(key string) interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.State[key]
}

// GetProps retrieves the props map for a specific component.
func (h *HydrationData) GetProps(componentName string) interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.Props[componentName]
}

// GetConfig retrieves a configuration value.
func (h *HydrationData) GetConfig(key string) interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.Config[key]
}

// Render generates the <script> tag that embeds the hydration data into the
// HTML page. The data is serialized as JSON and assigned to
// window.__GOSCRIPT_HYDRATION__. Empty namespaces are omitted from the
// output to keep the payload minimal.
//
// The returned string looks like:
//
//	<script>window.__GOSCRIPT_HYDRATION__={"state":{...},"props":{...},"config":{...}}</script>
//
// Usage:
//
//	hydration := NewHydrationData().Set("count", 0)
//	html := hydration.Render()
func (h *HydrationData) Render() string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Build a clean output map that omits empty namespaces
	out := make(map[string]interface{})
	if len(h.State) > 0 {
		out["state"] = h.State
	}
	if len(h.Props) > 0 {
		out["props"] = h.Props
	}
	if len(h.Config) > 0 {
		out["config"] = h.Config
	}

	if len(out) == 0 {
		return `<script>window.__GOSCRIPT_HYDRATION__={}</script>`
	}

	dataJSON, err := json.Marshal(out)
	if err != nil {
		dataJSON = []byte("{}")
	}

	return fmt.Sprintf(`<script>window.__GOSCRIPT_HYDRATION__=%s</script>`, string(dataJSON))
}

// RenderModule returns the hydration data as a JavaScript module-compatible
// inline script. Unlike Render() which sets a global variable, this version
// uses Object.assign to merge into an existing global if present, making it
// safe for code-split or multi-entry scenarios.
//
// Usage:
//
//	html := hydration.RenderModule()
//	// Output: <script>(window.__GOSCRIPT_HYDRATION__=window.__GOSCRIPT_HYDRATION__||{}),Object.assign(window.__GOSCRIPT_HYDRATION__,{"state":{...}})</script>
func (h *HydrationData) RenderModule() string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	out := make(map[string]interface{})
	if len(h.State) > 0 {
		out["state"] = h.State
	}
	if len(h.Props) > 0 {
		out["props"] = h.Props
	}
	if len(h.Config) > 0 {
		out["config"] = h.Config
	}

	if len(out) == 0 {
		return `<script>(window.__GOSCRIPT_HYDRATION__=window.__GOSCRIPT_HYDRATION__||{})</script>`
	}

	dataJSON, err := json.Marshal(out)
	if err != nil {
		dataJSON = []byte("{}")
	}

	return fmt.Sprintf(
		`<script>(window.__GOSCRIPT_HYDRATION__=window.__GOSCRIPT_HYDRATION__||{}),Object.assign(window.__GOSCRIPT_HYDRATION__,%s)</script>`,
		string(dataJSON),
	)
}

// RenderJSON returns the raw JSON string of the hydration data without the
// surrounding <script> tags. This is useful when embedding hydration data
// in a JSON response or in a data attribute.
func (h *HydrationData) RenderJSON() string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	out := make(map[string]interface{})
	if len(h.State) > 0 {
		out["state"] = h.State
	}
	if len(h.Props) > 0 {
		out["props"] = h.Props
	}
	if len(h.Config) > 0 {
		out["config"] = h.Config
	}

	dataJSON, err := json.Marshal(out)
	if err != nil {
		return "{}"
	}
	return string(dataJSON)
}

// Merge combines another HydrationData instance into this one. Values from
// other overwrite existing values with the same keys. The receiver is
// returned for chaining.
//
// Usage:
//
//	hydration.Merge(other).Set("extra", "value")
func (h *HydrationData) Merge(other *HydrationData) *HydrationData {
	if other == nil {
		return h
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	other.mu.RLock()
	defer other.mu.RUnlock()

	// Merge state
	for k, v := range other.State {
		h.State[k] = v
	}

	// Merge props
	for k, v := range other.Props {
		h.Props[k] = v
	}

	// Merge config
	for k, v := range other.Config {
		h.Config[k] = v
	}

	return h
}

// MergeState adds all key-value pairs from the given map into the state
// namespace.
func (h *HydrationData) MergeState(state map[string]interface{}) *HydrationData {
	if state == nil {
		return h
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for k, v := range state {
		h.State[k] = v
	}
	return h
}

// MergeConfig adds all key-value pairs from the given map into the config
// namespace.
func (h *HydrationData) MergeConfig(config map[string]interface{}) *HydrationData {
	if config == nil {
		return h
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for k, v := range config {
		h.Config[k] = v
	}
	return h
}

// Clear removes all hydration data from all namespaces.
func (h *HydrationData) Clear() *HydrationData {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.State = make(map[string]interface{})
	h.Props = make(map[string]interface{})
	h.Config = make(map[string]interface{})
	return h
}

// IsEmpty returns true if all three namespaces are empty.
func (h *HydrationData) IsEmpty() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.State) == 0 && len(h.Props) == 0 && len(h.Config) == 0
}

// HasState returns true if the state namespace is non-empty.
func (h *HydrationData) HasState() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.State) > 0
}

// HasProps returns true if the props namespace is non-empty.
func (h *HydrationData) HasProps() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.Props) > 0
}

// HasConfig returns true if the config namespace is non-empty.
func (h *HydrationData) HasConfig() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.Config) > 0
}

// FromStore creates a HydrationData pre-populated with the Store's current
// state. This provides a convenient bridge from the server-side Store to
// client-side hydration.
//
// Usage:
//
//	hydration := FromStore(myStore)
//	html := hydration.Render()
func FromStore(store *Store) *HydrationData {
	h := NewHydrationData()
	if store == nil {
		return h
	}
	store.mu.RLock()
	defer store.mu.RUnlock()
	for k, v := range store.state {
		h.State[k] = v
	}
	return h
}

// RenderHydrationBundle generates the complete set of script tags needed for
// a typical SSR page: first the hydration data, then the runtime. This is a
// convenience function that combines HydrationData.Render() with
// RenderGoScriptTag().
//
// Usage:
//
//	hydration := NewHydrationData().Set("user", user)
//	html := RenderHydrationBundle(hydration)
//	// Output: <script>window.__GOSCRIPT_HYDRATION__=...</script>\n<script src="/__goscript/runtime.js"></script>
func RenderHydrationBundle(h *HydrationData) string {
	if h == nil {
		return RenderGoScriptTag()
	}
	var sb strings.Builder
	sb.WriteString(h.Render())
	sb.WriteString("\n")
	sb.WriteString(RenderGoScriptTag())
	return sb.String()
}
