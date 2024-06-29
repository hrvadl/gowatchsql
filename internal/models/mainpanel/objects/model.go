package objects

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/message"
	"github.com/hrvadl/gowatchsql/internal/models/mainpanel/objects/detailspanel"
	"github.com/hrvadl/gowatchsql/internal/models/mainpanel/objects/infopanel"
	"github.com/hrvadl/gowatchsql/internal/platform/db"
	"github.com/hrvadl/gowatchsql/internal/service/sysexplorer"
	"github.com/hrvadl/gowatchsql/internal/service/tableexplorer"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

func NewModel() Model {
	return Model{
		info:    infopanel.NewModel(sysexplorer.New),
		details: detailspanel.NewModel(),
	}
}

type focus int

type state struct {
	focused focus
}

const (
	unfocused focus = iota
	infoFocused
	detailsFocused
)

type Model struct {
	state state

	width  int
	height int

	info    infopanel.Model
	details detailspanel.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSize(msg)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case message.DSNReady:
		return m.delegateToInfoModel(msg)
	case message.TableChosen:
		return m.handleTableChosen(msg)
	case message.Error:
		return m.handleError(msg)
	case message.MoveFocus:
		return m.handleMoveFocus(msg)
	default:
		return m, nil
	}
}

func (m Model) View() string {
	return lipgloss.NewStyle().Width(m.width).Height(m.height).Render(
		lipgloss.JoinHorizontal(lipgloss.Left, m.info.View(), m.details.View()),
	)
}

func (m Model) handleMoveFocus(msg message.MoveFocus) (tea.Model, tea.Cmd) {
	switch msg.Direction {
	case direction.Forward:
		return m.forwardStrategy(msg)
	case direction.Backwards:
		return m.backwardsStrategy(msg)
	default:
		return m, nil
	}
}

func (m Model) forwardStrategy(msg message.MoveFocus) (tea.Model, tea.Cmd) {
	switch m.state.focused {
	case unfocused:
		m.state.focused++
		return m.delegateToInfoModel(msg)
	case infoFocused:
		m.state.focused++
		return m.delegateToAllModels(msg)
	default:
		m.state.focused = unfocused
		unfocusPane := func() tea.Msg { return message.MoveFocus{Direction: direction.Forward} }
		m, cmds := m.delegateToDetailsModel(msg)
		return m, tea.Batch(cmds, unfocusPane)
	}
}

func (m Model) backwardsStrategy(msg message.MoveFocus) (tea.Model, tea.Cmd) {
	switch m.state.focused {
	case unfocused:
		m.state.focused = detailsFocused
		return m.delegateToDetailsModel(msg)
	case infoFocused:
		m.state.focused = unfocused
		unfocusPane := func() tea.Msg { return message.MoveFocus{Direction: direction.Backwards} }
		m, cmds := m.delegateToInfoModel(msg)
		return m, tea.Batch(cmds, unfocusPane)
	default:
		m.state.focused--
		return m.delegateToAllModels(msg)
	}
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyTab:
		return m.handleMoveFocus(message.MoveFocus{Direction: direction.Forward})
	case tea.KeyShiftTab:
		return m.handleMoveFocus(message.MoveFocus{Direction: direction.Backwards})
	default:
		return m.delegateToActiveModel(msg)
	}
}

func (m Model) handleWindowSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	const infoPanelWidth = 20
	m.width = msg.Width
	m.height = msg.Height

	info, infoCmd := m.delegateToInfoModel(tea.WindowSizeMsg{
		Width:  infoPanelWidth,
		Height: msg.Height,
	})

	details, detailsCmd := info.delegateToDetailsModel(tea.WindowSizeMsg{
		Width:  msg.Width - infoPanelWidth,
		Height: msg.Height,
	})

	return details, tea.Batch(infoCmd, detailsCmd)
}

func (m Model) handleTableChosen(msg message.TableChosen) (tea.Model, tea.Cmd) {
	m.details.SetTableExplorer(tableexplorer.New(db.Get()))
	return m.delegateToDetailsModel(msg)
}

func (m Model) handleError(msg message.Error) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) delegateToAllModels(msg tea.Msg) (tea.Model, tea.Cmd) {
	info, infoCmd := m.delegateToInfoModel(msg)
	details, detailsCmd := info.delegateToDetailsModel(msg)
	return details, tea.Batch(infoCmd, detailsCmd)
}

func (m Model) delegateToActiveModel(msg tea.Msg) (Model, tea.Cmd) {
	switch m.state.focused {
	case detailsFocused:
		return m.delegateToDetailsModel(msg)
	case infoFocused:
		return m.delegateToInfoModel(msg)
	default:
		return m, nil
	}
}

func (m Model) delegateToDetailsModel(msg tea.Msg) (Model, tea.Cmd) {
	model, cmd := m.details.Update(msg)
	m.details = model.(detailspanel.Model)
	return m, cmd
}

func (m Model) delegateToInfoModel(msg tea.Msg) (Model, tea.Cmd) {
	model, cmd := m.info.Update(msg)
	m.info = model.(infopanel.Model)
	return m, cmd
}