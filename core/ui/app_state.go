package ui

import (
	"sync"

	"fyne.io/fyne/v2"
)

type AppState struct {
	Window  fyne.Window
	Visible bool
	Mu      sync.Mutex
}
