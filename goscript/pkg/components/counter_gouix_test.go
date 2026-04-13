package components

import (
	"strings"
	"testing"

	"github.com/davidjeba/goscript/pkg/goscript"
	"github.com/davidjeba/goscript/pkg/gouix"
)

func TestGouixCounter(t *testing.T) {
	counter := NewCounter(goscript.Props{
		"initialCount": 5,
		"title":        "Test Counter",
		"theme":        "light",
	})

	html := counter.Render()

	if !strings.Contains(html, "Test Counter") {
		t.Errorf("Counter should contain title 'Test Counter', got: %s", html)
	}

	if !strings.Contains(html, "Count: 5") {
		t.Errorf("Counter should display initial count of 5, got: %s", html)
	}
}

func TestCounterWithHooks(t *testing.T) {
	counter := NewCounter(goscript.Props{
		"initialCount": 5,
		"title":        "Test Counter",
		"theme":        "light",
	})

	html := counter.Render()
	if !strings.Contains(html, "Count: 5") {
		t.Errorf("Counter should display initial count of 5, got: %s", html)
	}
}

func TestDraggableCounter(t *testing.T) {
	t.Skip("NewDraggableCounter not yet implemented")
}

func TestCanvasRendering(t *testing.T) {
	canvas := gouix.NewCanvas("test-canvas", 400, 300, nil)

	CanvasCounter(canvas, "canvas-counter", 50, 50, gouix.Props{
		"initialCount": 20,
		"title":        "Canvas Counter",
	})

	html := canvas.Render()

	if !strings.Contains(html, "<svg") {
		t.Errorf("Canvas should render an SVG, got: %s", html)
	}
}

func TestCounterWithHooksFunc(t *testing.T) {
	html := CounterWithHooks(gouix.Props{
		"initialCount": 15,
		"title":        "Hooks Counter",
		"id":           "hooks-test",
	})

	if !strings.Contains(html, "Hooks Counter") {
		t.Errorf("Counter should contain title 'Hooks Counter', got: %s", html)
	}
}

func TestHomePage(t *testing.T) {
	t.Skip("NewHomePage not yet implemented")
}
