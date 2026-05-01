package raycast

import (
	"errors"
	"testing"

	"github.com/cjairm/devgita/internal/apps"
	"github.com/cjairm/devgita/internal/testutil"
	"github.com/cjairm/devgita/pkg/constants"
)

func init() {
	testutil.InitLogger()
}

func TestNew(t *testing.T) {
	r := New()
	if r == nil {
		t.Fatal("New() returned nil")
	}
	if r.Cmd == nil {
		t.Fatal("New() returned Raycast with nil Cmd")
	}
}

func TestNameAndKind(t *testing.T) {
	r := &Raycast{}
	if r.Name() != constants.Raycast {
		t.Errorf("expected Name() %q, got %q", constants.Raycast, r.Name())
	}
	if r.Kind() != apps.KindDesktop {
		t.Errorf("expected Kind() KindDesktop, got %v", r.Kind())
	}
}

func TestInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	r := &Raycast{Cmd: mockApp.Cmd}

	err := r.Install()
	if err != nil {
		t.Fatalf("Install() returned error: %v", err)
	}

	if mockApp.Cmd.InstalledDesktopApp != "raycast" {
		t.Errorf("Expected InstalledDesktopApp to be 'raycast', got '%s'", mockApp.Cmd.InstalledDesktopApp)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestSoftInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	r := &Raycast{Cmd: mockApp.Cmd}

	t.Run("installs when not present", func(t *testing.T) {
		mockApp.Reset()
		mockApp.Cmd.DesktopAppInstalled = false

		err := r.SoftInstall()
		if err != nil {
			t.Fatalf("SoftInstall() failed: %v", err)
		}

		if mockApp.Cmd.MaybeInstalledDesktop != "raycast" {
			t.Errorf("Expected MaybeInstalledDesktop to be 'raycast', got '%s'", mockApp.Cmd.MaybeInstalledDesktop)
		}

		testutil.VerifyNoRealCommands(t, mockApp.Base)
	})

	t.Run("skips when already present", func(t *testing.T) {
		mockApp.Reset()
		mockApp.Cmd.DesktopAppInstalled = true

		err := r.SoftInstall()
		if err != nil {
			t.Fatalf("SoftInstall() failed: %v", err)
		}

		testutil.VerifyNoRealCommands(t, mockApp.Base)
	})
}

func TestForceInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	r := &Raycast{Cmd: mockApp.Cmd}

	if err := r.ForceInstall(); err != nil {
		t.Fatalf("ForceInstall() should succeed even when uninstall is not supported: %v", err)
	}
	if mockApp.Cmd.InstalledDesktopApp != constants.Raycast {
		t.Errorf("expected Install to be called, got %q", mockApp.Cmd.InstalledDesktopApp)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestForceConfigure(t *testing.T) {
	mockApp := testutil.NewMockApp()
	r := &Raycast{Cmd: mockApp.Cmd}

	err := r.ForceConfigure()
	if err != nil {
		t.Fatalf("ForceConfigure() returned error: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestSoftConfigure(t *testing.T) {
	mockApp := testutil.NewMockApp()
	r := &Raycast{Cmd: mockApp.Cmd}

	err := r.SoftConfigure()
	if err != nil {
		t.Fatalf("SoftConfigure() returned error: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestUninstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	r := &Raycast{Cmd: mockApp.Cmd}

	err := r.Uninstall()
	if err == nil {
		t.Fatal("Uninstall() should return error")
	}
	if !errors.Is(err, apps.ErrUninstallNotSupported) {
		t.Errorf("expected ErrUninstallNotSupported, got: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestExecuteCommand(t *testing.T) {
	mockApp := testutil.NewMockApp()
	r := &Raycast{Cmd: mockApp.Cmd}

	err := r.ExecuteCommand("--version")
	if err != nil {
		t.Fatalf("ExecuteCommand() returned error: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestUpdate(t *testing.T) {
	mockApp := testutil.NewMockApp()
	r := &Raycast{Cmd: mockApp.Cmd}

	err := r.Update()
	if err == nil {
		t.Fatal("Update() should return error")
	}
	if !errors.Is(err, apps.ErrUpdateNotSupported) {
		t.Errorf("expected ErrUpdateNotSupported, got: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}
