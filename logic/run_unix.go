//go:build linux || darwin

// run_unix.go
package logic

import (
	"os/exec"

	"github.com/adelylria/GoFinder/models"
)

func RunApplication(app models.Application) error {
	if app.Exec == "" {
		return nil
	}
	cmd := exec.Command(app.Exec)
	return cmd.Start()
}
