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

// Function to run all terminal installers
func RunTerminalInstallers(devgitaPath string) error {
	installFunctions := []func() error{
		upgradeHomebrew,
		// Function to install curl
		func() error {
			return common.InstallOrUpdateBrewPackage("curl")
		},
		// Function to install git
		func() error {
			return common.InstallOrUpdateBrewPackage("git")
		},
		// Function to install unzip
		func() error {
			return common.InstallOrUpdateBrewPackage("unzip")
		},
		func() error {
			return macos.InstallFastFetch(devgitaPath)
		},
		macos.InstallGitHubCli,
		macos.InstallLazyDocker,
		macos.InstallLazyGit,
		func() error {
			return macos.InstallNeovim(devgitaPath)
		},
		func() error {
			return macos.InstallTmux(devgitaPath)
		},
		func() error {
			// installs fzf, ripgrep, bat, eza, zoxide, btop, fd (fd-find), tldr
			packages := []string{"fzf", "ripgrep", "bat", "eza", "zoxide", "btop", "fd", "tldr"}
			for _, pkg := range packages {
				if err := common.InstallOrUpdateBrewPackage(pkg); err != nil {
					return err
				}
			}
			return nil
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
