// -------------------------
// TODO: Write documentation how to use this
// - Tmux documentation: https://github.com/tmux/tmux
// - Personal configuration: https://github.com/cjairm/devenv/tree/main/tmux
// - Releases: https://github.com/tmux/tmux/releases
// - Installing instructions: https://github.com/tmux/tmux/wiki/Installing
// -------------------------

package tmux

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cjairm/devgita/internal/apps"
	"github.com/cjairm/devgita/internal/apps/baseapp"
	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

var _ apps.App = (*Tmux)(nil)

const configFileName = ".tmux.conf"

type Tmux struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func (t *Tmux) Name() string       { return constants.Tmux }
func (t *Tmux) Kind() apps.AppKind { return apps.KindTerminal }

func New() *Tmux {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Tmux{Cmd: osCmd, Base: baseCmd}
}

func (t *Tmux) Install() error {
	return t.Cmd.InstallPackage(constants.Tmux)
}

func (t *Tmux) ForceInstall() error {
	return baseapp.Reinstall(t.Install, t.Uninstall)
}

func (t *Tmux) SoftInstall() error {
	return t.Cmd.MaybeInstallPackage(constants.Tmux)
}

func (t *Tmux) ForceConfigure() error {
	gc := &config.GlobalConfig{}
	if err := gc.Create(); err != nil {
		return fmt.Errorf("failed to create global config: %w", err)
	}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	gc.AddToInstalled(constants.Tmux, "package")
	if err := enableFeature(gc); err != nil {
		return fmt.Errorf("failed to enable tmux feature: %w", err)
	}
	return files.CopyFile(
		filepath.Join(paths.Paths.App.Configs.Tmux, "tmux.conf"),
		filepath.Join(paths.Paths.Home.Root, configFileName),
	)
}

func (t *Tmux) SoftConfigure() error {
	configFile := filepath.Join(paths.Paths.Home.Root, configFileName)
	isFilePresent := files.FileAlreadyExist(configFile)
	if isFilePresent {
		gc := &config.GlobalConfig{}
		if err := gc.Create(); err != nil {
			return fmt.Errorf("failed to create global config: %w", err)
		}
		if err := gc.Load(); err != nil {
			return fmt.Errorf("failed to load global config: %w", err)
		}
		if !gc.IsShellFeatureEnabled(constants.Tmux) {
			if err := enableFeature(gc); err != nil {
				return fmt.Errorf("failed to enable tmux feature: %w", err)
			}
		}
		return nil
	}
	return t.ForceConfigure()
}

func (t *Tmux) Uninstall() error {
	gc := &config.GlobalConfig{}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	if err := t.Cmd.UninstallPackage(constants.Tmux); err != nil {
		return fmt.Errorf("failed to uninstall tmux: %w", err)
	}
	_ = os.Remove(filepath.Join(paths.Paths.Home.Root, configFileName))
	gc.DisableShellFeature(constants.Tmux)
	if err := gc.RegenerateShellConfig(); err != nil {
		return fmt.Errorf("failed to regenerate shell config: %w", err)
	}
	gc.RemoveFromInstalled(constants.Tmux, "package")
	return gc.Save()
}

func (t *Tmux) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		Command: constants.Tmux,
		Args:    args,
	}
	if _, _, err := t.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to execute tmux command: %w", err)
	}
	return nil
}

func (t *Tmux) Update() error {
	return fmt.Errorf("%w for tmux", apps.ErrUpdateNotSupported)
}

// CreateSession creates a new detached tmux session in the given directory
func (t *Tmux) CreateSession(name, workdir string) error {
	return t.ExecuteCommand("new-session", "-d", "-s", name, "-c", workdir)
}

// KillSession terminates a tmux session
func (t *Tmux) KillSession(name string) error {
	return t.ExecuteCommand("kill-session", "-t", name)
}

// HasSession checks if a session exists
func (t *Tmux) HasSession(name string) bool {
	err := t.ExecuteCommand("has-session", "-t", name)
	return err == nil
}

// SendKeys sends keystrokes to a session
func (t *Tmux) SendKeys(session, keys string) error {
	return t.ExecuteCommand("send-keys", "-t", session, keys, "Enter")
}

// CreateWindow creates a new window in the current session
func (t *Tmux) CreateWindow(name, workdir string) error {
	return t.ExecuteCommand("new-window", "-n", name, "-c", workdir)
}

// HasWindow checks if a window exists in the current session
func (t *Tmux) HasWindow(name string) bool {
	if os.Getenv("TMUX") == "" {
		return false
	}
	execCommand := cmd.CommandParams{
		Command: constants.Tmux,
		Args:    []string{"list-windows", "-F", "#{window_name}"},
	}
	stdout, _, err := t.Base.ExecCommand(execCommand)
	if err != nil {
		return false
	}
	scanner := bufio.NewScanner(strings.NewReader(stdout))
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == name {
			return true
		}
	}
	return false
}

// KillWindow closes a specific window by name
func (t *Tmux) KillWindow(name string) error {
	return t.ExecuteCommand("kill-window", "-t", name)
}

// SendKeysToWindow sends keystrokes to a specific window
func (t *Tmux) SendKeysToWindow(window, keys string) error {
	return t.ExecuteCommand("send-keys", "-t", window, keys, "Enter")
}

// SelectWindow switches focus to a specific window by name
func (t *Tmux) SelectWindow(name string) error {
	return t.ExecuteCommand("select-window", "-t", name)
}

func enableFeature(gc *config.GlobalConfig) error {
	gc.EnableShellFeature(constants.Tmux)
	if err := gc.RegenerateShellConfig(); err != nil {
		return fmt.Errorf("failed to generate shell config: %w", err)
	}
	if err := gc.Save(); err != nil {
		return fmt.Errorf("failed to save global config: %w", err)
	}
	return nil
}
