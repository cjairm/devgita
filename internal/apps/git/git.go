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
	"path/filepath"
	"strings"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

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
	err := g.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall git: %w", err)
	}
	return g.Install()
}

func (g *Git) Uninstall() error {
	return fmt.Errorf("git uninstall not supported through devgita")
}

func (g *Git) ForceConfigure() error {
	return files.CopyDir(paths.Paths.App.Configs.Git, paths.Paths.Config.Git)
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
	return fmt.Errorf("git update not implemented through devgita")
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

// RemoveWorktree removes a worktree and optionally its associated branch
func (g *Git) RemoveWorktree(path string, deleteBranch bool, branchName string) error {
	// Remove the worktree directory
	if err := g.ExecuteCommand("worktree", "remove", path); err != nil {
		return err
	}

	// Delete the branch if requested
	if deleteBranch && branchName != "" {
		// Use force delete (-D) to avoid issues with unmerged changes
		if err := g.DeleteBranch(branchName, true); err != nil {
			// Log but don't fail - branch might not exist or might be current branch
			return fmt.Errorf("removed worktree but failed to delete branch '%s': %w", branchName, err)
		}
	}

	return nil
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
