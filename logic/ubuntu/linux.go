package ubuntu

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adelylria/GoFinder/logic/common"
	"github.com/adelylria/GoFinder/models"
)

type LinuxAppFinder struct{}

func (f LinuxAppFinder) Find() []models.Application {
	fmt.Println("Buscando aplicaciones en Linux...")
	return findLinuxApplications()
}

func findLinuxApplications() []models.Application {
	var apps []models.Application

	for _, dir := range common.GetAppDirs() {
		fmt.Println("Escaneando directorio:", dir)
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || !strings.HasSuffix(path, ".desktop") {
				return nil
			}
			fmt.Println("Encontrado .desktop:", path)

			if app, ok := parseDesktopFile(path); ok {
				fmt.Printf("  → %s -> %s\n", app.Name, app.Exec)
				apps = append(apps, app)
			}
			return nil
		})
	}
	return apps
}

func parseDesktopFile(path string) (models.Application, bool) {
	app := models.NewApplication()
	data, err := os.ReadFile(path)
	if err != nil {
		return app, false
	}

	inDesktopEntry := false
	lines := strings.SplitSeq(string(data), "\n")

	for raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if ok, name := parseSectionHeader(line); ok {
			inDesktopEntry = (name == "Desktop Entry")
			continue
		}

		if !inDesktopEntry {
			continue
		}

		if key, value, ok := splitKeyValue(line); ok {
			applyDesktopKey(&app, key, value)
		}
	}
	return app, common.IsValidApp(app)
}

func parseSectionHeader(line string) (bool, string) {
	if len(line) > 2 && line[0] == '[' && line[len(line)-1] == ']' {
		name := strings.TrimSpace(line[1 : len(line)-1])
		return true, name
	}
	return false, ""
}

func splitKeyValue(line string) (string, string, bool) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	return key, value, true
}

func applyDesktopKey(app *models.Application, key, value string) {
	switch key {
	case "Name":
		if app.Name == "" {
			app.Name = value
		}
	case "Exec":
		if app.Exec == "" {
			app.Exec = value
		}
	case "Icon":
		if app.Icon == "" {
			app.Icon = value
		}
	}
}
