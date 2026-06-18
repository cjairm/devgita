// jq JSON processor with devgita integration
//
// jq is a lightweight and flexible command-line JSON processor. It lets you
// slice, filter, map, and transform structured JSON data with ease.
//
// References:
// - jq Documentation: https://jqlang.github.io/jq/
// - jq Manual: https://jqlang.github.io/jq/manual/
//
// Common jq commands available through ExecuteCommand():
//   - jq '.' file.json          - Pretty-print JSON
//   - jq '.key' file.json       - Extract a field
//   - jq '.[] | .field' file.json - Iterate array, extract field
//   - jq -r '.key' file.json    - Raw string output (no quotes)
//   - jq -c '.' file.json       - Compact output

package jq

import (
	"fmt"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
)

type Jq struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func New() *Jq {
	return &Jq{
		Cmd:  cmd.NewCommand(),
		Base: cmd.NewBaseCommand(),
	}
}

func (j *Jq) Install() error {
	return j.Cmd.InstallPackage(constants.Jq)
}

func (j *Jq) SoftInstall() error {
	return j.Cmd.MaybeInstallPackage(constants.Jq)
}

func (j *Jq) ForceInstall() error {
	if err := j.Uninstall(); err != nil {
		return fmt.Errorf("failed to uninstall jq: %w", err)
	}
	return j.Install()
}

func (j *Jq) Uninstall() error {
	return fmt.Errorf("jq uninstall not supported through devgita")
}

func (j *Jq) ForceConfigure() error {
	return nil
}

func (j *Jq) SoftConfigure() error {
	return nil
}

func (j *Jq) ExecuteCommand(args ...string) error {
	_, _, err := j.Base.ExecCommand(cmd.CommandParams{
		IsSudo:  false,
		Command: constants.Jq,
		Args:    args,
	})
	if err != nil {
		return fmt.Errorf("failed to run jq command: %w", err)
	}
	return nil
}

func (j *Jq) Update() error {
	return fmt.Errorf("jq update not implemented through devgita")
}
