//go:build darwin

package logic

import (
	"errors"

	"fyne.io/fyne/v2"
	"github.com/adelylria/GoFinder/models"
)

func LoadAppIcon(app models.Application) fyne.Resource {
	return nil
}

func RunApplication(app models.Application) error {
	return errors.New("darwin is not supported")
}
