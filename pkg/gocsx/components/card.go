package components

import (
	"fmt"
	"strings"

	"github.com/davidjeba/goscript/pkg/gocsx/core"
)

// CardProps represents card props
type CardProps struct {
	// ID is the card ID
	ID string

	// Title is the card title
	Title string

	// Subtitle is the card subtitle
	Subtitle string

	// Body is the card body content
	Body string

	// Footer is the card footer content
	Footer string

	// Header is the card header content
	Header string

	// Image is the card image URL
	Image string

	// ImageAlt is the card image alt text
	ImageAlt string

	// ImagePosition is the card image position (top, bottom)
	ImagePosition string

	// ClassName is additional class names
	ClassName string

	// Attributes is additional HTML attributes
	Attributes map[string]string
}

// Card creates a card component
func Card(props CardProps) string {
	// Build class names
	var classes []string
	classes = append(classes, "card")

	if props.ClassName != "" {
		classes = append(classes, props.ClassName)
	}

	// Build attributes
	attributes := make(map[string]string)
	if props.ID != "" {
		attributes["id"] = props.ID
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

	// Build card content
	var content strings.Builder

	// Add image at top
	if props.Image != "" && (props.ImagePosition == "" || props.ImagePosition == "top") {
		alt := props.ImageAlt
		if alt == "" {
			alt = "Card image"
		}
		content.WriteString(fmt.Sprintf(`<img src="%s" alt="%s" class="card-img-top">`, props.Image, alt))
	}

	// Add header
	if props.Header != "" {
		content.WriteString(fmt.Sprintf(`<div class="card-header">%s</div>`, props.Header))
	}

	// Add body
	content.WriteString(`<div class="card-body">`)
	if props.Title != "" {
		content.WriteString(fmt.Sprintf(`<h5 class="card-title">%s</h5>`, props.Title))
	}
	if props.Subtitle != "" {
		content.WriteString(fmt.Sprintf(`<h6 class="card-subtitle mb-2 text-muted">%s</h6>`, props.Subtitle))
	}
	if props.Body != "" {
		content.WriteString(fmt.Sprintf(`<div class="card-text">%s</div>`, props.Body))
	}
	content.WriteString(`</div>`)

	// Add footer
	if props.Footer != "" {
		content.WriteString(fmt.Sprintf(`<div class="card-footer">%s</div>`, props.Footer))
	}

	// Add image at bottom
	if props.Image != "" && props.ImagePosition == "bottom" {
		alt := props.ImageAlt
		if alt == "" {
			alt = "Card image"
		}
		content.WriteString(fmt.Sprintf(`<img src="%s" alt="%s" class="card-img-bottom">`, props.Image, alt))
	}

	// Build card
	return fmt.Sprintf(`<div %s>%s</div>`, strings.Join(attributeStrings, " "), content.String())
}

// RegisterCardComponent registers the card component with Gocsx
func RegisterCardComponent(gocsx *core.Gocsx) *core.Component {
	// Define base classes
	baseClasses := []string{
		"card",
		"position-relative",
		"d-flex",
		"flex-column",
		"bg-white",
		"border",
		"rounded",
		"overflow-hidden",
	}

	// Define variant classes
	variantClasses := map[string][]string{
		"shadow": {
			"shadow",
		},
		"shadow-sm": {
			"shadow-sm",
		},
		"shadow-lg": {
			"shadow-lg",
		},
		"border-primary": {
			"border-primary",
		},
		"border-secondary": {
			"border-secondary",
		},
		"border-success": {
			"border-success",
		},
		"border-danger": {
			"border-danger",
		},
		"border-warning": {
			"border-warning",
		},
		"border-info": {
			"border-info",
		},
		"border-light": {
			"border-light",
		},
		"border-dark": {
			"border-dark",
		},
		"bg-primary": {
			"bg-primary",
			"text-white",
		},
		"bg-secondary": {
			"bg-secondary",
			"text-white",
		},
		"bg-success": {
			"bg-success",
			"text-white",
		},
		"bg-danger": {
			"bg-danger",
			"text-white",
		},
		"bg-warning": {
			"bg-warning",
		},
		"bg-info": {
			"bg-info",
			"text-white",
		},
		"bg-light": {
			"bg-light",
		},
		"bg-dark": {
			"bg-dark",
			"text-white",
		},
		"text-center": {
			"text-center",
		},
		"text-right": {
			"text-right",
		},
	}

	// Register component
	return gocsx.RegisterComponent("card", baseClasses, variantClasses)
}

// CardHeaderProps represents card header props
type CardHeaderProps struct {
	// ID is the card header ID
	ID string

	// Children is the card header content
	Children string

	// ClassName is additional class names
	ClassName string

	// Attributes is additional HTML attributes
	Attributes map[string]string
}

// CardHeader creates a card header component
func CardHeader(props CardHeaderProps) string {
	// Build class names
	var classes []string
	classes = append(classes, "card-header")

	if props.ClassName != "" {
		classes = append(classes, props.ClassName)
	}

	// Build attributes
	attributes := make(map[string]string)
	if props.ID != "" {
		attributes["id"] = props.ID
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

	// Build card header
	return fmt.Sprintf(`<div %s>%s</div>`, strings.Join(attributeStrings, " "), props.Children)
}

// CardBodyProps represents card body props
type CardBodyProps struct {
	// ID is the card body ID
	ID string

	// Children is the card body content
	Children string

	// ClassName is additional class names
	ClassName string

	// Attributes is additional HTML attributes
	Attributes map[string]string
}

// CardBody creates a card body component
func CardBody(props CardBodyProps) string {
	// Build class names
	var classes []string
	classes = append(classes, "card-body")

	if props.ClassName != "" {
		classes = append(classes, props.ClassName)
	}

	// Build attributes
	attributes := make(map[string]string)
	if props.ID != "" {
		attributes["id"] = props.ID
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

	// Build card body
	return fmt.Sprintf(`<div %s>%s</div>`, strings.Join(attributeStrings, " "), props.Children)
}

// CardFooterProps represents card footer props
type CardFooterProps struct {
	// ID is the card footer ID
	ID string

	// Children is the card footer content
	Children string

	// ClassName is additional class names
	ClassName string

	// Attributes is additional HTML attributes
	Attributes map[string]string
}

// CardFooter creates a card footer component
func CardFooter(props CardFooterProps) string {
	// Build class names
	var classes []string
	classes = append(classes, "card-footer")

	if props.ClassName != "" {
		classes = append(classes, props.ClassName)
	}

	// Build attributes
	attributes := make(map[string]string)
	if props.ID != "" {
		attributes["id"] = props.ID
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

	// Build card footer
	return fmt.Sprintf(`<div %s>%s</div>`, strings.Join(attributeStrings, " "), props.Children)
}

// CardTitleProps represents card title props
type CardTitleProps struct {
	// ID is the card title ID
	ID string

	// Children is the card title content
	Children string

	// ClassName is additional class names
	ClassName string

	// Attributes is additional HTML attributes
	Attributes map[string]string
}

// CardTitle creates a card title component
func CardTitle(props CardTitleProps) string {
	// Build class names
	var classes []string
	classes = append(classes, "card-title")

	if props.ClassName != "" {
		classes = append(classes, props.ClassName)
	}

	// Build attributes
	attributes := make(map[string]string)
	if props.ID != "" {
		attributes["id"] = props.ID
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

	// Build card title
	return fmt.Sprintf(`<h5 %s>%s</h5>`, strings.Join(attributeStrings, " "), props.Children)
}

// CardSubtitleProps represents card subtitle props
type CardSubtitleProps struct {
	// ID is the card subtitle ID
	ID string

	// Children is the card subtitle content
	Children string

	// ClassName is additional class names
	ClassName string

	// Attributes is additional HTML attributes
	Attributes map[string]string
}

// CardSubtitle creates a card subtitle component
func CardSubtitle(props CardSubtitleProps) string {
	// Build class names
	var classes []string
	classes = append(classes, "card-subtitle", "mb-2", "text-muted")

	if props.ClassName != "" {
		classes = append(classes, props.ClassName)
	}

	// Build attributes
	attributes := make(map[string]string)
	if props.ID != "" {
		attributes["id"] = props.ID
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

	// Build card subtitle
	return fmt.Sprintf(`<h6 %s>%s</h6>`, strings.Join(attributeStrings, " "), props.Children)
}

// CardTextProps represents card text props
type CardTextProps struct {
	// ID is the card text ID
	ID string

	// Children is the card text content
	Children string

	// ClassName is additional class names
	ClassName string

	// Attributes is additional HTML attributes
	Attributes map[string]string
}

// CardText creates a card text component
func CardText(props CardTextProps) string {
	// Build class names
	var classes []string
	classes = append(classes, "card-text")

	if props.ClassName != "" {
		classes = append(classes, props.ClassName)
	}

	// Build attributes
	attributes := make(map[string]string)
	if props.ID != "" {
		attributes["id"] = props.ID
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

	// Build card text
	return fmt.Sprintf(`<p %s>%s</p>`, strings.Join(attributeStrings, " "), props.Children)
}

// CardImageProps represents card image props
type CardImageProps struct {
	// ID is the card image ID
	ID string

	// Src is the card image source URL
	Src string

	// Alt is the card image alt text
	Alt string

	// Position is the card image position (top, bottom)
	Position string

	// ClassName is additional class names
	ClassName string

	// Attributes is additional HTML attributes
	Attributes map[string]string
}

// CardImage creates a card image component
func CardImage(props CardImageProps) string {
	// Build class names
	var classes []string
	
	if props.Position == "bottom" {
		classes = append(classes, "card-img-bottom")
	} else {
		classes = append(classes, "card-img-top")
	}

	if props.ClassName != "" {
		classes = append(classes, props.ClassName)
	}

	// Build attributes
	attributes := make(map[string]string)
	if props.ID != "" {
		attributes["id"] = props.ID
	}
	
	if props.Src != "" {
		attributes["src"] = props.Src
	}
	
	if props.Alt != "" {
		attributes["alt"] = props.Alt
	} else {
		attributes["alt"] = "Card image"
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

	// Build card image
	return fmt.Sprintf(`<img %s>`, strings.Join(attributeStrings, " "))
}