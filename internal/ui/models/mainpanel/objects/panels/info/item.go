package info

import (
	"github.com/charmbracelet/bubbles/list"

	"github.com/hrvadl/gowatchsql/internal/domain/engine"
)

func newItemsFromTable(t []engine.Table) []list.Item {
	items := make([]list.Item, 0, len(t))
	for _, tt := range t {
		items = append(items, tableItem{Table: tt})
	}
	return items
}

type tableItem struct {
	engine.Table
}

func (i tableItem) Title() string {
	return i.Name
}

func (i tableItem) Description() string {
	return i.Schema
}

func (i tableItem) FilterValue() string {
	return i.Name
}
