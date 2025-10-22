package desktop

import (
	"github.com/cjairm/devgita/internal/apps/aerospace"
	"github.com/cjairm/devgita/internal/apps/alacritty"
	"github.com/cjairm/devgita/internal/apps/fonts"
	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/logger"
	"github.com/cjairm/devgita/pkg/promptui"
	"github.com/cjairm/devgita/pkg/utils"
)

type Desktop struct {
	Cmd cmd.Command
}

func New() *Desktop {
	osCmd := cmd.NewCommand()
	return &Desktop{Cmd: osCmd}
}

func (d *Desktop) InstallAll() error {
	utils.PrintInfo("Installing docker (if no previously installed)...")
	err := d.InstallDocker()
	ifErrorDisplayMessage(err, "docker")

	utils.PrintInfo("Installing and setting up alacritty (if no previous configuration)...")
	err = d.InstallAlacritty()
	ifErrorDisplayMessage(err, "alacritty")

	utils.PrintInfo("Installing fonts (if no previously installed)...")
	f := fonts.New()
	f.MaybeInstallAll()

	utils.PrintInfo("Installing gimp (if no previously installed)...")
	err = d.InstallGimp()
	ifErrorDisplayMessage(err, "gimp")

	utils.PrintInfo("Installing brave (if no previously installed)...")
	err = d.InstallBrave()
	ifErrorDisplayMessage(err, "brave")

	utils.PrintInfo("Installing flameshot (if no previously installed)...")
	err = d.InstallFlameshot()
	ifErrorDisplayMessage(err, "flameshot")

	utils.PrintInfo("Installing aerospace (if no previously installed)...")
	err = d.InstallAerospace()
	ifErrorDisplayMessage(err, "aerospace")

	utils.PrintInfo("Installing raycast (if no previously installed)...")
	err = d.InstallRaycast()
	ifErrorDisplayMessage(err, "raycast")

	d.DisplayPrivacyInstructions()

	return nil
}

func (d *Desktop) InstallDocker() error {
	// - Quit Docker Desktop: Make sure Docker Desktop is not running. Right-click the Docker icon in the menu bar and select "Quit Docker Desktop."
	// - Open Finder: Navigate to the Applications folder.
	// - Locate Docker: Find the Docker.app application.
	// - Move to Trash: Drag Docker.app to the Trash or right-click and select "Move to Trash."

	// brew uninstall --cask docker && sudo rm -f /usr/local/bin/docker && sudo rm -f /usr/local/bin/docker-compose && sudo rm -f /usr/local/bin/docker-credential-desktop && sudo rm -f /usr/local/bin/docker-credential-ecr-login && sudo rm -f /usr/local/bin/docker-credential-osxkeychain && sudo rm -rf ~/Library/Containers/com.docker.docker && sudo rm -rf ~/Library/Application\ Support/Docker\ Desktop && sudo rm -rf ~/.docker && sudo rm -f /usr/local/bin/hub-tool && sudo rm -f /usr/local/bin/kubectl.docker && sudo rm -f /usr/local/etc/bash_completion.d/docker && sudo rm -f /usr/local/share/zsh/site-functions/_docker && sudo rm -f /usr/local/share/fish/vendor_completions.d/docker.fish
	return d.Cmd.MaybeInstallDesktopApp("docker")
}

func (d *Desktop) InstallAlacritty() error {
	a := alacritty.New()
	err := a.MaybeInstall()
	if err != nil {
		return err
	}
	err = a.MaybeSetupApp()
	if err != nil {
		return err
	}
	err = a.MaybeSetupFont()
	if err != nil {
		return err
	}
	err = a.MaybeSetupTheme()
	if err != nil {
		return err
	}
	err = a.UpdateConfigFilesWithCurrentHomeDir()
	if err != nil {
		return err
	}
	return nil
}

func (d *Desktop) InstallGimp() error {
	return d.Cmd.MaybeInstallDesktopApp("gimp")
}

func (d *Desktop) InstallBrave() error {
	return d.Cmd.MaybeInstallDesktopApp("brave-browser", "brave")
}

func (d *Desktop) InstallFlameshot() error {
	return d.Cmd.MaybeInstallDesktopApp("flameshot")
}

func (d *Desktop) InstallAerospace() error {
	a := aerospace.New()
	err := a.MaybeInstall()
	if err != nil {
		return err
	}
	err = a.MaybeSetup()
	if err != nil {
		return err
	}
	return nil
}

func (d *Desktop) InstallRaycast() error {
	return d.Cmd.MaybeInstallDesktopApp("raycast")
}

func (d *Desktop) DisplayPrivacyInstructions() error {
	instructions := `
1. Open System Preferences.
2. Go to Security & Privacy.
3. Click on the Privacy tab.
4. Select Full Disk Access from the left sidebar.
5. Click the lock icon in the bottom left corner to make changes and enter your password.
`
	return promptui.DisplayInstructions(
		"To enable full functionality of the applications, please do the following",
		instructions,
		false,
	)
}

func ifErrorDisplayMessage(err error, packageName string) {
	if err != nil {
		logger.L().Error("Error installing "+packageName, "error", err)
		utils.PrintError("Error installing " + packageName + ": ")
		utils.PrintWarning("Proceeding... (To halt the installation, press ctrl+c)")
	}
}
