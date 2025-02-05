package commands

type CommandParams struct {
	PreExecMsg  string
	PostExecMsg string
	IsSudo      bool
	Verbose     bool
	Command     string
	Args        []string
}

type Command interface {
	ExecCommand(cmd CommandParams) error
	PrintSomething()
}
