# Devgita Roadmap

Public roadmap of planned features, improvements, and discussion topics for future releases.

---

## 🔴 Priority: Neovim Dependencies

Add the following packages as terminal tools that must be installed for Neovim to function properly:

- `make`
- `gcc`
- `ripgrep`
- `fd-find` (APT) / `fd` (Homebrew)
- `tree-sitter-cli`
- `unzip`
- `xclip` (Linux only; macOS uses `pbcopy`/`pbpaste`)
- `neovim`

---

## Implemented Commands

### Configuration & Management

- 🟢 **`dg configure [app]`** — Re-applies configuration files for a named app without reinstalling
  - `--force` overwrites existing config files; default (soft mode) only applies if files are absent
  - Supports 19 apps; apps with no config return an info message and exit zero
  - Shipped in v0.10.0

- 🟢 **`dg uninstall [app/category]`** — Remove installed packages
  - Verifies packages were installed by us (not pre-existing)
  - Handles updates differently (leaves pre-existing packages in place)
  - Scope: Can pass `--app` or `--package` for single-package uninstall

---

## Planned Commands

The following commands are planned but not yet implemented:

### Configuration & Management

- **`dg list` / `dg installed`** — Show what's installed
  - Display installed packages with versions and timestamps
  - Show category grouping

- **`dg update [app] [options]`** — Check and apply updates
  - Complex due to breaking changes
  - Support app-specific flags: `dg update --neovim=[options] --aerospace=[options]`
  - Must handle version compatibility

- **`dg check-updates`** — Find available updates for all installed packages
  - Report which apps have updates available

- **`dg validate`** — Check configuration validity
  - Verify current configuration is valid
  - Check if all dependencies are met

### Customization

- **`dg change --theme=[options] --font=[options]`** — Modify environment
  - Change terminal theme (with background image updates)
  - Change font selection
  - Persist selections in global config

### Backup & Recovery

- **`dg backup [name]`** — Backup current configurations
  - Create snapshots of current installation state
  - Enable rollback workflows

- **`dg restore [backup]`** — Revert to previous configuration
  - Restore from saved backup point

### Tmux Utilities

- **`dg tmux --new-window="~/path/to/project"`** — Create tmux window shortcuts
  - Auto-generates window aliases (e.g., `tmhello-world` for `~/my-path/hello-world`)
  - Removes need for custom commands like `tmn`

### Worktree Enhancements

- Improved UX for `dg worktree` management
- Better navigation and selection workflows
- Related: See active cycle `2026-04-22-worktree-ux-improvements.md`

---

## Planned Categories & Tools

### AI & Development Tools

**New Category: `ai-tools`**

- **Ollama** — Local LLM inference engine
- Gemini CLI integration
- IDE AI integrations (Claude Code, Opencode configuration)
- Terminal app for running AI commands (TBD)

**Related Questions:**

- Should this be a separate category or extend existing categories?
- How to integrate with Opencode / Claude Code setup?

### Alternative Tools

Users should be able to choose between alternatives:

- **Multiplexer:** Tmux vs. Zellij
- **IDE:** Opencode vs. Claude Code (Claude Code preferred)
- **Window Manager:** Existing options, possibly add more

### Additional Apps (Backlog)

- **Postman** — API testing
- **DevPod** — Dev containers with SSH/Vim support
- **Snap packages** — Slack, Spotify, VS Code (with `sudo apt install snapd`)
- **Music players** — TBD
- **Email clients** — TBD
- **Language-specific tools:** PHP/Go/Python dev tools
- **System utilities:** Font management (fc-list support)

---

## Known Open Questions

### Development & Architecture

- **Mise management** — What's the best approach? Documentation is difficult to follow
- **TUI improvements** — Should we add [Bubbletea](https://github.com/charmbracelet/bubbletea) + [Lipgloss](https://github.com/charmbracelet/lipgloss) for better UI?
- **Git integration** — Should git-related commands be namespaced as `dg git clean --flags`, `dg git revert`, or just `dg clean-branch`?
- **NPM integration** — Should we offer `dg npm clean` for fresh installs and full cleanup?
- **Update strategy** — How to handle updates when apps can have breaking changes?
- **Shortcut creation** — Do we need to create shell shortcuts during installation?

### Features & Scope

- **Cleanup utilities** — Extend [existing cleanup script](https://github.com/cjairm/devgita/blob/038da72eec456e0a60c50dce2bc9ab615c795fb2/configs/bash/init.zsh#L15) to remove node_modules, docs, examples, unused files?
- **Opencode integrations** — Check [Opencode.nvim](https://github.com/nickjvandyke/opencode.nvim) for nvim integration
- **Super powers** — How to extend [super powers](https://github.com/obra/superpowers) for Opencode?
- **Auto Claude integration** — Reuse or build [Auto Claude](https://github.com/AndyMik90/Auto-Claude)?

### References & Templates

- https://www.aitmpl.com/agents — Agent templates
- https://huggingface.co/pricing — HuggingFace pricing options
- https://a2ui.org/ — AI UI patterns
- https://genai.owasp.org/resource/cheatsheet-a-practical-guide-for-securely-using-third-party-mcp-servers-1-0/ — MCP security

---

## Release Timeline

TBD — To be determined based on:

- Community feedback
- Priority user requests
- Available development resources
- Dependencies (e.g., external library updates)

---

## How to Contribute Ideas

Have a feature request or idea?

1. Check this roadmap — it might already be planned
2. Check [GitHub Issues](https://github.com/cjairm/devgita/issues) — it might be discussed
3. [Open a new issue](https://github.com/cjairm/devgita/issues/new) with the `feature-request` label
4. Or contribute directly — see `CONTRIBUTING.md`

---

## Status Legend

- 🟢 **Implemented** — Available in current release
- 🟡 **In Progress** — Active development in a cycle
- 🔵 **Planned** — Decided, scheduled for development
- ⚪ **Proposed** — Under discussion, not yet committed
- ❌ **Deferred** — Decided not to implement

---
