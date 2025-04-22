package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cjairm/devgita/pkg/utils"
)

type MacOSCommand struct {
	BaseCommand
}

func (m *MacOSCommand) MaybeInstallPackage(packageName string, alias ...string) error {
	return maybeInstall(packageName, alias, m.IsPackageInstalled, m.InstallPackage)
}

func (m *MacOSCommand) MaybeInstallDesktopApp(desktopAppName string, alias ...string) error {
	return maybeInstall(desktopAppName, alias, func(name string) (bool, error) {
		isInstalled, err := m.IsDesktopAppInstalled(name)
		if !isInstalled {
			isInstalled, err = desktopApplicationExist(name)
		}
		return isInstalled, err
	}, m.InstallDesktopApp)
}

func (m *MacOSCommand) MaybeInstallFont(fontName string, alias ...string) error {
	return maybeInstall(fontName, alias, func(name string) (bool, error) {
		isInstalled, err := m.IsDesktopAppInstalled(name)
		if !isInstalled {
			isInstalled, err = fontExist(name)
		}
		return isInstalled, err
	}, m.InstallDesktopApp)
}

func (m *MacOSCommand) InstallPackage(packageName string) error {
	cmd := CommandParams{
		PreExecMsg:  fmt.Sprintf("Installing %s...", strings.ToLower(packageName)),
		PostExecMsg: "",
		Verbose:     false,
		IsSudo:      false,
		Command:     "brew",
		Args:        []string{"install", packageName},
	}
	return ExecCommand(cmd)
}

func (m *MacOSCommand) InstallDesktopApp(packageName string) error {
	cmd := CommandParams{
		PreExecMsg:  fmt.Sprintf("Installing %s...", strings.ToLower(packageName)),
		PostExecMsg: "",
		Verbose:     false,
		IsSudo:      false,
		Command:     "brew",
		Args:        []string{"install", "--cask", packageName},
	}
	return ExecCommand(cmd)
}

func (m *MacOSCommand) UpgradePackage(packageName string) error {
	cmd := CommandParams{
		PreExecMsg:  fmt.Sprintf("Upgrading %s...", strings.ToLower(packageName)),
		PostExecMsg: "",
		Verbose:     false,
		IsSudo:      false,
		Command:     "brew",
		Args:        []string{"upgrade", packageName},
	}
	return ExecCommand(cmd)
}

func (m *MacOSCommand) UpgradePackageManager(verbose bool) error {
	cmd := CommandParams{
		PreExecMsg:  "Upgrating Homebrew",
		PostExecMsg: "",
		Verbose:     verbose,
		IsSudo:      false,
		Command:     "brew",
		Args:        []string{"upgrade"},
	}
	return ExecCommand(cmd)
}

func (m *MacOSCommand) UpdatePackageManager() error {
	cmd := CommandParams{
		PreExecMsg:  "Updating Homebrew",
		PostExecMsg: "",
		Verbose:     false,
		IsSudo:      false,
		Command:     "brew",
		Args:        []string{"update"},
	}
	return ExecCommand(cmd)
}

func (m *MacOSCommand) MaybeInstallPackageManager() error {
	isInstalled := m.IsPackageManagerInstalled()
	if isInstalled {
		return nil
	}
	return m.InstallPackageManager()
}

func (m *MacOSCommand) IsPackageManagerInstalled() bool {
	err := exec.Command("brew", "--version").Run()
	return err == nil
}

func (m *MacOSCommand) InstallPackageManager() error {
	cmd := CommandParams{
		PreExecMsg:  "Installing Homebrew",
		PostExecMsg: "Homebrew installed ✔",
		Verbose:     false,
		IsSudo:      false,
		Command:     "/bin/bash",
		Args: []string{
			"-c",
			"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)",
		},
	}
	return ExecCommand(cmd)
}

func (m *MacOSCommand) ValidateOSVersion() error {
	// Get the macOS version
	version, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		return err
	}
	// Trim whitespace and split the version string
	versionStr := strings.TrimSpace(string(version))
	versionParts := strings.Split(versionStr, ".")
	if len(versionParts) < 2 {
		return fmt.Errorf("Unable to parse macOS version")
	}
	// Convert the major and minor version to integers
	major, err := strconv.Atoi(versionParts[0])
	if err != nil {
		return err
	}
	minor, err := strconv.Atoi(versionParts[1])
	if err != nil {
		return err
	}
	// NOTE: (11/22/2024) Check if the version is at least 13.0 (macOS Sonoma)
	// Update to the latest version if necessary
	if major < utils.SupportedMacOSVersionNumber ||
		(major == utils.SupportedMacOSVersionNumber && minor < 0) {
		return fmt.Errorf(
			"OS requirement not met\nOS required: macOS %s (%d.0) or higher",
			utils.SupportedMacOSVersionName,
			utils.SupportedMacOSVersionNumber,
		)
	}
	return nil
}

func (m *MacOSCommand) IsPackageInstalled(packageName string) (bool, error) {
	cmd := exec.Command("brew", "list")
	return m.FindPackageInCommandOutput(cmd, packageName)
}

func (m *MacOSCommand) IsDesktopAppInstalled(desktopAppName string) (bool, error) {
	cmd := exec.Command("brew", "list", "--cask")
	return m.FindPackageInCommandOutput(cmd, desktopAppName)
}

func desktopApplicationExist(appName string) (bool, error) {
	applicationsPath := "/Applications"
	return checkFileExistsInDirectory(applicationsPath, appName)
}

func fontExist(appName string) (bool, error) {
	fontsPath := "~/Library/Fonts/"
	return checkFileExistsInDirectory(fontsPath, appName)
}

func checkFileExistsInDirectory(dirPath, name string) (bool, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return false, fmt.Errorf("Failed to read directory: %v", err)
	}
	for _, file := range files {
		lowerCaseName := strings.ToLower(file.Name())
		if strings.Contains(lowerCaseName, name) {
			return true, nil
		}
	}
	return false, nil
}

func maybeInstall(
	itemName string,
	alias []string,
	checkInstalled func(string) (bool, error),
	installFunc func(string) error,
) error {
	var isInstalled bool
	var err error
	pkgToInstall := itemName
	if len(alias) > 0 {
		pkgToInstall = alias[0]
	}
	isInstalled, err = checkInstalled(pkgToInstall)
	if err != nil {
		return err
	}
	if isInstalled {
		return nil
	}
	return installFunc(pkgToInstall)
}
