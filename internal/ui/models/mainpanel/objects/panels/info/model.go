package info

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/service/engine"
	"github.com/hrvadl/gowatchsql/internal/ui/color"
	"github.com/hrvadl/gowatchsql/internal/ui/message"
	"github.com/hrvadl/gowatchsql/internal/ui/styles"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

const margin = 1

func NewModel(ef ExplorerFactory) Model {
	item := list.NewDefaultDelegate()
	item.Styles = styles.NewForItemDelegate()
	l := newList(item)
	l.SetShowHelp(false)
	return Model{
		engineFactory: ef,
		list:          l,
	}
}

type ExplorerFactory interface {
	Create(dsn string) (engine.Explorer, error)
}

type Model struct {
	width  int
	height int

	list list.Model

	state state

	engineFactory ExplorerFactory
	explorer      engine.Explorer
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSize(msg.Width-margin*2, msg.Height-margin*2)
	case message.SelectedContext:
		return m.handleSelectedContext(msg)
	case message.Error:
		return m.handleError(msg)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case message.MoveFocus:
		return m.handleMoveFocus(msg)
	case message.FetchedTableList:
		return m.handleFetchedTableList(msg)
	default:
		return m, nil
	}
}

func (m Model) View() string {
	s := m.newStyles()
	switch m.state.status {
	case errored:
		return s.Render(m.state.err.Error())
	case loading:
		return s.Render("Loading...")
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
	default:
		return m.delegateToList(msg)
	}
}

func (m Model) handleSelectItem() (tea.Model, tea.Cmd) {
	chosen, ok := m.list.SelectedItem().(tableItem)
	if !ok {
		return m, nil
	}

	return m, func() tea.Msg { return message.SelectedTable{Name: chosen.Name} }
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
	m.state.status = errored
	return m, nil
}

func (m Model) handleFetchedTableList(msg message.FetchedTableList) (Model, tea.Cmd) {
	m.state.status = ready
	items := newItemsFromTable(msg.Tables)
	cmd := m.list.SetItems(items)
	return m, tea.Batch(cmd, m.commandSelectTable(msg.Tables))
}

func (m Model) handleSelectedContext(msg message.SelectedContext) (tea.Model, tea.Cmd) {
	explorer, err := m.engineFactory.Create(msg.DSN)
	if err != nil {
		m.state.err = err
		m.state.status = errored
		return m, nil
	}

	m.explorer = explorer
	m.state.status = loading
	return m, m.commandFetchTables
}

func (m Model) commandSelectTable(tables []engine.Table) tea.Cmd {
	if len(tables) == 0 {
		return nil
	}

	return func() tea.Msg {
		return message.SelectedTable{Name: tables[0].Name}
	}
}

func (m *Model) commandFetchTables() tea.Msg {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tables, err := m.explorer.GetTables(ctx)
	if err != nil {
		m.state.status = errored
		return message.Error{Err: err}
	}

	return message.FetchedTableList{Tables: tables}
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
	l.KeyMap.Quit = key.NewBinding(key.WithDisabled())

	return l
}
