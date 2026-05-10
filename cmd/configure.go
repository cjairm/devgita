/*
* Copyright © 2025 Carlos Mendez <carlos@hadaelectronics.com> | https://cjairm.me/
 */
package cmd

import (
	"errors"

	"github.com/cjairm/devgita/internal/apps"
	"github.com/cjairm/devgita/internal/apps/registry"
	"github.com/cjairm/devgita/pkg/utils"
	"github.com/spf13/cobra"
)

var configureForce bool

// getAppFn is the registry lookup; overridden in tests.
var getAppFn = func(name string) (apps.App, error) {
	return registry.GetApp(name)
}

var configureCmd = &cobra.Command{
	Use:   "configure [app]",
	Short: "Apply configuration files for a named app",
	Long: `Apply configuration files for a named app without reinstalling.

By default (soft mode), configuration is only applied if files do not already exist.
Use --force to overwrite existing configuration files.

Examples:
  dg configure git            # Apply git config if not already present
  dg configure neovim --force # Overwrite existing neovim config
  dg configure tmux           # Apply tmux config if not already present
`,
	Args: cobra.ExactArgs(1),
	RunE: runConfigure,
}

func init() {
	rootCmd.AddCommand(configureCmd)

	configureCmd.Flags().
		BoolVar(&configureForce, "force", false, "Overwrite existing configuration files")
}

func runConfigure(cmd *cobra.Command, args []string) error {
	appName := args[0]

	app, err := getAppFn(appName)
	if err != nil {
		return err
	}

	if configureForce {
		err = app.ForceConfigure()
	} else {
		err = app.SoftConfigure()
	}

	if errors.Is(err, apps.ErrConfigureNotSupported) {
		utils.PrintInfo("configure not supported for " + appName)
		return nil
	}
	if err != nil {
		return err
	}

	utils.PrintSuccess("configured " + appName)
	return nil
}
