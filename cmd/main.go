package main

import (
	"github.com/adelylria/GoFinder/hotkey"
	"github.com/adelylria/GoFinder/logic"
	"github.com/adelylria/GoFinder/ui"
)

func main() {
	hotkey.StartListening()
	defer hotkey.StopListening()

	apps := logic.FindApplications()
	ui.RunLauncher(apps)
}
