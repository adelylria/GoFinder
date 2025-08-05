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

var (
	// Para ExtractIconExW
	shell32           = windows.NewLazySystemDLL("shell32.dll")
	extractIconExProc = shell32.NewProc("ExtractIconExW")
	// Para GetDIBits
	gdi32         = windows.NewLazySystemDLL("gdi32.dll")
	procGetDIBits = gdi32.NewProc("GetDIBits")
)

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
		return 0, errors.New("no icon extracted")
	}
	return largeIcon, nil
}

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

func LoadAppIcon(app models.Application) fyne.Resource {
	ext := strings.ToLower(filepath.Ext(app.IconPath))
	if ext == ".ico" || ext == ".png" {
		file, err := os.Open(app.IconPath)
		if err == nil {
			defer file.Close()
			img, _, err := image.Decode(file)
			if err == nil {
				var buf bytes.Buffer
				png.Encode(&buf, img)
				return fyne.NewStaticResource(app.Name+".png", buf.Bytes())
			}
		}
	}

	if runtime.GOOS == "windows" && app.IconPath != "" {
		hIcon, err := ExtractIcon(app.IconPath, app.IconIdx)
		if err == nil && hIcon != 0 {
			return LoadIconFromHICON(hIcon, app.Name)
		}
	}

	return nil
}

func IconHandleToImage(hIcon win.HICON) (image.Image, error) {
	var iconInfo win.ICONINFO
	ok := win.GetIconInfo(hIcon, &iconInfo)
	if !ok {
		return nil, errors.New("GetIconInfo failed")
	}
	defer win.DeleteObject(win.HGDIOBJ(iconInfo.HbmColor))
	defer win.DeleteObject(win.HGDIOBJ(iconInfo.HbmMask))

	var bmpInfo win.BITMAP
	ret := win.GetObject(win.HGDIOBJ(iconInfo.HbmColor), unsafe.Sizeof(bmpInfo), unsafe.Pointer(&bmpInfo))
	if ret == 0 {
		return nil, errors.New("GetObject failed")
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
	bi.BiHeight = int32(-height)
	bi.BiPlanes = 1
	bi.BiBitCount = 32
	bi.BiCompression = win.BI_RGB

	bufSize := width * height * 4
	buf := make([]byte, bufSize)

	ret2, _, _ := procGetDIBits.Call(
		uintptr(hdcMem),
		uintptr(iconInfo.HbmColor),
		0,
		uintptr(height),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bi)),
		win.DIB_RGB_COLORS,
	)
	if ret2 == 0 {
		return nil, errors.New("GetDIBits failed")
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := (y*width + x) * 4
			b := buf[i+0]
			g := buf[i+1]
			r := buf[i+2]
			a := buf[i+3]
			img.SetRGBA(x, y, color.RGBA{r, g, b, a})
		}
	}

	return img, nil
}

func LoadIconFromHICON(hIcon win.HICON, appName string) fyne.Resource {
	img, err := IconHandleToImage(hIcon)
	if err != nil {
		return nil
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil
	}

	return fyne.NewStaticResource(appName+".png", buf.Bytes())
}
