package devlanguages

import (
	"context"
	"strings"

	"github.com/cjairm/devgita/internal/apps/mise"
	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/logger"
	"github.com/cjairm/devgita/pkg/promptui"
	"github.com/cjairm/devgita/pkg/utils"
)

type DevLanguages struct {
	Cmd cmd.Command
}

func New() *DevLanguages {
	osCmd := cmd.NewCommand()
	return &DevLanguages{Cmd: osCmd}
}

func (dl *DevLanguages) InstallWithMise(language, version string) error {
	m := mise.New()
	return m.UseGlobal(language, version)
}

// TODO: Add MaybeInstall with mise

func (dl *DevLanguages) InstallNative(language string) error {
	return dl.Cmd.MaybeInstallPackage(language)
}

func (dl *DevLanguages) AvailableLanguages() []string {
	return []string{
		"All",
		"None",
		"Done",
		"Node",
		"Go",
		"PHP",
		"Python",
	}
}

func (dl *DevLanguages) ChooseLanguages(ctx context.Context) (context.Context, error) {
	languages := dl.AvailableLanguages()
	selectedLanguages, err := promptui.MultiSelect("Select programming languages", languages)
	logger.L().Info("Selected languages: ", selectedLanguages)
	if err != nil {
		return nil, err
	}
	initialConfig := config.ContextConfig{}
	initialConfig.SelectedLanguages = selectedLanguages
	return config.WithConfig(ctx, initialConfig), nil
}

func (dl *DevLanguages) InstallChosen(ctx context.Context) {
	selections, ok := config.GetConfig(ctx)
	logger.L().Info("Installing chosen languages: ", selections.SelectedLanguages)
	if ok {
		if len(selections.SelectedLanguages) > 0 {
			for _, language := range selections.SelectedLanguages {
				switch strings.ToLower(language) {
				case "node":
					utils.PrintInfo("Installing Node LTS (if no previously installed)...")
					if err := dl.InstallWithMise("node", "lts"); err != nil {
						utils.PrintError("Error: Unable to install Node.")
					}
				case "go":
					utils.PrintInfo("Installing Go latest (if no previously installed)...")
					if err := dl.InstallWithMise("go", "latest"); err != nil {
						utils.PrintError("Error: Unable to install Go.")
					}
				case "python":
					utils.PrintInfo("Installing Python latest (if no previously installed)...")
					if err := dl.InstallWithMise("python", "latest"); err != nil {
						utils.PrintError("Error: Unable to install Python.")
					}
				case "php":
					utils.PrintInfo("Installing PHP latest (if no previously installed)...")
					if err := dl.InstallNative("php"); err != nil {
						utils.PrintError("Error: Unable to install PHP.")
					}
				}
			}
		}
	}
}
