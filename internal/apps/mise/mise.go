// -------------------------
// TODO: Compete documentation: https://mise.jdx.dev/
// -------------------------

package mise

import (
	"fmt"

	cmd "github.com/cjairm/devgita/internal/commands"
)

type Mise struct {
	Cmd cmd.Command
}

func New() *Mise {
	osCmd := cmd.NewCommand()
	return &Mise{Cmd: osCmd}
}

func Command(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		Verbose:     true,
		IsSudo:      false,
		Command:     "mise",
		Args:        args,
	}
	return cmd.ExecCommand(execCommand)
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
	return Command("use", "--global", language+"@"+version)
}
