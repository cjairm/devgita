// Package alacritty provides installation and configuration management for Alacritty terminal emulator.
// Alacritty is a fast, cross-platform terminal emulator written in Rust that uses GPU acceleration.
// This module follows the standardized devgita app interface for consistent lifecycle management.

package alacritty

import (
	"fmt"
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

type Alacritty struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

type ConfigureOptions struct {
	Font  string
	Theme string
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

func (a *Alacritty) ForceConfigure(options ConfigureOptions) error {
	gc := &config.GlobalConfig{}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	font := "default"
	if options.Font != "" {
		font = options.Font
	}
	theme := "default"
	if options.Theme != "" {
		theme = options.Theme
	}
	configFilePath := filepath.Join(
		paths.Paths.Config.Alacritty,
		fmt.Sprintf("%s.toml", constants.Alacritty),
	)
	tmplPath := filepath.Join(
		paths.Paths.App.Configs.Alacritty,
		fmt.Sprintf("%s.toml.tmpl", constants.Alacritty),
	)
	if err := files.GenerateFromTemplate(tmplPath, configFilePath, map[string]string{
		"Font":       font,
		"Theme":      theme,
		"ConfigPath": paths.Paths.Config.Root,
	}); err != nil {
		return fmt.Errorf("failed to generate alacritty configuration: %w", err)
	}
	if err := files.CopyFile(
		filepath.Join(paths.Paths.App.Configs.Alacritty, "starter.sh"),
		filepath.Join(paths.Paths.Config.Alacritty, "starter.sh"),
	); err != nil {
		return fmt.Errorf("failed to copy alacritty starter script: %w", err)
	}
	gc.AddToInstalled(constants.Alacritty, "desktop_app")
	if err := gc.Save(); err != nil {
		return fmt.Errorf("failed to save global config: %w", err)
	}
	return nil
}

func (a *Alacritty) SoftConfigure(options ConfigureOptions) error {
	gc := &config.GlobalConfig{}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	if gc.IsAlreadyInstalled(constants.Alacritty, "desktop_app") ||
		gc.IsInstalledByDevgita(constants.Alacritty, "desktop_app") {
		return nil
	}
	if err := a.ForceConfigure(options); err != nil {
		return fmt.Errorf("failed to configure alacritty: %w", err)
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
