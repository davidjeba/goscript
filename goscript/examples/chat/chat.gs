package main

import "goscript/dom"
import "goscript/state"
import "goscript/realtime"
import "goscript/fmt"

// Chat — a real-time chat interface demonstrating WebSocket connections,
// SSE notifications, event bus pub/sub, and live message rendering.
//
// Usage:
//
//	gopm compile chat.gs -o chat.js
//	<script src="/chat.js"></script>

// Message represents a single chat message.
struct Message {
    User      string
    Text      string
    Timestamp string
    IsSelf    bool
}

// ChatApp renders the chat interface with WebSocket and SSE integration.
func ChatApp() dom.Element {
    // --- State ---
    messages, setMessages := state.Use([]Message{})
    input, setInput := state.Use("")
    username, setUsername := state.Use("Anonymous")
    connected, setConnected := state.Use(false)
    onlineCount, setOnlineCount := state.Use(0)

    // --- WebSocket connection ---
    // Connect to the chat server via WebSocket for bi-directional messaging.
    ws := realtime.WebSocket("wss://example.com/chat")

    ws.OnConnect(func() {
        setConnected(true)
        // Send a join notification
        ws.Send(fmt.Sprintf(`{"type":"join","user":"%s"}`, username))
    })

    ws.OnDisconnect(func() {
        setConnected(false)
    })

    ws.OnMessage(func(data string) {
        // Parse incoming message and append to state
        msg := Message{
            User:      "Other",
            Text:      data,
            Timestamp: dom.Now(),
            IsSelf:    false,
        }
        setMessages(append(messages, msg))
    })

    // --- SSE for server notifications ---
    // Subscribe to typing indicators and system announcements.
    sse := realtime.SSE("/api/chat/events")

    sse.On("user_count", func(data string) {
        setOnlineCount(dom.ParseInt(data))
    })

    sse.On("user_typing", func(data string) {
        dom.SetText("#typing-indicator", fmt.Sprintf("%s is typing...", data))
    })

    sse.On("announcement", func(data string) {
        msg := Message{
            User:      "System",
            Text:      data,
            Timestamp: dom.Now(),
            IsSelf:    false,
        }
        setMessages(append(messages, msg))
    })

    // --- Event bus for cross-component communication ---
    // Other components can listen for "chat:send" events.
    bus := realtime.EventBus()

    // --- Helpers ---
    sendMessage := func(e dom.Event) {
        e.PreventDefault()
        text := dom.ValueOf("#chat-input")
        if text == "" {
            return
        }
        msg := Message{
            User:      username,
            Text:      text,
            Timestamp: dom.Now(),
            IsSelf:    true,
        }
        // Update local state immediately
        setMessages(append(messages, msg))
        // Send via WebSocket
        ws.Send(fmt.Sprintf(`{"type":"message","text":"%s"}`, text))
        // Emit on the event bus for other components
        bus.Emit("chat:send", text)
        // Clear input
        dom.SetValue("#chat-input", "")
        setInput("")
        // Stop typing indicator
        sse.EmitLocal("typing_stop", username)
    }

    // Notify server when user starts typing
    onTyping := func(e dom.Event) {
        sse.EmitLocal("typing_start", username)
    }

    // --- Render message bubbles ---
    var bubbles []dom.Element
    for _, m := range messages {
        bubbleClass := "chat-bubble other"
        if m.IsSelf {
            bubbleClass = "chat-bubble self"
        }
        // System messages get a special style
        if m.User == "System" {
            bubbleClass = "chat-bubble system"
        }
        bubble := dom.CreateElement("div", dom.Props{"class": bubbleClass},
            dom.CreateElement("span", dom.Props{"class": "msg-user"}, m.User),
            dom.CreateElement("span", dom.Props{"class": "msg-text"}, m.Text),
            dom.CreateElement("span", dom.Props{"class": "msg-time"}, m.Timestamp),
        )
        bubbles = append(bubbles, bubble)
    }

    // --- Status indicator ---
    statusText := "● Offline"
    statusClass := "status offline"
    if connected {
        statusText = fmt.Sprintf("● Online — %d users", onlineCount)
        statusClass = "status online"
    }

    // --- Layout ---
    return dom.CreateElement("div", dom.Props{"class": "chat-app"},
        // Header with connection status
        dom.CreateElement("div", dom.Props{"class": "chat-header"},
            dom.CreateElement("h2", nil, "💬 Goscript Chat"),
            dom.CreateElement("span", dom.Props{"class": statusClass}, statusText),
        ),

        // Message area
        dom.CreateElement("div", dom.Props{
            "id":    "chat-messages",
            "class": "chat-messages",
        }, bubbles),

        // Typing indicator
        dom.CreateElement("div", dom.Props{
            "id":    "typing-indicator",
            "class": "typing-indicator",
        }, ""),

        // Input form
        dom.CreateElement("form", dom.Props{
            "class":    "chat-form",
            "onsubmit": sendMessage,
        },
            dom.CreateElement("input", dom.Props{
                "id":          "chat-input",
                "type":        "text",
                "placeholder": "Type a message...",
                "class":       "chat-input",
                "oninput":     onTyping,
            }),
            dom.CreateElement("button", dom.Props{
                "type":  "submit",
                "class": "btn-send",
            }, "Send"),
        ),
    )
}

// Mount the chat app into #app.
func main() {
    dom.Mount("#app", ChatApp())
}
