package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DoneMsg struct {
	Err error
}

type spinnerModel struct {
	spinner spinner.Model
	title   string
	done    bool
	err     error
}

func newSpinnerModel(title string) spinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(ColorAccent)
	return spinnerModel{spinner: s, title: title}
}

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case DoneMsg:
		m.done = true
		m.err = msg.Err
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m spinnerModel) View() string {
	if m.done {
		if m.err != nil {
			return StyleError.Render("Failed: " + m.err.Error()) + "\n"
		}
		return StyleSuccess.Render("Done!") + "\n"
	}
	return fmt.Sprintf("%s %s\n", m.spinner.View(), m.title)
}

func RunWithSpinner(title string, fn func() error) error {
	if Quiet {
		return fn()
	}

	m := newSpinnerModel(title)
	p := tea.NewProgram(m)

	go func() {
		err := fn()
		p.Send(DoneMsg{Err: err})
	}()

	finalModel, err := p.Run()
	if err != nil {
		return err
	}
	result := finalModel.(spinnerModel)
	return result.err
}
