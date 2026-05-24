package configuration

import "testing"

func TestNormalizeConfig(t *testing.T) {
	cfg := Config{
		ToggleHotkey: KeyBinding{Modifier: "control", Key: "r"},
		QuitHotkey:   KeyBinding{Modifier: "wat", Key: "1"},
	}
	cfg.Normalize()

	if cfg.ToggleHotkey != (KeyBinding{Modifier: "Ctrl", Key: "R"}) {
		t.Fatalf("unexpected toggle hotkey: %#v", cfg.ToggleHotkey)
	}
	if cfg.QuitHotkey != DefaultConfig().QuitHotkey {
		t.Fatalf("unexpected quit fallback: %#v", cfg.QuitHotkey)
	}
	if cfg.ThemeName != DefaultConfig().ThemeName {
		t.Fatalf("unexpected theme fallback: %q", cfg.ThemeName)
	}
}
