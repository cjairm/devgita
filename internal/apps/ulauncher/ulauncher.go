package ulauncher

import (
	"fmt"

	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
)

type Ulauncher struct {
	Cmd commands.Command
}

func New() *Ulauncher {
	return &Ulauncher{Cmd: commands.NewCommand()}
}

func (u *Ulauncher) Install() error {
	return u.Cmd.InstallDesktopApp(constants.Ulauncher)
}

func (u *Ulauncher) SoftInstall() error {
	return u.Cmd.MaybeInstallDesktopApp(constants.Ulauncher)
}

func (u *Ulauncher) ForceInstall() error {
	if err := u.Uninstall(); err != nil {
		return fmt.Errorf("failed to uninstall ulauncher before force install: %w", err)
	}
	return u.Install()
}

func (u *Ulauncher) Uninstall() error {
	return fmt.Errorf("ulauncher uninstall not supported - manage via system package manager")
}

func (u *Ulauncher) ForceConfigure() error {
	// Ulauncher uses GUI-based configuration
	// No config files to manage
	return nil
}

func (u *Ulauncher) SoftConfigure() error {
	// Ulauncher uses GUI-based configuration
	// No config files to manage
	return nil
}

func (u *Ulauncher) ExecuteCommand(args ...string) error {
	// Ulauncher is a desktop application without CLI commands typically managed by devgita
	// Return success for interface compliance
	return nil
}

func (u *Ulauncher) Update() error {
	return fmt.Errorf("ulauncher update not implemented - use system package manager")
}
