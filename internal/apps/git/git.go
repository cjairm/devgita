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

// WorktreeInfo contains information about a git worktree
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
	if _, _, err := g.Base.ExecCommand(execCommand); err != nil {
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
	if _, _, err := g.Base.ExecCommand(execCommand); err != nil {
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

// CreateWorktree creates a new worktree with a branch
// Handles three cases:
// 1. Local branch exists: checkout that branch
// 2. Remote branch exists: create tracking branch (after fetch)
// 3. Neither exists: create new branch from HEAD
func (g *Git) CreateWorktree(path, branch string) error {
	// Fetch latest remote refs to ensure we see recent branches
	if err := g.FetchOrigin(); err != nil {
		// Log but don't fail - user might be offline or not have a remote
		// Continue with local/cached refs
	}

	// Check if local branch already exists
	localExists, err := g.BranchExists(branch)
	if err != nil {
		return fmt.Errorf("failed to check if local branch exists: %w", err)
	}

	// If local branch exists, checkout that branch
	if localExists {
		return g.ExecuteCommand("worktree", "add", path, branch)
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

	// Neither local nor remote exists, create new branch from HEAD
	return g.ExecuteCommand("worktree", "add", path, "-b", branch)
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
			return fmt.Errorf("removed worktree but failed to delete branch '%s': %w", branchName, err)
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
	for _, line := range strings.Split(stdout, "\n") {
		if strings.HasPrefix(line, "worktree ") {
			return strings.TrimPrefix(line, "worktree "), nil
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
// Returns one warning string per offending hook file, or nil if all clear.
func (g *Git) CheckHookCompatibility(repoRoot string) []string {
	hooksDir := g.hooksDir(repoRoot)

	hookFiles := []string{"pre-commit", "commit-msg", "prepare-commit-msg", "post-commit", "pre-push"}
	incompatiblePatterns := []string{"[ -d .git", "test -d .git"}

	var warnings []string
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

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if current.Path != "" {
				worktrees = append(worktrees, current)
				current = WorktreeInfo{}
			}
			continue
		}

		if strings.HasPrefix(line, "worktree ") {
			current.Path = strings.TrimPrefix(line, "worktree ")
		} else if strings.HasPrefix(line, "HEAD ") {
			current.Commit = strings.TrimPrefix(line, "HEAD ")
		} else if strings.HasPrefix(line, "branch ") {
			branch := strings.TrimPrefix(line, "branch ")
			// Remove refs/heads/ prefix if present
			branch = strings.TrimPrefix(branch, "refs/heads/")
			current.Branch = branch
		}
	}

	// Don't forget the last worktree if output doesn't end with blank line
	if current.Path != "" {
		worktrees = append(worktrees, current)
	}

	return worktrees
}
