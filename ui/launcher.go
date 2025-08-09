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

type Launcher struct {
	window        fyne.Window
	input         *KeyEventInterceptor
	list          *widget.List
	appMap        map[string]models.Application
	filteredIDs   []string
	selectedIndex int
}

// KeyEventInterceptor intercepta eventos de teclado para manejar la navegación
type KeyEventInterceptor struct {
	widget.Entry
	onKeyDown func()
	onKeyUp   func()
}

func NewLauncher(apps []models.Application) *Launcher {
	myApp := app.New()
	window := myApp.NewWindow("GoFinder")
	window.SetFixedSize(true)
	window.Resize(fyne.NewSize(600, 500))

	appMap := createAppMap(apps)

	return &Launcher{
		window:        window,
		appMap:        appMap,
		filteredIDs:   getAllAppIDs(appMap),
		selectedIndex: 0,
	}
}

// Inicia y muestra la interfaz de usuario
func (l *Launcher) Run() {
	l.initializeUI()
	l.window.ShowAndRun()
}

// Configura todos los componentes de la interfaz
func (l *Launcher) initializeUI() {
	l.input = l.createInputField()
	l.list = l.createAppList()
	l.setupEventHandlers()

	content := container.NewBorder(l.input, nil, nil, nil, l.list)
	l.window.SetContent(content)

	l.scheduleFocusInput()
}

// Crea y configura el campo de búsqueda
func (l *Launcher) createInputField() *KeyEventInterceptor {
	input := NewKeyEventInterceptor()
	input.SetPlaceHolder("Buscar aplicación...")
	return input
}

// Crea la lista de aplicaciones
func (l *Launcher) createAppList() *widget.List {
	return widget.NewList(
		l.getItemCount,
		l.createListItem,
		l.updateListItem,
	)
}

// Configura todos los manejadores de eventos
func (l *Launcher) setupEventHandlers() {
	// Configurar navegación con flechas
	l.input.onKeyDown = l.handleKeyDown
	l.input.onKeyUp = l.handleKeyUp

	// Eventos de cambio y envío
	l.input.OnChanged = l.handleInputChange
	l.input.OnSubmitted = l.handleInputSubmit

	// Eventos globales de teclado
	l.window.Canvas().SetOnTypedKey(l.handleGlobalKeyEvent)

	// Selección en lista
	l.list.OnSelected = l.handleListSelection
}

// Programa el enfoque en el campo de entrada
func (l *Launcher) scheduleFocusInput() {
	go func() {
		time.Sleep(100 * time.Millisecond)
		l.window.Canvas().Focus(l.input)
	}()
}

// --- Funciones para el widget.List ---

func (l *Launcher) getItemCount() int {
	return len(l.filteredIDs)
}

func (l *Launcher) createListItem() fyne.CanvasObject {
	bg := canvas.NewRectangle(theme.HoverColor())
	bg.Hide()

	icon := widget.NewIcon(nil)
	icon.Resize(fyne.NewSize(32, 32))

	label := widget.NewLabel("")
	label.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewStack(
		bg,
		container.NewHBox(icon, label),
	)
}

func (l *Launcher) updateListItem(id widget.ListItemID, obj fyne.CanvasObject) {
	if id >= len(l.filteredIDs) {
		return
	}

	appID := l.filteredIDs[id]
	app := l.appMap[appID]
	stack := obj.(*fyne.Container)

	bg := stack.Objects[0].(*canvas.Rectangle)
	contentContainer := stack.Objects[1].(*fyne.Container)
	icon := contentContainer.Objects[0].(*widget.Icon)
	label := contentContainer.Objects[1].(*widget.Label)

	label.SetText(app.Name)
	l.setAppIcon(icon, app)
	l.highlightSelectedItem(bg, id)
}

// --- Handlers de eventos ---

func (l *Launcher) handleKeyDown() {
	if l.selectedIndex < len(l.filteredIDs)-1 {
		l.selectedIndex++
		l.list.Refresh()
		l.list.ScrollTo(l.selectedIndex)
	}
}

func (l *Launcher) handleKeyUp() {
	if l.selectedIndex > 0 {
		l.selectedIndex--
		l.list.Refresh()
		l.list.ScrollTo(l.selectedIndex)
	}
}

func (l *Launcher) handleInputChange(text string) {
	l.filteredIDs = getFilteredIDs(text, l.appMap)
	l.selectedIndex = 0
	l.list.Refresh()
}

func (l *Launcher) handleInputSubmit(text string) {
	l.executeSelectedApp()
}

func (l *Launcher) handleGlobalKeyEvent(ev *fyne.KeyEvent) {
	switch ev.Name {
	case fyne.KeyReturn, fyne.KeyEnter:
		l.executeSelectedApp()
	case fyne.KeyEscape:
		l.window.Close()
	}
}

func (l *Launcher) handleListSelection(id widget.ListItemID) {
	l.selectedIndex = id
	l.executeSelectedApp()
}

// --- Funciones de lógica de aplicación ---

func (l *Launcher) executeSelectedApp() {
	if len(l.filteredIDs) == 0 {
		return
	}

	if l.selectedIndex >= len(l.filteredIDs) {
		l.selectedIndex = len(l.filteredIDs) - 1
	}

	appID := l.filteredIDs[l.selectedIndex]
	app := l.appMap[appID]
	log.Printf("Ejecutando: %s (%s)", app.Name, app.Exec)

	if err := logic.RunApplication(app); err != nil {
		log.Printf("Error al ejecutar %s: %v", app.Name, err)
	}
}

// --- Funciones auxiliares ---

func createAppMap(apps []models.Application) map[string]models.Application {
	appMap := make(map[string]models.Application)
	for _, app := range apps {
		appMap[app.ID] = app
	}
	return appMap
}

func getAllAppIDs(appMap map[string]models.Application) []string {
	ids := make([]string, 0, len(appMap))
	for id := range appMap {
		ids = append(ids, id)
	}
	return ids
}

func getFilteredIDs(filter string, appMap map[string]models.Application) []string {
	if filter == "" {
		return getAllAppIDs(appMap)
	}

	var ids []string
	lowerFilter := strings.ToLower(filter)

	for id, app := range appMap {
		if strings.Contains(strings.ToLower(app.Name), lowerFilter) {
			ids = append(ids, id)
		}
	}
	return ids
}

func (l *Launcher) setAppIcon(icon *widget.Icon, app models.Application) {
	if iconRes := logic.LoadAppIcon(app); iconRes != nil {
		icon.SetResource(iconRes)
		icon.Show()
	} else {
		icon.Hide()
	}
}

func (l *Launcher) highlightSelectedItem(bg *canvas.Rectangle, id int) {
	if id == l.selectedIndex {
		bg.Show()
	} else {
		bg.Hide()
	}
	bg.Refresh()
}

// --- Implementación de KeyEventInterceptor ---

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

// Punto de entrada para iniciar el lanzador
func RunLauncher(apps []models.Application) {
	launcher := NewLauncher(apps)
	launcher.Run()
}
