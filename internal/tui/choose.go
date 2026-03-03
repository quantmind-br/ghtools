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
	title     string
	subtitle  string
	done      bool
	cancelled bool
	height    int
	offset    int
	width     int
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
	case tea.WindowSizeMsg:
		titleLines := 0
		if m.title != "" {
			titleLines = 1
		}
		// Overhead: title(0-1) + header(1) + scroll-up(1) + scroll-down(1) + help(1) = titleLines + 4
		m.height = min(msg.Height-4-titleLines, len(m.options))
		if m.height < 3 {
			m.height = 3
		}
		m.width = msg.Width
		m.adjustScroll()
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
				m.adjustScroll()
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
				m.adjustScroll()
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

func (m *chooseModel) adjustScroll() {
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+m.height {
		m.offset = m.cursor - m.height + 1
	}
}

func (m chooseModel) View() string {
	var b strings.Builder

	if m.title != "" {
		titleStyle := lipgloss.NewStyle().Foreground(ColorSecondary).Bold(true)
		line := titleStyle.Render(m.title)
		if m.subtitle != "" {
			line += StyleMuted.Render("  " + m.subtitle)
		}
		b.WriteString(line + "\n")
	}

	headerStyle := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	header := m.header
	if m.width > 0 && len(header) > m.width {
		header = header[:m.width]
	}
	b.WriteString(headerStyle.Render(header) + "\n")

	if m.offset > 0 {
		b.WriteString(StyleMuted.Render("  ↑ more") + "\n")
	}

	end := min(m.offset+m.height, len(m.options))
	for i := m.offset; i < end; i++ {
		opt := m.options[i]

		prefix := "  "
		if i == m.cursor {
			prefix = StyleAccent.Render("> ")
		}
		if m.width > 0 {
			maxOptLen := m.width - 2
			if maxOptLen > 0 && len(opt) > maxOptLen {
				opt = opt[:maxOptLen]
			}
		}

		if i == m.cursor {
			opt = StyleSecondary.Render(opt)
		}
		fmt.Fprintf(&b, "%s%s\n", prefix, opt)
	}

	if m.offset+m.height < len(m.options) {
		b.WriteString(StyleMuted.Render("  ↓ more") + "\n")
	}

	helpText := "↑↓:nav  Enter:select  Esc:cancel"
	if len(m.shortcuts) > 0 {
		helpText += "  [key]:shortcut"
	}
	if m.width > 0 && len(helpText)+2 > m.width {
		helpText = helpText[:max(0, m.width-2)]
	}
	b.WriteString(StyleMuted.Render("  "+helpText) + "\n")

	return b.String()
}

func RunChooseWithTitle(title, subtitle, header string, options []string) (string, error) {
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
		title:     title,
		subtitle:  subtitle,
		height:    min(15, len(options)),
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
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
		height:    min(15, len(options)),
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
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
