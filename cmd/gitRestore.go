/*
Copyright © 2025 Carlos Mendez <carlos@hadaelectronics.com> | https://cjairm.me/
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	git "github.com/cjairm/devgita/internal/commands/git"
	"github.com/cjairm/devgita/pkg/utils"
	"github.com/spf13/cobra"
)

const GIT_RESTORE_CMD = "gr"

var gitRestoreCmd = &cobra.Command{
	Use:   GIT_RESTORE_CMD,
	Short: "Restore files from a specified branch or commit",
	Long: fmt.Sprintf(
		`The git-restore command allows you to restore files from a specified branch or commit in your Git repository. 
You can choose to restore specific files or all files from the source branch. 
By default, if no source branch is specified, the command will restore files from the "main" branch.

Usage:
  devgita %s --source-branch <branch-name> --files <file1> <file2> ...
  devgita %s --files <file1> <file2> ...

Flags:
  -s, --source-branch string   The name of the branch to fetch the document from (default "main")
  -f, --files                  The paths of the files to restore (required)`,
		GIT_RESTORE_CMD,
		GIT_RESTORE_CMD,
	),
	Run: runGitRestore,
}

func init() {
	rootCmd.AddCommand(gitRestoreCmd)
	gitRestoreCmd.Flags().
		StringP("source-branch", "s", "", "The name of the branch to fetch the document from (optional)")
	gitRestoreCmd.Flags().
		StringSliceP("files", "f", []string{}, "The paths of the files to restore (required)")
}

func runGitRestore(cmd *cobra.Command, args []string) {
	srcBranch, err := cmd.Flags().GetString("source-branch")
	utils.MaybeExitWithError(err)

	filePaths, err := cmd.Flags().GetStringSlice("files")
	utils.MaybeExitWithError(err)

	if len(filePaths) == 0 {
		utils.PrintInfo("No files to restore")
		os.Exit(1)
	}

	g := git.New()
	err = g.Restore(srcBranch, strings.Join(filePaths, " "))
	utils.MaybeExitWithError(err)

	successMessage := fmt.Sprintf("Files %s restored", srcBranch)
	utils.PrintSuccess(successMessage)
	os.Exit(0)
}
