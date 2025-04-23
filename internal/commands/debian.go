package commands

import (
	"fmt"
	"os/exec"

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
	fmt.Println("Executing `InstallPackage` on Debian")
	return nil
}

func (d *DebianCommand) InstallDesktopApp(packageName string) error {
	fmt.Println("Executing `InstallDesktopApp` on Debian")
	return nil
}

func (d *DebianCommand) UpgradePackage(packageName string) error {
	fmt.Println("Executing `UpgradePackage` on Debian")
	return nil
}

func (d *DebianCommand) UpgradePackageManager(verbose bool) error {
	fmt.Println("Executing `UpgradePackageManager` on Debian")
	return nil
}

func (d *DebianCommand) UpdatePackageManager() error {
	fmt.Println("Executing `UpdatePackageManager` on Debian")
	return nil
}

func (d *DebianCommand) MaybeInstallPackageManager() error {
	fmt.Println("Executing `MaybeInstallPackageManager` on Debian")
	return nil
}

func (d *DebianCommand) InstallPackageManager() error {
	fmt.Println("Executing `InstallPackageManager` on Debian")
	return nil
}

func (d *DebianCommand) IsPackageManagerInstalled() bool {
	fmt.Println("Executing `IsPackageManagerInstalled` on Debian")
	return false
}

func (d *DebianCommand) ValidateOSVersion() error {
	fmt.Println("Executing `ValidateOSVersion` on Debian")
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
