// OpenCode terminal-based AI code editor with devgita integration
//
// OpenCode is an AI-powered code editor that runs in the terminal, providing
// intelligent code completion, refactoring, and assistance. This module provides
// installation and configuration management for OpenCode with devgita integration.
//
// References:
// - OpenCode Documentation: https://opencode.ai/docs

package opencode

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cjairm/devgita/internal/apps"
	"github.com/cjairm/devgita/internal/apps/baseapp"
	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

var (
	_ apps.App                 = (*OpenCode)(nil)
	_ apps.SelectiveConfigurer = (*OpenCode)(nil)
)

const DEFAULT_THEME_NAME = "default"

type OpenCode struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func (o *OpenCode) Name() string       { return constants.OpenCode }
func (o *OpenCode) Kind() apps.AppKind { return apps.KindTerminal }

func New() *OpenCode {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &OpenCode{Cmd: osCmd, Base: baseCmd}
}

func (o *OpenCode) Install() error {
	return o.Cmd.InstallPackage(constants.OpenCode)
}

func (o *OpenCode) ForceInstall() error {
	return baseapp.Reinstall(o.Install, o.Uninstall)
}

func (o *OpenCode) SoftInstall() error {
	return o.Cmd.MaybeInstallPackage(constants.OpenCode)
}

func (o *OpenCode) Uninstall() error {
	gc := &config.GlobalConfig{}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	if err := o.Cmd.UninstallPackage(constants.OpenCode); err != nil {
		return fmt.Errorf("failed to uninstall opencode: %w", err)
	}
	_ = os.RemoveAll(paths.Paths.Config.OpenCode)
	gc.DisableShellFeature(constants.OpenCode)
	if err := gc.RegenerateShellConfig(); err != nil {
		return fmt.Errorf("failed to regenerate shell config: %w", err)
	}
	gc.RemoveFromInstalled(constants.OpenCode, "package")
	return gc.Save()
}

func (o *OpenCode) ForceConfigure() error {
	if err := os.RemoveAll(paths.Paths.Config.OpenCode); err != nil {
		return err
	}
	// Directory permissions should be 0755 not 0644. Directories need execute
	// permission to be entered.
	if err := os.MkdirAll(paths.Paths.Config.OpenCode, 0o755); err != nil {
		return err
	}
	gc := &config.GlobalConfig{}
	if err := gc.Create(); err != nil {
		return fmt.Errorf("failed to create global config: %w", err)
	}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	theme := DEFAULT_THEME_NAME
	configFilePath := filepath.Join(
		paths.Paths.Config.OpenCode,
		fmt.Sprintf("%s.json", constants.OpenCode),
	)
	tmplPath := filepath.Join(
		paths.Paths.App.Configs.OpenCode,
		fmt.Sprintf("%s.json.tmpl", constants.OpenCode),
	)
	if theme == DEFAULT_THEME_NAME {
		themesDir := filepath.Join(paths.Paths.Config.OpenCode, "themes")
		if err := os.MkdirAll(themesDir, 0o755); err != nil {
			return fmt.Errorf("failed to create themes directory: %w", err)
		}
		if err := files.CopyFile(
			filepath.Join(
				paths.Paths.App.Configs.OpenCode,
				"themes",
				fmt.Sprintf("%s.json", DEFAULT_THEME_NAME),
			),
			filepath.Join(themesDir, fmt.Sprintf("%s.json", DEFAULT_THEME_NAME)),
		); err != nil {
			return fmt.Errorf("failed to copy opencode config theme: %w", err)
		}
	}
	if err := files.GenerateFromTemplate(tmplPath, configFilePath, map[string]string{
		"Theme": theme,
	}); err != nil {
		return fmt.Errorf("failed to generate opencode configuration: %w", err)
	}
	if err := baseapp.SyncSharedParts(
		paths.Paths.Config.OpenCode,
		baseapp.SharedConfigParts,
	); err != nil {
		return fmt.Errorf("failed to copy opencode shared config: %w", err)
	}

	// The task-redirect plugin (and any future local OpenCode plugins) ships
	// from configs/opencode/plugin/, not configs/shared/ — plugins are an
	// OpenCode-specific mechanism, outside SharedConfigParts' skills/commands/
	// agents sync surface. OpenCode loads plugin files from
	// ~/.config/opencode/plugin/ (or the singular/plural "plugins" variant;
	// see task-redirect.js's header comment).
	// CopyDir creates its destination directory itself (see pkg/files.CopyDir),
	// same as every other CopyDir call site in this codebase — no explicit
	// MkdirAll needed here.
	pluginDst := filepath.Join(paths.Paths.Config.OpenCode, "plugin")
	if err := files.CopyDir(
		filepath.Join(paths.Paths.App.Configs.OpenCode, "plugin"),
		pluginDst,
	); err != nil {
		return fmt.Errorf("failed to copy opencode plugins: %w", err)
	}

	gc.ReconcileShellFeatures()
	gc.AddToInstalled(constants.OpenCode, "package")
	gc.Shell.Opencode = true
	if err := gc.Save(); err != nil {
		return fmt.Errorf("failed to save global config: %w", err)
	}
	if err := gc.RegenerateShellConfig(); err != nil {
		return fmt.Errorf("failed to regenerate shell config: %w", err)
	}
	return nil
}

func (o *OpenCode) SoftConfigure() error {
	markerFile := filepath.Join(
		paths.Paths.Config.OpenCode,
		fmt.Sprintf("%s.json", constants.OpenCode),
	)
	if files.FileAlreadyExist(markerFile) {
		// Config already exists, but ensure shell feature is enabled
		gc := &config.GlobalConfig{}
		if err := gc.Create(); err != nil {
			return fmt.Errorf("failed to create global config: %w", err)
		}
		if err := gc.Load(); err != nil {
			return fmt.Errorf("failed to load global config: %w", err)
		}
		if !gc.Shell.Opencode {
			gc.Shell.Opencode = true
			if err := gc.Save(); err != nil {
				return fmt.Errorf("failed to save global config: %w", err)
			}
			if err := gc.RegenerateShellConfig(); err != nil {
				return fmt.Errorf("failed to regenerate shell config: %w", err)
			}
		}
		return nil
	}
	return o.ForceConfigure()
}

// ConfigurableParts lists the shared config subtrees that --only can refresh.
func (o *OpenCode) ConfigurableParts() []string { return baseapp.SharedConfigParts }

// ForceConfigureParts overwrites only the named shared subtrees (skills,
// commands, agents) under the OpenCode config dir. Unlike full ForceConfigure
// it does not remove or regenerate opencode.json or themes, so a hand-edited
// config survives. This is the `--force --only=...` path.
func (o *OpenCode) ForceConfigureParts(parts []string) error {
	if err := os.MkdirAll(paths.Paths.Config.OpenCode, 0o755); err != nil {
		return err
	}
	return baseapp.SyncSharedParts(paths.Paths.Config.OpenCode, parts)
}

func (o *OpenCode) ExecuteCommand(args ...string) error {
	params := cmd.CommandParams{
		Command: constants.OpenCode,
		Args:    args,
	}
	_, _, err := o.Base.ExecCommand(params)
	if err != nil {
		return fmt.Errorf("opencode command execution failed: %w", err)
	}
	return nil
}

func (o *OpenCode) Update() error {
	return fmt.Errorf("%w for opencode", apps.ErrUpdateNotSupported)
}
