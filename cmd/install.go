/*
* Copyright © 2025 Carlos Mendez <carlos@hadaelectronics.com> | https://cjairm.me/
 */
package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/cjairm/devgita/internal/apps/databases"
	"github.com/cjairm/devgita/internal/apps/desktop"
	devlanguages "github.com/cjairm/devgita/internal/apps/devLanguages"
	"github.com/cjairm/devgita/internal/apps/git"
	"github.com/cjairm/devgita/internal/apps/terminal"
	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/paths"
	"github.com/cjairm/devgita/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	dryRun       bool
	forceInstall bool
	only         []string
	skip         []string
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install devgita and all required tools",
	Long: `Installs the devgita platform and sets up your development environment.

This command performs the following steps:
  - Validates your OS version
  - Installs essential dependencies (e.g., git, fc-*)
  - Clones the devgita repository
  - Installs terminal tools, programming languages, and databases
  - Optionally installs desktop applications and shell configuration

Supported platforms:
  - macOS (via Homebrew)
  - Debian/Ubuntu (via apt)

Flags:
  --dry-run        Show what would be installed, without making any changes
  --only <...>     Only install specific categories (e.g., terminal, languages, desktop)
  --skip <...>     Skip specific categories (e.g., terminal, languages, desktop)
  --force          Reinstall tools even if they are already present
`,
	Run: run,
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Flags
	installCmd.Flags().
		BoolVar(&dryRun, "dry-run", false, "Show what would be installed without doing it")
	installCmd.Flags().
		BoolVar(&forceInstall, "force", false, "Force reinstallation even if components are already installed")

	// Multi-value flags
	installCmd.Flags().
		StringSliceVar(&only, "only", []string{}, "Only install specific categories (comma or repeatable)")
	installCmd.Flags().
		StringSliceVar(&skip, "skip", []string{}, "Skip specific categories (comma or repeatable)")
}

func run(cmd *cobra.Command, args []string) {
	if dryRun {
		utils.PrintInfo("Dry run: nothing will be installed.")
	}

	onlySet := make(map[string]bool)
	for _, item := range only {
		onlySet[item] = true
	}

	skipSet := make(map[string]bool)
	for _, item := range skip {
		skipSet[item] = true
	}

	// Example of how to use the shouldInstall function
	// if shouldInstall("terminal", onlySet, skipSet) {
	// 	fmt.Println("🔧 Installing terminal tools...")
	// 	// installTerminal()
	// }

	var err error

	utils.PrintBold(constants.Devgita)
	utils.Print("=> Begin installation (or abort with ctrl+c)...", "")
	utils.Print("===============================================", "")

	ctx := context.Background()
	osCmd := commands.NewCommand()

	utils.PrintInfo("* Validate version")
	utils.MaybeExitWithError(osCmd.ValidateOSVersion(verbose))

	utils.PrintInfo("- Pre-install steps")
	utils.MaybeExitWithError(osCmd.MaybeInstallPackageManager())
	utils.MaybeExitWithError(osCmd.MaybeInstallPackage("git"))

	installDevgita()

	setupDevgitaConfig()

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

func setupDevgitaConfig() {
	utils.PrintInfo("- Setup global config")
	utils.MaybeExitWithError(config.CreateGlobalConfig())
	globalConfig, err := config.LoadGlobalConfig()
	utils.MaybeExitWithError(err)
	globalConfig.AppPath = paths.AppDir
	globalConfig.ConfigPath = filepath.Join(paths.ConfigDir, constants.AppName)
	utils.MaybeExitWithError(config.SetGlobalConfig(globalConfig))
}

func installDevgita() {
	utils.PrintInfo("- Install devgita")
	// Create folder if it doesn't exist
	utils.MaybeExitWithError(os.MkdirAll(paths.AppDir, 0755))
	// Clean folder before (re)installing
	utils.MaybeExitWithError(os.RemoveAll(paths.AppDir))
	// Install from repository
	g := git.New()
	utils.MaybeExitWithError(g.Clone(constants.DevgitaRepositoryUrl, paths.AppDir))
}

func shouldInstall(category string, only, skip map[string]bool) bool {
	if len(only) > 0 {
		return only[category]
	}
	return !skip[category]
}
