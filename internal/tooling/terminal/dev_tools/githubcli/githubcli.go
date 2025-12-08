// GitHub CLI (gh) tool with devgita integration
//
// GitHub CLI is the official command-line tool for GitHub, providing access to
// pull requests, issues, releases, repositories, gists, and more directly from
// the terminal. This module provides installation and configuration management
// for gh with devgita integration.
//
// References:
// - GitHub CLI Documentation: https://cli.github.com/manual/
// - GitHub CLI Repository: https://github.com/cli/cli
//
// Common gh commands available through ExecuteCommand():
//   - gh --version - Show gh version information
//   - gh auth login - Authenticate with GitHub
//   - gh auth status - View authentication status
//   - gh repo clone <repo> - Clone a repository
//   - gh pr list - List pull requests
//   - gh pr create - Create a new pull request
//   - gh pr checkout <number> - Check out a pull request locally
//   - gh issue list - List issues
//   - gh issue create - Create a new issue
//   - gh release list - List releases
//   - gh api <endpoint> - Make authenticated GitHub API requests

package githubcli

import (
	"fmt"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/pkg/constants"
)

type GithubCli struct {
	Cmd  cmd.Command
	Base cmd.BaseCommandExecutor
}

func New() *GithubCli {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &GithubCli{Cmd: osCmd, Base: baseCmd}
}

func (g *GithubCli) Install() error {
	return g.Cmd.InstallPackage(constants.GithubCli)
}

func (g *GithubCli) SoftInstall() error {
	return g.Cmd.MaybeInstallPackage(constants.GithubCli)
}

func (g *GithubCli) ForceInstall() error {
	err := g.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall gh: %w", err)
	}
	return g.Install()
}

func (g *GithubCli) Uninstall() error {
	return fmt.Errorf("gh uninstall not supported through devgita")
}

func (g *GithubCli) ForceConfigure() error {
	// GitHub CLI configuration is typically handled via:
	// - gh auth login (interactive authentication)
	// - gh config set <key> <value> (setting configuration values)
	// Configuration is usually handled via command-line operations
	// rather than copying config files
	return nil
}

func (g *GithubCli) SoftConfigure() error {
	// GitHub CLI configuration is typically handled via:
	// - gh auth login (interactive authentication)
	// - gh config set <key> <value> (setting configuration values)
	// Configuration is usually handled via command-line operations
	// rather than copying config files
	return nil
}

func (g *GithubCli) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		IsSudo:  false,
		Command: constants.GithubCli,
		Args:    args,
	}
	if _, _, err := g.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run gh command: %w", err)
	}
	return nil
}

func (g *GithubCli) Update() error {
	return fmt.Errorf("gh update not implemented through devgita")
}
