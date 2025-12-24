// Package gimp provides GIMP (GNU Image Manipulation Program) installation and management
// with devgita integration.
//
// GIMP is a free and open-source raster graphics editor used for image manipulation,
// image editing, free-form drawing, transcoding between different image file formats,
// and more specialized tasks. This module ensures GIMP is properly installed across
// macOS (Homebrew cask) and Debian/Ubuntu (apt) systems.

package gimp

import (
	"fmt"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
)

type Gimp struct {
	Cmd cmd.Command
}

func New() *Gimp {
	return &Gimp{Cmd: cmd.NewCommand()}
}

func (g *Gimp) Install() error {
	return g.Cmd.InstallDesktopApp(constants.Gimp)
}

func (g *Gimp) SoftInstall() error {
	return g.Cmd.MaybeInstallDesktopApp(constants.Gimp)
}

func (g *Gimp) ForceInstall() error {
	err := g.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall GIMP: %w", err)
	}
	return g.Install()
}

func (g *Gimp) ForceConfigure() error {
	// Desktop applications don't require configuration file management
	return nil
}

func (g *Gimp) SoftConfigure() error {
	// Desktop applications don't require configuration file management
	return nil
}

func (g *Gimp) Uninstall() error {
	return fmt.Errorf("uninstall not supported for GIMP")
}

func (g *Gimp) ExecuteCommand(args ...string) error {
	// Desktop applications don't have CLI commands to execute
	return nil
}

func (g *Gimp) Update() error {
	return fmt.Errorf("update not implemented for GIMP")
}
