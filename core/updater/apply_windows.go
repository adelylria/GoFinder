//go:build windows

package updater

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func ApplyDownloadedUpdate(downloadPath string) error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	scriptPath := filepath.Join(os.TempDir(), "gofinder-apply-update.ps1")
	if err := os.WriteFile(scriptPath, []byte(updateScript()), 0o600); err != nil {
		return err
	}

	cmd := exec.Command(
		"powershell",
		"-NoProfile",
		"-ExecutionPolicy",
		"Bypass",
		"-WindowStyle",
		"Hidden",
		"-File",
		scriptPath,
		downloadPath,
		exePath,
		strconv.Itoa(os.Getpid()),
	)
	return cmd.Start()
}

func updateScript() string {
	lines := []string{
		"param([string]$Source, [string]$Target, [int]$Pid)",
		"$ErrorActionPreference = 'Stop'",
		"Wait-Process -Id $Pid -ErrorAction SilentlyContinue",
		"Start-Sleep -Milliseconds 500",
		"Copy-Item -LiteralPath $Source -Destination $Target -Force",
		"Start-Process -FilePath $Target",
		"Remove-Item -LiteralPath $Source -Force -ErrorAction SilentlyContinue",
		"Remove-Item -LiteralPath $MyInvocation.MyCommand.Path -Force -ErrorAction SilentlyContinue",
	}
	return strings.Join(lines, "\r\n") + "\r\n"
}
