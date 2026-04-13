// Package api provides HTTP API helper functions for goscript .gs files.
// These functions are compiled to JavaScript calls to __gs.api by the GS compiler.
//
// All API helpers are thin wrappers around the browser's fetch() API with the
// following enhancements:
//   - Automatic GS-Request header so the Go server can identify goscript requests
//   - Descriptive error messages on non-OK HTTP responses
//   - Automatic interception of GS-Trigger and GS-State response headers
//
// This package serves as API documentation and type definitions for .gs developers.
// The functions listed here map to JavaScript operations in the goscript client runtime
// (pkg/gslib/runtime.js). You do not need to import this package in Go server code —
// it exists solely for .gs files.
//
// # Usage in .gs code
//
//	import "goscript/api"
//
//	func LoadUser(id string) {
//	    data := api.GetJSON("/api/users/" + id)
//	    // data is a Promise<Object>; use .then() in .gs
//	    name := data["name"]
//	    fmt.Println("User:", name)
//	}
//
//	func SaveUser(user map[string]interface{}) {
//	    result := api.PostJSON("/api/users", user)
//	    // result is a Promise<Object> with the server response
//	}
package api

// Response represents a fetch Response object returned by the Fetch function.
// In .gs code, this maps to the standard JavaScript Response object with
// methods like .json(), .text(), .status, .ok, etc.
type Response struct {
	Status     int
	StatusText string
	Ok         bool
	Headers    map[string]string
	Body       string
}

// ---------------------------------------------------------------------------
// JSON API Helpers
// ---------------------------------------------------------------------------

// GetJSON performs an HTTP GET request and parses the JSON response.
// Sets the GS-Request header automatically for server-side identification.
//
// In .gs code this compiles to __gs.getJSON(url).
//
// Parameters:
//   - url: the request URL (relative or absolute)
//
// Returns a Promise that resolves to the parsed JSON as a map[string]interface{}.
// Rejects with an error on non-OK HTTP responses.
//
// Example (.gs):
//
//	data := api.GetJSON("/api/items")
//	// In .gs: data.then(func(result) { ... })
//	fmt.Println(data["items"])
func GetJSON(url string) interface{} { return nil }

// PostJSON performs an HTTP POST request with a JSON-encoded body
// and parses the JSON response.
//
// Sets Content-Type: application/json and GS-Request headers automatically.
//
// In .gs code this compiles to __gs.postJSON(url, data).
//
// Parameters:
//   - url: the request URL
//   - data: the request body (any JSON-serializable value; typically a map)
//
// Returns a Promise that resolves to the parsed JSON response.
// Rejects with an error on non-OK HTTP responses.
//
// Example (.gs):
//
//	result := api.PostJSON("/api/users", map[string]interface{}{
//	    "name":  "Alice",
//	    "email": "alice@example.com",
//	})
//	fmt.Println("Created user:", result["id"])
func PostJSON(url string, data interface{}) interface{} { return nil }

// PutJSON performs an HTTP PUT request with a JSON-encoded body
// and parses the JSON response.
//
// Sets Content-Type: application/json and GS-Request headers automatically.
//
// In .gs code this compiles to __gs.putJSON(url, data).
//
// Parameters:
//   - url: the request URL
//   - data: the request body (any JSON-serializable value)
//
// Returns a Promise that resolves to the parsed JSON response.
// Rejects with an error on non-OK HTTP responses.
//
// Example (.gs):
//
//	result := api.PutJSON("/api/users/1", map[string]interface{}{
//	    "name": "Alice Updated",
//	})
func PutJSON(url string, data interface{}) interface{} { return nil }

// DeleteJSON performs an HTTP DELETE request and parses the JSON response.
// Sets the GS-Request header automatically.
//
// In .gs code this compiles to __gs.deleteJSON(url).
//
// Parameters:
//   - url: the request URL
//
// Returns a Promise that resolves to the parsed JSON response.
// Rejects with an error on non-OK HTTP responses.
//
// Example (.gs):
//
//	result := api.DeleteJSON("/api/users/1")
//	fmt.Println("Deleted:", result["success"])
func DeleteJSON(url string) interface{} { return nil }

// ---------------------------------------------------------------------------
// HTML API Helper
// ---------------------------------------------------------------------------

// PostHTML performs an HTTP POST request with a JSON body and returns
// the raw HTML/text response string. This is used when the server returns
// an HTML fragment for DOM swapping via the HTML-over-the-wire pattern.
//
// Sets Content-Type: application/json and GS-Request headers automatically.
//
// In .gs code this compiles to __gs.postHTML(url, data).
//
// Parameters:
//   - url: the request URL
//   - data: the request body (any JSON-serializable value)
//
// Returns a Promise that resolves to the HTML response string.
// Rejects with an error on non-OK HTTP responses.
//
// Example (.gs):
//
//	html := api.PostHTML("/api/fragments/card", map[string]interface{}{
//	    "title": "New Card",
//	})
//	// Insert the HTML fragment into the DOM
//	dom.SetInnerHTML(target, html)
func PostHTML(url string, data interface{}) interface{} { return nil }

// ---------------------------------------------------------------------------
// Generic Fetch
// ---------------------------------------------------------------------------

// Fetch performs a generic HTTP request using the browser's fetch API.
// Unlike the convenience helpers above, Fetch gives you full control over
// the request method, headers, and body.
//
// Note: Fetch does NOT automatically set the GS-Request header or intercept
// GS-Trigger/GS-State response headers. Use the convenience helpers (GetJSON,
// PostJSON, etc.) when possible for automatic header handling.
//
// In .gs code this compiles to fetch(url, options).
//
// Parameters:
//   - url: the request URL
//   - options: a map of fetch options:
//     - "method": "GET", "POST", "PUT", "DELETE", etc.
//     - "headers": map of header name to value
//     - "body": request body (string or object for JSON)
//
// Returns a Promise that resolves to a Response object.
//
// Example (.gs):
//
//	resp := api.Fetch("/api/data", map[string]interface{}{
//	    "method": "POST",
//	    "headers": map[string]interface{}{
//	        "Content-Type": "application/json",
//	        "X-Custom":     "value",
//	    },
//	    "body": "{\"key\":\"value\"}",
//	})
func Fetch(url string, options map[string]interface{}) interface{} { return nil }
