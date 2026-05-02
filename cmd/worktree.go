/*
* Copyright © 2025 Carlos Mendez <carlos@hadaelectronics.com> | https://cjairm.me/
 */
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/internal/tooling/worktree"
	"github.com/cjairm/devgita/pkg/utils"
	"github.com/spf13/cobra"
)

var worktreeCmd = &cobra.Command{
	Use:     "worktree",
	Aliases: []string{"wt"},
	Short:   "Manage git worktrees with tmux windows",
	Long: `Manage git worktrees with tmux windows (alias: wt).

Each worktree gets its own tmux window in the current session with an AI assistant running,
enabling parallel AI-assisted development across multiple branches.

Worktrees are stored in ~/.local/share/devgita/worktrees/<repo-slug>/,
and tmux windows are prefixed with "wt-" for easy identification.

Examples:
  dg worktree create feature-login              # Create worktree + window with default AI
  dg worktree create feature-login --ai claude  # Create with Claude Code
  dg wt c feature-login                         # Same, using short form
  dg wt l                                       # List all worktrees
  dg wt j                                       # Jump to worktree (fzf selection)
  dg wt rm                                      # Remove worktree (fzf selection)
  dg wt repair feature-login                    # Repair missing window
  dg wt prune                                   # Remove all worktrees`,
}

var worktreeCreateCmd = &cobra.Command{
	Use:     "create <name>",
	Aliases: []string{"c", "new"},
	Short:   "Create a new worktree with tmux window",
	Long: `Create a new git worktree with an associated tmux window (aliases: c, new).

This command:
  1. Creates a new git worktree in ~/.local/share/devgita/worktrees/<repo>/<name>
  2. Creates a new branch with the same name
  3. Creates a new tmux window named wt-<name> in the current session
  4. Launches the selected AI coder in the window

AI coder selection precedence:
  1. --ai flag
  2. DEVGITA_WORKTREE_AI environment variable
  3. worktree.default_ai in global_config.yaml
  4. Default: opencode

Valid AI coders: opencode (oc), claude (cc, claudecode)

After creation, switch to the window with:
  <prefix> + [window number] or <prefix> + w to see all windows`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		aiAlias := resolveAIAlias(aiFlag, &globalConfig)

		coder, err := worktree.ResolveAICoder(aiAlias)
		if err != nil {
			utils.MaybeExitWithError(err)
		}

		wm := worktree.New()
		if err := wm.Create(name, coder); err != nil {
			utils.MaybeExitWithError(err)
		}

		repoRoot, _ := wm.Git.GetRepoRoot()
		repoSlug := repoRoot
		if repoRoot != "" {
			repoSlug = repoRoot[findLastSlash(repoRoot)+1:]
		}
		utils.PrintSuccess(fmt.Sprintf("Created worktree: %s/%s", repoSlug, name))
		utils.PrintSuccess(fmt.Sprintf("Created tmux window: %s", worktree.GetWindowName(name)))
		utils.PrintInfo("Switch to window with: <prefix> + w")
	},
}

var worktreeListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "List all worktrees with window status",
	Long: `List all git worktrees managed by devgita with their tmux window status (aliases: l, ls).

Shows worktrees from all repos in ~/.local/share/devgita/worktrees/ along with:
  - Repo name
  - Branch name
  - Associated tmux window name
  - Whether the window is currently active`,
	Run: func(cmd *cobra.Command, args []string) {
		wm := worktree.New()

		statuses, err := wm.List()
		utils.MaybeExitWithError(err)

		if len(statuses) == 0 {
			utils.PrintInfo("No worktrees found")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "REPO\tWORKTREE\tBRANCH\tWINDOW\tSTATUS")
		for _, s := range statuses {
			status := "No window"
			if s.WindowActive {
				status = "Active"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				s.Repo, s.Name, s.Branch, s.TmuxWindow, status)
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
  2. Removes the git worktree
  3. Deletes the branch (force delete with -D)

If no name is provided, opens an interactive fzf picker to select a worktree.

Use --force to remove even if the worktree has uncommitted changes.

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

		if err := wm.Remove(name, forceFlag); err != nil {
			utils.MaybeExitWithError(err)
		}

		utils.PrintSuccess(fmt.Sprintf("Removed worktree: %s", name))
	},
}

var worktreeJumpCmd = &cobra.Command{
	Use:     "jump",
	Aliases: []string{"j"},
	Short:   "Jump to a worktree's tmux window",
	Long: `Jump to a worktree's tmux window using fzf selection (alias: j).

Opens an interactive fzf picker showing all worktrees across all repos.

Inside tmux:
  - enter: jump to selected worktree window
  - ctrl-x: delete worktree (with confirmation)
  - ctrl-r: repair worktree (recreate window and launch AI)

Outside tmux:
  - enter: print worktree path to stdout
  - ctrl-x: delete worktree (with confirmation)
  - ctrl-r: repair worktree (warning: tmux not running)

The fzf dialog also shows non-worktree tmux windows when inside a tmux session.
Non-worktree windows can only be jumped to (enter); delete/repair are no-op.

Example:
  dg wt j    # Opens fzf picker, then switches to selected window`,
	Run: func(cmd *cobra.Command, args []string) {
		wm := worktree.New()
		aiAlias := resolveAIAlias("", &globalConfig)

		if err := wm.Jump(aiAlias); err != nil {
			utils.MaybeExitWithError(err)
		}
	},
}

var worktreeRepairCmd = &cobra.Command{
	Use:   "repair <name>",
	Short: "Repair a worktree by recreating its tmux window",
	Long: `Repair a worktree by recreating its tmux window and launching the AI coder.

This command:
  1. Checks that the worktree directory exists
  2. Creates a new tmux window if missing
  3. Launches the selected AI coder in the window

AI coder selection follows the same precedence as create:
  1. --ai flag
  2. DEVGITA_WORKTREE_AI environment variable
  3. worktree.default_ai in global_config.yaml
  4. Default: opencode`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		aiAlias := resolveAIAlias(aiFlag, &globalConfig)

		coder, err := worktree.ResolveAICoder(aiAlias)
		if err != nil {
			utils.MaybeExitWithError(err)
		}

		wm := worktree.New()
		if err := wm.Repair(name, coder); err != nil {
			utils.MaybeExitWithError(err)
		}

		utils.PrintSuccess(fmt.Sprintf("Repaired worktree: %s", name))
		utils.PrintSuccess(fmt.Sprintf("Launched AI coder in window: %s", worktree.GetWindowName(name)))
	},
}

var worktreePruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove all worktrees",
	Long: `Remove all worktrees managed by devgita.

This command prompts for confirmation before removing all worktrees
across all repos in ~/.local/share/devgita/worktrees/.

Each worktree is removed using the same logic as 'dg wt remove':
  - Kills the tmux window if present
  - Removes the git worktree
  - Deletes the branch

Example:
  dg wt prune    # Prompts for confirmation, then removes all worktrees`,
	Run: func(cmd *cobra.Command, args []string) {
		wm := worktree.New()

		if err := wm.Prune(); err != nil {
			utils.MaybeExitWithError(err)
		}

		utils.PrintSuccess("All worktrees removed")
	},
}

var (
	aiFlag    string
	forceFlag bool
)

func init() {
	rootCmd.AddCommand(worktreeCmd)
	worktreeCmd.AddCommand(worktreeCreateCmd)
	worktreeCmd.AddCommand(worktreeListCmd)
	worktreeCmd.AddCommand(worktreeRemoveCmd)
	worktreeCmd.AddCommand(worktreeJumpCmd)
	worktreeCmd.AddCommand(worktreeRepairCmd)
	worktreeCmd.AddCommand(worktreePruneCmd)

	worktreeCreateCmd.Flags().StringVarP(&aiFlag, "ai", "a", "", "AI coder to launch (opencode, oc, claude, cc, claudecode)")
	worktreeRepairCmd.Flags().StringVarP(&aiFlag, "ai", "a", "", "AI coder to launch (opencode, oc, claude, cc, claudecode)")
	worktreeRemoveCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force removal even if worktree has uncommitted changes")
}

var globalConfig config.GlobalConfig

func resolveAIAlias(flagValue string, gc *config.GlobalConfig) string {
	if flagValue != "" {
		return flagValue
	}
	if env := os.Getenv("DEVGITA_WORKTREE_AI"); env != "" {
		return env
	}
	if gc.Worktree.DefaultAI != "" {
		return gc.Worktree.DefaultAI
	}
	return "opencode"
}

func findLastSlash(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			return i
		}
	}
	return 0
}
