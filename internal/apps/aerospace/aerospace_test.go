package aerospace

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
	app := &Aerospace{Cmd: mc}

	if err := app.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}
	if mc.InstalledDesktopApp != "nikitabobko/tap/aerospace" {
		t.Fatalf(
			"expected InstallDesktopApp(nikitabobko/tap/aerospace), got %q",
			mc.InstalledDesktopApp,
		)
	}
}

// SKIP: ForceInstall test as per guidelines
// func TestForceInstall(t *testing.T) {
// 	mc := commands.NewMockCommand()
// 	app := &Aerospace{Cmd: mc}
//
// 	if err := app.ForceInstall(); err != nil {
// 		t.Fatalf("ForceInstall error: %v", err)
// 	}
// 	// ForceInstall should call Install() which uses InstallDesktopApp
// 	if mc.InstalledDesktopApp != "nikitabobko/tap/aerospace" {
// 		t.Fatalf("expected InstallDesktopApp(nikitabobko/tap/aerospace), got %q", mc.InstalledDesktopApp)
// 	}
// }

func TestSoftInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Aerospace{Cmd: mc}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}
	if mc.MaybeInstalledDesktop != "nikitabobko/tap/aerospace" {
		t.Fatalf(
			"expected MaybeInstallDesktopApp(nikitabobko/tap/aerospace), got %q",
			mc.MaybeInstalledDesktop,
		)
	}
}

func TestUninstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Aerospace{Cmd: mc}

	err := app.Uninstall()
	if err == nil {
		t.Fatal("expected Uninstall to return error for unsupported operation")
	}
	if err.Error() != "aerospace uninstall not supported through devgita" {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestForceConfigure(t *testing.T) {
	// Create temp "app config" dir with a fake file as source
	src := t.TempDir()
	dst := t.TempDir()

	// Override global paths for the duration of the test
	oldAppDir, oldLocalDir := paths.AerospaceConfigAppDir, paths.AerospaceConfigLocalDir
	paths.AerospaceConfigAppDir, paths.AerospaceConfigLocalDir = src, dst
	t.Cleanup(func() {
		paths.AerospaceConfigAppDir, paths.AerospaceConfigLocalDir = oldAppDir, oldLocalDir
	})

	originalContent := "[workspace]\nkey = \"value\""
	if err := os.WriteFile(filepath.Join(src, "aerospace.toml"), []byte(originalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	mc := commands.NewMockCommand()
	app := &Aerospace{Cmd: mc}

	if err := app.ForceConfigure(); err != nil {
		t.Fatalf("ForceConfigure error: %v", err)
	}

	check := filepath.Join(dst, "aerospace.toml")
	if _, err := os.Stat(check); err != nil {
		t.Fatalf("expected copied file at %s: %v", check, err)
	}

	copiedContent, err := os.ReadFile(check)
	if err != nil {
		t.Fatalf("failed to read copied file: %v", err)
	}
	if string(copiedContent) != originalContent {
		t.Fatalf("content mismatch: expected %q, got %q", originalContent, string(copiedContent))
	}

	modifiedContent := "[workspace]\nkey = \"modified\""
	if err := os.WriteFile(check, []byte(modifiedContent), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := app.ForceConfigure(); err != nil {
		t.Fatalf("second ForceConfigure error: %v", err)
	}

	finalContent, err := os.ReadFile(check)
	if err != nil {
		t.Fatalf("failed to read file after second configure: %v", err)
	}
	if string(finalContent) == string(modifiedContent) {
		t.Fatalf(
			"ForceConfigure did not overwrite: expected %q, got %q",
			originalContent,
			string(finalContent),
		)
	}
}

func TestSoftConfigure(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	// Override global paths for the duration of the test
	oldAppDir, oldLocalDir := paths.AerospaceConfigAppDir, paths.AerospaceConfigLocalDir
	paths.AerospaceConfigAppDir, paths.AerospaceConfigLocalDir = src, dst
	t.Cleanup(func() {
		paths.AerospaceConfigAppDir, paths.AerospaceConfigLocalDir = oldAppDir, oldLocalDir
	})

	originalContent := "[workspace]\nkey = \"value\""
	if err := os.WriteFile(filepath.Join(src, "aerospace.toml"), []byte(originalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	mc := commands.NewMockCommand()
	app := &Aerospace{Cmd: mc}

	if err := app.SoftConfigure(); err != nil {
		t.Fatalf("SoftConfigure error: %v", err)
	}

	check := filepath.Join(dst, "aerospace.toml")
	if _, err := os.Stat(check); err != nil {
		t.Fatalf("expected copied file at %s: %v", check, err)
	}

	modifiedContent := "[workspace]\nkey = \"modified\""
	if err := os.WriteFile(check, []byte(modifiedContent), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := app.SoftConfigure(); err != nil {
		t.Fatalf("second SoftConfigure error: %v", err)
	}

	finalContent, err := os.ReadFile(check)
	if err != nil {
		t.Fatalf("failed to read file after second configure: %v", err)
	}
	if string(finalContent) == string(originalContent) {
		t.Fatalf(
			"SoftConfigure overwrote existing file: expected %q, got %q",
			modifiedContent,
			string(finalContent),
		)
	}
}

// SKIP: Uninstall test

// SKIP: Updates test
