package hotkey

// #cgo LDFLAGS: -luser32
// #include "hotkey.h"
import "C"
import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type HotkeyManager struct {
	ToggleHandler func()
	ExitHandler   func()
}

// KeyEventInterceptor intercepta eventos de teclado para manejar la navegación
type KeyEventInterceptor struct {
	widget.Entry
	OnKeyDown func()
	OnKeyUp   func()
}

//export handleHotkey
func handleHotkey(id C.int) {
	if hotkeyHandler != nil {
		hotkeyHandler(int(id))
	}
}

var hotkeyHandler func(int)

func SetupHotkey(handler func(int)) {
	hotkeyHandler = handler
	C.setupHotkey()
}

func (hm *HotkeyManager) ListenHotkeys() {
	go func() {
		SetupHotkey(func(id int) {
			switch id {
			case 1: // Ctrl+Alt+B
				if hm.ToggleHandler != nil {
					hm.ToggleHandler()
				}
			case 2: // Ctrl+Alt+Q
				if hm.ExitHandler != nil {
					hm.ExitHandler()
				}
			}
		})
	}()
}

func NewKeyEventInterceptor() *KeyEventInterceptor {
	e := &KeyEventInterceptor{}
	e.ExtendBaseWidget(e)
	return e
}

func (e *KeyEventInterceptor) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyDown:
		if e.OnKeyDown != nil {
			e.OnKeyDown()
		}
	case fyne.KeyUp:
		if e.OnKeyUp != nil {
			e.OnKeyUp()
		}
	default:
		e.Entry.TypedKey(key)
	}
}
