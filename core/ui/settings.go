package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"

	"github.com/adelylria/GoFinder/core/configuration"
	"github.com/adelylria/GoFinder/core/i18n"
)

func (l *Launcher) showSettingsDialog() {
	l.dialogsMu.Lock()
	if l.settingsOpen {
		l.dialogsMu.Unlock()
		return
	}

	l.prevContent = l.window.Content()
	l.settingsOpen = true
	l.dialogsMu.Unlock()

	l.showSettingsHome()
}

// closeSettingsView restores the previous window content when leaving settings.
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

	if prev == nil {
		l.initializeUI()
		return
	}

	l.window.SetContent(prev)
	go fyne.Do(func() {
		time.Sleep(50 * time.Millisecond)
		l.window.Canvas().Focus(l.input)
	})
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

func (l *Launcher) saveSettings(successMessage string) {
	if err := configuration.Save(l.config); err != nil {
		l.showSettingsToast(err.Error())
		return
	}
	if err := configuration.ApplyAutoStart(l.config.AutoStart); err != nil {
		l.showSettingsToast(err.Error())
		return
	}
	l.showSettingsToast(successMessage)
}
