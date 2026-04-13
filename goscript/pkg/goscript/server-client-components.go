package goscript

import (
        "fmt"
        "strings"
)

// ServerComponent renders entirely on the server and ships zero JavaScript to
// the client. It is ideal for SEO-critical pages, static content, and
// data-driven layouts where client-side interactivity is not required.
type ServerComponent struct {
        *BaseComponent
        name       string
        renderFunc func(Props) string
}

// NewServerComponent creates a new ServerComponent with the given name, render
// function, and initial props. The renderFunc receives Props and returns an
// HTML string. No client-side JavaScript is generated.
func NewServerComponent(name string, renderFunc func(Props) string, props Props) *ServerComponent {
        if props == nil {
                props = Props{}
        }
        base := NewBaseComponent(props, nil)
        return &ServerComponent{
                BaseComponent: base,
                name:          name,
                renderFunc:    renderFunc,
        }
}

// Render executes the server-side render function and returns the HTML string.
// This method produces pure HTML with no client-side hydration markers.
func (sc *ServerComponent) Render() string {
        return sc.renderFunc(sc.props)
}

// Name returns the component name.
func (sc *ServerComponent) Name() string {
        return sc.name
}

// ClientComponent includes a hydration marker and minimal client-side JavaScript
// runtime. It is intended for interactive UI elements that need browser-side
// event handling while still supporting server-side rendering.
type ClientComponent struct {
        *BaseComponent
        name       string
        props      Props
        events     map[string]string
        clientJS   string
}

// NewClientComponent creates a new ClientComponent with the given name and props.
// The component includes a hydration marker (`data-goscript-component`) and
// optionally a small JavaScript runtime for event handling.
func NewClientComponent(name string, props Props) *ClientComponent {
        if props == nil {
                props = Props{}
        }
        base := NewBaseComponent(props, nil)
        return &ClientComponent{
                BaseComponent: base,
                name:          name,
                props:         props,
                events:        make(map[string]string),
        }
}

// OnEvent registers a client-side event handler for the specified event type.
// The handler parameter is a string containing a JavaScript function body or a
// named function identifier that will be called when the event fires.
func (cc *ClientComponent) OnEvent(eventType string, handler string) {
        cc.events[eventType] = handler
}

// Render generates the server-rendered HTML output with hydration markers and
// embedded event handler attributes. The output includes a wrapper div with
// data attributes that the client-side runtime uses to attach interactivity.
func (cc *ClientComponent) Render() string {
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf(`<div data-goscript-component="%s"`, cc.name))

        // Serialize props as JSON data attribute
        if len(cc.props) > 0 {
                sb.WriteString(fmt.Sprintf(` data-goscript-props='{%s}'`, serializeProps(cc.props)))
        }

        // Attach event handler attributes
        for eventType, handler := range cc.events {
                sb.WriteString(fmt.Sprintf(` data-goscript-event-%s="%s"`, eventType, handler))
        }

        sb.WriteString(`>`)

        // Render children if present
        for _, child := range cc.children {
                sb.WriteString(renderChild(child))
        }

        sb.WriteString(`</div>`)

        // Append minimal client-side hydration script if events are registered
        if len(cc.events) > 0 {
                sb.WriteString(cc.renderHydrationScript())
        }

        return sb.String()
}

// renderHydrationScript generates a minimal JavaScript snippet that hydrates the
// client-side component by attaching event listeners based on data attributes.
func (cc *ClientComponent) renderHydrationScript() string {
        var sb strings.Builder
        sb.WriteString(`<script>(function(){`)
        sb.WriteString(fmt.Sprintf(`var el=document.querySelector('[data-goscript-component="%s"]');`, cc.name))
        sb.WriteString(`if(!el)return;`)
        sb.WriteString(`var props=el.getAttribute('data-goscript-props');`)

        for eventType := range cc.events {
                attrName := fmt.Sprintf(`data-goscript-event-%s`, eventType)
                sb.WriteString(fmt.Sprintf(`var %sFn=el.getAttribute('%s');`, eventType, attrName))
                sb.WriteString(fmt.Sprintf(`if(%sFn){el.addEventListener('%s',function(e){try{new Function(%sFn).call(el,e)}catch(ex){console.error('GoScript event error:',ex)}});}`, eventType, eventType, eventType))
        }

        sb.WriteString(`})();</script>`)
        return sb.String()
}

// Name returns the component name.
func (cc *ClientComponent) Name() string {
        return cc.name
}

// serializeProps converts Props into a compact JavaScript object literal string
// suitable for embedding in a data attribute.
func serializeProps(props Props) string {
        if len(props) == 0 {
                return ""
        }
        var sb strings.Builder
        first := true
        for key, value := range props {
                if !first {
                        sb.WriteString(",")
                }
                first = false

                switch v := value.(type) {
                case string:
                        sb.WriteString(fmt.Sprintf(`"%s":"%s"`, key, escapeJSString(v)))
                case bool:
                        sb.WriteString(fmt.Sprintf(`"%s":%t`, key, v))
                case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
                        sb.WriteString(fmt.Sprintf(`"%s":%v`, key, v))
                case nil:
                        sb.WriteString(fmt.Sprintf(`"%s":null`, key))
                default:
                        sb.WriteString(fmt.Sprintf(`"%s":"%v"`, key, escapeJSString(fmt.Sprintf("%v", v))))
                }
        }
        return sb.String()
}

// escapeJSString escapes a string for safe embedding in a JavaScript string literal.
func escapeJSString(s string) string {
        s = strings.Replace(s, `\`, `\\`, -1)
        s = strings.Replace(s, `"`, `\"`, -1)
        s = strings.Replace(s, `'`, `\'`, -1)
        s = strings.Replace(s, "\n", `\n`, -1)
        s = strings.Replace(s, "\r", `\r`, -1)
        s = strings.Replace(s, "\t", `\t`, -1)
        return s
}

// IsServerComponent returns true if the given Component is a ServerComponent.
// This type-assertion check enables runtime differentiation between server and
// client components for SSR optimization decisions.
func IsServerComponent(c Component) bool {
        _, ok := c.(*ServerComponent)
        return ok
}

// IsClientComponent returns true if the given Component is a ClientComponent.
// This type-assertion check enables the SSR engine to decide whether to include
// hydration markers and client-side JavaScript for the component.
func IsClientComponent(c Component) bool {
        _, ok := c.(*ClientComponent)
        return ok
}
