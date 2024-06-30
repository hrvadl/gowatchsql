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

type formState struct {
	page         focused
	inputFocused focused
}

type state struct {
	active bool
	form   formState
}
