package search

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/color"
	"github.com/hrvadl/gowatchsql/internal/message"
)

const (
	padding     = 1
	margin      = 1
	placeholder = "mysql://user:password@(db:3306)/database"
)

var inputStyles = lipgloss.NewStyle().MarginTop(margin)

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
		return m.handleUpdateSize(msg.Width-margin*2, msg.Height-margin)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	default:
		return m, nil
	}
}

func (m Model) View() string {
	const title = "Database"
	barStyles := m.newBarStyles()
	titleStyles := m.newTitleStyles()
	return barStyles.Render(titleStyles.Render(title), inputStyles.Render(m.input.View()))
}

func (m *Model) Focus() {
	m.input.Focus()
}

func (m *Model) Unfocus() {
	m.input.Blur()
}

func (m Model) Value() string {
	return strings.TrimSpace(m.input.Value())
}

func (m Model) IsFocused() bool {
	return m.input.Focused()
}

func (m Model) handleUpdateSize(w, h int) (tea.Model, tea.Cmd) {
	m.width = w
	m.height = h
	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.Type {
	case tea.KeyCtrlC:
		return nil, tea.Quit
	case tea.KeyEnter:
		m.Unfocus()
		return m, func() tea.Msg { return message.DSNReady{DSN: m.Value()} }
	default:
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
}

func (m Model) newBarStyles() lipgloss.Style {
	return lipgloss.
		NewStyle().
		Height(m.height).
		Width(m.width)
}

func (m Model) newTitleStyles() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(color.Border).
		Bold(true).
		Width(m.width)
}
