package newcontext

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/message"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

const margin = 1

func NewModel() Model {
	form, name, dsn, confirm := newForm()

	return Model{
		state: state{
			active: true,
			form: formState{
				inputFocused: nameInputFocused,
				page:         inputsPage,
			},
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
	dsn     *huh.Input
	name    *huh.Input
	confirm *huh.Confirm
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.form.Init(), m.form.NextField())
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
	case tea.KeyTab, tea.KeyShiftTab:
		return m.handleTab(msg)
	case tea.KeyEnter:
		return m.handleEnterKey(msg)
	default:
		return m.delegateToForm(msg)
	}
}

func (m Model) handleTab(msg tea.KeyMsg) (Model, tea.Cmd) {
	m, cmd := m.delegateToForm(msg)
	switch m.state.form.inputFocused {
	case nameInputFocused:
		m.state.form.inputFocused = dsnInputFocused
		return m, tea.Batch(cmd, m.form.NextField())
	case dsnInputFocused:
		m.state.form.inputFocused = nameInputFocused
		return m, tea.Batch(cmd, m.form.PrevField())
	default:
		return m, nil
	}
}

func (m Model) handleEnterKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch m.state.form.page {
	case inputsPage:
		return m.handleEnterKeyFirstPage(msg)
	case confirmationPage:
		return m.handleEnterKeySecondPage(msg)
	default:
		return m, nil
	}
}

func (m Model) handleEnterKeyFirstPage(msg tea.KeyMsg) (Model, tea.Cmd) {
	m, cmd := m.delegateToForm(msg)
	switch m.state.form.inputFocused {
	case nameInputFocused:
		m.state.form.inputFocused = dsnInputFocused
		return m, tea.Batch(cmd, m.form.NextField())
	case dsnInputFocused:
		m.state.form.inputFocused = nameInputFocused
		m.state.form.page = confirmationPage
		return m, tea.Batch(cmd, m.form.NextGroup(), m.confirm.Focus())
	default:
		return m, nil
	}
}

func (m Model) handleEnterKeySecondPage(msg tea.KeyMsg) (Model, tea.Cmd) {
	m, cmd := m.delegateToForm(msg)
	name, dsn, done := m.getValues()
	m, resetCmd := m.resetForm()

	doneMsg := message.NewContext{
		DSN:  dsn,
		Name: name,
		OK:   done,
	}

	return m, tea.Batch(cmd, resetCmd, func() tea.Msg { return doneMsg })
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
	form, cmd := m.form.Update(msg)
	m.form = form.(*huh.Form)
	return m, cmd
}

func (m Model) resetForm() (Model, tea.Cmd) {
	m.form, m.name, m.dsn, m.confirm = newForm()
	m.state = state{}
	return m.delegateToForm(tea.WindowSizeMsg{Height: m.height, Width: m.width})
}

func (m Model) getValues() (string, string, bool) {
	done, ok := m.confirm.GetValue().(bool)
	if !ok {
		return "", "", false
	}

	name, ok := m.name.GetValue().(string)
	if !ok {
		return "", "", false
	}

	dsn, ok := m.dsn.GetValue().(string)
	if !ok {
		return "", "", false
	}

	return strings.TrimSpace(name), strings.TrimSpace(dsn), done
}

func (m Model) newContainerStyles() lipgloss.Style {
	return lipgloss.
		NewStyle().
		Height(m.height).
		Width(m.width)
}

func newForm() (*huh.Form, *huh.Input, *huh.Input, *huh.Confirm) {
	nameInput := huh.NewInput().
		Key("name").
		Title("Human-readable name:")
	nameInput.Focus()

	dsnInput := huh.NewInput().
		Key("dsn").
		Title("DSN")

	confirm := huh.NewConfirm().
		Key("done").
		Title("Are you sure?").
		Affirmative("Yes").
		Negative("No")

	form := huh.NewForm(
		huh.NewGroup(
			nameInput,
			dsnInput,
		),
		huh.NewGroup(
			confirm,
		),
	)

	return form, nameInput, dsnInput, confirm
}
