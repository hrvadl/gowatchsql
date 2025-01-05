package message

import (
	"fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/hrvadl/gowatchsql/internal/domain/engine"
	"github.com/hrvadl/gowatchsql/internal/ui/command"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

func With(msg tea.Msg) func() tea.Msg {
	return func() tea.Msg {
		fmted := fmt.Sprintf("Sending message %T", msg)
		slog.Debug(fmted, slog.Any("message", msg))
		return msg
	}
}

type (
	CleanCommandLine struct{}

	BlockCommandLine struct{}

	UnblockCommandLine struct{}

	Command struct {
		Text command.Command
	}

	SelectedContext struct {
		Name string
		DSN  string
	}

	Error struct {
		Err error
	}

	FetchedTableList struct {
		Tables []engine.Table
	}

	FetchedRows struct {
		Rows [][]string
		Cols []string
	}

	FetchedColumns struct {
		Rows [][]string
		Cols []string
	}

	FetchedIndexes struct {
		Rows [][]string
		Cols []string
	}

	FetchedConstraints struct {
		Rows [][]string
		Cols []string
	}

	SelectedTable struct {
		Name string
	}

	ExecuteCommand struct {
		Cmd string
	}

	MoveFocus struct {
		Direction direction.Direction
	}

	NewContext struct {
		DSN  string
		Name string
		OK   bool
	}
)
