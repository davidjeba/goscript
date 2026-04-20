package goscript

import (
	"bytes"
	"encoding/json"
	"html/template"
)

// HydrationPayload contains the data needed to hydrate a UI tree.
type HydrationPayload struct {
	AppID    string                 `json:"appId"`
	Version  string                 `json:"version,omitempty"`
	State    interface{}            `json:"state"`
	Endpoint string                 `json:"endpoint,omitempty"`
	Meta     map[string]string      `json:"meta,omitempty"`
}

// RenderHydrationShell wraps HTML with hydration metadata.
func RenderHydrationShell(content string, payload HydrationPayload) (string, error) {
	if payload.AppID == "" {
		payload.AppID = "app"
	}

	stateJSON, err := json.Marshal(payload.State)
	if err != nil {
		return "", err
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="goscript-app" content="{{.AppID}}">
	<meta name="goscript-version" content="{{.Version}}">
	<script>window.__GOSCRIPT_STATE__ = {{.State}}</script>
</head>
<body>
	<div id="{{.AppID}}" data-goscript-hydrate="true">{{.Content}}</div>
	<script>
		window.__GOSCRIPT_META__ = {{.Meta}};
		window.__GOSCRIPT_ENDPOINT__ = "{{.Endpoint}}";
	</script>
	<script src="/static/client.js"></script>
</body>
</html>`

	t, err := template.New("hydration").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var metaJSON []byte
	if payload.Meta != nil {
		metaJSON, err = json.Marshal(payload.Meta)
		if err != nil {
			return "", err
		}
	} else {
		metaJSON = []byte("{}")
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, struct {
		AppID    string
		Version  string
		Content  template.HTML
		State    template.JS
		Meta     template.JS
		Endpoint string
	}{
		AppID:    payload.AppID,
		Version:  payload.Version,
		Content:  template.HTML(content),
		State:    template.JS(string(stateJSON)),
		Meta:     template.JS(string(metaJSON)),
		Endpoint: payload.Endpoint,
	})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// HydrateInfo produces a compact serializable structure for clients.
func HydrateInfo(appID string, state interface{}, endpoint string) HydrationPayload {
	return HydrationPayload{
		AppID:    appID,
		State:    state,
		Endpoint: endpoint,
	}
}
