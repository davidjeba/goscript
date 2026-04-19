package goscript

import (
	"bytes"
	"encoding/json"
	"html/template"
)

type SSREngine struct {
	store *Store
}

func NewSSREngine(store *Store) *SSREngine {
	return &SSREngine{store: store}
}

func (ssr *SSREngine) RenderToString(component Component) (string, error) {
	rendered := component.Render()

	// Create a template with the rendered component and state
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<script>window.__INITIAL_STATE__ = {{.State}}</script>
</head>
<body>
	<div id="app">{{.Content}}</div>
	<script src="/static/client.js"></script>
</body>
</html>`

	t, err := template.New("ssr").Parse(tmpl)
	if err != nil {
		return "", err
	}

	// Serialize the state
	state, err := json.Marshal(ssr.store.state)
	if err != nil {
		return "", err
	}

	// Execute the template
	var buf bytes.Buffer
	err = t.Execute(&buf, struct {
		Content template.HTML
		State   template.JS
	}{
		Content: template.HTML(rendered),
		State:   template.JS(state),
	})

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

