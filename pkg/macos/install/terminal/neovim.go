// -------------------------
// TODO: Write documentation how to use this
// - Kickstart documentation: https://github.com/nvim-lua/kickstart.nvim?tab=readme-ov-file
// - Personal configuration: https://github.com/cjairm/devenv/blob/main/nvim/init.lua
// - Releases: https://github.com/neovim/neovim/releases
//
// NOTE: Is it possible to install different themes?
// If so, see more here: https://linovox.com/the-best-color-schemes-for-neovim-nvim/
// -------------------------

package macos

import (
	"fmt"
	"os"

	"github.com/cjairm/devgita/pkg/common"
	"github.com/manifoldco/promptui"
)

func InstallNeovim() error {
	if common.FileAlreadyExist("/usr/local/bin/nvim") &&
		common.FileAlreadyExist("/usr/local/lib/nvim") &&
		common.FileAlreadyExist("/usr/local/share/nvim") {
		fmt.Println("Neovim is already installed!")
		return nil
	}

	dest := "/tmp"
	fileName := ""
	downloadURL := ""

	architecture := selectArchitecture()
	if architecture == "arm" {
		downloadURL = "https://github.com/neovim/neovim/releases/latest/download/nvim-macos-arm64.tar.gz"
		fileName = "nvim-macos-arm64"
	} else if architecture == "intel" {
		downloadURL = "https://github.com/neovim/neovim/releases/latest/download/nvim-macos-x86_64.tar.gz"
		fileName = "nvim-macos-x86_64"
	} else {
		fmt.Println("\033[31mError: Error selecting languages.\033[0m")
		fmt.Println("Installation stopped.")
		os.Exit(1)
	}
	tarFile := fmt.Sprintf("%s/%s.tar.gz", dest, fileName)
	extractedFile := fmt.Sprintf("%s/%s", dest, fileName)

	// Download the Neovim
	if common.FileAlreadyExist(tarFile) {
		fmt.Printf("File already downloaded: %s\n\n", tarFile)
	} else {
		fmt.Printf("Downloading Neovim...\n\n")
		if err := common.DownloadFile(downloadURL, tarFile); err != nil {
			return fmt.Errorf("Error downloading file: %v", err)
		}
	}

	// Extract the tar content
	if !common.FileAlreadyExist(extractedFile) {
		if err := common.ExecCommand("Extracting content", "Content extracted ✔", "tar", "-xf", tarFile, "-C", dest); err != nil {
			return fmt.Errorf("Error extracting content: %v", err)
		}
	}

	// Install Neovim binary
	if !common.FileAlreadyExist("/usr/local/bin/nvim") {
		if err := common.ExecCommand(
			"Installing Neovim binary",
			"Neovim binary installed ✔",
			"sudo",
			"install",
			fmt.Sprintf("%s/bin/nvim", extractedFile),
			"/usr/local/bin/nvim",
		); err != nil {
			return fmt.Errorf("Error installing Neovim binary: %v", err)
		}
	}

	// Copy lib to /usr/local
	if !common.FileAlreadyExist("/usr/local/lib/nvim") {
		if err := common.ExecCommand(
			"Copying Neovim lib",
			"lib copied ✔",
			"sudo",
			"cp",
			"-R",
			fmt.Sprintf("%s/lib", extractedFile),
			"/usr/local/",
		); err != nil {
			return fmt.Errorf("Error copying lib files: %v", err)
		}
	}

	// Copy share directories to /usr/local
	if !common.FileAlreadyExist("/usr/local/share/nvim") {
		if err := common.ExecCommand(
			"Copying Neovim share",
			"share copied ✔",
			"sudo",
			"cp",
			"-R",
			fmt.Sprintf("%s/share", extractedFile),
			"/usr/local/",
		); err != nil {
			return fmt.Errorf("Error copying share files: %v", err)
		}
	}

	// Clean up the downloaded files
	if err := common.ExecCommand(
		"Cleaning up downloaded files",
		"Downloaded files cleaned up ✔",
		"rm",
		"-rf",
		extractedFile,
		tarFile,
	); err != nil {
		return fmt.Errorf("Error removing temporary files: %v", err)
	}

	fmt.Println("Neovim installed successfully!")
	return nil
}

func selectArchitecture() string {
	options := []string{"arm", "intel"}
	prompt := promptui.Select{
		Label: "Select architecture (arm/intel)",
		Items: options,
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("\033[31mError: Error selecting architecture.\033[0m")
		fmt.Println("Installation stopped.")
		os.Exit(1)
	}
	return result
}
