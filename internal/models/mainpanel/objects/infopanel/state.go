package infopanel

type state struct {
	status status
	active bool
	err    error
}

type status int

const (
	pending status = iota
	Error
	ready
)
