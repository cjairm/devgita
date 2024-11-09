package macos

import (
	"fmt"
	"os"
	"os/exec"
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
	_, err := exec.LookPath("fastfetch")
	if err != nil {
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

	// Check if the configuration file already exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create the config directory if it doesn't exist
		if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
			return err
		}
		return copyFile(devgitaConfig, configFile)
	}
	fmt.Println("Configuration file for fastfetch already exists.")
	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, os.ModePerm)
}
