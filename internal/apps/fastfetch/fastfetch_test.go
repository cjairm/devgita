package fastfetch

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/cjairm/devgita/internal/apps"
	"github.com/cjairm/devgita/internal/testutil"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/paths"
)

func init() {
	// Initialize logger for tests
	testutil.InitLogger()
}

func TestNew(t *testing.T) {
	app := New()

	if app == nil {
		t.Fatal("New() returned nil")
	}
}

func TestNameAndKind(t *testing.T) {
	a := &Fastfetch{}
	if a.Name() != constants.Fastfetch {
		t.Errorf("expected Name() %q, got %q", constants.Fastfetch, a.Name())
	}
	if a.Kind() != apps.KindTerminal {
		t.Errorf("expected Kind() KindTerminal, got %v", a.Kind())
	}
}

func TestInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Fastfetch{Cmd: mockApp.Cmd}

	if err := app.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}
	if mockApp.Cmd.InstalledPkg != "fastfetch" {
		t.Fatalf(
			"expected InstallPackage(%s), got %q",
			"fastfetch",
			mockApp.Cmd.InstalledPkg,
		)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestForceInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Fastfetch{Cmd: mockApp.Cmd, Base: mockApp.Base}

	if err := app.ForceInstall(); err != nil {
		t.Fatalf("ForceInstall() should succeed even when uninstall is not supported: %v", err)
	}
	if mockApp.Cmd.InstalledPkg != constants.Fastfetch {
		t.Errorf("expected Install to be called, got %q", mockApp.Cmd.InstalledPkg)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestSoftInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Fastfetch{Cmd: mockApp.Cmd}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}
	if mockApp.Cmd.MaybeInstalled != "fastfetch" {
		t.Fatalf(
			"expected MaybeInstallPackage(%s), got %q",
			"fastfetch",
			mockApp.Cmd.MaybeInstalled,
		)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestUninstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Fastfetch{Cmd: mockApp.Cmd, Base: mockApp.Base}

	err := app.Uninstall()
	if err == nil {
		t.Fatal("expected Uninstall to return error for unsupported operation")
	}
	if !errors.Is(err, apps.ErrUninstallNotSupported) {
		t.Errorf("expected ErrUninstallNotSupported, got: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestUpdate(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Fastfetch{Cmd: mockApp.Cmd, Base: mockApp.Base}

	err := app.Update()
	if err == nil {
		t.Fatal("expected Update to return error")
	}
	if !errors.Is(err, apps.ErrUpdateNotSupported) {
		t.Errorf("expected ErrUpdateNotSupported, got: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestForceConfigure(t *testing.T) {
	// Create temp directories for testing
	tempSourceDir := t.TempDir()
	tempTargetDir := t.TempDir()

	// Override global paths for the duration of the test
	oldFastFetchConfigAppDir := paths.Paths.App.Configs.Fastfetch
	oldFastFetchConfigLocalDir := paths.Paths.Config.Fastfetch
	defer func() {
		paths.Paths.App.Configs.Fastfetch = oldFastFetchConfigAppDir
		paths.Paths.Config.Fastfetch = oldFastFetchConfigLocalDir
	}()
	paths.Paths.App.Configs.Fastfetch = tempSourceDir
	paths.Paths.Config.Fastfetch = tempTargetDir

	// Create a test config file in the source directory
	testConfigContent := `{
		"$schema": "https://github.com/fastfetch-cli/fastfetch/raw/dev/doc/json_schema.json",
		"logo": {
			"source": "auto",
			"padding": {
				"top": 1,
				"left": 4
			}
		},
		"display": {
		"separator": " ~ "
	    }
	}`
	configFile := filepath.Join(tempSourceDir, "config.jsonc")
	if err := os.WriteFile(configFile, []byte(testConfigContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	mockApp := testutil.NewMockApp()
	app := &Fastfetch{Cmd: mockApp.Cmd}

	if err := app.ForceConfigure(); err != nil {
		t.Fatalf("ForceConfigure error: %v", err)
	}

	// Verify that the config file was copied to the target directory
	targetConfigFile := filepath.Join(tempTargetDir, "config.jsonc")
	if _, err := os.Stat(targetConfigFile); os.IsNotExist(err) {
		t.Fatalf("Expected config file to be copied to target directory, but it doesn't exist")
	}

	// Verify the content was copied correctly
	copiedContent, err := os.ReadFile(targetConfigFile)
	if err != nil {
		t.Fatalf("Failed to read copied config file: %v", err)
	}

	if string(copiedContent) != testConfigContent {
		t.Fatalf(
			"Expected copied content to match original, but it doesn't. Original: %s, Copied: %s",
			testConfigContent,
			string(copiedContent),
		)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestSoftConfigure(t *testing.T) {
	// Test case 1: Configuration doesn't exist - should configure
	t.Run("ConfigureWhenNotExists", func(t *testing.T) {
		tempSourceDir := t.TempDir()
		tempTargetDir := t.TempDir()

		// Override global paths for the duration of the test
		oldFastFetchConfigAppDir := paths.Paths.App.Configs.Fastfetch
		oldFastFetchConfigLocalDir := paths.Paths.Config.Fastfetch
		defer func() {
			paths.Paths.App.Configs.Fastfetch = oldFastFetchConfigAppDir
			paths.Paths.Config.Fastfetch = oldFastFetchConfigLocalDir
		}()
		paths.Paths.App.Configs.Fastfetch = tempSourceDir
		paths.Paths.Config.Fastfetch = tempTargetDir

		// Create a test config file in the source directory
		testConfigContent := `{"display": {"separator": " ~ "}}`
		configFile := filepath.Join(tempSourceDir, "config.jsonc")
		if err := os.WriteFile(configFile, []byte(testConfigContent), 0644); err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}

		mockApp := testutil.NewMockApp()
		app := &Fastfetch{Cmd: mockApp.Cmd}

		if err := app.SoftConfigure(); err != nil {
			t.Fatalf("SoftConfigure error: %v", err)
		}

		// Verify that the config file was copied to the target directory
		targetConfigFile := filepath.Join(tempTargetDir, "config.jsonc")
		if _, err := os.Stat(targetConfigFile); os.IsNotExist(err) {
			t.Fatalf("Expected config file to be copied to target directory, but it doesn't exist")
		}

		testutil.VerifyNoRealCommands(t, mockApp.Base)
	})

	// Test case 2: Configuration already exists - should skip
	t.Run("SkipWhenExists", func(t *testing.T) {
		tempSourceDir := t.TempDir()
		tempTargetDir := t.TempDir()

		// Override global paths for the duration of the test
		oldFastFetchConfigAppDir := paths.Paths.App.Configs.Fastfetch
		oldFastFetchConfigLocalDir := paths.Paths.Config.Fastfetch
		defer func() {
			paths.Paths.App.Configs.Fastfetch = oldFastFetchConfigAppDir
			paths.Paths.Config.Fastfetch = oldFastFetchConfigLocalDir
		}()
		paths.Paths.App.Configs.Fastfetch = tempSourceDir
		paths.Paths.Config.Fastfetch = tempTargetDir

		// Create existing config file in target directory
		existingContent := `{"existing": "config"}`
		targetConfigFile := filepath.Join(tempTargetDir, "config.jsonc")
		if err := os.WriteFile(targetConfigFile, []byte(existingContent), 0644); err != nil {
			t.Fatalf("Failed to create existing config file: %v", err)
		}

		// Create different config file in source directory
		sourceContent := `{"new": "config"}`
		sourceConfigFile := filepath.Join(tempSourceDir, "config.jsonc")
		if err := os.WriteFile(sourceConfigFile, []byte(sourceContent), 0644); err != nil {
			t.Fatalf("Failed to create source config file: %v", err)
		}

		mockApp := testutil.NewMockApp()
		app := &Fastfetch{Cmd: mockApp.Cmd}

		if err := app.SoftConfigure(); err != nil {
			t.Fatalf("SoftConfigure error: %v", err)
		}

		// Verify content wasn't overwritten
		content, err := os.ReadFile(targetConfigFile)
		if err != nil {
			t.Fatalf("Failed to read target config file: %v", err)
		}

		if string(content) != existingContent {
			t.Fatalf(
				"Expected file content to remain unchanged, but it was modified. Original: %s, New: %s",
				existingContent,
				string(content),
			)
		}

		testutil.VerifyNoRealCommands(t, mockApp.Base)
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
