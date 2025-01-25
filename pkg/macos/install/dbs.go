package macos

import (
	"context"
	"fmt"
	"os"

	"github.com/cjairm/devgita/pkg/common"
)

var dbs = []string{
	"All",
	"None",
	"Done",
	"Redis",
	"SQLite",
	"MySQL",
	"PostgreSQL",
}

func ChooseDatabases(ctx context.Context) context.Context {
	selectedDatabases, err := common.MultiSelect("Select databases", dbs)
	if err != nil {
		fmt.Println("\033[31mError: Error selecting databases.\033[0m")
		fmt.Println("Installation stopped.")
		os.Exit(1)
	}
	initialConfig, ok := common.GetConfig(ctx)
	if ok == false {
		fmt.Println("\033[31mError: Error storing selections of databases.\033[0m")
		fmt.Println("Installation stopped.")
		os.Exit(1)
	}
	initialConfig.SelectedDbs = selectedDatabases
	fmt.Printf("\n")
	return common.WithConfig(ctx, initialConfig)
}

func InstallDatabases(ctx context.Context) {
	selections, ok := common.GetConfig(ctx)
	if ok {
		if len(selections.SelectedDbs) > 0 {
			fmt.Printf("Installing databases...\n\n")
			for _, db := range selections.SelectedDbs {
				if err := common.InstallOrUpdateBrewPackage(db); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

		} else {
			fmt.Printf("No databases installed...\n\n")
		}
	} else {
		fmt.Printf("\033[31mError: Skip database selections\033[0m\n")
	}
}
