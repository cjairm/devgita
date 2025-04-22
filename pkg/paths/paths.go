package paths

import (
	"os"
	"path/filepath"

	"github.com/cjairm/devgita/pkg/constants"
)

var (
	AppDir    = GetAppDir()
	CacheDir  = GetCacheDir()
	ConfigDir = GetConfigDir()
	DataDir   = GetDataDir()
	HomeDir   = GetHomeDir()

	// Configs from Devgita app
	AerospaceConfigAppDir = GetAppDir(constants.ConfigAppDirName, constants.Aerospace)
	AlacrittyConfigAppDir = GetAppDir(constants.ConfigAppDirName, constants.Alacritty)
	BashConfigAppDir      = GetAppDir(constants.ConfigAppDirName, constants.Bash)
	FastFetchConfigAppDir = GetAppDir(constants.ConfigAppDirName, constants.Fastfetch)
	FontConfigAppDir      = GetAppDir(constants.ConfigAppDirName, constants.Fonts)
	NeovimConfigAppDir    = GetAppDir(constants.ConfigAppDirName, constants.Neovim)
	ThemesConfigAppDir    = GetAppDir(constants.ConfigAppDirName, constants.Themes)
	TmuxConfigAppDir      = GetAppDir(constants.ConfigAppDirName, constants.Tmux)

	// Config in local (usually `.config` folder)
	AerospaceConfigLocalDir = GetConfigDir(constants.Aerospace)
	AlacrittyConfigLocalDir = GetConfigDir(constants.Alacritty)
	BashConfigLocalDir      = GetConfigDir(constants.Bash)
	FastFetchConfigLocalDir = GetConfigDir(constants.Fastfetch)
	FontConfigLocalDir      = GetConfigDir(constants.Fonts)
	NvimConfigLocalDir      = GetConfigDir(constants.Nvim)
	ThemesConfigLocalDir    = GetConfigDir(constants.Themes)
	TmuxConfigLocalDir      = GetConfigDir(constants.Tmux)

	// Fonts
	AlacrittyFontAppDir = GetAppDir(
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
func GetConfigDir(subDirs ...string) string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home := GetHomeDir()
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(append([]string{base}, subDirs...)...)
}

// Returns XDG_DATA_HOME or fallback to ~/.local/share
func GetDataDir(subDirs ...string) string {
	base := os.Getenv("XDG_DATA_HOME")
	if base == "" {
		home := GetHomeDir()
		base = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(append([]string{base}, subDirs...)...)
}

func GetAppDir(subDirs ...string) string {
	appDir := GetDataDir(constants.AppName)
	return filepath.Join(append([]string{appDir}, subDirs...)...)
}

func GetHomeDir(subDirs ...string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("could not determine home directory")
	}
	return filepath.Join(append([]string{home}, subDirs...)...)
}

// Returns XDG_CACHE_HOME or fallback to ~/.cache
func GetCacheDir(subDirs ...string) string {
	base := os.Getenv("XDG_CACHE_HOME")
	if base == "" {
		home := GetHomeDir()
		base = filepath.Join(home, ".cache")
	}
	return filepath.Join(append([]string{base}, subDirs...)...)
}
