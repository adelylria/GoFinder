package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/adelylria/GoFinder/core/hotkey"
)

// theme.go
// Archivo "core" para la configuración visual y métricas de la UI.
// Contiene tamaños, paddings, helpers para crear widgets con estilo
// y funciones utilitarias para aplicar el tema a la ventana.
//
// Está pensado para integrarse con el código del lanzador (launcher.go)
// y centralizar valores (ej: WindowSize, IconSize, ListItemHeight, etc.).

// ThemeConfig agrupa todas las métricas visuales y colores usados por la UI.
type ThemeConfig struct {
	WindowSize       fyne.Size   // tamaño por defecto de la ventana
	FixedWindow      bool        // si fijar el tamaño de la ventana
	InputHeight      float32     // altura del campo de búsqueda
	InputPlaceholder string      // placeholder por defecto para el entry
	ListItemHeight   float32     // altura de cada item en la lista
	Padding          float32     // padding general
	CornerRadius     float32     // radio de esquinas para rectángulos decorativos
	HighlightColor   color.Color // color de selección / resaltado
	DefaultIcon      fyne.Resource
}

// DefaultTheme devuelve un ThemeConfig con valores predeterminados,
// tomados del launcher existente (por ejemplo 600x500 y icon 32x32).
func DefaultTheme() *ThemeConfig {
	return &ThemeConfig{
		WindowSize:       fyne.NewSize(600, 500),
		FixedWindow:      true,
		InputHeight:      44,
		InputPlaceholder: "Buscar aplicación...",
		ListItemHeight:   56,
		Padding:          8,
		CornerRadius:     6,
		HighlightColor:   theme.Color(theme.ColorNameHover),
	}
}

// ApplyToWindow aplica los ajustes de tamaño/fijado/centrado a una ventana fyne.
func (t *ThemeConfig) ApplyToWindow(w fyne.Window) {
	if w == nil {
		return
	}
	w.SetFixedSize(t.FixedWindow)
	w.Resize(t.WindowSize)
	// Centrar si está disponible en esta plataforma
	if centerer, ok := w.(interface{ CenterOnScreen() }); ok {
		centerer.CenterOnScreen()
	}
}

// NewStyledSearchEntry crea un KeyEventInterceptor (Entry con captura de teclas)
// con las métricas del tema ya aplicadas. Asume que KeyEventInterceptor existe
// en el mismo paquete (launcher.go).
func (t *ThemeConfig) NewStyledSearchEntry() *hotkey.KeyEventInterceptor {
	entry := hotkey.NewKeyEventInterceptor()
	entry.SetPlaceHolder(t.InputPlaceholder)
	entry.MinSize().Min(fyne.NewSize(t.WindowSize.Width-(2*t.Padding), t.InputHeight))
	return entry
}

// NewStyledList crea un widget.List configurado para usar los tamaños del tema.
// getCount / createItem / updateItem se delegan a quien construye la lista.
func (t *ThemeConfig) NewStyledList(getCount func() int, createItem func() fyne.CanvasObject, updateItem func(id widget.ListItemID, obj fyne.CanvasObject)) *widget.List {
	widgetList := widget.NewList(getCount, createItem, updateItem)
	// Ajuste mínimo para que la lista ocupe el alto deseado
	min := fyne.NewSize(t.WindowSize.Width-(2*t.Padding), t.WindowSize.Height-(t.InputHeight+3*t.Padding))
	widgetList.MinSize().Min(min)
	return widgetList
}

// CreateListItemDefault devuelve un CanvasObject para usar como item en la lista
// usando la apariencia definida en el tema (icon + label y rectángulo de fondo).
func (t *ThemeConfig) CreateListItemDefault() fyne.CanvasObject {
	bg := canvas.NewRectangle(t.HighlightColor)
	bg.Hide()

	icon := widget.NewIcon(nil)

	label := widget.NewLabel("")
	label.TextStyle = fyne.TextStyle{Bold: true}

	content := container.NewHBox(icon, container.NewVBox(label))

	// Stack para poder tener un bg que se muestre/oculte
	stack := container.NewStack(bg, content)
	stack.Resize(fyne.NewSize(t.WindowSize.Width-(2*t.Padding), t.ListItemHeight))
	return stack
}

// UpdateListItemDefault actualiza los elementos típicos de un item creado con CreateListItemDefault.
// Se espera que obj sea el Container devuelto por CreateListItemDefault.
func (t *ThemeConfig) UpdateListItemDefault(id widget.ListItemID, obj fyne.CanvasObject, name string, iconRes fyne.Resource, selected bool) {
	stack, ok := obj.(*fyne.Container)
	if !ok || len(stack.Objects) < 2 {
		return
	}

	bg := stack.Objects[0].(*canvas.Rectangle)
	content := stack.Objects[1].(*fyne.Container)

	// content tiene HBox(icon, VBox(label))
	if len(content.Objects) > 0 {
		iconWidget, ok := content.Objects[0].(*widget.Icon)
		if ok {
			if iconRes != nil {
				iconWidget.SetResource(iconRes)
				iconWidget.Show()
			} else {
				iconWidget.Hide()
			}
		}
		// label dentro del VBox
		if len(content.Objects) > 1 {
			vbox := content.Objects[1].(*fyne.Container)
			if len(vbox.Objects) > 0 {
				label := vbox.Objects[0].(*widget.Label)
				label.SetText(name)
			}
		}
	}

	if selected {
		bg.Show()
	} else {
		bg.Hide()
	}
	bg.Refresh()
}

// ComputeListItemHeight devuelve la altura recomendada para un item según el tema.
func (t *ThemeConfig) ComputeListItemHeight() float32 {
	return t.ListItemHeight
}

// Helpers para spacing y contenedores cortos
func (t *ThemeConfig) Padded(content fyne.CanvasObject) *fyne.Container {
	return container.NewPadded(content)
}

func (t *ThemeConfig) NewBorderWithInputTop(input fyne.CanvasObject, body fyne.CanvasObject) *fyne.Container {
	// Border con input arriba y body en el centro
	return container.NewBorder(input, nil, nil, nil, body)
}

// Merge permite aplicar cambios puntuales de otro ThemeConfig (override).
// Cualquier campo no-nulo/positivo del "overrides" reemplaza al actual.
func (t *ThemeConfig) Merge(overrides *ThemeConfig) *ThemeConfig {
	if overrides == nil {
		return t
	}
	out := *t
	if overrides.WindowSize.Width > 0 && overrides.WindowSize.Height > 0 {
		out.WindowSize = overrides.WindowSize
	}
	if overrides.InputHeight > 0 {
		out.InputHeight = overrides.InputHeight
	}
	if overrides.InputPlaceholder != "" {
		out.InputPlaceholder = overrides.InputPlaceholder
	}
	if overrides.ListItemHeight > 0 {
		out.ListItemHeight = overrides.ListItemHeight
	}
	if overrides.Padding > 0 {
		out.Padding = overrides.Padding
	}
	if overrides.CornerRadius > 0 {
		out.CornerRadius = overrides.CornerRadius
	}
	if overrides.HighlightColor != nil {
		out.HighlightColor = overrides.HighlightColor
	}
	out.FixedWindow = overrides.FixedWindow
	return &out
}

// --- Nota ---
// Este archivo está pensado para que el resto de la UI (por ejemplo launcher.go)
// consuma helperes: DefaultTheme(), theme.ApplyToWindow(window),
// theme.NewStyledSearchEntry(), theme.NewStyledList(...), etc.
// De esta forma todos los tamaños y colores quedan centralizados y fáciles de
// ajustar.
