package constraints

type status int

type state struct {
	status status
}

const (
	emtpy status = iota
	loading
	errored
	ready
)
