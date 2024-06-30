package newcontext

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/ui/message"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

const margin = 1

func NewModel() Model {
	form, name, dsn, confirm := newForm()

	return Model{
		state: state{
			active: true,
		},
		form:    form,
		name:    name,
		dsn:     dsn,
		confirm: confirm,
	}
}

type Model struct {
	width  int
	height int
	state  state

	form    *huh.Form
	dsn     huh.Field
	name    huh.Field
	confirm huh.Field
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleUpdateSize(msg.Width-margin*2, msg.Height-margin*2)
	case message.MoveFocus:
		return m.handleMoveFocus(msg)
	default:
		return m.handleDefault(msg)
	}
}

func (m Model) View() string {
	s := m.newContainerStyles()
	return s.Render(m.form.View())
}

func (m Model) handleDefault(msg tea.Msg) (Model, tea.Cmd) {
	m, cmd := m.delegateToForm(msg)
	switch m.form.State {
	case huh.StateNormal:
		return m, cmd
	default:
		m, cmpCmd := m.handleFormCompleted()
		return m, tea.Batch(cmd, cmpCmd)
	}
}

func (m Model) handleFormCompleted() (Model, tea.Cmd) {
	msg := message.NewContext{
		Name: m.form.GetString("name"),
		DSN:  m.form.GetString("dsn"),
	}

	if done := m.form.GetBool("done"); done {
		msg.OK = true
	}

	return m, func() tea.Msg { return msg }
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
	m.form = m.form.WithWidth(w).WithHeight(h)
	return m.delegateToForm(tea.WindowSizeMsg{Height: h, Width: w})
}

func (m Model) delegateToForm(msg tea.Msg) (Model, tea.Cmd) {
	form, cmd := m.form.Update(msg)
	m.form = form.(*huh.Form)
	return m, cmd
}

func (m Model) newContainerStyles() lipgloss.Style {
	return lipgloss.
		NewStyle().
		Height(m.height).
		Width(m.width).
		Align(lipgloss.Center)
}
