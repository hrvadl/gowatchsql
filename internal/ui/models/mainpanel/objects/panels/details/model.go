package details

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/domain/engine"
	"github.com/hrvadl/gowatchsql/internal/ui/color"
	"github.com/hrvadl/gowatchsql/internal/ui/message"
	"github.com/hrvadl/gowatchsql/internal/ui/models/mainpanel/objects/panels/details/rows"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

const margin = 1

func NewModel(ef ExplorerFactory) Model {
	return Model{
		rows: rows.NewModel(ef),
	}
}

type ExplorerFactory interface {
	Create(ctx context.Context, name, dsn string) (engine.Explorer, error)
}

type Model struct {
	width  int
	height int

	state state
	rows  rows.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleUpdateSize(msg.Width-margin*2, msg.Height-margin*2)
	case message.SelectedContext:
		return m.delegateToAllModels(msg)
	case message.MoveFocus:
		return m.handleMoveFocus(msg)
	default:
		return m.delegateToActiveModel(msg)
	}
}

func (m Model) View() string {
	s := m.newContainerStyles()
	headerStyles := m.newHeaderStyles()
	header := headerStyles.Render(m.newTitle())

	var content string
	switch m.state.focused {
	case rowsFocused:
		content = m.rows.View()
	default:
		content = "Indexes"
	}

	return s.Render(lipgloss.JoinVertical(lipgloss.Top, header, content))
}

func (m Model) handleUpdateSize(width, height int) (Model, tea.Cmd) {
	m.width = width
	m.height = height
	return m.delegateToActiveModel(tea.WindowSizeMsg{Width: width - 5, Height: height - 5})
}

func (m Model) delegateToAllModels(msg tea.Msg) (Model, tea.Cmd) {
	m, rowsCmd := m.delegateToRowsModel(msg)
	return m, tea.Batch(rowsCmd)
}

func (m Model) delegateToActiveModel(msg tea.Msg) (Model, tea.Cmd) {
	switch m.state.focused {
	case rowsFocused:
		return m.delegateToRowsModel(msg)
	default:
		return m, nil
	}
}

func (m Model) delegateToRowsModel(msg tea.Msg) (Model, tea.Cmd) {
	rows, cmd := m.rows.Update(msg)
	m.rows = rows
	return m, cmd
}

func (m Model) newContainerStyles() lipgloss.Style {
	base := lipgloss.
		NewStyle().
		Height(m.height).
		Width(m.width).
		Border(lipgloss.NormalBorder())

	return base.BorderForeground(color.Border)
}

func (m Model) newHeaderStyles() lipgloss.Style {
	return lipgloss.NewStyle().
		Margin(1)
}

func (m Model) newTitle() string {
	return ""
}

func (m Model) handleMoveFocus(msg message.MoveFocus) (Model, tea.Cmd) {
	if msg.Direction == direction.Away {
		m.state.active = false
		return m, nil
	}

	m.state.active = !m.state.active
	return m, nil
}
