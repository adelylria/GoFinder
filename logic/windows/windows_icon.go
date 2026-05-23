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

	iconInfo, iconCleanup, err := getIconInfoCleanup(hIcon)
	if err != nil {
		return nil, err
	}
	defer iconCleanup()

	width, height := getIconDimensions(iconInfo)

	_, hdcMem, dcCleanup, err := prepareDCs()
	if err != nil {
		return nil, err
	}
	defer dcCleanup()

	bitsPtr, deselect, dibCleanup, err := createDIBSectionAndSelect(hdcMem, width, height)
	if err != nil {
		return nil, err
	}
	defer dibCleanup()
	defer deselect()

	if err := drawIconToHdc(hdcMem, hIcon, width, height); err != nil {
		return nil, err
	}

	img := bitsToNRGBA(bitsPtr, width, height)
	return img, nil
}

func getIconInfoCleanup(hIcon win.HICON) (win.ICONINFO, func(), error) {
	var iconInfo win.ICONINFO
	if !win.GetIconInfo(hIcon, &iconInfo) {
		return iconInfo, nil, errors.New("GetIconInfo falló")
	}
	cleanup := func() {
		if iconInfo.HbmColor != 0 {
			win.DeleteObject(win.HGDIOBJ(iconInfo.HbmColor))
		}
		if iconInfo.HbmMask != 0 {
			win.DeleteObject(win.HGDIOBJ(iconInfo.HbmMask))
		}
	}
	return iconInfo, cleanup, nil
}

func getIconDimensions(iconInfo win.ICONINFO) (int, int) {
	width, height := 32, 32
	if iconInfo.HbmColor != 0 {
		var bmp win.BITMAP
		if win.GetObject(win.HGDIOBJ(iconInfo.HbmColor), unsafe.Sizeof(bmp), unsafe.Pointer(&bmp)) != 0 {
			width = int(bmp.BmWidth)
			height = int(bmp.BmHeight)
		}
	}
	return width, height
}

func prepareDCs() (win.HDC, win.HDC, func(), error) {
	hdcScreen := win.GetDC(0)
	if hdcScreen == 0 {
		return 0, 0, nil, errors.New("GetDC falló")
	}
	hdcMem := win.CreateCompatibleDC(hdcScreen)
	if hdcMem == 0 {
		win.ReleaseDC(0, hdcScreen)
		return 0, 0, nil, errors.New("CreateCompatibleDC falló")
	}
	cleanup := func() {
		win.DeleteDC(hdcMem)
		win.ReleaseDC(0, hdcScreen)
	}
	return hdcScreen, hdcMem, cleanup, nil
}

func createDIBSectionAndSelect(hdcMem win.HDC, width, height int) (unsafe.Pointer, func(), func(), error) {
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
		return nil, nil, nil, fmt.Errorf("CreateDIBSection falló: %v", err)
	}
	cleanup := func() {
		win.DeleteObject(win.HGDIOBJ(hBitmap))
	}
	oldObj := win.SelectObject(hdcMem, win.HGDIOBJ(hBitmap))
	deselect := func() {
		win.SelectObject(hdcMem, oldObj)
	}
	return bitsPtr, deselect, cleanup, nil
}

func drawIconToHdc(hdcMem win.HDC, hIcon win.HICON, width, height int) error {
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
		return errors.New("DrawIconEx falló")
	}
	return nil
}

func bitsToNRGBA(bitsPtr unsafe.Pointer, width, height int) image.Image {
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
	return img
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

	lf, err := lnk.File(path)
	if err != nil {
		fmt.Printf("Error al analizar acceso directo %s: %v\n", path, err)
		return app
	}

	app.Exec = extractExecFromLnk(lf)
	app.Icon = extractIconFromLnk(lf, app.Exec)

	normalizeExec(&app)
	normalizeIcon(&app)

	return app
}

func extractExecFromLnk(lf lnk.LnkFile) string {
	if lf.LinkInfo.LocalBasePath != "" {
		exec := lf.LinkInfo.LocalBasePath
		if lf.LinkInfo.CommonPathSuffix != "" {
			exec = filepath.Join(exec, lf.LinkInfo.CommonPathSuffix)
		}
		return exec
	}
	if lf.StringData.NameString != "" && filepath.IsAbs(lf.StringData.NameString) {
		return lf.StringData.NameString
	}
	return ""
}

func extractIconFromLnk(lf lnk.LnkFile, defaultExec string) string {
	if lf.StringData.IconLocation != "" {
		return lf.StringData.IconLocation
	}
	if lf.LinkInfo.LocalBasePath != "" {
		return defaultExec
	}
	return ""
}

func normalizeExec(app *models.Application) {
	if app.Exec == "" {
		return
	}
	expandedExec := os.ExpandEnv(app.Exec)
	absExec, err := filepath.Abs(expandedExec)
	if err != nil {
		return
	}
	if _, statErr := os.Stat(absExec); os.IsNotExist(statErr) {
		fmt.Printf("Advertencia: La ruta de ejecución %s no existe\n", absExec)
		app.Exec = ""
	} else if os.IsPermission(statErr) {
		fmt.Printf("Advertencia: Sin permisos para acceder a %s\n", absExec)
		app.Exec = ""
	} else {
		app.Exec = absExec
	}
}

func normalizeIcon(app *models.Application) {
	if app.Icon == "" {
		return
	}
	expandedIcon := os.ExpandEnv(app.Icon)
	absIcon, err := filepath.Abs(expandedIcon)
	if err != nil {
		return
	}
	lower := strings.ToLower(absIcon)
	allowDllExe := strings.HasSuffix(lower, ".dll") || strings.HasSuffix(lower, ".exe")
	if _, statErr := os.Stat(absIcon); os.IsNotExist(statErr) && !allowDllExe {
		fmt.Printf("Advertencia: La ruta del icono %s no existe\n", absIcon)
		app.Icon = app.Exec
	} else if os.IsPermission(statErr) {
		fmt.Printf("Advertencia: Sin permisos para acceder al icono %s\n", absIcon)
		app.Icon = app.Exec
	} else {
		app.Icon = absIcon
	}
}
