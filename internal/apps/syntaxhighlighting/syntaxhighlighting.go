package syntaxhighlighting

import (
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/files"
)

type Syntaxhighlighting struct {
	Cmd  cmd.Command
	Base cmd.BaseCommand
}

func New() *Syntaxhighlighting {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Syntaxhighlighting{Cmd: osCmd, Base: *baseCmd}
}

func (a *Syntaxhighlighting) Install() error {
	return a.Cmd.InstallPackage("zsh-syntax-highlighting")
}

func (a *Syntaxhighlighting) MaybeInstall() error {
	return a.Cmd.MaybeInstallPackage("zsh-syntax-highlighting")
}

func (a *Syntaxhighlighting) Setup() error {
	devgitaCustomDir, err := a.Base.AppDir()
	if err != nil {
		return err
	}
	err = files.AddLineToFile(
		"source $(brew --prefix)/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh",
		devgitaCustomDir+"/devgita.zsh",
	)
	if err != nil {
		return err
	}
	return nil
}

func (a *Syntaxhighlighting) MaybeSetup() error {
	devgitaCustomDir, err := a.Base.AppDir()
	if err != nil {
		return err
	}
	devgitaConfigFile := filepath.Join(devgitaCustomDir, "devgita.zsh")
	isConfigured, err := files.ContentExistsInFile(
		devgitaConfigFile,
		"zsh-syntax-highlighting.zsh",
	)
	if isConfigured == true {
		return nil
	}
	return a.Setup()
}
