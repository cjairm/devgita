// Aerospace window manager with devgita integration
//
// Aerospace is a macOS-specific tiling window manager inspired by i3wm.
// This module provides installation and configuration management for Aerospace
// with devgita integration for seamless macOS window management setup.
//
// References:
// - Aerospace Documentation: https://github.com/nikitabobko/AeroSpace/blob/main/docs/guide.md
// - Aerospace Configuration: https://github.com/nikitabobko/AeroSpace/blob/main/docs/config-examples.md

package aerospace

import (
	"fmt"
	"path/filepath"

	"github.com/cjairm/devgita/internal/apps"
	"github.com/cjairm/devgita/internal/apps/baseapp"
	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

var _ apps.App = (*Aerospace)(nil)

type Aerospace struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func (a *Aerospace) Name() string       { return constants.Aerospace }
func (a *Aerospace) Kind() apps.AppKind { return apps.KindDesktop }

func New() *Aerospace {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Aerospace{Cmd: osCmd, Base: baseCmd}
}

func (a *Aerospace) Install() error {
	return a.Cmd.InstallDesktopApp("nikitabobko/tap/aerospace")
}

func (a *Aerospace) SoftInstall() error {
	return a.Cmd.MaybeInstallDesktopApp("nikitabobko/tap/aerospace", "AeroSpace")
}

func (a *Aerospace) ForceInstall() error {
	return baseapp.Reinstall(a.Install, a.Uninstall)
}

func (a *Aerospace) Uninstall() error {
	return fmt.Errorf("%w for aerospace", apps.ErrUninstallNotSupported)
}

func (a *Aerospace) ForceConfigure() error {
	gc := &config.GlobalConfig{}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	err := files.CopyDir(paths.Paths.App.Configs.Aerospace, paths.Paths.Config.Aerospace)
	if err != nil {
		return fmt.Errorf("failed to copy aerospace starter script: %w", err)
	}
	gc.AddToInstalled(constants.Aerospace, "desktop_app")
	if err := gc.Save(); err != nil {
		return fmt.Errorf("failed to save global config: %w", err)
	}
	return nil
}

func (a *Aerospace) SoftConfigure() error {
	aerospaceConfigFile := filepath.Join(
		paths.Paths.Config.Aerospace,
		fmt.Sprintf("%s.toml", constants.Aerospace),
	)
	if isFilePresent := files.FileAlreadyExist(aerospaceConfigFile); isFilePresent {
		return nil
	}
	return a.ForceConfigure()
}

func (a *Aerospace) ExecuteCommand(args ...string) error {
	// Aerospace is a tiling window manager; no terminal commands managed by devgita
	return nil
}

func (a *Aerospace) Update() error {
	return fmt.Errorf("%w for aerospace", apps.ErrUpdateNotSupported)
}
