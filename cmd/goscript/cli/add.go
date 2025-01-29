package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

func AddComponent(name string) error {
	// Create components directory if it doesn't exist
	componentsDir := filepath.Join("pkg", "components")
	if err := os.MkdirAll(componentsDir, 0755); err != nil {
		return fmt.Errorf("failed to create components directory: %v", err)
	}

	// Create component file
	filename := filepath.Join(componentsDir, fmt.Sprintf("%s.gsx", name))
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create component file: %v", err)
	}
	defer f.Close()

	// Write component template
	template := fmt.Sprintf(`package components

import (
	"github.com/davidjeba/goscript/pkg/goscript"
)

func %s(props goscript.Props) string {
	return goscript.createElement("div", nil,
		goscript.createElement("h1", nil, "%s Component"),
	)
}
`, name, name)

	if _, err := f.WriteString(template); err != nil {
		return fmt.Errorf("failed to write component template: %v", err)
	}

	fmt.Printf("Created new component: %s\n", filename)
	return nil
}

