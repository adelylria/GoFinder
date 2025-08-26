//go:build windows

package logic

import (
	"github.com/adelylria/GoFinder/logic/windows"
)

func init() {
	RegisterAppFinder("windows", windows.WindowsAppFinder{})
}
