package newcontext

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/color"
	"github.com/hrvadl/gowatchsql/internal/message"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

const margin = 1

func NewModel() Model {
	var (
		name = new(string)
		dsn  = new(string)
	)

	return Model{
		name: name,
		dsn:  dsn,
		state: state{
			active: true,
		},

		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Human-readable name:").
					Value(name),
				huh.NewInput().
					Title("DSN:").
					Value(dsn),
			),
		),
	}
}

type Model struct {
	width  int
	height int

	state state
	name  *string
	dsn   *string

	form *huh.Form
}

func (m Model) Init() (Model, tea.Cmd) {
	return m, m.form.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleUpdateSize(msg.Width-margin*2, msg.Height-margin*2)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case message.MoveFocus:
		return m.handleMoveFocus(msg)
	default:
		return m, nil
	}
}

func (m Model) View() string {
	s := m.newContainerStyles()
	return s.Render(m.form.View())
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	default:
		return m.delegateToForm(msg)
	}
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

func (m Model) handleUpdateSize(w, h int) (Model, tea.Cmd) {
	m.width = w
	m.height = h
	return m.delegateToForm(tea.WindowSizeMsg{Height: h, Width: w})
}

func (m Model) delegateToForm(msg tea.Msg) (Model, tea.Cmd) {
	slog.Info("delegate to form")
	form, cmd := m.form.Update(msg)
	m.form = form.(*huh.Form)
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
