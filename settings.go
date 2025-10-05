/*
 * settings.go
 *
 * Copyright (c) 2019 Jerry Reno (original Java version)
 * Copyright (c) 2025 Go port
 * This is public domain software, under the terms of the UNLICENSE
 * http://unlicense.org
 */

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Settings loads properties from embedded resources and home directory overrides
type Settings struct {
	filename      string
	homefile      string
	builtinProps  map[string]string
	overrideProps map[string]string
}

// NewSettings creates a new Settings instance
func NewSettings(filename string) *Settings {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	return &Settings{
		filename:      filename,
		homefile:      filepath.Join(home, filename),
		builtinProps:  make(map[string]string),
		overrideProps: make(map[string]string),
	}
}

// Load loads properties from the builtin resource and home directory
func (s *Settings) Load() {
	// Clear existing properties
	s.builtinProps = make(map[string]string)

	// Load builtin properties from resources directory
	builtinPath := filepath.Join("resources", s.filename)
	s.loadPropertiesFile(builtinPath, s.builtinProps)

	// Load override properties from home directory
	s.LoadOverride()
}

// LoadOverride loads properties from the home directory override file
func (s *Settings) LoadOverride() {
	s.overrideProps = make(map[string]string)

	if _, err := os.Stat(s.homefile); os.IsNotExist(err) {
		return
	}

	s.loadPropertiesFile(s.homefile, s.overrideProps)
}

// loadPropertiesFile loads a Java-style properties file
func (s *Settings) loadPropertiesFile(path string, props map[string]string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			props[key] = value
		}
	}
}

// GetString gets a string property with optional default
func (s *Settings) GetString(key string, defaultValue ...string) string {
	// Check override first
	if val, ok := s.overrideProps[key]; ok {
		return val
	}

	// Then check builtin
	if val, ok := s.builtinProps[key]; ok {
		return val
	}

	// Return default if provided
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return ""
}

// GetInt gets an integer property with optional default
func (s *Settings) GetInt(key string, defaultValue ...int) int {
	var defVal int
	if len(defaultValue) > 0 {
		defVal = defaultValue[0]
	}

	// Check override first
	if val, ok := s.overrideProps[key]; ok {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}

	// Then check builtin
	if val, ok := s.builtinProps[key]; ok {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}

	return defVal
}

// Print prints a message to stdout (for debugging)
func (s *Settings) Print(msg string) {
	fmt.Println(msg)
}
