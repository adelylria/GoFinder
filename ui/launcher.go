package ui

import (
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/adelylria/GoFinder/logic"
	"github.com/adelylria/GoFinder/models"
)

func RunLauncher(apps []models.Application) {
	myApp := app.New()
	myWindow := myApp.NewWindow("Buscador")
	myWindow.SetFixedSize(true)
	myWindow.Resize(fyne.NewSize(600, 500))

	input := widget.NewEntry()
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
		} else {
			//myWindow.Close()
		}
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

	// Manejar eventos de teclado
	handleKeyEvent := func(ev *fyne.KeyEvent) {
		switch ev.Name {
		case fyne.KeyUp:
			if selectedIndex > 0 {
				selectedIndex--
				listWidget.Refresh()
				listWidget.ScrollTo(selectedIndex)
			}
		case fyne.KeyDown:
			if selectedIndex < len(filteredIDs)-1 {
				selectedIndex++
				listWidget.Refresh()
				listWidget.ScrollTo(selectedIndex)
			}
		case fyne.KeyReturn, fyne.KeyEnter:
			executeSelected()
		}
	}

	// Configurar eventos de teclado
	input.OnChanged = func(text string) {
		updateFilter()
	}
	input.OnSubmitted = func(text string) {
		executeSelected()
	}
	myWindow.Canvas().SetOnTypedKey(handleKeyEvent)
	listWidget.OnSelected = func(id widget.ListItemID) {
		selectedIndex = id
		executeSelected()
	}

	content := container.NewBorder(input, nil, nil, nil, listWidget)
	myWindow.SetContent(content)
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
