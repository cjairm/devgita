# Devgita Product Specification

**Last Updated**: 2026-05-11  
**Owner**: @cjairm

---

## Overview

Devgita is a cross-platform development environment manager that automates installation and configuration of tools, runtimes, databases, and applications on macOS and Linux (Debian/Ubuntu).

**Core value proposition**: One command to install a complete, configured development environment instead of manual setup across 10+ tools and 100+ configuration files.

**Core features:**

- **Smart installation**: Detects existing packages to avoid conflicts
- **Global state tracking**: Maintains what was installed by devgita vs pre-existing
- **Interactive selection**: TUI-based multi-select for languages and databases
- **Safe operations**: Only manages devgita-installed packages
- **Configuration templates**: Consistent, reproducible configs across machines

---

## Architecture

```
devgita/
├── cmd/                     # Cobra CLI commands
├── internal/
│   ├── tooling/            # Category-based coordinators
│   │   ├── terminal/       # Dev tools, shell, editors
│   │   ├── languages/      # Runtime management via Mise
│   │   ├── databases/      # Database systems
│   │   └── worktree/       # Git worktree management
│   ├── apps/               # Individual app implementations (19 apps)
│   ├── commands/           # Platform-specific installers (Darwin, Debian)
│   ├── config/             # State management
│   └── tui/                # Interactive UI components
├── pkg/                    # Shared utilities (logger, paths, constants)
├── configs/                # Configuration templates (embedded at compile time)
└── docs/                   # Documentation
```

**Key patterns:**

- **Interface-based design** for cross-platform compatibility
- **Strategy pattern** for installation (AptStrategy, PPAStrategy, InstallScriptStrategy, etc.)
- **Factory pattern** for platform detection
- **Coordinator pattern** for category orchestration (see `internal/tooling/languages/` as reference)

---

## Features

### 1. Installation Command: `dg install`

Primary entry point for setting up a development environment.

#### Behavior

- Interactive mode: `dg install` launches interactive prompts for category selection
- Category filtering: `--only category1,category2`
- Category exclusion: `--skip category1`
- Verbose logging: `--verbose`

#### Categories

**Terminal Tools** (no selection, always installed)

- Essential: curl, unzip, git, gh (GitHub CLI)
- Shell: Zsh with Powerlevel10k, syntax highlighting, autosuggestions
- Editors: Neovim with LSP support
- Multiplexer: Tmux with custom configuration
- Modern utilities: fzf, ripgrep, bat, eza, zoxide, btop, fd, tldr, lazydocker, lazygit, fastfetch
- Runtime manager: Mise for language version management
- Fonts: JetBrains Mono and developer fonts

**Languages** (interactive selection)

- Node.js (LTS)
- Python (latest)
- Go (latest)
- PHP (native package)
- Rust (latest)
- [Others to be added]

**Databases** (interactive selection)

- PostgreSQL
- Redis
- MySQL
- MongoDB
- SQLite

**Desktop Applications** (category-level selection)

_macOS_:

- Docker Desktop
- Alacritty terminal
- Brave browser
- Aerospace window manager
- Raycast launcher
- GIMP
- Flameshot

_Linux (Debian/Ubuntu)_:

- Docker Desktop
- Alacritty terminal
- Brave browser
- i3 window manager
- Ulauncher launcher
- GIMP
- Flameshot

---

### 2. Configuration Management

#### Persistent State: `~/.config/devgita/`

**`global_config.yaml`**

- Tracks installed packages (name, version, category, timestamp)
- Prevents duplicate installations
- Used by other commands to detect what's already installed

**`devgita.zsh`**

- Shell integration script sourced from `~/.zshrc`
- Sets up Mise activation, aliases, and environment variables
- Platform-specific customizations

**App-specific configs**

- `neovim/init.lua` — Neovim configuration
- `tmux/.tmux.conf` — Tmux configuration
- `alacritty/alacritty.toml` — Terminal emulator config
- `git/.gitconfig` — Git configuration (extends user's existing config)
- Platform-specific: `i3/config` (Linux), `aerospace/aerospace.toml` (macOS)

#### Installation Idempotency

- Check `global_config.yaml` before installing
- Only install if not already present
- Skip system packages that conflict with existing installations (with user prompt)

---

### 3. Command Reference

**Current commands**:

#### `dg install`

See [Installation Command](#1-installation-command-dg-install) above.

#### `dg configure [app]`

Re-applies configuration files for a named app without reinstalling the app itself.

```
dg configure <app> [--force]
```

**Flags**:

- `--force` — Overwrite existing configuration files. Without this flag, configuration is only applied if files do not already exist (soft mode).

**Behavior**:

- Exact app name required (case-sensitive). Supported apps: `aerospace`, `alacritty`, `brave`, `claude`, `devgita`, `docker`, `fastfetch`, `flameshot`, `gimp`, `git`, `i3`, `lazydocker`, `lazygit`, `mise`, `neovim`, `opencode`, `raycast`, `tmux`, `ulauncher`.
- Apps that have no configuration to deploy (e.g., `brave`) return `ErrConfigureNotSupported` — the command prints an info message and exits zero.
- Unknown app names print a sorted list of supported apps and exit non-zero.

**Examples**:

```
dg configure git            # Apply git config if not already present
dg configure neovim --force # Overwrite existing neovim config
dg configure brave          # Info: configure not supported for brave (exit 0)
dg configure foo            # Error: unknown app "foo" + supported list (exit non-zero)
```

**Planned commands**: See [ROADMAP.md](ROADMAP.md) for planned features and future commands.

#### `dg completion [shell]`

Generates a shell completion script for the given shell.

```
dg completion [bash|zsh|fish|powershell]
```

**Behavior**:

- Prints the completion script to stdout; source it or add to your shell config.
- Example: `dg completion zsh > ~/.zsh/completions/_dg`
- Tab completion is pre-wired for `dg worktree remove <name>` and `dg worktree jump`.

#### `dg worktree`

Manages git worktrees with integrated tmux sessions and AI coders.

```
dg worktree <subcommand> [flags]
dg wt <subcommand> [flags]     # alias
```

**Subcommands**:

| Subcommand      | Description                                            |
| --------------- | ------------------------------------------------------ |
| `create <name>` | Create a new worktree + tmux window                    |
| `list`          | List all managed worktrees                             |
| `remove [name]` | Remove a worktree (interactive picker if name omitted) |
| `jump`          | fzf-powered picker to switch between worktree windows  |
| `repair <name>` | Recreate the tmux window for an existing worktree      |
| `prune`         | Remove **all** managed worktrees after confirmation    |

**Flags for `create` and `repair`**:

- `--ai <alias>` / `-a <alias>` — AI coder to launch in the window. Accepted aliases: `opencode`, `oc`, `claude`, `cc`, `claudecode`. Resolution order: flag → `DEVGITA_AI` env var → `worktree.default_ai` in `global_config.yaml`.

**Flag for `remove`**:

- `--force` / `-f` — Force removal even if the worktree has uncommitted changes.

**Examples**:

```
dg wt create feature-login                  # Create worktree, use default AI
dg wt create feature-login --ai claude      # Create with Claude Code
dg wt repair feature-login                  # Recreate missing tmux window
dg wt jump                                  # fzf picker (ctrl-d delete, ctrl-r repair)
dg wt prune                                 # Remove all worktrees (prompts for confirmation)
```

**Planned commands**: See [ROADMAP.md](ROADMAP.md) for planned features and future commands.

---

## Behavior & Edge Cases

### Installation Failures

- If an app installation fails partway through, document which apps were installed
- Provide clear error messages with steps to fix (e.g., "Permission denied, run: `sudo chown`")
- Do not partially commit state; either succeed or roll back to previous state

### Platform Differences

- Same feature set on macOS and Linux where possible
- Platform-specific apps clearly labeled in help text and documentation
- Same command syntax across platforms

### Version Management

- Languages installed via Mise (automatic version tracking)
- Mise enables multiple versions of same language
- Database versions follow platform package manager defaults
- User can override post-installation (e.g., `dg install node@20`)

### Configuration Files

- Templates provided in `configs/` directory
- User edits after installation are preserved (not overwritten)
- Re-configure command can update files (with user confirmation if file exists)

---

## Installation Flow (Current UX)

```
$ dg install

Welcome to Devgita! Let's set up your development environment.

[✓] Installing terminal tools...
    ├─ curl
    ├─ git
    ├─ Zsh + Powerlevel10k
    ├─ Neovim
    ├─ Tmux
    ├─ fzf, ripgrep, bat, ...
    └─ Mise

Select programming languages to install:
  ◉ Node.js (LTS)
  ○ Python
  ○ Go
  ○ PHP
  ○ Rust

[✓] Installing languages...

Select databases to install:
  ◉ PostgreSQL
  ○ Redis
  ○ MySQL
  ○ MongoDB
  ○ SQLite

[✓] Installing databases...

Install desktop applications?
  ◉ Docker Desktop
  ○ Alacritty
  ○ Brave Browser
  ○ [others...]

[✓] Installing desktop apps...

✓ Setup complete! Restart your shell to activate.
  source ~/.zshrc
```

---

## Error Handling

### Common Issues & Messages

| Error                | Message                                                         | Resolution                              |
| -------------------- | --------------------------------------------------------------- | --------------------------------------- |
| Missing dependency   | "Git not found. Install from: [link]"                           | Prompt user to install dependency first |
| Permission denied    | "Permission denied: [path]. Run: `sudo mkdir -p [path]`"        | Specific fix instructions               |
| Package conflict     | "Homebrew already has [package] installed. Use system version?" | Allow user to skip or force             |
| Unsupported platform | "Your system (macOS 12) is not supported. Requires macOS 13+"   | Clear version requirement               |

---

## Testing Strategy

### Unit Tests

- Config parsing and validation
- State tracking logic
- Command argument parsing
- Cross-platform path handling

### Integration Tests

- End-to-end installation flows
- Config file creation and updates
- Multi-category installation combinations
- Idempotency (running install twice gives same result)

### Manual Testing Checklist

- [ ] Fresh system install (clean virtual machine)
- [ ] Install + skip category combinations
- [ ] Re-run install (idempotency)
- [ ] Verify shell configuration is sourced correctly
- [ ] Verify all installed tools are in PATH

---

## Platform Scope & Constraints

**Supported platforms:**

- **macOS 13+** (Ventura or newer) with Homebrew; supports both Apple Silicon (M1/M2/M3+) and Intel
- **Debian 12+** (Bookworm) and **Ubuntu 24+** with APT

**Intentional constraints:**

- **CLI-only** — No graphical installation interfaces
- **Official package sources only** — Homebrew (macOS) and APT (Linux); no custom repositories
- **No Windows support** — macOS and Linux only

---

## Related Documents

- `CLAUDE.md` — Development guidelines and constraints
- `docs/decisions/` — Architectural decisions
- `docs/plans/cycles/` — Feature planning and cycles
- `README.md` — User-facing documentation
