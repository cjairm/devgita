// -------------------------
// NOTE: Git version control system with devgita integration
// - **Documentation**: https://git-scm.com/doc
// - **Devgita documentation**: Git is the distributed version control system. Here are some useful commands to get you started:
//   - Check status: `git status`
//   - Clone repository: `git clone <url> <directory>`
//   - Switch branch: `git checkout <branch>`
//   - Create and switch to new branch: `git checkout -b <branch>`
//   - Stage changes: `git add .`
//   - Commit changes: `git commit -m "message"`
//   - Push changes: `git push origin <branch>`
//   - Pull changes: `git pull origin <branch>`
//   - Apply stashed changes: `git stash pop`
//   - View commit history: `git log`
//   - Clean ignored files: `git clean -X -d -f`
// -------------------------

package git

import (
	"fmt"
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
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

func (g *Git) ForceInstall() error {
	return g.Cmd.InstallPackage(constants.Git)
}

func (g *Git) SoftInstall() error {
	return g.Cmd.MaybeInstallPackage(constants.Git)
}

func (g *Git) Uninstall() error {
	return fmt.Errorf("git uninstall not supported through devgita")
}

// Standardized configuration methods
func configure() error {
	return files.CopyDir(paths.GitConfigAppDir, paths.GitConfigLocalDir)
}

func (g *Git) ForceConfigure() error {
	return configure()
}

func (g *Git) SoftConfigure() error {
	configFile := filepath.Join(paths.GitConfigLocalDir, ".gitconfig")
	if files.FileAlreadyExist(configFile) {
		return nil
	}
	return configure()
}

// Standardized execution method
func (g *Git) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		IsSudo:      false,
		Command:     constants.Git,
		Args:        args,
	}
	if _, _, err := g.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run git command: %w", err)
	}
	return nil
}

// Git-specific methods
func (g *Git) Clone(url, dstPath string) error {
	return g.ExecuteCommand("clone", url, dstPath)
}

func (g *Git) DeleteBranch(branch string, isForced bool) error {
	deleteArg := "-d"
	if isForced {
		deleteArg = "-D"
	}
	return g.ExecuteCommand("branch", deleteArg, branch)
}

func (g *Git) DeepClean(url, dstPath string) error {
	// -X: This option tells Git to remove only the files that are ignored by Git (i.e., files that are listed in your .gitignore file). It does not remove untracked files that are not ignored.
	// -d: This option allows Git to remove untracked directories in addition to untracked files.
	// -f: This option stands for "force." Git requires this option to actually perform the clean operation, as a safety measure to prevent accidental data loss.
	return g.ExecuteCommand("clean", "-X", "-d", "-f")
}

func (g *Git) FetchOrigin() error {
	return g.ExecuteCommand("fetch", "origin")
}

func (g *Git) Pop(branch string) error {
	return g.ExecuteCommand("stash", "pop")
}

func (g *Git) Pull(branch string) error {
	if branch == "" {
		return g.ExecuteCommand("pull")
	}
	return g.ExecuteCommand("pull", "origin", branch)
}

func (g *Git) SwitchBranch(branch string) error {
	return g.ExecuteCommand("checkout", branch)
}

func (g *Git) Restore(branch, files string) error {
	if branch == "" {
		branch = "main"
	}
	return g.ExecuteCommand("restore", "--source", branch, "--", files)
}
