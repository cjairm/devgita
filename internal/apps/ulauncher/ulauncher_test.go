package ulauncher

import (
	"strings"
	"testing"

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

	err := u.ForceInstall()
	// Should fail because Uninstall is not supported
	if err == nil {
		t.Fatal("Expected ForceInstall to fail due to unsupported Uninstall")
	}

	if !strings.Contains(err.Error(), "uninstall not supported") {
		t.Errorf("Expected error about uninstall not supported, got: %v", err)
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

	if !strings.Contains(err.Error(), "not supported") {
		t.Errorf("Expected error about not supported, got: %v", err)
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

	if !strings.Contains(err.Error(), "not implemented") {
		t.Errorf("Expected error about not implemented, got: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}
