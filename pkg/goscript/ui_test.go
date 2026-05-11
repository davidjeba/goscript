package goscript

import (
	"strings"
	"testing"
)

func TestValidateForm(t *testing.T) {
	schema := FormSchema{
		"name": {
			Name:      "name",
			Type:      "string",
			Required:  true,
			MinLength: 3,
		},
		"enabled": {
			Name:    "enabled",
			Type:    "bool",
			Default: true,
		},
	}

	normalized := BindForm(schema, map[string]interface{}{
		"name": "admin",
	})

	if normalized["enabled"] != true {
		t.Fatalf("expected default enabled=true, got %v", normalized["enabled"])
	}

	if errs := ValidateForm(schema, map[string]interface{}{"name": "admin"}); len(errs) > 0 {
		t.Fatalf("unexpected validation errors: %+v", errs)
	}
}

func TestRenderHydrationShell(t *testing.T) {
	html, err := RenderHydrationShell("<div>Hello</div>", HydrationPayload{
		AppID: "admin",
		State: map[string]interface{}{"count": 1},
	})
	if err != nil {
		t.Fatalf("unexpected hydration error: %v", err)
	}

	if !strings.Contains(html, "data-goscript-hydrate=\"true\"") {
		t.Fatalf("hydration marker missing: %s", html)
	}
}

func TestDocumentIndex(t *testing.T) {
	index := NewDocumentIndex()
	index.IndexSource("main.go", "package main\nfunc Build() {}\nconst Version = 1\n")

	if symbol, ok := index.Lookup("Build"); !ok || symbol.Kind != "func" {
		t.Fatalf("expected func symbol, got %+v, ok=%v", symbol, ok)
	}

	if completions := index.Complete("B"); len(completions) == 0 {
		t.Fatalf("expected completion results")
	}
}
