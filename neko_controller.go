/*
 * neko_controller.go
 *
 * Copyright (c) 2019 Jerry Reno (original Java version)
 * Copyright (c) 2025 Go port
 * This is public domain software, under the terms of the UNLICENSE
 * http://unlicense.org
 */

package main

import (
	"fmt"
	"image"
	_ "image/gif"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

// Position enum
type Position int

const (
	PositionOver Position = iota
	PositionUnder
	PositionLeft
	PositionRight
)

// NekoController controls the Neko cat animation and behavior
type NekoController struct {
	settings *NekoSettings

	// State variables
	windowMode   bool
	ox, oy       int           // image position
	no           int           // image number
	init         int           // initialization counter
	state        int           // current state
	slp          time.Duration // sleep time
	ilc1, ilc2   int           // image loop counters
	mouseMoved   bool
	windowOffset image.Point
	nekoBounds   image.Rectangle
	w, h         int // size of icons

	// Images
	images []*image.Image

	// Timer
	ticker *time.Ticker
	done   chan bool

	// UI references
	window *glfw.Window

	// X11 connection for getting mouse position
	xconn  *xgb.Conn
	screen *xproto.ScreenInfo
}

// NewNekoController creates a new NekoController
func NewNekoController(settings *NekoSettings, window *glfw.Window) *NekoController {
	nc := &NekoController{
		settings:     settings,
		window:       window,
		windowMode:   false,
		init:         0,
		state:        0,
		slp:          0,
		windowOffset: image.Pt(-16, -30),
		done:         make(chan bool),
	}

	nc.loadKitten()
	nc.w = 32 // Default icon width
	nc.h = 32 // Default icon height

	// Initialize X11 connection for mouse tracking
	nc.initX11()

	// Start animation timer
	nc.ticker = time.NewTicker(50 * time.Millisecond)
	go nc.animationLoop()

	return nc
}

// initX11 initializes X11 connection for mouse position tracking
func (nc *NekoController) initX11() {
	var err error
	nc.xconn, err = xgb.NewConn()
	if err != nil {
		return
	}

	setup := xproto.Setup(nc.xconn)
	if setup != nil && len(setup.Roots) > 0 {
		nc.screen = &setup.Roots[0]
	}
}

// loadKitten loads all 32 cat animation images
func (nc *NekoController) loadKitten() {
	nc.images = make([]*image.Image, 33)

	// Load images 1-32
	for i := 1; i <= 32; i++ {
		imgPath := filepath.Join("resources", "images", fmt.Sprintf("%d.GIF", i))
		file, err := os.Open(imgPath)
		if err != nil {
			continue
		}

		img, _, err := image.Decode(file)
		file.Close()

		if err == nil {
			nc.images[i] = &img
		}
	}

	// Image 0 is a copy of image 25
	if nc.images[25] != nil {
		nc.images[0] = nc.images[25]
	}
}

// GetWidth returns the icon width
func (nc *NekoController) GetWidth() int {
	return nc.w
}

// GetHeight returns the icon height
func (nc *NekoController) GetHeight() int {
	return nc.h
}

// GetWindowMode returns whether the cat is in window mode
func (nc *NekoController) GetWindowMode() bool {
	return nc.windowMode
}

// SetWindowMode sets the window mode
func (nc *NekoController) SetWindowMode(windowed bool) {
	nc.windowMode = windowed
}

// getMousePosition gets the current mouse position using X11
func (nc *NekoController) getMousePosition() (int, int, bool) {
	if nc.xconn == nil || nc.screen == nil {
		return 0, 0, false
	}

	reply, err := xproto.QueryPointer(nc.xconn, nc.screen.Root).Reply()
	if err != nil {
		return 0, 0, false
	}

	return int(reply.RootX), int(reply.RootY), true
}

// calculateBounds calculates the valid area where Neko can move
func (nc *NekoController) calculateBounds() {
	if nc.windowMode {
		// In window mode, bounds are relative to the window
		winW, winH := nc.window.GetSize()
		nc.nekoBounds = image.Rect(
			nc.w/2, nc.h,
			winW-nc.w, winH-nc.h,
		)
	} else {
		// In free mode, use screen bounds
		// For simplicity, using a fixed screen size
		// In production, you'd query the actual screen size
		screenWidth := 1920
		screenHeight := 1080

		nc.nekoBounds = image.Rect(
			nc.w/2, nc.h,
			screenWidth-nc.w, screenHeight-nc.h,
		)
	}
}

// animationLoop is the main animation loop
func (nc *NekoController) animationLoop() {
	for {
		select {
		case <-nc.done:
			return
		case <-nc.ticker.C:
			nc.locateMouseAndAnimateCat()
		}
	}
}

// locateMouseAndAnimateCat locates the mouse and determines cat's action
func (nc *NekoController) locateMouseAndAnimateCat() {
	mx, my, ok := nc.getMousePosition()
	if !ok {
		return
	}

	mx += nc.settings.GetOffsetX()
	my += nc.settings.GetOffsetY()

	nc.calculateBounds()

	// Determine what the cat should do if mouse moves
	var pos *Position
	out := !(nc.nekoBounds.Min.X <= mx && mx <= nc.nekoBounds.Max.X &&
		nc.nekoBounds.Min.Y <= my && my <= nc.nekoBounds.Max.Y)

	x, y := mx, my
	if out {
		if y < nc.nekoBounds.Min.Y {
			y = nc.nekoBounds.Min.Y
			p := PositionOver
			pos = &p
		}
		if y > nc.nekoBounds.Max.Y {
			y = nc.nekoBounds.Max.Y
			p := PositionUnder
			pos = &p
		}
		if x < nc.nekoBounds.Min.X {
			x = nc.nekoBounds.Min.X
			p := PositionLeft
			pos = &p
		}
		if x > nc.nekoBounds.Max.X {
			x = nc.nekoBounds.Max.X
			p := PositionRight
			pos = &p
		}
	}

	// Calculate distance and angle
	dx := x - nc.ox
	dy := nc.oy - y
	dist := math.Sqrt(float64(dx*dx + dy*dy))
	theta := math.Atan2(float64(dy), float64(dx))

	// Determine if mouse moved
	nc.mouseMoved = nc.mouseMoved || dist > float64(nc.settings.GetTriggerDist())

	// Decrement sleep timer
	nc.slp -= 50 * time.Millisecond
	if nc.slp < 0 {
		nc.slp = 0
	}

	if nc.slp == 0 {
		nc.animateCat(pos, theta, dist)
	}
}

// animateCat handles cat animation based on state
func (nc *NekoController) animateCat(pos *Position, theta, dist float64) {
	doMove := false

	// State 0: Initialization
	if nc.state == 0 {
		if nc.init < 33 {
			doMove = true
			nc.slp = nc.settings.GetLoadDelay()
			nc.ox = nc.nekoBounds.Min.X + nc.nekoBounds.Dx()/2
			nc.oy = nc.nekoBounds.Min.Y + nc.nekoBounds.Dy()/2
			nc.no = nc.init
			nc.init++
		} else {
			nc.state = 1
		}
	} else if nc.state == 1 {
		// State 1: Chasing the mouse
		doMove = true
		nc.slp = nc.settings.GetRunDelay()

		run := float64(nc.settings.GetRunDist())
		if run > dist {
			run = dist
		}

		nc.ox = int(float64(nc.ox) + math.Cos(theta)*run)
		nc.oy = int(float64(nc.oy) - math.Sin(theta)*run)
		dist = dist - run

		if dist < float64(nc.settings.GetCatchDist()) {
			nc.state = 2
		}

		// Determine animation frame based on direction
		nc.updateDirectionalAnimation(theta)
		nc.mouseMoved = false

	} else {
		// State 2: Sleeping or preparing to sleep
		nc.handleSleepingState(pos)

		// State 3: Surprised (mouse moved while sleeping)
		if nc.mouseMoved {
			nc.slp = nc.settings.GetSurpriseDelay()
			nc.no = 32
			nc.ilc1 = 0
			nc.ilc2 = 0
			nc.state = 1
		}
	}

	// Update cat position
	if doMove {
		nc.moveCat()
	}

	// Render current frame
	nc.renderFrame()
}

// updateDirectionalAnimation updates the animation frame based on direction
func (nc *NekoController) updateDirectionalAnimation(theta float64) {
	pi := math.Pi

	if theta >= -pi/8 && theta <= pi/8 { // right
		if nc.no == 5 {
			nc.no = 6
		} else {
			nc.no = 5
		}
	} else if theta > pi/8 && theta < 3*pi/8 { // upper-right
		if nc.no == 3 {
			nc.no = 4
		} else {
			nc.no = 3
		}
	} else if theta >= 3*pi/8 && theta <= 5*pi/8 { // up
		if nc.no == 1 {
			nc.no = 2
		} else {
			nc.no = 1
		}
	} else if theta > 5*pi/8 && theta < 7*pi/8 { // upper-left
		if nc.no == 15 {
			nc.no = 16
		} else {
			nc.no = 15
		}
	} else if theta >= 7*pi/8 || theta <= -7*pi/8 { // left
		if nc.no == 13 {
			nc.no = 14
		} else {
			nc.no = 13
		}
	} else if theta > -7*pi/8 && theta < -5*pi/8 { // bottom-left
		if nc.no == 11 {
			nc.no = 12
		} else {
			nc.no = 11
		}
	} else if theta >= -5*pi/8 && theta <= -3*pi/8 { // down
		if nc.no == 9 {
			nc.no = 10
		} else {
			nc.no = 9
		}
	} else if theta > -3*pi/8 && theta < -pi/8 { // bottom-right
		if nc.no == 7 {
			nc.no = 8
		} else {
			nc.no = 7
		}
	}
}

// handleSleepingState handles the various sleeping/idle animations
func (nc *NekoController) handleSleepingState(pos *Position) {
	switch nc.no {
	case 0: // cat sit
		if pos != nil {
			nc.slp = nc.settings.GetSharpenDelay()
			switch *pos {
			case PositionOver:
				nc.no = 17
			case PositionUnder:
				nc.no = 21
			case PositionLeft:
				nc.no = 23
			case PositionRight:
				nc.no = 19
			default:
				nc.slp = nc.settings.GetSitDelay()
				nc.no = 31
			}
		} else {
			nc.slp = nc.settings.GetSitDelay()
			nc.no = 31
		}

	case 17, 18: // Scratching upward
		nc.handleScratchAnimation(17, 18)
	case 21, 22: // Scratching downward
		nc.handleScratchAnimation(21, 22)
	case 23, 24: // Scratching left
		nc.handleScratchAnimation(23, 24)
	case 19, 20: // Scratching right
		nc.handleScratchAnimation(19, 20)

	case 31: // cat lick
		nc.slp = nc.settings.GetSitDelay()
		nc.no = 25
		nc.ilc1++
		if nc.ilc1 == 6 {
			nc.slp = nc.settings.GetScratchDelay()
			nc.no = 27
			nc.ilc1 = 0
		}

	case 25:
		nc.slp = nc.settings.GetSitDelay()
		nc.no = 31

	case 27: // cat scratch
		nc.slp = nc.settings.GetScratchDelay()
		nc.no = 28

	case 28:
		nc.no = 27
		nc.ilc2++
		if nc.ilc2 == 4 {
			nc.no = 26
			nc.slp = nc.settings.GetYawnDelay()
			nc.ilc2 = 0
		}

	case 26: // cat yawn
		nc.no = 29
		nc.slp = nc.settings.GetSleepDelay()

	case 29: // cat sleep
		nc.no = 30
		nc.slp = nc.settings.GetSleepDelay()

	case 30:
		nc.no = 29
		nc.slp = nc.settings.GetSleepDelay()

	default:
		nc.no = 0
	}
}

// handleScratchAnimation handles the scratch animation cycles
func (nc *NekoController) handleScratchAnimation(frame1, frame2 int) {
	nc.slp = nc.settings.GetSharpenDelay()
	if nc.no == frame1 {
		nc.no = frame2
		nc.ilc1++
		if nc.ilc1 == 6 {
			nc.no = 27
			nc.ilc1 = 0
		}
	} else {
		nc.no = frame1
	}
}

// moveCat moves the cat to its current position
func (nc *NekoController) moveCat() {
	if nc.windowMode {
		// In window mode, position is relative to window
		// This would update the sprite position in the window
	} else {
		// In free mode, move the window itself
		if nc.window != nil {
			nc.window.SetPos(nc.ox+nc.windowOffset.X, nc.oy+nc.windowOffset.Y)
		}
	}
}

// renderFrame renders the current animation frame
func (nc *NekoController) renderFrame() {
	// This would render the current image (nc.images[nc.no])
	// Implementation depends on the rendering backend
}

// MoveCatInBox moves the cat inside the window
func (nc *NekoController) MoveCatInBox() {
	if nc.windowMode {
		winW, winH := nc.window.GetSize()
		nc.window.SetPos(nc.ox-winW/2, max(0, nc.oy-winH/2))
	}
	nc.calculateBounds()
	nc.moveCat()
}

// CatboxDeiconified handles window deiconification
func (nc *NekoController) CatboxDeiconified() {
	nc.calculateBounds()
	nc.ox = nc.nekoBounds.Min.X + nc.nekoBounds.Dx()/2
	nc.oy = nc.nekoBounds.Min.Y + nc.nekoBounds.Dy()/2
}

// CatboxMoved handles window movement
func (nc *NekoController) CatboxMoved() {
	if nc.windowMode {
		oldX := nc.nekoBounds.Min.X
		oldY := nc.nekoBounds.Min.Y

		nc.calculateBounds()

		dx := nc.nekoBounds.Min.X - oldX
		dy := nc.nekoBounds.Min.Y - oldY

		nc.ox += dx
		nc.oy += dy
	}
}

// Stop stops the animation loop
func (nc *NekoController) Stop() {
	nc.ticker.Stop()
	nc.done <- true
	if nc.xconn != nil {
		nc.xconn.Close()
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
