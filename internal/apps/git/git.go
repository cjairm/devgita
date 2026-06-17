// Git version control system with devgita integration
//
// Git is the distributed version control system that tracks changes in source code
// during software development. This module provides installation and configuration
// management for Git with devgita integration.
//
// References:
// - Git Documentation: https://git-scm.com/doc
// - Git Commands Reference: https://git-scm.com/docs
//
// Common Git commands available through ExecuteCommand():
//   - git status - Show working tree status
//   - git clone <url> <dir> - Clone repository
//   - git checkout <branch> - Switch branch
//   - git add . - Stage changes
//   - git commit -m "msg" - Commit changes
//   - git push origin <branch> - Push changes
//   - git pull origin <branch> - Pull changes
//   - git stash pop - Apply stashed changes
//   - git log - View commit history
//   - git clean -X -d -f - Clean ignored files

package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cjairm/devgita/internal/apps"
	"github.com/cjairm/devgita/internal/apps/baseapp"
	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

var _ apps.App = (*Git)(nil)

// WorktreeInfo contains information about a git worktree.
type WorktreeInfo struct {
	Path   string
	Branch string
	Commit string
}

type Git struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func (g *Git) Name() string       { return constants.Git }
func (g *Git) Kind() apps.AppKind { return apps.KindTerminal }

func New() *Git {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Git{Cmd: osCmd, Base: baseCmd}
}

func (g *Git) Install() error {
	return g.Cmd.InstallPackage(constants.Git)
}

func (g *Git) SoftInstall() error {
	return g.Cmd.MaybeInstallPackage(constants.Git)
}

func (g *Git) ForceInstall() error {
	return baseapp.Reinstall(g.Install, g.Uninstall)
}

func (g *Git) Uninstall() error {
	gc := &config.GlobalConfig{}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	if err := g.Cmd.UninstallPackage(constants.Git); err != nil {
		return fmt.Errorf("failed to uninstall git: %w", err)
	}
	_ = os.RemoveAll(paths.Paths.Config.Git)
	gc.RemoveFromInstalled(constants.Git, "package")
	return gc.Save()
}

func (g *Git) ForceConfigure() error {
	if err := files.CopyDir(paths.Paths.App.Configs.Git, paths.Paths.Config.Git); err != nil {
		return err
	}
	gc := &config.GlobalConfig{}
	if err := gc.Create(); err != nil {
		return fmt.Errorf("failed to create global config: %w", err)
	}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	gc.AddToInstalled(constants.Git, "package")
	return gc.Save()
}

func (g *Git) SoftConfigure() error {
	configFile := filepath.Join(paths.Paths.Config.Git, ".gitconfig")
	if files.FileAlreadyExist(configFile) {
		return nil
	}
	return files.CopyDir(paths.Paths.App.Configs.Git, paths.Paths.Config.Git)
}

func (g *Git) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		IsSudo:  false,
		Command: constants.Git,
		Args:    args,
	}
	if _, stderr, err := g.Base.ExecCommand(execCommand); err != nil {
		if stderr != "" {
			return fmt.Errorf("git: %s", stderr)
		}
		return fmt.Errorf("failed to run git command: %w", err)
	}
	return nil
}

// ExecuteCommandAt runs a git command with -C <dir> so it operates in the
// given directory regardless of the process's current working directory.
func (g *Git) ExecuteCommandAt(dir string, args ...string) error {
	fullArgs := append([]string{"-C", dir}, args...)
	execCommand := cmd.CommandParams{
		IsSudo:  false,
		Command: constants.Git,
		Args:    fullArgs,
	}
	if _, stderr, err := g.Base.ExecCommand(execCommand); err != nil {
		if stderr != "" {
			return fmt.Errorf("git: %s", stderr)
		}
		return fmt.Errorf("failed to run git command at %s: %w", dir, err)
	}
	return nil
}

func (g *Git) Clone(url, dstPath string) error {
	return g.ExecuteCommand("clone", url, dstPath)
}

func (g *Git) DeleteBranch(branch string, isForced bool) error {
	deleteArg := "-d"
	if isForced {
		deleteArg = "-D"
	}
	return g.ExecuteCommand("branch", deleteArg, branch)
}

func (g *Git) DeepClean(url, dstPath string) error {
	// -X: This option tells Git to remove only the files that are ignored by Git (i.e., files that are listed in your .gitignore file). It does not remove untracked files that are not ignored.
	// -d: This option allows Git to remove untracked directories in addition to untracked files.
	// -f: This option stands for "force." Git requires this option to actually perform the clean operation, as a safety measure to prevent accidental data loss.
	return g.ExecuteCommand("clean", "-X", "-d", "-f")
}

func (g *Git) FetchOrigin() error {
	return g.ExecuteCommand("fetch", "origin")
}

func (g *Git) Pop(branch string) error {
	return g.ExecuteCommand("stash", "pop")
}

func (g *Git) Pull(branch string) error {
	if branch == "" {
		return g.ExecuteCommand("pull")
	}
	return g.ExecuteCommand("pull", "origin", branch)
}

func (g *Git) SwitchBranch(branch string) error {
	return g.ExecuteCommand("checkout", branch)
}

func (g *Git) Restore(branch, files string) error {
	if branch == "" {
		branch = "main"
	}
	return g.ExecuteCommand("restore", "--source", branch, "--", files)
}

func (g *Git) Update() error {
	return fmt.Errorf("%w for git", apps.ErrUpdateNotSupported)
}

// BranchExists checks if a branch exists in the repository
func (g *Git) BranchExists(branch string) (bool, error) {
	execCommand := cmd.CommandParams{
		Command: constants.Git,
		Args:    []string{"branch", "--list", branch},
	}
	stdout, _, err := g.Base.ExecCommand(execCommand)
	if err != nil {
		return false, err
	}
	// If output contains the branch name, it exists
	return strings.TrimSpace(stdout) != "", nil
}

// RemoteBranchExists checks if a remote branch exists (e.g., origin/feature-A)
func (g *Git) RemoteBranchExists(branch string) (bool, error) {
	execCommand := cmd.CommandParams{
		Command: constants.Git,
		Args:    []string{"branch", "-r", "--list", fmt.Sprintf("origin/%s", branch)},
	}
	stdout, _, err := g.Base.ExecCommand(execCommand)
	if err != nil {
		return false, err
	}
	// If output contains the remote branch name, it exists
	return strings.TrimSpace(stdout) != "", nil
}

// DefaultBranch returns the repository's default branch name (e.g. "main").
// It resolves origin/HEAD when available and falls back to "main" so callers
// always get a usable branch name even on repos where origin/HEAD is unset.
func (g *Git) DefaultBranch() string {
	execCommand := cmd.CommandParams{
		Command: constants.Git,
		Args:    []string{"symbolic-ref", "--short", "refs/remotes/origin/HEAD"},
	}
	stdout, _, err := g.Base.ExecCommand(execCommand)
	if err == nil {
		// Output looks like "origin/main"; strip the remote prefix.
		ref := strings.TrimSpace(stdout)
		if i := strings.LastIndex(ref, "/"); i != -1 {
			ref = ref[i+1:]
		}
		if ref != "" {
			return ref
		}
	}
	return "main"
}

// CreateWorktree creates a new worktree with a branch
// Handles three cases:
// 1. Local branch exists: checkout that branch
// 2. Remote branch exists: create tracking branch (after fetch)
// 3. Neither exists: create new branch from the freshly-fetched default branch
func (g *Git) CreateWorktree(path, branch string) error {
	// Fetch latest remote refs to ensure we see recent branches.
	// Best-effort: ignore errors (user may be offline or have no remote).
	_ = g.FetchOrigin()

	// Check if local branch already exists
	localExists, err := g.BranchExists(branch)
	if err != nil {
		return fmt.Errorf("failed to check if local branch exists: %w", err)
	}

	// If local branch exists, check it out in the worktree, then bring it up to
	// date with its remote counterpart.
	if localExists {
		if err := g.ExecuteCommand("worktree", "add", path, branch); err != nil {
			return err
		}
		return g.syncExistingBranch(path, branch)
	}

	// Check if remote branch exists
	remoteExists, err := g.RemoteBranchExists(branch)
	if err != nil {
		return fmt.Errorf("failed to check if remote branch exists: %w", err)
	}

	// If remote branch exists, checkout it (creates tracking branch automatically)
	if remoteExists {
		return g.ExecuteCommand("worktree", "add", path, branch)
	}

	// Neither local nor remote exists: create a new branch. Prefer basing it on
	// the freshly-fetched default branch (origin/<default>) so new worktrees are
	// deterministic and never inherit a stale or unrelated HEAD. Fall back to
	// HEAD when the remote default isn't available (e.g. offline, no origin).
	defaultBranch := g.DefaultBranch()
	baseExists, err := g.RemoteBranchExists(defaultBranch)
	if err != nil {
		return fmt.Errorf("failed to check if remote default branch exists: %w", err)
	}
	if baseExists {
		base := fmt.Sprintf("origin/%s", defaultBranch)
		return g.ExecuteCommand("worktree", "add", path, "-b", branch, base)
	}
	return g.ExecuteCommand("worktree", "add", path, "-b", branch)
}

// syncExistingBranch brings a worktree's already-existing local branch up to
// date with its remote counterpart. It only fast-forwards, so unpushed local
// commits are never discarded. When there is no remote counterpart there is
// nothing to sync. When histories have diverged the fast-forward fails; we
// leave the branch untouched and warn the user how to reconcile manually
// (logger is suppressed below ERROR in normal runs, so we print directly).
func (g *Git) syncExistingBranch(path, branch string) error {
	remoteExists, err := g.RemoteBranchExists(branch)
	if err != nil {
		return fmt.Errorf("failed to check if remote branch exists: %w", err)
	}
	if !remoteExists {
		return nil
	}

	base := fmt.Sprintf("origin/%s", branch)
	if ffErr := g.ExecuteCommandAt(path, "merge", "--ff-only", base); ffErr != nil {
		fmt.Printf(
			"Warning: local branch %q diverged from %s and was not updated.\n"+
				"  The worktree was created at the branch's current local state.\n"+
				"  If the local branch has no unique work, sync it with:\n"+
				"    git -C %s fetch origin %s\n"+
				"    git -C %s reset --hard %s\n",
			branch, base, path, branch, path, base,
		)
	}
	return nil
}

// ListWorktrees returns parsed worktree information
func (g *Git) ListWorktrees() ([]WorktreeInfo, error) {
	execCommand := cmd.CommandParams{
		Command: constants.Git,
		Args:    []string{"worktree", "list", "--porcelain"},
	}
	stdout, _, err := g.Base.ExecCommand(execCommand)
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %w", err)
	}
	return parseWorktreeOutput(stdout), nil
}

// ListWorktreesAt lists worktrees for the git repository at the given directory.
// This avoids depending on the current working directory.
func (g *Git) ListWorktreesAt(dir string) ([]WorktreeInfo, error) {
	execCommand := cmd.CommandParams{
		Command: constants.Git,
		Args:    []string{"-C", dir, "worktree", "list", "--porcelain"},
	}
	stdout, _, err := g.Base.ExecCommand(execCommand)
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees at %s: %w", dir, err)
	}
	return parseWorktreeOutput(stdout), nil
}

// RemoveWorktree removes a worktree and optionally its associated branch.
// Resolves the main worktree first so the remove command doesn't run from
// within the worktree being deleted.
func (g *Git) RemoveWorktree(path string, deleteBranch bool, branchName string) error {
	// Find the main worktree by resolving git-common-dir from the target path
	mainWorktree, err := g.getMainWorktree(path)
	if err != nil {
		return fmt.Errorf("cannot resolve main worktree for %s: %w", path, err)
	}

	// Remove the worktree from the main worktree context
	if err := g.ExecuteCommandAt(mainWorktree, "worktree", "remove", path); err != nil {
		return err
	}

	// Delete the branch if requested
	if deleteBranch && branchName != "" {
		if err := g.ExecuteCommandAt(mainWorktree, "branch", "-D", branchName); err != nil {
			return fmt.Errorf(
				"removed worktree but failed to delete branch '%s': %w",
				branchName,
				err,
			)
		}
	}

	return nil
}

// getMainWorktree resolves the main worktree path from any worktree in the repo.
func (g *Git) getMainWorktree(fromPath string) (string, error) {
	execCommand := cmd.CommandParams{
		Command: constants.Git,
		Args:    []string{"-C", fromPath, "worktree", "list", "--porcelain"},
	}
	stdout, _, err := g.Base.ExecCommand(execCommand)
	if err != nil {
		return "", err
	}
	// First "worktree <path>" line is always the main worktree
	for line := range strings.SplitSeq(stdout, "\n") {
		if path, ok := strings.CutPrefix(line, "worktree "); ok {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find main worktree")
}

// GetRepoRoot returns the root directory of the current git repository
func (g *Git) GetRepoRoot() (string, error) {
	execCommand := cmd.CommandParams{
		Command: constants.Git,
		Args:    []string{"rev-parse", "--show-toplevel"},
	}
	stdout, _, err := g.Base.ExecCommand(execCommand)
	if err != nil {
		return "", fmt.Errorf("failed to get repo root: %w", err)
	}
	return strings.TrimSpace(stdout), nil
}

// IsWorktreeDirty checks if a worktree has uncommitted changes
func (g *Git) IsWorktreeDirty(path string) (bool, error) {
	execCommand := cmd.CommandParams{
		Command: constants.Git,
		Args:    []string{"-C", path, "status", "--porcelain"},
	}
	stdout, _, err := g.Base.ExecCommand(execCommand)
	if err != nil {
		return false, fmt.Errorf("failed to check worktree status: %w", err)
	}
	return strings.TrimSpace(stdout) != "", nil
}

// Diff returns the colored diff for a worktree (HEAD vs working tree + untracked files).
func (g *Git) Diff(path string) (string, error) {
	// Tracked changes: diff HEAD (covers both staged + unstaged)
	execCommand := cmd.CommandParams{
		Command: constants.Git,
		Args:    []string{"-C", path, "diff", "--color=always", "HEAD"},
	}
	stdout, _, err := g.Base.ExecCommand(execCommand)
	if err != nil {
		// Fallback for repos with no commits
		execCommand2 := cmd.CommandParams{
			Command: constants.Git,
			Args:    []string{"-C", path, "diff", "--color=always"},
		}
		stdout, _, err = g.Base.ExecCommand(execCommand2)
		if err != nil {
			return "", fmt.Errorf("git diff failed: %w", err)
		}
	}

	// Untracked files from status --porcelain
	statusCmd := cmd.CommandParams{
		Command: constants.Git,
		Args:    []string{"-C", path, "status", "--porcelain"},
	}
	statusOut, _, _ := g.Base.ExecCommand(statusCmd)
	var untracked []string
	for line := range strings.SplitSeq(statusOut, "\n") {
		if path, ok := strings.CutPrefix(line, "?? "); ok {
			untracked = append(untracked, path)
		}
	}
	if len(untracked) > 0 {
		stdout += "\nUntracked files:\n"
		for _, f := range untracked {
			stdout += "  " + f + "\n"
		}
	}
	return stdout, nil
}

// DiffStat returns aggregate stats: files changed, lines added, lines removed.
// Untracked files are counted in files but their line counts are 0.
func (g *Git) DiffStat(path string) (files, added, removed int, err error) {
	execCommand := cmd.CommandParams{
		Command: constants.Git,
		Args:    []string{"-C", path, "diff", "--numstat", "HEAD"},
	}
	stdout, _, execErr := g.Base.ExecCommand(execCommand)
	if execErr != nil {
		execCommand2 := cmd.CommandParams{
			Command: constants.Git,
			Args:    []string{"-C", path, "diff", "--numstat"},
		}
		stdout, _, execErr = g.Base.ExecCommand(execCommand2)
		if execErr != nil {
			return 0, 0, 0, fmt.Errorf("git diff --numstat failed: %w", execErr)
		}
	}
	for line := range strings.SplitSeq(strings.TrimSpace(stdout), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		files++
		if parts[0] != "-" {
			var a int
			_, _ = fmt.Sscanf(parts[0], "%d", &a)
			added += a
		}
		if parts[1] != "-" {
			var r int
			_, _ = fmt.Sscanf(parts[1], "%d", &r)
			removed += r
		}
	}

	// Count untracked files
	statusCmd := cmd.CommandParams{
		Command: constants.Git,
		Args:    []string{"-C", path, "status", "--porcelain"},
	}
	statusOut, _, _ := g.Base.ExecCommand(statusCmd)
	for line := range strings.SplitSeq(statusOut, "\n") {
		if _, ok := strings.CutPrefix(line, "?? "); ok {
			files++
		}
	}
	return files, added, removed, nil
}

// PruneWorktrees removes stale worktree entries
func (g *Git) PruneWorktrees() error {
	return g.ExecuteCommand("worktree", "prune")
}

// PruneWorktreesAt removes stale worktree entries, running from the given directory.
func (g *Git) PruneWorktreesAt(dir string) error {
	return g.ExecuteCommandAt(dir, "worktree", "prune")
}

// CheckHookCompatibility scans the repo's effective hooks directory for scripts
// that use `[ -d .git ]` or `test -d .git`. In a git worktree the .git entry is
// a FILE, not a directory, so those checks always fail and block git commit.
//
// Also detects Affiance hooks which have a known bug where the .git file regex
// fails to match due to trailing newlines, causing "no .git directory found".
// See: https://github.com/mariusbutuc/affiance/issues/XXX
//
// Returns one warning string per offending hook file, or nil if all clear.
func (g *Git) CheckHookCompatibility(repoRoot string) []string {
	hooksDir := g.hooksDir(repoRoot)

	hookFiles := []string{
		"pre-commit",
		"commit-msg",
		"prepare-commit-msg",
		"post-commit",
		"pre-push",
	}
	incompatiblePatterns := []string{"[ -d .git", "test -d .git"}

	var warnings []string

	// Check for Affiance hooks (known worktree incompatibility)
	affianceHook := filepath.Join(hooksDir, "affiance-hook")
	if _, err := os.Stat(affianceHook); err == nil {
		warnings = append(
			warnings,
			"affiance-hook (Affiance has a bug parsing .git files in worktrees; use --no-verify to bypass)",
		)
		// If Affiance is present, all hooks delegate to it, so skip individual checks
		return warnings
	}

	for _, hookFile := range hookFiles {
		content, err := os.ReadFile(filepath.Join(hooksDir, hookFile))
		if err != nil {
			continue
		}
		contentStr := string(content)
		for _, pattern := range incompatiblePatterns {
			if strings.Contains(contentStr, pattern) {
				warnings = append(warnings, fmt.Sprintf("%s (contains %q)", hookFile, pattern))
				break
			}
		}
	}
	return warnings
}

// hooksDir returns the effective hooks directory for repoRoot.
// Respects core.hooksPath if configured; falls back to <repoRoot>/.git/hooks.
func (g *Git) hooksDir(repoRoot string) string {
	execCommand := cmd.CommandParams{
		Command: constants.Git,
		Args:    []string{"-C", repoRoot, "config", "--get", "core.hooksPath"},
	}
	if stdout, _, err := g.Base.ExecCommand(execCommand); err == nil {
		p := strings.TrimSpace(stdout)
		if p != "" {
			if !filepath.IsAbs(p) {
				p = filepath.Join(repoRoot, p)
			}
			return p
		}
	}
	return filepath.Join(repoRoot, ".git", "hooks")
}

// parseWorktreeOutput parses the porcelain output of git worktree list
func parseWorktreeOutput(output string) []WorktreeInfo {
	var worktrees []WorktreeInfo
	var current WorktreeInfo

	for line := range strings.SplitSeq(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			if current.Path != "" {
				worktrees = append(worktrees, current)
				current = WorktreeInfo{}
			}
			continue
		}

		switch {
		case strings.HasPrefix(line, "worktree "):
			current.Path, _ = strings.CutPrefix(line, "worktree ")
		case strings.HasPrefix(line, "HEAD "):
			current.Commit, _ = strings.CutPrefix(line, "HEAD ")
		case strings.HasPrefix(line, "branch "):
			branchRef, _ := strings.CutPrefix(line, "branch ")
			current.Branch = strings.TrimPrefix(branchRef, "refs/heads/")
		}
	}

	// Don't forget the last worktree if output doesn't end with blank line
	if current.Path != "" {
		worktrees = append(worktrees, current)
	}

	return worktrees
}
