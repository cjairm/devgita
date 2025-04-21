package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	return d.FindPackageInCommandOutput(cmd, packageName)
}

func (d *DebianCommand) IsDesktopAppInstalled(desktopAppName string) (bool, error) {
	for _, dirType := range []string{"user", "system"} {
		appDir, err := getLinuxApplicationsDir(dirType)
		if err != nil {
			return false, err
		}
		isInstalled, err := d.CheckFileExistsInDirectory(appDir, desktopAppName)
		if isInstalled {
			return true, nil
		}
	}
	return false, nil
}

func getLinuxApplicationsDir(t string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if t == "user" {
		return filepath.Join(homeDir, ".local", "share", "applications"), nil
	} else if t == "system" {
		return filepath.Join("usr", "share", "applications"), nil
	} else {
		return "", fmt.Errorf("unsupported argument")
	}
}
