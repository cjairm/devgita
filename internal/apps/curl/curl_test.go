package curl

import (
	"strings"
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
	app := &Curl{Cmd: mc}

	if err := app.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}
	if mc.InstalledPkg != "curl" {
		t.Fatalf("expected InstallPackage(%s), got %q", "curl", mc.InstalledPkg)
	}
}

// SKIP: ForceInstall test as per guidelines
// func TestForceInstall(t *testing.T) {
// 	mc := commands.NewMockCommand()
// 	app := &Curl{Cmd: mc}
//
// 	if err := app.ForceInstall(); err != nil {
// 		t.Fatalf("ForceInstall error: %v", err)
// 	}
// 	// ForceInstall should call Install() which uses InstallPackage
// 	if mc.InstalledPkg != "curl" {
// 		t.Fatalf("expected InstallPackage(%s), got %q", "curl", mc.InstalledPkg)
// 	}
// }

func TestSoftInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Curl{Cmd: mc}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}
	if mc.MaybeInstalled != "curl" {
		t.Fatalf("expected MaybeInstallPackage(%s), got %q", "curl", mc.MaybeInstalled)
	}
}

func TestForceConfigure(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Curl{Cmd: mc}

	// Curl doesn't require separate configuration files
	// so ForceConfigure should return nil (no-op)
	if err := app.ForceConfigure(); err != nil {
		t.Fatalf("ForceConfigure error: %v", err)
	}
}

func TestSoftConfigure(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Curl{Cmd: mc}

	// Curl doesn't require separate configuration files
	// so SoftConfigure should return nil (no-op)
	if err := app.SoftConfigure(); err != nil {
		t.Fatalf("SoftConfigure error: %v", err)
	}
}

func TestExecuteCommand(t *testing.T) {
	// Note: In an ideal architecture, we would mock BaseCommand.ExecCommand
	// to avoid running actual system commands in tests. Currently, the BaseCommand
	// struct doesn't implement an interface that we can easily mock.
	//
	// For now, this test creates a Curl instance with MockCommand (which handles
	// the Command interface methods) and a real BaseCommand. Since curl is likely
	// not available in the test environment, the ExecCommand call will fail,
	// which is the expected behavior for this test setup.
	//
	// TODO: Create an interface for BaseCommand to enable proper mocking of ExecCommand

	mc := commands.NewMockCommand()
	app := &Curl{Cmd: mc}

	// Test that the method properly handles command execution
	// This will likely fail since curl isn't available in test environment,
	// but it verifies the method doesn't panic and properly wraps errors
	err := app.ExecuteCommand("--version")

	// Verify error handling - should fail gracefully without panicking
	if err == nil {
		t.Log("Curl command succeeded unexpectedly (curl must be available in test environment)")
	} else {
		// Expected case: curl not available, should get a wrapped error
		if !strings.Contains(err.Error(), "failed to run curl command") {
			t.Fatalf("Expected error to contain 'failed to run curl command', got: %v", err)
		}
		t.Logf("Curl command failed as expected (no curl available): %v", err)
	}
}

// SKIP: Uninstall test as per guidelines
// func TestUninstall(t *testing.T) {
// 	mc := commands.NewMockCommand()
// 	app := &Curl{Cmd: mc}
//
// 	err := app.Uninstall()
// 	if err == nil {
// 		t.Fatal("expected Uninstall to return error for unsupported operation")
// 	}
// 	if err.Error() != "curl uninstall not supported through devgita" {
// 		t.Fatalf("unexpected error message: %v", err)
// 	}
// }

// SKIP: Update test as per guidelines
// func TestUpdate(t *testing.T) {
// 	mc := commands.NewMockCommand()
// 	app := &Curl{Cmd: mc}
//
// 	err := app.Update()
// 	if err == nil {
// 		t.Fatal("expected Update to return error for unsupported operation")
// 	}
// 	if err.Error() != "curl update not implemented through devgita" {
// 		t.Fatalf("unexpected error message: %v", err)
// 	}
// }
