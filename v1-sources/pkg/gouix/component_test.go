package gouix

import (
        "strings"
        "testing"
)

// TestBaseComponent tests the base component functionality
func TestBaseComponent(t *testing.T) {
        // Create a base component
        component := NewBaseComponent("test", Props{
                "name": "Test Component",
                "value": 42,
        })
        
        // Test GetID
        if component.GetID() != "test" {
                t.Errorf("Expected ID to be 'test', got '%s'", component.GetID())
        }
        
        // Test GetProps
        props := component.GetProps()
        if props["name"] != "Test Component" {
                t.Errorf("Expected name to be 'Test Component', got '%v'", props["name"])
        }
        if props["value"] != 42 {
                t.Errorf("Expected value to be 42, got '%v'", props["value"])
        }
        
        // Test event handling
        called := false
        component.On("click", func(event Event) interface{} {
                called = true
                return nil
        })
        
        component.HandleEvent(Event{
                Type: "click",
                Target: "test",
        })
        
        if !called {
                t.Errorf("Expected event handler to be called")
        }
}

// TestHyperComponent tests the hyper(reactive) component functionality
func TestHyperComponent(t *testing.T) {
        // Create a hyper(reactive) component
        component := NewHyperComponent("test", Props{
                "name": "Test Component",
        }, map[string]interface{}{
                "count": 0,
        })
        
        // Test GetState
        if component.GetState("count") != 0 {
                t.Errorf("Expected count to be 0, got '%v'", component.GetState("count"))
        }
        
        // Test SetState
        component.SetState("count", 42)
        if component.GetState("count") != 42 {
                t.Errorf("Expected count to be 42, got '%v'", component.GetState("count"))
        }
}

// TestCreateElement tests the CreateElement function
func TestCreateElement(t *testing.T) {
        // Create a simple element
        html := CreateElement("div", Props{
                "class": "test",
                "id": "test-div",
        }, "Hello, World!")
        
        // Check the HTML
        if !strings.Contains(html, "<div") {
                t.Errorf("Expected HTML to contain '<div', got '%s'", html)
        }
        if !strings.Contains(html, "class=\"test\"") {
                t.Errorf("Expected HTML to contain 'class=\"test\"', got '%s'", html)
        }
        if !strings.Contains(html, "id=\"test-div\"") {
                t.Errorf("Expected HTML to contain 'id=\"test-div\"', got '%s'", html)
        }
        if !strings.Contains(html, "Hello, World!") {
                t.Errorf("Expected HTML to contain 'Hello, World!', got '%s'", html)
        }
        if !strings.Contains(html, "</div>") {
                t.Errorf("Expected HTML to contain '</div>', got '%s'", html)
        }
        
        // Create an element with style
        html = CreateElement("div", Props{
                "style": map[string]interface{}{
                        "color": "red",
                        "font-size": "16px",
                },
        }, "Styled Text")
        
        // Check the HTML
        if !strings.Contains(html, "style=\"") {
                t.Errorf("Expected HTML to contain 'style=\"', got '%s'", html)
        }
        if !strings.Contains(html, "color:red") {
                t.Errorf("Expected HTML to contain 'color:red', got '%s'", html)
        }
        if !strings.Contains(html, "font-size:16px") {
                t.Errorf("Expected HTML to contain 'font-size:16px', got '%s'", html)
        }
        
        // Create an element with children
        html = CreateElement("div", nil,
                CreateElement("h1", nil, "Title"),
                CreateElement("p", nil, "Paragraph"),
        )
        
        // Check the HTML
        if !strings.Contains(html, "<h1>Title</h1>") {
                t.Errorf("Expected HTML to contain '<h1>Title</h1>', got '%s'", html)
        }
        if !strings.Contains(html, "<p>Paragraph</p>") {
                t.Errorf("Expected HTML to contain '<p>Paragraph</p>', got '%s'", html)
        }
}

// TestFragment tests the Fragment function
func TestFragment(t *testing.T) {
        // Create a fragment
        html := Fragment(nil,
                CreateElement("h1", nil, "Title"),
                CreateElement("p", nil, "Paragraph"),
        )
        
        // Check the HTML
        if !strings.Contains(html, "<h1>Title</h1>") {
                t.Errorf("Expected HTML to contain '<h1>Title</h1>', got '%s'", html)
        }
        if !strings.Contains(html, "<p>Paragraph</p>") {
                t.Errorf("Expected HTML to contain '<p>Paragraph</p>', got '%s'", html)
        }
        
        // Make sure there's no wrapper element
        if strings.Contains(html, "<div") {
                t.Errorf("Expected HTML not to contain '<div', got '%s'", html)
        }
}