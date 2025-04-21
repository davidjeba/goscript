# Gocsx - A Comprehensive CSS Framework for Go

Gocsx (pronounced "gosix") is a powerful, utility-first CSS framework for Go that enables developers to build beautiful, responsive interfaces for web, mobile, and AR/VR applications. Drawing inspiration from Tailwind CSS but extending far beyond its capabilities, Gocsx provides a unified styling approach across multiple platforms.

![Gocsx Logo](https://via.placeholder.com/800x400?text=Gocsx)

## Features

### Core Features

- **Utility-First Approach**: Build complex designs without leaving your Go code
- **Cross-Platform Support**: One framework for web, mobile, and AR/VR applications
- **Type-Safe Styling**: Leverage Go's type system for safer styling
- **Component System**: Pre-built, customizable components
- **Responsive Design**: Built-in responsive utilities
- **Dark Mode**: First-class dark mode support
- **Theme System**: Powerful theming capabilities
- **CSS-in-Go**: Write CSS directly in your Go code

### Platform-Specific Features

#### Web
- **CSS Generation**: Optimized CSS output
- **CSS Variables**: Dynamic theming with CSS variables
- **Media Queries**: Responsive design utilities
- **Print Styles**: Optimized print styles
- **CSS Grid**: Advanced layout capabilities
- **Flexbox**: Flexible box layout
- **Animations**: CSS animations and transitions

#### Mobile
- **Native Styling**: Translates to native mobile styles
- **Touch Interactions**: Optimized for touch interfaces
- **Platform Adaptations**: Adapts to iOS and Android conventions
- **Responsive to Screen Size**: Adapts to different device sizes
- **Native Components**: Styles that work with native components

#### AR/VR
- **Spatial Layouts**: 3D positioning and layout
- **Immersive UI**: Styles for immersive interfaces
- **Gaze Interactions**: Support for gaze-based interactions
- **Hand Tracking**: Support for hand tracking interactions
- **3D Typography**: Typography optimized for 3D space

## Installation

```bash
go get github.com/davidjeba/goscript/pkg/gocsx
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/davidjeba/goscript/pkg/gocsx"
    "github.com/davidjeba/goscript/pkg/gocsx/components"
)

func main() {
    // Create a new Gocsx instance
    g := gocsx.New()

    // Create a button
    button := g.Button(components.ButtonProps{
        ID:       "my-button",
        Variant:  components.ButtonPrimary,
        Size:     components.ButtonSizeLarge,
        Children: "Click Me",
        OnClick:  "alert('Hello, Gocsx!')",
    })

    // Generate CSS
    css := g.GetCSS()

    // Use the button and CSS in your application
    fmt.Println(button)
    fmt.Println(css)
}
```

### Creating a Web Page

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/davidjeba/goscript/pkg/gocsx"
    "github.com/davidjeba/goscript/pkg/gocsx/components"
)

func main() {
    // Create a new Gocsx instance
    g := gocsx.New()

    // Create a button
    button := g.Button(components.ButtonProps{
        ID:       "my-button",
        Variant:  components.ButtonPrimary,
        Size:     components.ButtonSizeLarge,
        Children: "Click Me",
        OnClick:  "alert('Hello, Gocsx!')",
    })

    // Create a card
    card := g.Card(components.CardProps{
        ID:        "my-card",
        Title:     "Gocsx Card",
        Subtitle:  "A powerful CSS framework for Go",
        Body:      "This is a card component built with Gocsx.",
        Footer:    "Footer content",
        ClassName: "shadow",
    })

    // Create a page
    page := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Gocsx Demo</title>
    %s
</head>
<body>
    <div class="container py-5">
        <h1 class="text-center mb-5">Gocsx Demo</h1>
        <div class="row">
            <div class="col-md-6 offset-md-3">
                <div class="text-center mb-4">
                    %s
                </div>
                %s
            </div>
        </div>
    </div>
</body>
</html>
    `, g.GenerateStyleTag(), button, card)

    // Create a handler for the page
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html")
        fmt.Fprint(w, page)
    })

    // Start the server
    log.Println("Server starting on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Components

Gocsx comes with a set of pre-built components that you can use in your applications:

### Button

```go
button := g.Button(components.ButtonProps{
    ID:       "my-button",
    Variant:  components.ButtonPrimary,
    Size:     components.ButtonSizeLarge,
    Children: "Click Me",
    OnClick:  "alert('Hello, Gocsx!')",
    Disabled: false,
    FullWidth: false,
})
```

### Card

```go
card := g.Card(components.CardProps{
    ID:        "my-card",
    Title:     "Gocsx Card",
    Subtitle:  "A powerful CSS framework for Go",
    Body:      "This is a card component built with Gocsx.",
    Footer:    "Footer content",
    Image:     "https://example.com/image.jpg",
    ImageAlt:  "Example Image",
    ClassName: "shadow",
})
```

### Alert

```go
alert := g.Alert(components.AlertProps{
    ID:       "my-alert",
    Variant:  components.AlertSuccess,
    Children: "This is a success alert!",
    Dismissible: true,
})
```

### Form Controls

```go
input := g.Input(components.InputProps{
    ID:          "my-input",
    Type:        "text",
    Label:       "Username",
    Placeholder: "Enter your username",
    Required:    true,
    HelperText:  "Your username must be 5-20 characters long",
})

select := g.Select(components.SelectProps{
    ID:      "my-select",
    Label:   "Country",
    Options: []components.SelectOption{
        {Value: "us", Label: "United States"},
        {Value: "ca", Label: "Canada"},
        {Value: "mx", Label: "Mexico"},
    },
})

checkbox := g.Checkbox(components.CheckboxProps{
    ID:      "my-checkbox",
    Label:   "I agree to the terms and conditions",
    Checked: false,
})
```

## Utility Classes

Gocsx provides a wide range of utility classes for styling your components:

### Layout

- `container`, `row`, `col`, `col-{1-12}`, `col-md-{1-12}`
- `d-flex`, `d-inline-flex`, `d-block`, `d-inline-block`, `d-none`
- `flex-row`, `flex-column`, `flex-wrap`, `flex-nowrap`
- `justify-content-start`, `justify-content-center`, `justify-content-end`, `justify-content-between`, `justify-content-around`
- `align-items-start`, `align-items-center`, `align-items-end`, `align-items-baseline`, `align-items-stretch`

### Spacing

- `m-{0-5}`, `mt-{0-5}`, `mr-{0-5}`, `mb-{0-5}`, `ml-{0-5}`, `mx-{0-5}`, `my-{0-5}`
- `p-{0-5}`, `pt-{0-5}`, `pr-{0-5}`, `pb-{0-5}`, `pl-{0-5}`, `px-{0-5}`, `py-{0-5}`

### Typography

- `text-{primary, secondary, success, danger, warning, info, light, dark, white, muted}`
- `text-{left, center, right, justify}`
- `text-{lowercase, uppercase, capitalize}`
- `font-weight-{normal, bold}`
- `font-italic`

### Borders

- `border`, `border-{top, right, bottom, left}`
- `border-{primary, secondary, success, danger, warning, info, light, dark}`
- `rounded`, `rounded-{top, right, bottom, left, circle, pill}`

### Colors

- `bg-{primary, secondary, success, danger, warning, info, light, dark, white, transparent}`
- `text-{primary, secondary, success, danger, warning, info, light, dark, white, muted}`

### Sizing

- `w-{25, 50, 75, 100, auto}`, `h-{25, 50, 75, 100, auto}`
- `mw-100`, `mh-100`

### Shadows

- `shadow-sm`, `shadow`, `shadow-lg`, `shadow-none`

## Customization

### Theming

```go
// Create a custom theme
theme := &core.ThemeConfig{
    Colors: map[string]map[string]string{
        "primary": {
            "50":  "#f0f9ff",
            "100": "#e0f2fe",
            "500": "#0ea5e9",
            "900": "#0c4a6e",
        },
        "brand": {
            "500": "#ff0000",
        },
    },
    // Add more theme customizations here
}

// Create a new Gocsx instance with the custom theme
g := gocsx.New(core.WithTheme(theme))
```

### Custom Components

```go
// Define a custom component
type MyComponentProps struct {
    ID       string
    Title    string
    Content  string
    ClassName string
}

func MyComponent(g *gocsx.Gocsx, props MyComponentProps) string {
    // Build class names
    classes := []string{"my-component", "p-4", "border", "rounded"}
    
    if props.ClassName != "" {
        classes = append(classes, props.ClassName)
    }
    
    // Build the component
    return fmt.Sprintf(`
        <div id="%s" class="%s">
            <h3 class="my-component-title mb-3">%s</h3>
            <div class="my-component-content">%s</div>
        </div>
    `, props.ID, strings.Join(classes, " "), props.Title, props.Content)
}
```

## Platform-Specific Usage

### Web

```go
// Create a web-specific configuration
config := core.NewConfig(
    core.WithPlatform(core.PlatformConfig{
        Target: "web",
        Features: map[string]bool{
            "darkMode": true,
            "rtl": false,
        },
    }),
)

// Create a new Gocsx instance with the web configuration
g := gocsx.New(config)
```

### Mobile

```go
// Create a mobile-specific configuration
config := core.NewConfig(
    core.WithPlatform(core.PlatformConfig{
        Target: "mobile",
        Features: map[string]bool{
            "touchEvents": true,
            "nativeComponents": true,
        },
    }),
)

// Create a new Gocsx instance with the mobile configuration
g := gocsx.New(config)
```

### AR/VR

```go
// Create an AR/VR-specific configuration
config := core.NewConfig(
    core.WithPlatform(core.PlatformConfig{
        Target: "ar",
        Features: map[string]bool{
            "spatialLayout": true,
            "gazeInteraction": true,
        },
    }),
)

// Create a new Gocsx instance with the AR/VR configuration
g := gocsx.New(config)
```

## Why Gocsx?

### 100% Go

Gocsx is written entirely in Go, allowing you to build full-stack applications without context switching between languages.

### Cross-Platform

Unlike most CSS frameworks that are web-only, Gocsx is designed to work across web, mobile, and AR/VR platforms.

### Type Safety

Gocsx leverages Go's type system to provide type-safe styling, catching errors at compile time rather than runtime.

### Performance

Gocsx generates optimized CSS with minimal overhead, resulting in faster load times and better performance.

### Developer Experience

Gocsx provides a familiar API for developers coming from other CSS frameworks like Tailwind, Bootstrap, or Material UI.

## Comparison with Other Frameworks

### Gocsx vs Tailwind CSS

- **Language**: Gocsx uses Go, Tailwind uses JavaScript/CSS
- **Platforms**: Gocsx supports web, mobile, and AR/VR, Tailwind is web-only
- **Type Safety**: Gocsx has type safety, Tailwind does not
- **Components**: Gocsx has built-in components, Tailwind requires additional libraries
- **Customization**: Both have powerful customization options

### Gocsx vs Bootstrap

- **Approach**: Gocsx is utility-first, Bootstrap is component-first
- **Size**: Gocsx generates only the CSS you need, Bootstrap includes everything
- **Customization**: Gocsx has more granular customization options
- **Platforms**: Gocsx supports multiple platforms, Bootstrap is web-only
- **Language**: Gocsx uses Go, Bootstrap uses JavaScript/CSS

### Gocsx vs Material UI

- **Design System**: Gocsx is framework-agnostic, Material UI follows Google's Material Design
- **Language**: Gocsx uses Go, Material UI uses JavaScript/React
- **Platforms**: Gocsx supports multiple platforms, Material UI is web-focused
- **Customization**: Both have powerful customization options
- **Components**: Both have rich component libraries

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License