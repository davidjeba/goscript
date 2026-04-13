package goscript

import (
        "encoding/json"
        "fmt"
        "strings"
)

// OpenGraphMeta defines the Open Graph metadata for social media sharing.
// These meta tags control how a page appears when shared on platforms like
// Facebook, Twitter, LinkedIn, and Slack.
type OpenGraphMeta struct {
        Title       string
        Description string
        Type        string
        Image       string
        URL         string
        SiteName    string
        Locale      string
}

// TwitterMeta defines Twitter Card metadata for rich previews on Twitter/X.
type TwitterMeta struct {
        Card        string
        Title       string
        Description string
        Image       string
        Site        string
        Creator     string
}

// Metadata holds the complete metadata configuration for a page. It includes
// basic HTML metadata, Open Graph tags, Twitter Card tags, JSON-LD structured
// data, and robot directives.
type Metadata struct {
        Title       string
        Description string
        Canonical   string
        ThemeColor  string
        Keywords    []string
        Viewport    string
        Robots      robotsConfig
        OpenGraph   OpenGraphMeta
        Twitter     TwitterMeta
        JSONLD      []map[string]interface{}
        Extra       map[string]string
}

// robotsConfig holds the robots meta directive configuration.
type robotsConfig struct {
        index  bool
        follow bool
}

// MetadataBuilder provides a fluent API for constructing page Metadata.
// Each method returns the builder to allow method chaining.
type MetadataBuilder struct {
        metadata Metadata
}

// NewMetadata creates a new MetadataBuilder with sensible defaults.
func NewMetadata() *MetadataBuilder {
        return &MetadataBuilder{
                metadata: Metadata{
                        Keywords: make([]string, 0),
                        Viewport: "width=device-width, initial-scale=1",
                        Robots: robotsConfig{
                                index:  true,
                                follow: true,
                        },
                        Extra: make(map[string]string),
                },
        }
}

// SetTitle sets the page title. This is rendered as both <title> and og:title.
func (mb *MetadataBuilder) SetTitle(title string) *MetadataBuilder {
        mb.metadata.Title = title
        return mb
}

// SetDescription sets the page description used in the meta description tag
// and og:description.
func (mb *MetadataBuilder) SetDescription(desc string) *MetadataBuilder {
        mb.metadata.Description = desc
        return mb
}

// SetCanonical sets the canonical URL for the page to help search engines
// identify the primary version of duplicate content.
func (mb *MetadataBuilder) SetCanonical(url string) *MetadataBuilder {
        mb.metadata.Canonical = url
        return mb
}

// SetThemeColor sets the theme-color meta tag used by browsers to customize
// the browser chrome color on mobile devices.
func (mb *MetadataBuilder) SetThemeColor(color string) *MetadataBuilder {
        mb.metadata.ThemeColor = color
        return mb
}

// AddKeywords appends keywords to the page's keyword list.
func (mb *MetadataBuilder) AddKeywords(keywords ...string) *MetadataBuilder {
        mb.metadata.Keywords = append(mb.metadata.Keywords, keywords...)
        return mb
}

// SetOpenGraph sets the complete Open Graph metadata configuration.
func (mb *MetadataBuilder) SetOpenGraph(og OpenGraphMeta) *MetadataBuilder {
        mb.metadata.OpenGraph = og
        return mb
}

// SetTwitter sets the complete Twitter Card metadata configuration.
func (mb *MetadataBuilder) SetTwitter(tw TwitterMeta) *MetadataBuilder {
        mb.metadata.Twitter = tw
        return mb
}

// AddJSONLD appends a JSON-LD structured data block to the page metadata.
// Multiple JSON-LD blocks can be added for different structured data types
// (e.g., Article, Product, Organization).
func (mb *MetadataBuilder) AddJSONLD(data map[string]interface{}) *MetadataBuilder {
        mb.metadata.JSONLD = append(mb.metadata.JSONLD, data)
        return mb
}

// SetRobots configures the robots meta directives for indexing and following.
func (mb *MetadataBuilder) SetRobots(index, follow bool) *MetadataBuilder {
        mb.metadata.Robots = robotsConfig{index: index, follow: follow}
        return mb
}

// SetViewport sets the viewport meta tag content.
func (mb *MetadataBuilder) SetViewport(width, initialScale string) *MetadataBuilder {
        mb.metadata.Viewport = fmt.Sprintf("width=%s, initial-scale=%s", width, initialScale)
        return mb
}

// SetExtra adds an arbitrary extra meta tag with the given name and content.
func (mb *MetadataBuilder) SetExtra(name, content string) *MetadataBuilder {
        mb.metadata.Extra[name] = content
        return mb
}

// Build finalizes and returns the constructed Metadata.
func (mb *MetadataBuilder) Build() *Metadata {
        return &mb.metadata
}

// Render generates the complete HTML <head> content from the Metadata.
// It produces title, meta description, viewport, canonical, robots, theme-color,
// keywords, Open Graph tags, Twitter Card tags, JSON-LD scripts, and any
// extra meta tags.
func (m *Metadata) Render() string {
        var sb strings.Builder

        // Title
        if m.Title != "" {
                sb.WriteString(fmt.Sprintf("    <title>%s</title>\n", escapeHTML(m.Title)))
        }

        // Meta description
        if m.Description != "" {
                sb.WriteString(fmt.Sprintf(`    <meta name="description" content="%s">`+"\n", escapeHTML(m.Description)))
        }

        // Viewport
        if m.Viewport != "" {
                sb.WriteString(fmt.Sprintf(`    <meta name="viewport" content="%s">`+"\n", escapeHTML(m.Viewport)))
        }

        // Canonical URL
        if m.Canonical != "" {
                sb.WriteString(fmt.Sprintf(`    <link rel="canonical" href="%s">`+"\n", escapeHTML(m.Canonical)))
        }

        // Robots
        robots := "index, follow"
        if !m.Robots.index {
                robots = "noindex"
        }
        if !m.Robots.follow {
                if robots != "" {
                        robots += ", "
                }
                robots += "nofollow"
        }
        sb.WriteString(fmt.Sprintf(`    <meta name="robots" content="%s">`+"\n", robots))

        // Theme color
        if m.ThemeColor != "" {
                sb.WriteString(fmt.Sprintf(`    <meta name="theme-color" content="%s">`+"\n", escapeHTML(m.ThemeColor)))
        }

        // Keywords
        if len(m.Keywords) > 0 {
                sb.WriteString(fmt.Sprintf(`    <meta name="keywords" content="%s">`+"\n", escapeHTML(strings.Join(m.Keywords, ", "))))
        }

        // Open Graph tags
        if m.OpenGraph.Title != "" || m.OpenGraph.Description != "" {
                sb.WriteString("\n")
                if m.OpenGraph.Type != "" {
                        sb.WriteString(fmt.Sprintf(`    <meta property="og:type" content="%s">`+"\n", escapeHTML(m.OpenGraph.Type)))
                }
                if m.OpenGraph.Title != "" {
                        sb.WriteString(fmt.Sprintf(`    <meta property="og:title" content="%s">`+"\n", escapeHTML(m.OpenGraph.Title)))
                }
                if m.OpenGraph.Description != "" {
                        sb.WriteString(fmt.Sprintf(`    <meta property="og:description" content="%s">`+"\n", escapeHTML(m.OpenGraph.Description)))
                }
                if m.OpenGraph.Image != "" {
                        sb.WriteString(fmt.Sprintf(`    <meta property="og:image" content="%s">`+"\n", escapeHTML(m.OpenGraph.Image)))
                }
                if m.OpenGraph.URL != "" {
                        sb.WriteString(fmt.Sprintf(`    <meta property="og:url" content="%s">`+"\n", escapeHTML(m.OpenGraph.URL)))
                }
                if m.OpenGraph.SiteName != "" {
                        sb.WriteString(fmt.Sprintf(`    <meta property="og:site_name" content="%s">`+"\n", escapeHTML(m.OpenGraph.SiteName)))
                }
                if m.OpenGraph.Locale != "" {
                        sb.WriteString(fmt.Sprintf(`    <meta property="og:locale" content="%s">`+"\n", escapeHTML(m.OpenGraph.Locale)))
                }
        }

        // Twitter Card tags
        if m.Twitter.Card != "" || m.Twitter.Title != "" {
                sb.WriteString("\n")
                if m.Twitter.Card != "" {
                        sb.WriteString(fmt.Sprintf(`    <meta name="twitter:card" content="%s">`+"\n", escapeHTML(m.Twitter.Card)))
                }
                if m.Twitter.Title != "" {
                        sb.WriteString(fmt.Sprintf(`    <meta name="twitter:title" content="%s">`+"\n", escapeHTML(m.Twitter.Title)))
                }
                if m.Twitter.Description != "" {
                        sb.WriteString(fmt.Sprintf(`    <meta name="twitter:description" content="%s">`+"\n", escapeHTML(m.Twitter.Description)))
                }
                if m.Twitter.Image != "" {
                        sb.WriteString(fmt.Sprintf(`    <meta name="twitter:image" content="%s">`+"\n", escapeHTML(m.Twitter.Image)))
                }
                if m.Twitter.Site != "" {
                        sb.WriteString(fmt.Sprintf(`    <meta name="twitter:site" content="%s">`+"\n", escapeHTML(m.Twitter.Site)))
                }
                if m.Twitter.Creator != "" {
                        sb.WriteString(fmt.Sprintf(`    <meta name="twitter:creator" content="%s">`+"\n", escapeHTML(m.Twitter.Creator)))
                }
        }

        // JSON-LD structured data
        if len(m.JSONLD) > 0 {
                sb.WriteString("\n")
                for _, data := range m.JSONLD {
                        jsonBytes, err := json.MarshalIndent(data, "    ", "  ")
                        if err != nil {
                                continue
                        }
                        sb.WriteString(fmt.Sprintf(`    <script type="application/ld+json">`+"\n"))
                        sb.WriteString(string(jsonBytes))
                        sb.WriteString("\n    </script>\n")
                }
        }

        // Extra meta tags
        for name, content := range m.Extra {
                if content != "" {
                        sb.WriteString(fmt.Sprintf(`    <meta name="%s" content="%s">`+"\n", escapeHTML(name), escapeHTML(content)))
                }
        }

        return sb.String()
}

// escapeHTML escapes special characters in a string for safe embedding in
// HTML attribute values and text content.
func escapeHTML(s string) string {
        s = strings.Replace(s, "&", "&amp;", -1)
        s = strings.Replace(s, "<", "&lt;", -1)
        s = strings.Replace(s, ">", "&gt;", -1)
        s = strings.Replace(s, `"`, "&quot;", -1)
        s = strings.Replace(s, "'", "&#39;", -1)
        return s
}

// DefaultMetadata returns a pre-configured MetadataBuilder with commonly used
// defaults suitable for most web pages.
func DefaultMetadata() *MetadataBuilder {
        return NewMetadata().
                SetTitle("GoScript App").
                SetDescription("Built with GoScript 2.0 — Full-Stack Go Web Framework").
                SetViewport("device-width", "1.0").
                SetThemeColor("#ffffff").
                SetRobots(true, true).
                SetOpenGraph(OpenGraphMeta{
                        Type: "website",
                }).
                SetTwitter(TwitterMeta{
                        Card: "summary_large_image",
                })
}
