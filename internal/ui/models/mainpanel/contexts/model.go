package contexts

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hrvadl/gowatchsql/internal/ui/color"
	"github.com/hrvadl/gowatchsql/internal/ui/message"
	"github.com/hrvadl/gowatchsql/internal/ui/models/mainpanel/contexts/newcontext"
	"github.com/hrvadl/gowatchsql/internal/ui/styles"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

const margin = 1

func NewModel() Model {
	item := list.NewDefaultDelegate()
	item.Styles = styles.NewForItemDelegate()
	list := newList(item)

	return Model{
		list: list,
		state: state{
			active:     false,
			formActive: false,
		},
		newCtx: newcontext.NewModel(),
	}
}

type Model struct {
	width  int
	height int
	list   list.Model
	state  state
	newCtx newcontext.Model
}

func (m Model) Init() tea.Cmd {
	return nil
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
		return s.Render(m.list.View())
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

	newItems := append(m.list.Items(), newItemFromContext(msg))
	return m, tea.Batch(cmd, m.list.SetItems(newItems))
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

	if ctx, ok := m.list.SelectedItem().(ctxItem); ok {
		return m, func() tea.Msg { return message.SelectedContext{DSN: ctx.Description(), Name: ctx.Title()} }
	}
	return m, nil
}

func (m Model) handleKeyRunes(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "N":
		return m.handleToggleForm()
	default:
		return m.delegateToActive(msg)
	}
}

func (m Model) handleToggleForm() (Model, tea.Cmd) {
	m.state.formActive = !m.state.formActive
	if m.state.formActive {
		return m, func() tea.Msg { return message.BlockCommandLine{} }
	}

	return m, func() tea.Msg { return message.UnblockCommandLine{} }
}

func (m Model) handleError(msg message.Error) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) handleWindowSize(w, h int) (Model, tea.Cmd) {
	m.width = w
	m.height = h
	m.list.SetSize(w, h)
	m.list.Styles.TitleBar = m.list.Styles.TitleBar.Width(w)
	msg := tea.WindowSizeMsg{Height: h, Width: w}
	m, cmd := m.delegateToAll(msg)
	return m, cmd
}

func (m Model) delegateToAll(msg tea.Msg) (Model, tea.Cmd) {
	m, listCmd := m.delegateToList(msg)
	m, ctxCmd := m.delegateToNewContextModel(msg)
	return m, tea.Batch(listCmd, ctxCmd)
}

func (m Model) delegateToActive(msg tea.Msg) (Model, tea.Cmd) {
	switch m.state.formActive {
	case true:
		return m.delegateToNewContextModel(msg)
	default:
		return m.delegateToList(msg)
	}
}

func (m Model) delegateToNewContextModel(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.newCtx, cmd = m.newCtx.Update(msg)
	return m, cmd
}

func (m Model) delegateToList(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
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

func newList(item list.ItemDelegate) list.Model {
	const defaultTitle = "Contexts"

	l := list.New([]list.Item{}, item, 0, 0)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.InfiniteScrolling = true
	l.Styles = styles.NewForList()
	l.Title = defaultTitle
	l.KeyMap.Quit = key.NewBinding(key.WithDisabled())

	return l
}
