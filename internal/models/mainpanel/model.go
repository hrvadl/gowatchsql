package mainpanel

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/hrvadl/gowatchsql/internal/command"
	"github.com/hrvadl/gowatchsql/internal/message"
	"github.com/hrvadl/gowatchsql/internal/models/mainpanel/contexts"
	"github.com/hrvadl/gowatchsql/internal/models/mainpanel/contexts/newcontext"
	"github.com/hrvadl/gowatchsql/internal/models/mainpanel/objects"
	"github.com/hrvadl/gowatchsql/internal/models/mainpanel/queryrun"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

func NewModel() Model {
	return Model{
		objects:    objects.NewModel(),
		contexts:   contexts.NewModel(),
		newContext: newcontext.NewModel(),
		queryrun:   queryrun.NewModel(),
	}
}

type Model struct {
	objects    objects.Model
	contexts   contexts.Model
	newContext newcontext.Model
	queryrun   queryrun.Model
	state      state
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.delegateToAllModels(msg)
	case tea.KeyMsg:
		return m.delegateToActiveModel(msg)
	case message.MoveFocus:
		slog.Info("Moving focus to main panel")
		return m.delegateToActiveModel(msg)
	case message.SelectedDB, message.TableChosen:
		return m.delegateToObjectsModel(msg)
	case message.Command:
		return m.handleCommand(msg)
	case message.Error:
		return m.handleError(msg)
	default:
		return m.delegateToActiveModel(msg)
	}
}

func (m Model) View() string {
	switch m.state.active {
	case objectsActive:
		return m.objects.View()
	case contextsActive:
		return m.contexts.View()
	case newContextActive:
		return m.newContext.View()
	case queryRunActive:
		return m.queryrun.View()
	default:
		return "Idk that view"
	}
}

func (m Model) Help() string {
	switch m.state.active {
	case objectsActive:
		return m.objects.Help()
	case contextsActive:
		return m.contexts.Help()
	default:
		return "TODO: change me"
	}
}

func (m Model) handleCommand(msg message.Command) (Model, tea.Cmd) {
	switch msg.Text {
	case command.Tables:
		m.state.active = objectsActive
	case command.Context:
		m.state.active = contextsActive
	case command.NewContext:
		m.state.active = newContextActive
	}
	return m, func() tea.Msg { return message.MoveFocus{Direction: direction.Forward} }
}

func (m Model) handleError(msg message.Error) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) delegateToAllModels(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	m, objCmd := m.delegateToObjectsModel(msg)
	m, contextsCmd := m.delegateToContextsModel(msg)
	m, newContextsCmd := m.delegateToNewContextsModel(msg)
	m, queryRunCmd := m.delegateToQueryRunModel(msg)
	return m, tea.Batch(objCmd, contextsCmd, newContextsCmd, queryRunCmd)
}

func (m Model) delegateToActiveModel(msg tea.Msg) (Model, tea.Cmd) {
	switch m.state.active {
	case objectsActive:
		return m.delegateToObjectsModel(msg)
	case contextsActive:
		return m.delegateToContextsModel(msg)
	case newContextActive:
		return m.delegateToNewContextsModel(msg)
	case queryRunActive:
		return m.delegateToQueryRunModel(msg)
	default:
		return m, nil
	}
}

func (m Model) delegateToObjectsModel(msg tea.Msg) (Model, tea.Cmd) {
	model, cmd := m.objects.Update(msg)
	m.objects = model
	return m, cmd
}

func (m Model) delegateToContextsModel(msg tea.Msg) (Model, tea.Cmd) {
	model, cmd := m.contexts.Update(msg)
	m.contexts = model
	return m, cmd
}

func (m Model) delegateToNewContextsModel(msg tea.Msg) (Model, tea.Cmd) {
	model, cmd := m.newContext.Update(msg)
	m.newContext = model
	return m, cmd
}

func (m Model) delegateToQueryRunModel(msg tea.Msg) (Model, tea.Cmd) {
	model, cmd := m.queryrun.Update(msg)
	m.queryrun = model
	return m, cmd
}
