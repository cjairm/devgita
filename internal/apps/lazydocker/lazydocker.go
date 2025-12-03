// Lazydocker terminal UI for Docker container and image management with devgita integration
//
// Lazydocker is a simple terminal UI for both docker and docker-compose, written in Go
// with the gocui library. It provides an interactive interface to manage Docker containers,
// images, volumes, and networks, all from the comfort of the terminal.
//
// References:
// - Lazydocker Repository: https://github.com/jesseduffield/lazydocker
// - Lazydocker Documentation: https://github.com/jesseduffield/lazydocker/blob/master/docs/Config.md
//
// Common lazydocker commands available through ExecuteCommand():
//   - lazydocker - Launch interactive TUI
//   - lazydocker --version - Show lazydocker version information
//   - lazydocker --config - Show configuration file path
//   - lazydocker --help - Display help information

package lazydocker

import (
	"fmt"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
)

type Lazydocker struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

var packageName = fmt.Sprintf("jesseduffield/%s/%s", constants.Lazydocker, constants.Lazydocker)

func New() *Lazydocker {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Lazydocker{Cmd: osCmd, Base: baseCmd}
}

func (l *Lazydocker) Install() error {
	return l.Cmd.InstallPackage(packageName)
}

func (l *Lazydocker) SoftInstall() error {
	return l.Cmd.MaybeInstallPackage(packageName, constants.Lazydocker)
}

func (l *Lazydocker) ForceInstall() error {
	err := l.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall lazydocker: %w", err)
	}
	return l.Install()
}

func (l *Lazydocker) Uninstall() error {
	return fmt.Errorf("lazydocker uninstall not supported through devgita")
}

func (l *Lazydocker) ForceConfigure() error {
	// Lazydocker configuration is optional and user-specific
	// Configuration file typically located at ~/.config/lazydocker/config.yml
	// For now, no default configuration is applied
	return nil
}

func (l *Lazydocker) SoftConfigure() error {
	// Lazydocker configuration is optional and user-specific
	// Configuration file typically located at ~/.config/lazydocker/config.yml
	// For now, no default configuration is applied
	return nil
}

func (l *Lazydocker) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		IsSudo:  false,
		Command: constants.Lazydocker,
		Args:    args,
	}
	if _, _, err := l.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run lazydocker command: %w", err)
	}
	return nil
}

func (l *Lazydocker) Update() error {
	return fmt.Errorf("lazydocker update not implemented through devgita")
}
