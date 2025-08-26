//go:build linux || darwin

package logic

import (
	"github.com/adelylria/GoFinder/logic/ubuntu"
)

func init() {
	// Linux
	RegisterAppFinder("linux", ubuntu.LinuxAppFinder{})

	// Darwin (puede usar el mismo finder o uno propio más adelante)
	RegisterAppFinder("darwin", ubuntu.LinuxAppFinder{})
}
