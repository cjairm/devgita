# Devgita - Cross-platform Development Environment Manager

This is a Go CLI tool that automates the installation and configuration of development environments across macOS and Debian/Ubuntu systems. The tool helps developers quickly set up their preferred development stack with consistent configurations.

## Project Overview

Devgita (`dg`) is designed to eliminate the manual setup pain when configuring new development machines or maintaining existing ones. It provides both automated installation and manual configuration options for maximum flexibility.

### Core Philosophy

- **Cross-platform support**: Works on both macOS (via Homebrew) and Debian/Ubuntu (via apt/dpkg)
- **Smart installation**: Detects existing packages to avoid conflicts and duplicates
- **Configuration management**: Maintains global state of what was installed vs pre-existing
- **Safe operations**: Only uninstalls what Devgita installed, leaving user packages intact
- **Rollback capability**: Can revert failed installations to keep systems clean

## Project Architecture

```
devgita/
├── cmd/ # Cobra CLI commands and entry points
│ ├── install.go # Main installation command logic
│ └── root.go # Root command with help and global flags
├── internal/
│ ├── apps/ # App-specific installation modules
│ │ ├── aerospace/ # Window manager for macOS
│ │ ├── alacritty/ # Terminal emulator
│ │ ├── neovim/ # Text editor
│ │ ├── tmux/ # Terminal multiplexer
│ │ └── ... # Other development tools
│ ├── commands/ # Platform-specific command implementations
│ │ ├── base.go # Common installation patterns and utilities
│ │ ├── macos.go # macOS-specific commands (brew, etc.)
│ │ ├── debian.go # Debian/Ubuntu commands (apt, dpkg)
│ │ └── factory.go # Platform detection and command factory
│ ├── config/ # Configuration management
│ │ ├── fromFile.go # Global config loading/saving
│ │ └── fromContext.go # Runtime configuration
│ └── tui/ # Terminal user interface components (work not implemented, but planned)
├── pkg/ # Reusable utilities
│ ├── constants/ # App names, paths, and global constants
│ ├── files/ # File operations (copy, backup, etc.)
│ ├── paths/ # Cross-platform path resolution
│ └── utils/ # CLI utilities and error handling
├── configs/ # Default configuration templates
│ ├── alacritty/ # Terminal emulator configs
│ ├── neovim/ # Vim configuration with plugins
│ ├── tmux/ # Terminal multiplexer setup
│ └── themes/ # Color schemes and themes
└── logger/ # Structured logging with verbose mode
```

## Supported Categories

### Development Tools

- **Neovim**: Modern text editor with LSP support and plugin ecosystem
- **Git**: Version control with sensible defaults and aliases
- **Mise**: Runtime version manager for multiple languages

### Terminal & Shell

- **Alacritty**: GPU-accelerated terminal emulator
- **Tmux**: Terminal multiplexer for session management
- **Zsh enhancements**: Autosuggestions, syntax highlighting, Powerlevel10k

### Desktop Applications

- **Aerospace**: Tiling window manager for macOS
- **Fastfetch**: System information display

### Development Languages

- **Node.js**: JavaScript runtime
- **Python**: Programming language with pip
- **Go**: Systems programming language
- **Rust**: Memory-safe systems language

### Databases

- **PostgreSQL**: Relational database
- **Redis**: In-memory data store
- **MongoDB**: Document database

### Fonts & Themes

- **JetBrains Mono**: Developer-focused monospace font
- **Fira Code**: Font with programming ligatures
- **Custom themes**: Color schemes for terminals and editors

## Installation Patterns

### Smart Installation Logic

1. **Detection**: Check if package already exists on system
2. **Classification**: Mark as "installed by devgita" vs "pre-existing"
3. **Installation**: Only install if not present
4. **Configuration**: Apply custom configs from `configs/` directory
5. **Tracking**: Update global manifest for future reference

### Platform Abstraction

- **Factory pattern**: `commands.NewCommand()` returns platform-specific implementation
- **Interface compliance**: All platforms implement common installation interface
- **Path resolution**: Cross-platform paths via `pkg/paths` utilities

## Configuration Management

### Global State Tracking

```yaml
installed:
  packages: ["neovim", "alacritty"]
  fonts: ["JetBrainsMono"]
  themes: ["tokyonight"]
already_installed:
  packages: ["git", "curl"] # Pre-existing, won't be uninstalled
```

### Configuration Templates

- Stored in `configs/` directory
- Applied after successful installation
- Can be reconfigured with `dg configure --force`
- Supports themes and fonts switching

## Command Structure

### Primary Commands

- `dg install` - Interactive installation with TUI selection
- `dg configure` - Apply/update configurations
- `dg uninstall` - Remove devgita-managed packages only
- `dg list` - Show installed packages and their status
- `dg change --theme/--font` - Switch themes or fonts

### Safety Features

- `--soft` mode for detection-only runs
- Confirmation prompts for destructive operations
- Backup creation before major changes
- Rollback on installation failures

## Development Workflow

### Adding New Applications

1. Create module in `internal/apps/newapp/`
2. Implement installation logic for both platforms
3. Add configuration templates to `configs/newapp/`
4. Register in main installation flow
5. Add tests following existing patterns

### Testing Strategy

- Unit tests for platform-specific logic
- Integration tests with temporary directories
- Mock interfaces for external dependencies
- Test both successful and failure scenarios

This architecture enables Devgita to be a reliable, cross-platform development environment manager that respects existing system configurations while providing powerful automation capabilities.
