/*
Copyright © 2025 Carlos Mendez <carlos@hadaelectronics.com> | https://cjairm.me/
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/cjairm/devgita/pkg/common"
	"github.com/spf13/cobra"
)

const CMD_NAME = "gc"

var gitCleanCmd = &cobra.Command{
	Use:   CMD_NAME,
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
  -f, --force-clean	   boolean  Force the deletion of the branch`, CMD_NAME, CMD_NAME),
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
	gitCleanCmd.SetHelpFunc(common.PrompCustomHelp)
}

func runGitClean(cmd *cobra.Command, args []string) {
	dstBranch, err := cmd.Flags().GetString("destination-branch")
	if err != nil {
		common.PrintError(err)
		os.Exit(1)
	}
	if dstBranch == "" {
		dstBranch = "main"
	}
	err = common.GitCommand("checkout", dstBranch)
	if err != nil {
		common.PrintError(err)
		os.Exit(1)
	}
	err = common.GitCommand("fetch", "origin")
	if err != nil {
		common.PrintError(err)
		os.Exit(1)
	}
	err = common.GitCommand("pull", "origin", dstBranch)
	if err != nil {
		common.PrintError(err)
		os.Exit(1)
	}
	branchToClean, err := cmd.Flags().GetString("branch-to-clean")
	if err != nil {
		common.PrintError(err)
		os.Exit(1)
	}
	forceClean, err := cmd.Flags().GetBool("force-clean")
	if err != nil {
		common.PrintError(err)
		os.Exit(1)
	}
	deleteArg := "-d"
	if forceClean {
		deleteArg = "-D"
	}
	err = common.GitCommand("branch", deleteArg, branchToClean)
	if err != nil {
		common.PrintError(err)
		os.Exit(1)
	}
	return
}
