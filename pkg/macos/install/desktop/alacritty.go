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

func InstallAlacritty(devgitaPath string) error {
	if err := common.MaybeInstallBrewCask("alacritty"); err != nil {
		return err
	}
	// Configure Alacritty
	if err := configAlacritty(devgitaPath); err != nil {
		return fmt.Errorf("Error copying alacritty config: %v", err)
	}
	// Update config file to match installer path
	if err := updateHomeDirPath(); err != nil {
		return fmt.Errorf("Error updating alacritty config: %v", err)
	}
	return nil
}

func configAlacritty(devgitaPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("Error getting home directory: %w", err)
	}
	configDir := filepath.Join(homeDir, ".config", "alacritty")
	configFile := filepath.Join(configDir, "alacritty.toml")
	devgitaConfig := filepath.Join(
		devgitaPath,
		"configs",
		"alacritty",
	)
	devgitaFontDir := filepath.Join(
		devgitaPath,
		"configs",
		"fonts",
		"alacritty",
		"default",
	)
	devgitaThemeDir := filepath.Join(
		devgitaPath,
		"configs",
		"themes",
		"alacritty",
		"default",
	)
	if common.FileAlreadyExist(configFile) {
		fmt.Printf("%s already exist\n\n", configDir)
		return nil
	}
	// Moves general config
	if err := common.MoveContents(devgitaConfig, configDir); err != nil {
		return fmt.Errorf("error setting up alacritty: %w", err)
	}
	// Moves font config
	if err := common.MoveContents(devgitaFontDir, configDir); err != nil {
		return fmt.Errorf("error setting up alacritty's font: %w", err)
	}
	// Moves theme config
	if err := common.MoveContents(devgitaThemeDir, configDir); err != nil {
		return fmt.Errorf("error setting up alacritty's theme: %w", err)
	}
	fmt.Printf("Alacritty configuration set successfully!\n\n")
	return nil
}

func updateHomeDirPath() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("Error getting home directory: %w", err)
	}
	configFile := filepath.Join(
		homeDir,
		".config", "alacritty",
		"alacritty.toml",
	)
	return common.UpdateFile(configFile, "<HOME-PATH>")
}
