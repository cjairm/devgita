package fonts

import cmd "github.com/cjairm/devgita/internal/commands"

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

func (f *Fonts) MaybeInstall(fontName string) error {
	return f.Cmd.MaybeInstallFont("", fontName, false)
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

func (f *Fonts) MaybeInstallAll() {
	availableFonts := f.Available()
	if len(availableFonts) > 0 {
		for _, font := range availableFonts {
			f.MaybeInstall(font)
		}
	}
}
