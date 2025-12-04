// Zoxide smart directory navigation tool with devgita integration
//
// Zoxide is a smarter cd command that learns your habits and allows you to navigate
// to frequently and recently used directories with just a few keystrokes. It tracks
// your most used directories and provides fuzzy matching for quick navigation.
//
// References:
// - Zoxide Documentation: https://github.com/ajeetdsouza/zoxide
// - Zoxide Wiki: https://github.com/ajeetdsouza/zoxide/wiki
//
// Common zoxide commands available through ExecuteCommand():
//   - zoxide --version - Show zoxide version information
//   - zoxide query <keywords> - Search for directories matching keywords
//   - zoxide add <path> - Add a directory to the database
//   - zoxide remove <path> - Remove a directory from the database
//   - zoxide query --list - List all tracked directories
//   - zoxide query --interactive - Interactive directory selection
//   - zoxide import <file> - Import directories from file
//   - zoxide init zsh - Generate shell initialization script
//
// Shell integration:
//   After installation, add to your shell config:
//   - zsh: eval "$(zoxide init zsh)"
//   - bash: eval "$(zoxide init bash)"
//   - fish: zoxide init fish | source
//
//   This enables the 'z' command for smart navigation:
//   - z foo - Jump to directory matching 'foo'
//   - zi foo - Interactive selection when multiple matches

package zoxide

import (
	"fmt"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
)

type Zoxide struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func New() *Zoxide {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Zoxide{Cmd: osCmd, Base: baseCmd}
}

func (z *Zoxide) Install() error {
	return z.Cmd.InstallPackage(constants.Zoxide)
}

func (z *Zoxide) SoftInstall() error {
	return z.Cmd.MaybeInstallPackage(constants.Zoxide)
}

func (z *Zoxide) ForceInstall() error {
	err := z.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall zoxide: %w", err)
	}
	return z.Install()
}

func (z *Zoxide) Uninstall() error {
	return fmt.Errorf("zoxide uninstall not supported through devgita")
}

func (z *Zoxide) ForceConfigure() error {
	// Zoxide typically doesn't require separate configuration files
	// Shell integration is handled via 'zoxide init' command
	// Configuration is usually handled via shell initialization scripts

	// TODO: Replace `cd` with this app.
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

func (z *Zoxide) SoftConfigure() error {
	// Zoxide typically doesn't require separate configuration files
	// Shell integration is handled via 'zoxide init' command
	// Configuration is usually handled via shell initialization scripts

	// TODO: Replace `cd` with this app.
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

func (z *Zoxide) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		IsSudo:  false,
		Command: constants.Zoxide,
		Args:    args,
	}
	if _, _, err := z.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run zoxide command: %w", err)
	}
	return nil
}

func (z *Zoxide) Update() error {
	return fmt.Errorf("zoxide update not implemented through devgita")
}
