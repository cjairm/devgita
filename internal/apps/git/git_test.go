package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/apps"
	"github.com/cjairm/devgita/internal/commands"
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
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()
	testutil.IsolateXDGDirs(t)

	oldConfigGit := paths.Paths.Config.Git
	oldAppConfigsGit := paths.Paths.App.Configs.Git
	t.Cleanup(func() {
		paths.Paths.Config.Git = oldConfigGit
		paths.Paths.App.Configs.Git = oldAppConfigsGit
	})
	paths.Paths.Config.Git = filepath.Join(t.TempDir(), "git-config")
	paths.Paths.App.Configs.Git = filepath.Join(t.TempDir(), "git-app-configs")

	app := &Git{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

	if err := app.ForceInstall(); err != nil {
		t.Fatalf("ForceInstall() error: %v", err)
	}
	if tc.MockApp.Cmd.InstalledPkg != constants.Git {
		t.Errorf("expected Install to be called, got %q", tc.MockApp.Cmd.InstalledPkg)
	}

	testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
}

func TestSoftInstall(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Git{Cmd: mockApp.Cmd}

	if err := app.SoftInstall(); err != nil {
		t.Fatalf("SoftInstall error: %v", err)
	}
	if mockApp.Cmd.MaybeInstalled != constants.Git {
		t.Fatalf(
			"expected MaybeInstallPackage(%s), got %q",
			constants.Git,
			mockApp.Cmd.MaybeInstalled,
		)
	}

	testutil.VerifyNoRealCommands(t, mockApp.Base)
}

func TestUninstall(t *testing.T) {
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()
	testutil.IsolateXDGDirs(t)

	oldConfigGit := paths.Paths.Config.Git
	oldAppConfigsGit := paths.Paths.App.Configs.Git
	t.Cleanup(func() {
		paths.Paths.Config.Git = oldConfigGit
		paths.Paths.App.Configs.Git = oldAppConfigsGit
	})
	paths.Paths.Config.Git = filepath.Join(t.TempDir(), "git-config")
	paths.Paths.App.Configs.Git = filepath.Join(t.TempDir(), "git-app-configs")

	app := &Git{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

	if err := app.Uninstall(); err != nil {
		t.Fatalf("Uninstall error: %v", err)
	}
	if tc.MockApp.Cmd.UninstalledPkg != constants.Git {
		t.Errorf(
			"expected UninstallPackage(%s), got %q",
			constants.Git,
			tc.MockApp.Cmd.UninstalledPkg,
		)
	}

	testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
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
	tc := testutil.SetupCompleteTest(t)
	defer tc.Cleanup()
	testutil.IsolateXDGDirs(t)

	src := filepath.Join(tc.AppDir, "git-src")
	dst := filepath.Join(tc.ConfigDir, "git")
	if err := os.MkdirAll(src, 0o755); err != nil {
		t.Fatal(err)
	}

	oldAppDir, oldLocalDir := paths.Paths.App.Configs.Git, paths.Paths.Config.Git
	t.Cleanup(func() {
		paths.Paths.App.Configs.Git, paths.Paths.Config.Git = oldAppDir, oldLocalDir
	})
	paths.Paths.App.Configs.Git, paths.Paths.Config.Git = src, dst

	originalContent := "[user]\n\tname = Test User"
	if err := os.WriteFile(
		filepath.Join(src, ".gitconfig"),
		[]byte(originalContent),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	app := &Git{Cmd: tc.MockApp.Cmd, Base: tc.MockApp.Base}

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

	testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
}

func TestSoftConfigure(t *testing.T) {
	testutil.IsolateXDGDirs(t)

	src := t.TempDir()
	dst := t.TempDir()

	oldAppDir, oldLocalDir := paths.Paths.App.Configs.Git, paths.Paths.Config.Git
	paths.Paths.App.Configs.Git, paths.Paths.Config.Git = src, dst
	t.Cleanup(func() {
		paths.Paths.App.Configs.Git, paths.Paths.Config.Git = oldAppDir, oldLocalDir
	})

	originalContent := "[user]\n\tname = Test User"
	if err := os.WriteFile(
		filepath.Join(src, ".gitconfig"),
		[]byte(originalContent),
		0o644,
	); err != nil {
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
		t.Fatalf(
			"SoftConfigure overwrote existing file: expected %q, got %q",
			modifiedContent,
			string(finalContent),
		)
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
		mockApp.Base.SetExecCommandResult(
			"",
			"command not found",
			fmt.Errorf("command not found: git"),
		)

		err := app.ExecuteCommand("--invalid-flag")
		if err == nil {
			t.Fatal("Expected ExecuteCommand to return error")
		}
		if !strings.Contains(err.Error(), "git: command not found") {
			t.Fatalf("Expected error to contain 'git: command not found', got: %v", err)
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

func TestBranchExists(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}

	t.Run("branch exists returns true", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult("  feature-a\n", "", nil)

		exists, err := app.BranchExists("feature-a")
		if err != nil {
			t.Fatalf("BranchExists failed: %v", err)
		}
		if !exists {
			t.Error("expected branch to exist")
		}
		last := mockApp.Base.GetLastExecCommandCall()
		if last == nil {
			t.Fatal("no ExecCommand call recorded")
		}
		expectedArgs := []string{"branch", "--list", "feature-a"}
		if len(last.Args) != len(expectedArgs) {
			t.Fatalf("expected args %v, got %v", expectedArgs, last.Args)
		}
	})

	t.Run("branch does not exist returns false", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult("", "", nil)

		exists, err := app.BranchExists("no-such-branch")
		if err != nil {
			t.Fatalf("BranchExists failed: %v", err)
		}
		if exists {
			t.Error("expected branch to not exist")
		}
	})

	t.Run("exec error propagates", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult("", "fatal: error", fmt.Errorf("git error"))

		_, err := app.BranchExists("any")
		if err == nil {
			t.Error("expected error to propagate")
		}
	})
}

func TestListWorktreesAt(t *testing.T) {
	const porcelain = "worktree /repo\nHEAD abc123\nbranch refs/heads/main\n\n" +
		"worktree /repo/.worktrees/feature-a\nHEAD def456\nbranch refs/heads/feature-a\n\n"

	t.Run("lists worktrees from given directory", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		mockApp.Base.SetExecCommandResult(porcelain, "", nil)
		app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}

		worktrees, err := app.ListWorktreesAt("/repo")
		if err != nil {
			t.Fatalf("ListWorktreesAt failed: %v", err)
		}
		if len(worktrees) != 2 {
			t.Fatalf("expected 2 worktrees, got %d", len(worktrees))
		}
		if worktrees[1].Branch != "feature-a" {
			t.Errorf("expected branch 'feature-a', got %q", worktrees[1].Branch)
		}

		last := mockApp.Base.GetLastExecCommandCall()
		if last == nil {
			t.Fatal("no ExecCommand call recorded")
		}
		// Must use -C flag with the given directory
		if len(last.Args) < 2 || last.Args[0] != "-C" || last.Args[1] != "/repo" {
			t.Errorf("expected -C /repo as first args, got %v", last.Args)
		}
	})

	t.Run("exec error returns error", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		mockApp.Base.SetExecCommandResult("", "fatal", fmt.Errorf("not a git repo"))
		app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}

		_, err := app.ListWorktreesAt("/not-a-repo")
		if err == nil {
			t.Error("expected error for non-git directory")
		}
	})
}

func TestPruneWorktreesAt(t *testing.T) {
	t.Run("runs prune in given directory", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		mockApp.Base.SetExecCommandResult("", "", nil)
		app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}

		if err := app.PruneWorktreesAt("/repo"); err != nil {
			t.Fatalf("PruneWorktreesAt failed: %v", err)
		}

		last := mockApp.Base.GetLastExecCommandCall()
		if last == nil {
			t.Fatal("no ExecCommand call recorded")
		}
		if len(last.Args) < 2 || last.Args[0] != "-C" || last.Args[1] != "/repo" {
			t.Errorf("expected -C /repo as first args, got %v", last.Args)
		}
		// Must include worktree prune
		found := false
		for _, arg := range last.Args {
			if arg == "prune" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected 'prune' in args, got %v", last.Args)
		}
	})

	t.Run("returns error on failure", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		mockApp.Base.SetExecCommandResult("", "error", fmt.Errorf("not a git repo"))
		app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}

		if err := app.PruneWorktreesAt("/bad"); err == nil {
			t.Error("expected error")
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

	t.Run("new branch bases off origin default branch when available", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		// Call sequence inside CreateWorktree:
		// 1. fetch origin
		// 2. BranchExists(branch)        -> none (empty)
		// 3. RemoteBranchExists(branch)  -> none (empty)
		// 4. DefaultBranch symbolic-ref  -> "origin/main"
		// 5. RemoteBranchExists("main")  -> exists
		// 6. worktree add ... -b ... origin/main
		mockApp.Base.SetExecCommandResults(
			commands.ExecCommandResult("", "", nil),
			commands.ExecCommandResult("", "", nil),
			commands.ExecCommandResult("", "", nil),
			commands.ExecCommandResult("origin/main\n", "", nil),
			commands.ExecCommandResult("  origin/main\n", "", nil),
			commands.ExecCommandResult("", "", nil),
		)

		if err := app.CreateWorktree("/path/to/worktree", "feature-branch"); err != nil {
			t.Fatalf("CreateWorktree failed: %v", err)
		}

		lastCall := mockApp.Base.GetLastExecCommandCall()
		if lastCall == nil {
			t.Fatal("Expected a worktree add call")
		}
		if lastCall.Args[len(lastCall.Args)-1] != "origin/main" {
			t.Fatalf("Expected new branch to base off 'origin/main', got args: %v", lastCall.Args)
		}
		hasNewBranchFlag := false
		for _, arg := range lastCall.Args {
			if arg == "-b" {
				hasNewBranchFlag = true
			}
		}
		if !hasNewBranchFlag {
			t.Fatalf("Expected -b flag for creating new branch, got args: %v", lastCall.Args)
		}
	})

	t.Run("new branch falls back to HEAD when origin default missing", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		// origin/main does not exist (offline / no remote): all checks empty.
		mockApp.Base.SetExecCommandResults(
			commands.ExecCommandResult("", "", nil),
		)

		if err := app.CreateWorktree("/path/to/worktree", "feature-branch"); err != nil {
			t.Fatalf("CreateWorktree failed: %v", err)
		}

		lastCall := mockApp.Base.GetLastExecCommandCall()
		if lastCall == nil {
			t.Fatal("Expected a worktree add call")
		}
		// Should end at "-b feature-branch" with no base ref appended.
		last := lastCall.Args[len(lastCall.Args)-1]
		if last != "feature-branch" {
			t.Fatalf("Expected fallback to HEAD (no base ref), got args: %v", lastCall.Args)
		}
	})

	t.Run("existing local branch is fast-forwarded to remote", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		// 1. fetch origin
		// 2. BranchExists(branch)        -> exists
		// 3. worktree add path branch
		// 4. RemoteBranchExists(branch)  -> exists
		// 5. merge --ff-only origin/branch
		mockApp.Base.SetExecCommandResults(
			commands.ExecCommandResult("", "", nil),
			commands.ExecCommandResult("  feature-branch\n", "", nil),
			commands.ExecCommandResult("", "", nil),
			commands.ExecCommandResult("  origin/feature-branch\n", "", nil),
			commands.ExecCommandResult("Updating abc..def\n", "", nil),
		)

		if err := app.CreateWorktree("/path/to/worktree", "feature-branch"); err != nil {
			t.Fatalf("CreateWorktree failed: %v", err)
		}

		lastCall := mockApp.Base.GetLastExecCommandCall()
		if lastCall == nil {
			t.Fatal("Expected a merge call")
		}
		joined := strings.Join(lastCall.Args, " ")
		if !strings.Contains(joined, "merge --ff-only origin/feature-branch") {
			t.Fatalf(
				"Expected ff-only merge against origin/feature-branch, got args: %v",
				lastCall.Args,
			)
		}
	})

	t.Run("diverged local branch warns but still succeeds", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResults(
			commands.ExecCommandResult("", "", nil), // fetch
			commands.ExecCommandResult(
				"  feature-branch\n",
				"",
				nil,
			), // BranchExists -> exists
			commands.ExecCommandResult("", "", nil), // worktree add
			commands.ExecCommandResult(
				"  origin/feature-branch\n",
				"",
				nil,
			), // RemoteBranchExists -> exists
			commands.ExecCommandResult(
				"",
				"fatal: Not possible to fast-forward",
				fmt.Errorf("diverged"),
			),
		)

		// Diverged ff-merge must not fail the operation: the worktree was created.
		if err := app.CreateWorktree("/path/to/worktree", "feature-branch"); err != nil {
			t.Fatalf("Expected success despite divergence, got: %v", err)
		}
	})

	t.Run("existing local branch with no remote skips sync", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResults(
			commands.ExecCommandResult("", "", nil),               // fetch
			commands.ExecCommandResult("  local-only\n", "", nil), // BranchExists -> exists
			commands.ExecCommandResult("", "", nil),               // worktree add
			commands.ExecCommandResult("", "", nil),               // RemoteBranchExists -> none
		)

		if err := app.CreateWorktree("/path/to/worktree", "local-only"); err != nil {
			t.Fatalf("CreateWorktree failed: %v", err)
		}

		// No merge should have been attempted (last call is RemoteBranchExists).
		lastCall := mockApp.Base.GetLastExecCommandCall()
		if lastCall == nil {
			t.Fatal("Expected at least one call")
		}
		for _, arg := range lastCall.Args {
			if arg == "merge" {
				t.Fatalf(
					"Expected no merge for a branch with no remote, got args: %v",
					lastCall.Args,
				)
			}
		}
	})

	t.Run("creation error on worktree add", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult(
			"",
			"fatal: worktree exists",
			fmt.Errorf("worktree exists"),
		)

		err := app.CreateWorktree("/path/to/worktree", "existing-branch")
		if err == nil {
			t.Fatal("Expected error but got none")
		}
		if !strings.Contains(err.Error(), "worktree exists") {
			t.Errorf("Expected error to contain 'worktree exists', got: %v", err)
		}
	})
}

func TestDefaultBranch(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}

	t.Run("resolves origin/HEAD", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResults(
			commands.ExecCommandResult("origin/develop\n", "", nil),
		)
		if got := app.DefaultBranch(); got != "develop" {
			t.Fatalf("Expected 'develop', got %q", got)
		}
	})

	t.Run("falls back to main when origin/HEAD unset", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResults(
			commands.ExecCommandResult("", "fatal: ref not found", fmt.Errorf("no origin/HEAD")),
		)
		if got := app.DefaultBranch(); got != "main" {
			t.Fatalf("Expected fallback 'main', got %q", got)
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
			t.Errorf(
				"Expected path '/Users/test/repo/.worktrees/feature', got %q",
				worktrees[1].Path,
			)
		}
		if worktrees[1].Branch != "feature" {
			t.Errorf("Expected branch 'feature', got %q", worktrees[1].Branch)
		}
	})

	t.Run("list error", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult(
			"",
			"fatal: not a git repository",
			fmt.Errorf("not a repo"),
		)

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
		// First call is getMainWorktree (worktree list), subsequent calls succeed with same output
		mockApp.Base.SetExecCommandResult(
			"worktree /main/repo\nHEAD abc123\nbranch refs/heads/main\n",
			"",
			nil,
		)

		if err := app.RemoveWorktree("/path/to/worktree", false, ""); err != nil {
			t.Fatalf("RemoveWorktree failed: %v", err)
		}
		if mockApp.Base.GetExecCommandCallCount() != 2 {
			t.Fatalf("Expected 2 calls, got %d", mockApp.Base.GetExecCommandCallCount())
		}

		// First call: getMainWorktree
		firstCall := mockApp.Base.ExecCommandCalls[0]
		expectedFirst := []string{"-C", "/path/to/worktree", "worktree", "list", "--porcelain"}
		if len(firstCall.Args) != len(expectedFirst) {
			t.Fatalf(
				"Expected %d args for first call, got %d",
				len(expectedFirst),
				len(firstCall.Args),
			)
		}

		// Second call: worktree remove from main repo
		secondCall := mockApp.Base.ExecCommandCalls[1]
		expectedSecond := []string{"-C", "/main/repo", "worktree", "remove", "/path/to/worktree"}
		if len(secondCall.Args) != len(expectedSecond) {
			t.Fatalf(
				"Expected %d args for second call, got %d: %v",
				len(expectedSecond),
				len(secondCall.Args),
				secondCall.Args,
			)
		}
		for i, arg := range expectedSecond {
			if secondCall.Args[i] != arg {
				t.Fatalf("Expected arg[%d] to be %q, got %q", i, arg, secondCall.Args[i])
			}
		}
	})

	t.Run("successful removal with branch deletion", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult(
			"worktree /main/repo\nHEAD abc123\nbranch refs/heads/main\n",
			"",
			nil,
		)

		if err := app.RemoveWorktree("/path/to/worktree", true, "feature-branch"); err != nil {
			t.Fatalf("RemoveWorktree failed: %v", err)
		}
		if mockApp.Base.GetExecCommandCallCount() != 3 {
			t.Fatalf("Expected 3 calls, got %d", mockApp.Base.GetExecCommandCallCount())
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
		mockApp.Base.SetExecCommandResult(
			"",
			"fatal: not a git repository",
			fmt.Errorf("not a repo"),
		)

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

func TestIsWorktreeDirty(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}

	t.Run("clean worktree", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult("", "", nil)

		dirty, err := app.IsWorktreeDirty("/path/to/worktree")
		if err != nil {
			t.Fatalf("IsWorktreeDirty failed: %v", err)
		}
		if dirty {
			t.Error("Expected clean worktree")
		}

		lastCall := mockApp.Base.GetLastExecCommandCall()
		if lastCall == nil {
			t.Fatal("No ExecCommand call recorded")
		}
		expectedArgs := []string{"-C", "/path/to/worktree", "status", "--porcelain"}
		if len(lastCall.Args) != len(expectedArgs) {
			t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(lastCall.Args))
		}
		for i, arg := range expectedArgs {
			if lastCall.Args[i] != arg {
				t.Fatalf("Expected arg[%d] to be %q, got %q", i, arg, lastCall.Args[i])
			}
		}
	})

	t.Run("dirty worktree", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult("M file.go\n", "", nil)

		dirty, err := app.IsWorktreeDirty("/path/to/worktree")
		if err != nil {
			t.Fatalf("IsWorktreeDirty failed: %v", err)
		}
		if !dirty {
			t.Error("Expected dirty worktree")
		}
	})
}

func TestCheckHookCompatibility(t *testing.T) {
	t.Run("no hooks directory returns no warnings", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}
		mockApp.Base.SetExecCommandResult("", "exit status 1", fmt.Errorf("exit status 1"))

		warnings := app.CheckHookCompatibility(t.TempDir())
		if len(warnings) != 0 {
			t.Errorf("Expected no warnings, got %v", warnings)
		}
	})

	t.Run("hook with [ -d .git pattern triggers warning", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}
		mockApp.Base.SetExecCommandResult("", "exit status 1", fmt.Errorf("exit status 1"))

		tmpDir := t.TempDir()
		hooksDir := filepath.Join(tmpDir, ".git", "hooks")
		if err := os.MkdirAll(hooksDir, 0o755); err != nil {
			t.Fatal(err)
		}
		hookContent := "#!/bin/bash\n[ -d .git ] || { echo 'no .git directory found'; exit 1; }\n"
		if err := os.WriteFile(
			filepath.Join(hooksDir, "pre-commit"),
			[]byte(hookContent),
			0o755,
		); err != nil {
			t.Fatal(err)
		}

		warnings := app.CheckHookCompatibility(tmpDir)
		if len(warnings) != 1 {
			t.Fatalf("Expected 1 warning, got %d: %v", len(warnings), warnings)
		}
		if !strings.Contains(warnings[0], "pre-commit") {
			t.Errorf("Expected warning to mention pre-commit, got %q", warnings[0])
		}
		if !strings.Contains(warnings[0], `[ -d .git`) {
			t.Errorf("Expected warning to mention the pattern, got %q", warnings[0])
		}
	})

	t.Run("hook with test -d .git pattern triggers warning", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}
		mockApp.Base.SetExecCommandResult("", "exit status 1", fmt.Errorf("exit status 1"))

		tmpDir := t.TempDir()
		hooksDir := filepath.Join(tmpDir, ".git", "hooks")
		if err := os.MkdirAll(hooksDir, 0o755); err != nil {
			t.Fatal(err)
		}
		hookContent := "#!/bin/bash\ntest -d .git || exit 1\n"
		if err := os.WriteFile(
			filepath.Join(hooksDir, "commit-msg"),
			[]byte(hookContent),
			0o755,
		); err != nil {
			t.Fatal(err)
		}

		warnings := app.CheckHookCompatibility(tmpDir)
		if len(warnings) != 1 {
			t.Fatalf("Expected 1 warning, got %d: %v", len(warnings), warnings)
		}
		if !strings.Contains(warnings[0], "commit-msg") {
			t.Errorf("Expected warning to mention commit-msg, got %q", warnings[0])
		}
	})

	t.Run("hook using git rev-parse is compatible", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}
		mockApp.Base.SetExecCommandResult("", "exit status 1", fmt.Errorf("exit status 1"))

		tmpDir := t.TempDir()
		hooksDir := filepath.Join(tmpDir, ".git", "hooks")
		if err := os.MkdirAll(hooksDir, 0o755); err != nil {
			t.Fatal(err)
		}
		hookContent := "#!/bin/bash\ngit_dir=$(git rev-parse --git-dir)\necho \"git dir: $git_dir\"\n"
		if err := os.WriteFile(
			filepath.Join(hooksDir, "pre-commit"),
			[]byte(hookContent),
			0o755,
		); err != nil {
			t.Fatal(err)
		}

		warnings := app.CheckHookCompatibility(tmpDir)
		if len(warnings) != 0 {
			t.Errorf("Expected no warnings for compatible hook, got %v", warnings)
		}
	})

	t.Run("respects custom core.hooksPath", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}

		tmpDir := t.TempDir()
		huskyDir := filepath.Join(tmpDir, ".husky")
		if err := os.MkdirAll(huskyDir, 0o755); err != nil {
			t.Fatal(err)
		}
		hookContent := "#!/bin/sh\n[ -d .git ] || exit 1\n"
		if err := os.WriteFile(
			filepath.Join(huskyDir, "pre-commit"),
			[]byte(hookContent),
			0o755,
		); err != nil {
			t.Fatal(err)
		}
		mockApp.Base.SetExecCommandResult(".husky\n", "", nil)

		warnings := app.CheckHookCompatibility(tmpDir)
		if len(warnings) != 1 {
			t.Fatalf("Expected 1 warning for husky hook, got %d: %v", len(warnings), warnings)
		}
		if !strings.Contains(warnings[0], "pre-commit") {
			t.Errorf("Expected pre-commit warning, got %q", warnings[0])
		}
	})

	t.Run("multiple hooks each trigger one warning", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}
		mockApp.Base.SetExecCommandResult("", "exit status 1", fmt.Errorf("exit status 1"))

		tmpDir := t.TempDir()
		hooksDir := filepath.Join(tmpDir, ".git", "hooks")
		if err := os.MkdirAll(hooksDir, 0o755); err != nil {
			t.Fatal(err)
		}
		badHook := "#!/bin/bash\n[ -d .git ] || exit 1\n"
		for _, name := range []string{"pre-commit", "commit-msg"} {
			if err := os.WriteFile(
				filepath.Join(hooksDir, name),
				[]byte(badHook),
				0o755,
			); err != nil {
				t.Fatal(err)
			}
		}

		warnings := app.CheckHookCompatibility(tmpDir)
		if len(warnings) != 2 {
			t.Errorf("Expected 2 warnings, got %d: %v", len(warnings), warnings)
		}
	})

	t.Run("affiance hooks trigger warning and skip other checks", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}
		mockApp.Base.SetExecCommandResult("", "exit status 1", fmt.Errorf("exit status 1"))

		tmpDir := t.TempDir()
		hooksDir := filepath.Join(tmpDir, ".git", "hooks")
		if err := os.MkdirAll(hooksDir, 0o755); err != nil {
			t.Fatal(err)
		}

		// Create affiance-hook file (indicates Affiance is installed)
		affianceContent := "#!/bin/bash\nhook=`basename \"$0\"`\nnode $DIR/affiance-hook.js \"$hook\" \"$@\"\n"
		if err := os.WriteFile(
			filepath.Join(hooksDir, "affiance-hook"),
			[]byte(affianceContent),
			0o755,
		); err != nil {
			t.Fatal(err)
		}

		// Also create a pre-commit with bad pattern - should be ignored when Affiance is present
		badHook := "#!/bin/bash\n[ -d .git ] || exit 1\n"
		if err := os.WriteFile(
			filepath.Join(hooksDir, "pre-commit"),
			[]byte(badHook),
			0o755,
		); err != nil {
			t.Fatal(err)
		}

		warnings := app.CheckHookCompatibility(tmpDir)
		if len(warnings) != 1 {
			t.Fatalf("Expected 1 warning for Affiance, got %d: %v", len(warnings), warnings)
		}
		if !strings.Contains(warnings[0], "affiance") {
			t.Errorf("Expected warning to mention affiance, got %q", warnings[0])
		}
		if !strings.Contains(warnings[0], "--no-verify") {
			t.Errorf("Expected warning to suggest --no-verify, got %q", warnings[0])
		}
	})
}

func TestPruneWorktrees(t *testing.T) {
	mockApp := testutil.NewMockApp()
	app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}

	t.Run("successful prune", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult("", "", nil)

		if err := app.PruneWorktrees(); err != nil {
			t.Fatalf("PruneWorktrees failed: %v", err)
		}

		lastCall := mockApp.Base.GetLastExecCommandCall()
		if lastCall == nil {
			t.Fatal("No ExecCommand call recorded")
		}
		expectedArgs := []string{"worktree", "prune"}
		if len(lastCall.Args) != len(expectedArgs) {
			t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(lastCall.Args))
		}
		for i, arg := range expectedArgs {
			if lastCall.Args[i] != arg {
				t.Fatalf("Expected arg[%d] to be %q, got %q", i, arg, lastCall.Args[i])
			}
		}
	})

	t.Run("prune error", func(t *testing.T) {
		mockApp.Base.ResetExecCommand()
		mockApp.Base.SetExecCommandResult(
			"",
			"fatal: not a git repository",
			fmt.Errorf("not a repo"),
		)

		err := app.PruneWorktrees()
		if err == nil {
			t.Fatal("Expected error but got none")
		}
		if !strings.Contains(err.Error(), "git: fatal: not a git repository") {
			t.Fatalf(
				"Expected error message to contain 'git: fatal: not a git repository', got: %v",
				err,
			)
		}
	})
}

func TestDiff(t *testing.T) {
	t.Run("diff HEAD succeeds returns content", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		// The mock returns the same value for every call.
		mockApp.Base.SetExecCommandResult("diff --git a/foo.go b/foo.go\n+added line\n", "", nil)

		app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}
		content, err := app.Diff("/tmp/repo")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if !strings.Contains(content, "+added line") {
			t.Errorf("Expected diff content, got %q", content)
		}
	})

	t.Run("diff HEAD fails returns error", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		// Both diff HEAD and fallback diff fail
		mockApp.Base.SetExecCommandResult("", "fatal: bad revision 'HEAD'", fmt.Errorf("git error"))

		app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}
		_, err := app.Diff("/tmp/repo")
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

func TestDiffStat(t *testing.T) {
	t.Run("numstat parses tracked changes", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		// Mock returns same value for all calls. DiffStat calls numstat then status --porcelain.
		// Since the mock returns "5\t2\tfoo.go\n3\t1\tbar.go\n" for both calls, the status
		// call returns no "?? " lines, so only tracked file counts apply.
		mockApp.Base.SetExecCommandResult("5\t2\tfoo.go\n3\t1\tbar.go\n", "", nil)

		app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}
		files, added, removed, err := app.DiffStat("/tmp/repo")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if files < 2 {
			t.Errorf("Expected at least 2 files, got %d", files)
		}
		if added < 8 {
			t.Errorf("Expected at least 8 added lines, got %d", added)
		}
		if removed < 3 {
			t.Errorf("Expected at least 3 removed lines, got %d", removed)
		}
	})

	t.Run("numstat fails returns error", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		// Both numstat HEAD and fallback fail
		mockApp.Base.SetExecCommandResult("", "", fmt.Errorf("git error"))

		app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}
		_, _, _, err := app.DiffStat("/tmp/repo")
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})

	t.Run("empty output returns zero stats", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		mockApp.Base.SetExecCommandResult("", "", nil)

		app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}
		files, added, removed, err := app.DiffStat("/tmp/repo")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if files != 0 || added != 0 || removed != 0 {
			t.Errorf("Expected 0/0/0, got %d/%d/%d", files, added, removed)
		}
	})

	t.Run("untracked files increment file count", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		// Call 1: numstat HEAD — two tracked files
		// Call 2: status --porcelain — one untracked file
		mockApp.Base.SetExecCommandResults(
			commands.ExecCommandResult("3\t1\tfoo.go\n", "", nil),
			commands.ExecCommandResult("?? newfile.go\n", "", nil),
		)

		app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}
		files, added, removed, err := app.DiffStat("/tmp/repo")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if files != 2 {
			t.Errorf("expected 2 files (1 tracked + 1 untracked), got %d", files)
		}
		if added != 3 {
			t.Errorf("expected 3 added lines, got %d", added)
		}
		if removed != 1 {
			t.Errorf("expected 1 removed line, got %d", removed)
		}
	})

	t.Run("binary files counted but contribute zero lines", func(t *testing.T) {
		mockApp := testutil.NewMockApp()
		// Binary files show as "-\t-\tfile" in numstat
		mockApp.Base.SetExecCommandResults(
			commands.ExecCommandResult("-\t-\timage.png\n5\t2\tfoo.go\n", "", nil),
			commands.ExecCommandResult("", "", nil),
		)

		app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}
		files, added, removed, err := app.DiffStat("/tmp/repo")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if files != 2 {
			t.Errorf("expected 2 files (1 binary + 1 text), got %d", files)
		}
		if added != 5 {
			t.Errorf("expected 5 added (binary contributes 0), got %d", added)
		}
		if removed != 2 {
			t.Errorf("expected 2 removed (binary contributes 0), got %d", removed)
		}
	})
}

func TestDiffIncludesUntrackedFiles(t *testing.T) {
	mockApp := testutil.NewMockApp()
	// Call 1: diff --color=always HEAD — some tracked diff
	// Call 2: status --porcelain — one untracked file
	mockApp.Base.SetExecCommandResults(
		commands.ExecCommandResult("+added line\n", "", nil),
		commands.ExecCommandResult("?? brand_new.go\n", "", nil),
	)

	app := &Git{Cmd: mockApp.Cmd, Base: mockApp.Base}
	content, err := app.Diff("/tmp/repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(content, "brand_new.go") {
		t.Errorf("expected untracked file 'brand_new.go' in diff output, got:\n%s", content)
	}
	if !strings.Contains(content, "Untracked files:") {
		t.Errorf("expected 'Untracked files:' section in diff output")
	}
}
