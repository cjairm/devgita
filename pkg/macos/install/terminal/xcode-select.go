package macos

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/cjairm/devgita/pkg/common"
)

func MaybeInstallXcode() error {
	isInstalled, err := isXcodeInstalled()
	if err != nil {
		return err
	}
	if isInstalled {
		fmt.Printf("Xcode is already installed\n\n")
		return nil
	} else {
		return installXcode()
	}

}

func installXcode() error {
	cmd := common.CommandInfo{
		PreExecutionMessage:  "Installing Xcode Command Line Tools",
		PostExecutionMessage: "Xcode Command Line Tools installed ✔",
		IsSudo:               false,
		Command:              "xcode-select",
		Args:                 []string{"--install"},
	}
	return common.ExecCommand(cmd)
}

func isXcodeInstalled() (bool, error) {
	cmd := exec.Command("xcode-select", "-p")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("error running xcode-select: %v", err)
	}

	// Check if the output contains the expected path for Xcode
	xcodePath := strings.TrimSpace(out.String())
	if xcodePath == "/Applications/Xcode.app/Contents/Developer" ||
		strings.HasPrefix(xcodePath, "/Applications/Xcode.app/Contents/Developer") {
		return true, nil
	}

	return false, nil
}
