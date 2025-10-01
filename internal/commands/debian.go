package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cjairm/devgita/logger"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/paths"
	"github.com/cjairm/devgita/pkg/utils"
)

type DebianCommand struct {
	BaseCommand
}

func (d *DebianCommand) MaybeInstallPackage(packageName string, alias ...string) error {
	return d.MaybeInstall(
		packageName,
		alias,
		d.IsPackageInstalled,
		d.InstallPackage,
		nil,
		"package",
	)
}

func (d *DebianCommand) MaybeInstallDesktopApp(desktopAppName string, alias ...string) error {
	return d.MaybeInstall(desktopAppName, alias, func(name string) (bool, error) {
		return d.IsDesktopAppInstalled(name)
	}, d.InstallDesktopApp, nil, "desktop_app")
}

func (d *DebianCommand) MaybeInstallFont(
	url, fontFileName string,
	runCache bool,
	alias ...string,
) error {
	return d.MaybeInstall(fontFileName, alias, func(name string) (bool, error) {
		isInstalled, err := d.IsFontPresent(name)
		return isInstalled, err
	}, d.InstallPackage, func(name string) error {
		return d.InstallFontFromURL(url, name, runCache)
	}, "font")
}

func (d *DebianCommand) InstallPackage(packageName string) error {
	return d.installWithApt(packageName)
}

func (d *DebianCommand) InstallDesktopApp(packageName string) error {
	err := d.installWithApt(packageName)
	if err == nil {
		return nil
	}
	return d.installWithSnap(packageName)
}

func (d *DebianCommand) IsPackageManagerInstalled() bool {
	logger.L().Debug("executing: which apt")
	_, err := exec.LookPath("apt")
	if err != nil {
		return false
	}
	return true
}

func (d *DebianCommand) MaybeInstallPackageManager() error {
	isInstalled := d.IsPackageManagerInstalled()
	if isInstalled {
		return nil
	}
	return d.InstallPackageManager()
}

func (d *DebianCommand) InstallPackageManager() error {
	logger.L().Info("APT is preinstalled on Debian/Ubuntu systems")
	// APT is preinstalled on Debian/Ubuntu systems.
	return nil
}

func (d *DebianCommand) ValidateOSVersion() error {
	utils.PrintSecondary("Getting Linux version")

	content, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return fmt.Errorf("failed to read OS release info: %w", err)
	}

	utils.PrintSecondary("Parsing OS version")

	var name, versionStr string
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			name = strings.Trim(strings.SplitN(line, "=", 2)[1], `"`)
		} else if strings.HasPrefix(line, "VERSION_ID=") {
			versionStr = strings.Trim(strings.SplitN(line, "=", 2)[1], `"`)
		}
	}
	logger.L().Debugw("OS version info", "name", name, "version", versionStr)

	if name == "" || versionStr == "" {
		err := fmt.Errorf("unable to parse OS version information")
		return err
	}

	utils.PrintSecondary("Extracting major and minor version")
	versionParts := strings.Split(versionStr, ".")
	if len(versionParts) < 1 {
		err := fmt.Errorf("invalid version format: %s", versionStr)
		return err
	}

	major, err := strconv.Atoi(versionParts[0])
	if err != nil {
		return fmt.Errorf("invalid major version: %w", err)
	}
	logger.L().Debugw("OS major version", "os", name, "major_version", major)

	// Check supported versions for Debian and Ubuntu
	switch name {
	case "debian":
		logger.L().Debugw("supported_debian_version", constants.SupportedDebianVersionNumber)
		if major < constants.SupportedDebianVersionNumber {
			err := fmt.Errorf("OS requirement not met\nOS required: Debian %s (%d.0) or higher",
				constants.SupportedDebianVersionName,
				constants.SupportedDebianVersionNumber,
			)
			return err
		}
	case "ubuntu":
		logger.L().Debugw("supported_ubuntu_version", constants.SupportedUbuntuVersionNumber)
		if major < constants.SupportedUbuntuVersionNumber {
			err := fmt.Errorf("OS requirement not met\nOS required: Ubuntu %s (%d.0) or higher",
				constants.SupportedUbuntuVersionName,
				constants.SupportedUbuntuVersionNumber,
			)
			return err
		}
	}

	utils.PrintSecondary(fmt.Sprintf("✅ %s %s is supported", name, versionStr))
	return nil
}

func (d *DebianCommand) IsPackageInstalled(packageName string) (bool, error) {
	logger.L().Debug("executing: dpkg -l")
	cmd := exec.Command("dpkg", "-l")
	return d.IsPackagePresent(cmd, packageName)
}

func (d *DebianCommand) IsDesktopAppInstalled(appName string) (bool, error) {
	for _, applicationsDir := range []string{paths.UserApplicationsDir, paths.SystemApplicationsDir} {
		isInstalled, err := d.IsDesktopAppPresent(applicationsDir, appName)
		if err != nil {
			return false, err
		}
		if isInstalled {
			return true, nil
		}
	}
	return false, nil
}

func (d *DebianCommand) installWithApt(packageName string) error {
	logger.L().Debug(fmt.Sprintf("executing: apt install -y %s", packageName))
	cmd := CommandParams{
		PreExecMsg:  fmt.Sprintf("Installing %s...", strings.ToLower(packageName)),
		PostExecMsg: "",
		IsSudo:      true,
		Command:     "apt",
		Args:        []string{"install", "-y", packageName},
	}
	if _, err := d.ExecCommand(cmd); err != nil {
		return fmt.Errorf("failed to install %s: %w", packageName, err)
	}
	return nil
}

func (d *DebianCommand) installWithSnap(packageName string) error {
	logger.L().Debug(fmt.Sprintf("executing: snap install %s", packageName))
	cmd := CommandParams{
		PreExecMsg:  fmt.Sprintf("Installing %s via Snap...", strings.ToLower(packageName)),
		PostExecMsg: "",
		IsSudo:      true,
		Command:     "snap",
		Args:        []string{"install", packageName},
	}
	if _, err := d.ExecCommand(cmd); err != nil {
		return fmt.Errorf("failed to install %s via Snap: %w", packageName, err)
	}
	return nil
}
