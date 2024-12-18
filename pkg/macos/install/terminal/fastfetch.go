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
	if err := checkIfFastfetchIsInstalled(); err != nil {
		return fmt.Errorf("Error checking fastfetch: %w", err)
	}
	if err := setupFastFetch(devgitaPath); err != nil {
		return fmt.Errorf("Error setting up config: %w", err)
	}
	return nil
}

func checkIfFastfetchIsInstalled() error {
	if !common.IsCommandInstalled("fastfetch") {
		return installFastfetch()
	}
	return nil
}

func installFastfetch() error {
	return common.ExecCommand(
		"Installing fastfetch",
		"fastfetch installed ✔",
		"brew",
		"install",
		"fastfetch",
	)
}

func setupFastFetch(devgitaPath string) error {
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "fastfetch")
	configFile := filepath.Join(configDir, "config.jsonc")
	devgitaConfig := filepath.Join(
		devgitaPath,
		"pkg",
		"configs",
		"fastfetch.jsonc",
	)

	return common.MkdirOrCopyFile(configFile, configDir, devgitaConfig, "fastfetch config")
}
