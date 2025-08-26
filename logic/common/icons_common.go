package common

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"sync"

	"fyne.io/fyne/v2"
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

// ---- Helpers to load image files ----
func LoadImageFileToResource(path, nameHint string) fyne.Resource {
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

	return fyne.NewStaticResource(SanitizeResourceName(nameHint)+".png", buf.Bytes())
}

func LoadICOToResource(path, nameHint string) fyne.Resource {
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

	return fyne.NewStaticResource(SanitizeResourceName(nameHint)+".png", buf.Bytes())
}
