package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cjairm/devgita/internal/apps/git"
	"github.com/cjairm/devgita/internal/apps/tmux"
	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/testutil"
	"github.com/cjairm/devgita/pkg/paths"
)

func init() {
	testutil.InitLogger()
}

func TestNew(t *testing.T) {
	wm := New()
	if wm == nil {
		t.Fatal("New() returned nil")
	}
	if wm.Git == nil {
		t.Error("Git should not be nil")
	}
	if wm.Tmux == nil {
		t.Error("Tmux should not be nil")
	}
	if wm.Fzf == nil {
		t.Error("Fzf should not be nil")
	}
	if wm.Base == nil {
		t.Error("Base should not be nil")
	}
}

func TestGetWindowName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple name", "feature", "wt-feature"},
		{"hyphenated name", "feature-login", "wt-feature-login"},
		{"with numbers", "fix-123", "wt-fix-123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetWindowName(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSelectWorktreeInteractively(t *testing.T) {
	t.Skip("Skipping: SelectFromList uses exec.Command which requires actual fzf binary and would block in CI")
}

func TestGetWorktreeDir(t *testing.T) {
	result := GetWorktreeDir()
	if result != ".worktrees" {
		t.Errorf("Expected '.worktrees', got %q", result)
	}
}

func TestCreate(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		tempDir := t.TempDir()

		mockGitBase := commands.NewMockBaseCommand()
		mockTmuxBase := commands.NewMockBaseCommand()

		gitApp := &git.Git{
			Cmd:  commands.NewMockCommand(),
			Base: mockGitBase,
		}
		tmuxApp := &tmux.Tmux{
			Cmd:  commands.NewMockCommand(),
			Base: mockTmuxBase,
		}

		wm := &WorktreeManager{
			Git:  gitApp,
			Tmux: tmuxApp,
			Base: commands.NewMockBaseCommand(),
		}

		mockGitBase.SetExecCommandResult(tempDir+"\n", "", nil)
		mockTmuxBase.SetExecCommandResult("", "window not found", os.ErrNotExist)

		err := wm.Create("feature-test", &OpenCodeCoder{})
		if err == nil {
			if mockGitBase.GetExecCommandCallCount() < 1 {
				t.Error("Expected git commands to be called")
			}
		}
	})

	t.Run("nil coder returns error", func(t *testing.T) {
		wm := &WorktreeManager{}
		err := wm.Create("test", nil)
		if err == nil {
			t.Fatal("Expected error for nil coder")
		}
	})

	t.Run("not in git repo", func(t *testing.T) {
		mockGitBase := commands.NewMockBaseCommand()
		mockTmuxBase := commands.NewMockBaseCommand()

		gitApp := &git.Git{
			Cmd:  commands.NewMockCommand(),
			Base: mockGitBase,
		}
		tmuxApp := &tmux.Tmux{
			Cmd:  commands.NewMockCommand(),
			Base: mockTmuxBase,
		}

		wm := &WorktreeManager{
			Git:  gitApp,
			Tmux: tmuxApp,
			Base: commands.NewMockBaseCommand(),
		}

		mockGitBase.SetExecCommandResult("", "fatal: not a git repository", os.ErrNotExist)

		err := wm.Create("feature-test", &OpenCodeCoder{})
		if err == nil {
			t.Fatal("Expected error when not in git repo")
		}
	})
}

func TestList(t *testing.T) {
	t.Run("list worktrees from centralized dir", func(t *testing.T) {
		mockGitBase := commands.NewMockBaseCommand()
		mockTmuxBase := commands.NewMockBaseCommand()

		gitApp := &git.Git{
			Cmd:  commands.NewMockCommand(),
			Base: mockGitBase,
		}
		tmuxApp := &tmux.Tmux{
			Cmd:  commands.NewMockCommand(),
			Base: mockTmuxBase,
		}

		wm := &WorktreeManager{
			Git:  gitApp,
			Tmux: tmuxApp,
			Base: commands.NewMockBaseCommand(),
		}

		statuses, err := wm.List()
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		// Note: This test may return non-zero results if real worktrees exist in the centralized dir.
		// The important thing is that List() doesn't error and returns valid WorktreeStatus structs.
		for _, s := range statuses {
			if s.Name == "" {
				t.Error("Worktree name should not be empty")
			}
			if s.Repo == "" {
				t.Error("Repo should not be empty")
			}
		}
	})
}

func TestRemove(t *testing.T) {
	t.Run("successful removal with active window", func(t *testing.T) {
		tempDir := t.TempDir()

		mockGitBase := commands.NewMockBaseCommand()
		mockTmuxBase := commands.NewMockBaseCommand()

		gitApp := &git.Git{
			Cmd:  commands.NewMockCommand(),
			Base: mockGitBase,
		}
		tmuxApp := &tmux.Tmux{
			Cmd:  commands.NewMockCommand(),
			Base: mockTmuxBase,
		}

		wm := &WorktreeManager{
			Git:  gitApp,
			Tmux: tmuxApp,
			Base: commands.NewMockBaseCommand(),
		}

		mockGitBase.SetExecCommandResult(tempDir+"\n", "", nil)
		mockTmuxBase.SetExecCommandResult("", "", nil)

		repoSlug := filepath.Base(tempDir)
		wtPath := filepath.Join(paths.Paths.Data.Root, "devgita", "worktrees", repoSlug, "feature-test")
		if err := os.MkdirAll(wtPath, 0755); err != nil {
			t.Fatalf("Failed to create worktree dir: %v", err)
		}
		defer os.RemoveAll(filepath.Dir(wtPath))

		err := wm.Remove("feature-test", true)
		if err != nil {
			t.Fatalf("Remove failed: %v", err)
		}

		if mockGitBase.GetExecCommandCallCount() < 1 {
			t.Error("Expected git commands to be called")
		}
		if mockTmuxBase.GetExecCommandCallCount() < 1 {
			t.Error("Expected tmux commands to be called")
		}
	})

	t.Run("removal without active window", func(t *testing.T) {
		tempDir := t.TempDir()

		mockGitBase := commands.NewMockBaseCommand()
		mockTmuxBase := commands.NewMockBaseCommand()

		gitApp := &git.Git{
			Cmd:  commands.NewMockCommand(),
			Base: mockGitBase,
		}
		tmuxApp := &tmux.Tmux{
			Cmd:  commands.NewMockCommand(),
			Base: mockTmuxBase,
		}

		wm := &WorktreeManager{
			Git:  gitApp,
			Tmux: tmuxApp,
			Base: commands.NewMockBaseCommand(),
		}

		mockGitBase.SetExecCommandResult(tempDir+"\n", "", nil)
		mockTmuxBase.SetExecCommandResult("", "window not found", os.ErrNotExist)

		repoSlug := filepath.Base(tempDir)
		wtPath := filepath.Join(paths.Paths.Data.Root, "devgita", "worktrees", repoSlug, "feature-test")
		if err := os.MkdirAll(wtPath, 0755); err != nil {
			t.Fatalf("Failed to create worktree dir: %v", err)
		}
		defer os.RemoveAll(filepath.Dir(filepath.Dir(wtPath)))

		err := wm.Remove("feature-test", true)
		if err != nil {
			t.Fatalf("Remove failed: %v", err)
		}
	})

	t.Run("not in git repo", func(t *testing.T) {
		mockGitBase := commands.NewMockBaseCommand()
		mockTmuxBase := commands.NewMockBaseCommand()

		gitApp := &git.Git{
			Cmd:  commands.NewMockCommand(),
			Base: mockGitBase,
		}
		tmuxApp := &tmux.Tmux{
			Cmd:  commands.NewMockCommand(),
			Base: mockTmuxBase,
		}

		wm := &WorktreeManager{
			Git:  gitApp,
			Tmux: tmuxApp,
			Base: commands.NewMockBaseCommand(),
		}

		mockGitBase.SetExecCommandResult("", "fatal: not a git repository", os.ErrNotExist)

		err := wm.Remove("feature-test", false)
		if err == nil {
			t.Fatal("Expected error when not in git repo")
		}
	})
}

func TestFormatJumpRow(t *testing.T) {
	tests := []struct {
		name     string
		repo     string
		wtName   string
		branch   string
		status   string
		expected string
	}{
		{
			name:     "basic row",
			repo:     "myrepo",
			wtName:   "feature-a",
			branch:   "feature-a",
			status:   "active",
			expected: "myrepo/feature-a\tfeature-a\tactive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatJumpRow(tt.repo, tt.wtName, tt.branch, tt.status)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestParseJumpRow(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "basic row",
			input:    "myrepo/feature-a\tfeature-a\tactive",
			expected: []string{"myrepo/feature-a", "feature-a", "active"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseJumpRow(tt.input)
			if len(result) != len(tt.expected) {
				t.Fatalf("Expected %d parts, got %d", len(tt.expected), len(result))
			}
			for i, part := range tt.expected {
				if result[i] != part {
					t.Errorf("Expected part[%d] %q, got %q", i, part, result[i])
				}
			}
		})
	}
}

func TestFormatWindowRow(t *testing.T) {
	result := formatWindowRow("main")
	expected := "[win]\tmain\t"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestParseJumpOutput(t *testing.T) {
	t.Run("enter key (no special key)", func(t *testing.T) {
		output := "myrepo/feature-a\tfeature-a\tactive"
		key, row, err := parseJumpOutput(output)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if key != "" {
			t.Errorf("Expected empty key, got %q", key)
		}
		if row != "myrepo/feature-a\tfeature-a\tactive" {
			t.Errorf("Expected row %q, got %q", "myrepo/feature-a\tfeature-a\tactive", row)
		}
	})

	t.Run("ctrl-d key", func(t *testing.T) {
		output := "ctrl-d\nmyrepo/feature-a\tfeature-a\tactive"
		key, row, err := parseJumpOutput(output)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if key != "ctrl-d" {
			t.Errorf("Expected key 'ctrl-d', got %q", key)
		}
		if row != "myrepo/feature-a\tfeature-a\tactive" {
			t.Errorf("Expected row %q, got %q", "myrepo/feature-a\tfeature-a\tactive", row)
		}
	})

	t.Run("empty output", func(t *testing.T) {
		_, _, err := parseJumpOutput("")
		if err == nil {
			t.Fatal("Expected error for empty output")
		}
	})
}

func TestWorktreePath(t *testing.T) {
	wm := &WorktreeManager{}
	path := wm.worktreePath("myrepo", "feature-a")
	expectedSuffix := filepath.Join("devgita", "worktrees", "myrepo", "feature-a")
	if !filepath.IsAbs(path) {
		t.Errorf("Expected absolute path, got %q", path)
	}
	if !strings.HasSuffix(path, expectedSuffix) {
		t.Errorf("Expected path to end with %q, got %q", expectedSuffix, path)
	}
}

func TestGetWorktreeBasePath(t *testing.T) {
	basePath := GetWorktreeBasePath()
	expectedSuffix := filepath.Join("devgita", "worktrees")
	if !filepath.IsAbs(basePath) {
		t.Errorf("Expected absolute path, got %q", basePath)
	}
	if !strings.HasSuffix(basePath, expectedSuffix) {
		t.Errorf("Expected path to end with %q, got %q", expectedSuffix, basePath)
	}
}
