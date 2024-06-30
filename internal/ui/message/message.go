package message

import (
	"github.com/hrvadl/gowatchsql/internal/service/sysexplorer"
	"github.com/hrvadl/gowatchsql/internal/ui/command"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

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
		Tables []sysexplorer.Table
	}

	FetchedTableContent struct {
		Rows [][]string
		Cols []string
	}

	SelectedTable struct {
		Name string
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
