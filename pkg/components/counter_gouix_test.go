package components

import (
	"strings"
	"testing"

	"github.com/davidjeba/goscript/pkg/gouix"
)

func TestCounter(t *testing.T) {
	// Create a counter
	counter := NewCounter("test-counter", gouix.Props{
		"initialCount": 5,
		"title":        "Test Counter",
		"theme":        "light",
	})

	// Render the counter
	html := counter.Render()

	// Check that the counter renders correctly
	if !strings.Contains(html, "Test Counter") {
		t.Errorf("Counter should contain title 'Test Counter'")
	}

	if !strings.Contains(html, ">5<") {
		t.Errorf("Counter should display initial count of 5")
	}

	// Test increment event
	counter.HandleEvent(gouix.Event{
		Type:   "increment",
		Target: "test-counter",
	})

	// Re-render and check that the count increased
	html = counter.Render()
	if !strings.Contains(html, ">6<") {
		t.Errorf("Counter should display count of 6 after increment")
	}

	// Test decrement event
	counter.HandleEvent(gouix.Event{
		Type:   "decrement",
		Target: "test-counter",
	})

	// Re-render and check that the count decreased
	html = counter.Render()
	if !strings.Contains(html, ">5<") {
		t.Errorf("Counter should display count of 5 after decrement")
	}

	// Test reset event
	counter.HandleEvent(gouix.Event{
		Type:   "reset",
		Target: "test-counter",
	})

	// Re-render and check that the count reset
	html = counter.Render()
	if !strings.Contains(html, ">5<") {
		t.Errorf("Counter should display initial count of 5 after reset")
	}
}

func TestDraggableCounter(t *testing.T) {
	// Create a draggable counter
	counter := NewDraggableCounter("test-draggable", gouix.Props{
		"initialCount": 10,
		"title":        "Draggable Counter",
		"theme":        "dark",
	})

	// Render the counter
	html := counter.Render()

	// Check that the counter is draggable
	if !strings.Contains(html, "draggable=\"true\"") {
		t.Errorf("Draggable counter should have draggable attribute")
	}

	// Check that the counter has the correct class
	if !strings.Contains(html, "class=\"counter draggable\"") {
		t.Errorf("Draggable counter should have 'draggable' class")
	}

	// Check that the counter renders correctly
	if !strings.Contains(html, "Draggable Counter") {
		t.Errorf("Counter should contain title 'Draggable Counter'")
	}

	if !strings.Contains(html, ">10<") {
		t.Errorf("Counter should display initial count of 10")
	}
}

func TestCanvasRendering(t *testing.T) {
	// Create a canvas
	canvas := gouix.NewCanvas("test-canvas", 400, 300, nil)

	// Add a counter to the canvas
	CanvasCounter(canvas, "canvas-counter", 50, 50, gouix.Props{
		"initialCount": 20,
		"title":        "Canvas Counter",
	})

	// Render the canvas
	html := canvas.Render()

	// Check that the canvas renders correctly
	if !strings.Contains(html, "<svg id=\"test-canvas\" width=\"400\" height=\"300\"") {
		t.Errorf("Canvas should render an SVG with correct dimensions")
	}

	// Check that the counter elements are rendered
	if !strings.Contains(html, "Canvas Counter") {
		t.Errorf("Canvas should contain the counter title")
	}

	if !strings.Contains(html, "20") {
		t.Errorf("Canvas should display the initial count")
	}
}

func TestCounterWithHooks(t *testing.T) {
	// Render a counter with hooks
	html := CounterWithHooks(gouix.Props{
		"initialCount": 15,
		"title":        "Hooks Counter",
		"id":           "hooks-test",
	})

	// Check that the counter renders correctly
	if !strings.Contains(html, "Hooks Counter") {
		t.Errorf("Counter should contain title 'Hooks Counter'")
	}

	if !strings.Contains(html, ">15<") {
		t.Errorf("Counter should display initial count of 15")
	}

	// Check that the counter has the correct ID
	if !strings.Contains(html, "id=\"hooks-test\"") {
		t.Errorf("Counter should have the correct ID")
	}
}

func TestHomePage(t *testing.T) {
	// Create a home page
	home := NewHomePage("test-home", nil)

	// Render the home page
	html := home.Render()

	// Check that the home page renders correctly
	if !strings.Contains(html, "GoUIX Demo") {
		t.Errorf("Home page should contain title 'GoUIX Demo'")
	}

	// Check that the counters are rendered
	if !strings.Contains(html, "Counter 1") {
		t.Errorf("Home page should contain 'Counter 1'")
	}

	if !strings.Contains(html, "Counter 2") {
		t.Errorf("Home page should contain 'Counter 2'")
	}

	if !strings.Contains(html, "Draggable Counter") {
		t.Errorf("Home page should contain 'Draggable Counter'")
	}

	// Check that the canvas is rendered
	if !strings.Contains(html, "<svg id=\"home-canvas\"") {
		t.Errorf("Home page should contain a canvas")
	}

	// Test adding a counter
	home.HandleEvent(gouix.Event{
		Type:   "addCounter",
		Target: "test-home",
	})

	// Re-render and check that a new counter was added
	html = home.Render()
	if !strings.Contains(html, "Counter 3") {
		t.Errorf("Home page should contain 'Counter 3' after adding a counter")
	}
}