# Devgita Product Specification

**Last Updated**: 2026-07-15  
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
- Per-app targeting: `--only appname` or `--skip appname` (registry apps only)
- Verbose logging: `--verbose`

#### Per-App Targeting

`--only` and `--skip` accept both category names and individual app names from the registry.

**Granularity levels:**

```
dg install --only terminal          # full terminal category (existing behavior)
dg install --only neovim            # single app by name — only neovim installed
dg install --skip git               # skip git; install everything else normally
dg install --only terminal --skip lazygit  # full terminal minus lazygit
dg install --only neovim --only docker    # neovim (terminal) + docker (desktop) only
```

**Behavior when an app filter is active** (`--only <appname>`):

- Only the specified registry apps are installed in that coordinator
- `InstallDevTools` and `InstallCoreLibs` are skipped (user asked for a specific app, not a full setup)
- Fonts installation is also skipped in the desktop coordinator

**Individually targetable apps** (registry-managed, 18 apps):

| Coordinator | Apps                                                                         |
| ----------- | ---------------------------------------------------------------------------- |
| terminal    | claude, fastfetch, git, lazydocker, lazygit, mise, neovim, opencode, tmux    |
| desktop     | aerospace, alacritty, brave, docker, flameshot, gimp, i3, raycast, ulauncher |

**Note on alacritty:** `alacritty` has `KindTerminal` in the registry but is installed by the desktop coordinator. Use `--only alacritty` (not `--only terminal`) to target it specifically.

**Not individually targetable** (no registry entry):

- Core libs: autoconf, bison, ncurses, openssl, etc.
- Dev tools: bat, fzf, ripgrep, zoxide, etc.
- Languages: node, python, go, rust, php (use interactive selection)
- Databases: postgresql, redis, mysql, etc. (use interactive selection)

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
- Tab completion is pre-wired for `dg worktree remove <name>`.

#### `dg worktree`

Manages git worktrees with integrated tmux sessions and AI coders.

```
dg worktree <subcommand> [flags]
dg wt <subcommand> [flags]     # alias
```

**Subcommands**:

| Subcommand      | Description                                                                                                                                                                                                                                                                                                                                                                              |
| --------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `create <name>` | Create a new worktree + tmux window                                                                                                                                                                                                                                                                                                                                                      |
| `list`          | List all managed worktrees                                                                                                                                                                                                                                                                                                                                                               |
| `remove [name]` | Remove a worktree (interactive picker if name omitted)                                                                                                                                                                                                                                                                                                                                   |
| `ui` / `dash`   | Full-screen TUI dashboard (NERDTree-style tree + branch-diff pane). The pane labels its comparison (`<default-branch> @<merge-base> ← <branch>`), renders one styled header per file with per-file +/- counts, and is focusable (Space) for vim-style scrolling; replaces `jump`. Press `n` to create a worktree without leaving the dashboard — see "Creating from the dashboard" below |
| `repair <name>` | Recreate the tmux window for an existing worktree                                                                                                                                                                                                                                                                                                                                        |
| `prune`         | Remove **all** managed worktrees after confirmation                                                                                                                                                                                                                                                                                                                                      |

**Flags for `create` and `repair`**:

- `--ai <alias>` / `-a <alias>` — AI coder to launch in the window. Accepted aliases: `opencode`, `oc`, `claude`, `cc`, `claudecode`. Resolution order: flag → `DEVGITA_WORKTREE_AI` env var → `worktree.default_ai` in `global_config.yaml`.

**Flag for `create`**:

- `--repo <path>` / `-r <path>` — Path to the repository (`~` is expanded), so the command works
  from any directory. The window opens in a tmux session named after the repo — created when
  missing, reused otherwise — and the attached client switches to it when run inside tmux.
  Without the flag, the repo is the one containing the current directory and the window opens
  in the current session.

**Flag for `remove`**:

- `--force` / `-f` — Force removal even if the worktree has uncommitted changes.

**Examples**:

```
dg wt create feature-login                  # Create worktree, use default AI
dg wt create feature-login --ai claude      # Create with Claude Code
dg wt new fix-auth --repo ~/code/api        # Create for another repo; window opens in its session
dg wt ui                                    # Open TUI dashboard (j/k nav, Enter attach, n create, d delete, D delete + kill session, r repair,
                                            #   Space focuses the diff pane: j/k scroll, [/] jump between files, g/G top/bottom, Esc back)
dg wt repair feature-login                  # Recreate missing tmux window
dg wt prune                                 # Remove all worktrees (prompts for confirmation)
```

**Creating from the dashboard (`n`)**:

- `n` opens a floating repo picker over the dashboard — the background stays visible, matching
  the `?` help overlay. Candidates are ranked: the repo containing the directory `dg wt ui` was
  launched from first (when that directory is inside a git repo — otherwise this source is
  skipped), then the repo under the cursor, then repos from the recent-repos store
  (most-recently-used first), then `zoxide query -l` results when zoxide is installed. Typing
  filters the list; if the query matches nothing, Enter validates it directly as a free-typed
  repo path instead.
- Enter on a repo opens a floating name prompt. Enter creates the worktree (same as `create
--repo`) and attaches into the new window — the TUI exits, identical to pressing Enter on an
  existing row. If the create's pre-flight hook-compatibility check finds warnings, they're
  shown as a status message and a second Enter confirms; any other key cancels the confirm.
- A failed create (invalid path, duplicate name, etc.) is shown as a status message; the
  dashboard keeps running rather than exiting.
- Esc at the repo picker or the name prompt returns to the dashboard unchanged.
- Every successful create — from `ui`, `create`, or `new` — records the repo's root path in
  `global_config.yaml`'s `worktree.recent_repos` (MRU-ordered, capped at 20). This is what lets
  the picker offer a repo that currently has zero worktrees.

**Planned commands**: See [ROADMAP.md](ROADMAP.md) for planned features and future commands.

#### `dg list`

`dg list` is the single inventory command (there is no separate `dg validate`; drift
checking lives inside the dashboard's problems-only view):

- **Data model** — every item Devgita has tracked in `~/.config/devgita/global_config.yaml`
  (both what it installed and what it found pre-existing) is live-checked against the system,
  producing a three-state status per item:
  - `OK` — the presence check ran and found the item.
  - `MISSING` — the check ran and definitively did not find the item.
  - `UNKNOWN` — the check itself failed to run (e.g. `brew`/`dpkg` unavailable). A failed
    check is never conflated with a missing item, so a flaky or unavailable package manager
    can't misreport drift.
  - The `themes` and `terminal_tools` categories are tracked but have no live presence check
    implemented yet (no current code path populates them); if a future feature starts
    tracking items there, they report `UNKNOWN` rather than being silently misreported as
    `OK` or `MISSING`.
- **Dashboard** — renders the collected items as an interactive, grouped list. Opens
  automatically in a terminal; falls back to plain-text output for piped, CI, or `--plain`
  invocations. Keybindings: `j`/`k` move, `h`/`l` collapse/expand a group, `/` enter a text
  filter, `p` toggle problems-only (`MISSING`/`UNKNOWN` items), `g` toggle between grouping
  by category and grouping by status, `?` open the keybinding help overlay, `q` quit.
  While problems-only is active the pane title shows the mode and the `p` hint flips to
  "show all"; if nothing is missing, the pane says so explicitly instead of rendering an
  empty list.

```
dg list [--category <name>] [--plain]
dg installed [--category <name>] [--plain]   # alias
```

**Flags**:

- `--category <name>` — Filter to a single bucket. Valid values: `packages`, `desktop_apps`,
  `fonts`, `themes`, `terminal_tools`, `dev_languages`, `databases`.
- `--plain` — Force plain-text output even when run in a terminal.

**Behavior**:

- In a terminal, opens the dashboard unfiltered (every tracked item, all categories).
- Piped output, CI, or `--plain`: prints one table per non-empty category (name only, no
  live status check); empty categories are omitted. The "Already on this machine (not
  installed by Devgita)" section only prints if it has entries. An empty config prints a
  clear message instead of a blank screen. An unrecognized `--category` value prints an
  error listing the valid category names.
- This is still the MVP: name + category only in plain mode. Per-item version and
  install-timestamp tracking requires a `global_config.yaml` schema change and is planned as
  a future release (see [ROADMAP.md](ROADMAP.md)).

**Examples**:

```
dg list                             # Interactive dashboard in a terminal
dg list --category=terminal_tools   # Only the terminal tools bucket
dg installed                        # Same as 'dg list'
```

#### `dg task`

Developer utility commands for git branch management and npm dependency management.
These commands are callable by both agents (Claude Code, CI, any non-interactive process)
and humans (via the `dge` shell wrapper or directly).

```
dg task <subcommand> [args]
dg t <subcommand> [args]   # alias
```

| Subcommand            | Args       | Description                                                                            |
| --------------------- | ---------- | -------------------------------------------------------------------------------------- |
| `refresh-branch`      | `[target]` | Checkout target (default: `main`), pull, return to previous branch, merge              |
| `reset-main-branch`   | —          | Checkout `main`, hard-reset to `origin/main`                                           |
| `delete-branch`       | `[target]` | Checkout target (default: `main`), fetch, pick a branch via fzf to force-delete        |
| `reinstall-libraries` | —          | `git clean -Xdf`, remove `node_modules/`, `npm install`, remove `tsconfig.tsbuildinfo` |
| `reinstall-library`   | `<name>`   | Remove `node_modules/<name>`, run `npm install`                                        |

**Review scope subcommands** (compact, noise-filtered git context for agents — the
`review-threads` pattern applied to git; `git` plumbing is fetched, Go formatters render):

| Subcommand     | Args / Flags    | Description                                                                                                                                                                                                                                                                                         |
| -------------- | --------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `review-scope` | —               | Fetch origin (bounded, best-effort), then print branch, default branch, ahead/behind, commit subjects, and a per-file stat table. Lockfile-style noise (`package-lock.json`, `go.sum`, `*.min.js`, …) is excluded from the table and noted separately with its own counts — never silently dropped. |
| `branch-diff`  | `--file <path>` | Diff against the merge-base with the default branch, same default exclusions applied in one `git diff` call. Does **not** fetch (reuses `review-scope`'s comparison base within the same review session). `--file` bypasses exclusions for that one file.                                           |

**Pull request subcommands** (via `gh`; data-returning ones are formatted by `jq`
into compact, LLM-oriented output — `gh` fetches/acts, `jq` renders):

| Subcommand              | Args / Flags                                  | Description                                                          |
| ----------------------- | --------------------------------------------- | -------------------------------------------------------------------- |
| `review-threads`        | `--pr N`, `--state unresolved\|resolved\|all` | Render PR review threads as compact markdown (default: unresolved)   |
| `resolve-thread`        | `<id>`                                        | Mark a review thread resolved                                        |
| `unresolve-thread`      | `<id>`                                        | Reopen a resolved review thread                                      |
| `reply-thread`          | `<id> <body>`                                 | Reply to a review thread                                             |
| `create-pr`             | `--title` (req), `--body`, `--base`           | Open a PR from the current branch; prints the URL                    |
| `update-pr-description` | `--pr N`, `--body` (req)                      | Replace a PR's description                                           |
| `approve-pr`            | `--pr N`, `--body`                            | Approve a PR                                                         |
| `request-changes-pr`    | `--pr N`, `--body` (req)                      | Request changes on a PR                                              |
| `request-review`        | `--pr N`, `<reviewer>...` (req)               | Re-request review (adds reviewers back to the requested list)        |
| `comment-pr`            | `--pr N`, `--body` (req)                      | Post a top-level PR comment                                          |
| `merge-pr`              | `--pr N`, `--method squash\|merge\|rebase`    | Merge a PR (default: squash)                                         |
| `pr-view`               | `--pr N`                                      | Compact PR summary (number, title, state, mergeable, review, branch) |
| `pr-checks`             | `--pr N`                                      | CI check status, one line per check (returns data even when red)     |
| `current-pr`            | —                                             | PR number for the current branch                                     |
| `current-repo`          | —                                             | Current repository as `owner/name`                                   |

For every PR subcommand, `--pr` defaults to the current branch's PR when omitted.
Review-thread output is paginated across all threads (`gh api graphql --paginate`).

`dge` (the shell function in `devgita.zsh`) is now a thin wrapper that forwards to `dg task`:

```sh
dge() {
  if [[ $# -eq 0 ]]; then dg task --help; return; fi
  dg task "$@"
}
```

Agents should prefer `dg task` directly; humans can use either `dg task` or `dge`.

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
