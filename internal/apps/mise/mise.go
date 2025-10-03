// -------------------------
// TODO: Compete documentation: https://mise.jdx.dev/
// -------------------------

package mise

import (
	"fmt"

	cmd "github.com/cjairm/devgita/internal/commands"
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

func (m *Mise) MaybeInstall() error {
	return m.Cmd.MaybeInstallPackage("mise")
}

func (m *Mise) Setup() error {
	return fmt.Errorf("No configuration file")
}

func (m *Mise) MaybeSetup() error {
	return fmt.Errorf("No configuration file")
}

func (m *Mise) UseGlobal(language, version string) error {
	if language == "" {
		return fmt.Errorf("`language` is required")
	}
	if version == "" {
		return fmt.Errorf("`version` is required")
	}
	return m.Run("use", "--global", language+"@"+version)
}

func (m *Mise) Run(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		IsSudo:      false,
		Command:     "mise",
		Args:        args,
	}
	if _, _, err := m.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run mise command: %w", err)
	}
	return nil
}
