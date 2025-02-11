/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime"

	commands "github.com/cjairm/devgita/internal"
	git "github.com/cjairm/devgita/internal/commands/git"
	"github.com/cjairm/devgita/internal/commands/powerlevel10k"
	bash "github.com/cjairm/devgita/internal/commands/zsh"
	"github.com/cjairm/devgita/pkg/common"
	"github.com/cjairm/devgita/pkg/debian"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/macos"
	macosInstall "github.com/cjairm/devgita/pkg/macos/install"
	macosTerminal "github.com/cjairm/devgita/pkg/macos/install/terminal"
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
	devgitaInstallPath, err := utils.GetDevgitaPath()
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
	g := git.NewGit()
	utils.MaybeExitWithError(files.CleanDestinationDir(devgitaInstallPath))
	utils.MaybeExitWithError(g.Clone(utils.DevgitaRepositoryUrl, devgitaInstallPath))

	// utils.PrintInfo("Preparing to install essential tools and packages...")
	// t := terminal.NewTerminal()
	// t.InstallAll()
	//
	// utils.PrintInfo("Installing development languages")
	// dl := devlanguages.NewDevLanguages()
	// ctx, err = dl.ChooseLanguages(ctx)
	// utils.MaybeExitWithError(err)
	// dl.InstallChosen(ctx)
	//
	// utils.PrintInfo("Installing databases")
	// db := databases.NewDatabases()
	// ctx, err = db.ChooseDatabases(ctx)
	// utils.MaybeExitWithError(err)
	// db.InstallChosen(ctx)
	//
	// utils.PrintInfo("Preparing to install desktop apps...")
	// d := desktop.NewDesktop()
	// d.InstallAll()

	utils.PrintInfo("Configuring custom commands...")
	b := bash.NewBash()
	err = b.CopyCustomConfig()
	utils.MaybeExitWithError(err)

	utils.PrintInfo("Installing terminal theme...")
	p := powerlevel10k.NewPowerLevel10k()
	p.MaybeInstall()
	p.MaybeSetup()
	p.Reconfigure()

	os.Exit(0)

	switch runtime.GOOS {
	case "darwin":
		// IMPORTANT....!!!
		// IF you are using an M mac computer... make sure you active Rosseta first
		/////////////////////////
		macos.PreInstall()

		fmt.Printf("Checking version...\n\n")
		macos.CheckVersion()

		fmt.Printf("Cloning repo...\n\n")
		if err := common.CloneDevgita(devgitaInstallPath); err != nil {
			fmt.Printf("\033[31mError: %s\033[0m\n", err.Error())
			fmt.Println("Installation stopped.")
			os.Exit(1)
		}

		ctx = macosTerminal.ChooseLanguages(ctx)
		ctx = macosTerminal.ChooseDatabases(ctx)

		fmt.Printf("Starting installation...\n\n")
		macosInstall.RunTerminalInstallers(devgitaInstallPath)

		macosTerminal.InstallDatabases(ctx)
		macosTerminal.InstallDevLanguages(ctx)

		macosInstall.RunDesktopInstallers(devgitaInstallPath)
		common.Reboot()
	case "linux":
		debian.PreInstall()
		// Check if common.CloneDevgita works here
	default:
		fmt.Printf("Unsupported operating system: %s\n", runtime.GOOS)
	}
}
