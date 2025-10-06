package paths

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
)

// This allows swapping it during tests
var FileAlreadyExist = files.FileAlreadyExist

var (
	AppDir          = GetAppDir()
	CacheDir        = GetCacheDir()
	ConfigDir       = GetConfigDir()
	DataDir         = GetDataDir()
	HomeDir         = GetHomeDir()
	ShellConfigFile = GetShellConfigFile()

	// System apps
	SystemApplicationsDir = GetSystemApplicationsDir(runtime.GOOS == "darwin")

	// System fonts
	SystemFontsDir = GetSystemFontsDir(runtime.GOOS == "darwin")

	// User apps
	UserApplicationsDir = GetUserApplicationsDir(runtime.GOOS == "darwin")

	// User fonts
	UserFontsDir = GetUserFontsDir(runtime.GOOS == "darwin")

	// Configs from Devgita app
	AerospaceConfigAppDir = GetAppDir(constants.ConfigAppDirName, constants.Aerospace)
	AlacrittyConfigAppDir = GetAppDir(constants.ConfigAppDirName, constants.Alacritty)
	BashConfigAppDir      = GetAppDir(constants.ConfigAppDirName, constants.Bash)
	FastFetchConfigAppDir = GetAppDir(constants.ConfigAppDirName, constants.Fastfetch)
	FontConfigAppDir      = GetAppDir(constants.ConfigAppDirName, constants.Fonts)
	GitConfigAppDir       = GetAppDir(constants.ConfigAppDirName, constants.Git)
	NeovimConfigAppDir    = GetAppDir(constants.ConfigAppDirName, constants.Neovim)
	ThemesConfigAppDir    = GetAppDir(constants.ConfigAppDirName, constants.Themes)
	TmuxConfigAppDir      = GetAppDir(constants.ConfigAppDirName, constants.Tmux)

	// Config in local (usually `.config` folder)
	AerospaceConfigLocalDir = GetConfigDir(constants.Aerospace)
	AlacrittyConfigLocalDir = GetConfigDir(constants.Alacritty)
	BashConfigLocalDir      = GetConfigDir(constants.Bash)
	FastFetchConfigLocalDir = GetConfigDir(constants.Fastfetch)
	FontConfigLocalDir      = GetConfigDir(constants.Fonts)
	GitConfigLocalDir       = GetConfigDir(constants.Git)
	NvimConfigLocalDir      = GetConfigDir(constants.Nvim)
	ThemesConfigLocalDir    = GetConfigDir(constants.Themes)
	TmuxConfigLocalDir      = GetConfigDir(constants.Tmux)

	// Fonts
	AlacrittyFontsAppDir = GetAppDir(
		constants.ConfigAppDirName,
		constants.Fonts,
		constants.Alacritty,
	)

	// Themes
	AlacrittyThemesAppDir = GetAppDir(
		constants.ConfigAppDirName,
		constants.Themes,
		constants.Alacritty,
	)
)

// Returns XDG_CONFIG_HOME or fallback to ~/.config
func GetConfigDir(subPath ...string) string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home := GetHomeDir()
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(append([]string{base}, subPath...)...)
}

// Returns XDG_DATA_HOME or fallback to ~/.local/share
func GetDataDir(subPath ...string) string {
	base := os.Getenv("XDG_DATA_HOME")
	if base == "" {
		home := GetHomeDir()
		base = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(append([]string{base}, subPath...)...)
}

func GetAppDir(subPath ...string) string {
	appDir := GetDataDir(constants.AppName)
	return filepath.Join(append([]string{appDir}, subPath...)...)
}

func GetHomeDir(subPath ...string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("could not determine home directory")
	}
	return filepath.Join(append([]string{home}, subPath...)...)
}

// Returns XDG_CACHE_HOME or fallback to ~/.cache
func GetCacheDir(subPath ...string) string {
	base := os.Getenv("XDG_CACHE_HOME")
	if base == "" {
		home := GetHomeDir()
		base = filepath.Join(home, ".cache")
	}
	return filepath.Join(append([]string{base}, subPath...)...)
}

// Returns user-level applications dir
func GetUserApplicationsDir(isMac bool, subPath ...string) string {
	if isMac {
		base := "/Applications"
		return filepath.Join(append([]string{base}, subPath...)...)
	}
	// Linux (XDG-compliant user apps)
	return GetDataDir(append([]string{"applications"}, subPath...)...)
}

// Returns system-level applications dir
func GetSystemApplicationsDir(isMac bool, subPath ...string) string {
	if isMac {
		base := "/Applications"
		return filepath.Join(append([]string{base}, subPath...)...)
	}
	// Linux system-wide applications dirs
	// NOTE: /usr/share/applications is more common, but /usr/local/share/applications is also valid
	// You could return both or let the caller choose
	base := "/usr/share/applications"
	return filepath.Join(append([]string{base}, subPath...)...)
}

func GetShellConfigFile() string {
	shellConfigFiles := []string{
		filepath.Join(HomeDir, ".zshrc"),
		filepath.Join(HomeDir, ".bashrc"),
		filepath.Join(HomeDir, ".bash_profile"),
		filepath.Join(HomeDir, ".profile"),
		filepath.Join(ConfigDir, "fish", "config.fish"),
	}
	for _, filepath := range shellConfigFiles {
		if FileAlreadyExist(filepath) {
			return filepath
		}
	}
	// If none exist, default to .zshrc
	return filepath.Join(HomeDir, ".zshrc")
}

// Returns user-level fonts dir
func GetUserFontsDir(isMac bool, subPath ...string) string {
	if isMac {
		base := filepath.Join(HomeDir, "Library", "Fonts")
		return filepath.Join(append([]string{base}, subPath...)...)
	}
	// Linux user fonts (XDG-compliant)
	return GetDataDir(append([]string{"fonts"}, subPath...)...)
}

// Returns system-level fonts dir
func GetSystemFontsDir(isMac bool, subPath ...string) string {
	if isMac {
		base := filepath.Join("/Library", "Fonts")
		return filepath.Join(append([]string{base}, subPath...)...)
	}
	// Linux system fonts (common default)
	base := "/usr/share/fonts"
	return filepath.Join(append([]string{base}, subPath...)...)
}
