// Package router provides client-side routing functions for goscript .gs files.
// These functions are compiled to JavaScript calls to __gs.router by the GS compiler.
//
// The router uses the browser's History API (pushState / popstate) for clean,
// hash-free URLs. Navigation does not cause full page reloads — instead, it
// dispatches events that your .gs code can listen for and respond to.
//
// This package serves as API documentation and type definitions for .gs developers.
// The functions listed here map to JavaScript operations in the goscript client runtime
// (pkg/gslib/runtime.js). You do not need to import this package in Go server code —
// it exists solely for .gs files.
//
// # Usage in .gs code
//
//	import "goscript/router"
//	import "goscript/dom"
//	import "goscript/fmt"
//
//	func main() {
//	    // Navigate programmatically
//	    router.Navigate("/about")
//
//	    // Read current URL info
//	    path := router.UsePathname()
//	    params := router.UseParams()
//	    fmt.Println("Current path:", path)
//	    fmt.Println("Query params:", params)
//
//	    // Create a client-side navigation link
//	    link := router.Link("/dashboard", map[string]interface{}{
//	        "className": "nav-link",
//	    })
//	    dom.SetTextContent(link, "Dashboard")
//	    dom.AppendChild(nav, link)
//
//	    // Listen for route changes
//	    router.OnRouteChange(func() {
//	        path := router.UsePathname()
//	        fmt.Println("Route changed to:", path)
//	    })
//	}
package router

// ---------------------------------------------------------------------------
// Navigation
// ---------------------------------------------------------------------------

// Navigate pushes a new URL onto the browser history stack and dispatches
// a popstate event so route listeners can react. The page does NOT reload.
//
// In .gs code this compiles to __gs.navigate(path).
//
// Parameters:
//   - path: the new URL path (e.g. "/about", "/users/123")
//
// Example (.gs):
//
//	router.Navigate("/dashboard")
//	router.Navigate("/users/" + userID)
func Navigate(path string) {}

// Back navigates to the previous page in the browser history.
// Equivalent to clicking the browser's back button.
//
// In .gs code this compiles to window.history.back().
//
// Example (.gs):
//
//	router.Back()
func Back() {}

// Forward navigates to the next page in the browser history.
// Equivalent to clicking the browser's forward button.
//
// In .gs code this compiles to window.history.forward().
//
// Example (.gs):
//
//	router.Forward()
func Forward() {}

// ---------------------------------------------------------------------------
// Route Information
// ---------------------------------------------------------------------------

// UsePathname returns the current URL pathname (without query string or hash).
//
// In .gs code this compiles to __gs.usePathname().
//
// Returns the pathname string (e.g. "/users/123").
//
// Example (.gs):
//
//	path := router.UsePathname()
//	if path == "/login" {
//	    showLoginForm()
//	}
func UsePathname() string { return "" }

// UseParams returns the URL search parameters as a map.
// Each key is a query parameter name, and each value is the parameter's
// string value. If a parameter appears multiple times, the last value is used.
//
// In .gs code this compiles to __gs.useParams().
//
// Returns a map[string]string of query parameters.
//
// Example (.gs):
//
//	// Given URL: /search?q=hello&page=2
//	params := router.UseParams()
//	fmt.Println("Search:", params["q"])    // "hello"
//	fmt.Println("Page:", params["page"])   // "2"
func UseParams() map[string]string { return nil }

// UseQuery is an alias for UseParams. Returns the URL search parameters
// as a map, identical to UseParams.
//
// In .gs code this compiles to __gs.useQuery().
//
// Returns a map[string]string of query parameters.
//
// Example (.gs):
//
//	query := router.UseQuery()
//	filter := query["filter"]
func UseQuery() map[string]string { return nil }

// ---------------------------------------------------------------------------
// Link Component
// ---------------------------------------------------------------------------

// Link creates an HTML anchor (<a>) element that uses the client-side router
// instead of causing a full page navigation. Clicking the link calls
// router.Navigate(href) internally, preventing the default browser behavior.
//
// In .gs code this compiles to __gs.Link(href, props).
//
// Parameters:
//   - href: the target path (e.g. "/about", "/users/123")
//   - props: optional HTML attributes for the anchor element.
//     Common keys: "className", "id", "style", "title", etc.
//     Pass nil if no additional attributes are needed.
//
// Returns an Element (the created <a> node).
//
// Example (.gs):
//
//	link := router.Link("/dashboard", map[string]interface{}{
//	    "className": "nav-link primary",
//	    "title":     "Go to Dashboard",
//	})
//	dom.SetTextContent(link, "Dashboard")
//	dom.AppendChild(nav, link)
//
//	// Simple link with no extra props
//	homeLink := router.Link("/", nil)
//	dom.SetTextContent(homeLink, "Home")
func Link(href string, props map[string]interface{}) interface{} { return nil }

// ---------------------------------------------------------------------------
// Route Change Listener
// ---------------------------------------------------------------------------

// OnRouteChange registers a handler to be called whenever the route changes.
// Route changes occur when Navigate is called, when Link elements are clicked,
// or when the browser's back/forward buttons are used.
//
// In .gs code this compiles to window.addEventListener('popstate', handler).
//
// Parameters:
//   - handler: a function called with no arguments on every route change
//
// Example (.gs):
//
//	router.OnRouteChange(func() {
//	    path := router.UsePathname()
//	    switch path {
//	    case "/":
//	        renderHome()
//	    case "/about":
//	        renderAbout()
//	    default:
//	        renderNotFound()
//	    }
//	})
func OnRouteChange(handler func()) {}
