package syntaxhighlighting

import (
	"errors"
	"fmt"
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

type Syntaxhighlighting struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func New() *Syntaxhighlighting {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Syntaxhighlighting{Cmd: osCmd, Base: baseCmd}
}

func (a *Syntaxhighlighting) Install() error {
	return a.Cmd.InstallPackage(constants.Syntaxhighlighting)
}

func (a *Syntaxhighlighting) ForceInstall() error {
	err := a.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall syntaxhighlighting before force install: %w", err)
	}
	return a.Install()
}

func (a *Syntaxhighlighting) SoftInstall() error {
	return a.Cmd.MaybeInstallPackage(constants.Syntaxhighlighting)
}

func (a *Syntaxhighlighting) ForceConfigure() error {
	return files.AddLineToFile(
		"source $(brew --prefix)/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh",
		filepath.Join(paths.AppDir, "devgita.zsh"),
	)
}

func (a *Syntaxhighlighting) SoftConfigure() error {
	isConfigured, err := files.ContentExistsInFile(
		filepath.Join(paths.AppDir, "devgita.zsh"),
		constants.Syntaxhighlighting+".zsh",
	)
	if err != nil {
		return err
	}
	if isConfigured {
		return nil
	}
	return a.ForceConfigure()
}

func (a *Syntaxhighlighting) Uninstall() error {
	return errors.New(constants.Syntaxhighlighting + " uninstall is not supported")
}

func (a *Syntaxhighlighting) ExecuteCommand(args ...string) error {
	return nil
}

func (a *Syntaxhighlighting) Update() error {
	return errors.New(
		constants.Syntaxhighlighting + " update is not implemented - use system package manager",
	)
}
