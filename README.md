# Devgita

Cross-platform development environment manager for macOS and Debian/Ubuntu systems.

Automates installation and configuration of:
- Terminal tools (Neovim, Tmux, shell utilities)
- Programming language runtimes (Node.js, Python, Go, PHP, Rust)
- Database systems (PostgreSQL, Redis, MySQL, MongoDB, SQLite)
- Desktop applications (Docker, browsers, window managers)

## Installation

### Quick Install (Recommended)

Install devgita with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/cjairm/devgita/main/install.sh | bash
```

This will:
1. Download the appropriate binary for your OS and architecture
2. Install to `~/.local/bin/devgita`
3. Configure your PATH automatically
4. Verify the installation

After installation, restart your shell or run:
```bash
source ~/.zshrc  # or ~/.bashrc for bash users
```

### Manual Installation

If you prefer to review the script first:

```bash
# Download the installer
curl -fsSL -o install.sh https://raw.githubusercontent.com/cjairm/devgita/main/install.sh

# Review the script
cat install.sh

# Run the installer
bash install.sh
```

### Supported Platforms

- **macOS**: 13+ (Ventura or newer)
  - Apple Silicon (M1/M2/M3): `darwin-arm64`
  - Intel chips: `darwin-amd64`
- **Linux**: Debian 12+ (Bookworm) or Ubuntu 24+
  - AMD64 architecture: `linux-amd64`

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

## What Gets Installed

### Terminal Tools

Core development tools and shell environment:

- **Essential tools**: curl, unzip, git, gh (GitHub CLI)
- **Shell**: Zsh with autosuggestions, syntax highlighting, Powerlevel10k theme
- **Editors**: Neovim with LSP support
- **Multiplexer**: Tmux with custom configuration
- **Modern utilities**: 
  - fzf (fuzzy finder)
  - ripgrep (fast search)
  - bat (better cat)
  - eza (better ls)
  - zoxide (smart cd)
  - btop (system monitor)
  - fd (better find)
  - tldr (quick command help)
  - lazydocker (Docker TUI)
  - lazygit (Git TUI)
  - fastfetch (system info)
- **Runtime manager**: Mise for language version management
- **Fonts**: JetBrains Mono and other developer fonts

### Languages (Interactive Selection)

Programming language runtimes managed via Mise:

- **Node.js** (LTS version)
- **Python** (latest)
- **Go** (latest)
- **PHP** (native package manager)
- **Rust** (latest)

### Databases (Interactive Selection)

Database systems:

- **PostgreSQL**
- **Redis**
- **MySQL**
- **MongoDB**
- **SQLite**

### Desktop Applications

GUI applications for development:

**macOS:**
- Docker Desktop
- Alacritty terminal
- Brave browser
- Aerospace window manager
- Raycast launcher
- GIMP
- Flameshot

**Debian/Ubuntu:**
- Docker Desktop
- Alacritty terminal
- Brave browser
- i3 window manager
- Ulauncher launcher
- GIMP
- Flameshot

## Configuration

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

## Building from Source

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

## Development

### Project Structure

```
devgita/
├── cmd/                     # CLI commands
├── internal/
│   ├── apps/               # Individual app modules
│   ├── tooling/            # Category coordinators
│   ├── commands/           # Platform-specific installers
│   ├── config/             # State management
│   └── tui/                # Interactive UI
├── pkg/                    # Shared utilities
├── configs/                # Embedded configuration templates
├── docs/                   # Documentation
├── install.sh              # Zero-dependency installer
└── Makefile                # Build system
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/apps/neovim/

# Verbose output
go test -v ./...
```

### Code Quality

```bash
# Static analysis
go vet ./...

# Format code
go fmt ./...

# Or use make
make lint
```

## Troubleshooting

### Installation Issues

**Binary not found after installation:**
```bash
# Restart your shell
exec $SHELL

# Or manually source shell config
source ~/.zshrc  # or ~/.bashrc
```

**Permission denied:**
```bash
# Ensure install directory is writable
mkdir -p ~/.local/bin
chmod 755 ~/.local/bin
```

**Unsupported platform:**
- Devgita requires macOS 13+ or Debian 12+/Ubuntu 24+
- Only amd64 and arm64 architectures are supported

### Runtime Issues

**Mise commands not found:**
```bash
# Restart shell to activate mise
exec $SHELL

# Or manually activate
eval "$(mise activate zsh)"
```

**Language version not found:**
```bash
# List installed versions
mise list

# Install specific version
mise install node@20
mise use --global node@20
```

## Contributing

Contributions welcome! See development notes below for project roadmap and areas needing work.

## License

[Add license information]

---

## Development Notes

### Pending Work

- ~Debian/Ubuntu compatibility~
- ~Global document to see what we have installed already~
  - ~Check if languages are already installed - to avoid duplicates/collisions~
- Command `re-configure` app
  - Reinstall configuration files / should we uninstall and install app?
- Command `uninstall` app, font or theme.
  - make sure we only uninstall stuff related in the config file.
- Command to select a new `theme`
  - Update background images, too?
- Command to select a new `font`
- Help to display available commands
- MANUAL - Super important
- Revert installed programs if any error
  - We need to confirm if program was installed by us. If update leave the program there
- REMOVE GIT STUFF or any unrelated to installing new apps
- ~Add `verbose` option for all commands~
- Allow select between tmux and zellij
- Allow select between opencode and claudecode
- Install Ollama? So create a new category? AI tools?
- Search for a good terminal app to run AI commands (??? what's this???)
- [See](https://github.com/cjairm/devgita/blob/038da72eec456e0a60c50dce2bc9ab615c795fb2/configs/bash/init.zsh#L15) and extend to zip un remove all unneeded docs example node_modules, etc

### Planned Commands

- `dg install ...` (`--soft` that does `maybeInstall`)
- ~`dg reinstall ...` this will be replaced by `dg install ... --force`. We should uninstall first and then install again.~
- `dg configure ...` this will update changes to the configuration if made through the `dg install` command.
- ~`dg re-configure ..` this will be replaced by `dg configure ... --force`. We should remove files and recreate them.~
- `dg uninstall ...`
- `gd update ... [--neovim=[...options] --aerospace=[...options] ...flags] ...apps`
- `dg list` or `dg installed`
- `dg check-updates`
- `dg backup ...` - This would allow users to create backups of their current configurations before making changes.
- `dg restore ...` - This would allow users to revert to a previous configuration if needed.
- `dg validate ...` - This could check if the current configuration is valid and if all dependencies are met.
- `dg change --theme=[...options] --font=[...options]`

Note. We should optionally be able to pass `--app` or `--package` to only do it for one app/package

#### Considerations

- **Error Handling**: Ensure that all commands have robust error handling to provide meaningful feedback to users in case something goes wrong.
- **Logging**: Implement logging for actions taken by the commands, which can help in troubleshooting and understanding user actions.
- **User Prompts**: For commands that make significant changes (like uninstalling or reinstalling), consider adding user prompts to confirm actions.
- **Dependency Management**: Ensure that when installing or uninstalling, dependencies are also managed appropriately to avoid leaving orphaned packages.
- **Documentation**: Make sure that the MANUAL is comprehensive and includes examples for each command to help users understand how to use them effectively.

### Optional Apps/Packages

- Postman
- fc-list ([maybe required](https://github.com/cjairm/devgita/commit/c01797defb5e95a5ccce4206d46f435f9c513215)?)
- We must add [devpod](https://devpod.sh/) Uses dev containers underneath. Ssh in and run Vim.... to work with AI agent safely
- install snap? Also maybe required `sudo apt install snapd`
  - `Slack, Spotify, VS Code... etc`
- Music
- Emails (?)
- By default enable git: https://github.com/cjairm/devgita/commit/11c4a54d1d79cbe977a8053dbb755322cc83376e

### Questions

- Do we need to create shortcuts?
- What's that best way to handle mise? it's useful, but the documentation is difficult to follow
- (?) Git related - Move all git to `dg git clean --flags`, `dg git revert --flags`, or `dg clean-branch` ?
- (?) npm related - fully clean? fresh-installs? `dg npm clean`
- Command to update (we may need to solve issues). If we want to update apps can be difficult; we need to handle breaking changes
- Maybe a TUI? - https://github.com/charmbracelet/lipgloss
  - go get github.com/charmbracelet/bubbletea
  - go get github.com/charmbracelet/lipgloss
  - go get github.com/charmbracelet/bubbles/spinner
- Add new tmux window commands? `dg tmux --new-window="~/my-path/hello-world"` // Will add `tmhello-world`. This will remove the necessity for some custom commands (`tmn`?)
- AI related:
  - build or reuse https://github.com/AndyMik90/Auto-Claude ?
  - templates https://www.aitmpl.com/agents ?
  - ideas https://huggingface.co/pricing ?
  - Gemini CLI?
    - extension: https://developers.googleblog.com/conductor-introducing-context-driven-development-for-gemini-cli/
  - docs:
    - https://a2ui.org/
    - https://genai.owasp.org/resource/cheatsheet-a-practical-guide-for-securely-using-third-party-mcp-servers-1-0/
  - Super powers for Opencode (?) - https://github.com/obra/superpowers
  - Something to integrate to Opencode (?) - https://github.com/affaan-m/everything-claude-code
  - Download, Review, Answer back - with UI to keep track of addressed issues in github
  - Opencode configuration with nvim - https://github.com/nickjvandyke/opencode.nvim

### Developer Commands

Command to update devgita - cli (download latest version)
Add alias of `devgita` to `cli` (for easy access) (-or `dg`-)

```bash
# For M chips (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o cli_mac_arm64

# For Intel chips
GOOS=darwin GOARCH=amd64 go build -o cli_mac_amd64
```

When pushes happen, try:

```bash
# For example, to test an individual file
go test <file-path> -v

# To test all files in the current directory and subdirectories
go test ./...

# To test with coverage
go test <file-path> -cover

# To test with race detector
go test <file-path> -race
```

### Useful Commands

#### Git

##### Steps to Uncommit a PR and Merge Main

**Step 1: Soft reset to before your commit(s)**

```bash
git reset --soft <commit-hash-before-your-work>
```

**Step 2: Stash the staged changes**

```bash
git stash
```

**Step 3: Merge main**

```bash
git merge main
```

**Step 4: If merge has conflicts, resolve them and commit**

```bash
# Resolve conflicts in files, then:
git add <conflicted-files>
git merge --continue
```

**Step 5: Pop stashed changes back**

```bash
git stash pop
```

**Step 6: (Optional) Unstage everything to keep it uncommitted**

```bash
git restore --staged .
```

**Alternative single-line workflow (when no conflicts):**

```bash
git reset --soft <commit-hash> && git stash && git merge main && git stash pop
```

##### Steps to create a clean branch

**1. Fetch latest changes**

```bash
git fetch origin
```

**2. Create a new clean branch from main**

```bash
git checkout -b <BRANCH-NAME> origin/main
```

**3. Cherry-pick specific commits from your old branch**

```bash
git cherry-pick <COMMIT-HASH-1> <COMMIT-HASH-2> ... <COMMIT-HASH-N>
```

**4. Push the new branch to origin**

```bash
git push -u origin <BRANCH-NAME>
```

##### Steps to create a clean branch with squashed changes

**1. Fetch latest changes**

```bash
git fetch origin
```

**2. Create a new clean branch from main**

```bash
git checkout -b <BRANCH-NAME> origin/main
```

**3. Squash merge changes from your source branch**

```bash
git merge --squash origin/<SOURCE-BRANCH>
```

**4. Commit the squashed changes**

```bash
git commit -m "Squash merge <SOURCE-BRANCH>"
```

**5. Push the new branch to origin**

```bash
git push -u origin <BRANCH-NAME>
```

#### Kubectl

**GET SECRET ENV VARS (from running pod):**

```bash
kubectl exec -n <NAMESPACE> <POD-NAME> -- env | sort
# Note: Only works if pod status is Running, not CrashLoopBackOff
```

**SECRETS (view secret references and decode them):**

Step 1 - See which secrets are referenced:

```bash
kubectl get pod <POD-NAME> -o json | jq -r '.spec.containers[0].envFrom'
```

Step 2 - Decode each secret's values:

```bash
kubectl get secret <SECRET-NAME> -o json | jq -r '.data | to_entries | .[] | "\(.key)=\(.value | @base64d)"'
```

**DIRECT ENV VARS (defined in deployment, not from secrets):**

```bash
kubectl get pod <POD-NAME> -o json | jq '.spec.containers[0].env'
```
