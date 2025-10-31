package fonts

import (
	"errors"
	"fmt"

	cmd "github.com/cjairm/devgita/internal/commands"
)

type Fonts struct {
	Cmd cmd.Command
}

func New() *Fonts {
	osCmd := cmd.NewCommand()
	return &Fonts{Cmd: osCmd}
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
	return f.Cmd.MaybeInstallFont("", fontName, false)
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
	return []string{
		"font-hack-nerd-font",
		"font-meslo-lg-nerd-font",
		"font-caskaydia-mono-nerd-font",
		"font-fira-mono",
		"font-jetbrains-mono-nerd-font",
	}
}

func (f *Fonts) SoftInstallAll() {
	availableFonts := f.Available()
	if len(availableFonts) > 0 {
		for _, font := range availableFonts {
			f.SoftInstall(font)
		}
	}
}
