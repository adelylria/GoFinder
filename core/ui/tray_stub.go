//go:build !windows

package ui

import "github.com/adelylria/GoFinder/models"

func startSystemTray(state *models.AppState, icon []byte) {
	// System tray is not implemented for non-Windows platforms in this version.
}
