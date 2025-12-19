// Package alacritty provides installation and configuration management for Alacritty terminal emulator.
// Alacritty is a fast, cross-platform terminal emulator written in Rust that uses GPU acceleration.
// This module follows the standardized devgita app interface for consistent lifecycle management.

package alacritty

import (
	"fmt"
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

type Alacritty struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func New() *Alacritty {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Alacritty{Cmd: osCmd, Base: baseCmd}
}

func (a *Alacritty) Install() error {
	return a.Cmd.InstallDesktopApp(constants.Alacritty)
}

func (a *Alacritty) SoftInstall() error {
	return a.Cmd.MaybeInstallDesktopApp(constants.Alacritty)
}

func (a *Alacritty) ForceInstall() error {
	err := a.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall alacritty: %w", err)
	}
	return a.Install()
}

func (a *Alacritty) ForceConfigureApp() error {
	return files.CopyDir(paths.Paths.App.Configs.Alacritty, paths.Paths.Config.Alacritty)
}

func (a *Alacritty) ForceConfigureFont() error {
	return files.CopyDir(
		filepath.Join(paths.Paths.App.Fonts.Alacritty, "default"),
		paths.Paths.Config.Alacritty,
	)
}

func (a *Alacritty) ForceConfigureTheme() error {
	return files.CopyDir(
		filepath.Join(paths.Paths.App.Themes.Alacritty, "default"),
		paths.Paths.Config.Alacritty,
	)
}

func (a *Alacritty) SoftConfigureApp() error {
	return maybeConfigure(a.ForceConfigureApp, paths.Paths.Config.Alacritty, "alacritty.toml")
}

func (a *Alacritty) SoftConfigureFont() error {
	return maybeConfigure(a.ForceConfigureFont, paths.Paths.Config.Alacritty, "font.toml")
}

func (a *Alacritty) SoftConfigureTheme() error {
	return maybeConfigure(a.ForceConfigureTheme, paths.Paths.Config.Alacritty, "theme.toml")
}

func (a *Alacritty) ForceConfigure() error {
	err := a.ForceConfigureApp()
	if err != nil {
		return fmt.Errorf("failed to setup app configuration: %w", err)
	}
	err = a.ForceConfigureFont()
	if err != nil {
		return fmt.Errorf("failed to setup font configuration: %w", err)
	}
	err = a.ForceConfigureTheme()
	if err != nil {
		return fmt.Errorf("failed to setup theme configuration: %w", err)
	}
	err = a.UpdateConfigFilesWithCurrentHomeDir()
	if err != nil {
		return fmt.Errorf("failed to update config files with home directory: %w", err)
	}
	return nil
}

func (a *Alacritty) SoftConfigure() error {
	err := a.SoftConfigureApp()
	if err != nil {
		return fmt.Errorf("failed to setup app configuration: %w", err)
	}
	err = a.SoftConfigureFont()
	if err != nil {
		return fmt.Errorf("failed to setup font configuration: %w", err)
	}
	err = a.SoftConfigureTheme()
	if err != nil {
		return fmt.Errorf("failed to setup theme configuration: %w", err)
	}
	err = a.UpdateConfigFilesWithCurrentHomeDir()
	if err != nil {
		return fmt.Errorf("failed to update config files with home directory: %w", err)
	}
	return nil
}

func (a *Alacritty) Uninstall() error {
	return fmt.Errorf("uninstall not implemented for alacritty")
}

func (a *Alacritty) ExecuteCommand(args ...string) error {
	// No alacritty commands in terminal
	return nil
}

func (a *Alacritty) Update() error {
	return fmt.Errorf("update not implemented for alacritty")
}

func (a *Alacritty) UpdateConfigFilesWithCurrentHomeDir() error {
	alacrittyConfigFile := filepath.Join(paths.Paths.Config.Alacritty, "alacritty.toml")
	return files.UpdateFile(alacrittyConfigFile, "<ALACRITTY-CONFIG-PATH>", paths.Paths.Config.Root)
}

func maybeConfigure(setupFunc func() error, localConfig string, fileSegments ...string) error {
	filePath := localConfig
	for _, segment := range fileSegments {
		filePath = filepath.Join(filePath, segment)
	}
	if isFilePresent := files.FileAlreadyExist(filePath); isFilePresent {
		return nil
	}
	return setupFunc()
}
