package goscript

import (
	"fmt"
	"strings"
)

type Props map[string]interface{}

type Component interface {
	Render() string
}

type FunctionalComponent func(props Props) string

func (f FunctionalComponent) Render() string {
	return f(nil)
}

func createElement(component interface{}, props Props, children ...interface{}) string {
	var result strings.Builder

	switch c := component.(type) {
	case string:
		result.WriteString("<")
		result.WriteString(c)
		
		for key, value := range props {
			result.WriteString(fmt.Sprintf(" %s=\"%v\"", key, value))
		}
		
		if len(children) == 0 {
			result.WriteString("/>")
		} else {
			result.WriteString(">")
			for _, child := range children {
				switch ch := child.(type) {
				case Component:
					result.WriteString(ch.Render())
				case string:
					result.WriteString(ch)
				}
			}
			result.WriteString("</")
			result.WriteString(c)
			result.WriteString(">")
		}
	case Component:
		result.WriteString(c.Render())
	case FunctionalComponent:
		result.WriteString(c(props))
	}

	return result.String()
}

// Example of a stateful component
func Counter(props Props) string {
	count, setCount := UseState("count", 0)

	incrementCount := func() {
		setCount(count.(int) + 1)
	}

	return createElement("div", nil,
		createElement("p", nil, fmt.Sprintf("Count: %d", count)),
		createElement("button", Props{"onClick": incrementCount}, "Increment"),
	)
}

