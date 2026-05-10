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
	"os/exec"
	"path/filepath"
	"regexp"
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

type pendingDeleteInfo struct {
	repoSlug string
	name     string
}

// WorktreeManager coordinates git worktrees with tmux windows
type WorktreeManager struct {
	Git           *git.Git
	Tmux          *tmux.Tmux
	Fzf           *fzf.Fzf
	Base          cmd.BaseCommandExecutor
	pendingDelete *pendingDeleteInfo
	fzfRun        func(rows []string, header string) (string, error)
}

// New creates a new WorktreeManager instance
func New() *WorktreeManager {
	wm := &WorktreeManager{
		Git:           git.New(),
		Tmux:          tmux.New(),
		Fzf:           fzf.New(),
		Base:          cmd.NewBaseCommand(),
		pendingDelete: nil,
	}
	wm.fzfRun = wm.execFzf
	return wm
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
		// Check if directory actually exists on disk
		if _, err := os.Stat(wtPath); os.IsNotExist(err) {
			// Directory missing but git still tracks it - auto-prune and continue
			if pruneErr := w.Git.PruneWorktrees(); pruneErr != nil {
				return fmt.Errorf("stale worktree entry detected but failed to prune: %w", pruneErr)
			}
			// After pruning, continue with creation
		} else {
			// Directory exists, suggest repair
			return fmt.Errorf("worktree '%s' exists but has no active window; use `dg wt repair %s`", name, name)
		}
	}
	if !state.WtExists && state.WindowExists {
		return fmt.Errorf("orphan window '%s' exists; run `tmux kill-window -t %s` manually", windowName, windowName)
	}

	if err := os.MkdirAll(filepath.Dir(wtPath), 0755); err != nil {
		return fmt.Errorf("failed to create worktree directory: %w", err)
	}

	if err := w.Git.CreateWorktree(wtPath, name); err != nil {
		if strings.Contains(err.Error(), "is a missing but already registered") {
			if pruneErr := w.Git.PruneWorktrees(); pruneErr == nil {
				if retryErr := w.Git.CreateWorktree(wtPath, name); retryErr == nil {
					return w.createWindowAndLaunch(windowName, wtPath, coder)
				}
			}
		}
		return fmt.Errorf("failed to create worktree: %w", err)
	}

	return w.createWindowAndLaunch(windowName, wtPath, coder)
}

func (w *WorktreeManager) createWindowAndLaunch(windowName, wtPath string, coder AICoder) error {
	if err := w.Tmux.CreateWindow(windowName, wtPath); err != nil {
		_ = w.Git.RemoveWorktree(wtPath, true, "")
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

	worktrees, err := w.Git.ListWorktrees()
	if err == nil {
		for _, wt := range worktrees {
			if wt.Path == state.WtPath {
				state.WtExists = true
				break
			}
		}
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

// ListNames returns just the worktree names across all repos for shell completion.
func (w *WorktreeManager) ListNames() ([]string, error) {
	statuses, err := w.List()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(statuses))
	for _, s := range statuses {
		names = append(names, s.Name)
	}
	return names, nil
}

// findRepoForWorktree searches the centralized base path for a worktree by name
// and returns the repo slug that owns it. Returns "" if not found or ambiguous.
func (w *WorktreeManager) findRepoForWorktree(name string) string {
	entries, err := os.ReadDir(GetWorktreeBasePath())
	if err != nil {
		return ""
	}
	var matches []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(GetWorktreeBasePath(), e.Name(), name)); err == nil {
			matches = append(matches, e.Name())
		}
	}
	if len(matches) == 1 {
		return matches[0]
	}
	return ""
}

// Remove removes a worktree and its tmux window.
// When the named worktree is not in the current repo, it searches the
// centralized base path so cross-repo removal works from any directory.
func (w *WorktreeManager) Remove(name string, force bool) error {
	repoRoot, err := w.Git.GetRepoRoot()
	if err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	repoSlug := filepath.Base(repoRoot)
	state, err := w.worktreeState(repoSlug, name)
	if err != nil {
		return err
	}

	// Not found in current repo — search centralized base path.
	if !state.WtExists && !state.WindowExists {
		if slug := w.findRepoForWorktree(name); slug != "" {
			repoSlug = slug
			state, err = w.worktreeState(repoSlug, name)
			if err != nil {
				return err
			}
		}
	}

	wtPath := w.worktreePath(repoSlug, name)
	windowName := windowPrefix + name

	if !state.WtExists && !state.WindowExists {
		return fmt.Errorf("nothing to remove for worktree '%s'", name)
	}

	if state.WtExists && !force {
		dirty, err := w.Git.IsWorktreeDirty(wtPath)
		// If the dirty check errors (e.g. stale/broken worktree not registered
		// with git), allow removal rather than blocking on an unverifiable state.
		if err == nil && dirty {
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
			// Fallback for stale worktrees: directory exists on disk but git
			// doesn't know about it. Remove the directory directly then prune.
			if rmErr := os.RemoveAll(wtPath); rmErr != nil {
				return fmt.Errorf("failed to remove worktree: %w", err)
			}
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

	state, err := w.worktreeState(repoSlug, name)
	if err != nil {
		return err
	}

	// If worktree doesn't exist at all (neither on disk nor in git), error out
	if !state.WtExists {
		return fmt.Errorf("no worktree '%s' to repair", name)
	}

	// If directory doesn't exist on disk but git knows about it, prune and error
	if _, err := os.Stat(wtPath); os.IsNotExist(err) {
		// Prune stale worktree entries
		if pruneErr := w.Git.PruneWorktrees(); pruneErr != nil {
			return fmt.Errorf("worktree '%s' directory missing and failed to prune: %w", name, pruneErr)
		}
		return fmt.Errorf("worktree '%s' directory was missing; pruned stale entry. Run `dg wt new %s` to recreate", name, name)
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
	if !confirmFromTTY() {
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

// removeByRepo removes a worktree by repo slug and name.
// Mirrors the same tolerant logic as Remove.
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

	if state.WtExists && !force {
		dirty, err := w.Git.IsWorktreeDirty(wtPath)
		if err == nil && dirty {
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
			if rmErr := os.RemoveAll(wtPath); rmErr != nil {
				return fmt.Errorf("failed to remove worktree: %w", err)
			}
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

// confirmAndRemove implements double-confirm delete pattern (like opencode).
// First ctrl-d highlights the row red (inline) and moves it to the top so the
// cursor lands on it; second ctrl-d on the same row actually deletes.
func (w *WorktreeManager) confirmAndRemove(rows []string, repoSlug, name string) error {
	// If there's already a pending delete for this worktree, execute it
	if w.pendingDelete != nil && w.pendingDelete.repoSlug == repoSlug && w.pendingDelete.name == name {
		w.pendingDelete = nil
		return w.removeByRepo(repoSlug, name, false)
	}

	// Set this worktree as pending delete
	w.pendingDelete = &pendingDeleteInfo{repoSlug: repoSlug, name: name}

	output, err := w.runFzfWithExpect(buildConfirmRows(rows, repoSlug, name))
	if err != nil {
		w.pendingDelete = nil // Clear pending on error/cancel
		if err.Error() == "selection cancelled" {
			return nil
		}
		return err
	}

	key, row, err := parseJumpOutput(output)
	if err != nil {
		w.pendingDelete = nil
		return err
	}

	newRepoSlug, newName, err := parseRepoAndName(row)
	if err != nil {
		return err
	}

	if key == "ctrl-d" && newRepoSlug == repoSlug && newName == name {
		w.pendingDelete = nil
		return w.removeByRepo(newRepoSlug, newName, false)
	}

	// If user pressed ctrl-d on a different worktree, start over with that one
	if key == "ctrl-d" {
		// Recursively handle the new pending delete
		w.pendingDelete = nil
		return w.confirmAndRemove(rows, newRepoSlug, newName)
	}

	// If user pressed ctrl-r, repair the worktree
	if key == "ctrl-r" {
		w.pendingDelete = nil
		coder, err := ResolveAICoder("")
		if err != nil {
			return err
		}
		return w.Repair(newName, coder)
	}

	// Any other action (enter, etc.) cancels the pending delete
	w.pendingDelete = nil
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

		repoSlug, name, err := parseRepoAndName(row)
		if err != nil {
			return err
		}
		windowName := windowPrefix + name

		switch key {
		case "":
			if w.Tmux.HasWindow(windowName) {
				return w.Tmux.SelectWindow(windowName)
			}
			fmt.Printf("Window '%s' not found. Use 'dg wt repair %s' to recreate it.\n", windowName, name)
			return nil
		case "ctrl-d":
			return w.confirmAndRemove(rows, repoSlug, name)
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

	repoSlug, name, err := parseRepoAndName(row)
	if err != nil {
		return err
	}

	switch key {
	case "":
		fmt.Println(w.worktreePath(repoSlug, name))
		return nil
	case "ctrl-d":
		return w.confirmAndRemove(rows, repoSlug, name)
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

// parseRepoAndName extracts repoSlug and worktree name from a jump row.
// Row format is "repo/name\tbranch\tstatus"; parts[0] is the combined "repo/name".
func parseRepoAndName(row string) (repoSlug, name string, err error) {
	parts := parseJumpRow(row)
	if len(parts) < 1 {
		return "", "", fmt.Errorf("invalid selection")
	}
	combined := strings.SplitN(parts[0], "/", 2)
	if len(combined) < 2 {
		return "", "", fmt.Errorf("invalid selection: expected repo/name format, got %q", parts[0])
	}
	return combined[0], combined[1], nil
}

// confirmFromTTY reads a y/n answer directly from /dev/tty so it works even
// after fzf has consumed the process stdin (e.g. inside a tmux display-popup).
// In tmux popup mode, we skip confirmation since the user already pressed ctrl-d.
func confirmFromTTY() bool {
	if os.Getenv("TMUX") != "" {
		return true
	}
	tty, err := os.Open("/dev/tty")
	if err != nil {
		var response string
		fmt.Scanln(&response) //nolint:errcheck
		return strings.ToLower(strings.TrimSpace(response)) == "y"
	}
	defer tty.Close()
	var response string
	fmt.Fscan(tty, &response)
	return strings.ToLower(strings.TrimSpace(response)) == "y"
}

// runFzfWithExpect delegates to w.fzfRun so tests can inject a mock.
func (w *WorktreeManager) runFzfWithExpect(rows []string, header ...string) (string, error) {
	h := "enter: jump | ctrl-d: delete | ctrl-r: repair"
	if len(header) > 0 && header[0] != "" {
		h = header[0]
	}
	return w.fzfRun(rows, h)
}

// execFzf is the real fzf execution assigned to fzfRun in New().
// fzf renders its UI to /dev/tty and writes the selection to stdout,
// which Output() captures.
func (w *WorktreeManager) execFzf(rows []string, header string) (string, error) {
	fzfCmd := exec.Command("fzf",
		"--height=60%",
		"--reverse",
		"--ansi",
		"--header", header,
		"--expect", "ctrl-d,ctrl-r",
		"--with-nth", "1,2,3",
		"--delimiter", "\t",
	)
	fzfCmd.Stdin = strings.NewReader(strings.Join(rows, "\n"))

	output, err := fzfCmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
			return "", fmt.Errorf("selection cancelled")
		}
		return "", fmt.Errorf("fzf failed: %w", err)
	}

	return string(output), nil
}

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(s string) string {
	return ansiEscape.ReplaceAllString(s, "")
}

// buildConfirmRows returns display rows for the delete-confirmation step.
// The pending item's name column is highlighted red and placed first so that
// fzf's cursor lands on it when it reopens (--reverse shows item 0 at top).
func buildConfirmRows(rows []string, repoSlug, name string) []string {
	pendingKey := repoSlug + "/" + name
	displayRows := make([]string, 0, len(rows))
	var pendingDisplayRow string
	for _, row := range rows {
		parts := strings.SplitN(row, "\t", 2)
		if parts[0] == pendingKey {
			rest := ""
			if len(parts) == 2 {
				rest = "\t" + parts[1]
			}
			pendingDisplayRow = "\033[41m" + parts[0] + "\033[0m" + rest
		} else {
			displayRows = append(displayRows, row)
		}
	}
	if pendingDisplayRow != "" {
		displayRows = append([]string{pendingDisplayRow}, displayRows...)
	}
	return displayRows
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
		return "", stripANSI(lines[0]), nil
	}

	return lines[0], stripANSI(lines[1]), nil
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
