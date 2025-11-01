package syntaxhighlighting_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/cjairm/devgita/internal/apps/syntaxhighlighting"
	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/logger"
	"github.com/cjairm/devgita/pkg/paths"
)

func init() {
	// Initialize logger for tests
	logger.Init(false)
}

func TestNew(t *testing.T) {
	t.Helper()

	app := syntaxhighlighting.New()
	if app == nil {
		t.Error("Expected New() to return a non-nil Syntaxhighlighting instance")
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

			app := &syntaxhighlighting.Syntaxhighlighting{
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
			if mockCmd.InstalledPkg != "zsh-syntax-highlighting" {
				t.Errorf(
					"Expected package 'zsh-syntax-highlighting', got '%s'",
					mockCmd.InstalledPkg,
				)
			}
		})
	}
}

func TestForceInstall(t *testing.T) {
	t.Helper()

	app := syntaxhighlighting.New()

	// ForceInstall should always fail because Uninstall returns error
	err := app.ForceInstall()

	if err == nil {
		t.Error("Expected ForceInstall to return an error (because Uninstall is not supported)")
	}

	expectedMsg := "failed to uninstall syntaxhighlighting before force install: zsh-syntax-highlighting uninstall is not supported"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
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

			app := &syntaxhighlighting.Syntaxhighlighting{
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
			if mockCmd.MaybeInstalled != "zsh-syntax-highlighting" {
				t.Errorf(
					"Expected package 'zsh-syntax-highlighting', got '%s'",
					mockCmd.MaybeInstalled,
				)
			}
		})
	}
}

func TestForceConfigure(t *testing.T) {
	t.Helper()

	tempDir := t.TempDir()

	// Override global paths for the duration of the test
	oldAppDir := paths.AppDir
	defer func() { paths.AppDir = oldAppDir }()
	paths.AppDir = tempDir

	// Create the devgita.zsh file in temp directory
	zshFile := filepath.Join(tempDir, "devgita.zsh")
	file, err := os.Create(zshFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	app := syntaxhighlighting.New()

	err = app.ForceConfigure()
	if err != nil {
		t.Errorf("ForceConfigure returned error: %v", err)
	}

	// Verify content was added to the file
	content, err := os.ReadFile(zshFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	expectedContent := "source $(brew --prefix)/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh"
	if !contains(string(content), expectedContent) {
		t.Errorf(
			"Expected file to contain %q, but it didn't. Content: %s",
			expectedContent,
			string(content),
		)
	}
}

func TestSoftConfigure(t *testing.T) {
	t.Helper()

	// Test case 1: Configuration doesn't exist - should configure
	t.Run("ConfigureWhenNotExists", func(t *testing.T) {
		tempDir := t.TempDir()

		// Override global paths for the duration of the test
		oldAppDir := paths.AppDir
		defer func() { paths.AppDir = oldAppDir }()
		paths.AppDir = tempDir

		// Create the devgita.zsh file in temp directory
		zshFile := filepath.Join(tempDir, "devgita.zsh")
		file, err := os.Create(zshFile)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		file.Close()

		app := syntaxhighlighting.New()

		err = app.SoftConfigure()
		if err != nil {
			t.Errorf("SoftConfigure returned error: %v", err)
		}

		// Verify content was added to the file
		content, err := os.ReadFile(zshFile)
		if err != nil {
			t.Fatalf("Failed to read test file: %v", err)
		}

		expectedContent := "source $(brew --prefix)/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh"
		if !contains(string(content), expectedContent) {
			t.Errorf(
				"Expected file to contain %q, but it didn't. Content: %s",
				expectedContent,
				string(content),
			)
		}
	})

	// Test case 2: Configuration already exists - should skip
	t.Run("SkipWhenExists", func(t *testing.T) {
		tempDir := t.TempDir()

		// Override global paths for the duration of the test
		oldAppDir := paths.AppDir
		defer func() { paths.AppDir = oldAppDir }()
		paths.AppDir = tempDir

		// Create the devgita.zsh file with existing content
		zshFile := filepath.Join(tempDir, "devgita.zsh")
		existingContent := "# Existing content\nsource $(brew --prefix)/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh\n# More content"
		err := os.WriteFile(zshFile, []byte(existingContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		app := syntaxhighlighting.New()

		err = app.SoftConfigure()
		if err != nil {
			t.Errorf("SoftConfigure returned error: %v", err)
		}

		// Verify content wasn't duplicated
		content, err := os.ReadFile(zshFile)
		if err != nil {
			t.Fatalf("Failed to read test file: %v", err)
		}

		// Content should be exactly the same as before
		if string(content) != existingContent {
			t.Errorf("Expected content to remain unchanged, but got: %s", string(content))
		}
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr ||
		len(s) >= len(substr) && s[:len(substr)] == substr ||
		(len(s) > len(substr) && stringContains(s, substr))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
