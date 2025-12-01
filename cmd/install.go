/*
* Copyright © 2025 Carlos Mendez <carlos@hadaelectronics.com> | https://cjairm.me/
 */
package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/cjairm/devgita/internal/apps/git"
	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/internal/tooling/databases"
	"github.com/cjairm/devgita/internal/tooling/desktop"
	"github.com/cjairm/devgita/internal/tooling/languages"
	"github.com/cjairm/devgita/internal/tooling/terminal"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/logger"
	"github.com/cjairm/devgita/pkg/paths"
	"github.com/cjairm/devgita/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	only []string
	skip []string
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
  --only <...>     Only install specific categories (e.g., terminal, languages, desktop)
  --skip <...>     Skip specific categories (e.g., terminal, languages, desktop)
`,
	Run: run,
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Multi-value flags
	installCmd.Flags().
		StringSliceVar(&only, "only", []string{}, "Only install specific categories (comma or repeatable)")
	installCmd.Flags().
		StringSliceVar(&skip, "skip", []string{}, "Skip specific categories (comma or repeatable)")
}

func run(cmd *cobra.Command, args []string) {
	onlySet := make(map[string]bool)
	for _, item := range only {
		onlySet[item] = true
	}

	skipSet := make(map[string]bool)
	for _, item := range skip {
		skipSet[item] = true
	}

	logger.L().Debugw("flags", "only", onlySet, "skip", skipSet, "verbose", verbose)

	utils.PrintBold(constants.Devgita)
	utils.Print("=> Begin installation (or abort with ctrl+c)...", "")
	utils.Print("===============================================", "")

	ctx := context.Background()
	osCmd := commands.NewCommand()

	utils.PrintInfo("Validating version...")
	utils.MaybeExitWithError(osCmd.ValidateOSVersion())

	utils.PrintInfo("Installing essential tools to begin...")
	utils.MaybeExitWithError(osCmd.MaybeInstallPackageManager())
	g := git.New()
	utils.MaybeExitWithError(g.SoftInstall())

	installDevgita(g)
	setupDevgita()

	installTerminalTools(onlySet, skipSet)
	installLanguages(ctx, onlySet, skipSet)
	installDatabases(ctx, onlySet, skipSet)
	installDesktopTools(onlySet, skipSet)
}

func setupDevgita() {
	utils.PrintInfo("- Setup global config")
	globalConfig := &config.GlobalConfig{}
	utils.MaybeExitWithError(globalConfig.Create())
	utils.MaybeExitWithError(globalConfig.Load())
	globalConfig.AppPath = paths.AppDir
	globalConfig.ConfigPath = filepath.Join(paths.ConfigDir, constants.AppName)
	utils.MaybeExitWithError(globalConfig.Save())
}

func installDevgita(g *git.Git) {
	utils.PrintInfo("- Install devgita")
	// Create folder if it doesn't exist
	utils.MaybeExitWithError(os.MkdirAll(paths.AppDir, 0755))
	// Clean folder before (re)installing
	utils.MaybeExitWithError(os.RemoveAll(paths.AppDir))
	// Install from repository
	utils.MaybeExitWithError(g.Clone(constants.DevgitaRepositoryUrl, paths.AppDir))
}

func shouldInstall(category string, only, skip map[string]bool) bool {
	if len(only) > 0 {
		return only[category]
	}
	return !skip[category]
}

func installTerminalTools(onlySet, skipSet map[string]bool) {
	utils.PrintInfo("Preparing to install essential tools and packages...")
	if shouldInstall("terminal", onlySet, skipSet) {
		t := terminal.New()
		t.InstallAll()
		err := t.ConfigureZsh()
		utils.MaybeExitWithError(err)
	} else {
		utils.PrintInfo("Skipping terminal tools installation")
	}
}

func installLanguages(ctx context.Context, onlySet, skipSet map[string]bool) {
	utils.PrintInfo("Installing dev languages")
	if shouldInstall("languages", onlySet, skipSet) {
		utils.PrintInfo("Installing development languages")
		l := languages.New()
		ctx, err := l.ChooseLanguages(ctx)
		utils.MaybeExitWithError(err)
		l.InstallChosen(ctx)
	} else {
		utils.PrintInfo("Skipping development languages installation")
	}
}

func installDatabases(ctx context.Context, onlySet, skipSet map[string]bool) {
	utils.PrintInfo("Installing databases")
	if shouldInstall("databases", onlySet, skipSet) {
		utils.PrintInfo("Installing databases")
		d := databases.New()
		ctx, err := d.ChooseDatabases(ctx)
		utils.MaybeExitWithError(err)
		d.InstallChosen(ctx)
	} else {
		utils.PrintInfo("Skipping databases installation")
	}
}

func installDesktopTools(onlySet, skipSet map[string]bool) {
	utils.PrintInfo("Installing desktop applications")
	if shouldInstall("desktop", onlySet, skipSet) {
		utils.PrintInfo("Installing desktop applications")
		desktopTool := desktop.New()
		desktopTool.InstallAll()
	} else {
		utils.PrintInfo("Skipping desktop applications installation")
	}
}
