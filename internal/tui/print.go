package tui

import "fmt"

var Quiet bool

func PrintError(msg string) {
	fmt.Println(StyleError.Render("ERROR") + " " + msg)
}

func PrintSuccess(msg string) {
	if Quiet {
		return
	}
	fmt.Println(StyleSuccess.Render("OK") + " " + msg)
}

func PrintInfo(msg string) {
	if Quiet {
		return
	}
	fmt.Println(StyleInfo.Render("INFO") + " " + msg)
}

func PrintWarning(msg string) {
	if Quiet {
		return
	}
	fmt.Println(StyleWarning.Render("WARN") + " " + msg)
}

func ShowHeader(title string, subtitle string) {
	content := title
	if subtitle != "" {
		content += "\n" + subtitle
	}
	termWidth, _ := GetTerminalSize()
	dynWidth := termWidth - 6
	if dynWidth > 60 {
		dynWidth = 60
	}
	fmt.Println()
	fmt.Println(StyleHeader.Width(dynWidth).Render(content))
	fmt.Println()
}

func ShowSection(title string, content string) {
	fmt.Println()
	fmt.Println(StyleMuted.Render("=== " + title + " ==="))
	fmt.Println()
	fmt.Println(content)
}

func ShowEmptyState(message string) {
	fmt.Println()
	ShowHeader("No Results", message)
}
