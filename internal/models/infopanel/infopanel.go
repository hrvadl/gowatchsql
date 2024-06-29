package infopanel

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/color"
	"github.com/hrvadl/gowatchsql/internal/message"
	"github.com/hrvadl/gowatchsql/internal/service/sysexplorer"
)

const (
	padding = 4
	margin  = 1
)

var titleStyle = lipgloss.NewStyle().
	Foreground(color.Text).
	Bold(true)

func NewModel(ef ExplorerFactory) Model {
	item := list.NewDefaultDelegate()
	setupItemStyles(&item.Styles)
	l := newList(item)
	l.SetShowHelp(false)
	return Model{
		state:           Pending,
		explorerFactory: ef,
		list:            l,
	}
}

type ExplorerFactory = func(dsn string) (*sysexplorer.Explorer, error)

type Explorer interface {
	GetTables() ([]sysexplorer.Table, error)
}

type Model struct {
	width  int
	height int

	list list.Model

	chosen string
	state  State
	err    error

	explorerFactory ExplorerFactory
	explorer        Explorer
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleUpdateSize(msg.Width-margin*2, msg.Height-margin*3)
	case message.DSNReady:
		return m.handleDSNReady(msg)
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
		Width(m.width)

	switch m.state {
	case Error:
		return s.Render(m.err.Error())
	default:
		return s.Render(m.list.View())
	}
}

func (m Model) Help() string {
	return "info help"
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		return m.handleSelectItem()
	default:
		return m.delegateToList(msg)
	}
}

func (m Model) handleSelectItem() (tea.Model, tea.Cmd) {
	chosen := m.list.SelectedItem().(tableItem)
	m.chosen = chosen.Table.Name
	return m, func() tea.Msg { return message.TableChosen{Name: m.chosen} }
}

func (m Model) delegateToList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	lm, err := m.list.Update(msg)
	if err != nil {
		panic("failed to update list")
	}

	m.list = lm
	return m, nil
}

func (m Model) handleUpdateSize(w, h int) (tea.Model, tea.Cmd) {
	m.width = w
	m.height = h
	m.list.SetSize(w-3, h)
	m.list.Styles.TitleBar = m.list.Styles.TitleBar.Width(w)
	return m, nil
}

func (m Model) handleError(msg message.Error) (tea.Model, tea.Cmd) {
	m.err = msg.Err
	m.state = Error
	return m, nil
}

func (m Model) handleDSNReady(msg message.DSNReady) (tea.Model, tea.Cmd) {
	explorer, err := m.explorerFactory(msg.DSN)
	if err != nil {
		m.err = err
		m.state = Error
		return m, nil
	}

	m.explorer = explorer
	tables, err := m.explorer.GetTables()
	if err != nil {
		return m.handleError(message.Error{Err: err})
	}

	m.state = Ready
	_ = m.list.SetItems(newItemsFromTable(tables))
	return m.handleSelectItem()
}

func setupItemStyles(st *list.DefaultItemStyles) {
	st.SelectedTitle = st.SelectedTitle.Foreground(color.MainAccent).
		BorderForeground(color.MainAccent)

	st.SelectedDesc = st.SelectedDesc.Foreground(color.SecondaryAccent).
		BorderForeground(color.MainAccent)
}

func newList(item list.ItemDelegate) list.Model {
	const defaultTitle = "Tables"

	l := list.New([]list.Item{}, item, 0, 0)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowPagination(false)
	l.Styles.Title = titleStyle

	l.InfiniteScrolling = true
	l.Title = defaultTitle

	l.Styles.TitleBar = lipgloss.NewStyle().
		Bold(true).
		Foreground(color.Text).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(color.Border)

	return l
}
