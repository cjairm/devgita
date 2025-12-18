package constants

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
	Aerospace          = "aerospace"
	Alacritty          = "alacritty"
	Autoconf           = "autoconf"
	Bison              = "bison"
	Bat                = "bat"
	Btop               = "btop"
	Curl               = "curl"
	Eza                = "eza"
	Fastfetch          = "fastfetch"
	FdFind             = "fd" // Equivalent to "fd-find" in Linux package managers
	Fzf                = "fzf"
	Gdbm               = "gdbm"
	Git                = "git"
	GithubCli          = "gh"
	Jemalloc           = "jemalloc"
	LazyDocker         = "lazydocker"
	LazyGit            = "lazygit"
	Libffi             = "libffi"
	Libyaml            = "libyaml"
	Mupdf              = "mupdf"
	Ncurses            = "ncurses"
	Neovim             = "neovim"
	Mise               = "mise"
	Nvim               = "nvim"
	OpenSSL            = "openssl"
	PkgConfig          = "pkg-config"
	Powerlevel10k      = "powerlevel10k"
	Readline           = "readline"
	Ripgrep            = "rg"
	Syntaxhighlighting = "zsh-syntax-highlighting"
	Tldr               = "tldr"
	Tmux               = "tmux"
	Unzip              = "unzip"
	Vips               = "vips"
	Zlib               = "zlib"
	Zoxide             = "zoxide"
	ZshAutosuggestions = "zsh-autosuggestions"

	// Other available folders for installation
	Bash      = "bash"
	Fonts     = "fonts"
	Themes    = "themes"
	Templates = "templates"

	// Template files
	DevgitaShellTemplate = "devgita.zsh.tmpl"
)

var Devgita = `
    .___                .__  __          
  __| _/_______  ______ |__|/  |______   
 / __ |/ __ \  \/ / ___\|  \   __\__  \  
/ /_/ \  ___/\   / /_/  >  ||  |  / __ \_
\____ |\___  >\_/\___  /|__||__| (____  /
     \/    \/   /_____/               \/ 
@cjairm / https://cjairm.me/
`
