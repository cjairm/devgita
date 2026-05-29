# Claude Code (`claude`)

Devgita installs [Claude Code](https://claude.com/claude-code), Anthropic's
terminal AI CLI, as a first-class terminal tool and deploys a curated config to
`~/.claude/`.

- **Module:** `internal/apps/claude/`
- **Config source:** `configs/claude/` (+ shared content in `configs/shared/`)
- **Install:** official script (`curl -fsSL https://claude.ai/install.sh | bash`)

## What gets deployed

`ForceConfigure` copies the following into `~/.claude/`:

| Source                                     | Destination                           | Notes                                 |
| ------------------------------------------ | ------------------------------------- | ------------------------------------- |
| `configs/claude/settings.json`             | `~/.claude/settings.json`             | theme, permissions, statusline, hooks |
| `configs/claude/statusline.sh`             | `~/.claude/statusline.sh`             | `chmod 0755`                          |
| `configs/claude/format.sh`                 | `~/.claude/format.sh`                 | `chmod 0755`                          |
| `configs/claude/themes/`                   | `~/.claude/themes/`                   |                                       |
| `configs/shared/{skills,commands,agents}/` | `~/.claude/{skills,commands,agents}/` | shared with OpenCode                  |

## Formatting & linting (PostToolUse hook)

`settings.json` registers a single `PostToolUse` hook on `Edit|Write` that runs
`~/.claude/format.sh`. After Claude edits a file, the script:

1. Reads the edited path from the hook payload on stdin (`.tool_input.file_path`).
2. **Formats** the file in place with the tools matching its extension.
3. **Lints** the formatted result and feeds any findings back to Claude as
   `hookSpecificOutput.additionalContext` (non-blocking — Claude self-corrects on
   its next turn).

> **Dependency: neovim/Mason.** The hook deliberately **reuses the formatter and
> linter binaries that neovim installs via Mason** at
> `~/.local/share/nvim/mason/bin`, rather than installing a separate toolchain.
> This keeps one source of truth for tool versions across the editor and Claude.
> The trade-off: if neovim (and its Mason tools) is not installed via devgita,
> `format.sh` silently no-ops — each tool is guarded by an executable check, so a
> missing binary is skipped, not an error.

| Extension                                 | Formatters                  | Linter (→ Claude) |
| ----------------------------------------- | --------------------------- | ----------------- |
| `.js .jsx .ts .tsx .mjs .cjs`             | eslint_d --fix, prettier    | eslint_d          |
| `.py`                                     | isort, black                | flake8            |
| `.go`                                     | goimports, gofumpt, golines | golangci-lint     |
| `.md .markdown`                           | prettier                    | —                 |
| `.json .css .scss .less .html .yaml .yml` | prettier                    | —                 |
| `.lua`                                    | stylua                      | —                 |
| `.sh .bash`                               | shfmt                       | —                 |

Notes:

- The hook prints **only** JSON to stdout (Claude parses stdout for JSON on exit
  0); all formatter/linter chatter is routed to `/dev/null`.
- `golangci-lint` adds a few seconds of latency per `.go` edit (longer on first
  run while it builds its cache).
- `eslint_d` on a project without an ESLint config will report a config error as
  a "finding"; this is expected outside configured JS/TS repos.

## Statusline

`statusline.sh` renders model, context bar, git branch/status, rate limits, and
duration. For directories that are **git worktrees** — either a _linked_ worktree
(`git --git-dir` ≠ `--git-common-dir`) **or** any directory under devgita's
worktree base (`$XDG_DATA_HOME/devgita/worktrees`, e.g. standalone clones placed
there) — it shows a compact `wt:<repo>` label instead of the full path, since the
branch already conveys the worktree identity.
