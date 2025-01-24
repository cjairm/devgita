package common

import (
	"bytes"
	"fmt"
	"os/exec"
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
