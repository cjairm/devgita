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
	if err := common.BrewGlobalUpdate(); err != nil {
		log.Fatalf("Failed to update Homebrew: %v", err)
		return err
	}
	if err := common.InstallOrUpdateBrewPackage("git"); err != nil {
		log.Fatalf("Failed to install git: %v", err)
		return err
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
