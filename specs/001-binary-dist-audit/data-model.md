# Data Model: Binary Distribution with Embedded Configs

**Feature**: `001-binary-dist-audit`

## Entities

### GlobalConfig (existing — no schema changes)

Persisted at `~/.config/devgita/global_config.yaml`. Tracks all devgita state.

| Field | Type | Description |
|-------|------|-------------|
| `app_path` | string | Path to devgita app directory (`~/.config/devgita/`) |
| `config_path` | string | Path to devgita config directory |
| `installed.packages` | []string | Packages installed by devgita |
| `installed.desktop_apps` | []string | Desktop apps installed by devgita |
| `installed.fonts` | []string | Fonts installed by devgita |
| `installed.themes` | []string | Themes installed by devgita |
| `installed.terminal_tools` | []string | Terminal tools installed by devgita |
| `installed.dev_languages` | []string | Languages installed by devgita |
| `installed.databases` | []string | Databases installed by devgita |
| `already_installed.*` | []string | Pre-existing items (mirrors installed structure) |
| `current_font` | string | Active font name |
| `current_theme` | string | Active theme name |
| `shortcuts` | map[string]string | Custom shortcuts |
| `shell.is_mac` | bool | **NEW**: Platform flag for template rendering |
| `shell.mise` | bool | Mise shell integration enabled |
| `shell.zoxide` | bool | Zoxide shell integration enabled |
| `shell.zsh_autosuggestions` | bool | Zsh autosuggestions enabled |
| `shell.zsh_syntax_highlighting` | bool | Zsh syntax highlighting enabled |
| `shell.powerlevel10k` | bool | Powerlevel10k enabled |
| `shell.extended_capabilities` | bool | Extended capabilities enabled |

**Note**: The only schema addition is `shell.is_mac`. This is backward-compatible — YAML unmarshaling ignores missing fields and defaults booleans to `false`.

### Embedded Configs (new — compile-time)

The `embed.FS` variable containing the `configs/` directory tree. Not persisted — compiled into binary.

| Directory | Platform | Contains |
|-----------|----------|----------|
| `configs/aerospace/` | macOS only | `aerospace.toml` |
| `configs/alacritty/` | Both | `alacritty.toml.tmpl`, `starter.sh` |
| `configs/fastfetch/` | Both | `config.jsonc` |
| `configs/git/` | Both | `.gitconfig` (NEW) |
| `configs/i3/` | Debian only | `config` (NEW) |
| `configs/neovim/` | Both | Full Lua config tree (18+ files) |
| `configs/opencode/` | Both | Template, themes, agents, commands |
| `configs/templates/` | Both | `devgita.zsh.tmpl`, `global_config.yaml` |
| `configs/tmux/` | Both | `tmux.conf` |

### Platform Equivalent Map (logical — implemented as code)

| Function | macOS App | Debian App |
|----------|-----------|------------|
| Tiling window manager | Aerospace | i3 |
| Application launcher | Raycast | Ulauncher |
| Dev tools (Xcode) | Xcode CLI tools | N/A (skipped) |

## State Transitions

### Config Extraction (new flow)

```
Binary started → Check ~/.config/devgita/configs/ exists?
  → No  → Extract embed.FS to ~/.config/devgita/configs/
  → Yes → Skip extraction (idempotent)
```

### Install Flow (updated)

```
dg install
  → Validate OS version
  → Bootstrap package manager (Homebrew/apt)
  → Extract embedded configs (NEW — replaces git clone)
  → Configure devgita (GlobalConfig, shell config)
  → [if terminal] Install terminal tools + configure
  → [if languages] Choose + install languages
  → [if databases] Choose + install databases
  → [if desktop] Install desktop apps (platform-gated)
```
