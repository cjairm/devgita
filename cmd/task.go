/*
 * Copyright © 2025 Carlos Mendez <carlos@hadaelectronics.com> | https://cjairm.me/
 */
package cmd

import (
	"github.com/cjairm/devgita/internal/tooling/task"
	"github.com/spf13/cobra"
)

// taskHelpFunc restores standard Cobra help for the task subtree. The root sets
// a branded help func (utils.PrompCustomHelp) that prints only Use+Long and is
// inherited by children — which hides subcommands and flags. Agents re-reading
// `dg task --help` or `dg task <sub> --help` need the full listing, so this
// renders the long/short description followed by the default usage block
// (Available Commands, Flags, Examples).
func taskHelpFunc(cmd *cobra.Command, args []string) {
	if cmd.Long != "" {
		cmd.Println(cmd.Long)
		cmd.Println()
	} else if cmd.Short != "" {
		cmd.Println(cmd.Short)
		cmd.Println()
	}
	cmd.Print(cmd.UsageString())
}

// taskRunner is the interface used by task subcommands, enabling injection in tests.
type taskRunner interface {
	RefreshBranch(target string) error
	ResetMainBranch() error
	ReinstallLibraries() error
	ReinstallLibrary(name string) error
	DeleteBranch(target string) error
}

// newTaskManager is the factory used by task subcommands; overridden in tests.
var newTaskManager = func() taskRunner { return task.New() }

var taskCmd = &cobra.Command{
	Use:     "task",
	Aliases: []string{"t"},
	Short:   "Developer utilities (git, npm, GitHub PRs) callable by agents and humans",
	Long: `Developer utility commands callable by agents (Claude Code, CI, any
non-interactive process) and humans (via the dge() shell wrapper or directly).

Three families:
  - git branch:  refresh-branch, reset-main-branch, delete-branch
  - npm deps:    reinstall-libraries, reinstall-library
  - GitHub PRs:  review-threads, resolve/unresolve/reply-thread, submit-review,
                 create-pr, update-pr-description, approve-pr, request-changes-pr,
                 comment-pr, merge-pr, pr-view, pr-checks, current-pr, current-repo

PR data commands return compact, LLM-oriented output (gh fetches, jq renders).
Run "dg task <subcommand> --help" for flags and examples.`,
	Example: `  dg task review-threads --state unresolved
  dg task pr-view
  dg task refresh-branch
  dg task reinstall-library lodash`,
}

var taskRefreshBranchCmd = &cobra.Command{
	Use:   "refresh-branch [target]",
	Short: "Checkout target branch, pull, return to previous branch, and merge",
	Long: `Checkout the target branch (default: main), pull latest changes from origin,
return to the previous branch (git checkout -), and merge target into it.

This is equivalent to the dge refresh-branch shell utility.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := ""
		if len(args) > 0 {
			target = args[0]
		}
		return newTaskManager().RefreshBranch(target)
	},
}

var taskResetMainBranchCmd = &cobra.Command{
	Use:   "reset-main-branch",
	Short: "Checkout main and hard-reset to origin/main",
	Long: `Checkout the main branch and reset it hard to origin/main, discarding
any local commits or changes on main.

This is equivalent to the dge reset-main-branch shell utility.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return newTaskManager().ResetMainBranch()
	},
}

var taskDeleteBranchCmd = &cobra.Command{
	Use:   "delete-branch [target]",
	Short: "Checkout target, pull, then pick a branch to force-delete via fzf",
	Long: `Checkout the target branch (default: main), fetch, and pull, then open an
interactive fzf picker to select a local branch to force-delete (git branch -D).

This is equivalent to the dge delete-branch shell utility.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := ""
		if len(args) > 0 {
			target = args[0]
		}
		return newTaskManager().DeleteBranch(target)
	},
}

var taskReinstallLibrariesCmd = &cobra.Command{
	Use:   "reinstall-libraries",
	Short: "Clean git-ignored files, remove node_modules, and run npm install",
	Long: `Run git clean -Xdf, remove node_modules/, run npm install, and remove
tsconfig.tsbuildinfo. Useful for fixing dependency issues in Node.js projects.

This is equivalent to the dge reinstall-libraries shell utility.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return newTaskManager().ReinstallLibraries()
	},
}

var taskReinstallLibraryCmd = &cobra.Command{
	Use:   "reinstall-library <name>",
	Short: "Remove a specific node_modules package and run npm install",
	Long: `Remove node_modules/<name> and re-run npm install. Useful for fixing
a single corrupted or mis-linked package without reinstalling everything.

This is equivalent to the dge reinstall-library shell utility.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return newTaskManager().ReinstallLibrary(args[0])
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)
	// Standard Cobra help for the whole task subtree (overrides the branded
	// root help func, which children would otherwise inherit and which hides
	// subcommands/flags). Children inherit this from taskCmd.
	taskCmd.SetHelpFunc(taskHelpFunc)
	taskCmd.AddCommand(taskRefreshBranchCmd)
	taskCmd.AddCommand(taskResetMainBranchCmd)
	taskCmd.AddCommand(taskDeleteBranchCmd)
	taskCmd.AddCommand(taskReinstallLibrariesCmd)
	taskCmd.AddCommand(taskReinstallLibraryCmd)
}
