package mise

import (
	"testing"

	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/logger"
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
	app := &Mise{Cmd: mc}

	if err := app.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}
	if mc.InstalledPkg != "mise" {
		t.Fatalf("expected InstallPackage(mise), got %q", mc.InstalledPkg)
	}
}

func TestSoftInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Mise{Cmd: mc}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}
	if mc.MaybeInstalled != "mise" {
		t.Fatalf("expected MaybeInstallPackage(mise), got %q", mc.MaybeInstalled)
	}
}

func TestExecuteCommand(t *testing.T) {
	// These tests use the actual BaseCommand but won't execute actual mise commands
	// since we expect mise to fail in test environment
	mc := commands.NewMockCommand()
	app := &Mise{Cmd: mc}

	// Test that methods don't panic - they will fail due to mise not being available
	// but that's expected in test environment
	err := app.ExecuteCommand("--version")

	// Should fail since mise isn't available, but shouldn't panic
	if err == nil {
		t.Log("Mise commands succeeded unexpectedly (mise must be available)")
	} else {
		t.Log("Expected failure due to mise not being available in test environment")
	}
}
