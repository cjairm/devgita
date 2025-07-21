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
	return d.MaybeInstall(packageName, alias, d.IsPackageInstalled, d.InstallPackage, nil)
}

func (d *DebianCommand) MaybeInstallDesktopApp(desktopAppName string, alias ...string) error {
	return d.MaybeInstall(desktopAppName, alias, func(name string) (bool, error) {
		return d.IsDesktopAppInstalled(name)
	}, d.InstallDesktopApp, nil)
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
	})
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
	// APT is preinstalled on Debian/Ubuntu systems.
	return nil
}

func (d *DebianCommand) ValidateOSVersion(verbose bool) error {
	logger.L().Debug("Reading /etc/os-release for Linux version info")
	utils.PrintSecondary("Getting Linux version")

	content, err := os.ReadFile("/etc/os-release")
	if err != nil {
		logger.L().Errorw("failed to read OS release info", "error", err)
		return fmt.Errorf("failed to read OS release info: %w", err)
	}

	logger.L().Debug("Parsing OS version from /etc/os-release")
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

	if name == "" || versionStr == "" {
		err := fmt.Errorf("unable to parse OS version information")
		logger.L().Error(err.Error())
		return err
	}

	logger.L().Debug("Extracting major version from OS version string")
	utils.PrintSecondary("Extracting major and minor version from OS version")

	versionParts := strings.Split(versionStr, ".")
	if len(versionParts) < 1 {
		err := fmt.Errorf("invalid version format: %s", versionStr)
		logger.L().Error(err.Error())
		return err
	}
	major, err := strconv.Atoi(versionParts[0])
	if err != nil {
		logger.L().Errorw("invalid major version", "error", err)
		return fmt.Errorf("invalid major version: %w", err)
	}

	// Check supported versions for Debian and Ubuntu
	if name == "debian" && major < constants.SupportedDebianVersionNumber {
		err := fmt.Errorf(
			"OS requirement not met\nOS required: Debian %s (%d.0) or higher",
			constants.SupportedDebianVersionName,
			constants.SupportedDebianVersionNumber,
		)
		logger.L().Warnw("unsupported Debian version", "version", versionStr)
		return err
	} else if name == "ubuntu" && major < constants.SupportedUbuntuVersionNumber {
		err := fmt.Errorf(
			"OS requirement not met\nOS required: Ubuntu %s (%d.0) or higher",
			constants.SupportedUbuntuVersionName,
			constants.SupportedUbuntuVersionNumber,
		)
		logger.L().Warnw("unsupported Ubuntu version", "version", versionStr)
		return err
	}

	logger.L().Infow("Linux OS version supported", "name", name, "version", versionStr)
	utils.PrintSecondary(fmt.Sprintf("OS version is supported: %s %s", name, versionStr))

	return nil
}

func (d *DebianCommand) IsPackageInstalled(packageName string) (bool, error) {
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
