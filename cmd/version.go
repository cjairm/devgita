/*
 * Copyright © 2025 Carlos Mendez <carlos@hadaelectronics.com> | https://cjairm.me/
 */
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version is set during build via ldflags
	Version = "dev"
	// Commit is set during build via ldflags
	Commit = "unknown"
	// BuildDate is set during build via ldflags
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of devgita",
	Long:  `All software has versions. This is devgita's.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("devgita %s (commit: %s, built: %s)\n", Version, Commit, BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Add --version flag to root command
	rootCmd.Flags().BoolP("version", "v", false, "Print version information")

	// Check for --version flag before running any command
	rootCmd.ParseFlags(os.Args[1:])
	if versionFlag, _ := rootCmd.Flags().GetBool("version"); versionFlag {
		fmt.Printf("devgita %s (commit: %s, built: %s)\n", Version, Commit, BuildDate)
		os.Exit(0)
	}
}
