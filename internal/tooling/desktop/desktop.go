package desktop

import (
	"fmt"

	"github.com/cjairm/devgita/internal/apps/aerospace"
	"github.com/cjairm/devgita/internal/apps/alacritty"
	"github.com/cjairm/devgita/internal/apps/brave"
	"github.com/cjairm/devgita/internal/apps/docker"
	"github.com/cjairm/devgita/internal/apps/flameshot"
	"github.com/cjairm/devgita/internal/apps/fonts"
	"github.com/cjairm/devgita/internal/apps/gimp"
	"github.com/cjairm/devgita/internal/apps/raycast"
	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/logger"
	"github.com/cjairm/devgita/pkg/promptui"
	"github.com/cjairm/devgita/pkg/utils"
)

type Desktop struct {
	Cmd  cmd.Command
	Base cmd.BaseCommand
}

func New() *Desktop {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Desktop{Cmd: osCmd, Base: *baseCmd}
}

func (d *Desktop) InstallAndConfigure() error {
	err := d.InstallAlacritty()
	displayMessage(err, constants.Alacritty)

	err = d.InstallAerospace()
	displayMessage(err, constants.Aerospace)

	utils.PrintInfo("Installing fonts (if no previously installed)...")
	f := fonts.New()
	f.SoftInstallAll()

	if d.Base.Platform.IsMac() {
		d.DisplayPrivacyInstructions()
	}

	return nil
}

func (d *Desktop) InstallDesktopAppsWithoutConfiguration() {
	// should install docker, gimp, brave, flameshot, raycast
	desktopApps := []struct {
		name string
		app  interface {
			SoftInstall() error
		}
	}{
		{constants.Docker, docker.New()},
		{constants.Gimp, gimp.New()},
		{constants.Brave, brave.New()},
		{constants.Flameshot, flameshot.New()},
		{constants.Raycast, raycast.New()},
	}
	for _, desktopApp := range desktopApps {
		if err := desktopApp.app.SoftInstall(); err != nil {
			displayMessage(err, desktopApp.name)
			continue
		}
	}
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
