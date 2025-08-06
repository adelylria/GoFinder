package logic

import (
	"os/exec"
	"runtime"
	"syscall"

	"github.com/adelylria/GoFinder/models"
)

func RunApplication(app models.Application) error {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/C", "start", "", app.Exec)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: false}
		return cmd.Start()
	}

	cmd := exec.Command(app.Exec)
	return cmd.Start()
}
