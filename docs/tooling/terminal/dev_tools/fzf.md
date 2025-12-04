# Fzf Module Documentation

## Overview

The Fzf module provides installation and command execution management for fzf fuzzy finder with devgita integration. It follows the standardized devgita app interface while providing fzf-specific operations for interactive file searching, command history filtering, directory navigation, and process selection.

## App Purpose

Fzf (fuzzy finder) is a general-purpose command-line fuzzy finder written in Go. It's an interactive Unix filter for command-line that can be used with any list: files, command history, processes, hostnames, bookmarks, git commits, etc. This module ensures fzf is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems and provides high-level operations for launching the interactive fuzzy search interface.

## Lifecycle Summary

1. **Installation**: Install fzf package via platform package managers (Homebrew/apt)
2. **Configuration**: Fzf uses optional environment variables and shell integration (no default configuration applied by devgita)
3. **Execution**: Provide high-level fzf operations for launching the fuzzy finder and integrating with shell workflows

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Fzf instance with platform-specific commands             |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install fzf                               |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute fzf               | Runs fzf with provided arguments                                     |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
fzf := fzf.New()
err := fzf.Install()
```

- **Purpose**: Standard fzf installation
- **Behavior**: Uses `InstallPackage()` to install fzf package
- **Use case**: Initial fzf installation or explicit reinstall

### ForceInstall()

```go
fzf := fzf.New()
err := fzf.ForceInstall()
```

- **Purpose**: Force fzf installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh fzf installation or fix corrupted installation

### SoftInstall()

```go
fzf := fzf.New()
err := fzf.SoftInstall()
```

- **Purpose**: Install fzf only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing fzf installations

### Uninstall()

```go
err := fzf.Uninstall()
```

- **Purpose**: Remove fzf installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Fuzzy finder tools are typically managed at the system level

### Update()

```go
err := fzf.Update()
```

- **Purpose**: Update fzf installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Fzf updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := fzf.ForceConfigure()
err := fzf.SoftConfigure()
```

- **Purpose**: Apply fzf configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: Fzf configuration is optional and handled via environment variables and shell integration. Configuration is typically done through `FZF_*` environment variables and shell key bindings rather than config file copying.

## Execution Methods

### ExecuteCommand()

```go
err := fzf.ExecuteCommand()                         // Launch interactive finder
err := fzf.ExecuteCommand("--version")              // Show version
err := fzf.ExecuteCommand("--query", "test")        // Start with query
err := fzf.ExecuteCommand("--reverse")              // Reverse layout
```

- **Purpose**: Execute fzf commands with provided arguments
- **Parameters**: Variable arguments passed directly to fzf binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### Fzf-Specific Operations

The fzf CLI provides extensive interactive fuzzy finding capabilities:

#### Basic Usage

```bash
# Interactive file finder
find . -type f | fzf

# Interactive directory finder
find . -type d | fzf

# Command history search
history | fzf

# Process selector
ps aux | fzf

# Git commit browser
git log --oneline | fzf
```

#### Command-Line Options

```bash
# Show version information
fzf --version

# Start with initial query
fzf --query "search term"

# Reverse layout (top to bottom)
fzf --reverse

# Multi-select mode
fzf --multi

# Enable preview window
fzf --preview 'cat {}'

# Custom prompt
fzf --prompt "Select file: "

# Filter mode (non-interactive)
echo -e "one\ntwo\nthree" | fzf --filter "tw"
```

#### Layout and Appearance

```bash
# Set height
fzf --height 40%

# Border styles
fzf --border
fzf --border=rounded
fzf --border=sharp

# Color scheme
fzf --color=dark
fzf --color=light
fzf --color='fg:188,bg:233,hl:103,fg+:222,bg+:234'

# Info style
fzf --info=inline
fzf --info=hidden
```

#### Preview Window

```bash
# Enable preview with command
fzf --preview 'cat {}'
fzf --preview 'head -n 50 {}'

# Preview window position
fzf --preview-window=right:50%
fzf --preview-window=down:40%
fzf --preview-window=hidden

# Syntax highlighting in preview
fzf --preview 'bat --color=always {}'
```

#### Search Behavior

```bash
# Exact match mode
fzf --exact

# Case-sensitive search
fzf +i

# Search algorithm
fzf --algo=v2
fzf --algo=v1

# Sorting
fzf --no-sort
fzf --tac  # Reverse input order
```

#### Key Bindings

```bash
# Execute command on selection
fzf --bind 'enter:execute(vim {})'

# Custom key bindings
fzf --bind 'ctrl-y:execute-silent(echo {} | pbcopy)'
fzf --bind 'ctrl-/:toggle-preview'

# Multiple actions
fzf --bind 'ctrl-a:select-all+accept'
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **Launch Finder**: `New()` → `SoftInstall()` → `ExecuteCommand()`
4. **Version Check**: `New()` → `ExecuteCommand("--version")`

## Constants and Paths

### Relevant Constants

- **Package name**: Referenced via `constants.Fzf` (typically "fzf")
- Used by all installation methods for consistent package reference

### Configuration Approach

- **Optional configuration**: Fzf uses optional environment variables and shell integration
- **Environment variables**: `FZF_DEFAULT_COMMAND`, `FZF_DEFAULT_OPTS`, `FZF_CTRL_T_COMMAND`, etc.
- **No default config**: Devgita does not apply default configuration for fzf
- **User customization**: Users can configure fzf via shell rc files and environment variables

### Environment Variables

While devgita doesn't apply default configuration, users can customize fzf behavior:

```bash
# Default command for finding files
export FZF_DEFAULT_COMMAND='fd --type f'

# Default options
export FZF_DEFAULT_OPTS='--height 40% --reverse --border'

# CTRL-T command
export FZF_CTRL_T_COMMAND="$FZF_DEFAULT_COMMAND"

# CTRL-T options
export FZF_CTRL_T_OPTS='--preview "bat --color=always --style=numbers {}"'

# ALT-C command
export FZF_ALT_C_COMMAND='fd --type d'

# ALT-C options
export FZF_ALT_C_OPTS='--preview "tree -C {} | head -200"'
```

## Implementation Notes

- **CLI Tool Nature**: Fzf is an interactive command-line tool without traditional config file requirements
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since fzf uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since fzf uses environment variables for configuration
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as fzf updates should be handled by system package managers

## Usage Examples

### Basic Fzf Operations

```go
fzf := fzf.New()

// Install fzf
err := fzf.SoftInstall()
if err != nil {
    return err
}

// Launch interactive finder (reads from stdin)
err = fzf.ExecuteCommand()

// Check version
err = fzf.ExecuteCommand("--version")
```

### Advanced Usage

```go
// Start with query
err := fzf.ExecuteCommand("--query", "search")

// Multi-select mode
err = fzf.ExecuteCommand("--multi")

// With preview
err = fzf.ExecuteCommand("--preview", "cat {}")

// Reverse layout
err = fzf.ExecuteCommand("--reverse", "--border")
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Finder Won't Launch**: Verify input is provided via stdin or file argument
3. **Preview Not Working**: Check if preview command is available (e.g., `cat`, `bat`)
4. **Commands Don't Work**: Verify fzf is installed and accessible in PATH
5. **Shell Integration Missing**: Install fzf shell key bindings separately

### Platform Considerations

- **macOS**: Installed via Homebrew package manager
- **Linux**: Installed via apt package manager
- **Shell Integration**: Shell key bindings (CTRL-T, CTRL-R, ALT-C) require separate setup
- **Terminal Support**: Works best with modern terminal emulators that support colors

### Prerequisites

Before using fzf effectively:
- Shell integration for key bindings (optional)
- Preview tools like `bat`, `tree`, or `cat` for enhanced previews
- Terminal with color support for optimal visual experience

### Shell Integration

Fzf provides useful shell key bindings:

```bash
# CTRL-T: Paste selected files/directories into command line
# CTRL-R: Search command history
# ALT-C: cd into selected directory
```

To enable these, add to shell configuration:

```bash
# For bash
eval "$(fzf --bash)"

# For zsh
source <(fzf --zsh)

# For fish
fzf --fish | source
```

## Common Use Cases

### File Navigation

```bash
# Open file in vim
vim $(fzf)

# cd into directory
cd $(find . -type d | fzf)

# Preview files before opening
fzf --preview 'bat --color=always {}'
```

### Git Integration

```bash
# Checkout git branch
git checkout $(git branch | fzf)

# View git commit
git log --oneline | fzf --preview 'git show {1}'

# Interactive git add
git status --short | fzf --multi | awk '{print $2}' | xargs git add
```

### Process Management

```bash
# Kill process interactively
ps aux | fzf | awk '{print $2}' | xargs kill

# Monitor processes
watch "ps aux | fzf"
```

### Command History

```bash
# Search and execute command from history
eval $(history | fzf | sed 's/^[ 0-9]*//')
```

## External References

- **Fzf Repository**: https://github.com/junegunn/fzf
- **Installation Guide**: https://github.com/junegunn/fzf#installation
- **Advanced Usage**: https://github.com/junegunn/fzf/wiki
- **Examples**: https://github.com/junegunn/fzf/wiki/examples

## Integration with Devgita

Fzf integrates with devgita's terminal category:

- **Installation**: Installed as part of terminal dev tools setup
- **Configuration**: User-managed environment variables (no default applied)
- **Usage**: Available system-wide after installation
- **Updates**: Managed through system package manager
- **Shell Integration**: User configures key bindings separately if desired

This module provides essential fuzzy finding capabilities for interactive file searching, command history filtering, and development workflow enhancement within the devgita ecosystem.

