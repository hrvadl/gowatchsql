package objects

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/platform/db"
	"github.com/hrvadl/gowatchsql/internal/service/sysexplorer"
	"github.com/hrvadl/gowatchsql/internal/service/tableexplorer"
	"github.com/hrvadl/gowatchsql/internal/ui/message"
	"github.com/hrvadl/gowatchsql/internal/ui/models/mainpanel/objects/detailspanel"
	"github.com/hrvadl/gowatchsql/internal/ui/models/mainpanel/objects/infopanel"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

// @TODO: gracefully close connections
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

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSize(msg)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case message.SelectedTable:
		return m.handleTableChosen(msg)
	case message.FetchedTableContent:
		return m.delegateToDetailsModel(msg)
	case message.SelectedContext, message.FetchedTableList:
		return m.delegateToInfoModel(msg)
	case message.Error:
		return m.handleError(msg)
	case message.MoveFocus:
		return m.handleMoveFocus(msg)
	default:
		return m.delegateToActiveModel(msg)
	}
}

func (m Model) View() string {
	return m.newStyles().Render(
		lipgloss.JoinHorizontal(lipgloss.Left, m.info.View(), m.details.View()),
	)
}

func (m Model) Help() string {
	switch m.state.focused {
	case infoFocused:
		return m.info.Help()
	case detailsFocused:
		return m.details.Help()
	default:
		return "error: unknown view"
	}
}

func (m Model) handleMoveFocus(msg message.MoveFocus) (Model, tea.Cmd) {
	switch msg.Direction {
	case direction.Forward:
		return m.forwardStrategy(msg)
	case direction.Backwards:
		return m.backwardsStrategy(msg)
	case direction.Away:
		m.state.focused = unfocused
		return m.delegateToAllModels(msg)
	default:
		return m, nil
	}
}

func (m Model) forwardStrategy(msg message.MoveFocus) (Model, tea.Cmd) {
	switch m.state.focused {
	case unfocused:
		m.state.focused = infoFocused
		return m.delegateToInfoModel(msg)
	case infoFocused:
		m.state.focused = detailsFocused
		return m.delegateToAllModels(msg)
	default:
		m.state.focused = infoFocused
		return m.delegateToAllModels(msg)
	}
}

func (m Model) backwardsStrategy(msg message.MoveFocus) (Model, tea.Cmd) {
	switch m.state.focused {
	case unfocused:
		m.state.focused = detailsFocused
		return m.delegateToDetailsModel(msg)
	case infoFocused:
		m.state.focused = detailsFocused
		return m.delegateToAllModels(msg)
	default:
		m.state.focused = infoFocused
		return m.delegateToAllModels(msg)
	}
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyTab:
		return m.handleMoveFocus(message.MoveFocus{Direction: direction.Forward})
	case tea.KeyShiftTab:
		return m.handleMoveFocus(message.MoveFocus{Direction: direction.Backwards})
	default:
		return m.delegateToActiveModel(msg)
	}
}

func (m Model) handleWindowSize(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	const infoPanelWidth = 20
	m.width = msg.Width
	m.height = msg.Height

	m, infoCmd := m.delegateToInfoModel(tea.WindowSizeMsg{
		Width:  infoPanelWidth,
		Height: msg.Height,
	})

	m, detailsCmd := m.delegateToDetailsModel(tea.WindowSizeMsg{
		Width:  msg.Width - infoPanelWidth,
		Height: msg.Height,
	})

	return m, tea.Batch(infoCmd, detailsCmd)
}

func (m Model) handleTableChosen(msg message.SelectedTable) (Model, tea.Cmd) {
	m.details.SetTableExplorer(tableexplorer.New(db.Get()))
	return m.delegateToDetailsModel(msg)
}

func (m Model) handleError(msg message.Error) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) delegateToAllModels(msg tea.Msg) (Model, tea.Cmd) {
	m, infoCmd := m.delegateToInfoModel(msg)
	m, detailsCmd := m.delegateToDetailsModel(msg)
	return m, tea.Batch(infoCmd, detailsCmd)
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

func (m Model) newStyles() lipgloss.Style {
	return lipgloss.NewStyle().Width(m.width).Height(m.height)
}
