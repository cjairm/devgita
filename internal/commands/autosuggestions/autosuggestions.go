package autosuggestions

import (
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal"
	commands "github.com/cjairm/devgita/internal"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/utils"
)

type Autosuggestions struct {
	Cmd cmd.Command
}

func NewAutosuggestions() *Autosuggestions {
	osCmd := cmd.NewCommand()
	return &Autosuggestions{Cmd: osCmd}
}

func (a *Autosuggestions) Install() error {
	return a.Cmd.InstallPackage("zsh-autosuggestions")
}

func (a *Autosuggestions) MaybeInstall() error {
	return a.Cmd.MaybeInstallPackage("zsh-autosuggestions")
}

func (a *Autosuggestions) Setup() error {
	devgitaCustomDir, err := utils.GetDevgitaConfigDir()
	if err != nil {
		return err
	}
	err = commands.AddLineToFile(
		"source $(brew --prefix)/share/zsh-autosuggestions/zsh-autosuggestions.zsh",
		devgitaCustomDir+"/devgita.zsh",
	)
	if err != nil {
		return err
	}
	return nil
}

func (a *Autosuggestions) MaybeSetup() error {
	devgitaCustomDir, err := utils.GetDevgitaConfigDir()
	if err != nil {
		return err
	}
	devgitaConfigFile := filepath.Join(devgitaCustomDir, "devgita.zsh")
	isConfigured, err := files.ContentExistsInFile(
		devgitaConfigFile,
		"zsh-autosuggestions.zsh",
	)
	if isConfigured == true {
		return nil
	}
	return a.Setup()
}
