package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/adelylria/GoFinder/core/i18n"
)

// hotkeysSection builds the hotkeys controls (toggle/quit bindings).
func (l *Launcher) hotkeysSection(initializing *bool) fyne.CanvasObject {
	toggleModifier, toggleKey := l.hotkeyControls(
		l.config.ToggleHotkey.Modifier,
		l.config.ToggleHotkey.Key,
		initializing,
		func(value string) { l.config.ToggleHotkey.Modifier = value },
		func(value string) { l.config.ToggleHotkey.Key = value },
	)
	quitModifier, quitKey := l.hotkeyControls(
		l.config.QuitHotkey.Modifier,
		l.config.QuitHotkey.Key,
		initializing,
		func(value string) { l.config.QuitHotkey.Modifier = value },
		func(value string) { l.config.QuitHotkey.Key = value },
	)

	return container.NewVBox(
		widget.NewLabelWithStyle(i18n.T(i18n.SettingsHotkeys), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		settingsHotkeyRow(i18n.T(i18n.SettingsToggle), toggleModifier, toggleKey),
		settingsHotkeyRow(i18n.T(i18n.SettingsQuit), quitModifier, quitKey),
	)
}

func (l *Launcher) hotkeyControls(
	modifier string,
	key string,
	initializing *bool,
	setModifier func(string),
	setKey func(string),
) (*widget.Select, *widget.Select) {
	modifierSelect := settingsSelect([]string{"Alt", "Ctrl", "Shift", "Ctrl+Alt"}, modifier, initializing, func(value string) {
		setModifier(value)
		l.saveSettings(i18n.T(i18n.SettingsHotkeysSaved))
	})
	keySelect := settingsSelect(letterOptions(), key, initializing, func(value string) {
		setKey(value)
		l.saveSettings(i18n.T(i18n.SettingsHotkeysSaved))
	})

	return modifierSelect, keySelect
}

func settingsSelect(opts []string, selected string, initializing *bool, onChange func(string)) *widget.Select {
	selectWidget := widget.NewSelect(opts, func(value string) {
		if *initializing {
			return
		}
		onChange(value)
	})
	selectWidget.SetSelected(selected)

	return selectWidget
}

func settingsHotkeyRow(label string, modifier *widget.Select, key *widget.Select) fyne.CanvasObject {
	return container.NewGridWithColumns(3, widget.NewLabel(label), modifier, key)
}

// generalSection builds the general settings controls.
func (l *Launcher) generalSection(initializing *bool) fyne.CanvasObject {
	autoStart := widget.NewCheck(i18n.T(i18n.SettingsAutoStart), func(value bool) {
		if *initializing {
			return
		}
		l.config.AutoStart = value
		l.saveSettings(i18n.T(i18n.SettingsAutoSaved))
	})
	autoStart.SetChecked(l.config.AutoStart)

	startHidden := widget.NewCheck(i18n.T(i18n.SettingsHidden), func(value bool) {
		if *initializing {
			return
		}
		l.config.StartHidden = value
		l.startHidden = value
		l.saveSettings(i18n.T(i18n.SettingsHiddenSaved))
	})
	startHidden.SetChecked(l.config.StartHidden)

	return container.NewVBox(
		widget.NewLabelWithStyle(i18n.T(i18n.SettingsGeneral), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		autoStart,
		startHidden,
	)
}
