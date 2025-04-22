package autosuggestions

import (
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

type Autosuggestions struct {
	Cmd  cmd.Command
	Base cmd.BaseCommand
}

func New() *Autosuggestions {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Autosuggestions{Cmd: osCmd, Base: *baseCmd}
}

func (a *Autosuggestions) Install() error {
	return a.Cmd.InstallPackage("zsh-autosuggestions")
}

func (a *Autosuggestions) MaybeInstall() error {
	return a.Cmd.MaybeInstallPackage("zsh-autosuggestions")
}

func (a *Autosuggestions) Setup() error {
	return files.AddLineToFile(
		"source $(brew --prefix)/share/zsh-autosuggestions/zsh-autosuggestions.zsh",
		filepath.Join(paths.AppDir, "devgita.zsh"),
	)
}

func (a *Autosuggestions) MaybeSetup() error {
	isConfigured, err := files.ContentExistsInFile(
		filepath.Join(paths.AppDir, "devgita.zsh"),
		"zsh-autosuggestions.zsh",
	)
	if err != nil {
		return err
	}
	if isConfigured == true {
		return nil
	}
	return a.Setup()
}
