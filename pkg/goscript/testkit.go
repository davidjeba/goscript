package goscript

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

// Snapshot captures a renderable output for tests.
type Snapshot struct {
	Name    string
	Content string
}

// SnapshotComponent renders a component for snapshot testing.
func SnapshotComponent(name string, component Component) Snapshot {
	return Snapshot{
		Name:    name,
		Content: component.Render(),
	}
}

// RenderRoute renders an HTTP route and returns the response body.
func RenderRoute(router http.Handler, method, path string) (int, string) {
	req := httptest.NewRequest(method, path, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr.Code, rr.Body.String()
}

// NormalizeMarkup collapses whitespace for stable comparisons.
func NormalizeMarkup(input string) string {
	fields := strings.Fields(input)
	return strings.Join(fields, " ")
}

