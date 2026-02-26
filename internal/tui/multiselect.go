package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MultiSelectItem struct {
	Label string
	Value string
}

type multiSelectModel struct {
	items     []MultiSelectItem
	filtered  []int
	selected  map[int]bool
	cursor    int
	filter    textinput.Model
	header    string
	done      bool
	cancelled bool
	height    int
	offset    int
}

func newMultiSelectModel(header string, items []MultiSelectItem) multiSelectModel {
	ti := textinput.New()
	ti.Placeholder = "Type to filter..."
	ti.Focus()
	ti.PromptStyle = lipgloss.NewStyle().Foreground(ColorPrimary)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(ColorAccent)

	filtered := make([]int, len(items))
	for i := range items {
		filtered[i] = i
	}

	return multiSelectModel{
		items:    items,
		filtered: filtered,
		selected: make(map[int]bool),
		filter:   ti,
		header:   header,
		height:   15,
	}
}

func (m multiSelectModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m multiSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = min(msg.Height-6, 20)
		if m.height < 5 {
			m.height = 5
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			m.done = true
			return m, tea.Quit
		case "enter":
			m.done = true
			return m, tea.Quit
		case "tab", " ":
			if len(m.filtered) > 0 {
				idx := m.filtered[m.cursor]
				m.selected[idx] = !m.selected[idx]
				if !m.selected[idx] {
					delete(m.selected, idx)
				}
				if m.cursor < len(m.filtered)-1 {
					m.cursor++
					m.adjustScroll()
				}
			}
			return m, nil
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m.adjustScroll()
			}
			return m, nil
		case "down", "j":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
				m.adjustScroll()
			}
			return m, nil
		case "ctrl+a":
			for _, idx := range m.filtered {
				m.selected[idx] = true
			}
			return m, nil
		case "ctrl+d":
			m.selected = make(map[int]bool)
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.filter, cmd = m.filter.Update(msg)

	// Re-filter
	query := strings.ToLower(m.filter.Value())
	m.filtered = nil
	for i, item := range m.items {
		if query == "" || strings.Contains(strings.ToLower(item.Label), query) {
			m.filtered = append(m.filtered, i)
		}
	}
	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
	}
	m.adjustScroll()

	return m, cmd
}

func (m *multiSelectModel) adjustScroll() {
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+m.height {
		m.offset = m.cursor - m.height + 1
	}
}

func (m multiSelectModel) View() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	b.WriteString(headerStyle.Render(m.header) + "\n")
	b.WriteString(m.filter.View() + "\n")

	selectedCount := len(m.selected)
	countStyle := lipgloss.NewStyle().Foreground(ColorMuted)
	b.WriteString(countStyle.Render(fmt.Sprintf("  %d/%d selected  Tab:toggle  Ctrl+A:all  Enter:confirm", selectedCount, len(m.items))) + "\n")

	end := min(m.offset+m.height, len(m.filtered))
	for i := m.offset; i < end; i++ {
		idx := m.filtered[i]
		item := m.items[idx]

		cursor := "  "
		if i == m.cursor {
			cursor = StyleAccent.Render("> ")
		}

		check := "[ ] "
		if m.selected[idx] {
			check = StyleSuccess.Render("[x] ")
		}

		label := item.Label
		if i == m.cursor {
			label = StyleSecondary.Render(label)
		}

		b.WriteString(cursor + check + label + "\n")
	}

	if len(m.filtered) == 0 {
		b.WriteString(StyleMuted.Render("  No matches") + "\n")
	}

	return b.String()
}

func RunMultiSelect(header string, items []MultiSelectItem) ([]MultiSelectItem, error) {
	m := newMultiSelectModel(header, items)
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	result := finalModel.(multiSelectModel)
	if result.cancelled {
		return nil, nil
	}

	var selected []MultiSelectItem
	for idx := range result.selected {
		selected = append(selected, result.items[idx])
	}
	return selected, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
