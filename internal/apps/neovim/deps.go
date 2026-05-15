package neovim

import (
	"fmt"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/logger"
)

// InstallDeps installs all system packages required for Neovim to function correctly.
// Hard deps (make, gcc, ripgrep, fd-find, unzip, xclip) return an error if they fail.
// tree-sitter-cli is soft: primary package manager → npm fallback → warn-and-continue.
// Called by Install() and SoftInstall() before the Neovim binary is installed.
func InstallDeps(base cmd.BaseCommandExecutor, c cmd.Command) error {
	// make — required by Neovim plugin ecosystem (e.g. telescope-fzf-native)
	if err := c.MaybeInstallPackage(constants.Make); err != nil {
		return fmt.Errorf("failed to install make: %w", err)
	}

	// gcc — required for compiling Neovim plugins on Linux
	// On macOS, Xcode CLT provides clang; brew gcc is unnecessary
	if !base.IsMac() {
		if err := c.MaybeInstallPackage(constants.Gcc); err != nil {
			return fmt.Errorf("failed to install gcc: %w", err)
		}
	}

	// ripgrep — live grep search in Neovim (e.g. telescope live_grep)
	if err := c.MaybeInstallPackage(constants.Ripgrep); err != nil {
		return fmt.Errorf("failed to install ripgrep: %w", err)
	}

	// fd-find — fast file finder used by Neovim plugins (e.g. telescope find_files)
	if err := c.MaybeInstallPackage(constants.FdFind); err != nil {
		return fmt.Errorf("failed to install fd-find: %w", err)
	}

	// unzip — required for extracting Neovim plugin archives
	if err := c.MaybeInstallPackage(constants.Unzip); err != nil {
		return fmt.Errorf("failed to install unzip: %w", err)
	}

	// xclip — clipboard integration for Neovim on Linux
	// macOS uses pbcopy/pbpaste which are built-in
	if !base.IsMac() {
		if err := c.MaybeInstallPackage(constants.Xclip); err != nil {
			return fmt.Errorf("failed to install xclip: %w", err)
		}
	}

	// tree-sitter-cli — best-effort; fallback to npm if primary package manager fails
	installTreeSitter(base, c)

	return nil
}

// installTreeSitter attempts to install tree-sitter-cli via the primary package
// manager (brew on macOS, apt on Linux), falling back to npm install -g.
// On Debian Bookworm (stable), tree-sitter-cli is only in trixie/sid, so the
// apt path is expected to fail and npm is the real install path.
// If both fail, a warning is logged and nil is returned — Neovim still installs.
func installTreeSitter(base cmd.BaseCommandExecutor, c cmd.Command) {
	primaryErr := c.MaybeInstallPackage(constants.TreeSitterCli)
	if primaryErr == nil {
		return
	}
	logger.L().Warnw("Primary tree-sitter-cli install failed, trying npm fallback",
		"error", primaryErr)

	_, stderr, err := base.ExecCommand(cmd.CommandParams{
		Command: "npm",
		Args:    []string{"install", "-g", "tree-sitter-cli"},
	})
	if err != nil {
		logger.L().Warnw("npm fallback for tree-sitter-cli also failed — skipping",
			"error", err, "stderr", stderr)
		return
	}

	// Track npm-installed tree-sitter-cli in global_config for transparent state
	gc := &config.GlobalConfig{}
	if loadErr := gc.Load(); loadErr != nil {
		logger.L().Warnw("Could not load global config to track tree-sitter-cli",
			"error", loadErr)
		return
	}
	gc.AddToInstalled(constants.TreeSitterCli, "package")
	if saveErr := gc.Save(); saveErr != nil {
		logger.L().Warnw("Could not save global config after tracking tree-sitter-cli",
			"error", saveErr)
	}
}
