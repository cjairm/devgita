package baseapp

import (
	"errors"
	"fmt"
	"testing"

	"github.com/cjairm/devgita/internal/apps"
	"github.com/cjairm/devgita/internal/testutil"
)

func init() { testutil.InitLogger() }

func TestReinstall_UninstallSucceeds(t *testing.T) {
	uninstalled := false
	installed := false

	err := Reinstall(
		func() error { installed = true; return nil },
		func() error { uninstalled = true; return nil },
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !uninstalled {
		t.Error("expected uninstall to be called")
	}
	if !installed {
		t.Error("expected install to be called")
	}
}

func TestReinstall_UninstallReturnsNotSupported(t *testing.T) {
	installed := false

	err := Reinstall(
		func() error { installed = true; return nil },
		func() error { return fmt.Errorf("wrapped: %w", apps.ErrUninstallNotSupported) },
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !installed {
		t.Error("expected install to still be called when uninstall is not supported")
	}
}

func TestReinstall_UninstallReturnsOtherError(t *testing.T) {
	installErr := errors.New("some uninstall failure")
	installed := false

	err := Reinstall(
		func() error { installed = true; return nil },
		func() error { return installErr },
	)
	if !errors.Is(err, installErr) {
		t.Fatalf("expected uninstall error to propagate, got: %v", err)
	}
	if installed {
		t.Error("install should not be called when uninstall returns a real error")
	}
}

func TestReinstall_InstallReturnsError(t *testing.T) {
	installErr := errors.New("install failed")

	err := Reinstall(
		func() error { return installErr },
		func() error { return apps.ErrUninstallNotSupported },
	)
	if !errors.Is(err, installErr) {
		t.Fatalf("expected install error to propagate, got: %v", err)
	}
}
