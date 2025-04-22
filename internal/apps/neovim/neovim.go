// -------------------------
// TODO: Write documentation how to use this
// - Kickstart documentation: https://github.com/nvim-lua/kickstart.nvim?tab=readme-ov-file
// - Personal configuration: https://github.com/cjairm/devenv/blob/main/nvim/init.lua
// - Releases: https://github.com/neovim/neovim/releases
//
// NOTE: Is it possible to install different themes?
// If so, see more here: https://linovox.com/the-best-color-schemes-for-neovim-nvim/
// -------------------------

package neovim

import (
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

type Neovim struct {
	Cmd  cmd.Command
	Base cmd.BaseCommand
}

func New() *Neovim {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Neovim{Cmd: osCmd, Base: *baseCmd}
}

func Command(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		Verbose:     true,
		IsSudo:      false,
		Command:     constants.Nvim,
		Args:        args,
	}
	return cmd.ExecCommand(execCommand)
}

func (n *Neovim) Install() error {
	return n.Cmd.InstallPackage("neovim")
}

func (n *Neovim) MaybeInstall() error {
	return n.Cmd.MaybeInstallPackage("neovim")
}

func (n *Neovim) Setup() error {
	return files.CopyDir(paths.NeovimConfigAppDir, paths.NvimConfigLocalDir)
}

func (n *Neovim) MaybeSetup() error {
	isFilePresent := files.FileAlreadyExist(filepath.Join(paths.NvimConfigLocalDir, "init.lua"))
	if isFilePresent {
		return nil
	}
	return n.Setup()
}
