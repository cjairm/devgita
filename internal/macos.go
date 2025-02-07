package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type MacOSCommand struct{}

func (m *MacOSCommand) MaybeInstallPackage(packageName string, alias ...string) error {
	var isInstalled bool
	var err error
	pkgToInstall := packageName
	if len(alias) > 0 {
		pkgToInstall = alias[0]
	}
	isInstalled, err = isPackageInstalled(pkgToInstall)
	if err != nil {
		return err
	}
	if isInstalled {
		return nil
	}
	return m.InstallPackage(packageName)
}

func (m *MacOSCommand) MaybeInstallDesktopApp(desktopAppName string, alias ...string) error {
	var isInstalled bool
	var err error
	pkgToInstall := desktopAppName
	if len(alias) > 0 {
		pkgToInstall = alias[0]
	}
	isInstalled, err = isDesktopAppInstalled(pkgToInstall)
	if !isInstalled {
		isInstalled, err = desktopApplicationExist(pkgToInstall)
	}
	if err != nil {
		return err
	}
	if isInstalled {
		return nil
	}
	return m.InstallDesktopApp(pkgToInstall)
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

func (m *MacOSCommand) UpgradePackageManager() error {
	cmd := CommandParams{
		PreExecMsg:  "Upgrating Homebrew",
		PostExecMsg: "",
		Verbose:     false,
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

func (m *MacOSCommand) GitCommand(args ...string) error {
	cmd := CommandParams{
		PreExecMsg:  "",
		PostExecMsg: "",
		Verbose:     true,
		IsSudo:      false,
		Command:     "git",
		Args:        args,
	}
	return ExecCommand(cmd)
}

func isPackageInstalled(packageName string) (bool, error) {
	cmd := exec.Command("brew", "list")
	return findPackageInCommandOutput(cmd, packageName)
}

func isDesktopAppInstalled(desktopAppName string) (bool, error) {
	cmd := exec.Command("brew", "list", "--cask")
	return findPackageInCommandOutput(cmd, desktopAppName)
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

func findPackageInCommandOutput(cmd *exec.Cmd, packageName string) (bool, error) {
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("Failed running brew command: %v", err)
	}
	for _, line := range bytes.Split(out.Bytes(), []byte{'\n'}) {
		if string(line) == packageName {
			return true, nil
		}
	}
	return false, nil
}
