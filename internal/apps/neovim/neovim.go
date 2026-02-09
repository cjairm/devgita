// -------------------------
// TODO: Write documentation how to use this
// - Kickstart documentation: https://github.com/nvim-lua/kickstart.nvim?tab=readme-ov-file
// - Personal configuration: https://github.com/cjairm/devenv/blob/main/nvim/init.lua
// - Releases: https://github.com/neovim/neovim/releases
// - Download app directly instead of using `brew`?
//
// NOTE: install different themes...?
// -------------------------

package neovim

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

type Neovim struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func New() *Neovim {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Neovim{Cmd: osCmd, Base: baseCmd}
}

func (n *Neovim) Install() error {
	return n.Cmd.InstallPackage(constants.Neovim)
}

func (n *Neovim) ForceInstall() error {
	err := n.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall Neovim before force install: %w", err)
	}
	return n.Install()
}

func (n *Neovim) SoftInstall() error {
	return n.Cmd.MaybeInstallPackage(constants.Neovim)
}

func (n *Neovim) ForceConfigure() error {
	if err := n.checkVersion(); err != nil {
		return fmt.Errorf("failed to check Neovim version: %w", err)
	}
	return files.CopyDir(paths.Paths.App.Configs.Neovim, paths.Paths.Config.Nvim)
}

func (n *Neovim) SoftConfigure() error {
	isDirPresent := files.DirAlreadyExist(paths.Paths.Config.Nvim)
	isDirEmpty := files.IsDirEmpty(paths.Paths.Config.Nvim)
	if isDirPresent && !isDirEmpty {
		return nil
	}
	return n.ForceConfigure()
}

func (n *Neovim) Uninstall() error {
	return fmt.Errorf("uninstall operation not supported for Neovim")
}

func (n *Neovim) ExecuteCommand(args ...string) error {
	baseCmd := getBaseCmd(args...)
	_, stderr, err := n.Base.ExecCommand(baseCmd)
	if err != nil {
		return fmt.Errorf("failed to check neovim version: %w, stderr: %s", err, stderr)
	}
	return nil
}

func (n *Neovim) Update() error {
	return fmt.Errorf("update operation not implemented for Neovim")
}

func (n *Neovim) checkVersion() error {
	baseCmd := getBaseCmd("--version")
	stdout, stderr, err := n.Base.ExecCommand(baseCmd)
	if err != nil {
		return fmt.Errorf("failed to check neovim version: %w, stderr: %s", err, stderr)
	}
	scanner := bufio.NewScanner(strings.NewReader(stdout))
	for scanner.Scan() {
		line := scanner.Text()
		if versionStr, found := strings.CutPrefix(line, "NVIM v"); found {
			versionStr = strings.Fields(versionStr)[0]
			if isVersionEqualOrHigher(versionStr, constants.SupportedVersion.Neovim.Number) {
				return nil
			}
			return fmt.Errorf(
				"neovim version %s is too old, requires %s",
				versionStr,
				constants.SupportedVersion.Neovim.Number,
			)
		}
	}
	return fmt.Errorf("could not parse Neovim version from output")
}

func isVersionEqualOrHigher(currentVersion, requiredVersion string) bool {
	currentParts := strings.Split(currentVersion, ".")
	requiredParts := strings.Split(requiredVersion, ".")
	for i, requiredPartStr := range requiredParts {
		if i >= len(currentParts) {
			return false // Current version has fewer parts
		}
		currentPart, err := strconv.Atoi(currentParts[i])
		if err != nil {
			return false
		}
		requiredPart, err := strconv.Atoi(requiredPartStr)
		if err != nil {
			return false
		}
		if currentPart < requiredPart {
			return false
		}
	}
	return true
}

func getBaseCmd(args ...string) cmd.CommandParams {
	return cmd.CommandParams{
		Command: constants.Nvim,
		Args:    args,
	}
}
