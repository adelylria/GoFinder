package logic

import (
	"fmt"
	"runtime"

	"github.com/adelylria/GoFinder/models"
)

type AppFinder interface {
	Find() []models.Application
}

var appFinders = make(map[string]AppFinder)

func RegisterAppFinder(os string, finder AppFinder) {
	appFinders[os] = finder
}

func FindApplications() []models.Application {
	finder, exists := appFinders[runtime.GOOS]
	if !exists {
		fmt.Printf("Sistema operativo no soportado: %s\n", runtime.GOOS)
		return []models.Application{}
	}
	return finder.Find()
}
