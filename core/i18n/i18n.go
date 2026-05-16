package i18n

import (
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
)

var (
	mu          sync.RWMutex
	currentLang = DetectLanguage()
)

var translations = map[Language]map[string]string{
	English: {
		SearchPlaceholder: "Search application...",
		TrayTooltip:       "GoFinder",
		TrayToggleTitle:   "Show/Hide",
		TrayToggleTooltip: "Show or hide GoFinder",
		TrayMinimizeTitle: "Minimize",
		TrayMinimizeTip:   "Hide GoFinder in the system tray",
		TrayQuitTitle:     "Quit",
		TrayQuitTooltip:   "Close GoFinder",
		AppExitMessage:    "Exiting...",
		LogRunningApp:     "Running: %s (%s)",
		LogRunAppError:    "Error running %s: %v",
	},
	Spanish: {
		SearchPlaceholder: "Buscar aplicación...",
		TrayTooltip:       "GoFinder",
		TrayToggleTitle:   "Mostrar/Ocultar",
		TrayToggleTooltip: "Mostrar u ocultar GoFinder",
		TrayMinimizeTitle: "Minimizar",
		TrayMinimizeTip:   "Ocultar GoFinder en la bandeja",
		TrayQuitTitle:     "Salir",
		TrayQuitTooltip:   "Cerrar GoFinder",
		AppExitMessage:    "Saliendo...",
		LogRunningApp:     "Ejecutando: %s (%s)",
		LogRunAppError:    "Error al ejecutar %s: %v",
	},
	Catalan: {
		SearchPlaceholder: "Cerca una aplicació...",
		TrayTooltip:       "GoFinder",
		TrayToggleTitle:   "Mostra/Amaga",
		TrayToggleTooltip: "Mostra o amaga GoFinder",
		TrayMinimizeTitle: "Minimitza",
		TrayMinimizeTip:   "Amaga GoFinder a la safata del sistema",
		TrayQuitTitle:     "Surt",
		TrayQuitTooltip:   "Tanca GoFinder",
		AppExitMessage:    "Sortint...",
		LogRunningApp:     "Executant: %s (%s)",
		LogRunAppError:    "Error en executar %s: %v",
	},
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
