package claude

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cjairm/devgita/internal/apps"
	"github.com/cjairm/devgita/internal/apps/baseapp"
	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

var _ apps.App = (*Claude)(nil)

type Claude struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func (c *Claude) Name() string       { return constants.Claude }
func (c *Claude) Kind() apps.AppKind { return apps.KindTerminal }

func New() *Claude {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Claude{Cmd: osCmd, Base: baseCmd}
}

func (c *Claude) Install() error {
	params := cmd.CommandParams{
		Command: "sh",
		Args:    []string{"-c", "curl -fsSL https://claude.ai/install.sh | bash"},
	}
	_, _, err := c.Base.ExecCommand(params)
	if err != nil {
		return fmt.Errorf("failed to install claude: %w", err)
	}
	return nil
}

func (c *Claude) ForceInstall() error {
	return baseapp.Reinstall(c.Install, c.Uninstall)
}

func (c *Claude) SoftInstall() error {
	if _, err := exec.LookPath(constants.Claude); err == nil {
		return nil
	}
	return c.Install()
}

func (c *Claude) Uninstall() error {
	gc := &config.GlobalConfig{}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	if _, _, err := c.Base.ExecCommand(cmd.CommandParams{
		Command: "npm",
		Args:    []string{"uninstall", "-g", "@anthropic-ai/claude-code"},
	}); err != nil {
		return fmt.Errorf("failed to uninstall claude: %w", err)
	}
	_ = os.RemoveAll(paths.Paths.Config.Claude)
	gc.DisableShellFeature(constants.Claude)
	if err := gc.RegenerateShellConfig(); err != nil {
		return fmt.Errorf("failed to regenerate shell config: %w", err)
	}
	gc.RemoveFromInstalled(constants.Claude, "package")
	return gc.Save()
}

func (c *Claude) ForceConfigure() error {
	if err := os.MkdirAll(paths.Paths.Config.Claude, 0755); err != nil {
		return err
	}

	gc := &config.GlobalConfig{}
	if err := gc.Create(); err != nil {
		return fmt.Errorf("failed to create global config: %w", err)
	}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}

	if err := files.CopyFile(
		filepath.Join(paths.Paths.App.Configs.Claude, "settings.json"),
		filepath.Join(paths.Paths.Config.Claude, "settings.json"),
	); err != nil {
		return fmt.Errorf("failed to copy claude settings: %w", err)
	}

	statuslineDst := filepath.Join(paths.Paths.Config.Claude, "statusline.sh")
	if err := files.CopyFile(
		filepath.Join(paths.Paths.App.Configs.Claude, "statusline.sh"),
		statuslineDst,
	); err != nil {
		return fmt.Errorf("failed to copy claude statusline: %w", err)
	}
	if err := os.Chmod(statuslineDst, 0755); err != nil {
		return fmt.Errorf("failed to chmod statusline.sh: %w", err)
	}

	if err := files.CopyDir(
		filepath.Join(paths.Paths.App.Configs.Claude, "themes"),
		filepath.Join(paths.Paths.Config.Claude, "themes"),
	); err != nil {
		return fmt.Errorf("failed to copy claude themes: %w", err)
	}

	for _, dir := range []string{"skills", "commands", "agents"} {
		src := filepath.Join(paths.Paths.App.Configs.Shared, dir)
		dst := filepath.Join(paths.Paths.Config.Claude, dir)
		if err := os.MkdirAll(dst, 0755); err != nil {
			return fmt.Errorf("failed to create claude %s dir: %w", dir, err)
		}
		if err := files.CopyDir(src, dst); err != nil {
			return fmt.Errorf("failed to copy claude %s: %w", dir, err)
		}
	}

	gc.AddToInstalled(constants.Claude, "package")
	gc.Shell.Claude = true
	if err := gc.Save(); err != nil {
		return fmt.Errorf("failed to save global config: %w", err)
	}
	if err := gc.RegenerateShellConfig(); err != nil {
		return fmt.Errorf("failed to regenerate shell config: %w", err)
	}
	return nil
}

func (c *Claude) SoftConfigure() error {
	markerFile := filepath.Join(paths.Paths.Config.Claude, "settings.json")
	if files.FileAlreadyExist(markerFile) {
		// Config already exists, but ensure shell feature is enabled
		gc := &config.GlobalConfig{}
		if err := gc.Create(); err != nil {
			return fmt.Errorf("failed to create global config: %w", err)
		}
		if err := gc.Load(); err != nil {
			return fmt.Errorf("failed to load global config: %w", err)
		}
		if !gc.Shell.Claude {
			gc.Shell.Claude = true
			if err := gc.Save(); err != nil {
				return fmt.Errorf("failed to save global config: %w", err)
			}
			if err := gc.RegenerateShellConfig(); err != nil {
				return fmt.Errorf("failed to regenerate shell config: %w", err)
			}
		}
		return nil
	}
	return c.ForceConfigure()
}

func (c *Claude) ExecuteCommand(args ...string) error {
	params := cmd.CommandParams{
		Command: constants.Claude,
		Args:    args,
	}
	_, _, err := c.Base.ExecCommand(params)
	if err != nil {
		return fmt.Errorf("claude command execution failed: %w", err)
	}
	return nil
}

func (c *Claude) Update() error {
	return fmt.Errorf("%w — re-run: curl -fsSL https://claude.ai/install.sh | bash", apps.ErrUpdateNotSupported)
}
