package main

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/term"
)

// Terminal escape codes
const (
	clearScreen = "\033[2J"
	hideCursor  = "\033[?25l"
	showCursor  = "\033[?25h"
	altScreen   = "\033[?1049h"
	mainScreen  = "\033[?1049l"
	resetStyle  = "\033[0m"
	blueColor   = "\033[34m"    // Blue text color
	redColor    = "\033[31m"    // Red text color
	mouseOn     = "\033[?1000h" // Enable basic mouse tracking
	mouseOff    = "\033[?1000l" // Disable mouse tracking
)

func moveCursor(row, col int) string {
	return fmt.Sprintf("\033[%d;%dH", row, col)
}

// fixNewlines converts \n to \r\n for raw terminal mode
func fixNewlines(s string) string {
	return strings.ReplaceAll(s, "\n", "\r\n")
}

// setupTerminal configures the terminal for raw mode and returns the previous state
func setupTerminal() (*term.State, error) {
	fd := int(syscall.Stdin)
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return nil, fmt.Errorf("failed to set terminal to raw mode: %w", err)
	}
	return oldState, nil
}

// restoreTerminal restores the terminal to its previous state
func restoreTerminal(oldState *term.State) error {
	fd := int(syscall.Stdin)
	if err := term.Restore(fd, oldState); err != nil {
		return fmt.Errorf("failed to restore terminal: %w", err)
	}
	return nil
}

func getTerminalSize() (width, height int) {
	// Default fallback
	width, height = defaultTermWidth, defaultTermHeight

	// Use syscall to get actual terminal size
	type winsize struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}

	ws := &winsize{}
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if errno == 0 && ws.Col > 0 && ws.Row > 0 {
		width = int(ws.Col)
		height = int(ws.Row)
	}

	return width, height
}
