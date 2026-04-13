// Package ui provides pre-built UI component functions for goscript .gs files.
// These functions are compiled to JavaScript calls to __gs.ui by the GS compiler.
//
// The ui package offers common UI patterns — modals, toasts, tooltips, dropdowns,
// tabs, accordions, and collapsibles — as easy-to-use functions that return DOM
// elements ready to be mounted into the page.
//
// Each component function returns an Element that can be appended to the DOM,
// and components integrate with the goscript state and event systems.
//
// This package serves as API documentation and type definitions for .gs developers.
// The functions listed here map to JavaScript operations in the goscript client runtime
// (pkg/gslib/runtime.js). You do not need to import this package in Go server code —
// it exists solely for .gs files.
//
// # Usage in .gs code
//
//	import "goscript/ui"
//	import "goscript/dom"
//
//	func main() {
//	    // Show a toast notification
//	    toast := ui.Toast("Saved successfully!", "success")
//	    dom.AppendChild(document.body, toast)
//
//	    // Show a modal dialog
//	    modal := ui.Modal(
//	        dom.CreateElement("p"), // content
//	        map[string]interface{}{
//	            "title": "Confirm Action",
//	            "onClose": func() { fmt.Println("closed") },
//	        },
//	    )
//	    dom.AppendChild(document.body, modal)
//	}
package ui

// ModalOptions configures a Modal component.
type ModalOptions struct {
	Title       string      // Modal title text (displayed in header)
	CloseOnOverlay bool      // Close when clicking the overlay (default: true)
	OnClose     func()      // Callback when the modal is closed
	ClassNames  string      // Additional CSS class names for the modal container
}

// ToastType represents the visual style of a Toast notification.
// Common values: "success", "error", "warning", "info".
type ToastType string

const (
	ToastSuccess ToastType = "success"
	ToastError   ToastType = "error"
	ToastWarning ToastType = "warning"
	ToastInfo    ToastType = "info"
)

// TooltipOptions configures a Tooltip component.
type TooltipOptions struct {
	Position string // Tooltip position: "top", "bottom", "left", "right" (default: "top")
	Delay    int    // Show delay in milliseconds (default: 0)
}

// DropdownItem represents a single item in a Dropdown menu.
type DropdownItem struct {
	Label    string      // Display text
	Value    interface{} // Associated value (passed to onClick)
	Disabled bool        // Whether the item is disabled
	OnClick  func()      // Click handler for this item
}

// TabsItem represents a single tab panel.
type TabsItem struct {
	Label    string      // Tab label text
	Content  interface{} // Tab content (Element or string)
	Disabled bool        // Whether the tab is disabled
}

// AccordionItem represents a single collapsible section in an Accordion.
type AccordionItem struct {
	Title   string      // Section title (always visible)
	Content interface{} // Section content (Element or string, collapsible)
	Open    bool        // Whether the section is initially open (default: false)
}

// ---------------------------------------------------------------------------
// Modal
// ---------------------------------------------------------------------------

// Modal creates a modal dialog component with the given content and options.
//
// The modal includes:
//   - An overlay that covers the viewport
//   - A centered dialog container
//   - An optional title bar with a close button
//   - The provided content element
//
// In .gs code this compiles to a function that builds the modal DOM structure.
//
// Parameters:
//   - content: the main content element (or Element) to display inside the modal
//   - options: configuration options (may be nil for defaults)
//
// Returns an Element representing the modal (including the overlay).
// Append this to document.body to display it.
//
// Example (.gs):
//
//	content := dom.CreateElement("div")
//	dom.SetInnerHTML(content, "<p>Are you sure?</p>")
//
//	modal := ui.Modal(content, map[string]interface{}{
//	    "title": "Confirm",
//	    "closeOnOverlay": true,
//	    "onClose": func() {
//	        fmt.Println("Modal closed")
//	    },
//	})
//	dom.AppendChild(document.body, modal)
func Modal(content interface{}, options map[string]interface{}) interface{} { return nil }

// ---------------------------------------------------------------------------
// Toast
// ---------------------------------------------------------------------------

// Toast creates a toast notification element.
//
// Toasts are lightweight, non-blocking notifications that typically appear
// at the top or bottom of the screen and auto-dismiss after a few seconds.
//
// In .gs code this compiles to a function that builds the toast DOM structure.
//
// Parameters:
//   - message: the notification text
//   - toastType: the toast type/variant. Use one of the ToastType constants:
//     "success", "error", "warning", "info". Pass "" for a default style.
//
// Returns an Element representing the toast notification.
//
// Example (.gs):
//
//	toast := ui.Toast("Changes saved!", "success")
//	dom.AppendChild(document.body, toast)
//
//	errToast := ui.Toast("Failed to save. Please try again.", "error")
//	dom.AppendChild(document.body, errToast)
func Toast(message string, toastType string) interface{} { return nil }

// ---------------------------------------------------------------------------
// Tooltip
// ---------------------------------------------------------------------------

// Tooltip attaches a tooltip to an element. The tooltip appears on hover
// (or after a configurable delay) and shows the specified text.
//
// In .gs code this compiles to a function that sets up the tooltip behavior
// on the target element.
//
// Parameters:
//   - element: the target element to attach the tooltip to
//   - text: the tooltip text to display on hover
//
// Example (.gs):
//
//	btn := dom.CreateElement("button")
//	dom.SetTextContent(btn, "Hover me")
//	ui.Tooltip(btn, "Click to submit the form")
func Tooltip(element interface{}, text string) {}

// ---------------------------------------------------------------------------
// Dropdown
// ---------------------------------------------------------------------------

// Dropdown creates a dropdown menu component.
//
// The dropdown consists of a trigger element and a list of items that appears
// when the trigger is clicked. Clicking outside the dropdown closes it.
//
// In .gs code this compiles to a function that builds the dropdown DOM structure
// with toggle behavior.
//
// Parameters:
//   - trigger: the trigger element (typically a button) that toggles the dropdown
//   - items: a slice of maps, each representing a dropdown item:
//     - "label": string (display text)
//     - "value": interface{} (associated value)
//     - "disabled": bool (optional, default false)
//     - "onClick": func() (optional click handler)
//
// Returns an Element wrapping the trigger and dropdown menu.
//
// Example (.gs):
//
//	trigger := dom.CreateElement("button")
//	dom.SetTextContent(trigger, "Options")
//
//	items := []map[string]interface{}{
//	    {"label": "Edit", "onClick": func() { fmt.Println("Edit clicked") }},
//	    {"label": "Delete", "onClick": func() { fmt.Println("Delete clicked") }},
//	    {"label": "Disabled", "disabled": true},
//	}
//
//	dropdown := ui.Dropdown(trigger, items)
//	dom.AppendChild(toolbar, dropdown)
func Dropdown(trigger interface{}, items []map[string]interface{}) interface{} { return nil }

// ---------------------------------------------------------------------------
// Tabs
// ---------------------------------------------------------------------------

// Tabs creates a tabbed interface component.
//
// The tabs component displays a row of tab labels with a content panel below
// that shows the active tab's content. Clicking a tab label switches the
// visible content.
//
// In .gs code this compiles to a function that builds the tabs DOM structure
// with click-to-switch behavior.
//
// Parameters:
//   - items: a slice of maps, each representing a tab:
//     - "label": string (tab label text)
//     - "content": interface{} (Element or string — the tab's content)
//     - "disabled": bool (optional, default false)
//   - activeIndex: the zero-based index of the initially active tab (default: 0)
//
// Returns an Element containing the tab bar and content panel.
//
// Example (.gs):
//
//	items := []map[string]interface{}{
//	    {"label": "Overview", "content": "<p>Overview content</p>"},
//	    {"label": "Details", "content": dom.CreateElement("div")},
//	    {"label": "Settings", "content": "<p>Settings panel</p>", "disabled": true},
//	}
//
//	tabs := ui.Tabs(items, 0)
//	dom.AppendChild(main, tabs)
func Tabs(items []map[string]interface{}, activeIndex int) interface{} { return nil }

// ---------------------------------------------------------------------------
// Accordion
// ---------------------------------------------------------------------------

// Accordion creates an accordion component with collapsible sections.
//
// Each section has a title (always visible) and content (collapsible).
// Clicking a section title toggles its content open/closed. Multiple
// sections can be open simultaneously.
//
// In .gs code this compiles to a function that builds the accordion DOM
// structure with toggle behavior.
//
// Parameters:
//   - items: a slice of maps, each representing an accordion section:
//     - "title": string (section title)
//     - "content": interface{} (Element or string — the section's content)
//     - "open": bool (optional, default false — whether initially open)
//
// Returns an Element containing all accordion sections.
//
// Example (.gs):
//
//	items := []map[string]interface{}{
//	    {
//	        "title":   "Getting Started",
//	        "content": "<p>Welcome to goscript! Here's how to begin...</p>",
//	        "open":    true,
//	    },
//	    {
//	        "title":   "Installation",
//	        "content": "<p>Install with: go install github.com/davidjeba/goscript@latest</p>",
//	    },
//	    {
//	        "title":   "FAQ",
//	        "content": "<p>Frequently asked questions...</p>",
//	    },
//	}
//
//	accordion := ui.Accordion(items)
//	dom.AppendChild(sidebar, accordion)
func Accordion(items []map[string]interface{}) interface{} { return nil }

// ---------------------------------------------------------------------------
// Collapse
// ---------------------------------------------------------------------------

// Collapse creates a collapsible section with a trigger element and content.
//
// Clicking the trigger toggles the visibility of the content. The content
// slides open/closed with a CSS transition (if supported).
//
// In .gs code this compiles to a function that builds the collapse DOM
// structure with toggle behavior.
//
// Parameters:
//   - content: the collapsible content (Element or string)
//   - trigger: the element that toggles the collapse when clicked
//
// Returns an Element wrapping the trigger and collapsible content.
//
// Example (.gs):
//
//	content := dom.CreateElement("div")
//	dom.SetInnerHTML(content, "<p>Detailed information that can be collapsed.</p>")
//
//	trigger := dom.CreateElement("button")
//	dom.SetTextContent(trigger, "Show Details")
//
//	collapse := ui.Collapse(content, trigger)
//	dom.AppendChild(main, collapse)
func Collapse(content interface{}, trigger interface{}) interface{} { return nil }
