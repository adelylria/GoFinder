//go:build !windows

package updater

import "fmt"

func ApplyDownloadedUpdate(downloadPath string) error {
	return fmt.Errorf("auto-update is only available on Windows")
}
