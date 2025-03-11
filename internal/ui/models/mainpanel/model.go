package mainpanel

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/hrvadl/gowatchsql/internal/domain/engine"
	"github.com/hrvadl/gowatchsql/internal/platform/cfg"
	"github.com/hrvadl/gowatchsql/internal/ui/command"
	"github.com/hrvadl/gowatchsql/internal/ui/message"
	"github.com/hrvadl/gowatchsql/internal/ui/models/mainpanel/contexts"
	"github.com/hrvadl/gowatchsql/internal/ui/models/mainpanel/objects"
	"github.com/hrvadl/gowatchsql/internal/ui/models/mainpanel/queryrun"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

//go:generate mockgen -destination=mocks/mock_factory.go -package=mocks . ExplorerFactory
type ExplorerFactory interface {
	Create(ctx context.Context, name, dsn string) (engine.Explorer, error)
}

//go:generate mockgen -destination=mocks/mock_repo.go -package=mocks . ConnectionsRepo
type ConnectionsRepo interface {
	GetConnections(context.Context) []cfg.Connection
	DeleteConnection(ctx context.Context, dsn string) error
}

func NewModel(explorerFactory ExplorerFactory, connections ConnectionsRepo) Model {
	return Model{
		objects:  objects.NewModel(explorerFactory),
		contexts: contexts.NewModel(connections),
		queryrun: queryrun.NewModel(explorerFactory),
	}
}

type Model struct {
	objects  objects.Model
	contexts *contexts.Model
	queryrun queryrun.Model
	state    state
}

func (m Model) Init() tea.Cmd {
	b := tea.Batch(m.objects.Init(), m.contexts.Init(), m.queryrun.Init())
	return b
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	mm, cmd := m.update(msg)
	return mm, cmd
}

func (m Model) update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.delegateToAllModels(msg)
	case tea.KeyMsg:
		return m.delegateToActiveModel(msg)
	case message.MoveFocus:
		return m.delegateToActiveModel(msg)
	case message.SelectedContext,
		message.SelectedTable,
		message.FetchedRows,
		message.FetchedColumns,
		message.FetchedTableList,
		message.FetchedIndexes,
		message.FetchedConstraints:
		return m.delegateToAllModels(msg)
	case message.Command:
		return m.handleCommand(msg)
	case message.Error:
		return m.delegateToActiveModel(msg)
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
	case command.Query:
		m.state.active = queryRunActive
	case command.Context:
		m.state.active = contextsActive
	case command.Exit:
		return m, tea.Quit
	}
	return m, message.With(message.MoveFocus{Direction: direction.Forward})
}

func (m Model) delegateToAllModels(msg tea.Msg) (Model, tea.Cmd) {
	m, objCmd := m.delegateToObjectsModel(msg)
	m, contextsCmd := m.delegateToContextsModel(msg)
	m, queryRunCmd := m.delegateToQueryRunModel(msg)
	return m, tea.Batch(objCmd, contextsCmd, queryRunCmd)
}

func (m Model) delegateToActiveModel(msg tea.Msg) (Model, tea.Cmd) {
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

func (m Model) delegateToObjectsModel(msg tea.Msg) (Model, tea.Cmd) {
	model, cmd := m.objects.Update(msg)
	m.objects = model
	return m, cmd
}

func (m Model) delegateToContextsModel(msg tea.Msg) (Model, tea.Cmd) {
	model, cmd := m.contexts.Update(msg)
	m.contexts = model.(*contexts.Model)
	return m, cmd
}

func (m Model) delegateToQueryRunModel(msg tea.Msg) (Model, tea.Cmd) {
	model, cmd := m.queryrun.Update(msg)
	m.queryrun = model
	return m, cmd
}
