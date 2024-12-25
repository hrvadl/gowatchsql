package xtable

import (
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"

	"github.com/hrvadl/gowatchsql/internal/ui/color"
)

const (
	scrollLeft  = "h"
	scrollRight = "l"
)

func New(cols []string, entries [][]string) table.Model {
	columns := toColumns(cols, getColumnWidth(entries, cols))
	rows := toRows(entries)

	keymap := table.DefaultKeyMap()
	keymap.ScrollLeft = key.NewBinding(key.WithKeys(scrollLeft))
	keymap.ScrollRight = key.NewBinding(key.WithKeys(scrollRight))

	return table.New(columns).
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
}

func toColumns(cols []string, width []int) []table.Column {
	t := make([]table.Column, 0)

	for i, v := range cols {
		//		if i == 0 {
		//			t = append(t, table.NewFlexColumn(strconv.Itoa(i), v, 2))
		//		} else {
		t = append(t, table.NewColumn(strconv.Itoa(i), v, 100))
		//}
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

func getColumnWidth(rows [][]string, columns []string) []int {
	widths := make([]int, len(columns))
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

	return widths
}

func newTableStyles() lipgloss.Style {
	styleBase := lipgloss.NewStyle().
		Foreground(color.Text).
		BorderForeground(color.Border).
		Align(lipgloss.Right)

	return styleBase
}
