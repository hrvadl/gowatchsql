package queryrun

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/domain/engine"
	"github.com/hrvadl/gowatchsql/internal/ui/color"
	"github.com/hrvadl/gowatchsql/internal/ui/message"
	"github.com/hrvadl/gowatchsql/internal/ui/models/mainpanel/objects/panels/details/rows"
)

const (
	padding     = 1
	margin      = 1
	placeholder = "SELECT * FROM"
)

var inputStyles = lipgloss.NewStyle().MarginTop(margin).PaddingRight(1).Foreground(color.Border)

type ExplorerFactory interface {
	Create(ctx context.Context, name, dsn string) (engine.Explorer, error)
}

func NewModel(ef ExplorerFactory) Model {
	input := textinput.New()
	input.Focus()
	input.Placeholder = placeholder
	input.PromptStyle = lipgloss.NewStyle().Foreground(color.MainAccent)

	return Model{
		input:           input,
		rows:            rows.NewModel(ef),
		explorerFactory: ef,
	}
}

type Model struct {
	width           int
	height          int
	state           state
	input           textinput.Model
	explorerFactory ExplorerFactory
	explorer        engine.Explorer
	rows            rows.Model
	table           string
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
		return m.handleMoveFocus()
	case message.SelectedTable:
		m.state.err = nil
		m.table = msg.Name
		return m.delegateToRows(msg)
	case message.SelectedContext:
		return m.handleSelectedContext(msg)
	case message.FetchedRows:
		return m.delegateToRows(msg)
	case message.Error:
		return m.handleError(msg)
	default:
		return m.delegateToActiveModel(msg)
	}
}

func (m Model) View() string {
	barStyles := m.newBarStyles()
	titleStyles := m.newTitleStyles()

	titleText := "Query prompt"
	if m.table != "" {
		titleText += " - " + m.table
	}

	inputStyles := inputStyles
	m.input.TextStyle = m.input.TextStyle.Foreground(color.Border)
	if m.state.focused == promptFocused {
		m.input.TextStyle = m.input.TextStyle.Foreground(color.Text)
	}

	slog.Info("Rendering query view", slog.Any("err", m.state.err))
	if m.state.err != nil {
		m.input.Placeholder = m.state.err.Error()
		m.input.PlaceholderStyle = m.input.PlaceholderStyle.Foreground(color.Error)
	}

	title := titleStyles.Render(titleText)
	input := inputStyles.Render(m.input.View())
	return barStyles.Render(title, lipgloss.JoinVertical(lipgloss.Top, input, m.rows.View()))
}

func (m Model) Help() string {
	return "searchbar help"
}

func (m Model) Value() string {
	return strings.TrimSpace(m.input.Value())
}

func (m Model) handleSelectedContext(msg message.SelectedContext) (Model, tea.Cmd) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	explorer, err := m.explorerFactory.Create(ctx, msg.Name, msg.DSN)
	if err != nil {
		return m, nil
	}

	m.explorer = explorer

	return m.delegateToRows(msg)
}

func (m Model) delegateToActiveModel(msg tea.Msg) (Model, tea.Cmd) {
	switch m.state.focused {
	case promptFocused:
		return m.delegateToPrompt(msg)
	case tableFocused:
		return m.delegateToRows(msg)
	}

	return m, nil
}

func (m Model) delegateToAllModels(msg tea.Msg) (Model, tea.Cmd) {
	m, rowsCmd := m.delegateToRows(msg)
	m, inputCmd := m.delegateToPrompt(msg)
	return m, tea.Batch(rowsCmd, inputCmd)
}

func (m Model) delegateToPrompt(msg tea.Msg) (Model, tea.Cmd) {
	input, cmd := m.input.Update(msg)
	m.input = input
	return m, cmd
}

func (m Model) delegateToRows(msg tea.Msg) (Model, tea.Cmd) {
	rows, cmd := m.rows.Update(msg)
	m.rows = rows
	return m, cmd
}

func (m Model) handleFocus() (Model, tea.Cmd) {
	if m.state.focused == promptFocused {
		m.input.Focus()
	}

	m.state.active = true
	return m, nil
}

func (m Model) handleUnfocus() (Model, tea.Cmd) {
	m.input.Blur()
	m.state.active = false
	return m, nil
}

func (m Model) handleMoveTabFocus() (Model, tea.Cmd) {
	if m.state.focused == promptFocused {
		m.state.focused = tableFocused
		m.input.Blur()
		return m, nil
	}

	m.input.Focus()
	m.state.focused = promptFocused
	return m, nil
}

func (m Model) handleMoveFocus() (Model, tea.Cmd) {
	if !m.state.active {
		return m.handleFocus()
	}

	return m.handleUnfocus()
}

func (m Model) handleUpdateSize(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m.width = msg.Width - 2
	m.height = msg.Height - 2

	m.input.Width = msg.Width - 10
	return m.delegateToAllModels(tea.WindowSizeMsg{Height: msg.Height - 5, Width: msg.Width - 5})
}

func (m Model) handleError(msg message.Error) (Model, tea.Cmd) {
	m.state.err = msg.Err
	return m, nil
}

func (m Model) handleKeyEnter() (Model, tea.Cmd) {
	m.input.Blur()
	query := m.Value()

	if m.state.focused != promptFocused {
		return m, nil
	}

	m.input.SetValue("")
	m.input.Placeholder = query
	m.state.focused = tableFocused

	return m, func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := m.explorer.Execute(ctx, query); err != nil {
			slog.Error("Execute query", slog.Any("err", err), slog.Any("state", m.state))
			return message.Error{Err: err}
		}

		return message.SelectedTable{Name: m.table}
	}
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return Model{}, tea.Quit
	case tea.KeyEsc:
		return m.handleMoveFocus()
	case tea.KeyShiftTab, tea.KeyTab:
		return m.handleMoveTabFocus()
	case tea.KeyEnter:
		return m.handleKeyEnter()
	default:
		return m.delegateToActiveModel(msg)
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
