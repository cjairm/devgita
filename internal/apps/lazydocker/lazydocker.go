// LazyDocker terminal UI for Docker container and image management with devgita integration
//
// LazyDocker is a simple terminal UI for both docker and docker-compose, written in Go
// with the gocui library. It provides an interactive interface to manage Docker containers,
// images, volumes, and networks, all from the comfort of the terminal.
//
// References:
// - LazyDocker Repository: https://github.com/jesseduffield/lazydocker
// - LazyDocker Documentation: https://github.com/jesseduffield/lazydocker/blob/master/docs/Config.md
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
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/constants"
)

type LazyDocker struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

var packageName = fmt.Sprintf("jesseduffield/%s/%s", constants.LazyDocker, constants.LazyDocker)

func New() *LazyDocker {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &LazyDocker{Cmd: osCmd, Base: baseCmd}
}

func (ld *LazyDocker) Install() error {
	return ld.Cmd.InstallPackage(packageName)
}

func (ld *LazyDocker) SoftInstall() error {
	return ld.Cmd.MaybeInstallPackage(packageName, constants.LazyDocker)
}

func (ld *LazyDocker) ForceInstall() error {
	err := ld.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall lazydocker: %w", err)
	}
	return ld.Install()
}

func (ld *LazyDocker) Uninstall() error {
	return fmt.Errorf("lazydocker uninstall not supported through devgita")
}

func (ld *LazyDocker) ForceConfigure() error {
	gc := &config.GlobalConfig{}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	gc.EnableShellFeature(constants.LazyDocker)
	if err := gc.RegenerateShellConfig(); err != nil {
		return fmt.Errorf("failed to generate shell config: %w", err)
	}
	if err := gc.Save(); err != nil {
		return fmt.Errorf("failed to save global config: %w", err)
	}
	return nil
}

func (ld *LazyDocker) SoftConfigure() error {
	gc := &config.GlobalConfig{}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	if gc.IsShellFeatureEnabled(constants.LazyDocker) {
		return nil
	}
	return ld.ForceConfigure()
}

func (ld *LazyDocker) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		Command: constants.LazyDocker,
		Args:    args,
	}
	if _, _, err := ld.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run lazydocker command: %w", err)
	}
	return nil
}

func (ld *LazyDocker) Update() error {
	return fmt.Errorf("lazydocker update not implemented through devgita")
}
