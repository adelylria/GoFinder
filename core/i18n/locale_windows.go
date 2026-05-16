//go:build windows

package i18n

import (
	"syscall"
	"unsafe"
)

const localeNameMaxLength = 85

var (
	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	procGetUserDefaultLocaleName = kernel32.NewProc("GetUserDefaultLocaleName")
)

func systemLocale() string {
	buffer := make([]uint16, localeNameMaxLength)
	ret, _, _ := procGetUserDefaultLocaleName.Call(
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
	)
	if ret == 0 {
		return ""
	}
	return syscall.UTF16ToString(buffer)
}
