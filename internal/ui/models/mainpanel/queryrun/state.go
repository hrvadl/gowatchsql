package queryrun

type focused int

const (
	promptFocused focused = iota
	tableFocused
)

type state struct {
	active  bool
	focused focused
	err     error
}
