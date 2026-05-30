package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/adelylria/GoFinder/core/i18n"
)

// appearanceSection builds the theme mode selector and preview tiles.
func (l *Launcher) appearanceSection(initializing *bool) fyne.CanvasObject {
	modeRadio := l.themeModeRadio(initializing)
	previews := container.NewGridWithColumns(3)

	for _, name := range themeOptions() {
		if isThemeModeOption(name) {
			continue
		}
		previews.Add(l.createPreviewTile(name, initializing))
	}

	return container.NewVBox(
		widget.NewLabelWithStyle(i18n.T(i18n.SettingsTheme), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		modeRadio,
		previews,
	)
}

func (l *Launcher) themeModeRadio(initializing *bool) *widget.RadioGroup {
	modeOptions := []string{i18n.T(i18n.ThemeSystem), i18n.T(i18n.ThemeLight), i18n.T(i18n.ThemeDark)}
	modeRadio := widget.NewRadioGroup(modeOptions, func(label string) {
		if *initializing {
			return
		}
		l.config.ThemeName = themeNameFromLabel(label)
		applyAppTheme(l.config.ThemeName)
		l.saveSettings(i18n.T(i18n.SettingsThemeSaved))
	})
	modeRadio.Horizontal = true
	modeRadio.SetSelected(l.selectedThemeModeLabel())

	return modeRadio
}

func (l *Launcher) selectedThemeModeLabel() string {
	if isThemeModeOption(l.config.ThemeName) {
		return themeLabel(l.config.ThemeName)
	}

	return i18n.T(i18n.ThemeSystem)
}

func isThemeModeOption(name string) bool {
	return name == "system" || name == "light" || name == "dark"
}

// createPreviewTile builds a small preview tile for a preset theme.
func (l *Launcher) createPreviewTile(name string, initializing *bool) fyne.CanvasObject {
	colors := themePreviewColors(themeForName(name))
	preview := themePreview(colors)
	swatches := container.NewGridWithColumns(5,
		themeSwatch(colors.background),
		themeSwatch(colors.button),
		themeSwatch(colors.foreground),
		themeSwatch(colors.primary),
		themeSwatch(colors.selection),
	)

	title := themePreviewTitle(themeLabel(name), colors.foreground)
	applyBtn := widget.NewButton(i18n.T(i18n.SettingsApply), func() {
		if *initializing {
			return
		}
		l.config.ThemeName = name
		applyAppTheme(name)
		l.saveSettings(i18n.T(i18n.SettingsThemeSaved))
	})

	content := container.NewVBox(title, preview, swatches, applyBtn)
	return container.NewStack(themePreviewCardBackground(colors), container.NewPadded(content))
}

type themePreviewPalette struct {
	background color.Color
	button     color.Color
	foreground color.Color
	primary    color.Color
	hover      color.Color
	selection  color.Color
}

func themePreviewColors(t fyne.Theme) themePreviewPalette {
	variant := theme.VariantLight

	return themePreviewPalette{
		background: t.Color(theme.ColorNameBackground, variant),
		button:     t.Color(theme.ColorNameButton, variant),
		foreground: t.Color(theme.ColorNameForeground, variant),
		primary:    t.Color(theme.ColorNamePrimary, variant),
		hover:      t.Color(theme.ColorNameHover, variant),
		selection:  t.Color(theme.ColorNameSelection, variant),
	}
}

func themePreviewTitle(text string, fill color.Color) fyne.CanvasObject {
	title := canvas.NewText(text, fill)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.TextSize = theme.TextSize()

	return title
}

func themePreviewCardBackground(colors themePreviewPalette) *canvas.Rectangle {
	cardBg := canvas.NewRectangle(colors.button)
	cardBg.CornerRadius = 6
	cardBg.StrokeColor = colors.selection
	cardBg.StrokeWidth = 1

	return cardBg
}

func themePreview(colors themePreviewPalette) fyne.CanvasObject {
	previewBg := canvas.NewRectangle(colors.background)
	previewBg.CornerRadius = 5
	previewBg.StrokeColor = theme.Color(theme.ColorNameShadow)
	previewBg.StrokeWidth = 1
	previewBg.SetMinSize(fyne.NewSize(150, 56))

	content := container.NewVBox(
		themePreviewLine(colors.primary, 42, 10),
		themePreviewLine(colors.foreground, 80, 8),
		themePreviewLine(colors.hover, 118, 8),
	)

	return container.NewStack(previewBg, container.NewPadded(content))
}

func themePreviewLine(fill color.Color, width float32, height float32) fyne.CanvasObject {
	line := canvas.NewRectangle(fill)
	line.CornerRadius = 3
	line.SetMinSize(fyne.NewSize(width, height))

	return line
}

func themeSwatch(fill color.Color) fyne.CanvasObject {
	swatch := canvas.NewRectangle(fill)
	swatch.CornerRadius = 4
	swatch.StrokeColor = theme.Color(theme.ColorNameSeparator)
	swatch.StrokeWidth = 1
	swatch.SetMinSize(fyne.NewSize(22, 18))

	return swatch
}
