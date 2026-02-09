// Package mise provides runtime environment management via mise.jdx.dev
//
// Mise is a polyglot runtime manager that replaces tools like nvm, rbenv, pyenv, etc.
// It manages language versions globally and per-project, supporting Node.js, Python, Go, Ruby, and more.
//
// Reference: https://mise.jdx.dev/

package mise

import (
	"fmt"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/constants"
)

type Mise struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func New() *Mise {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Mise{Cmd: osCmd, Base: baseCmd}
}

func (m *Mise) Install() error {
	return m.Cmd.InstallPackage(constants.Mise)
}

func (m *Mise) ForceInstall() error {
	err := m.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall mise before force install: %w", err)
	}
	return m.Install()
}

func (m *Mise) SoftInstall() error {
	return m.Cmd.MaybeInstallPackage(constants.Mise)
}

// ForceConfigure enables mise shell integration and regenerates shell config
func (m *Mise) ForceConfigure() error {
	gc := &config.GlobalConfig{}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	gc.EnableShellFeature(constants.Mise)
	if err := gc.RegenerateShellConfig(); err != nil {
		return fmt.Errorf("failed to generate shell config: %w", err)
	}
	if err := gc.Save(); err != nil {
		return fmt.Errorf("failed to save global config: %w", err)
	}
	return nil
}

func (m *Mise) SoftConfigure() error {
	gc := &config.GlobalConfig{}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	if gc.IsShellFeatureEnabled(constants.Mise) {
		return nil
	}
	return m.ForceConfigure()
}

func (m *Mise) Uninstall() error {
	gc := &config.GlobalConfig{}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	gc.DisableShellFeature(constants.Mise)
	if err := gc.RegenerateShellConfig(); err != nil {
		return fmt.Errorf("failed to generate shell config: %w", err)
	}
	// TODO: We still uninstall the app or remove downloaded doc - see `Install`
	return gc.Save()
}

func (m *Mise) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		Command: constants.Mise,
		Args:    args,
	}
	if _, _, err := m.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run mise command: %w", err)
	}
	return nil
}

func (m *Mise) Update() error {
	return fmt.Errorf("update not implemented for mise")
}

func (m *Mise) UseGlobal(language, version string) error {
	if language == "" {
		return fmt.Errorf("`language` is required")
	}
	if version == "" {
		return fmt.Errorf("`version` is required")
	}
	return m.ExecuteCommand("use", "--global", fmt.Sprintf("%s@%s", language, version))
}
