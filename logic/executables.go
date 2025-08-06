package logic

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"

	"github.com/adelylria/GoFinder/models"
)

// FindApplications detecta programas visibles instalados (menú inicio en Windows o .desktop en Linux)
func FindApplications() []models.Application {
	if runtime.GOOS == "windows" {
		fmt.Println("Buscando aplicaciones en Windows...")
		return findWindowsApplications()
	}
	fmt.Println("Buscando aplicaciones en Linux...")
	return findLinuxApplications()
}

// ----------------------------------------
// WINDOWS
// ----------------------------------------

func findWindowsApplications() []models.Application {
	dirs := []string{
		os.Getenv("APPDATA") + `\Microsoft\Windows\Start Menu\Programs`,
		os.Getenv("ProgramData") + `\Microsoft\Windows\Start Menu\Programs`,
	}

	var apps []models.Application

	for _, dir := range dirs {
		fmt.Println("Escaneando directorio:", dir)

		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".lnk") {
				return nil
			}
			fmt.Println("Encontrado acceso directo:", path)

			app := resolveShortcut(path)
			if app.Exec != "" {
				iconPath, iconIndex := SplitIconLocation(app.Icon)
				iconPath = os.ExpandEnv(iconPath)

				app.IconPath = iconPath
				app.IconIdx = iconIndex

				// Solo añadir si tiene icono válido
				if LoadAppIcon(app) != nil {
					apps = append(apps, app)
				} else {
					fmt.Println("Descartada app sin icono válido:", app.Name)
				}
			}
			return nil
		})
	}
	return apps
}

func resolveShortcut(path string) models.Application {
	app := models.NewApplication()
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	shell, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		fmt.Println("Error creando WScript.Shell:", err)
		return app
	}
	defer shell.Release()

	wshell, err := shell.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		fmt.Println("Error obteniendo IDispatch:", err)
		return app
	}
	defer wshell.Release()

	link, err := oleutil.CallMethod(wshell, "CreateShortcut", path)
	if err != nil {
		fmt.Println("Error creando Shortcut COM object:", err)
		return app
	}
	defer link.Clear()
	dispatch := link.ToIDispatch()

	target, err1 := oleutil.GetProperty(dispatch, "TargetPath")
	iconLoc, err2 := oleutil.GetProperty(dispatch, "IconLocation")

	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	app.Name = name

	if err1 == nil && target.VT != ole.VT_NULL {
		app.Exec = target.ToString()
	}

	if err2 == nil && iconLoc.VT != ole.VT_NULL {
		app.Icon = iconLoc.ToString()
	}

	return app
}

// ----------------------------------------
// LINUX / UNIX
// ----------------------------------------

func findLinuxApplications() []models.Application {
	paths := []string{
		"/usr/share/applications",
		filepath.Join(os.Getenv("HOME"), ".local/share/applications"),
	}
	var apps []models.Application

	for _, dir := range paths {
		fmt.Println("Escaneando directorio:", dir)

		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || !strings.HasSuffix(path, ".desktop") {
				return nil
			}
			fmt.Println("Encontrado .desktop:", path)

			app := parseDesktopFile(path)
			if app.Name != "" && app.Exec != "" {
				fmt.Printf("  → %s -> %s\n", app.Name, app.Exec)
				apps = append(apps, app)
			}
			return nil
		})
	}
	return apps
}

func parseDesktopFile(path string) models.Application {
	app := models.NewApplication()
	data, err := os.ReadFile(path)
	if err != nil {
		return app
	}
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "Name=") && app.Name == "" {
			app.Name = strings.TrimPrefix(line, "Name=")
		} else if strings.HasPrefix(line, "Exec=") && app.Exec == "" {
			exec := strings.TrimPrefix(line, "Exec=")
			// Quitar parámetros y variables (%F, %U, etc.)
			app.Exec = strings.Split(exec, " ")[0]
		} else if strings.HasPrefix(line, "Icon=") && app.Icon == "" {
			app.Icon = strings.TrimPrefix(line, "Icon=")
		}
	}

	return app
}
