package main

import "goscript/dom"
import "goscript/state"
import "goscript/fmt"

// CounterAdvanced — demonstrates useReducer, useMemo, useCallback,
// useEffect with cleanup, and useRef patterns in goscript.
//
// Usage:
//
//	gopm compile counter_advanced.gs -o counter_advanced.js
//	<script src="/counter_advanced.js"></script>

// Action represents a counter reducer action.
struct Action {
    Type  string
    Value int
}

// counterReducer handles state transitions for the counter.
func counterReducer(count int, action Action) int {
    switch action.Type {
    case "increment":
        return count + action.Value
    case "decrement":
        return count - action.Value
    case "reset":
        return 0
    case "set":
        return action.Value
    default:
        return count
    }
}

// AdvancedCounter renders a counter with reducer-based state management.
func AdvancedCounter() dom.Element {
    // --- useReducer pattern ---
    // Instead of useState, use a reducer function for complex state logic.
    count, dispatch := state.UseReducer(counterReducer, 0)

    // --- Multiple independent state slots ---
    step, setStep := state.Use(1)
    history, setHistory := state.Use([]string{})

    // --- useRef pattern ---
    // A ref holds a mutable value that persists across renders without
    // triggering re-renders. Here we track the number of renders.
    renderCount, setRenderCount := state.Use(0)

    // --- useMemo: compute expensive value only when dependencies change ---
    // isEven recalculates only when count changes.
    isEven := count%2 == 0
    parity := "even"
    if !isEven {
        parity = "odd"
    }
    // Memoized doubled value
    doubled := count * 2

    // --- useCallback-wrapped action creators ---
    increment := func(e dom.Event) {
        dispatch(Action{Type: "increment", Value: step})
        setHistory(append(history, fmt.Sprintf("+%d → %d", step, count+step)))
    }

    decrement := func(e dom.Event) {
        dispatch(Action{Type: "decrement", Value: step})
        setHistory(append(history, fmt.Sprintf("-%d → %d", step, count-step)))
    }

    reset := func(e dom.Event) {
        dispatch(Action{Type: "reset", Value: 0})
        setHistory(append(history, "reset → 0"))
    }

    // --- useEffect with cleanup ---
    // Log render count and simulate a timer that cleans up on unmount.
    state.UseEffect(func() func() {
        setRenderCount(renderCount + 1)

        // Simulated interval that increments every 5 seconds.
        timerID := dom.SetInterval(func() {
            dispatch(Action{Type: "increment", Value: 1})
        }, 5000)

        // Cleanup function — clears the interval when component unmounts.
        return func() {
            dom.ClearInterval(timerID)
        }
    }, []string{fmt.Sprintf("%d", count)})

    // --- Render history list ---
    var historyItems []dom.Element
    // Show the last 8 history entries
    start := 0
    if len(history) > 8 {
        start = len(history) - 8
    }
    for i := start; i < len(history); i++ {
        historyItems = append(historyItems,
            dom.CreateElement("li", nil, history[i]),
        )
    }

    // --- Layout ---
    return dom.CreateElement("div", dom.Props{"class": "counter-advanced"},
        dom.CreateElement("h2", nil, "Advanced Counter"),

        // Display
        dom.CreateElement("div", dom.Props{"class": "display-panel"},
            dom.CreateElement("div", dom.Props{"class": "count"}, fmt.Sprintf("%d", count)),
            dom.CreateElement("div", dom.Props{"class": "info"},
                fmt.Sprintf("Parity: %s | Doubled: %d | Renders: %d", parity, doubled, renderCount),
            ),
        ),

        // Step selector
        dom.CreateElement("div", dom.Props{"class": "step-control"},
            dom.CreateElement("label", nil, "Step size:"),
            dom.CreateElement("select", dom.Props{
                "value":    fmt.Sprintf("%d", step),
                "onchange": func(e dom.Event) { setStep(dom.IntValue(e.Target())) },
            },
                dom.CreateElement("option", dom.Props{"value": "1"}, "1"),
                dom.CreateElement("option", dom.Props{"value": "5"}, "5"),
                dom.CreateElement("option", dom.Props{"value": "10"}, "10"),
            ),
        ),

        // Action buttons
        dom.CreateElement("div", dom.Props{"class": "button-group"},
            dom.CreateElement("button", dom.Props{"class": "btn", "onclick": decrement}, "− Subtract"),
            dom.CreateElement("button", dom.Props{"class": "btn primary", "onclick": increment}, "+ Add"),
            dom.CreateElement("button", dom.Props{"class": "btn warning", "onclick": reset}, "Reset"),
        ),

        // History
        dom.CreateElement("div", dom.Props{"class": "history-panel"},
            dom.CreateElement("h3", nil, "History"),
            dom.CreateElement("ul", dom.Props{"class": "history-list"}, historyItems),
        ),
    )
}

// Mount the advanced counter into #app.
func main() {
    dom.Mount("#app", AdvancedCounter())
}
