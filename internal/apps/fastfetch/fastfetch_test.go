package fastfetch

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/logger"
	"github.com/cjairm/devgita/pkg/paths"
)

func init() {
	// Initialize logger for tests
	logger.Init(false)
}

func TestNew(t *testing.T) {
	app := New()

	if app == nil {
		t.Fatal("New() returned nil")
	}
}

func TestInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Fastfetch{Cmd: mc}

	if err := app.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}
	if mc.InstalledPkg != "fastfetch" {
		t.Fatalf(
			"expected InstallPackage(%s), got %q",
			"fastfetch",
			mc.InstalledPkg,
		)
	}
}

func TestSoftInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Fastfetch{Cmd: mc}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}
	if mc.MaybeInstalled != "fastfetch" {
		t.Fatalf(
			"expected MaybeInstallPackage(%s), got %q",
			"fastfetch",
			mc.MaybeInstalled,
		)
	}
}

func TestForceConfigure(t *testing.T) {
	// Create temp directories for testing
	tempSourceDir := t.TempDir()
	tempTargetDir := t.TempDir()

	// Override global paths for the duration of the test
	oldFastFetchConfigAppDir := paths.FastFetchConfigAppDir
	oldFastFetchConfigLocalDir := paths.FastFetchConfigLocalDir
	defer func() {
		paths.FastFetchConfigAppDir = oldFastFetchConfigAppDir
		paths.FastFetchConfigLocalDir = oldFastFetchConfigLocalDir
	}()
	paths.FastFetchConfigAppDir = tempSourceDir
	paths.FastFetchConfigLocalDir = tempTargetDir

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

	mc := commands.NewMockCommand()
	app := &Fastfetch{Cmd: mc}

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
}

func TestSoftConfigure(t *testing.T) {
	// Test case 1: Configuration doesn't exist - should configure
	t.Run("ConfigureWhenNotExists", func(t *testing.T) {
		tempSourceDir := t.TempDir()
		tempTargetDir := t.TempDir()

		// Override global paths for the duration of the test
		oldFastFetchConfigAppDir := paths.FastFetchConfigAppDir
		oldFastFetchConfigLocalDir := paths.FastFetchConfigLocalDir
		defer func() {
			paths.FastFetchConfigAppDir = oldFastFetchConfigAppDir
			paths.FastFetchConfigLocalDir = oldFastFetchConfigLocalDir
		}()
		paths.FastFetchConfigAppDir = tempSourceDir
		paths.FastFetchConfigLocalDir = tempTargetDir

		// Create a test config file in the source directory
		testConfigContent := `{"display": {"separator": " ~ "}}`
		configFile := filepath.Join(tempSourceDir, "config.jsonc")
		if err := os.WriteFile(configFile, []byte(testConfigContent), 0644); err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}

		mc := commands.NewMockCommand()
		app := &Fastfetch{Cmd: mc}

		if err := app.SoftConfigure(); err != nil {
			t.Fatalf("SoftConfigure error: %v", err)
		}

		// Verify that the config file was copied to the target directory
		targetConfigFile := filepath.Join(tempTargetDir, "config.jsonc")
		if _, err := os.Stat(targetConfigFile); os.IsNotExist(err) {
			t.Fatalf("Expected config file to be copied to target directory, but it doesn't exist")
		}
	})

	// Test case 2: Configuration already exists - should skip
	t.Run("SkipWhenExists", func(t *testing.T) {
		tempSourceDir := t.TempDir()
		tempTargetDir := t.TempDir()

		// Override global paths for the duration of the test
		oldFastFetchConfigAppDir := paths.FastFetchConfigAppDir
		oldFastFetchConfigLocalDir := paths.FastFetchConfigLocalDir
		defer func() {
			paths.FastFetchConfigAppDir = oldFastFetchConfigAppDir
			paths.FastFetchConfigLocalDir = oldFastFetchConfigLocalDir
		}()
		paths.FastFetchConfigAppDir = tempSourceDir
		paths.FastFetchConfigLocalDir = tempTargetDir

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

		mc := commands.NewMockCommand()
		app := &Fastfetch{Cmd: mc}

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
