package app

import (
	goscript "github.com/davidjeba/goscript/pkg/goscript"
)

// Layout returns the root layout component for the application
func Layout() goscript.Component {
	return goscript.FunctionalComponent(func(props goscript.Props) string {
		metadata := goscript.NewMetadata().
			SetTitle("GoScript App").
			SetDescription("Built with GoScript 2.0").
			SetThemeColor("#10b981").
			Build()

		return `<!DOCTYPE html>
<html lang="en">
<head>` + metadata.Render() + `</head>
<body>
<div id="__goscript_app">
{{children}}
</div>
</body>
</html>`
	})
}
