package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type confirmModel struct {
	prompt    string
	value     bool
	done      bool
	cancelled bool
}

func (m confirmModel) Init() tea.Cmd { return nil }

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			m.done = true
			return m, tea.Quit
		case "y", "Y":
			m.value = true
			m.done = true
			return m, tea.Quit
		case "n", "N":
			m.value = false
			m.done = true
			return m, tea.Quit
		case "enter":
			m.done = true
			return m, tea.Quit
		case "left", "h":
			m.value = true
		case "right", "l":
			m.value = false
		}
	}
	return m, nil
}

func (m confirmModel) View() string {
	var b strings.Builder

	promptStyle := lipgloss.NewStyle().Foreground(ColorWarning).Bold(true)
	b.WriteString(promptStyle.Render(m.prompt) + " ")

	yes := "Yes"
	no := "No"
	if m.value {
		yes = StyleSuccess.Render("[Yes]")
		no = StyleMuted.Render(" No ")
	} else {
		yes = StyleMuted.Render(" Yes ")
		no = StyleError.Render("[No]")
	}
	b.WriteString(fmt.Sprintf("%s / %s", yes, no))
	b.WriteString("\n" + StyleMuted.Render("  y/n: choose  Enter: confirm") + "\n")

	return b.String()
}

func RunConfirm(prompt string, defaultYes bool) (bool, error) {
	if Quiet {
		return defaultYes, nil
	}

	m := confirmModel{prompt: prompt, value: defaultYes}
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return false, err
	}

	result := finalModel.(confirmModel)
	if result.cancelled {
		return false, fmt.Errorf("cancelled")
	}

	return result.value, nil
}
