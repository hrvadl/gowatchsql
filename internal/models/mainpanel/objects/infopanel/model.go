package infopanel

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/color"
	"github.com/hrvadl/gowatchsql/internal/message"
	"github.com/hrvadl/gowatchsql/internal/service/sysexplorer"
	"github.com/hrvadl/gowatchsql/internal/styles"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

const margin = 1

func NewModel(ef ExplorerFactory) Model {
	item := list.NewDefaultDelegate()
	item.Styles = styles.NewForItemDelegate()
	l := newList(item)
	l.SetShowHelp(false)
	return Model{
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
	state  state

	explorerFactory ExplorerFactory
	explorer        Explorer
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSize(msg.Width-margin*2, msg.Height-margin*2)
	case message.SelectedDB:
		return m.handleSelectedDB(msg)
	case message.Error:
		return m.handleError(msg)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case message.MoveFocus:
		return m.handleMoveFocus(msg)
	default:
		return m, nil
	}
}

func (m Model) View() string {
	s := m.newStyles()
	switch m.state.status {
	case Error:
		return s.Render(m.state.err.Error())
	default:
		return s.Render(m.list.View())
	}
}

func (m Model) handleMoveFocus(msg message.MoveFocus) (tea.Model, tea.Cmd) {
	if msg.Direction == direction.Away {
		m.state.active = false
		return m, nil
	}

	m.state.active = !m.state.active
	return m, nil
}

func (m Model) Help() string {
	return "info help"
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		return m.handleSelectItem()
	case tea.KeyEsc:
		return m, nil
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
	lm, cmd := m.list.Update(msg)
	m.list = lm
	return m, cmd
}

func (m Model) handleWindowSize(w, h int) (tea.Model, tea.Cmd) {
	m.width = w
	m.height = h
	m.list.SetSize(w-3, h)
	m.list.Styles.TitleBar = m.list.Styles.TitleBar.Width(w)
	return m, nil
}

func (m Model) handleError(msg message.Error) (tea.Model, tea.Cmd) {
	m.state.err = msg.Err
	m.state.status = Error
	return m, nil
}

func (m Model) handleSelectedDB(msg message.SelectedDB) (tea.Model, tea.Cmd) {
	explorer, err := m.explorerFactory(msg.DSN)
	if err != nil {
		m.state.err = err
		m.state.status = Error
		return m, nil
	}

	m.explorer = explorer
	tables, err := m.explorer.GetTables()
	if err != nil {
		return m.handleError(message.Error{Err: err})
	}

	m.state.status = ready
	_ = m.list.SetItems(newItemsFromTable(tables))
	return m.handleSelectItem()
}

func (m Model) newStyles() lipgloss.Style {
	base := lipgloss.
		NewStyle().
		Height(m.height).
		Width(m.width).Border(lipgloss.NormalBorder())

	if m.state.active {
		return base.Border(lipgloss.ThickBorder()).
			BorderForeground(color.MainAccent)
	}

	return base.BorderForeground(color.Border)
}

func newList(item list.ItemDelegate) list.Model {
	const defaultTitle = "Tables 📋"

	l := list.New([]list.Item{}, item, 0, 0)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.InfiniteScrolling = true
	l.Styles = styles.NewForList()
	l.Title = defaultTitle

	return l
}
