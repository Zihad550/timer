package main

import (
	"fmt"
	"strings"
	"time"
)

func formatHMS(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	total := int(d.Round(time.Second).Seconds())
	h := total / 3600
	m := (total % 3600) / 60
	s := total % 60
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

func renderBigTime(timeStr string, termWidth, termHeight int) string {
	// Calculate if we can fit big text
	totalWidth := len(timeStr)*(glyphWidth+glyphSpacing) - glyphSpacing

	// If too small, return simple text
	if termWidth < totalWidth+4 || termHeight < glyphHeight+2 {
		return timeStr
	}

	// Build the big text line by line
	lines := make([]string, 0, glyphHeight)
	for row := 0; row < glyphHeight; row++ {
		// Pre-allocate capacity for the line builder
		var line strings.Builder
		line.Grow(totalWidth)

		for i, ch := range timeStr {
			glyph, ok := glyphs[ch]
			if !ok {
				glyph = glyphs[' ']
			}
			line.WriteString(glyph[row])
			if i < len(timeStr)-1 {
				line.WriteString("  ") // 2 spaces between characters
			}
		}
		lines = append(lines, line.String())
	}

	return strings.Join(lines, "\n")
}

func centerText(text string, width, height int) string {
	lines := strings.Split(text, "\n")

	// Calculate vertical centering
	vOffset := (height - len(lines)) / 2
	if vOffset < 0 {
		vOffset = 0
	}

	// Pre-allocate builder capacity
	var result strings.Builder
	estimatedSize := (vOffset + len(lines)) * (width + 1) // rough estimate
	result.Grow(estimatedSize)

	// Add vertical padding
	for i := 0; i < vOffset; i++ {
		result.WriteString("\n")
	}

	// Center each line horizontally
	for _, line := range lines {
		lineLen := len([]rune(line))
		hOffset := (width - lineLen) / 2
		if hOffset < 0 {
			hOffset = 0
		}
		result.WriteString(strings.Repeat(" ", hOffset))
		result.WriteString(line)
		result.WriteString("\n")
	}

	return result.String()
}
