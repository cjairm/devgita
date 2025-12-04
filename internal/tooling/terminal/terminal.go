package terminal

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cjairm/devgita/internal/apps/autosuggestions"
	"github.com/cjairm/devgita/internal/apps/curl"
	"github.com/cjairm/devgita/internal/apps/fastfetch"
	"github.com/cjairm/devgita/internal/apps/githubcli"
	"github.com/cjairm/devgita/internal/apps/lazydocker"
	"github.com/cjairm/devgita/internal/apps/lazygit"
	"github.com/cjairm/devgita/internal/apps/mise"
	"github.com/cjairm/devgita/internal/apps/neovim"
	"github.com/cjairm/devgita/internal/apps/powerlevel10k"
	"github.com/cjairm/devgita/internal/apps/syntaxhighlighting"
	"github.com/cjairm/devgita/internal/apps/tmux"
	"github.com/cjairm/devgita/internal/apps/unzip"
	"github.com/cjairm/devgita/internal/commands"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/bat"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/eza"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/fzf"
	"github.com/cjairm/devgita/internal/tooling/terminal/dev_tools/ripgrep"
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

func (t *Terminal) InstallAll() error {
	err := t.DisplayGithubInstructions()
	ifErrorDisplayMessage(err, "instructions")

	utils.PrintInfo("Installing curl (if no previously installed)...")
	c := curl.New()
	ifErrorDisplayMessage(c.SoftInstall(), "curl")

	utils.PrintInfo("Installing and setting up fastfetch (if no previous configuration)...")
	err = t.InstallFastFetch()
	ifErrorDisplayMessage(err, "fastfetch")

	utils.PrintInfo("Installing unzip (if no previously installed)...")
	uz := unzip.New()
	ifErrorDisplayMessage(uz.SoftInstall(), "unzip")

	utils.PrintInfo("Installing gh (if no previously installed)...")
	gh := githubcli.New()
	ifErrorDisplayMessage(gh.SoftInstall(), "gh")

	utils.PrintInfo("Installing lazydocker (if no previously installed)...")
	ld := lazydocker.New()
	ifErrorDisplayMessage(ld.SoftInstall(), "lazydocker")

	utils.PrintInfo("Installing lazygit (if no previously installed)...")
	lg := lazygit.New()
	ifErrorDisplayMessage(lg.SoftInstall(), "lazygit")

	utils.PrintInfo("Installing and setting up neovim (if no previous configuration)...")
	err = t.InstallNeovim()
	ifErrorDisplayMessage(err, "neovim")

	utils.PrintInfo("Installing and setting up tmux (if no previous configuration)...")
	err = t.InstallTmux()
	ifErrorDisplayMessage(err, "tmux")

	utils.PrintInfo(
		"Installing fzf, ripgrep, bat, eza, zoxide, btop, fd-find, and tldr (if no previously installed)...",
	)
	err = t.InstallDevTools()
	ifErrorDisplayMessage(err, "fzf, ripgrep, bat, eza, zoxide, btop, fd-find, and tldr")

	if t.Base.Platform.IsMac() {
		utils.PrintInfo("Installing xcode (if no previously installed)...")
		err = t.InstallXCode()
		ifErrorDisplayMessage(err, "xcode")
	}

	utils.PrintInfo(
		"Installing pkg-config, autoconf, bison, rust, openssl, readline, zlib, libyaml, ncurses, libffi, gdbm, jemalloc, vips, and mupdf (if no previously installed)...",
	)
	err = t.InstallCoreLibs()
	ifErrorDisplayMessage(
		err,
		"pkg-config, autoconf, bison, rust, openssl, readline, zlib, libyaml, ncurses, libffi, gdbm, jemalloc, vips, and mupdf",
	)

	utils.PrintInfo("Installing mise (if no previously installed)...")
	m := mise.New()
	ifErrorDisplayMessage(m.SoftInstall(), "mise")

	return nil
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

func (t *Terminal) InstallFastFetch() error {
	ff := fastfetch.New()
	if err := ff.SoftInstall(); err != nil {
		return err
	}
	if err := ff.SoftConfigure(); err != nil {
		return err
	}
	return nil
}

func (t *Terminal) InstallNeovim() error {
	nv := neovim.New()
	err := nv.SoftInstall()
	if err != nil {
		return err
	}
	err = nv.SoftConfigure()
	if err != nil {
		return err
	}
	return nil
}

func (t *Terminal) InstallTmux() error {
	tm := tmux.New()
	err := tm.SoftInstall()
	if err != nil {
		return err
	}
	err = tm.SoftConfigure()
	if err != nil {
		return err
	}
	return nil
}

// installs fzf, ripgrep, bat, eza, zoxide, btop, fd-find, tldr
func (t *Terminal) InstallDevTools() error {
	devtools := []struct {
		name string
		app  interface{ SoftInstall() error }
	}{
		{constants.Fzf, fzf.New()},
		{constants.Ripgrep, ripgrep.New()},
		{constants.Bat, bat.New()},
		{constants.Eza, eza.New()},
	}
	for _, devtool := range devtools {
		ifErrorDisplayMessage(devtool.app.SoftInstall(), devtool.name)
	}

	packages := []string{"zoxide", "btop", "fd", "tldr"}
	for _, pkg := range packages {
		if err := t.Cmd.MaybeInstallPackage(pkg); err != nil {
			return err
		}
	}
	return nil
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

// installs libs pkg-config, autoconf, bison, rust, openssl, readline, zlib, libyaml, ncurses, libffi, gdbm, jemalloc, vips, mupdf
func (t *Terminal) InstallCoreLibs() error {
	libs := []string{
		"pkg-config",
		"autoconf",
		"bison",
		"rust",
		"openssl",
		"readline",
		"zlib",
		"libyaml",
		"ncurses",
		"libffi",
		"gdbm",
		"jemalloc",
		"vips",
		"mupdf",
	}
	for _, lib := range libs {
		if err := t.Cmd.MaybeInstallPackage(lib); err != nil {
			return err
		}
	}
	return nil
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

func ifErrorDisplayMessage(err error, packageName string) {
	if err != nil {
		logger.L().Errorw("Error installing ", "package_name", packageName, "error", err)
		utils.PrintWarning(
			fmt.Sprintf(
				"Install (%s) errored... To halt the installation, press ctrl+c or use --debug flag to see more details",
				packageName,
			),
		)
	}
}
