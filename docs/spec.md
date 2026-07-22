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

| Subcommand      | Description                                            |
| --------------- | ------------------------------------------------------ |
| `create <name>` | Create a new worktree + tmux window                    |
| `list`          | List all managed worktrees                             |
| `remove [name]` | Remove a worktree (interactive picker if name omitted) |
| `repair <name>` | Recreate the tmux window for an existing worktree      |
| `prune`         | Remove **all** managed worktrees after confirmation    |

**Flags for `create` and `repair`**:

- `--ai <alias>` / `-a <alias>` — AI coder to launch in the window. Accepted aliases: `opencode`, `oc`, `claude`, `cc`, `claudecode`.
- `--layout <name>` / `-l <name>` — Window layout to build. Valid names: `opencode`, `claude`, `claude-nvim`, `nvim` (see "Window layouts" below). Mutually exclusive with `--ai` — passing both is a cobra error before either command runs.

  **Resolution order** (highest wins; each rule below only applies when none of the rules above it fired):

  1. `--layout` flag — explicit layout name, wins over everything.
  2. `--ai` flag — derived into a single-pane layout running that coder.
  3. `DEVGITA_WORKTREE_AI` env var — derived into a single-pane layout.
  4. `worktree.default_layout` in `global_config.yaml`.
  5. `worktree.default_ai` in `global_config.yaml` — derived into a single-pane layout.
  6. Default: `opencode`, single-pane.

  `repair` uses the exact same resolution order as `create` — it does **not** remember the layout a worktree was originally created with. If the window is missing, it is rebuilt from scratch using whatever `--layout`/`--ai`/env/config resolves to at that moment. If the window already exists (e.g. only one pane in it was closed), `repair` only relaunches the AI coder in the existing window — it does not add or recreate missing panes, since there's no way to tell whether the surviving panes already match the requested layout.

**Window layouts**:

A layout is a named, ordered set of tmux panes built when a worktree's window is created or repaired. Built-in layouts (no config required):

| Layout        | Panes                                                |
| ------------- | ---------------------------------------------------- |
| `opencode`    | Single pane running OpenCode                         |
| `claude`      | Single pane running Claude Code                      |
| `claude-nvim` | Claude Code and Neovim side by side (vertical split) |
| `nvim`        | Single pane running Neovim only                      |

Before any tmux window is touched, every pane's underlying tool is checked for installation; a layout referencing a missing tool fails with an actionable error and the worktree is not created. If a multi-pane window fails to build partway through (e.g. a later pane's split fails), the partially built window is killed and the worktree is rolled back — never left half-created.

**Worktree scan and layout config keys** (`global_config.yaml`, under `worktree:`):

- `search_paths` — list of directories to scan for git repositories to offer in the `n`/`N` repo picker (see "Creating from the dashboard" below). Default: empty, which disables the scan entirely — this is the only off-switch. The scan walks each path with `filepath.WalkDir`, stops descending at a repo's `.git` boundary (so nested repos/submodules are not listed as separate entries), and skips `node_modules`, `.cache`, `vendor`, `target`, `dist`, and `.git` directories encountered during the walk (a configured root itself is still scanned even if its name matches one of these, e.g. a root literally named `vendor`).
- `scan_depth` — max directory depth below each search path to descend. Unset, `0`, or negative all mean the default of `4` — there is no separate "unlimited" or "disabled via depth" mode; use an empty `search_paths` to disable scanning.
- `default_layout` — default window layout name (see the resolution order above). Default: empty, which means rule 5 (`default_ai`-derived single-pane layout) or the built-in `opencode` fallback applies instead.

**Flag for `create`**:

- `--repo <path>` / `-r <path>` — Path to the repository (`~` is expanded), so the command works
  from any directory. The window opens in a tmux session named after the repo — created when
  missing, reused otherwise — and the attached client switches to it when run inside tmux.
  Without the flag, the repo is the one containing the current directory and the window opens
  in the current session.

**Flag for `remove`**:

- `--force` / `-f` — Force removal even if the worktree has uncommitted changes.

**Adopting an existing branch (`create`)**: if a branch named `<name>` already exists locally,
`create` adopts it into the worktree instead of failing. If that branch is currently checked out
in the main clone, git refuses to check it out again in the new worktree — `create` frees it
first by switching the main clone to the repo's default branch (printing a one-line note so the
switch isn't a surprise), then proceeds. If the main clone's checkout of that branch has
uncommitted changes, `create` refuses up front with an error telling you to commit or stash them
first, rather than risk carrying or losing that work.

**Examples**:

```
dg wt create feature-login                  # Create worktree, use default AI/layout
dg wt create feature-login --ai claude      # Create with Claude Code
dg wt create feature-login --layout nvim    # Create with the nvim-only layout
dg wt new fix-auth --repo ~/code/api        # Create for another repo; window opens in its session
dg wt repair feature-login                  # Recreate missing tmux window (rebuilds current layout resolution, not the original)
dg wt prune                                 # Remove all worktrees (prompts for confirmation)
```

**Creating from the dashboard (`n` / `N`)**:

- `n` opens a floating repo picker over the dashboard — the background stays visible, matching
  the `?` help overlay. Candidates are ranked: the repo containing the directory `dg ws` was
  launched from first (when that directory is inside a git repo — otherwise this source is
  skipped), then the repo under the cursor, then repos from the recent-repos store
  (most-recently-used first), then repos found by scanning `worktree.search_paths` (see above —
  contributes nothing until configured), then `zoxide query -l` results when zoxide is
  installed. Typing filters the list; if the query matches nothing, Enter validates it directly
  as a free-typed repo path instead.
- Enter on a repo opens a floating name prompt. For `n`, Enter on the name creates the worktree
  immediately using the resolved default layout (same as `create --repo` with no `--layout`/
  `--ai`) and attaches into the new window — the TUI exits, identical to pressing Enter on an
  existing row.
- `N` follows the same repo-pick → name-prompt flow as `n`, but after the name is entered it
  opens one more floating picker: a layout picker listing the built-in layout names
  (`opencode`, `claude`, `claude-nvim`, `nvim`), cursor pre-positioned on the resolved default so
  accepting it is a single Enter. Picking a layout (or free-typing an unlisted name — `ResolveLayout`
  validates it and reports an unknown name the same way the CLI does) creates the worktree with
  that layout and attaches, same as `n`.
- If the create's pre-flight hook-compatibility check finds warnings, they're shown as a status
  message and a second Enter confirms; any other key cancels the confirm.
- A failed create (invalid path, duplicate name, unknown/uninstalled layout, etc.) is shown as a
  status message; the dashboard keeps running rather than exiting.
- Esc at the repo picker, the name prompt, or (for `N`) the layout picker returns to the
  dashboard unchanged; no worktree is created.
- Every successful create — from `ui`, `create`, or `new` — records the repo's root path in
  `global_config.yaml`'s `worktree.recent_repos` (MRU-ordered, capped at 20). This is what lets
  the picker offer a repo that currently has zero worktrees.

**In-progress feedback**: any dashboard action that builds or tears down tmux state shows a
status line at the bottom while it runs, so the dashboard is never silently unresponsive during
the (occasionally slow) git/tmux work: `creating worktree: <name> (<layout>)…` on create,
`repairing: <name> (<layout>)…` on `r` (and on attaching to a worktree whose window is missing,
which auto-repairs), and `deleting: <name>…` on the confirming `d`/`D` press. The layout is named
by its tools (`claude`, `opencode`, `neovim`, `claude + neovim`). Each message is replaced by the
result the moment the action finishes (and is moot inside tmux on a create/attach that succeeds,
since the TUI then attaches and exits).

**Planned commands**: See [ROADMAP.md](ROADMAP.md) for planned features and future commands.

#### `dg ws`

```
dg ws
dg workspace   # alias
```

Unified full-screen TUI dashboard — the single entry point to the worktree/session UI (the
old `dg wt ui` subcommand has been removed). Scoped to **workspaces** rather than worktrees
only. Every top-level row in the dashboard is exactly one of two kinds:

- **Repo workspace** (worktree-backed): a repo with git worktrees, sourced from the same
  worktree scan `dg wt list` uses. Expandable to its worktree rows via `h`/`l` (or `z`
  to toggle every repo at once), shown with a `▼`/`▶` chevron and an `N trees` badge. Shown even
  when its repo-slug tmux session isn't live.
- **Session workspace**: a standalone tmux session with no worktree-backed (`wt-`) window,
  sourced from `tmux list-sessions`. A leaf row, labeled `session`.

The two kinds carry different marker shapes so they're distinguishable at a glance, not just
by their label: worktree rows use a circle (`●` running / `○` not), session rows use a square
(`■` attached / `□` detached) — in both, a filled glyph means active and the color matches
(green active, dim inactive). Both kinds share the existing worktree-row keys (`j`/`k` nav,
`h`/`l` fold, `z` toggle-all, `n`/`N` create a worktree, `/` filter, `?` help, `q` quit).
Session rows add:

- `enter` — switch the attached tmux client to the session (guarded: only works inside tmux,
  same guard message as attaching to a worktree) and quit the dashboard.
- `d` `d` — kill the session (two-press confirm, same "press again" hint style as worktree
  delete).
- `s` (works from any row, not just a session row) — opens a two-step create flow:
  1. **Pick a folder** — a fuzzy picker with `root` (the user's home `~`) pinned at the top,
     then the same ranked repo candidates the worktree flow offers, and — like that flow — a
     free-typed path is also accepted. A session isn't tied to a repo, so the chosen folder is
     validated as an existing directory only (not a git repo).
  2. **Name prompt** — on Enter with a name, creates the session in the chosen folder. Enter with
     a **blank** name auto-generates a `devgita-<character>` name (Dragon Ball characters, e.g.
     `devgita-goku`), checked against the live tmux sessions so a blank-name create never collides
     with an existing one.

  Inside tmux, the client switches to the new session and the dashboard quits; outside tmux, the
  session is created detached and reported (`session created: <name>`) without switching. A
  duplicate typed name surfaces tmux's own "duplicate session" error on the status line — there's
  no separate pre-check.

- `D`/`r` are worktree-only actions and are no-ops on a session row.

Bare `ctrl+t` (no tmux prefix) opens `dg ws` (see `configs/tmux/tmux.conf`) — it previously
opened tmux's native `choose-tree -Zs` popup, which this replaces. This is the only key bound
to the dashboard.

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
