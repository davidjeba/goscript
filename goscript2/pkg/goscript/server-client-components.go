package goscript

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// RenderMode determines where a component renders
type RenderMode int

const (
	RenderModeServer RenderMode = iota // Server-only (default)
	RenderModeClient                   // Client-side with hydration
)

// ComponentMetadata holds runtime metadata for a component
type ComponentMetadata struct {
	Name        string
	RenderMode  RenderMode
	PropsSchema PropsSchema
	HasState    bool
	HasEffects  bool
}

// PropsSchema defines the JSON schema for component props
type PropsSchema map[string]PropSchemaField

type PropSchemaField struct {
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default,omitempty"`
	Description string      `json:"description,omitempty"`
}

// ServerComponent is a component that ONLY renders on the server.
// It never ships JS to the client.
type ServerComponent struct {
	BaseComponent
	metadata ComponentMetadata
	renderFn func(Props) string
	data     interface{}
}

// NewServerComponent creates a server-only component
func NewServerComponent(name string, renderFn func(Props) string, data interface{}) *ServerComponent {
	return &ServerComponent{
		metadata: ComponentMetadata{
			Name:       name,
			RenderMode: RenderModeServer,
		},
		renderFn: renderFn,
		data:     data,
	}
}

// Render executes the server render function — NEVER reaches the browser
func (sc *ServerComponent) Render() string {
	return sc.renderFn(sc.GetProps())
}

// ClientMetadata returns metadata safe to send to the client
func (sc *ServerComponent) ClientMetadata() ComponentMetadata {
	return ComponentMetadata{
		Name:        sc.metadata.Name,
		RenderMode:  sc.metadata.RenderMode,
		PropsSchema: sc.metadata.PropsSchema,
	}
}

// Serialize serializes server component output for the client
func (sc *ServerComponent) Serialize() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"html":  sc.Render(),
		"props": sc.GetProps(),
	})
}

// ClientComponent is a component that hydrates on the client.
// It ships minimal JS for interactivity.
type ClientComponent struct {
	BaseComponent
	metadata      ComponentMetadata
	eventHandlers map[string]string
	serializable  bool
}

// NewClientComponent creates a client-side component
func NewClientComponent(name string, props Props) *ClientComponent {
	base := NewBaseComponent(props, nil)
	return &ClientComponent{
		BaseComponent: *base,
		metadata: ComponentMetadata{
			Name:       name,
			RenderMode: RenderModeClient,
		},
		eventHandlers: make(map[string]string),
		serializable:  true,
	}
}

// OnEvent registers a client-side event handler
func (cc *ClientComponent) OnEvent(event string, handler string) *ClientComponent {
	cc.eventHandlers[event] = handler
	return cc
}

// Render generates HTML with hydration markers
func (cc *ClientComponent) Render() string {
	propsJSON, _ := json.Marshal(cc.GetProps())
	eventsJSON, _ := json.Marshal(cc.eventHandlers)

	return fmt.Sprintf(
		`<div data-gs-client="%s" data-gs-props='%s' data-gs-events='%s'></div><script>__gs_hydrate("%s",%s,%s)</script>`,
		cc.metadata.Name,
		string(propsJSON),
		string(eventsJSON),
		cc.metadata.Name,
		string(propsJSON),
		string(eventsJSON),
	)
}

// GetData returns the server data associated with a client component
func (cc *ClientComponent) GetData() interface{} {
	return nil
}

// IsSerializable returns whether the component can be serialized for transfer
func (cc *ClientComponent) IsSerializable() bool {
	return cc.serializable
}

// HydrateFromProps hydrates the component from serialized props
func (cc *ClientComponent) HydrateFromProps(propsJSON []byte) error {
	var props Props
	if err := json.Unmarshal(propsJSON, &props); err != nil {
		return err
	}
	cc.props = props
	return nil
}

// AsBridge creates a bridge component that can be shared between server and client
func AsBridge(sc *ServerComponent) *ClientComponent {
	propsCopy := make(Props)
	for k, v := range sc.GetProps() {
		propsCopy[k] = v
	}
	return NewClientComponent(sc.metadata.Name, propsCopy)
}

// MatchProps validates that the given props match the expected schema
func MatchProps(schema PropsSchema, props Props) bool {
	for field, def := range schema {
		_, exists := props[field]
		if def.Required && !exists {
			return false
		}
		if exists {
			pv := reflect.ValueOf(props[field])
			switch def.Type {
			case "string":
				if pv.Kind() != reflect.String {
					return false
				}
			case "number":
				if pv.Kind() != reflect.Float64 && pv.Kind() != reflect.Int {
					return false
				}
			case "boolean":
				if pv.Kind() != reflect.Bool {
					return false
				}
			}
		}
	}
	return true
}
