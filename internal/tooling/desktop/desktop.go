package desktop

import (
	"fmt"

	"github.com/cjairm/devgita/internal/apps/aerospace"
	"github.com/cjairm/devgita/internal/apps/alacritty"
	"github.com/cjairm/devgita/internal/apps/brave"
	"github.com/cjairm/devgita/internal/apps/docker"
	"github.com/cjairm/devgita/internal/apps/fonts"
	"github.com/cjairm/devgita/internal/apps/gimp"
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

func (d *Desktop) InstallAndConfigure() error {
	dkr := docker.New()
	displayMessage(dkr.SoftInstall(), "docker")

	err := d.InstallAlacritty()
	displayMessage(err, "alacritty")

	utils.PrintInfo("Installing fonts (if no previously installed)...")
	f := fonts.New()
	f.SoftInstallAll()

	gimp := gimp.New()
	displayMessage(gimp.SoftInstall(), "gimp")

	b := brave.New()
	displayMessage(b.SoftInstall(), "brave")

	utils.PrintInfo("Installing flameshot (if no previously installed)...")
	err = d.InstallFlameshot()
	displayMessage(err, "flameshot")

	utils.PrintInfo("Installing aerospace (if no previously installed)...")
	err = d.InstallAerospace()
	displayMessage(err, "aerospace")

	utils.PrintInfo("Installing raycast (if no previously installed)...")
	err = d.InstallRaycast()
	displayMessage(err, "raycast")

	d.DisplayPrivacyInstructions()

	return nil
}

func (d *Desktop) InstallAlacritty() error {
	a := alacritty.New()
	err := a.SoftInstall()
	if err != nil {
		return err
	}
	err = a.SoftConfigure(alacritty.ConfigureOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (d *Desktop) InstallFlameshot() error {
	return d.Cmd.MaybeInstallDesktopApp("flameshot")
}

func (d *Desktop) InstallAerospace() error {
	a := aerospace.New()
	err := a.SoftInstall()
	if err != nil {
		return err
	}
	err = a.SoftConfigure()
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

func displayMessage(err error, desktopAppName string, displayOnlyErrors ...bool) {
	if err != nil {
		logger.L().Errorw("Error installing ", "desktop_app", desktopAppName, "error", err)
		utils.PrintWarning(
			fmt.Sprintf(
				"Install (%s) errored... To halt the installation, press ctrl+c or use --debug flag to see more details",
				desktopAppName,
			),
		)
	} else {
		if displayOnlyErrors != nil && displayOnlyErrors[0] == true {
			return
		}
		msg := fmt.Sprintf("Installing %s (if no previously installed)...", desktopAppName)
		utils.PrintInfo(msg)
	}
}
