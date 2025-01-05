package queryrun

import (
	"context"
	"log/slog"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/domain/engine"
	"github.com/hrvadl/gowatchsql/internal/ui/color"
	"github.com/hrvadl/gowatchsql/internal/ui/command"
	"github.com/hrvadl/gowatchsql/internal/ui/message"
	"github.com/hrvadl/gowatchsql/internal/ui/models/mainpanel/objects/panels/details/rows"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

const (
	padding     = 1
	margin      = 1
	placeholder = "SELECT * FROM"
)

var inputStyles = lipgloss.NewStyle().MarginTop(margin)

type ExplorerFactory interface {
	Create(ctx context.Context, name, dsn string) (engine.Explorer, error)
}

func NewModel(ef ExplorerFactory) Model {
	input := textinput.New()
	input.Focus()
	input.Placeholder = placeholder
	input.PromptStyle = lipgloss.NewStyle().Foreground(color.MainAccent)
	return Model{
		input: input,
		rows:  rows.NewModel(ef),
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
	rows   rows.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleUpdateSize(msg)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case message.MoveFocus:
		return m.handleFocus()
	case message.SelectedContext,
		message.SelectedTable,
		message.FetchedRows:
		return m.delegateToRows(msg)
	default:
		return m, nil
	}
}

func (m Model) View() string {
	barStyles := m.newBarStyles()
	titleStyles := m.newTitleStyles()

	title := titleStyles.Render("Query prompt")
	input := inputStyles.Render(m.input.View())
	return barStyles.Render(title, lipgloss.JoinVertical(lipgloss.Top, input, m.rows.View()))
}

func (m Model) Help() string {
	return "searchbar help"
}

func (m Model) Value() string {
	return strings.TrimSpace(m.input.Value())
}

func (m Model) delegateToRows(msg tea.Msg) (Model, tea.Cmd) {
	slog.Info("Delegate to rows", slog.Any("msg", msg))
	rows, cmd := m.rows.Update(msg)
	m.rows = rows
	return m, cmd
}

func (m Model) handleFocus() (Model, tea.Cmd) {
	m.input.Focus()
	m.state.active = true
	return m, nil
}

func (m Model) handleMoveFocus(to direction.Direction) (Model, tea.Cmd) {
	m.input.Blur()
	m.state.active = false
	return m, message.With(message.MoveFocus{Direction: to})
}

func (m Model) handleUpdateSize(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height
	return m.delegateToRows(msg)
}

func (m Model) handleKeyEnter() (Model, tea.Cmd) {
	m.input.Blur()
	m.state.active = false
	val := m.Value()
	m.input.SetValue("")
	m.input.Placeholder = val
	return m, message.With(message.Command{Text: command.Command(val)})
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
