package i3

import (
	"fmt"
	"path/filepath"

	"github.com/cjairm/devgita/internal/apps"
	"github.com/cjairm/devgita/internal/apps/baseapp"
	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/logger"
	"github.com/cjairm/devgita/pkg/paths"
)

var _ apps.App = (*I3)(nil)

type I3 struct {
	Cmd commands.Command
}

func (i *I3) Name() string       { return constants.I3 }
func (i *I3) Kind() apps.AppKind { return apps.KindDesktop }

func New() *I3 {
	return &I3{Cmd: commands.NewCommand()}
}

func (i *I3) Install() error {
	return i.Cmd.InstallPackage(constants.I3)
}

func (i *I3) SoftInstall() error {
	return i.Cmd.MaybeInstallPackage(constants.I3)
}

func (i *I3) ForceInstall() error {
	return baseapp.Reinstall(i.Install, i.Uninstall)
}

func (i *I3) Uninstall() error {
	return fmt.Errorf("%w for i3 — manage via system package manager", apps.ErrUninstallNotSupported)
}

func (i *I3) ForceConfigure() error {
	// Copy i3 config from app configs to local i3 config directory
	if err := files.CopyDir(paths.Paths.App.Configs.I3, paths.Paths.Config.I3); err != nil {
		return fmt.Errorf("failed to copy i3 config: %w", err)
	}
	logger.L().Infow("i3 configuration applied", "source", paths.Paths.App.Configs.I3, "dest", paths.Paths.Config.I3)
	return nil
}

func (i *I3) SoftConfigure() error {
	// Check for marker file (config) in i3 config directory
	markerFile := filepath.Join(paths.Paths.Config.I3, "config")
	if files.FileAlreadyExist(markerFile) {
		logger.L().Infow("i3 config already exists", "path", markerFile)
		return nil
	}
	return i.ForceConfigure()
}

func (i *I3) ExecuteCommand(args ...string) error {
	// i3 commands could be useful for window management automation
	// For now, return nil for interface compliance
	return nil
}

func (i *I3) Update() error {
	return fmt.Errorf("%w for i3 — use system package manager", apps.ErrUpdateNotSupported)
}
