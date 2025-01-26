package common

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func InstallOrUpdateBrewPackage(packageName string, alias ...string) error {
	var isInstalled bool
	var err error
	if len(alias) > 0 {
		isInstalled, err = isBrewPackageInstalled(alias[0])
	} else {
		isInstalled, err = isBrewPackageInstalled(packageName)
	}
	if err != nil {
		return err
	}

	if isInstalled {
		if err := BrewUpgrade(packageName); err != nil {
			return err
		}
	} else {
		if err := BrewInstall(packageName); err != nil {
			return err
		}
	}
	return nil
}

// TODO: Reuse logic from `InstallOrUpdateBrewPackage`
func MaybeInstallBrewCask(packageName string, alias ...string) error {
	var isInstalled bool
	var err error
	if len(alias) > 0 {
		isInstalled, err = isBrewCaskPackageInstalled(alias[0])
		if !isInstalled {
			isInstalled, err = desktopApplicationExist(alias[0])
		}
	} else {
		isInstalled, err = isBrewCaskPackageInstalled(packageName)
		if !isInstalled {
			isInstalled, err = desktopApplicationExist(packageName)
		}
	}
	if err != nil {
		return err
	}
	if isInstalled {
		fmt.Printf("%s is already installed\n\n", packageName)
		return nil
	} else {
		return BrewInstallCask(packageName)
	}

}

func isBrewPackageInstalled(packageName string) (bool, error) {
	cmd := exec.Command("brew", "list")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("error running brew list: %v", err)
	}
	for _, line := range bytes.Split(out.Bytes(), []byte{'\n'}) {
		if string(line) == packageName {
			return true, nil
		}
	}
	return false, nil
}

// TODO: Reuse logic from `isBrewCaskPackageInstalled`
func isBrewCaskPackageInstalled(packageName string) (bool, error) {
	cmd := exec.Command("brew", "list", "--cask")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("error running brew list: %v", err)
	}
	for _, line := range bytes.Split(out.Bytes(), []byte{'\n'}) {
		if string(line) == packageName {
			return true, nil
		}
	}
	return false, nil
}

func desktopApplicationExist(appName string) (bool, error) {
	applicationsPath := "/Applications"
	files, err := os.ReadDir(applicationsPath)
	if err != nil {
		return false, fmt.Errorf("Failed to read directory: %v", err)
	}
	for _, file := range files {
		if strings.Contains(strings.ToLower(file.Name()), appName) {
			return true, nil
		}
	}
	return false, nil
}

func BrewInstall(packageName string) error {
	cmd := CommandInfo{
		PreExecutionMessage:  fmt.Sprintf("Installing %s using Homebrew...", packageName),
		PostExecutionMessage: fmt.Sprintf("%s installed successfully ✔", packageName),
		IsSudo:               false,
		Command:              "brew",
		Args:                 []string{"install", packageName},
	}
	return ExecCommand(cmd)
}

// TODO: Reuse logic from `BrewInstall`
func BrewInstallCask(packageName string) error {
	cmd := CommandInfo{
		PreExecutionMessage:  fmt.Sprintf("Installing %s using Homebrew...", packageName),
		PostExecutionMessage: fmt.Sprintf("%s installed successfully ✔", packageName),
		IsSudo:               false,
		Command:              "brew",
		Args:                 []string{"install", "--cask", packageName},
	}
	return ExecCommand(cmd)
}

func BrewUpgrade(packageName string) error {
	cmd := CommandInfo{
		PreExecutionMessage:  fmt.Sprintf("Upgrading %s using Homebrew...", packageName),
		PostExecutionMessage: fmt.Sprintf("%s upgraded successfully ✔\n", packageName),
		IsSudo:               false,
		Command:              "brew",
		Args:                 []string{"upgrade", packageName},
	}
	return ExecCommand(cmd)
}

func BrewGlobalUpgrade() error {
	cmd := CommandInfo{
		PreExecutionMessage:  "Upgrating Homebrew",
		PostExecutionMessage: "Homebrew upgrated ✔",
		IsSudo:               false,
		Command:              "brew",
		Args:                 []string{"upgrade"},
	}
	return ExecCommand(cmd)
}

func BrewGlobalUpdate() error {
	cmd := CommandInfo{
		PreExecutionMessage:  "Updating Homebrew",
		PostExecutionMessage: "Homebrew updated ✔",
		IsSudo:               false,
		Command:              "brew",
		Args:                 []string{"update"},
	}
	return ExecCommand(cmd)
}
