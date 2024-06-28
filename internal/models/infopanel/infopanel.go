package infopanel

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/service/sysexplorer"
)

const (
	padding = 4
	margin  = 1
)

func NewModel(ef ExplorerFactory) Model {
	return Model{
		state:           Pending,
		explorerFactory: ef,
	}
}

type ExplorerFactory = func(dsn string) (*sysexplorer.Explorer, error)

type Explorer interface {
	GetTables() ([]sysexplorer.Table, error)
}

type Model struct {
	width  int
	height int

	tables []sysexplorer.Table

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
	case DSNReadyMsg:
		return m.handleDSNReady(msg)
	case ErrorMsg:
		return m.handleError(msg)
	default:
		return m, nil
	}
}

func (m Model) View() string {
	s := lipgloss.
		NewStyle().
		Height(m.height).
		Width(m.width).
		Border(lipgloss.NormalBorder())

	switch m.state {
	case Ready:
		return s.Render(m.tables[0].Name)
	case Error:
		return s.Render(m.err.Error())
	default:
		return s.Render()
	}
}

func (m Model) handleUpdateSize(w, h int) (tea.Model, tea.Cmd) {
	m.width = w
	m.height = h
	return m, nil
}

func (m Model) handleError(msg ErrorMsg) (tea.Model, tea.Cmd) {
	m.err = msg.Err
	m.state = Error
	return m, nil
}

func (m Model) handleDSNReady(msg DSNReadyMsg) (tea.Model, tea.Cmd) {
	explorer, err := m.explorerFactory(msg.DSN)
	if err != nil {
		m.err = err
		m.state = Error
		return m, nil
	}

	m.explorer = explorer
	tables, err := m.explorer.GetTables()
	if err != nil {
		return m.handleError(ErrorMsg{err})
	}

	m.state = Ready
	m.tables = tables
	return m, nil
}
