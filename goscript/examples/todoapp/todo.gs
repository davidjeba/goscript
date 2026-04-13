package main

import "goscript/dom"
import "goscript/state"
import "goscript/fmt"

// TodoApp — a full todo list application built with goscript.
// Demonstrates: useState, DOM creation, event handling, form submission,
// list rendering, toggling completion, and filtering.
//
// Usage:
//
//	gopm compile todo.gs -o todo.js
//	<script src="/todo.js"></script>

// Todo represents a single todo item.
struct Todo {
    ID    int
    Text  string
    Done  bool
}

// TodoApp renders the complete todo list component.
func TodoApp() dom.Element {
    // --- State ---
    // Track the list of todos, the input value, and the active filter.
    todos, setTodos := state.Use([]Todo{})
    input, setInput := state.Use("")
    filter, setFilter := state.Use("all") // "all", "active", "done"

    // --- Derived state ---
    // Count active (non-done) todos for the summary bar.
    // In a real app this would use useMemo; here we compute inline.
    activeCount := 0
    for _, t := range todos {
        if !t.Done {
            activeCount = activeCount + 1
        }
    }

    // --- Helpers ---
    // addTodo handles form submission to create a new todo.
    addTodo := func(e dom.Event) {
        e.PreventDefault()
        text := dom.ValueOf("#todo-input")
        if text == "" {
            return
        }
        newTodo := Todo{ID: len(todos) + 1, Text: text, Done: false}
        setTodos(append(todos, newTodo))
        dom.SetValue("#todo-input", "")
        setInput("")
    }

    // toggleTodo flips the Done flag on a todo by ID.
    toggleTodo := func(id int) func(dom.Event) {
        return func(e dom.Event) {
            updated := []Todo{}
            for _, t := range todos {
                if t.ID == id {
                    updated = append(updated, Todo{ID: t.ID, Text: t.Text, Done: !t.Done})
                } else {
                    updated = append(updated, t)
                }
            }
            setTodos(updated)
        }
    }

    // removeTodo deletes a todo by ID.
    removeTodo := func(id int) func(dom.Event) {
        return func(e dom.Event) {
            filtered := []Todo{}
            for _, t := range todos {
                if t.ID != id {
                    filtered = append(filtered, t)
                }
            }
            setTodos(filtered)
        }
    }

    // clearCompleted removes all done todos.
    clearCompleted := func(e dom.Event) {
        remaining := []Todo{}
        for _, t := range todos {
            if !t.Done {
                remaining = append(remaining, t)
            }
        }
        setTodos(remaining)
    }

    // --- Render todo items ---
    var items []dom.Element
    for _, t := range todos {
        // Apply filter
        if filter == "active" && t.Done {
            continue
        }
        if filter == "done" && !t.Done {
            continue
        }
        className := "todo-item"
        if t.Done {
            className = className + " completed"
        }
        item := dom.CreateElement("li", dom.Props{"class": className},
            dom.CreateElement("span", dom.Props{
                "class":   "todo-text",
                "onclick": toggleTodo(t.ID),
            }, t.Text),
            dom.CreateElement("button", dom.Props{
                "class":   "btn-delete",
                "onclick": removeTodo(t.ID),
            }, "✕"),
        )
        items = append(items, item)
    }

    // --- Filter bar ---
    filterButton := func(label string, value string) dom.Element {
        active := ""
        if filter == value {
            active = " active"
        }
        return dom.CreateElement("button", dom.Props{
            "class":   "btn-filter" + active,
            "onclick": func(e dom.Event) { setFilter(value) },
        }, label)
    }

    // --- Layout ---
    return dom.CreateElement("div", dom.Props{"class": "todo-app"},
        // Header
        dom.CreateElement("h1", nil, "📝 Goscript Todo App"),

        // Input form
        dom.CreateElement("form", dom.Props{
            "id":       "todo-form",
            "class":    "todo-form",
            "onsubmit": addTodo,
        },
            dom.CreateElement("input", dom.Props{
                "id":          "todo-input",
                "type":        "text",
                "placeholder": "What needs to be done?",
                "class":       "todo-input",
            }),
            dom.CreateElement("button", dom.Props{
                "type":  "submit",
                "class": "btn-add",
            }, "Add"),
        ),

        // Todo list
        dom.CreateElement("ul", dom.Props{"class": "todo-list"}, items),

        // Footer with count and filters
        dom.CreateElement("div", dom.Props{"class": "todo-footer"},
            dom.CreateElement("span", dom.Props{"class": "todo-count"},
                fmt.Sprintf("%d item(s) left", activeCount),
            ),
            dom.CreateElement("div", dom.Props{"class": "filter-group"},
                filterButton("All", "all"),
                filterButton("Active", "active"),
                filterButton("Done", "done"),
            ),
            dom.CreateElement("button", dom.Props{
                "class":   "btn-clear",
                "onclick": clearCompleted,
            }, "Clear completed"),
        ),
    )
}

// Mount the todo app into #app on page load.
func main() {
    dom.Mount("#app", TodoApp())
}
