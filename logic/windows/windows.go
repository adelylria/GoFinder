//go:build windows
// +build windows

package windows

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adelylria/GoFinder/core/global"
	"github.com/adelylria/GoFinder/logic/common"
	"github.com/adelylria/GoFinder/models"
)

type WindowsAppFinder struct{}

func (f WindowsAppFinder) Find() []models.Application {
	fmt.Println("Buscando aplicaciones en Windows...")
	return findWindowsApplications()
}

func ProcessWindowsShortcut(path string, seen map[string]bool) *models.Application {
	app := resolveWindowsShortcut(path)

	if !common.IsValidApp(app) {
		fmt.Printf("Descartada (no válida): %s -> %s\n", app.Name, app.Exec)
		return nil
	}

	baseExec := filepath.Base(app.Exec)
	if global.ExcludedApps[baseExec] || global.ExcludedApps[app.Name] ||
		strings.Contains(strings.ToLower(app.Name), "uninstall") ||
		strings.Contains(strings.ToLower(app.Name), "settings") {
		fmt.Printf("Descartada (excluida): %s -> %s\n", app.Name, app.Exec)
		return nil
	}

	if seen[app.Exec] {
		fmt.Printf("Descartada (duplicada): %s -> %s\n", app.Name, app.Exec)
		return nil
	}
	seen[app.Exec] = true

	iconPath, iconIndex := common.ParseIconLocation(app.Icon)
	app.IconPath = iconPath
	app.IconIdx = iconIndex

	if app.IconPath != "" && !filepath.IsAbs(app.IconPath) {
		app.IconPath = filepath.Join(filepath.Dir(path), app.IconPath)
	}

	if app.IconPath == "" {
		app.IconPath = app.Exec
	}

	return &app
}

func findWindowsApplications() []models.Application {
	var apps []models.Application
	seen := make(map[string]bool)
	desktopDir := filepath.Join(os.Getenv("USERPROFILE"), "Desktop")

	for _, dir := range common.GetAppDirs() {
		fmt.Println("Escaneando directorio:", dir)
		addShortcutsFromDir(&apps, dir, desktopDir, seen)
	}
	return apps
}

func addShortcutsFromDir(apps *[]models.Application, dir string, desktopDir string, seen map[string]bool) {
	absDir, _ := filepath.Abs(dir)
	absDesktop, _ := filepath.Abs(desktopDir)

	if absDir == absDesktop {
		processDesktopDir(apps, dir, seen)
		return
	}

	if err := walkAndAddShortcuts(apps, dir, seen); err != nil {
		fmt.Printf("Error recorriendo %s: %v\n", dir, err)
	}
}

func processDesktopDir(apps *[]models.Application, dir string, seen map[string]bool) {
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("Error leyendo escritorio: %v\n", err)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := strings.ToLower(file.Name())
		if !strings.HasSuffix(name, ".lnk") {
			continue
		}
		addIfShortcutPath(apps, filepath.Join(dir, file.Name()), seen)
	}
}

func walkAndAddShortcuts(apps *[]models.Application, dir string, seen map[string]bool) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accediendo %s: %v\n", path, err)
			return nil
		}
		if info.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".lnk") {
			return nil
		}
		addIfShortcutPath(apps, path, seen)
		return nil
	})
}

func addIfShortcutPath(apps *[]models.Application, path string, seen map[string]bool) {
	if app := ProcessWindowsShortcut(path, seen); app != nil {
		*apps = append(*apps, *app)
		fmt.Printf("Añadida: %s -> %s\n", app.Name, app.Exec)
	}
}
