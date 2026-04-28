# Devgita

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![macOS](https://img.shields.io/badge/macOS-13+-black?logo=apple)](README.md#supported-platforms)
[![Linux](https://img.shields.io/badge/Linux-Debian%2012%2B-orange?logo=linux)](README.md#supported-platforms)

One command to set up a complete, configured development environment. Automates installation and configuration of terminal tools, language runtimes, database systems, and desktop applications on macOS and Linux.

**Core value proposition:** Instead of manually installing 40+ tools across 10+ installers and editing 100+ configuration files, run `dg install` and get a fully configured dev environment in minutes.

---

## 🚀 Quick Start

**1. Install devgita:**

```bash
curl -fsSL https://raw.githubusercontent.com/cjairm/devgita/main/install.sh | bash
```

**2. Restart your shell:**

```bash
source ~/.zshrc  # or source ~/.bashrc for bash
```

**3. Set up your environment:**

```bash
dg install
```

That's it! You'll be prompted to select which programming languages, databases, and desktop apps to install.

---

## ✨ What Gets Installed

**Terminal Tools** (always included)

- Shell: Zsh with Powerlevel10k, syntax highlighting, autosuggestions
- Editors: Neovim with LSP support
- Multiplexer: Tmux with custom configuration
- Modern utilities: fzf, ripgrep, bat, eza, zoxide, btop, fd, tldr, lazydocker, lazygit, fastfetch
- Runtime manager: Mise for managing multiple language versions
- Fonts: JetBrains Mono and developer fonts

**Programming Languages** (interactive selection)

- Node.js (LTS)
- Python (latest)
- Go (latest)
- PHP
- Rust (latest)

**Database Systems** (interactive selection)

- PostgreSQL
- Redis
- MySQL
- MongoDB
- SQLite

**Desktop Applications** (platform-specific)

- **macOS:** Docker Desktop, Alacritty, Brave, Aerospace, Raycast, GIMP, Flameshot
- **Linux:** Docker Desktop, Alacritty, Brave, i3, Ulauncher, GIMP, Flameshot

---

## 📋 Requirements

### System Requirements

- **macOS:** 13+ (Ventura or newer)
  - Apple Silicon (M1/M2/M3+) or Intel chips
- **Linux:** Debian 12+ (Bookworm) or Ubuntu 24+
  - x86_64 architecture
- **Free disk space:** ~5GB for all tools and databases
- **Internet connection:** Required for downloading packages

### Prerequisites

- **curl** — For downloading the installer
- **bash** — For running the installation script
- **git** — Automatically installed or required for use

---

## Installation

### Quick Install (Recommended)

Install devgita with a single command. This will:

1. Download the appropriate binary for your OS and architecture
2. Install to `~/.local/bin/devgita`
3. Configure your PATH automatically

```bash
curl -fsSL https://raw.githubusercontent.com/cjairm/devgita/main/install.sh | bash
source ~/.zshrc  # or source ~/.bashrc
```

### Manual Installation (Verify First)

If you prefer to review the installer script:

```bash
# Download the installer
curl -fsSL -o install.sh https://raw.githubusercontent.com/cjairm/devgita/main/install.sh

# Review it
cat install.sh

# Run it
bash install.sh
```

## Usage

### Install Development Environment

Run the interactive installer to set up your development environment:

```bash
devgita install
```

Or use the short alias:

```bash
dg install
```

This will guide you through:

1. Installing terminal tools and shell configuration
2. Selecting programming languages to install
3. Selecting database systems to install
4. Installing desktop applications

### Category-Specific Installation

Install only specific categories:

```bash
# Install only terminal tools
dg install --only terminal

# Install only languages (interactive selection)
dg install --only languages

# Install only databases (interactive selection)
dg install --only databases

# Install only desktop applications
dg install --only desktop

# Combine multiple categories
dg install --only terminal,languages
```

### Skip Categories

Install everything except specific categories:

```bash
# Install everything except desktop apps
dg install --skip desktop

# Skip multiple categories
dg install --skip databases,desktop
```

### Available Commands

- `dg install` - Install and configure development environment
  - `--only <categories>` - Install only specified categories (terminal, languages, databases, desktop)
  - `--skip <categories>` - Install everything except specified categories
  - `--verbose` - Enable verbose logging
- `dg --version` - Show version information
- `dg --help` - Show help message

## ⚙️ Configuration

Devgita manages configurations in `~/.config/devgita/`:

- **Global state**: `global_config.yaml` tracks installed packages
- **Shell integration**: `devgita.zsh` sourced from your shell config
- **App configs**: Application-specific configuration templates

### Customization

After installation, you can customize:

- Neovim: `~/.config/nvim/init.lua`
- Tmux: `~/.tmux.conf`
- Alacritty: `~/.config/alacritty/alacritty.toml`
- Git: `~/.config/git/.gitconfig`
- i3 (Linux): `~/.config/i3/config`
- Aerospace (macOS): `~/.config/aerospace/aerospace.toml`

## ❓ Troubleshooting

### Common Issues

**Binary not found after installation:**

```bash
# Restart your shell to pick up PATH changes
exec $SHELL

# Or manually source shell config
source ~/.zshrc  # or ~/.bashrc
```

**Permission denied when installing:**

```bash
# Ensure install directory is writable
mkdir -p ~/.local/bin
chmod 755 ~/.local/bin
```

**Unsupported platform:**

- Devgita requires macOS 13+ or Debian 12+/Ubuntu 24+
- Only amd64 and arm64 (Apple Silicon) architectures are supported

**Mise commands not found after installation:**

```bash
# Mise needs to be activated in your shell
exec $SHELL
# Or manually activate:
eval "$(mise activate zsh)"  # or bash, fish, etc.
```

**Language version not available:**

```bash
# List installed versions
mise list

# Install specific version
mise install node@20
mise use --global node@20
```

For more help, see the [full documentation](docs/) or [open an issue](https://github.com/cjairm/devgita/issues).

---

## 📚 Documentation

- **[Product Specification](docs/spec.md)** — Detailed feature documentation, edge cases, testing strategy
- **[Architecture](docs/guides/cross-platform-installation.md)** — Cross-platform installation strategy and package mappings
- **[App Guides](docs/apps/)** — Individual configuration guides for each installed app
- **[Roadmap](ROADMAP.md)** — Planned features, future commands, and open questions
- **[Contributing Guide](CONTRIBUTING.md)** — Development setup, testing, release process

---

## 🤝 Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for:

- Development setup and build instructions
- Testing and code quality standards
- Git workflow and PR process
- Release procedures

Quick links:

- [Open an issue](https://github.com/cjairm/devgita/issues)
- [View the roadmap](ROADMAP.md)
- [See development guide](CONTRIBUTING.md)

---

## 📄 License

[MIT License](LICENSE) — Use devgita freely in personal and commercial projects.

---

## Advanced Usage

### Building from Source

### Prerequisites

- Go 1.21 or newer
- Git

### Build Commands

```bash
# Build for current platform
make build

# Build all platform binaries
make all

# Build specific platforms
make build-darwin-arm64    # macOS Apple Silicon
make build-darwin-amd64    # macOS Intel
make build-linux-amd64     # Linux/Debian/Ubuntu

# Run tests
make test

# Code quality checks
make lint

# Clean build artifacts
make clean
```

### Testing Local Builds

Test a local build before releasing:

```bash
# Build for your platform
make build

# Install using local binary
bash install.sh --local ./devgita-darwin-arm64
```
