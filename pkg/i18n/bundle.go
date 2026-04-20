package i18n

import (
	"sort"
	"strings"
)

// Direction describes layout direction for a locale.
type Direction string

const (
	DirectionLTR Direction = "ltr"
	DirectionRTL Direction = "rtl"
)

// Message is a localizable string entry.
type Message struct {
	ID          string `json:"id"`
	Text        string `json:"text"`
	Description string `json:"description,omitempty"`
}

// Catalog stores messages for a locale.
type Catalog struct {
	Locale   string             `json:"locale"`
	Messages map[string]Message `json:"messages"`
}

// Bundle stores all locale catalogs.
type Bundle struct {
	DefaultLocale string             `json:"defaultLocale"`
	Catalogs      map[string]Catalog `json:"catalogs"`
}

// NewBundle creates a new message bundle.
func NewBundle(defaultLocale string) *Bundle {
	if defaultLocale == "" {
		defaultLocale = "en"
	}

	return &Bundle{
		DefaultLocale: normalizeLocale(defaultLocale),
		Catalogs:      make(map[string]Catalog),
	}
}

// AddCatalog registers a locale catalog.
func (b *Bundle) AddCatalog(catalog Catalog) {
	if catalog.Locale == "" {
		return
	}
	if catalog.Messages == nil {
		catalog.Messages = map[string]Message{}
	}

	catalog.Locale = normalizeLocale(catalog.Locale)
	b.Catalogs[catalog.Locale] = catalog
}

// AvailableLocales returns the locales in sorted order.
func (b *Bundle) AvailableLocales() []string {
	locales := make([]string, 0, len(b.Catalogs))
	for locale := range b.Catalogs {
		locales = append(locales, locale)
	}
	sort.Strings(locales)
	return locales
}

// DirectionForLocale returns the writing direction for a locale.
func DirectionForLocale(locale string) Direction {
	locale = normalizeLocale(locale)
	prefix := locale
	if idx := strings.Index(prefix, "-"); idx >= 0 {
		prefix = prefix[:idx]
	}

	switch prefix {
	case "ar", "fa", "he", "ur":
		return DirectionRTL
	default:
		return DirectionLTR
	}
}

func normalizeLocale(locale string) string {
	locale = strings.TrimSpace(strings.ReplaceAll(locale, "_", "-"))
	if locale == "" {
		return ""
	}
	return strings.ToLower(locale)
}
