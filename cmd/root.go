/*
Copyright © 2024 Carlos Mendez <carlos@hadaelectronics.com>
*/
package cmd

import (
	"os"

	commands "github.com/cjairm/devgita/internal"
	"github.com/cjairm/devgita/pkg/utils"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "devgita",
	Short: "Command-line tool for macOS Ventura that automates the setup of development environments, streamlining installations of essential apps like Raycast, Homebrew, and iTerm2 with clear documentation.",
	Long:  utils.Devgita,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	utils.MaybeExitWithError(err)
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.devgita.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.SetHelpFunc(utils.PrompCustomHelp)

	utils.MaybeExitWithError(setDevgitaPath())
}

func setDevgitaPath() error {
	bc := commands.NewBaseCommand()
	devgitaInstallPath, err := bc.AppDir()
	if err != nil {
		return nil
	}
	err = os.MkdirAll(devgitaInstallPath, 0755)
	if err != nil {
		return nil
	}
	return nil
}
