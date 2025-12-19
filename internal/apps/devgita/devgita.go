package devgita

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cjairm/devgita/internal/apps/git"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/logger"
	"github.com/cjairm/devgita/pkg/paths"
)

type Devgita struct {
	git git.Git
}

var configDirPath = filepath.Join(paths.Paths.Config.Root, constants.App.Name)
var globalConfigPath = filepath.Join(configDirPath, constants.App.File.GlobalConfig)

func New() *Devgita {
	return &Devgita{git: *git.New()}
}

func (dg *Devgita) Install() error {
	if err := os.MkdirAll(paths.Paths.App.Root, 0755); err != nil {
		return err
	}
	if err := os.RemoveAll(paths.Paths.App.Root); err != nil {
		return err
	}
	if err := dg.git.Clone(constants.App.Repository.URL, paths.Paths.App.Root); err != nil {
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
	if !files.DirAlreadyExist(paths.Paths.App.Root) {
		logger.L().Info("Devgita repository not found at %s, nothing to uninstall", paths.Paths.App.Root)
	} else if files.IsDirEmpty(paths.Paths.App.Root) {
		logger.L().
			Info("Devgita repository directory is empty at %s, removing directory", paths.Paths.App.Root)
		if err := os.Remove(paths.Paths.App.Root); err != nil {
			return fmt.Errorf("failed to remove empty devgita directory: %w", err)
		}
	} else {
		if err := os.RemoveAll(paths.Paths.App.Root); err != nil {
			return fmt.Errorf("failed to uninstall devgita repository: %w", err)
		}
	}

	if files.FileAlreadyExist(globalConfigPath) {
		logger.L().Info("Removing global config file at %s", globalConfigPath)
		if err := os.Remove(globalConfigPath); err != nil {
			return fmt.Errorf("failed to remove global config file: %w", err)
		}
	} else {
		logger.L().Info("Global config file not found at %s", globalConfigPath)
	}

	if files.DirAlreadyExist(configDirPath) && files.IsDirEmpty(configDirPath) {
		logger.L().Info("Removing empty config directory at %s", configDirPath)
		if err := os.Remove(configDirPath); err != nil {
			return fmt.Errorf("failed to remove empty config directory: %w", err)
		}
	}

	return nil
}

func (dg *Devgita) ForceConfigure() error {
	globalConfig := &config.GlobalConfig{}
	if err := globalConfig.Create(); err != nil {
		return fmt.Errorf("failed to create global config: %w", err)
	}
	if err := globalConfig.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	globalConfig.AppPath = paths.Paths.App.Root
	globalConfig.ConfigPath = configDirPath
	if err := globalConfig.Save(); err != nil {
		return fmt.Errorf("failed to save global config: %w", err)
	}
	return nil
}

func (dg *Devgita) SoftConfigure() error {
	if files.FileAlreadyExist(globalConfigPath) {
		logger.L().Info("Global config already exists at %s", globalConfigPath)
		return nil
	}
	return dg.ForceConfigure()
}
