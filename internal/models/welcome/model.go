package welcome

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/models/detailspanel"
	"github.com/hrvadl/gowatchsql/internal/models/infopanel"
	"github.com/hrvadl/gowatchsql/internal/models/search"
	"github.com/hrvadl/gowatchsql/internal/service/sysexplorer"
)

func NewModel(log *slog.Logger) Model {
	return Model{
		log:          log,
		searchbar:    search.New(),
		infopanel:    infopanel.NewModel(sysexplorer.New),
		detailspanel: detailspanel.NewModel(),
	}
}

type Model struct {
	searchbar    search.Model
	infopanel    infopanel.Model
	detailspanel detailspanel.Model

	log *slog.Logger
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleUpdateSize(msg)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case tea.QuitMsg:
		return nil, tea.Quit
	}
	return m, nil
}

func (m Model) View() string {
	// root:secret@(0.0.0.0:3306)/test
	mainPane := lipgloss.JoinHorizontal(lipgloss.Left, m.infopanel.View(), m.detailspanel.View())
	return lipgloss.JoinVertical(lipgloss.Top, m.searchbar.View(), mainPane)
}

func (m *Model) handleUpdateSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	const (
		searchBarHeight = 3
		infoPanelWidth  = 20
	)

	searchbar, cmd1 := m.searchbar.Update(tea.WindowSizeMsg{
		Width:  msg.Width,
		Height: searchBarHeight,
	})
	search, ok := searchbar.(search.Model)
	if !ok {
		panic("unknown model")
	}

	m.searchbar = search

	panel, cmd := m.infopanel.Update(tea.WindowSizeMsg{
		Width:  infoPanelWidth,
		Height: msg.Height - searchBarHeight,
	})
	infopanel, ok := panel.(infopanel.Model)
	if !ok {
		panic("unknown model")
	}

	m.infopanel = infopanel

	panel, cmd2 := m.detailspanel.Update(tea.WindowSizeMsg{
		Width:  msg.Width - infoPanelWidth,
		Height: msg.Height - searchBarHeight,
	})

	detailspanel, ok := panel.(detailspanel.Model)
	if !ok {
		panic("unknown model")
	}

	m.detailspanel = detailspanel

	return m, tea.Batch(cmd1, cmd, cmd2)
}

func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return nil, tea.Quit
	case tea.KeyEnter:
		return m.handleEnterKey(msg)
	default:
		return m.delegateKeyPressHandler(msg)
	}
}

func (m *Model) delegateKeyPressHandler(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if !m.searchbar.IsFocused() {
		return m, nil
	}

	searchbar, cmd := m.searchbar.Update(msg)
	search, ok := searchbar.(search.Model)
	if !ok {
		panic("undefined model")
	}

	m.searchbar = search
	return m, cmd
}

func (m *Model) handleEnterKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if !m.searchbar.IsFocused() {
		return m, nil
	}

	searchbar, cmd := m.searchbar.Update(msg)
	search, ok := searchbar.(search.Model)
	if !ok {
		panic("undefined model")
	}

	m.searchbar = search

	panel, cmd2 := m.infopanel.Update(infopanel.DSNReadyMsg{DSN: m.searchbar.Value()})
	infopanel, ok := panel.(infopanel.Model)
	if !ok {
		panic("unknown model")
	}

	m.infopanel = infopanel

	return m, tea.Batch(cmd, cmd2)
}
