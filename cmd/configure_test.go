package cmd

import (
	"fmt"
	"testing"

	"github.com/cjairm/devgita/internal/apps"
	"github.com/cjairm/devgita/internal/testutil"
)

func init() { testutil.InitLogger() }

// mockConfigureApp is a minimal apps.App that records which configure method was called.
type mockConfigureApp struct {
	forceCalled bool
	softCalled  bool
	forceErr    error
	softErr     error
}

func (m *mockConfigureApp) Name() string                   { return "mock" }
func (m *mockConfigureApp) Kind() apps.AppKind             { return apps.KindTerminal }
func (m *mockConfigureApp) Install() error                 { return nil }
func (m *mockConfigureApp) ForceInstall() error            { return nil }
func (m *mockConfigureApp) SoftInstall() error             { return nil }
func (m *mockConfigureApp) Uninstall() error               { return nil }
func (m *mockConfigureApp) Update() error                  { return nil }
func (m *mockConfigureApp) ExecuteCommand(...string) error { return nil }

func (m *mockConfigureApp) ForceConfigure() error {
	m.forceCalled = true
	return m.forceErr
}

func (m *mockConfigureApp) SoftConfigure() error {
	m.softCalled = true
	return m.softErr
}

func setupConfigureCmd(t *testing.T, mock apps.App) func() {
	t.Helper()
	origApp := getAppFn
	origRefresh := refreshEmbeddedConfigs
	getAppFn = func(name string) (apps.App, error) { return mock, nil }
	refreshEmbeddedConfigs = func() error { return nil }
	return func() {
		getAppFn = origApp
		refreshEmbeddedConfigs = origRefresh
	}
}

func TestConfigure_SoftPath(t *testing.T) {
	mock := &mockConfigureApp{}
	restore := setupConfigureCmd(t, mock)
	defer restore()

	configureForce = false
	err := runConfigure(configureCmd, []string{"git"})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !mock.softCalled {
		t.Error("expected SoftConfigure to be called")
	}
	if mock.forceCalled {
		t.Error("expected ForceConfigure NOT to be called")
	}
}

func TestConfigure_ForcePath(t *testing.T) {
	mock := &mockConfigureApp{}
	restore := setupConfigureCmd(t, mock)
	defer restore()

	configureForce = true
	defer func() { configureForce = false }()

	err := runConfigure(configureCmd, []string{"git"})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !mock.forceCalled {
		t.Error("expected ForceConfigure to be called")
	}
	if mock.softCalled {
		t.Error("expected SoftConfigure NOT to be called")
	}
}

func TestConfigure_UnknownApp(t *testing.T) {
	origApp := getAppFn
	origRefresh := refreshEmbeddedConfigs
	getAppFn = func(name string) (apps.App, error) {
		return nil, fmt.Errorf("unknown app %q", name)
	}
	refreshEmbeddedConfigs = func() error { return nil }
	defer func() {
		getAppFn = origApp
		refreshEmbeddedConfigs = origRefresh
	}()

	configureForce = false
	err := runConfigure(configureCmd, []string{"notanapp"})
	if err == nil {
		t.Fatal("expected non-nil error for unknown app")
	}
}

func TestConfigure_NotSupported(t *testing.T) {
	mock := &mockConfigureApp{softErr: fmt.Errorf("%w for mock", apps.ErrConfigureNotSupported)}
	restore := setupConfigureCmd(t, mock)
	defer restore()

	configureForce = false
	err := runConfigure(configureCmd, []string{"brave"})
	if err != nil {
		t.Fatalf("ErrConfigureNotSupported should exit zero, got: %v", err)
	}
}
