package neovim

import (
	"os"
	"path/filepath"
	"strings"
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
	app := &Neovim{Cmd: mc}

	if err := app.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}
	if mc.InstalledPkg != constants.Neovim {
		t.Fatalf("expected InstallPackage(%s), got %q", constants.Neovim, mc.InstalledPkg)
	}
}

func TestSoftInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Neovim{Cmd: mc}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}
	if mc.MaybeInstalled != constants.Neovim {
		t.Fatalf("expected MaybeInstallPackage(%s), got %q", constants.Neovim, mc.MaybeInstalled)
	}
}

func TestForceConfigure(t *testing.T) {
	// Create temp "app config" dir with a fake file as source
	src := t.TempDir()
	dst := t.TempDir()

	// Override global paths for the duration of the test
	oldAppDir, oldLocalDir := paths.Paths.App.Configs.Neovim, paths.Paths.Config.Nvim
	paths.Paths.App.Configs.Neovim, paths.Paths.Config.Nvim = src, dst
	t.Cleanup(func() {
		paths.Paths.App.Configs.Neovim, paths.Paths.Config.Nvim = oldAppDir, oldLocalDir
	})

	originalContent := "-- Neovim init.lua\nvim.g.mapleader = ' '"
	if err := os.WriteFile(filepath.Join(src, "init.lua"), []byte(originalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	mc := commands.NewMockCommand()
	app := &Neovim{Cmd: mc}

	// Test ForceConfigure - may succeed or fail depending on nvim availability
	err := app.ForceConfigure()

	if err != nil {
		// If nvim is not available, we expect a version check error
		if !strings.Contains(err.Error(), "failed to check Neovim version") {
			t.Fatalf("unexpected error (expected version check error): %v", err)
		}
		t.Logf("ForceConfigure failed as expected (no nvim binary): %v", err)
		return
	}

	// If nvim is available, configuration should succeed and files should be copied
	t.Logf("ForceConfigure succeeded (nvim binary available)")

	// Verify that the configuration was copied
	check := filepath.Join(dst, "init.lua")
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
}

func TestSoftConfigure(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	// Override global paths for the duration of the test
	oldAppDir, oldLocalDir := paths.Paths.App.Configs.Neovim, paths.Paths.Config.Nvim
	paths.Paths.App.Configs.Neovim, paths.Paths.Config.Nvim = src, dst
	t.Cleanup(func() {
		paths.Paths.App.Configs.Neovim, paths.Paths.Config.Nvim = oldAppDir, oldLocalDir
	})

	originalContent := "-- Neovim init.lua\nvim.g.mapleader = ' '"
	if err := os.WriteFile(filepath.Join(src, "init.lua"), []byte(originalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	mc := commands.NewMockCommand()
	app := &Neovim{Cmd: mc}

	// First call should attempt to configure
	err := app.SoftConfigure()

	if err != nil {
		// If nvim is not available, we expect a version check error
		if !strings.Contains(err.Error(), "failed to check Neovim version") {
			t.Fatalf("unexpected error (expected version check error): %v", err)
		}
		t.Logf("SoftConfigure failed as expected (no nvim binary): %v", err)
	} else {
		// If nvim is available, configuration should succeed
		t.Logf("SoftConfigure succeeded (nvim binary available)")

		// Verify that the configuration was copied
		check := filepath.Join(dst, "init.lua")
		if _, err := os.Stat(check); err != nil {
			t.Fatalf("expected copied file at %s: %v", check, err)
		}
	}

	// Create the marker file to simulate existing config (or ensure it exists)
	if err := os.MkdirAll(dst, 0o755); err != nil {
		t.Fatal(err)
	}
	markerContent := "-- Existing config"
	if err := os.WriteFile(filepath.Join(dst, "init.lua"), []byte(markerContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Second call should skip configuration since init.lua exists
	err = app.SoftConfigure()
	if err != nil {
		t.Fatalf("SoftConfigure should skip when init.lua exists: %v", err)
	}

	// Verify the existing file wasn't changed
	finalContent, err := os.ReadFile(filepath.Join(dst, "init.lua"))
	if err != nil {
		t.Fatalf("failed to read init.lua: %v", err)
	}
	if string(finalContent) != markerContent {
		t.Fatalf(
			"SoftConfigure overwrote existing file: expected %q, got %q",
			markerContent,
			string(finalContent),
		)
	}
}
