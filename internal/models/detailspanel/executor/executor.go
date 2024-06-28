package executor

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hrvadl/gowatchsql/internal/message"
)

func NewModel() Model {
	return Model{}
}

type Model struct {
	width  int
	height int

	input textinput.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleUpdateSize(msg.Width, msg.Height)
	default:
		return m, nil
	}
}

func (m Model) View() string {
}

func (m Model) Value() string {
	return m.input.Value()
}

func (m *Model) Focus() {}

func (m *Model) Unfocus() {}

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
