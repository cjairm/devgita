package fonts

import (
	"errors"
	"fmt"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/logger"
)

type Fonts struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func New() *Fonts {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Fonts{Cmd: osCmd, Base: baseCmd}
}

func (f *Fonts) Install(fontName string) error {
	return f.Cmd.InstallDesktopApp(fontName)
}

func (f *Fonts) ForceInstall(fontName string) error {
	err := f.Uninstall(fontName)
	if err != nil {
		return fmt.Errorf("failed to uninstall fonts: %w", err)
	}
	return f.Install(fontName)
}

func (f *Fonts) SoftInstall(fontName string) error {
	if f.Base.IsMac() {
		// macOS: Homebrew handles font installation, URL is ignored
		return f.Cmd.MaybeInstallFont("", fontName, false)
	}

	// Debian: download tar.xz from GitHub Nerd Fonts releases
	fc := constants.GetFontConfigByPackageName(fontName)
	if fc == nil {
		logger.L().Warnw("No Debian font config found, skipping", "font", fontName)
		return nil
	}

	fontURL := constants.GetNerdFontURL(fc.ArchiveName)
	return f.Cmd.MaybeInstallFont(fontURL, fc.InstallName, true)
}

func (f *Fonts) ForceConfigure() error {
	return nil
}

func (f *Fonts) SoftConfigure() error {
	return nil
}

func (f *Fonts) Uninstall(_fontName string) error {
	return errors.New("font uninstallation is not supported")
}

func (f *Fonts) ExecuteCommand(args ...string) error {
	return nil
}

func (f *Fonts) Update() error {
	return errors.New("font updates are not implemented - use system package manager")
}

func (f *Fonts) Available() []string {
	fontConfigs := constants.GetFontConfigs()
	names := make([]string, len(fontConfigs))
	for i, fc := range fontConfigs {
		names[i] = fc.PackageName
	}
	return names
}

func (f *Fonts) SoftInstallAll() {
	availableFonts := f.Available()
	for _, font := range availableFonts {
		if err := f.SoftInstall(font); err != nil {
			logger.L().Warnw("Font installation failed, continuing", "font", font, "error", err)
		}
	}
}
