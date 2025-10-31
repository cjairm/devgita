// Package neovim provides installation and configuration management for Neovim editor with devgita integration.
// It follows the standardized devgita app interface while providing Neovim-specific operations for editor
// setup, LSP configuration, and development environment customization.
//
// References:
// - Kickstart documentation: https://github.com/nvim-lua/kickstart.nvim
// - Neovim releases: https://github.com/neovim/neovim/releases
// - Configuration examples: https://github.com/cjairm/devenv/blob/main/nvim/init.lua
// - Color schemes: https://linovox.com/the-best-color-schemes-for-neovim-nvim/

package neovim

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

type Neovim struct {
	Cmd  cmd.Command
	Base cmd.BaseCommand
}

func New() *Neovim {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Neovim{Cmd: osCmd, Base: *baseCmd}
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
	if err := n.CheckVersion(); err != nil {
		return fmt.Errorf("failed to check Neovim version: %w", err)
	}
	return files.CopyDir(paths.NeovimConfigAppDir, paths.NvimConfigLocalDir)
}

func (n *Neovim) SoftConfigure() error {
	isFilePresent := files.FileAlreadyExist(filepath.Join(paths.NvimConfigLocalDir, "init.lua"))
	if isFilePresent {
		return nil
	}
	return n.ForceConfigure()
}

func (n *Neovim) Uninstall() error {
	return fmt.Errorf("uninstall operation not supported for Neovim")
}

func (n *Neovim) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		IsSudo:      false,
		Command:     constants.Nvim,
		Args:        args,
	}
	if _, _, err := n.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to execute Neovim command: %w", err)
	}
	return nil
}

func (n *Neovim) Update() error {
	return fmt.Errorf("update operation not implemented for Neovim")
}

func (n *Neovim) CheckVersion() error {
	cmd := exec.Command(constants.Nvim, "--version")
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "NVIM v") {
			versionStr := strings.TrimPrefix(line, "NVIM v")
			versionStr = strings.Fields(versionStr)[0]
			if isVersionEqualOrHigher(versionStr, constants.NeovimVersion) {
				return nil
			}
		}
	}
	return fmt.Errorf("could not parse Neovim version")
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

// Deprecated: Use SoftInstall() instead
func (n *Neovim) MaybeInstall() error {
	return n.SoftInstall()
}

// Deprecated: Use ForceConfigure() instead
func (n *Neovim) Setup() error {
	return n.ForceConfigure()
}

// Deprecated: Use SoftConfigure() instead
func (n *Neovim) MaybeSetup() error {
	return n.SoftConfigure()
}

// Deprecated: Use ExecuteCommand() instead
func (n *Neovim) Run(args ...string) error {
	return n.ExecuteCommand(args...)
}
