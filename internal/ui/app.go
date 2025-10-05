/*
 * app.go - Simple UI implementation placeholder
 *
 * Copyright (c) 2025 Go port
 * This is public domain software, under the terms of the UNLICENSE
 * http://unlicense.org
 */

package ui

import (
	"fmt"

	"github.com/glreno/oneko-go/internal/config"
)

// NekoApp represents the Neko application
type NekoApp struct {
	settings *config.NekoSettings
}

// NewNekoApp creates a new Neko application
func NewNekoApp(settings *config.NekoSettings) *NekoApp {
	return &NekoApp{
		settings: settings,
	}
}

// Run runs the application
func (app *NekoApp) Run() error {
	title := app.settings.GetTitle()
	if title == "" {
		title = "Neko"
	}

	fmt.Printf("Neko application: %s\n", title)
	fmt.Println("UI implementation TODO - will be replaced with Fyne")
	fmt.Printf("Settings loaded: trigger=%d, catch=%d, run=%d\n",
		app.settings.GetTriggerDist(),
		app.settings.GetCatchDist(),
		app.settings.GetRunDist())

	return nil
}
