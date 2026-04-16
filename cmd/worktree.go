/*
* Copyright © 2025 Carlos Mendez <carlos@hadaelectronics.com> | https://cjairm.me/
 */
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/cjairm/devgita/internal/tooling/worktree"
	"github.com/cjairm/devgita/pkg/utils"
	"github.com/spf13/cobra"
)

var worktreeCmd = &cobra.Command{
	Use:   "worktree",
	Short: "Manage git worktrees with tmux sessions",
	Long: `Create and manage git worktrees with isolated tmux sessions.

Each worktree gets its own tmux session with OpenCode running,
enabling parallel AI-assisted development across multiple branches.

Worktrees are created in the .worktrees/ directory of your repository,
and tmux sessions are prefixed with "wt-" for easy identification.

Examples:
  dg worktree create feature-login    # Create worktree + session
  dg worktree list                    # List all worktrees
  dg worktree remove feature-login    # Remove worktree + session`,
}

var worktreeCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new worktree with tmux session",
	Long: `Create a new git worktree with an associated tmux session.

This command:
  1. Creates a new git worktree in .worktrees/<name>
  2. Creates a new branch with the same name
  3. Creates a detached tmux session named wt-<name>
  4. Launches OpenCode in the session

After creation, attach to the session with:
  tmux attach -t wt-<name>`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		wm := worktree.New()

		if err := wm.Create(name); err != nil {
			utils.MaybeExitWithError(err)
		}

		utils.PrintSuccess(fmt.Sprintf("Created worktree: %s/%s", worktree.GetWorktreeDir(), name))
		utils.PrintSuccess(fmt.Sprintf("Created tmux session: %s", worktree.GetSessionName(name)))
		utils.PrintInfo(fmt.Sprintf("Attach with: tmux attach -t %s", worktree.GetSessionName(name)))
	},
}

var worktreeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all worktrees with session status",
	Long: `List all git worktrees managed by devgita with their tmux session status.

Shows worktrees in the .worktrees/ directory along with:
  - Branch name
  - Associated tmux session name
  - Whether the session is currently active`,
	Run: func(cmd *cobra.Command, args []string) {
		wm := worktree.New()

		statuses, err := wm.List()
		utils.MaybeExitWithError(err)

		if len(statuses) == 0 {
			utils.PrintInfo("No worktrees found in .worktrees/")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "WORKTREE\tBRANCH\tSESSION\tSTATUS")
		for _, s := range statuses {
			status := "No session"
			if s.SessionActive {
				status = "Active"
			}
			fmt.Fprintf(w, "%s/%s\t%s\t%s\t%s\n",
				worktree.GetWorktreeDir(), s.Name, s.Branch, s.TmuxSession, status)
		}
		w.Flush()
	},
}

var worktreeRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a worktree and its tmux session",
	Long: `Remove a git worktree and kill its associated tmux session.

This command:
  1. Kills the tmux session wt-<name> if it exists
  2. Removes the git worktree from .worktrees/<name>
  3. Deletes the branch (if not merged, use git branch -D manually)

Warning: Any uncommitted changes in the worktree will be lost.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		wm := worktree.New()

		if err := wm.Remove(name); err != nil {
			utils.MaybeExitWithError(err)
		}

		utils.PrintSuccess(fmt.Sprintf("Removed worktree: %s/%s", worktree.GetWorktreeDir(), name))
	},
}

func init() {
	rootCmd.AddCommand(worktreeCmd)
	worktreeCmd.AddCommand(worktreeCreateCmd)
	worktreeCmd.AddCommand(worktreeListCmd)
	worktreeCmd.AddCommand(worktreeRemoveCmd)
}
