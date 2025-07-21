package commands

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/paths"
	"github.com/cjairm/devgita/pkg/utils"
)

type MacOSCommand struct {
	BaseCommand
}

func (m *MacOSCommand) MaybeInstallPackage(packageName string, alias ...string) error {
	return m.MaybeInstall(packageName, alias, m.IsPackageInstalled, m.InstallPackage, nil)
}

func (m *MacOSCommand) MaybeInstallDesktopApp(desktopAppName string, alias ...string) error {
	return m.MaybeInstall(desktopAppName, alias, func(name string) (bool, error) {
		isInstalled, err := m.IsDesktopAppInstalled(name)
		if !isInstalled {
			isInstalled, err = m.IsDesktopAppPresent(paths.UserApplicationsDir, name)
		}
		return isInstalled, err
	}, m.InstallDesktopApp, nil)
}

func (m *MacOSCommand) MaybeInstallFont(
	url, fontFileName string,
	runCache bool,
	alias ...string,
) error {
	return m.MaybeInstall(fontFileName, alias, func(name string) (bool, error) {
		isInstalled, err := m.IsDesktopAppInstalled(name)
		if !isInstalled {
			isInstalled, err = m.IsFontPresent(name)
		}
		return isInstalled, err
	}, m.InstallDesktopApp, nil)
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
	if _, err := m.ExecCommand(cmd); err != nil {
		return fmt.Errorf("failed to install package %s: %w", packageName, err)
	}
	return nil
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
	if _, err := m.ExecCommand(cmd); err != nil {
		return fmt.Errorf("failed to install desktop app %s: %w", packageName, err)
	}
	return nil
}

func (m *MacOSCommand) IsPackageManagerInstalled() bool {
	err := exec.Command("brew", "--version").Run()
	return err == nil
}

func (m *MacOSCommand) MaybeInstallPackageManager() error {
	isInstalled := m.IsPackageManagerInstalled()
	if isInstalled {
		return nil
	}
	return m.InstallPackageManager()
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
	if _, err := m.ExecCommand(cmd); err != nil {
		return fmt.Errorf("failed to install Homebrew: %w", err)
	}
	return nil
}

func (m *MacOSCommand) ValidateOSVersion(verbose bool) error {
	cmd := CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		Verbose:     verbose,
		IsSudo:      false,
		Command:     "sw_vers",
		Args:        []string{"-productVersion"},
	}
	if verbose {
		utils.PrintSecondary("Getting macOS version")
	}
	version, err := m.BaseCommand.ExecCommand(cmd)
	if err != nil {
		return err
	}

	if verbose {
		utils.PrintSecondary("Parsing OS version")
	}
	// Trim whitespace and split the version string
	versionStr := strings.TrimSpace(string(version))
	versionParts := strings.Split(versionStr, ".")
	if len(versionParts) < 2 {
		return fmt.Errorf("invalid macOS version format: %s", versionStr)
	}
	if verbose {
		utils.PrintSecondary("Extracting major and minor version from OS version")
	}
	// Convert the major and minor version to integers
	major, err := strconv.Atoi(versionParts[0])
	if err != nil {
		return fmt.Errorf("invalid major version: %w", err)
	}
	minor, err := strconv.Atoi(versionParts[1])
	if err != nil {
		return fmt.Errorf("invalid minor version: %w", err)
	}
	// NOTE: (11/22/2024) Check if the version is at least 13.0 (macOS Sonoma)
	// Update to the latest version if necessary
	if major < constants.SupportedMacOSVersionNumber ||
		(major == constants.SupportedMacOSVersionNumber && minor < 0) {
		return fmt.Errorf(
			"OS requirement not met\nOS required: macOS %s (%d.0) or higher",
			constants.SupportedMacOSVersionName,
			constants.SupportedMacOSVersionNumber,
		)
	}
	if verbose {
		utils.PrintSecondary(fmt.Sprintf("OS version is supported: macOS %s", versionStr))
	}
	return nil
}

func (m *MacOSCommand) IsPackageInstalled(packageName string) (bool, error) {
	cmd := exec.Command("brew", "list")
	return m.IsPackagePresent(cmd, packageName)
}

func (m *MacOSCommand) IsDesktopAppInstalled(desktopAppName string) (bool, error) {
	cmd := exec.Command("brew", "list", "--cask")
	return m.IsPackagePresent(cmd, desktopAppName)
}
