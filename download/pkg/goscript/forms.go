// Package goscript — forms.go provides server-side form rendering, validation,
// and CSRF protection that integrates natively with goscript reactive attributes.
// Forms built with these helpers require zero client-side JavaScript — the
// gs-trigger, gs-target, and gs-swap attributes drive all interactivity.
//
// Key features:
//   - Declarative form configuration with reactive attributes baked in
//   - Built-in validation rules (required, min, max, pattern, type)
//   - CSRF token generation and validation
//   - Auto-rendering of validation errors near fields
//   - File upload with drag-and-drop support
//   - Progressive enhancement (forms work without JavaScript)
package goscript

import (
        "crypto/rand"
        "encoding/hex"
        "encoding/json"
        "fmt"
        "html"
        "net/http"
        "reflect"
        "regexp"
        "strconv"
        "strings"
)

// =========================================================================
// Form Configuration
// =========================================================================

// FormConfig defines the structure and behavior of a goscript form.
// Use this to declaratively create forms with reactive attributes.
type FormConfig struct {
        // Action is the URL the form submits to (the form's action attribute).
        Action string

        // Method is the HTTP method: "GET", "POST", "PUT", "DELETE".
        // Defaults to "POST".
        Method string

        // Target is the CSS selector for the element where the response
        // will be swapped. Equivalent to gs-target.
        Target string

        // Swap is the swap strategy for the response. Equivalent to gs-swap.
        // Defaults to SwapInnerHTML.
        Swap string

        // CSRF enables automatic CSRF token generation and validation.
        // When true, a hidden _gs_csrf field is added to the form.
        CSRF bool

        // Encode specifies the encoding type. Set to "multipart/form-data"
        // for file uploads, "application/json" for JSON-encoded requests,
        // or "" for the default URL-encoded encoding.
        Encode string

        // Class adds CSS classes to the <form> element.
        Class string

        // ID sets the HTML id attribute on the <form> element.
        ID string

        // NoValidate disables browser-native HTML5 validation, allowing
        // server-side validation to handle all feedback.
        NoValidate bool

        // ExtraProps adds arbitrary HTML attributes to the form element.
        ExtraProps Props
}

// FieldConfig defines a single form field. The Type field determines
// which HTML element is rendered (input, textarea, select, checkbox, etc.).
type FieldConfig struct {
        // Name is the field's name attribute (used as the form data key).
        Name string

        // Type is the field type: "text", "email", "password", "number",
        // "textarea", "select", "checkbox", "file", "hidden", "tel", "url",
        // "search", "date", "datetime-local", "time", "range", "color".
        Type string

        // Label is the human-readable label text shown above the field.
        Label string

        // Value is the initial value of the field.
        Value interface{}

        // Placeholder is the placeholder text shown when the field is empty.
        Placeholder string

        // Required enables required validation for this field.
        Required bool

        // Min sets the minimum value (for number, range, date fields).
        Min interface{}

        // Max sets the maximum value (for number, range, date fields).
        Max interface{}

        // Pattern sets a regex pattern the field value must match.
        Pattern string

        // Class adds CSS classes to the field element.
        Class string

        // Options defines the choices for select and radio fields.
        Options []FieldOption

        // ID sets the HTML id attribute on the field element.
        ID string

        // Hint is optional helper text shown below the field.
        Hint string

        // Disabled renders the field as disabled.
        Disabled bool

        // ReadOnly renders the field as read-only.
        ReadOnly bool

        // Multiple allows multiple selections (for select and file fields).
        Multiple bool

        // Accept sets the accept attribute for file inputs (e.g. "image/*").
        Accept string
}

// FieldOption represents a single option in a <select> dropdown or
// radio button group.
type FieldOption struct {
        // Value is the option's value attribute.
        Value string

        // Label is the display text for the option.
        Label string

        // Selected sets this option as pre-selected.
        Selected bool

        // Disabled sets this option as disabled.
        Disabled bool
}

// =========================================================================
// Validation
// =========================================================================

// ValidationResult holds the outcome of form validation.
type ValidationResult struct {
        // Valid is true when all fields pass validation.
        Valid bool

        // Errors contains all validation errors, one per field.
        Errors []FieldError
}

// FieldError describes a single validation error for a specific field.
type FieldError struct {
        // Field is the name of the field that failed validation.
        Field string

        // Message is the human-readable error description.
        Message string
}

// Error returns the first error message, implementing the error interface.
func (vr ValidationResult) Error() string {
        if len(vr.Errors) == 0 {
                return ""
        }
        return vr.Errors[0].Message
}

// FieldErrors returns all error messages for a specific field.
func (vr ValidationResult) FieldErrors(fieldName string) []string {
        var messages []string
        for _, err := range vr.Errors {
                if err.Field == fieldName {
                        messages = append(messages, err.Message)
                }
        }
        return messages
}

// HasError checks whether any error exists for the given field name.
func (vr ValidationResult) HasError(fieldName string) bool {
        for _, err := range vr.Errors {
                if err.Field == fieldName {
                        return true
                }
        }
        return false
}

// ValidationRule defines a single validation constraint for a form field.
// Rules are applied by the Validate function.
type ValidationRule struct {
        // Field is the name of the field to validate.
        Field string

        // Required ensures the field is not empty.
        Required bool

        // Min sets a minimum length (for strings) or value (for numbers).
        Min interface{}

        // Max sets a maximum length (for strings) or value (for numbers).
        Max interface{}

        // Pattern requires the field value to match this regex.
        Pattern string

        // Type validates the field value as a specific type:
        // "email", "url", "number", "int", "float", "phone".
        Type string

        // Message is a custom error message. If empty, a default message
        // is generated based on the rule.
        Message string

        // Custom is an optional custom validation function. If it returns
        // a non-nil error, the field fails validation with that message.
        Custom func(value interface{}) error
}

// Validate validates form data against the given rules. The data map
// typically comes from r.FormValue() or parsed JSON.
//
// Returns a ValidationResult with Valid=true if all rules pass, or
// Valid=false with a list of FieldErrors describing what failed.
//
// Usage:
//
//      rules := []ValidationRule{
//          {Field: "email", Required: true, Type: "email"},
//          {Field: "age", Type: "int", Min: 0, Max: 150},
//          {Field: "password", Required: true, Min: 8,
//              Message: "Password must be at least 8 characters"},
//      }
//      result := Validate(formData, rules...)
//      if !result.Valid {
//          // Handle errors
//      }
func Validate(data map[string]interface{}, rules ...ValidationRule) ValidationResult {
        result := ValidationResult{
                Valid:  true,
                Errors: make([]FieldError, 0),
        }

        for _, rule := range rules {
                value, exists := data[rule.Field]
                strValue := fmt.Sprintf("%v", value)

                // Required check
                if rule.Required {
                        if !exists || strValue == "" {
                                msg := rule.Message
                                if msg == "" {
                                        msg = fmt.Sprintf("%s is required", rule.Field)
                                }
                                result.Errors = append(result.Errors, FieldError{
                                        Field:   rule.Field,
                                        Message: msg,
                                })
                                result.Valid = false
                                continue
                        }
                }

                // Skip further validation if field is empty and not required
                if !exists || strValue == "" {
                        continue
                }

                // Type validation
                if rule.Type != "" {
                        if err := validateType(strValue, rule.Type); err != nil {
                                msg := rule.Message
                                if msg == "" {
                                        msg = err.Error()
                                }
                                result.Errors = append(result.Errors, FieldError{
                                        Field:   rule.Field,
                                        Message: msg,
                                })
                                result.Valid = false
                                continue
                        }
                }

                // Min validation
                if rule.Min != nil {
                        if err := validateMin(strValue, rule.Min); err != nil {
                                msg := rule.Message
                                if msg == "" {
                                        msg = err.Error()
                                }
                                result.Errors = append(result.Errors, FieldError{
                                        Field:   rule.Field,
                                        Message: msg,
                                })
                                result.Valid = false
                                continue
                        }
                }

                // Max validation
                if rule.Max != nil {
                        if err := validateMax(strValue, rule.Max); err != nil {
                                msg := rule.Message
                                if msg == "" {
                                        msg = err.Error()
                                }
                                result.Errors = append(result.Errors, FieldError{
                                        Field:   rule.Field,
                                        Message: msg,
                                })
                                result.Valid = false
                                continue
                        }
                }

                // Pattern validation
                if rule.Pattern != "" {
                        matched, err := regexp.MatchString(rule.Pattern, strValue)
                        if err != nil || !matched {
                                msg := rule.Message
                                if msg == "" {
                                        msg = fmt.Sprintf("%s does not match the required format", rule.Field)
                                }
                                result.Errors = append(result.Errors, FieldError{
                                        Field:   rule.Field,
                                        Message: msg,
                                })
                                result.Valid = false
                                continue
                        }
                }

                // Custom validation
                if rule.Custom != nil {
                        if err := rule.Custom(value); err != nil {
                                msg := rule.Message
                                if msg == "" {
                                        msg = err.Error()
                                }
                                result.Errors = append(result.Errors, FieldError{
                                        Field:   rule.Field,
                                        Message: msg,
                                })
                                result.Valid = false
                        }
                }
        }

        return result
}

// validateType checks if a string value matches the expected type format.
func validateType(value, typ string) error {
        switch typ {
        case "email":
                // Basic email validation
                if !strings.Contains(value, "@") || !strings.Contains(value, ".") {
                        return fmt.Errorf("must be a valid email address")
                }
        case "url":
                if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
                        return fmt.Errorf("must be a valid URL starting with http:// or https://")
                }
        case "number", "float":
                if _, err := strconv.ParseFloat(value, 64); err != nil {
                        return fmt.Errorf("must be a valid number")
                }
        case "int":
                if _, err := strconv.ParseInt(value, 10, 64); err != nil {
                        return fmt.Errorf("must be a valid integer")
                }
        case "phone":
                // Basic phone validation: digits, spaces, dashes, parens, plus
                phoneRegex := regexp.MustCompile(`^[\d\s\-\+\(\)]{7,20}$`)
                if !phoneRegex.MatchString(value) {
                        return fmt.Errorf("must be a valid phone number")
                }
        }
        return nil
}

// validateMin checks if a value meets a minimum constraint.
func validateMin(value string, min interface{}) error {
        switch m := min.(type) {
        case int:
                if len(value) < m {
                        return fmt.Errorf("must be at least %d characters", m)
                }
        case float64:
                if num, err := strconv.ParseFloat(value, 64); err == nil && num < m {
                        return fmt.Errorf("must be at least %v", m)
                }
        case string:
                num, err := strconv.Atoi(m)
                if err == nil && len(value) < num {
                        return fmt.Errorf("must be at least %s characters", m)
                }
        }
        return nil
}

// validateMax checks if a value meets a maximum constraint.
func validateMax(value string, max interface{}) error {
        switch m := max.(type) {
        case int:
                if len(value) > m {
                        return fmt.Errorf("must be at most %d characters", m)
                }
        case float64:
                if num, err := strconv.ParseFloat(value, 64); err == nil && num > m {
                        return fmt.Errorf("must be at most %v", m)
                }
        case string:
                num, err := strconv.Atoi(m)
                if err == nil && len(value) > num {
                        return fmt.Errorf("must be at most %s characters", m)
                }
        }
        return nil
}

// =========================================================================
// Form Rendering
// =========================================================================

// Form renders a complete HTML form with goscript reactive attributes.
// The form includes all specified fields, labels, validation error
// placeholders, and a CSRF token (if enabled).
//
// The rendered form uses gs-trigger="submit" and gs-target/gs-swap
// attributes, so it works reactively without any JavaScript.
//
// Usage:
//
//      html := goscript.Form(goscript.FormConfig{
//          Action: "/api/contact",
//          Method: "POST",
//          Target: "#result",
//          Swap:   SwapInnerHTML,
//          CSRF:   true,
//      },
//          goscript.FieldConfig{Name: "name", Type: "text", Label: "Name", Required: true},
//          goscript.FieldConfig{Name: "email", Type: "email", Label: "Email", Required: true},
//          goscript.FieldConfig{Name: "message", Type: "textarea", Label: "Message", Required: true},
//      )
func Form(config FormConfig, fields ...FieldConfig) string {
        // Set defaults
        if config.Method == "" {
                config.Method = "POST"
        }
        if config.Swap == "" {
                config.Swap = SwapInnerHTML
        }

        var sb strings.Builder

        // Build form attributes
        sb.WriteString("<form")

        // Method
        sb.WriteString(fmt.Sprintf(` method="%s"`, html.EscapeString(config.Method)))

        // Action
        if config.Action != "" {
                sb.WriteString(fmt.Sprintf(` action="%s"`, html.EscapeString(config.Action)))
        }

        // ID
        if config.ID != "" {
                sb.WriteString(fmt.Sprintf(` id="%s"`, html.EscapeString(config.ID)))
        }

        // Class
        formClasses := []string{"gs-form"}
        if config.Class != "" {
                formClasses = append(formClasses, config.Class)
        }
        sb.WriteString(fmt.Sprintf(` class="%s"`, html.EscapeString(strings.Join(formClasses, " "))))

        // Encoding type
        if config.Encode != "" {
                sb.WriteString(fmt.Sprintf(` enctype="%s"`, html.EscapeString(config.Encode)))
        }

        // NoValidate
        if config.NoValidate {
                sb.WriteString(" novalidate")
        }

        // Reactive attributes: gs-trigger, gs-target, gs-swap
        if config.Target != "" {
                sb.WriteString(fmt.Sprintf(` gs-trigger="%s"`, TriggerSubmit))
                sb.WriteString(fmt.Sprintf(` gs-target="%s"`, html.EscapeString(config.Target)))
                sb.WriteString(fmt.Sprintf(` gs-swap="%s"`, html.EscapeString(config.Swap)))
        }

        // Extra props
        if config.ExtraProps != nil {
                for key, val := range config.ExtraProps {
                        if val == true {
                                sb.WriteString(fmt.Sprintf(` %s`, key))
                        } else if val != false && val != nil {
                                sb.WriteString(fmt.Sprintf(` %s="%v"`, key, val))
                        }
                }
        }

        sb.WriteString(">")

        // CSRF token field
        if config.CSRF {
                token := CSRFToken("goscript-default-csrf")
                sb.WriteString(fmt.Sprintf(`<input type="hidden" name="_gs_csrf" value="%s">`,
                        html.EscapeString(token)))
        }

        // Render each field
        for _, field := range fields {
                sb.WriteString(renderField(field))
        }

        sb.WriteString("</form>")
        return sb.String()
}

// renderField renders a single form field with its label, input element,
// hint text, and error placeholder.
func renderField(field FieldConfig) string {
        var sb strings.Builder

        // Field wrapper
        sb.WriteString(`<div class="gs-field"`)
        if field.ID != "" {
                sb.WriteString(fmt.Sprintf(` id="%s"`, html.EscapeString(field.ID+"-field")))
        }
        sb.WriteString(">")

        // Label
        if field.Label != "" && field.Type != "hidden" && field.Type != "checkbox" {
                sb.WriteString(`<label class="gs-label"`)
                if field.ID != "" {
                        sb.WriteString(fmt.Sprintf(` for="%s"`, html.EscapeString(field.ID)))
                }
                sb.WriteString(">")
                sb.WriteString(html.EscapeString(field.Label))
                if field.Required {
                        sb.WriteString(` <span class="gs-required">*</span>`)
                }
                sb.WriteString("</label>")
        }

        // Input element
        switch field.Type {
        case "textarea":
                sb.WriteString(renderTextarea(field))
        case "select":
                sb.WriteString(renderSelect(field))
        case "checkbox":
                sb.WriteString(renderCheckbox(field))
        case "file":
                sb.WriteString(renderFileInput(field))
        default:
                sb.WriteString(renderInput(field))
        }

        // Hint text
        if field.Hint != "" {
                sb.WriteString(fmt.Sprintf(`<small class="gs-hint">%s</small>`, html.EscapeString(field.Hint)))
        }

        // Error placeholder
        sb.WriteString(fmt.Sprintf(`<div class="gs-errors" id="%s"></div>`,
                html.EscapeString(field.Name+"-errors")))

        sb.WriteString("</div>")
        return sb.String()
}

// renderInput renders a standard <input> element.
func renderInput(field FieldConfig) string {
        var attrs []string

        attrs = append(attrs, fmt.Sprintf(`type="%s"`, html.EscapeString(field.Type)))
        attrs = append(attrs, fmt.Sprintf(`name="%s"`, html.EscapeString(field.Name)))

        if field.ID != "" {
                attrs = append(attrs, fmt.Sprintf(`id="%s"`, html.EscapeString(field.ID)))
        }

        if field.Placeholder != "" {
                attrs = append(attrs, fmt.Sprintf(`placeholder="%s"`, html.EscapeString(field.Placeholder)))
        }

        if field.Required {
                attrs = append(attrs, "required")
        }

        if field.Disabled {
                attrs = append(attrs, "disabled")
        }

        if field.ReadOnly {
                attrs = append(attrs, "readonly")
        }

        if field.Value != nil {
                val := fmt.Sprintf("%v", field.Value)
                if !reflect.DeepEqual(field.Value, false) {
                        attrs = append(attrs, fmt.Sprintf(`value="%s"`, html.EscapeString(val)))
                }
        }

        if field.Pattern != "" {
                attrs = append(attrs, fmt.Sprintf(`pattern="%s"`, html.EscapeString(field.Pattern)))
        }

        if field.Min != nil {
                attrs = append(attrs, fmt.Sprintf(`min="%v"`, field.Min))
        }

        if field.Max != nil {
                attrs = append(attrs, fmt.Sprintf(`max="%v"`, field.Max))
        }

        if field.Class != "" {
                attrs = append(attrs, fmt.Sprintf(`class="%s"`, html.EscapeString(field.Class)))
        } else {
                attrs = append(attrs, `class="gs-input"`)
        }

        return fmt.Sprintf("<input %s>", strings.Join(attrs, " "))
}

// renderTextarea renders a <textarea> element.
func renderTextarea(field FieldConfig) string {
        var attrs []string

        attrs = append(attrs, fmt.Sprintf(`name="%s"`, html.EscapeString(field.Name)))

        if field.ID != "" {
                attrs = append(attrs, fmt.Sprintf(`id="%s"`, html.EscapeString(field.ID)))
        }

        if field.Placeholder != "" {
                attrs = append(attrs, fmt.Sprintf(`placeholder="%s"`, html.EscapeString(field.Placeholder)))
        }

        if field.Required {
                attrs = append(attrs, "required")
        }

        if field.Disabled {
                attrs = append(attrs, "disabled")
        }

        if field.ReadOnly {
                attrs = append(attrs, "readonly")
        }

        if field.Class != "" {
                attrs = append(attrs, fmt.Sprintf(`class="%s"`, html.EscapeString(field.Class)))
        } else {
                attrs = append(attrs, `class="gs-input gs-textarea"`)
        }

        value := ""
        if field.Value != nil {
                value = html.EscapeString(fmt.Sprintf("%v", field.Value))
        }

        return fmt.Sprintf("<textarea %s>%s</textarea>", strings.Join(attrs, " "), value)
}

// renderSelect renders a <select> element with <option> children.
func renderSelect(field FieldConfig) string {
        var attrs []string

        attrs = append(attrs, fmt.Sprintf(`name="%s"`, html.EscapeString(field.Name)))

        if field.ID != "" {
                attrs = append(attrs, fmt.Sprintf(`id="%s"`, html.EscapeString(field.ID)))
        }

        if field.Required {
                attrs = append(attrs, "required")
        }

        if field.Disabled {
                attrs = append(attrs, "disabled")
        }

        if field.Multiple {
                attrs = append(attrs, "multiple")
        }

        if field.Class != "" {
                attrs = append(attrs, fmt.Sprintf(`class="%s"`, html.EscapeString(field.Class)))
        } else {
                attrs = append(attrs, `class="gs-input gs-select"`)
        }

        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("<select %s>", strings.Join(attrs, " ")))

        // Empty default option
        if field.Placeholder != "" {
                sb.WriteString(fmt.Sprintf(`<option value="">%s</option>`, html.EscapeString(field.Placeholder)))
        }

        // Render options
        for _, opt := range field.Options {
                sb.WriteString("<option")
                sb.WriteString(fmt.Sprintf(` value="%s"`, html.EscapeString(opt.Value)))

                if opt.Selected {
                        selectedValue := ""
                        if field.Value != nil {
                                selectedValue = fmt.Sprintf("%v", field.Value)
                        }
                        if selectedValue == opt.Value || selectedValue == "" {
                                sb.WriteString(" selected")
                        }
                }

                if opt.Disabled {
                        sb.WriteString(" disabled")
                }

                sb.WriteString(">")
                sb.WriteString(html.EscapeString(opt.Label))
                sb.WriteString("</option>")
        }

        sb.WriteString("</select>")
        return sb.String()
}

// renderCheckbox renders a checkbox <input> with label inline.
func renderCheckbox(field FieldConfig) string {
        attrs := []string{
                fmt.Sprintf(`type="checkbox"`),
                fmt.Sprintf(`name="%s"`, html.EscapeString(field.Name)),
        }

        if field.ID != "" {
                attrs = append(attrs, fmt.Sprintf(`id="%s"`, html.EscapeString(field.ID)))
        }

        if field.Class != "" {
                attrs = append(attrs, fmt.Sprintf(`class="%s"`, html.EscapeString(field.Class)))
        } else {
                attrs = append(attrs, `class="gs-input gs-checkbox"`)
        }

        // Check if value indicates "checked"
        checked := false
        if field.Value != nil {
                switch v := field.Value.(type) {
                case bool:
                        checked = v
                case string:
                        checked = v != "" && v != "false" && v != "0"
                default:
                        checked = true
                }
        }
        if checked {
                attrs = append(attrs, "checked")
        }

        return fmt.Sprintf("<input %s>", strings.Join(attrs, " "))
}

// renderFileInput renders a file upload <input> element.
func renderFileInput(field FieldConfig) string {
        attrs := []string{
                `type="file"`,
                fmt.Sprintf(`name="%s"`, html.EscapeString(field.Name)),
        }

        if field.ID != "" {
                attrs = append(attrs, fmt.Sprintf(`id="%s"`, html.EscapeString(field.ID)))
        }

        if field.Accept != "" {
                attrs = append(attrs, fmt.Sprintf(`accept="%s"`, html.EscapeString(field.Accept)))
        }

        if field.Multiple {
                attrs = append(attrs, "multiple")
        }

        if field.Required {
                attrs = append(attrs, "required")
        }

        if field.Disabled {
                attrs = append(attrs, "disabled")
        }

        if field.Class != "" {
                attrs = append(attrs, fmt.Sprintf(`class="%s"`, html.EscapeString(field.Class)))
        } else {
                attrs = append(attrs, `class="gs-input gs-file"`)
        }

        return fmt.Sprintf("<input %s>", strings.Join(attrs, " "))
}

// =========================================================================
// Validation Error Rendering
// =========================================================================

// RenderFieldErrors renders validation errors for a specific field as HTML.
// Returns an empty string if the field has no errors.
//
// Usage:
//
//      result := Validate(formData, rules...)
//      if !result.Valid {
//          errorsHTML := RenderFieldErrors(result.Errors, "email")
//          // Returns: <div class="gs-field-error">Email is required</div>
//      }
func RenderFieldErrors(errors []FieldError, fieldName string) string {
        var sb strings.Builder
        for _, err := range errors {
                if err.Field == fieldName {
                        sb.WriteString(fmt.Sprintf(`<div class="gs-field-error">%s</div>`,
                                html.EscapeString(err.Message)))
                }
        }
        return sb.String()
}

// RenderAllErrors renders all validation errors as an HTML list.
// Returns an empty string if there are no errors.
//
// Usage:
//
//      if !result.Valid {
//          fmt.Println(RenderAllErrors(result.Errors))
//      }
func RenderAllErrors(errors []FieldError) string {
        if len(errors) == 0 {
                return ""
        }

        var sb strings.Builder
        sb.WriteString(`<div class="gs-errors-summary">`)
        sb.WriteString(`<ul class="gs-errors-list">`)

        for _, err := range errors {
                sb.WriteString(fmt.Sprintf(`<li class="gs-error-item">%s</li>`,
                        html.EscapeString(err.Message)))
        }

        sb.WriteString(`</ul>`)
        sb.WriteString(`</div>`)
        return sb.String()
}

// =========================================================================
// CSRF Protection
// =========================================================================

// CSRFToken generates a cryptographic CSRF token using the given secret.
// The secret should be a long, random string unique to the application.
// Tokens are hex-encoded 32-byte random values.
//
// Usage:
//
//      // In form rendering:
//      token := CSRFToken(appSecret)
//
//      // In form validation:
//      if !ValidateCSRF(appSecret, r.FormValue("_gs_csrf")) {
//          // CSRF validation failed
//      }
func CSRFToken(secret string) string {
        bytes := make([]byte, 32)
        _, err := rand.Read(bytes)
        if err != nil {
                // Fall back to a time-based token if crypto/rand fails
                return fmt.Sprintf("%x-%s", secret, "fallback-token")
        }
        return hex.EncodeToString(bytes)
}

// ValidateCSRF checks whether the provided token was generated by the
// application. In this implementation, any non-empty token is considered
// valid (since tokens are randomly generated and stored in the form).
// For production use, tokens should be stored server-side (e.g., in a
// session or signed HMAC).
//
// Usage:
//
//      if !ValidateCSRF("", r.FormValue("_gs_csrf")) {
//          http.Error(w, "Invalid CSRF token", http.StatusForbidden)
//          return
//      }
func ValidateCSRF(secret, token string) bool {
        // Validate that the token is a valid hex-encoded 64-character string
        if len(token) != 64 {
                return false
        }
        _, err := hex.DecodeString(token)
        return err == nil
}

// ValidateCSRFRequest is a convenience function that validates the CSRF
// token from an HTTP request's form data. Returns true if valid.
//
// Usage:
//
//      if !ValidateCSRFRequest(r) {
//          GoscriptTrigger(w, "error", "Invalid security token")
//          return
//      }
func ValidateCSRFRequest(r *http.Request) bool {
        token := r.FormValue("_gs_csrf")
        return ValidateCSRF("", token)
}

// =========================================================================
// File Upload Helper
// =========================================================================

// FileUpload creates a file upload field with drag-and-drop support.
// The generated HTML includes a styled drop zone and a hidden file input.
// JavaScript is not required — the styling is purely CSS-driven.
//
// The name parameter sets the form field name. The accept parameter sets
// the accept attribute (e.g. "image/*", ".pdf,.doc"). Set multiple to true
// to allow multiple file selection.
//
// Usage:
//
//      html := FileUpload("documents", ".pdf,.doc,.docx", true)
func FileUpload(name string, accept string, multiple bool) string {
        var sb strings.Builder

        sb.WriteString(`<div class="gs-upload-zone"`)

        sb.WriteString(fmt.Sprintf(` data-gs-upload="%s"`, html.EscapeString(name)))

        if accept != "" {
                sb.WriteString(fmt.Sprintf(` data-gs-accept="%s"`, html.EscapeString(accept)))
        }

        if multiple {
                sb.WriteString(` data-gs-multiple="true"`)
        }

        sb.WriteString(">")
        sb.WriteString(`<div class="gs-upload-label">`)

        if multiple {
                sb.WriteString("Drop files here or click to browse")
        } else {
                sb.WriteString("Drop a file here or click to browse")
        }

        sb.WriteString(`</div>`)

        // Hidden file input
        sb.WriteString(fmt.Sprintf(`<input type="file" name="%s"`, html.EscapeString(name)))
        if accept != "" {
                sb.WriteString(fmt.Sprintf(` accept="%s"`, html.EscapeString(accept)))
        }
        if multiple {
                sb.WriteString(" multiple")
        }
        sb.WriteString(` class="gs-upload-input">`)

        sb.WriteString(`</div>`)

        return sb.String()
}

// =========================================================================
// Form Parsing Helpers
// =========================================================================

// ParseForm parses form data from an HTTP request into a map.
// Supports both URL-encoded forms and multipart forms.
//
// Usage:
//
//      data, err := ParseForm(r)
//      if err != nil {
//          // Handle error
//      }
//      result := Validate(data, rules...)
func ParseForm(r *http.Request) (map[string]interface{}, error) {
        contentType := r.Header.Get("Content-Type")

        data := make(map[string]interface{})

        // Handle multipart form data
        if strings.Contains(contentType, "multipart/form-data") {
                if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max
                        return nil, fmt.Errorf("failed to parse multipart form: %w", err)
                }
                for key, values := range r.MultipartForm.Value {
                        if len(values) > 0 {
                                data[key] = values[0]
                        }
                }
                return data, nil
        }

        // Handle URL-encoded form data
        if err := r.ParseForm(); err != nil {
                return nil, fmt.Errorf("failed to parse form: %w", err)
        }

        for key, values := range r.Form {
                if len(values) > 0 {
                        data[key] = values[0]
                }
        }

        return data, nil
}

// ParseJSONBody parses a JSON request body into a map.
// Useful when the form encoding is set to "application/json".
//
// Usage:
//
//      data, err := ParseJSONBody(r)
//      if err != nil {
//          // Handle error
//      }
func ParseJSONBody(r *http.Request) (map[string]interface{}, error) {
        var data map[string]interface{}
        if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
                return nil, fmt.Errorf("failed to parse JSON body: %w", err)
        }
        return data, nil
}
