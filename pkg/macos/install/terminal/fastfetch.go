// -------------------------
// NOTE: Write documentation or create icon to open and get information of this Mac
// - Documentation: https://github.com/fastfetch-cli/fastfetch
// -------------------------

package macos

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cjairm/devgita/pkg/common"
)

func InstallFastFetch(devgitaPath string) error {
	if err := common.InstallOrUpdateBrewPackage("fastfetch"); err != nil {
		return err
	}
	if err := configureFastFetch(devgitaPath); err != nil {
		return fmt.Errorf("Error copying fastfetch config: %v", err)
	}
	return nil

}

func configureFastFetch(devgitaPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("Error getting home directory: %w", err)
	}
	configDir := filepath.Join(homeDir, ".config", "fastfetch")
	configFile := filepath.Join(configDir, "config.jsonc")
	devgitaConfig := filepath.Join(
		devgitaPath,
		"configs",
		"fastfetch",
	)
	if common.FileAlreadyExist(configFile) {
		fmt.Printf("%s already exist\n\n", configFile)
		return nil
	}
	if err := common.MoveContents(devgitaConfig, configDir); err != nil {
		return fmt.Errorf("error setting up fastfetch: %w", err)
	}
	fmt.Printf("FastFetch configuration set successfully!\n\n")
	return nil
}
