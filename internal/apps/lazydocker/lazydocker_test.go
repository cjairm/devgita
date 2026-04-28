package lazydocker

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/testutil"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/downloader"
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

var expectedPackageName = fmt.Sprintf(
	"jesseduffield/%s/%s",
	constants.LazyDocker,
	constants.LazyDocker,
)

func TestInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	mb := commands.NewMockBaseCommand()
	mb.IsMacResult = true // simulate macOS so test uses Homebrew path
	app := &LazyDocker{Cmd: mc, Base: mb}

	if err := app.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}
	if mc.InstalledPkg != expectedPackageName {
		t.Fatalf("expected InstallPackage(%s), got %q", expectedPackageName, mc.InstalledPkg)
	}
}

func TestInstallDebian(t *testing.T) {
	mc := commands.NewMockCommand()
	mb := commands.NewMockBaseCommand()
	mb.IsMacResult = false // simulate Debian/Linux
	// Both ExecCommand calls (tar + install) succeed
	mb.SetExecCommandResult("", "", nil)

	app := &LazyDocker{
		Cmd:  mc,
		Base: mb,
		fetchVersion: func(owner, repo string) (string, error) {
			if owner != "jesseduffield" || repo != "lazydocker" {
				t.Errorf("unexpected version fetch: owner=%s repo=%s", owner, repo)
			}
			return "0.23.1", nil
		},
		downloadFn: func(_ context.Context, url, _ string, _ downloader.RetryConfig) error {
			if !strings.Contains(url, "0.23.1") {
				t.Errorf("download URL does not contain version: %s", url)
			}
			return nil
		},
	}

	if err := app.Install(); err != nil {
		t.Fatalf("Install (Debian) error: %v", err)
	}

	// Expect 2 ExecCommand calls: tar (extract) + sudo install
	if mb.GetExecCommandCallCount() != 2 {
		t.Fatalf("expected 2 ExecCommand calls, got %d", mb.GetExecCommandCallCount())
	}
	calls := mb.ExecCommandCalls
	if calls[0].Command != "tar" {
		t.Errorf("expected first command 'tar', got %q", calls[0].Command)
	}
	if calls[1].Command != "install" || !calls[1].IsSudo {
		t.Errorf("expected second command 'install' with IsSudo=true, got command=%q IsSudo=%v",
			calls[1].Command, calls[1].IsSudo)
	}
}

// SKIP: ForceInstall test as per guidelines
// ForceInstall calls Uninstall (which returns error) before Install
// Testing this creates false negatives
// func TestForceInstall(t *testing.T) {
// 	mc := commands.NewMockCommand()
// 	app := &LazyDocker{Cmd: mc}
//
// 	if err := app.ForceInstall(); err != nil {
// 		t.Fatalf("ForceInstall error: %v", err)
// 	}
// 	if mc.InstalledPkg != "lazydocker" {
// 		t.Fatalf("expected InstallPackage(%s), got %q", "lazydocker", mc.InstalledPkg)
// 	}
// }

func TestSoftInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	mb := commands.NewMockBaseCommand()
	mb.IsMacResult = true // macOS path uses MaybeInstallPackage
	app := &LazyDocker{Cmd: mc, Base: mb}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}
	if mc.MaybeInstalled != expectedPackageName {
		t.Fatalf("expected MaybeInstallPackage(%s), got %q", expectedPackageName, mc.MaybeInstalled)
	}
}

func TestSoftInstallDebian_AlreadyInstalled(t *testing.T) {
	mc := commands.NewMockCommand()
	mb := commands.NewMockBaseCommand()
	mb.IsMacResult = false

	orig := commands.LookPathFn
	commands.LookPathFn = func(string) (string, error) { return "/usr/local/bin/lazydocker", nil }
	defer func() { commands.LookPathFn = orig }()

	app := &LazyDocker{Cmd: mc, Base: mb}
	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall (already installed) error: %v", err)
	}
	if mb.GetExecCommandCallCount() != 0 {
		t.Fatalf("expected 0 ExecCommand calls, got %d", mb.GetExecCommandCallCount())
	}
}

func TestSoftInstallDebian_NotInstalled(t *testing.T) {
	mc := commands.NewMockCommand()
	mb := commands.NewMockBaseCommand()
	mb.IsMacResult = false
	mb.SetExecCommandResult("", "", nil)

	orig := commands.LookPathFn
	commands.LookPathFn = func(string) (string, error) { return "", fmt.Errorf("not found") }
	defer func() { commands.LookPathFn = orig }()

	app := &LazyDocker{
		Cmd:          mc,
		Base:         mb,
		fetchVersion: func(_, _ string) (string, error) { return "0.23.1", nil },
		downloadFn:   func(_ context.Context, _, _ string, _ downloader.RetryConfig) error { return nil },
	}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall (not installed) error: %v", err)
	}
	if mb.GetExecCommandCallCount() != 2 {
		t.Fatalf("expected 2 ExecCommand calls, got %d", mb.GetExecCommandCallCount())
	}
}

func TestForceConfigure(t *testing.T) {
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	app := &LazyDocker{Cmd: tc.MockApp.Cmd}

	// Test ForceConfigure - should enable shell feature
	err := app.ForceConfigure()
	if err != nil {
		t.Fatalf("ForceConfigure error: %v", err)
	}

	// Verify shell config was generated
	content, err := os.ReadFile(tc.ZshConfigPath)
	if err != nil {
		t.Fatalf("Failed to read shell config: %v", err)
	}

	if !strings.Contains(string(content), "# LazyDocker enabled") {
		t.Error("Expected shell config to contain LazyDocker feature")
	}

	testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
}

func TestSoftConfigure(t *testing.T) {
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()

	app := &LazyDocker{Cmd: tc.MockApp.Cmd}

	// First call should configure
	err := app.SoftConfigure()
	if err != nil {
		t.Fatalf("SoftConfigure error: %v", err)
	}

	// Verify shell config was generated
	content, err := os.ReadFile(tc.ZshConfigPath)
	if err != nil {
		t.Fatalf("Failed to read shell config: %v", err)
	}

	if !strings.Contains(string(content), "# LazyDocker enabled") {
		t.Error("Expected shell config to contain LazyDocker feature on first call")
	}

	// Second call should skip (feature already enabled)
	err = app.SoftConfigure()
	if err != nil {
		t.Fatalf("SoftConfigure should not error on second call: %v", err)
	}

	testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
}

func TestExecuteCommand(t *testing.T) {
	mc := commands.NewMockCommand()
	mockBase := commands.NewMockBaseCommand()
	app := &LazyDocker{Cmd: mc, Base: mockBase}

	// Test 1: Successful execution
	t.Run("successful execution", func(t *testing.T) {
		mockBase.SetExecCommandResult("lazydocker version 0.20.0", "", nil)

		err := app.ExecuteCommand("--version")
		if err != nil {
			t.Fatalf("ExecuteCommand failed: %v", err)
		}

		// Verify ExecCommand was called once
		if mockBase.GetExecCommandCallCount() != 1 {
			t.Fatalf("Expected 1 ExecCommand call, got %d", mockBase.GetExecCommandCallCount())
		}

		// Verify command parameters
		lastCall := mockBase.GetLastExecCommandCall()
		if lastCall == nil {
			t.Fatal("No ExecCommand call recorded")
		}
		if lastCall.Command != "lazydocker" {
			t.Fatalf("Expected command 'lazydocker', got %q", lastCall.Command)
		}
		if len(lastCall.Args) != 1 || lastCall.Args[0] != "--version" {
			t.Fatalf("Expected args ['--version'], got %v", lastCall.Args)
		}
		if lastCall.IsSudo {
			t.Fatal("Expected IsSudo to be false")
		}
	})

	// Test 2: Error handling
	t.Run("command execution error", func(t *testing.T) {
		mockBase.ResetExecCommand()
		mockBase.SetExecCommandResult(
			"",
			"command not found",
			fmt.Errorf("command not found: lazydocker"),
		)

		err := app.ExecuteCommand("--invalid-flag")
		if err == nil {
			t.Fatal("Expected ExecuteCommand to return error")
		}
		if !strings.Contains(err.Error(), "failed to run lazydocker command") {
			t.Fatalf("Expected error to contain 'failed to run lazydocker command', got: %v", err)
		}

		// Verify the error was properly wrapped
		if !strings.Contains(err.Error(), "command not found: lazydocker") {
			t.Fatalf("Expected error to contain original error message, got: %v", err)
		}
	})

	// Test 3: No arguments (launch TUI)
	t.Run("launch without arguments", func(t *testing.T) {
		mockBase.ResetExecCommand()
		mockBase.SetExecCommandResult("TUI launched", "", nil)

		err := app.ExecuteCommand()
		if err != nil {
			t.Fatalf("ExecuteCommand failed: %v", err)
		}

		lastCall := mockBase.GetLastExecCommandCall()
		if len(lastCall.Args) != 0 {
			t.Fatalf("Expected no args, got %v", lastCall.Args)
		}
	})
}

// SKIP: Uninstall test as per guidelines
// func TestUninstall(t *testing.T) {
// 	mc := commands.NewMockCommand()
// 	app := &LazyDocker{Cmd: mc}
//
// 	err := app.Uninstall()
// 	if err == nil {
// 		t.Fatal("expected Uninstall to return error for unsupported operation")
// 	}
// 	if err.Error() != "lazydocker uninstall not supported through devgita" {
// 		t.Fatalf("unexpected error message: %v", err)
// 	}
// }

// SKIP: Update test as per guidelines
// func TestUpdate(t *testing.T) {
// 	mc := commands.NewMockCommand()
// 	app := &LazyDocker{Cmd: mc}
//
// 	err := app.Update()
// 	if err == nil {
// 		t.Fatal("expected Update to return error for unsupported operation")
// 	}
// 	if err.Error() != "lazydocker update not implemented through devgita" {
// 		t.Fatalf("unexpected error message: %v", err)
// 	}
// }
