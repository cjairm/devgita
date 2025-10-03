package git

import (
	"fmt"

	cmd "github.com/cjairm/devgita/internal/commands"
)

type Git struct {
	Cmd  cmd.Command
	Base cmd.BaseCommand
}

func New() *Git {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Git{Cmd: osCmd, Base: *baseCmd}
}

func (g *Git) Install() error {
	return g.Cmd.MaybeInstallPackage("git")
}

func (g *Git) Clone(url, dstPath string) error {
	return g.Run("clone", url, dstPath)
}

func (g *Git) DeleteBranch(branch string, isForced bool) error {
	deleteArg := "-d"
	if isForced {
		deleteArg = "-D"
	}
	return g.Run("branch", deleteArg, branch)
}

func (g *Git) DeepClean(url, dstPath string) error {
	// -X: This option tells Git to remove only the files that are ignored by Git (i.e., files that are listed in your .gitignore file). It does not remove untracked files that are not ignored.
	// -d: This option allows Git to remove untracked directories in addition to untracked files.
	// -f: This option stands for "force." Git requires this option to actually perform the clean operation, as a safety measure to prevent accidental data loss.
	return g.Run("clean", "-X", "-d", "-f")
}

func (g *Git) FetchOrigin() error {
	return g.Run("fetch", "origin")
}

func (g *Git) Pop(branch string) error {
	return g.Run("pop")
}

func (g *Git) Pull(branch string) error {
	if branch == "" {
		return g.Run("pull")
	}
	return g.Run("pull", "origin", branch)
}

func (g *Git) SwitchBranch(branch string) error {
	return g.Run("checkout", branch)
}

func (g *Git) Stash(branch string) error {
	return g.Run("stash")
}

func (g *Git) Restore(branch, files string) error {
	if branch == "" {
		branch = "main"
	}
	return g.Run(fmt.Sprintf("restore --source %s --", branch), files)
}

func (g *Git) Run(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		IsSudo:      false,
		Command:     "git",
		Args:        args,
	}
	if _, _, err := g.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to execute git command: %w", err)
	}
	return nil
}
