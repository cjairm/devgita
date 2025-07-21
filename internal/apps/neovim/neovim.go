// -------------------------
// TODO: Write documentation how to use this
// - Kickstart documentation: https://github.com/nvim-lua/kickstart.nvim?tab=readme-ov-file
// - Personal configuration: https://github.com/cjairm/devenv/blob/main/nvim/init.lua
// - Releases: https://github.com/neovim/neovim/releases
// - Check version before setup
// - Download app directly instead of using `brew`
//
// NOTE: Is it possible to install different themes?
// If so, see more here: https://linovox.com/the-best-color-schemes-for-neovim-nvim/
// -------------------------

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
	return n.Cmd.InstallPackage("neovim")
}

func (n *Neovim) MaybeInstall() error {
	return n.Cmd.MaybeInstallPackage("neovim")
}

func (n *Neovim) Setup() error {
	if err := n.CheckVersion(); err != nil {
		return err
	}
	return files.CopyDir(paths.NeovimConfigAppDir, paths.NvimConfigLocalDir)
}

func (n *Neovim) MaybeSetup() error {
	isFilePresent := files.FileAlreadyExist(filepath.Join(paths.NvimConfigLocalDir, "init.lua"))
	if isFilePresent {
		return nil
	}
	return n.Setup()
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

func (n *Neovim) Run(args ...string) error {
	execCommand := cmd.CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		IsSudo:      false,
		Command:     constants.Nvim,
		Args:        args,
	}
	if _, err := n.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run Neovim: %w", err)
	}
	return nil
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
