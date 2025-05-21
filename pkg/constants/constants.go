package constants

import "fmt"

const (
	Reset  = "\033[0m"
	Gray   = "\033[90m" // Muted color
	Red    = "\033[31m" // Error color
	Green  = "\033[32m" // Success color
	Blue   = "\033[34m" // Informative color
	Yellow = "\033[33m" // Warning color
	Bold   = "\033[1m"

	// Supported Mac version
	SupportedMacOSVersionNumber = 13
	SupportedMacOSVersionName   = "Ventura"

	// Supported Debian version
	SupportedDebianVersionNumber = 12
	SupportedDebianVersionName   = "Bookworm"

	// Supported Ubuntu version
	SupportedUbuntuVersionNumber = 24
	SupportedUbuntuVersionName   = "Noble Numbat"

	// Supported Neovim version
	NeovimVersion = "0.11.1"

	// App specific
	AppName              = "devgita"
	ConfigAppDirName     = "configs"
	DevgitaRepositoryUrl = "https://github.com/cjairm/devgita.git"
	GlobalConfigFile     = "global_config.yaml"

	// App names
	Aerospace = "aerospace"
	Alacritty = "alacritty"
	Fastfetch = "fastfetch"
	Neovim    = "neovim"
	Nvim      = "nvim"
	Tmux      = "tmux"

	// Other available folders for installation
	Bash   = "bash"
	Fonts  = "fonts"
	Themes = "themes"
)

var Devgita = fmt.Sprint(`
    .___                .__  __          
  __| _/_______  ______ |__|/  |______   
 / __ |/ __ \  \/ / ___\|  \   __\__  \  
/ /_/ \  ___/\   / /_/  >  ||  |  / __ \_
\____ |\___  >\_/\___  /|__||__| (____  /
     \/    \/   /_____/               \/ 
@cjairm
`)
