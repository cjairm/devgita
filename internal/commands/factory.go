package commands

import "runtime"

type Command interface {
	MaybeInstallPackage(packageName string, alias ...string) error
	MaybeInstallDesktopApp(desktopAppName string, alias ...string) error
	MaybeInstallFont(url, fontName string, runCache bool, alias ...string) error
	InstallPackage(packageName string) error
	InstallDesktopApp(packageName string) error
	ValidateOSVersion() error

	// PackageManager
	MaybeInstallPackageManager() error
	InstallPackageManager() error
	IsPackageManagerInstalled() bool

	// Utils
	IsPackageInstalled(packageName string) (bool, error)
	IsDesktopAppInstalled(desktopAppName string) (bool, error)
}

func NewCommand() Command {
	switch runtime.GOOS {
	case "darwin":
		return &MacOSCommand{
			BaseCommand: *NewBaseCommand(),
		}
	// TODO: Is it possible to detect the distribution of Linux?
	case "linux":
		return &DebianCommand{
			BaseCommand: *NewBaseCommand(),
		}
	default:
		panic("unsupported operating system")
	}
}
