package resource

import (
	"os"
	"sync"

	"fyne.io/fyne/v2"
	"github.com/adelylria/GoFinder/core/logger"
)

var (
	iconPath = "assets/GoFinder.ico"
	iconName = "GoFinder.ico"

	once      sync.Once
	cachedRes fyne.Resource
)

func GetAppIcon() fyne.Resource {
	once.Do(func() {
		data, err := os.ReadFile(iconPath)
		if err != nil {
			logger.LoggerErr.Println("No se pudo leer el icono:", err)
			return
		}
		cachedRes = fyne.NewStaticResource(iconName, data)
	})
	return cachedRes
}
