package goscript

import (
	"fmt"
	"reflect"
	"regexp"
)

// FormField describes a field in a UI form.
type FormField struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Required  bool   `json:"required,omitempty"`
	MinLength int    `json:"minLength,omitempty"`
	MaxLength int    `json:"maxLength,omitempty"`
	Pattern   string `json:"pattern,omitempty"`
	Default   interface{} `json:"default,omitempty"`
}

// FormSchema describes a form validation contract.
type FormSchema map[string]FormField

// FormError describes a validation issue.
type FormError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// BindForm applies defaults and returns a normalized payload.
func BindForm(schema FormSchema, values map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	for key, field := range schema {
		if values != nil {
			if value, ok := values[key]; ok {
				out[key] = value
				continue
			}
		}

		if field.Default != nil {
			out[key] = field.Default
		}
	}

	if values != nil {
		for key, value := range values {
			if _, ok := out[key]; !ok {
				out[key] = value
			}
		}
	}

	return out
}

// ValidateForm validates payload values against a schema.
func ValidateForm(schema FormSchema, values map[string]interface{}) []FormError {
	normalized := BindForm(schema, values)
	errors := make([]FormError, 0)

	for key, field := range schema {
		value, exists := normalized[key]
		if field.Required && !exists {
			errors = append(errors, FormError{
				Field:   key,
				Message: "field is required",
			})
			continue
		}

		if !exists {
			continue
		}

		switch field.Type {
		case "string":
			str, ok := value.(string)
			if !ok {
				errors = append(errors, FormError{Field: key, Message: "must be a string"})
				continue
			}

			if field.MinLength > 0 && len(str) < field.MinLength {
				errors = append(errors, FormError{Field: key, Message: fmt.Sprintf("must be at least %d characters", field.MinLength)})
			}

			if field.MaxLength > 0 && len(str) > field.MaxLength {
				errors = append(errors, FormError{Field: key, Message: fmt.Sprintf("must be at most %d characters", field.MaxLength)})
			}

			if field.Pattern != "" {
				matched, err := regexp.MatchString(field.Pattern, str)
				if err != nil {
					errors = append(errors, FormError{Field: key, Message: fmt.Sprintf("invalid pattern: %v", err)})
				} else if !matched {
					errors = append(errors, FormError{Field: key, Message: "value does not match pattern"})
				}
			}
		case "int":
			if !isInteger(value) {
				errors = append(errors, FormError{Field: key, Message: "must be an integer"})
			}
		case "bool":
			if _, ok := value.(bool); !ok {
				errors = append(errors, FormError{Field: key, Message: "must be a boolean"})
			}
		case "any":
			// no-op
		default:
			if field.Type != "" {
				if reflect.TypeOf(value) == nil {
					continue
				}
			}
		}
	}

	return errors
}

func isInteger(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	default:
		return false
	}
}

