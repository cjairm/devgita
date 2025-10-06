package commands_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/logger"
	"github.com/cjairm/devgita/pkg/paths"
	"gopkg.in/yaml.v3"
)

type FakePlatform struct {
	Linux bool
	Mac   bool
}

func (f FakePlatform) IsLinux() bool { return f.Linux }
func (f FakePlatform) IsMac() bool   { return f.Mac }

func createFile(t *testing.T, dir, name string) {
	t.Helper()
	err := os.WriteFile(filepath.Join(dir, name), []byte{}, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file %q: %v", name, err)
	}
}

func fakeCmdWithOutput(output string) *exec.Cmd {
	return exec.Command("bash", "-c", "echo -e \""+output+"\"")
}

func TestIsDesktopAppPresent(t *testing.T) {
	t.Run("Linux with matching .desktop file", func(t *testing.T) {
		tmpDir := t.TempDir()
		createFile(t, tmpDir, "myapp.desktop")

		b := commands.NewBaseCommandCustom(FakePlatform{Linux: true})
		found, err := b.IsDesktopAppPresent(tmpDir, "myapp")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !found {
			t.Errorf("Expected to find the desktop app, but did not")
		}
	})

	t.Run("Mac with matching file", func(t *testing.T) {
		tmpDir := t.TempDir()
		createFile(t, tmpDir, "myapp")

		b := commands.NewBaseCommandCustom(FakePlatform{Mac: true})
		found, err := b.IsDesktopAppPresent(tmpDir, "myapp")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !found {
			t.Errorf("Expected to find the app file, but did not")
		}
	})

	t.Run("Linux: no match with wrong extension", func(t *testing.T) {
		tmpDir := t.TempDir()
		createFile(t, tmpDir, "myapp.txt")

		b := commands.NewBaseCommandCustom(FakePlatform{Linux: true})
		found, err := b.IsDesktopAppPresent(tmpDir, "myapp")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if found {
			t.Errorf("Expected not to find a desktop app, but did")
		}
	})

	t.Run("Mac: partial match not found", func(t *testing.T) {
		tmpDir := t.TempDir()
		createFile(t, tmpDir, "unrelated")

		b := commands.NewBaseCommandCustom(FakePlatform{Mac: true})
		found, err := b.IsDesktopAppPresent(tmpDir, "myapp")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if found {
			t.Errorf("Expected not to find the app file, but did")
		}
	})

	t.Run("Linux: case-insensitive match with .desktop", func(t *testing.T) {
		tmpDir := t.TempDir()
		createFile(t, tmpDir, "MyApp.Desktop")

		b := commands.NewBaseCommandCustom(FakePlatform{Linux: true})
		found, err := b.IsDesktopAppPresent(tmpDir, "myapp")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !found {
			t.Errorf("Expected to find desktop app with case-insensitive match")
		}
	})

	t.Run("Linux: directory read error", func(t *testing.T) {
		b := commands.NewBaseCommandCustom(FakePlatform{Linux: true})
		found, err := b.IsDesktopAppPresent("/nonexistent/path", "app")
		if err == nil {
			t.Fatalf("Expected error reading nonexistent directory, got nil")
		}
		if found {
			t.Errorf("Expected not to find the app file, but did")
		}
	})
}

func TestIsPackagePresent_Mac(t *testing.T) {
	b := commands.NewBaseCommandCustom(FakePlatform{Mac: true})

	t.Run("Exact match in brew output", func(t *testing.T) {
		cmd := fakeCmdWithOutput("pkg1\nmytool\npkg2")
		found, err := b.IsPackagePresent(cmd, "mytool")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !found {
			t.Errorf("Expected to find package 'mytool', but did not")
		}
	})

	t.Run("No match in brew output", func(t *testing.T) {
		cmd := fakeCmdWithOutput("pkg1\npkg2")
		found, err := b.IsPackagePresent(cmd, "missing")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if found {
			t.Errorf("Expected not to find package 'missing', but did")
		}
	})
}

func TestIsPackagePresent_Linux(t *testing.T) {
	b := commands.NewBaseCommandCustom(FakePlatform{Linux: true})

	t.Run("Match in dpkg output", func(t *testing.T) {
		// Simulate `dpkg -l` format: status name version
		cmd := fakeCmdWithOutput("ii  mytool  1.0.0\nrc  oldtool  0.9.0")
		found, err := b.IsPackagePresent(cmd, "mytool")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !found {
			t.Errorf("Expected to find package 'mytool', but did not")
		}
	})

	t.Run("No match in dpkg output", func(t *testing.T) {
		cmd := fakeCmdWithOutput("ii  someother  1.0.0")
		found, err := b.IsPackagePresent(cmd, "missing")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if found {
			t.Errorf("Expected not to find package 'missing', but did")
		}
	})
}

func TestIsPackagePresent_CommandError(t *testing.T) {
	b := commands.NewBaseCommandCustom(FakePlatform{Linux: true})

	// Using an invalid command to trigger an error
	cmd := exec.Command("false") // This always exits with non-zero status
	_, err := b.IsPackagePresent(cmd, "anything")
	if err == nil {
		t.Fatalf("Expected error from failed command, got nil")
	}
}

func TestIsFontPresent(t *testing.T) {
	t.Run("Linux: fc-list detects font", func(t *testing.T) {
		b := commands.NewBaseCommandCustom(FakePlatform{Linux: true})
		// Only works if `fc-list` and the font actually exist in system
		found, err := b.IsFontPresent("DejaVu Sans")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		// Can't assert found == true unless we're sure the font is installed,
		// so just log the result
		t.Logf("DejaVu Sans found: %v", found)
	})

	t.Run("Fallback to directory scan - font present", func(t *testing.T) {
		tmpUser := t.TempDir()
		tmpSystem := t.TempDir()

		paths.UserFontsDir = tmpUser
		paths.SystemFontsDir = tmpSystem

		createFile(t, tmpSystem, "myfont.ttf")

		b := commands.NewBaseCommandCustom(FakePlatform{Linux: true})
		commands.LookPathFn = func(string) (string, error) {
			return "", exec.ErrNotFound
		}
		commands.CommandFn = func(name string, args ...string) *exec.Cmd {
			t.Fatalf("fc-list should not be called during fallback")
			return nil
		}

		found, err := b.IsFontPresent("myfont")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !found {
			t.Errorf("Expected to find the font via fallback directory scan, but did not")
		}
	})

	t.Run("Fallback to directory scan - font not found", func(t *testing.T) {
		tmpUserDir := t.TempDir()
		tmpSystemDir := t.TempDir()

		// override paths
		originalUser := paths.UserFontsDir
		originalSystem := paths.SystemFontsDir
		paths.UserFontsDir = tmpUserDir
		paths.SystemFontsDir = tmpSystemDir
		defer func() {
			paths.UserFontsDir = originalUser
			paths.SystemFontsDir = originalSystem
		}()

		b := commands.NewBaseCommandCustom(FakePlatform{Linux: true})
		found, err := b.IsFontPresent("UnknownFont")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if found {
			t.Errorf("Expected not to find nonexistent font")
		}
	})

	t.Run("Fallback to directory scan - font present with valid extension", func(t *testing.T) {
		tmpUser := t.TempDir()
		tmpSystem := t.TempDir()

		paths.UserFontsDir = tmpUser
		paths.SystemFontsDir = tmpSystem

		createFile(t, tmpUser, "FancyFont.OTF")

		b := commands.NewBaseCommandCustom(FakePlatform{Linux: true})
		commands.LookPathFn = func(string) (string, error) {
			return "", exec.ErrNotFound
		}
		commands.CommandFn = func(name string, args ...string) *exec.Cmd {
			t.Fatalf("fc-list should not be called during fallback")
			return nil
		}

		found, err := b.IsFontPresent("fancyfont")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !found {
			t.Errorf("Expected to find font via fallback directory scan, but did not")
		}
	})
}

// Mock install function that tracks calls for MaybeInstall tests
type MockInstaller struct {
	InstallCalled bool
	InstalledItem string
	ShouldFail    bool
}

func (m *MockInstaller) Install(item string) error {
	m.InstallCalled = true
	m.InstalledItem = item
	if m.ShouldFail {
		return fmt.Errorf("mock install error")
	}
	return nil
}

func (m *MockInstaller) Check(item string) (bool, error) {
	// Default: item not installed
	return false, nil
}

type MockInstallerPresent struct {
	MockInstaller
}

func (m *MockInstallerPresent) Check(item string) (bool, error) {
	// Item is already installed
	return true, nil
}

// Helper function to setup test environment for MaybeInstall tests
func setupMaybeInstallTest(t *testing.T, testConfig *config.GlobalConfig) (cleanup func()) {
	// Initialize logger for tests
	logger.Init(false)

	// Set up temporary config - each test gets its own unique directory
	tempDir := t.TempDir() // This creates a unique temp dir per test
	configDir := filepath.Join(tempDir, constants.AppName)
	os.MkdirAll(configDir, 0755)
	configPath := filepath.Join(configDir, constants.GlobalConfigFile)

	// Marshal and write the test config
	data, err := yaml.Marshal(testConfig)
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Override global config path
	originalConfigDir := paths.ConfigDir
	paths.ConfigDir = tempDir

	return func() {
		paths.ConfigDir = originalConfigDir
	}
}

func TestMaybeInstall_ItemNotInstalled_InstallsSuccessfully(t *testing.T) {
	// Create test config
	testConfig := &config.GlobalConfig{
		Installed: config.InstalledConfig{
			Packages: []string{},
		},
		AlreadyInstalled: config.AlreadyInstalledConfig{
			Packages: []string{},
		},
	}

	cleanup := setupMaybeInstallTest(t, testConfig)
	defer cleanup()

	b := commands.NewBaseCommandCustom(FakePlatform{Mac: true})
	mockInstaller := &MockInstaller{}

	err := b.MaybeInstall(
		"test-package",
		[]string{},
		mockInstaller.Check,
		mockInstaller.Install,
		nil,
		"package",
	)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !mockInstaller.InstallCalled {
		t.Errorf("Expected install to be called")
	}

	if mockInstaller.InstalledItem != "test-package" {
		t.Errorf("Expected to install 'test-package', got '%s'", mockInstaller.InstalledItem)
	}

	// Verify item was added to installed config
	var updatedConfig config.GlobalConfig
	updatedConfig.Load()
	if len(updatedConfig.Installed.Packages) != 1 ||
		updatedConfig.Installed.Packages[0] != "test-package" {
		t.Errorf("Expected 'test-package' to be added to installed config")
	}
}

func TestMaybeInstall_ItemAlreadyInstalledByDevgita_SkipsInstall(t *testing.T) {
	// Create test config with item already installed by devgita
	testConfig := &config.GlobalConfig{
		Installed: config.InstalledConfig{
			Packages: []string{"test-package"},
		},
		AlreadyInstalled: config.AlreadyInstalledConfig{
			Packages: []string{},
		},
	}

	cleanup := setupMaybeInstallTest(t, testConfig)
	defer cleanup()

	b := commands.NewBaseCommandCustom(FakePlatform{Mac: true})
	mockInstaller := &MockInstaller{}

	err := b.MaybeInstall(
		"test-package",
		[]string{},
		mockInstaller.Check,
		mockInstaller.Install,
		nil,
		"package",
	)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if mockInstaller.InstallCalled {
		t.Errorf("Expected install NOT to be called for already installed item")
	}
}

func TestMaybeInstall_ItemPreExisting_TracksAsAlreadyInstalled(t *testing.T) {
	// Create test config with empty lists
	testConfig := &config.GlobalConfig{
		Installed: config.InstalledConfig{
			Packages: []string{},
		},
		AlreadyInstalled: config.AlreadyInstalledConfig{
			Packages: []string{},
		},
	}

	// Set up temporary config
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, constants.AppName)
	os.MkdirAll(configDir, 0755)
	configPath := filepath.Join(configDir, constants.GlobalConfigFile)

	data, _ := yaml.Marshal(testConfig)
	os.WriteFile(configPath, data, 0644)

	// Override global config path
	originalConfigDir := paths.ConfigDir
	paths.ConfigDir = tempDir
	defer func() {
		paths.ConfigDir = originalConfigDir
	}()

	b := commands.NewBaseCommandCustom(FakePlatform{Mac: true})
	mockInstaller := &MockInstallerPresent{} // This will return true for Check()

	err := b.MaybeInstall(
		"pre-existing-package",
		[]string{},
		mockInstaller.Check,
		mockInstaller.Install,
		nil,
		"package",
	)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if mockInstaller.InstallCalled {
		t.Errorf("Expected install NOT to be called for pre-existing item")
	}

	// Verify item was added to already installed config
	var updatedConfig config.GlobalConfig
	updatedConfig.Load()
	if len(updatedConfig.AlreadyInstalled.Packages) != 1 ||
		updatedConfig.AlreadyInstalled.Packages[0] != "pre-existing-package" {
		t.Errorf("Expected 'pre-existing-package' to be added to already installed config")
	}
}

func TestMaybeInstall_DifferentItemTypes(t *testing.T) {
	testCases := []struct {
		name     string
		itemType string
		itemName string
	}{
		{"Font", "font", "test-font"},
		{"DesktopApp", "desktop_app", "test-app"},
		{"TerminalTool", "terminal_tool", "test-tool"},
		{"Theme", "theme", "test-theme"},
		{"DevLanguage", "dev_language", "test-language"},
		{"Database", "database", "test-db"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test config
			testConfig := &config.GlobalConfig{
				Installed:        config.InstalledConfig{},
				AlreadyInstalled: config.AlreadyInstalledConfig{},
			}

			// Set up temporary config
			tempDir := t.TempDir()
			configDir := filepath.Join(tempDir, constants.AppName)
			os.MkdirAll(configDir, 0755)
			configPath := filepath.Join(configDir, constants.GlobalConfigFile)

			data, _ := yaml.Marshal(testConfig)
			os.WriteFile(configPath, data, 0644)

			// Override global config path
			originalConfigDir := paths.ConfigDir
			paths.ConfigDir = tempDir
			defer func() {
				paths.ConfigDir = originalConfigDir
			}()

			b := commands.NewBaseCommandCustom(FakePlatform{Mac: true})
			mockInstaller := &MockInstaller{}

			err := b.MaybeInstall(
				tc.itemName,
				[]string{},
				mockInstaller.Check,
				mockInstaller.Install,
				nil,
				tc.itemType,
			)

			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			if !mockInstaller.InstallCalled {
				t.Errorf("Expected install to be called for %s", tc.itemType)
			}

			// Verify item was added to correct section of installed config
			var updatedConfig config.GlobalConfig
			updatedConfig.Load()
			found := false
			switch tc.itemType {
			case "font":
				found = len(updatedConfig.Installed.Fonts) == 1 &&
					updatedConfig.Installed.Fonts[0] == tc.itemName
			case "desktop_app":
				found = len(updatedConfig.Installed.DesktopApps) == 1 &&
					updatedConfig.Installed.DesktopApps[0] == tc.itemName
			case "terminal_tool":
				found = len(updatedConfig.Installed.TerminalTools) == 1 &&
					updatedConfig.Installed.TerminalTools[0] == tc.itemName
			case "theme":
				found = len(updatedConfig.Installed.Themes) == 1 &&
					updatedConfig.Installed.Themes[0] == tc.itemName
			case "dev_language":
				found = len(updatedConfig.Installed.DevLanguages) == 1 &&
					updatedConfig.Installed.DevLanguages[0] == tc.itemName
			case "database":
				found = len(updatedConfig.Installed.Databases) == 1 &&
					updatedConfig.Installed.Databases[0] == tc.itemName
			}

			if !found {
				t.Errorf(
					"Expected '%s' to be added to installed %s config",
					tc.itemName,
					tc.itemType,
				)
			}
		})
	}
}
