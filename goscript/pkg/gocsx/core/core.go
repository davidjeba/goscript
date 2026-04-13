// Package core provides the core types and engine for the Gocsx CSS framework.
package core

// Component represents a Gocsx UI component with a name and styling.
type Component struct {
	Name      string
	ClassName string
	HTML      string
	Styles    string
}

// Config holds global configuration for a Gocsx instance.
type Config struct {
	Prefix       string
	Theme        string
	DarkMode     bool
	Minify       bool
	CustomColors map[string]string
}

// Gocsx is the core engine for generating CSS and component markup.
type Gocsx struct {
	Config    *Config
	classes   []string
	styles    []string
	compIndex map[string]*Component
}

// New creates a new core Gocsx engine with the given options.
func New(options ...func(*Config)) *Gocsx {
	cfg := &Config{
		Prefix:   "gocsx",
		Theme:    "default",
		DarkMode: false,
		Minify:   false,
	}
	for _, opt := range options {
		opt(cfg)
	}
	return &Gocsx{
		Config:    cfg,
		classes:   make([]string, 0),
		styles:    make([]string, 0),
		compIndex: make(map[string]*Component),
	}
}

// PlatformAdapter is an interface for platform-specific rendering adapters.
type PlatformAdapter interface {
	Name() string
	Render(css string) string
}

// RegisterPlatformAdapter registers a platform-specific rendering adapter.
func (g *Gocsx) RegisterPlatformAdapter(name string, adapter PlatformAdapter) {
	g.compIndex["__adapter_"+name] = &Component{Name: name, HTML: ""}
}

// GetCSS returns the accumulated CSS output.
func (g *Gocsx) GetCSS() string {
	result := ""
	for _, s := range g.styles {
		result += s + "\n"
	}
	return result
}

// GenerateStyleTag returns an HTML <style> tag wrapping the CSS output.
func (g *Gocsx) GenerateStyleTag() string {
	return "<style>\n" + g.GetCSS() + "\n</style>"
}

// NewClassList creates a new ClassList for building CSS class strings.
func (g *Gocsx) NewClassList() *ClassList {
	return &ClassList{}
}

// ClassList is a builder for CSS class strings.
type ClassList struct {
	classes []string
}

// Add adds class names to the ClassList.
func (cl *ClassList) Add(classes ...string) *ClassList {
	cl.classes = append(cl.classes, classes...)
	return cl
}

// String returns the space-separated class string.
func (cl *ClassList) String() string {
	result := ""
	for i, c := range cl.classes {
		if i > 0 {
			result += " "
		}
		result += c
	}
	return result
}
