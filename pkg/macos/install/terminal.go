package macos

import (
	"fmt"
	"os"

	"github.com/cjairm/devgita/pkg/common"
	macos "github.com/cjairm/devgita/pkg/macos/install/terminal"
)

// Function to run all terminal installers
func RunTerminalInstallers(devgitaPath string) error {
	installFunctions := []func() error{
		// Function to upgrade Homebrew
		func() error {
			err := common.BrewGlobalUpgrade()
			if err != nil {
				fmt.Println("Please try `brew doctor`. It may fix the issue")
				fmt.Println("Installation stopped.")
				os.Exit(1)
			}
			return nil
		},
		// Function to install curl
		func() error {
			return common.InstallOrUpdateBrewPackage("curl")
		},
		// Function to install unzip
		func() error {
			return common.InstallOrUpdateBrewPackage("unzip")
		},
		// Installs and configures fast-fetch
		func() error {
			return macos.InstallFastFetch(devgitaPath)
		},
		// Function to install GitHub CLI
		func() error {
			return common.InstallOrUpdateBrewPackage("gh")
		},
		// Function to install Lazy Docker
		func() error {
			return common.InstallOrUpdateBrewPackage(
				"jesseduffield/lazydocker/lazydocker",
				"lazydocker",
			)
		},
		// Function to install Lazy Git
		func() error {
			return common.InstallOrUpdateBrewPackage("lazygit")
		},
		// Installs and configures neovim
		func() error {
			return macos.InstallNeovim(devgitaPath)
		},
		// Installs and configures tmux
		func() error {
			return macos.InstallTmux(devgitaPath)
		},
		// installs fzf, ripgrep, bat, eza, zoxide, btop, fd-find, tldr
		func() error {
			packages := []string{"fzf", "ripgrep", "bat", "eza", "zoxide", "btop", "fd", "tldr"}
			for _, pkg := range packages {
				if err := common.InstallOrUpdateBrewPackage(pkg); err != nil {
					return err
				}
			}
			return nil
		},
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
		macos.MaybeInstallXcode,
		// installs libs pkg-config, autoconf, bison, rust, openssl, readline, zlib, libyaml, ncurses, libffi, gdbm, jemalloc, vips, imagemagick, mupdf
		func() error {
			libs := []string{
				"pkg-config",
				"autoconf",
				"bison",
				"rust",
				"openssl",
				"readline",
				"zlib",
				"libyaml",
				"ncurses",
				"libffi",
				"gdbm",
				"jemalloc",
				"vips",
				"imagemagick",
				"mupdf",
			}
			for _, lib := range libs {
				if err := common.InstallOrUpdateBrewPackage(lib); err != nil {
					return err
				}
			}
			return nil
		},
		// Function to install Mise
		// TODO: Compete documentation: https://mise.jdx.dev/
		func() error {
			return common.InstallOrUpdateBrewPackage("mise")
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
