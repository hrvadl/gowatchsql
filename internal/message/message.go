package message

import (
	"github.com/hrvadl/gowatchsql/pkg/direction"
)

type (
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
