package main

import (
	"github.com/adelylria/GoFinder/core/ui"
	"github.com/adelylria/GoFinder/logic"
)

func main() {
	apps := logic.FindApplications()
	ui.RunLauncher(apps)
}
