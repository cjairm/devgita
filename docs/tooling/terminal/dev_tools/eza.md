# Eza Module Documentation

## Overview

The Eza module provides modern ls replacement tool installation and command execution with devgita integration. It follows the standardized devgita app interface while providing eza-specific operations for enhanced directory listing, file information display, and tree visualization.

## App Purpose

Eza is a modern, maintained replacement for the venerable `ls` command, with more features and better defaults. It uses colours to distinguish file types and metadata, knows about symlinks, extended attributes, and git. This module ensures eza is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems and provides high-level operations for common directory listing and file browsing tasks.

## Lifecycle Summary

1. **Installation**: Install eza package via platform package managers (Homebrew/apt)
2. **Configuration**: Eza typically doesn't require separate configuration files - operations are handled via command-line arguments and shell aliases
3. **Execution**: Provide high-level eza operations for directory listing, file information, and tree views

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Eza instance with platform-specific commands             |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install eza                               |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute eza commands      | Runs eza with provided arguments                                      |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
eza := eza.New()
err := eza.Install()
```

- **Purpose**: Standard eza installation
- **Behavior**: Uses `InstallPackage()` to install eza package
- **Use case**: Initial eza installation or explicit reinstall

### ForceInstall()

```go
eza := eza.New()
err := eza.ForceInstall()
```

- **Purpose**: Force eza installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh eza installation or fix corrupted installation

### SoftInstall()

```go
eza := eza.New()
err := eza.SoftInstall()
```

- **Purpose**: Install eza only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing eza installations

### Uninstall()

```go
err := eza.Uninstall()
```

- **Purpose**: Remove eza installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Directory listing tools are typically managed at the system level

### Update()

```go
err := eza.Update()
```

- **Purpose**: Update eza installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Eza updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := eza.ForceConfigure()
err := eza.SoftConfigure()
```

- **Purpose**: Apply eza configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: Eza doesn't use traditional config files; operation parameters are passed via command-line arguments and shell aliases

## Execution Methods

### ExecuteCommand()

```go
err := eza.ExecuteCommand("--version")
err := eza.ExecuteCommand("-l", "-a", "-h")
err := eza.ExecuteCommand("-T", "--git-ignore")
```

- **Purpose**: Execute eza commands with provided arguments
- **Parameters**: Variable arguments passed directly to eza binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### Eza-Specific Operations

The eza CLI provides extensive directory listing and file information capabilities:

#### Basic Listing

```bash
# List directory contents with colors
eza

# Long format listing
eza -l

# Show hidden files
eza -a

# Long format with all files
eza -la

# Human-readable file sizes
eza -lh

# Sort by modification time
eza -l --sort=modified
```

#### Enhanced Features

```bash
# Show file icons
eza --icons

# Show git status
eza --git

# Show extended attributes
eza -l@

# Show file headers
eza -lh --header

# Group directories first
eza -l --group-directories-first
```

#### Tree Views

```bash
# Tree view
eza -T

# Tree view with depth limit
eza -T --level=2

# Tree view with git status
eza -T --git

# Tree view ignoring git ignored files
eza -T --git-ignore
```

#### Color and Display

```bash
# Force color output
eza --color=always

# No color output
eza --color=never

# Show only directories
eza -D

# Show only files
eza -f

# One file per line
eza -1
```

#### Time and Sorting

```bash
# Sort by name
eza -l --sort=name

# Sort by size
eza -l --sort=size

# Sort by modification time (newest first)
eza -l --sort=modified --reverse

# Show access time
eza -l --time=accessed

# Show creation time
eza -l --time=created
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **Listing Operations**: `New()` → `SoftInstall()` → `ExecuteCommand()` with eza arguments
4. **Tree Operations**: `New()` → `ExecuteCommand()` with tree parameters

## Constants and Paths

### Relevant Constants

- **Package name**: Referenced via `constants.Eza` (typically "eza")
- Used by all installation methods for consistent package reference

### Configuration Approach

- **No traditional config files**: Eza operations are configured via command-line arguments
- **Shell aliases**: Users typically configure shell aliases for common eza commands
- **Runtime configuration**: Parameters passed directly to `ExecuteCommand()`

### Common Shell Aliases

While devgita doesn't create these automatically, users often configure:

```bash
# Common eza aliases
alias ls='eza'
alias ll='eza -l'
alias la='eza -la'
alias lt='eza -T'
alias l='eza -lah'
alias lg='eza -l --git'
```

## Implementation Notes

- **CLI Tool Nature**: Unlike typical applications, eza is a command-line directory listing tool without persistent configuration
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since eza uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since eza doesn't use config files
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as eza updates should be handled by system package managers

## Usage Examples

### Basic Directory Listing

```go
eza := eza.New()

// Install eza
err := eza.SoftInstall()
if err != nil {
    return err
}

// List directory contents
err = eza.ExecuteCommand()

// Long format listing
err = eza.ExecuteCommand("-l")

// Show hidden files with icons
err = eza.ExecuteCommand("-la", "--icons")
```

### Advanced Operations

```go
// Check eza version
err := eza.ExecuteCommand("--version")

// Tree view with git status
err = eza.ExecuteCommand("-T", "--git")

// Long format with git status and icons
err = eza.ExecuteCommand("-l", "--git", "--icons", "--header")

// Human-readable sizes, sorted by modification time
err = eza.ExecuteCommand("-lh", "--sort=modified", "--reverse")
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Icons Not Showing**: Requires a Nerd Font to be installed and configured in terminal
3. **Git Status Not Working**: Requires git repository context
4. **Color Issues**: Terminal may not support 256 colors or true color
5. **Commands Don't Work**: Verify eza is installed and accessible in PATH

### Platform Considerations

- **macOS**: Installed via Homebrew package manager
- **Linux**: Installed via apt package manager (may require additional repository configuration)
- **Icons Support**: Requires Nerd Font installation (e.g., JetBrains Mono Nerd Font)
- **Git Integration**: Works in directories with git repositories

### Feature Requirements

- **Icons**: Requires `--icons` flag and Nerd Font terminal font
- **Git Status**: Requires git repository and `--git` flag
- **Extended Attributes**: Platform-specific (macOS: `@`, Linux: varies)
- **Color Output**: Most modern terminals support color; use `--color=always` to force

## Key Features

### Visual Enhancements
- **Colored Output**: File types distinguished by color
- **Icons**: Optional file type icons (requires Nerd Font)
- **Grid Display**: Automatic column sizing
- **Headers**: Optional column headers in long format

### Git Integration
- **Repository Status**: Shows modified, staged, and untracked files
- **Branch Awareness**: Recognizes git repository context
- **Ignore Support**: Respects `.gitignore` patterns with `--git-ignore`

### Flexible Display
- **Multiple Formats**: Grid, long, one-per-line, tree
- **Sorting Options**: Name, size, time, extension
- **Filtering**: Files only, directories only, patterns
- **Tree Views**: Recursive directory visualization

### Performance
- **Parallel Processing**: Fast directory traversal
- **Efficient Memory**: Low resource usage
- **Large Directories**: Handles thousands of files efficiently

## Migration from `ls`

Eza is designed as a drop-in replacement for `ls`:

### Basic Equivalents

| ls command | eza equivalent    | Description                 |
| ---------- | ----------------- | --------------------------- |
| `ls`       | `eza`             | Basic listing               |
| `ls -l`    | `eza -l`          | Long format                 |
| `ls -la`   | `eza -la`         | Long with hidden files      |
| `ls -lh`   | `eza -lh`         | Human-readable sizes        |
| `ls -R`    | `eza -T`          | Recursive (tree view)       |
| `ls -t`    | `eza --sort=time` | Sort by modification time   |

### Enhanced Features

Eza adds features beyond `ls`:

- **Git integration** (`--git`)
- **Icons** (`--icons`)
- **Tree view** (`-T` or `--tree`)
- **Extended attributes** (`@`)
- **Better color schemes**
- **More flexible sorting** (`--sort=<field>`)

## External References

- **Eza Repository**: https://github.com/eza-community/eza
- **Installation Guide**: https://github.com/eza-community/eza/blob/main/INSTALL.md
- **Command Reference**: https://github.com/eza-community/eza/blob/main/man/eza.1.md
- **Comparison with ls**: https://github.com/eza-community/eza#command-line-options

## Integration with Devgita

Eza integrates with devgita's terminal category:

- **Installation**: Installed as part of terminal tools setup
- **Configuration**: User-managed via shell aliases (no default applied)
- **Usage**: Available system-wide after installation
- **Updates**: Managed through system package manager
- **Fonts**: Works best with Nerd Fonts (also in devgita's fonts category)

This module provides essential directory listing and file browsing capabilities with modern enhancements, significantly improving command-line productivity within the devgita ecosystem.
