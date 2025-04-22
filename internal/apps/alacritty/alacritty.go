package alacritty

import (
	"path/filepath"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
)

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
	return files.CopyDir(paths.AlacrittyConfigAppDir, paths.AlacrittyConfigLocalDir)
}

func (a *Alacritty) SetupFont() error {
	return files.CopyDir(
		filepath.Join(paths.AlacrittyFontAppDir, "default"),
		paths.AlacrittyConfigLocalDir,
	)
}

func (a *Alacritty) SetupTheme() error {
	return files.CopyDir(
		filepath.Join(paths.AlacrittyThemesAppDir, "default"),
		paths.AlacrittyConfigLocalDir,
	)
}

func (a *Alacritty) MaybeSetupApp() error {
	return maybeSetup(a.SetupApp, paths.AlacrittyConfigLocalDir, "alacritty.toml")
}

func (a *Alacritty) MaybeSetupFont() error {
	return maybeSetup(a.SetupFont, paths.AlacrittyConfigLocalDir, "font.toml")
}

func (a *Alacritty) MaybeSetupTheme() error {
	return maybeSetup(a.SetupFont, paths.AlacrittyConfigLocalDir, "theme.toml")
}

func (a *Alacritty) UpdateConfigFilesWithCurrentHomeDir() error {
	alacrittyConfigFile := filepath.Join(paths.AlacrittyConfigLocalDir, "alacritty.toml")
	return files.UpdateFile(alacrittyConfigFile, "<ALACRITTY-CONFIG-PATH>", paths.ConfigDir)
}

func maybeSetup(setupFunc func() error, localConfig string, fileSegments ...string) error {
	filePath := localConfig
	for _, segment := range fileSegments {
		filePath = filepath.Join(filePath, segment)
	}
	if isFilePresent := files.FileAlreadyExist(filePath); isFilePresent {
		return nil
	}
	return setupFunc()
}
