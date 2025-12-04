// Eza modern ls replacement tool with devgita integration
//
// Eza is a modern, maintained replacement for ls with improved features including
// colors, icons, git integration, and tree views. This module provides installation
// and command execution management for eza with devgita integration.
//
// References:
// - Eza Repository: https://github.com/eza-community/eza
// - Eza Documentation: https://github.com/eza-community/eza/blob/main/README.md
//
// Common eza commands available through ExecuteCommand():
//   - eza --version - Show eza version information
//   - eza - List directory contents with colors
//   - eza -l - Long format listing
//   - eza -T - Tree view
//   - eza -a - Show hidden files
//   - eza --git - Show git status
//   - eza --icons - Show file icons
//   - eza -lah - Long format with hidden files and human-readable sizes

package eza

import (
	"fmt"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
)

type Eza struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func New() *Eza {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &Eza{Cmd: osCmd, Base: baseCmd}
}

func (e *Eza) Install() error {
	return e.Cmd.InstallPackage("eza")
}

func (e *Eza) SoftInstall() error {
	return e.Cmd.MaybeInstallPackage("eza")
}

func (e *Eza) ForceInstall() error {
	err := e.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall eza: %w", err)
	}
	return e.Install()
}

func (e *Eza) Uninstall() error {
	return fmt.Errorf("eza uninstall not supported through devgita")
}

func (e *Eza) ForceConfigure() error {
	// Eza typically doesn't require separate configuration files
	// Configuration is usually handled via command-line arguments and aliases

	// TODO: Replace `ls` with this app. Use `text/template`
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

func (e *Eza) SoftConfigure() error {
	// Eza typically doesn't require separate configuration files
	// Configuration is usually handled via command-line arguments and aliases

	// TODO: Replace `ls` with this app. Use `text/template`
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

func (e *Eza) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		IsSudo:  false,
		Command: constants.Eza,
		Args:    args,
	}
	if _, _, err := e.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run eza command: %w", err)
	}
	return nil
}

func (e *Eza) Update() error {
	return fmt.Errorf("eza update not implemented through devgita")
}
