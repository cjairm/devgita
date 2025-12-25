// Package brave provides installation and management for Brave browser desktop application.
//
// Brave is a privacy-focused web browser built on Chromium. This module follows the
// standardized devgita app interface for desktop applications.

package brave

import (
	"fmt"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
)

type Brave struct {
	Cmd cmd.Command
}

func New() *Brave {
	return &Brave{Cmd: cmd.NewCommand()}
}

func (b *Brave) Install() error {
	return b.Cmd.InstallDesktopApp(fmt.Sprintf("%s-browser", constants.Brave))
}

func (b *Brave) ForceInstall() error {
	if err := b.Uninstall(); err != nil {
		return fmt.Errorf("brave.ForceInstall: uninstall failed: %w", err)
	}
	if err := b.Install(); err != nil {
		return fmt.Errorf("brave.ForceInstall: install failed: %w", err)
	}
	return nil
}

func (b *Brave) SoftInstall() error {
	return b.Cmd.MaybeInstallDesktopApp(
		fmt.Sprintf("%s-browser", constants.Brave),
		constants.Brave,
	)
}

func (b *Brave) ForceConfigure() error {
	// No configuration needed for GUI-based browser
	return nil
}

func (b *Brave) SoftConfigure() error {
	// No configuration needed for GUI-based browser
	return nil
}

func (b *Brave) Uninstall() error {
	return fmt.Errorf("uninstall not supported for Brave")
}

func (b *Brave) ExecuteCommand(args ...string) error {
	// No CLI commands for desktop GUI application
	return nil
}

func (b *Brave) Update() error {
	return fmt.Errorf("update not supported for Brave")
}
