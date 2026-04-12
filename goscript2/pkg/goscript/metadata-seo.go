package goscript

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
)

// Metadata defines the metadata for a page
type Metadata struct {
	Title         string
	Description   string
	Canonical     string
	Keywords      []string
	Authors       []string
	OpenGraph     OpenGraphMeta
	Twitter       TwitterMeta
	Robots        RobotsMeta
	Viewport      string
	Charset       string
	ThemeColor    string
	Icons         []IconMeta
	Manifest      string
	AlternateLang []AlternateLangMeta
	JSONLD        []map[string]interface{}
	Scripts       []ScriptMeta
	Styles        []StyleMeta
}

// OpenGraphMeta defines Open Graph metadata
type OpenGraphMeta struct {
	Title       string
	Description string
	URL         string
	Type        string
	Image       string
	ImageWidth  int
	ImageHeight int
	SiteName    string
	Locale      string
}

// TwitterMeta defines Twitter Card metadata
type TwitterMeta struct {
	Card        string // "summary", "summary_large_image"
	Title       string
	Description string
	Image       string
	Site        string
	Creator     string
}

// RobotsMeta defines robots directives
type RobotsMeta struct {
	Index       bool
	Follow      bool
	NoArchive   bool
	NoSnippet   bool
	MaxSnippet  int
	MaxPreview  int
	MaxImage    int
}

// IconMeta defines a favicon or touch icon
type IconMeta struct {
	Rel  string
	Href string
	Sizes string
	Type string
}

// ScriptMeta defines an external script tag
type ScriptMeta struct {
	Src     string
	Defer   bool
	Async   bool
	Type    string
	Integrity string
}

// StyleMeta defines an external stylesheet link
type StyleMeta struct {
	Href       string
	Rel        string
	Integrity  string
	CrossOrigin string
}

// AlternateLangMeta defines an alternate language link
type AlternateLangMeta struct {
	HrefLang string
	Href     string
}

// MetadataBuilder provides a fluent API for constructing metadata
type MetadataBuilder struct {
	m *Metadata
}

// NewMetadata creates a new metadata builder
func NewMetadata() *MetadataBuilder {
	return &MetadataBuilder{
		m: &Metadata{
			Keywords:      make([]string, 0),
			Authors:       make([]string, 0),
			Icons:         make([]IconMeta, 0),
			JSONLD:        make([]map[string]interface{}, 0),
			Scripts:       make([]ScriptMeta, 0),
			Styles:        make([]StyleMeta, 0),
			AlternateLang: make([]AlternateLangMeta, 0),
			Viewport:      "width=device-width, initial-scale=1",
			Charset:       "utf-8",
			Robots:        RobotsMeta{Index: true, Follow: true},
		},
	}
}

// SetTitle sets the page title
func (b *MetadataBuilder) SetTitle(title string) *MetadataBuilder {
	b.m.Title = title
	return b
}

// SetDescription sets the page description
func (b *MetadataBuilder) SetDescription(desc string) *MetadataBuilder {
	b.m.Description = desc
	return b
}

// SetCanonical sets the canonical URL
func (b *MetadataBuilder) SetCanonical(url string) *MetadataBuilder {
	b.m.Canonical = url
	return b
}

// AddKeywords adds keywords
func (b *MetadataBuilder) AddKeywords(keywords ...string) *MetadataBuilder {
	b.m.Keywords = append(b.m.Keywords, keywords...)
	return b
}

// SetOpenGraph configures Open Graph metadata
func (b *MetadataBuilder) SetOpenGraph(og OpenGraphMeta) *MetadataBuilder {
	b.m.OpenGraph = og
	return b
}

// SetTwitter configures Twitter Card metadata
func (b *MetadataBuilder) SetTwitter(tw TwitterMeta) *MetadataBuilder {
	b.m.Twitter = tw
	return b
}

// AddJSONLD adds structured data
func (b *MetadataBuilder) AddJSONLD(data map[string]interface{}) *MetadataBuilder {
	b.m.JSONLD = append(b.m.JSONLD, data)
	return b
}

// SetThemeColor sets the theme color
func (b *MetadataBuilder) SetThemeColor(color string) *MetadataBuilder {
	b.m.ThemeColor = color
	return b
}

// SetViewport sets the viewport meta content
func (b *MetadataBuilder) SetViewport(viewport string) *MetadataBuilder {
	b.m.Viewport = viewport
	return b
}

// SetManifest sets the web app manifest URL
func (b *MetadataBuilder) SetManifest(url string) *MetadataBuilder {
	b.m.Manifest = url
	return b
}

// SetRobots configures robots directives
func (b *MetadataBuilder) SetRobots(robots RobotsMeta) *MetadataBuilder {
	b.m.Robots = robots
	return b
}

// Build returns the configured metadata
func (b *MetadataBuilder) Build() *Metadata {
	return b.m
}

// Render generates the complete HTML <head> section
func (m *Metadata) Render() string {
	var sb strings.Builder

	if m.Charset != "" {
		sb.WriteString(fmt.Sprintf(`<meta charset="%s">`, m.Charset))
	}
	if m.Viewport != "" {
		sb.WriteString(fmt.Sprintf(`<meta name="viewport" content="%s">`, template.HTMLEscapeString(m.Viewport)))
	}
	if m.Title != "" {
		sb.WriteString(fmt.Sprintf(`<title>%s</title>`, template.HTMLEscapeString(m.Title)))
	}
	if m.Description != "" {
		sb.WriteString(fmt.Sprintf(`<meta name="description" content="%s">`, template.HTMLEscapeString(m.Description)))
	}
	if m.Canonical != "" {
		sb.WriteString(fmt.Sprintf(`<link rel="canonical" href="%s">`, template.HTMLEscapeString(m.Canonical)))
	}
	if m.ThemeColor != "" {
		sb.WriteString(fmt.Sprintf(`<meta name="theme-color" content="%s">`, template.HTMLEscapeString(m.ThemeColor)))
	}
	if m.Manifest != "" {
		sb.WriteString(fmt.Sprintf(`<link rel="manifest" href="%s">`, template.HTMLEscapeString(m.Manifest)))
	}

	// Keywords
	if len(m.Keywords) > 0 {
		sb.WriteString(fmt.Sprintf(`<meta name="keywords" content="%s">`, template.HTMLEscapeString(strings.Join(m.Keywords, ", "))))
	}

	// Authors
	for _, author := range m.Authors {
		sb.WriteString(fmt.Sprintf(`<meta name="author" content="%s">`, template.HTMLEscapeString(author)))
	}

	// Icons
	for _, icon := range m.Icons {
		sb.WriteString(fmt.Sprintf(`<link rel="%s" href="%s"`, template.HTMLEscapeString(icon.Rel), template.HTMLEscapeString(icon.Href)))
		if icon.Sizes != "" {
			sb.WriteString(fmt.Sprintf(` sizes="%s"`, template.HTMLEscapeString(icon.Sizes)))
		}
		if icon.Type != "" {
			sb.WriteString(fmt.Sprintf(` type="%s"`, template.HTMLEscapeString(icon.Type)))
		}
		sb.WriteString(`>`)
	}

	// Robots
	if !m.Robots.Index || !m.Robots.Follow {
		var directives []string
		if !m.Robots.Index {
			directives = append(directives, "noindex")
		}
		if !m.Robots.Follow {
			directives = append(directives, "nofollow")
		}
		if m.Robots.NoArchive {
			directives = append(directives, "noarchive")
		}
		if m.Robots.NoSnippet {
			directives = append(directives, "nosnippet")
		}
		if len(directives) > 0 {
			sb.WriteString(fmt.Sprintf(`<meta name="robots" content="%s">`, strings.Join(directives, ", ")))
		}
	}

	// Open Graph
	if m.OpenGraph.Title != "" {
		sb.WriteString(fmt.Sprintf(`<meta property="og:title" content="%s">`, template.HTMLEscapeString(m.OpenGraph.Title)))
	}
	if m.OpenGraph.Description != "" {
		sb.WriteString(fmt.Sprintf(`<meta property="og:description" content="%s">`, template.HTMLEscapeString(m.OpenGraph.Description)))
	}
	if m.OpenGraph.URL != "" {
		sb.WriteString(fmt.Sprintf(`<meta property="og:url" content="%s">`, template.HTMLEscapeString(m.OpenGraph.URL)))
	}
	if m.OpenGraph.Type != "" {
		sb.WriteString(fmt.Sprintf(`<meta property="og:type" content="%s">`, template.HTMLEscapeString(m.OpenGraph.Type)))
	}
	if m.OpenGraph.Image != "" {
		sb.WriteString(fmt.Sprintf(`<meta property="og:image" content="%s">`, template.HTMLEscapeString(m.OpenGraph.Image)))
	}
	if m.OpenGraph.SiteName != "" {
		sb.WriteString(fmt.Sprintf(`<meta property="og:site_name" content="%s">`, template.HTMLEscapeString(m.OpenGraph.SiteName)))
	}
	if m.OpenGraph.Locale != "" {
		sb.WriteString(fmt.Sprintf(`<meta property="og:locale" content="%s">`, template.HTMLEscapeString(m.OpenGraph.Locale)))
	}

	// Twitter
	if m.Twitter.Card != "" {
		sb.WriteString(fmt.Sprintf(`<meta name="twitter:card" content="%s">`, template.HTMLEscapeString(m.Twitter.Card)))
	}
	if m.Twitter.Title != "" {
		sb.WriteString(fmt.Sprintf(`<meta name="twitter:title" content="%s">`, template.HTMLEscapeString(m.Twitter.Title)))
	}
	if m.Twitter.Description != "" {
		sb.WriteString(fmt.Sprintf(`<meta name="twitter:description" content="%s">`, template.HTMLEscapeString(m.Twitter.Description)))
	}
	if m.Twitter.Image != "" {
		sb.WriteString(fmt.Sprintf(`<meta name="twitter:image" content="%s">`, template.HTMLEscapeString(m.Twitter.Image)))
	}
	if m.Twitter.Site != "" {
		sb.WriteString(fmt.Sprintf(`<meta name="twitter:site" content="%s">`, template.HTMLEscapeString(m.Twitter.Site)))
	}
	if m.Twitter.Creator != "" {
		sb.WriteString(fmt.Sprintf(`<meta name="twitter:creator" content="%s">`, template.HTMLEscapeString(m.Twitter.Creator)))
	}

	// Alternate languages
	for _, alt := range m.AlternateLang {
		sb.WriteString(fmt.Sprintf(`<link rel="alternate" hreflang="%s" href="%s">`, template.HTMLEscapeString(alt.HrefLang), template.HTMLEscapeString(alt.Href)))
	}

	// JSON-LD structured data
	for _, ld := range m.JSONLD {
		b, _ := json.Marshal(ld)
		sb.WriteString(fmt.Sprintf(`<script type="application/ld+json">%s</script>`, string(b)))
	}

	// External scripts
	for _, script := range m.Scripts {
		attrs := make([]string, 0)
		attrs = append(attrs, fmt.Sprintf(`src="%s"`, template.HTMLEscapeString(script.Src)))
		if script.Defer {
			attrs = append(attrs, "defer")
		}
		if script.Async {
			attrs = append(attrs, "async")
		}
		if script.Type != "" {
			attrs = append(attrs, fmt.Sprintf(`type="%s"`, template.HTMLEscapeString(script.Type)))
		}
		if script.Integrity != "" {
			attrs = append(attrs, fmt.Sprintf(`integrity="%s"`, template.HTMLEscapeString(script.Integrity)))
		}
		sb.WriteString(fmt.Sprintf(`<script %s></script>`, strings.Join(attrs, " ")))
	}

	// External styles
	for _, style := range m.Styles {
		rel := style.Rel
		if rel == "" {
			rel = "stylesheet"
		}
		sb.WriteString(fmt.Sprintf(`<link rel="%s" href="%s"`, template.HTMLEscapeString(rel), template.HTMLEscapeString(style.Href)))
		if style.Integrity != "" {
			sb.WriteString(fmt.Sprintf(` integrity="%s"`, template.HTMLEscapeString(style.Integrity)))
		}
		if style.CrossOrigin != "" {
			sb.WriteString(fmt.Sprintf(` crossorigin="%s"`, template.HTMLEscapeString(style.CrossOrigin)))
		}
		sb.WriteString(`>`)
	}

	return sb.String()
}

// ToMap converts metadata to a map for JSON serialization
func (m *Metadata) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
	if m.Title != "" {
		result["title"] = m.Title
	}
	if m.Description != "" {
		result["description"] = m.Description
	}
	if m.Canonical != "" {
		result["canonical"] = m.Canonical
	}
	if len(m.Keywords) > 0 {
		result["keywords"] = m.Keywords
	}
	return result
}
