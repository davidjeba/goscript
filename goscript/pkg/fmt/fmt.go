// Package fmt provides formatted I/O functions for goscript .gs files.
// These functions are compiled to JavaScript calls to __gs.fmt by the GS compiler.
//
// This package provides a subset of Go's standard fmt package functionality,
// adapted for the browser environment. It supports Go-style format verbs (%s, %d, %v)
// and is the primary way to output text and format strings in .gs code.
//
// This package serves as API documentation and type definitions for .gs developers.
// The functions listed here map to JavaScript operations in the goscript client runtime
// (pkg/gslib/runtime.js). You do not need to import this package in Go server code —
// it exists solely for .gs files.
//
// # Format Verbs
//
// The formatting functions support these Go-style verbs:
//   - %s  — string
//   - %d  — integer/number
//   - %v  — any value (default format)
//   - %%  — literal percent sign
//
// The verb count must match the argument count.
//
// # Usage in .gs code
//
//	import "goscript/fmt"
//
//	func main() {
//	    fmt.Println("Hello, World!")
//	    fmt.Printf("Name: %s, Age: %d\n", "Alice", 30)
//
//	    msg := fmt.Sprintf("User %s has %d items", "Bob", 5)
//	    fmt.Println(msg)
//
//	    err := fmt.Errorf("failed to load: %s", "network timeout")
//	    fmt.Println(err)
//	}
package fmt

// ---------------------------------------------------------------------------
// Output Functions
// ---------------------------------------------------------------------------

// Println prints its arguments to the browser console (console.log) followed
// by a newline. Arguments are separated by spaces.
//
// In .gs code this compiles to console.log(...args).
//
// Parameters:
//   - args: one or more values to print
//
// Example (.gs):
//
//	fmt.Println("Hello")
//	fmt.Println("count =", 42)
//	fmt.Println("user:", name, "age:", age)
func Println(args ...interface{}) {}

// Printf formats its arguments according to the format string and prints the
// result to the browser console (console.log). A newline is appended automatically.
//
// Supports Go-style format verbs: %s (string), %d (number), %v (any value), %% (literal %).
//
// In .gs code this compiles to console.log(__gs.fmt.sprintf(format, ...args)).
//
// Parameters:
//   - format: a format string with %s, %d, %v placeholders
//   - args: values to interpolate into the format string
//
// Example (.gs):
//
//	fmt.Printf("Hello, %s! You are %d years old.\n", "Alice", 30)
//	fmt.Printf("x=%d y=%d\n", 10, 20)
func Printf(format string, args ...interface{}) {}

// ---------------------------------------------------------------------------
// String Formatting
// ---------------------------------------------------------------------------

// Sprintf formats its arguments according to the format string and returns
// the result as a string. It does NOT print anything.
//
// Supports Go-style format verbs: %s (string), %d (number), %v (any value), %% (literal %).
//
// In .gs code this compiles to __gs.sprintf(format, ...args).
//
// Parameters:
//   - format: a format string with %s, %d, %v placeholders
//   - args: values to interpolate into the format string
//
// Returns the formatted string.
//
// Example (.gs):
//
//	msg := fmt.Sprintf("Hello, %s!", "World")
//	dom.SetTextContent(el, msg)
//
//	url := fmt.Sprintf("/api/users/%d", userID)
//	data := api.GetJSON(url)
func Sprintf(format string, args ...interface{}) string { return "" }

// ---------------------------------------------------------------------------
// Error Formatting
// ---------------------------------------------------------------------------

// Errorf formats its arguments according to the format string and returns
// the result as an error string. This is the .gs equivalent of Go's
// fmt.Errorf.
//
// Supports Go-style format verbs: %s (string), %d (number), %v (any value), %% (literal %).
//
// In .gs code this compiles to a string created via __gs.sprintf(format, ...args).
// The returned string can be used as an error value in .gs code.
//
// Parameters:
//   - format: a format string with %s, %d, %v placeholders
//   - args: values to interpolate into the format string
//
// Returns the formatted error string.
//
// Example (.gs):
//
//	err := fmt.Errorf("user %s not found", username)
//	if err != "" {
//	    dom.SetTextContent(el, "Error: "+err)
//	}
//
//	err = fmt.Errorf("HTTP %d: %s", statusCode, statusText)
func Errorf(format string, args ...interface{}) string { return "" }
