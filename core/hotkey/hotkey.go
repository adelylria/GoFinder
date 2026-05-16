package hotkey

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type HotkeyManager struct {
	ToggleHandler func()
	ExitHandler   func()
	PrefsHandler  func()
	AboutHandler  func()
	ToggleHotkey  KeyBinding
	ExitHotkey    KeyBinding
}

func (hm *HotkeyManager) SetMenuHandlers(prefs, about func()) {
	hm.PrefsHandler = prefs
	hm.AboutHandler = about
}

type KeyBinding struct {
	Modifier string
	Key      string
}

type KeyEventInterceptor struct {
	widget.Entry
	OnKeyDown    func()
	OnKeyUp      func()
	OnMenuQuit   func()
	OnMenuPrefs  func()
	OnMenuAbout  func()
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

func (e *KeyEventInterceptor) TypedShortcut(shortcut fyne.Shortcut) {
	if handleMenuShortcut(shortcut, e.OnMenuQuit, e.OnMenuPrefs, e.OnMenuAbout) {
		return
	}
	e.Entry.TypedShortcut(shortcut)
}

func NewHotkeyManager(toggle func(), exit func(), bindings ...KeyBinding) *HotkeyManager {
	toggleHotkey := KeyBinding{Modifier: "Alt", Key: "R"}
	exitHotkey := KeyBinding{Modifier: "Alt", Key: "Q"}
	if len(bindings) > 0 {
		toggleHotkey = bindings[0]
	}
	if len(bindings) > 1 {
		exitHotkey = bindings[1]
	}
	return &HotkeyManager{
		ToggleHandler: toggle,
		ExitHandler:   exit,
		ToggleHotkey:  toggleHotkey,
		ExitHotkey:    exitHotkey,
	}
}
