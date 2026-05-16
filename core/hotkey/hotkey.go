package hotkey

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type HotkeyManager struct {
	ToggleHandler func()
	ExitHandler   func()
}

type KeyEventInterceptor struct {
	widget.Entry
	OnKeyDown func()
	OnKeyUp   func()
}

// NewKeyEventInterceptor crea el Entry personalizado para eventos de teclado
func NewKeyEventInterceptor() *KeyEventInterceptor {
	e := &KeyEventInterceptor{}
	e.ExtendBaseWidget(e)
	return e
}

func (e *KeyEventInterceptor) TypedKey(ev *fyne.KeyEvent) {
	switch ev.Name {
	case fyne.KeyDown:
		if e.OnKeyDown != nil {
			e.OnKeyDown()
			return
		}
	case fyne.KeyUp:
		if e.OnKeyUp != nil {
			e.OnKeyUp()
			return
		}
	}

	e.Entry.TypedKey(ev)
}

func NewHotkeyManager(toggle func(), exit func()) *HotkeyManager {
	return &HotkeyManager{
		ToggleHandler: toggle,
		ExitHandler:   exit,
	}
}
