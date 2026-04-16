package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/logger"
	"github.com/cjairm/devgita/pkg/paths"
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
	app := &Git{Cmd: mc}

	if err := app.Install(); err != nil {
		t.Fatalf("Install error: %v", err)
	}
	if mc.InstalledPkg != constants.Git {
		t.Fatalf("expected InstallPackage(%s), got %q", constants.Git, mc.InstalledPkg)
	}
}

// SKIP: ForceInstall test as per guidelines
// func TestForceInstall(t *testing.T) {
// 	mc := commands.NewMockCommand()
// 	app := &Git{Cmd: mc}
//
// 	if err := app.ForceInstall(); err != nil {
// 		t.Fatalf("ForceInstall error: %v", err)
// 	}
// 	// ForceInstall should call Install() which uses InstallPackage
// 	if mc.InstalledPkg != constants.Git {
// 		t.Fatalf("expected InstallPackage(%s), got %q", constants.Git, mc.InstalledPkg)
// 	}
// }

func TestSoftInstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Git{Cmd: mc}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}
	if mc.MaybeInstalled != constants.Git {
		t.Fatalf("expected MaybeInstallPackage(%s), got %q", constants.Git, mc.MaybeInstalled)
	}
}

func TestUninstall(t *testing.T) {
	mc := commands.NewMockCommand()
	app := &Git{Cmd: mc}

	err := app.Uninstall()
	if err == nil {
		t.Fatal("expected Uninstall to return error for unsupported operation")
	}
	if err.Error() != "git uninstall not supported through devgita" {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestForceConfigure(t *testing.T) {
	// Create temp "app config" dir with a fake file as source
	src := t.TempDir()
	dst := t.TempDir()

	// Override global paths for the duration of the test
	oldAppDir, oldLocalDir := paths.Paths.App.Configs.Git, paths.Paths.Config.Git
	paths.Paths.App.Configs.Git, paths.Paths.Config.Git = src, dst
	t.Cleanup(func() {
		paths.Paths.App.Configs.Git, paths.Paths.Config.Git = oldAppDir, oldLocalDir
	})

	originalContent := "[user]\n\tname = Test User"
	if err := os.WriteFile(filepath.Join(src, ".gitconfig"), []byte(originalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	mc := commands.NewMockCommand()
	app := &Git{Cmd: mc}

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
		t.Fatalf(
			"ForceConfigure did not overwrite: expected %q, got %q",
			originalContent,
			string(finalContent),
		)
	}
}

func TestSoftConfigure(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	// Override global paths for the duration of the test
	oldAppDir, oldLocalDir := paths.Paths.App.Configs.Git, paths.Paths.Config.Git
	paths.Paths.App.Configs.Git, paths.Paths.Config.Git = src, dst
	t.Cleanup(func() {
		paths.Paths.App.Configs.Git, paths.Paths.Config.Git = oldAppDir, oldLocalDir
	})

	originalContent := "[user]\n\tname = Test User"
	if err := os.WriteFile(filepath.Join(src, ".gitconfig"), []byte(originalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	mc := commands.NewMockCommand()
	app := &Git{Cmd: mc}

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
		t.Fatalf(
			"SoftConfigure overwrote existing file: expected %q, got %q",
			modifiedContent,
			string(finalContent),
		)
	}
}

func TestExecuteCommand(t *testing.T) {
	mc := commands.NewMockCommand()
	mockBase := commands.NewMockBaseCommand()
	app := &Git{Cmd: mc, Base: mockBase}

	// Test 1: Successful execution
	t.Run("successful execution", func(t *testing.T) {
		mockBase.SetExecCommandResult("git version 2.39.0", "", nil)

		err := app.ExecuteCommand("--version")
		if err != nil {
			t.Fatalf("ExecuteCommand failed: %v", err)
		}

		// Verify ExecCommand was called once
		if mockBase.GetExecCommandCallCount() != 1 {
			t.Fatalf("Expected 1 ExecCommand call, got %d", mockBase.GetExecCommandCallCount())
		}

		// Verify command parameters
		lastCall := mockBase.GetLastExecCommandCall()
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

	// Test 2: Error handling
	t.Run("command execution error", func(t *testing.T) {
		mockBase.ResetExecCommand()
		mockBase.SetExecCommandResult("", "command not found", fmt.Errorf("command not found: git"))

		err := app.ExecuteCommand("--invalid-flag")
		if err == nil {
			t.Fatal("Expected ExecuteCommand to return error")
		}
		if !strings.Contains(err.Error(), "failed to run git command") {
			t.Fatalf("Expected error to contain 'failed to run git command', got: %v", err)
		}
	})

	// Test 3: Clone command
	t.Run("clone command", func(t *testing.T) {
		mockBase.ResetExecCommand()
		mockBase.SetExecCommandResult("Cloning into...", "", nil)

		err := app.Clone("https://github.com/user/repo.git", "/tmp/repo")
		if err != nil {
			t.Fatalf("Clone failed: %v", err)
		}

		lastCall := mockBase.GetLastExecCommandCall()
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

// SKIP: Uninstall test

// SKIP: Updates test

func TestRemoteBranchExists(t *testing.T) {
	mc := commands.NewMockCommand()
	mockBase := commands.NewMockBaseCommand()
	app := &Git{Cmd: mc, Base: mockBase}

	t.Run("remote branch exists", func(t *testing.T) {
		mockBase.ResetExecCommand()
		mockBase.SetExecCommandResult("  origin/feature-A\n", "", nil)

		exists, err := app.RemoteBranchExists("feature-A")
		if err != nil {
			t.Fatalf("RemoteBranchExists failed: %v", err)
		}
		if !exists {
			t.Error("Expected remote branch to exist")
		}

		lastCall := mockBase.GetLastExecCommandCall()
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
		mockBase.ResetExecCommand()
		mockBase.SetExecCommandResult("", "", nil)

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
	mc := commands.NewMockCommand()
	mockBase := commands.NewMockBaseCommand()
	app := &Git{Cmd: mc, Base: mockBase}

	t.Run("new branch creation - neither local nor remote exists", func(t *testing.T) {
		mockBase.ResetExecCommand()
		// Mock will return empty strings (branch doesn't exist) for all checks
		// This simulates: fetch succeeds, local check fails, remote check fails, worktree add succeeds
		mockBase.SetExecCommandResult("", "", nil)

		err := app.CreateWorktree("/path/to/worktree", "feature-branch")
		if err != nil {
			t.Fatalf("CreateWorktree failed: %v", err)
		}

		// Should have made multiple calls: fetch, local check, remote check, worktree add
		callCount := mockBase.GetExecCommandCallCount()
		if callCount < 3 {
			t.Fatalf("Expected at least 3 command calls, got %d", callCount)
		}

		// Check the final worktree add command used -b for new branch
		lastCall := mockBase.GetLastExecCommandCall()
		if lastCall.Command != "git" {
			t.Fatalf("Expected command 'git', got %q", lastCall.Command)
		}

		// Verify it's a worktree add command with -b flag
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
		mockBase.ResetExecCommand()
		// Simulate worktree add failure
		mockBase.SetExecCommandResult("", "fatal: worktree exists", fmt.Errorf("worktree exists"))

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
	mc := commands.NewMockCommand()
	mockBase := commands.NewMockBaseCommand()
	app := &Git{Cmd: mc, Base: mockBase}

	t.Run("successful list", func(t *testing.T) {
		mockBase.ResetExecCommand()
		porcelainOutput := `worktree /Users/test/repo
HEAD abc123def456
branch refs/heads/main

worktree /Users/test/repo/.worktrees/feature
HEAD def456abc789
branch refs/heads/feature
`
		mockBase.SetExecCommandResult(porcelainOutput, "", nil)

		worktrees, err := app.ListWorktrees()
		if err != nil {
			t.Fatalf("ListWorktrees failed: %v", err)
		}

		if len(worktrees) != 2 {
			t.Fatalf("Expected 2 worktrees, got %d", len(worktrees))
		}

		// Check main worktree
		if worktrees[0].Path != "/Users/test/repo" {
			t.Errorf("Expected path '/Users/test/repo', got %q", worktrees[0].Path)
		}
		if worktrees[0].Branch != "main" {
			t.Errorf("Expected branch 'main', got %q", worktrees[0].Branch)
		}
		if worktrees[0].Commit != "abc123def456" {
			t.Errorf("Expected commit 'abc123def456', got %q", worktrees[0].Commit)
		}

		// Check feature worktree
		if worktrees[1].Path != "/Users/test/repo/.worktrees/feature" {
			t.Errorf("Expected path '/Users/test/repo/.worktrees/feature', got %q", worktrees[1].Path)
		}
		if worktrees[1].Branch != "feature" {
			t.Errorf("Expected branch 'feature', got %q", worktrees[1].Branch)
		}
	})

	t.Run("list error", func(t *testing.T) {
		mockBase.ResetExecCommand()
		mockBase.SetExecCommandResult("", "fatal: not a git repository", fmt.Errorf("not a repo"))

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
	mc := commands.NewMockCommand()
	mockBase := commands.NewMockBaseCommand()
	app := &Git{Cmd: mc, Base: mockBase}

	t.Run("successful removal without branch deletion", func(t *testing.T) {
		mockBase.ResetExecCommand()
		mockBase.SetExecCommandResult("", "", nil)

		err := app.RemoveWorktree("/path/to/worktree", false, "")
		if err != nil {
			t.Fatalf("RemoveWorktree failed: %v", err)
		}

		// Should only call worktree remove
		if mockBase.GetExecCommandCallCount() != 1 {
			t.Fatalf("Expected 1 call, got %d", mockBase.GetExecCommandCallCount())
		}

		lastCall := mockBase.GetLastExecCommandCall()
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
		mockBase.ResetExecCommand()
		mockBase.SetExecCommandResult("", "", nil)

		err := app.RemoveWorktree("/path/to/worktree", true, "feature-branch")
		if err != nil {
			t.Fatalf("RemoveWorktree failed: %v", err)
		}

		// Should call worktree remove + branch delete
		if mockBase.GetExecCommandCallCount() != 2 {
			t.Fatalf("Expected 2 calls, got %d", mockBase.GetExecCommandCallCount())
		}
	})

	t.Run("removal error", func(t *testing.T) {
		mockBase.ResetExecCommand()
		mockBase.SetExecCommandResult("", "fatal: worktree not found", fmt.Errorf("not found"))

		err := app.RemoveWorktree("/nonexistent/path", false, "")
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

func TestGetRepoRoot(t *testing.T) {
	mc := commands.NewMockCommand()
	mockBase := commands.NewMockBaseCommand()
	app := &Git{Cmd: mc, Base: mockBase}

	t.Run("successful get repo root", func(t *testing.T) {
		mockBase.ResetExecCommand()
		mockBase.SetExecCommandResult("/Users/test/my-repo\n", "", nil)

		root, err := app.GetRepoRoot()
		if err != nil {
			t.Fatalf("GetRepoRoot failed: %v", err)
		}

		if root != "/Users/test/my-repo" {
			t.Errorf("Expected '/Users/test/my-repo', got %q", root)
		}

		lastCall := mockBase.GetLastExecCommandCall()
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
		mockBase.ResetExecCommand()
		mockBase.SetExecCommandResult("", "fatal: not a git repository", fmt.Errorf("not a repo"))

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
