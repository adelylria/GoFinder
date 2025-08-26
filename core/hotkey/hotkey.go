package hotkey

import "fyne.io/fyne/v2/widget"

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

func NewHotkeyManager(toggle func(), exit func()) *HotkeyManager {
	return &HotkeyManager{
		ToggleHandler: toggle,
		ExitHandler:   exit,
	}
}
