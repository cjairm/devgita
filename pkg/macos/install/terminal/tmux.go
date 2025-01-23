// -------------------------
// TODO: Write documentation how to use this
// - Tmux documentation: https://github.com/tmux/tmux
// - Personal configuration: https://github.com/cjairm/devenv/tree/main/tmux
// - Releases: https://github.com/tmux/tmux/releases
// - Installing instructions: https://github.com/tmux/tmux/wiki/Installing
// -------------------------

package macos

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cjairm/devgita/pkg/common"
)

func InstallTmux(devgitaPath string) error {
	if common.IsCommandInstalled("tmux") {
		if err := upgradeTmux(); err != nil {
			return err
		}
	} else {
		if err := installTmux(); err != nil {
			return err
		}
	}
	// Copy .tmux.conf to the home directory
	if err := copyTmuxConfig(devgitaPath); err != nil {
		return fmt.Errorf("Error copying .tmux.conf: %v", err)
	}
	return nil
}

// installTmux installs Tmux using Homebrew.
func installTmux() error {
	cmd := common.CommandInfo{
		PreExecutionMessage:  "Installing Tmux using Homebrew...",
		PostExecutionMessage: "Tmux installed successfully ✔",
		IsSudo:               false,
		Command:              "brew",
		Args:                 []string{"install", "tmux"},
	}
	return common.ExecCommand(cmd)
}

// upgradeTmux upgrades Tmux using Homebrew.
func upgradeTmux() error {
	cmd := common.CommandInfo{
		PreExecutionMessage:  "Upgrading Tmux using Homebrew...",
		PostExecutionMessage: "Tmux upgraded ✔",
		IsSudo:               false,
		Command:              "brew",
		Args:                 []string{"upgrade", "tmux"},
	}
	return common.ExecCommand(cmd)
}

func copyTmuxConfig(devgitaPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("Error getting home directory: %w", err)
	}
	destinationFile := filepath.Join(homeDir, ".tmux.conf")
	if common.FileAlreadyExist(destinationFile) {
		fmt.Println("Tmux configuration already exists!")
		return nil
	}
	// Define the source directory for .tmux.conf
	sourceDir := filepath.Join(devgitaPath, "pkg", "configs", "tmux")
	sourceFile := filepath.Join(sourceDir, ".tmux.conf")
	// Check if the source file exists
	if common.FileAlreadyExist(sourceFile) {
		// Copy the .tmux.conf file to the home directory
		if err := common.CopyFile(sourceFile, destinationFile); err != nil {
			return fmt.Errorf("Error copying .tmux.conf: %w", err)
		}
		fmt.Println(".tmux.conf copied to home directory successfully!")
	}
	return nil
}
