/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	commands "github.com/cjairm/devgita/internal"
	"github.com/cjairm/devgita/internal/commands/databases"
	"github.com/cjairm/devgita/internal/commands/desktop"
	devlanguages "github.com/cjairm/devgita/internal/commands/devLanguages"
	git "github.com/cjairm/devgita/internal/commands/git"
	"github.com/cjairm/devgita/internal/commands/terminal"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/utils"
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
	bc := commands.NewBaseCommand()
	devgitaInstallPath, err := bc.GetDevgitaAppDir()
	utils.MaybeExitWithError(err)

	utils.PrintBold(utils.Devgita)
	utils.Print("=> Begin installation (or abort with ctrl+c)...", "")
	utils.Print("===============================================", "")

	ctx := context.Background()

	osCmd := commands.NewCommand()

	utils.PrintInfo("- Validate versions before start installation")
	utils.MaybeExitWithError(osCmd.ValidateOSVersion())

	utils.PrintInfo("- Check if package manager exist and install if needed")
	utils.MaybeExitWithError(osCmd.MaybeInstallPackageManager())

	utils.PrintInfo("- Install git dependency")
	utils.MaybeExitWithError(osCmd.MaybeInstallPackage("git"))

	utils.PrintInfo("- Install devgita")
	g := git.New()
	utils.MaybeExitWithError(files.CleanDestinationDir(devgitaInstallPath))
	utils.MaybeExitWithError(g.Clone(utils.DevgitaRepositoryUrl, devgitaInstallPath))

	utils.PrintInfo("Preparing to install essential tools and packages...")
	t := terminal.New()
	t.InstallAll()

	utils.PrintInfo("Installing development languages")
	dl := devlanguages.New()
	ctx, err = dl.ChooseLanguages(ctx)
	utils.MaybeExitWithError(err)
	dl.InstallChosen(ctx)

	utils.PrintInfo("Installing databases")
	db := databases.New()
	ctx, err = db.ChooseDatabases(ctx)
	utils.MaybeExitWithError(err)
	db.InstallChosen(ctx)

	utils.PrintInfo("Preparing to install desktop apps...")
	d := desktop.New()
	d.InstallAll()

	err = t.ConfigureZsh()
	utils.MaybeExitWithError(err)
}
