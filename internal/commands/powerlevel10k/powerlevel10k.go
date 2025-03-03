package powerlevel10k

import (
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal"
	commands "github.com/cjairm/devgita/internal"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/files"
)

type PowerLevel10k struct {
	Cmd cmd.Command
}

func New() *PowerLevel10k {
	osCmd := cmd.NewCommand()
	return &PowerLevel10k{Cmd: osCmd}
}

func Command(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		Verbose:     true,
		IsSudo:      false,
		Command:     "p10k",
		Args:        args,
	}
	return cmd.ExecCommand(execCommand)
}

func (p *PowerLevel10k) Install() error {
	return p.Cmd.InstallPackage("powerlevel10k")
}

func (p *PowerLevel10k) MaybeInstall() error {
	return p.Cmd.MaybeInstallPackage("powerlevel10k")
}

func (p *PowerLevel10k) Setup() error {
	devgitaCustomDir, err := config.GetDevgitaConfigDir()
	if err != nil {
		return err
	}
	err = commands.AddLineToFile(
		"source $(brew --prefix)/share/powerlevel10k/powerlevel10k.zsh-theme",
		devgitaCustomDir+"/devgita.zsh",
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *PowerLevel10k) MaybeSetup() error {
	devgitaCustomDir, err := config.GetDevgitaConfigDir()
	if err != nil {
		return err
	}
	devgitaConfigFile := filepath.Join(devgitaCustomDir, "devgita.zsh")
	isConfigured, err := files.ContentExistsInFile(
		devgitaConfigFile,
		"powerlevel10k.zsh-theme",
	)
	if isConfigured == true {
		return nil
	}
	return p.Setup()
}

func (p *PowerLevel10k) Reconfigure() error {
	return Command("configure")
}
