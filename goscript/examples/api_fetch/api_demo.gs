package main

import "goscript/dom"
import "goscript/state"
import "goscript/api"
import "goscript/fmt"

// ApiDemo — demonstrates fetching data from an API with loading states,
// error handling, POST requests, and JSON parsing in goscript.
//
// Usage:
//
//	gopm compile api_demo.gs -o api_demo.js
//	<script src="/api_demo.js"></script>

// Post represents a blog post from the API.
struct Post {
    ID     int
    UserID int
    Title  string
    Body   string
}

// TodoItem represents a todo item from the API.
struct TodoItem {
    ID        int
    Title     string
    Completed bool
}

// ApiApp renders the API interaction demo.
func ApiApp() dom.Element {
    // --- State for GET request ---
    posts, setPosts := state.Use([]Post{})
    postsLoading, setPostsLoading := state.Use(false)
    postsError, setPostsError := state.Use("")

    // --- State for POST request ---
    postTitle, setPostTitle := state.Use("")
    postBody, setPostBody := state.Use("")
    postResult, setPostResult := state.Use("")
    postSending, setPostSending := state.Use(false)

    // --- State for filtered todos ---
    todos, setTodos := state.Use([]TodoItem{})
    todosLoading, setTodosLoading := state.Use(false)

    // --- Fetch posts (GET) ---
    fetchPosts := func(e dom.Event) {
        setPostsLoading(true)
        setPostsError("")

        data, err := api.Get("/api/posts", api.Options{
            Headers: map[string]string{"Accept": "application/json"},
        })

        if err != nil {
            setPostsError(fmt.Sprintf("Failed to load posts: %s", err))
            setPostsLoading(false)
            return
        }

        parsed, parseErr := api.ParseJSON(data)
        if parseErr != nil {
            setPostsError(fmt.Sprintf("Failed to parse JSON: %s", parseErr))
            setPostsLoading(false)
            return
        }

        // Convert parsed data to Post slice
        var result []Post
        for _, item := range parsed {
            p := Post{
                ID:    item["id"].(int),
                Title: item["title"].(string),
                Body:  item["body"].(string),
            }
            result = append(result, p)
        }
        setPosts(result)
        setPostsLoading(false)
    }

    // --- Create a post (POST) ---
    createPost := func(e dom.Event) {
        e.PreventDefault()
        if postTitle == "" {
            return
        }

        setPostSending(true)
        setPostResult("")

        body := map[string]string{
            "title": postTitle,
            "body":  postBody,
        }

        resp, err := api.Post("/api/posts", api.Options{
            Headers: map[string]string{"Content-Type": "application/json"},
            Body:    api.ToJSON(body),
        })

        if err != nil {
            setPostResult(fmt.Sprintf("Error: %s", err))
            setPostSending(false)
            return
        }

        setPostResult(fmt.Sprintf("Created! Response: %s", resp))
        setPostSending(false)
        setPostTitle("")
        setPostBody("")
    }

    // --- Fetch todos with query params ---
    fetchTodos := func(e dom.Event) {
        setTodosLoading(true)

        data, err := api.Get("/api/todos", api.Options{
            Query: map[string]string{"_limit": "5", "completed": "false"},
        })

        if err != nil {
            setTodosLoading(false)
            return
        }

        parsed, _ := api.ParseJSON(data)
        var result []TodoItem
        for _, item := range parsed {
            t := TodoItem{
                ID:        item["id"].(int),
                Title:     item["title"].(string),
                Completed: item["completed"].(bool),
            }
            result = append(result, t)
        }
        setTodos(result)
        setTodosLoading(false)
    }

    // --- Render post list ---
    var postItems []dom.Element
    for _, p := range posts {
        postItems = append(postItems,
            dom.CreateElement("div", dom.Props{"class": "post-card"},
                dom.CreateElement("h3", nil, fmt.Sprintf("#%d %s", p.ID, p.Title)),
                dom.CreateElement("p", nil, p.Body),
            ),
        )
    }

    // --- Render todo items ---
    var todoItems []dom.Element
    for _, t := range todos {
        todoItems = append(todoItems,
            dom.CreateElement("div", dom.Props{"class": "todo-item"},
                dom.CreateElement("input", dom.Props{
                    "type":  "checkbox",
                    "class": "todo-check",
                }),
                dom.CreateElement("span", nil, fmt.Sprintf("%s (ID: %d)", t.Title, t.ID)),
            ),
        )
    }

    // --- Layout ---
    return dom.CreateElement("div", dom.Props{"class": "api-demo"},
        dom.CreateElement("h2", nil, "🌐 API Demo"),

        // Section 1: GET posts
        dom.CreateElement("section", dom.Props{"class": "demo-section"},
            dom.CreateElement("h3", nil, "GET — Fetch Posts"),
            dom.CreateElement("button", dom.Props{
                "class":   "btn",
                "onclick": fetchPosts,
            }, "Fetch Posts"),
            dom.CreateElement("div", dom.Props{"class": "results"}, postItems),
            dom.CreateElementIf(postsLoading, dom.CreateElement("p", dom.Props{"class": "loading"}, "Loading...")),
            dom.CreateElementIf(postsError != "", dom.CreateElement("p", dom.Props{"class": "error"}, postsError)),
        ),

        // Section 2: POST create
        dom.CreateElement("section", dom.Props{"class": "demo-section"},
            dom.CreateElement("h3", nil, "POST — Create Post"),
            dom.CreateElement("form", dom.Props{"class": "form", "onsubmit": createPost},
                dom.CreateElement("input", dom.Props{
                    "type":        "text",
                    "placeholder": "Post title",
                    "class":       "input",
                }),
                dom.CreateElement("textarea", dom.Props{
                    "placeholder": "Post body",
                    "class":       "textarea",
                }, ""),
                dom.CreateElement("button", dom.Props{
                    "type":  "submit",
                    "class": "btn primary",
                }, "Create Post"),
            ),
            dom.CreateElementIf(postSending, dom.CreateElement("p", dom.Props{"class": "loading"}, "Sending...")),
            dom.CreateElementIf(postResult != "", dom.CreateElement("p", dom.Props{"class": "result"}, postResult)),
        ),

        // Section 3: GET with query params
        dom.CreateElement("section", dom.Props{"class": "demo-section"},
            dom.CreateElement("h3", nil, "GET — Fetch Todos (with query params)"),
            dom.CreateElement("button", dom.Props{
                "class":   "btn",
                "onclick": fetchTodos,
            }, "Fetch Active Todos"),
            dom.CreateElement("div", dom.Props{"class": "results"}, todoItems),
            dom.CreateElementIf(todosLoading, dom.CreateElement("p", dom.Props{"class": "loading"}, "Loading...")),
        ),
    )
}

// Mount the API demo into #app.
func main() {
    dom.Mount("#app", ApiApp())
}
