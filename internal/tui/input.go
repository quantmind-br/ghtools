package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type inputModel struct {
	input     textinput.Model
	header    string
	done      bool
	cancelled bool
}

func (m inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			m.done = true
			return m, tea.Quit
		case "enter":
			m.done = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	headerStyle := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	return headerStyle.Render(m.header) + "\n" + m.input.View() + "\n" +
		StyleMuted.Render("  Enter: confirm  Esc: cancel") + "\n"
}

func RunInput(header string, placeholder string, defaultValue string) (string, error) {
	if Quiet {
		if defaultValue != "" {
			return defaultValue, nil
		}
		return placeholder, nil
	}

	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.SetValue(defaultValue)
	ti.Focus()
	ti.PromptStyle = lipgloss.NewStyle().Foreground(ColorAccent)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(ColorAccent)

	m := inputModel{input: ti, header: header}
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	result := finalModel.(inputModel)
	if result.cancelled {
		return "", fmt.Errorf("cancelled")
	}

	return result.input.Value(), nil
}
