package aerospace

import (
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

type Aerospace struct {
	Cmd  cmd.Command
	Base cmd.BaseCommand
}

func New() *Aerospace {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Aerospace{Cmd: osCmd, Base: *baseCmd}
}

func (a *Aerospace) Install() error {
	return a.Cmd.InstallDesktopApp("nikitabobko/tap/aerospace")
}

func (a *Aerospace) MaybeInstall() error {
	return a.Cmd.MaybeInstallDesktopApp("nikitabobko/tap/aerospace")
}

func (a *Aerospace) Setup() error {
	return files.CopyDir(paths.AerospaceConfigAppDir, paths.AerospaceConfigLocalDir)
}

func (a *Aerospace) MaybeSetup() error {
	aerospaceConfigFile := filepath.Join(paths.AerospaceConfigLocalDir, "aerospace.toml")
	if isFilePresent := files.FileAlreadyExist(aerospaceConfigFile); isFilePresent {
		return nil
	}
	return a.Setup()
}
