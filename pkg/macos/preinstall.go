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
	cmd := common.CommandInfo{
		PreExecutionMessage:  "Installing Homebrew",
		PostExecutionMessage: "Homebrew installed ✔",
		IsSudo:               false,
		Command:              "bash",
		Args: []string{
			"-c",
			"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)",
		},
	}
	return common.ExecCommand(cmd)
}

func updateHomebrew() error {
	cmd := common.CommandInfo{
		PreExecutionMessage:  "Updating Homebrew",
		PostExecutionMessage: "Homebrew updated ✔",
		IsSudo:               false,
		Command:              "brew",
		Args: []string{
			"update",
		},
	}
	return common.ExecCommand(cmd)
}

func isGitInstalled() bool {
	err := exec.Command("git", "--version").Run()
	return err == nil
}

func installGit() error {
	cmd := common.CommandInfo{
		PreExecutionMessage:  "Installing Git",
		PostExecutionMessage: "Git installed ✔",
		IsSudo:               false,
		Command:              "brew",
		Args: []string{
			"install",
			"git",
		},
	}
	return common.ExecCommand(cmd)
}
