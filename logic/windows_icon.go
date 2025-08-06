package logic

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"fyne.io/fyne/v2"
	"github.com/adelylria/GoFinder/models"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

// Cache global de iconos
var iconCache = make(map[string]fyne.Resource)

var (
	// DLLs y procedimientos para manejo de iconos en Windows
	shell32           = windows.NewLazySystemDLL("shell32.dll")
	extractIconExProc = shell32.NewProc("ExtractIconExW")
	gdi32             = windows.NewLazySystemDLL("gdi32.dll")
	procGetDIBits     = gdi32.NewProc("GetDIBits")
)

// ExtractIcon extrae un icono de un archivo en Windows
func ExtractIcon(path string, index int) (win.HICON, error) {
	var largeIcon, smallIcon win.HICON

	pPath, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}

	ret, _, err := extractIconExProc.Call(
		uintptr(unsafe.Pointer(pPath)),
		uintptr(index),
		uintptr(unsafe.Pointer(&largeIcon)),
		uintptr(unsafe.Pointer(&smallIcon)),
		uintptr(1),
	)
	if ret == 0 {
		return 0, errors.New("no se pudo extraer el icono")
	}
	return largeIcon, nil
}

// SplitIconLocation separa la ruta del icono y su índice
func SplitIconLocation(iconLocation string) (string, int) {
	parts := strings.Split(iconLocation, ",")
	iconPath := parts[0]
	iconIndex := 0
	if len(parts) > 1 {
		i, err := strconv.Atoi(parts[1])
		if err == nil {
			iconIndex = i
		}
	}
	return iconPath, iconIndex
}

// LoadAppIcon carga el icono de una aplicación, usando caché si está disponible
func LoadAppIcon(app models.Application) fyne.Resource {
	// Verificar si ya está en caché
	if cached, ok := iconCache[app.ID]; ok {
		return cached
	}

	var resource fyne.Resource

	// Intentar cargar desde archivo de imagen
	ext := strings.ToLower(filepath.Ext(app.IconPath))
	if ext == ".ico" || ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
		file, err := os.Open(app.IconPath)
		if err == nil {
			defer file.Close()
			img, _, err := image.Decode(file)
			if err == nil {
				var buf bytes.Buffer
				if err := png.Encode(&buf, img); err == nil {
					resource = fyne.NewStaticResource(app.Name+".png", buf.Bytes())
				}
			}
		}
	}

	// Para Windows: extraer icono de archivos EXE, DLL, etc.
	if resource == nil && runtime.GOOS == "windows" && app.IconPath != "" {
		hIcon, err := ExtractIcon(app.IconPath, app.IconIdx)
		if err == nil && hIcon != 0 {
			resource = LoadIconFromHICON(hIcon, app.Name)
		}
	}

	// Almacenar en caché si se encontró un recurso válido
	if resource != nil {
		iconCache[app.ID] = resource
	}

	return resource
}

// IconHandleToImage convierte un HICON de Windows en una imagen Go
func IconHandleToImage(hIcon win.HICON) (image.Image, error) {
	var iconInfo win.ICONINFO
	if !win.GetIconInfo(hIcon, &iconInfo) {
		return nil, errors.New("error al obtener información del icono")
	}
	defer win.DeleteObject(win.HGDIOBJ(iconInfo.HbmColor))
	defer win.DeleteObject(win.HGDIOBJ(iconInfo.HbmMask))

	var bmpInfo win.BITMAP
	if win.GetObject(win.HGDIOBJ(iconInfo.HbmColor), unsafe.Sizeof(bmpInfo), unsafe.Pointer(&bmpInfo)) == 0 {
		return nil, errors.New("error al obtener información del bitmap")
	}
	width := int(bmpInfo.BmWidth)
	height := int(bmpInfo.BmHeight)

	hdc := win.GetDC(0)
	defer win.ReleaseDC(0, hdc)

	hdcMem := win.CreateCompatibleDC(hdc)
	defer win.DeleteDC(hdcMem)

	oldBmp := win.SelectObject(hdcMem, win.HGDIOBJ(iconInfo.HbmColor))
	defer win.SelectObject(hdcMem, oldBmp)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	var bi win.BITMAPINFOHEADER
	bi.BiSize = uint32(unsafe.Sizeof(bi))
	bi.BiWidth = int32(width)
	bi.BiHeight = int32(-height) // Negativo para top-down
	bi.BiPlanes = 1
	bi.BiBitCount = 32
	bi.BiCompression = win.BI_RGB

	bufSize := width * height * 4
	buf := make([]byte, bufSize)

	ret, _, _ := procGetDIBits.Call(
		uintptr(hdcMem),
		uintptr(iconInfo.HbmColor),
		0,
		uintptr(height),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bi)),
		win.DIB_RGB_COLORS,
	)
	if ret == 0 {
		return nil, errors.New("error al obtener bits de imagen")
	}

	// Copiar datos de píxeles a la imagen
	for y := range height {
		for x := range width {
			i := (y*width + x) * 4
			b := buf[i]
			g := buf[i+1]
			r := buf[i+2]
			a := buf[i+3]
			img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: a})
		}
	}

	return img, nil
}

// LoadIconFromHICON convierte un HICON en un recurso Fyne
func LoadIconFromHICON(hIcon win.HICON, appName string) fyne.Resource {
	img, err := IconHandleToImage(hIcon)
	if err != nil {
		return nil
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil
	}

	return fyne.NewStaticResource(appName+".png", buf.Bytes())
}
