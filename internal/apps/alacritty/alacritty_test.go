package alacritty

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
	app := &Alacritty{Cmd: mc}

	if err := app.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}
	if mc.InstalledDesktopApp != constants.Alacritty {
		t.Fatalf(
			"expected InstallDesktopApp(%s), got %q",
			constants.Alacritty,
			mc.InstalledDesktopApp,
		)
	}
}

// SKIP: ForceInstall test as per guidelines
// func TestForceInstall(t *testing.T) {
// 	mc := commands.NewMockCommand()
// 	app := &Alacritty{Cmd: mc}
//
// 	if err := app.ForceInstall(); err != nil {
// 		t.Fatalf("ForceInstall error: %v", err)
// 	}
// 	// ForceInstall should call Install() which uses InstallDesktopApp
// 	if mc.InstalledDesktopApp != constants.Alacritty {
// 		t.Fatalf("expected InstallDesktopApp(%s), got %q", constants.Alacritty, mc.InstalledDesktopApp)
// 	}
// }

func TestSoftInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Alacritty{Cmd: mc}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}
	if mc.MaybeInstalledDesktop != constants.Alacritty {
		t.Fatalf(
			"expected MaybeInstallDesktopApp(%s), got %q",
			constants.Alacritty,
			mc.MaybeInstalledDesktop,
		)
	}
}

// SKIP: Uninstall test as per guidelines
// func TestUninstall(t *testing.T) {
// 	mc := commands.NewMockCommand()
// 	app := &Alacritty{Cmd: mc}
//
// 	err := app.Uninstall()
// 	if err == nil {
// 		t.Fatal("expected Uninstall to return error for unsupported operation")
// 	}
// 	if err.Error() != "uninstall not implemented for alacritty" {
// 		t.Fatalf("unexpected error message: %v", err)
// 	}
// }

func TestForceConfigure(t *testing.T) {
	// Create temp directories for testing
	appDir := t.TempDir()
	localDir := t.TempDir()
	fontsDir := t.TempDir()
	themesDir := t.TempDir()

	// Override global paths for the duration of the test
	oldAppDir := paths.AlacrittyConfigAppDir
	oldLocalDir := paths.AlacrittyConfigLocalDir
	oldFontsDir := paths.AlacrittyFontsAppDir
	oldThemesDir := paths.AlacrittyThemesAppDir
	oldConfigDir := paths.ConfigDir

	paths.AlacrittyConfigAppDir = appDir
	paths.AlacrittyConfigLocalDir = localDir
	paths.AlacrittyFontsAppDir = fontsDir
	paths.AlacrittyThemesAppDir = themesDir
	paths.ConfigDir = "/test/config"

	t.Cleanup(func() {
		paths.AlacrittyConfigAppDir = oldAppDir
		paths.AlacrittyConfigLocalDir = oldLocalDir
		paths.AlacrittyFontsAppDir = oldFontsDir
		paths.AlacrittyThemesAppDir = oldThemesDir
		paths.ConfigDir = oldConfigDir
	})

	// Create test config files
	configContent := "window:\n  opacity: 0.9\n  option_as_alt: both"
	if err := os.WriteFile(filepath.Join(appDir, "alacritty.toml"), []byte(configContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create default font and theme directories with test files
	defaultFontDir := filepath.Join(fontsDir, "default")
	defaultThemeDir := filepath.Join(themesDir, "default")
	if err := os.MkdirAll(defaultFontDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(defaultThemeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	fontContent := "size = 14\nfamily = \"JetBrains Mono\""
	themeContent := "background = \"#1e1e1e\"\nforeground = \"#ffffff\""
	if err := os.WriteFile(filepath.Join(defaultFontDir, "font.toml"), []byte(fontContent), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(defaultThemeDir, "theme.toml"), []byte(themeContent), 0o644); err != nil {
		t.Fatal(err)
	}

	mc := commands.NewMockCommand()
	app := &Alacritty{Cmd: mc}

	if err := app.ForceConfigure(); err != nil {
		t.Fatalf("ForceConfigure error: %v", err)
	}

	// Check that config files were copied
	checkFiles := []string{"alacritty.toml", "font.toml", "theme.toml"}
	for _, file := range checkFiles {
		checkPath := filepath.Join(localDir, file)
		if _, err := os.Stat(checkPath); err != nil {
			t.Fatalf("expected copied file at %s: %v", checkPath, err)
		}
	}
}

func TestSoftConfigure(t *testing.T) {
	// Create temp directories for testing
	appDir := t.TempDir()
	localDir := t.TempDir()
	fontsDir := t.TempDir()
	themesDir := t.TempDir()

	// Override global paths for the duration of the test
	oldAppDir := paths.AlacrittyConfigAppDir
	oldLocalDir := paths.AlacrittyConfigLocalDir
	oldFontsDir := paths.AlacrittyFontsAppDir
	oldThemesDir := paths.AlacrittyThemesAppDir
	oldConfigDir := paths.ConfigDir

	paths.AlacrittyConfigAppDir = appDir
	paths.AlacrittyConfigLocalDir = localDir
	paths.AlacrittyFontsAppDir = fontsDir
	paths.AlacrittyThemesAppDir = themesDir
	paths.ConfigDir = "/test/config"

	t.Cleanup(func() {
		paths.AlacrittyConfigAppDir = oldAppDir
		paths.AlacrittyConfigLocalDir = oldLocalDir
		paths.AlacrittyFontsAppDir = oldFontsDir
		paths.AlacrittyThemesAppDir = oldThemesDir
		paths.ConfigDir = oldConfigDir
	})

	// Create test config files
	configContent := "window:\n  opacity: 0.9\n  option_as_alt: both"
	if err := os.WriteFile(filepath.Join(appDir, "alacritty.toml"), []byte(configContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create default font and theme directories with test files
	defaultFontDir := filepath.Join(fontsDir, "default")
	defaultThemeDir := filepath.Join(themesDir, "default")
	if err := os.MkdirAll(defaultFontDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(defaultThemeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	fontContent := "size = 14\nfamily = \"JetBrains Mono\""
	themeContent := "background = \"#1e1e1e\"\nforeground = \"#ffffff\""
	if err := os.WriteFile(filepath.Join(defaultFontDir, "font.toml"), []byte(fontContent), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(defaultThemeDir, "theme.toml"), []byte(themeContent), 0o644); err != nil {
		t.Fatal(err)
	}

	mc := commands.NewMockCommand()
	app := &Alacritty{Cmd: mc}

	// Test 1: Should configure when no existing config
	if err := app.SoftConfigure(); err != nil {
		t.Fatalf("SoftConfigure error: %v", err)
	}

	// Check that config files were copied
	checkFiles := []string{"alacritty.toml", "font.toml", "theme.toml"}
	for _, file := range checkFiles {
		checkPath := filepath.Join(localDir, file)
		if _, err := os.Stat(checkPath); err != nil {
			t.Fatalf("expected copied file at %s: %v", checkPath, err)
		}
	}

	// Test 2: Should not overwrite when config already exists
	// Modify the existing config
	modifiedContent := "window:\n  opacity: 0.5\n  option_as_alt: none"
	configPath := filepath.Join(localDir, "alacritty.toml")
	if err := os.WriteFile(configPath, []byte(modifiedContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Run SoftConfigure again
	if err := app.SoftConfigure(); err != nil {
		t.Fatalf("second SoftConfigure error: %v", err)
	}

	// Check that the file was not overwritten
	finalContent, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}
	if string(finalContent) != modifiedContent {
		t.Fatalf(
			"SoftConfigure overwrote existing config: expected %q, got %q",
			modifiedContent,
			string(finalContent),
		)
	}
}
