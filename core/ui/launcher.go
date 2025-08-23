package ui

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"

	"github.com/adelylria/GoFinder/core/global"
	"github.com/adelylria/GoFinder/core/hotkey"
	"github.com/adelylria/GoFinder/core/resource"
	"github.com/adelylria/GoFinder/logic"
	"github.com/adelylria/GoFinder/logic/windows"
	"github.com/adelylria/GoFinder/models"
)

// Launcher es el componente principal de la UI del lanzador.
// Ahora incorpora el ThemeConfig (core) para construir los widgets
// con apariencia y tamaños centralizados.
type Launcher struct {
	window        fyne.Window
	input         *hotkey.KeyEventInterceptor
	list          *widget.List
	appMap        map[string]models.Application
	filteredIDs   []string
	selectedIndex int
	theme         *ThemeConfig
}

// NewLauncher crea el lanzador e inyecta el theme core.
func NewLauncher(apps []models.Application) *Launcher {
	myApp := app.New()

	myApp.SetIcon(resource.GetEmbedAppIcon())
	window := myApp.NewWindow("GoFinder")

	// Theme central
	t := DefaultTheme()
	t.ApplyToWindow(window)

	appMap := createAppMap(apps)

	appState := &models.AppState{
		Window:  window,
		Visible: true,
	}

	// Configurar el HotkeyManager
	hm := &hotkey.HotkeyManager{
		ToggleHandler: func() {
			toggleWindowVisibility(appState)
		},
		ExitHandler: exitApplication,
	}
	hm.ListenHotkeys()

	return &Launcher{
		window:        window,
		appMap:        appMap,
		filteredIDs:   getAllAppIDs(appMap),
		selectedIndex: 0,
		theme:         t,
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

	content := l.theme.NewBorderWithInputTop(l.input, l.list)
	l.window.SetContent(content)

	l.scheduleFocusInput()
}

// Crea y configura el campo de búsqueda usando el theme
func (l *Launcher) createInputField() *hotkey.KeyEventInterceptor {
	input := l.theme.NewStyledSearchEntry()
	return input
}

// Crea la lista de aplicaciones usando el theme (create/update delegados)
func (l *Launcher) createAppList() *widget.List {
	return l.theme.NewStyledList(
		l.getItemCount,
		func() fyne.CanvasObject { return l.theme.CreateListItemDefault() },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			// seguridad por si hay cambios en filteredIDs
			if id >= len(l.filteredIDs) {
				// limpiar item
				l.theme.UpdateListItemDefault(id, obj, "", nil, false)
				return
			}
			appID := l.filteredIDs[id]
			app := l.appMap[appID]
			res := logic.LoadAppIcon(app)
			selected := (id == l.selectedIndex)
			l.theme.UpdateListItemDefault(id, obj, app.Name, res, selected)
		},
	)
}

// Configura todos los manejadores de eventos
func (l *Launcher) setupEventHandlers() {
	// Configurar navegación con flechas
	l.input.OnKeyDown = l.handleKeyDown
	l.input.OnKeyUp = l.handleKeyUp

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
	go fyne.DoAndWait(func() {
		time.Sleep(global.UIInteractionDelay)
		l.window.Canvas().Focus(l.input)
	})
}

// --- Funciones para el widget.List ---

func (l *Launcher) getItemCount() int {
	return len(l.filteredIDs)
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

	go fyne.Do(func() {
		time.Sleep(global.UIInteractionDelay)
		l.window.Canvas().Focus(l.input)
	})
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

	if err := windows.RunApplication(app); err != nil {
		log.Printf("Error al ejecutar %s: %v", app.Name, err)
	}

	l.clearList()
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

func (l *Launcher) clearList() {
	l.input.SetText("")
	l.filteredIDs = getAllAppIDs(l.appMap)

	go fyne.Do(func() {
		time.Sleep(global.UIInteractionDelay)
		l.list.Unselect(l.selectedIndex)
		l.selectedIndex = 0
		l.list.Refresh()
		l.list.ScrollTo(l.selectedIndex)
	})
}

// Punto de entrada para iniciar el lanzador
func RunLauncher(apps []models.Application) {
	launcher := NewLauncher(apps)
	launcher.Run()
}

func toggleWindowVisibility(state *models.AppState) {
	state.Mu.Lock()
	shouldShow := !state.Visible
	state.Mu.Unlock()

	go fyne.Do(func() {
		if shouldShow {
			state.Window.Show()
		} else {
			state.Window.Hide()
		}
	})

	state.Mu.Lock()
	state.Visible = shouldShow
	state.Mu.Unlock()
}

func exitApplication() {
	fmt.Println("Saliendo...")
	os.Exit(0)
}
