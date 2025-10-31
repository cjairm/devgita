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
	"github.com/cjairm/devgita/pkg/constants"
)

type Mise struct {
	Cmd  cmd.Command
	Base cmd.BaseCommand
}

func New() *Mise {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Mise{Cmd: osCmd, Base: *baseCmd}
}

func (m *Mise) Install() error {
	return m.Cmd.InstallPackage("mise")
}

func (m *Mise) ForceInstall() error {
	err := m.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall mise before force install: %w", err)
	}
	return m.Install()
}

func (m *Mise) SoftInstall() error {
	return m.Cmd.MaybeInstallPackage("mise")
}

func (m *Mise) ForceConfigure() error {
	return fmt.Errorf(
		"mise configuration not implemented - runtime management handled via UseGlobal",
	)
}

func (m *Mise) SoftConfigure() error {
	return fmt.Errorf(
		"mise configuration not implemented - runtime management handled via UseGlobal",
	)
}

func (m *Mise) Uninstall() error {
	return fmt.Errorf("uninstall not implemented for mise")
}

func (m *Mise) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		IsSudo:      false,
		Command:     constants.Mise,
		Args:        args,
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
	return m.ExecuteCommand("use", "--global", language+"@"+version)
}
