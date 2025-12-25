package flameshot

import (
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/testutil"
)

func init() {
	testutil.InitLogger()
}

func TestNew(t *testing.T) {
	flameshot := New()

	if flameshot == nil {
		t.Fatal("New() returned nil")
	}

	if flameshot.Cmd == nil {
		t.Error("Expected Cmd to be initialized")
	}
}

func TestInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	flameshot := &Flameshot{Cmd: mockApp.Cmd}

	err := flameshot.Install()
	if err != nil {
		t.Fatalf("Install() failed: %v", err)
	}

	if mockApp.Cmd.InstalledDesktopApp == "" {
		t.Error("Expected Flameshot to be installed")
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestForceInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	flameshot := &Flameshot{Cmd: mockApp.Cmd}

	err := flameshot.ForceInstall()

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
	flameshot := &Flameshot{Cmd: mockApp.Cmd}

	t.Run("installs when not present", func(t *testing.T) {
		mockApp.Reset()
		mockApp.Cmd.DesktopAppInstalled = false

		err := flameshot.SoftInstall()
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

		err := flameshot.SoftInstall()
		if err != nil {
			t.Fatalf("SoftInstall() failed: %v", err)
		}

		testutil.VerifyNoRealCommands(t, mockApp.Base)
	})
}

func TestForceConfigure(t *testing.T) {
	mockApp := testutil.NewMockApp()
	flameshot := &Flameshot{Cmd: mockApp.Cmd}

	err := flameshot.ForceConfigure()
	if err != nil {
		t.Fatalf("ForceConfigure() failed: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestSoftConfigure(t *testing.T) {
	mockApp := testutil.NewMockApp()
	flameshot := &Flameshot{Cmd: mockApp.Cmd}

	err := flameshot.SoftConfigure()
	if err != nil {
		t.Fatalf("SoftConfigure() failed: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestUninstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	flameshot := &Flameshot{Cmd: mockApp.Cmd}

	err := flameshot.Uninstall()

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
	flameshot := &Flameshot{Cmd: mockApp.Cmd}

	err := flameshot.ExecuteCommand("--version")
	if err != nil {
		t.Fatalf("ExecuteCommand() failed: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestUpdate(t *testing.T) {
	mockApp := testutil.NewMockApp()
	flameshot := &Flameshot{Cmd: mockApp.Cmd}

	err := flameshot.Update()

	if err == nil {
		t.Fatal("Expected Update() to return error")
	}

	if !strings.Contains(err.Error(), "not supported") {
		t.Errorf("Expected error to mention 'not supported', got: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}
