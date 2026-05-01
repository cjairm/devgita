package ulauncher

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
	u := New()
	if u == nil {
		t.Fatal("New() returned nil")
	}
	if u.Cmd == nil {
		t.Fatal("Cmd is nil")
	}
}

func TestNameAndKind(t *testing.T) {
	u := &Ulauncher{}
	if u.Name() != constants.Ulauncher {
		t.Errorf("expected Name() %q, got %q", constants.Ulauncher, u.Name())
	}
	if u.Kind() != apps.KindDesktop {
		t.Errorf("expected Kind() KindDesktop, got %v", u.Kind())
	}
}

func TestInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	u := &Ulauncher{Cmd: mockApp.Cmd}

	err := u.Install()
	if err != nil {
		t.Fatalf("Install() failed: %v", err)
	}

	if mockApp.Cmd.InstalledDesktopApp != constants.Ulauncher {
		t.Errorf("Expected desktop app %s, got %s", constants.Ulauncher, mockApp.Cmd.InstalledDesktopApp)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestSoftInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	u := &Ulauncher{Cmd: mockApp.Cmd}

	err := u.SoftInstall()
	if err != nil {
		t.Fatalf("SoftInstall() failed: %v", err)
	}

	if mockApp.Cmd.MaybeInstalledDesktop != constants.Ulauncher {
		t.Errorf("Expected desktop app %s, got %s", constants.Ulauncher, mockApp.Cmd.MaybeInstalledDesktop)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestForceInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	u := &Ulauncher{Cmd: mockApp.Cmd}

	if err := u.ForceInstall(); err != nil {
		t.Fatalf("ForceInstall() should succeed even when uninstall is not supported: %v", err)
	}
	if mockApp.Cmd.InstalledDesktopApp != constants.Ulauncher {
		t.Errorf("expected Install to be called, got %q", mockApp.Cmd.InstalledDesktopApp)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestUninstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	u := &Ulauncher{Cmd: mockApp.Cmd}

	err := u.Uninstall()
	if err == nil {
		t.Fatal("Expected Uninstall to return error")
	}
	if !errors.Is(err, apps.ErrUninstallNotSupported) {
		t.Errorf("expected ErrUninstallNotSupported, got: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestForceConfigure(t *testing.T) {
	mockApp := testutil.NewMockApp()
	u := &Ulauncher{Cmd: mockApp.Cmd}

	err := u.ForceConfigure()
	if err != nil {
		t.Fatalf("ForceConfigure() failed: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestSoftConfigure(t *testing.T) {
	mockApp := testutil.NewMockApp()
	u := &Ulauncher{Cmd: mockApp.Cmd}

	err := u.SoftConfigure()
	if err != nil {
		t.Fatalf("SoftConfigure() failed: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestExecuteCommand(t *testing.T) {
	mockApp := testutil.NewMockApp()
	u := &Ulauncher{Cmd: mockApp.Cmd}

	err := u.ExecuteCommand("--version")
	if err != nil {
		t.Fatalf("ExecuteCommand() failed: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestUpdate(t *testing.T) {
	mockApp := testutil.NewMockApp()
	u := &Ulauncher{Cmd: mockApp.Cmd}

	err := u.Update()
	if err == nil {
		t.Fatal("Expected Update to return error")
	}
	if !errors.Is(err, apps.ErrUpdateNotSupported) {
		t.Errorf("expected ErrUpdateNotSupported, got: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}
