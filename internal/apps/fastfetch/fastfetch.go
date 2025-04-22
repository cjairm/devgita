// -------------------------
// NOTE: Write documentation or create icon to open and get information of this Mac
// - Documentation: https://github.com/fastfetch-cli/fastfetch
// -------------------------

package fastfetch

import (
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

type Fastfetch struct {
	Cmd  cmd.Command
	Base cmd.BaseCommand
}

func New() *Fastfetch {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Fastfetch{Cmd: osCmd, Base: *baseCmd}
}

func Command(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		Verbose:     true,
		IsSudo:      false,
		Command:     constants.Fastfetch,
		Args:        args,
	}
	return cmd.ExecCommand(execCommand)
}

func (f *Fastfetch) Install() error {
	return f.Cmd.InstallPackage("fastfetch")
}

func (f *Fastfetch) MaybeInstall() error {
	return f.Cmd.MaybeInstallPackage("fastfetch")
}

func (f *Fastfetch) Setup() error {
	return files.CopyDir(paths.FastFetchConfigAppDir, paths.FastFetchConfigLocalDir)
}

func (f *Fastfetch) MaybeSetup() error {
	fastfetchConfigFile := filepath.Join(paths.FastFetchConfigLocalDir, "config.jsonc")
	isFilePresent := files.FileAlreadyExist(fastfetchConfigFile)
	if isFilePresent {
		return nil
	}
	return f.Setup()
}
