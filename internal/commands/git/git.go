package git

import (
	"fmt"

	cmd "github.com/cjairm/devgita/internal"
)

type Git struct {
	Cmd cmd.Command
}

func New() *Git {
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

func (g *Git) Install() error {
	return g.Cmd.MaybeInstallPackage("git")
}

func (g *Git) Clone(url, dstPath string) error {
	return Command("clone", url, dstPath)
}

func (g *Git) DeleteBranch(branch string, isForced bool) error {
	deleteArg := "-d"
	if isForced {
		deleteArg = "-D"
	}
	return Command("branch", deleteArg, branch)
}

func (g *Git) DeepClean(url, dstPath string) error {
	// -X: This option tells Git to remove only the files that are ignored by Git (i.e., files that are listed in your .gitignore file). It does not remove untracked files that are not ignored.
	// -d: This option allows Git to remove untracked directories in addition to untracked files.
	// -f: This option stands for "force." Git requires this option to actually perform the clean operation, as a safety measure to prevent accidental data loss.
	return Command("clean", "-X", "-d", "-f")
}

func (g *Git) FetchOrigin() error {
	return Command("fetch", "origin")
}

func (g *Git) Pop(branch string) error {
	return Command("pop")
}

func (g *Git) Pull(branch string) error {
	if branch == "" {
		return Command("pull")
	}
	return Command("pull", "origin", branch)
}

func (g *Git) SwitchBranch(branch string) error {
	return Command("checkout", branch)
}

func (g *Git) Stash(branch string) error {
	return Command("stash")
}

func (g *Git) Restore(branch, files string) error {
	if branch == "" {
		branch = "main"
	}
	return Command(fmt.Sprintf("restore --source %s --", branch), files)
}
