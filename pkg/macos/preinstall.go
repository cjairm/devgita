package macos

import (
	"log"
	"os/exec"

	"github.com/cjairm/devgita/pkg/common"
)

func PreInstall() error {
	if !isHomebrewInstalled() {
		if err := installHomebrew(); err != nil {
			log.Fatalf("Failed to install Homebrew: %v", err)
			return err
		}
	}

	if err := updateHomebrew(); err != nil {
		log.Fatalf("Failed to update Homebrew: %v", err)
		return err
	}

	if !isGitInstalled() {
		if err := installGit(); err != nil {
			log.Fatalf("Failed to install Git: %v", err)
			return err
		}
	}
	return nil
}

func isHomebrewInstalled() bool {
	err := exec.Command("brew", "--version").Run()
	return err == nil
}

func installHomebrew() error {
	return common.ExecCommand(
		"Installing Homebrew",
		"Homebrew installed ✔",
		"bash",
		"-c",
		"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)",
	)
}

func updateHomebrew() error {
	return common.ExecCommand("Updating Homebrew", "Homebrew updated ✔", "brew", "update")
}

func isGitInstalled() bool {
	err := exec.Command("git", "--version").Run()
	return err == nil
}

func installGit() error {
	return common.ExecCommand("Installing Git", "Git installed ✔", "brew", "install", "git")
}
