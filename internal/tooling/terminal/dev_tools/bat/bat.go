// Bat - A cat clone with syntax highlighting and Git integration
//
// Bat is a cat clone with syntax highlighting, Git integration, and automatic paging.
// This module provides installation and configuration management for bat with devgita integration.
//
// References:
// - Bat Repository: https://github.com/sharkdp/bat
// - Bat Documentation: https://github.com/sharkdp/bat#readme
//
// Common bat commands available through ExecuteCommand():
//   - bat --version - Show bat version information
//   - bat <file> - Display file with syntax highlighting
//   - bat -l <lang> <file> - Force specific language syntax
//   - bat --list-languages - Show all supported languages
//   - bat --theme <theme> - Use specific color theme
//   - bat --list-themes - Show all available themes
//   - bat -n <file> - Show line numbers
//   - bat -A <file> - Show all characters including non-printable
//   - bat --style <components> - Configure which UI elements to display

package bat

import (
	"fmt"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
)

type Bat struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func New() *Bat {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Bat{Cmd: osCmd, Base: baseCmd}
}

func (b *Bat) Install() error {
	return b.Cmd.InstallPackage(constants.Bat)
}

func (b *Bat) SoftInstall() error {
	return b.Cmd.MaybeInstallPackage(constants.Bat)
}

func (b *Bat) ForceInstall() error {
	err := b.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall bat: %w", err)
	}
	return b.Install()
}

func (b *Bat) Uninstall() error {
	return fmt.Errorf("bat uninstall not supported through devgita")
}

func (b *Bat) ForceConfigure() error {
	// Bat typically doesn't require separate configuration files for basic usage
	// Configuration is usually handled via command-line arguments or optional config file
	// Users can create ~/.config/bat/config if desired for custom settings

	// TODO: Replace `cat` with this app.
	//
	// Ex, export HOME={{.Home}}
	//
	// func main() {
	// 	tmpl, err := template.ParseFiles("myfile.zsh")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	data := map[string]string{
	// 		"Home": "/User/Somethin/haha",
	// 	}
	//
	// 	outputFile, err := os.Create("myfile.generated.zsh")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer outputFile.Close()
	//
	// 	err = tmpl.Execute(outputFile, data)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	return nil
}

func (b *Bat) SoftConfigure() error {
	// Bat typically doesn't require separate configuration files for basic usage
	// Configuration is usually handled via command-line arguments or optional config file
	// Users can create ~/.config/bat/config if desired for custom settings

	// TODO: Replace `cat` with this app.
	//
	// Ex, export HOME={{.Home}}
	//
	// func main() {
	// 	tmpl, err := template.ParseFiles("myfile.zsh")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	data := map[string]string{
	// 		"Home": "/User/Somethin/haha",
	// 	}
	//
	// 	outputFile, err := os.Create("myfile.generated.zsh")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer outputFile.Close()
	//
	// 	err = tmpl.Execute(outputFile, data)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	return nil
}

func (b *Bat) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		IsSudo:  false,
		Command: constants.Bat,
		Args:    args,
	}
	if _, _, err := b.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run bat command: %w", err)
	}
	return nil
}

func (b *Bat) Update() error {
	return fmt.Errorf("bat update not implemented through devgita")
}
