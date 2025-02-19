package syntaxhighlighting

import (
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal"
	commands "github.com/cjairm/devgita/internal"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/files"
)

type Syntaxhighlighting struct {
	Cmd cmd.Command
}

func New() *Syntaxhighlighting {
	osCmd := cmd.NewCommand()
	return &Syntaxhighlighting{Cmd: osCmd}
}

func (a *Syntaxhighlighting) Install() error {
	return a.Cmd.InstallPackage("zsh-syntax-highlighting")
}

func (a *Syntaxhighlighting) MaybeInstall() error {
	return a.Cmd.MaybeInstallPackage("zsh-syntax-highlighting")
}

func (a *Syntaxhighlighting) Setup() error {
	devgitaCustomDir, err := config.GetDevgitaConfigDir()
	if err != nil {
		return err
	}
	err = commands.AddLineToFile(
		"source $(brew --prefix)/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh",
		devgitaCustomDir+"/devgita.zsh",
	)
	if err != nil {
		return err
	}
	return nil
}

func (a *Syntaxhighlighting) MaybeSetup() error {
	devgitaCustomDir, err := config.GetDevgitaConfigDir()
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
