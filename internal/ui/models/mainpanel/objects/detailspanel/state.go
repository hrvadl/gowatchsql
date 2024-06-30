package detailspanel

type status int

type state struct {
	active   bool
	status   status
	showHelp bool
}

const (
	errored status = iota
	ready
)
