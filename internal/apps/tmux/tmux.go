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
	return t.Cmd.InstallPackage("tmux")
}

func (t *Tmux) MaybeInstall() error {
	return t.Cmd.MaybeInstallPackage("tmux")
}

func (t *Tmux) Setup() error {
	return files.CopyFile(
		filepath.Join(paths.TmuxConfigAppDir, ".tmux.conf"),
		filepath.Join(paths.HomeDir, ".tmux.conf"),
	)
}

func (t *Tmux) MaybeSetup() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	tmuxConfigFile := filepath.Join(homeDir, ".tmux.conf")
	isFilePresent := files.FileAlreadyExist(tmuxConfigFile)
	if isFilePresent {
		return nil
	}
	return t.Setup()
}

func (t *Tmux) Run(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		IsSudo:      false,
		Command:     "tmux",
		Args:        args,
	}
	if _, _, err := t.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run tmux command: %w", err)
	}
	return nil
}
