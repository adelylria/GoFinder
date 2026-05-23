//go:build windows

// run_windows.go
package logic

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/adelylria/GoFinder/models"
)

func RunApplication(app models.Application) error {
	// Build a conservative PATH containing only system directories
	systemRoot := os.Getenv("SystemRoot")
	if systemRoot == "" {
		systemRoot = `C:\Windows`
	}

	safeDirs := []string{
		filepath.Join(systemRoot, "System32"),
		systemRoot,
	}

	// Add extras if present
	extras := []string{
		filepath.Join(systemRoot, "System32", "Wbem"),
		filepath.Join(systemRoot, "System32", "WindowsPowerShell", "v1.0"),
	}
	for _, d := range extras {
		if fi, err := os.Stat(d); err == nil && fi.IsDir() {
			safeDirs = append(safeDirs, d)
		}
	}

	// Prepare environment: keep existing vars except PATH, then set a safe PATH
	env := os.Environ()
	newEnv := make([]string, 0, len(env)+1)
	for _, e := range env {
		if strings.HasPrefix(strings.ToUpper(e), "PATH=") {
			continue
		}
		newEnv = append(newEnv, e)
	}
	newEnv = append(newEnv, "PATH="+strings.Join(safeDirs, string(os.PathListSeparator)))

	cmd := exec.Command("cmd", "/C", "start", "", app.Exec)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: false}
	cmd.Env = newEnv
	return cmd.Start()
}
