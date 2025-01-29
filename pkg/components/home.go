package components

import (
	"github.com/davidjeba/goscript/pkg/goscript"
)

func Home(props goscript.Props) string {
	return goscript.CreateElement("div", nil,
		goscript.CreateElement("h1", nil, "Welcome to GoScript"),
		goscript.CreateElement("p", nil, "This is the home component."),
	)
}

