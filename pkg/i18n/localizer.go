package i18n

import "strings"

// Localizer resolves messages for one locale with fallbacks.
type Localizer struct {
	bundle    *Bundle
	locale    string
	fallbacks []string
}

// NewLocalizer creates a locale resolver.
func NewLocalizer(bundle *Bundle, locale string, fallbacks ...string) *Localizer {
	normalizedFallbacks := make([]string, 0, len(fallbacks))
	for _, fallback := range fallbacks {
		normalized := normalizeLocale(fallback)
		if normalized != "" {
			normalizedFallbacks = append(normalizedFallbacks, normalized)
		}
	}

	return &Localizer{
		bundle:    bundle,
		locale:    normalizeLocale(locale),
		fallbacks: normalizedFallbacks,
	}
}

// Locale returns the currently selected locale.
func (l *Localizer) Locale() string {
	return l.locale
}

// Direction returns the layout direction for the active locale.
func (l *Localizer) Direction() Direction {
	return DirectionForLocale(l.locale)
}

// T resolves a message by ID and interpolates variables of the form {{name}}.
func (l *Localizer) T(id string, vars map[string]string) string {
	if l == nil || l.bundle == nil {
		return id
	}

	for _, locale := range l.lookupLocales() {
		catalog, ok := l.bundle.Catalogs[locale]
		if !ok {
			continue
		}

		message, ok := catalog.Messages[id]
		if !ok {
			continue
		}

		return interpolate(message.Text, vars)
	}

	return id
}

func (l *Localizer) lookupLocales() []string {
	locales := make([]string, 0, 2+len(l.fallbacks))
	if l.locale != "" {
		locales = append(locales, l.locale)
	}
	locales = append(locales, l.fallbacks...)

	if l.bundle != nil && l.bundle.DefaultLocale != "" {
		locales = append(locales, l.bundle.DefaultLocale)
	}

	seen := map[string]bool{}
	filtered := make([]string, 0, len(locales))
	for _, locale := range locales {
		if locale == "" || seen[locale] {
			continue
		}
		seen[locale] = true
		filtered = append(filtered, locale)
	}

	return filtered
}

func interpolate(text string, vars map[string]string) string {
	if len(vars) == 0 {
		return text
	}

	result := text
	for key, value := range vars {
		token := "{{" + strings.TrimSpace(key) + "}}"
		result = strings.ReplaceAll(result, token, value)
	}

	return result
}
