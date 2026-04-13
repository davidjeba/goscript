// Package components provides pre-built UI components for the Gocsx framework.
package components

// ButtonVariant controls the visual style of a button.
type ButtonVariant string

const (
	ButtonPrimary   ButtonVariant = "primary"
	ButtonSecondary ButtonVariant = "secondary"
	ButtonSuccess   ButtonVariant = "success"
	ButtonDanger    ButtonVariant = "danger"
	ButtonWarning   ButtonVariant = "warning"
	ButtonInfo      ButtonVariant = "info"
	ButtonLight     ButtonVariant = "light"
	ButtonDark      ButtonVariant = "dark"
)

// ButtonSize controls the size of a button.
type ButtonSize string

const (
	ButtonSizeSmall  ButtonSize = "sm"
	ButtonSizeMedium ButtonSize = "md"
	ButtonSizeLarge  ButtonSize = "lg"
)

// ButtonProps defines the properties for creating a button component.
type ButtonProps struct {
	ID        string
	Variant   ButtonVariant
	Size      ButtonSize
	Children  string
	OnClick   string
	Disabled  bool
	ClassName string
}

// Button creates an HTML button element with Gocsx styling.
func Button(props ButtonProps) string {
	disabled := ""
	if props.Disabled {
		disabled = " disabled"
	}
	class := "btn"
	if props.Variant != "" {
		class += " btn-" + string(props.Variant)
	}
	if props.Size != "" {
		class += " btn-" + string(props.Size)
	}
	if props.ClassName != "" {
		class += " " + props.ClassName
	}
	return "<button" + attr("id", props.ID) + attr("class", class) + disabled + ">" + props.Children + "</button>"
}

// CardProps defines the properties for creating a card component.
type CardProps struct {
	ID        string
	Title     string
	Subtitle  string
	Body      string
	Footer    string
	Image     string
	ClassName string
}

// Card creates an HTML card element with Gocsx styling.
func Card(props CardProps) string {
	class := "card"
	if props.ClassName != "" {
		class += " " + props.ClassName
	}
	html := "<div" + attr("id", props.ID) + attr("class", class) + ">"
	if props.Image != "" {
		html += "<img src=\"" + props.Image + "\" class=\"card-img-top\" alt=\"\" />"
	}
	if props.Title != "" || props.Subtitle != "" {
		html += "<div class=\"card-body\">"
		if props.Title != "" {
			html += "<h5 class=\"card-title\">" + props.Title + "</h5>"
		}
		if props.Subtitle != "" {
			html += "<h6 class=\"card-subtitle\">" + props.Subtitle + "</h6>"
		}
		if props.Body != "" {
			html += "<p class=\"card-text\">" + props.Body + "</p>"
		}
		html += "</div>"
	}
	if props.Footer != "" {
		html += "<div class=\"card-footer\">" + props.Footer + "</div>"
	}
	html += "</div>"
	return html
}

// CardHeaderProps defines the properties for a card header.
type CardHeaderProps struct {
	Children  string
	ClassName string
}

// CardHeader creates an HTML card header element.
func CardHeader(props CardHeaderProps) string {
	class := "card-header"
	if props.ClassName != "" {
		class += " " + props.ClassName
	}
	return "<div class=\"" + class + "\">" + props.Children + "</div>"
}

// CardBodyProps defines the properties for a card body.
type CardBodyProps struct {
	Children  string
	ClassName string
}

// CardBody creates an HTML card body element.
func CardBody(props CardBodyProps) string {
	class := "card-body"
	if props.ClassName != "" {
		class += " " + props.ClassName
	}
	return "<div class=\"" + class + "\">" + props.Children + "</div>"
}

// CardFooterProps defines the properties for a card footer.
type CardFooterProps struct {
	Children  string
	ClassName string
}

// CardFooter creates an HTML card footer element.
func CardFooter(props CardFooterProps) string {
	class := "card-footer"
	if props.ClassName != "" {
		class += " " + props.ClassName
	}
	return "<div class=\"" + class + "\">" + props.Children + "</div>"
}

// CardTitleProps defines the properties for a card title.
type CardTitleProps struct {
	Children  string
	ClassName string
}

// CardTitle creates an HTML card title element.
func CardTitle(props CardTitleProps) string {
	class := "card-title"
	if props.ClassName != "" {
		class += " " + props.ClassName
	}
	return "<h5 class=\"" + class + "\">" + props.Children + "</h5>"
}

// CardSubtitleProps defines the properties for a card subtitle.
type CardSubtitleProps struct {
	Children  string
	ClassName string
}

// CardSubtitle creates an HTML card subtitle element.
func CardSubtitle(props CardSubtitleProps) string {
	class := "card-subtitle"
	if props.ClassName != "" {
		class += " " + props.ClassName
	}
	return "<h6 class=\"" + class + "\">" + props.Children + "</h6>"
}

// CardTextProps defines the properties for card text.
type CardTextProps struct {
	Children  string
	ClassName string
}

// CardText creates an HTML card text element.
func CardText(props CardTextProps) string {
	class := "card-text"
	if props.ClassName != "" {
		class += " " + props.ClassName
	}
	return "<p class=\"" + class + "\">" + props.Children + "</p>"
}

// CardImageProps defines the properties for a card image.
type CardImageProps struct {
	Src       string
	Alt       string
	ClassName string
}

// CardImage creates an HTML card image element.
func CardImage(props CardImageProps) string {
	class := "card-img-top"
	if props.ClassName != "" {
		class += " " + props.ClassName
	}
	return "<img src=\"" + props.Src + "\" alt=\"" + props.Alt + "\" class=\"" + class + "\" />"
}

// RegisterButtonComponent is a placeholder for future component registration.
func RegisterButtonComponent() string {
	return "button-registered"
}

// RegisterCardComponent is a placeholder for future component registration.
func RegisterCardComponent() string {
	return "card-registered"
}

// attr returns an HTML attribute string if value is non-empty.
func attr(name, value string) string {
	if value == "" {
		return ""
	}
	return " " + name + "=\"" + value + "\""
}
