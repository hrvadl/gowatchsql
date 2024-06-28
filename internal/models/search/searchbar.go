package search

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding     = 1
	margin      = 1
	placeholder = "mysql://user:password@(db:3306)/database"
)

func New() Model {
	input := textinput.New()
	input.Focus()
	input.Placeholder = placeholder
	return Model{
		input: input,
	}
}

type Model struct {
	width  int
	height int
	input  textinput.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleUpdateSize(msg.Width-margin*2, msg.Height-margin)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	default:
	}
	return m, nil
}

func (m Model) View() string {
	s := lipgloss.
		NewStyle().
		Height(m.height).
		Width(m.width).
		Padding(0, padding).
		Border(lipgloss.NormalBorder())

	ts := lipgloss.NewStyle().MarginTop(margin).Render(m.input.View())
	return s.Render("Input your DSN:", ts)
}

func (m Model) Value() string {
	return m.input.Value()
}

func (m Model) IsFocused() bool {
	return m.input.Focused()
}

func (m *Model) handleUpdateSize(w, h int) {
	m.width = w
	m.height = h
}

func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.Type {
	case tea.KeyCtrlC:
		return nil, tea.Quit
	case tea.KeyEnter:
		m.input.Blur()
	default:
		m.input, cmd = m.input.Update(msg)
	}
	return *m, cmd
}
