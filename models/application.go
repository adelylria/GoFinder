package models

import (
	"github.com/google/uuid"
)

type Application struct {
	ID       string // UUID único
	Name     string
	Exec     string
	Icon     string
	IconPath string
	IconIdx  int
}

func NewApplication() Application {
	return Application{
		ID: uuid.New().String(),
	}
}
