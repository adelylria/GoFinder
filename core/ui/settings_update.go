package ui

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/adelylria/GoFinder/core/i18n"
	"github.com/adelylria/GoFinder/core/resource"
	"github.com/adelylria/GoFinder/core/singleinstance"
	"github.com/adelylria/GoFinder/core/updater"
	"github.com/adelylria/GoFinder/core/version"
)

type updateSettingsView struct {
	currentLabel  *widget.Label
	latestLabel   *widget.Label
	commitLabel   *widget.Label
	buildLabel    *widget.Label
	repoRow       fyne.CanvasObject
	statusLabel   *widget.Label
	checkButton   *widget.Button
	updateButton  *widget.Button
	releaseButton *widget.Button
	info          updater.ReleaseInfo
}

func (l *Launcher) showUpdateSettings() {
	view := l.newUpdateSettingsView()
	l.setSettingsContent(i18n.T(i18n.SettingsUpdates), l.showSettingsHome, view.content())
	l.checkForUpdates(view)
}

func (l *Launcher) newUpdateSettingsView() *updateSettingsView {
	// helper to extract a short title before the ':' in translation strings
	extractTitle := func(translated string) string {
		if idx := strings.Index(translated, ":"); idx >= 0 {
			return strings.TrimSpace(translated[:idx])
		}
		return translated
	}

	titleRepo := widget.NewLabel(extractTitle(fmt.Sprintf(i18n.T(i18n.UpdateRepository), "")))
	parsedRepo, _ := url.Parse(version.RepoURL)
	repoLink := widget.NewHyperlink(version.RepoURL, parsedRepo)
	repoLink.Text = version.RepoURL
	repoLink.Alignment = fyne.TextAlignTrailing
	repoIcon := widget.NewButtonWithIcon("", resource.GetEmbedGithubIcon(), func() {
		l.openURL(version.RepoURL)
	})
	repoIcon.Importance = widget.LowImportance

	// For repository show only the GitHub icon at the right; label on the left
	repoLinkIcon := widget.NewButtonWithIcon("", resource.GetEmbedGithubIcon(), func() {
		l.openURL(version.RepoURL)
	})
	repoLinkIcon.Importance = widget.LowImportance
	repoRow := container.NewHBox(widget.NewLabelWithStyle(titleRepo.Text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), layout.NewSpacer(), repoLinkIcon)

	// Create value labels and arrange them in two-column rows for a cleaner layout
	currentVal := widget.NewLabel(version.Version)
	latestVal := widget.NewLabel("-")
	commitVal := widget.NewLabel(version.Commit)
	buildVal := widget.NewLabel(version.BuildDate)

	currentRow := container.NewHBox(widget.NewLabelWithStyle(extractTitle(fmt.Sprintf(i18n.T(i18n.UpdateCurrent), "")), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), layout.NewSpacer(), currentVal)
	latestRow := container.NewHBox(widget.NewLabelWithStyle(extractTitle(fmt.Sprintf(i18n.T(i18n.UpdateLatest), "")), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), layout.NewSpacer(), latestVal)
	commitRow := container.NewHBox(widget.NewLabelWithStyle(extractTitle(fmt.Sprintf(i18n.T(i18n.UpdateCommit), "")), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), layout.NewSpacer(), commitVal)
	buildRow := container.NewHBox(widget.NewLabelWithStyle(extractTitle(fmt.Sprintf(i18n.T(i18n.UpdateBuildDate), "")), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), layout.NewSpacer(), buildVal)

	view := &updateSettingsView{
		currentLabel: currentVal,
		latestLabel:  latestVal,
		commitLabel:  commitVal,
		buildLabel:   buildVal,
		repoRow: container.NewVBox(
			currentRow,
			widget.NewSeparator(),
			latestRow,
			widget.NewSeparator(),
			commitRow,
			widget.NewSeparator(),
			buildRow,
			widget.NewSeparator(),
			repoRow,
		),
		statusLabel: widget.NewLabel(""),
	}

	view.checkButton = widget.NewButtonWithIcon(i18n.T(i18n.UpdateCheck), theme.ViewRefreshIcon(), func() {
		l.checkForUpdates(view)
	})
	view.updateButton = widget.NewButtonWithIcon(i18n.T(i18n.UpdateDownload), theme.DownloadIcon(), func() {
		l.confirmAndApplyUpdate(view.info)
	})
	view.releaseButton = widget.NewButtonWithIcon(i18n.T(i18n.UpdateOpenRelease), resource.GetEmbedGithubIcon(), func() {
		l.openURL(view.info.ReleaseURL)
	})
	view.updateButton.Disable()
	view.releaseButton.Disable()

	return view
}

func (v *updateSettingsView) content() fyne.CanvasObject {
	header := widget.NewLabelWithStyle(i18n.T(i18n.SettingsUpdates), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Wrap main content in padded VBox for better spacing
	main := container.NewPadded(v.repoRow)

	actions := container.NewHBox(v.checkButton, v.releaseButton, v.updateButton)

	return container.NewVBox(
		header,
		main,
		v.statusLabel,
		actions,
	)
}

func (l *Launcher) checkForUpdates(view *updateSettingsView) {
	l.setUpdateChecking(view)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
		defer cancel()

		info, err := updater.CheckLatest(ctx, version.Version)
		fyne.Do(func() {
			l.applyUpdateCheckResult(view, info, err)
		})
	}()
}

func (l *Launcher) setUpdateChecking(view *updateSettingsView) {
	view.statusLabel.SetText(i18n.T(i18n.UpdateChecking))
	view.checkButton.Disable()
	view.updateButton.Disable()
	view.releaseButton.Disable()
}

func (l *Launcher) applyUpdateCheckResult(view *updateSettingsView, info updater.ReleaseInfo, err error) {
	if err != nil {
		view.statusLabel.SetText(err.Error())
		view.checkButton.Enable()
		return
	}

	view.info = info
	view.latestLabel.SetText(fmt.Sprintf(i18n.T(i18n.UpdateLatest), info.LatestVersion))
	view.checkButton.Enable()
	view.releaseButton.Enable()
	if !info.UpdateAvailable {
		view.statusLabel.SetText(i18n.T(i18n.UpdateUnavailable))
		return
	}

	view.statusLabel.SetText(i18n.T(i18n.UpdateAvailable))
	if info.DownloadURL != "" {
		view.updateButton.Enable()
	}
}

func (l *Launcher) confirmAndApplyUpdate(info updater.ReleaseInfo) {
	dialog.NewConfirm(
		i18n.T(i18n.UpdateConfirmTitle),
		i18n.T(i18n.UpdateConfirmMessage),
		func(ok bool) {
			if ok {
				go l.downloadAndApplyUpdate(info)
			}
		},
		l.window,
	).Show()
}

func (l *Launcher) downloadAndApplyUpdate(info updater.ReleaseInfo) {
	l.showSettingsToast(i18n.T(i18n.UpdateDownloading))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	downloadPath, err := updater.Download(ctx, info.DownloadURL)
	if err != nil {
		l.showSettingsToast(err.Error())
		return
	}
	if err := updater.ApplyDownloadedUpdate(downloadPath); err != nil {
		l.showSettingsToast(err.Error())
		return
	}

	l.showSettingsToast(i18n.T(i18n.UpdateRestarting))
	time.Sleep(800 * time.Millisecond)
	singleinstance.Release()
	os.Exit(0)
}

func (l *Launcher) openURL(rawURL string) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		l.showSettingsToast(err.Error())
		return
	}
	if err := fyne.CurrentApp().OpenURL(parsed); err != nil {
		l.showSettingsToast(err.Error())
	}
}
