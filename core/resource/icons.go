package resource

import (
	// Import required to enable the //go:embed directive below (appIconBytes).
	// The embed package must be imported even if not referenced directly.
	_ "embed"
	"os"
	"sync"

	"fyne.io/fyne/v2"
	"github.com/adelylria/GoFinder/core/logger"
)

var (
	iconName = "GoFinder.ico"

	//go:embed assets/GoFinder.ico
	appIconBytes []byte
	once         sync.Once
	cachedRes    fyne.Resource
)

func GetIcon(iconPath, iconName string) fyne.Resource {
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

func GetEmbedAppIcon() fyne.Resource {
	if len(appIconBytes) == 0 {
		return nil
	}
	return fyne.NewStaticResource(iconName, appIconBytes)
}

func GetEmbedAppIconBytes() []byte {
	if len(appIconBytes) == 0 {
		return nil
	}
	return append([]byte(nil), appIconBytes...)
}
