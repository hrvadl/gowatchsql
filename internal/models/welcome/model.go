package welcome

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/color"
	"github.com/hrvadl/gowatchsql/internal/message"
	"github.com/hrvadl/gowatchsql/internal/models/command"
	"github.com/hrvadl/gowatchsql/internal/models/mainpanel"
	"github.com/hrvadl/gowatchsql/pkg/overlay"
)

func NewModel(log *slog.Logger) Model {
	return Model{
		log:     log,
		command: command.NewModel(),
		main:    mainpanel.NewModel(),
	}
}

type focus int

type state struct {
	active    focus
	showModal bool
}

const (
	cmdFocused focus = iota
	mainFocused
)

type Model struct {
	command command.Model
	main    mainpanel.Model

	state state

	height int
	width  int

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
	case message.MoveFocus:
		return m.handleMoveFocus(msg)
	case message.DSNReady, message.TableChosen:
		return m.delegateToMainPanel(msg)
	}
	return m, nil
}

func (m Model) View() string {
	window := lipgloss.JoinVertical(lipgloss.Top, m.command.View(), m.main.View())
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
	const searchBarHeight = 4

	m.height = msg.Height
	m.width = msg.Width

	searchbar, searchCmd := m.command.Update(tea.WindowSizeMsg{
		Width:  msg.Width,
		Height: searchBarHeight,
	})
	m.command = searchbar.(command.Model)

	main, mainCmd := m.main.Update(tea.WindowSizeMsg{
		Width:  msg.Width,
		Height: msg.Height - searchBarHeight - 1,
	})
	m.main = main.(mainpanel.Model)

	return m, tea.Batch(searchCmd, mainCmd)
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return nil, tea.Quit
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

func (m Model) handleMoveFocus(msg message.MoveFocus) (tea.Model, tea.Cmd) {
	switch m.state.active {
	case cmdFocused:
		m.state.active++
		m, cmd := m.delegateToMainPanel(msg)
		return m, tea.Batch(cmd)
	case mainFocused:
		m.state.active--
		m, cmd := m.delegateToCommand(msg)
		return m, tea.Batch(cmd)
	default:
		return m, nil
	}
}

func (m Model) handleShowPopup() (tea.Model, tea.Cmd) {
	m.state.showModal = true
	return m, nil
}

func (m Model) handleHidePopup() (tea.Model, tea.Cmd) {
	m.state.showModal = false
	switch m.state.active {
	case cmdFocused:
		return m, nil
	default:
		return m, nil
	}
}

func (m *Model) delegateKeyPressHandler(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state.active {
	case cmdFocused:
		return m.delegateToCommand(msg)
	case mainFocused:
		return m.delegateToMainPanel(msg)
	default:
		return m, nil
	}
}

func (m Model) delegateToCommand(msg tea.Msg) (Model, tea.Cmd) {
	sb, cmd := m.command.Update(msg)
	m.command = sb.(command.Model)
	return m, cmd
}

func (m Model) delegateToMainPanel(msg tea.Msg) (Model, tea.Cmd) {
	main, cmd := m.main.Update(msg)
	m.main = main.(mainpanel.Model)
	return m, cmd
}

func (m Model) getHelpPopupContent() string {
	switch m.state.active {
	case cmdFocused:
		return m.command.Help()
	case mainFocused:
		return m.main.Help()
	default:
		return ""
	}
}

func (m Model) newPopupStyles() lipgloss.Style {
	return lipgloss.NewStyle().
		Width(m.width/2).
		Height(m.height/2).
		MaxWidth(120).
		MaxHeight(120).
		Border(lipgloss.NormalBorder(), true, true, true, true).
		BorderForeground(color.MainAccent)
}
