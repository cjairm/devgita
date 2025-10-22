# Devgita - Cross-platform Development Environment Manager

A Go CLI tool that automates installation and configuration of development environments across macOS (Homebrew) and Debian/Ubuntu (apt) systems.

## Core Features

- **Smart installation**: Detects existing packages to avoid conflicts
- **Global state tracking**: Maintains what was installed vs pre-existing
- **Interactive selection**: TUI-based multi-select for languages and databases
- **Safe operations**: Only manages devgita-installed packages
- **Configuration templates**: Consistent configs across machines

## Architecture

```
devgita/
├── cmd/                     # Cobra CLI commands
├── internal/
│   ├── tooling/            # Category-based coordinators
│   │   ├── terminal/       # Dev tools, shell, editors
│   │   ├── languages/      # Runtime management via Mise
│   │   ├── databases/      # Database systems
│   │   └── desktop/        # GUI applications
│   ├── apps/               # Individual app implementations
│   ├── commands/           # Platform-specific installers
│   ├── config/             # State management
│   └── tui/                # Interactive UI components
├── pkg/                    # Shared utilities
├── configs/                # Configuration templates
└── docs/                   # Documentation
```

## Installation Categories

**Terminal** (`dg install --only terminal`)
- Core tools: curl, unzip, gh, lazydocker, lazygit, fzf, ripgrep, bat, eza, zoxide, btop, fd, tldr
- Shell: Zsh with autosuggestions, syntax highlighting, Powerlevel10k theme
- Editors: Neovim with LSP, Tmux multiplexer
- Runtime manager: Mise for language versions
- System libraries: pkg-config, autoconf, bison, rust, openssl, etc.

**Languages** (`dg install --only languages`) - Interactive TUI selection
- Node.js (LTS), Go (latest), Python (latest) via Mise
- PHP via native package manager

**Databases** (`dg install --only databases`) - Interactive TUI selection  
- PostgreSQL, Redis, MySQL, SQLite

**Desktop** (`dg install --only desktop`)
- Development: Docker Desktop, Alacritty terminal
- Productivity: Brave browser, Flameshot, Raycast (macOS)
- Design: GIMP, Aerospace window manager (macOS)
- Fonts: JetBrains Mono and developer fonts

## Installation Flow

1. **OS validation** - Check macOS 13+/Ventura+ or Debian 12+/Ubuntu 24+
2. **Package manager setup** - Install Homebrew (macOS) or update apt (Debian/Ubuntu)  
3. **Devgita repository** - Clone to `~/.config/devgita/`
4. **Category installation** - Run coordinators based on flags
5. **Configuration** - Apply templates from `configs/` directory

### Smart Installation Logic
- `MaybeInstallPackage()` - Check if package exists before installing
- Global config tracks "installed by devgita" vs "pre-existing"
- Uses Mise for language runtimes (Node.js, Go, Python)
- Platform abstraction via command factory pattern

## Commands

```bash
dg install                    # Interactive installation with TUI
dg install --only terminal   # Install only terminal tools
dg install --only languages  # Interactive language selection
dg install --skip desktop    # Install everything except desktop apps
```

**Planned commands:**
- `dg configure` - Apply/update configurations
- `dg uninstall` - Remove devgita-managed packages only
- `dg list` - Show installed packages and status
- `dg change --theme/--font` - Switch themes or fonts

## Configuration Management

**Global state tracking** in `~/.config/devgita/global_config.yaml`:
```yaml
installed:
  packages: ["neovim", "alacritty", "docker"]
  languages: ["node@lts", "go@latest", "python@latest"]
  databases: ["postgresql", "redis"]
already_installed:
  packages: ["git", "curl"] # Pre-existing, won't be uninstalled
```

**Configuration templates**:
- Stored in `configs/` with app-specific subdirectories
- Applied via `MaybeSetup()` methods after installation
- Context-aware based on user selections
- Interactive setup for GitHub SSH keys and macOS privacy permissions

## Development

**Adding new apps:**
1. Create `internal/apps/newapp/` with standard interface
2. Add platform-specific installation logic
3. Create configuration templates in `configs/newapp/`
4. Register in appropriate tooling category
5. Add tests with mock interfaces

**Architecture patterns:**
- Interface-based design for cross-platform compatibility
- Factory pattern for platform detection
- Context propagation for user selections
- Centralized error handling via `utils.MaybeExitWithError()`
