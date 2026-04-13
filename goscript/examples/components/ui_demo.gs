package main

import "goscript/dom"
import "goscript/state"
import "goscript/ui"
import "goscript/fmt"

// UIDemo — demonstrates goscript's built-in UI components:
// Modal, Toast notifications, Tooltip, Dropdown, Tabs, and Accordion.
// Each component is self-contained and demonstrates a common UI pattern.
//
// Usage:
//
//	gopm compile ui_demo.gs -o ui_demo.js
//	<script src="/ui_demo.js"></script>

// --- Modal Demo ---
func ModalDemo() dom.Element {
    showModal, setShowModal := state.Use(false)

    openModal := func(e dom.Event) { setShowModal(true) }
    closeModal := func(e dom.Event) { setShowModal(false) }

    return dom.CreateElement("div", dom.Props{"class": "demo-section"},
        dom.CreateElement("h3", nil, "Modal"),
        dom.CreateElement("button", dom.Props{
            "class":   "btn",
            "onclick": openModal,
        }, "Open Modal"),
        ui.Modal(ui.ModalProps{
            Open:    showModal,
            OnClose: closeModal,
            Title:   "Confirm Action",
        },
            dom.CreateElement("p", nil, "Are you sure you want to proceed? This action cannot be undone."),
            dom.CreateElement("div", dom.Props{"class": "modal-actions"},
                dom.CreateElement("button", dom.Props{
                    "class":   "btn secondary",
                    "onclick": closeModal,
                }, "Cancel"),
                dom.CreateElement("button", dom.Props{
                    "class":   "btn primary",
                    "onclick": closeModal,
                }, "Confirm"),
            ),
        ),
    )
}

// --- Toast Demo ---
func ToastDemo() dom.Element {
    toasts, setToasts := state.Use([]map[string]string{})

    addToast := func(message string, toastType string) func(dom.Event) {
        return func(e dom.Event) {
            t := map[string]string{
                "id":      fmt.Sprintf("%d", dom.NowTimestamp()),
                "message": message,
                "type":    toastType,
            }
            setToasts(append(toasts, t))

            // Auto-dismiss after 3 seconds
            go func() {
                dom.Sleep(3000)
                remaining := []map[string]string{}
                for _, existing := range toasts {
                    if existing["id"] != t["id"] {
                        remaining = append(remaining, existing)
                    }
                }
                setToasts(remaining)
            }()
        }
    }

    var toastElements []dom.Element
    for _, t := range toasts {
        toastElements = append(toastElements,
            ui.Toast(ui.ToastProps{
                Message: t["message"],
                Type:    t["type"],
            }),
        )
    }

    return dom.CreateElement("div", dom.Props{"class": "demo-section"},
        dom.CreateElement("h3", nil, "Toast Notifications"),
        dom.CreateElement("div", dom.Props{"class": "button-group"},
            dom.CreateElement("button", dom.Props{
                "class":   "btn success",
                "onclick": addToast("Operation succeeded!", "success"),
            }, "Success Toast"),
            dom.CreateElement("button", dom.Props{
                "class":   "btn danger",
                "onclick": addToast("Something went wrong.", "error"),
            }, "Error Toast"),
            dom.CreateElement("button", dom.Props{
                "class":   "btn info",
                "onclick": addToast("Here is some info.", "info"),
            }, "Info Toast"),
        ),
        dom.CreateElement("div", dom.Props{"class": "toast-container"}, toastElements),
    )
}

// --- Tooltip Demo ---
func TooltipDemo() dom.Element {
    return dom.CreateElement("div", dom.Props{"class": "demo-section"},
        dom.CreateElement("h3", nil, "Tooltips"),
        dom.CreateElement("div", dom.Props{"class": "tooltip-row"},
            ui.Tooltip(ui.TooltipProps{
                Text: "Saves your work to the server",
                Side: "top",
            },
                dom.CreateElement("button", dom.Props{"class": "btn"}, "Save"),
            ),
            ui.Tooltip(ui.TooltipProps{
                Text: "Copies content to clipboard",
                Side: "bottom",
            },
                dom.CreateElement("button", dom.Props{"class": "btn"}, "Copy"),
            ),
            ui.Tooltip(ui.TooltipProps{
                Text: "Opens the help documentation",
                Side: "right",
            },
                dom.CreateElement("button", dom.Props{"class": "btn"}, "Help"),
            ),
        ),
    )
}

// --- Dropdown Demo ---
func DropdownDemo() dom.Element {
    selected, setSelected := state.Use("Choose an option")
    dropdownOpen, setDropdownOpen := state.Use(false)

    toggleDropdown := func(e dom.Event) {
        setDropdownOpen(!dropdownOpen)
    }

    selectOption := func(value string) func(dom.Event) {
        return func(e dom.Event) {
            setSelected(value)
            setDropdownOpen(false)
        }
    }

    options := []string{"React", "Vue", "Svelte", "Goscript", "Angular"}

    var optionElements []dom.Element
    for _, opt := range options {
        optClass := "dropdown-option"
        if selected == opt {
            optClass = optClass + " selected"
        }
        optionElements = append(optionElements,
            dom.CreateElement("div", dom.Props{
                "class":   optClass,
                "onclick": selectOption(opt),
            }, opt),
        )
    }

    dropdownClass := "dropdown"
    if dropdownOpen {
        dropdownClass = dropdownClass + " open"
    }

    return dom.CreateElement("div", dom.Props{"class": "demo-section"},
        dom.CreateElement("h3", nil, "Dropdown"),
        ui.Dropdown(ui.DropdownProps{
            Label:    "Favorite Framework",
            Selected: selected,
            Open:     dropdownOpen,
            OnToggle: toggleDropdown,
        }, optionElements),
    )
}

// --- Tabs Demo ---
func TabsDemo() dom.Element {
    activeTab, setActiveTab := state.Use("tab1")

    tabs := []map[string]string{
        {"id": "tab1", "label": "Overview"},
        {"id": "tab2", "label": "Details"},
        {"id": "tab3", "label": "Settings"},
    }

    var tabButtons []dom.Element
    for _, tab := range tabs {
        isActive := activeTab == tab["id"]
        tabButtons = append(tabButtons,
            dom.CreateElement("button", dom.Props{
                "class":   "tab-btn" + dom.If(isActive, " active", ""),
                "onclick": func(e dom.Event) { setActiveTab(tab["id"]) },
            }, tab["label"]),
        )
    }

    tabContent := ""
    switch activeTab {
    case "tab1":
        tabContent = "This is the overview tab. It provides a high-level summary of the component."
    case "tab2":
        tabContent = "Detailed information lives here. You can find specs, usage examples, and API docs."
    case "tab3":
        tabContent = "Configure the component settings. Toggle features, change themes, set defaults."
    }

    return dom.CreateElement("div", dom.Props{"class": "demo-section"},
        dom.CreateElement("h3", nil, "Tabs"),
        ui.Tabs(ui.TabsProps{Active: activeTab, OnChange: setActiveTab},
            tabButtons,
            dom.CreateElement("div", dom.Props{"class": "tab-content"}, tabContent),
        ),
    )
}

// --- Accordion Demo ---
func AccordionDemo() dom.Element {
    openItems, setOpenItems := state.Use(map[string]bool{"acc1": true})

    toggle := func(id string) func(dom.Event) {
        return func(e dom.Event) {
            current := openItems[id]
            openItems[id] = !current
            setOpenItems(openItems)
        }
    }

    sections := []map[string]string{
        {"id": "acc1", "title": "What is Goscript?", "body": "Goscript is a Go web framework that lets you write .gs files (Go syntax) that compile to JavaScript. It provides reactive state management, DOM helpers, and a component model."},
        {"id": "acc2", "title": "How does it work?", "body": "Write your client-side code in Go syntax using func, struct, and familiar patterns. The goscript compiler (gopm) transpiles .gs files to optimized JavaScript that runs in the browser."},
        {"id": "acc3", "title": "Can I use it in production?", "body": "Goscript is designed for production use. It supports SSR, hydration, WebSocket real-time features, and a full component lifecycle."},
    }

    var accordionItems []dom.Element
    for _, s := range sections {
        accordionItems = append(accordionItems,
            ui.AccordionItem(ui.AccordionItemProps{
                ID:     s["id"],
                Title:  s["title"],
                Open:   openItems[s["id"]],
                OnToggle: toggle(s["id"]),
            }, s["body"]),
        )
    }

    return dom.CreateElement("div", dom.Props{"class": "demo-section"},
        dom.CreateElement("h3", nil, "Accordion"),
        ui.Accordion(nil, accordionItems),
    )
}

// --- Main App ---
// UIDemoApp composes all UI component demos into a single page.
func UIDemoApp() dom.Element {
    return dom.CreateElement("div", dom.Props{"class": "ui-demo"},
        dom.CreateElement("h2", nil, "🎨 Goscript UI Components"),
        dom.CreateElement("p", dom.Props{"class": "subtitle"}, "Built-in component library demonstrations"),

        ModalDemo(),
        ToastDemo(),
        TooltipDemo(),
        DropdownDemo(),
        TabsDemo(),
        AccordionDemo(),
    )
}

// Mount the UI demo into #app.
func main() {
    dom.Mount("#app", UIDemoApp())
}
