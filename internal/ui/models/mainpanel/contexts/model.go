package contexts

import (
	"context"
	"slices"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/platform/cfg"
	"github.com/hrvadl/gowatchsql/internal/ui/color"
	"github.com/hrvadl/gowatchsql/internal/ui/message"
	"github.com/hrvadl/gowatchsql/internal/ui/models/mainpanel/contexts/createmodal"
	"github.com/hrvadl/gowatchsql/internal/ui/styles"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

const margin = 1

type ConnectionsReppo interface {
	GetConnections(context.Context) []cfg.Connection
}

func NewModel(connections ConnectionsReppo) *Model {
	item := list.NewDefaultDelegate()
	item.Styles = styles.NewForItemDelegate()
	list := newList(item, []list.Item{})

	return &Model{
		List: list,
		state: state{
			active:     false,
			formActive: false,
		},
		newCtx:      createmodal.NewModel(),
		connections: connections,
	}
}

type Model struct {
	width       int
	height      int
	List        list.Model
	state       state
	newCtx      createmodal.Model
	connections ConnectionsReppo
}

func (m *Model) Init() tea.Cmd {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	connections := m.connections.GetConnections(ctx)
	if len(connections) == 0 {
		return nil
	}

	listItems := make([]list.Item, 0, len(connections))
	for _, connection := range connections {
		msg := message.NewContext{OK: true, DSN: connection.DSN, Name: connection.Name}
		listItems = append(listItems, newItemFromContext(msg))
	}

	item := list.NewDefaultDelegate()
	item.Styles = styles.NewForItemDelegate()
	m.List = newList(item, listItems)

	return message.With(message.SelectedContext{Name: connections[0].Name, DSN: connections[0].DSN})
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSize(msg.Width-margin*2, msg.Height-margin*3)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case message.Error:
		return m.handleError(msg)
	case message.MoveFocus:
		return m.handleMoveFocus(msg)
	case message.NewContext:
		return m.handleNewContext(msg)
	default:
		return m.delegateToActive(msg)
	}
}

func (m Model) View() string {
	s := m.newContainerStyles()
	if !m.state.formActive {
		return s.Render(m.List.View())
	}

	return s.Padding(1, 2).Render(m.newCtx.View())
}

func (m Model) Help() string {
	return "help from contexts"
}

func (m Model) handleNewContext(msg message.NewContext) (Model, tea.Cmd) {
	m, cmd := m.handleToggleForm()
	if !msg.OK {
		return m, cmd
	}

	newItems := append(m.List.Items(), newItemFromContext(msg))
	return m, tea.Batch(cmd, m.List.SetItems(newItems))
}

func (m Model) handleMoveFocus(msg message.MoveFocus) (Model, tea.Cmd) {
	switch msg.Direction {
	case direction.Away:
		m.state.active = false
	default:
		m.state.active = true
	}
	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		return m.handleKeyEnter(msg)
	case tea.KeyRunes:
		return m.handleKeyRunes(msg)
	default:
		return m.delegateToActive(msg)
	}
}

func (m Model) handleKeyEnter(msg tea.KeyMsg) (Model, tea.Cmd) {
	if m.state.formActive {
		return m.delegateToActive(msg)
	}

	if ctx, ok := m.List.SelectedItem().(ctxItem); ok {
		return m, message.With(message.SelectedContext{DSN: ctx.Description(), Name: ctx.Title()})
	}
	return m, nil
}

func (m Model) handleKeyRunes(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case toggleCreateForm:
		return m.handleToggleForm()
	default:
		return m.delegateKeypress(msg)
	}
}

func (m Model) handleDeleteContext() (Model, tea.Cmd) {
	items := m.List.Items()
	idx := m.List.Index()
	return m, m.List.SetItems(slices.Delete(items, idx, idx+1))
}

func (m Model) handleToggleForm() (Model, tea.Cmd) {
	m.state.formActive = !m.state.formActive
	if m.state.formActive {
		return m, message.With(message.BlockCommandLine{})
	}

	return m, message.With(message.UnblockCommandLine{})
}

func (m Model) handleError(msg message.Error) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) handleWindowSize(w, h int) (Model, tea.Cmd) {
	m.width = w
	m.height = h
	m.List.SetSize(w, h)
	m.List.Styles.TitleBar = m.List.Styles.TitleBar.Width(w)
	msg := tea.WindowSizeMsg{Height: h, Width: w}
	m, cmd := m.delegateToAll(msg)
	return m, cmd
}

func (m Model) delegateToAll(msg tea.Msg) (Model, tea.Cmd) {
	m, listCmd := m.delegateToList(msg)
	m, ctxCmd := m.delegateToNewContextModel(msg)
	return m, tea.Batch(listCmd, ctxCmd)
}

func (m Model) delegateKeypress(msg tea.KeyMsg) (Model, tea.Cmd) {
	if m.state.formActive {
		return m.delegateToNewContextModel(msg)
	}

	switch msg.String() {
	case deleteContext:
		return m.handleDeleteContext()
	default:
		return m.delegateToActive(msg)
	}
}

func (m Model) delegateToActive(msg tea.Msg) (Model, tea.Cmd) {
	if m.state.formActive {
		return m.delegateToNewContextModel(msg)
	}

	return m.delegateToList(msg)
}

func (m Model) delegateToNewContextModel(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.newCtx, cmd = m.newCtx.Update(msg)
	return m, cmd
}

func (m Model) delegateToList(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) newContainerStyles() lipgloss.Style {
	base := lipgloss.
		NewStyle().
		Height(m.height).
		Width(m.width).
		Border(lipgloss.ThickBorder())

	if m.state.active {
		return base.BorderForeground(color.MainAccent)
	}

	return base.BorderForeground(color.Border)
}

func newList(item list.ItemDelegate, rows []list.Item) list.Model {
	const defaultTitle = "Contexts"

	l := list.New(rows, item, 0, 0)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.InfiniteScrolling = true
	l.Styles = styles.NewForList()
	l.Title = defaultTitle
	l.KeyMap.Quit = key.NewBinding(key.WithDisabled())

	return l
}
