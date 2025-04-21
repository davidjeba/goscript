# GoStore: High-Performance State Management for GoScript

GoStore is a comprehensive state management solution for GoScript applications, designed to handle form data, interactive elements, and application state with exceptional performance.

## Core Features

### In-Memory Data Collection

GoStore provides efficient in-memory data storage and retrieval with optimized data structures for frontend applications.

```go
// Create a new store
store := gostore.New()

// Set values
store.Set("user", map[string]interface{}{
    "name": "John Doe",
    "email": "john@example.com",
    "preferences": map[string]interface{}{
        "theme": "dark",
        "notifications": true,
    },
})

// Get values
user := store.Get("user")
theme := store.GetPath("user.preferences.theme") // "dark"

// Watch for changes
store.Watch("user.preferences.theme", func(oldValue, newValue interface{}) {
    fmt.Printf("Theme changed from %s to %s\n", oldValue, newValue)
})
```

### Form State Management

Comprehensive utilities for managing form state, including input values, validation, and submission.

```go
// Create a form store
form := gostore.NewForm("registration")

// Initialize form fields
form.InitFields(map[string]interface{}{
    "username": "",
    "email": "",
    "password": "",
    "confirmPassword": "",
})

// Add validation rules
form.AddValidation("username", gostore.Required(), gostore.MinLength(3))
form.AddValidation("email", gostore.Required(), gostore.Email())
form.AddValidation("password", gostore.Required(), gostore.MinLength(8))
form.AddValidation("confirmPassword", gostore.Required(), gostore.Matches("password"))

// Handle form submission
form.OnSubmit(func(values map[string]interface{}, isValid bool) {
    if isValid {
        // Process form submission
        api.RegisterUser(values)
    }
})

// Connect form to UI
gosky.Component("RegistrationForm", func(props map[string]interface{}) string {
    return `
        <form id="registration" onsubmit="{{.form.HandleSubmit}}">
            <div class="form-group">
                <label for="username">Username</label>
                <input 
                    type="text" 
                    id="username" 
                    value="{{.form.values.username}}" 
                    onchange="{{.form.HandleChange}}" 
                />
                {{if .form.errors.username}}
                    <div class="error">{{.form.errors.username}}</div>
                {{end}}
            </div>
            <!-- Other form fields -->
            <button type="submit">Register</button>
        </form>
    `
}, map[string]interface{}{
    "form": form,
})
```

### Validation System

Powerful validation system with built-in validators and support for custom validation logic.

```go
// Built-in validators
gostore.Required()
gostore.MinLength(5)
gostore.MaxLength(100)
gostore.Email()
gostore.URL()
gostore.Pattern(`^[A-Z][a-z]+$`)
gostore.Min(0)
gostore.Max(100)
gostore.OneOf([]string{"option1", "option2", "option3"})

// Custom validators
gostore.Custom(func(value interface{}, formValues map[string]interface{}) (bool, string) {
    // Custom validation logic
    if someCondition {
        return false, "Custom error message"
    }
    return true, ""
})

// Async validators
gostore.Async(func(value interface{}, callback func(bool, string)) {
    // Async validation, e.g., checking if username is available
    go func() {
        available := checkUsernameAvailability(value.(string))
        if available {
            callback(true, "")
        } else {
            callback(false, "Username is already taken")
        }
    }()
})
```

### State Persistence

Optional persistence for state recovery across page reloads or application restarts.

```go
// Create a persistent store
store := gostore.New(gostore.Options{
    Persistence: gostore.PersistenceOptions{
        Enabled: true,
        Storage: gostore.LocalStorage,
        Key: "app-state",
        Encrypt: true,
        EncryptionKey: "app-secret-key",
        ExpiresIn: 24 * time.Hour,
    },
})

// State will be automatically persisted and restored
```

### Change Tracking

Efficient tracking of state changes for optimized rendering and undo/redo functionality.

```go
// Enable change tracking
store := gostore.New(gostore.Options{
    TrackChanges: true,
})

// Make changes
store.Set("count", 1)
store.Set("count", 2)
store.Set("count", 3)

// Undo changes
store.Undo() // count is now 2
store.Undo() // count is now 1

// Redo changes
store.Redo() // count is now 2

// Get change history
history := store.GetHistory()
```

### Optimistic Updates

Support for optimistic UI updates with automatic rollback on failure.

```go
// Perform an optimistic update
store.Optimistic("todos", func(todos []interface{}) []interface{} {
    // Optimistically add a new todo
    newTodo := map[string]interface{}{
        "id": "temp-" + uuid.New().String(),
        "text": "New todo item",
        "completed": false,
    }
    return append(todos, newTodo)
}, func() {
    // Actual API call
    api.AddTodo("New todo item", func(success bool, result interface{}) {
        if !success {
            // Automatically rolls back if this callback returns an error
            return fmt.Errorf("failed to add todo")
        }
        // Update with the real ID from the server
        store.UpdateWhere("todos", func(todo map[string]interface{}) bool {
            return strings.HasPrefix(todo["id"].(string), "temp-")
        }, func(todo map[string]interface{}) map[string]interface{} {
            todo["id"] = result.(map[string]interface{})["id"]
            return todo
        })
        return nil
    })
})
```

## Advanced Features

### Computed Properties

Define properties that are automatically derived from other state values.

```go
store := gostore.New()

// Set base values
store.Set("items", []interface{}{
    map[string]interface{}{"price": 10, "quantity": 2},
    map[string]interface{}{"price": 15, "quantity": 1},
    map[string]interface{}{"price": 20, "quantity": 3},
})

// Define computed property
store.Compute("totalPrice", func(state map[string]interface{}) interface{} {
    items := state["items"].([]interface{})
    total := 0.0
    for _, item := range items {
        itemMap := item.(map[string]interface{})
        total += itemMap["price"].(float64) * itemMap["quantity"].(float64)
    }
    return total
})

// Access computed property
fmt.Println(store.Get("totalPrice")) // 95.0

// Computed property automatically updates when dependencies change
store.UpdateAt("items.0.quantity", 3)
fmt.Println(store.Get("totalPrice")) // 105.0
```

### Middleware System

Extend store functionality with middleware for logging, analytics, etc.

```go
// Create a store with middleware
store := gostore.New(gostore.Options{
    Middleware: []gostore.Middleware{
        // Logging middleware
        func(next gostore.DispatchFunction) gostore.DispatchFunction {
            return func(action gostore.Action) interface{} {
                fmt.Printf("Dispatching action: %s\n", action.Type)
                result := next(action)
                fmt.Printf("New state: %v\n", store.GetState())
                return result
            }
        },
        // Analytics middleware
        func(next gostore.DispatchFunction) gostore.DispatchFunction {
            return func(action gostore.Action) interface{} {
                if action.Type == "SET" {
                    analytics.TrackEvent("state_change", map[string]interface{}{
                        "path": action.Payload["path"],
                    })
                }
                return next(action)
            }
        },
    },
})
```

### Selectors

Create optimized selectors for derived data with memoization.

```go
// Create a selector
completedTodosSelector := gostore.CreateSelector(
    func(state map[string]interface{}) interface{} {
        return state["todos"]
    },
    func(todos []interface{}) interface{} {
        completed := []interface{}{}
        for _, todo := range todos {
            todoMap := todo.(map[string]interface{})
            if todoMap["completed"].(bool) {
                completed = append(completed, todo)
            }
        }
        return completed
    },
)

// Use the selector
completedTodos := completedTodosSelector(store.GetState())
```

## Integration with GoScript Ecosystem

GoStore is designed to work seamlessly with other GoScript components:

- **GoSky**: For state management during SSR and hydration
- **Gocsx**: For state-driven styling
- **GoConnect**: For managing API request state
- **Jetpack**: For monitoring state performance

## Getting Started

```bash
# Install GoStore using GOPM
gopm get github.com/davidjeba/goscript/pkg/gostore

# Create a new project with GoStore
gopm init myproject --template gostore
```

## Performance Considerations

GoStore is optimized for frontend performance:

- **Immutable Data Structures**: Efficient change detection
- **Selective Updates**: Only update components affected by state changes
- **Batched Updates**: Multiple state changes are batched for performance
- **Memory Optimization**: Efficient memory usage for large state trees
- **Lazy Evaluation**: Computed properties are evaluated only when needed