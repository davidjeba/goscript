package components

import (
	"github.com/davidjeba/goscript/pkg/goscript"
)

func Home(props goscript.Props) string {
	return goscript.createElement("div", nil,
		goscript.createElement("h1", nil, "Welcome to GoScript"),
		goscript.createElement("p", nil, "This is a synthetic page for deployment"),
	)
}

