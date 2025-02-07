package utils

import "fmt"

const (
	Reset = "\033[0m"
	Gray  = "\033[90m"
	Red   = "\033[31m"
	Green = "\033[32m"
	Bold  = "\033[1m"
)

var Devgita = fmt.Sprintf(`
%s
    .___                .__  __          
  __| _/_______  ______ |__|/  |______   
 / __ |/ __ \  \/ / ___\|  \   __\__  \  
/ /_/ \  ___/\   / /_/  >  ||  |  / __ \_
\____ |\___  >\_/\___  /|__||__| (____  /
     \/    \/   /_____/               \/ 
@cjairm
%s`, Bold, Reset)
