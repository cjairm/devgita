package devgita

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/cjairm/devgita/internal/apps"
	"github.com/cjairm/devgita/internal/apps/baseapp"
	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/internal/embedded"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/logger"
	"github.com/cjairm/devgita/pkg/paths"
)

var _ apps.App = (*Devgita)(nil)

const DevgitaExtended = "extended_capabilities"

type Devgita struct {
	Base            commands.BaseCommandExecutor
	ExtractEmbedded embedded.ExtractFunc
}

func (dg *Devgita) Name() string       { return constants.DevgitaApp }
func (dg *Devgita) Kind() apps.AppKind { return apps.KindMeta }

func getConfigDirPath() string {
	return filepath.Join(paths.Paths.Config.Root, constants.App.Name)
}

func getGlobalConfigPath() string {
	return filepath.Join(getConfigDirPath(), constants.App.File.GlobalConfig)
}

func getZshConfigPath() string {
	return filepath.Join(paths.Paths.App.Root, fmt.Sprintf("%s.zsh", constants.App.Name))
}

func New() *Devgita {
	return &Devgita{
		Base:            commands.NewBaseCommand(),
		ExtractEmbedded: embedded.DefaultExtractor,
	}
}

func (dg *Devgita) Install() error {
	// Create configs directory inside app root
	configsDir := filepath.Join(paths.Paths.App.Root, "configs")

	// Clean up existing configs directory if it exists
	if files.DirAlreadyExist(configsDir) {
		if err := os.RemoveAll(configsDir); err != nil {
			return fmt.Errorf("failed to remove existing configs directory: %w", err)
		}
	}

	// Extract embedded configs to the app directory
	if err := dg.ExtractEmbedded(configsDir); err != nil {
		return fmt.Errorf("failed to extract embedded configs: %w", err)
	}

	return nil
}

func (dg *Devgita) SoftInstall() error {
	configsDir := filepath.Join(paths.Paths.App.Root, "configs")

	// Check if configs/ subdirectory exists and is non-empty
	if files.DirAlreadyExist(configsDir) && !files.IsDirEmpty(configsDir) {
		logger.L().Infow("Devgita configs already installed", "path", configsDir)
		return nil
	}

	return dg.Install()
}

func (dg *Devgita) ForceInstall() error {
	return baseapp.Reinstall(dg.Install, dg.Uninstall)
}

func (dg *Devgita) Uninstall() error {
	// Clean up extracted configs directory
	configsDir := filepath.Join(paths.Paths.App.Root, "configs")
	if files.DirAlreadyExist(configsDir) {
		logger.L().Debugw("Removing extracted configs directory", "path", configsDir)
		if err := os.RemoveAll(configsDir); err != nil {
			return fmt.Errorf("failed to remove configs directory: %w", err)
		}
	} else {
		logger.L().Debugw("Configs directory not found", "path", configsDir)
	}

	// Remove app root directory if empty
	if files.DirAlreadyExist(paths.Paths.App.Root) {
		if files.IsDirEmpty(paths.Paths.App.Root) {
			logger.L().Debugw("App directory is empty, removing", "path", paths.Paths.App.Root)
			if err := os.Remove(paths.Paths.App.Root); err != nil {
				return fmt.Errorf("failed to remove empty app directory: %w", err)
			}
		}
	}

	// Remove global config file
	if files.FileAlreadyExist(getGlobalConfigPath()) {
		logger.L().Debugw("Removing global config file", "path", getGlobalConfigPath())
		if err := os.Remove(getGlobalConfigPath()); err != nil {
			return fmt.Errorf("failed to remove global config file: %w", err)
		}
	} else {
		logger.L().Debugw("Global config file not found", "path", getGlobalConfigPath())
	}

	// Remove zsh config file
	if files.FileAlreadyExist(getZshConfigPath()) {
		logger.L().Debugw("Removing zsh config file", "path", getZshConfigPath())
		if err := os.Remove(getZshConfigPath()); err != nil {
			return fmt.Errorf("failed to remove zsh config file: %w", err)
		}
	} else {
		logger.L().Debugw("zsh file not found", "path", getZshConfigPath())
	}

	// Remove config directory if empty
	if files.DirAlreadyExist(getConfigDirPath()) && files.IsDirEmpty(getConfigDirPath()) {
		logger.L().Debugw("Config directory is empty, removing", "path", getConfigDirPath())
		if err := os.Remove(getConfigDirPath()); err != nil {
			return fmt.Errorf("failed to remove empty config directory: %w", err)
		}
	}

	return nil
}

func (dg *Devgita) ForceConfigure() error {
	// Re-extract embedded configs so deployed configs/claude, configs/opencode, etc.
	// always match the current binary (not the original install).
	if err := dg.Install(); err != nil {
		return fmt.Errorf("failed to refresh embedded configs: %w", err)
	}

	gc := &config.GlobalConfig{}
	if err := gc.Create(); err != nil {
		return fmt.Errorf("failed to create global config: %w", err)
	}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	gc.AppPath = paths.Paths.App.Root
	gc.ConfigPath = getConfigDirPath()
	// Set platform flag for shell template conditionals
	gc.Shell.IsMac = runtime.GOOS == "darwin"
	gc.EnableShellFeature(DevgitaExtended)
	if err := gc.Save(); err != nil {
		return fmt.Errorf("failed to save global config: %w", err)
	}
	// NOTE: This should be regenerated per app, but creating it just in case
	if err := gc.RegenerateShellConfig(); err != nil {
		return fmt.Errorf("failed to create global config file: %w", err)
	}
	devgitaConfigLine := fmt.Sprintf("source %s", getZshConfigPath())
	if err := dg.Base.MaybeSetup(devgitaConfigLine, getZshConfigPath()); err != nil {
		return err
	}
	return nil
}

func (dg *Devgita) SoftConfigure() error {
	if !files.FileAlreadyExist(getGlobalConfigPath()) ||
		!files.FileAlreadyExist(getZshConfigPath()) {
		return dg.ForceConfigure()
	}
	gc := &config.GlobalConfig{}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	if !gc.IsShellFeatureEnabled(DevgitaExtended) {
		if err := enableExtendedCapabilities(gc); err != nil {
			return fmt.Errorf("failed to enable extended capabilities: %w", err)
		}
	}
	return nil
}

func (dg *Devgita) ExecuteCommand(_ ...string) error {
	return fmt.Errorf("%w for devgita", apps.ErrExecuteNotSupported)
}

func (dg *Devgita) Update() error {
	return fmt.Errorf("%w for devgita", apps.ErrUpdateNotSupported)
}

func enableExtendedCapabilities(gc *config.GlobalConfig) error {
	gc.EnableShellFeature(DevgitaExtended)
	if err := gc.RegenerateShellConfig(); err != nil {
		return fmt.Errorf("failed to generate shell config: %w", err)
	}
	if err := gc.Save(); err != nil {
		return fmt.Errorf("failed to save global config: %w", err)
	}
	return nil
}
