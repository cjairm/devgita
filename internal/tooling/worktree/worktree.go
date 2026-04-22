// Worktree coordinator manages git worktrees with tmux window integration
//
// Each worktree gets its own tmux window with OpenCode running, enabling
// parallel AI-assisted development across multiple branches within the same session.
// This follows the "one session per folder" workflow where worktrees are managed
// as separate windows rather than separate sessions.
//
// References:
// - Git Worktree Documentation: https://git-scm.com/docs/git-worktree
// - Tmux Manual: https://man7.org/linux/man-pages/man1/tmux.1.html

package worktree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cjairm/devgita/internal/apps/git"
	"github.com/cjairm/devgita/internal/apps/tmux"
	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/fzf"
	"github.com/cjairm/devgita/pkg/constants"
)

const (
	// worktreeDir is the directory within the repo where worktrees are created
	worktreeDir = ".worktrees"
	// windowPrefix is prepended to worktree names for tmux windows
	windowPrefix = "wt-"
)

// WorktreeStatus contains information about a worktree and its associated window
type WorktreeStatus struct {
	Name         string
	Path         string
	Branch       string
	TmuxWindow   string
	WindowActive bool
}

// WorktreeManager coordinates git worktrees with tmux windows
type WorktreeManager struct {
	Git  *git.Git
	Tmux *tmux.Tmux
	Fzf  *fzf.Fzf
	Base cmd.BaseCommandExecutor
}

// New creates a new WorktreeManager instance
func New() *WorktreeManager {
	return &WorktreeManager{
		Git:  git.New(),
		Tmux: tmux.New(),
		Fzf:  fzf.New(),
		Base: cmd.NewBaseCommand(),
	}
}

// Create creates a new worktree with tmux window and launches OpenCode
func (w *WorktreeManager) Create(name string) error {
	// 1. Validate we're in a git repo
	repoRoot, err := w.Git.GetRepoRoot()
	if err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	// 2. Check worktree doesn't exist
	wtPath := filepath.Join(repoRoot, worktreeDir, name)
	if _, err := os.Stat(wtPath); err == nil {
		return fmt.Errorf("worktree '%s' already exists", name)
	}

	// 3. Check window doesn't exist
	windowName := windowPrefix + name
	if w.Tmux.HasWindow(windowName) {
		return fmt.Errorf("tmux window '%s' already exists", windowName)
	}

	// 4. Create worktree
	if err := w.Git.CreateWorktree(wtPath, name); err != nil {
		return fmt.Errorf("failed to create worktree: %w", err)
	}

	// 5. Create tmux window
	if err := w.Tmux.CreateWindow(windowName, wtPath); err != nil {
		// Rollback: remove worktree and delete branch if window creation fails
		_ = w.Git.RemoveWorktree(wtPath, true, name)
		return fmt.Errorf("failed to create tmux window: %w", err)
	}

	// 6. Launch OpenCode in window
	if err := w.Tmux.SendKeysToWindow(windowName, constants.OpenCode); err != nil {
		return fmt.Errorf("failed to launch opencode: %w", err)
	}

	return nil
}

// List returns all worktrees with their window status
func (w *WorktreeManager) List() ([]WorktreeStatus, error) {
	worktrees, err := w.Git.ListWorktrees()
	if err != nil {
		return nil, err
	}

	var statuses []WorktreeStatus
	for _, wt := range worktrees {
		// Skip main worktree (not in .worktrees/)
		if !strings.Contains(wt.Path, worktreeDir) {
			continue
		}
		name := filepath.Base(wt.Path)
		windowName := windowPrefix + name
		statuses = append(statuses, WorktreeStatus{
			Name:         name,
			Path:         wt.Path,
			Branch:       wt.Branch,
			TmuxWindow:   windowName,
			WindowActive: w.Tmux.HasWindow(windowName),
		})
	}
	return statuses, nil
}

// Remove removes a worktree and its tmux window
func (w *WorktreeManager) Remove(name string) error {
	repoRoot, err := w.Git.GetRepoRoot()
	if err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	wtPath := filepath.Join(repoRoot, worktreeDir, name)
	windowName := windowPrefix + name

	// Kill tmux window if exists
	if w.Tmux.HasWindow(windowName) {
		if err := w.Tmux.KillWindow(windowName); err != nil {
			return fmt.Errorf("failed to kill tmux window: %w", err)
		}
	}

	// Remove worktree and delete the associated branch
	if err := w.Git.RemoveWorktree(wtPath, true, name); err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
	}

	return nil
}

// GetWindowName returns the tmux window name for a given worktree name
func GetWindowName(name string) string {
	return windowPrefix + name
}

// GetWorktreeDir returns the worktree directory name
func GetWorktreeDir() string {
	return worktreeDir
}

// SelectWorktreeInteractively presents an fzf picker with available worktrees
// and returns the selected worktree name. Returns error if no worktrees exist
// or user cancels selection.
func (w *WorktreeManager) SelectWorktreeInteractively(prompt string) (string, error) {
	statuses, err := w.List()
	if err != nil {
		return "", fmt.Errorf("failed to list worktrees: %w", err)
	}

	if len(statuses) == 0 {
		return "", fmt.Errorf("no worktrees available")
	}

	// Extract worktree names for fzf
	names := make([]string, len(statuses))
	for i, s := range statuses {
		names[i] = s.Name
	}

	selected, err := w.Fzf.SelectFromList(names, prompt)
	if err != nil {
		return "", fmt.Errorf("selection failed: %w", err)
	}

	return selected, nil
}
