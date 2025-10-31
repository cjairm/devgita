package powerlevel10k

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
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
	app := &PowerLevel10k{Cmd: mc}

	if err := app.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}
	if mc.InstalledPkg != constants.Powerlevel10k {
		t.Fatalf(
			"expected InstallPackage(%s), got %q",
			constants.Powerlevel10k,
			mc.InstalledPkg,
		)
	}
}

func TestSoftInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &PowerLevel10k{Cmd: mc}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}
	if mc.MaybeInstalled != constants.Powerlevel10k {
		t.Fatalf(
			"expected MaybeInstallPackage(%s), got %q",
			constants.Powerlevel10k,
			mc.MaybeInstalled,
		)
	}
}

func TestForceConfigure(t *testing.T) {
	// Create temp directory for testing
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

	mc := commands.NewMockCommand()
	app := &PowerLevel10k{Cmd: mc}

	if err := app.ForceConfigure(); err != nil {
		t.Fatalf("ForceConfigure error: %v", err)
	}

	// Verify content was added to the file
	content, err := os.ReadFile(zshFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	expectedContent := "source $(brew --prefix)/share/powerlevel10k/powerlevel10k.zsh-theme"
	if !contains(string(content), expectedContent) {
		t.Fatalf(
			"Expected file to contain %q, but it didn't. Content: %s",
			expectedContent,
			string(content),
		)
	}
}

func TestSoftConfigure(t *testing.T) {
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

		mc := commands.NewMockCommand()
		app := &PowerLevel10k{Cmd: mc}

		if err := app.SoftConfigure(); err != nil {
			t.Fatalf("SoftConfigure error: %v", err)
		}

		// Verify content was added to the file
		content, err := os.ReadFile(zshFile)
		if err != nil {
			t.Fatalf("Failed to read test file: %v", err)
		}

		expectedContent := "source $(brew --prefix)/share/powerlevel10k/powerlevel10k.zsh-theme"
		if !contains(string(content), expectedContent) {
			t.Fatalf(
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

		// Create the devgita.zsh file with existing powerlevel10k config
		zshFile := filepath.Join(tempDir, "devgita.zsh")
		existingContent := "# Some existing content\nsource $(brew --prefix)/share/powerlevel10k/powerlevel10k.zsh-theme\n# More content"
		if err := os.WriteFile(zshFile, []byte(existingContent), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		mc := commands.NewMockCommand()
		app := &PowerLevel10k{Cmd: mc}

		if err := app.SoftConfigure(); err != nil {
			t.Fatalf("SoftConfigure error: %v", err)
		}

		// Verify content wasn't duplicated
		content, err := os.ReadFile(zshFile)
		if err != nil {
			t.Fatalf("Failed to read test file: %v", err)
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
