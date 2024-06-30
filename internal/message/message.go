package message

import (
	"github.com/hrvadl/gowatchsql/internal/command"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

type (
	BlockCommandLine struct{}

	UnblockCommandLine struct{}

	Command struct {
		Text command.Command
	}

	SelectedDB struct {
		DSN string
	}

	Error struct {
		Err error
	}

	TableChosen struct {
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
