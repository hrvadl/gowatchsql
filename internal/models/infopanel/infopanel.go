package infopanel

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/message"
	"github.com/hrvadl/gowatchsql/internal/service/sysexplorer"
)

const (
	padding = 4
	margin  = 1
)

func NewModel(ef ExplorerFactory) Model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowPagination(false)

	l.Title = "Tables"

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

	state State
	err   error

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
	case Ready:
		return s.Render(m.list.View())
	case Error:
		return s.Render(m.err.Error())
	default:
		return s.Render()
	}
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	cmd := m.list.SetItems(newItemsFromTable(tables))

	loadTableCmd := func() tea.Msg {
		if len(tables) == 0 {
			return nil
		}
		return message.TableChosen{Name: tables[0].Name}
	}

	return m, tea.Batch(cmd, loadTableCmd)
}
