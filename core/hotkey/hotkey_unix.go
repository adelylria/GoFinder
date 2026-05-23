//go:build !windows
// +build !windows

package hotkey

func SetupHotkey(toggle KeyBinding, exit KeyBinding, handler func(int)) {
	// Hotkey setup is not implemented for Unix-like systems in this version.
}

func (hm *HotkeyManager) ListenHotkeys() {
	// Hotkey listening is not implemented for Unix-like systems in this version.
}

func NewUnixHotkeyManager(toggle, exit func()) *HotkeyManager {
	return &HotkeyManager{
		ToggleHandler: toggle,
		ExitHandler:   exit,
	}
}
