// Package raycast provides installation and management for Raycast productivity launcher.
//
// Raycast is a blazingly fast, extendable launcher for macOS. This module follows the
// standardized devgita app interface for desktop applications.

package raycast

import (
	"fmt"

	"github.com/cjairm/devgita/internal/apps"
	"github.com/cjairm/devgita/internal/apps/baseapp"
	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
)

var _ apps.App = (*Raycast)(nil)

type Raycast struct {
	Cmd cmd.Command
}

func (r *Raycast) Name() string       { return constants.Raycast }
func (r *Raycast) Kind() apps.AppKind { return apps.KindDesktop }

func New() *Raycast {
	return &Raycast{Cmd: cmd.NewCommand()}
}

func (r *Raycast) Install() error {
	return r.Cmd.InstallDesktopApp(constants.Raycast)
}

func (r *Raycast) ForceInstall() error {
	return baseapp.Reinstall(r.Install, r.Uninstall)
}

func (r *Raycast) SoftInstall() error {
	return r.Cmd.MaybeInstallDesktopApp(constants.Raycast)
}

func (r *Raycast) ForceConfigure() error {
	// No configuration needed for GUI-based launcher
	return nil
}

func (r *Raycast) SoftConfigure() error {
	// No configuration needed for GUI-based launcher
	return nil
}

func (r *Raycast) Uninstall() error {
	return fmt.Errorf("%w for raycast", apps.ErrUninstallNotSupported)
}

func (r *Raycast) ExecuteCommand(args ...string) error {
	// No CLI commands for desktop GUI application
	return nil
}

func (r *Raycast) Update() error {
	return fmt.Errorf("%w for raycast", apps.ErrUpdateNotSupported)
}
