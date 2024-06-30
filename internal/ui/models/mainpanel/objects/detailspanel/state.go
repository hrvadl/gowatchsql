package detailspanel

type status int

type state struct {
	active bool
	status status
}

const (
	emtpy status = iota
	loading
	errored
	ready
)
