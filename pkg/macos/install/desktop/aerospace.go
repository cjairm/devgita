package macos

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cjairm/devgita/pkg/common"
)

func InstallAerospace(devgitaPath string) error {
	if err := common.MaybeInstallBrewCask("nikitabobko/tap/aerospace"); err != nil {
		return err
	}
	// Configure Aerospace
	if err := configAerospace(devgitaPath); err != nil {
		return fmt.Errorf("Error copying aerospace config: %v", err)
	}
	// // Update config file to match installer path
	// if err := updateHomeDirPath(); err != nil {
	// 	return fmt.Errorf("Error updating alacritty config: %v", err)
	// }
	return nil
}

func configAerospace(devgitaPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("Error getting home directory: %w", err)
	}
	configDir := filepath.Join(homeDir, ".config", "aerospace")
	configFile := filepath.Join(configDir, "aerospace.toml")
	devgitaConfig := filepath.Join(
		devgitaPath,
		"configs",
		"aerospace",
	)
	if common.FileAlreadyExist(configFile) {
		fmt.Printf("%s already exist\n\n", configDir)
		return nil
	}
	// Moves general config
	if err := common.MoveContents(devgitaConfig, configDir); err != nil {
		return fmt.Errorf("error setting up alacritty: %w", err)
	}
	fmt.Printf("Aerospace configuration set successfully!\n\n")
	return nil
}

// func updateHomeDirPath() error {
// 	homeDir, err := os.UserHomeDir()
// 	if err != nil {
// 		return fmt.Errorf("Error getting home directory: %w", err)
// 	}
// 	configFile := filepath.Join(
// 		homeDir,
// 		".config", "alacritty",
// 		"alacritty.toml",
// 	)
// 	return common.UpdateFile(configFile, "<HOME-PATH>")
// }
