package common

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func CloneDevgita() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %w", err)
	}

	devgitaPath := filepath.Join(homeDir, ".local", "share", "devgita")

	fmt.Printf("Cloning Devgita\n")
	err = os.RemoveAll(devgitaPath)
	if err != nil {
		return fmt.Errorf("error removing directory: %w", err)
	}

	cmd := exec.Command("git", "clone", "https://github.com/cjairm/devgita.git", devgitaPath)
	cmd.Stdout = nil // Redirect stdout to nil (equivalent to >/dev/null)
	cmd.Stderr = nil // Redirect stderr to nil (optional, can handle errors)

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error cloning repository: %w", err)
	}

	fmt.Printf("Devgita cloned ✔\n\n")
	return nil
}