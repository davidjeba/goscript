package gouix

import (
        "strings"
        "testing"
)

// TestCanvas tests the Canvas functionality
func TestCanvas(t *testing.T) {
        // Create a canvas
        canvas := NewCanvas("test-canvas", 800, 600, nil)
        
        // Test GetID
        if canvas.GetID() != "test-canvas" {
                t.Errorf("Expected ID to be 'test-canvas', got '%s'", canvas.GetID())
        }
        
        // Test GetSize
        size := canvas.GetSize()
        if size.Width != 800 {
                t.Errorf("Expected width to be 800, got '%v'", size.Width)
        }
        if size.Height != 600 {
                t.Errorf("Expected height to be 600, got '%v'", size.Height)
        }
        
        // Test Render
        html := canvas.Render()
        if !strings.Contains(html, "<svg") {
                t.Errorf("Expected HTML to contain '<svg', got '%s'", html)
        }
        if !strings.Contains(html, "width=\"800\"") {
                t.Errorf("Expected HTML to contain 'width=\"800\"', got '%s'", html)
        }
        if !strings.Contains(html, "height=\"600\"") {
                t.Errorf("Expected HTML to contain 'height=\"600\"', got '%s'", html)
        }
}

// TestRectangle tests the Rectangle functionality
func TestRectangle(t *testing.T) {
        // Create a rectangle
        rect := Rectangle("test-rect", 50, 100, 200, 150, Props{
                "fill": "red",
                "stroke": "black",
        })
        
        // Test GetID
        if rect.GetID() != "test-rect" {
                t.Errorf("Expected ID to be 'test-rect', got '%s'", rect.GetID())
        }
        
        // Test GetPosition
        pos := rect.GetPosition()
        if pos.X != 50 {
                t.Errorf("Expected X to be 50, got '%v'", pos.X)
        }
        if pos.Y != 100 {
                t.Errorf("Expected Y to be 100, got '%v'", pos.Y)
        }
        
        // Test GetSize
        size := rect.GetSize()
        if size.Width != 200 {
                t.Errorf("Expected width to be 200, got '%v'", size.Width)
        }
        if size.Height != 150 {
                t.Errorf("Expected height to be 150, got '%v'", size.Height)
        }
        
        // Test Render
        html := rect.Render()
        if !strings.Contains(html, "<rect") {
                t.Errorf("Expected HTML to contain '<rect', got '%s'", html)
        }
        if !strings.Contains(html, "x=\"50.000000\"") {
                t.Errorf("Expected HTML to contain 'x=\"50.000000\"', got '%s'", html)
        }
        if !strings.Contains(html, "y=\"100.000000\"") {
                t.Errorf("Expected HTML to contain 'y=\"100.000000\"', got '%s'", html)
        }
        if !strings.Contains(html, "width=\"200.000000\"") {
                t.Errorf("Expected HTML to contain 'width=\"200.000000\"', got '%s'", html)
        }
        if !strings.Contains(html, "height=\"150.000000\"") {
                t.Errorf("Expected HTML to contain 'height=\"150.000000\"', got '%s'", html)
        }
        if !strings.Contains(html, "fill=\"red\"") {
                t.Errorf("Expected HTML to contain 'fill=\"red\"', got '%s'", html)
        }
        if !strings.Contains(html, "stroke=\"black\"") {
                t.Errorf("Expected HTML to contain 'stroke=\"black\"', got '%s'", html)
        }
}

// TestCircle tests the Circle functionality
func TestCircle(t *testing.T) {
        // Create a circle
        circle := Circle("test-circle", 100, 200, 50, Props{
                "fill": "blue",
        })
        
        // Test GetID
        if circle.GetID() != "test-circle" {
                t.Errorf("Expected ID to be 'test-circle', got '%s'", circle.GetID())
        }
        
        // Test GetPosition
        pos := circle.GetPosition()
        if pos.X != 50 {
                t.Errorf("Expected X to be 50, got '%v'", pos.X)
        }
        if pos.Y != 150 {
                t.Errorf("Expected Y to be 150, got '%v'", pos.Y)
        }
        
        // Test Render
        html := circle.Render()
        if !strings.Contains(html, "<circle") {
                t.Errorf("Expected HTML to contain '<circle', got '%s'", html)
        }
        if !strings.Contains(html, "cx=\"100.000000\"") {
                t.Errorf("Expected HTML to contain 'cx=\"100.000000\"', got '%s'", html)
        }
        if !strings.Contains(html, "cy=\"200.000000\"") {
                t.Errorf("Expected HTML to contain 'cy=\"200.000000\"', got '%s'", html)
        }
        if !strings.Contains(html, "r=\"50.000000\"") {
                t.Errorf("Expected HTML to contain 'r=\"50.000000\"', got '%s'", html)
        }
        if !strings.Contains(html, "fill=\"blue\"") {
                t.Errorf("Expected HTML to contain 'fill=\"blue\"', got '%s'", html)
        }
}

// TestText tests the Text functionality
func TestText(t *testing.T) {
        // Create text
        text := Text("test-text", 150, 250, "Hello, World!", Props{
                "fill": "black",
                "font-size": "24px",
        })
        
        // Test GetID
        if text.GetID() != "test-text" {
                t.Errorf("Expected ID to be 'test-text', got '%s'", text.GetID())
        }
        
        // Test GetPosition
        pos := text.GetPosition()
        if pos.X != 150 {
                t.Errorf("Expected X to be 150, got '%v'", pos.X)
        }
        if pos.Y != 250 {
                t.Errorf("Expected Y to be 250, got '%v'", pos.Y)
        }
        
        // Test Render
        html := text.Render()
        if !strings.Contains(html, "<text") {
                t.Errorf("Expected HTML to contain '<text', got '%s'", html)
        }
        if !strings.Contains(html, "x=\"150.000000\"") {
                t.Errorf("Expected HTML to contain 'x=\"150.000000\"', got '%s'", html)
        }
        if !strings.Contains(html, "y=\"260.000000\"") {
                t.Errorf("Expected HTML to contain 'y=\"260.000000\"', got '%s'", html)
        }
        if !strings.Contains(html, "Hello, World!") {
                t.Errorf("Expected HTML to contain 'Hello, World!', got '%s'", html)
        }
        if !strings.Contains(html, "fill=\"black\"") {
                t.Errorf("Expected HTML to contain 'fill=\"black\"', got '%s'", html)
        }
}

// TestCanvasWithElements tests adding elements to a canvas
func TestCanvasWithElements(t *testing.T) {
        // Create a canvas
        canvas := NewCanvas("test-canvas", 800, 600, nil)
        
        // Add a rectangle
        rect := Rectangle("test-rect", 50, 100, 200, 150, Props{
                "fill": "red",
        })
        canvas.AddElement(rect)
        
        // Add a circle
        circle := Circle("test-circle", 300, 200, 50, Props{
                "fill": "blue",
        })
        canvas.AddElement(circle)
        
        // Test Render
        html := canvas.Render()
        if !strings.Contains(html, "<rect") {
                t.Errorf("Expected HTML to contain '<rect', got '%s'", html)
        }
        if !strings.Contains(html, "<circle") {
                t.Errorf("Expected HTML to contain '<circle', got '%s'", html)
        }
        
        // Test FindElementByID
        element := canvas.FindElementByID("test-rect")
        if element == nil {
                t.Errorf("Expected to find element with ID 'test-rect'")
        }
        if element.GetID() != "test-rect" {
                t.Errorf("Expected ID to be 'test-rect', got '%s'", element.GetID())
        }
        
        // Test RemoveElement
        canvas.RemoveElement("test-rect")
        
        // Test that the element was removed
        element = canvas.FindElementByID("test-rect")
        if element != nil {
                t.Errorf("Expected element with ID 'test-rect' to be removed")
        }
        
        // Test that the other element is still there
        element = canvas.FindElementByID("test-circle")
        if element == nil {
                t.Errorf("Expected to find element with ID 'test-circle'")
        }
}