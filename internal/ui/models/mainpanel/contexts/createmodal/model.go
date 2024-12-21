package createmodal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/ui/message"
)

const margin = 1

func NewModel() Model {
	return Model{
		form: newForm(),
	}
}

type Model struct {
	width  int
	height int
	form   *huh.Form
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleUpdateSize(msg.Width-margin*2, msg.Height-margin*2)
	case tea.KeyMsg:
		return m.handleKeyMessage(msg)
	default:
		return m.handleDefault(msg)
	}
}

func (m Model) View() string {
	s := m.newContainerStyles()
	return s.Render(m.form.View())
}

func (m Model) handleKeyMessage(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		return m.handleFormCompleted()
	default:
		return m.handleDefault(msg)
	}
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

	m.form.State = huh.StateNormal
	m.form = newForm()
	return m, tea.Batch(m.form.Init(), func() tea.Msg { return msg })
}

func (m Model) handleUpdateSize(w, h int) (Model, tea.Cmd) {
	m.width = w
	m.height = h
	m.form = m.form.WithWidth(w).WithHeight(h)
	return m.delegateToForm(tea.WindowSizeMsg{Height: h, Width: w})
}

func (m Model) delegateToForm(msg tea.Msg) (Model, tea.Cmd) {
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}
	return m, cmd
}

func (m Model) newContainerStyles() lipgloss.Style {
	return lipgloss.
		NewStyle().
		Height(m.height).
		Width(m.width)
}
