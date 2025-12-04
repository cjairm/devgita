# Zoxide Module Documentation

## Overview

The Zoxide module provides smart directory navigation tool installation and command execution with devgita integration. It follows the standardized devgita app interface while providing zoxide-specific operations for intelligent directory jumping, frequency-based navigation, and fuzzy path matching.

## App Purpose

Zoxide is a smarter cd command that learns your habits and allows you to navigate to frequently and recently used directories with just a few keystrokes. It tracks your most used directories using a ranking algorithm (frecency) and provides fuzzy matching for quick navigation. This module ensures zoxide is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems and provides high-level operations for smart directory navigation and database management.

## Lifecycle Summary

1. **Installation**: Install zoxide package via platform package managers (Homebrew/apt)
2. **Configuration**: Zoxide doesn't require traditional configuration files - shell integration is handled via `zoxide init` command
3. **Execution**: Provide high-level zoxide operations for directory navigation, database queries, and shell integration

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Zoxide instance with platform-specific commands          |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install zoxide                            |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute zoxide commands   | Runs zoxide with provided arguments                                  |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
zoxide := zoxide.New()
err := zoxide.Install()
```

- **Purpose**: Standard zoxide installation
- **Behavior**: Uses `InstallPackage()` to install zoxide package
- **Use case**: Initial zoxide installation or explicit reinstall

### ForceInstall()

```go
zoxide := zoxide.New()
err := zoxide.ForceInstall()
```

- **Purpose**: Force zoxide installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh zoxide installation or fix corrupted installation

### SoftInstall()

```go
zoxide := zoxide.New()
err := zoxide.SoftInstall()
```

- **Purpose**: Install zoxide only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing zoxide installations

### Uninstall()

```go
err := zoxide.Uninstall()
```

- **Purpose**: Remove zoxide installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Directory navigation tools are typically managed at the system level

### Update()

```go
err := zoxide.Update()
```

- **Purpose**: Update zoxide installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Zoxide updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := zoxide.ForceConfigure()
err := zoxide.SoftConfigure()
```

- **Purpose**: Apply zoxide configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: Zoxide doesn't use traditional config files; shell integration is handled via `zoxide init` command which generates shell-specific initialization scripts

## Execution Methods

### ExecuteCommand()

```go
err := zoxide.ExecuteCommand("--version")
err := zoxide.ExecuteCommand("query", "projects")
err := zoxide.ExecuteCommand("add", "/path/to/directory")
```

- **Purpose**: Execute zoxide commands with provided arguments
- **Parameters**: Variable arguments passed directly to zoxide binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### Zoxide-Specific Operations

The zoxide CLI provides extensive directory navigation and database management capabilities:

#### Directory Query Operations

```bash
# Query for directories matching keywords
zoxide query projects

# Interactive query (select from multiple matches)
zoxide query --interactive projects

# List all tracked directories
zoxide query --list

# List with scores (shows frecency ranking)
zoxide query --list --score
```

#### Database Management

```bash
# Add directory to database
zoxide add /path/to/directory

# Remove directory from database
zoxide remove /path/to/directory

# Import directories from file
zoxide import ~/.z  # Import from z/fasd database

# Clean up removed directories
zoxide remove --interactive
```

#### Shell Integration

```bash
# Generate shell initialization script
zoxide init zsh     # For Zsh
zoxide init bash    # For Bash
zoxide init fish    # For Fish
zoxide init powershell  # For PowerShell

# With custom hooks
zoxide init zsh --hook prompt
zoxide init bash --hook pwd

# Generate with command alias
zoxide init zsh --cmd j  # Use 'j' instead of 'z'
```

#### Version and Help

```bash
# Show version information
zoxide --version

# Display help
zoxide --help
zoxide query --help
zoxide add --help
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **Navigation Operations**: `New()` → `SoftInstall()` → `ExecuteCommand()` with zoxide arguments
4. **Shell Integration**: `New()` → `ExecuteCommand("init", "zsh")` to generate init script

## Constants and Paths

### Relevant Constants

- **Package name**: `"zoxide"` used directly for installation (referenced via `constants.Zoxide`)
- Used by all installation methods for consistent package reference

### Configuration Approach

- **No traditional config files**: Zoxide doesn't use configuration files
- **Shell integration**: Configured via shell initialization scripts generated by `zoxide init`
- **Database location**: Zoxide stores its database at:
  - macOS/Linux: `$XDG_DATA_HOME/zoxide` or `~/.local/share/zoxide`
  - Windows: `%LOCALAPPDATA%\zoxide`
- **Environment variables**: 
  - `_ZO_DATA_DIR`: Override database location
  - `_ZO_ECHO`: Print matched directory before navigating
  - `_ZO_EXCLUDE_DIRS`: Directories to exclude from tracking
  - `_ZO_FZF_OPTS`: Custom fzf options for interactive mode
  - `_ZO_MAXAGE`: Maximum age of database entries
  - `_ZO_RESOLVE_SYMLINKS`: Follow symbolic links

## Implementation Notes

- **Smart Navigation Nature**: Unlike typical applications, zoxide is a directory navigation tool without traditional configuration
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since zoxide uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since zoxide uses shell initialization
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as zoxide updates should be handled by system package managers

## Usage Examples

### Basic Installation and Setup

```go
zoxide := zoxide.New()

// Install zoxide
err := zoxide.SoftInstall()
if err != nil {
    return err
}

// Check version
err = zoxide.ExecuteCommand("--version")

// Generate shell initialization script
err = zoxide.ExecuteCommand("init", "zsh")
```

### Database Operations

```go
// Query for directories
err := zoxide.ExecuteCommand("query", "projects")

// Add directory to database
err = zoxide.ExecuteCommand("add", "/home/user/workspace")

// List all tracked directories
err = zoxide.ExecuteCommand("query", "--list")

// Interactive selection
err = zoxide.ExecuteCommand("query", "--interactive", "docs")
```

### Shell Integration Setup

After installing zoxide, users should add the following to their shell configuration:

**Zsh (`~/.zshrc`):**
```bash
eval "$(zoxide init zsh)"
```

**Bash (`~/.bashrc`):**
```bash
eval "$(zoxide init bash)"
```

**Fish (`~/.config/fish/config.fish`):**
```fish
zoxide init fish | source
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Commands Don't Work**: Verify zoxide is installed and accessible in PATH
3. **Shell Integration Not Working**: Ensure shell initialization script is properly sourced
4. **Database Not Updating**: Check that zoxide hook is properly installed in shell
5. **Permissions Issues**: Verify write access to database directory

### Platform Considerations

- **macOS**: Installed via Homebrew package manager
- **Linux**: Installed via apt package manager or from releases
- **Windows**: Can be installed via Scoop or Chocolatey
- **Database Location**: Varies by platform and XDG Base Directory specification

### Shell Integration

After installation, zoxide requires shell integration to work properly:

- **Automatic tracking**: Shell hook tracks directory changes and updates database
- **'z' command**: Jump to directories with `z <keywords>`
- **'zi' command**: Interactive selection with `zi <keywords>`
- **Tab completion**: Shell-specific completion for directory names

### Frecency Algorithm

Zoxide ranks directories using "frecency" (frequency + recency):

- **Frequency**: How often you visit a directory
- **Recency**: How recently you visited a directory
- **Combined score**: Higher scores for frequently and recently used directories
- **Automatic cleanup**: Old, unused entries are gradually removed

### Best Practices

1. **Let it learn**: Use zoxide regularly to build up the database
2. **Use fuzzy matching**: Don't need exact directory names, partial matches work
3. **Interactive mode**: Use `zi` when multiple matches are possible
4. **Exclude directories**: Set `_ZO_EXCLUDE_DIRS` to ignore certain paths
5. **Custom aliases**: Use `zoxide init --cmd <alias>` to customize command name

## Integration with Devgita

Zoxide integrates with devgita's terminal category:

- **Installation**: Installed as part of terminal tools setup
- **Configuration**: Shell integration via user's shell config files
- **Usage**: Available system-wide after installation and shell configuration
- **Updates**: Managed through system package manager

## External References

- **Zoxide Repository**: https://github.com/ajeetdsouza/zoxide
- **Installation Guide**: https://github.com/ajeetdsouza/zoxide#installation
- **Configuration Options**: https://github.com/ajeetdsouza/zoxide/wiki/Configuration
- **Shell Integration**: https://github.com/ajeetdsouza/zoxide/wiki/Shell-integration
- **Comparison with z/autojump**: https://github.com/ajeetdsouza/zoxide#comparison

This module provides essential smart directory navigation capabilities for development workflows, significantly improving command-line productivity by reducing time spent navigating directory structures within the devgita ecosystem.
