package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/adelylria/GoFinder/core/configuration"
	"github.com/adelylria/GoFinder/core/i18n"
)

func (l *Launcher) showSettingsDialog() {
	content := container.NewPadded(l.buildSettingsForm())
	dialog.ShowCustom(i18n.T(i18n.MenuPreferences), i18n.T(i18n.DialogClose), content, l.window)
}

func (l *Launcher) showAboutDialog() {
	dialog.ShowInformation(i18n.T(i18n.MenuAbout), i18n.T(i18n.AboutText), l.window)
}

func (l *Launcher) buildSettingsForm() fyne.CanvasObject {
	modifiers := []string{"Alt", "Ctrl", "Shift", "Ctrl+Alt"}
	keys := letterOptions()
	status := widget.NewLabel("")
	initializing := true

	toggleModifier := widget.NewSelect(modifiers, func(value string) {
		if initializing {
			return
		}
		l.config.ToggleHotkey.Modifier = value
		l.saveSettings(status, true)
	})
	toggleModifier.SetSelected(l.config.ToggleHotkey.Modifier)

	toggleKey := widget.NewSelect(keys, func(value string) {
		if initializing {
			return
		}
		l.config.ToggleHotkey.Key = value
		l.saveSettings(status, true)
	})
	toggleKey.SetSelected(l.config.ToggleHotkey.Key)

	quitModifier := widget.NewSelect(modifiers, func(value string) {
		if initializing {
			return
		}
		l.config.QuitHotkey.Modifier = value
		l.saveSettings(status, true)
	})
	quitModifier.SetSelected(l.config.QuitHotkey.Modifier)

	quitKey := widget.NewSelect(keys, func(value string) {
		if initializing {
			return
		}
		l.config.QuitHotkey.Key = value
		l.saveSettings(status, true)
	})
	quitKey.SetSelected(l.config.QuitHotkey.Key)

	autoStart := widget.NewCheck(i18n.T(i18n.SettingsAutoStart), func(value bool) {
		if initializing {
			return
		}
		l.config.AutoStart = value
		l.saveSettings(status, false)
	})
	autoStart.SetChecked(l.config.AutoStart)

	startHidden := widget.NewCheck(i18n.T(i18n.SettingsHidden), func(value bool) {
		if initializing {
			return
		}
		l.config.StartHidden = value
		l.startHidden = value
		l.saveSettings(status, false)
	})
	startHidden.SetChecked(l.config.StartHidden)
	initializing = false

	hotkeys := container.NewVBox(
		widget.NewLabelWithStyle(i18n.T(i18n.SettingsHotkeys), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewGridWithColumns(3,
			widget.NewLabel(i18n.T(i18n.SettingsToggle)),
			toggleModifier,
			toggleKey,
		),
		container.NewGridWithColumns(3,
			widget.NewLabel(i18n.T(i18n.SettingsQuit)),
			quitModifier,
			quitKey,
		),
	)

	general := container.NewVBox(
		widget.NewLabelWithStyle(i18n.T(i18n.SettingsGeneral), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		autoStart,
		startHidden,
	)

	return container.NewVBox(hotkeys, widget.NewSeparator(), general, status)
}

func (l *Launcher) saveSettings(status *widget.Label, needsRestart bool) {
	if err := configuration.Save(l.config); err != nil {
		status.SetText(err.Error())
		return
	}
	if err := configuration.ApplyAutoStart(l.config.AutoStart); err != nil {
		status.SetText(err.Error())
		return
	}
	if needsRestart {
		status.SetText(i18n.T(i18n.SettingsRestart))
		return
	}
	status.SetText(i18n.T(i18n.SettingsSaved))
}
