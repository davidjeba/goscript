package components

import (
	"strings"
	"testing"

	"github.com/davidjeba/goscript/pkg/goscript"
)

func TestCounter(t *testing.T) {
	// Create a counter with initial count of 5
	counter := NewCounter(goscript.Props{
		"initialCount": 5,
		"title":        "Test Counter",
	})

	// Render the counter
	html := counter.Render()

	// Check that the title is correct
	if !strings.Contains(html, "Test Counter") {
		t.Errorf("Counter title not found in rendered HTML: %s", html)
	}

	// Check that the count is correct
	if !strings.Contains(html, "Count: 5") {
		t.Errorf("Initial count not found in rendered HTML: %s", html)
	}

	// Test increment
	counter.Increment()
	html = counter.Render()
	if !strings.Contains(html, "Count: 6") {
		t.Errorf("Count not incremented correctly: %s", html)
	}

	// Test decrement
	counter.Decrement()
	html = counter.Render()
	if !strings.Contains(html, "Count: 5") {
		t.Errorf("Count not decremented correctly: %s", html)
	}
}

func TestFunctionalCounter(t *testing.T) {
	// Render a functional counter with initial count of 10
	html := FunctionalCounter(goscript.Props{
		"initialCount": 10,
		"title":        "Test Functional Counter",
	})

	// Check that the title is correct
	if !strings.Contains(html, "Test Functional Counter") {
		t.Errorf("Counter title not found in rendered HTML: %s", html)
	}

	// Check that the count is correct
	if !strings.Contains(html, "Count: 10") {
		t.Errorf("Initial count not found in rendered HTML: %s", html)
	}
}

func TestCounterPropValidation(t *testing.T) {
	// Create a counter with invalid props
	counter := NewCounter(goscript.Props{
		"initialCount": "not a number", // Should be an int
	})

	// Validate props
	errors := counter.ValidateProps()

	// Check that there's an error for initialCount
	if len(errors) == 0 {
		t.Errorf("Expected prop validation error, but got none")
	}

	// Create a counter with valid props
	counter = NewCounter(goscript.Props{
		"initialCount": 5,
		"title":        "Valid Counter",
	})

	// Validate props
	errors = counter.ValidateProps()

	// Check that there are no errors
	if len(errors) > 0 {
		t.Errorf("Expected no prop validation errors, but got: %v", errors)
	}
}