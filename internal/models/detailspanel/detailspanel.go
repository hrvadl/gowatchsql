package detailspanel

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding = 1
	margin  = 1
)

func NewModel() Model {
	return Model{}
}

type Model struct {
	width  int
	height int

	table string

	state State
	err   error
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleUpdateSize(msg.Width-margin*2, msg.Height-margin*3)
	case TableChosenMsg:
		return m.handleTableChosen(msg)
	case ErrorMsg:
		return m.handleError(msg)
	default:
		return m, nil
	}
}

func (m Model) View() string {
	s := lipgloss.
		NewStyle().
		Height(m.height).
		Width(m.width).
		Padding(0, padding).
		Border(lipgloss.NormalBorder())
	return s.Render()
}

func (m Model) handleTableChosen(msg TableChosenMsg) (tea.Model, tea.Cmd) {
	m.table = msg.Name

	return m, nil
}

func (m Model) handleError(msg ErrorMsg) (tea.Model, tea.Cmd) {
	m.err = msg.Error
	m.state = Error
	return m, nil
}

func (m Model) handleUpdateSize(w, h int) (tea.Model, tea.Cmd) {
	m.width = w
	m.height = h
	return m, nil
}
