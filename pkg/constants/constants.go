package constants

const (
	Reset  = "\033[0m"
	Gray   = "\033[90m" // Muted color
	Red    = "\033[31m" // Error color
	Green  = "\033[32m" // Success color
	Blue   = "\033[34m" // Informative color
	Yellow = "\033[33m" // Warning color
	Bold   = "\033[1m"
)

type OsVersion struct {
	Number int
	Name   string
}

type AppVersion struct {
	Number string
}

// SupportedVersion contains version requirements for different platforms and applications
var SupportedVersion = struct {
	MacOS  OsVersion
	Debian OsVersion
	Ubuntu OsVersion
	Neovim AppVersion
}{
	MacOS: OsVersion{
		Number: 13,
		Name:   "Ventura",
	},
	Debian: OsVersion{
		Number: 12,
		Name:   "Bookworm",
	},
	Ubuntu: OsVersion{
		Number: 24,
		Name:   "Noble Numbat",
	},
	Neovim: AppVersion{
		Number: "0.11.1",
	},
}

// App contains application-specific constants
var App = struct {
	Name       string
	Repository struct {
		Name string
		URL  string
	}
	Dir struct {
		Configs string
	}
	File struct {
		GlobalConfig string
		ZshConfig    string
	}
	Template struct {
		ShellConfig string
	}
}{
	Name: "devgita",
	Repository: struct {
		Name string
		URL  string
	}{
		Name: "devgita",
		URL:  "https://github.com/cjairm/devgita.git",
	},
	Dir: struct {
		Configs string
	}{
		Configs: "configs",
	},
	File: struct {
		GlobalConfig string
		ZshConfig    string
	}{
		GlobalConfig: "global_config.yaml",
		ZshConfig:    ".devgita.zsh",
	},
	Template: struct {
		ShellConfig string
	}{
		ShellConfig: "devgita.zsh.tmpl",
	},
}

const (
	// App names
	Aerospace = "aerospace"
	Alacritty = "alacritty"
	Autoconf  = "autoconf"
	Bison     = "bison"
	Bat       = "bat"
	Brave     = "brave"
	Btop      = "btop"
	Curl      = "curl"
	// TODO: Rename this to `Devgita` and update current Devgita graphic to match the new name. (This is to avoid confusion between the app and the repository)
	DevgitaApp         = "devgita"
	Docker             = "docker"
	Eza                = "eza"
	Fastfetch          = "fastfetch"
	FdFind             = "fd" // Equivalent to "fd-find" in Linux package managers
	FontConfig         = "fontconfig"
	Flameshot          = "flameshot"
	Fzf                = "fzf"
	Gdbm               = "gdbm"
	Gimp               = "gimp"
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
	OpenCode           = "opencode"
	OpenSSL            = "openssl"
	PkgConfig          = "pkg-config"
	Powerlevel10k      = "powerlevel10k"
	Raycast            = "raycast"
	Readline           = "readline"
	Ripgrep            = "rg"
	Syntaxhighlighting = "zsh-syntax-highlighting"
	Tldr               = "tldr"
	Tmux               = "tmux"
	Unzip              = "unzip"
	Vips               = "vips"
	Xcode              = "xcode"
	Zlib               = "zlib"
	Zoxide             = "zoxide"
	ZshAutosuggestions = "zsh-autosuggestions"

	// Other available folders for installation
	Fonts     = "fonts"
	Themes    = "themes"
	Templates = "templates"

	// Programming languages (Mise core tools)
	Bun    = "bun"
	Deno   = "deno"
	Elixir = "elixir"
	Erlang = "erlang"
	Go     = "go"
	Java   = "java"
	Node   = "node"
	Python = "python"
	Ruby   = "ruby"
	Rust   = "rust"

	// Programming languages (native package manager)
	PHP = "php"

	// Databases (native package manager)
	MongoDB    = "mongodb"
	MySQL      = "mysql"
	PostgreSQL = "postgresql"
	Redis      = "redis"
	SQLite     = "sqlite"
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
