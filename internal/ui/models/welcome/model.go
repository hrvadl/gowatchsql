package welcome

import (
	"context"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/domain/engine"
	"github.com/hrvadl/gowatchsql/internal/platform/cfg"
	"github.com/hrvadl/gowatchsql/internal/ui/color"
	"github.com/hrvadl/gowatchsql/internal/ui/message"
	"github.com/hrvadl/gowatchsql/internal/ui/models/command"
	"github.com/hrvadl/gowatchsql/internal/ui/models/mainpanel"
	"github.com/hrvadl/gowatchsql/pkg/direction"
	"github.com/hrvadl/gowatchsql/pkg/overlay"
)

type ExplorerFactory interface {
	Create(ctx context.Context, name, dsn string) (engine.Explorer, error)
}

type ConnectionsRepo interface {
	GetConnections(context.Context) []cfg.Connection
}

func NewModel(log *slog.Logger, ef ExplorerFactory, connections ConnectionsRepo) Model {
	return Model{
		log:     log,
		command: command.NewModel(),
		main:    mainpanel.NewModel(ef, connections),
	}
}

type Model struct {
	command command.Model
	main    mainpanel.Model

	state state

	modalY int
	modalX int

	height int
	width  int

	log *slog.Logger
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.main.Init(), m.command.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleUpdateSize(msg)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case tea.QuitMsg:
		return m, tea.Quit
	case message.MoveFocus:
		return m.handleMoveFocus(msg)
	case message.Command, message.FetchedTableList, message.FetchedRows, message.FetchedColumns, message.SelectedTable:
		return m.delegateToMainPanel(msg)
	case message.SelectedContext:
		return m.delegateToAll(msg)
	case message.BlockCommandLine:
		return m.handleBlockCommandLine()
	case message.UnblockCommandLine:
		return m.handleUnblockCommandLine()
	default:
		return m.delegateToActive(msg)
	}
}

func (m Model) View() string {
	window := lipgloss.JoinVertical(lipgloss.Top, m.command.View(), m.main.View())
	popupStyles := m.newPopupStyles()

	if m.state.showModal {
		return overlay.Place(
			m.modalX,
			m.modalY,
			popupStyles.Render(m.getHelpPopupContent()),
			window,
			true,
		)
	}

	return window
}

func (m Model) handleBlockCommandLine() (Model, tea.Cmd) {
	m.state.blockModal = true
	return m, nil
}

func (m Model) handleUnblockCommandLine() (Model, tea.Cmd) {
	m.state.blockModal = false
	return m, nil
}

func (m Model) handleUpdateSize(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	const searchBarHeight = 4

	m.height = msg.Height
	m.width = msg.Width

	m.modalX = m.width / 4
	m.modalY = m.height / 4

	searchbar, searchCmd := m.command.Update(tea.WindowSizeMsg{
		Width:  msg.Width,
		Height: searchBarHeight,
	})
	m.command = searchbar

	main, mainCmd := m.main.Update(tea.WindowSizeMsg{
		Width:  msg.Width,
		Height: msg.Height - searchBarHeight - 1,
	})
	m.main = main

	return m, tea.Batch(searchCmd, mainCmd)
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		return m.handleEscape(msg)
	default:
		return m.handleKeyRunes(msg)
	}
}

func (m Model) handleKeyRunes(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "?":
		return m.handleShowPopup()
	case ":":
		return m.handleImmediateMoveFocus(msg)
	default:
		return m.delegateToActive(msg)
	}
}

func (m Model) handleImmediateMoveFocus(msg tea.KeyMsg) (Model, tea.Cmd) {
	if m.state.blockModal {
		return m.delegateToActive(msg)
	}

	switch m.state.active {
	case cmdFocused:
		return m.delegateToActive(msg)
	default:
		return m.handleUnfocusMainPanel()
	}
}

func (m Model) handleMoveFocus(msg message.MoveFocus) (Model, tea.Cmd) {
	switch m.state.active {
	case cmdFocused:
		return m.handleFocusMainPanel(msg)
	case mainFocused:
		return m.handleFocusCommand(msg)
	default:
		return m, nil
	}
}

func (m Model) handleFocusCommand(msg message.MoveFocus) (Model, tea.Cmd) {
	m.state.active = cmdFocused
	model, cmd := m.delegateToCommand(msg)
	return model, tea.Batch(cmd)
}

func (m Model) handleFocusMainPanel(msg message.MoveFocus) (Model, tea.Cmd) {
	m.state.active = mainFocused
	model, cmd := m.delegateToMainPanel(msg)
	return model, tea.Batch(cmd)
}

func (m Model) handleUnfocusMainPanel() (Model, tea.Cmd) {
	m.state.active = cmdFocused
	m, mainCmd := m.delegateToMainPanel(message.MoveFocus{Direction: direction.Away})
	m, cmdCmd := m.handleFocusCommand(message.MoveFocus{})
	return m, tea.Batch(mainCmd, cmdCmd)
}

func (m Model) handleShowPopup() (Model, tea.Cmd) {
	m.state.showModal = true
	return m, nil
}

func (m Model) handleEscape(msg tea.KeyMsg) (Model, tea.Cmd) {
	if !m.state.showModal {
		return m.delegateToActive(msg)
	}

	m.state.showModal = false
	return m, nil
}

func (m Model) delegateToAll(msg tea.Msg) (Model, tea.Cmd) {
	m, cmdCmd := m.delegateToCommand(msg)
	m, mainCmd := m.delegateToMainPanel(msg)
	return m, tea.Batch(cmdCmd, mainCmd)
}

func (m Model) delegateToActive(msg tea.Msg) (Model, tea.Cmd) {
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
	m.command = sb
	return m, cmd
}

func (m Model) delegateToMainPanel(msg tea.Msg) (Model, tea.Cmd) {
	main, cmd := m.main.Update(msg)
	m.main = main
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
