package main

import "time"

const version = "dev"

// Configuration constants
const (
	// Ticker intervals for different duration ranges
	tickIntervalFast   = 100 * time.Millisecond // For durations < 1 minute
	tickIntervalMedium = 500 * time.Millisecond // For durations 1-10 minutes
	tickIntervalSlow   = 1 * time.Second        // For durations > 10 minutes

	// Warning threshold for countdown timer
	warningThreshold = 5 * time.Minute

	// Big text glyph dimensions
	glyphWidth   = 8
	glyphHeight  = 7
	glyphSpacing = 1

	// Visual spacing (terminal line height cannot be changed, but we can adjust visual perception)

	// Keyboard input buffer size
	keyBufferSize = 10

	// Default terminal dimensions
	defaultTermWidth  = 80
	defaultTermHeight = 24
)
