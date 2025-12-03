// Package fzf provides installation and command execution management for fzf
// fuzzy finder with devgita integration. It follows the standardized devgita
// app interface while providing fzf-specific operations for interactive file
// searching, command history filtering, and directory navigation.
package fzf

import (
	"fmt"

	cmd "github.com/cjairm/devgita/internal/commands"
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
	// Placeholder: Fzf uses environment variables for configuration
	// Shell integration is typically added to shell rc files
	// No traditional config file copying required
	return nil
}

func (f *Fzf) SoftConfigure() error {
	// Placeholder: Check if fzf shell integration is already configured
	// Typically checked via presence of FZF_* environment variables
	// or fzf key bindings in shell configuration
	return nil
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
