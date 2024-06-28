package message

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
)
