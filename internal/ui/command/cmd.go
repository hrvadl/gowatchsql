package command

type Command string

const (
	Context Command = "contexts"
	Query   Command = "query"
	Tables  Command = "tables"
	Exit    Command = "exit"
)
