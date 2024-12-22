package columns

import (
	"context"
	"log/slog"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	xtable "github.com/evertras/bubble-table/table"

	"github.com/hrvadl/gowatchsql/internal/domain/engine"
	"github.com/hrvadl/gowatchsql/internal/ui/color"
	"github.com/hrvadl/gowatchsql/internal/ui/message"
)

const margin = 1

type Column = string

type Row = []string

func NewModel(factory ExplorerFactory) Model {
	return Model{
		engineFactory: factory,
	}
}

type ExplorerFactory interface {
	Create(ctx context.Context, name, dsn string) (engine.Explorer, error)
}

type Model struct {
	width  int
	height int

	chosen        string
	engineFactory ExplorerFactory
	explorer      engine.Explorer
	table         xtable.Model

	state state
	err   error
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleUpdateSize(msg.Width-margin*2, msg.Height-margin*2)
	case message.SelectedTable:
		return m.handleTableChosen(msg)
	case message.Error:
		return m.handleError(msg)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case message.FetchedColumns:
		return m.handleFetchedTableContent(msg)
	case message.SelectedContext:
		return m.handleSelectedContext(msg)
	default:
		return m.delegateToTable(msg)
	}
}

func (m Model) View() string {
	s := m.newContainerStyles()

	content := m.table.View()
	switch m.state.status {
	case loading:
		content = "Loading..."
	case errored:
		content = m.err.Error()
	}

	return s.Render(content)
}

func (m Model) Help() string {
	return "Details help"
}

func (m Model) delegateToTable(msg tea.Msg) (Model, tea.Cmd) {
	table, cmd := m.table.Update(msg)
	m.table = table
	return m, cmd
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) handleFetchedTableContent(msg message.FetchedColumns) (Model, tea.Cmd) {
	m.state.status = ready
	m.table = xtable.New(m.mapToColumns(msg.Cols)).
		WithRows(m.mapToRows(msg.Rows)).
		WithRowStyleFunc(func(rsfi xtable.RowStyleFuncInput) lipgloss.Style {
			if !rsfi.IsHighlighted {
				return rsfi.Row.Style
			}

			return rsfi.Row.Style.Background(color.MainAccent)
		}).
		WithMaxTotalWidth(m.width - 1).
		WithHorizontalFreezeColumnCount(1).
		WithBaseStyle(m.newTableStyles()).
		Focused(true)

	return m, nil
}

func (m Model) handleSelectedContext(msg message.SelectedContext) (Model, tea.Cmd) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	explorer, err := m.engineFactory.Create(ctx, msg.Name, msg.DSN)
	if err != nil {
		m.err = err
		m.state.status = errored
		return m, nil
	}

	m.explorer = explorer
	m.state.status = loading

	return m, nil
}

func (m Model) handleTableChosen(msg message.SelectedTable) (Model, tea.Cmd) {
	m.chosen = msg.Name
	m.state.status = loading
	return m, m.commandFetchTableContent(msg.Name)
}

func (m Model) handleError(msg message.Error) (Model, tea.Cmd) {
	m.err = msg.Err
	m.state.status = errored
	return m, nil
}

func (m Model) handleUpdateSize(w, h int) (Model, tea.Cmd) {
	m.width = w
	m.height = h
	m.table = m.table.WithMaxTotalWidth(w - 1)
	return m, nil
}

func (m *Model) commandFetchTableContent(table string) tea.Cmd {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	return func() tea.Msg {
		defer cancel()

		rows, cols, err := m.explorer.GetColumns(ctx, table)
		if err != nil {
			m.state.status = errored
			slog.Error("Got error fetching table content", slog.Any("err", err))
			return message.Error{Err: err}
		}

		return message.FetchedColumns{Rows: rows, Cols: cols}
	}
}

func (m Model) mapToColumns(cols []string) []xtable.Column {
	t := make([]xtable.Column, 0)

	var useDefaultWidth bool
	if len(strings.Join(cols, "")) < m.width {
		useDefaultWidth = true
	}

	for i, v := range cols {
		width := len(v)
		if useDefaultWidth {
			width = m.width / len(cols)
		}

		t = append(t, xtable.NewColumn(strconv.Itoa(i), v, width))
	}
	return t
}

func (m Model) mapToRows(entries []Row) []xtable.Row {
	rows := make([]xtable.Row, 0)

	for _, row := range entries {
		rowData := make(map[string]any)
		for i, data := range row {
			rowData[strconv.Itoa(i)] = data
		}
		rows = append(rows, xtable.NewRow(rowData))
	}

	return rows
}

func (m Model) newContainerStyles() lipgloss.Style {
	base := lipgloss.
		NewStyle().
		Height(m.height).
		Width(m.width)

	return base.BorderForeground(color.Border)
}

func (m Model) newTableStyles() lipgloss.Style {
	styleBase := lipgloss.NewStyle().
		Foreground(color.Text).
		BorderForeground(color.Border).
		Align(lipgloss.Right)

	return styleBase
}
