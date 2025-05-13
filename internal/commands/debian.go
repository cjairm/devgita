package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/paths"
)

type DebianCommand struct {
	BaseCommand
}

func (d *DebianCommand) MaybeInstallPackage(packageName string, alias ...string) error {
	fmt.Println("Executing `MaybeInstallPackage` on Debian")
	return nil
}

func (d *DebianCommand) MaybeInstallDesktopApp(desktopAppName string, alias ...string) error {
	fmt.Println("Executing `MaybeInstallDesktopApp` on Debian")
	return nil
}

func (d *DebianCommand) MaybeInstallFont(desktopAppName string, alias ...string) error {
	fmt.Println("Executing `MaybeInstallFont` on Debian")
	return nil
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

func (d *DebianCommand) ValidateOSVersion() error {
	content, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return fmt.Errorf("failed to read OS release info: %w", err)
	}

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
		return fmt.Errorf("unable to parse OS version information")
	}

	versionParts := strings.Split(versionStr, ".")
	if len(versionParts) < 1 {
		return fmt.Errorf("invalid version format")
	}
	major, err := strconv.Atoi(versionParts[0])
	if err != nil {
		return fmt.Errorf("invalid major version: %w", err)
	}

	if name == "debian" && major < constants.SupportedDebianVersionNumber {
		return fmt.Errorf(
			"OS requirement not met\nOS required: Debian %s (%d.0) or higher",
			constants.SupportedDebianVersionName,
			constants.SupportedDebianVersionNumber,
		)
	} else if name == "ubuntu" && major < constants.SupportedUbuntuVersionNumber {
		return fmt.Errorf(
			"OS requirement not met\nOS required: Ubuntu %s (%d.0) or higher",
			constants.SupportedUbuntuVersionName,
			constants.SupportedUbuntuVersionNumber,
		)
	}

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
		Verbose:     false,
		IsSudo:      true,
		Command:     "apt",
		Args:        []string{"install", "-y", packageName},
	}
	return d.ExecCommand(cmd)
}

func (d *DebianCommand) installWithSnap(packageName string) error {
	cmd := CommandParams{
		PreExecMsg:  fmt.Sprintf("Installing %s...", strings.ToLower(packageName)),
		PostExecMsg: "",
		Verbose:     false,
		IsSudo:      true,
		Command:     "apt",
		Args:        []string{"install", "-y", packageName},
	}
	return d.ExecCommand(cmd)
}
