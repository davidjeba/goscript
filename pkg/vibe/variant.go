package vibe

import "fmt"

// StyleMap represents an animatable style target.
type StyleMap map[string]interface{}

// Variant groups a named style target with an optional transition.
type Variant struct {
	Name       string     `json:"name"`
	Style      StyleMap   `json:"style"`
	Transition Transition `json:"transition,omitempty"`
}

// VariantSet stores variants by name.
type VariantSet map[string]Variant

// GestureTargets maps UI states to variants.
type GestureTargets struct {
	Hover  string `json:"hover,omitempty"`
	Tap    string `json:"tap,omitempty"`
	Focus  string `json:"focus,omitempty"`
	Drag   string `json:"drag,omitempty"`
	InView string `json:"inView,omitempty"`
}

// MotionProps represents the declarative animation surface of a component.
type MotionProps struct {
	Initial     StyleMap       `json:"initial,omitempty"`
	Animate     StyleMap       `json:"animate,omitempty"`
	Exit        StyleMap       `json:"exit,omitempty"`
	Layout      bool           `json:"layout,omitempty"`
	LayoutID    string         `json:"layoutId,omitempty"`
	Variants    VariantSet     `json:"variants,omitempty"`
	Gestures    GestureTargets `json:"gestures,omitempty"`
	Transition  Transition     `json:"transition,omitempty"`
	WhileHover  StyleMap       `json:"whileHover,omitempty"`
	WhileTap    StyleMap       `json:"whileTap,omitempty"`
	WhileFocus  StyleMap       `json:"whileFocus,omitempty"`
	WhileDrag   StyleMap       `json:"whileDrag,omitempty"`
	WhileInView StyleMap       `json:"whileInView,omitempty"`
}

// MergeStyles merges style maps from left to right.
func MergeStyles(styles ...StyleMap) StyleMap {
	merged := StyleMap{}
	for _, style := range styles {
		for key, value := range style {
			merged[key] = value
		}
	}
	return merged
}

// ResolveVariant resolves one or more named variants into a single style map and transitions list.
func ResolveVariant(variants VariantSet, names ...string) (StyleMap, []Transition, error) {
	styles := make([]StyleMap, 0, len(names))
	transitions := make([]Transition, 0, len(names))

	for _, name := range names {
		variant, ok := variants[name]
		if !ok {
			return nil, nil, fmt.Errorf("unknown vibe variant %q", name)
		}

		styles = append(styles, variant.Style)
		if variant.Transition != (Transition{}) {
			transitions = append(transitions, variant.Transition.Normalize())
		}
	}

	return MergeStyles(styles...), transitions, nil
}
