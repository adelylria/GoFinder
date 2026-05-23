//go:build linux

package logic

import "github.com/adelylria/GoFinder/logic/ubuntu"

func init() {
	RegisterAppFinder("linux", ubuntu.LinuxAppFinder{})
}
