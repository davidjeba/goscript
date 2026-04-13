package main

import "goscript/dom"
import "goscript/state"
import "goscript/fmt"

// StoreDemo — demonstrates goscript's global state store pattern.
// Shows how multiple components can share state, subscribe to changes,
// dispatch actions, and hydrate initial state from the server.
//
// Usage:
//
//	gopm compile store_demo.gs -o store_demo.js
//	<script src="/store_demo.js"></script>

// AppState represents the global application state.
struct AppState struct {
    Count      int
    Username   string
    Theme      string
    LastAction string
}

// Action represents a store action.
struct StoreAction struct {
    Type  string
    Value interface{}
}

// --- Store setup ---
// storeReducer handles all state transitions for the global store.
func storeReducer(state AppState, action StoreAction) AppState {
    switch action.Type {
    case "INCREMENT":
        return AppState{
            Count:      state.Count + 1,
            Username:   state.Username,
            Theme:      state.Theme,
            LastAction: "increment",
        }
    case "DECREMENT":
        return AppState{
            Count:      state.Count - 1,
            Username:   state.Username,
            Theme:      state.Theme,
            LastAction: "decrement",
        }
    case "SET_USERNAME":
        return AppState{
            Count:      state.Count,
            Username:   action.Value.(string),
            Theme:      state.Theme,
            LastAction: "set_username",
        }
    case "SET_THEME":
        return AppState{
            Count:      state.Count,
            Username:   state.Username,
            Theme:      action.Value.(string),
            LastAction: "set_theme",
        }
    case "RESET":
        return AppState{Count: 0, Username: "Guest", Theme: "light", LastAction: "reset"}
    default:
        return state
    }
}

// CreateStore initializes the global state store with the given reducer
// and initial state.
func CreateStore(reducer func(AppState, StoreAction) AppState, initial AppState) map[string]interface{} {
    store := map[string]interface{}{
        "state":    initial,
        "reducer":  reducer,
        "listeners": []func(AppState){},
    }
    return store
}

// StoreGet returns the current state from the store.
func StoreGet(store map[string]interface{}) AppState {
    return store["state"].(AppState)
}

// StoreDispatch sends an action through the reducer and notifies listeners.
func StoreDispatch(store map[string]interface{}, action StoreAction) {
    reducer := store["reducer"].(func(AppState, StoreAction) AppState)
    current := store["state"].(AppState)
    newState := reducer(current, action)
    store["state"] = newState

    // Notify all subscribed listeners
    listeners := store["listeners"].([]func(AppState))
    for _, listener := range listeners {
        listener(newState)
    }
}

// StoreSubscribe adds a listener that fires on every state change.
// Returns an unsubscribe function.
func StoreSubscribe(store map[string]interface{}, listener func(AppState)) func() {
    listeners := store["listeners"].([]func(AppState))
    listeners = append(listeners, listener)
    store["listeners"] = listeners

    // Return unsubscribe function
    return func() {
        remaining := []func(AppState){}
        for _, l := range listeners {
            if fmt.Sprintf("%p", l) != fmt.Sprintf("%p", listener) {
                remaining = append(remaining, l)
            }
        }
        store["listeners"] = remaining
    }
}

// StoreHydrate merges server-provided state into the store.
func StoreHydrate(store map[string]interface{}, serverState AppState) {
    StoreDispatch(store, StoreAction{Type: "HYDRATE", Value: serverState})
}

// --- Components ---

// CounterPanel reads count from the store and renders increment/decrement.
func CounterPanel(store map[string]interface{}) dom.Element {
    count, setCount := state.Use(0)
    lastAction, setLastAction := state.Use("none")

    // Subscribe to store changes — update local state when store changes
    StoreSubscribe(store, func(newState AppState) {
        setCount(newState.Count)
        setLastAction(newState.LastAction)
    })

    increment := func(e dom.Event) {
        StoreDispatch(store, StoreAction{Type: "INCREMENT"})
    }

    decrement := func(e dom.Event) {
        StoreDispatch(store, StoreAction{Type: "DECREMENT"})
    }

    return dom.CreateElement("div", dom.Props{"class": "panel"},
        dom.CreateElement("h3", nil, "Counter Panel"),
        dom.CreateElement("div", dom.Props{"class": "count-display"}, fmt.Sprintf("%d", count)),
        dom.CreateElement("p", dom.Props{"class": "action-label"},
            fmt.Sprintf("Last action: %s", lastAction),
        ),
        dom.CreateElement("div", dom.Props{"class": "button-group"},
            dom.CreateElement("button", dom.Props{"class": "btn", "onclick": decrement}, "− 1"),
            dom.CreateElement("button", dom.Props{"class": "btn", "onclick": increment}, "+ 1"),
        ),
    )
}

// ProfilePanel reads username from the store.
func ProfilePanel(store map[string]interface{}) dom.Element {
    username, setUsername := state.Use("")

    StoreSubscribe(store, func(newState AppState) {
        setUsername(newState.Username)
    })

    return dom.CreateElement("div", dom.Props{"class": "panel"},
        dom.CreateElement("h3", nil, "Profile Panel"),
        dom.CreateElement("p", nil, fmt.Sprintf("Hello, %s!", username)),
        dom.CreateElement("input", dom.Props{
            "type":        "text",
            "placeholder": "Change username",
            "class":       "input",
            "onchange": func(e dom.Event) {
                StoreDispatch(store, StoreAction{Type: "SET_USERNAME", Value: dom.ValueOf(e.Target())})
            },
        }),
    )
}

// ThemePanel reads theme from the store and provides toggle buttons.
func ThemePanel(store map[string]interface{}) dom.Element {
    theme, setTheme := state.Use("light")

    StoreSubscribe(store, func(newState AppState) {
        setTheme(newState.Theme)
    })

    setThemeAction := func(t string) func(dom.Event) {
        return func(e dom.Event) {
            StoreDispatch(store, StoreAction{Type: "SET_THEME", Value: t})
            dom.SetAttribute("body", "data-theme", t)
        }
    }

    return dom.CreateElement("div", dom.Props{"class": "panel"},
        dom.CreateElement("h3", nil, "Theme Panel"),
        dom.CreateElement("p", nil, fmt.Sprintf("Current theme: %s", theme)),
        dom.CreateElement("div", dom.Props{"class": "button-group"},
            dom.CreateElement("button", dom.Props{"class": "btn", "onclick": setThemeAction("light")}, "☀️ Light"),
            dom.CreateElement("button", dom.Props{"class": "btn", "onclick": setThemeAction("dark")}, "🌙 Dark"),
        ),
    )
}

// StoreDemoApp composes all panels and initializes the store.
func StoreDemoApp() dom.Element {
    // --- Initialize store with default state ---
    store := CreateStore(storeReducer, AppState{
        Count:      0,
        Username:   "Guest",
        Theme:      "light",
        LastAction: "init",
    })

    // --- Hydrate from server state if available ---
    // In production, the server injects __GS_HYDRATE__ into the page.
    serverState := dom.GetHydrationData()
    if serverState != nil {
        hydrated := AppState{
            Count:      serverState["count"].(int),
            Username:   serverState["username"].(string),
            Theme:      serverState["theme"].(string),
            LastAction: "hydrated",
        }
        StoreHydrate(store, hydrated)
    }

    resetAll := func(e dom.Event) {
        StoreDispatch(store, StoreAction{Type: "RESET"})
    }

    return dom.CreateElement("div", dom.Props{"class": "store-demo"},
        dom.CreateElement("h2", nil, "🗄️ Goscript State Store"),
        dom.CreateElement("p", dom.Props{"class": "subtitle"}, "Global store with shared state across components"),

        // All three panels share the same store
        dom.CreateElement("div", dom.Props{"class": "panels-grid"},
            CounterPanel(store),
            ProfilePanel(store),
            ThemePanel(store),
        ),

        // Global reset
        dom.CreateElement("div", dom.Props{"class": "demo-footer"},
            dom.CreateElement("button", dom.Props{
                "class":   "btn danger",
                "onclick": resetAll,
            }, "Reset All State"),
        ),
    )
}

// Mount the store demo into #app.
func main() {
    dom.Mount("#app", StoreDemoApp())
}
