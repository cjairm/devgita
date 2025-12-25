// Package raycast provides installation and management for Raycast productivity launcher.
//
// Raycast is a blazingly fast, extendable launcher for macOS. This module follows the
// standardized devgita app interface for desktop applications.

package raycast

import (
	"fmt"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
)

type Raycast struct {
	Cmd cmd.Command
}

func New() *Raycast {
	return &Raycast{Cmd: cmd.NewCommand()}
}

func (r *Raycast) Install() error {
	return r.Cmd.InstallDesktopApp(constants.Raycast)
}

func (r *Raycast) ForceInstall() error {
	if err := r.Uninstall(); err != nil {
		return fmt.Errorf("raycast.ForceInstall: uninstall failed: %w", err)
	}
	if err := r.Install(); err != nil {
		return fmt.Errorf("raycast.ForceInstall: install failed: %w", err)
	}
	return nil
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
	return fmt.Errorf("uninstall not supported for Raycast")
}

func (r *Raycast) ExecuteCommand(args ...string) error {
	// No CLI commands for desktop GUI application
	return nil
}

func (r *Raycast) Update() error {
	return fmt.Errorf("update not supported for Raycast")
}
