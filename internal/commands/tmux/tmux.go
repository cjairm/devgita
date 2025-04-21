// -------------------------
// TODO: Write documentation how to use this
// - Tmux documentation: https://github.com/tmux/tmux
// - Personal configuration: https://github.com/cjairm/devenv/tree/main/tmux
// - Releases: https://github.com/tmux/tmux/releases
// - Installing instructions: https://github.com/tmux/tmux/wiki/Installing
// -------------------------

package tmux

import (
	"os"
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal"
	"github.com/cjairm/devgita/pkg/files"
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

func Command(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		Verbose:     true,
		IsSudo:      false,
		Command:     "tmux",
		Args:        args,
	}
	return cmd.ExecCommand(execCommand)
}

func (t *Tmux) Install() error {
	return t.Cmd.InstallPackage("tmux")
}

func (t *Tmux) MaybeInstall() error {
	return t.Cmd.MaybeInstallPackage("tmux")
}

func (t *Tmux) Setup() error {
	return t.Base.CopyAppConfigFileToHomeDir("tmux", ".tmux.conf")
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
