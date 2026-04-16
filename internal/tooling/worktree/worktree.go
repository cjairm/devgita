// Worktree coordinator manages git worktrees with tmux session integration
//
// Each worktree gets its own tmux session with OpenCode running, enabling
// parallel AI-assisted development across multiple branches.
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
	"github.com/cjairm/devgita/pkg/constants"
)

const (
	// worktreeDir is the directory within the repo where worktrees are created
	worktreeDir = ".worktrees"
	// sessionPrefix is prepended to worktree names for tmux sessions
	sessionPrefix = "wt-"
)

// WorktreeStatus contains information about a worktree and its associated session
type WorktreeStatus struct {
	Name          string
	Path          string
	Branch        string
	TmuxSession   string
	SessionActive bool
}

// WorktreeManager coordinates git worktrees with tmux sessions
type WorktreeManager struct {
	Git  *git.Git
	Tmux *tmux.Tmux
	Base cmd.BaseCommandExecutor
}

// New creates a new WorktreeManager instance
func New() *WorktreeManager {
	return &WorktreeManager{
		Git:  git.New(),
		Tmux: tmux.New(),
		Base: cmd.NewBaseCommand(),
	}
}

// Create creates a new worktree with tmux session and launches OpenCode
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

	// 3. Check session doesn't exist
	sessionName := sessionPrefix + name
	if w.Tmux.HasSession(sessionName) {
		return fmt.Errorf("tmux session '%s' already exists", sessionName)
	}

	// 4. Create worktree
	if err := w.Git.CreateWorktree(wtPath, name); err != nil {
		return fmt.Errorf("failed to create worktree: %w", err)
	}

	// 5. Create tmux session
	if err := w.Tmux.CreateSession(sessionName, wtPath); err != nil {
		// Rollback: remove worktree and delete branch if session creation fails
		_ = w.Git.RemoveWorktree(wtPath, true, name)
		return fmt.Errorf("failed to create tmux session: %w", err)
	}

	// 6. Launch OpenCode in session
	if err := w.Tmux.SendKeys(sessionName, constants.OpenCode); err != nil {
		return fmt.Errorf("failed to launch opencode: %w", err)
	}

	return nil
}

// List returns all worktrees with their session status
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
		sessionName := sessionPrefix + name
		statuses = append(statuses, WorktreeStatus{
			Name:          name,
			Path:          wt.Path,
			Branch:        wt.Branch,
			TmuxSession:   sessionName,
			SessionActive: w.Tmux.HasSession(sessionName),
		})
	}
	return statuses, nil
}

// Remove removes a worktree and its tmux session
func (w *WorktreeManager) Remove(name string) error {
	repoRoot, err := w.Git.GetRepoRoot()
	if err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	wtPath := filepath.Join(repoRoot, worktreeDir, name)
	sessionName := sessionPrefix + name

	// Kill tmux session if exists
	if w.Tmux.HasSession(sessionName) {
		if err := w.Tmux.KillSession(sessionName); err != nil {
			return fmt.Errorf("failed to kill tmux session: %w", err)
		}
	}

	// Remove worktree and delete the associated branch
	if err := w.Git.RemoveWorktree(wtPath, true, name); err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
	}

	return nil
}

// GetSessionName returns the tmux session name for a given worktree name
func GetSessionName(name string) string {
	return sessionPrefix + name
}

// GetWorktreeDir returns the worktree directory name
func GetWorktreeDir() string {
	return worktreeDir
}
