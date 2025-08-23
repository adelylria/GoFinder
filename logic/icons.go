package logic

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"github.com/adelylria/GoFinder/logic/common"
	"github.com/adelylria/GoFinder/logic/windows"
	"github.com/adelylria/GoFinder/models"
	"github.com/fyne-io/image/ico"
)

// ---- Icon Cache ----
var (
	IconCache = make(map[string]fyne.Resource)
	CacheMu   sync.RWMutex
)

func CacheGet(key string) (fyne.Resource, bool) {
	CacheMu.RLock()
	defer CacheMu.RUnlock()
	r, ok := IconCache[key]
	return r, ok
}

func CacheSet(key string, res fyne.Resource) {
	CacheMu.Lock()
	defer CacheMu.Unlock()
	IconCache[key] = res
}

// ---- Icon Loading Workflow ----
func LoadAppIcon(app models.Application) fyne.Resource {
	// Sistema no-Windows
	if runtime.GOOS != "windows" {
		if app.IconPath == "" {
			return nil
		}
		if cached, ok := CacheGet(app.IconPath); ok {
			return cached
		}
		res := loadImageFileToResource(app.IconPath, app.Name)
		if res != nil {
			CacheSet(app.IconPath, res)
		}
		return res
	}

	// Windows: usar caché con clave única
	cacheKey := fmt.Sprintf("%s|%d", app.IconPath, app.IconIdx)
	if cached, ok := CacheGet(cacheKey); ok {
		return cached
	}

	// 1. Intentar cargar como archivo de imagen
	if app.IconPath != "" {
		cleanPath := strings.Trim(app.IconPath, `"`)
		ext := strings.ToLower(filepath.Ext(cleanPath))
		switch ext {
		case ".png", ".jpg", ".jpeg":
			if res := loadImageFileToResource(cleanPath, app.Name); res != nil {
				CacheSet(cacheKey, res)
				return res
			}
		case ".ico":
			if res := loadICOToResource(cleanPath, app.Name); res != nil {
				CacheSet(cacheKey, res)
				return res
			}
		}
	}

	// 2. Extraer icono de ejecutable/archivo usando el índice especificado
	if app.IconPath != "" {
		cleanPath := strings.Trim(app.IconPath, `"`)

		// Primero con el índice original
		if hIcon, err := windows.ExtractIconEx(cleanPath, app.IconIdx); err == nil {
			if res := windows.LoadIconFromHICON(hIcon, app.Name); res != nil {
				CacheSet(cacheKey, res)
				return res
			}
		}

		// Si falla, intentar con índice 0
		if app.IconIdx != 0 {
			if hIcon, err := windows.ExtractIconEx(cleanPath, 0); err == nil {
				if res := windows.LoadIconFromHICON(hIcon, app.Name); res != nil {
					CacheSet(cacheKey, res)
					return res
				}
			}
		}
	}

	// 3. Intentar con el ejecutable principal
	if app.Exec != "" {
		cleanExec := strings.Trim(app.Exec, `"`)
		if hIcon, err := windows.ExtractIconEx(cleanExec, 0); err == nil {
			if res := windows.LoadIconFromHICON(hIcon, app.Name); res != nil {
				CacheSet(cacheKey, res)
				return res
			}
		}
	}

	// 4. Fallback: obtener icono del sistema
	if app.IconPath != "" {
		cleanPath := strings.Trim(app.IconPath, `"`)
		if hIcon, err := windows.SHGetFileIcon(cleanPath); err == nil {
			if res := windows.LoadIconFromHICON(hIcon, app.Name); res != nil {
				CacheSet(cacheKey, res)
				return res
			}
		}
	}

	// 5. Fallback final: icono del ejecutable usando SHGetFileInfo
	if app.Exec != "" {
		cleanExec := strings.Trim(app.Exec, `"`)
		if hIcon, err := windows.SHGetFileIcon(cleanExec); err == nil {
			if res := windows.LoadIconFromHICON(hIcon, app.Name); res != nil {
				CacheSet(cacheKey, res)
				return res
			}
		}
	}

	return nil
}

// ---- File Loaders ----
func loadImageFileToResource(path, nameHint string) fyne.Resource {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil
	}

	return fyne.NewStaticResource(common.SanitizeResourceName(nameHint)+".png", buf.Bytes())
}

func loadICOToResource(path, nameHint string) fyne.Resource {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	img, err := ico.Decode(file)
	if err != nil {
		return nil
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil
	}

	return fyne.NewStaticResource(common.SanitizeResourceName(nameHint)+".png", buf.Bytes())
}
