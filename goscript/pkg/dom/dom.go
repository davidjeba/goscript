// Package dom provides DOM manipulation functions for goscript .gs files.
// These functions are compiled to JavaScript calls to __gs.dom by the GS compiler.
//
// This package serves as API documentation and type definitions for .gs developers.
// The functions listed here map directly to JavaScript DOM operations at compile time.
// You do not need to import this package in Go server code — it exists solely for
// .gs files where the GS compiler translates these calls into browser DOM APIs.
//
// # Element Types
//
// The DOM types below (Element, HTMLElement, TextNode) represent browser DOM nodes.
// In .gs code, these are opaque values returned by DOM functions and passed to other
// DOM functions. The GS compiler handles the type translation automatically.
//
// # Usage in .gs code
//
//	import "goscript/dom"
//
//	func main() {
//	    el := dom.GetElementByID("app")
//	    dom.SetInnerHTML(el, "<h1>Hello</h1>")
//
//	    div := dom.CreateElement("div")
//	    dom.AddClass(div, "container")
//	    dom.AppendChild(el, div)
//
//	    btn := dom.CreateElement("button")
//	    dom.SetTextContent(btn, "Click me")
//	    dom.AddEventListener(btn, "click", func() {
//	        fmt.Println("clicked!")
//	    })
//	    dom.AppendChild(div, btn)
//
//	    dom.Show(el)
//	}
package dom

// Element represents a generic DOM element node.
// In compiled .gs code, this maps to a JavaScript HTMLElement or Element object.
type Element interface{}

// HTMLElement represents a DOM element that supports HTML-specific properties
// such as innerHTML, className, and style. In compiled .gs code, this maps
// to a JavaScript HTMLElement object.
type HTMLElement interface{}

// TextNode represents a DOM text node. Created by CreateTextNode or
// implicitly when using SetTextContent. Maps to JavaScript Text node.
type TextNode interface{}

// Func is a generic function signature used for DOM event handlers
// and callback functions in .gs code.
type Func func(...interface{}) interface{}

// ---------------------------------------------------------------------------
// Element Selection
// ---------------------------------------------------------------------------

// GetElementByID returns the first element with the specified ID.
// In .gs code this compiles to document.getElementById(id).
//
// Parameters:
//   - id: the element ID string (without the # prefix)
//
// Returns the matching HTMLElement, or nil if not found.
//
// Example (.gs):
//
//	el := dom.GetElementByID("app")
func GetElementByID(id string) HTMLElement { return nil }

// QuerySelector returns the first element matching the CSS selector.
// In .gs code this compiles to document.querySelector(selector).
//
// Parameters:
//   - selector: a valid CSS selector string (e.g. ".class", "#id", "div > p")
//
// Returns the first matching Element, or nil if no match.
//
// Example (.gs):
//
//	el := dom.QuerySelector(".card")
func QuerySelector(selector string) Element { return nil }

// QuerySelectorAll returns all elements matching the CSS selector.
// In .gs code this compiles to document.querySelectorAll(selector).
//
// Parameters:
//   - selector: a valid CSS selector string
//
// Returns a slice of all matching Element values.
//
// Example (.gs):
//
//	items := dom.QuerySelectorAll("li.active")
//	for i := 0; i < len(items); i++ {
//	    dom.AddClass(items[i], "highlighted")
//	}
func QuerySelectorAll(selector string) []Element { return nil }

// ---------------------------------------------------------------------------
// Element Creation
// ---------------------------------------------------------------------------

// CreateElement creates a new HTML element with the specified tag name.
// In .gs code this compiles to document.createElement(tag).
//
// Parameters:
//   - tag: the HTML tag name (e.g. "div", "span", "button")
//
// Returns the newly created Element.
//
// Example (.gs):
//
//	div := dom.CreateElement("div")
//	dom.SetAttribute(div, "class", "container")
func CreateElement(tag string) Element { return nil }

// CreateTextNode creates a new text node with the given content.
// In .gs code this compiles to document.createTextNode(text).
//
// Parameters:
//   - text: the text content
//
// Returns the new TextNode.
//
// Example (.gs):
//
//	txt := dom.CreateTextNode("Hello, World!")
//	dom.AppendChild(el, txt)
func CreateTextNode(text string) TextNode { return nil }

// ---------------------------------------------------------------------------
// HTML Content
// ---------------------------------------------------------------------------

// SetInnerHTML sets the inner HTML of an element.
// In .gs code this compiles to el.innerHTML = html.
//
// Parameters:
//   - el: the target element
//   - html: the HTML string to set as inner content
//
// Warning: be careful with user-provided HTML to avoid XSS vulnerabilities.
//
// Example (.gs):
//
//	dom.SetInnerHTML(el, "<h1>Welcome</h1><p>Hello!</p>")
func SetInnerHTML(el Element, html string) {}

// GetInnerHTML returns the inner HTML of an element as a string.
// In .gs code this compiles to el.innerHTML.
//
// Parameters:
//   - el: the target element
//
// Returns the inner HTML string.
//
// Example (.gs):
//
//	html := dom.GetInnerHTML(el)
//	fmt.Println(html)
func GetInnerHTML(el Element) string { return "" }

// SetTextContent sets the text content of an element, replacing all children.
// In .gs code this compiles to el.textContent = text.
//
// Parameters:
//   - el: the target element
//   - text: the plain text content
//
// Example (.gs):
//
//	dom.SetTextContent(el, "Hello, World!")
func SetTextContent(el Element, text string) {}

// GetTextContent returns the text content of an element and its descendants.
// In .gs code this compiles to el.textContent.
//
// Parameters:
//   - el: the target element
//
// Returns the concatenated text content of all descendant text nodes.
//
// Example (.gs):
//
//	text := dom.GetTextContent(el)
func GetTextContent(el Element) string { return "" }

// ---------------------------------------------------------------------------
// Attributes
// ---------------------------------------------------------------------------

// SetAttribute sets an attribute on an element.
// In .gs code this compiles to el.setAttribute(name, value).
//
// Parameters:
//   - el: the target element
//   - name: the attribute name
//   - value: the attribute value
//
// Example (.gs):
//
//	dom.SetAttribute(el, "data-id", "123")
//	dom.SetAttribute(el, "disabled", "")
func SetAttribute(el Element, name string, value string) {}

// GetAttribute returns the value of an element's attribute.
// In .gs code this compiles to el.getAttribute(name).
//
// Parameters:
//   - el: the target element
//   - name: the attribute name
//
// Returns the attribute value, or an empty string if the attribute does not exist.
//
// Example (.gs):
//
//	href := dom.GetAttribute(link, "href")
func GetAttribute(el Element, name string) string { return "" }

// ---------------------------------------------------------------------------
// CSS Classes
// ---------------------------------------------------------------------------

// AddClass adds one or more CSS classes to an element.
// In .gs code this compiles to el.classList.add(className).
//
// Parameters:
//   - el: the target element
//   - className: one or more space-separated CSS class names
//
// Example (.gs):
//
//	dom.AddClass(el, "active visible")
func AddClass(el Element, className string) {}

// RemoveClass removes one or more CSS classes from an element.
// In .gs code this compiles to el.classList.remove(className).
//
// Parameters:
//   - el: the target element
//   - className: one or more space-separated CSS class names
//
// Example (.gs):
//
//	dom.RemoveClass(el, "hidden")
func RemoveClass(el Element, className string) {}

// ToggleClass toggles a CSS class on an element.
// In .gs code this compiles to el.classList.toggle(className).
//
// Parameters:
//   - el: the target element
//   - className: the CSS class name to toggle
//
// Returns true if the class is now present, false if removed.
//
// Example (.gs):
//
//	isActive := dom.ToggleClass(el, "active")
func ToggleClass(el Element, className string) bool { return false }

// HasClass checks whether an element has a specific CSS class.
// In .gs code this compiles to el.classList.contains(className).
//
// Parameters:
//   - el: the target element
//   - className: the CSS class name to check
//
// Returns true if the element has the class.
//
// Example (.gs):
//
//	if dom.HasClass(el, "open") {
//	    dom.SetInnerHTML(el, "Expanded!")
//	}
func HasClass(el Element, className string) bool { return false }

// ---------------------------------------------------------------------------
// Inline Styles
// ---------------------------------------------------------------------------

// SetStyle sets an inline CSS style property on an element.
// In .gs code this compiles to el.style[property] = value.
//
// Parameters:
//   - el: the target element
//   - property: the CSS property name in camelCase (e.g. "backgroundColor", "fontSize")
//   - value: the CSS value (e.g. "red", "16px", "1rem solid black")
//
// Example (.gs):
//
//	dom.SetStyle(el, "backgroundColor", "#f0f0f0")
//	dom.SetStyle(el, "display", "none")
func SetStyle(el Element, property string, value string) {}

// GetStyle returns the computed value of a CSS property for an element.
// In .gs code this compiles to getComputedStyle(el)[property].
//
// Parameters:
//   - el: the target element
//   - property: the CSS property name in camelCase
//
// Returns the computed style value as a string.
//
// Example (.gs):
//
//	color := dom.GetStyle(el, "color")
func GetStyle(el Element, property string) string { return "" }

// ---------------------------------------------------------------------------
// DOM Tree Manipulation
// ---------------------------------------------------------------------------

// AppendChild appends a child node to a parent element.
// In .gs code this compiles to parent.appendChild(child).
//
// Parameters:
//   - parent: the parent element
//   - child: the child node (Element or TextNode) to append
//
// Example (.gs):
//
//	div := dom.CreateElement("div")
//	dom.AppendChild(container, div)
func AppendChild(parent Element, child interface{}) {}

// RemoveChild removes a child node from a parent element.
// In .gs code this compiles to parent.removeChild(child).
//
// Parameters:
//   - parent: the parent element
//   - child: the child node to remove
//
// Returns the removed child node.
//
// Example (.gs):
//
//	dom.RemoveChild(parent, oldDiv)
func RemoveChild(parent Element, child interface{}) {}

// InsertBefore inserts a new node before a reference node in the parent.
// In .gs code this compiles to parent.insertBefore(newNode, referenceNode).
//
// Parameters:
//   - parent: the parent element
//   - newNode: the node to insert
//   - referenceNode: the existing node before which newNode is inserted
//
// Example (.gs):
//
//	dom.InsertBefore(parent, newEl, existingEl)
func InsertBefore(parent Element, newNode interface{}, referenceNode interface{}) {}

// ReplaceChild replaces an existing child node with a new node.
// In .gs code this compiles to parent.replaceChild(newChild, oldChild).
//
// Parameters:
//   - parent: the parent element
//   - newChild: the replacement node
//   - oldChild: the node being replaced
//
// Example (.gs):
//
//	dom.ReplaceChild(parent, newDiv, oldDiv)
func ReplaceChild(parent Element, newChild interface{}, oldChild interface{}) {}

// CloneNode creates a deep or shallow copy of an element.
// In .gs code this compiles to el.cloneNode(deep).
//
// Parameters:
//   - el: the element to clone
//   - deep: if true, the element and all its descendants are cloned;
//     if false, only the element itself is cloned
//
// Returns the cloned Element.
//
// Example (.gs):
//
//	copy := dom.CloneNode(el, true)
//	dom.AppendChild(container, copy)
func CloneNode(el Element, deep bool) Element { return nil }

// ---------------------------------------------------------------------------
// DOM Traversal
// ---------------------------------------------------------------------------

// GetParent returns the parent element of the given element.
// In .gs code this compiles to el.parentNode.
//
// Parameters:
//   - el: the child element
//
// Returns the parent Element, or nil if the element has no parent.
//
// Example (.gs):
//
//	parent := dom.GetParent(el)
func GetParent(el Element) Element { return nil }

// GetChildren returns all child elements of the given element.
// In .gs code this compiles to el.children.
//
// Parameters:
//   - el: the parent element
//
// Returns a slice of child Element values.
//
// Example (.gs):
//
//	children := dom.GetParent(el)
//	for i := 0; i < len(children); i++ {
//	    fmt.Println(children[i])
//	}
func GetChildren(el Element) []Element { return nil }

// GetNextSibling returns the next sibling element.
// In .gs code this compiles to el.nextElementSibling.
//
// Parameters:
//   - el: the reference element
//
// Returns the next sibling Element, or nil if there is none.
//
// Example (.gs):
//
//	next := dom.GetNextSibling(el)
func GetNextSibling(el Element) Element { return nil }

// GetPreviousSibling returns the previous sibling element.
// In .gs code this compiles to el.previousElementSibling.
//
// Parameters:
//   - el: the reference element
//
// Returns the previous sibling Element, or nil if there is none.
//
// Example (.gs):
//
//	prev := dom.GetPreviousSibling(el)
func GetPreviousSibling(el Element) Element { return nil }

// ---------------------------------------------------------------------------
// Events
// ---------------------------------------------------------------------------

// AddEventListener registers an event handler on an element.
// In .gs code this compiles to el.addEventListener(event, handler).
//
// Parameters:
//   - el: the target element
//   - event: the event name (e.g. "click", "submit", "input", "keydown")
//   - handler: a function to call when the event fires
//
// Example (.gs):
//
//	dom.AddEventListener(btn, "click", func() {
//	    dom.SetTextContent(btn, "Clicked!")
//	})
func AddEventListener(el Element, event string, handler Func) {}

// RemoveEventListener removes a previously registered event handler.
// In .gs code this compiles to el.removeEventListener(event, handler).
//
// Parameters:
//   - el: the target element
//   - event: the event name
//   - handler: the same function reference that was passed to AddEventListener
//
// Example (.gs):
//
//	dom.RemoveEventListener(btn, "click", handleClick)
func RemoveEventListener(el Element, event string, handler Func) {}

// ---------------------------------------------------------------------------
// Visibility & Focus
// ---------------------------------------------------------------------------

// Show makes an element visible by setting display to its default value.
// In .gs code this compiles to el.style.display = ''.
//
// Parameters:
//   - el: the element to show
//
// Example (.gs):
//
//	dom.Show(modal)
func Show(el Element) {}

// Hide hides an element by setting display to "none".
// In .gs code this compiles to el.style.display = 'none'.
//
// Parameters:
//   - el: the element to hide
//
// Example (.gs):
//
//	dom.Hide(modal)
func Hide(el Element) {}

// Toggle toggles the visibility of an element.
// In .gs code this compiles to toggling el.style.display between '' and 'none'.
//
// Parameters:
//   - el: the element to toggle
//
// Example (.gs):
//
//	dom.Toggle(dropdown)
func Toggle(el Element) {}

// Focus sets focus on an element (typically an input or textarea).
// In .gs code this compiles to el.focus().
//
// Parameters:
//   - el: the element to focus
//
// Example (.gs):
//
//	dom.Focus(inputEl)
func Focus(el Element) {}

// Blur removes focus from an element.
// In .gs code this compiles to el.blur().
//
// Parameters:
//   - el: the element to blur
//
// Example (.gs):
//
//	dom.Blur(inputEl)
func Blur(el Element) {}

// ---------------------------------------------------------------------------
// Scrolling & Form Values
// ---------------------------------------------------------------------------

// ScrollTo scrolls an element (or the viewport if el is nil) to the given position.
// In .gs code this compiles to el.scrollTo(x, y) or window.scrollTo(x, y).
//
// Parameters:
//   - el: the element to scroll (nil for the viewport)
//   - x: the horizontal scroll position in pixels
//   - y: the vertical scroll position in pixels
//
// Example (.gs):
//
//	dom.ScrollTo(nil, 0, 500) // scroll to top of page
func ScrollTo(el Element, x int, y int) {}

// GetValue returns the current value of a form input element.
// In .gs code this compiles to el.value.
//
// Parameters:
//   - el: a form input element (input, select, textarea)
//
// Returns the current value string.
//
// Example (.gs):
//
//	name := dom.GetValue(nameInput)
//	fmt.Println("Name:", name)
func GetValue(el Element) string { return "" }

// SetValue sets the value of a form input element.
// In .gs code this compiles to el.value = value.
//
// Parameters:
//   - el: a form input element
//   - value: the value to set
//
// Example (.gs):
//
//	dom.SetValue(input, "Hello, World!")
func SetValue(el Element, value string) {}
