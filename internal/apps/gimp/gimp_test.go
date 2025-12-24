package gimp

import (
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/testutil"
)

func init() {
	testutil.InitLogger()
}

func TestNew(t *testing.T) {
	gimp := New()

	if gimp == nil {
		t.Fatal("New() returned nil")
	}

	if gimp.Cmd == nil {
		t.Error("Expected Cmd to be initialized")
	}
}

func TestInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	gimp := &Gimp{Cmd: mockApp.Cmd}

	err := gimp.Install()
	if err != nil {
		t.Fatalf("Install() failed: %v", err)
	}

	if mockApp.Cmd.InstalledDesktopApp == "" {
		t.Error("Expected GIMP to be installed")
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestForceInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	gimp := &Gimp{Cmd: mockApp.Cmd}

	err := gimp.ForceInstall()

	if err == nil {
		t.Fatal("Expected ForceInstall() to fail when Uninstall() is not supported")
	}

	if !strings.Contains(err.Error(), "uninstall") {
		t.Errorf("Expected error to mention uninstall, got: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestSoftInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	gimp := &Gimp{Cmd: mockApp.Cmd}

	t.Run("installs when not present", func(t *testing.T) {
		mockApp.Reset()
		mockApp.Cmd.DesktopAppInstalled = false

		err := gimp.SoftInstall()
		if err != nil {
			t.Fatalf("SoftInstall() failed: %v", err)
		}

		if mockApp.Cmd.MaybeInstalledDesktop == "" {
			t.Error("Expected MaybeInstallDesktopApp to be called")
		}

		testutil.VerifyNoRealCommands(t, mockApp.Base)
	})

	t.Run("skips when already present", func(t *testing.T) {
		mockApp.Reset()
		mockApp.Cmd.DesktopAppInstalled = true

		err := gimp.SoftInstall()
		if err != nil {
			t.Fatalf("SoftInstall() failed: %v", err)
		}

		testutil.VerifyNoRealCommands(t, mockApp.Base)
	})
}

func TestForceConfigure(t *testing.T) {
	mockApp := testutil.NewMockApp()
	gimp := &Gimp{Cmd: mockApp.Cmd}

	err := gimp.ForceConfigure()
	if err != nil {
		t.Fatalf("ForceConfigure() failed: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestSoftConfigure(t *testing.T) {
	mockApp := testutil.NewMockApp()
	gimp := &Gimp{Cmd: mockApp.Cmd}

	err := gimp.SoftConfigure()
	if err != nil {
		t.Fatalf("SoftConfigure() failed: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestUninstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	gimp := &Gimp{Cmd: mockApp.Cmd}

	err := gimp.Uninstall()

	if err == nil {
		t.Fatal("Expected Uninstall() to return error")
	}

	if !strings.Contains(err.Error(), "not supported") {
		t.Errorf("Expected error to mention 'not supported', got: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestExecuteCommand(t *testing.T) {
	mockApp := testutil.NewMockApp()
	gimp := &Gimp{Cmd: mockApp.Cmd}

	err := gimp.ExecuteCommand("--version")
	if err != nil {
		t.Fatalf("ExecuteCommand() failed: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestUpdate(t *testing.T) {
	mockApp := testutil.NewMockApp()
	gimp := &Gimp{Cmd: mockApp.Cmd}

	err := gimp.Update()

	if err == nil {
		t.Fatal("Expected Update() to return error")
	}

	if !strings.Contains(err.Error(), "not implemented") {
		t.Errorf("Expected error to mention 'not implemented', got: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}
