//go:build windows
// +build windows

package hotkey

/*
#cgo CFLAGS: -I./hotkey_windows
#cgo LDFLAGS: -luser32
#include "hotkey.h"
*/
import "C"
import "fmt"

var hotkeyHandler func(int)

//export handleHotkey
func handleHotkey(id C.int) {
	if hotkeyHandler != nil {
		hotkeyHandler(int(id))
	}
}

func SetupHotkey(handler func(int)) {
	hotkeyHandler = handler
	C.setupHotkey()
}

func (hm *HotkeyManager) ListenHotkeys() {
	go func() {
		SetupHotkey(func(id int) {
			switch id {
			case 1:
				if hm.ToggleHandler != nil {
					hm.ToggleHandler()
				}
			case 2:
				if hm.ExitHandler != nil {
					hm.ExitHandler()
				}
			default:
				fmt.Printf("unknown hotkey id: %d\n", id)
			}
		})
	}()
}
