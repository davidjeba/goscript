package goscript

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// SuspenseBoundary represents a lazy-loaded section of a page
type SuspenseBoundary struct {
	ID       string
	Fallback Component
	Loader   func(ctx context.Context) (Component, error)
}

// StreamChunk represents a chunk of HTML to be sent to the client
type StreamChunk struct {
	ID      string
	Type    string // "html", "suspense-start", "suspense-resolve", "error"
	Content string
}

// StreamSSREngine provides streaming server-side rendering
type StreamSSREngine struct {
	store        *Store
	chunkChannel chan StreamChunk
	flusher      func(io.Writer)
}

// NewStreamSSREngine creates a new streaming SSR engine
func NewStreamSSREngine(store *Store) *StreamSSREngine {
	return &StreamSSREngine{
		store:        store,
		chunkChannel: make(chan StreamChunk, 100),
	}
}

// RenderStream renders a component with streaming support.
// It sends chunks of HTML as they become ready, enabling progressive loading.
func (s *StreamSSREngine) RenderStream(
	w http.ResponseWriter,
	r *http.Request,
	component Component,
	suspenseBoundaries []SuspenseBoundary,
) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		s.renderFallback(w, component)
		return
	}

	// Set streaming headers
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	ctx := r.Context()

	// Phase 1: Send the shell with initial HTML
	_, _ = io.WriteString(w, "<!DOCTYPE html><html><head>")
	_, _ = io.WriteString(w, s.generateHead())
	_, _ = io.WriteString(w, "</head><body>")
	_, _ = io.WriteString(w, s.generateInitialStoreScript())
	_, _ = io.WriteString(w, `<div id="__goscript_app">`)
	flusher.Flush()

	// Phase 2: Stream the main component
	mainHTML := component.Render()
	_, _ = io.WriteString(w, mainHTML)
	flusher.Flush()

	// Phase 3: Stream Suspense boundaries in parallel
	var wg sync.WaitGroup
	for _, boundary := range suspenseBoundaries {
		wg.Add(1)
		go func(b SuspenseBoundary) {
			defer wg.Done()

			// Send the fallback immediately
			fallbackHTML := ""
			if b.Fallback != nil {
				fallbackHTML = b.Fallback.Render()
			}
			chunk := StreamChunk{
				ID:      b.ID,
				Type:    "suspense-start",
				Content: fmt.Sprintf(`<div id="%s" data-suspense>%s</div>`, b.ID, fallbackHTML),
			}
			s.chunkChannel <- chunk

			// Load the actual component
			loaded, err := b.Loader(ctx)
			if err != nil {
				errChunk := StreamChunk{
					ID:      b.ID,
					Type:    "error",
					Content: fmt.Sprintf(`<div id="%s" data-error>%s</div>`, b.ID, err.Error()),
				}
				s.chunkChannel <- errChunk
				return
			}

			html := loaded.Render()
			resolveChunk := StreamChunk{
				ID:      b.ID,
				Type:    "suspense-resolve",
				Content: fmt.Sprintf(`<div id="%s">%s</div>`, b.ID, html),
			}
			s.chunkChannel <- resolveChunk
		}(boundary)
	}

	// Phase 4: Stream chunks as they arrive
	go func() {
		wg.Wait()
		close(s.chunkChannel)
	}()

	for chunk := range s.chunkChannel {
		_, _ = io.WriteString(w, chunk.Content)
		flusher.Flush()
	}

	// Phase 5: Close the HTML document
	_, _ = io.WriteString(w, `</div>`)
	_, _ = io.WriteString(w, s.generateBootScript())
	_, _ = io.WriteString(w, "</body></html>")
	flusher.Flush()
}

// renderFallback renders the entire page without streaming
func (s *StreamSSREngine) renderFallback(w http.ResponseWriter, component Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var sb strings.Builder
	sb.WriteString("<!DOCTYPE html><html><head>")
	sb.WriteString(s.generateHead())
	sb.WriteString("</head><body>")
	sb.WriteString(`<div id="__goscript_app">`)
	if component != nil {
		sb.WriteString(component.Render())
	}
	sb.WriteString(`</div>`)
	sb.WriteString(s.generateBootScript())
	sb.WriteString("</body></html>")

	w.Write([]byte(sb.String()))
}

// generateHead produces the HTML <head> content
func (s *StreamSSREngine) generateHead() string {
	return `<meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>GoScript App</title>`
}

// generateInitialStoreScript serializes the store state for hydration
func (s *StreamSSREngine) generateInitialStoreScript() string {
	if s.store == nil {
		return ""
	}
	stateJSON := "{}"
	if s.store != nil {
		b, err := json.Marshal(s.store.state)
		if err == nil {
			stateJSON = string(b)
		}
	}
	return fmt.Sprintf(`<script>window.__GS_INITIAL_STATE__=%s;</script>`, stateJSON)
}

// generateBootScript produces the client-side boot script
func (s *StreamSSREngine) generateBootScript() string {
	return `<script>
(function(){
  var app = document.getElementById("__goscript_app");
  if (window.__gs_init) window.__gs_init(app);
})();
</script>`
}

// StreamHTML writes an HTML chunk to the response writer
func (s *StreamSSREngine) StreamHTML(w io.Writer, html string) {
	if s.flusher != nil {
		s.flusher(w)
	}
	io.WriteString(w, html)
}

// StreamReader streams lines from an io.Reader as HTML chunks
func (s *StreamSSREngine) StreamReader(ctx context.Context, r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			line := scanner.Text()
			chunk := StreamChunk{
				ID:      fmt.Sprintf("stream-%d", time.Now().UnixNano()),
				Type:    "html",
				Content: line + "\n",
			}
			s.chunkChannel <- chunk
			if flusher, ok := w.(http.Flusher); ok {
				io.WriteString(w, chunk.Content)
				flusher.Flush()
			}
		}
	}
	return scanner.Err()
}
