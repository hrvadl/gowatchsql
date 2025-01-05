package queryrun

type focused int

const (
	promptFocused focused = iota
	tableFocues
)

type state struct {
	active  bool
	focused focused
}
