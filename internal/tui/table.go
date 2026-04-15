package tui

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

func GetTerminalSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80, 24 // Default fallback
	}
	return width, height
}

func CalculateDynamicWidths(headers []string, rows [][]string, terminalWidth int) []int {
	widths := make([]int, len(headers))

	// Start with header widths
	for i, h := range headers {
		widths[i] = len(h)
	}

	// Find max content width for each column
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) {
				cellLen := len(cell)
				if cellLen > widths[i] {
					widths[i] = cellLen
				}
			}
		}
	}

	// Calculate total width needed
	totalNeeded := 0
	for _, w := range widths {
		totalNeeded += w + 1 // +1 for space between columns
	}

	// If content fits, use it. Otherwise, scale down proportionally
	if totalNeeded <= terminalWidth {
		return widths
	}

	// Scale down proportionally
	scale := float64(terminalWidth-10) / float64(totalNeeded)
	for i := range widths {
		w := int(float64(widths[i]) * scale)
		if w < 8 { // Minimum width
			widths[i] = 8
		} else {
			widths[i] = w
		}
	}

	return widths
}

func PrintTable(headers []string, widths []int, rows [][]string) {
	// Print header
	var headerParts []string
	for i, h := range headers {
		headerParts = append(headerParts, fmt.Sprintf("%-*s", widths[i], h))
	}
	fmt.Println(StyleSecondary.Render(strings.Join(headerParts, " ")))

	// Print separator
	totalWidth := 0
	for _, w := range widths {
		totalWidth += w + 1
	}
	fmt.Println(StyleSecondary.Render(strings.Repeat("-", totalWidth)))

	// Print rows
	for _, row := range rows {
		var parts []string
		for i, cell := range row {
			if i < len(widths) {
				parts = append(parts, fmt.Sprintf("%-*s", widths[i], Truncate(cell, widths[i])))
			} else {
				parts = append(parts, cell)
			}
		}
		fmt.Println(strings.Join(parts, " "))
	}
}
