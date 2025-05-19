/*
Copyright © 2025 Carlos Mendez <carlos@hadaelectronics.com> | https://cjairm.me/
*/
package cmd

import (
	"context"
	"os"

	"github.com/cjairm/devgita/internal/apps/databases"
	"github.com/cjairm/devgita/internal/apps/desktop"
	devlanguages "github.com/cjairm/devgita/internal/apps/devLanguages"
	git "github.com/cjairm/devgita/internal/apps/git"
	"github.com/cjairm/devgita/internal/apps/terminal"
	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/paths"
	"github.com/cjairm/devgita/pkg/utils"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install devgita and all required tools",
	Long: `Installs the devgita platform and sets up your development environment.

This command performs the following steps:
  - Validates your OS version
  - Installs essential dependencies (ex. git, fc-*)
  - Clones the devgita repository
  - Installs terminal tools, programming languages, and databases
  - Sets up desktop apps and configures your shell

It supports both macOS (via Homebrew) and Debian/Ubuntu systems (via apt).`,
	Run: run,
}

func init() {
	rootCmd.AddCommand(installCmd)

	// --------------------------------------------------------------------
	// A --dry-run flag (show what would be installed without doing it)
	// A --verbose or --debug flag for verbose logs
	// Allow installing individual categories like --only-languages or --only-terminal
}

func run(cmd *cobra.Command, args []string) {
	var err error

	utils.PrintBold(constants.Devgita)
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
	// Create folder if it doesn't exist
	utils.MaybeExitWithError(os.MkdirAll(paths.AppDir, 0755))
	// Clean folder before (re)installing
	utils.MaybeExitWithError(os.RemoveAll(paths.AppDir))
	// Install from repository
	g := git.New()
	utils.MaybeExitWithError(g.Clone(constants.DevgitaRepositoryUrl, paths.AppDir))

	utils.MaybeExitWithError(config.CreateGlobalConfig())
	configFile, err := config.LoadGlobalConfig()
	utils.MaybeExitWithError(err)
	configFile.AppPath = paths.AppDir
	utils.MaybeExitWithError(config.SetGlobalConfig(configFile))

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
