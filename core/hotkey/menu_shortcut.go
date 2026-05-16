package hotkey

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

func handleMenuShortcut(shortcut fyne.Shortcut, onQuit, onPrefs, onAbout func()) bool {
	custom, ok := shortcut.(*desktop.CustomShortcut)
	if !ok {
		return false
	}

	switch {
	case custom.Modifier == fyne.KeyModifierControl && custom.KeyName == fyne.KeyQ:
		if onQuit != nil {
			onQuit()
		}
		return true
	case custom.Modifier == fyne.KeyModifierControl && custom.KeyName == fyne.KeyComma:
		if onPrefs != nil {
			onPrefs()
		}
		return true
	case custom.Modifier == 0 && custom.KeyName == fyne.KeyF1:
		if onAbout != nil {
			onAbout()
		}
		return true
	default:
		return false
	}
}
