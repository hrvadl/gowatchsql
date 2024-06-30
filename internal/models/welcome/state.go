package welcome

type focus int

const (
	cmdFocused focus = iota
	mainFocused
)

type state struct {
	active    focus
	showModal bool
}
