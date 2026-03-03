package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type chooseModel struct {
	options   []string
	shortcuts map[string]int // key -> option index
	cursor    int
	header    string
	done      bool
	cancelled bool
}

// parseShortcuts extracts keyboard shortcuts from options like "[L] List Repositories"
func parseShortcuts(options []string) map[string]int {
	shortcuts := make(map[string]int)
	for i, opt := range options {
		// Match patterns like [L], [S], etc.
		if len(opt) >= 4 && opt[0] == '[' && opt[2] == ']' {
			key := strings.ToLower(opt[1:2])
			shortcuts[key] = i
		}
	}
	return shortcuts
}

func (m chooseModel) Init() tea.Cmd { return nil }

func (m chooseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			m.cancelled = true
			m.done = true
			return m, tea.Quit
		case "enter":
			m.done = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		default:
			// Check for shortcut key press
			key := strings.ToLower(msg.String())
			if idx, ok := m.shortcuts[key]; ok {
				m.cursor = idx
				m.done = true
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m chooseModel) View() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	b.WriteString(headerStyle.Render(m.header) + "\n\n")

	for i, opt := range m.options {
		cursor := "  "
		if i == m.cursor {
			cursor = StyleAccent.Render("> ")
			opt = StyleSecondary.Render(opt)
		}
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, opt))
	}

	helpText := "  Up/Down: navigate  Enter: select  Esc: cancel"
	if len(m.shortcuts) > 0 {
		helpText += "  Shortcuts: [key] to select"
	}
	b.WriteString("\n" + StyleMuted.Render(helpText) + "\n")

	return b.String()
}

func RunChoose(header string, options []string) (string, error) {
	if Quiet {
		if len(options) > 0 {
			return options[0], nil
		}
		return "", fmt.Errorf("no options")
	}

	m := chooseModel{
		options:   options,
		shortcuts: parseShortcuts(options),
		header:    header,
	}
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	result := finalModel.(chooseModel)
	if result.cancelled {
		return "", fmt.Errorf("cancelled")
	}

	return result.options[result.cursor], nil
}
