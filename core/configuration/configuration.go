package configuration

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const appName = "GoFinder"

type KeyBinding struct {
	Modifier string `json:"modifier"`
	Key      string `json:"key"`
}

type Config struct {
	ToggleHotkey KeyBinding `json:"toggle_hotkey"`
	QuitHotkey   KeyBinding `json:"quit_hotkey"`
	AutoStart    bool       `json:"auto_start"`
	StartHidden  bool       `json:"start_hidden"`
	ThemeName    string     `json:"theme_name"`
}

func DefaultConfig() Config {
	return Config{
		ToggleHotkey: KeyBinding{Modifier: "Alt", Key: "R"},
		QuitHotkey:   KeyBinding{Modifier: "Alt", Key: "Q"},
		AutoStart:    false,
		StartHidden:  false,
		ThemeName:    "system",
	}
}

func Load() (Config, error) {
	cfg := DefaultConfig()
	path, err := configPath()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	cfg.Normalize()
	return cfg, nil
}

func Save(cfg Config) error {
	cfg.Normalize()
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (c *Config) Normalize() {
	defaults := DefaultConfig()
	c.ToggleHotkey = normalizeKeyBinding(c.ToggleHotkey, defaults.ToggleHotkey)
	c.QuitHotkey = normalizeKeyBinding(c.QuitHotkey, defaults.QuitHotkey)
	c.ThemeName = normalizeThemeName(c.ThemeName, defaults.ThemeName)
}

func normalizeKeyBinding(binding, fallback KeyBinding) KeyBinding {
	modifier := normalizeModifier(binding.Modifier)
	key := normalizeKey(binding.Key)
	if modifier == "" {
		modifier = fallback.Modifier
	}
	if key == "" {
		key = fallback.Key
	}
	return KeyBinding{Modifier: modifier, Key: key}
}

func normalizeModifier(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case "alt":
		return "Alt"
	case "ctrl", "control":
		return "Ctrl"
	case "shift":
		return "Shift"
	case "ctrl+alt", "alt+ctrl":
		return "Ctrl+Alt"
	default:
		return ""
	}
}

func normalizeKey(value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	if len(value) != 1 {
		return ""
	}
	if value[0] < 'A' || value[0] > 'Z' {
		return ""
	}
	return value
}

func normalizeThemeName(value, fallback string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case "system", "light", "dark", "ocean", "forest", "contrast", "paper", "pastel", "solar", "midnight":
		return normalized
	default:
		return fallback
	}
}

func configPath() (string, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cfgDir, appName, "config.json"), nil
}
