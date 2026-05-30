//go:build !windows

package ui

func startSystemTray(state *AppState, icon []byte) {
	// System tray is not implemented for non-Windows platforms in this version.
}
