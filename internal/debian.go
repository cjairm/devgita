package commands

import (
	"fmt"
)

type DebianCommand struct{}

func (d *DebianCommand) MaybeInstallPackage(packageName string, alias ...string) error {
	fmt.Println("Executing `MaybeInstallPackage` on Debian")
	return nil
}

func (d *DebianCommand) MaybeInstallDesktopApp(desktopAppName string, alias ...string) error {
	fmt.Println("Executing `MaybeInstallDesktopApp` on Debian")
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

func (d *DebianCommand) IsMac() bool {
	return false
}

func (d *DebianCommand) IsLinux() bool {
	return true
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
