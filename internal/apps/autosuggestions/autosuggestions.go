// Autosuggestions module provides installation and configuration management for zsh-autosuggestions with devgita integration.
// zsh-autosuggestions suggests commands as you type based on history and completions.

package autosuggestions

import (
	"fmt"
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

type Autosuggestions struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func New() *Autosuggestions {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Autosuggestions{Cmd: osCmd, Base: baseCmd}
}

func (a *Autosuggestions) Install() error {
	return a.Cmd.InstallPackage("zsh-autosuggestions")
}

func (a *Autosuggestions) ForceInstall() error {
	err := a.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall autosuggestions: %w", err)
	}
	return a.Install()
}

func (a *Autosuggestions) SoftInstall() error {
	return a.Cmd.MaybeInstallPackage("zsh-autosuggestions")
}

func (a *Autosuggestions) ForceConfigure() error {
	return files.AddLineToFile(
		"source $(brew --prefix)/share/zsh-autosuggestions/zsh-autosuggestions.zsh",
		filepath.Join(paths.AppDir, "devgita.zsh"),
	)
}

func (a *Autosuggestions) SoftConfigure() error {
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
	return a.ForceConfigure()
}

func (a *Autosuggestions) Uninstall() error {
	return fmt.Errorf("uninstall not implemented for autosuggestions")
}

func (a *Autosuggestions) ExecuteCommand(args ...string) error {
	return nil
}

func (a *Autosuggestions) Update() error {
	return fmt.Errorf("update not implemented for autosuggestions")
}
