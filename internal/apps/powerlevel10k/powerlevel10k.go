package powerlevel10k

import (
	"fmt"
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

type PowerLevel10k struct {
	Cmd  cmd.Command
	Base cmd.BaseCommand
}

func New() *PowerLevel10k {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &PowerLevel10k{Cmd: osCmd, Base: *baseCmd}
}

func (p *PowerLevel10k) Install() error {
	return p.Cmd.InstallPackage("powerlevel10k")
}

func (p *PowerLevel10k) MaybeInstall() error {
	return p.Cmd.MaybeInstallPackage("powerlevel10k")
}

func (p *PowerLevel10k) Setup() error {
	return files.AddLineToFile(
		"source $(brew --prefix)/share/powerlevel10k/powerlevel10k.zsh-theme",
		filepath.Join(paths.AppDir, "devgita.zsh"),
	)
}

func (p *PowerLevel10k) MaybeSetup() error {
	isConfigured, err := files.ContentExistsInFile(
		filepath.Join(paths.AppDir, "devgita.zsh"),
		"powerlevel10k.zsh-theme",
	)
	if err != nil {
		return err
	}
	if isConfigured == true {
		return nil
	}
	return p.Setup()
}

func (p *PowerLevel10k) Reconfigure() error {
	return p.Run("configure")
}

func (p *PowerLevel10k) Run(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		Verbose:     true,
		IsSudo:      false,
		Command:     "p10k",
		Args:        args,
	}
	if _, err := p.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run powerlevel10k command: %w", err)
	}
	return nil
}
