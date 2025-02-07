package commands

import "runtime"

type Command interface {
	MaybeInstallPackage(packageName string, alias ...string) error
	MaybeInstallDesktopApp(desktopAppName string, alias ...string) error
	InstallPackage(packageName string) error
	InstallDesktopApp(packageName string) error
	UpgradePackage(packageName string) error
	UpgradePackageManager() error
	UpdatePackageManager() error

	// Git commands
	GitCommand(args ...string) error
}

func NewCommand() Command {
	switch runtime.GOOS {
	case "darwin":
		return &MacOSCommand{}
	// TODO: Is it possible to detect the distribution of Linux?
	case "linux":
		return &DebianCommand{}
	default:
		panic("unsupported operating system")
	}
}
