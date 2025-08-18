package logic

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/adelylria/GoFinder/models"
)

func GetAppDirs() []string {
	if runtime.GOOS == "windows" {
		return []string{
			filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs"),
			filepath.Join(os.Getenv("ProgramData"), "Microsoft", "Windows", "Start Menu", "Programs"),
			filepath.Join(os.Getenv("USERPROFILE"), "Desktop"),
		}
	}
	return []string{
		"/usr/share/applications",
		filepath.Join(os.Getenv("HOME"), ".local/share/applications"),
	}
}

// isValidApp verifica si una aplicación tiene los campos mínimos requeridos
func IsValidApp(app models.Application) bool {
	if runtime.GOOS == "windows" {
		return app.Name != "" && app.Exec != "" && strings.EqualFold(filepath.Ext(app.Exec), ".exe")
	}

	return app.Name != "" && app.Exec != ""
}

func parseIconLocation(iconLoc string) (string, int) {
	if iconLoc == "" {
		return "", 0
	}

	parts := strings.Split(iconLoc, ",")
	if len(parts) < 2 {
		return os.ExpandEnv(strings.Trim(parts[0], `"`)), 0
	}

	pathPart := strings.Join(parts[:len(parts)-1], ",")
	indexPart := parts[len(parts)-1]

	pathPart = strings.Trim(pathPart, `"`)

	idx, err := strconv.Atoi(strings.TrimSpace(indexPart))
	if err != nil {
		return os.ExpandEnv(pathPart), 0
	}

	return os.ExpandEnv(pathPart), idx
}

func SplitIconLocation(iconLoc string) (string, int) {
	if iconLoc == "" {
		return "", 0
	}

	if strings.HasPrefix(iconLoc, `"`) {
		endQuote := strings.Index(iconLoc[1:], `"`)
		if endQuote != -1 {
			pathPart := iconLoc[1 : endQuote+1]
			remaining := strings.TrimSpace(iconLoc[endQuote+2:])

			if strings.HasPrefix(remaining, ",") {
				indexPart := strings.TrimSpace(remaining[1:])
				idx, _ := strconv.Atoi(indexPart)
				return pathPart, idx
			}
			return pathPart, 0
		}
	}

	lastComma := strings.LastIndex(iconLoc, ",")
	if lastComma == -1 {
		return iconLoc, 0
	}

	pathPart := strings.TrimSpace(iconLoc[:lastComma])
	indexPart := strings.TrimSpace(iconLoc[lastComma+1:])

	idx, err := strconv.Atoi(indexPart)
	if err != nil {
		return pathPart, 0
	}

	return pathPart, idx
}
