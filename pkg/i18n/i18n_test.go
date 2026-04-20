package i18n

import "testing"

func TestLocalizerFallbackAndInterpolation(t *testing.T) {
	bundle := NewBundle("en")
	bundle.AddCatalog(Catalog{
		Locale: "en",
		Messages: map[string]Message{
			"welcome": {ID: "welcome", Text: "Welcome, {{name}}"},
		},
	})
	bundle.AddCatalog(Catalog{
		Locale: "ar",
		Messages: map[string]Message{
			"welcome": {ID: "welcome", Text: "مرحبا، {{name}}"},
		},
	})

	localizer := NewLocalizer(bundle, "ar")
	got := localizer.T("welcome", map[string]string{"name": "Asha"})
	if got != "مرحبا، Asha" {
		t.Fatalf("unexpected localized message: %q", got)
	}

	fallback := NewLocalizer(bundle, "fr")
	got = fallback.T("welcome", map[string]string{"name": "Asha"})
	if got != "Welcome, Asha" {
		t.Fatalf("unexpected fallback message: %q", got)
	}
}

func TestDirectionForLocale(t *testing.T) {
	if DirectionForLocale("ar-SA") != DirectionRTL {
		t.Fatalf("expected ar-SA to be rtl")
	}
	if DirectionForLocale("en-US") != DirectionLTR {
		t.Fatalf("expected en-US to be ltr")
	}
}
