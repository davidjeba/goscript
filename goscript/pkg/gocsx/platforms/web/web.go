// Package web provides the web platform adapter for the Gocsx framework.
package web

import (
	"github.com/davidjeba/goscript/pkg/gocsx/core"
)

// WebAdapter adapts Gocsx output for web browser rendering.
type WebAdapter struct {
	config *core.Config
}

// NewWebAdapter creates a new web platform adapter.
func NewWebAdapter(config *core.Config) *WebAdapter {
	return &WebAdapter{config: config}
}

// Name returns the adapter name.
func (w *WebAdapter) Name() string {
	return "web"
}

// Render adapts CSS for web platform output.
func (w *WebAdapter) Render(css string) string {
	return css
}
