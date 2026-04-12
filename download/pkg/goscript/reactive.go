// Package goscript — reactive.go provides Go-side helpers for building reactive
// components that the goscript client runtime (__gs) understands. These functions
// generate HTML attributes (gs-trigger, gs-target, gs-swap, etc.) that drive the
// HTML-over-the-wire pattern. Developers use these in Go to build fully reactive
// UIs without writing any JavaScript.
//
// The reactive attributes are processed by the client runtime's reactive engine,
// which attaches event listeners and performs DOM swaps when server responses
// arrive. This is conceptually similar to HTMX but tightly integrated with
// the GoScript component model and type system.
package goscript

import (
	"fmt"
	"strings"
)

// =========================================================================
// Goscript Attribute Constants
// =========================================================================

const (
	// GsTrigger is the attribute that specifies what triggers a reactive request.
	// Values: "click", "submit", "change", "load", "input", "focus", "blur",
	//         "every 2s", "intersect", "morph"
	// The trigger value can include a URL: gs-trigger="click /api/users"
	GsTrigger = "gs-trigger"

	// GsTarget is the CSS selector of the element where the response content
	// will be swapped. Can be "this" to refer to the element itself, a CSS
	// selector like "#result", or a component reference like "component:name".
	GsTarget = "gs-target"

	// GsSwap determines how the response content is inserted into the target.
	// See SwapInnerHTML, SwapOuterHTML, etc. for available strategies.
	GsSwap = "gs-swap"

	// GsIndicator is a CSS selector for an element to show during the request
	// (loading spinner, overlay, etc.). Hidden when the request completes.
	GsIndicator = "gs-indicator"

	// GsDisabled is a comma-separated list of CSS selectors for elements to
	// disable while the request is in flight. Useful for preventing double-submits.
	GsDisabled = "gs-disabled"

	// GsPushUrl specifies a URL to push into the browser's history stack after
	// a successful swap. Enables SPA-like navigation without JavaScript.
	GsPushUrl = "gs-push-url"

	// GsConfirm shows a browser confirmation dialog before the request is sent.
	// If the user cancels, the request is aborted.
	GsConfirm = "gs-confirm"

	// GsBoost enhances a normal <a> or <form> to use goscript reactivity.
	// The element behaves normally without JavaScript but upgrades to use
	// reactive swaps when the runtime is loaded (progressive enhancement).
	GsBoost = "gs-boost"
)

// =========================================================================
// Swap Strategy Constants
// =========================================================================

const (
	// SwapInnerHTML replaces the target element's innerHTML with the response.
	SwapInnerHTML = "innerHTML"

	// SwapOuterHTML replaces the entire target element (including the tag itself)
	// with the response content.
	SwapOuterHTML = "outerHTML"

	// SwapBeforeEnd appends the response content before the target's closing tag.
	SwapBeforeEnd = "beforeend"

	// SwapAfterEnd inserts the response content immediately after the target element.
	SwapAfterEnd = "afterend"

	// SwapBeforeBegin inserts the response content immediately before the target element.
	SwapBeforeBegin = "beforebegin"

	// SwapAfterBegin inserts the response content after the target's opening tag.
	SwapAfterBegin = "afterbegin"

	// SwapDelete removes the target element from the DOM after the response.
	// The response content is ignored.
	SwapDelete = "delete"

	// SwapMorph performs an intelligent diff between the current content and the
	// response, updating only changed nodes. Preserves focus and input state.
	SwapMorph = "morph"

	// SwapNone performs no DOM swap. Useful when you only want to trigger
	// events or update state.
	SwapNone = "none"
)

// =========================================================================
// Trigger Type Constants
// =========================================================================

const (
	// TriggerClick fires on a mouse click.
	TriggerClick = "click"

	// TriggerSubmit fires on form submission.
	TriggerSubmit = "submit"

	// TriggerChange fires when an input value changes and loses focus.
	TriggerChange = "change"

	// TriggerLoad fires immediately when the element is parsed by the runtime.
	TriggerLoad = "load"

	// TriggerInput fires on every keystroke or input event.
	TriggerInput = "input"

	// TriggerFocus fires when the element receives focus.
	TriggerFocus = "focus"

	// TriggerBlur fires when the element loses focus.
	TriggerBlur = "blur"

	// TriggerEvery starts a polling interval. Format: "every <duration>"
	// where duration is "2s", "500ms", "1m", etc.
	TriggerEvery = "every"

	// TriggerIntersect fires when the element enters the viewport (uses
	// IntersectionObserver). Fires at most once per element.
	TriggerIntersect = "intersect"

	// TriggerMorph performs a morph swap on load.
	TriggerMorph = "morph"
)

// =========================================================================
// Reactive Props Builders
// =========================================================================

// GoscriptProps creates a Props map from alternating key-value pairs of
// goscript reactive attributes. It is the primary way to build reactive
// attribute sets for elements.
//
// Each pair consists of a constant key (e.g. GsTrigger) and a string value.
// Odd trailing arguments are silently ignored.
//
// Usage:
//
//	GoscriptProps(GsTrigger, "click", GsTarget, "#result", GsSwap, "innerHTML")
//	// => Props{"gs-trigger": "click", "gs-target": "#result", "gs-swap": "innerHTML"}
//
//	GoscriptProps(GsTrigger, "click /api/users", GsSwap, SwapInnerHTML)
func GoscriptProps(keyValuePairs ...interface{}) Props {
	props := make(Props)
	for i := 0; i+1 < len(keyValuePairs); i += 2 {
		key := fmt.Sprintf("%v", keyValuePairs[i])
		val := fmt.Sprintf("%v", keyValuePairs[i+1])
		props[key] = val
	}
	return props
}

// =========================================================================
// Reactive Element Constructors
// =========================================================================

// ReactiveElement creates a Props map with goscript reactive attributes
// suitable for passing to CreateElement. It is a convenience wrapper
// around GoscriptProps that also accepts children.
//
// The tag parameter is the HTML tag name. Children are returned alongside
// the props as a convenience, but the caller should pass them separately
// to CreateElement.
//
// Usage:
//
//	props := ReactiveElement("button", GsTrigger, "click", GsTarget, "#result")
//	html := CreateElement("button", props, "Load Data")
func ReactiveElement(tag string, keyValuePairs ...interface{}) Props {
	return GoscriptProps(keyValuePairs...)
}

// OnClick creates Props for a click-triggered reactive element. The URL
// is included in the gs-trigger value, and the target and swap are set
// as separate attributes.
//
// Usage:
//
//	props := OnClick("/api/users", "#user-list", SwapInnerHTML)
//	html := CreateElement("button", props, "Load Users")
func OnClick(url string, target string, swap string) Props {
	return Props{
		GsTrigger: fmt.Sprintf("%s %s", TriggerClick, url),
		GsTarget:  target,
		GsSwap:    swap,
	}
}

// OnSubmit creates Props for a form with reactive submission. The form
// will intercept the submit event and send a request to the given URL,
// swapping the response into the target element.
//
// Usage:
//
//	props := OnSubmit("/api/contact", "#result", SwapInnerHTML)
//	html := CreateElement("form", props, fields...)
func OnSubmit(url string, target string, swap string) Props {
	return Props{
		GsTrigger: fmt.Sprintf("%s %s", TriggerSubmit, url),
		GsTarget:  target,
		GsSwap:    swap,
	}
}

// OnLoad creates Props for auto-loading content when the element is first
// processed by the runtime. Useful for lazy-loading sections, counters,
// or initial data fetches.
//
// Usage:
//
//	props := OnLoad("/api/stats", "#stats-panel", SwapInnerHTML)
//	html := CreateElement("div", props)
func OnLoad(url string, target string, swap string) Props {
	return Props{
		GsTrigger: fmt.Sprintf("%s %s", TriggerLoad, url),
		GsTarget:  target,
		GsSwap:    swap,
	}
}

// Poll creates Props for periodic polling. The interval string uses
// Go-style duration suffixes: "2s" (seconds), "500ms" (milliseconds),
// "1m" (minutes). The URL is placed after the interval in the gs-trigger
// attribute.
//
// Usage:
//
//	props := Poll("/api/notifications", "5s", "#notifications", SwapInnerHTML)
//	html := CreateElement("div", props)
func Poll(url string, interval string, target string, swap string) Props {
	return Props{
		GsTrigger: fmt.Sprintf("%s %s %s", TriggerEvery, interval, url),
		GsTarget:  target,
		GsSwap:    swap,
	}
}

// LazyLoad creates Props for loading content when the element enters the
// browser viewport. Uses IntersectionObserver under the hood. The URL
// is fetched at most once per element.
//
// Usage:
//
//	props := LazyLoad("/api/comments", "#comments", SwapInnerHTML)
//	html := CreateElement("div", props, "Loading comments...")
func LazyLoad(url string, target string, swap string) Props {
	return Props{
		GsTrigger: fmt.Sprintf("%s %s", TriggerIntersect, url),
		GsTarget:  target,
		GsSwap:    swap,
	}
}

// Boost returns Props that enhance a normal <a> or <form> element to use
// goscript reactivity. The element works normally without JavaScript but
// upgrades to reactive swaps when the runtime is loaded (progressive
// enhancement).
//
// Usage:
//
//	props := Boost()
//	html := CreateElement("a", props, "href", "/about", "About Us")
func Boost() Props {
	return Props{
		GsBoost: "true",
	}
}

// WithConfirm adds a confirmation dialog to existing reactive Props.
// The user must confirm before the request is sent. Returns a new Props
// map with the gs-confirm attribute added.
//
// Usage:
//
//	props := OnClick("/api/delete", "#result", SwapInnerHTML)
//	props = WithConfirm(props, "Are you sure you want to delete this?")
func WithConfirm(props Props, message string) Props {
	result := make(Props)
	for k, v := range props {
		result[k] = v
	}
	result[GsConfirm] = message
	return result
}

// WithIndicator adds a loading indicator to existing reactive Props.
// The indicator CSS selector points to an element that will be shown
// during the request and hidden when it completes.
//
// Usage:
//
//	props := OnClick("/api/data", "#result", SwapInnerHTML)
//	props = WithIndicator(props, "#loading-spinner")
func WithIndicator(props Props, selector string) Props {
	result := make(Props)
	for k, v := range props {
		result[k] = v
	}
	result[GsIndicator] = selector
	return result
}

// WithDisabled adds elements to disable during a reactive request.
// Accepts a variadic list of CSS selectors. Returns a new Props map
// with the gs-disabled attribute set to a comma-separated selector list.
//
// Usage:
//
//	props := OnSubmit("/api/form", "#result", SwapInnerHTML)
//	props = WithDisabled(props, "#submit-btn", "#cancel-btn")
func WithDisabled(props Props, selectors ...string) Props {
	result := make(Props)
	for k, v := range props {
		result[k] = v
	}
	result[GsDisabled] = strings.Join(selectors, ", ")
	return result
}

// WithPushURL adds URL history pushing to existing reactive Props.
// After a successful swap, the browser URL will be updated to the
// specified path.
//
// Usage:
//
//	props := OnClick("/users/42", "#content", SwapInnerHTML)
//	props = WithPushURL(props, "/users/42")
func WithPushURL(props Props, url string) Props {
	result := make(Props)
	for k, v := range props {
		result[k] = v
	}
	result[GsPushUrl] = url
	return result
}

// =========================================================================
// Reactive Element Helpers
// =========================================================================

// ReactiveButton creates a complete button element string with click-triggered
// reactivity. This is a convenience function that combines CreateElement with
// OnClick props.
//
// Usage:
//
//	html := ReactiveButton("Load Users", "/api/users", "#result", SwapInnerHTML)
func ReactiveButton(text string, url string, target string, swap string) string {
	props := OnClick(url, target, swap)
	return CreateElement("button", props, text)
}

// ReactiveLink creates a complete anchor element with click-triggered
// reactivity. The link uses goscript reactivity instead of a full page
// navigation.
//
// Usage:
//
//	html := ReactiveLink("View Profile", "/api/profile", "#main", SwapInnerHTML)
func ReactiveLink(text string, url string, target string, swap string) string {
	props := OnClick(url, target, swap)
	return CreateElement("a", props, text)
}

// ReactiveDiv creates a div element with reactive attributes. Commonly
// used as a swap target or a lazy-load container.
//
// Usage:
//
//	html := ReactiveDiv("data-container", OnLoad, "/api/data", "#data-container", SwapInnerHTML)
func ReactiveDiv(id string, trigger string, url string, target string, swap string) string {
	props := Props{
		"id":        id,
		GsTrigger:  fmt.Sprintf("%s %s", trigger, url),
		GsTarget:   target,
		GsSwap:     swap,
	}
	return CreateElement("div", props)
}

// ReactiveForm creates a form element with reactive submission. The form
// will intercept the submit event and send a request to the action URL.
//
// Usage:
//
//	html := ReactiveForm("/api/contact", "#result", SwapInnerHTML, formFields...)
func ReactiveForm(action string, target string, swap string, children ...interface{}) string {
	props := OnSubmit(action, target, swap)
	return CreateElement("form", props, children...)
}

// MergeReactiveProps merges multiple Props maps into a single Props map.
// Later maps override earlier maps for duplicate keys. This is useful for
// combining base reactive props with additional customization.
//
// Usage:
//
//	base := OnClick("/api/data", "#result", SwapInnerHTML)
//	extra := Props{"class": "btn-primary", "id": "load-btn"}
//	merged := MergeReactiveProps(base, extra)
func MergeReactiveProps(props ...Props) Props {
	result := make(Props)
	for _, p := range props {
		for k, v := range p {
			result[k] = v
		}
	}
	return result
}
