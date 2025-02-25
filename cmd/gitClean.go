/*
Copyright © 2025 Carlos Mendez <carlos@hadaelectronics.com> | https://cjairm.me/
*/
package cmd

import (
	"fmt"
	"os"

	git "github.com/cjairm/devgita/internal/commands/git"
	"github.com/cjairm/devgita/pkg/utils"
	"github.com/spf13/cobra"
)

const GIT_CLEAN_CMD = "gc"

var gitCleanCmd = &cobra.Command{
	Use:   GIT_CLEAN_CMD,
	Short: "Clean up local Git branches",
	Long: fmt.Sprintf(`The git-clean command allows you to delete a specified local Git branch. 
You can choose to force the deletion of the branch, even if it has unmerged changes. 
By default, the command checks out the specified destination branch (default is 'main') 
and updates it with the latest changes from the remote repository before deleting the target branch.

Usage:
  devgita %s --branch-to-clean <branch-name> [--destination-branch <branch-name>] [--force-clean]
  devgita %s -c <branch-name> [-d <branch-name>] [-f]

Flags:
  -c, --branch-to-clean    string   The name of the branch to clean (required)
  -d, --destination-branch string   The name of the destination branch after cleaning (default "main")
  -f, --force-clean	   boolean  Force the deletion of the branch`, GIT_CLEAN_CMD, GIT_CLEAN_CMD),
	Run: runGitClean,
}

func init() {
	rootCmd.AddCommand(gitCleanCmd)
	gitCleanCmd.Flags().
		StringP("branch-to-clean", "c", "", "The name of the branch to clean (required)")
	gitCleanCmd.MarkFlagRequired("branch-to-clean")
	gitCleanCmd.Flags().
		StringP("destination-branch", "d", "", "The name of the destination branch after cleaning")
	gitCleanCmd.Flags().BoolP("force-clean", "f", false, "Force the deletion of the branch")
}

func runGitClean(cmd *cobra.Command, args []string) {
	dstBranch, err := cmd.Flags().GetString("destination-branch")
	utils.MaybeExitWithError(err)
	if dstBranch == "" {
		dstBranch = "main"
	}

	g := git.New()
	err = g.SwitchBranch(dstBranch)
	utils.MaybeExitWithError(err)

	err = g.FetchOrigin()
	utils.MaybeExitWithError(err)

	err = g.Pull(dstBranch)
	utils.MaybeExitWithError(err)

	branchToClean, err := cmd.Flags().GetString("branch-to-clean")
	utils.MaybeExitWithError(err)

	forceClean, err := cmd.Flags().GetBool("force-clean")
	utils.MaybeExitWithError(err)

	err = g.Clean(branchToClean, forceClean)
	utils.MaybeExitWithError(err)

	successMessage := fmt.Sprintf("Branch %s cleaned successfully", branchToClean)
	utils.PrintSuccess(successMessage)
	os.Exit(0)
}
