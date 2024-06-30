package message

import (
	"github.com/hrvadl/gowatchsql/internal/command"
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

type (
	Command struct {
		Text command.Command
	}

	DSNReady struct {
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
)
