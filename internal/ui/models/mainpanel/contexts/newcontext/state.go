package newcontext

type focused int

const (
	nameInputFocused focused = iota
	dsnInputFocused
)

const (
	inputsPage focused = iota
	confirmationPage
)

type state struct {
	active bool
}
