// Package xcode provides installation and command execution management for
// Xcode Command Line Tools with devgita integration.
//
// Xcode Command Line Tools is a macOS-specific package that provides essential
// development tools including compilers (clang, gcc), build tools (make), SDKs,
// and headers needed for software development on macOS. This module ensures
// Xcode Command Line Tools are properly installed on macOS systems.
//
// References:
//   - Project overview: docs/project-overview.md
//   - Testing patterns: docs/guides/testing-patterns.md
//   - Error handling: docs/guides/error-handling.md
//   - Module documentation: docs/tooling/terminal/core/xcode.md
//
// Common xcode-select commands available through ExecuteCommand():
//   - xcode-select --version         - Show xcode-select version
//   - xcode-select --print-path      - Print active developer directory path
//   - xcode-select --install         - Install Xcode Command Line Tools
//   - xcode-select --switch <path>   - Switch active developer directory
//   - xcode-select --reset           - Reset to default developer directory

package xcode

import (
	"fmt"
	"strings"

	cmd "github.com/cjairm/devgita/internal/commands"
)

type XcodeCommandLineTools struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func New() *XcodeCommandLineTools {
	return &XcodeCommandLineTools{
		Cmd:  cmd.NewCommand(),
		Base: cmd.NewBaseCommand(),
	}
}

// Install installs Xcode Command Line Tools on macOS.
// Returns error on non-macOS platforms or if installation fails.
func (x *XcodeCommandLineTools) Install() error {
	installed, err := x.isInstalled()
	if err != nil {
		return fmt.Errorf("failed to check xcode command line tools installation: %w", err)
	}
	if installed {
		return nil
	}
	// Install using xcode-select --install
	_, _, err = x.Base.ExecCommand(cmd.CommandParams{
		PreExecMsg: "Installing Xcode Command Line Tools",
		Command:    "xcode-select",
		Args:       []string{"--install"},
	})
	if err != nil {
		return fmt.Errorf("failed to install xcode command line tools: %w", err)
	}
	return nil
}

// ForceInstall forces installation by uninstalling first, then installing.
// Note: Uninstall is not supported, so this will return an error.
func (x *XcodeCommandLineTools) ForceInstall() error {
	if err := x.Uninstall(); err != nil {
		return fmt.Errorf("failed to uninstall before force install: %w", err)
	}
	return x.Install()
}

// SoftInstall installs Xcode Command Line Tools only if not already present.
func (x *XcodeCommandLineTools) SoftInstall() error {
	installed, err := x.isInstalled()
	if err != nil {
		return fmt.Errorf("failed to check xcode command line tools installation: %w", err)
	}
	if installed {
		return nil
	}
	return x.Install()
}

// ForceConfigure applies configuration.
// Xcode Command Line Tools don't require separate configuration files.
// Returns nil (no-op) for interface compliance.
func (x *XcodeCommandLineTools) ForceConfigure() error {
	// Xcode Command Line Tools don't require configuration files
	// No-op for interface compliance
	return nil
}

// SoftConfigure applies configuration only if not already configured.
// Xcode Command Line Tools don't require separate configuration files.
// Returns nil (no-op) for interface compliance.
func (x *XcodeCommandLineTools) SoftConfigure() error {
	// Xcode Command Line Tools don't require configuration files
	// No-op for interface compliance
	return nil
}

// Uninstall removes Xcode Command Line Tools installation.
// Returns error as uninstall is not supported through devgita.
func (x *XcodeCommandLineTools) Uninstall() error {
	return fmt.Errorf("xcode command line tools uninstall not supported through devgita")
}

// ExecuteCommand executes xcode-select commands with provided arguments.
// Common commands: --version, --print-path, --switch, --reset
func (x *XcodeCommandLineTools) ExecuteCommand(args ...string) error {
	_, _, err := x.Base.ExecCommand(cmd.CommandParams{
		Command: "xcode-select",
		Args:    args,
	})
	if err != nil {
		return fmt.Errorf("failed to run xcode-select command: %w", err)
	}
	return nil
}

// Update updates Xcode Command Line Tools installation.
// Returns error as update is not implemented - use system update mechanisms.
func (x *XcodeCommandLineTools) Update() error {
	return fmt.Errorf("xcode command line tools update not implemented - use macOS software update")
}

// isInstalled checks if Xcode Command Line Tools are installed.
// Returns true if installed, false otherwise.
func (x *XcodeCommandLineTools) isInstalled() (bool, error) {
	stdout, _, err := x.Base.ExecCommand(cmd.CommandParams{
		Command: "xcode-select",
		Args:    []string{"-p"},
		IsSudo:  false,
	})
	if err != nil {
		return false, fmt.Errorf("error running xcode-select: %v", err)
	}
	xcodePath := strings.ToLower(strings.TrimSpace(stdout))
	return strings.Contains(xcodePath, "xcode.app") ||
		strings.Contains(xcodePath, "commandlinetools"), nil
}
