package detailspanel

type status int

type state struct {
	status status
	active bool
}

const (
	Error status = iota
	pending
	ready
)
