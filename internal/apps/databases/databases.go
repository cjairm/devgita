package databases

import (
	"context"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/logger"
	"github.com/cjairm/devgita/pkg/promptui"
	"github.com/cjairm/devgita/pkg/utils"
)

type Databases struct {
	Cmd cmd.Command
}

func New() *Databases {
	osCmd := cmd.NewCommand()
	return &Databases{Cmd: osCmd}
}

func (d *Databases) InstallNative(db string) error {
	return d.Cmd.MaybeInstallPackage(db)
}

func (d *Databases) AvailableDatabases() []string {
	return []string{
		"All",
		"None",
		"Done",
		"Redis",
		"SQLite",
		"MySQL",
		"PostgreSQL",
	}
}

func (d *Databases) ChooseDatabases(ctx context.Context) (context.Context, error) {
	dbs := d.AvailableDatabases()
	selectedDatabases, err := promptui.MultiSelect("Select databases", dbs)
	logger.L().Info("Selected databases: ", selectedDatabases)
	if err != nil {
		return nil, err
	}
	initialConfig := config.ContextConfig{}
	initialConfig.SelectedDbs = selectedDatabases
	return config.WithConfig(ctx, initialConfig), nil
}

func (d *Databases) InstallChosen(ctx context.Context) {
	selections, ok := config.GetConfig(ctx)
	logger.L().Info("Installing chosen databases: ", selections.SelectedDbs)
	if ok {
		if len(selections.SelectedDbs) > 0 {
			for _, db := range selections.SelectedDbs {
				utils.PrintInfo("Installing " + db + " (if no previously installed)...")
				if err := d.InstallNative(db); err != nil {
					utils.PrintError("Error: Unable to install Python.")
				}
			}
		}
	}
}
