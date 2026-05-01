// Fastfetch is a neofetch-like tool for fetching system information and displaying
// it in a visually appealing way. It is written mainly in C, with a focus on performance
// and customizability. Currently, it supports Linux, macOS, Windows 7+, Android,
// FreeBSD, OpenBSD, NetBSD, DragonFly, Haiku, and illumos (SunOS).
//
// Documentation: https://github.com/fastfetch-cli/fastfetch

package fastfetch

import (
	"fmt"
	"path/filepath"

	"github.com/cjairm/devgita/internal/apps"
	"github.com/cjairm/devgita/internal/apps/baseapp"
	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

var _ apps.App = (*Fastfetch)(nil)

type Fastfetch struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func (f *Fastfetch) Name() string       { return constants.Fastfetch }
func (f *Fastfetch) Kind() apps.AppKind { return apps.KindTerminal }

func New() *Fastfetch {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Fastfetch{Cmd: osCmd, Base: baseCmd}
}

func (f *Fastfetch) Install() error {
	return f.Cmd.InstallPackage(constants.Fastfetch)
}

func (f *Fastfetch) ForceInstall() error {
	return baseapp.Reinstall(f.Install, f.Uninstall)
}

func (f *Fastfetch) SoftInstall() error {
	return f.Cmd.MaybeInstallPackage(constants.Fastfetch)
}

func (f *Fastfetch) ForceConfigure() error {
	return files.CopyDir(paths.Paths.App.Configs.Fastfetch, paths.Paths.Config.Fastfetch)
}

func (f *Fastfetch) SoftConfigure() error {
	fastfetchConfigFile := filepath.Join(paths.Paths.Config.Fastfetch, "config.jsonc")
	isFilePresent := files.FileAlreadyExist(fastfetchConfigFile)
	if isFilePresent {
		return nil
	}
	return f.ForceConfigure()
}

func (f *Fastfetch) Uninstall() error {
	return fmt.Errorf("%w for fastfetch", apps.ErrUninstallNotSupported)
}

func (f *Fastfetch) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		IsSudo:  false,
		Command: constants.Fastfetch,
		Args:    args,
	}
	if _, _, err := f.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run fastfetch command: %w", err)
	}
	return nil
}

func (f *Fastfetch) Update() error {
	return fmt.Errorf("%w for fastfetch", apps.ErrUpdateNotSupported)
}
