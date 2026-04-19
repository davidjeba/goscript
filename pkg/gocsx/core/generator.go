package core

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// Generator is responsible for generating CSS from Gocsx classes
type Generator struct {
	// Configuration for the generator
	Config *Config

	// Map of generated CSS rules
	Rules map[string]string

	// Map of utility functions
	Utilities map[string]UtilityFunction

	// Map of component styles
	Components map[string]ComponentStyle

	// Map of variants
	Variants map[string]VariantFunction
}

// UtilityFunction is a function that generates CSS for a utility class
type UtilityFunction func(value string, config *Config) string

// ComponentStyle represents a component style
type ComponentStyle struct {
	// Base styles for the component
	Base string

	// Variant styles for the component
	Variants map[string]map[string]string
}

// VariantFunction is a function that applies a variant to a CSS rule
type VariantFunction func(css string, config *Config) string

// NewGenerator creates a new CSS generator
func NewGenerator(config *Config) *Generator {
	if config == nil {
		config = DefaultConfig()
	}

	return &Generator{
		Config:     config,
		Rules:      make(map[string]string),
		Utilities:  make(map[string]UtilityFunction),
		Components: make(map[string]ComponentStyle),
		Variants:   make(map[string]VariantFunction),
	}
}

// RegisterUtility registers a utility function
func (g *Generator) RegisterUtility(name string, fn UtilityFunction) {
	g.Utilities[name] = fn
}

// RegisterComponent registers a component style
func (g *Generator) RegisterComponent(name string, style ComponentStyle) {
	g.Components[name] = style
}

// RegisterVariant registers a variant function
func (g *Generator) RegisterVariant(name string, fn VariantFunction) {
	g.Variants[name] = fn
}

// GenerateCSS generates CSS for the given classes
func (g *Generator) GenerateCSS(classes []string) string {
	// Process each class
	for _, class := range classes {
		g.processClass(class)
	}

	// Sort the rules by key for consistent output
	var keys []string
	for key := range g.Rules {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Build the CSS
	var buf bytes.Buffer
	for _, key := range keys {
		buf.WriteString(fmt.Sprintf(".%s%s {\n", g.Config.Prefix, key))
		buf.WriteString(g.Rules[key])
		buf.WriteString("}\n")
	}

	return buf.String()
}

// processClass processes a single class and adds it to the rules
func (g *Generator) processClass(class string) {
	// Skip empty classes
	if class == "" {
		return
	}

	// Check if this is a component
	if component, ok := g.Components[class]; ok {
		g.Rules[class] = component.Base
		return
	}

	// Parse the class
	parts := strings.Split(class, ":")
	var baseClass string
	var variants []string

	if len(parts) > 1 {
		baseClass = parts[len(parts)-1]
		variants = parts[:len(parts)-1]
	} else {
		baseClass = parts[0]
	}

	// Parse the utility
	utilityParts := strings.Split(baseClass, "-")
	if len(utilityParts) < 2 {
		return
	}

	utilityName := utilityParts[0]
	utilityValue := strings.Join(utilityParts[1:], "-")

	// Generate the CSS for the utility
	utility, ok := g.Utilities[utilityName]
	if !ok {
		return
	}

	css := utility(utilityValue, g.Config)
	if css == "" {
		return
	}

	// Apply variants
	for i := len(variants) - 1; i >= 0; i-- {
		variant := variants[i]
		variantFn, ok := g.Variants[variant]
		if !ok {
			continue
		}

		css = variantFn(css, g.Config)
	}

	// Add the rule
	g.Rules[class] = css
}

// GenerateUtilities generates all utility classes
func (g *Generator) GenerateUtilities() string {
	var classes []string

	// Generate classes for each utility
	for utilityName, utilityFn := range g.Utilities {
		// Get the values for this utility
		values := g.getUtilityValues(utilityName)

		// Generate a class for each value
		for _, value := range values {
			classes = append(classes, fmt.Sprintf("%s-%s", utilityName, value))
		}
	}

	return g.GenerateCSS(classes)
}

// getUtilityValues gets the values for a utility
func (g *Generator) getUtilityValues(utilityName string) []string {
	switch utilityName {
	case "text":
		// Get color values
		var values []string
		for colorName := range g.Config.Theme.Colors {
			for shade := range g.Config.Theme.Colors[colorName] {
				values = append(values, fmt.Sprintf("%s-%s", colorName, shade))
			}
		}
		return values
	case "bg":
		// Get color values
		var values []string
		for colorName := range g.Config.Theme.Colors {
			for shade := range g.Config.Theme.Colors[colorName] {
				values = append(values, fmt.Sprintf("%s-%s", colorName, shade))
			}
		}
		return values
	case "p", "px", "py", "pt", "pr", "pb", "pl":
		// Get spacing values
		var values []string
		for spacing := range g.Config.Theme.Spacing {
			values = append(values, spacing)
		}
		return values
	case "m", "mx", "my", "mt", "mr", "mb", "ml":
		// Get spacing values
		var values []string
		for spacing := range g.Config.Theme.Spacing {
			values = append(values, spacing)
		}
		return values
	case "w", "h":
		// Get spacing values
		var values []string
		for spacing := range g.Config.Theme.Spacing {
			values = append(values, spacing)
		}
		// Add percentage values
		values = append(values, "1/2", "1/3", "2/3", "1/4", "3/4", "full", "screen", "auto")
		return values
	case "rounded":
		// Get border radius values
		var values []string
		for radius := range g.Config.Theme.BorderRadius {
			values = append(values, radius)
		}
		return values
	case "shadow":
		// Get shadow values
		var values []string
		for shadow := range g.Config.Theme.Shadows {
			values = append(values, shadow)
		}
		return values
	case "font":
		// Get font family values
		var values []string
		for font := range g.Config.Theme.Typography["fontFamily"] {
			values = append(values, font)
		}
		return values
	case "text-size":
		// Get font size values
		var values []string
		for size := range g.Config.Theme.Typography["fontSize"] {
			values = append(values, size)
		}
		return values
	case "font-weight":
		// Get font weight values
		var values []string
		for weight := range g.Config.Theme.Typography["fontWeight"] {
			values = append(values, weight)
		}
		return values
	case "leading":
		// Get line height values
		var values []string
		for height := range g.Config.Theme.Typography["lineHeight"] {
			values = append(values, height)
		}
		return values
	case "z":
		// Get z-index values
		var values []string
		for z := range g.Config.Theme.ZIndex {
			values = append(values, z)
		}
		return values
	case "duration":
		// Get duration values
		var values []string
		for duration := range g.Config.Theme.Durations {
			values = append(values, duration)
		}
		return values
	case "ease":
		// Get easing values
		var values []string
		for ease := range g.Config.Theme.Easings {
			values = append(values, ease)
		}
		return values
	default:
		return []string{}
	}
}

// RegisterDefaultUtilities registers the default utility functions
func (g *Generator) RegisterDefaultUtilities() {
	// Text color
	g.RegisterUtility("text", func(value string, config *Config) string {
		parts := strings.Split(value, "-")
		if len(parts) != 2 {
			return ""
		}

		colorName := parts[0]
		shade := parts[1]

		if colors, ok := config.Theme.Colors[colorName]; ok {
			if color, ok := colors[shade]; ok {
				return fmt.Sprintf("  color: %s;\n", color)
			}
		}

		return ""
	})

	// Background color
	g.RegisterUtility("bg", func(value string, config *Config) string {
		parts := strings.Split(value, "-")
		if len(parts) != 2 {
			return ""
		}

		colorName := parts[0]
		shade := parts[1]

		if colors, ok := config.Theme.Colors[colorName]; ok {
			if color, ok := colors[shade]; ok {
				return fmt.Sprintf("  background-color: %s;\n", color)
			}
		}

		return ""
	})

	// Padding
	g.RegisterUtility("p", func(value string, config *Config) string {
		if spacing, ok := config.Theme.Spacing[value]; ok {
			return fmt.Sprintf("  padding: %s;\n", spacing)
		}
		return ""
	})

	// Padding X
	g.RegisterUtility("px", func(value string, config *Config) string {
		if spacing, ok := config.Theme.Spacing[value]; ok {
			return fmt.Sprintf("  padding-left: %s;\n  padding-right: %s;\n", spacing, spacing)
		}
		return ""
	})

	// Padding Y
	g.RegisterUtility("py", func(value string, config *Config) string {
		if spacing, ok := config.Theme.Spacing[value]; ok {
			return fmt.Sprintf("  padding-top: %s;\n  padding-bottom: %s;\n", spacing, spacing)
		}
		return ""
	})

	// Margin
	g.RegisterUtility("m", func(value string, config *Config) string {
		if spacing, ok := config.Theme.Spacing[value]; ok {
			return fmt.Sprintf("  margin: %s;\n", spacing)
		}
		return ""
	})

	// Margin X
	g.RegisterUtility("mx", func(value string, config *Config) string {
		if spacing, ok := config.Theme.Spacing[value]; ok {
			return fmt.Sprintf("  margin-left: %s;\n  margin-right: %s;\n", spacing, spacing)
		}
		return ""
	})

	// Margin Y
	g.RegisterUtility("my", func(value string, config *Config) string {
		if spacing, ok := config.Theme.Spacing[value]; ok {
			return fmt.Sprintf("  margin-top: %s;\n  margin-bottom: %s;\n", spacing, spacing)
		}
		return ""
	})

	// Width
	g.RegisterUtility("w", func(value string, config *Config) string {
		if value == "full" {
			return "  width: 100%;\n"
		} else if value == "screen" {
			return "  width: 100vw;\n"
		} else if value == "auto" {
			return "  width: auto;\n"
		} else if strings.Contains(value, "/") {
			parts := strings.Split(value, "/")
			if len(parts) == 2 {
				numerator := parts[0]
				denominator := parts[1]
				return fmt.Sprintf("  width: calc(%s / %s * 100%%);\n", numerator, denominator)
			}
		} else if spacing, ok := config.Theme.Spacing[value]; ok {
			return fmt.Sprintf("  width: %s;\n", spacing)
		}
		return ""
	})

	// Height
	g.RegisterUtility("h", func(value string, config *Config) string {
		if value == "full" {
			return "  height: 100%;\n"
		} else if value == "screen" {
			return "  height: 100vh;\n"
		} else if value == "auto" {
			return "  height: auto;\n"
		} else if strings.Contains(value, "/") {
			parts := strings.Split(value, "/")
			if len(parts) == 2 {
				numerator := parts[0]
				denominator := parts[1]
				return fmt.Sprintf("  height: calc(%s / %s * 100%%);\n", numerator, denominator)
			}
		} else if spacing, ok := config.Theme.Spacing[value]; ok {
			return fmt.Sprintf("  height: %s;\n", spacing)
		}
		return ""
	})

	// Border radius
	g.RegisterUtility("rounded", func(value string, config *Config) string {
		if radius, ok := config.Theme.BorderRadius[value]; ok {
			return fmt.Sprintf("  border-radius: %s;\n", radius)
		}
		return ""
	})

	// Shadow
	g.RegisterUtility("shadow", func(value string, config *Config) string {
		if shadow, ok := config.Theme.Shadows[value]; ok {
			return fmt.Sprintf("  box-shadow: %s;\n", shadow)
		}
		return ""
	})

	// Font family
	g.RegisterUtility("font", func(value string, config *Config) string {
		if fontFamily, ok := config.Theme.Typography["fontFamily"][value]; ok {
			return fmt.Sprintf("  font-family: %s;\n", fontFamily)
		}
		return ""
	})

	// Font size
	g.RegisterUtility("text-size", func(value string, config *Config) string {
		if fontSize, ok := config.Theme.Typography["fontSize"][value]; ok {
			return fmt.Sprintf("  font-size: %s;\n", fontSize)
		}
		return ""
	})

	// Font weight
	g.RegisterUtility("font-weight", func(value string, config *Config) string {
		if fontWeight, ok := config.Theme.Typography["fontWeight"][value]; ok {
			return fmt.Sprintf("  font-weight: %s;\n", fontWeight)
		}
		return ""
	})

	// Line height
	g.RegisterUtility("leading", func(value string, config *Config) string {
		if lineHeight, ok := config.Theme.Typography["lineHeight"][value]; ok {
			return fmt.Sprintf("  line-height: %s;\n", lineHeight)
		}
		return ""
	})

	// Z-index
	g.RegisterUtility("z", func(value string, config *Config) string {
		if zIndex, ok := config.Theme.ZIndex[value]; ok {
			if value == "auto" {
				return "  z-index: auto;\n"
			}
			return fmt.Sprintf("  z-index: %d;\n", zIndex)
		}
		return ""
	})

	// Display
	g.RegisterUtility("display", func(value string, config *Config) string {
		validValues := map[string]bool{
			"block":        true,
			"inline":       true,
			"inline-block": true,
			"flex":         true,
			"inline-flex":  true,
			"grid":         true,
			"inline-grid":  true,
			"table":        true,
			"hidden":       true,
		}

		if validValues[value] {
			if value == "hidden" {
				return "  display: none;\n"
			}
			return fmt.Sprintf("  display: %s;\n", value)
		}
		return ""
	})

	// Flex direction
	g.RegisterUtility("flex", func(value string, config *Config) string {
		switch value {
		case "row":
			return "  flex-direction: row;\n"
		case "row-reverse":
			return "  flex-direction: row-reverse;\n"
		case "col":
			return "  flex-direction: column;\n"
		case "col-reverse":
			return "  flex-direction: column-reverse;\n"
		case "1":
			return "  flex: 1 1 0%;\n"
		case "auto":
			return "  flex: 1 1 auto;\n"
		case "initial":
			return "  flex: 0 1 auto;\n"
		case "none":
			return "  flex: none;\n"
		}
		return ""
	})

	// Justify content
	g.RegisterUtility("justify", func(value string, config *Config) string {
		switch value {
		case "start":
			return "  justify-content: flex-start;\n"
		case "end":
			return "  justify-content: flex-end;\n"
		case "center":
			return "  justify-content: center;\n"
		case "between":
			return "  justify-content: space-between;\n"
		case "around":
			return "  justify-content: space-around;\n"
		case "evenly":
			return "  justify-content: space-evenly;\n"
		}
		return ""
	})

	// Align items
	g.RegisterUtility("items", func(value string, config *Config) string {
		switch value {
		case "start":
			return "  align-items: flex-start;\n"
		case "end":
			return "  align-items: flex-end;\n"
		case "center":
			return "  align-items: center;\n"
		case "baseline":
			return "  align-items: baseline;\n"
		case "stretch":
			return "  align-items: stretch;\n"
		}
		return ""
	})

	// Text align
	g.RegisterUtility("text-align", func(value string, config *Config) string {
		switch value {
		case "left":
			return "  text-align: left;\n"
		case "center":
			return "  text-align: center;\n"
		case "right":
			return "  text-align: right;\n"
		case "justify":
			return "  text-align: justify;\n"
		}
		return ""
	})

	// Position
	g.RegisterUtility("position", func(value string, config *Config) string {
		switch value {
		case "static":
			return "  position: static;\n"
		case "relative":
			return "  position: relative;\n"
		case "absolute":
			return "  position: absolute;\n"
		case "fixed":
			return "  position: fixed;\n"
		case "sticky":
			return "  position: sticky;\n"
		}
		return ""
	})

	// Overflow
	g.RegisterUtility("overflow", func(value string, config *Config) string {
		switch value {
		case "auto":
			return "  overflow: auto;\n"
		case "hidden":
			return "  overflow: hidden;\n"
		case "visible":
			return "  overflow: visible;\n"
		case "scroll":
			return "  overflow: scroll;\n"
		case "x-auto":
			return "  overflow-x: auto;\n"
		case "y-auto":
			return "  overflow-y: auto;\n"
		case "x-hidden":
			return "  overflow-x: hidden;\n"
		case "y-hidden":
			return "  overflow-y: hidden;\n"
		}
		return ""
	})

	// Animation duration
	g.RegisterUtility("duration", func(value string, config *Config) string {
		if duration, ok := config.Theme.Durations[value]; ok {
			return fmt.Sprintf("  transition-duration: %s;\n", duration)
		}
		return ""
	})

	// Animation easing
	g.RegisterUtility("ease", func(value string, config *Config) string {
		if easing, ok := config.Theme.Easings[value]; ok {
			return fmt.Sprintf("  transition-timing-function: %s;\n", easing)
		}
		return ""
	})
}

// RegisterDefaultVariants registers the default variant functions
func (g *Generator) RegisterDefaultVariants() {
	// Hover variant
	g.RegisterVariant("hover", func(css string, config *Config) string {
		return fmt.Sprintf("  &:hover {\n%s  }\n", indentCSS(css))
	})

	// Focus variant
	g.RegisterVariant("focus", func(css string, config *Config) string {
		return fmt.Sprintf("  &:focus {\n%s  }\n", indentCSS(css))
	})

	// Active variant
	g.RegisterVariant("active", func(css string, config *Config) string {
		return fmt.Sprintf("  &:active {\n%s  }\n", indentCSS(css))
	})

	// Disabled variant
	g.RegisterVariant("disabled", func(css string, config *Config) string {
		return fmt.Sprintf("  &:disabled {\n%s  }\n", indentCSS(css))
	})

	// Dark mode variant
	g.RegisterVariant("dark", func(css string, config *Config) string {
		return fmt.Sprintf("  @media (prefers-color-scheme: dark) {\n%s  }\n", indentCSS(css))
	})

	// Responsive variants
	for breakpoint, width := range g.Config.Breakpoints {
		breakpointName := breakpoint
		breakpointWidth := width
		
		g.RegisterVariant(breakpointName, func(css string, config *Config) string {
			return fmt.Sprintf("  @media (min-width: %dpx) {\n%s  }\n", breakpointWidth, indentCSS(css))
		})
	}
}

// indentCSS indents CSS by two spaces
func indentCSS(css string) string {
	lines := strings.Split(css, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = "    " + line
		}
	}
	return strings.Join(lines, "\n")
}