package ui

import (
	"log"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/adelylria/GoFinder/logic"
	"github.com/adelylria/GoFinder/models"
)

// KeyEventInterceptor intercepta eventos de teclado para manejar la navegación
type KeyEventInterceptor struct {
	widget.Entry
	onKeyDown func()
	onKeyUp   func()
}

func NewKeyEventInterceptor() *KeyEventInterceptor {
	e := &KeyEventInterceptor{}
	e.ExtendBaseWidget(e)
	return e
}

func (e *KeyEventInterceptor) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyDown:
		if e.onKeyDown != nil {
			e.onKeyDown()
		}
	case fyne.KeyUp:
		if e.onKeyUp != nil {
			e.onKeyUp()
		}
	default:
		e.Entry.TypedKey(key)
	}
}

func RunLauncher(apps []models.Application) {
	myApp := app.New()
	myWindow := myApp.NewWindow("Buscador")
	myWindow.SetFixedSize(true)
	myWindow.Resize(fyne.NewSize(600, 500))

	input := NewKeyEventInterceptor()
	input.SetPlaceHolder("Buscar aplicación...")

	// Crear mapa de aplicaciones por ID
	appMap := make(map[string]models.Application)
	for _, app := range apps {
		appMap[app.ID] = app
	}

	// Variables de estado
	var (
		filteredIDs   []string
		selectedIndex int
		listWidget    *widget.List
	)

	// Actualizar lista filtrada
	updateFilter := func() {
		filteredIDs = getFilteredIDs(input.Text, appMap)
		selectedIndex = 0 // Resetear selección al filtrar
		if listWidget != nil {
			listWidget.Refresh()
		}
	}

	// Ejecutar la aplicación seleccionada
	executeSelected := func() {
		if len(filteredIDs) == 0 {
			return
		}

		if selectedIndex >= len(filteredIDs) {
			selectedIndex = len(filteredIDs) - 1
		}

		appID := filteredIDs[selectedIndex]
		app := appMap[appID]
		log.Printf("Ejecutando: %s (%s)", app.Name, app.Exec)

		if err := logic.RunApplication(app); err != nil {
			log.Printf("Error al ejecutar %s: %v", app.Name, err)
		}
		// if i want put else to hide
	}

	// Inicializar con todos los resultados
	updateFilter()

	// Crear lista personalizada con resaltado de selección
	listWidget = widget.NewList(
		func() int {
			return len(filteredIDs)
		},
		func() fyne.CanvasObject {
			// Fondo para resaltar selección (invisible por defecto)
			bg := canvas.NewRectangle(theme.HoverColor())
			bg.Hide()

			icon := widget.NewIcon(nil)
			icon.Resize(fyne.NewSize(32, 32))

			label := widget.NewLabel("")
			label.TextStyle = fyne.TextStyle{Bold: true}

			// Contenedor con fondo para resaltado
			return container.NewStack(
				bg,
				container.NewHBox(icon, label),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(filteredIDs) {
				return
			}

			appID := filteredIDs[id]
			app := appMap[appID]

			// Obtener los elementos de la interfaz
			stack := obj.(*fyne.Container)
			bg := stack.Objects[0].(*canvas.Rectangle)
			contentContainer := stack.Objects[1].(*fyne.Container)

			icon := contentContainer.Objects[0].(*widget.Icon)
			label := contentContainer.Objects[1].(*widget.Label)

			label.SetText(app.Name)
			if iconRes := logic.LoadAppIcon(app); iconRes != nil {
				icon.SetResource(iconRes)
				icon.Show()
			} else {
				icon.Hide()
			}

			// Resaltar elemento seleccionado
			if id == selectedIndex {
				bg.Show()
			} else {
				bg.Hide()
			}
			bg.Refresh()
		},
	)

	// Configurar callbacks para flechas
	input.onKeyDown = func() {
		if selectedIndex < len(filteredIDs)-1 {
			selectedIndex++
			listWidget.Refresh()
			listWidget.ScrollTo(selectedIndex)
		}
	}

	input.onKeyUp = func() {
		if selectedIndex > 0 {
			selectedIndex--
			listWidget.Refresh()
			listWidget.ScrollTo(selectedIndex)
		}
	}

	// Manejar eventos de teclado globales
	handleKeyEvent := func(ev *fyne.KeyEvent) {
		switch ev.Name {
		case fyne.KeyReturn, fyne.KeyEnter:
			executeSelected()
		case fyne.KeyEscape:
			myWindow.Close()
		}
	}

	// Configurar eventos de teclado
	input.OnChanged = func(text string) {
		updateFilter()
	}
	input.OnSubmitted = func(text string) {
		executeSelected()
	}

	// Capturar eventos de teclado a nivel de ventana
	myWindow.Canvas().SetOnTypedKey(handleKeyEvent)

	listWidget.OnSelected = func(id widget.ListItemID) {
		selectedIndex = id
		executeSelected()
	}

	content := container.NewBorder(input, nil, nil, nil, listWidget)
	myWindow.SetContent(content)

	// Enfocar el input después de un breve retraso
	go func() {
		time.Sleep(100 * time.Millisecond)
		myWindow.Canvas().Focus(input)
	}()

	myWindow.ShowAndRun()
}

// getFilteredIDs devuelve los IDs de aplicaciones que coinciden con el filtro
func getFilteredIDs(filter string, apps map[string]models.Application) []string {
	var ids []string
	lowerFilter := strings.ToLower(filter)

	for id, app := range apps {
		if filter == "" || strings.Contains(strings.ToLower(app.Name), lowerFilter) {
			ids = append(ids, id)
		}
	}
	return ids
}
