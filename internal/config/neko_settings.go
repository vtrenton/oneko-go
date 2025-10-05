/*
 * neko_settings.go
 *
 * Copyright (c) 2019 Jerry Reno (original Java version)
 * Copyright (c) 2025 Go port
 * This is public domain software, under the terms of the UNLICENSE
 * http://unlicense.org
 */

package config

import "time"

// NekoSettings manages all configuration for the Neko application
type NekoSettings struct {
	settings *Settings

	// Distance thresholds
	triggerDist int
	catchDist   int
	runDist     int

	// Mouse offset
	offsetX int
	offsetY int

	// Animation delays (in milliseconds)
	minDelay      time.Duration
	runDelay      time.Duration
	sitDelay      time.Duration
	scratchDelay  time.Duration
	sharpenDelay  time.Duration
	loadDelay     time.Duration
	sleepDelay    time.Duration
	yawnDelay     time.Duration
	surpriseDelay time.Duration
}

// Settings keys
const (
	keyHello            = "hello"
	keyTitle            = "windowTitle"
	keyTriggerDist      = "triggerDistance"
	keyCatchDist        = "catchDistance"
	keyRunDist          = "runDistancePerFrame"
	keyOffsetX          = "offsetx"
	keyOffsetY          = "offsety"
	keyMaxFramerate     = "maxFramerate"
	keyRunFramerate     = "runFramerate"
	keySitFramerate     = "sitFramerate"
	keySharpenFramerate = "sharpenFramerate"
	keyScratchFramerate = "scratchFramerate"
	keyLoadFramerate    = "loadFramerate"
	keySleepDelay       = "sleepDelay"
	keyYawnDelay        = "yawnDelay"
	keySurpriseDelay    = "surpriseDelay"
)

// NewNekoSettings creates a new NekoSettings instance
func NewNekoSettings() *NekoSettings {
	ns := &NekoSettings{
		settings: NewSettings("neko.properties"),
	}
	ns.Load()
	return ns
}

// GetTitle returns the window title
func (ns *NekoSettings) GetTitle() string {
	return ns.settings.GetString(keyTitle)
}

// getDelay converts frames-per-second to milliseconds duration
func (ns *NekoSettings) getDelay(key string) time.Duration {
	fps := ns.settings.GetInt(key, 100)
	if fps <= 0 {
		fps = 1
	}
	return time.Duration(1000/fps) * time.Millisecond
}

// Getter methods
func (ns *NekoSettings) GetTriggerDist() int             { return ns.triggerDist }
func (ns *NekoSettings) GetCatchDist() int               { return ns.catchDist }
func (ns *NekoSettings) GetRunDist() int                 { return ns.runDist }
func (ns *NekoSettings) GetOffsetX() int                 { return ns.offsetX }
func (ns *NekoSettings) GetOffsetY() int                 { return ns.offsetY }
func (ns *NekoSettings) GetMinDelay() time.Duration      { return ns.minDelay }
func (ns *NekoSettings) GetRunDelay() time.Duration      { return ns.runDelay }
func (ns *NekoSettings) GetSitDelay() time.Duration      { return ns.sitDelay }
func (ns *NekoSettings) GetScratchDelay() time.Duration  { return ns.scratchDelay }
func (ns *NekoSettings) GetSharpenDelay() time.Duration  { return ns.sharpenDelay }
func (ns *NekoSettings) GetLoadDelay() time.Duration     { return ns.loadDelay }
func (ns *NekoSettings) GetSleepDelay() time.Duration    { return ns.sleepDelay }
func (ns *NekoSettings) GetYawnDelay() time.Duration     { return ns.yawnDelay }
func (ns *NekoSettings) GetSurpriseDelay() time.Duration { return ns.surpriseDelay }

// Load loads all settings
func (ns *NekoSettings) Load() {
	ns.settings.Load()

	// Print hello message if configured
	hello := ns.settings.GetString(keyHello)
	if hello != "" {
		ns.settings.Print(hello)
	}

	// Load distance settings
	ns.triggerDist = ns.settings.GetInt(keyTriggerDist, 16)
	ns.catchDist = ns.settings.GetInt(keyCatchDist, 16)
	ns.runDist = ns.settings.GetInt(keyRunDist, 16)
	ns.offsetX = ns.settings.GetInt(keyOffsetX, 0)
	ns.offsetY = ns.settings.GetInt(keyOffsetY, 0)

	// Load delay settings
	ns.sleepDelay = time.Duration(ns.settings.GetInt(keySleepDelay, 1000)) * time.Millisecond
	ns.yawnDelay = time.Duration(ns.settings.GetInt(keyYawnDelay, 1000)) * time.Millisecond
	ns.surpriseDelay = time.Duration(ns.settings.GetInt(keySurpriseDelay, 1000)) * time.Millisecond

	// Convert framerate settings to delays
	ns.minDelay = ns.getDelay(keyMaxFramerate)
	ns.runDelay = ns.getDelay(keyRunFramerate)
	ns.sitDelay = ns.getDelay(keySitFramerate)
	ns.scratchDelay = ns.getDelay(keyScratchFramerate)
	ns.sharpenDelay = ns.getDelay(keySharpenFramerate)
	ns.loadDelay = ns.getDelay(keyLoadFramerate)
}
