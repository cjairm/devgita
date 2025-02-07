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

func (d *DebianCommand) UpgradePackage(packageName string) error {
	fmt.Println("Executing `UpgradePackage` on Debian")
	return nil
}

func (d *DebianCommand) UpgradePackageManager() error {
	fmt.Println("Executing `UpgradePackageManager` on Debian")
	return nil
}

func (d *DebianCommand) UpdatePackageManager() error {
	fmt.Println("Executing `UpdatePackageManager` on Debian")
	return nil
}

func (d *DebianCommand) GitCommand(args ...string) error {
	fmt.Println("Executing `GitCommand` on Debian")
	return nil
}
