//go:build windows

// run_windows.go
package logic

import (
	"os/exec"
	"syscall"

	"github.com/adelylria/GoFinder/models"
)

func RunApplication(app models.Application) error {
	cmd := exec.Command("cmd", "/C", "start", "", app.Exec)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: false}
	return cmd.Start()
}
