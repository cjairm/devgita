// Package docker provides installation and configuration management for Docker Desktop.
// Docker Desktop is a containerization platform that enables developers to build, ship, and run
// distributed applications using containers. This module follows the standardized devgita app
// interface for consistent lifecycle management.

package docker

import (
	"fmt"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
)

type Docker struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func New() *Docker {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Docker{Cmd: osCmd, Base: baseCmd}
}

func (d *Docker) Install() error {
	return d.Cmd.InstallDesktopApp(constants.Docker)
}

func (d *Docker) SoftInstall() error {
	return d.Cmd.MaybeInstallDesktopApp(constants.Docker)
}

func (d *Docker) ForceInstall() error {
	err := d.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall docker: %w", err)
	}
	return d.Install()
}

func (d *Docker) ForceConfigure() error {
	// Docker Desktop doesn't require separate configuration files
	// Configuration is managed through Docker Desktop GUI or daemon.json
	return nil
}

func (d *Docker) SoftConfigure() error {
	// Docker Desktop doesn't require separate configuration files
	// Configuration is managed through Docker Desktop GUI or daemon.json
	return nil
}

func (d *Docker) Uninstall() error {
	// Docker Desktop requires comprehensive cleanup including:
	// 1. Quit Docker Desktop application
	// 2. Uninstall via Homebrew cask (macOS) or package manager (Linux)
	// 3. Remove binaries: docker, docker-compose, docker-credential-*
	// 4. Remove configuration and data directories
	// 5. Remove shell completions
	//
	// This is a complex operation that requires elevated privileges and
	// interactive confirmation. It should be handled manually by users following
	// Docker's official uninstall documentation.
	//
	// Recommended manual steps for macOS:
	// - Quit Docker Desktop: Right-click Docker icon → "Quit Docker Desktop"
	// - Uninstall cask: brew uninstall --cask docker
	// - Remove binaries: sudo rm -f /usr/local/bin/docker*
	// - Remove containers: sudo rm -rf ~/Library/Containers/com.docker.docker
	// - Remove app support: sudo rm -rf ~/Library/Application\ Support/Docker\ Desktop
	// - Remove Docker config: sudo rm -rf ~/.docker
	// - Remove completions: sudo rm -f /usr/local/etc/bash_completion.d/docker
	//                       sudo rm -f /usr/local/share/zsh/site-functions/_docker
	//                       sudo rm -f /usr/local/share/fish/vendor_completions.d/docker.fish

	// === CUSTOM ===
	// - Quit Docker Desktop: Make sure Docker Desktop is not running. Right-click the Docker icon in the menu bar and select "Quit Docker Desktop."
	// - Open Finder: Navigate to the Applications folder.
	// - Locate Docker: Find the Docker.app application.
	// - Move to Trash: Drag Docker.app to the Trash or right-click and select "Move to Trash."

	// brew uninstall --cask docker && sudo rm -f /usr/local/bin/docker && sudo rm -f /usr/local/bin/docker-compose && sudo rm -f /usr/local/bin/docker-credential-desktop && sudo rm -f /usr/local/bin/docker-credential-ecr-login && sudo rm -f /usr/local/bin/docker-credential-osxkeychain && sudo rm -rf ~/Library/Containers/com.docker.docker && sudo rm -rf ~/Library/Application\ Support/Docker\ Desktop && sudo rm -rf ~/.docker && sudo rm -f /usr/local/bin/hub-tool && sudo rm -f /usr/local/bin/kubectl.docker && sudo rm -f /usr/local/etc/bash_completion.d/docker && sudo rm -f /usr/local/share/zsh/site-functions/_docker && sudo rm -f /usr/local/share/fish/vendor_completions.d/docker.fish
	return fmt.Errorf("uninstall not implemented for docker - requires manual cleanup")
}

func (d *Docker) ExecuteCommand(args ...string) error {
	params := cmd.CommandParams{
		Command: constants.Docker,
		Args:    args,
	}
	_, _, err := d.Base.ExecCommand(params)
	if err != nil {
		return fmt.Errorf("failed to execute docker command: %w", err)
	}
	return nil
}

func (d *Docker) Update() error {
	return fmt.Errorf("update not implemented for docker")
}
