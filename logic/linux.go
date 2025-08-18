package logic

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adelylria/GoFinder/models"
)

type linuxAppFinder struct{}

func (f linuxAppFinder) Find() []models.Application {
	fmt.Println("Buscando aplicaciones en Linux...")
	return findLinuxApplications()
}

func findLinuxApplications() []models.Application {
	var apps []models.Application

	for _, dir := range GetAppDirs() {
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

	for line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}

		if len(line) > 2 && line[0] == '[' && line[len(line)-1] == ']' {
			inDesktopEntry = (line == "[Desktop Entry]")
			continue
		}

		if !inDesktopEntry {
			continue
		}

		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

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
	}
	return app, IsValidApp(app)
}
