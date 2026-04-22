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
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "List all worktrees with window status",
	Long: `List all git worktrees managed by devgita with their tmux window status (aliases: l, ls).

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
	Use:     "remove [name]",
	Aliases: []string{"rm", "r"},
	Short:   "Remove a worktree and its tmux window",
	Long: `Remove a git worktree and kill its associated tmux window (aliases: rm, r).

This command:
  1. Kills the tmux window wt-<name> if it exists
  2. Removes the git worktree from .worktrees/<name>
  3. Deletes the branch (if not merged, use git branch -D manually)

If no name is provided, opens an interactive fzf picker to select a worktree.

Warning: Any uncommitted changes in the worktree will be lost.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		wm := worktree.New()
		var name string

		if len(args) == 0 {
			selected, err := wm.SelectWorktreeInteractively("Select worktree to remove:")
			if err != nil {
				utils.MaybeExitWithError(err)
			}
			name = selected
		} else {
			name = args[0]
		}

		if err := wm.Remove(name); err != nil {
			utils.MaybeExitWithError(err)
		}

		utils.PrintSuccess(fmt.Sprintf("Removed worktree: %s/%s", worktree.GetWorktreeDir(), name))
	},
}

var worktreeJumpCmd = &cobra.Command{
	Use:     "jump",
	Aliases: []string{"j"},
	Short:   "Jump to a worktree's tmux window",
	Long: `Jump to a worktree's tmux window using fzf selection (alias: j).

Opens an interactive fzf picker showing all available worktrees.
After selection, switches the current tmux session to that worktree's window.

Requires:
  - Running inside a tmux session
  - fzf installed

Example:
  dg wt j    # Opens fzf picker, then switches to selected window`,
	Run: func(cmd *cobra.Command, args []string) {
		wm := worktree.New()

		selected, err := wm.SelectWorktreeInteractively("Select worktree to jump to:")
		if err != nil {
			utils.MaybeExitWithError(err)
		}

		windowName := worktree.GetWindowName(selected)

		// Check if window exists
		if !wm.Tmux.HasWindow(windowName) {
			utils.PrintError(fmt.Sprintf("Window '%s' not found. The worktree exists but has no active window.", windowName))
			utils.PrintInfo("Use 'dg wt create' to recreate the window, or manually create it.")
			return
		}

		if err := wm.Tmux.SelectWindow(windowName); err != nil {
			utils.MaybeExitWithError(fmt.Errorf("failed to switch to window: %w", err))
		}

		utils.PrintSuccess(fmt.Sprintf("Switched to window: %s", windowName))
	},
}

func init() {
	rootCmd.AddCommand(worktreeCmd)
	worktreeCmd.AddCommand(worktreeCreateCmd)
	worktreeCmd.AddCommand(worktreeListCmd)
	worktreeCmd.AddCommand(worktreeRemoveCmd)
	worktreeCmd.AddCommand(worktreeJumpCmd)
}
