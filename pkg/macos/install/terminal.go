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
		// To have access to github wihtout passwords
		// see documentation: https://docs.github.com/en/authentication/connecting-to-github-with-ssh/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent
		// ssh-keygen -t rsa -b 4096
		// eval "$(ssh-agent -s)"
		// ssh-add -K ~/.ssh/id_rsa
		// GITHUB STEPS
		// pbcopy < ~/.ssh/id_rsa.pub
		// To test it out! : ssh -T git@github.com
		//
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
		macos.MaybeInstallXcode,
		// installs libs pkg-config, autoconf, bison, rust, openssl, readline, zlib, libyaml, ncurses, libffi, gdbm, jemalloc, vips, mupdf
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
