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
	fmt.Println()
	fmt.Println(StyleHeader.Render(content))
	fmt.Println()
}
