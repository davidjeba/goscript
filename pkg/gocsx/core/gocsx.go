package core

import (
	"fmt"
	"strings"
	"sync"
)

// Gocsx is the main entry point for the Gocsx framework
type Gocsx struct {
	// Configuration
	Config *Config

	// CSS generator
	Generator *Generator

	// Class cache
	classCache map[string]bool
	cacheMutex sync.RWMutex

	// Generated CSS
	css string
	cssMutex sync.RWMutex

	// Platform adapters
	platformAdapters map[string]PlatformAdapter
}

// PlatformAdapter is an interface for platform-specific adapters
type PlatformAdapter interface {
	// TransformCSS transforms CSS for the specific platform
	TransformCSS(css string) string

	// TransformClass transforms a class name for the specific platform
	TransformClass(class string) string

	// GetPlatformSpecificClasses returns platform-specific classes
	GetPlatformSpecificClasses() []string
}

// New creates a new Gocsx instance
func New(options ...func(*Config)) *Gocsx {
	config := NewConfig(options...)
	generator := NewGenerator(config)

	// Register default utilities and variants
	generator.RegisterDefaultUtilities()
	generator.RegisterDefaultVariants()

	return &Gocsx{
		Config:           config,
		Generator:        generator,
		classCache:       make(map[string]bool),
		platformAdapters: make(map[string]PlatformAdapter),
	}
}

// RegisterPlatformAdapter registers a platform adapter
func (g *Gocsx) RegisterPlatformAdapter(platform string, adapter PlatformAdapter) {
	g.platformAdapters[platform] = adapter
}

// GetPlatformAdapter gets a platform adapter
func (g *Gocsx) GetPlatformAdapter(platform string) (PlatformAdapter, bool) {
	adapter, ok := g.platformAdapters[platform]
	return adapter, ok
}

// AddClasses adds classes to the cache
func (g *Gocsx) AddClasses(classes ...string) {
	g.cacheMutex.Lock()
	defer g.cacheMutex.Unlock()

	for _, class := range classes {
		g.classCache[class] = true
	}

	// Regenerate CSS
	g.regenerateCSS()
}

// HasClass checks if a class is in the cache
func (g *Gocsx) HasClass(class string) bool {
	g.cacheMutex.RLock()
	defer g.cacheMutex.RUnlock()

	return g.classCache[class]
}

// GetCSS gets the generated CSS
func (g *Gocsx) GetCSS() string {
	g.cssMutex.RLock()
	defer g.cssMutex.RUnlock()

	return g.css
}

// regenerateCSS regenerates the CSS
func (g *Gocsx) regenerateCSS() {
	// Get all classes
	var classes []string
	for class := range g.classCache {
		classes = append(classes, class)
	}

	// Add platform-specific classes
	if adapter, ok := g.GetPlatformAdapter(g.Config.Platform.Target); ok {
		platformClasses := adapter.GetPlatformSpecificClasses()
		classes = append(classes, platformClasses...)
	}

	// Generate CSS
	css := g.Generator.GenerateCSS(classes)

	// Transform CSS for the platform
	if adapter, ok := g.GetPlatformAdapter(g.Config.Platform.Target); ok {
		css = adapter.TransformCSS(css)
	}

	// Update CSS
	g.cssMutex.Lock()
	g.css = css
	g.cssMutex.Unlock()
}

// ClassList represents a list of classes
type ClassList struct {
	// Classes
	classes []string

	// Gocsx instance
	gocsx *Gocsx
}

// NewClassList creates a new class list
func (g *Gocsx) NewClassList() *ClassList {
	return &ClassList{
		classes: []string{},
		gocsx:   g,
	}
}

// Add adds classes to the list
func (c *ClassList) Add(classes ...string) *ClassList {
	c.classes = append(c.classes, classes...)
	c.gocsx.AddClasses(classes...)
	return c
}

// AddIf adds a class if the condition is true
func (c *ClassList) AddIf(condition bool, class string) *ClassList {
	if condition {
		c.Add(class)
	}
	return c
}

// AddUnless adds a class unless the condition is true
func (c *ClassList) AddUnless(condition bool, class string) *ClassList {
	if !condition {
		c.Add(class)
	}
	return c
}

// AddWhen adds classes based on a map of conditions
func (c *ClassList) AddWhen(conditions map[string]bool) *ClassList {
	for class, condition := range conditions {
		c.AddIf(condition, class)
	}
	return c
}

// Remove removes classes from the list
func (c *ClassList) Remove(classes ...string) *ClassList {
	for _, class := range classes {
		for i, c := range c.classes {
			if c == class {
				c.classes = append(c.classes[:i], c.classes[i+1:]...)
				break
			}
		}
	}
	return c
}

// Toggle toggles a class
func (c *ClassList) Toggle(class string) *ClassList {
	for i, c := range c.classes {
		if c == class {
			c.classes = append(c.classes[:i], c.classes[i+1:]...)
			return c
		}
	}
	c.Add(class)
	return c
}

// Has checks if the list has a class
func (c *ClassList) Has(class string) bool {
	for _, c := range c.classes {
		if c == class {
			return true
		}
	}
	return false
}

// String returns the class list as a string
func (c *ClassList) String() string {
	return strings.Join(c.classes, " ")
}

// Component represents a Gocsx component
type Component struct {
	// Name of the component
	Name string

	// Base classes
	BaseClasses []string

	// Variant classes
	VariantClasses map[string][]string

	// Gocsx instance
	gocsx *Gocsx
}

// NewComponent creates a new component
func (g *Gocsx) NewComponent(name string, baseClasses []string, variantClasses map[string][]string) *Component {
	component := &Component{
		Name:          name,
		BaseClasses:   baseClasses,
		VariantClasses: variantClasses,
		gocsx:         g,
	}

	// Add all classes to the cache
	g.AddClasses(baseClasses...)
	for _, classes := range variantClasses {
		g.AddClasses(classes...)
	}

	return component
}

// GetClasses gets the classes for a component with variants
func (c *Component) GetClasses(variants ...string) string {
	classList := c.gocsx.NewClassList()

	// Add base classes
	classList.Add(c.BaseClasses...)

	// Add variant classes
	for _, variant := range variants {
		if classes, ok := c.VariantClasses[variant]; ok {
			classList.Add(classes...)
		}
	}

	return classList.String()
}

// GetClassList gets the class list for a component with variants
func (c *Component) GetClassList(variants ...string) *ClassList {
	classList := c.gocsx.NewClassList()

	// Add base classes
	classList.Add(c.BaseClasses...)

	// Add variant classes
	for _, variant := range variants {
		if classes, ok := c.VariantClasses[variant]; ok {
			classList.Add(classes...)
		}
	}

	return classList
}

// RegisterComponent registers a component with the generator
func (g *Gocsx) RegisterComponent(name string, baseClasses []string, variantClasses map[string][]string) *Component {
	return g.NewComponent(name, baseClasses, variantClasses)
}

// cx is a shorthand function for creating a class list
func (g *Gocsx) cx(classes ...string) string {
	classList := g.NewClassList()
	classList.Add(classes...)
	return classList.String()
}

// cxIf is a shorthand function for conditionally adding a class
func (g *Gocsx) cxIf(condition bool, trueClass, falseClass string) string {
	if condition {
		return trueClass
	}
	return falseClass
}

// GenerateStyleTag generates a style tag with the CSS
func (g *Gocsx) GenerateStyleTag() string {
	return fmt.Sprintf("<style>\n%s</style>", g.GetCSS())
}

// GenerateStylesheet generates a stylesheet with the CSS
func (g *Gocsx) GenerateStylesheet(filename string) error {
	// Implementation depends on the platform
	if adapter, ok := g.GetPlatformAdapter(g.Config.Platform.Target); ok {
		css := adapter.TransformCSS(g.GetCSS())
		// Write CSS to file
		// This is platform-specific and would be implemented in the adapter
	}
	return nil
}