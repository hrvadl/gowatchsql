package mainpanel

type active int

const (
	objectsActive active = iota
	contextsActive
	newContextActive
	queryRunActive
)

type state struct {
	active active
}
