package alacritty

import (
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/files"
)

const alacrittyDir = "alacritty"

type Alacritty struct {
	Cmd cmd.Command
}

func NewAlacritty() *Alacritty {
	osCmd := cmd.NewCommand()
	return &Alacritty{Cmd: osCmd}
}

func (a *Alacritty) Install() error {
	return a.Cmd.InstallDesktopApp("alacritty")
}

func (a *Alacritty) MaybeInstall() error {
	return a.Cmd.MaybeInstallDesktopApp("alacritty")
}

func (a *Alacritty) SetupApp() error {
	path := []string{alacrittyDir}
	return files.MoveFromConfigsToLocalConfig(path, path)
}

func (a *Alacritty) SetupFont() error {
	devgitaPath := []string{"fonts", alacrittyDir, "default"}
	localPath := []string{alacrittyDir}
	return files.MoveFromConfigsToLocalConfig(devgitaPath, localPath)
}

func (a *Alacritty) SetupTheme() error {
	devgitaPath := []string{"themes", alacrittyDir, "default"}
	localPath := []string{alacrittyDir}
	return files.MoveFromConfigsToLocalConfig(devgitaPath, localPath)
}

func (a *Alacritty) MaybeSetupApp() error {
	return maybeSetup(a.SetupApp, []string{alacrittyDir, "alacritty.toml"})
}

func (a *Alacritty) MaybeSetupFont() error {
	return maybeSetup(a.SetupFont, []string{alacrittyDir, "font.toml"})
}

func (a *Alacritty) MaybeSetupTheme() error {
	return maybeSetup(a.SetupTheme, []string{alacrittyDir, "theme.toml"})
}

func (a *Alacritty) UpdateConfigFilesWithCurrentHomeDir() error {
	localConfig, err := config.GetLocalConfigPath()
	if err != nil {
		return err
	}
	alacrittyConfigFile := filepath.Join(localConfig, alacrittyDir, "alacritty.toml")
	return files.UpdateFile(alacrittyConfigFile, "<ALACRITTY-CONFIG-PATH>", localConfig)
}

func maybeSetup(setupFunc func() error, fileSegments []string) error {
	localConfig, err := config.GetLocalConfigPath()
	if err != nil {
		return err
	}
	filePath := localConfig
	for _, segment := range fileSegments {
		filePath = filepath.Join(filePath, segment)
	}
	isFilePresent := files.FileAlreadyExist(filePath)
	if isFilePresent {
		return nil
	}
	return setupFunc()
}
