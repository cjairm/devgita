package docker

import (
	"fmt"
	"testing"

	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/testutil"
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

func TestInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Docker{Cmd: mc}

	if err := app.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}
	if mc.InstalledDesktopApp != constants.Docker {
		t.Fatalf(
			"expected InstallDesktopApp(%s), got %q",
			constants.Docker,
			mc.InstalledDesktopApp,
		)
	}
}

// SKIP: ForceInstall test as per guidelines
// func TestForceInstall(t *testing.T) {
// 	mc := commands.NewMockCommand()
// 	app := &Docker{Cmd: mc}
//
// 	if err := app.ForceInstall(); err != nil {
// 		t.Fatalf("ForceInstall error: %v", err)
// 	}
// 	// ForceInstall should call Install() which uses InstallDesktopApp
// 	if mc.InstalledDesktopApp != constants.Docker {
// 		t.Fatalf("expected InstallDesktopApp(%s), got %q", constants.Docker, mc.InstalledDesktopApp)
// 	}
// }

func TestSoftInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Docker{Cmd: mc}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}
	if mc.MaybeInstalledDesktop != constants.Docker {
		t.Fatalf(
			"expected MaybeInstallDesktopApp(%s), got %q",
			constants.Docker,
			mc.MaybeInstalledDesktop,
		)
	}
}

// SKIP: Uninstall test as per guidelines
// func TestUninstall(t *testing.T) {
// 	mc := commands.NewMockCommand()
// 	app := &Docker{Cmd: mc}
//
// 	err := app.Uninstall()
// 	if err == nil {
// 		t.Fatal("expected Uninstall to return error for unsupported operation")
// 	}
// 	if !strings.Contains(err.Error(), "uninstall not implemented") {
// 		t.Fatalf("unexpected error message: %v", err)
// 	}
// }

func TestForceConfigure(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Docker{Cmd: mc}

	if err := app.ForceConfigure(); err != nil {
		t.Fatalf("ForceConfigure error: %v", err)
	}

	// Docker doesn't require configuration files, so this should succeed without operations
}

func TestSoftConfigure(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Docker{Cmd: mc}

	if err := app.SoftConfigure(); err != nil {
		t.Fatalf("SoftConfigure error: %v", err)
	}

	// Docker doesn't require configuration files, so this should succeed without operations
}

func TestExecuteCommand(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Docker{Cmd: mockApp.Cmd, Base: mockApp.Base}

	t.Run("successful execution", func(t *testing.T) {
		mockApp.Base.SetExecCommandResult("Docker version 24.0.0", "", nil)

		err := app.ExecuteCommand("--version")
		if err != nil {
			t.Fatalf("ExecuteCommand failed: %v", err)
		}

		// Verify ExecCommand was called once
		if mockApp.Base.GetExecCommandCallCount() != 1 {
			t.Fatalf("Expected 1 ExecCommand call, got %d", mockApp.Base.GetExecCommandCallCount())
		}

		// Verify command parameters
		lastCall := mockApp.Base.GetLastExecCommandCall()
		if lastCall == nil {
			t.Fatal("No ExecCommand call recorded")
		}
		if lastCall.Command != constants.Docker {
			t.Fatalf("Expected command %s, got %q", constants.Docker, lastCall.Command)
		}
		if len(lastCall.Args) != 1 || lastCall.Args[0] != "--version" {
			t.Fatalf("Expected args [--version], got %v", lastCall.Args)
		}
	})

	t.Run("error handling", func(t *testing.T) {
		mockApp.Reset()
		mockApp.Base.SetExecCommandResult("", "command failed", fmt.Errorf("exit status 1"))

		err := app.ExecuteCommand("--invalid")
		if err == nil {
			t.Fatal("Expected error from ExecuteCommand")
		}
	})
}

// SKIP: Update test as per guidelines
// func TestUpdate(t *testing.T) {
// 	mc := commands.NewMockCommand()
// 	app := &Docker{Cmd: mc}
//
// 	err := app.Update()
// 	if err == nil {
// 		t.Fatal("expected Update to return error for unsupported operation")
// 	}
// 	if err.Error() != "update not implemented for docker" {
// 		t.Fatalf("unexpected error message: %v", err)
// 	}
// }
