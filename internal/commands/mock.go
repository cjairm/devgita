package commands

import "os/exec"

// MockCommand provides a mock implementation of the Command interface for testing
type MockCommand struct {
	InstalledPkg          string
	UninstalledPkg        string
	MaybeInstalled        string
	InstalledDesktopApp   string
	MaybeInstalledDesktop string
	FontURL               string
	FontName              string

	// Error fields to simulate various failure scenarios
	InstallError        error
	UninstallError      error
	MaybeInstallError   error
	DesktopInstallError error
	FontInstallError    error
	ValidationError     error

	// State tracking
	PackageManagerInstalled bool
	PackageInstalled        bool
	DesktopAppInstalled     bool
}

// NewMockCommand creates a new MockCommand with sensible defaults
func NewMockCommand() *MockCommand {
	return &MockCommand{
		PackageManagerInstalled: true,
		PackageInstalled:        false,
		DesktopAppInstalled:     false,
	}
}

func (m *MockCommand) InstallPackage(pkg string) error {
	m.InstalledPkg = pkg
	return m.InstallError
}

func (m *MockCommand) UninstallPackage(pkg string) error {
	m.UninstalledPkg = pkg
	return m.UninstallError
}

func (m *MockCommand) MaybeInstallPackage(pkg string, alias ...string) error {
	m.MaybeInstalled = pkg
	return m.MaybeInstallError
}

func (m *MockCommand) InstallDesktopApp(packageName string) error {
	m.InstalledDesktopApp = packageName
	return m.DesktopInstallError
}

func (m *MockCommand) MaybeInstallDesktopApp(desktopAppName string, alias ...string) error {
	m.MaybeInstalledDesktop = desktopAppName
	return m.DesktopInstallError
}

func (m *MockCommand) MaybeInstallFont(url, fontName string, runCache bool, alias ...string) error {
	m.FontURL = url
	m.FontName = fontName
	return m.FontInstallError
}

func (m *MockCommand) ValidateOSVersion() error {
	return m.ValidationError
}

func (m *MockCommand) MaybeInstallPackageManager() error {
	return m.ValidationError
}

func (m *MockCommand) InstallPackageManager() error {
	return m.ValidationError
}

func (m *MockCommand) IsPackageManagerInstalled() bool {
	return m.PackageManagerInstalled
}

func (m *MockCommand) IsPackageInstalled(packageName string) (bool, error) {
	return m.PackageInstalled, nil
}

func (m *MockCommand) IsDesktopAppInstalled(desktopAppName string) (bool, error) {
	return m.DesktopAppInstalled, nil
}

// Helper methods for testing

// Reset clears all tracked state for reuse in multiple tests
func (m *MockCommand) Reset() {
	m.InstalledPkg = ""
	m.UninstalledPkg = ""
	m.MaybeInstalled = ""
	m.InstalledDesktopApp = ""
	m.MaybeInstalledDesktop = ""
	m.FontURL = ""
	m.FontName = ""

	m.InstallError = nil
	m.UninstallError = nil
	m.MaybeInstallError = nil
	m.DesktopInstallError = nil
	m.FontInstallError = nil
	m.ValidationError = nil
}

// SetError configures error scenarios for different operations
func (m *MockCommand) SetError(operation string, err error) {
	switch operation {
	case "install":
		m.InstallError = err
	case "uninstall":
		m.UninstallError = err
	case "maybe-install":
		m.MaybeInstallError = err
	case "desktop":
		m.DesktopInstallError = err
	case "font":
		m.FontInstallError = err
	case "validation":
		m.ValidationError = err
	}
}

// MockBaseCommand provides a mock implementation for BaseCommand methods
// This allows tests to avoid running actual system commands
type MockBaseCommand struct {
	// Tracks all ExecCommand calls for verification
	ExecCommandCalls []CommandParams

	// Return values for ExecCommand
	ExecCommandStdout string
	ExecCommandStderr string
	ExecCommandError  error

	// Return values for presence checks
	IsDesktopAppPresentResult bool
	IsPackagePresentResult    bool
	IsFontPresentResult       bool

	// Return values for other methods
	SetupError          error
	MaybeSetupError     error
	MaybeInstallError   error
	InstallFontURLError error
}

// NewMockBaseCommand creates a new MockBaseCommand with sensible defaults
func NewMockBaseCommand() *MockBaseCommand {
	return &MockBaseCommand{
		ExecCommandCalls:          []CommandParams{},
		ExecCommandStdout:         "",
		ExecCommandStderr:         "",
		ExecCommandError:          nil,
		IsDesktopAppPresentResult: false,
		IsPackagePresentResult:    false,
		IsFontPresentResult:       false,
		SetupError:                nil,
		MaybeSetupError:           nil,
		MaybeInstallError:         nil,
		InstallFontURLError:       nil,
	}
}

// ExecCommand mocks the BaseCommand.ExecCommand method
// It records the call parameters and returns the configured mock values
func (m *MockBaseCommand) ExecCommand(cmd CommandParams) (string, string, error) {
	m.ExecCommandCalls = append(m.ExecCommandCalls, cmd)
	return m.ExecCommandStdout, m.ExecCommandStderr, m.ExecCommandError
}

// Setup mocks the BaseCommand.Setup method
func (m *MockBaseCommand) Setup(line string) error {
	return m.SetupError
}

// MaybeSetup mocks the BaseCommand.MaybeSetup method
func (m *MockBaseCommand) MaybeSetup(line, toSearch string) error {
	return m.MaybeSetupError
}

// IsDesktopAppPresent mocks the BaseCommand.IsDesktopAppPresent method
func (m *MockBaseCommand) IsDesktopAppPresent(dirPath, appName string) (bool, error) {
	return m.IsDesktopAppPresentResult, nil
}

// IsPackagePresent mocks the BaseCommand.IsPackagePresent method
func (m *MockBaseCommand) IsPackagePresent(cmd *exec.Cmd, packageName string) (bool, error) {
	return m.IsPackagePresentResult, nil
}

// IsFontPresent mocks the BaseCommand.IsFontPresent method
func (m *MockBaseCommand) IsFontPresent(fontName string) (bool, error) {
	return m.IsFontPresentResult, nil
}

// MaybeInstall mocks the BaseCommand.MaybeInstall method
func (m *MockBaseCommand) MaybeInstall(itemName string, alias []string, checkInstalled func(string) (bool, error), installFunc func(string) error, installURLFunc func(string) error, itemType string) error {
	return m.MaybeInstallError
}

// InstallFontFromURL mocks the BaseCommand.InstallFontFromURL method
func (m *MockBaseCommand) InstallFontFromURL(url, fontFileName string, runCache bool) error {
	return m.InstallFontURLError
}

// Reset clears all tracked state for reuse in multiple tests
func (m *MockBaseCommand) ResetExecCommand() {
	m.ExecCommandCalls = []CommandParams{}
	m.ExecCommandStdout = ""
	m.ExecCommandStderr = ""
	m.ExecCommandError = nil
}

// SetExecCommandResult configures the return values for ExecCommand
func (m *MockBaseCommand) SetExecCommandResult(stdout, stderr string, err error) {
	m.ExecCommandStdout = stdout
	m.ExecCommandStderr = stderr
	m.ExecCommandError = err
}

// GetLastExecCommandCall returns the most recent ExecCommand call parameters
// Returns nil if no calls have been made
func (m *MockBaseCommand) GetLastExecCommandCall() *CommandParams {
	if len(m.ExecCommandCalls) == 0 {
		return nil
	}
	return &m.ExecCommandCalls[len(m.ExecCommandCalls)-1]
}

// GetExecCommandCallCount returns the number of times ExecCommand was called
func (m *MockBaseCommand) GetExecCommandCallCount() int {
	return len(m.ExecCommandCalls)
}
