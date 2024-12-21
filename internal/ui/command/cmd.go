package command

type Command string

const (
	Context Command = "contexts"
	Tables  Command = "tables"
	Exit    Command = "exit"
)
