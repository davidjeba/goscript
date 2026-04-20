package goscript

import "reflect"

// MatchCase represents one branch in a pattern match.
type MatchCase struct {
	Equals    interface{}
	Kind      reflect.Kind
	Predicate func(interface{}) bool
	Then      func(interface{}) interface{}
}

// Match evaluates the first branch that fits the value.
func Match(value interface{}, cases ...MatchCase) interface{} {
	for _, c := range cases {
		if c.Predicate != nil && c.Predicate(value) {
			if c.Then != nil {
				return c.Then(value)
			}
			return value
		}

		if c.Equals != nil && reflect.DeepEqual(value, c.Equals) {
			if c.Then != nil {
				return c.Then(value)
			}
			return value
		}

		if c.Kind != reflect.Invalid {
			if value != nil && reflect.TypeOf(value).Kind() == c.Kind {
				if c.Then != nil {
					return c.Then(value)
				}
				return value
			}
		}
	}

	return nil
}

