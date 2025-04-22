package constants

import "fmt"

const (
	Reset                       = "\033[0m"
	Gray                        = "\033[90m"
	Red                         = "\033[31m" // Error color
	Green                       = "\033[32m" // Success color
	Blue                        = "\033[34m" // Informative color
	Yellow                      = "\033[33m" // Warning color
	Bold                        = "\033[1m"
	SupportedMacOSVersionNumber = 13
	SupportedMacOSVersionName   = "Ventura"

	// App specific
	AppName              = "devgita"
	ConfigAppDirName     = "configs"
	DevgitaRepositoryUrl = "https://github.com/cjairm/devgita.git"

	// App names
	Aerospace = "aerospace"
	Alacritty = "alacritty"
	Fastfetch = "fastfetch"
	Neovim    = "neovim"
	Nvim    = "nvim"
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
