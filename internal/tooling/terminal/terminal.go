package terminal

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cjairm/devgita/internal/apps/fastfetch"
	"github.com/cjairm/devgita/internal/apps/lazydocker"
	"github.com/cjairm/devgita/internal/apps/lazygit"
	"github.com/cjairm/devgita/internal/apps/mise"
	"github.com/cjairm/devgita/internal/apps/neovim"
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
	"github.com/cjairm/devgita/pkg/files"
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

func (t *Terminal) InstallAll() {
	err := t.DisplayGithubInstructions()
	displayMessage(err, "instructions", true)
	if t.Base.Platform.IsMac() {
		displayMessage(t.InstallXCode(), "xcode")
	}
	t.InstallTerminalApps()
	t.InstallDevTools()
	t.InstallCoreLibs()
}

func (t *Terminal) ConfigureZsh() error {
	var err error

	utils.PrintInfo("Adding config custom files...")
	isDevgitaConfigFilePresent := files.FileAlreadyExist(
		filepath.Join(paths.AppDir, "devgita.zsh"),
	)
	if !isDevgitaConfigFilePresent {
		err = files.CopyDir(
			paths.BashConfigAppDir,
			filepath.Join(paths.ConfigDir, constants.AppName),
		)
		if err != nil {
			return err
		}
	}

	utils.PrintInfo("Installing terminal theme...")
	p := powerlevel10k.New()
	err = p.SoftInstall()
	if err != nil {
		return err
	}
	err = p.SoftConfigure()
	if err != nil {
		return err
	}

	utils.PrintInfo("Installing zsh-autosuggestions...")
	za := autosuggestions.New()
	err = za.SoftInstall()
	if err != nil {
		return err
	}
	err = za.SoftConfigure()
	if err != nil {
		return err
	}

	utils.PrintInfo("Installing zsh-syntax-highlighting...")
	sh := syntaxhighlighting.New()
	err = sh.SoftInstall()
	if err != nil {
		return err
	}
	err = sh.SoftConfigure()
	if err != nil {
		return err
	}

	utils.PrintInfo("Sourcing custom files...")
	// TODO: Create a fuction that builds these strings. If $HOME exists, use it
	// if not, use the full path
	err = t.Base.MaybeSetup("source $HOME/.config/devgita/aliases.zsh", "aliases.zsh")
	if err != nil {
		return err
	}
	err = t.Base.MaybeSetup("source $HOME/.config/devgita/init.zsh", "init.zsh")
	if err != nil {
		return err
	}
	return t.Base.MaybeSetup("source $HOME/.config/devgita/devgita.zsh", "devgita.zsh")
}

func (t *Terminal) InstallTerminalApps() {
	// should install lazydocker, lazygit, mise, fastfetch, neovim, tmux
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
}

func (t *Terminal) InstallDevTools() {
	// should install curl, fzf, ripgrep, bat, eza, gh, zoxide, btop, fd-find, tldr
	devtools := []struct {
		name string
		app  interface{ SoftInstall() error }
	}{
		{constants.Curl, curl.New()},
		{constants.GithubCli, githubcli.New()},
		{constants.Fzf, fzf.New()},
		{constants.Ripgrep, ripgrep.New()},
		{constants.Bat, bat.New()},
		{constants.Eza, eza.New()},
		{constants.Zoxide, zoxide.New()},
		{constants.Btop, btop.New()},
		{constants.FdFind, fdfind.New()},
		{constants.Tldr, tldr.New()},
	}
	for _, devtool := range devtools {
		displayMessage(devtool.app.SoftInstall(), devtool.name)
	}
}

func (t *Terminal) InstallCoreLibs() {
	// installs libs pkg-config, autoconf, bison, openssl, readline, zlib,
	//               libyaml, ncurses, libffi, gdbm, jemalloc, vips, mupdf,
	//	         unzip
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

// This dedicated to MacOS
func (t *Terminal) InstallXCode() error {
	isInstalled, err := isXcodeInstalled()
	if err != nil {
		return err
	}
	if isInstalled {
		return nil
	}
	cmd := commands.CommandParams{
		PreExecMsg:  "Installing Xcode Command Line Tools",
		PostExecMsg: "",
		IsSudo:      false,
		Command:     "xcode-select",
		Args:        []string{"--install"},
	}
	if _, _, err := t.Base.ExecCommand(cmd); err != nil {
		return fmt.Errorf("failed to install Xcode Command Line Tools: %w", err)
	}
	return nil
}

func isXcodeInstalled() (bool, error) {
	cmd := exec.Command("xcode-select", "-p")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("error running xcode-select: %v", err)
	}
	// Check if the output contains the expected path for Xcode
	xcodePath := strings.ToLower(strings.TrimSpace(out.String()))
	if strings.Contains(xcodePath, "xcode.app") || strings.Contains(xcodePath, "commandlinetools") {
		return true, nil
	}
	return false, nil
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
