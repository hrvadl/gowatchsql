package newcontext

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/color"
)

const margin = 1

func NewModel() Model {
	return Model{}
}

type Model struct {
	width  int
	height int
}

func (m Model) Init() (Model, tea.Cmd) {
	return m, nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleUpdateSize(msg.Width-margin*2, msg.Height-margin*2)
	default:
		return m, nil
	}
}

func (m Model) View() string {
	s := m.newContainerStyles()
	return s.Render("hello from new context")
}

func (m Model) handleUpdateSize(w, h int) (Model, tea.Cmd) {
	m.width = w
	m.height = h
	return m, nil
}

func (m Model) newContainerStyles() lipgloss.Style {
	return lipgloss.
		NewStyle().
		Height(m.height).
		Width(m.width).
		Border(lipgloss.ThickBorder()).
		BorderForeground(color.MainAccent)
}
