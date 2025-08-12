package resource

import (
	"os"

	"github.com/adelylria/GoFinder/core/logger"
)

func GetAppIcon() []byte {
	iconPath := "assets/GoFinder.ico"
	iconData, err := os.ReadFile(iconPath)
	if err != nil {
		logger.LoggerErr.Println("No se pudo leer el icono: ", err)
		return nil
	}
	return iconData
}
