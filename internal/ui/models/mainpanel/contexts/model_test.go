package contexts

import (
	"bytes"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"go.uber.org/mock/gomock"

	"github.com/hrvadl/gowatchsql/internal/platform/cfg"
	"github.com/hrvadl/gowatchsql/internal/ui/message"
	"github.com/hrvadl/gowatchsql/internal/ui/models/mainpanel/contexts/mocks"
	"github.com/hrvadl/gowatchsql/pkg/xtest"
)

func TestNewConnectionAppears(t *testing.T) {
	xtest.SkipIntegrationIfRequired(t)

	repo := mocks.NewMockConnectionsRepo(gomock.NewController(t))
	repo.EXPECT().GetConnections(gomock.Any()).Return([]cfg.Connection{})

	m := NewModel(repo)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(300, 100))

	tm.Send(message.NewContext{DSN: "DSN", Name: "new naaame", OK: true})
	teatest.WaitFor(
		t, tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("new naaame"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*3),
	)
}

func TestConnectionSelected(t *testing.T) {
	xtest.SkipIntegrationIfRequired(t)

	repo := mocks.NewMockConnectionsRepo(gomock.NewController(t))
	repo.EXPECT().GetConnections(gomock.Any()).Return([]cfg.Connection{})

	m := NewModel(repo)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(300, 100))

	tm.Send(message.NewContext{DSN: "DSN", Name: "pg", OK: true})
	tm.Send(message.NewContext{DSN: "DSN", Name: "mysql", OK: true})
	tm.Send(message.NewContext{DSN: "DSN", Name: "sqlite", OK: true})
	tm.Send(message.SelectedContext{DSN: "DSN", Name: "sqlite"})

	teatest.WaitFor(
		t, tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("sqlite"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*3),
	)
}

func TestConnectionDeleted(t *testing.T) {
	xtest.SkipIntegrationIfRequired(t)

	repo := mocks.NewMockConnectionsRepo(gomock.NewController(t))
	repo.EXPECT().GetConnections(gomock.Any()).Return([]cfg.Connection{})
	repo.EXPECT().DeleteConnection(gomock.Any(), "DSN").MaxTimes(1).Return(nil)

	m := NewModel(repo)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(300, 100))

	tm.Send(message.NewContext{DSN: "DSN", Name: "pg", OK: true})
	tm.Send(message.NewContext{DSN: "DSN", Name: "mysql", OK: true})
	tm.Send(message.NewContext{DSN: "DSN", Name: "sqlite", OK: true})
	tm.Send(message.SelectedContext{DSN: "DSN", Name: "sqlite"})

	teatest.WaitFor(
		t, tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("sqlite"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*3),
	)

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	teatest.WaitFor(
		t, tm.Output(),
		func(bts []byte) bool {
			return !bytes.Contains(bts, []byte("sqlite"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*3),
	)
}
