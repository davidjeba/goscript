// Package realtime provides real-time communication functions for goscript .gs files.
// These functions are compiled to JavaScript calls to __gs.realtime by the GS compiler.
//
// The package supports two real-time communication patterns:
//   - Server-Sent Events (SSE): one-way server-to-client streaming
//   - WebSocket: bidirectional full-duplex communication
//
// Additionally, it provides an event bus (On/Emit/Off) for publish/subscribe
// messaging within the client. The Go server can trigger client-side events
// via the GS-Trigger HTTP response header, which the runtime intercepts
// automatically on all API responses.
//
// This package serves as API documentation and type definitions for .gs developers.
// The functions listed here map to JavaScript operations in the goscript client runtime
// (pkg/gslib/runtime.js). You do not need to import this package in Go server code —
// it exists solely for .gs files.
//
// # Usage in .gs code
//
//	import "goscript/realtime"
//
//	func StreamUpdates() {
//	    // Server-Sent Events
//	    conn := realtime.SSE("/api/events", func(data interface{}) {
//	        fmt.Println("Server event:", data)
//	        dom.SetTextContent(statusEl, data.(string))
//	    })
//
//	    // Close when done
//	    conn.Close()
//	}
//
//	func Chat() {
//	    // WebSocket
//	    ws := realtime.WebSocket("/ws/chat", map[string]interface{}{
//	        "onmessage": func(data interface{}) {
//	            fmt.Println("Message:", data)
//	        },
//	        "onopen": func() {
//	            fmt.Println("Connected!")
//	        },
//	    })
//
//	    // Send a message
//	    ws.Send(map[string]interface{}{
//	        "text": "Hello, World!",
//	    })
//
//	    // Close connection
//	    ws.Close()
//	}
package realtime

// SSEConnection represents an active Server-Sent Events connection.
// It provides a Close method to terminate the connection.
type SSEConnection struct {
	Close func()
}

// WebSocketConnection represents an active WebSocket connection.
// It provides Send and Close methods for communication.
type WebSocketConnection struct {
	Send  func(data interface{})
	Close func()
}

// Unsubscriber is a function that cancels a previously registered event subscription.
// Call it to stop receiving events for a given handler.
type Unsubscriber func()

// ---------------------------------------------------------------------------
// Server-Sent Events (SSE)
// ---------------------------------------------------------------------------

// SSE opens a Server-Sent Events connection to the specified URL.
// The handler function is called for each message received from the server.
// Messages are automatically parsed from JSON; if parsing fails, the raw
// string is passed to the handler.
//
// The connection is automatically closed on error (the browser's built-in
// auto-reconnect behavior is disabled to avoid log spam).
//
// In .gs code this compiles to __gs.sse(url, handler).
//
// Parameters:
//   - url: the SSE endpoint URL (e.g. "/api/events", "/stream/notifications")
//   - handler: a function called with each message. The data argument is
//     the parsed JSON object (or raw string if JSON parsing fails).
//
// Returns an SSEConnection with a Close() method to terminate the connection.
//
// Example (.gs):
//
//	conn := realtime.SSE("/api/stock-prices", func(data interface{}) {
//	    price := data.(map[string]interface{})
//	    dom.SetTextContent(priceEl, fmt.Sprintf("$%s", price["value"]))
//	})
//
//	// Later, close the connection
//	conn.Close()
func SSE(url string, handler func(data interface{})) SSEConnection { return SSEConnection{} }

// ---------------------------------------------------------------------------
// WebSocket
// ---------------------------------------------------------------------------

// WebSocket opens a bidirectional WebSocket connection to the specified URL.
//
// The handlers map can contain the following optional keys:
//   - "onopen":    func() — called when the connection is established
//   - "onmessage": func(data interface{}) — called for each received message
//   - "onclose":   func() — called when the connection is closed
//   - "onerror":   func() — called when an error occurs
//
// Messages are automatically parsed from JSON; if parsing fails, the raw
// string is passed to onmessage. When sending data via Send(), objects are
// automatically JSON-stringified, and strings are sent as-is.
//
// Note: WebSocket connections require manual reconnect logic. The runtime
// does not auto-reconnect WebSockets.
//
// In .gs code this compiles to __gs.ws(url, handlers).
//
// Parameters:
//   - url: the WebSocket endpoint URL (e.g. "/ws/chat", "wss://example.com/ws")
//   - handlers: a map of event handler functions
//
// Returns a WebSocketConnection with Send(data) and Close() methods.
//
// Example (.gs):
//
//	ws := realtime.WebSocket("/ws/chat", map[string]interface{}{
//	    "onopen": func() {
//	        fmt.Println("WebSocket connected")
//	    },
//	    "onmessage": func(data interface{}) {
//	        msg := data.(map[string]interface{})
//	        fmt.Println("Received:", msg["text"])
//	    },
//	    "onclose": func() {
//	        fmt.Println("WebSocket disconnected")
//	    },
//	    "onerror": func() {
//	        fmt.Println("WebSocket error")
//	    },
//	})
//
//	// Send a message
//	ws.Send(map[string]interface{}{
//	    "type": "chat",
//	    "text": "Hello!",
//	})
//
//	// Send raw string
//	ws.Send("ping")
//
//	// Close the connection
//	ws.Close()
func WebSocket(url string, handlers map[string]interface{}) WebSocketConnection {
	return WebSocketConnection{}
}

// ---------------------------------------------------------------------------
// Event Bus (Publish/Subscribe)
// ---------------------------------------------------------------------------

// On subscribes to a named event. The handler is called whenever the event
// is emitted (via Emit, or automatically via GS-Trigger response headers).
//
// Errors in handlers are caught and logged so one bad listener doesn't
// break the event chain.
//
// In .gs code this compiles to __gs.on(event, handler).
//
// Parameters:
//   - event: the event name (e.g. "cart:updated", "notification:received")
//   - handler: a function called with the event detail/payload
//
// Returns an Unsubscriber function. Call it to stop receiving events.
//
// Example (.gs):
//
//	unsub := realtime.On("user:login", func(detail interface{}) {
//	    user := detail.(map[string]interface{})
//	    fmt.Println("User logged in:", user["name"])
//	    updateUI(user)
//	})
//
//	// Later, unsubscribe
//	unsub()
func On(event string, handler func(detail interface{})) Unsubscriber {
	return func() {}
}

// Emit publishes a named event to all subscribers. Errors in handlers are
// caught and logged to the console.
//
// In .gs code this compiles to __gs.emit(event, detail).
//
// Parameters:
//   - event: the event name
//   - detail: the event payload (any JSON-serializable value)
//
// Example (.gs):
//
//	realtime.Emit("notification:new", map[string]interface{}{
//	    "message": "New order received!",
//	    "type":    "order",
//	})
//
//	// Emit with no payload
//	realtime.Emit("refresh", nil)
func Emit(event string, detail interface{}) {}

// Off unsubscribes a specific handler from a named event.
// This is the manual equivalent of calling the Unsubscriber returned by On.
//
// In .gs code this compiles to __gs.off(event, handler).
//
// Parameters:
//   - event: the event name
//   - handler: the same function reference passed to On
//
// Example (.gs):
//
//	handleNotification := func(detail interface{}) {
//	    fmt.Println("Notification:", detail)
//	}
//
//	realtime.On("notification", handleNotification)
//	// Later...
//	realtime.Off("notification", handleNotification)
func Off(event string, handler func(detail interface{})) {}
