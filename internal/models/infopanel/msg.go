package infopanel

type (
	DSNReadyMsg struct {
		DSN string
	}

	ErrorMsg struct {
		Err error
	}
)
