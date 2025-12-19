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

// Paths contains all directory path structures
var Paths = struct {
	App struct {
		Root    string
		Configs struct {
			Aerospace string
			Alacritty string
			Bash      string
			Fastfetch string
			Fonts     string
			Git       string
			Neovim    string
			Templates string
			Themes    string
			Tmux      string
		}
		Fonts struct {
			Alacritty string
		}
		Themes struct {
			Alacritty string
		}
	}
	Cache struct {
		Root string
	}
	Config struct {
		Root      string
		Aerospace string
		Alacritty string
		Bash      string
		Fastfetch string
		Fonts     string
		Git       string
		Nvim      string
		Themes    string
		Tmux      string
	}
	Data struct {
		Root string
	}
	Home struct {
		Root string
	}
	System struct {
		Applications string
		Fonts        string
	}
	User struct {
		Applications string
		Fonts        string
	}
}{
	App: struct {
		Root    string
		Configs struct {
			Aerospace string
			Alacritty string
			Bash      string
			Fastfetch string
			Fonts     string
			Git       string
			Neovim    string
			Templates string
			Themes    string
			Tmux      string
		}
		Fonts struct {
			Alacritty string
		}
		Themes struct {
			Alacritty string
		}
	}{
		Root: GetDataDir(constants.App.Name),
		Configs: struct {
			Aerospace string
			Alacritty string
			Bash      string
			Fastfetch string
			Fonts     string
			Git       string
			Neovim    string
			Templates string
			Themes    string
			Tmux      string
		}{
			Aerospace: GetAppDir(constants.App.Dir.Configs, constants.Aerospace),
			Alacritty: GetAppDir(constants.App.Dir.Configs, constants.Alacritty),
			Bash:      GetAppDir(constants.App.Dir.Configs, constants.Bash),
			Fastfetch: GetAppDir(constants.App.Dir.Configs, constants.Fastfetch),
			Fonts:     GetAppDir(constants.App.Dir.Configs, constants.Fonts),
			Git:       GetAppDir(constants.App.Dir.Configs, constants.Git),
			Neovim:    GetAppDir(constants.App.Dir.Configs, constants.Neovim),
			Templates: GetAppDir(constants.App.Dir.Configs, constants.Templates),
			Themes:    GetAppDir(constants.App.Dir.Configs, constants.Themes),
			Tmux:      GetAppDir(constants.App.Dir.Configs, constants.Tmux),
		},
		Fonts: struct {
			Alacritty string
		}{
			Alacritty: GetAppDir(constants.App.Dir.Configs, constants.Fonts, constants.Alacritty),
		},
		Themes: struct {
			Alacritty string
		}{
			Alacritty: GetAppDir(constants.App.Dir.Configs, constants.Themes, constants.Alacritty),
		},
	},
	Cache: struct {
		Root string
	}{
		Root: GetCacheDir(),
	},
	Config: struct {
		Root      string
		Aerospace string
		Alacritty string
		Bash      string
		Fastfetch string
		Fonts     string
		Git       string
		Nvim      string
		Themes    string
		Tmux      string
	}{
		Root:      GetConfigDir(),
		Aerospace: GetConfigDir(constants.Aerospace),
		Alacritty: GetConfigDir(constants.Alacritty),
		Bash:      GetConfigDir(constants.Bash),
		Fastfetch: GetConfigDir(constants.Fastfetch),
		Fonts:     GetConfigDir(constants.Fonts),
		Git:       GetConfigDir(constants.Git),
		Nvim:      GetConfigDir(constants.Nvim),
		Themes:    GetConfigDir(constants.Themes),
		Tmux:      GetConfigDir(constants.Tmux),
	},
	Data: struct {
		Root string
	}{
		Root: GetDataDir(),
	},
	Home: struct {
		Root string
	}{
		Root: GetHomeDir(),
	},
	System: struct {
		Applications string
		Fonts        string
	}{
		Applications: GetSystemApplicationsDir(runtime.GOOS == "darwin"),
		Fonts:        GetSystemFontsDir(runtime.GOOS == "darwin"),
	},
	User: struct {
		Applications string
		Fonts        string
	}{
		Applications: GetUserApplicationsDir(runtime.GOOS == "darwin"),
		Fonts:        GetUserFontsDir(runtime.GOOS == "darwin"),
	},
}

// Files contains all file path structures
var Files = struct {
	ShellConfig string
}{
	ShellConfig: GetShellConfigFile(),
}

// Public API functions

// Returns XDG_CONFIG_HOME or fallback to ~/.config
// Reads environment variables dynamically to support testing
func GetConfigDir(subPath ...string) string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic("could not determine home directory")
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(append([]string{base}, subPath...)...)
}

// Returns XDG_DATA_HOME or fallback to ~/.local/share
// Reads environment variables dynamically to support testing
func GetDataDir(subPath ...string) string {
	base := os.Getenv("XDG_DATA_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic("could not determine home directory")
		}
		base = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(append([]string{base}, subPath...)...)
}

func GetAppDir(subPath ...string) string {
	appRoot := GetDataDir(constants.App.Name)
	return filepath.Join(append([]string{appRoot}, subPath...)...)
}

func GetHomeDir(subPath ...string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("could not determine home directory")
	}
	return filepath.Join(append([]string{home}, subPath...)...)
}

// Returns XDG_CACHE_HOME or fallback to ~/.cache
// Reads environment variables dynamically to support testing
func GetCacheDir(subPath ...string) string {
	base := os.Getenv("XDG_CACHE_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic("could not determine home directory")
		}
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
	home, err := os.UserHomeDir()
	if err != nil {
		panic("could not determine home directory")
	}
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		configDir = filepath.Join(home, ".config")
	}
	shellConfigFiles := []string{
		filepath.Join(home, ".zshrc"),
		filepath.Join(home, ".bashrc"),
		filepath.Join(home, ".bash_profile"),
		filepath.Join(home, ".profile"),
		filepath.Join(configDir, "fish", "config.fish"),
	}
	for _, filepath := range shellConfigFiles {
		if FileAlreadyExist(filepath) {
			return filepath
		}
	}
	// If none exist, default to .zshrc
	return filepath.Join(home, ".zshrc")
}

// Returns user-level fonts dir
func GetUserFontsDir(isMac bool, subPath ...string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("could not determine home directory")
	}

	if isMac {
		base := filepath.Join(home, "Library", "Fonts")
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
