package details

type state struct {
	active  bool
	focused focused
}

type focused int

const (
	rowsFocused focused = iota
	columnsFocused
)
