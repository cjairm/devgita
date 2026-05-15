package claude

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
	testutil.InitLogger()
}

func TestNew(t *testing.T) {
	app := New()
	if app == nil {
		t.Fatal("New() returned nil")
	}
	if app.Cmd == nil {
		t.Error("Expected Cmd to be initialized")
	}
	if app.Base == nil {
		t.Error("Expected Base to be initialized")
	}
}

func TestNameAndKind(t *testing.T) {
	app := &Claude{}
	if app.Name() != constants.Claude {
		t.Errorf("expected Name() %q, got %q", constants.Claude, app.Name())
	}
	if app.Kind() != apps.KindTerminal {
		t.Errorf("expected Kind() KindTerminal, got %v", app.Kind())
	}
}

func TestInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Claude{Cmd: mockApp.Cmd, Base: mockApp.Base}

	if err := app.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}

	last := mockApp.Base.GetLastExecCommandCall()
	if last == nil {
		t.Fatal("Expected ExecCommand to be called")
	}
	if last.Command != "sh" {
		t.Errorf("Expected command 'sh', got %q", last.Command)
	}
}

func TestForceInstall(t *testing.T) {
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	claudeConfigDir := filepath.Join(tc.ConfigDir, ".claude")
	oldClaudeDir := paths.Paths.Config.Claude
	t.Cleanup(func() { paths.Paths.Config.Claude = oldClaudeDir })
	paths.Paths.Config.Claude = claudeConfigDir

	app := &Claude{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

	if err := app.ForceInstall(); err != nil {
		t.Fatalf("ForceInstall error: %v", err)
	}

	// ForceInstall runs Uninstall (sh rm) then Install (sh curl)
	// Both use Base.ExecCommand, so we expect 2 calls
	calls := tc.MockApp.Base.ExecCommandCalls
	if len(calls) < 2 {
		t.Fatalf("expected at least 2 ExecCommand calls, got %d", len(calls))
	}
	if calls[0].Command != "sh" {
		t.Errorf("expected first command 'sh' (uninstall), got %q", calls[0].Command)
	}
	last := calls[len(calls)-1]
	if last.Command != "sh" {
		t.Errorf("Expected install script via 'sh', got %q", last.Command)
	}
}

func TestUninstall(t *testing.T) {
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	claudeConfigDir := filepath.Join(tc.ConfigDir, ".claude")
	if err := os.MkdirAll(claudeConfigDir, 0755); err != nil {
		t.Fatal(err)
	}
	oldClaudeDir := paths.Paths.Config.Claude
	t.Cleanup(func() { paths.Paths.Config.Claude = oldClaudeDir })
	paths.Paths.Config.Claude = claudeConfigDir

	app := &Claude{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

	if err := app.Uninstall(); err != nil {
		t.Fatalf("Uninstall error: %v", err)
	}

	// Verify sh rm uninstall was called
	last := tc.MockApp.Base.GetLastExecCommandCall()
	if last == nil || last.Command != "sh" {
		t.Fatalf("expected sh uninstall command, got %v", last)
	}
	if len(last.Args) < 2 || last.Args[1] != "rm -f ~/.local/bin/claude && rm -rf ~/.local/share/claude" {
		t.Errorf("expected sh -c rm uninstall command, got args %v", last.Args)
	}

	// Config dir should be removed
	if _, err := os.Stat(claudeConfigDir); err == nil {
		t.Error("expected claude config dir to be removed")
	}
}

func TestUpdate(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Claude{Cmd: mockApp.Cmd, Base: mockApp.Base}

	err := app.Update()
	if err == nil {
		t.Fatal("Expected Update to return error")
	}
	if !errors.Is(err, apps.ErrUpdateNotSupported) {
		t.Errorf("expected ErrUpdateNotSupported, got: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestForceConfigure(t *testing.T) {
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	appConfigDir := filepath.Join(tc.AppDir, "configs", "claude")
	userConfigDir := filepath.Join(tc.ConfigDir, "..", ".claude")

	// Create source structure
	for _, sub := range []string{"themes"} {
		if err := os.MkdirAll(filepath.Join(appConfigDir, sub), 0755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(appConfigDir, "settings.json"), []byte(`{"theme":"default"}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(appConfigDir, "statusline.sh"), []byte(`#!/bin/bash`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(appConfigDir, "themes", "default.json"), []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}

	// Create shared dirs
	sharedDir := filepath.Join(tc.AppDir, "configs", "shared")
	for _, sub := range []string{"skills", "commands", "agents"} {
		if err := os.MkdirAll(filepath.Join(sharedDir, sub), 0755); err != nil {
			t.Fatal(err)
		}
	}

	paths.Paths.App.Configs.Claude = appConfigDir
	paths.Paths.App.Configs.Shared = sharedDir
	paths.Paths.Config.Claude = userConfigDir

	app := &Claude{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

	if err := app.ForceConfigure(); err != nil {
		t.Fatalf("ForceConfigure error: %v", err)
	}

	// settings.json deployed
	if _, err := os.Stat(filepath.Join(userConfigDir, "settings.json")); err != nil {
		t.Errorf("Expected settings.json at %s: %v", userConfigDir, err)
	}

	// statusline.sh deployed and executable
	statuslinePath := filepath.Join(userConfigDir, "statusline.sh")
	info, err := os.Stat(statuslinePath)
	if err != nil {
		t.Fatalf("Expected statusline.sh: %v", err)
	}
	if info.Mode()&0111 == 0 {
		t.Error("Expected statusline.sh to be executable")
	}

	// themes deployed
	if _, err := os.Stat(filepath.Join(userConfigDir, "themes")); err != nil {
		t.Errorf("Expected themes dir: %v", err)
	}

	// shared dirs deployed
	for _, dir := range []string{"skills", "commands", "agents"} {
		if _, err := os.Stat(filepath.Join(userConfigDir, dir)); err != nil {
			t.Errorf("Expected %s dir: %v", dir, err)
		}
	}

	testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
	testutil.VerifyNoRealConfigChanges(t)
}

func TestSoftConfigure_AlreadyConfigured(t *testing.T) {
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	userConfigDir := filepath.Join(tc.ConfigDir, "..", ".claude")
	paths.Paths.Config.Claude = userConfigDir

	// Pre-create marker file so SoftConfigure skips
	if err := os.MkdirAll(userConfigDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(userConfigDir, "settings.json"), []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}

	app := &Claude{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

	if err := app.SoftConfigure(); err != nil {
		t.Fatalf("SoftConfigure error: %v", err)
	}

	testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
}

func TestExecuteCommand(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Claude{Cmd: mockApp.Cmd, Base: mockApp.Base}

	if err := app.ExecuteCommand("--version"); err != nil {
		t.Fatalf("ExecuteCommand error: %v", err)
	}

	last := mockApp.Base.GetLastExecCommandCall()
	if last == nil || last.Command != constants.Claude {
		t.Errorf("Expected command %q, got %v", constants.Claude, last)
	}
	if len(last.Args) == 0 || last.Args[0] != "--version" {
		t.Errorf("Expected args [--version], got %v", last.Args)
	}
}
