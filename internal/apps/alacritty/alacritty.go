package alacritty

import (
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/files"
)

const alacrittyDir = "alacritty"

type Alacritty struct {
	Cmd  cmd.Command
	Base cmd.BaseCommand
}

func New() *Alacritty {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Alacritty{Cmd: osCmd, Base: *baseCmd}
}

func (a *Alacritty) Install() error {
	return a.Cmd.InstallDesktopApp("alacritty")
}

func (a *Alacritty) MaybeInstall() error {
	return a.Cmd.MaybeInstallDesktopApp("alacritty")
}

func (a *Alacritty) SetupApp() error {
	path := []string{alacrittyDir}
	return a.Base.CopyAppConfigDirToLocalConfigDir(path, path)
}

func (a *Alacritty) SetupFont() error {
	devgitaPath := []string{"fonts", alacrittyDir, "default"}
	localPath := []string{alacrittyDir}
	return a.Base.CopyAppConfigDirToLocalConfigDir(devgitaPath, localPath)
}

func (a *Alacritty) SetupTheme() error {
	devgitaPath := []string{"themes", alacrittyDir, "default"}
	localPath := []string{alacrittyDir}
	return a.Base.CopyAppConfigDirToLocalConfigDir(devgitaPath, localPath)
}

func (a *Alacritty) MaybeSetupApp() error {
	localConfig, err := a.Base.ConfigDir()
	if err != nil {
		return err
	}
	return maybeSetup(a.SetupApp, localConfig, []string{alacrittyDir, "alacritty.toml"})
}

func (a *Alacritty) MaybeSetupFont() error {
	localConfig, err := a.Base.ConfigDir()
	if err != nil {
		return err
	}
	return maybeSetup(a.SetupFont, localConfig, []string{alacrittyDir, "font.toml"})
}

func (a *Alacritty) MaybeSetupTheme() error {
	localConfig, err := a.Base.ConfigDir()
	if err != nil {
		return err
	}
	return maybeSetup(a.SetupTheme, localConfig, []string{alacrittyDir, "theme.toml"})
}

func (a *Alacritty) UpdateConfigFilesWithCurrentHomeDir() error {
	localConfig, err := a.Base.ConfigDir()
	if err != nil {
		return err
	}
	alacrittyConfigFile := filepath.Join(localConfig, alacrittyDir, "alacritty.toml")
	return files.UpdateFile(alacrittyConfigFile, "<ALACRITTY-CONFIG-PATH>", localConfig)
}

func maybeSetup(setupFunc func() error, localConfig string, fileSegments []string) error {
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
