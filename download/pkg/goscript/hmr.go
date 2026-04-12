// Package goscript provides Hot Module Replacement (HMR) for the GoScript development
// server. It watches source directories for file changes and broadcasts update
// notifications to connected browser clients via WebSocket connections, enabling
// a rapid feedback loop during development without full page reloads.
package goscript

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// hmrMessage represents a WebSocket message sent from the dev server to the
// browser client to notify it of file changes or server restarts.
type hmrMessage struct {
	Type    string `json:"type"`
	Path    string `json:"path,omitempty"`
	Action  string `json:"action,omitempty"`
	Message string `json:"message,omitempty"`
}

// DevServer provides a development server with Hot Module Replacement support.
// It wraps an existing http.Handler with file watching and WebSocket-based HMR
// notifications, allowing browser clients to update incrementally when source
// files change.
type DevServer struct {
	port       int
	router     http.Handler
	watchers   map[string]bool
	callbacks  []func(path string)
	clients    map[*websocket.Conn]bool
	clientsMu  sync.RWMutex
	wsUpgrader websocket.Upgrader
	httpServer *http.Server
	mutex      sync.Mutex
}

// NewDevServer creates a new DevServer on the specified port that serves the
// given router. The server provides HMR WebSocket connections at the /__hmr
// endpoint and proxies all other requests to the provided handler.
func NewDevServer(port int, router http.Handler) *DevServer {
	return &DevServer{
		port:     port,
		router:   router,
		watchers: make(map[string]bool),
		callbacks: make([]func(path string), 0),
		clients:  make(map[*websocket.Conn]bool),
		wsUpgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// Watch adds directories to the file watch list. Changes to any file in these
// directories will trigger HMR notifications to connected clients. The watcher
// uses a polling mechanism compatible with all platforms.
func (ds *DevServer) Watch(dirs ...string) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	for _, dir := range dirs {
		absPath, err := filepath.Abs(dir)
		if err != nil {
			log.Printf("[HMR] warning: cannot resolve path %s: %v", dir, err)
			continue
		}
		if _, exists := ds.watchers[absPath]; !exists {
			ds.watchers[absPath] = true
		}
	}
}

// OnFileChange registers a callback function that is invoked whenever a file
// change is detected in any of the watched directories. Multiple callbacks can
// be registered and they are invoked in registration order.
func (ds *DevServer) OnFileChange(callback func(path string)) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	ds.callbacks = append(ds.callbacks, callback)
}

// Start begins serving the development server and starts file watching. The
// server runs until the provided context is cancelled. It listens for HTTP
// requests and manages WebSocket connections for HMR.
func (ds *DevServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	// Mount the HMR WebSocket endpoint
	mux.HandleFunc("/__hmr", ds.handleHMRWebSocket)

	// Mount the user's router for all other paths
	mux.Handle("/", ds.router)

	addr := fmt.Sprintf(":%d", ds.port)
	ds.httpServer = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Start file watching in a background goroutine
	go ds.watchFiles(ctx)

	// Start the HTTP server in a background goroutine
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("[HMR] Dev server starting on http://localhost:%d", ds.port)
		if err := ds.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
		close(serverErr)
	}()

	// Start a goroutine to handle client cleanup
	go ds.cleanupClients(ctx)

	select {
	case <-ctx.Done():
		log.Println("[HMR] Shutting down dev server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ds.httpServer.Shutdown(shutdownCtx)
		return nil
	case err := <-serverErr:
		return err
	}
}

// handleHMRWebSocket handles incoming WebSocket upgrade requests for the HMR channel.
func (ds *DevServer) handleHMRWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := ds.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[HMR] WebSocket upgrade failed: %v", err)
		return
	}

	ds.clientsMu.Lock()
	ds.clients[conn] = true
	ds.clientsMu.Unlock()

	log.Printf("[HMR] Client connected from %s (total: %d)", conn.RemoteAddr(), ds.clientCount())

	// Send a welcome message
	welcomeMsg := hmrMessage{
		Type:    "connected",
		Message: "GoScript HMR connected",
	}
	conn.WriteJSON(welcomeMsg)

	// Read loop to keep the connection alive and detect disconnects
	defer func() {
		ds.clientsMu.Lock()
		delete(ds.clients, conn)
		ds.clientsMu.Unlock()
		conn.Close()
		log.Printf("[HMR] Client disconnected (total: %d)", ds.clientCount())
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// broadcast sends an HMR message to all connected WebSocket clients.
func (ds *DevServer) broadcast(msg hmrMessage) {
	ds.clientsMu.RLock()
	defer ds.clientsMu.RUnlock()

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[HMR] Failed to marshal message: %v", err)
		return
	}

	for client := range ds.clients {
		err := client.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Printf("[HMR] Failed to send to client: %v", err)
		}
	}
}

// clientCount returns the current number of connected WebSocket clients.
func (ds *DevServer) clientCount() int {
	ds.clientsMu.RLock()
	defer ds.clientsMu.RUnlock()
	return len(ds.clients)
}

// watchFiles polls watched directories for changes at 500ms intervals. When a
// change is detected, it invokes all registered callbacks and broadcasts an
// HMR notification to connected clients.
func (ds *DevServer) watchFiles(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	snapshots := make(map[string]os.FileInfo)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ds.mutex.Lock()
			dirs := make([]string, 0, len(ds.watchers))
			for dir := range ds.watchers {
				dirs = append(dirs, dir)
			}
			callbacks := make([]func(path string), len(ds.callbacks))
			copy(callbacks, ds.callbacks)
			ds.mutex.Unlock()

			for _, dir := range dirs {
				entries, err := os.ReadDir(dir)
				if err != nil {
					continue
				}

				for _, entry := range entries {
					if entry.IsDir() {
						continue
					}

					fullPath := filepath.Join(dir, entry.Name())
					info, err := entry.Info()
					if err != nil {
						continue
					}

					key := fullPath
					prevInfo, exists := snapshots[key]
					if exists {
						if !prevInfo.ModTime().Equal(info.ModTime()) || prevInfo.Size() != info.Size() {
							snapshots[key] = info
							ds.handleFileChange(fullPath, callbacks)
						}
					} else {
						snapshots[key] = info
					}
				}
			}
		}
	}
}

// handleFileChange processes a detected file change by invoking callbacks and
// broadcasting HMR notifications. Go source files trigger a rebuild notification,
// while other files trigger a simple update notification.
func (ds *DevServer) handleFileChange(path string, callbacks []func(path string)) {
	log.Printf("[HMR] File changed: %s", path)

	// Invoke registered callbacks
	for _, cb := range callbacks {
		cb(path)
	}

	// Determine the action based on file type
	action := "update"
	if strings.HasSuffix(path, ".go") {
		action = "reload"
	}

	// Broadcast to all connected clients
	ds.broadcast(hmrMessage{
		Type:   "change",
		Path:   path,
		Action: action,
	})
}

// cleanupClients periodically checks for stale WebSocket connections and removes them.
func (ds *DevServer) cleanupClients(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ds.clientsMu.Lock()
			for client := range ds.clients {
				// Send a ping to check if the connection is alive
				err := client.WriteMessage(websocket.PingMessage, nil)
				if err != nil {
					client.Close()
					delete(ds.clients, client)
				}
			}
			ds.clientsMu.Unlock()
		}
	}
}
