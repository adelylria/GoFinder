//go:build windows
// +build windows

package windows

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"fyne.io/fyne/v2"
	"github.com/adelylria/GoFinder/logic/common"
	"github.com/adelylria/GoFinder/models"
	"github.com/lxn/win"
	lnk "github.com/parsiya/golnk"
	"golang.org/x/sys/windows"
)

var (
	kernel32             = windows.NewLazySystemDLL("kernel32.dll")
	user32               = windows.NewLazySystemDLL("user32.dll")
	gdi32                = windows.NewLazySystemDLL("gdi32.dll")
	shell32              = windows.NewLazySystemDLL("shell32.dll")
	procDrawIconEx       = user32.NewProc("DrawIconEx")
	procCreateDIBSection = gdi32.NewProc("CreateDIBSection")
	procExtractIconExW   = shell32.NewProc("ExtractIconExW")
)

// ExtractIconExW
func ExtractIconEx(path string, index int) (win.HICON, error) {
	if path == "" {
		return 0, errors.New("ruta vacía")
	}

	cleanPath := strings.Trim(path, `"`)

	fullPath := cleanPath
	if !strings.HasPrefix(cleanPath, `\\?\`) && len(cleanPath) > 260 {
		fullPath = `\\?\` + cleanPath
	}

	pPath, err := syscall.UTF16PtrFromString(fullPath)
	if err != nil {
		return 0, err
	}

	var largeIcon, smallIcon win.HICON
	ret, _, _ := procExtractIconExW.Call(
		uintptr(unsafe.Pointer(pPath)),
		uintptr(index),
		uintptr(unsafe.Pointer(&largeIcon)),
		uintptr(unsafe.Pointer(&smallIcon)),
		1,
	)

	if ret <= 0 || (largeIcon == 0 && smallIcon == 0) {
		return 0, fmt.Errorf("ExtractIconExW falló para %s (index %d)", cleanPath, index)
	}

	if largeIcon != 0 {
		if smallIcon != 0 {
			win.DestroyIcon(smallIcon)
		}
		return largeIcon, nil
	}
	return smallIcon, nil
}

func SHGetFileIcon(path string) (win.HICON, error) {
	cleanPath := strings.Trim(path, `"`)
	if cleanPath == "" {
		return 0, errors.New("ruta vacía")
	}

	fullPath := cleanPath
	if !strings.HasPrefix(cleanPath, `\\?\`) && len(cleanPath) > 260 {
		fullPath = `\\?\` + cleanPath
	}

	pPath, err := syscall.UTF16PtrFromString(fullPath)
	if err != nil {
		return 0, err
	}

	var shfi win.SHFILEINFO
	flags := uint32(win.SHGFI_ICON | win.SHGFI_LARGEICON | win.SHGFI_USEFILEATTRIBUTES)
	ret := win.SHGetFileInfo(pPath, 0x80, &shfi, uint32(unsafe.Sizeof(shfi)), flags)

	if ret == 0 || shfi.HIcon == 0 {
		return 0, fmt.Errorf("SHGetFileInfo no devolvió icono para %s", cleanPath)
	}
	return shfi.HIcon, nil
}

// ---- HICON to Image Conversion ----
func IconHandleToImage(hIcon win.HICON) (image.Image, error) {
	if hIcon == 0 {
		return nil, errors.New("hIcon inválido")
	}

	var iconInfo win.ICONINFO
	if !win.GetIconInfo(hIcon, &iconInfo) {
		return nil, errors.New("GetIconInfo falló")
	}
	defer func() {
		if iconInfo.HbmColor != 0 {
			win.DeleteObject(win.HGDIOBJ(iconInfo.HbmColor))
		}
		if iconInfo.HbmMask != 0 {
			win.DeleteObject(win.HGDIOBJ(iconInfo.HbmMask))
		}
	}()

	width, height := 32, 32
	if iconInfo.HbmColor != 0 {
		var bmp win.BITMAP
		if win.GetObject(win.HGDIOBJ(iconInfo.HbmColor), unsafe.Sizeof(bmp), unsafe.Pointer(&bmp)) != 0 {
			width = int(bmp.BmWidth)
			height = int(bmp.BmHeight)
		}
	}

	hdcScreen := win.GetDC(0)
	if hdcScreen == 0 {
		return nil, errors.New("GetDC falló")
	}
	defer win.ReleaseDC(0, hdcScreen)

	hdcMem := win.CreateCompatibleDC(hdcScreen)
	if hdcMem == 0 {
		return nil, errors.New("CreateCompatibleDC falló")
	}
	defer win.DeleteDC(hdcMem)

	var bi win.BITMAPINFO
	bi.BmiHeader = win.BITMAPINFOHEADER{
		BiSize:        uint32(unsafe.Sizeof(win.BITMAPINFOHEADER{})),
		BiWidth:       int32(width),
		BiHeight:      int32(-height),
		BiPlanes:      1,
		BiBitCount:    32,
		BiCompression: win.BI_RGB,
	}

	var bitsPtr unsafe.Pointer
	hBitmap, _, err := procCreateDIBSection.Call(
		uintptr(hdcMem),
		uintptr(unsafe.Pointer(&bi)),
		uintptr(win.DIB_RGB_COLORS),
		uintptr(unsafe.Pointer(&bitsPtr)),
		0,
		0,
	)
	if hBitmap == 0 {
		return nil, fmt.Errorf("CreateDIBSection falló: %v", err)
	}
	defer win.DeleteObject(win.HGDIOBJ(hBitmap))

	oldObj := win.SelectObject(hdcMem, win.HGDIOBJ(hBitmap))
	defer win.SelectObject(hdcMem, oldObj)

	// Dibujar icono
	ret, _, _ := procDrawIconEx.Call(
		uintptr(hdcMem),
		0,
		0,
		uintptr(hIcon),
		uintptr(width),
		uintptr(height),
		0,
		0,
		uintptr(win.DI_NORMAL),
	)
	if ret == 0 {
		return nil, errors.New("DrawIconEx falló")
	}

	// Convertir a imagen Go
	byteLen := width * height * 4
	raw := unsafe.Slice((*byte)(bitsPtr), byteLen)
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := (y*width + x) * 4
			b := raw[i]
			g := raw[i+1]
			r := raw[i+2]
			a := raw[i+3]
			img.SetNRGBA(x, y, color.NRGBA{R: r, G: g, B: b, A: a})
		}
	}

	return img, nil
}

// ---- Resource Helpers ----
func LoadIconFromHICON(hIcon win.HICON, nameHint string) fyne.Resource {
	if hIcon == 0 {
		return nil
	}

	img, err := IconHandleToImage(hIcon)
	win.DestroyIcon(hIcon)
	if err != nil {
		return nil
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil
	}

	resName := common.SanitizeResourceName(nameHint)
	return fyne.NewStaticResource(resName+".png", buf.Bytes())
}

// resolveWindowsShortcut resuelve un acceso directo de Windows
func resolveWindowsShortcut(path string) models.Application {
	app := models.NewApplication()
	app.Name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	lnk, err := lnk.File(path)
	if err != nil {
		fmt.Printf("Error al analizar acceso directo %s: %v\n", path, err)
		return app
	}

	// Obtener ruta de destino
	if lnk.LinkInfo.LocalBasePath != "" {
		app.Exec = lnk.LinkInfo.LocalBasePath
		if lnk.LinkInfo.CommonPathSuffix != "" {
			app.Exec = filepath.Join(app.Exec, lnk.LinkInfo.CommonPathSuffix)
		}
	} else if lnk.StringData.NameString != "" {
		if filepath.IsAbs(lnk.StringData.NameString) {
			app.Exec = lnk.StringData.NameString
		}
	}

	// Obtener ubicación del icono
	if lnk.StringData.IconLocation != "" {
		app.Icon = lnk.StringData.IconLocation
	} else if lnk.LinkInfo.LocalBasePath != "" {
		app.Icon = app.Exec
	}

	// Normalizar y validar rutas
	if app.Exec != "" {
		expandedExec := os.ExpandEnv(app.Exec)
		expandedExec, err = filepath.Abs(expandedExec)
		if err == nil {
			if _, err := os.Stat(expandedExec); os.IsNotExist(err) {
				fmt.Printf("Advertencia: La ruta de ejecución %s no existe\n", expandedExec)
				app.Exec = ""
			} else if os.IsPermission(err) {
				fmt.Printf("Advertencia: Sin permisos para acceder a %s\n", expandedExec)
				app.Exec = ""
			} else {
				app.Exec = expandedExec
			}
		}
	}

	if app.Icon != "" {
		expandedIcon := os.ExpandEnv(app.Icon)
		expandedIcon, err = filepath.Abs(expandedIcon)
		if err == nil {
			// Permitir .dll y .exe para íconos
			if _, err := os.Stat(expandedIcon); os.IsNotExist(err) && !strings.HasSuffix(strings.ToLower(expandedIcon), ".dll") && !strings.HasSuffix(strings.ToLower(expandedIcon), ".exe") {
				fmt.Printf("Advertencia: La ruta del icono %s no existe\n", expandedIcon)
				app.Icon = app.Exec // Usar ejecutable como respaldo
			} else if os.IsPermission(err) {
				fmt.Printf("Advertencia: Sin permisos para acceder al icono %s\n", expandedIcon)
				app.Icon = app.Exec
			} else {
				app.Icon = expandedIcon
			}
		}
	}

	return app
}
