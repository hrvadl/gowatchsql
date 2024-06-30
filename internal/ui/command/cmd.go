package command

type Command string

const (
	Context    Command = "contexts"
	Tables     Command = "tables"
	NewContext Command = "new context"
	Exit       Command = "exit"
)
