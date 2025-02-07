package git

import cmd "github.com/cjairm/devgita/internal"

type Git struct {
	Cmd cmd.Command
}

func NewGit() *Git {
	osCmd := cmd.NewCommand()
	return &Git{Cmd: osCmd}
}

func (g *Git) Install() error {
	return g.Cmd.InstallPackage("git")
}

func (g *Git) SwitchBranch(branch string) error {
	return g.Cmd.GitCommand("checkout", branch)
}

func (g *Git) FetchOrigin() error {
	return g.Cmd.GitCommand("fetch", "origin")
}

func (g *Git) Pull(branch string) error {
	return g.Cmd.GitCommand("pull", "origin", branch)
}

func (g *Git) Clean(branch string, isForced bool) error {
	deleteArg := "-d"
	if isForced {
		deleteArg = "-D"
	}
	return g.Cmd.GitCommand("branch", deleteArg, branch)
}
