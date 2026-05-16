package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"

	"github.com/adelylria/GoFinder/core/i18n"
)

func (l *Launcher) configureNativeMenu() {
	l.setupNativeMainMenu()
}

func (l *Launcher) setupNativeMainMenu() {
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

	l.window.SetMainMenu(fyne.NewMainMenu(fileMenu, configMenu, helpMenu))
}
