package detailspanel

import (
	"math"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/ui/color"
	"github.com/hrvadl/gowatchsql/internal/ui/message"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

const margin = 1

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
	case message.SelectedTable:
		return m.handleTableChosen(msg)
	case message.Error:
		return m.handleError(msg)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case message.MoveFocus:
		return m.handleMoveFocus(msg)
	case message.FetchedTableContent:
		return m.handleFetchedTableContent(msg)
	default:
		return m, nil
	}
}

func (m Model) View() string {
	s := m.newContainerStyles()
	headerStyles := m.newHeaderStyles()
	header := headerStyles.Render(m.newTitle())

	content := m.table.View()
	switch m.state.status {
	case loading:
		content = "Loading..."
	case errored:
		content = m.err.Error()
	}

	return s.Render(lipgloss.JoinVertical(lipgloss.Top, header, content))
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

func (m Model) handleFetchedTableContent(msg message.FetchedTableContent) (Model, tea.Cmd) {
	m.state.status = ready
	m.table = table.New(
		table.WithColumns(m.mapToColumns(msg.Cols)),
		table.WithRows(m.mapToRows(msg.Rows)),
		table.WithFocused(true),
		table.WithWidth(m.width-1),
		table.WithHeight(m.height-10),
		table.WithStyles(m.newTableStyles()),
	)
	return m, nil
}

func (m Model) handleTableChosen(msg message.SelectedTable) (tea.Model, tea.Cmd) {
	m.chosen = msg.Name
	m.state.status = loading
	return m, m.commandFetchTableContent(msg.Name)
}

func (m Model) handleError(msg message.Error) (tea.Model, tea.Cmd) {
	m.err = msg.Err
	m.state.status = errored
	return m, nil
}

func (m Model) handleUpdateSize(w, h int) (tea.Model, tea.Cmd) {
	m.width = w
	m.height = h
	m.table.SetWidth(w - 1)
	m.table.SetHeight(h - 10)
	return m, nil
}

func (m *Model) commandFetchTableContent(table string) tea.Cmd {
	return func() tea.Msg {
		rows, cols, err := m.tableExplorer.GetAll(table)
		if err != nil {
			m.state.status = errored
			return message.Error{Err: err}
		}

		return message.FetchedTableContent{Rows: rows, Cols: cols}
	}
}

func (m Model) mapToColumns(cols []string) []table.Column {
	t := make([]table.Column, 0)
	width := round(m.width, len(cols))
	for _, v := range cols {
		t = append(t, table.Column{
			Title: v,
			Width: width,
		})
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
	s.Cell = s.Cell.MaxWidth(m.width + 1)
	s.Header = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(color.Border).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(color.Text).
		Background(color.MainAccent).
		Bold(true)
	return s
}

func (m Model) newTitle() string {
	const base = "Table "
	if m.chosen == "" {
		return base + "📃"
	}
	return base + m.chosen + " 📃"
}

func round(i, ii int) int {
	return int(math.Round(float64(i) / float64(ii)))
}
