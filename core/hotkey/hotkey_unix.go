//go:build !windows
// +build !windows

package hotkey

// En Linux/macOS no hacemos nada. Solo stubs.

func SetupHotkey(handler func(int))      {}
func (hm *HotkeyManager) ListenHotkeys() {}

func NewUnixHotkeyManager(toggle, exit func()) *HotkeyManager {
	return &HotkeyManager{
		ToggleHandler: toggle,
		ExitHandler:   exit,
	}
}
