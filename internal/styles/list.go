package styles

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/color"
)

var titleStyle = lipgloss.NewStyle().
	Foreground(color.Text).
	Bold(true)

func NewForList() list.Styles {
	s := list.DefaultStyles()
	s.Title = titleStyle
	s.TitleBar = lipgloss.NewStyle().
		Bold(true).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(color.Border)
	return s
}

func NewForItemDelegate() list.DefaultItemStyles {
	st := list.NewDefaultItemStyles()
	st.SelectedTitle = st.SelectedTitle.Foreground(color.MainAccent).
		BorderForeground(color.MainAccent)

	st.SelectedDesc = st.SelectedDesc.Foreground(color.SecondaryAccent).
		BorderForeground(color.MainAccent)
	return st
}
