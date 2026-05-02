// Worktree coordinator manages git worktrees with tmux window integration
//
// Each worktree gets its own tmux window with an AI assistant running, enabling
// parallel AI-assisted development across multiple branches within the same session.
// This follows the "one session per folder" workflow where worktrees are managed
// as separate windows rather than separate sessions.
//
// References:
// - Git Worktree Documentation: https://git-scm.com/docs/git-worktree
// - Tmux Manual: https://man7.org/linux/man-pages/man1/tmux.1.html

package worktree

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cjairm/devgita/internal/apps/git"
	"github.com/cjairm/devgita/internal/apps/tmux"
	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/fzf"
	"github.com/cjairm/devgita/pkg/paths"
)

const (
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
	Repo         string
}

// WorktreeState holds the current state of a worktree
type WorktreeState struct {
	WtPath       string
	WindowName   string
	WtExists     bool
	WindowExists bool
	BranchExists bool
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

// worktreePath returns ~/.local/share/devgita/worktrees/<repo-slug>/<name>
func (w *WorktreeManager) worktreePath(repoSlug, name string) string {
	return filepath.Join(paths.Paths.Data.Root, "devgita", "worktrees", repoSlug, name)
}

// GetWorktreeBasePath returns the base path for all devgita worktrees
func GetWorktreeBasePath() string {
	return filepath.Join(paths.Paths.Data.Root, "devgita", "worktrees")
}

// Create creates a new worktree with tmux window and launches the specified AI coder
func (w *WorktreeManager) Create(name string, coder AICoder) error {
	if coder == nil {
		return fmt.Errorf("AI coder is required")
	}

	if err := coder.EnsureInstalled(); err != nil {
		return err
	}

	repoRoot, err := w.Git.GetRepoRoot()
	if err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	repoSlug := filepath.Base(repoRoot)
	wtPath := w.worktreePath(repoSlug, name)
	windowName := windowPrefix + name

	state, err := w.worktreeState(repoSlug, name)
	if err != nil {
		return err
	}

	if state.WtExists && state.WindowExists {
		return fmt.Errorf("worktree '%s' already exists and has an active window; use `dg wt jump %s`", name, name)
	}
	if state.WtExists && !state.WindowExists {
		return fmt.Errorf("worktree '%s' exists but has no active window; use `dg wt repair %s`", name, name)
	}
	if !state.WtExists && state.WindowExists {
		return fmt.Errorf("orphan window '%s' exists; run `tmux kill-window -t %s` manually", windowName, windowName)
	}

	if err := os.MkdirAll(filepath.Dir(wtPath), 0755); err != nil {
		return fmt.Errorf("failed to create worktree directory: %w", err)
	}

	if err := w.Git.CreateWorktree(wtPath, name); err != nil {
		return fmt.Errorf("failed to create worktree: %w", err)
	}

	if err := w.Tmux.CreateWindow(windowName, wtPath); err != nil {
		_ = w.Git.RemoveWorktree(wtPath, true, name)
		return fmt.Errorf("failed to create tmux window: %w", err)
	}

	if err := w.Tmux.SendKeysToWindow(windowName, coder.Command()); err != nil {
		return fmt.Errorf("failed to launch %s: %w", coder.Name(), err)
	}

	return nil
}

// worktreeState checks the current state of a worktree
func (w *WorktreeManager) worktreeState(repoSlug, name string) (WorktreeState, error) {
	state := WorktreeState{
		WtPath:     w.worktreePath(repoSlug, name),
		WindowName: windowPrefix + name,
	}

	if _, err := os.Stat(state.WtPath); err == nil {
		state.WtExists = true
	}

	if w.Tmux.HasWindow(state.WindowName) {
		state.WindowExists = true
	}

	if state.WtExists {
		branchExists, err := w.Git.BranchExists(name)
		if err == nil {
			state.BranchExists = branchExists
		}
	}

	return state, nil
}

// List returns all worktrees with their window status across all repos
func (w *WorktreeManager) List() ([]WorktreeStatus, error) {
	basePath := GetWorktreeBasePath()

	entries, err := os.ReadDir(basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []WorktreeStatus{}, nil
		}
		return nil, err
	}

	var statuses []WorktreeStatus
	for _, repoEntry := range entries {
		if !repoEntry.IsDir() {
			continue
		}

		repoSlug := repoEntry.Name()
		repoDir := filepath.Join(basePath, repoSlug)

		wtEntries, err := os.ReadDir(repoDir)
		if err != nil {
			continue
		}

		for _, wtEntry := range wtEntries {
			if !wtEntry.IsDir() {
				continue
			}

			name := wtEntry.Name()
			wtPath := filepath.Join(repoDir, name)
			windowName := windowPrefix + name

			branch := ""
			worktrees, err := w.Git.ListWorktrees()
			if err == nil {
				for _, wt := range worktrees {
					if wt.Path == wtPath {
						branch = wt.Branch
						break
					}
				}
			}

			statuses = append(statuses, WorktreeStatus{
				Name:         name,
				Path:         wtPath,
				Branch:       branch,
				TmuxWindow:   windowName,
				WindowActive: w.Tmux.HasWindow(windowName),
				Repo:         repoSlug,
			})
		}
	}

	return statuses, nil
}

// Remove removes a worktree and its tmux window
func (w *WorktreeManager) Remove(name string, force bool) error {
	repoRoot, err := w.Git.GetRepoRoot()
	if err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	repoSlug := filepath.Base(repoRoot)
	wtPath := w.worktreePath(repoSlug, name)
	windowName := windowPrefix + name

	state, err := w.worktreeState(repoSlug, name)
	if err != nil {
		return err
	}

	if !state.WtExists && !state.WindowExists {
		return fmt.Errorf("nothing to remove for worktree '%s'", name)
	}

	if state.WtExists && !force {
		dirty, err := w.Git.IsWorktreeDirty(wtPath)
		if err != nil {
			return fmt.Errorf("failed to check worktree status: %w", err)
		}
		if dirty {
			return fmt.Errorf("worktree '%s' has uncommitted changes; use --force to remove anyway", name)
		}
	}

	if state.WindowExists {
		if err := w.Tmux.KillWindow(windowName); err != nil {
			// Log but don't fail
		}
	}

	if state.WtExists {
		if err := w.Git.RemoveWorktree(wtPath, true, name); err != nil {
			return fmt.Errorf("failed to remove worktree: %w", err)
		}
		if err := w.Git.PruneWorktrees(); err != nil {
			// Log but don't fail
		}
	} else {
		if err := w.Git.PruneWorktrees(); err != nil {
			// Log but don't fail
		}
	}

	return nil
}

// Repair creates missing window for an existing worktree and re-sends the AI command
func (w *WorktreeManager) Repair(name string, coder AICoder) error {
	if coder == nil {
		return fmt.Errorf("AI coder is required")
	}

	if err := coder.EnsureInstalled(); err != nil {
		return err
	}

	repoRoot, err := w.Git.GetRepoRoot()
	if err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	repoSlug := filepath.Base(repoRoot)
	wtPath := w.worktreePath(repoSlug, name)
	windowName := windowPrefix + name

	if _, err := os.Stat(wtPath); os.IsNotExist(err) {
		return fmt.Errorf("no worktree '%s' to repair", name)
	}

	if !w.Tmux.HasWindow(windowName) {
		if err := w.Tmux.CreateWindow(windowName, wtPath); err != nil {
			return fmt.Errorf("failed to create tmux window: %w", err)
		}
	}

	if err := w.Tmux.SendKeysToWindow(windowName, coder.Command()); err != nil {
		return fmt.Errorf("failed to launch %s: %w", coder.Name(), err)
	}

	return nil
}

// Prune removes all worktrees in the centralized directory
func (w *WorktreeManager) Prune() error {
	statuses, err := w.List()
	if err != nil {
		return err
	}

	if len(statuses) == 0 {
		fmt.Println("Nothing to prune")
		return nil
	}

	fmt.Println("The following worktrees will be removed:")
	for _, s := range statuses {
		fmt.Printf("  - %s/%s\n", s.Repo, s.Name)
	}

	fmt.Print("Remove all? [y/N]: ")
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		return fmt.Errorf("cancelled")
	}

	if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
		return fmt.Errorf("cancelled")
	}

	var errors []string
	for _, s := range statuses {
		if err := w.removeByRepo(s.Repo, s.Name, false); err != nil {
			errors = append(errors, fmt.Sprintf("%s/%s: %v", s.Repo, s.Name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to remove some worktrees:\n  %s", strings.Join(errors, "\n  "))
	}

	return nil
}

// removeByRepo removes a worktree by repo slug and name
func (w *WorktreeManager) removeByRepo(repoSlug, name string, force bool) error {
	wtPath := w.worktreePath(repoSlug, name)
	windowName := windowPrefix + name

	state, err := w.worktreeState(repoSlug, name)
	if err != nil {
		return err
	}

	if !state.WtExists && !state.WindowExists {
		return nil
	}

	if state.WindowExists {
		if err := w.Tmux.KillWindow(windowName); err != nil {
			// Log but don't fail
		}
	}

	if state.WtExists {
		if err := w.Git.RemoveWorktree(wtPath, true, name); err != nil {
			return fmt.Errorf("failed to remove worktree: %w", err)
		}
		if err := w.Git.PruneWorktrees(); err != nil {
			// Log but don't fail
		}
	} else {
		if err := w.Git.PruneWorktrees(); err != nil {
			// Log but don't fail
		}
	}

	return nil
}

// Jump presents an fzf dialog to jump to a worktree
func (w *WorktreeManager) Jump(resolvedCoder string) error {
	statuses, err := w.List()
	if err != nil {
		return err
	}

	if len(statuses) == 0 {
		return fmt.Errorf("no worktrees found")
	}

	rows := make([]string, 0, len(statuses))
	for _, s := range statuses {
		status := "active"
		if !s.WindowActive {
			status = "inactive"
		}
		rows = append(rows, formatJumpRow(s.Repo, s.Name, s.Branch, status))
	}

	isInTmux := os.Getenv("TMUX") != ""

	if isInTmux {
		windows, err := w.listNonWorktreeWindows()
		if err == nil {
			for _, win := range windows {
				rows = append(rows, formatWindowRow(win))
			}
		}
	}

	if len(rows) == 0 {
		return fmt.Errorf("no worktrees found")
	}

	output, err := w.runFzfWithExpect(rows)
	if err != nil {
		return err
	}

	key, row, err := parseJumpOutput(output)
	if err != nil {
		return err
	}

	if isInTmux {
		if strings.HasPrefix(row, "[win]") {
			if key == "" {
				parts := strings.SplitN(row, "\t", 3)
				if len(parts) >= 2 {
					windowName := strings.TrimSpace(parts[1])
					return w.Tmux.SelectWindow(windowName)
				}
			}
			return nil
		}

		parts := parseJumpRow(row)
		if len(parts) < 2 {
			return fmt.Errorf("invalid selection")
		}

		repoSlug := parts[0]
		name := parts[1]
		windowName := windowPrefix + name

		switch key {
		case "":
			if w.Tmux.HasWindow(windowName) {
				return w.Tmux.SelectWindow(windowName)
			}
			fmt.Printf("Window '%s' not found. Use 'dg wt repair %s' to recreate it.\n", windowName, name)
			return nil
		case "ctrl-x":
			fmt.Printf("Remove worktree '%s/%s'? [y/N]: ", repoSlug, name)
			var response string
			if _, err := fmt.Scanln(&response); err != nil || (strings.ToLower(response) != "y" && strings.ToLower(response) != "yes") {
				return nil
			}
			return w.removeByRepo(repoSlug, name, false)
		case "ctrl-r":
			coder, err := ResolveAICoder(resolvedCoder)
			if err != nil {
				return err
			}
			return w.Repair(name, coder)
		default:
			return nil
		}
	}

	parts := parseJumpRow(row)
	if len(parts) < 2 {
		return fmt.Errorf("invalid selection")
	}

	repoSlug := parts[0]
	name := parts[1]
	wtPath := w.worktreePath(repoSlug, name)

	switch key {
	case "":
		fmt.Println(wtPath)
		return nil
	case "ctrl-x":
		fmt.Printf("Remove worktree '%s/%s'? [y/N]: ", repoSlug, name)
		var response string
		if _, err := fmt.Scanln(&response); err != nil || (strings.ToLower(response) != "y" && strings.ToLower(response) != "yes") {
			return nil
		}
		return w.removeByRepo(repoSlug, name, false)
	case "ctrl-r":
		coder, err := ResolveAICoder(resolvedCoder)
		if err != nil {
			return err
		}
		if os.Getenv("TMUX") == "" {
			fmt.Println("Warning: tmux not running, window not created")
		}
		return w.Repair(name, coder)
	default:
		return nil
	}
}

// formatJumpRow formats a worktree row for fzf display
func formatJumpRow(repo, name, branch, status string) string {
	return fmt.Sprintf("%s/%s\t%s\t%s", repo, name, branch, status)
}

// formatWindowRow formats a tmux window row for fzf display
func formatWindowRow(name string) string {
	return fmt.Sprintf("[win]\t%s\t", name)
}

// parseJumpRow parses a worktree row back into its components
func parseJumpRow(row string) []string {
	return strings.SplitN(row, "\t", 3)
}

// runFzfWithExpect runs fzf with --expect flag
func (w *WorktreeManager) runFzfWithExpect(rows []string) (string, error) {
	execCommand := cmd.CommandParams{
		Command: "fzf",
		Args: []string{
			"--header", "enter: jump | ctrl-x: delete | ctrl-r: repair",
			"--expect", "ctrl-x,ctrl-r",
			"--with-nth", "1,2,3",
			"--delimiter", "\t",
			"--reverse",
		},
	}

	stdout, _, err := w.Base.ExecCommand(execCommand)
	if err != nil {
		return "", fmt.Errorf("fzf failed: %w", err)
	}

	return stdout, nil
}

// parseJumpOutput parses the fzf output into key and row
func parseJumpOutput(output string) (key, row string, err error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	lines := []string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	if len(lines) == 0 {
		return "", "", fmt.Errorf("no selection made")
	}

	if len(lines) == 1 {
		return "", lines[0], nil
	}

	return lines[0], lines[1], nil
}

// listNonWorktreeWindows returns tmux windows that are not worktree-owned
func (w *WorktreeManager) listNonWorktreeWindows() ([]string, error) {
	execCommand := cmd.CommandParams{
		Command: "tmux",
		Args:    []string{"list-windows", "-F", "#{window_name}"},
	}

	stdout, _, err := w.Base.ExecCommand(execCommand)
	if err != nil {
		return nil, err
	}

	var windows []string
	scanner := bufio.NewScanner(strings.NewReader(stdout))
	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if name != "" && !strings.HasPrefix(name, windowPrefix) {
			windows = append(windows, name)
		}
	}

	return windows, nil
}

// GetWindowName returns the tmux window name for a given worktree name
func GetWindowName(name string) string {
	return windowPrefix + name
}

// GetWorktreeDir returns the worktree directory name (deprecated, use worktreePath instead)
func GetWorktreeDir() string {
	return ".worktrees"
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
