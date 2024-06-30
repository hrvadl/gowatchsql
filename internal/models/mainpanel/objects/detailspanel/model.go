package detailspanel

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/color"
	"github.com/hrvadl/gowatchsql/internal/message"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

const (
	padding = 1
	margin  = 1
)

func NewModel() Model {
	return Model{}
}

type Column = string

type Row = []string

type TableExplorer interface {
	GetAll(string) ([]Row, []Column, error)
}

type Model struct {
	width  int
	height int

	chosen        string
	tableExplorer TableExplorer
	table         table.Model

	state state
	err   error
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleUpdateSize(msg.Width-margin*2, msg.Height-margin*2)
	case message.TableChosen:
		return m.handleTableChosen(msg)
	case message.Error:
		return m.handleError(msg)
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
	headerStyles := m.newHeaderStyles()
	header := headerStyles.Render(m.newTitle())
	return s.Render(lipgloss.JoinVertical(lipgloss.Top, header, m.table.View()))
}

func (m *Model) SetTableExplorer(te TableExplorer) {
	m.tableExplorer = te
}

func (m Model) Help() string {
	return "Details help"
}

func (m Model) handleMoveFocus(msg message.MoveFocus) (tea.Model, tea.Cmd) {
	if msg.Direction == direction.Away {
		m.state.active = false
		return m, nil
	}

	m.state.active = !m.state.active
	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) handleTableChosen(msg message.TableChosen) (tea.Model, tea.Cmd) {
	m.chosen = msg.Name
	rows, cols, err := m.tableExplorer.GetAll(msg.Name)
	if err != nil {
		return m.handleError(message.Error{Err: err})
	}

	m.table = table.New(
		table.WithColumns(m.mapToColumns(cols)),
		table.WithRows(m.mapToRows(rows)),
		table.WithFocused(true),
		table.WithWidth(m.width-10),
		table.WithHeight(m.height-10),
		table.WithStyles(m.newTableStyles()),
	)

	slog.Info("created table")
	return m, nil
}

func (m Model) handleError(msg message.Error) (tea.Model, tea.Cmd) {
	m.err = msg.Err
	m.state.status = errored
	return m, nil
}

func (m Model) handleUpdateSize(w, h int) (tea.Model, tea.Cmd) {
	m.width = w
	m.height = h
	return m, nil
}

func (m Model) mapToColumns(cols []string) []table.Column {
	t := make([]table.Column, 0)
	for _, v := range cols {
		t = append(t, table.Column{Title: v, Width: m.width / len(cols)})
	}
	return t
}

func (m Model) mapToRows(entries []Row) []table.Row {
	rows := make([]table.Row, 0)

	for _, row := range entries {
		rows = append(rows, table.Row(row))
	}

	return rows
}

func (m Model) newHeaderStyles() lipgloss.Style {
	return lipgloss.
		NewStyle().
		Width(m.width).
		Bold(true).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(color.Border)
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

func (m Model) newTableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(color.Border).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(color.Text).
		Background(color.MainAccent).
		Bold(false)
	return s
}

func (m Model) newTitle() string {
	const base = "Table "
	if m.chosen == "" {
		return base + "📃"
	}
	return base + m.chosen + " 📃"
}
