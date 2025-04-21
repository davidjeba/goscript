package components

import (
	"fmt"
	"strings"

	"github.com/davidjeba/goscript/pkg/gocsx/core"
)

// ButtonVariant represents a button variant
type ButtonVariant string

const (
	// ButtonPrimary is the primary button variant
	ButtonPrimary ButtonVariant = "primary"
	// ButtonSecondary is the secondary button variant
	ButtonSecondary ButtonVariant = "secondary"
	// ButtonSuccess is the success button variant
	ButtonSuccess ButtonVariant = "success"
	// ButtonDanger is the danger button variant
	ButtonDanger ButtonVariant = "danger"
	// ButtonWarning is the warning button variant
	ButtonWarning ButtonVariant = "warning"
	// ButtonInfo is the info button variant
	ButtonInfo ButtonVariant = "info"
	// ButtonLight is the light button variant
	ButtonLight ButtonVariant = "light"
	// ButtonDark is the dark button variant
	ButtonDark ButtonVariant = "dark"
	// ButtonLink is the link button variant
	ButtonLink ButtonVariant = "link"
	// ButtonOutlinePrimary is the outline primary button variant
	ButtonOutlinePrimary ButtonVariant = "outline-primary"
	// ButtonOutlineSecondary is the outline secondary button variant
	ButtonOutlineSecondary ButtonVariant = "outline-secondary"
	// ButtonOutlineSuccess is the outline success button variant
	ButtonOutlineSuccess ButtonVariant = "outline-success"
	// ButtonOutlineDanger is the outline danger button variant
	ButtonOutlineDanger ButtonVariant = "outline-danger"
	// ButtonOutlineWarning is the outline warning button variant
	ButtonOutlineWarning ButtonVariant = "outline-warning"
	// ButtonOutlineInfo is the outline info button variant
	ButtonOutlineInfo ButtonVariant = "outline-info"
	// ButtonOutlineLight is the outline light button variant
	ButtonOutlineLight ButtonVariant = "outline-light"
	// ButtonOutlineDark is the outline dark button variant
	ButtonOutlineDark ButtonVariant = "outline-dark"
)

// ButtonSize represents a button size
type ButtonSize string

const (
	// ButtonSizeSmall is the small button size
	ButtonSizeSmall ButtonSize = "sm"
	// ButtonSizeMedium is the medium button size
	ButtonSizeMedium ButtonSize = "md"
	// ButtonSizeLarge is the large button size
	ButtonSizeLarge ButtonSize = "lg"
)

// ButtonProps represents button props
type ButtonProps struct {
	// ID is the button ID
	ID string

	// Variant is the button variant
	Variant ButtonVariant

	// Size is the button size
	Size ButtonSize

	// Disabled is whether the button is disabled
	Disabled bool

	// FullWidth is whether the button is full width
	FullWidth bool

	// OnClick is the click handler
	OnClick string

	// Children is the button content
	Children string

	// ClassName is additional class names
	ClassName string

	// Attributes is additional HTML attributes
	Attributes map[string]string
}

// Button creates a button component
func Button(props ButtonProps) string {
	// Build class names
	var classes []string
	classes = append(classes, "btn")

	if props.Variant != "" {
		classes = append(classes, fmt.Sprintf("btn-%s", props.Variant))
	} else {
		classes = append(classes, "btn-primary")
	}

	if props.Size != "" {
		classes = append(classes, fmt.Sprintf("btn-%s", props.Size))
	}

	if props.FullWidth {
		classes = append(classes, "btn-block")
	}

	if props.ClassName != "" {
		classes = append(classes, props.ClassName)
	}

	// Build attributes
	attributes := make(map[string]string)
	if props.ID != "" {
		attributes["id"] = props.ID
	}

	if props.Disabled {
		attributes["disabled"] = "disabled"
	}

	if props.OnClick != "" {
		attributes["onclick"] = props.OnClick
	}

	// Add custom attributes
	for key, value := range props.Attributes {
		attributes[key] = value
	}

	// Build attribute string
	var attributeStrings []string
	attributeStrings = append(attributeStrings, fmt.Sprintf(`class="%s"`, strings.Join(classes, " ")))
	for key, value := range attributes {
		attributeStrings = append(attributeStrings, fmt.Sprintf(`%s="%s"`, key, value))
	}

	// Build button
	return fmt.Sprintf(`<button %s>%s</button>`, strings.Join(attributeStrings, " "), props.Children)
}

// RegisterButtonComponent registers the button component with Gocsx
func RegisterButtonComponent(gocsx *core.Gocsx) *core.Component {
	// Define base classes
	baseClasses := []string{
		"btn",
		"inline-block",
		"font-weight-normal",
		"text-center",
		"align-middle",
		"border",
		"rounded",
		"py-2",
		"px-4",
		"text-decoration-none",
	}

	// Define variant classes
	variantClasses := map[string][]string{
		"primary": {
			"bg-primary",
			"text-white",
			"border-primary",
			"hover:bg-primary-dark",
		},
		"secondary": {
			"bg-secondary",
			"text-white",
			"border-secondary",
			"hover:bg-secondary-dark",
		},
		"success": {
			"bg-success",
			"text-white",
			"border-success",
			"hover:bg-success-dark",
		},
		"danger": {
			"bg-danger",
			"text-white",
			"border-danger",
			"hover:bg-danger-dark",
		},
		"warning": {
			"bg-warning",
			"text-dark",
			"border-warning",
			"hover:bg-warning-dark",
		},
		"info": {
			"bg-info",
			"text-white",
			"border-info",
			"hover:bg-info-dark",
		},
		"light": {
			"bg-light",
			"text-dark",
			"border-light",
			"hover:bg-light-dark",
		},
		"dark": {
			"bg-dark",
			"text-white",
			"border-dark",
			"hover:bg-dark-dark",
		},
		"outline-primary": {
			"bg-transparent",
			"text-primary",
			"border-primary",
			"hover:bg-primary",
			"hover:text-white",
		},
		"outline-secondary": {
			"bg-transparent",
			"text-secondary",
			"border-secondary",
			"hover:bg-secondary",
			"hover:text-white",
		},
		"outline-success": {
			"bg-transparent",
			"text-success",
			"border-success",
			"hover:bg-success",
			"hover:text-white",
		},
		"outline-danger": {
			"bg-transparent",
			"text-danger",
			"border-danger",
			"hover:bg-danger",
			"hover:text-white",
		},
		"outline-warning": {
			"bg-transparent",
			"text-warning",
			"border-warning",
			"hover:bg-warning",
			"hover:text-dark",
		},
		"outline-info": {
			"bg-transparent",
			"text-info",
			"border-info",
			"hover:bg-info",
			"hover:text-white",
		},
		"outline-light": {
			"bg-transparent",
			"text-light",
			"border-light",
			"hover:bg-light",
			"hover:text-dark",
		},
		"outline-dark": {
			"bg-transparent",
			"text-dark",
			"border-dark",
			"hover:bg-dark",
			"hover:text-white",
		},
		"link": {
			"bg-transparent",
			"text-primary",
			"border-transparent",
			"hover:text-primary-dark",
			"hover:underline",
		},
		"sm": {
			"py-1",
			"px-2",
			"text-sm",
		},
		"lg": {
			"py-3",
			"px-6",
			"text-lg",
		},
		"disabled": {
			"opacity-50",
			"cursor-not-allowed",
		},
		"block": {
			"w-full",
		},
	}

	// Register component
	return gocsx.RegisterComponent("button", baseClasses, variantClasses)
}