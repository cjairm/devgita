package ulauncher

import (
	"fmt"

	"github.com/cjairm/devgita/internal/apps"
	"github.com/cjairm/devgita/internal/apps/baseapp"
	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
)

var _ apps.App = (*Ulauncher)(nil)

type Ulauncher struct {
	Cmd commands.Command
}

func (u *Ulauncher) Name() string       { return constants.Ulauncher }
func (u *Ulauncher) Kind() apps.AppKind { return apps.KindDesktop }

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
	return baseapp.Reinstall(u.Install, u.Uninstall)
}

func (u *Ulauncher) Uninstall() error {
	return fmt.Errorf("%w for ulauncher — manage via system package manager", apps.ErrUninstallNotSupported)
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
	return fmt.Errorf("%w for ulauncher — use system package manager", apps.ErrUpdateNotSupported)
}
