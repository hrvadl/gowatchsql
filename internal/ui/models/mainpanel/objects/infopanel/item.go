package infopanel

import (
	"github.com/charmbracelet/bubbles/list"

	"github.com/hrvadl/gowatchsql/internal/service/sysexplorer"
)

func newItemsFromTable(t []sysexplorer.Table) []list.Item {
	items := make([]list.Item, 0, len(t))
	for _, tt := range t {
		items = append(items, tableItem{Table: tt})
	}
	return items
}

type tableItem struct {
	sysexplorer.Table
}

func (i tableItem) Title() string       { return i.Name }
func (i tableItem) Description() string { return i.Schema }
func (i tableItem) FilterValue() string { return i.Name }
