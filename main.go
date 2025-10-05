/*
 * main.go - Neko the cat
 *
 * Copyright (c) 2019 Jerry Reno (original Java version)
 * Copyright (c) 2010 Werner Randelshofer (original Java port)
 * Copyright (c) 2025 Go port
 * This is public domain software, under the terms of the UNLICENSE
 * http://unlicense.org
 *
 * This program loads 32 images of Neko and animates them.
 * Neko will chase your mouse cursor around the desktop.
 * Once she's over it and the mouse doesn't move, she'll prepare to take a nap.
 * If the mouse goes outside the desktop, she will reach the border and try
 * to dig for it. She'll eventually give up and fall asleep.
 */

package main

import (
	"log"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	windowWidth  = 64 * 16
	windowHeight = 64 * 9
	iconWidth    = 64
	iconHeight   = 64
)

// Neko is the main application struct
type Neko struct {
	settings   *NekoSettings
	controller *NekoController
	window     *glfw.Window
	windowMode bool
}

// NewNeko creates a new Neko instance
func NewNeko() *Neko {
	n := &Neko{
		settings:   NewNekoSettings(),
		windowMode: false,
	}

	n.initComponents()
	return n
}

// initComponents initializes the window and controller
func (n *Neko) initComponents() {
	// Initialize GLFW
	if err := glfw.Init(); err != nil {
		log.Fatal("Failed to initialize GLFW:", err)
	}

	// Set window hints
	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.Decorated, glfw.False)
	glfw.WindowHint(glfw.Floating, glfw.True)
	glfw.WindowHint(glfw.TransparentFramebuffer, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.False)

	// Get title from settings
	title := n.settings.GetTitle()
	if title == "" {
		title = "Neko"
	}

	// Create window
	var err error
	n.window, err = glfw.CreateWindow(iconWidth, iconHeight, title, nil, nil)
	if err != nil {
		log.Fatal("Failed to create window:", err)
	}

	// Set window properties
	n.window.SetPos(100, 100)

	// Make the window's context current
	n.window.MakeContextCurrent()

	// Set callbacks
	n.window.SetMouseButtonCallback(n.mouseButtonCallback)
	n.window.SetPosCallback(n.windowPosCallback)
	n.window.SetIconifyCallback(n.windowIconifyCallback)

	// Create controller
	n.controller = NewNekoController(n.settings, n.window)

	// Start in free mode
	n.SetWindowMode(false)
}

// mouseButtonCallback handles mouse clicks
func (n *Neko) mouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		n.SetWindowMode(!n.windowMode)
	}
}

// windowPosCallback handles window position changes
func (n *Neko) windowPosCallback(w *glfw.Window, xpos int, ypos int) {
	if n.controller != nil {
		n.controller.CatboxMoved()
	}
}

// windowIconifyCallback handles window iconify/restore
func (n *Neko) windowIconifyCallback(w *glfw.Window, iconified bool) {
	if !iconified && n.controller != nil {
		n.controller.CatboxDeiconified()
	}
}

// SetWindowMode sets the window mode (windowed vs free-floating)
func (n *Neko) SetWindowMode(windowed bool) {
	n.windowMode = windowed
	n.controller.SetWindowMode(windowed)
	n.settings.Load()

	if windowed {
		// Windowed mode - show a bordered window
		title := n.settings.GetTitle()
		if title == "" {
			title = "Neko"
		}
		n.window.SetTitle(title)
		n.window.SetSize(windowWidth, windowHeight)
		glfw.WindowHint(glfw.Decorated, glfw.True)
		n.window.Show()
	} else {
		// Free mode - borderless, transparent window
		n.window.SetSize(iconWidth, iconHeight)
		glfw.WindowHint(glfw.Decorated, glfw.False)
		n.window.Show()
	}

	n.controller.MoveCatInBox()
}

// Run runs the main application loop
func (n *Neko) Run() {
	// Show the window
	n.window.Show()

	// Main loop
	for !n.window.ShouldClose() {
		// Poll events
		glfw.PollEvents()

		// Clear the screen
		// In a full implementation, you would render the current
		// animation frame here using OpenGL

		// Swap buffers
		n.window.SwapBuffers()
	}

	// Cleanup
	n.controller.Stop()
	n.window.Destroy()
	glfw.Terminate()
}

func main() {
	// Lock OS thread for GLFW
	runtime.LockOSThread()

	// Create and run Neko
	neko := NewNeko()
	neko.Run()
}
