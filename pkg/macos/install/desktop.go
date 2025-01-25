package macos

import (
	"fmt"
	"os"

	"github.com/cjairm/devgita/pkg/common"
	macosDesktop "github.com/cjairm/devgita/pkg/macos/install/desktop"
)

// Function to run all terminal installers
func RunDesktopInstallers(devgitaPath string) error {
	installFunctions := []func() error{
		// Installs Docker
		func() error {
			// - Quit Docker Desktop: Make sure Docker Desktop is not running. Right-click the Docker icon in the menu bar and select "Quit Docker Desktop."
			// - Open Finder: Navigate to the Applications folder.
			// - Locate Docker: Find the Docker.app application.
			// - Move to Trash: Drag Docker.app to the Trash or right-click and select "Move to Trash."

			// 1. Open System Preferences.
			// 2. Go to Security & Privacy.
			// 3. Click on the Privacy tab.
			// 4. Select Full Disk Access from the left sidebar.
			// 5. Click the lock icon in the bottom left corner to make changes and enter your password.

			// brew uninstall --cask docker && sudo rm -f /usr/local/bin/docker && sudo rm -f /usr/local/bin/docker-compose && sudo rm -f /usr/local/bin/docker-credential-desktop && sudo rm -f /usr/local/bin/docker-credential-ecr-login && sudo rm -f /usr/local/bin/docker-credential-osxkeychain && sudo rm -rf ~/Library/Containers/com.docker.docker && sudo rm -rf ~/Library/Application\ Support/Docker\ Desktop && sudo rm -rf ~/.docker && sudo rm -f /usr/local/bin/hub-tool && sudo rm -f /usr/local/bin/kubectl.docker && sudo rm -f /usr/local/etc/bash_completion.d/docker && sudo rm -f /usr/local/share/zsh/site-functions/_docker && sudo rm -f /usr/local/share/fish/vendor_completions.d/docker.fish
			return common.MaybeInstallBrewCask("docker")
		},
		// Installs Alacritty
		func() error {
			return macosDesktop.InstallAlacritty(devgitaPath)
		},
		// Installs nerd fonts
		func() error {
			return common.MaybeInstallBrewCask("font-hack-nerd-font")
		},
		// Installs (more) nerd fonts
		func() error {
			return common.MaybeInstallBrewCask("font-meslo-lg-nerd-font")
		},
		// Installs GIMP
		func() error {
			return common.MaybeInstallBrewCask("gimp")
		},
		// Installs Brave
		func() error {
			return common.MaybeInstallBrewCask("brave-browser")
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
