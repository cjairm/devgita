package bash

import (
	"os"
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal"
	commands "github.com/cjairm/devgita/internal"
	"github.com/cjairm/devgita/pkg/files"
)

type Bash struct {
	Cmd  cmd.Command
	Base cmd.BaseCommand
}

func New() *Bash {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Bash{Cmd: osCmd, Base: *baseCmd}
}

func (b *Bash) CopyCustomConfig() error {
	err := files.MoveFromConfigsToLocalConfig([]string{"bash"}, []string{"devgita"})
	if err != nil {
		return err
	}
	return nil
}

func (b *Bash) MaybeCopyCustomConfig() error {
	devgitaCustomDir, err := b.Base.GetDevgitaAppDir("")
	if err != nil {
		return err
	}
	devgitaConfigFile := filepath.Join(devgitaCustomDir, "devgita.zsh")
	isFilePresent := files.FileAlreadyExist(devgitaConfigFile)
	if isFilePresent {
		return nil
	}
	return b.CopyCustomConfig()
}

func (b *Bash) SetupCustom(line string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	zshConfigFile := filepath.Join(homeDir, ".zshrc")
	return commands.AddLineToFile(line, zshConfigFile)
}

func (b *Bash) MaybeSetupCustom(line, toSearch string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	zshConfigFile := filepath.Join(homeDir, ".zshrc")
	isAlreadySetup, err := files.ContentExistsInFile(zshConfigFile, toSearch)
	if err != nil {
		return err
	}
	if isAlreadySetup == true {
		return nil
	}
	return b.SetupCustom(line)
}
