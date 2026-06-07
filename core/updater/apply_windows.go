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

	// Prefer the system PowerShell executable by full path (under %SystemRoot%\System32),
	// and ensure the child process gets a safe PATH containing only system directories
	systemRoot := os.Getenv("SystemRoot")
	powPath := filepath.Join(systemRoot, "System32", "WindowsPowerShell", "v1.0", "powershell.exe")
	powExe := "powershell"
	if _, err := os.Stat(powPath); err == nil {
		powExe = powPath
	}

	cmd := exec.Command(
		powExe,
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

	// Restrict PATH to System32 and PowerShell folder only (these are non-writable by standard users)
	if systemRoot == "" {
		systemRoot = "C:\\Windows"
	}
	safePath := filepath.Join(systemRoot, "System32") + ";" + filepath.Join(systemRoot, "System32", "WindowsPowerShell", "v1.0")
	// Preserve other minimal env variables if helpful
	cmd.Env = append([]string{
		"PATH=" + safePath,
		"SystemRoot=" + systemRoot,
	}, os.Environ()...)
	// Start the process from temp dir
	cmd.Dir = os.TempDir()

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
