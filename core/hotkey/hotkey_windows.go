//go:build windows
// +build windows

package hotkey

/*
#cgo windows CFLAGS: -I./hotkey_windows
#cgo windows LDFLAGS: -luser32
#include "hotkey_windows/hotkey.h"
*/
import "C"
import (
	"fmt"

	"fyne.io/fyne/v2"
)

var hotkeyHandler func(int)

//export handleHotkey
func handleHotkey(id C.int) {
	if hotkeyHandler != nil {
		hotkeyHandler(int(id))
	}
}

func SetupHotkey(toggle KeyBinding, exit KeyBinding, handler func(int)) {
	hotkeyHandler = handler
	toggle = normalizeHotkeyBinding(toggle, KeyBinding{Modifier: "Alt", Key: "R"})
	exit = normalizeHotkeyBinding(exit, KeyBinding{Modifier: "Alt", Key: "Q"})
	C.setupHotkeys(
		C.uint(hotkeyModifier(toggle.Modifier)),
		C.uint(hotkeyKey(toggle.Key)),
		C.uint(hotkeyModifier(exit.Modifier)),
		C.uint(hotkeyKey(exit.Key)),
	)
}

func (hm *HotkeyManager) ListenHotkeys() {
	go func() {
		exitHandler := func() {
			if hm.ExitHandler != nil {
				fyne.Do(hm.ExitHandler)
			}
		}

		handlers := map[int]func(){
			1: func() {
				if hm.ToggleHandler != nil {
					fyne.Do(hm.ToggleHandler)
				}
			},
			2: exitHandler,
			3: exitHandler,
			4: func() {
				if hm.PrefsHandler != nil {
					fyne.Do(hm.PrefsHandler)
				}
			},
			5: func() {
				if hm.AboutHandler != nil {
					fyne.Do(hm.AboutHandler)
				}
			},
		}

		SetupHotkey(hm.ToggleHotkey, hm.ExitHotkey, func(id int) {
			if h, ok := handlers[id]; ok {
				h()
				return
			}
			fmt.Printf("unknown hotkey id: %d\n", id)
		})
	}()
}

func normalizeHotkeyBinding(binding, fallback KeyBinding) KeyBinding {
	if hotkeyModifier(binding.Modifier) == 0 {
		binding.Modifier = fallback.Modifier
	}
	if hotkeyKey(binding.Key) == 0 {
		binding.Key = fallback.Key
	}
	return binding
}

func hotkeyModifier(value string) uint32 {
	switch value {
	case "Ctrl":
		return 0x0002
	case "Shift":
		return 0x0004
	case "Ctrl+Alt":
		return 0x0002 | 0x0001
	case "Alt":
		return 0x0001
	default:
		return 0
	}
}

func hotkeyKey(value string) uint32 {
	if len(value) != 1 {
		return 0
	}
	key := value[0]
	if key >= 'a' && key <= 'z' {
		key -= 'a' - 'A'
	}
	if key < 'A' || key > 'Z' {
		return 0
	}
	return uint32(key)
}
