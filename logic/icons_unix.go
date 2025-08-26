//go:build linux || darwin

package logic

import (
	"fyne.io/fyne/v2"
	"github.com/adelylria/GoFinder/logic/common"
	"github.com/adelylria/GoFinder/models"
)

func LoadAppIcon(app models.Application) fyne.Resource {
	if app.IconPath == "" {
		return nil
	}
	if cached, ok := common.CacheGet(app.IconPath); ok {
		return cached
	}
	res := common.LoadImageFileToResource(app.IconPath, app.Name)
	if res != nil {
		common.CacheSet(app.IconPath, res)
	}
	return res
}
