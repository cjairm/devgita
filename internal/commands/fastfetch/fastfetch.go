// -------------------------
// NOTE: Write documentation or create icon to open and get information of this Mac
// - Documentation: https://github.com/fastfetch-cli/fastfetch
// -------------------------

package fastfetch

import (
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/files"
)

const fastfetchDir = "fastfetch"

type Fastfetch struct {
	Cmd cmd.Command
}

func New() *Fastfetch {
	osCmd := cmd.NewCommand()
	return &Fastfetch{Cmd: osCmd}
}

func Command(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		Verbose:     true,
		IsSudo:      false,
		Command:     "fastfetch",
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
	filePath := []string{fastfetchDir}
	return files.MoveFromConfigsToLocalConfig(filePath, filePath)
}

func (f *Fastfetch) MaybeSetup() error {
	localConfig, err := config.GetLocalConfigPath()
	if err != nil {
		return err
	}
	fastfetchConfigFile := filepath.Join(localConfig, "fastfetch", "config.jsonc")
	isFilePresent := files.FileAlreadyExist(fastfetchConfigFile)
	if isFilePresent {
		return nil
	}
	return f.Setup()
}
