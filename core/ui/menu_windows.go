//go:build windows

package ui

import (
	"sync"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver"

	"github.com/lxn/win"
)

type wndProcInfo struct {
	orig uintptr
	cb   func(h syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr
}

var menuWndProcs sync.Map // win.HWND -> wndProcInfo

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

	var callback func(h syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr
	callback = func(h syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
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

		infoI, ok := menuWndProcs.Load(win.HWND(h))
		if !ok {
			return win.DefWindowProc(win.HWND(h), msg, wParam, lParam)
		}
		info := infoI.(wndProcInfo)
		return win.CallWindowProc(info.orig, win.HWND(h), msg, wParam, lParam)
	}

	proc := syscall.NewCallback(callback)
	// store wndProcInfo (including the Go callback) to keep the function alive
	menuWndProcs.Store(hwnd, wndProcInfo{orig: orig, cb: callback})

	win.SetWindowLongPtr(hwnd, win.GWLP_WNDPROC, proc)
}
