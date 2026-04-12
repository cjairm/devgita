package neovim

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/testutil"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/downloader"
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

func TestInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	mb := commands.NewMockBaseCommand()
	mb.IsMacResult = true // simulate macOS so test uses Homebrew path
	app := &Neovim{Cmd: mc, Base: mb}

	if err := app.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}
	if mc.InstalledPkg != constants.Neovim {
		t.Fatalf("expected InstallPackage(%s), got %q", constants.Neovim, mc.InstalledPkg)
	}
}

func TestSoftInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	mb := commands.NewMockBaseCommand()
	mb.IsMacResult = true // simulate macOS so test uses Homebrew path
	app := &Neovim{Cmd: mc, Base: mb}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}
	if mc.MaybeInstalled != constants.Neovim {
		t.Fatalf("expected MaybeInstallPackage(%s), got %q", constants.Neovim, mc.MaybeInstalled)
	}
}

func TestInstallDebian(t *testing.T) {
	mc := commands.NewMockCommand()
	mb := commands.NewMockBaseCommand()
	mb.IsMacResult = false // simulate Debian/Linux
	// All ExecCommand calls (tar + install + 2×cp) succeed
	mb.SetExecCommandResult("", "", nil)

	app := &Neovim{
		Cmd:  mc,
		Base: mb,
		downloadFn: func(_ context.Context, url, _ string, _ downloader.RetryConfig) error {
			if !strings.Contains(url, constants.SupportedVersion.Neovim.Number) {
				t.Errorf("download URL does not contain version %s: %s",
					constants.SupportedVersion.Neovim.Number, url)
			}
			return nil
		},
	}

	if err := app.Install(); err != nil {
		t.Fatalf("Install (Debian) error: %v", err)
	}

	// Expect 4 ExecCommand calls: tar + sudo install + sudo cp lib + sudo cp share
	if mb.GetExecCommandCallCount() != 4 {
		t.Fatalf("expected 4 ExecCommand calls, got %d", mb.GetExecCommandCallCount())
	}
	calls := mb.ExecCommandCalls
	if calls[0].Command != "tar" {
		t.Errorf("expected first command 'tar', got %q", calls[0].Command)
	}
	if calls[1].Command != "install" || !calls[1].IsSudo {
		t.Errorf("expected second command 'install' with IsSudo=true, got command=%q IsSudo=%v",
			calls[1].Command, calls[1].IsSudo)
	}
	if calls[2].Command != "cp" || !calls[2].IsSudo {
		t.Errorf("expected third command 'cp' with IsSudo=true, got command=%q IsSudo=%v",
			calls[2].Command, calls[2].IsSudo)
	}
	if calls[3].Command != "cp" || !calls[3].IsSudo {
		t.Errorf("expected fourth command 'cp' with IsSudo=true, got command=%q IsSudo=%v",
			calls[3].Command, calls[3].IsSudo)
	}
}

func TestForceConfigure(t *testing.T) {
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	// Create neovim config source
	src := filepath.Join(tc.AppDir, "neovim")
	dst := filepath.Join(tc.ConfigDir, "nvim")

	if err := os.MkdirAll(src, 0755); err != nil {
		t.Fatal(err)
	}

	// Override global paths
	paths.Paths.App.Configs.Neovim = src
	paths.Paths.Config.Nvim = dst

	originalContent := "-- Neovim init.lua\nvim.g.mapleader = ' '"
	if err := os.WriteFile(filepath.Join(src, "init.lua"), []byte(originalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Mock successful version check with version output
	tc.MockApp.Base.SetExecCommandResult("NVIM v0.11.1\nBuild type: Release", "", nil)
	app := &Neovim{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

	// Test ForceConfigure - should succeed with mocked version check
	err := app.ForceConfigure()

	if err != nil {
		t.Fatalf("ForceConfigure failed: %v", err)
	}

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

	// Verify shell config was generated
	shellContent, err := os.ReadFile(tc.ZshConfigPath)
	if err != nil {
		t.Fatalf("Failed to read shell config: %v", err)
	}

	if !strings.Contains(string(shellContent), "# Neovim enabled") {
		t.Error("Expected shell config to contain Neovim feature")
	}

	// Verify version check was called once (this is expected)
	if tc.MockApp.Base.GetExecCommandCallCount() != 1 {
		t.Errorf("Expected 1 command call (version check), got %d", tc.MockApp.Base.GetExecCommandCallCount())
	}
}

func TestSoftConfigure(t *testing.T) {
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	// Create neovim config source
	src := filepath.Join(tc.AppDir, "neovim")
	dst := filepath.Join(tc.ConfigDir, "nvim")

	if err := os.MkdirAll(src, 0755); err != nil {
		t.Fatal(err)
	}

	// Override global paths
	paths.Paths.App.Configs.Neovim = src
	paths.Paths.Config.Nvim = dst

	originalContent := "-- Neovim init.lua\nvim.g.mapleader = ' '"
	if err := os.WriteFile(filepath.Join(src, "init.lua"), []byte(originalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Mock successful version check with version output
	tc.MockApp.Base.SetExecCommandResult("NVIM v0.11.1\nBuild type: Release", "", nil)
	app := &Neovim{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

	// First call should attempt to configure
	err := app.SoftConfigure()

	if err != nil {
		t.Fatalf("SoftConfigure failed: %v", err)
	}

	// Verify that the configuration was copied
	check := filepath.Join(dst, "init.lua")
	if _, err := os.Stat(check); err != nil {
		t.Fatalf("expected copied file at %s: %v", check, err)
	}

	// Verify shell config was generated
	shellContent, err := os.ReadFile(tc.ZshConfigPath)
	if err != nil {
		t.Fatalf("Failed to read shell config: %v", err)
	}

	if !strings.Contains(string(shellContent), "# Neovim enabled") {
		t.Error("Expected shell config to contain Neovim feature on first call")
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
	// But it should still enable the shell feature if not already enabled
	err = app.SoftConfigure()
	if err != nil {
		t.Fatalf("SoftConfigure should not error on second call: %v", err)
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

	// Verify version check was called once in first SoftConfigure (this is expected)
	if tc.MockApp.Base.GetExecCommandCallCount() != 1 {
		t.Errorf("Expected 1 command call (version check), got %d", tc.MockApp.Base.GetExecCommandCallCount())
	}
}
