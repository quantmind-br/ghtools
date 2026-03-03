package tui

import "github.com/charmbracelet/lipgloss"

var (
	ColorPrimary   = lipgloss.Color("99")  // Soft purple
	ColorSecondary = lipgloss.Color("39")  // Cyan
	ColorAccent    = lipgloss.Color("212") // Pink
	ColorSuccess   = lipgloss.Color("78")  // Green
	ColorWarning   = lipgloss.Color("220") // Yellow/Gold
	ColorError     = lipgloss.Color("196") // Red
	ColorInfo      = lipgloss.Color("75")  // Light blue
	ColorMuted     = lipgloss.Color("240") // Gray

	StylePrimary   = lipgloss.NewStyle().Foreground(ColorPrimary)
	StyleSecondary = lipgloss.NewStyle().Foreground(ColorSecondary)
	StyleAccent    = lipgloss.NewStyle().Foreground(ColorAccent)
	StyleSuccess   = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true)
	StyleWarning   = lipgloss.NewStyle().Foreground(ColorWarning).Bold(true)
	StyleError     = lipgloss.NewStyle().Foreground(ColorError).Bold(true)
	StyleInfo      = lipgloss.NewStyle().Foreground(ColorInfo)
	StyleMuted     = lipgloss.NewStyle().Foreground(ColorMuted)

	StyleHeader = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Foreground(ColorSecondary).
			Align(lipgloss.Center).
			Padding(1, 2).
			MarginLeft(2)

	StyleBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSecondary).
			Padding(1, 2).
			MarginLeft(2)
)
