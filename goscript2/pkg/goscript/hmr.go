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
)

// DevServer provides a development server with HMR and live reload
type DevServer struct {
        port         int
        router       *AppRouter
        watchPaths   []string
        wsClients    map[string]*WSClient
        hmrEnabled   bool
        liveReload   bool
        mu           sync.RWMutex
        buildErrors  []BuildError
        lastBuild    time.Time
        onFileChange func(path string)
}

// WSClient represents a connected WebSocket client
type WSClient struct {
        ID       string
        messages chan []byte
        done     chan struct{}
}

// BuildError represents a build error
type BuildError struct {
        File     string `json:"file"`
        Line     int    `json:"line"`
        Message  string `json:"message"`
        Severity string `json:"severity"`
}

// HMRMessage represents a Hot Module Replacement message
type HMRMessage struct {
        Type    string      `json:"type"`
        Path    string      `json:"path,omitempty"`
        Hash    string      `json:"hash,omitempty"`
        Error   *BuildError `json:"error,omitempty"`
        Modules []string    `json:"modules,omitempty"`
}

// NewDevServer creates a new development server
func NewDevServer(port int, router *AppRouter) *DevServer {
        return &DevServer{
                port:       port,
                router:     router,
                wsClients:  make(map[string]*WSClient),
                hmrEnabled: true,
                liveReload: true,
        }
}

// Start begins the development server with HMR
func (ds *DevServer) Start(ctx context.Context) error {
        mux := http.NewServeMux()

        // HMR WebSocket endpoint (stub - would use gorilla/websocket in production)
        mux.HandleFunc("/__goscript_hmr", ds.handleHMR)
        // Dev overlay
        mux.HandleFunc("/__goscript_dev", ds.handleDevOverlay)
        // App routes
        mux.Handle("/", ds.router)

        addr := fmt.Sprintf(":%d", ds.port)
        log.Printf("GoScript Dev Server running at http://localhost:%d\n", ds.port)
        log.Printf("   HMR enabled: %v | Live reload: %v\n", ds.hmrEnabled, ds.liveReload)

        server := &http.Server{Addr: addr, Handler: mux}

        if len(ds.watchPaths) > 0 {
                go ds.watchFiles(ctx)
        }
        go ds.heartbeat(ctx)

        go func() {
                <-ctx.Done()
                shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
                defer cancel()
                server.Shutdown(shutdownCtx)
        }()

        return server.ListenAndServe()
}

// handleHMR handles WebSocket upgrade requests for HMR
func (ds *DevServer) handleHMR(w http.ResponseWriter, r *http.Request) {
        // WebSocket upgrade stub - in production, use gorilla/websocket or nhooyr.io/websocket
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(HMRMessage{
                Type: "connected",
                Path: r.URL.Path,
        })
}

// handleDevOverlay serves the development error overlay
func (ds *DevServer) handleDevOverlay(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`<!DOCTYPE html><html><head><title>GoScript Dev Overlay</title></head><body><div id="dev-overlay"><h2>GoScript Development</h2><p>HMR: enabled</p></div></body></html>`))
}

// watchFiles monitors the file system for changes
func (ds *DevServer) watchFiles(ctx context.Context) {
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()

        for {
                select {
                case <-ctx.Done():
                        return
                case <-ticker.C:
                        for _, watchPath := range ds.watchPaths {
                                filepath.Walk(watchPath, func(path string, info os.FileInfo, err error) error {
                                        if err != nil || info.IsDir() {
                                                return nil
                                        }
                                        if strings.HasSuffix(path, ".go") {
                                                modTime := info.ModTime()
                                                if modTime.After(ds.lastBuild) {
                                                        ds.onFileChangeDetected(path)
                                                }
                                        }
                                        return nil
                                })
                        }
                }
        }
}

// heartbeat sends periodic keepalive messages to connected clients
func (ds *DevServer) heartbeat(ctx context.Context) {
        ticker := time.NewTicker(30 * time.Second)
        defer ticker.Stop()

        for {
                select {
                case <-ctx.Done():
                        return
                case <-ticker.C:
                        ds.mu.RLock()
                        for _, client := range ds.wsClients {
                                select {
                                case client.messages <- []byte(`{"type":"ping"}`):
                                default:
                                }
                        }
                        ds.mu.RUnlock()
                }
        }
}

// broadcast sends a message to all connected WebSocket clients
func (ds *DevServer) broadcast(msg HMRMessage) {
        data, _ := json.Marshal(msg)
        ds.mu.RLock()
        defer ds.mu.RUnlock()

        for _, client := range ds.wsClients {
                select {
                case client.messages <- data:
                default:
                }
        }
}

// onFileChangeDetected handles file change events
func (ds *DevServer) onFileChangeDetected(path string) {
        log.Printf("File changed: %s", path)

        // Trigger rebuild
        ds.buildErrors = ds.rebuild(path)

        // Broadcast HMR update
        if ds.hmrEnabled {
                hash := fmt.Sprintf("%x", time.Now().UnixNano())[:8]
                msg := HMRMessage{Type: "update", Path: path, Hash: hash}
                if len(ds.buildErrors) > 0 {
                        msg.Type = "error"
                        msg.Error = &ds.buildErrors[0]
                }
                ds.broadcast(msg)
        }

        if ds.onFileChange != nil {
                ds.onFileChange(path)
        }
}

// rebuild checks for compilation errors (stub - in production would use go/packages)
func (ds *DevServer) rebuild(changedFile string) []BuildError {
        // Stub implementation - in production this would use go/packages
        // to actually compile and check for errors
        ds.mu.Lock()
        ds.lastBuild = time.Now()
        ds.mu.Unlock()

        return nil
}

// SetWatchPaths configures the directories to watch for file changes
func (ds *DevServer) SetWatchPaths(paths []string) {
        ds.watchPaths = paths
}

// OnFileChange sets a callback for file change events
func (ds *DevServer) OnFileChange(fn func(path string)) {
        ds.onFileChange = fn
}

// GetBuildErrors returns the current build errors
func (ds *DevServer) GetBuildErrors() []BuildError {
        ds.mu.RLock()
        defer ds.mu.RUnlock()
        return ds.buildErrors
}

// InjectDevScript injects the HMR client script into HTML
func InjectDevScript(html string, port int) string {
        script := fmt.Sprintf(`<script>
(function(){
  var ws = new WebSocket("ws://localhost:%d/__goscript_hmr");
  ws.onmessage = function(e) {
    var msg = JSON.parse(e.data);
    if (msg.type === "update") {
      console.log("[GoScript HMR] Module updated:", msg.path);
      if (window.__goscript_hmr) window.__goscript_hmr(msg);
    }
    if (msg.type === "error") {
      console.error("[GoScript] Build error:", msg.error);
      if (window.__goscript_overlay) window.__goscript_overlay(msg.error);
    }
  };
  ws.onclose = function() { setTimeout(function(){ location.reload(); }, 2000); };
})();
</script>`, port)

        // Inject before closing </head> or at end of <head>
        if idx := strings.Index(html, "</head>"); idx != -1 {
                return html[:idx] + script + "\n" + html[idx:]
        }
        return html + script
}
