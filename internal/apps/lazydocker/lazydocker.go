// LazyDocker terminal UI for Docker container and image management with devgita integration
//
// LazyDocker is a simple terminal UI for both docker and docker-compose, written in Go
// with the gocui library. It provides an interactive interface to manage Docker containers,
// images, volumes, and networks, all from the comfort of the terminal.
//
// References:
// - LazyDocker Repository: https://github.com/jesseduffield/lazydocker
// - LazyDocker Documentation: https://github.com/jesseduffield/lazydocker/blob/master/docs/Config.md
//
// Common lazydocker commands available through ExecuteCommand():
//   - lazydocker - Launch interactive TUI
//   - lazydocker --version - Show lazydocker version information
//   - lazydocker --config - Show configuration file path
//   - lazydocker --help - Display help information

package lazydocker

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

type LazyDocker struct {
	Cmd          cmd.Command
	Base         cmd.BaseCommandExecutor
	fetchVersion func(owner, repo string) (string, error)                                      // injectable for tests
	downloadFn   func(ctx context.Context, url, dest string, cfg downloader.RetryConfig) error // injectable for tests
}

func (ld *LazyDocker) getVersion(owner, repo string) (string, error) {
	if ld.fetchVersion != nil {
		return ld.fetchVersion(owner, repo)
	}
	return gh.FetchLatestRelease(owner, repo)
}

var packageName = fmt.Sprintf("jesseduffield/%s/%s", constants.LazyDocker, constants.LazyDocker)

func New() *LazyDocker {
	osCmd := cmd.NewCommand()
	baseCmd := cmd.NewBaseCommand()
	return &LazyDocker{Cmd: osCmd, Base: baseCmd}
}

func (ld *LazyDocker) Install() error {
	if ld.Base.IsMac() {
		return ld.Cmd.InstallPackage(packageName)
	}
	return ld.installDebianLazydocker()
}

func (ld *LazyDocker) SoftInstall() error {
	if ld.Base.IsMac() {
		return ld.Cmd.MaybeInstallPackage(packageName, constants.LazyDocker)
	}
	// On Debian: lazydocker is not in apt — skip if already in PATH, otherwise install from GitHub
	if _, err := cmd.LookPathFn(constants.LazyDocker); err == nil {
		logger.L().Infow("lazydocker already installed, skipping")
		return nil
	}
	return ld.installDebianLazydocker()
}

// linuxArch returns the architecture string used in lazydocker release artifact names.
func linuxArch() string {
	if runtime.GOARCH == "arm64" {
		return "arm64"
	}
	return "x86_64"
}

// installDebianLazydocker fetches the latest lazydocker release from GitHub, downloads
// the Linux tar.gz for the current architecture, and installs the binary to /usr/local/bin
func (ld *LazyDocker) installDebianLazydocker() error {
	version, err := ld.getVersion("jesseduffield", "lazydocker")
	if err != nil {
		return fmt.Errorf("failed to fetch lazydocker version: %w", err)
	}

	url := fmt.Sprintf(
		"https://github.com/jesseduffield/lazydocker/releases/download/v%s/lazydocker_%s_Linux_%s.tar.gz",
		version, version, linuxArch(),
	)
	logger.L().Infow("Downloading lazydocker for Debian", "version", version, "url", url)

	if err := cmd.InstallGitHubBinary(ld.Base, constants.LazyDocker, url, ld.downloadFn); err != nil {
		return fmt.Errorf("lazydocker installation failed: %w", err)
	}

	logger.L().Infow("lazydocker installed successfully for Debian", "version", version)
	return nil
}

func (ld *LazyDocker) ForceInstall() error {
	err := ld.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to uninstall lazydocker: %w", err)
	}
	return ld.Install()
}

func (ld *LazyDocker) Uninstall() error {
	return fmt.Errorf("lazydocker uninstall not supported through devgita")
}

func (ld *LazyDocker) ForceConfigure() error {
	gc := &config.GlobalConfig{}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	gc.EnableShellFeature(constants.LazyDocker)
	if err := gc.RegenerateShellConfig(); err != nil {
		return fmt.Errorf("failed to generate shell config: %w", err)
	}
	if err := gc.Save(); err != nil {
		return fmt.Errorf("failed to save global config: %w", err)
	}
	return nil
}

func (ld *LazyDocker) SoftConfigure() error {
	gc := &config.GlobalConfig{}
	if err := gc.Load(); err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}
	if gc.IsShellFeatureEnabled(constants.LazyDocker) {
		return nil
	}
	return ld.ForceConfigure()
}

func (ld *LazyDocker) ExecuteCommand(args ...string) error {
	execCommand := cmd.CommandParams{
		Command: constants.LazyDocker,
		Args:    args,
	}
	if _, _, err := ld.Base.ExecCommand(execCommand); err != nil {
		return fmt.Errorf("failed to run lazydocker command: %w", err)
	}
	return nil
}

func (ld *LazyDocker) Update() error {
	return fmt.Errorf("lazydocker update not implemented through devgita")
}

