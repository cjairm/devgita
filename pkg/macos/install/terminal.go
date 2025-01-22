package macos

import (
	"fmt"
	"os"

	"github.com/cjairm/devgita/pkg/common"
	macos "github.com/cjairm/devgita/pkg/macos/install/terminal"
)

// Function to upgrade Homebrew
func upgradeHomebrew() error {
	cmd := common.CommandInfo{
		PreExecutionMessage:  "Upgrating Homebrew",
		PostExecutionMessage: "Homebrew upgrated ✔",
		IsSudo:               false,
		Command:              "brew",
		Args:                 []string{"upgrade"},
	}
	err := common.ExecCommand(cmd)
	if err != nil {
		fmt.Println("Please try `brew doctor`. It may fix the issue")
		fmt.Println("Installation stopped.")
		os.Exit(1)
	}
	return nil
}

// Function to install curl
func installCurl() error {
	cmd := common.CommandInfo{
		PreExecutionMessage:  "Installing curl",
		PostExecutionMessage: "curl installed ✔",
		IsSudo:               false,
		Command:              "brew",
		Args:                 []string{"install", "curl"},
	}
	return common.ExecCommand(cmd)
}

// Function to install git
func installGit() error {
	cmd := common.CommandInfo{
		PreExecutionMessage:  "Installing git",
		PostExecutionMessage: "git installed ✔",
		IsSudo:               false,
		Command:              "brew",
		Args:                 []string{"install", "git"},
	}
	return common.ExecCommand(cmd)
}

// Function to install unzip
func installUnzip() error {
	cmd := common.CommandInfo{
		PreExecutionMessage:  "Installing unzip",
		PostExecutionMessage: "unzip installed ✔",
		IsSudo:               false,
		Command:              "brew",
		Args:                 []string{"install", "unzip"},
	}
	return common.ExecCommand(cmd)
}

// Function to run all terminal installers
func RunTerminalInstallers(devgitaPath string) error {
	installFunctions := []func() error{
		upgradeHomebrew,
		installCurl,
		installGit,
		installUnzip,
		func() error {
			return macos.InstallFastFetch(devgitaPath)
		},
		macos.InstallGitHubCli,
		macos.InstallLazyDocker,
		macos.InstallLazyGit,
		func() error {
			return macos.InstallNeovim(devgitaPath)
		},
	}
	for _, installFunc := range installFunctions {
		if err := installFunc(); err != nil {
			fmt.Printf("\033[31mError: %s\033[0m\n", err.Error())
			fmt.Println("Installation stopped.")
			os.Exit(1)
		}
	}
	return nil
}
