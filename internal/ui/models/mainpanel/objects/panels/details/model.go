package details

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/domain/engine"
	"github.com/hrvadl/gowatchsql/internal/ui/color"
	"github.com/hrvadl/gowatchsql/internal/ui/message"
	"github.com/hrvadl/gowatchsql/internal/ui/models/mainpanel/objects/panels/details/columns"
	"github.com/hrvadl/gowatchsql/internal/ui/models/mainpanel/objects/panels/details/rows"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

const margin = 1

func NewModel(ef ExplorerFactory) Model {
	return Model{
		rows:    rows.NewModel(ef),
		columns: columns.NewModel(ef),
	}
}

type ExplorerFactory interface {
	Create(ctx context.Context, name, dsn string) (engine.Explorer, error)
}

type Model struct {
	width  int
	height int

	state   state
	rows    rows.Model
	columns columns.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleUpdateSize(msg.Width-margin*2, msg.Height-margin*2)
	case message.SelectedContext, message.SelectedTable:
		return m.delegateToAllModels(msg)
	case message.MoveFocus:
		return m.handleMoveFocus(msg)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case message.FetchedRows:
		return m.delegateToRowsModel(msg)
	case message.FetchedColumns:
		return m.delegateToColumnsModel(msg)
	default:
		return m.delegateToActiveModel(msg)
	}
}

func (m Model) View() string {
	s := m.newContainerStyles()

	headerStyles := lipgloss.NewStyle().
		Width(m.width).
		BorderBottom(true).
		BorderForeground(color.Border)

	rowsTab := m.newTabStyles(m.state.focused == rowsFocused).Render("Rows")
	columnsTab := m.newTabStyles(m.state.focused == columnsFocused).Render("Columns")

	header := headerStyles.Render(lipgloss.JoinHorizontal(lipgloss.Left, rowsTab, columnsTab))

	var content string
	switch m.state.focused {
	case rowsFocused:
		content = m.rows.View()
	default:
		content = m.columns.View()
	}

	return s.Render(lipgloss.JoinVertical(lipgloss.Top, header, content))
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case moveFocusLeft:
		return m.handleMoveTabFocus(direction.Backwards)
	case moveFocusRight:
		return m.handleMoveTabFocus(direction.Forward)
	default:
		return m.delegateToActiveModel(msg)
	}
}

func (m Model) handleMoveTabFocus(to direction.Direction) (Model, tea.Cmd) {
	if to == direction.Forward && m.state.focused == columnsFocused {
		m.state.focused = rowsFocused
		return m, nil
	}

	if to == direction.Backwards && m.state.focused == rowsFocused {
		m.state.focused = columnsFocused
		return m, nil
	}

	if to == direction.Forward {
		m.state.focused++
		return m, nil
	}

	m.state.focused--
	return m, nil
}

func (m Model) handleUpdateSize(width, height int) (Model, tea.Cmd) {
	m.width = width
	m.height = height
	return m.delegateToAllModels(tea.WindowSizeMsg{Width: width, Height: height - 5})
}

func (m Model) delegateToAllModels(msg tea.Msg) (Model, tea.Cmd) {
	m, rowsCmd := m.delegateToRowsModel(msg)
	m, columnsCmd := m.delegateToColumnsModel(msg)
	return m, tea.Batch(rowsCmd, columnsCmd)
}

func (m Model) delegateToActiveModel(msg tea.Msg) (Model, tea.Cmd) {
	switch m.state.focused {
	case rowsFocused:
		return m.delegateToRowsModel(msg)
	case columnsFocused:
		return m.delegateToColumnsModel(msg)
	default:
		return m, nil
	}
}

func (m Model) delegateToColumnsModel(msg tea.Msg) (Model, tea.Cmd) {
	columns, cmd := m.columns.Update(msg)
	m.columns = columns
	return m, cmd
}

func (m Model) delegateToRowsModel(msg tea.Msg) (Model, tea.Cmd) {
	rows, cmd := m.rows.Update(msg)
	m.rows = rows
	return m, cmd
}

func (m Model) newTabStyles(active bool) lipgloss.Style {
	base := lipgloss.
		NewStyle().
		Width(20).
		Border(lipgloss.NormalBorder()).
		Align(lipgloss.Center)

	if active {
		return base.BorderForeground(color.MainAccent)
	}

	return base.BorderForeground(color.Border)
}

func (m Model) newContainerStyles() lipgloss.Style {
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

func (m Model) handleMoveFocus(msg message.MoveFocus) (Model, tea.Cmd) {
	if msg.Direction == direction.Away {
		m.state.active = false
		return m, nil
	}

	m.state.active = !m.state.active
	return m, nil
}
