package tui

import (
	"fmt"
	"strings"
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
