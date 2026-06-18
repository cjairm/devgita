package cmd

import (
	"fmt"
	"testing"
)

// mockTaskRunner records calls to each task method.
type mockTaskRunner struct {
	refreshBranchArg    string
	refreshBranchCalled bool
	refreshBranchErr    error

	resetMainCalled bool
	resetMainErr    error

	reinstallLibsCalled bool
	reinstallLibsErr    error

	reinstallLibArg    string
	reinstallLibCalled bool
	reinstallLibErr    error

	deleteBranchArg    string
	deleteBranchCalled bool
	deleteBranchErr    error
}

func (m *mockTaskRunner) RefreshBranch(target string) error {
	m.refreshBranchCalled = true
	m.refreshBranchArg = target
	return m.refreshBranchErr
}

func (m *mockTaskRunner) ResetMainBranch() error {
	m.resetMainCalled = true
	return m.resetMainErr
}

func (m *mockTaskRunner) ReinstallLibraries() error {
	m.reinstallLibsCalled = true
	return m.reinstallLibsErr
}

func (m *mockTaskRunner) ReinstallLibrary(name string) error {
	m.reinstallLibCalled = true
	m.reinstallLibArg = name
	return m.reinstallLibErr
}

func (m *mockTaskRunner) DeleteBranch(target string) error {
	m.deleteBranchCalled = true
	m.deleteBranchArg = target
	return m.deleteBranchErr
}

func setupTaskMock(t *testing.T, mock taskRunner) func() {
	t.Helper()
	orig := newTaskManager
	newTaskManager = func() taskRunner { return mock }
	return func() { newTaskManager = orig }
}

func TestTask_RefreshBranch(t *testing.T) {
	t.Run("no args defaults to empty target", func(t *testing.T) {
		mock := &mockTaskRunner{}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskRefreshBranchCmd.RunE(taskRefreshBranchCmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !mock.refreshBranchCalled {
			t.Error("expected RefreshBranch to be called")
		}
		if mock.refreshBranchArg != "" {
			t.Errorf("expected empty target, got %q", mock.refreshBranchArg)
		}
	})

	t.Run("passes target arg", func(t *testing.T) {
		mock := &mockTaskRunner{}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskRefreshBranchCmd.RunE(taskRefreshBranchCmd, []string{"develop"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.refreshBranchArg != "develop" {
			t.Errorf("expected target 'develop', got %q", mock.refreshBranchArg)
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := &mockTaskRunner{refreshBranchErr: fmt.Errorf("git failed")}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskRefreshBranchCmd.RunE(taskRefreshBranchCmd, []string{})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestTask_ResetMainBranch(t *testing.T) {
	t.Run("calls ResetMainBranch", func(t *testing.T) {
		mock := &mockTaskRunner{}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskResetMainBranchCmd.RunE(taskResetMainBranchCmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !mock.resetMainCalled {
			t.Error("expected ResetMainBranch to be called")
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := &mockTaskRunner{resetMainErr: fmt.Errorf("reset failed")}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskResetMainBranchCmd.RunE(taskResetMainBranchCmd, []string{})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestTask_ReinstallLibraries(t *testing.T) {
	t.Run("calls ReinstallLibraries", func(t *testing.T) {
		mock := &mockTaskRunner{}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskReinstallLibrariesCmd.RunE(taskReinstallLibrariesCmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !mock.reinstallLibsCalled {
			t.Error("expected ReinstallLibraries to be called")
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := &mockTaskRunner{reinstallLibsErr: fmt.Errorf("clean failed")}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskReinstallLibrariesCmd.RunE(taskReinstallLibrariesCmd, []string{})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestTask_ReinstallLibrary(t *testing.T) {
	t.Run("passes library name", func(t *testing.T) {
		mock := &mockTaskRunner{}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskReinstallLibraryCmd.RunE(taskReinstallLibraryCmd, []string{"lodash"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !mock.reinstallLibCalled {
			t.Error("expected ReinstallLibrary to be called")
		}
		if mock.reinstallLibArg != "lodash" {
			t.Errorf("expected 'lodash', got %q", mock.reinstallLibArg)
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := &mockTaskRunner{reinstallLibErr: fmt.Errorf("rm failed")}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskReinstallLibraryCmd.RunE(taskReinstallLibraryCmd, []string{"lodash"})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestTask_DeleteBranch(t *testing.T) {
	t.Run("no args passes empty target", func(t *testing.T) {
		mock := &mockTaskRunner{}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskDeleteBranchCmd.RunE(taskDeleteBranchCmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !mock.deleteBranchCalled {
			t.Error("expected DeleteBranch to be called")
		}
		if mock.deleteBranchArg != "" {
			t.Errorf("expected empty target, got %q", mock.deleteBranchArg)
		}
	})

	t.Run("passes target arg", func(t *testing.T) {
		mock := &mockTaskRunner{}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskDeleteBranchCmd.RunE(taskDeleteBranchCmd, []string{"develop"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.deleteBranchArg != "develop" {
			t.Errorf("expected 'develop', got %q", mock.deleteBranchArg)
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		mock := &mockTaskRunner{deleteBranchErr: fmt.Errorf("delete failed")}
		restore := setupTaskMock(t, mock)
		defer restore()

		err := taskDeleteBranchCmd.RunE(taskDeleteBranchCmd, []string{})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
