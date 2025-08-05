package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/adelylria/GoFinder/hotkey"
	"github.com/adelylria/GoFinder/logic"
)

func main() {
	hotkey.StartListening()

	myApp := app.New()
	myWindow := myApp.NewWindow("Buscador")

	input := widget.NewEntry()
	input.SetPlaceHolder("Buscar aplicación...")

	apps := logic.FindApplications()
	filtered := apps // lista filtrada

	// Función para actualizar resultados
	updateList := func(list *widget.List) {
		filtered = nil
		text := strings.ToLower(input.Text)
		for _, app := range apps {
			if strings.Contains(strings.ToLower(app.Name), text) {
				filtered = append(filtered, app)
			}
		}
		list.Refresh()
	}

	// Lista con iconos y nombres
	list := widget.NewList(
		func() int {
			return len(filtered)
		},
		func() fyne.CanvasObject {
			icon := widget.NewIcon(nil)
			icon.Resize(fyne.NewSize(32, 32)) // Tamaño más grande del icono

			label := widget.NewLabel("")
			label.TextStyle = fyne.TextStyle{Bold: true}
			label.Wrapping = fyne.TextWrap(fyne.TextTruncateOff)

			return container.NewHBox(
				icon,
				container.NewVBox(label),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			app := filtered[id]
			items := obj.(*fyne.Container).Objects
			icon := items[0].(*widget.Icon)
			labelContainer := items[1].(*fyne.Container)
			label := labelContainer.Objects[0].(*widget.Label)

			label.SetText(app.Name)

			if iconRes := logic.LoadAppIcon(app); iconRes != nil {
				icon.SetResource(iconRes)
				icon.Show()
			} else {
				icon.Hide()
			}
		},
	)

	input.OnChanged = func(text string) {
		updateList(list)
	}

	content := container.NewBorder(input, nil, nil, nil, list)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 500))
	myWindow.ShowAndRun()
	hotkey.StopListening()
}
