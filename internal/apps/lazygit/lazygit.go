// LazyGit terminal UI for Git repository management with devgita integration
//
// LazyGit is a simple terminal UI for git commands, written in Go with the gocui library.
// It provides an interactive interface to manage Git repositories, branches, commits, and
// staging operations, all from the comfort of the terminal.
//
// References:
// - LazyGit Repository: https://github.com/jesseduffield/lazygit
// - LazyGit Documentation: https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md
//
// Common lazygit commands available through ExecuteCommand():
//   - lazygit - Launch interactive TUI
//   - lazygit --version - Show lazygit version information
//   - lazygit --config - Show configuration file path
//   - lazygit --help - Display help information

package lazygit

import (
	"fmt"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
)

type LazyGit struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func New() *LazyGit {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &LazyGit{Cmd: osCmd, Base: baseCmd}
}

func (lg *LazyGit) Install() error {
	return lg.Cmd.InstallPackage(constants.LazyGit)
}

func (lg *LazyGit) SoftInstall() error {
	return lg.Cmd.MaybeInstallPackage(constants.LazyGit)
}

func (lg *LazyGit) ForceInstall() error {
	err := lg.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall lazygit: %w", err)
	}
	return lg.Install()
}

func (lg *LazyGit) Uninstall() error {
	return fmt.Errorf("lazygit uninstall not supported through devgita")
}

func (lg *LazyGit) ForceConfigure() error {
	// LazyGit configuration is optional and user-specific
	// Configuration file typically located at ~/.config/lazygit/config.yml
	// For now, no default configuration is applied
	return nil
}

func (lg *LazyGit) SoftConfigure() error {
	// LazyGit configuration is optional and user-specific
	// Configuration file typically located at ~/.config/lazygit/config.yml
	// For now, no default configuration is applied
	return nil
}

func (lg *LazyGit) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		IsSudo:  false,
		Command: constants.LazyGit,
		Args:    args,
	}
	if _, _, err := lg.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run lazygit command: %w", err)
	}
	return nil
}

func (lg *LazyGit) Update() error {
	return fmt.Errorf("lazygit update not implemented through devgita")
}
