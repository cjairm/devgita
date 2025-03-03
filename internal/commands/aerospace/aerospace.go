package aerospace

import (
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/files"
)

const aerospaceDir = "aerospace"

type Aerospace struct {
	Cmd cmd.Command
}

func New() *Aerospace {
	osCmd := cmd.NewCommand()
	return &Aerospace{Cmd: osCmd}
}

func (a *Aerospace) Install() error {
	return a.Cmd.InstallDesktopApp("nikitabobko/tap/aerospace")
}

func (a *Aerospace) MaybeInstall() error {
	return a.Cmd.MaybeInstallDesktopApp("nikitabobko/tap/aerospace")
}

func (a *Aerospace) Setup() error {
	configPath := []string{aerospaceDir}
	return files.MoveFromConfigsToLocalConfig(configPath, configPath)
}

func (a *Aerospace) MaybeSetup() error {
	localConfig, err := config.GetLocalConfigPath()
	if err != nil {
		return err
	}
	aerospaceConfigFile := filepath.Join(localConfig, aerospaceDir, "aerospace.toml")
	isFilePresent := files.FileAlreadyExist(aerospaceConfigFile)
	if isFilePresent {
		return nil
	}
	return a.Setup()
}
