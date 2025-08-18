package logic

import (
	"fmt"

	"github.com/adelylria/GoFinder/models"
)

type darwinAppFinder struct{}

func (f darwinAppFinder) Find() []models.Application {
	fmt.Println("Buscando aplicaciones en macOS/iOS...")
	// Implementación futura
	return []models.Application{}
}
