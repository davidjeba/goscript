package goscript

import (
        "context"
        "fmt"
        "net/http"
        "os"
        "path/filepath"
        "strings"
        "sync"
        "time"
)

// RenderMode defines the rendering strategy for a page.
type RenderMode int

const (
        // RenderSSG indicates static site generation: the page is pre-rendered at
        // build time and served as a static HTML file.
        RenderSSG RenderMode = iota
        // RenderSSR indicates server-side rendering: the page is rendered on every
        // request for fully dynamic content.
        RenderSSR
        // RenderISR indicates incremental static regeneration: the page is served
        // from a static cache and revalidated in the background at the specified interval.
        RenderISR
)

// PageConfig defines the configuration for a single page in the SSG/ISR system.
type PageConfig struct {
        Path       string
        Component  Component
        RenderMode RenderMode
        Revalidate time.Duration
        Params     []map[string]string
}

// isrEntry holds the cached HTML and metadata for an ISR page.
type isrEntry struct {
        HTML         string
        LastRendered time.Time
}

// SSGEngine handles static site generation and incremental static regeneration.
// It can pre-render pages at build time (SSG), serve them statically, and
// revalidate ISR pages on a configurable timer.
type SSGEngine struct {
        outputDir     string
        pages         []PageConfig
        staticCache   map[string]string
        isrCache      map[string]*isrEntry
        store         *Store
        mutex         sync.RWMutex
        isrCancel     context.CancelFunc
}

// NewSSGEngine creates a new SSGEngine that writes generated HTML files to the
// specified output directory. If the directory does not exist, it will be created.
func NewSSGEngine(outputDir string) *SSGEngine {
        return &SSGEngine{
                outputDir:   outputDir,
                pages:       make([]PageConfig, 0),
                staticCache: make(map[string]string),
                isrCache:    make(map[string]*isrEntry),
                store:       NewStore(),
        }
}

// AddPage registers a page configuration with the engine. Pages can be configured
// for SSG, SSR, or ISR rendering modes.
func (e *SSGEngine) AddPage(config PageConfig) {
        e.mutex.Lock()
        defer e.mutex.Unlock()
        e.pages = append(e.pages, config)
}

// Build performs a full static build for all SSG pages and pre-renders ISR pages.
// It creates the output directory structure and writes HTML files for each page.
// Pages with Params will be generated for each parameter set.
func (e *SSGEngine) Build(ctx context.Context) error {
        e.mutex.Lock()
        defer e.mutex.Unlock()

        // Ensure the output directory exists
        if err := os.MkdirAll(e.outputDir, 0755); err != nil {
                return fmt.Errorf("failed to create output directory: %w", err)
        }

        for _, page := range e.pages {
                switch page.RenderMode {
                case RenderSSG, RenderISR:
                        if len(page.Params) > 0 {
                                // Generate a page for each parameter set
                                for _, params := range page.Params {
                                        paramPath := buildParamPath(page.Path, params)
                                        html := page.Component.Render()
                                        filePath := filepath.Join(e.outputDir, paramPath, "index.html")
                                        if err := writeFile(filePath, html); err != nil {
                                                return fmt.Errorf("failed to write %s: %w", filePath, err)
                                        }
                                        e.staticCache[paramPath] = html
                                        if page.RenderMode == RenderISR {
                                                e.isrCache[paramPath] = &isrEntry{
                                                        HTML:         html,
                                                        LastRendered: time.Now(),
                                                }
                                        }
                                }
                        } else {
                                html := page.Component.Render()
                                filePath := filepath.Join(e.outputDir, page.Path, "index.html")
                                if err := writeFile(filePath, html); err != nil {
                                        return fmt.Errorf("failed to write %s: %w", filePath, err)
                                }
                                e.staticCache[page.Path] = html
                                if page.RenderMode == RenderISR {
                                        e.isrCache[page.Path] = &isrEntry{
                                                HTML:         html,
                                                LastRendered: time.Now(),
                                        }
                                }
                        }
                case RenderSSR:
                        // SSR pages are not pre-rendered at build time
                }
        }

        // Start ISR background revalidation
        e.startISRRevalidation(ctx)

        return nil
}

// ServeSSG serves a pre-built static page for the given request path. If the
// page is found in the static cache, it is served directly. For ISR pages that
// have exceeded their revalidation interval, a stale-while-revalidate strategy
// is used: the stale page is served immediately while a background re-render is
// triggered.
func (e *SSGEngine) ServeSSG(w http.ResponseWriter, r *http.Request) {
        path := r.URL.Path
        e.mutex.RLock()
        defer e.mutex.RUnlock()

        // Check ISR cache first for revalidation eligibility
        if entry, ok := e.isrCache[path]; ok {
                // Check if the page needs revalidation
                page := e.findPage(path)
                if page != nil && time.Since(entry.LastRendered) > page.Revalidate {
                        // Serve stale content while revalidating in the background
                        w.Header().Set("Content-Type", "text/html; charset=utf-8")
                        w.Header().Set("X-GoScript-Cache", "stale")
                        fmt.Fprint(w, entry.HTML)

                        // Trigger background revalidation
                        go e.revalidatePage(path, page)
                        return
                }

                w.Header().Set("Content-Type", "text/html; charset=utf-8")
                w.Header().Set("X-GoScript-Cache", "hit")
                fmt.Fprint(w, entry.HTML)
                return
        }

        // Check static cache
        if html, ok := e.staticCache[path]; ok {
                w.Header().Set("Content-Type", "text/html; charset=utf-8")
                w.Header().Set("X-GoScript-Cache", "hit")
                fmt.Fprint(w, html)
                return
        }

        // Try to serve from the file system
        filePath := filepath.Join(e.outputDir, path, "index.html")
        data, err := os.ReadFile(filePath)
        if err != nil {
                w.Header().Set("Content-Type", "text/html; charset=utf-8")
                w.Header().Set("X-GoScript-Cache", "miss")
                w.WriteHeader(http.StatusNotFound)
                fmt.Fprint(w, "<h1>404 - Page Not Found</h1>")
                return
        }

        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        w.Header().Set("X-GoScript-Cache", "disk")
        fmt.Fprint(w, string(data))
}

// startISRRevalidation starts a background goroutine that periodically checks
// ISR pages and re-renders those that have exceeded their revalidation interval.
func (e *SSGEngine) startISRRevalidation(ctx context.Context) {
        childCtx, cancel := context.WithCancel(ctx)
        e.isrCancel = cancel

        go func() {
                ticker := time.NewTicker(10 * time.Second)
                defer ticker.Stop()

                for {
                        select {
                        case <-childCtx.Done():
                                return
                        case <-ticker.C:
                                e.revalidateStalePages()
                        }
                }
        }()
}

// revalidateStalePages checks all ISR pages and re-renders those that have
// exceeded their revalidation interval.
func (e *SSGEngine) revalidateStalePages() {
        e.mutex.Lock()
        defer e.mutex.Unlock()

        for path, entry := range e.isrCache {
                page := e.findPage(path)
                if page == nil {
                        continue
                }
                if time.Since(entry.LastRendered) > page.Revalidate {
                        html := page.Component.Render()
                        e.isrCache[path] = &isrEntry{
                                HTML:         html,
                                LastRendered: time.Now(),
                        }
                        e.staticCache[path] = html

                        filePath := filepath.Join(e.outputDir, path, "index.html")
                        _ = writeFile(filePath, html)
                }
        }
}

// revalidatePage re-renders a single ISR page in the background.
func (e *SSGEngine) revalidatePage(path string, page *PageConfig) {
        html := page.Component.Render()

        e.mutex.Lock()
        e.isrCache[path] = &isrEntry{
                HTML:         html,
                LastRendered: time.Now(),
        }
        e.staticCache[path] = html
        e.mutex.Unlock()

        filePath := filepath.Join(e.outputDir, path, "index.html")
        _ = writeFile(filePath, html)
}

// findPage locates a PageConfig by its path.
func (e *SSGEngine) findPage(path string) *PageConfig {
        for i := range e.pages {
                if e.pages[i].Path == path {
                        return &e.pages[i]
                }
        }
        return nil
}

// buildParamPath replaces dynamic segments in the path with parameter values.
// For example, "/posts/:id" with params {"id": "42"} becomes "/posts/42".
func buildParamPath(pathTemplate string, params map[string]string) string {
        result := pathTemplate
        for key, value := range params {
                result = strings.Replace(result, ":"+key, value, 1)
                result = strings.Replace(result, "*"+key, value, 1)
        }
        return result
}

// writeFile writes content to a file, creating parent directories as needed.
func writeFile(path, content string) error {
        dir := filepath.Dir(path)
        if err := os.MkdirAll(dir, 0755); err != nil {
                return err
        }
        return os.WriteFile(path, []byte(content), 0644)
}

// GetStats returns statistics about the SSG engine's cache state.
func (e *SSGEngine) GetStats() map[string]interface{} {
        e.mutex.RLock()
        defer e.mutex.RUnlock()

        return map[string]interface{}{
                "total_pages":    len(e.pages),
                "static_cached":  len(e.staticCache),
                "isr_cached":     len(e.isrCache),
                "ssg_pages":      countPagesByMode(e.pages, RenderSSG),
                "ssr_pages":      countPagesByMode(e.pages, RenderSSR),
                "isr_pages":      countPagesByMode(e.pages, RenderISR),
        }
}

// countPagesByMode counts pages by their RenderMode.
func countPagesByMode(pages []PageConfig, mode RenderMode) int {
        count := 0
        for _, p := range pages {
                if p.RenderMode == mode {
                        count++
                }
        }
        return count
}
