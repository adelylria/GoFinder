//go:build windows

// run_windows.go
package logic

import (
	"path/filepath"

	"github.com/adelylria/GoFinder/models"
	"golang.org/x/sys/windows"
)

func RunApplication(app models.Application) error {
	if app.Exec == "" {
		return nil
	}

	verb, err := windows.UTF16PtrFromString("open")
	if err != nil {
		return err
	}
	file, err := windows.UTF16PtrFromString(app.Exec)
	if err != nil {
		return err
	}

	cwd, err := windows.UTF16PtrFromString(filepath.Dir(app.Exec))
	if err != nil {
		return err
	}

	const showNormal = 1
	return windows.ShellExecute(0, verb, file, nil, cwd, showNormal)
}
