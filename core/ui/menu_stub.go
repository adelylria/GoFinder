//go:build !windows

package ui

func (l *Launcher) applyNativeMenuPlatformHooks() {
	// No platform-specific menu hooks needed for non-Windows platforms in this version.
}
