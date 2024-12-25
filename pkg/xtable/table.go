package xtable

import (
	"math"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"

	"github.com/hrvadl/gowatchsql/internal/ui/color"
)

const (
	scrollLeft  = "h"
	scrollRight = "l"
)

type Model struct {
	base    table.Model
	columns []string
	rows    [][]string
	width   int
}

func New(cols []string, entries [][]string) Model {
	xtable := Model{}
	columns := toColumns(cols, xtable.getColumnWidth(entries, cols))
	rows := toRows(entries)

	keymap := table.DefaultKeyMap()
	keymap.ScrollLeft = key.NewBinding(key.WithKeys(scrollLeft))
	keymap.ScrollRight = key.NewBinding(key.WithKeys(scrollRight))

	table := table.New(columns).
		WithRows(rows).
		WithBaseStyle(newTableStyles()).
		WithRowStyleFunc(func(rsfi table.RowStyleFuncInput) lipgloss.Style {
			if !rsfi.IsHighlighted {
				return rsfi.Row.Style
			}

			return rsfi.Row.Style.Background(color.MainAccent)
		}).
		WithKeyMap(keymap).
		WithHorizontalFreezeColumnCount(1).
		Focused(true)

	xtable.base = table
	xtable.rows = entries
	xtable.columns = cols

	return xtable
}

func (t Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	base, cmd := t.base.Update(msg)
	t.base = base
	return t, cmd
}

func (t Model) Init() tea.Cmd {
	return t.base.Init()
}

func (t Model) View() string {
	return t.base.View()
}

func (t Model) KeyMap() table.KeyMap {
	return t.base.KeyMap()
}

func (t Model) WithTargetWidth(w int) Model {
	t.base = t.base.WithTargetWidth(w)
	return t
}

func (t Model) WithMaxTotalWidth(w int) Model {
	t.width = w
	columns := toColumns(t.columns, t.getColumnWidth(t.rows, t.columns))
	rows := toRows(t.rows)

	t.base = t.base.WithColumns(columns).WithRows(rows).WithMaxTotalWidth(w)
	return t
}

func (t Model) getColumnWidth(rows [][]string, columns []string) []int {
	var (
		widths     = make([]int, len(columns))
		totalWidth int
		flexPixels int
	)

	for i, col := range columns {
		widths[i] = len(col)
	}

	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	for _, col := range widths {
		totalWidth += col
	}

	if diff := t.width - totalWidth; diff > 0 {
		flexPixels = int(math.Round(float64(diff) / float64(len(columns))))
	}

	for i, w := range widths {
		widths[i] = w + flexPixels + 2
	}

	return widths
}

func toColumns(cols []string, width []int) []table.Column {
	t := make([]table.Column, 0)

	for i, v := range cols {
		t = append(t, table.NewColumn(strconv.Itoa(i), v, width[i]))
	}
	return t
}

func toRows(entries [][]string) []table.Row {
	rows := make([]table.Row, 0)

	for _, row := range entries {
		rowData := make(map[string]any)
		for i, data := range row {
			rowData[strconv.Itoa(i)] = data
		}
		rows = append(rows, table.NewRow(rowData))
	}

	return rows
}

func newTableStyles() lipgloss.Style {
	styleBase := lipgloss.NewStyle().
		Foreground(color.Text).
		BorderForeground(color.Border).
		Align(lipgloss.Right)

	return styleBase
}
