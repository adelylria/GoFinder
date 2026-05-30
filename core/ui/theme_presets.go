package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/adelylria/GoFinder/core/i18n"
)

type themePreset struct {
	name       string
	base       fyne.Theme
	background color.Color
	button     color.Color
	foreground color.Color
	primary    color.Color
	hover      color.Color
	focus      color.Color
	selection  color.Color
	shadow     color.Color
}

type goFinderTheme struct {
	themePreset
}

// fixedVariantTheme forces a specific theme variant (light/dark) while delegating
// all rendering to the wrapped base theme.
type fixedVariantTheme struct {
	base    fyne.Theme
	variant fyne.ThemeVariant
}

func (t fixedVariantTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return t.base.Color(name, t.variant)
}
func (t fixedVariantTheme) Font(style fyne.TextStyle) fyne.Resource    { return t.base.Font(style) }
func (t fixedVariantTheme) Icon(name fyne.ThemeIconName) fyne.Resource { return t.base.Icon(name) }
func (t fixedVariantTheme) Size(name fyne.ThemeSizeName) float32       { return t.base.Size(name) }

func applyAppTheme(name string) {
	app := fyne.CurrentApp()
	if app == nil {
		return
	}
	app.Settings().SetTheme(themeForName(name))
}

func themeForName(name string) fyne.Theme {
	switch name {
	case "system":
		return theme.DefaultTheme()
	case "light":
		return fixedVariantTheme{base: theme.DefaultTheme(), variant: theme.VariantLight}
	case "dark":
		return fixedVariantTheme{base: theme.DefaultTheme(), variant: theme.VariantDark}
	case "paper":
		return goFinderTheme{themePreset{
			name:       "paper",
			base:       fixedVariantTheme{base: theme.DefaultTheme(), variant: theme.VariantLight},
			background: color.NRGBA{R: 250, G: 250, B: 247, A: 255},
			button:     color.NRGBA{R: 240, G: 240, B: 236, A: 255},
			foreground: color.NRGBA{R: 20, G: 20, B: 23, A: 255},
			primary:    color.NRGBA{R: 14, G: 120, B: 255, A: 255},
			hover:      color.NRGBA{R: 230, G: 230, B: 226, A: 255},
			focus:      color.NRGBA{R: 0, G: 102, B: 255, A: 255},
			selection:  color.NRGBA{R: 200, G: 230, B: 255, A: 255},
			shadow:     color.NRGBA{R: 0, G: 0, B: 0, A: 40},
		}}
	case "pastel":
		return goFinderTheme{themePreset{
			name:       "pastel",
			base:       fixedVariantTheme{base: theme.DefaultTheme(), variant: theme.VariantLight},
			background: color.NRGBA{R: 252, G: 248, B: 255, A: 255},
			button:     color.NRGBA{R: 245, G: 240, B: 250, A: 255},
			foreground: color.NRGBA{R: 44, G: 45, B: 49, A: 255},
			primary:    color.NRGBA{R: 138, G: 112, B: 255, A: 255},
			hover:      color.NRGBA{R: 240, G: 235, B: 245, A: 255},
			focus:      color.NRGBA{R: 160, G: 130, B: 255, A: 255},
			selection:  color.NRGBA{R: 230, G: 220, B: 255, A: 255},
			shadow:     color.NRGBA{R: 0, G: 0, B: 0, A: 40},
		}}
	case "solar":
		return goFinderTheme{themePreset{
			name:       "solar",
			base:       fixedVariantTheme{base: theme.DefaultTheme(), variant: theme.VariantLight},
			background: color.NRGBA{R: 255, G: 250, B: 240, A: 255},
			button:     color.NRGBA{R: 255, G: 245, B: 230, A: 255},
			foreground: color.NRGBA{R: 30, G: 30, B: 30, A: 255},
			primary:    color.NRGBA{R: 255, G: 162, B: 38, A: 255},
			hover:      color.NRGBA{R: 255, G: 235, B: 210, A: 255},
			focus:      color.NRGBA{R: 255, G: 190, B: 80, A: 255},
			selection:  color.NRGBA{R: 255, G: 230, B: 200, A: 255},
			shadow:     color.NRGBA{R: 0, G: 0, B: 0, A: 40},
		}}
	case "ocean":
		return goFinderTheme{themePreset{
			name:       "ocean",
			base:       fixedVariantTheme{base: theme.DefaultTheme(), variant: theme.VariantDark},
			background: color.NRGBA{R: 17, G: 32, B: 44, A: 255},
			button:     color.NRGBA{R: 31, G: 52, B: 66, A: 255},
			foreground: color.NRGBA{R: 232, G: 241, B: 246, A: 255},
			primary:    color.NRGBA{R: 47, G: 166, B: 185, A: 255},
			hover:      color.NRGBA{R: 38, G: 74, B: 88, A: 255},
			focus:      color.NRGBA{R: 71, G: 190, B: 210, A: 255},
			selection:  color.NRGBA{R: 39, G: 97, B: 112, A: 255},
			shadow:     color.NRGBA{R: 4, G: 8, B: 12, A: 190},
		}}
	case "forest":
		return goFinderTheme{themePreset{
			name:       "forest",
			base:       fixedVariantTheme{base: theme.DefaultTheme(), variant: theme.VariantDark},
			background: color.NRGBA{R: 25, G: 34, B: 28, A: 255},
			button:     color.NRGBA{R: 38, G: 53, B: 43, A: 255},
			foreground: color.NRGBA{R: 235, G: 242, B: 235, A: 255},
			primary:    color.NRGBA{R: 94, G: 173, B: 112, A: 255},
			hover:      color.NRGBA{R: 51, G: 74, B: 58, A: 255},
			focus:      color.NRGBA{R: 121, G: 202, B: 139, A: 255},
			selection:  color.NRGBA{R: 57, G: 105, B: 69, A: 255},
			shadow:     color.NRGBA{R: 6, G: 10, B: 7, A: 190},
		}}
	case "midnight":
		return goFinderTheme{themePreset{
			name:       "midnight",
			base:       fixedVariantTheme{base: theme.DefaultTheme(), variant: theme.VariantDark},
			background: color.NRGBA{R: 10, G: 14, B: 20, A: 255},
			button:     color.NRGBA{R: 20, G: 26, B: 33, A: 255},
			foreground: color.NRGBA{R: 220, G: 225, B: 230, A: 255},
			primary:    color.NRGBA{R: 90, G: 150, B: 230, A: 255},
			hover:      color.NRGBA{R: 25, G: 33, B: 42, A: 255},
			focus:      color.NRGBA{R: 110, G: 180, B: 255, A: 255},
			selection:  color.NRGBA{R: 40, G: 60, B: 90, A: 255},
			shadow:     color.NRGBA{R: 0, G: 0, B: 0, A: 200},
		}}
	case "contrast":
		return goFinderTheme{themePreset{
			name:       "contrast",
			base:       fixedVariantTheme{base: theme.DefaultTheme(), variant: theme.VariantDark},
			background: color.NRGBA{R: 12, G: 12, B: 12, A: 255},
			button:     color.NRGBA{R: 36, G: 36, B: 36, A: 255},
			foreground: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
			primary:    color.NRGBA{R: 255, G: 196, B: 30, A: 255},
			hover:      color.NRGBA{R: 64, G: 64, B: 64, A: 255},
			focus:      color.NRGBA{R: 255, G: 234, B: 83, A: 255},
			selection:  color.NRGBA{R: 120, G: 90, B: 0, A: 255},
			shadow:     color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		}}
	default:
		return theme.DefaultTheme()
	}
}

func (t goFinderTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return t.background
	case theme.ColorNameButton:
		return t.button
	case theme.ColorNameForeground:
		return t.foreground
	case theme.ColorNamePrimary, theme.ColorNameHyperlink:
		return t.primary
	case theme.ColorNameHover:
		return t.hover
	case theme.ColorNameFocus:
		return t.focus
	case theme.ColorNameSelection:
		return t.selection
	case theme.ColorNameShadow:
		return t.shadow
	case theme.ColorNameForegroundOnPrimary:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	default:
		return t.base.Color(name, variant)
	}
}

func (t goFinderTheme) Font(style fyne.TextStyle) fyne.Resource {
	return t.base.Font(style)
}

func (t goFinderTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.base.Icon(name)
}

func (t goFinderTheme) Size(name fyne.ThemeSizeName) float32 {
	return t.base.Size(name)
}

func themeOptions() []string {
	return []string{"system", "light", "dark", "paper", "pastel", "solar", "ocean", "forest", "midnight", "contrast"}
}

func themeLabel(name string) string {
	switch name {
	case "light":
		return i18n.T(i18n.ThemeLight)
	case "dark":
		return i18n.T(i18n.ThemeDark)
	case "ocean":
		return i18n.T(i18n.ThemeOcean)
	case "forest":
		return i18n.T(i18n.ThemeForest)
	case "contrast":
		return i18n.T(i18n.ThemeContrast)
	case "paper":
		return i18n.T(i18n.ThemePaper)
	case "pastel":
		return i18n.T(i18n.ThemePastel)
	case "solar":
		return i18n.T(i18n.ThemeSolar)
	case "midnight":
		return i18n.T(i18n.ThemeMidnight)
	default:
		return i18n.T(i18n.ThemeSystem)
	}
}

func themeNameFromLabel(label string) string {
	for _, option := range themeOptions() {
		if themeLabel(option) == label {
			return option
		}
	}
	return "system"
}
