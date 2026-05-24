package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"

	"github.com/adelylria/GoFinder/core/configuration"
	"github.com/adelylria/GoFinder/core/i18n"
)

func (l *Launcher) showSettingsDialog() {
	l.dialogsMu.Lock()
	if l.settingsOpen {
		// already showing settings in the main window; nothing to do
		l.dialogsMu.Unlock()
		return
	}
	// store previous content so we can restore it when exiting settings
	l.prevContent = l.window.Content()
	l.settingsOpen = true
	l.dialogsMu.Unlock()

	// Build the settings form and wrap it in a scrollable area
	form := l.buildSettingsForm()
	scroll := container.NewVScroll(form)
	scroll.SetMinSize(fyne.NewSize(480, 320))

	// Back button to restore previous content
	back := widget.NewButton("←", func() {
		l.closeSettingsView()
	})
	title := widget.NewLabelWithStyle(i18n.T(i18n.MenuPreferences), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	header := container.NewBorder(nil, nil, back, nil, container.NewHBox(title))

	content := container.NewBorder(header, nil, nil, nil, scroll)
	l.window.SetContent(content)
}

// closeSettingsView restores the previous window content when leaving settings
func (l *Launcher) closeSettingsView() {
	l.dialogsMu.Lock()
	if !l.settingsOpen {
		l.dialogsMu.Unlock()
		return
	}
	l.settingsOpen = false
	prev := l.prevContent
	l.prevContent = nil
	l.dialogsMu.Unlock()

	if prev != nil {
		l.window.SetContent(prev)
		// try to refocus the input field
		go fyne.Do(func() {
			time.Sleep(50 * time.Millisecond)
			l.window.Canvas().Focus(l.input)
		})
	} else {
		l.initializeUI()
	}
}

func (l *Launcher) showAboutDialog() {
	l.dialogsMu.Lock()
	if l.aboutDialog != nil {
		l.aboutDialog.Show()
		l.dialogsMu.Unlock()
		return
	}
	l.dialogsMu.Unlock()

	d := dialog.NewInformation(i18n.T(i18n.MenuAbout), i18n.T(i18n.AboutText), l.window)
	d.SetOnClosed(func() {
		l.dialogsMu.Lock()
		l.aboutDialog = nil
		l.dialogsMu.Unlock()
	})

	l.dialogsMu.Lock()
	l.aboutDialog = d
	l.dialogsMu.Unlock()
	d.Show()
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

	// Theme mode (system / light / dark) as a horizontal radio group
	modeOptions := []string{i18n.T(i18n.ThemeSystem), i18n.T(i18n.ThemeLight), i18n.T(i18n.ThemeDark)}
	modeRadio := widget.NewRadioGroup(modeOptions, func(label string) {
		if initializing {
			return
		}
		// only map the three main modes here
		l.config.ThemeName = themeNameFromLabel(label)
		applyAppTheme(l.config.ThemeName)
		l.saveSettings(status, false)
	})
	modeRadio.Horizontal = true

	// select initial mode if config is one of the main modes
	switch l.config.ThemeName {
	case "system", "light", "dark":
		modeRadio.SetSelected(themeLabel(l.config.ThemeName))
	default:
		modeRadio.SetSelected(i18n.T(i18n.ThemeSystem))
	}

	// Small preview tiles for each preset (excluding the three main modes)
	previews := container.NewGridWithColumns(3)
	for _, name := range themeOptions() {
		if name == "system" || name == "light" || name == "dark" {
			continue
		}
		t := themeForName(name)
		bg := canvas.NewRectangle(t.Color(theme.ColorNameBackground, theme.VariantLight))
		fg := canvas.NewRectangle(t.Color(theme.ColorNameForeground, theme.VariantLight))
		primary := canvas.NewRectangle(t.Color(theme.ColorNamePrimary, theme.VariantLight))

		bg.SetMinSize(fyne.NewSize(140, 48))
		fg.SetMinSize(fyne.NewSize(140, 10))
		primary.SetMinSize(fyne.NewSize(140, 10))

		title := widget.NewLabel(themeLabel(name))
		applyBtn := widget.NewButton(i18n.T(i18n.SettingsTheme), func() {
			if initializing {
				return
			}
			l.config.ThemeName = name
			applyAppTheme(name)
			l.saveSettings(status, false)
		})

		tile := container.NewVBox(title, bg, container.NewHBox(primary, fg), applyBtn)
		previews.Add(tile)
	}

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

	appearance := container.NewVBox(
		widget.NewLabelWithStyle(i18n.T(i18n.SettingsAppearance), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(i18n.T(i18n.SettingsTheme), fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
		modeRadio,
		widget.NewLabelWithStyle(i18n.T(i18n.SettingsAppearance), fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
		previews,
	)

	return container.NewVBox(hotkeys, widget.NewSeparator(), general, widget.NewSeparator(), appearance, status)
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
