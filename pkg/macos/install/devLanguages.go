package macos

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cjairm/devgita/pkg/common"
)

var languages = []string{
	"All",
	"None",
	"Done",
	"Node",
	"Go",
	"PHP",
	"Python",
}

func ChooseLanguages(ctx context.Context) context.Context {
	selectedLanguages, err := common.MultiSelect("Select programming languages", languages)
	if err != nil {
		fmt.Println("\033[31mError: Error selecting languages.\033[0m")
		fmt.Println("Installation stopped.")
		os.Exit(1)
	}
	initialConfig := common.Config{}
	initialConfig.SelectedLanguages = selectedLanguages
	fmt.Printf("\n")
	return common.WithConfig(ctx, initialConfig)
}

func InstallDevLanguages(ctx context.Context) {
	selections, ok := common.GetConfig(ctx)
	if ok {
		if len(selections.SelectedLanguages) > 0 {
			fmt.Printf("Installing languages...\n\n")
			for _, language := range selections.SelectedLanguages {
				switch strings.ToLower(language) {
				case "node":
					err := installNode()
					if err != nil {
						fmt.Println("\033[31mError: Unable to install Node.\033[0m")
					}
				case "go":
					err := installGo()
					if err != nil {
						fmt.Println("\033[31mError: Unable to install Go.\033[0m")
					}
				case "python":
					err := installPython()
					if err != nil {
						fmt.Println("\033[31mError: Unable to install Python.\033[0m")
					}
				case "php":
					err := common.InstallOrUpdateBrewPackage("php")
					if err != nil {
						fmt.Println("\033[31mError: Unable to install PHP.\033[0m")
					}
				}
			}

		} else {
			fmt.Printf("No databases installed...\n\n")
		}
	} else {
		fmt.Printf("\033[31mError: Skip language selections\033[0m\n")
	}
}

func installNode() error {
	cmd := common.CommandInfo{
		PreExecutionMessage:  "Installing Node using Mise",
		PostExecutionMessage: "Node installed ✔",
		IsSudo:               false,
		Command:              "mise",
		Args:                 []string{"use", "--global", "node@lts"},
	}
	return common.ExecCommand(cmd)
}

func installGo() error {
	cmd := common.CommandInfo{
		PreExecutionMessage:  "Installing Golang using Mise",
		PostExecutionMessage: "Golang installed ✔",
		IsSudo:               false,
		Command:              "mise",
		Args:                 []string{"use", "--global", "go@latest"},
	}
	return common.ExecCommand(cmd)
}

func installPython() error {
	cmd := common.CommandInfo{
		PreExecutionMessage:  "Installing Python using Mise",
		PostExecutionMessage: "Python installed ✔",
		IsSudo:               false,
		Command:              "mise",
		Args:                 []string{"use", "--global", "python@latest"},
	}
	return common.ExecCommand(cmd)
}
