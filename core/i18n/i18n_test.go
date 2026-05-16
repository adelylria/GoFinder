package i18n

import "testing"

func TestNormalizeLanguage(t *testing.T) {
	tests := map[string]Language{
		"en-US": English,
		"es_ES": Spanish,
		"ca-ES": Catalan,
		"fr-FR": English,
		"":      English,
	}

	for value, expected := range tests {
		if got := NormalizeLanguage(value); got != expected {
			t.Fatalf("NormalizeLanguage(%q) = %q, want %q", value, got, expected)
		}
	}
}

func TestTranslationFallback(t *testing.T) {
	SetLanguage(Catalan)
	if got := T(TrayQuitTitle); got != "Surt" {
		t.Fatalf("T(%q) = %q", TrayQuitTitle, got)
	}

	if got := T("missing.key"); got != "missing.key" {
		t.Fatalf("missing key = %q", got)
	}
}
