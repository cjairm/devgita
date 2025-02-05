/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/cjairm/devgita/pkg/common"
	"github.com/spf13/cobra"
)

var gitCleanCmd = &cobra.Command{
	Use:   "git-clean",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		dstBranch, err := cmd.Flags().GetString("destination-branch")
		if err != nil {
			fmt.Printf("\033[31mError: %s\033[0m\n", err.Error())
			os.Exit(1)
		}
		if dstBranch == "" {
			dstBranch = "main"
		}
		err = runGitCommand("checkout", dstBranch)
		if err != nil {
			fmt.Printf("\033[31mError: %s\033[0m\n", err.Error())
			os.Exit(1)
		}
		err = runGitCommand("fetch", "origin")
		if err != nil {
			fmt.Printf("\033[31mError: %s\033[0m\n", err.Error())
			os.Exit(1)
		}
		err = runGitCommand("pull", "origin", dstBranch)
		if err != nil {
			fmt.Printf("\033[31mError: %s\033[0m\n", err.Error())
			os.Exit(1)
		}
		branchToClean, err := cmd.Flags().GetString("branch-to-clean")
		if err != nil {
			fmt.Printf("\033[31mError: %s\033[0m\n", err.Error())
			os.Exit(1)
		}
		forceClean, err := cmd.Flags().GetBool("force-clean")
		if err != nil {
			fmt.Printf("\033[31mError: %s\033[0m\n", err.Error())
			os.Exit(1)
		}
		deleteArg := "-d"
		if forceClean {
			deleteArg = "-D"
		}
		err = runGitCommand("branch", deleteArg, branchToClean)
		if err != nil {
			fmt.Printf("\033[31mError: %s\033[0m\n", err.Error())
			os.Exit(1)
		}
		return
	},
}

func init() {
	rootCmd.AddCommand(gitCleanCmd)
	gitCleanCmd.Flags().
		StringP("branch-to-clean", "c", "", "The name of the branch to clean (required)")
	gitCleanCmd.Flags().
		StringP("destination-branch", "d", "", "The name of the destination branch after cleaning")
	gitCleanCmd.Flags().BoolP("force-clean", "f", false, "Force the deletion of the branch")

	gitCleanCmd.MarkFlagRequired("branch-to-clean")
}

func runGitCommand(gitArgs ...string) error {
	if len(gitArgs) == 0 {
		return fmt.Errorf("No command provided")
	}
	switch runtime.GOOS {
	case "darwin":
		cmd := common.CommandInfo{
			PreExecutionMessage:  "",
			PostExecutionMessage: "",
			IsSudo:               false,
			Command:              "git",
			Args:                 gitArgs,
		}
		return common.ExecCommand(cmd)
	case "linux":
		fmt.Println("Linux")
	default:
		fmt.Printf("Unsupported operating system: %s\n", runtime.GOOS)
	}
	return nil
}
