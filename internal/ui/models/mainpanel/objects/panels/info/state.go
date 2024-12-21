package info

type state struct {
	status status
	active bool
	err    error
}

type status int

const (
	empty status = iota
	loading
	errored
	ready
)
