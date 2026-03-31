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
	"github.com/cjairm/devgita/internal/apps/i3"
	"github.com/cjairm/devgita/internal/apps/raycast"
	"github.com/cjairm/devgita/internal/apps/ulauncher"
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

	// Platform-specific window managers
	if d.Base.Platform.IsMac() {
		err = d.InstallAerospace()
		displayMessage(err, constants.Aerospace)
	} else {
		// Debian/Ubuntu: Install i3 tiling window manager
		err = d.InstallI3()
		displayMessage(err, constants.I3)
	}

	utils.PrintInfo("Installing fonts (if no previously installed)...")
	f := fonts.New()
	f.SoftInstallAll()

	// Install desktop apps (platform-gated where appropriate)
	d.InstallDesktopAppsWithoutConfiguration()

	if d.Base.Platform.IsMac() {
		d.DisplayPrivacyInstructions()
	}

	return nil
}

func (d *Desktop) InstallDesktopAppsWithoutConfiguration() {
	// Cross-platform apps: docker, gimp, brave, flameshot
	crossPlatformApps := []struct {
		name string
		app  interface {
			SoftInstall() error
		}
	}{
		{constants.Docker, docker.New()},
		{constants.Gimp, gimp.New()},
		{constants.Brave, brave.New()},
		{constants.Flameshot, flameshot.New()},
	}

	for _, desktopApp := range crossPlatformApps {
		if err := desktopApp.app.SoftInstall(); err != nil {
			displayMessage(err, desktopApp.name)
			continue
		}
	}

	// Platform-specific launchers
	if d.Base.Platform.IsMac() {
		// macOS: Raycast launcher
		r := raycast.New()
		if err := r.SoftInstall(); err != nil {
			displayMessage(err, constants.Raycast)
		}
	} else {
		// Debian/Ubuntu: Ulauncher
		u := ulauncher.New()
		if err := u.SoftInstall(); err != nil {
			displayMessage(err, constants.Ulauncher)
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

func (d *Desktop) InstallI3() error {
	i := i3.New()
	err := i.SoftInstall()
	if err != nil {
		return err
	}
	err = i.SoftConfigure()
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
