/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"runtime"

	"github.com/cjairm/devgita/pkg/common"
	"github.com/cjairm/devgita/pkg/debian"
	"github.com/cjairm/devgita/pkg/macos"
	macosInstall "github.com/cjairm/devgita/pkg/macos/install"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "",
	Long:  ``,
	Run:   run,
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println(common.Devgita)
	fmt.Printf("=> Begin installation (or abort with ctrl+c)... \n\n")

	os := runtime.GOOS
	switch os {
	case "darwin":
		macos.PreInstall()

		fmt.Printf("Checking version...\n\n")
		macosInstall.CheckVersion()

		fmt.Printf("Cloning repo...\n\n")
		if err := common.CloneDevgita(); err != nil {
			panic(err)
		}

		fmt.Printf("Starting installation...\n\n")
	case "linux":
		debian.PreInstall()
		// Check if common.CloneDevgita works here
	default:
		fmt.Printf("Unsupported operating system: %s\n", os)
	}
}
