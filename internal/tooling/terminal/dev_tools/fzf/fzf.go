// Package fzf provides installation and command execution management for fzf
// fuzzy finder with devgita integration. It follows the standardized devgita
// app interface while providing fzf-specific operations for interactive file
// searching, command history filtering, and directory navigation.
package fzf

import (
	"fmt"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/constants"
)

type Fzf struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func New() *Fzf {
	return &Fzf{
		Cmd:  cmd.NewCommand(),
		Base: cmd.NewBaseCommand(),
	}
}

func (f *Fzf) Install() error {
	return f.Cmd.InstallPackage(constants.Fzf)
}

func (f *Fzf) ForceInstall() error {
	if err := f.Uninstall(); err != nil {
		return err
	}
	return f.Install()
}

func (f *Fzf) SoftInstall() error {
	return f.Cmd.MaybeInstallPackage(constants.Fzf)
}

func (f *Fzf) ForceConfigure() error {
	gc := &config.GlobalConfig{}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	gc.EnableShellFeature(constants.Fzf)
	if err := gc.RegenerateShellConfig(); err != nil {
		return fmt.Errorf("failed to generate shell config: %w", err)
	}
	if err := gc.Save(); err != nil {
		return fmt.Errorf("failed to save global config: %w", err)
	}
	return nil
}

func (f *Fzf) SoftConfigure() error {
	gc := &config.GlobalConfig{}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	if gc.IsShellFeatureEnabled(constants.Fzf) {
		return nil
	}
	return f.ForceConfigure()
}

func (f *Fzf) Uninstall() error {
	return fmt.Errorf("fzf uninstall not supported through devgita")
}

func (f *Fzf) ExecuteCommand(args ...string) error {
	_, _, err := f.Base.ExecCommand(cmd.CommandParams{
		Command: constants.Fzf,
		Args:    args,
		IsSudo:  false,
	})
	if err != nil {
		return fmt.Errorf("failed to run fzf command: %w", err)
	}
	return nil
}

func (f *Fzf) Update() error {
	return fmt.Errorf("fzf update not implemented - use system package manager")
}
