package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// getTickerInterval returns the appropriate ticker interval based on duration
func getTickerInterval(duration time.Duration) time.Duration {
	if duration == 0 {
		// Counter mode - use fast interval for smooth display
		return tickIntervalFast
	}
	if duration < time.Minute {
		return tickIntervalFast
	}
	if duration < 10*time.Minute {
		return tickIntervalMedium
	}
	return tickIntervalSlow
}

func runTimer(duration time.Duration, useFullscreen bool) error {
	// Determine if counter mode (duration == 0)
	isCounter := duration == 0

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Enter alt screen if fullscreen
	if useFullscreen {
		fmt.Print(altScreen)
		defer fmt.Print(mainScreen)
	}

	// Hide cursor
	fmt.Print(hideCursor)
	defer fmt.Print(showCursor)

	// Configure terminal for raw mode
	oldState, err := setupTerminal()
	if err != nil {
		return err
	}
	defer func() {
		if err := restoreTerminal(oldState); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to restore terminal: %v\n", err)
		}
	}()

	// Channel for keyboard input
	keysCh := make(chan byte, keyBufferSize)
	defer close(keysCh)

	// Start keyboard reader goroutine (blocking read, low CPU)
	fd := int(syscall.Stdin)
	go func() {
		buf := make([]byte, 1)
		for {
			// Blocking read - this will wait for input without consuming CPU
			n, err := syscall.Read(fd, buf)
			if err != nil {
				// Only exit on real errors
				if err != syscall.EINTR {
					return
				}
				continue
			}
			if n > 0 {
				select {
				case keysCh <- buf[0]:
				default:
					// Drop key if channel is full
				}
			}
		}
	}()

	start := time.Now()
	// Use adaptive ticker interval based on duration
	tickInterval := getTickerInterval(duration)
	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	// Pause state
	var paused bool
	var pauseStart time.Time
	var totalPausedDuration time.Duration

	// Cache for rendered output
	var lastRenderedSec int64 = -1
	var cachedOutput string

	// Initial render - show the starting time immediately
	var initialDisplayTime time.Duration
	if isCounter {
		// Counter mode - start at 00:00
		initialDisplayTime = 0
	} else {
		// Timer mode - show full duration
		initialDisplayTime = duration
	}

	// Render initial state
	timeStr := formatHMS(initialDisplayTime)
	if useFullscreen {
		width, height := getTerminalSize()
		bigText := renderBigTime(timeStr, width, height)
		centeredText := centerText(bigText, width, height)
		cachedOutput = centeredText
		fmt.Print(clearScreen + moveCursor(1, 1) + fixNewlines(cachedOutput))
	} else {
		cachedOutput = fmt.Sprintf("\r%s   ", timeStr)
		fmt.Print(cachedOutput)
	}
	lastRenderedSec = int64(initialDisplayTime.Seconds())

	for {
		select {
		case <-sigCh:
			return nil

		case key := <-keysCh:
			// Handle keyboard input
			switch key {
			case 0x20: // Space key - pause/unpause
				if paused {
					// Unpause
					totalPausedDuration += time.Since(pauseStart)
					paused = false
				} else {
					// Pause
					paused = true
					pauseStart = time.Now()
				}
				// Force re-render
				lastRenderedSec = -1

			case 'q', 'Q': // q or Q - quit
				fmt.Print("\r\nquitting...\r\n")
				return nil

			case 0x03: // Ctrl+C
				return nil
			}

		case <-ticker.C:
			// Calculate effective elapsed time (excluding paused duration)
			elapsed := time.Since(start) - totalPausedDuration
			if paused {
				elapsed -= time.Since(pauseStart)
			}

			var displayTime time.Duration
			var currentSec int64

			if isCounter {
				// Counter mode - count up
				displayTime = elapsed
				currentSec = int64(elapsed.Seconds())
				// Never exit automatically in counter mode
			} else {
				// Timer mode - count down
				if elapsed >= duration {
					// Timer finished
					fmt.Print("\r\nfinished!\r\n")
					return nil
				}
				displayTime = duration - elapsed
				currentSec = int64(displayTime.Seconds())
			}

			// Re-render when second changes OR when paused state changes
			if currentSec != lastRenderedSec || lastRenderedSec == -1 {
				lastRenderedSec = currentSec

				// Format time
				timeStr := formatHMS(displayTime)

				// Determine color based on state and time remaining
				var color string
				if paused {
					color = blueColor
				} else if !isCounter && displayTime < warningThreshold {
					// Only show red warning in timer mode
					color = redColor
				} else {
					color = ""
				}

				if useFullscreen {
					// Get terminal size
					width, height := getTerminalSize()

					// Render big text
					bigText := renderBigTime(timeStr, width, height)

					// Center the output first
					centeredText := centerText(bigText, width, height)

					// Apply color (paused = blue, <5min = red, else = default)
					if color != "" {
						cachedOutput = color + centeredText + resetStyle
					} else {
						cachedOutput = centeredText
					}
				} else {
					// Simple inline display
					if color != "" {
						cachedOutput = fmt.Sprintf("\r%s%s%s   ", color, timeStr, resetStyle)
					} else {
						cachedOutput = fmt.Sprintf("\r%s   ", timeStr)
					}
				}
			}

			// Output the cached rendering (fix newlines for raw mode)
			if useFullscreen {
				fmt.Print(clearScreen + moveCursor(1, 1) + fixNewlines(cachedOutput))
			} else {
				fmt.Print(cachedOutput)
			}
		}
	}
}
