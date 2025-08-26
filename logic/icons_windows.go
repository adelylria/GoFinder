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

	if app.IconPath != "" {
		cleanPath := strings.Trim(app.IconPath, `"`)
		ext := strings.ToLower(filepath.Ext(cleanPath))
		switch ext {
		case ".png", ".jpg", ".jpeg":
			if res := common.LoadImageFileToResource(cleanPath, app.Name); res != nil {
				common.CacheSet(cacheKey, res)
				return res
			}
		case ".ico":
			if res := common.LoadICOToResource(cleanPath, app.Name); res != nil {
				common.CacheSet(cacheKey, res)
				return res
			}
		}
	}

	if app.IconPath != "" {
		cleanPath := strings.Trim(app.IconPath, `"`)
		if hIcon, err := windows.ExtractIconEx(cleanPath, app.IconIdx); err == nil {
			if res := windows.LoadIconFromHICON(hIcon, app.Name); res != nil {
				common.CacheSet(cacheKey, res)
				return res
			}
		}
	}

	if app.IconPath != "" && app.IconIdx != 0 {
		cleanPath := strings.Trim(app.IconPath, `"`)
		if hIcon, err := windows.ExtractIconEx(cleanPath, 0); err == nil {
			if res := windows.LoadIconFromHICON(hIcon, app.Name); res != nil {
				common.CacheSet(cacheKey, res)
				return res
			}
		}
	}

	if app.Exec != "" {
		cleanExec := strings.Trim(app.Exec, `"`)
		if hIcon, err := windows.ExtractIconEx(cleanExec, 0); err == nil {
			if res := windows.LoadIconFromHICON(hIcon, app.Name); res != nil {
				common.CacheSet(cacheKey, res)
				return res
			}
		}
	}

	if app.Exec != "" {
		cleanExec := strings.Trim(app.Exec, `"`)
		if hIcon, err := windows.SHGetFileIcon(cleanExec); err == nil {
			if res := windows.LoadIconFromHICON(hIcon, app.Name); res != nil {
				common.CacheSet(cacheKey, res)
				return res
			}
		}
	}

	return nil
}
