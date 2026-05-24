package ui

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/adelylria/GoFinder/core/configuration"
	"github.com/adelylria/GoFinder/core/global"
	"github.com/adelylria/GoFinder/core/i18n"
	"github.com/adelylria/GoFinder/core/resource"
	"github.com/adelylria/GoFinder/core/singleinstance"
	"github.com/adelylria/GoFinder/logic"
	"github.com/adelylria/GoFinder/models"

	hotkey "github.com/adelylria/GoFinder/core/hotkey"
)

// Launcher es el componente principal de la UI del lanzador.
// Ahora incorpora el ThemeConfig (core) para construir los widgets
// con apariencia y tamaños centralizados.
type Launcher struct {
	window         fyne.Window
	input          *hotkey.KeyEventInterceptor
	list           *widget.List
	appMap         map[string]models.Application
	filteredIDs    []string
	selectedIndex  int
	theme          *ThemeConfig
	config         configuration.Config
	startHidden    bool
	hotkeys        *hotkey.HotkeyManager
	dialogsMu      sync.Mutex
	settingsDialog dialog.Dialog
	aboutDialog    dialog.Dialog
	// prevContent stores the previous window content when navigating to settings
	prevContent fyne.CanvasObject
	// settingsOpen indicates whether the settings view is currently shown
	settingsOpen bool
}

// NewLauncher crea el lanzador e inyecta el theme core.
func NewLauncher(apps []models.Application) *Launcher {
	myApp := app.New()

	cfg, err := configuration.Load()
	if err != nil {
		log.Printf("Error cargando configuración: %v", err)
		cfg = configuration.DefaultConfig()
	}
	applyAppTheme(cfg.ThemeName)

	myApp.SetIcon(resource.GetEmbedAppIcon())
	window := myApp.NewWindow("GoFinder")
	if err := configuration.ApplyAutoStart(cfg.AutoStart); err != nil {
		log.Printf("Error configurando inicio automático: %v", err)
	}

	// Theme central
	t := DefaultTheme()
	t.ApplyToWindow(window)

	appMap := createAppMap(apps)

	appState := &models.AppState{
		Window:  window,
		Visible: !cfg.StartHidden,
	}

	hm := hotkey.NewHotkeyManager(
		func() { toggleWindowVisibility(appState) },
		quitApplication,
		toHotkeyBinding(cfg.ToggleHotkey),
		toHotkeyBinding(cfg.QuitHotkey),
	)

	singleinstance.SetActivationHandler(func() {
		setWindowVisible(appState, true)
	})
	startSystemTray(appState, resource.GetEmbedAppIconBytes())

	return &Launcher{
		window:        window,
		appMap:        appMap,
		filteredIDs:   getAllAppIDs(appMap),
		selectedIndex: 0,
		theme:         t,
		config:        cfg,
		startHidden:   cfg.StartHidden,
		hotkeys:       hm,
	}
}

// Inicia y muestra la interfaz de usuario
func (l *Launcher) Run() {
	l.initializeUI()
	if !l.startHidden {
		l.window.Show()
	}
	l.applyNativeMenuPlatformHooks()
	fyne.CurrentApp().Run()
}

// Configura todos los componentes de la interfaz
func (l *Launcher) initializeUI() {
	l.input = l.createInputField()
	l.list = l.createAppList()
	l.setupEventHandlers()

	content := l.theme.NewBorderWithInputTop(l.input, l.list)
	l.window.SetContent(content)
	l.configureNativeMenu()

	l.setupMenuHotkeys()

	if !l.startHidden {
		l.scheduleFocusInput()
	}
}

func (l *Launcher) setupMenuHotkeys() {
	if l.hotkeys == nil {
		return
	}

	l.hotkeys.SetMenuHandlers(l.showSettingsDialog, l.showAboutDialog)
	l.hotkeys.ListenHotkeys()

	l.input.OnMenuQuit = quitApplication
	l.input.OnMenuPrefs = l.showSettingsDialog
	l.input.OnMenuAbout = l.showAboutDialog
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
	log.Printf(i18n.T(i18n.LogRunningApp), app.Name, app.Exec)

	if err := logic.RunApplication(app); err != nil {
		log.Printf(i18n.T(i18n.LogRunAppError), app.Name, err)
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

	setWindowVisible(state, shouldShow)
}

func setWindowVisible(state *models.AppState, visible bool) {
	state.Mu.Lock()
	if state.Visible == visible {
		state.Mu.Unlock()
		return
	}
	state.Visible = visible
	state.Mu.Unlock()

	fyne.Do(func() {
		if visible {
			state.Window.Show()
			state.Window.RequestFocus()
		} else {
			state.Window.Hide()
		}
	})
}

func quitApplication() {
	fmt.Println(i18n.T(i18n.AppExitMessage))
	singleinstance.Release()
	os.Exit(0)
}

func toHotkeyBinding(binding configuration.KeyBinding) hotkey.KeyBinding {
	return hotkey.KeyBinding{
		Modifier: binding.Modifier,
		Key:      binding.Key,
	}
}

func letterOptions() []string {
	letters := make([]string, 0, 26)
	for ch := 'A'; ch <= 'Z'; ch++ {
		letters = append(letters, string(ch))
	}
	return letters
}
