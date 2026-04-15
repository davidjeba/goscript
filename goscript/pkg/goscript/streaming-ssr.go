package goscript

import (
        "fmt"
        "net/http"
        "strings"
        "sync"
)

// SuspenseBoundary defines a region of a page that can be streamed independently.
// While the async content is loading, the Fallback component is rendered as a
// placeholder. Once the Loader function completes, the boundary is replaced with
// the resolved Component.
type SuspenseBoundary struct {
        ID       string
        Fallback Component
        Loader   func(ctx interface{}) (Component, error)
}

// StreamSSREngine enhances the existing v1 SSREngine with chunked HTML streaming
// via http.Flusher. It supports Suspense boundaries that allow portions of a page
// to be streamed as they become available, while showing skeleton/fallback content
// for sections that are still loading.
type StreamSSREngine struct {
        store *Store
}

// NewStreamSSREngine creates a new streaming SSR engine backed by the given Store.
// The store is used to inject initial state into the streamed HTML.
func NewStreamSSREngine(store *Store) *StreamSSREngine {
        return &StreamSSREngine{store: store}
}

// RenderStream performs chunked HTML streaming to the http.ResponseWriter. It renders
// the provided Component immediately, then streams each SuspenseBoundary's content
// as it becomes available. Each boundary is wrapped in a <div data-suspense-id="...">
// marker so the client can patch the DOM when the real content arrives.
func (sse *StreamSSREngine) RenderStream(w http.ResponseWriter, r *http.Request, component Component, boundaries []SuspenseBoundary) {
        // Set streaming-friendly headers
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        w.Header().Set("Transfer-Encoding", "chunked")
        w.Header().Set("X-Content-Type-Options", "nosniff")

        flusher, canFlush := w.(http.Flusher)
        if !canFlush {
                // Fall back to non-streaming render if flusher is not available
                html := sse.renderFullHTML(component, boundaries)
                fmt.Fprint(w, html)
                return
        }

        // Stream the document head
        head := sse.renderHead(r)
        fmt.Fprint(w, head)
        flusher.Flush()

        // Stream the opening body and initial component content
        fmt.Fprint(w, `<body><div id="app">`)
        flusher.Flush()

        // Render the main component
        componentHTML := component.Render()
        fmt.Fprint(w, componentHTML)
        flusher.Flush()

        // Stream each SuspenseBoundary: first the fallback, then the resolved content
        var wg sync.WaitGroup
        results := make([]string, len(boundaries))
        for i, boundary := range boundaries {
                wg.Add(1)
                idx := i
                b := boundary

                go func() {
                        defer wg.Done()
                        loaderResult, err := b.Loader(r.Context())
                        if err != nil {
                                results[idx] = renderErrorBoundary(b.ID, err.Error())
                                return
                        }
                        results[idx] = renderSuspenseResolved(b.ID, loaderResult.Render())
                }()
        }

        // First pass: render all fallbacks
        for _, boundary := range boundaries {
                fallbackHTML := renderSuspenseFallback(boundary.ID, boundary.Fallback)
                fmt.Fprint(w, fallbackHTML)
                flusher.Flush()
        }

        // Wait for all loaders to complete
        wg.Wait()

        // Second pass: stream the resolved content for each boundary
        for i := range boundaries {
                fmt.Fprint(w, results[i])
                flusher.Flush()
        }

        // Stream the closing tags and initial state script
        fmt.Fprint(w, `</div>`)
        fmt.Fprint(w, sse.renderStateScript())
        fmt.Fprint(w, `</body></html>`)
        flusher.Flush()
}

// renderHead generates the HTML document head with metadata placeholders.
func (sse *StreamSSREngine) renderHead(r *http.Request) string {
        return `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>GoScript App</title>
<style>
.goscript-skeleton {
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: goscript-shimmer 1.5s infinite;
  border-radius: 4px;
}
@keyframes goscript-shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
</style>
</head>`
}

// renderStateScript serializes the store's state into a <script> tag that
// provides __INITIAL_STATE__ to the client for hydration.
func (sse *StreamSSREngine) renderStateScript() string {
        if sse.store == nil {
                return ""
        }

        var sb strings.Builder
        sb.WriteString(`<script>window.__INITIAL_STATE__ = {`)
        first := true
        for key, value := range sse.store.state {
                if !first {
                        sb.WriteString(",")
                }
                first = false
                sb.WriteString(fmt.Sprintf(`"%s": %q`, key, fmt.Sprintf("%v", value)))
        }
        sb.WriteString(`};</script>`)
        return sb.String()
}

// renderSuspenseFallback renders a SuspenseBoundary's fallback content wrapped in
// a data-suspense-id marker div.
func renderSuspenseFallback(id string, fallback Component) string {
        fallbackHTML := ""
        if fallback != nil {
                fallbackHTML = fallback.Render()
        }
        return fmt.Sprintf(`<div data-suspense-id="%s" data-suspense-status="pending">%s</div>`, id, fallbackHTML)
}

// renderSuspenseResolved renders a script block that replaces a pending SuspenseBoundary
// with its resolved content on the client side.
func renderSuspenseResolved(id string, content string) string {
        // Escape content for safe embedding in a JavaScript string
        escaped := strings.Replace(content, `\`, `\\`, -1)
        escaped = strings.Replace(escaped, "`", "\\x60", -1)
        escaped = strings.Replace(escaped, "</", "<\\/", -1)

        script := "(function() {\n"
        script += fmt.Sprintf("  var el = document.querySelector('[data-suspense-id=\"%s\"]');\n", id)
        script += "  if (el && el.getAttribute('data-suspense-status') === 'pending') {\n"
        script += fmt.Sprintf("    el.innerHTML = `%s`;\n", escaped)
        script += "    el.setAttribute('data-suspense-status', 'resolved');\n"
        script += "  }\n"
        script += "})();\n"

        return "<script>\n" + script + "</script>"
}

// renderErrorBoundary renders an error state inside a SuspenseBoundary marker.
func renderErrorBoundary(id string, errMsg string) string {
        escaped := strings.Replace(errMsg, "`", "\\x60", -1)
        return fmt.Sprintf(`<script>
(function() {
  var el = document.querySelector('[data-suspense-id="%s"]');
  if (el && el.getAttribute('data-suspense-status') === 'pending') {
    el.innerHTML = '<div class="error-boundary" style="color:red;padding:1rem;border:1px solid red;border-radius:4px;">Error: %s</div>';
    el.setAttribute('data-suspense-status', 'error');
  }
})();
</script>`, id, escaped)
}

// renderFullHTML falls back to rendering a complete non-streamed HTML document.
// This is used when the ResponseWriter does not support flushing.
func (sse *StreamSSREngine) renderFullHTML(component Component, boundaries []SuspenseBoundary) string {
        var sb strings.Builder

        sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>GoScript App</title>
</head>
<body><div id="app">`)

        sb.WriteString(component.Render())

        for _, boundary := range boundaries {
                fallbackHTML := ""
                if boundary.Fallback != nil {
                        fallbackHTML = boundary.Fallback.Render()
                }
                sb.WriteString(fmt.Sprintf(`<div data-suspense-id="%s">%s</div>`, boundary.ID, fallbackHTML))
        }

        sb.WriteString(`</div>`)
        sb.WriteString(sse.renderStateScript())
        sb.WriteString(`</body></html>`)

        return sb.String()
}
