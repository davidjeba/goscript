package gocsx

import (
	"github.com/davidjeba/goscript/pkg/gocsx/components"
	"github.com/davidjeba/goscript/pkg/gocsx/core"
	"github.com/davidjeba/goscript/pkg/gocsx/platforms/web"
)

// Gocsx is the main entry point for the Gocsx framework
type Gocsx struct {
	// Core instance
	Core *core.Gocsx

	// Components
	Button *core.Component
	Card   *core.Component
}

// New creates a new Gocsx instance
func New(options ...func(*core.Config)) *Gocsx {
	// Create core instance
	coreInstance := core.New(options...)

	// Register platform adapters
	webAdapter := web.NewWebAdapter(coreInstance.Config)
	coreInstance.RegisterPlatformAdapter("web", webAdapter)

	// Create Gocsx instance
	gocsx := &Gocsx{
		Core: coreInstance,
	}

	// Register components
	gocsx.Button = components.RegisterButtonComponent(coreInstance)
	gocsx.Card = components.RegisterCardComponent(coreInstance)

	return gocsx
}

// GetCSS gets the generated CSS
func (g *Gocsx) GetCSS() string {
	return g.Core.GetCSS()
}

// GenerateStyleTag generates a style tag with the CSS
func (g *Gocsx) GenerateStyleTag() string {
	return g.Core.GenerateStyleTag()
}

// Button creates a button component
func (g *Gocsx) Button(props components.ButtonProps) string {
	return components.Button(props)
}

// Card creates a card component
func (g *Gocsx) Card(props components.CardProps) string {
	return components.Card(props)
}

// CardHeader creates a card header component
func (g *Gocsx) CardHeader(props components.CardHeaderProps) string {
	return components.CardHeader(props)
}

// CardBody creates a card body component
func (g *Gocsx) CardBody(props components.CardBodyProps) string {
	return components.CardBody(props)
}

// CardFooter creates a card footer component
func (g *Gocsx) CardFooter(props components.CardFooterProps) string {
	return components.CardFooter(props)
}

// CardTitle creates a card title component
func (g *Gocsx) CardTitle(props components.CardTitleProps) string {
	return components.CardTitle(props)
}

// CardSubtitle creates a card subtitle component
func (g *Gocsx) CardSubtitle(props components.CardSubtitleProps) string {
	return components.CardSubtitle(props)
}

// CardText creates a card text component
func (g *Gocsx) CardText(props components.CardTextProps) string {
	return components.CardText(props)
}

// CardImage creates a card image component
func (g *Gocsx) CardImage(props components.CardImageProps) string {
	return components.CardImage(props)
}

// cx is a shorthand function for creating a class list
func (g *Gocsx) cx(classes ...string) string {
	return g.Core.NewClassList().Add(classes...).String()
}

// cxIf is a shorthand function for conditionally adding a class
func (g *Gocsx) cxIf(condition bool, trueClass, falseClass string) string {
	if condition {
		return trueClass
	}
	return falseClass
}