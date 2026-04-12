package app

import (
	"fmt"

	goscript "github.com/davidjeba/goscript/pkg/goscript"
)

// Page returns the home page component
func Page() goscript.Component {
	return goscript.FunctionalComponent(func(props goscript.Props) string {
		title := "Welcome"
		if t, ok := props["title"].(string); ok {
			title = t
		}

		return fmt.Sprintf(`
<main>
  <h1>%s to GoScript</h1>
  <p>This is a page rendered by GoScript 2.0.</p>
</main>
`, title)
	})
}
