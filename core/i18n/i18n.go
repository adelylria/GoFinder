package i18n

import (
	"embed"
	"encoding/json"
	"os"
	"strings"
	"sync"
)

type Language string

const (
	English Language = "en"
	Spanish Language = "es"
	Catalan Language = "ca"
)

const (
	SearchPlaceholder = "search.placeholder"
	TrayTooltip       = "tray.tooltip"
	TrayToggleTitle   = "tray.toggle.title"
	TrayToggleTooltip = "tray.toggle.tooltip"
	TrayMinimizeTitle = "tray.minimize.title"
	TrayMinimizeTip   = "tray.minimize.tooltip"
	TrayQuitTitle     = "tray.quit.title"
	TrayQuitTooltip   = "tray.quit.tooltip"
	AppExitMessage    = "app.exit.message"
	LogRunningApp     = "log.running_app"
	LogRunAppError    = "log.run_app_error"
	SettingsToggle    = "settings.toggle"
	SettingsQuit      = "settings.quit"
	SettingsAutoStart = "settings.autostart"
	SettingsHidden    = "settings.hidden"
	SettingsSaved     = "settings.saved"
	SettingsRestart   = "settings.restart"
	SettingsHotkeys   = "settings.hotkeys"
	SettingsGeneral   = "settings.general"
	MenuFile          = "menu.file"
	MenuConfig        = "menu.config"
	MenuHelp          = "menu.help"
	MenuExit          = "menu.exit"
	MenuPreferences   = "menu.preferences"
	MenuAbout         = "menu.about"
	DialogClose       = "dialog.close"
	AboutText         = "about.text"
)

var (
	mu          sync.RWMutex
	currentLang = DetectLanguage()
)

//go:embed locales/en.json locales/es.json locales/ca.json
var localeFS embed.FS

var translations = make(map[Language]map[string]string)

func init() {
	// load embedded JSON locales
	files := map[Language]string{
		English: "locales/en.json",
		Spanish: "locales/es.json",
		Catalan: "locales/ca.json",
	}
	for lang, file := range files {
		data, err := localeFS.ReadFile(file)
		if err != nil {
			// if embed read fails, keep translations empty for that lang
			continue
		}
		var m map[string]string
		if err := json.Unmarshal(data, &m); err != nil {
			continue
		}
		translations[lang] = m
	}
	// ensure English exists as a fallback
	if _, ok := translations[English]; !ok {
		translations[English] = map[string]string{}
	}
}

func DetectLanguage() Language {
	if override := os.Getenv("GOFINDER_LANG"); override != "" {
		return NormalizeLanguage(override)
	}
	return NormalizeLanguage(systemLocale())
}

func NormalizeLanguage(value string) Language {
	normalized := strings.ToLower(strings.TrimSpace(value))
	normalized = strings.ReplaceAll(normalized, "_", "-")

	switch {
	case strings.HasPrefix(normalized, "ca"):
		return Catalan
	case strings.HasPrefix(normalized, "es"):
		return Spanish
	case strings.HasPrefix(normalized, "en"):
		return English
	default:
		return English
	}
}

func CurrentLanguage() Language {
	mu.RLock()
	defer mu.RUnlock()
	return currentLang
}

func SetLanguage(language Language) {
	mu.Lock()
	currentLang = NormalizeLanguage(string(language))
	mu.Unlock()
}

func T(key string) string {
	mu.RLock()
	language := currentLang
	mu.RUnlock()

	if value, ok := translations[language][key]; ok {
		return value
	}
	if value, ok := translations[English][key]; ok {
		return value
	}
	return key
}
