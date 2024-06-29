package mainpanel

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hrvadl/gowatchsql/internal/message"
	"github.com/hrvadl/gowatchsql/internal/models/mainpanel/contexts"
	"github.com/hrvadl/gowatchsql/internal/models/mainpanel/objects"
	"github.com/hrvadl/gowatchsql/internal/models/mainpanel/queryrun"
)

func NewModel() Model {
	return Model{
		objects:  objects.NewModel(),
		contexts: contexts.NewModel(),
		queryrun: queryrun.NewModel(),
	}
}

type active int

const (
	objectsActive active = iota
	contextsActive
	queryRunActive
)

type state struct {
	active active
}

type Model struct {
	objects  objects.Model
	contexts contexts.Model
	queryrun queryrun.Model
	state    state
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.delegateToAllModels(msg)
	case tea.KeyMsg:
		return m.delegateToActiveModel(msg)
	case message.MoveFocus:
		return m.delegateToActiveModel(msg)
	case message.DSNReady, message.TableChosen:
		return m.delegateToObjectsModel(msg)
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
	case queryRunActive:
		return m.queryrun.View()
	default:
		return "Idk that view"
	}
}

func (m Model) Help() string {
	return "help"
}

func (m Model) handleError(msg message.Error) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m Model) delegateToAllModels(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	obj, objCmd := m.objects.Update(msg)
	ctx, contextsCmd := m.contexts.Update(msg)
	queryRun, queryRunCmd := m.queryrun.Update(msg)

	m.objects = obj.(objects.Model)
	m.contexts = ctx.(contexts.Model)
	m.queryrun = queryRun.(queryrun.Model)

	return m, tea.Batch(objCmd, contextsCmd, queryRunCmd)
}

func (m Model) delegateToActiveModel(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state.active {
	case objectsActive:
		return m.delegateToObjectsModel(msg)
	case contextsActive:
		return m.delegateToContextsModel(msg)
	case queryRunActive:
		return m.delegateToQueryRunModel(msg)
	default:
		return m, nil
	}
}

func (m Model) delegateToObjectsModel(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := m.objects.Update(msg)
	m.objects = model.(objects.Model)
	return m, cmd
}

func (m Model) delegateToContextsModel(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := m.contexts.Update(msg)
	m.contexts = model.(contexts.Model)
	return m, cmd
}

func (m Model) delegateToQueryRunModel(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := m.queryrun.Update(msg)
	m.queryrun = model.(queryrun.Model)
	return m, cmd
}
