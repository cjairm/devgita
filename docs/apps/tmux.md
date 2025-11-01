# Tmux Module Documentation

## Overview

The Tmux module provides installation and configuration management for tmux terminal multiplexer with devgita integration. It follows the standardized devgita app interface while providing tmux-specific operations for session management, window control, and development environment customization.

## App Purpose

Tmux is a terminal multiplexer that allows you to switch easily between several programs in one terminal, detach them (they keep running in the background) and reattach them to a different terminal. This module ensures tmux is properly installed and configured with devgita's optimized settings for development workflows, including custom key bindings, vim integration, and enhanced navigation.

## Lifecycle Summary

1. **Installation**: Install tmux package via platform package managers (Homebrew/apt)
2. **Configuration**: Apply devgita's tmux configuration template with developer-focused settings
3. **Execution**: Provide high-level tmux operations for session and window management

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Tmux instance with platform-specific commands            |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install tmux                              |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | Overwrites existing .tmux.conf with devgita defaults                 |
| `SoftConfigure()`  | Conditional configuration | Preserves existing .tmux.conf if present                             |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute tmux commands     | Runs tmux with provided arguments                                     |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
tmux := tmux.New()
err := tmux.Install()
```

- **Purpose**: Standard tmux installation
- **Behavior**: Uses `InstallPackage()` to install tmux package
- **Use case**: Initial tmux installation or explicit reinstall

### ForceInstall()

```go
tmux := tmux.New()
err := tmux.ForceInstall()
```

- **Purpose**: Force tmux installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh tmux installation or fix corrupted installation

### SoftInstall()

```go
tmux := tmux.New()
err := tmux.SoftInstall()
```

- **Purpose**: Install tmux only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing tmux installations

### Uninstall()

```go
err := tmux.Uninstall()
```

- **Purpose**: Remove tmux installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Terminal multiplexers are typically managed at the system level

### Update()

```go
err := tmux.Update()
```

- **Purpose**: Update tmux installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Tmux updates are typically handled by the system package manager

## Configuration Methods

### Configuration Paths

- **Source**: `paths.TmuxConfigAppDir/.tmux.conf` (devgita's tmux config)
- **Destination**: `~/.tmux.conf` (user's home directory)
- **Marker file**: `.tmux.conf` in user's home directory

### ForceConfigure()

```go
err := tmux.ForceConfigure()
```

- **Purpose**: Apply tmux configuration regardless of existing files
- **Behavior**: Copies .tmux.conf from app dir to home directory, overwriting existing
- **Use case**: Reset to devgita defaults, apply config updates

### SoftConfigure()

```go
err := tmux.SoftConfigure()
```

- **Purpose**: Apply tmux configuration only if not already configured
- **Behavior**: Checks for existing `.tmux.conf` in home directory; if exists, does nothing
- **Marker logic**: Uses `os.UserHomeDir()` to locate and check for existing configuration
- **Use case**: Initial setup that preserves user customizations

## Execution Methods

### ExecuteCommand()

```go
err := tmux.ExecuteCommand("new-session", "-d", "-s", "dev")
err := tmux.ExecuteCommand("attach-session", "-t", "dev")
err := tmux.ExecuteCommand("list-sessions")
```

- **Purpose**: Execute tmux commands with provided arguments
- **Parameters**: Variable arguments passed directly to tmux binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### Tmux-Specific Operations

The tmux CLI provides extensive session and window management capabilities:

#### Session Management

```bash
# Create new session
tmux new-session -s mysession

# Create detached session
tmux new-session -d -s background

# Attach to session
tmux attach-session -t mysession

# List sessions
tmux list-sessions

# Kill session
tmux kill-session -t mysession
```

#### Window Management

```bash
# Create new window
tmux new-window

# Create window with name
tmux new-window -n "editor"

# List windows
tmux list-windows

# Select window by number
tmux select-window -t 1

# Kill window
tmux kill-window -t 1
```

#### Pane Operations

```bash
# Split window horizontally
tmux split-window -h

# Split window vertically
tmux split-window -v

# Select pane
tmux select-pane -t 0

# Kill pane
tmux kill-pane -t 0

# List panes
tmux list-panes
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()`
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()`
3. **Update Configuration**: `New()` → `SoftInstall()` → `ForceConfigure()`
4. **Tmux Operations**: `New()` → `ExecuteCommand()` with specific tmux arguments

## Constants and Paths

### Relevant Constants

- `constants.Tmux`: Package name ("tmux") for installation
- Used by all installation methods for consistent package reference

### Configuration Paths

- `paths.TmuxConfigAppDir`: Source directory for devgita's tmux configuration template
- `paths.HomeDir`: Target directory for tmux configuration (fallback)
- Home directory resolution: Uses `os.UserHomeDir()` for cross-platform compatibility

## Implementation Notes

- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since tmux uninstall is not supported
- **Configuration Strategy**: Uses file existence check to determine if configuration exists
- **Home Directory**: Uses `os.UserHomeDir()` instead of `paths.HomeDir` for configuration detection
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as tmux updates should be handled by system package managers

## Configuration Structure

The tmux configuration (`.tmux.conf`) includes developer-focused enhancements:

### Key Bindings

```bash
# Change prefix from C-b to C-Space
unbind C-b
set -g prefix C-Space
bind-key C-Space send-prefix

# Reload configuration
unbind r
bind r source-file ~/.tmux.conf

# Pane navigation (vim-style)
bind -n C-h select-pane -L
bind -n C-j select-pane -D
bind -n C-k select-pane -U
bind -n C-l select-pane -R

# Window management
bind v split-window -h -c "#{pane_current_path}"
bind - split-window -v -c "#{pane_current_path}"
bind w new-window -c "#{pane_current_path}"
bind h previous-window
bind l next-window
```

### Session Settings

```bash
# Start window and pane numbering at 1
set -g base-index 1
set-window-option -g pane-base-index 1

# Mouse support
set -g mouse on

# History limit
set -g history-limit 100000

# Escape time (for vim compatibility)
set -s escape-time 0

# Terminal color support
set -ga terminal-overrides ",xterm-256color*:Tc"
```

### Vim Integration

```bash
# Vi mode for copy mode
set-window-option -g mode-keys vi

# Copy mode bindings
bind -T copy-mode-vi v send-keys -X begin-selection
bind -T copy-mode-vi y send-keys -X copy-pipe-and-cancel "xsel --clipboard"

# Smart pane switching with vim awareness
is_vim="ps -o state= -o comm= -t '#{pane_tty}' \
    | grep -iqE '^[^TXZ ]+ +(\\S+\\/)?g?(view|n?vim?x?)(diff)?$'"
bind -n C-h if-shell "$is_vim" "send-keys C-h"  "select-pane -L"
bind -n C-j if-shell "$is_vim" "send-keys C-j"  "select-pane -D"
bind -n C-k if-shell "$is_vim" "send-keys C-k"  "select-pane -U"
bind -n C-l if-shell "$is_vim" "send-keys C-l"  "select-pane -R"
```

### Visual Styling

```bash
# Solarized dark color scheme
set-option -g status-style fg=yellow,bg=black
set-window-option -g window-status-style fg=brightblue,bg=default
set-window-option -g window-status-current-style fg=brightred,bg=default
set-option -g pane-border-style fg=black
set-option -g pane-active-border-style fg=brightgreen
set-option -g message-style fg=brightred,bg=black
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Configuration Not Applied**: Check file permissions in home directory
3. **Commands Don't Work**: Verify tmux is installed and accessible in PATH
4. **Key Bindings Not Working**: Restart tmux or source configuration with `tmux source-file ~/.tmux.conf`
5. **Vim Integration Issues**: Ensure vim-tmux-navigator plugin is installed in vim/neovim

### Platform Considerations

- **macOS**: Installed via Homebrew package manager
- **Linux**: Installed via apt package manager
- **Configuration Location**: Always in user's home directory (`~/.tmux.conf`)
- **Copy/Paste**: Linux uses `xsel --clipboard`, macOS behavior may differ

### Key Features

- **Prefix Key**: Changed from `C-b` to `C-Space` for easier access
- **Vim Navigation**: Full vim-style navigation between panes
- **Mouse Support**: Enabled for modern terminal usage
- **Visual Mode**: Vi-style copy mode with system clipboard integration
- **Smart Splitting**: New panes/windows inherit current path
- **Session Tree**: Quick session/window browser with `C-t`

## External References

- **Tmux Documentation**: https://github.com/tmux/tmux
- **Personal Configuration**: https://github.com/cjairm/devenv/tree/main/tmux
- **Releases**: https://github.com/tmux/tmux/releases
- **Installation Guide**: https://github.com/tmux/tmux/wiki/Installing
- **Vim Integration**: https://github.com/christoomey/vim-tmux-navigator

This module provides essential terminal multiplexing capabilities for creating a powerful, vim-integrated development environment within the devgita ecosystem.