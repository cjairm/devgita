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

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

const DEFAULT_THEME_NAME = "default"

type OpenCode struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

type ConfigureOptions struct {
	Theme string
}

func New() *OpenCode {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &OpenCode{Cmd: osCmd, Base: baseCmd}
}

func (o *OpenCode) Install() error {
	return o.Cmd.InstallPackage(constants.OpenCode)
}

func (o *OpenCode) ForceInstall() error {
	err := o.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall opencode: %w", err)
	}
	return o.Install()
}

func (o *OpenCode) SoftInstall() error {
	return o.Cmd.MaybeInstallPackage(constants.OpenCode)
}

func (o *OpenCode) Uninstall() error {
	return fmt.Errorf("opencode uninstall not supported through devgita")
}

func (o *OpenCode) ForceConfigure(opts ...ConfigureOptions) error {
	var options ConfigureOptions
	if len(opts) > 0 {
		options = opts[0]
	}

	if err := os.RemoveAll(paths.Paths.Config.OpenCode); err != nil {
		return err
	}
	// Directory permissions should be 0755 not 0644. Directories need execute
	// permission to be entered.
	if err := os.MkdirAll(paths.Paths.Config.OpenCode, 0755); err != nil {
		return err
	}
	gc := &config.GlobalConfig{}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	theme := DEFAULT_THEME_NAME
	if options.Theme != "" {
		theme = options.Theme
	}
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
		if err := os.MkdirAll(themesDir, 0755); err != nil {
			return fmt.Errorf("failed to create themes directory: %w", err)
		}
		if err := files.CopyFile(
			filepath.Join(paths.Paths.App.Configs.OpenCode, "themes", fmt.Sprintf("%s.json", DEFAULT_THEME_NAME)),
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
	gc.AddToInstalled(constants.OpenCode, "package")
	if err := gc.Save(); err != nil {
		return fmt.Errorf("failed to save global config: %w", err)
	}
	return nil
}

func (o *OpenCode) SoftConfigure(opts ...ConfigureOptions) error {
	var options ConfigureOptions
	if len(opts) > 0 {
		options = opts[0]
	}

	gc := &config.GlobalConfig{}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	if gc.IsAlreadyInstalled(constants.OpenCode, "package") ||
		gc.IsInstalledByDevgita(constants.OpenCode, "package") {
		return nil
	}
	markerFile := filepath.Join(
		paths.Paths.Config.OpenCode,
		fmt.Sprintf("%s.json", constants.OpenCode),
	)
	if files.FileAlreadyExist(markerFile) {
		return nil
	}
	return o.ForceConfigure(options)
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
	return fmt.Errorf("opencode update not implemented - use system package manager")
}
