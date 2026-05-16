package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"github.com/adelylria/GoFinder/core/i18n"
)

func (l *Launcher) createMenuBar() fyne.CanvasObject {
	exitItem := fyne.NewMenuItem(i18n.T(i18n.MenuExit), quitApplication)
	exitItem.Shortcut = &desktop.CustomShortcut{
		KeyName:  fyne.KeyQ,
		Modifier: fyne.KeyModifierControl,
	}

	prefItem := fyne.NewMenuItem(i18n.T(i18n.MenuPreferences), l.showSettingsDialog)
	prefItem.Shortcut = &desktop.CustomShortcut{
		KeyName:  fyne.KeyComma,
		Modifier: fyne.KeyModifierControl,
	}

	aboutItem := fyne.NewMenuItem(i18n.T(i18n.MenuAbout), l.showAboutDialog)
	aboutItem.Shortcut = &desktop.CustomShortcut{
		KeyName: fyne.KeyF1,
	}

	fileMenu := fyne.NewMenu(i18n.T(i18n.MenuFile), exitItem)
	configMenu := fyne.NewMenu(i18n.T(i18n.MenuConfig), prefItem)
	helpMenu := fyne.NewMenu(i18n.T(i18n.MenuHelp), aboutItem)

	return container.NewHBox(
		l.newMenuButton(fileMenu),
		l.newMenuButton(configMenu),
		l.newMenuButton(helpMenu),
	)
}

func (l *Launcher) newMenuButton(menu *fyne.Menu) fyne.CanvasObject {
	var btn *widget.Button
	btn = widget.NewButton(menu.Label, func() {
		l.showPopUpMenu(btn, menu)
	})
	btn.Importance = widget.LowImportance
	return btn
}

func (l *Launcher) showPopUpMenu(anchor fyne.CanvasObject, menu *fyne.Menu) {
	pop := widget.NewPopUpMenu(menu, l.window.Canvas())
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(anchor)
	pop.ShowAtPosition(pos.Add(fyne.NewPos(0, anchor.Size().Height)))
}
