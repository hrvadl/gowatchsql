package mainpanel

import (
	"bytes"
	"testing"
	"time"

	"github.com/charmbracelet/x/exp/teatest"
	"go.uber.org/mock/gomock"

	"github.com/hrvadl/gowatchsql/internal/domain/engine"
	"github.com/hrvadl/gowatchsql/internal/platform/cfg"
	"github.com/hrvadl/gowatchsql/internal/ui/message"
	"github.com/hrvadl/gowatchsql/internal/ui/models/mainpanel/mocks"
	"github.com/hrvadl/gowatchsql/pkg/xtest"
)

func TestInterfaceShowsAfterConnection(t *testing.T) {
	xtest.SkipIntegrationIfRequired(t)

	exp := engine.NewMockExplorer(gomock.NewController(t))
	exp.EXPECT().GetTables(gomock.Any()).Return([]engine.Table{
		{
			Name:   "table1",
			Schema: "schema",
		},
		{
			Name:   "table2",
			Schema: "schema",
		},
	}, nil)

	exp.EXPECT().GetRows(gomock.Any(), "table1").MinTimes(1).Return([]engine.Row{
		{
			"row1",
		},
		{
			"row2",
		},
	},
		[]engine.Column{
			"col1",
		},
		nil)

	exp.EXPECT().GetIndexes(gomock.Any(), "table1").MinTimes(1).Return([]engine.Row{
		{
			"index1",
		},
		{
			"index2",
		},
	},
		[]engine.Column{
			"col1",
		},
		nil)

	exp.EXPECT().GetColumns(gomock.Any(), "table1").MinTimes(1).Return([]engine.Row{
		{
			"c1",
		},
		{
			"c2",
		},
	},
		[]engine.Column{
			"col1",
		},
		nil)

	exp.EXPECT().GetConstraints(gomock.Any(), "table1").MinTimes(1).Return([]engine.Row{
		{
			"ct1",
		},
		{
			"ct2",
		},
	},
		[]engine.Column{
			"col1",
		},
		nil)

	ef := mocks.NewMockExplorerFactory(gomock.NewController(t))
	ef.EXPECT().
		Create(gomock.Any(), "new naaame", "DSN").MinTimes(1).
		Return(exp, nil)

	repo := mocks.NewMockConnectionsRepo(gomock.NewController(t))
	repo.EXPECT().GetConnections(gomock.Any()).Return([]cfg.Connection{})

	m := NewModel(ef, repo)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(300, 100))

	tm.Send(message.SelectedContext{DSN: "DSN", Name: "new naaame"})

	teatest.WaitFor(
		t, tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("Rows")) &&
				bytes.Contains(bts, []byte("Indexes")) &&
				bytes.Contains(bts, []byte("Tables")) &&
				bytes.Contains(bts, []byte("Constraints"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*3),
	)
}

func TestAllTablesAreShownAfterConnection(t *testing.T) {
	xtest.SkipIntegrationIfRequired(t)

	exp := engine.NewMockExplorer(gomock.NewController(t))
	exp.EXPECT().GetTables(gomock.Any()).Return([]engine.Table{
		{
			Name:   "table1",
			Schema: "schema",
		},
		{
			Name:   "table2",
			Schema: "schema",
		},
	}, nil)

	exp.EXPECT().GetRows(gomock.Any(), "table1").MinTimes(1).Return([]engine.Row{
		{
			"row1",
		},
		{
			"row2",
		},
	},
		[]engine.Column{
			"col1",
		},
		nil)

	exp.EXPECT().GetIndexes(gomock.Any(), "table1").MinTimes(1).Return([]engine.Row{
		{
			"index1",
		},
		{
			"index2",
		},
	},
		[]engine.Column{
			"col1",
		},
		nil)

	exp.EXPECT().GetColumns(gomock.Any(), "table1").MinTimes(1).Return([]engine.Row{
		{
			"c1",
		},
		{
			"c2",
		},
	},
		[]engine.Column{
			"col1",
		},
		nil)

	exp.EXPECT().GetConstraints(gomock.Any(), "table1").MinTimes(1).Return([]engine.Row{
		{
			"ct1",
		},
		{
			"ct2",
		},
	},
		[]engine.Column{
			"col1",
		},
		nil)

	ef := mocks.NewMockExplorerFactory(gomock.NewController(t))
	ef.EXPECT().
		Create(gomock.Any(), "new naaame", "DSN").MinTimes(1).
		Return(exp, nil)

	repo := mocks.NewMockConnectionsRepo(gomock.NewController(t))
	repo.EXPECT().GetConnections(gomock.Any()).Return([]cfg.Connection{})

	m := NewModel(ef, repo)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(300, 100))

	tm.Send(message.SelectedContext{DSN: "DSN", Name: "new naaame"})

	teatest.WaitFor(
		t, tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("table1")) &&
				bytes.Contains(bts, []byte("table2"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*3),
	)
}
