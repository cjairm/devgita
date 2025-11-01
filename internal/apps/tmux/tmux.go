// -------------------------
// TODO: Write documentation how to use this
// - Tmux documentation: https://github.com/tmux/tmux
// - Personal configuration: https://github.com/cjairm/devenv/tree/main/tmux
// - Releases: https://github.com/tmux/tmux/releases
// - Installing instructions: https://github.com/tmux/tmux/wiki/Installing
// -------------------------

package tmux

import (
	"fmt"
	"os"
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

type Tmux struct {
	Cmd  cmd.Command
	Base cmd.BaseCommand
}

func New() *Tmux {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Tmux{Cmd: osCmd, Base: *baseCmd}
}

func (t *Tmux) Install() error {
	return t.Cmd.InstallPackage(constants.Tmux)
}

func (t *Tmux) ForceInstall() error {
	if err := t.Uninstall(); err != nil {
		return fmt.Errorf("failed to uninstall tmux before force install: %w", err)
	}
	return t.Install()
}

func (t *Tmux) SoftInstall() error {
	return t.Cmd.MaybeInstallPackage(constants.Tmux)
}

func (t *Tmux) ForceConfigure() error {
	return files.CopyFile(
		filepath.Join(paths.TmuxConfigAppDir, ".tmux.conf"),
		filepath.Join(paths.HomeDir, ".tmux.conf"),
	)
}

func (t *Tmux) SoftConfigure() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	tmuxConfigFile := filepath.Join(homeDir, ".tmux.conf")
	isFilePresent := files.FileAlreadyExist(tmuxConfigFile)
	if isFilePresent {
		return nil
	}
	return t.ForceConfigure()
}

func (t *Tmux) Uninstall() error {
	return fmt.Errorf("tmux uninstall is not supported")
}

func (t *Tmux) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		IsSudo:      false,
		Command:     constants.Tmux,
		Args:        args,
	}
	if _, _, err := t.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to execute tmux command: %w", err)
	}
	return nil
}

func (t *Tmux) Update() error {
	return fmt.Errorf("tmux update is not implemented")
}
