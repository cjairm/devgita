/*
 * Copyright © 2025 Carlos Mendez <carlos@hadaelectronics.com> | https://cjairm.me/
 */
package cmd

import (
	"fmt"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("devgita %s (commit: %s, built: %s)\n", Version, Commit, BuildDate)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
