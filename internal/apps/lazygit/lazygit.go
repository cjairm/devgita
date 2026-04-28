// LazyGit terminal UI for Git repository management with devgita integration
//
// LazyGit is a simple terminal UI for git commands, written in Go with the gocui library.
// It provides an interactive interface to manage Git repositories, branches, commits, and
// staging operations, all from the comfort of the terminal.
//
// References:
// - LazyGit Repository: https://github.com/jesseduffield/lazygit
// - LazyGit Documentation: https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md
//
// Common lazygit commands available through ExecuteCommand():
//   - lazygit - Launch interactive TUI
//   - lazygit --version - Show lazygit version information
//   - lazygit --config - Show configuration file path
//   - lazygit --help - Display help information

package lazygit

import (
	"context"
	"fmt"
	"runtime"

	cmd "github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/config"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/downloader"
	gh "github.com/cjairm/devgita/pkg/github"
	"github.com/cjairm/devgita/pkg/logger"
)

type LazyGit struct {
	Cmd          cmd.Command
	Base         cmd.BaseCommandExecutor
	fetchVersion func(owner, repo string) (string, error)                                      // injectable for tests
	downloadFn   func(ctx context.Context, url, dest string, cfg downloader.RetryConfig) error // injectable for tests
}

func (lg *LazyGit) getVersion(owner, repo string) (string, error) {
	if lg.fetchVersion != nil {
		return lg.fetchVersion(owner, repo)
	}
	return gh.FetchLatestRelease(owner, repo)
}

func New() *LazyGit {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &LazyGit{Cmd: osCmd, Base: baseCmd}
}

func (lg *LazyGit) Install() error {
	if lg.Base.IsMac() {
		return lg.Cmd.InstallPackage(constants.LazyGit)
	}
	return lg.installDebianLazygit()
}

func (lg *LazyGit) SoftInstall() error {
	if lg.Base.IsMac() {
		return lg.Cmd.MaybeInstallPackage(constants.LazyGit)
	}
	// On Debian: lazygit is not in apt — skip if already in PATH, otherwise install from GitHub
	if _, err := cmd.LookPathFn(constants.LazyGit); err == nil {
		logger.L().Infow("lazygit already installed, skipping")
		return nil
	}
	return lg.installDebianLazygit()
}

// linuxArch returns the architecture string used in lazygit release artifact names.
func linuxArch() string {
	if runtime.GOARCH == "arm64" {
		return "arm64"
	}
	return "x86_64"
}

// installDebianLazygit fetches the latest lazygit release from GitHub, downloads
// the Linux tar.gz for the current architecture, and installs the binary to /usr/local/bin
func (lg *LazyGit) installDebianLazygit() error {
	version, err := lg.getVersion("jesseduffield", "lazygit")
	if err != nil {
		return fmt.Errorf("failed to fetch lazygit version: %w", err)
	}

	url := fmt.Sprintf(
		"https://github.com/jesseduffield/lazygit/releases/download/v%s/lazygit_%s_Linux_%s.tar.gz",
		version, version, linuxArch(),
	)
	logger.L().Infow("Downloading lazygit for Debian", "version", version, "url", url)

	if err := cmd.InstallGitHubBinary(lg.Base, constants.LazyGit, url, lg.downloadFn); err != nil {
		return fmt.Errorf("lazygit installation failed: %w", err)
	}

	logger.L().Infow("lazygit installed successfully for Debian", "version", version)
	return nil
}

func (lg *LazyGit) ForceInstall() error {
	err := lg.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall lazygit: %w", err)
	}
	return lg.Install()
}

func (lg *LazyGit) Uninstall() error {
	return fmt.Errorf("lazygit uninstall not supported through devgita")
}

func (lg *LazyGit) ForceConfigure() error {
	gc := &config.GlobalConfig{}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	gc.EnableShellFeature(constants.LazyGit)
	if err := gc.RegenerateShellConfig(); err != nil {
		return fmt.Errorf("failed to generate shell config: %w", err)
	}
	if err := gc.Save(); err != nil {
		return fmt.Errorf("failed to save global config: %w", err)
	}
	return nil
}

func (lg *LazyGit) SoftConfigure() error {
	gc := &config.GlobalConfig{}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	if gc.IsShellFeatureEnabled(constants.LazyGit) {
		return nil
	}
	return lg.ForceConfigure()
}

func (lg *LazyGit) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		Command: constants.LazyGit,
		Args:    args,
	}
	if _, _, err := lg.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run lazygit command: %w", err)
	}
	return nil
}

func (lg *LazyGit) Update() error {
	return fmt.Errorf("lazygit update not implemented through devgita")
}
