//go:build windows

package configuration

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const runKeyPath = `Software\Microsoft\Windows\CurrentVersion\Run`

func ApplyAutoStart(enabled bool) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, runKeyPath, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	if !enabled {
		if err := key.DeleteValue(appName); err != nil && err != registry.ErrNotExist {
			return err
		}
		return nil
	}

	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	return key.SetStringValue(appName, quoteWindowsArg(exePath))
}

func quoteWindowsArg(value string) string {
	return fmt.Sprintf(`"%s"`, strings.ReplaceAll(value, `"`, `\"`))
}
