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
	"path/filepath"

	"github.com/cjairm/devgita/pkg/common"
	"github.com/manifoldco/promptui"
)

func InstallNeovim(devgitaPath string) error {
	if common.FileAlreadyExist("/usr/local/bin/nvim") &&
		common.FileAlreadyExist("/usr/local/lib/nvim") &&
		common.FileAlreadyExist("/usr/local/share/nvim") || common.IsCommandInstalled("nvim") {
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
		cmd := common.CommandInfo{
			PreExecutionMessage:  "Extracting content",
			PostExecutionMessage: "Content extracted ✔",
			IsSudo:               false,
			Command:              "tar",
			Args:                 []string{"-xf", tarFile, "-C", dest},
		}
		if err := common.ExecCommand(cmd); err != nil {
			return fmt.Errorf("Error extracting content: %v", err)
		}
	}

	// Install Neovim binary
	if !common.FileAlreadyExist("/usr/local/bin/nvim") {
		cmd := common.CommandInfo{
			PreExecutionMessage:  "Installing Neovim binary",
			PostExecutionMessage: "Neovim binary installed ✔",
			IsSudo:               true,
			Command:              "install",
			Args: []string{
				fmt.Sprintf("%s/bin/nvim", extractedFile),
				"/usr/local/bin/nvim",
			},
		}
		if err := common.ExecCommand(cmd); err != nil {
			return fmt.Errorf("Error installing Neovim binary: %v", err)
		}
	}

	// Copy lib to /usr/local
	if !common.FileAlreadyExist("/usr/local/lib/nvim") {
		cmd := common.CommandInfo{
			PreExecutionMessage:  "Copying Neovim /lib",
			PostExecutionMessage: "/lib copied ✔",
			IsSudo:               true,
			Command:              "cp",
			Args: []string{
				"-R",
				fmt.Sprintf("%s/lib", extractedFile),
				"/usr/local/",
			},
		}
		if err := common.ExecCommand(cmd); err != nil {
			return fmt.Errorf("Error copying lib files: %v", err)
		}
	}

	// Copy share directories to /usr/local
	if !common.FileAlreadyExist("/usr/local/share/nvim") {
		cmd := common.CommandInfo{
			PreExecutionMessage:  "Copying Neovim /share",
			PostExecutionMessage: "/share copied ✔",
			IsSudo:               true,
			Command:              "cp",
			Args: []string{
				"-R",
				fmt.Sprintf("%s/share", extractedFile),
				"/usr/local/",
			},
		}
		if err := common.ExecCommand(cmd); err != nil {
			return fmt.Errorf("Error copying share files: %v", err)
		}
	}

	cmd := common.CommandInfo{
		PreExecutionMessage:  "Cleaning up downloaded files",
		PostExecutionMessage: "Downloaded files cleaned up ✔",
		IsSudo:               false,
		Command:              "rm",
		Args: []string{
			"-rf",
			extractedFile,
			tarFile,
		},
	}
	// Clean up the downloaded files
	if err := common.ExecCommand(cmd); err != nil {
		return fmt.Errorf("Error removing temporary files: %v", err)
	}

	// Configure Neovim
	configNeovim(devgitaPath)

	fmt.Println("Neovim installed successfully!")
	return nil
}

func configNeovim(devgitaPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("Error getting home directory: %w", err)
	}

	nvimConfigDir := filepath.Join(homeDir, ".config", "nvim")
	// Only attempt to set configuration if Neovim has never been run
	if common.FileAlreadyExist(nvimConfigDir) {
		fmt.Println("Neovim configuration already exists!")
		return nil
	}

	if err := os.MkdirAll(nvimConfigDir, os.ModePerm); err != nil {
		return fmt.Errorf("Error creating Neovim config directory: %w", err)
	}

	// Define the source directory
	sourceDir := filepath.Join(devgitaPath, "pkg", "configs", "neovim")

	// Copy contents from sourceDir to nvimConfigDir
	if err := common.CopyDirectory(sourceDir, nvimConfigDir); err != nil {
		return fmt.Errorf("Error copying files: %w", err)
	}

	fmt.Println("Neovim configuration set successfully!")

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
