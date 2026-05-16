//go:build !windows
// +build !windows

package hotkey

func SetupHotkey(toggle KeyBinding, exit KeyBinding, handler func(int)) {}

func (hm *HotkeyManager) ListenHotkeys() {}

func NewUnixHotkeyManager(toggle, exit func()) *HotkeyManager {
	return &HotkeyManager{
		ToggleHandler: toggle,
		ExitHandler:   exit,
	}
}
