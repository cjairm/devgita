package constants

import "fmt"

const (
	AppName                    = "devgita"
	Reset                       = "\033[0m"
	Gray                        = "\033[90m"
	Red                         = "\033[31m" // Error color
	Green                       = "\033[32m" // Success color
	Blue                        = "\033[34m" // Informative color
	Yellow                      = "\033[33m" // Warning color
	Bold                        = "\033[1m"
	SupportedMacOSVersionNumber = 13
	SupportedMacOSVersionName   = "Ventura"
	DevgitaRepositoryUrl        = "https://github.com/cjairm/devgita.git"

	// Paths

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
