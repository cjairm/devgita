package devgita

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cjairm/devgita/internal/apps/git"
	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/logger"
	"github.com/cjairm/devgita/pkg/paths"
)

type Devgita struct {
	Git  git.Git
	Base commands.BaseCommandExecutor
}

var configDirPath = filepath.Join(paths.Paths.Config.Root, constants.App.Name)
var globalConfigPath = filepath.Join(configDirPath, constants.App.File.GlobalConfig)
var zshConfigPath = filepath.Join(configDirPath, constants.App.File.ZshConfig)

func New() *Devgita {
	return &Devgita{Git: *git.New(), Base: commands.NewBaseCommand()}
}

func (dg *Devgita) Install() error {
	if err := os.MkdirAll(paths.Paths.App.Root, 0755); err != nil {
		return err
	}
	if err := os.RemoveAll(paths.Paths.App.Root); err != nil {
		return err
	}
	if err := dg.Git.Clone(constants.App.Repository.URL, paths.Paths.App.Root); err != nil {
		return err
	}
	return nil
}

func (dg *Devgita) SoftInstall() error {
	if files.DirAlreadyExist(paths.Paths.App.Root) && !files.IsDirEmpty(paths.Paths.App.Root) {
		logger.L().Info("Devgita repository already installed at %s", paths.Paths.App.Root)
		return nil
	}
	return dg.Install()
}

func (dg *Devgita) ForceInstall() error {
	err := dg.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall devgita: %w", err)
	}
	return dg.Install()
}

func (dg *Devgita) Uninstall() error {
	// TODO: main .zsh file should be cleaned up, too. (Remove source line)
	if !files.DirAlreadyExist(paths.Paths.App.Root) {
		logger.L().
			Debug("Devgita repository not found at %s, nothing to uninstall", paths.Paths.App.Root)
	} else if files.IsDirEmpty(paths.Paths.App.Root) {
		logger.L().
			Debug("Devgita repository directory is empty at %s, removing directory", paths.Paths.App.Root)
		if err := os.Remove(paths.Paths.App.Root); err != nil {
			return fmt.Errorf("failed to remove empty devgita directory: %w", err)
		}
	} else {
		if err := os.RemoveAll(paths.Paths.App.Root); err != nil {
			return fmt.Errorf("failed to uninstall devgita repository: %w", err)
		}
	}
	if files.FileAlreadyExist(globalConfigPath) {
		logger.L().Debug("Removing global config file at %s", globalConfigPath)
		if err := os.Remove(globalConfigPath); err != nil {
			return fmt.Errorf("failed to remove global config file: %w", err)
		}
	} else {
		logger.L().Debug("Global config file not found at %s", globalConfigPath)
	}
	if files.FileAlreadyExist(zshConfigPath) {
		logger.L().Debug("Removing zsh config file at %s", zshConfigPath)
		if err := os.Remove(zshConfigPath); err != nil {
			return fmt.Errorf("failed to remove zsh config file: %w", err)
		}
	} else {
		logger.L().Debug("zsh file not found at %s", zshConfigPath)
	}
	return nil
}

func (dg *Devgita) ForceConfigure() error {
	gc := &config.GlobalConfig{}
	if err := gc.Create(); err != nil {
		return fmt.Errorf("failed to create global config: %w", err)
	}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	gc.AppPath = paths.Paths.App.Root
	gc.ConfigPath = configDirPath
	if err := gc.Save(); err != nil {
		return fmt.Errorf("failed to save global config: %w", err)
	}
	// NOTE: This should be regenerated per app, but creating it just in case
	if err := gc.RegenerateShellConfig(); err != nil {
		return fmt.Errorf("failed to create global config file: %w", err)
	}
	devgitaConfigLine := fmt.Sprintf("source %s", zshConfigPath)
	if err := dg.Base.MaybeSetup(devgitaConfigLine, zshConfigPath); err != nil {
		return err
	}
	return nil
}

func (dg *Devgita) SoftConfigure() error {
	if files.FileAlreadyExist(globalConfigPath) && files.FileAlreadyExist(zshConfigPath) {
		logger.L().Info("Global config already exists at %s", globalConfigPath)
		return nil
	}
	return dg.ForceConfigure()
}
