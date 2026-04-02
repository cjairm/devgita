package tmux_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/apps/tmux"
	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/testutil"
	"github.com/cjairm/devgita/pkg/paths"
)

func init() {
	// Initialize logger for tests
	testutil.InitLogger()
}

func TestNew(t *testing.T) {
	t.Helper()

	app := tmux.New()
	if app == nil {
		t.Error("Expected New() to return a non-nil Tmux instance")
	}
}

func TestInstall(t *testing.T) {
	t.Helper()

	tests := []struct {
		name        string
		shouldError bool
		installErr  error
	}{
		{
			name:        "successful installation",
			shouldError: false,
			installErr:  nil,
		},
		{
			name:        "installation failure",
			shouldError: true,
			installErr:  errors.New("installation failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			mockCmd := commands.NewMockCommand()
			mockCmd.InstallError = tt.installErr

			app := &tmux.Tmux{
				Cmd: mockCmd,
			}

			err := app.Install()

			if tt.shouldError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Verify the correct package was passed
			if mockCmd.InstalledPkg != "tmux" {
				t.Errorf(
					"Expected package 'tmux', got '%s'",
					mockCmd.InstalledPkg,
				)
			}
		})
	}
}

func TestSoftInstall(t *testing.T) {
	t.Helper()

	tests := []struct {
		name        string
		shouldError bool
		installErr  error
	}{
		{
			name:        "successful soft installation",
			shouldError: false,
			installErr:  nil,
		},
		{
			name:        "soft installation failure",
			shouldError: true,
			installErr:  errors.New("soft installation failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			mockCmd := commands.NewMockCommand()
			mockCmd.MaybeInstallError = tt.installErr

			app := &tmux.Tmux{
				Cmd: mockCmd,
			}

			err := app.SoftInstall()

			if tt.shouldError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Verify the correct package was passed
			if mockCmd.MaybeInstalled != "tmux" {
				t.Errorf(
					"Expected package 'tmux', got '%s'",
					mockCmd.MaybeInstalled,
				)
			}
		})
	}
}

func TestForceConfigure(t *testing.T) {
	t.Helper()

	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	// Create source directory with tmux config
	sourceDir := filepath.Join(tc.AppDir, "tmux")
	err := os.MkdirAll(sourceDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create destination directory
	destDir := tc.ConfigDir
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create destination directory: %v", err)
	}

	paths.Paths.App.Configs.Tmux = sourceDir
	paths.Paths.Home.Root = destDir

	// Create source tmux.conf file (without leading dot in source)
	sourceConfig := filepath.Join(sourceDir, "tmux.conf")
	configContent := "# Test tmux configuration\nset -g default-terminal \"screen-256color\""
	err = os.WriteFile(sourceConfig, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create source config: %v", err)
	}

	app := tmux.New()

	err = app.ForceConfigure()
	if err != nil {
		t.Errorf("ForceConfigure returned error: %v", err)
	}

	// Verify config was copied to destination
	destConfig := filepath.Join(destDir, ".tmux.conf")
	content, err := os.ReadFile(destConfig)
	if err != nil {
		t.Fatalf("Failed to read destination config: %v", err)
	}

	if string(content) != configContent {
		t.Errorf(
			"Expected config content %q, got %q",
			configContent,
			string(content),
		)
	}

	// Verify shell config was generated
	shellContent, err := os.ReadFile(tc.ZshConfigPath)
	if err != nil {
		t.Fatalf("Failed to read shell config: %v", err)
	}

	if !strings.Contains(string(shellContent), "# Tmux enabled") {
		t.Error("Expected shell config to contain Tmux feature")
	}

	testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
}

func TestSoftConfigure(t *testing.T) {
	t.Helper()

	// Test case 1: Configuration doesn't exist - should configure
	t.Run("ConfigureWhenNotExists", func(t *testing.T) {
		tc := testutil.SetupCompleteTest(t)
		defer tc.Cleanup()

		// Create source directory with tmux config
		sourceDir := filepath.Join(tc.AppDir, "tmux")
		err := os.MkdirAll(sourceDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create source directory: %v", err)
		}

		destDir := tc.ConfigDir

		paths.Paths.App.Configs.Tmux = sourceDir
		paths.Paths.Home.Root = destDir

		// Create source tmux.conf file (without leading dot in source)
		sourceConfig := filepath.Join(sourceDir, "tmux.conf")
		configContent := "# Test tmux configuration\nset -g default-terminal \"screen-256color\""
		err = os.WriteFile(sourceConfig, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create source config: %v", err)
		}

		// Mock the HOME environment variable since SoftConfigure uses os.UserHomeDir()
		oldHome := os.Getenv("HOME")
		defer func() {
			if oldHome != "" {
				os.Setenv("HOME", oldHome)
			} else {
				os.Unsetenv("HOME")
			}
		}()
		os.Setenv("HOME", destDir)

		app := tmux.New()

		err = app.SoftConfigure()
		if err != nil {
			t.Errorf("SoftConfigure returned error: %v", err)
		}

		// Verify config was copied to destination
		destConfig := filepath.Join(destDir, ".tmux.conf")
		content, err := os.ReadFile(destConfig)
		if err != nil {
			t.Fatalf("Failed to read destination config: %v", err)
		}

		if string(content) != configContent {
			t.Errorf(
				"Expected config content %q, got %q",
				configContent,
				string(content),
			)
		}

		// Verify shell config was generated
		shellContent, err := os.ReadFile(tc.ZshConfigPath)
		if err != nil {
			t.Fatalf("Failed to read shell config: %v", err)
		}

		if !strings.Contains(string(shellContent), "# Tmux enabled") {
			t.Error("Expected shell config to contain Tmux feature on first call")
		}

		testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
	})

	// Test case 2: Configuration already exists - should skip file copy but enable shell feature
	t.Run("SkipWhenExists", func(t *testing.T) {
		tc := testutil.SetupCompleteTest(t)
		defer tc.Cleanup()

		// Create home directory with existing .tmux.conf
		homeDir := tc.ConfigDir

		// Set Home path before creating the config file
		paths.Paths.Home.Root = homeDir

		existingConfig := filepath.Join(homeDir, ".tmux.conf")
		existingContent := "# Existing tmux configuration\nset -g mouse on"
		err := os.WriteFile(existingConfig, []byte(existingContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create existing config: %v", err)
		}

		// Create source directory (though it shouldn't be used for file copy)
		sourceDir := filepath.Join(tc.AppDir, "tmux")
		err = os.MkdirAll(sourceDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create source directory: %v", err)
		}
		paths.Paths.App.Configs.Tmux = sourceDir

		// Create a different source config to prove it's not copied (without leading dot in source)
		sourceConfig := filepath.Join(sourceDir, "tmux.conf")
		sourceContent := "# Different config that should not be copied"
		err = os.WriteFile(sourceConfig, []byte(sourceContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create source config: %v", err)
		}

		// Mock the UserHomeDir for this test by temporarily setting HOME env var
		oldHome := os.Getenv("HOME")
		defer func() {
			if oldHome != "" {
				os.Setenv("HOME", oldHome)
			} else {
				os.Unsetenv("HOME")
			}
		}()
		os.Setenv("HOME", homeDir)

		app := tmux.New()

		err = app.SoftConfigure()
		if err != nil {
			t.Errorf("SoftConfigure returned error: %v", err)
		}

		// Verify content wasn't changed
		contentAfter, err := os.ReadFile(existingConfig)
		if err != nil {
			t.Fatalf("Failed to read config after test: %v", err)
		}

		if string(contentAfter) != existingContent {
			t.Errorf("Expected config to remain unchanged, but it was modified")
		}

		if string(contentAfter) == sourceContent {
			t.Errorf(
				"Config was overwritten with source content when it should have been preserved",
			)
		}

		// Verify shell config was still generated (feature should be enabled even when file exists)
		shellContent, err := os.ReadFile(tc.ZshConfigPath)
		if err != nil {
			t.Fatalf("Failed to read shell config: %v", err)
		}

		if !strings.Contains(string(shellContent), "# Tmux enabled") {
			t.Error("Expected shell config to contain Tmux feature even when config file exists")
		}

		testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
	})
}
