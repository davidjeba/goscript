package main

import "goscript/dom"
import "goscript/state"

// Counter — a goscript client component
// This .gs file compiles to JavaScript via: gopm compile counter.gs -o counter.js
// The compiled JavaScript uses the goscript runtime (__gs) for DOM and state.
//
// Usage:
//
//	gopm compile counter.gs -o counter.js
//
// Then include counter.js in your HTML:
//
//	<script src="/counter.js"></script>

// Counter returns a reactive counter DOM element.
// It uses goscript's built-in state management to track the count
// and automatically re-render when the state changes.
func Counter() dom.Element {
    count, setCount := state.Use(0)

    return dom.CreateElement("div", dom.Props{
        "class": "gs-counter",
    },
        dom.CreateElement("h2", nil, "Goscript .GS Counter"),
        dom.CreateElement("p", dom.Props{
            "class": "subtitle",
        }, "Client-side component compiled from .gs"),
        dom.CreateElement("div", dom.Props{
            "id":    "gs-count",
            "class": "count-display",
        }, count),
        dom.CreateElement("div", dom.Props{
            "class": "button-group",
        },
            dom.CreateElement("button", dom.Props{
                "class":   "btn decrement",
                "onclick": func(e dom.Event) {
                    setCount(count - 1)
                },
            }, "− Decrement"),
            dom.CreateElement("button", dom.Props{
                "class":   "btn increment",
                "onclick": func(e dom.Event) {
                    setCount(count + 1)
                },
            }, "+ Increment"),
            dom.CreateElement("button", dom.Props{
                "class":   "btn reset",
                "onclick": func(e dom.Event) {
                    setCount(0)
                },
            }, "Reset"),
        ),
    )
}

// Mount the counter component into the #app element when the page loads.
func main() {
    dom.Mount("#app", Counter())
}
