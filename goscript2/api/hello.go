package api

import (
	goscript "github.com/davidjeba/goscript/pkg/goscript"
)

// HelloHandler returns a greeting API handler
func HelloHandler(ctx *goscript.APIContext) (interface{}, error) {
	name := ctx.Query["name"]
	if name == "" {
		name = "World"
	}
	return map[string]interface{}{
		"message": "Hello, " + name + "!",
		"version": "2.0.0",
	}, nil
}
