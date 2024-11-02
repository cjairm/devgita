package macos

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var MACOS_VERSION_NUMBER = 14
var MACOS_VERSION_NAME = "Sonoma"

func CheckVersion() {
	// Check if the OS is macOS
	if runtime.GOOS != "darwin" {
		fmt.Println(
			"\033[31mError: Unable to determine OS. This script is intended for macOS.\033[0m",
		)
		fmt.Println("Installation stopped.")
		os.Exit(1)
	}

	// Get the macOS version
	version, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		fmt.Println("\033[31mError: Unable to determine macOS version.\033[0m")
		fmt.Println("Installation stopped.")
		os.Exit(1)
	}

	// Trim whitespace and split the version string
	versionStr := strings.TrimSpace(string(version))
	versionParts := strings.Split(versionStr, ".")
	if len(versionParts) < 2 {
		fmt.Println("\033[31mError: Unable to parse macOS version.\033[0m")
		fmt.Println("Installation stopped.")
		os.Exit(1)
	}

	// Convert the major and minor version to integers
	major, err := strconv.Atoi(versionParts[0])
	if err != nil {
		fmt.Println("\033[31mError: Unable to parse major version.\033[0m")
		fmt.Println("Installation stopped.")
		os.Exit(1)
	}

	minor, err := strconv.Atoi(versionParts[1])
	if err != nil {
		fmt.Println("\033[31mError: Unable to parse minor version.\033[0m")
		fmt.Println("Installation stopped.")
		os.Exit(1)
	}

	// Check if the version is at least 14.0 (macOS Sonoma)
	if major < MACOS_VERSION_NUMBER || (major == MACOS_VERSION_NUMBER && minor < 0) {
		fmt.Printf("\033[31mError: OS requirement not met\033[0m\n")
		fmt.Printf("You are currently running: macOS %s\n", versionStr)
		fmt.Printf(
			"OS required: macOS %s (%d.0) or higher\n",
			MACOS_VERSION_NAME,
			MACOS_VERSION_NUMBER,
		)
		fmt.Println("Installation stopped.")
		os.Exit(1)
	}

	fmt.Printf("OS requirement met ✔\n\n")
}
