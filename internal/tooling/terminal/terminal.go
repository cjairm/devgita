package terminal

import (
	"fmt"

	"github.com/cjairm/devgita/internal/apps/fastfetch"
	"github.com/cjairm/devgita/internal/apps/lazydocker"
	"github.com/cjairm/devgita/internal/apps/lazygit"
	"github.com/cjairm/devgita/internal/apps/mise"
	"github.com/cjairm/devgita/internal/apps/neovim"
	"github.com/cjairm/devgita/internal/apps/opencode"
	"github.com/cjairm/devgita/internal/apps/tmux"
	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/tooling/terminal/core/autoconf"
	"github.com/cjairm/devgita/internal/tooling/terminal/core/bison"
	"github.com/cjairm/devgita/internal/tooling/terminal/core/gdbm"
	"github.com/cjairm/devgita/internal/tooling/terminal/core/jemalloc"
	"github.com/cjairm/devgita/internal/tooling/terminal/core/libffi"
	"github.com/cjairm/devgita/internal/tooling/terminal/core/libyaml"
	"github.com/cjairm/devgita/internal/tooling/terminal/core/mupdf"
	"github.com/cjairm/devgita/internal/tooling/terminal/core/ncurses"
	"github.com/cjairm/devgita/internal/tooling/terminal/core/openssl"
	"github.com/cjairm/devgita/internal/tooling/terminal/core/pkgconfig"
	"github.com/cjairm/devgita/internal/tooling/terminal/core/readline"
	"github.com/cjairm/devgita/internal/tooling/terminal/core/unzip"
	"github.com/cjairm/devgita/internal/tooling/terminal/core/vips"
	"github.com/cjairm/devgita/internal/tooling/terminal/core/xcode"
	"github.com/cjairm/devgita/internal/tooling/terminal/core/zlib"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/autosuggestions"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/bat"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/btop"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/curl"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/eza"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/fdfind"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/fzf"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/githubcli"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/powerlevel10k"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/ripgrep"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/syntaxhighlighting"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/tldr"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/zoxide"
	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/logger"
	"github.com/cjairm/devgita/pkg/paths"
	"github.com/cjairm/devgita/pkg/promptui"
	"github.com/cjairm/devgita/pkg/utils"
)

type Terminal struct {
	Cmd  commands.Command
	Base commands.BaseCommand
}

func New() *Terminal {
	osCmd := commands.NewCommand()
	baseCmd := commands.NewBaseCommand()
	return &Terminal{Cmd: osCmd, Base: *baseCmd}
}

func (t *Terminal) InstallAndConfigure() {
	err := t.DisplayGithubInstructions()
	displayMessage(err, "instructions", true)
	if t.Base.Platform.IsMac() {
		xc := xcode.New()
		displayMessage(xc.SoftInstall(), "xcode")
	}
	t.InstallTerminalApps()
	t.InstallDevTools()
	t.InstallCoreLibs()

	if _, _, err := t.Base.ExecCommand(commands.CommandParams{
		Command: "source",
		Args:    []string{paths.Files.ShellConfig},
	}); err != nil {
		utils.PrintWarning(fmt.Sprintf(
			"Failed to source %s: %v",
			paths.Files.ShellConfig, err))
	}
}

func (t *Terminal) InstallTerminalApps() {
	// should install:
	// - fastfetch
	// - lazydocker
	// - lazygit
	// - mise
	// - neovim
	// - opencode
	// - tmux
	terminalApps := []struct {
		name string
		app  interface {
			SoftInstall() error
			SoftConfigure() error
		}
	}{
		{constants.Fastfetch, fastfetch.New()},
		{constants.Neovim, neovim.New()},
		{constants.Tmux, tmux.New()},
	}
	for _, terminalApp := range terminalApps {
		if err := terminalApp.app.SoftInstall(); err != nil {
			displayMessage(err, terminalApp.name)
			continue
		}
		if err := terminalApp.app.SoftConfigure(); err != nil {
			displayMessage(err, terminalApp.name, true)
			continue
		}
	}

	tuis := []struct {
		name string
		app  interface{ SoftInstall() error }
	}{
		{constants.LazyDocker, lazydocker.New()},
		{constants.LazyGit, lazygit.New()},
	}
	for _, tui := range tuis {
		displayMessage(tui.app.SoftInstall(), tui.name)
	}

	m := mise.New()
	displayMessage(m.SoftInstall(), constants.Mise)

	o := opencode.New()
	displayMessage(o.SoftInstall(), constants.OpenCode)
}

func (t *Terminal) InstallDevTools() {
	// should install bat, btop, curl, eza, fd-find, fzf, gh, powerlevel10k,
	// ripgrep, syntaxhighlighting, tldr, zoxide, zsh-autosuggestions
	devtools := []struct {
		name string
		app  interface{ SoftInstall() error }
	}{
		{constants.Bat, bat.New()},
		{constants.Btop, btop.New()},
		{constants.Curl, curl.New()},
		{constants.Eza, eza.New()},
		{constants.FdFind, fdfind.New()},
		{constants.Fzf, fzf.New()},
		{constants.GithubCli, githubcli.New()},
		{constants.Powerlevel10k, powerlevel10k.New()},
		{constants.Ripgrep, ripgrep.New()},
		{constants.Syntaxhighlighting, syntaxhighlighting.New()},
		{constants.Tldr, tldr.New()},
		{constants.Zoxide, zoxide.New()},
		{constants.ZshAutosuggestions, autosuggestions.New()},
	}
	for _, devtool := range devtools {
		displayMessage(devtool.app.SoftInstall(), devtool.name)
	}
}

func (t *Terminal) InstallCoreLibs() {
	// installs libs pkg-config, autoconf, bison, openssl, readline, zlib,
	// libyaml, ncurses, libffi, gdbm, jemalloc, vips, mupdf, unzip
	libs := []struct {
		name string
		app  interface{ SoftInstall() error }
	}{
		{constants.PkgConfig, pkgconfig.New()},
		{constants.Autoconf, autoconf.New()},
		{constants.Bison, bison.New()},
		{constants.OpenSSL, openssl.New()},
		{constants.Readline, readline.New()},
		{constants.Zlib, zlib.New()},
		{constants.Libyaml, libyaml.New()},
		{constants.Ncurses, ncurses.New()},
		{constants.Libffi, libffi.New()},
		{constants.Gdbm, gdbm.New()},
		{constants.Jemalloc, jemalloc.New()},
		{constants.Vips, vips.New()},
		{constants.Mupdf, mupdf.New()},
		{constants.Unzip, unzip.New()},
	}
	for _, lib := range libs {
		displayMessage(lib.app.SoftInstall(), lib.name)
	}
}

func (t *Terminal) DisplayGithubInstructions() error {
	instructions := `
		1. Generate a new SSH key:
		   ssh-keygen -t rsa -b 4096
		2. Start the SSH agent:
		   eval "$(ssh-agent -s)"
		3. Add your SSH key to the agent:
		   ssh-add -K ~/.ssh/id_rsa
		4. Copy the SSH key to your clipboard:
		   pbcopy < ~/.ssh/id_rsa.pub
		5. Go to github.com and store it.
		6. To test it out, run:
		   ssh -T git@github.com

		See documentation: https://docs.github.com/en/authentication/connecting-to-github-with-ssh/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent
	`
	return promptui.DisplayInstructions(
		"Before continue, it'd be nice if you get access to your GitHub repositories",
		instructions,
		false,
	)
}

func displayMessage(err error, name string, displayOnlyErrors ...bool) {
	if err != nil {
		logger.L().Errorw("Error installing ", "package_name", name, "error", err)
		utils.PrintWarning(
			fmt.Sprintf(
				"Install (%s) errored... To halt the installation, press ctrl+c or use --debug flag to see more details",
				name,
			),
		)
	} else {
		if displayOnlyErrors != nil && displayOnlyErrors[0] == true {
			return
		}
		msg := fmt.Sprintf("Installing %s (if no previously installed)...", name)
		utils.PrintInfo(msg)
	}
}
