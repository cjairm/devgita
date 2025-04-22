package syntaxhighlighting

import (
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
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
	return files.AddLineToFile(
		"source $(brew --prefix)/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh",
		filepath.Join(paths.AppDir, "devgita.zsh"),
	)
}

func (a *Syntaxhighlighting) MaybeSetup() error {
	isConfigured, err := files.ContentExistsInFile(
		filepath.Join(paths.AppDir, "devgita.zsh"),
		"zsh-syntax-highlighting.zsh",
	)
	if err != nil {
		return err
	}
	if isConfigured == true {
		return nil
	}
	return a.Setup()
}
