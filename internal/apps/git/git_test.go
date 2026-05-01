package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/apps"
	"github.com/cjairm/devgita/internal/testutil"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/paths"
)

func init() {
	testutil.InitLogger()
}

func TestNew(t *testing.T) {
	app := New()
	if app == nil {
		t.Fatal("New() returned nil")
	}
}

func TestNameAndKind(t *testing.T) {
	g := &Git{}
	if g.Name() != constants.Git {
		t.Errorf("expected Name() %q, got %q", constants.Git, g.Name())
	}
	if g.Kind() != apps.KindTerminal {
		t.Errorf("expected Kind() KindTerminal, got %v", g.Kind())
	}
}

func TestInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Git{Cmd: mockApp.Cmd}

	if err := app.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}
	if mockApp.Cmd.InstalledPkg != constants.Git {
		t.Fatalf("expected InstallPackage(%s), got %q", constants.Git, mockApp.Cmd.InstalledPkg)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestForceInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}

	if err := app.ForceInstall(); err != nil {
		t.Fatalf("ForceInstall() should succeed even when uninstall is not supported: %v", err)
	}
	if mockApp.Cmd.InstalledPkg != constants.Git {
		t.Errorf("expected Install to be called, got %q", mockApp.Cmd.InstalledPkg)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestSoftInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Git{Cmd: mockApp.Cmd}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}
	if mockApp.Cmd.MaybeInstalled != constants.Git {
		t.Fatalf("expected MaybeInstallPackage(%s), got %q", constants.Git, mockApp.Cmd.MaybeInstalled)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestUninstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}

	err := app.Uninstall()
	if err == nil {
		t.Fatal("expected Uninstall to return error for unsupported operation")
	}
	if !errors.Is(err, apps.ErrUninstallNotSupported) {
		t.Errorf("expected ErrUninstallNotSupported, got: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestUpdate(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}

	err := app.Update()
	if err == nil {
		t.Fatal("expected Update to return error")
	}
	if !errors.Is(err, apps.ErrUpdateNotSupported) {
		t.Errorf("expected ErrUpdateNotSupported, got: %v", err)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestForceConfigure(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	oldAppDir, oldLocalDir := paths.Paths.App.Configs.Git, paths.Paths.Config.Git
	paths.Paths.App.Configs.Git, paths.Paths.Config.Git = src, dst
	t.Cleanup(func() {
		paths.Paths.App.Configs.Git, paths.Paths.Config.Git = oldAppDir, oldLocalDir
	})

	originalContent := "[user]\n\tname = Test User"
	if err := os.WriteFile(filepath.Join(src, ".gitconfig"), []byte(originalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	mockApp := testutil.NewMockApp()
	app := &Git{Cmd: mockApp.Cmd}

	if err := app.ForceConfigure(); err != nil {
		t.Fatalf("ForceConfigure error: %v", err)
	}

	check := filepath.Join(dst, ".gitconfig")
	if _, err := os.Stat(check); err != nil {
		t.Fatalf("expected copied file at %s: %v", check, err)
	}

	copiedContent, err := os.ReadFile(check)
	if err != nil {
		t.Fatalf("failed to read copied file: %v", err)
	}
	if string(copiedContent) != originalContent {
		t.Fatalf("content mismatch: expected %q, got %q", originalContent, string(copiedContent))
	}

	modifiedContent := "[user]\n\tname = Modified User"
	if err := os.WriteFile(check, []byte(modifiedContent), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := app.ForceConfigure(); err != nil {
		t.Fatalf("second ForceConfigure error: %v", err)
	}

	finalContent, err := os.ReadFile(check)
	if err != nil {
		t.Fatalf("failed to read file after second configure: %v", err)
	}
	if string(finalContent) == string(modifiedContent) {
		t.Fatalf("ForceConfigure did not overwrite: expected %q, got %q", originalContent, string(finalContent))
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestSoftConfigure(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	oldAppDir, oldLocalDir := paths.Paths.App.Configs.Git, paths.Paths.Config.Git
	paths.Paths.App.Configs.Git, paths.Paths.Config.Git = src, dst
	t.Cleanup(func() {
		paths.Paths.App.Configs.Git, paths.Paths.Config.Git = oldAppDir, oldLocalDir
	})

	originalContent := "[user]\n\tname = Test User"
	if err := os.WriteFile(filepath.Join(src, ".gitconfig"), []byte(originalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	mockApp := testutil.NewMockApp()
	app := &Git{Cmd: mockApp.Cmd}

	if err := app.SoftConfigure(); err != nil {
		t.Fatalf("SoftConfigure error: %v", err)
	}

	check := filepath.Join(dst, ".gitconfig")
	if _, err := os.Stat(check); err != nil {
		t.Fatalf("expected copied file at %s: %v", check, err)
	}

	modifiedContent := "[user]\n\tname = Modified User"
	if err := os.WriteFile(check, []byte(modifiedContent), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := app.SoftConfigure(); err != nil {
		t.Fatalf("second SoftConfigure error: %v", err)
	}

	finalContent, err := os.ReadFile(check)
	if err != nil {
		t.Fatalf("failed to read file after second configure: %v", err)
	}
	if string(finalContent) == string(originalContent) {
		t.Fatalf("SoftConfigure overwrote existing file: expected %q, got %q", modifiedContent, string(finalContent))
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestExecuteCommand(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}

	t.Run("successful execution", func(t *testing.T) {
		mockApp.Base.SetExecCommandResult("git version 2.39.0", "", nil)

		if err := app.ExecuteCommand("--version"); err != nil {
			t.Fatalf("ExecuteCommand failed: %v", err)
		}

		if mockApp.Base.GetExecCommandCallCount() != 1 {
			t.Fatalf("Expected 1 ExecCommand call, got %d", mockApp.Base.GetExecCommandCallCount())
		}

		lastCall := mockApp.Base.GetLastExecCommandCall()
		if lastCall == nil {
			t.Fatal("No ExecCommand call recorded")
		}
		if lastCall.Command != "git" {
			t.Fatalf("Expected command 'git', got %q", lastCall.Command)
		}
		if len(lastCall.Args) != 1 || lastCall.Args[0] != "--version" {
			t.Fatalf("Expected args ['--version'], got %v", lastCall.Args)
		}
		if lastCall.IsSudo {
			t.Fatal("Expected IsSudo to be false")
		}
	})

	t.Run("command execution error", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult("", "command not found", fmt.Errorf("command not found: git"))

		err := app.ExecuteCommand("--invalid-flag")
		if err == nil {
			t.Fatal("Expected ExecuteCommand to return error")
		}
		if !strings.Contains(err.Error(), "failed to run git command") {
			t.Fatalf("Expected error to contain 'failed to run git command', got: %v", err)
		}
	})

	t.Run("clone command", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult("Cloning into...", "", nil)

		if err := app.Clone("https://github.com/user/repo.git", "/tmp/repo"); err != nil {
			t.Fatalf("Clone failed: %v", err)
		}

		lastCall := mockApp.Base.GetLastExecCommandCall()
		expectedArgs := []string{"clone", "https://github.com/user/repo.git", "/tmp/repo"}
		if len(lastCall.Args) != len(expectedArgs) {
			t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(lastCall.Args))
		}
		for i, arg := range expectedArgs {
			if lastCall.Args[i] != arg {
				t.Fatalf("Expected arg[%d] to be %q, got %q", i, arg, lastCall.Args[i])
			}
		}
	})
}

func TestRemoteBranchExists(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}

	t.Run("remote branch exists", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult("  origin/feature-A\n", "", nil)

		exists, err := app.RemoteBranchExists("feature-A")
		if err != nil {
			t.Fatalf("RemoteBranchExists failed: %v", err)
		}
		if !exists {
			t.Error("Expected remote branch to exist")
		}

		lastCall := mockApp.Base.GetLastExecCommandCall()
		if lastCall.Command != "git" {
			t.Fatalf("Expected command 'git', got %q", lastCall.Command)
		}
		expectedArgs := []string{"branch", "-r", "--list", "origin/feature-A"}
		if len(lastCall.Args) != len(expectedArgs) {
			t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(lastCall.Args))
		}
		for i, arg := range expectedArgs {
			if lastCall.Args[i] != arg {
				t.Fatalf("Expected arg[%d] to be %q, got %q", i, arg, lastCall.Args[i])
			}
		}
	})

	t.Run("remote branch does not exist", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult("", "", nil)

		exists, err := app.RemoteBranchExists("feature-B")
		if err != nil {
			t.Fatalf("RemoteBranchExists failed: %v", err)
		}
		if exists {
			t.Error("Expected remote branch to not exist")
		}
	})
}

func TestCreateWorktree(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}

	t.Run("new branch creation - neither local nor remote exists", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult("", "", nil)

		if err := app.CreateWorktree("/path/to/worktree", "feature-branch"); err != nil {
			t.Fatalf("CreateWorktree failed: %v", err)
		}

		callCount := mockApp.Base.GetExecCommandCallCount()
		if callCount < 3 {
			t.Fatalf("Expected at least 3 command calls, got %d", callCount)
		}

		lastCall := mockApp.Base.GetLastExecCommandCall()
		if lastCall.Command != "git" {
			t.Fatalf("Expected command 'git', got %q", lastCall.Command)
		}

		hasWorktreeAdd := false
		hasNewBranchFlag := false
		for i, arg := range lastCall.Args {
			if arg == "worktree" && i+1 < len(lastCall.Args) && lastCall.Args[i+1] == "add" {
				hasWorktreeAdd = true
			}
			if arg == "-b" {
				hasNewBranchFlag = true
			}
		}
		if !hasWorktreeAdd {
			t.Fatal("Expected 'git worktree add' command")
		}
		if !hasNewBranchFlag {
			t.Fatal("Expected -b flag for creating new branch")
		}
	})

	t.Run("creation error on worktree add", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult("", "fatal: worktree exists", fmt.Errorf("worktree exists"))

		err := app.CreateWorktree("/path/to/worktree", "existing-branch")
		if err == nil {
			t.Fatal("Expected error but got none")
		}
		if !strings.Contains(err.Error(), "worktree exists") {
			t.Errorf("Expected error to contain 'worktree exists', got: %v", err)
		}
	})
}

func TestListWorktrees(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}

	t.Run("successful list", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		porcelainOutput := `worktree /Users/test/repo
HEAD abc123def456
branch refs/heads/main

worktree /Users/test/repo/.worktrees/feature
HEAD def456abc789
branch refs/heads/feature
`
		mockApp.Base.SetExecCommandResult(porcelainOutput, "", nil)

		worktrees, err := app.ListWorktrees()
		if err != nil {
			t.Fatalf("ListWorktrees failed: %v", err)
		}
		if len(worktrees) != 2 {
			t.Fatalf("Expected 2 worktrees, got %d", len(worktrees))
		}

		if worktrees[0].Path != "/Users/test/repo" {
			t.Errorf("Expected path '/Users/test/repo', got %q", worktrees[0].Path)
		}
		if worktrees[0].Branch != "main" {
			t.Errorf("Expected branch 'main', got %q", worktrees[0].Branch)
		}
		if worktrees[0].Commit != "abc123def456" {
			t.Errorf("Expected commit 'abc123def456', got %q", worktrees[0].Commit)
		}

		if worktrees[1].Path != "/Users/test/repo/.worktrees/feature" {
			t.Errorf("Expected path '/Users/test/repo/.worktrees/feature', got %q", worktrees[1].Path)
		}
		if worktrees[1].Branch != "feature" {
			t.Errorf("Expected branch 'feature', got %q", worktrees[1].Branch)
		}
	})

	t.Run("list error", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult("", "fatal: not a git repository", fmt.Errorf("not a repo"))

		_, err := app.ListWorktrees()
		if err == nil {
			t.Fatal("Expected error but got none")
		}
		if !strings.Contains(err.Error(), "failed to list worktrees") {
			t.Fatalf("Expected error message to contain 'failed to list worktrees', got: %v", err)
		}
	})
}

func TestRemoveWorktree(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}

	t.Run("successful removal without branch deletion", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult("", "", nil)

		if err := app.RemoveWorktree("/path/to/worktree", false, ""); err != nil {
			t.Fatalf("RemoveWorktree failed: %v", err)
		}
		if mockApp.Base.GetExecCommandCallCount() != 1 {
			t.Fatalf("Expected 1 call, got %d", mockApp.Base.GetExecCommandCallCount())
		}

		lastCall := mockApp.Base.GetLastExecCommandCall()
		if lastCall == nil {
			t.Fatal("No ExecCommand call recorded")
		}
		expectedArgs := []string{"worktree", "remove", "/path/to/worktree"}
		if len(lastCall.Args) != len(expectedArgs) {
			t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(lastCall.Args))
		}
		for i, arg := range expectedArgs {
			if lastCall.Args[i] != arg {
				t.Fatalf("Expected arg[%d] to be %q, got %q", i, arg, lastCall.Args[i])
			}
		}
	})

	t.Run("successful removal with branch deletion", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult("", "", nil)

		if err := app.RemoveWorktree("/path/to/worktree", true, "feature-branch"); err != nil {
			t.Fatalf("RemoveWorktree failed: %v", err)
		}
		if mockApp.Base.GetExecCommandCallCount() != 2 {
			t.Fatalf("Expected 2 calls, got %d", mockApp.Base.GetExecCommandCallCount())
		}
	})

	t.Run("removal error", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult("", "fatal: worktree not found", fmt.Errorf("not found"))

		if err := app.RemoveWorktree("/nonexistent/path", false, ""); err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

func TestGetRepoRoot(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}

	t.Run("successful get repo root", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult("/Users/test/my-repo\n", "", nil)

		root, err := app.GetRepoRoot()
		if err != nil {
			t.Fatalf("GetRepoRoot failed: %v", err)
		}
		if root != "/Users/test/my-repo" {
			t.Errorf("Expected '/Users/test/my-repo', got %q", root)
		}

		lastCall := mockApp.Base.GetLastExecCommandCall()
		if lastCall == nil {
			t.Fatal("No ExecCommand call recorded")
		}
		expectedArgs := []string{"rev-parse", "--show-toplevel"}
		for i, arg := range expectedArgs {
			if lastCall.Args[i] != arg {
				t.Fatalf("Expected arg[%d] to be %q, got %q", i, arg, lastCall.Args[i])
			}
		}
	})

	t.Run("not a git repo", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult("", "fatal: not a git repository", fmt.Errorf("not a repo"))

		_, err := app.GetRepoRoot()
		if err == nil {
			t.Fatal("Expected error but got none")
		}
		if !strings.Contains(err.Error(), "failed to get repo root") {
			t.Fatalf("Expected error message to contain 'failed to get repo root', got: %v", err)
		}
	})
}

func TestParseWorktreeOutput(t *testing.T) {
	t.Run("multiple worktrees", func(t *testing.T) {
		output := `worktree /Users/test/repo
HEAD abc123
branch refs/heads/main

worktree /Users/test/repo/.worktrees/feature
HEAD def456
branch refs/heads/feature
`
		worktrees := parseWorktreeOutput(output)
		if len(worktrees) != 2 {
			t.Fatalf("Expected 2 worktrees, got %d", len(worktrees))
		}
		if worktrees[0].Path != "/Users/test/repo" {
			t.Errorf("Expected path '/Users/test/repo', got %q", worktrees[0].Path)
		}
		if worktrees[0].Commit != "abc123" {
			t.Errorf("Expected commit 'abc123', got %q", worktrees[0].Commit)
		}
		if worktrees[0].Branch != "main" {
			t.Errorf("Expected branch 'main', got %q", worktrees[0].Branch)
		}
	})

	t.Run("empty output", func(t *testing.T) {
		worktrees := parseWorktreeOutput("")
		if len(worktrees) != 0 {
			t.Errorf("Expected 0 worktrees for empty output, got %d", len(worktrees))
		}
	})

	t.Run("single worktree without trailing newline", func(t *testing.T) {
		output := `worktree /Users/test/repo
HEAD abc123
branch refs/heads/main`
		worktrees := parseWorktreeOutput(output)
		if len(worktrees) != 1 {
			t.Fatalf("Expected 1 worktree, got %d", len(worktrees))
		}
		if worktrees[0].Path != "/Users/test/repo" {
			t.Errorf("Expected path '/Users/test/repo', got %q", worktrees[0].Path)
		}
	})

	t.Run("detached HEAD", func(t *testing.T) {
		output := `worktree /Users/test/repo
HEAD abc123
detached
`
		worktrees := parseWorktreeOutput(output)
		if len(worktrees) != 1 {
			t.Fatalf("Expected 1 worktree, got %d", len(worktrees))
		}
		if worktrees[0].Branch != "" {
			t.Errorf("Expected empty branch for detached HEAD, got %q", worktrees[0].Branch)
		}
	})
}
