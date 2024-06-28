package detailspanel

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/message"
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

	tableExplorer TableExplorer
	table         table.Model

	state State
	err   error
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleUpdateSize(msg.Width-margin*2, msg.Height-margin*3)
	case message.TableChosen:
		return m.handleTableChosen(msg)
	case message.Error:
		return m.handleError(msg)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	default:
		return m, nil
	}
}

func (m Model) View() string {
	s := lipgloss.
		NewStyle().
		Height(m.height).
		Width(m.width).
		MaxWidth(m.width)

	return s.Render(m.table.View())
}

func (m *Model) SetTableExplorer(te TableExplorer) {
	m.tableExplorer = te
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) handleTableChosen(msg message.TableChosen) (tea.Model, tea.Cmd) {
	rows, cols, err := m.tableExplorer.GetAll(msg.Name)
	if err != nil {
		return m.handleError(message.Error{Err: err})
	}

	t := table.New(
		table.WithColumns(m.mapToColumns(cols)),
		table.WithRows(m.mapToRows(rows)),
		table.WithFocused(true),
		table.WithWidth(m.width-10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	m.table = t

	return m, nil
}

func (m Model) handleError(msg message.Error) (tea.Model, tea.Cmd) {
	m.err = msg.Err
	m.state = Error
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
