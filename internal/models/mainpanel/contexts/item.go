package contexts

import (
	"github.com/charmbracelet/bubbles/list"

	"github.com/hrvadl/gowatchsql/internal/message"
)

func newItemFromContext(ctx message.NewContext) list.Item {
	return ctxItem{
		name: ctx.Name,
		dsn:  ctx.DSN,
	}
}

type ctxItem struct {
	name string
	dsn  string
}

func (i ctxItem) Title() string       { return i.name }
func (i ctxItem) Description() string { return i.dsn }
func (i ctxItem) FilterValue() string { return i.name }
