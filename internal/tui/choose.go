package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type chooseModel struct {
	options   []string
	cursor    int
	header    string
	title     string
	done      bool
	cancelled bool
	height    int
	offset    int
	width     int
}

func (m chooseModel) Init() tea.Cmd { return nil }

func (m chooseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		extra := 4
		if m.title != "" {
			extra = 8
		}
		m.height = min(msg.Height-extra, 20)
		if m.height < 5 {
			m.height = 5
		}
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
		titleWidth := min(m.width-4, 60)
		if titleWidth < 20 {
			titleWidth = 20
		}
		titleStyle := StyleHeader.Width(titleWidth)
		b.WriteString(titleStyle.Render(m.title) + "\n\n")
	}

	headerStyle := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	b.WriteString(headerStyle.Render(m.header) + "\n\n")

	if m.offset > 0 {
		b.WriteString(StyleMuted.Render("  ↑ more") + "\n")
	}

	end := min(m.offset+m.height, len(m.options))
	for i := m.offset; i < end; i++ {
		opt := m.options[i]
		cursor := "  "
		if i == m.cursor {
			cursor = StyleAccent.Render("> ")
			opt = StyleSecondary.Render(opt)
		}
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, opt))
	}

	if end < len(m.options) {
		b.WriteString(StyleMuted.Render("  ↓ more") + "\n")
	}

	b.WriteString("\n" + StyleMuted.Render("  Up/Down: navigate  Enter: select  Esc: cancel") + "\n")

	return b.String()
}

func RunChoose(header string, options []string) (string, error) {
	return RunChooseWithTitle("", header, options)
}

func RunChooseWithTitle(title, header string, options []string) (string, error) {
	if Quiet {
		if len(options) > 0 {
			return options[0], nil
		}
		return "", fmt.Errorf("no options")
	}

	m := chooseModel{options: options, header: header, title: title, height: 15}
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
