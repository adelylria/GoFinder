package main

import (
	"github.com/adelylria/GoFinder/logic"
	"github.com/adelylria/GoFinder/ui"
)

func main() {
	apps := logic.FindApplications()
	ui.RunLauncher(apps)
}
