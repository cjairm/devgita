package powerlevel10k

import (
	"fmt"
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

type PowerLevel10k struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func New() *PowerLevel10k {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &PowerLevel10k{Cmd: osCmd, Base: baseCmd}
}

func (p *PowerLevel10k) Install() error {
	return p.Cmd.InstallPackage(constants.Powerlevel10k)
}

func (p *PowerLevel10k) ForceInstall() error {
	if err := p.Uninstall(); err != nil {
		return fmt.Errorf("failed to uninstall powerlevel10k before force install: %w", err)
	}
	return p.Install()
}

func (p *PowerLevel10k) SoftInstall() error {
	return p.Cmd.MaybeInstallPackage(constants.Powerlevel10k)
}

func (p *PowerLevel10k) ForceConfigure() error {
	return files.AddLineToFile(
		"source $(brew --prefix)/share/powerlevel10k/powerlevel10k.zsh-theme",
		filepath.Join(paths.AppDir, "devgita.zsh"),
	)
}

func (p *PowerLevel10k) SoftConfigure() error {
	isConfigured, err := files.ContentExistsInFile(
		filepath.Join(paths.AppDir, "devgita.zsh"),
		"powerlevel10k.zsh-theme",
	)
	if err != nil {
		return err
	}
	if isConfigured {
		return nil
	}
	return p.ForceConfigure()
}

func (p *PowerLevel10k) Uninstall() error {
	return fmt.Errorf("uninstall not supported for %s", constants.Powerlevel10k)
}

func (p *PowerLevel10k) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		IsSudo:      false,
		Command:     "p10k",
		Args:        args,
	}
	if _, _, err := p.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run powerlevel10k command: %w", err)
	}
	return nil
}

func (p *PowerLevel10k) Update() error {
	return fmt.Errorf("update not implemented for %s", constants.Powerlevel10k)
}

func (p *PowerLevel10k) Reconfigure() error {
	return p.ExecuteCommand("configure")
}
