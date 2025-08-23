package logic

import (
	"fmt"

	"github.com/adelylria/GoFinder/models"
)

type DarwinAppFinder struct{}

func (f DarwinAppFinder) Find() []models.Application {
	fmt.Println("Buscando aplicaciones en macOS/iOS...")
	// Implementación futura
	return []models.Application{}
}
