//go:build windows

package logic

import (
	"fmt"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/adelylria/GoFinder/logic/common"
	"github.com/adelylria/GoFinder/logic/windows"
	"github.com/adelylria/GoFinder/models"
)

func LoadAppIcon(app models.Application) fyne.Resource {
	cacheKey := fmt.Sprintf("%s|%d", app.IconPath, app.IconIdx)
	if cached, ok := common.CacheGet(cacheKey); ok {
		return cached
	}

	if res := loadFromIconPath(app); res != nil {
		common.CacheSet(cacheKey, res)
		return res
	}

	if res := extractFromIconPath(app); res != nil {
		common.CacheSet(cacheKey, res)
		return res
	}

	if res := extractFromExec(app); res != nil {
		common.CacheSet(cacheKey, res)
		return res
	}

	return nil
}

func loadFromIconPath(app models.Application) fyne.Resource {
	if app.IconPath == "" {
		return nil
	}
	cleanPath := strings.Trim(app.IconPath, `"`)
	ext := strings.ToLower(filepath.Ext(cleanPath))
	switch ext {
	case ".png", ".jpg", ".jpeg":
		return common.LoadImageFileToResource(cleanPath, app.Name)
	case ".ico":
		return common.LoadICOToResource(cleanPath, app.Name)
	default:
		return nil
	}
}

func extractIconsFromPath(path string, index int, name string) fyne.Resource {
	if path == "" {
		return nil
	}
	clean := strings.Trim(path, `"`)
	if hIcon, err := windows.ExtractIconEx(clean, index); err == nil {
		return windows.LoadIconFromHICON(hIcon, name)
	}
	return nil
}

func extractFromIconPath(app models.Application) fyne.Resource {
	if app.IconPath == "" {
		return nil
	}
	if res := extractIconsFromPath(app.IconPath, app.IconIdx, app.Name); res != nil {
		return res
	}
	if app.IconIdx != 0 {
		if res := extractIconsFromPath(app.IconPath, 0, app.Name); res != nil {
			return res
		}
	}
	return nil
}

func extractFromExec(app models.Application) fyne.Resource {
	if app.Exec == "" {
		return nil
	}
	cleanExec := strings.Trim(app.Exec, `"`)
	if res := extractIconsFromPath(cleanExec, 0, app.Name); res != nil {
		return res
	}
	if hIcon, err := windows.SHGetFileIcon(cleanExec); err == nil {
		return windows.LoadIconFromHICON(hIcon, app.Name)
	}
	return nil
}
