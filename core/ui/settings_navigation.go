package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/adelylria/GoFinder/core/i18n"
)

func (l *Launcher) showSettingsHome() {
	body := container.NewVBox(
		l.settingsNavCard(
			i18n.T(i18n.SettingsAppearance),
			i18n.T(i18n.SettingsThemeDesc),
			theme.ColorPaletteIcon(),
			l.showAppearanceSettings,
		),
		l.settingsNavCard(
			i18n.T(i18n.SettingsGeneral),
			i18n.T(i18n.SettingsConfigDesc),
			theme.SettingsIcon(),
			l.showConfigurationSettings,
		),
	)

	l.setSettingsRootContent(i18n.T(i18n.MenuPreferences), body)
}

func (l *Launcher) showAppearanceSettings() {
	initializing := true
	appearance := l.appearanceSection(&initializing)
	initializing = false

	l.setSettingsContent(
		i18n.T(i18n.SettingsAppearance),
		l.showSettingsHome,
		appearance,
	)
}

func (l *Launcher) showConfigurationSettings() {
	initializing := true
	general := l.generalSection(&initializing)
	hotkeys := l.hotkeysSection(&initializing)
	initializing = false

	l.setSettingsContent(
		i18n.T(i18n.SettingsGeneral),
		l.showSettingsHome,
		container.NewVBox(general, widget.NewSeparator(), hotkeys),
	)
}

func (l *Launcher) setSettingsContent(title string, onBack func(), body fyne.CanvasObject) {
	l.setSettingsContentWithIcon(title, theme.NavigateBackIcon(), onBack, body)
}

func (l *Launcher) setSettingsRootContent(title string, body fyne.CanvasObject) {
	l.setSettingsContentWithIcon(title, theme.HomeIcon(), l.closeSettingsView, body)
}

func (l *Launcher) setSettingsContentWithIcon(title string, icon fyne.Resource, onTap func(), body fyne.CanvasObject) {
	scroll := container.NewVScroll(body)
	scroll.SetMinSize(fyne.NewSize(480, 320))

	back := widget.NewButtonWithIcon("", icon, onTap)
	titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	header := container.NewBorder(nil, nil, back, nil, titleLabel)
	toast := l.newSettingsToast()

	toastOverlay := container.NewStack(
		scroll,
		container.NewBorder(nil, container.NewPadded(container.NewCenter(toast)), nil, nil),
	)

	content := container.NewBorder(header, nil, nil, nil, toastOverlay)
	l.window.SetContent(content)
}

func (l *Launcher) newSettingsToast() fyne.CanvasObject {
	// Same styling pattern as settingsNavCard: ColorNameButton bg + widget.Label text.
	// widget.Label automatically uses ColorNameForeground, adapting to dark/light themes.
	label := widget.NewLabel("")
	label.Alignment = fyne.TextAlignCenter

	bg := canvas.NewRectangle(theme.Color(theme.ColorNameButton))
	bg.CornerRadius = 12
	bg.StrokeColor = theme.Color(theme.ColorNameSeparator)
	bg.StrokeWidth = 1
	// Force a reasonable pill size so container.NewCenter renders it compactly.
	bg.SetMinSize(fyne.NewSize(420, 40))

	toast := container.NewStack(bg, container.NewPadded(label))
	toast.Hide()

	l.settingsToast = toast
	l.settingsToastBg = bg
	l.settingsToastLabel = label

	return toast
}

func (l *Launcher) showSettingsToast(message string) {
	l.dialogsMu.Lock()
	l.settingsToastVersion++
	version := l.settingsToastVersion
	toast := l.settingsToast
	bg := l.settingsToastBg
	label := l.settingsToastLabel
	l.dialogsMu.Unlock()

	if toast == nil || label == nil {
		return
	}

	fyne.Do(func() {
		// Re-read theme colors so the toast always matches the active theme.
		if bg != nil {
			bg.FillColor = theme.Color(theme.ColorNameButton)
			bg.StrokeColor = theme.Color(theme.ColorNameSeparator)
			bg.Refresh()
		}
		label.SetText(message)
		toast.Show()
		toast.Refresh()
	})
	go l.hideSettingsToastLater(version)
}

func (l *Launcher) hideSettingsToastLater(version uint64) {
	time.Sleep(time.Second)

	fyne.Do(func() {
		l.dialogsMu.Lock()
		defer l.dialogsMu.Unlock()
		if version != l.settingsToastVersion || l.settingsToast == nil {
			return
		}
		l.settingsToast.Hide()
	})
}

func (l *Launcher) settingsNavCard(title string, description string, icon fyne.Resource, onTap func()) fyne.CanvasObject {
	cardBg := settingsCardBackground(8)
	titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	descriptionLabel := widget.NewLabel(description)
	descriptionLabel.Wrapping = fyne.TextWrapWord

	iconWidget := widget.NewIcon(icon)
	text := container.NewVBox(titleLabel, descriptionLabel)
	content := container.NewBorder(nil, nil, iconWidget, nil, text)
	card := container.NewStack(cardBg, container.NewPadded(content))

	return newSettingsCardButton(card, onTap)
}

func settingsCardBackground(radius float32) *canvas.Rectangle {
	cardBg := canvas.NewRectangle(theme.Color(theme.ColorNameButton))
	cardBg.CornerRadius = radius
	cardBg.StrokeColor = theme.Color(theme.ColorNameSeparator)
	cardBg.StrokeWidth = 1

	return cardBg
}

type settingsCardButton struct {
	widget.BaseWidget
	content fyne.CanvasObject
	onTap   func()
}

func newSettingsCardButton(content fyne.CanvasObject, onTap func()) *settingsCardButton {
	card := &settingsCardButton{content: content, onTap: onTap}
	card.ExtendBaseWidget(card)

	return card
}

func (c *settingsCardButton) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.content)
}

func (c *settingsCardButton) Tapped(*fyne.PointEvent) {
	if c.onTap == nil {
		return
	}
	c.onTap()
}

func (c *settingsCardButton) TappedSecondary(*fyne.PointEvent) {}
