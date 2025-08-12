package models

import (
	"sync"

	"fyne.io/fyne/v2"
	"github.com/google/uuid"
)

type Application struct {
	ID       string // UUID único
	Name     string
	Exec     string
	Icon     string // crudo: path,index o nombre
	IconPath string // expandido
	IconIdx  int
}

type AppState struct {
	Window  fyne.Window
	Visible bool
	Mu      sync.Mutex
}

func NewApplication() Application {
	return Application{
		ID: uuid.New().String(), // Genera un nuevo UUID
	}
}
