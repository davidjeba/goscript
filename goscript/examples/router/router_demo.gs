package main

import "goscript/dom"
import "goscript/state"
import "goscript/router"
import "goscript/fmt"

// RouterDemo — a single-page application demonstrating client-side routing
// with goscript. Shows route matching, navigation, link highlighting,
// path parameters, and route change handling.
//
// Usage:
//
//	gopm compile router_demo.gs -o router_demo.js
//	<script src="/router_demo.js"></script>

// Page represents a route definition.
struct Page {
    Path  string
    Title string
}

// RouterApp renders the SPA with navigation and routed content.
func RouterApp() dom.Element {
    // --- Current route state ---
    pathname, setPathname := state.Use("/")
    params, setParams := state.Use(map[string]string{})
    navigationCount, setNavCount := state.Use(0)

    // --- Route definitions ---
    pages := []Page{
        {Path: "/", Title: "Home"},
        {Path: "/about", Title: "About"},
        {Path: "/users", Title: "Users"},
        {Path: "/users/:id", Title: "User Profile"},
        {Path: "/settings", Title: "Settings"},
    }

    // --- Route change handler ---
    // useEffect watches for pathname changes and updates params.
    state.UseEffect(func() func() {
        setNavCount(navigationCount + 1)

        // Extract route parameters from the current path
        newParams := map[string]string{}

        // Check for /users/:id pattern
        if dom.MatchPath(pathname, "/users/:id") {
            newParams["id"] = dom.PathParam(pathname, "id")
            setParams(newParams)
        }

        return nil // no cleanup needed
    }, []string{pathname})

    // --- Navigation helper ---
    navigate := func(path string) func(dom.Event) {
        return func(e dom.Event) {
            e.PreventDefault()
            router.Navigate(path)
            setPathname(path)
        }
    }

    // --- Link component ---
    // Creates a navigation link with active state highlighting.
    navLink := func(path string, label string) dom.Element {
        activeClass := ""
        if pathname == path {
            activeClass = " active"
        }
        // Also match parent paths for nested routes
        if path == "/users" && dom.HasPrefix(pathname, "/users/") {
            activeClass = " active"
        }
        return dom.CreateElement("a", dom.Props{
            "href":     path,
            "class":    "nav-link" + activeClass,
            "onclick":  navigate(path),
        }, label)
    }

    // --- Page content renderer ---
    // Renders the appropriate content based on the current route.
    renderContent := func() dom.Element {
        switch pathname {
        case "/":
            return dom.CreateElement("div", dom.Props{"class": "page home"},
                dom.CreateElement("h1", nil, "🏠 Home"),
                dom.CreateElement("p", nil, "Welcome to the Goscript router demo."),
                dom.CreateElement("p", nil, "Click the navigation links above to explore different pages."),
                dom.CreateElement("p", nil, fmt.Sprintf("You have navigated %d times.", navigationCount)),
            )
        case "/about":
            return dom.CreateElement("div", dom.Props{"class": "page about"},
                dom.CreateElement("h1", nil, "📖 About"),
                dom.CreateElement("p", nil, "Goscript is a Go web framework that compiles .gs files to JavaScript."),
                dom.CreateElement("p", nil, "This router supports path parameters, nested routes, and active link highlighting."),
            )
        case "/users":
            return dom.CreateElement("div", dom.Props{"class": "page users"},
                dom.CreateElement("h1", nil, "👥 Users"),
                dom.CreateElement("ul", dom.Props{"class": "user-list"},
                    dom.CreateElement("li", nil,
                        dom.CreateElement("a", dom.Props{
                            "href":    "/users/1",
                            "onclick": navigate("/users/1"),
                        }, "User 1 — Alice"),
                    ),
                    dom.CreateElement("li", nil,
                        dom.CreateElement("a", dom.Props{
                            "href":    "/users/2",
                            "onclick": navigate("/users/2"),
                        }, "User 2 — Bob"),
                    ),
                    dom.CreateElement("li", nil,
                        dom.CreateElement("a", dom.Props{
                            "href":    "/users/3",
                            "onclick": navigate("/users/3"),
                        }, "User 3 — Carol"),
                    ),
                ),
            )
        case "/settings":
            return dom.CreateElement("div", dom.Props{"class": "page settings"},
                dom.CreateElement("h1", nil, "⚙️ Settings"),
                dom.CreateElement("p", nil, "Application settings would go here."),
            )
        default:
            // Handle /users/:id
            if dom.HasPrefix(pathname, "/users/") {
                userID := params["id"]
                if userID == "" {
                    userID = "unknown"
                }
                return dom.CreateElement("div", dom.Props{"class": "page profile"},
                    dom.CreateElement("h1", nil, fmt.Sprintf("👤 User %s", userID)),
                    dom.CreateElement("p", nil, fmt.Sprintf("Viewing profile for user ID: %s", userID)),
                    dom.CreateElement("a", dom.Props{
                        "href":    "/users",
                        "class":   "back-link",
                        "onclick": navigate("/users"),
                    }, "← Back to users"),
                )
            }
            // 404 fallback
            return dom.CreateElement("div", dom.Props{"class": "page not-found"},
                dom.CreateElement("h1", nil, "404 — Not Found"),
                dom.CreateElement("p", nil, fmt.Sprintf("The page %s does not exist.", pathname)),
                dom.CreateElement("a", dom.Props{
                    "href":    "/",
                    "onclick": navigate("/"),
                }, "← Go home"),
            )
        }
    }

    // --- Layout ---
    return dom.CreateElement("div", dom.Props{"class": "router-app"},
        // Navigation bar
        dom.CreateElement("nav", dom.Props{"class": "navbar"},
            dom.CreateElement("div", dom.Props{"class": "nav-brand"}, "Goscript Router"),
            dom.CreateElement("div", dom.Props{"class": "nav-links"},
                navLink("/", "Home"),
                navLink("/about", "About"),
                navLink("/users", "Users"),
                navLink("/settings", "Settings"),
            ),
        ),

        // Breadcrumb showing current path
        dom.CreateElement("div", dom.Props{"class": "breadcrumb"},
            dom.CreateElement("span", nil, fmt.Sprintf("Current route: %s", pathname)),
        ),

        // Page content area
        dom.CreateElement("main", dom.Props{"class": "content"}, renderContent()),
    )
}

// Mount the router app into #app.
func main() {
    // Initialize the router with the current browser URL
    router.Init()
    dom.Mount("#app", RouterApp())
}
