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
	Use:     "worktree",
	Aliases: []string{"wt"},
	Short:   "Manage git worktrees with tmux windows",
	Long: `Manage git worktrees with tmux windows (alias: wt).

Each worktree gets its own tmux window in the current session with OpenCode running,
enabling parallel AI-assisted development across multiple branches.

Worktrees are created in the .worktrees/ directory of your repository,
and tmux windows are prefixed with "wt-" for easy identification.

Examples:
  dg worktree create feature-login    # Create worktree + window
  dg wt c feature-login               # Same, using short form
  dg wt l                             # List all worktrees
  dg wt j                             # Jump to worktree (fzf selection)
  dg wt rm                            # Remove worktree (fzf selection)`,
}

var worktreeCreateCmd = &cobra.Command{
	Use:     "create <name>",
	Aliases: []string{"c", "new"},
	Short:   "Create a new worktree with tmux window",
	Long: `Create a new git worktree with an associated tmux window (aliases: c, new).

This command:
  1. Creates a new git worktree in .worktrees/<name>
  2. Creates a new branch with the same name
  3. Creates a new tmux window named wt-<name> in the current session
  4. Launches OpenCode in the window

After creation, switch to the window with:
  <prefix> + [window number] or <prefix> + w to see all windows`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		wm := worktree.New()

		if err := wm.Create(name); err != nil {
			utils.MaybeExitWithError(err)
		}

		utils.PrintSuccess(fmt.Sprintf("Created worktree: %s/%s", worktree.GetWorktreeDir(), name))
		utils.PrintSuccess(fmt.Sprintf("Created tmux window: %s", worktree.GetWindowName(name)))
		utils.PrintInfo("Switch to window with: <prefix> + w")
	},
}

var worktreeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all worktrees with window status",
	Long: `List all git worktrees managed by devgita with their tmux window status.

Shows worktrees in the .worktrees/ directory along with:
  - Branch name
  - Associated tmux window name
  - Whether the window is currently active`,
	Run: func(cmd *cobra.Command, args []string) {
		wm := worktree.New()

		statuses, err := wm.List()
		utils.MaybeExitWithError(err)

		if len(statuses) == 0 {
			utils.PrintInfo("No worktrees found in .worktrees/")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "WORKTREE\tBRANCH\tWINDOW\tSTATUS")
		for _, s := range statuses {
			status := "No window"
			if s.WindowActive {
				status = "Active"
			}
			fmt.Fprintf(w, "%s/%s\t%s\t%s\t%s\n",
				worktree.GetWorktreeDir(), s.Name, s.Branch, s.TmuxWindow, status)
		}
		w.Flush()
	},
}

var worktreeRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a worktree and its tmux window",
	Long: `Remove a git worktree and kill its associated tmux window.

This command:
  1. Kills the tmux window wt-<name> if it exists
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
