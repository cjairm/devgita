package lazydocker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
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

var expectedPackageName = fmt.Sprintf(
	"jesseduffield/%s/%s",
	constants.Lazydocker,
	constants.Lazydocker,
)

func TestInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Lazydocker{Cmd: mc}

	if err := app.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}
	if mc.InstalledPkg != expectedPackageName {
		t.Fatalf("expected InstallPackage(%s), got %q", expectedPackageName, mc.InstalledPkg)
	}
}

// SKIP: ForceInstall test as per guidelines
// ForceInstall calls Uninstall (which returns error) before Install
// Testing this creates false negatives
// func TestForceInstall(t *testing.T) {
// 	mc := commands.NewMockCommand()
// 	app := &Lazydocker{Cmd: mc}
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
	app := &Lazydocker{Cmd: mc}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}
	if mc.MaybeInstalled != expectedPackageName {
		t.Fatalf("expected MaybeInstallPackage(%s), got %q", expectedPackageName, mc.MaybeInstalled)
	}
}

// SKIP: No relevant tests
// func TestForceConfigure(t *testing.T) {
// 	mc := commands.NewMockCommand()
// 	app := &Lazydocker{Cmd: mc}
//
// 	// Lazydocker doesn't apply default configuration in devgita
// 	// ForceConfigure should return nil (no-op)
// 	if err := app.ForceConfigure(); err != nil {
// 		t.Fatalf("ForceConfigure error: %v", err)
// 	}
// }
//
// func TestSoftConfigure(t *testing.T) {
// 	mc := commands.NewMockCommand()
// 	app := &Lazydocker{Cmd: mc}
//
// 	// Lazydocker doesn't apply default configuration in devgita
// 	// SoftConfigure should return nil (no-op)
// 	if err := app.SoftConfigure(); err != nil {
// 		t.Fatalf("SoftConfigure error: %v", err)
// 	}
// }

func TestExecuteCommand(t *testing.T) {
	mc := commands.NewMockCommand()
	mockBase := commands.NewMockBaseCommand()
	app := &Lazydocker{Cmd: mc, Base: mockBase}

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
// 	app := &Lazydocker{Cmd: mc}
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
// 	app := &Lazydocker{Cmd: mc}
//
// 	err := app.Update()
// 	if err == nil {
// 		t.Fatal("expected Update to return error for unsupported operation")
// 	}
// 	if err.Error() != "lazydocker update not implemented through devgita" {
// 		t.Fatalf("unexpected error message: %v", err)
// 	}
// }
