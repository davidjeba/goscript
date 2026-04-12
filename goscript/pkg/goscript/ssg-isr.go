package goscript

import (
        "context"
        "fmt"
        "io/ioutil"
        "net/http"
        "os"
        "path/filepath"
        "strings"
        "sync"
        "time"
)

// PageRenderMode determines how a page is rendered
type PageRenderMode int

const (
        RenderSSG PageRenderMode = iota // Static site generation at build time
        RenderSSR                       // Server-side rendering on every request
        RenderISR                       // Incremental static regeneration
)

// PageConfig defines how a page should be rendered
type PageConfig struct {
        Path       string
        Component  Component
        RenderMode PageRenderMode
        Revalidate time.Duration // ISR revalidation interval
        Params     []map[string]string
        Headers    map[string]string
        Priority   float64
}

// StaticPage represents a pre-rendered HTML page
type StaticPage struct {
        Path         string
        HTML         string
        CreatedAt    time.Time
        RevalidateAt time.Time
        Metadata     PageMetadata
}

// PageMetadata stores SEO metadata for a page
type PageMetadata struct {
        Title       string
        Description string
        Canonical   string
        OGImage     string
        NoIndex     bool
        JSONLD      map[string]interface{}
}

// toHeaders converts PageMetadata to HTTP headers
func (pm PageMetadata) toHeaders() map[string]string {
        headers := make(map[string]string)
        if pm.Title != "" {
                headers["X-Page-Title"] = pm.Title
        }
        if pm.Canonical != "" {
                headers["Link"] = fmt.Sprintf(`<%s>; rel="canonical"`, pm.Canonical)
        }
        if pm.NoIndex {
                headers["X-Robots-Tag"] = "noindex"
        }
        return headers
}

// SSGEngine handles static site generation
type SSGEngine struct {
        pages       map[string]*StaticPage
        configs     []PageConfig
        outputDir   string
        mu          sync.RWMutex
        onRevalidate func(path string)
}

// NewSSGEngine creates a new SSG engine
func NewSSGEngine(outputDir string) *SSGEngine {
        return &SSGEngine{
                pages:     make(map[string]*StaticPage),
                outputDir: outputDir,
        }
}

// RegisterPage adds a page configuration for generation
func (e *SSGEngine) RegisterPage(config PageConfig) {
        e.mu.Lock()
        defer e.mu.Unlock()
        e.configs = append(e.configs, config)
}

// Build generates all static pages
func (e *SSGEngine) Build(ctx context.Context) error {
        _ = os.MkdirAll(e.outputDir, 0755)

        for _, config := range e.configs {
                select {
                case <-ctx.Done():
                        return ctx.Err()
                default:
                }

                switch config.RenderMode {
                case RenderSSG:
                        if len(config.Params) > 0 {
                                for _, params := range config.Params {
                                        path := e.buildPath(config.Path, params)
                                        html := config.Component.Render()
                                        page := &StaticPage{
                                                Path:      path,
                                                HTML:      html,
                                                CreatedAt: time.Now(),
                                        }
                                        e.savePage(page)
                                }
                        } else {
                                html := config.Component.Render()
                                page := &StaticPage{
                                        Path:      config.Path,
                                        HTML:      html,
                                        CreatedAt: time.Now(),
                                }
                                e.savePage(page)
                        }

                case RenderISR:
                        html := config.Component.Render()
                        page := &StaticPage{
                                Path:         config.Path,
                                HTML:         html,
                                CreatedAt:    time.Now(),
                                RevalidateAt: time.Now().Add(config.Revalidate),
                        }
                        e.savePage(page)

                case RenderSSR:
                        // SSR pages are not pre-rendered, skip in build
                }
        }

        return nil
}

// ServePage serves a page with ISR support
func (e *SSGEngine) ServePage(w http.ResponseWriter, r *http.Request, path string) {
        page, exists := e.GetPage(path)
        if !exists {
                http.NotFound(w, r)
                return
        }

        // ISR: Check if revalidation is needed (stale-while-revalidate)
        if !page.RevalidateAt.IsZero() && page.RevalidateAt.Before(time.Now()) {
                go e.revalidate(page)
        }

        for k, v := range page.Metadata.toHeaders() {
                w.Header().Set(k, v)
        }
        w.Header().Set("Content-Type", "text/html")
        w.Header().Set("X-GoScript-Cache", "HIT")
        w.Write([]byte(page.HTML))
}

// GetPage retrieves a pre-rendered page by path
func (e *SSGEngine) GetPage(path string) (*StaticPage, bool) {
        e.mu.RLock()
        defer e.mu.RUnlock()
        page, ok := e.pages[path]
        return page, ok
}

// GenerateSitemap creates a sitemap.xml from all registered pages
func (e *SSGEngine) GenerateSitemap(baseURL string) string {
        e.mu.RLock()
        defer e.mu.RUnlock()

        var sb strings.Builder
        sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)

        for path, page := range e.pages {
                lastmod := page.CreatedAt.Format("2006-01-02")
                priority := "0.8"
                for _, config := range e.configs {
                        if config.Path == path && config.Priority > 0 {
                                priority = fmt.Sprintf("%.1f", config.Priority)
                        }
                }
                sb.WriteString(fmt.Sprintf(
                        `<url><loc>%s%s</loc><lastmod>%s</lastmod><priority>%s</priority></url>`,
                        baseURL, path, lastmod, priority,
                ))
        }

        sb.WriteString(`</urlset>`)
        return sb.String()
}

// GenerateRobotsTxt creates a robots.txt file
func (e *SSGEngine) GenerateRobotsTxt(baseURL string, disallowed []string) string {
        var sb strings.Builder
        sb.WriteString("User-agent: *\n")
        sb.WriteString(fmt.Sprintf("Sitemap: %s/sitemap.xml\n", baseURL))
        for _, path := range disallowed {
                sb.WriteString(fmt.Sprintf("Disallow: %s\n", path))
        }
        return sb.String()
}

// buildPath substitutes params into a path pattern
func (e *SSGEngine) buildPath(pattern string, params map[string]string) string {
        path := pattern
        for key, value := range params {
                path = strings.Replace(path, ":"+key, value, -1)
                path = strings.Replace(path, "*"+key, value, -1)
        }
        return path
}

// savePage stores a pre-rendered page
func (e *SSGEngine) savePage(page *StaticPage) {
        e.mu.Lock()
        defer e.mu.Unlock()
        e.pages[page.Path] = page

        if e.outputDir != "" {
                fullPath := filepath.Join(e.outputDir, page.Path)
                dir := filepath.Dir(fullPath)
                _ = os.MkdirAll(dir, 0755)

                filePath := fullPath
                if !strings.HasSuffix(filePath, ".html") {
                        filePath = filepath.Join(fullPath, "index.html")
                }
                _ = ioutil.WriteFile(filePath, []byte(page.HTML), 0644)
        }
}

// revalidate triggers background revalidation of a stale page
func (e *SSGEngine) revalidate(page *StaticPage) {
        for _, config := range e.configs {
                if config.Path == page.Path && config.Component != nil {
                        html := config.Component.Render()
                        e.mu.Lock()
                        page.HTML = html
                        page.RevalidateAt = time.Now().Add(config.Revalidate)
                        e.mu.Unlock()

                        if e.onRevalidate != nil {
                                e.onRevalidate(page.Path)
                        }
                }
        }
}

// OnRevalidate sets a callback for when a page is revalidated
func (e *SSGEngine) OnRevalidate(fn func(path string)) {
        e.onRevalidate = fn
}

// Stats returns statistics about the generated pages
func (e *SSGEngine) Stats() map[string]interface{} {
        e.mu.RLock()
        defer e.mu.RUnlock()

        sgCount := 0
        ssrCount := 0
        isrCount := 0
        for _, config := range e.configs {
                switch config.RenderMode {
                case RenderSSG:
                        sgCount++
                case RenderSSR:
                        ssrCount++
                case RenderISR:
                        isrCount++
                }
        }

        return map[string]interface{}{
                "total_pages":    len(e.pages),
                "ssg_pages":      sgCount,
                "ssr_pages":      ssrCount,
                "isr_pages":      isrCount,
                "output_dir":     e.outputDir,
                "registered_configs": len(e.configs),
        }
}
