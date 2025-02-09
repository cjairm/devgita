package git

import cmd "github.com/cjairm/devgita/internal"

type Git struct {
	Cmd cmd.Command
}

func NewGit() *Git {
	osCmd := cmd.NewCommand()
	return &Git{Cmd: osCmd}
}

func Command(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		Verbose:     true,
		IsSudo:      false,
		Command:     "git",
		Args:        args,
	}
	return cmd.ExecCommand(execCommand)
}

func (g *Git) Clean(branch string, isForced bool) error {
	deleteArg := "-d"
	if isForced {
		deleteArg = "-D"
	}
	return Command("branch", deleteArg, branch)
}

func (g *Git) Clone(url, dstPath string) error {
	return Command("clone", url, dstPath)
}

func (g *Git) FetchOrigin() error {
	return Command("fetch", "origin")
}

func (g *Git) Install() error {
	return g.Cmd.MaybeInstallPackage("git")
}

func (g *Git) Pull(branch string) error {
	return Command("pull", "origin", branch)
}

func (g *Git) SwitchBranch(branch string) error {
	return Command("checkout", branch)
}
