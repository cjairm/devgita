// -------------------------
// TODO: Write documentation how to use this
// - Kickstart documentation: https://github.com/nvim-lua/kickstart.nvim?tab=readme-ov-file
// - Personal configuration: https://github.com/cjairm/devenv/blob/main/nvim/init.lua
// - Releases: https://github.com/neovim/neovim/releases
//
// NOTE: Is it possible to install different themes?
// If so, see more here: https://linovox.com/the-best-color-schemes-for-neovim-nvim/
// -------------------------

package macos

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cjairm/devgita/pkg/common"
)

func InstallNeovim(devgitaPath string) error {
	if err := common.InstallOrUpdateBrewPackage("neovim"); err != nil {
		return err
	}
	// Configure Neovim
	if err := configNeovim(devgitaPath); err != nil {
		return fmt.Errorf("Error copying neovim config: %v", err)
	}
	return nil
}

func configNeovim(devgitaPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("Error getting home directory: %w", err)
	}
	configDir := filepath.Join(homeDir, ".config", "nvim")
	configFile := filepath.Join(configDir, "init.lua")
	devgitaConfig := filepath.Join(
		devgitaPath,
		"configs",
		"neovim",
	)
	if common.FileAlreadyExist(configFile) {
		fmt.Printf("%s already exist\n\n", configDir)
		return nil
	}
	if err := common.MoveContents(devgitaConfig, configDir); err != nil {
		return fmt.Errorf("error setting up fastfetch: %w", err)
	}
	fmt.Printf("Neovim configuration set successfully!\n\n")
	return nil
}
