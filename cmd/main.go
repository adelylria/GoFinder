package main

import (
	"github.com/adelylria/GoFinder/core/singleinstance"
	"github.com/adelylria/GoFinder/core/ui"
	"github.com/adelylria/GoFinder/logic"
)

func main() {
	if !singleinstance.EnsureFirstInstance() {
		return
	}
	defer singleinstance.Release()

	apps := logic.FindApplications()
	ui.RunLauncher(apps)
}
