//go:build windows

package ui

import (
	"sync"
	"syscall"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver"

	"github.com/lxn/win"
)

var menuWndProcs sync.Map // win.HWND -> original WndProc

func (l *Launcher) applyNativeMenuPlatformHooks() {
	disableMenuBarAltFocus(l.window)
}


func disableMenuBarAltFocus(w fyne.Window) {
	native, ok := w.(driver.NativeWindow)
	if !ok {
		return
	}
	native.RunNative(func(ctx any) {
		c := ctx.(driver.WindowsWindowContext)
		hookMenuBarAltFocus(win.HWND(c.HWND))
	})
}

func hookMenuBarAltFocus(hwnd win.HWND) {
	if _, exists := menuWndProcs.Load(hwnd); exists {
		return
	}

	orig := win.GetWindowLongPtr(hwnd, win.GWLP_WNDPROC)
	menuWndProcs.Store(hwnd, orig)

	proc := syscall.NewCallback(func(h syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
		switch msg {
		case win.WM_SYSCOMMAND:
			if wParam&0xFFF0 == win.SC_KEYMENU {
				return 0
			}
		case win.WM_SYSKEYDOWN, win.WM_SYSKEYUP:
			if wParam == win.VK_MENU {
				return 0
			}
		case win.WM_KEYDOWN:
			if wParam == win.VK_F10 {
				return 0
			}
		}

		origProc, ok := menuWndProcs.Load(win.HWND(h))
		if !ok {
			return win.DefWindowProc(win.HWND(h), msg, wParam, lParam)
		}
		return win.CallWindowProc(origProc.(uintptr), win.HWND(h), msg, wParam, lParam)
	})

	win.SetWindowLongPtr(hwnd, win.GWLP_WNDPROC, proc)
	_ = unsafe.Pointer(proc)
}
