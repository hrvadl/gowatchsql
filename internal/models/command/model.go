package command

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/color"
	"github.com/hrvadl/gowatchsql/internal/command"
	"github.com/hrvadl/gowatchsql/internal/message"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

const (
	padding     = 1
	margin      = 1
	placeholder = "mysql://user:password@(db:3306)/database"
)

var inputStyles = lipgloss.NewStyle().MarginTop(margin)

func NewModel() Model {
	input := textinput.New()
	input.Focus()
	input.Placeholder = placeholder
	return Model{
		input: input,
		state: state{
			active: true,
		},
	}
}

type Model struct {
	width  int
	height int
	state  state
	input  textinput.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleUpdateSize(msg.Width-margin*2, msg.Height-margin)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case message.MoveFocus:
		return m.handleFocus()
	default:
		return m, nil
	}
}

func (m Model) View() string {
	const title = "Command prompt"
	barStyles := m.newBarStyles()
	titleStyles := m.newTitleStyles()
	return barStyles.Render(titleStyles.Render(title), inputStyles.Render(m.input.View()))
}

func (m Model) Help() string {
	return "searchbar help"
}

func (m Model) Value() string {
	return strings.TrimSpace(m.input.Value())
}

func (m Model) handleFocus() (Model, tea.Cmd) {
	m.input.Focus()
	m.state.active = true
	return m, nil
}

func (m Model) handleMoveFocus(to direction.Direction) (Model, tea.Cmd) {
	m.input.Blur()
	m.state.active = false
	return m, func() tea.Msg { return message.MoveFocus{Direction: to} }
}

func (m Model) handleUpdateSize(w, h int) (Model, tea.Cmd) {
	m.width = w
	m.height = h
	return m, nil
}

func (m Model) handleKeyEnter() (Model, tea.Cmd) {
	m.input.Blur()
	m.state.active = false
	return m, func() tea.Msg { return message.Command{Text: command.Command(m.Value())} }
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.Type {
	case tea.KeyCtrlC:
		return Model{}, tea.Quit
	case tea.KeyEsc:
		return m.handleMoveFocus(direction.Forward)
	case tea.KeyShiftTab:
		return m.handleMoveFocus(direction.Backwards)
	case tea.KeyEnter:
		return m.handleKeyEnter()
	default:
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
}

func (m Model) newBarStyles() lipgloss.Style {
	base := lipgloss.
		NewStyle().
		Height(m.height).
		Width(m.width).
		Border(lipgloss.NormalBorder())

	if m.state.active {
		return base.Border(lipgloss.ThickBorder()).
			BorderForeground(color.MainAccent)
	}

	return base.BorderForeground(color.Border)
}

func (m Model) newTitleStyles() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(color.Border).
		Bold(true).
		Width(m.width)
}
