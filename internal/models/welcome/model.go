package welcome

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/color"
	"github.com/hrvadl/gowatchsql/internal/message"
	"github.com/hrvadl/gowatchsql/internal/models/detailspanel"
	"github.com/hrvadl/gowatchsql/internal/models/infopanel"
	"github.com/hrvadl/gowatchsql/internal/models/search"
	"github.com/hrvadl/gowatchsql/internal/platform/db"
	"github.com/hrvadl/gowatchsql/internal/service/sysexplorer"
	"github.com/hrvadl/gowatchsql/internal/service/tableexplorer"
	"github.com/hrvadl/gowatchsql/pkg/overlay"
)

func NewModel(log *slog.Logger) Model {
	return Model{
		log:          log,
		searchbar:    search.New(),
		infopanel:    infopanel.NewModel(sysexplorer.New),
		detailspanel: detailspanel.NewModel(),
	}
}

type focus int

type state struct {
	focused   focus
	showModal bool
}

const (
	searchFocused focus = iota
	infoFocused
	detailsFocused
)

type Model struct {
	searchbar    search.Model
	infopanel    infopanel.Model
	detailspanel detailspanel.Model
	state        state

	height int
	width  int

	log *slog.Logger
}

func (m Model) Init() tea.Cmd {
	return nil
}

// TODO: add shift+tab
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleUpdateSize(msg)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case tea.QuitMsg:
		return nil, tea.Quit
	case message.DSNReady:
		return m.delegateToInfoPanel(msg)
	case message.TableChosen:
		m.detailspanel.SetTableExplorer(tableexplorer.New(db.Get()))
		return m.delegateToDetailsPanel(msg)
	}
	return m, nil
}

func (m Model) View() string {
	searchView, infoView, detailsView := m.getFocusedView()
	mainPane := lipgloss.JoinHorizontal(lipgloss.Left, infoView, detailsView)
	window := lipgloss.JoinVertical(lipgloss.Top, searchView, mainPane)
	popupStyles := m.newPopupStyles()

	if m.state.showModal {
		return overlay.Place(
			(m.width-popupStyles.GetWidth())/2,
			(m.height-popupStyles.GetHeight())/2,
			popupStyles.Render(m.getHelpPopupContent()),
			window,
			true,
		)
	}

	return window
}

func (m Model) handleUpdateSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	const (
		searchBarHeight = 4
		infoPanelWidth  = 20
	)

	m.height = msg.Height
	m.width = msg.Width

	searchbar, cmd1 := m.searchbar.Update(tea.WindowSizeMsg{
		Width:  msg.Width,
		Height: searchBarHeight,
	})
	m.searchbar = searchbar.(search.Model)

	panel, cmd := m.infopanel.Update(tea.WindowSizeMsg{
		Width:  infoPanelWidth,
		Height: msg.Height - searchBarHeight,
	})
	m.infopanel = panel.(infopanel.Model)

	panel, cmd2 := m.detailspanel.Update(tea.WindowSizeMsg{
		Width:  msg.Width - infoPanelWidth,
		Height: msg.Height - searchBarHeight,
	})
	m.detailspanel = panel.(detailspanel.Model)

	return m, tea.Batch(cmd1, cmd, cmd2)
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return nil, tea.Quit
	case tea.KeyTab:
		return m.handleMoveFocusForward()
	case tea.KeyShiftTab:
		return m.handleMoveFocusBackwards()
	case tea.KeyEsc:
		return m.handleHidePopup()
	default:
		return m.handleKeyRunes(msg)
	}
}

func (m Model) handleKeyRunes(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "?":
		return m.handleShowPopup()
	default:
		return m.delegateKeyPressHandler(msg)
	}
}

func (m Model) handleMoveFocusBackwards() (tea.Model, tea.Cmd) {
	if m.state.showModal {
		return m, nil
	}

	switch m.state.focused {
	case searchFocused:
		m.state.focused = 2
		m.searchbar.Unfocus()
	case infoFocused:
		m.state.focused--
		m.searchbar.Focus()
	case detailsFocused:
		m.searchbar.Unfocus()
		m.state.focused--
	}

	return m, nil
}

func (m Model) handleMoveFocusForward() (tea.Model, tea.Cmd) {
	if m.state.showModal {
		return m, nil
	}

	switch m.state.focused {
	case searchFocused:
		m.state.focused++
		m.searchbar.Unfocus()
	case infoFocused:
		m.state.focused++
		m.searchbar.Unfocus()
	case detailsFocused:
		m.searchbar.Focus()
		m.state.focused = 0
	}

	return m, nil
}

func (m *Model) delegateKeyPressHandler(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state.focused {
	case searchFocused:
		return m.delegateToSearchbar(msg)
	case infoFocused:
		return m.delegateToInfoPanel(msg)
	case detailsFocused:
		return m.delegateToDetailsPanel(msg)
	}

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

func (m Model) handleShowPopup() (tea.Model, tea.Cmd) {
	m.state.showModal = true
	m.searchbar.Unfocus()
	return m, nil
}

func (m Model) handleHidePopup() (tea.Model, tea.Cmd) {
	m.state.showModal = false

	switch m.state.focused {
	case searchFocused:
		m.searchbar.Focus()
		return m, nil
	default:
		return m, nil
	}
}

func (m Model) delegateToSearchbar(msg tea.Msg) (tea.Model, tea.Cmd) {
	sb, cmd := m.searchbar.Update(msg)
	m.searchbar = sb.(search.Model)
	return m, cmd
}

func (m Model) delegateToInfoPanel(msg tea.Msg) (tea.Model, tea.Cmd) {
	ifp, cmd := m.infopanel.Update(msg)
	m.infopanel = ifp.(infopanel.Model)
	return m, cmd
}

func (m Model) delegateToDetailsPanel(msg tea.Msg) (tea.Model, tea.Cmd) {
	dp, cmd := m.detailspanel.Update(msg)
	m.detailspanel = dp.(detailspanel.Model)
	return m, cmd
}

func (m Model) getFocusedView() (string, string, string) {
	bordered := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(color.Border)
	focused := bordered.
		BorderForeground(color.MainAccent)

	searchView := m.searchbar.View()
	infoView := m.infopanel.View()
	detailsView := m.detailspanel.View()

	if m.state.showModal {
		return bordered.Render(searchView),
			bordered.Render(infoView),
			bordered.Render(detailsView)
	}

	switch m.state.focused {
	case searchFocused:
		searchView = focused.Render(searchView)
		infoView = bordered.Render(infoView)
		detailsView = bordered.Render(detailsView)
	case infoFocused:
		infoView = focused.Render(infoView)
		searchView = bordered.Render(searchView)
		detailsView = bordered.Render(detailsView)
	case detailsFocused:
		detailsView = focused.Render(detailsView)
		searchView = bordered.Render(searchView)
		infoView = bordered.Render(infoView)
	}

	return searchView, infoView, detailsView
}

func (m Model) getHelpPopupContent() string {
	switch m.state.focused {
	case searchFocused:
		return m.searchbar.Help()
	case infoFocused:
		return m.infopanel.Help()
	case detailsFocused:
		return m.detailspanel.Help()
	default:
		return ""
	}
}

func (m Model) newPopupStyles() lipgloss.Style {
	return lipgloss.NewStyle().
		Width(m.width-15).
		Height(m.height-15).
		Border(lipgloss.NormalBorder(), true).
		BorderForeground(color.MainAccent)
}
