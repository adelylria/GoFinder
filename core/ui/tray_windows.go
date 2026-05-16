//go:build windows

package ui

import (
	"sync"

	"fyne.io/systray"
	"github.com/adelylria/GoFinder/core/i18n"
	"github.com/adelylria/GoFinder/models"
)

var trayOnce sync.Once

func startSystemTray(state *models.AppState, icon []byte) {
	trayOnce.Do(func() {
		go systray.Run(
			func() { setupSystemTray(state, icon) },
			func() {},
		)
	})
}

func setupSystemTray(state *models.AppState, icon []byte) {
	if len(icon) > 0 {
		systray.SetIcon(icon)
	}
	systray.SetTooltip(i18n.T(i18n.TrayTooltip))

	toggleItem := systray.AddMenuItem(i18n.T(i18n.TrayToggleTitle), i18n.T(i18n.TrayToggleTooltip))
	minimizeItem := systray.AddMenuItem(i18n.T(i18n.TrayMinimizeTitle), i18n.T(i18n.TrayMinimizeTip))
	systray.AddSeparator()
	quitItem := systray.AddMenuItem(i18n.T(i18n.TrayQuitTitle), i18n.T(i18n.TrayQuitTooltip))

	go func() {
		for {
			select {
			case <-toggleItem.ClickedCh:
				toggleWindowVisibility(state)
			case <-minimizeItem.ClickedCh:
				setWindowVisible(state, false)
			case <-quitItem.ClickedCh:
				systray.Quit()
				quitApplication()
			}
		}
	}()
}
