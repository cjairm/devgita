package macos

import (
	"fmt"
	"os"

	"github.com/cjairm/devgita/pkg/common"
)

// Function to upgrade Homebrew
func upgradeHomebrew() error {
	err := common.ExecCommand("Upgrating Homebrew", "Homebrew upgrated ✔", "brew", "upgrade")
	if err != nil {
		fmt.Println("Please try `brew doctor` to fix this issue")
		fmt.Println("Installation stopped.")
		os.Exit(1)
	}
	return nil
}

// Function to install curl
func installCurl() error {
	return common.ExecCommand("Installing curl", "curl installed ✔", "brew", "install", "curl")
}

// Function to install git
func installGit() error {
	return common.ExecCommand("Installing git", "git installed ✔", "brew", "install", "git")
}

// Function to install unzip
func installUnzip() error {
	return common.ExecCommand("Installing unzip", "unzip installed ✔", "brew", "install", "unzip")
}

// Function to run all terminal installers
func RunInstallers() error {
	installFunctions := []func() error{
		upgradeHomebrew,
		installCurl,
		installGit,
		installUnzip,
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