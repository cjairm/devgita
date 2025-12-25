// Package flameshot provides installation and management for Flameshot screenshot tool.
//
// Flameshot is a powerful yet simple to use screenshot software. This module follows
// the standardized devgita app interface for desktop applications.

package flameshot

import (
	"fmt"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
)

type Flameshot struct {
	Cmd cmd.Command
}

func New() *Flameshot {
	return &Flameshot{Cmd: cmd.NewCommand()}
}

func (f *Flameshot) Install() error {
	return f.Cmd.InstallDesktopApp(constants.Flameshot)
}

func (f *Flameshot) ForceInstall() error {
	if err := f.Uninstall(); err != nil {
		return fmt.Errorf("flameshot.ForceInstall: uninstall failed: %w", err)
	}
	if err := f.Install(); err != nil {
		return fmt.Errorf("flameshot.ForceInstall: install failed: %w", err)
	}
	return nil
}

func (f *Flameshot) SoftInstall() error {
	return f.Cmd.MaybeInstallDesktopApp(constants.Flameshot)
}

func (f *Flameshot) ForceConfigure() error {
	// No configuration needed for GUI-based screenshot tool
	return nil
}

func (f *Flameshot) SoftConfigure() error {
	// No configuration needed for GUI-based screenshot tool
	return nil
}

func (f *Flameshot) Uninstall() error {
	return fmt.Errorf("uninstall not supported for Flameshot")
}

func (f *Flameshot) ExecuteCommand(args ...string) error {
	// No CLI commands for desktop GUI application
	return nil
}

func (f *Flameshot) Update() error {
	return fmt.Errorf("update not supported for Flameshot")
}
