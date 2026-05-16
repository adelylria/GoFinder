//go:build !windows

package ui

import "github.com/adelylria/GoFinder/models"

func startSystemTray(state *models.AppState, icon []byte) {}
