package commands

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
