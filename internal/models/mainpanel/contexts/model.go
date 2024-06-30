package contexts

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/color"
	"github.com/hrvadl/gowatchsql/internal/command"
	"github.com/hrvadl/gowatchsql/internal/message"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

const margin = 1

func NewModel() Model {
	return Model{
		list:  list.New(nil, list.NewDefaultDelegate(), 0, 0),
		state: state{active: true},
	}
}

type Model struct {
	width  int
	height int
	list   list.Model
	state  state
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSize(msg.Width-margin*2, msg.Height-margin*3)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case message.Error:
		return m.handleError(msg)
	case message.MoveFocus:
		return m.handleMoveFocus(msg)
	default:
		return m, nil
	}
}

func (m Model) View() string {
	s := m.newContainerStyles()
	return s.Render(m.list.View())
}

func (m Model) Help() string {
	return "help from contexts"
}

func (m Model) handleMoveFocus(msg message.MoveFocus) (Model, tea.Cmd) {
	switch msg.Direction {
	case direction.Away:
		m.state.active = false
	default:
		m.state.active = !m.state.active
	}
	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyRunes:
		return m.handleKeyRunes(msg)
	default:
		return m, nil
	}
}

func (m Model) handleKeyRunes(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "N":
		return m, func() tea.Msg { return message.Command{Text: command.NewContext} }
	default:
		return m.delegateToList(msg)
	}
}

func (m Model) handleError(msg message.Error) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) handleWindowSize(w, h int) (Model, tea.Cmd) {
	m.width = w
	m.height = h
	m.list.SetSize(w, h)
	return m, nil
}

func (m Model) delegateToList(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) newContainerStyles() lipgloss.Style {
	base := lipgloss.
		NewStyle().
		Height(m.height).
		Width(m.width).
		Border(lipgloss.ThickBorder())

	if m.state.active {
		return base.BorderForeground(color.MainAccent)
	}

	return base.BorderForeground(color.Border)
}
