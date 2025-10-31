// Fastfetch module provides installation and configuration management for fastfetch with devgita integration.
// Fastfetch is a neofetch-like tool for fetching system information and displaying it prettily.
// Documentation: https://github.com/fastfetch-cli/fastfetch

package fastfetch

import (
	"fmt"
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

type Fastfetch struct {
	Cmd  cmd.Command
	Base cmd.BaseCommand
}

func New() *Fastfetch {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Fastfetch{Cmd: osCmd, Base: *baseCmd}
}

func (f *Fastfetch) Install() error {
	return f.Cmd.InstallPackage("fastfetch")
}

func (f *Fastfetch) ForceInstall() error {
	err := f.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall fastfetch: %w", err)
	}
	return f.Install()
}

func (f *Fastfetch) SoftInstall() error {
	return f.Cmd.MaybeInstallPackage("fastfetch")
}

func (f *Fastfetch) ForceConfigure() error {
	return files.CopyDir(paths.FastFetchConfigAppDir, paths.FastFetchConfigLocalDir)
}

func (f *Fastfetch) SoftConfigure() error {
	fastfetchConfigFile := filepath.Join(paths.FastFetchConfigLocalDir, "config.jsonc")
	isFilePresent := files.FileAlreadyExist(fastfetchConfigFile)
	if isFilePresent {
		return nil
	}
	return f.ForceConfigure()
}

func (f *Fastfetch) Uninstall() error {
	return fmt.Errorf("uninstall not implemented for fastfetch")
}

func (f *Fastfetch) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		IsSudo:      false,
		Command:     constants.Fastfetch,
		Args:        args,
	}
	if _, _, err := f.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run fastfetch command: %w", err)
	}
	return nil
}

func (f *Fastfetch) Update() error {
	return fmt.Errorf("update not implemented for fastfetch")
}
