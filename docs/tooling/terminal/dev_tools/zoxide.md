# Zoxide Module Documentation

## Overview

The Zoxide module provides smart directory navigation tool installation and command execution with devgita integration. It follows the standardized devgita app interface while providing zoxide-specific operations for intelligent directory jumping, frequency-based navigation, and fuzzy path matching.

**Note**: This module uses devgita's template-based configuration management system. Configuration changes are tracked in GlobalConfig and shell configuration is regenerated from templates rather than direct file manipulation.

## App Purpose

Zoxide is a smarter cd command that learns your habits and allows you to navigate to frequently and recently used directories with just a few keystrokes. It tracks your most used directories using a ranking algorithm (frecency) and provides fuzzy matching for quick navigation. This module ensures zoxide is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems and provides high-level operations for smart directory navigation and database management.

## Lifecycle Summary

1. **Installation**: Install zoxide package via platform package managers (Homebrew/apt)
2. **Configuration**: Enable zoxide shell integration in GlobalConfig and regenerate shell configuration via templates
3. **Execution**: Provide high-level zoxide operations for directory navigation, database queries, and shell integration

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Zoxide instance with platform-specific commands          |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install zoxide                            |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | Enables shell integration and regenerates shell configuration        |
| `SoftConfigure()`  | Conditional configuration | Checks GlobalConfig; enables only if not already enabled             |
| `Uninstall()`      | Remove shell integration  | **Fully supported** - Disables feature and regenerates shell config  |
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

- **Purpose**: Remove zoxide shell integration
- **Behavior**: Loads GlobalConfig → Disables feature → Regenerates shell config → Saves
- **Use case**: Remove zoxide from devgita shell without uninstalling the package
- **Note**: This disables the shell integration feature but does not remove the package itself

### Update()

```go
err := zoxide.Update()
```

- **Purpose**: Update zoxide installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Zoxide updates are typically handled by the system package manager

## Configuration Methods

### Configuration Strategy

This module uses **template-based configuration management** via GlobalConfig:

- **Template**: `configs/templates/devgita.zsh.tmpl` contains conditional sections
- **GlobalConfig**: `~/.config/devgita/global_config.yaml` tracks enabled features
- **Generated file**: `~/.config/devgita/devgita.zsh` (regenerated from template)
- **Feature tracking**: `shell.zoxide` boolean field in GlobalConfig

### ForceConfigure()

```go
err := zoxide.ForceConfigure()
```

- **Purpose**: Enable zoxide shell integration in shell configuration
- **Behavior**: 
  1. Loads GlobalConfig from disk
  2. Enables `shell.zoxide` feature
  3. Regenerates `devgita.zsh` from template
  4. Saves GlobalConfig back to disk
- **Use case**: Enable zoxide shell integration or re-apply configuration

### SoftConfigure()

```go
err := zoxide.SoftConfigure()
```

- **Purpose**: Enable zoxide shell integration only if not already enabled
- **Behavior**: 
  1. Loads GlobalConfig from disk
  2. Checks if `shell.zoxide` is already enabled
  3. If enabled, returns nil (no operation)
  4. If not enabled, calls `ForceConfigure()`
- **Use case**: Initial setup that preserves existing configuration state

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

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()`
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()`
3. **Navigation Operations**: `New()` → `SoftInstall()` → `SoftConfigure()` → `ExecuteCommand()` with zoxide arguments
4. **Shell Integration**: Automatically enabled via template-based GlobalConfig management
5. **Remove Integration**: `New()` → `Uninstall()`

## Constants and Paths

### Relevant Constants

- `constants.Zoxide`: Package name ("zoxide") used for installation and feature tracking
- Used consistently across all methods for package management and GlobalConfig operations

### Configuration Paths

- `paths.TemplatesAppDir`: Source directory for shell configuration templates
- `paths.AppDir`: Target directory for generated shell configuration
- Template file: `filepath.Join(paths.TemplatesAppDir, "devgita.zsh.tmpl")`
- Generated file: `filepath.Join(paths.AppDir, "devgita.zsh")`
- GlobalConfig file: `~/.config/devgita/global_config.yaml`

### Database Location

Zoxide stores its database at:
- macOS/Linux: `$XDG_DATA_HOME/zoxide` or `~/.local/share/zoxide`
- Windows: `%LOCALAPPDATA%\zoxide`

### Environment Variables

- `_ZO_DATA_DIR`: Override database location
- `_ZO_ECHO`: Print matched directory before navigating
- `_ZO_EXCLUDE_DIRS`: Directories to exclude from tracking
- `_ZO_FZF_OPTS`: Custom fzf options for interactive mode
- `_ZO_MAXAGE`: Maximum age of database entries
- `_ZO_RESOLVE_SYMLINKS`: Follow symbolic links

## Implementation Notes

- **Smart Navigation Nature**: Unlike typical applications, zoxide is a directory navigation tool with shell integration
- **Template-Based Configuration**: Uses GlobalConfig and template regeneration instead of direct file manipulation
- **Load-Modify-Regenerate-Save Pattern**: Each configuration method follows this transaction pattern
- **Fresh GlobalConfig Instances**: Each method creates a new `&config.GlobalConfig{}` and loads from disk to prevent stale data
- **Stateless Configuration**: GlobalConfig represents disk state, not app instance state
- **ForceInstall Logic**: Calls `Uninstall()` first, which now properly disables the shell integration feature
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as zoxide updates should be handled by system package managers

## Template Integration

### Template Structure

The zoxide shell integration is defined in `configs/templates/devgita.zsh.tmpl`:

```bash
{{if .Zoxide}}
# Zoxide - Smarter cd command
if command -v zoxide &> /dev/null; then
  eval "$(zoxide init zsh)"
fi
{{end}}
```

### GlobalConfig Tracking

The feature state is tracked in `~/.config/devgita/global_config.yaml`:

```yaml
shell:
  zoxide: true  # Enabled
  # ... other shell features
```

### Generated Configuration

When enabled, the generated `devgita.zsh` contains:

```bash
# Zoxide - Smarter cd command
if command -v zoxide &> /dev/null; then
  eval "$(zoxide init zsh)"
fi
```

This approach:
- Provides single source of truth (template file)
- Enables clean enable/disable operations
- Prevents configuration conflicts
- Makes tracking and version control easier
- Ensures consistent regeneration

## Usage Examples

### Basic Installation and Setup

```go
zoxide := zoxide.New()

// Install zoxide
err := zoxide.SoftInstall()
if err != nil {
    return err
}

// Enable shell integration
err = zoxide.SoftConfigure()
if err != nil {
    return err
}
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

After installing zoxide and configuring it via devgita, users should source the devgita shell configuration:

**Zsh (`~/.zshrc`):**
```bash
source ~/.config/devgita/devgita.zsh
```

This automatically enables:
- The 'z' command for smart navigation: `z foo`
- The 'zi' command for interactive selection: `zi foo`
- Automatic directory tracking via shell hooks

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Commands Don't Work**: Verify zoxide is installed and accessible in PATH
3. **Shell Integration Not Working**: Ensure devgita shell configuration is properly sourced
4. **Database Not Updating**: Check that zoxide hook is properly installed in shell via template
5. **Permissions Issues**: Verify write access to database directory
6. **GlobalConfig Load Errors**: Ensure `~/.config/devgita/global_config.yaml` is valid YAML
7. **Template Not Found**: Verify `configs/templates/devgita.zsh.tmpl` exists in devgita repository

### Platform Considerations

- **macOS**: Installed via Homebrew package manager
- **Linux**: Installed via apt package manager or from releases
- **Windows**: Can be installed via Scoop or Chocolatey (not officially supported by devgita)
- **Database Location**: Varies by platform and XDG Base Directory specification

### Shell Integration

After installation and configuration, zoxide provides:

- **Automatic tracking**: Shell hook tracks directory changes and updates database
- **'z' command**: Jump to directories with `z <keywords>`
- **'zi' command**: Interactive selection with `zi <keywords>`
- **Tab completion**: Shell-specific completion for directory names
- **Alias support**: Template includes `alias cd="z"` for seamless integration

### Template System Benefits

- **Single Source of Truth**: Template file is the only place shell configuration is defined
- **Trackable**: Git can track template changes easily
- **Predictable**: Regeneration always produces same output for same inputs
- **No Conflicts**: No string manipulation or file appending/removing
- **Clean Uninstall**: Disabling a feature regenerates without it

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
5. **Custom aliases**: Template provides sensible defaults including `cd="z"`

## Integration with Devgita

Zoxide integrates with devgita's terminal category:

- **Installation**: Installed as part of terminal tools setup
- **Configuration**: Shell integration via template-based GlobalConfig management
- **Usage**: Available system-wide after sourcing devgita shell configuration
- **Updates**: Managed through system package manager
- **Enable/Disable**: Controlled via GlobalConfig feature flags

## External References

- **Zoxide Repository**: https://github.com/ajeetdsouza/zoxide
- **Installation Guide**: https://github.com/ajeetdsouza/zoxide#installation
- **Configuration Options**: https://github.com/ajeetdsouza/zoxide/wiki/Configuration
- **Shell Integration**: https://github.com/ajeetdsouza/zoxide/wiki/Shell-integration
- **Comparison with z/autojump**: https://github.com/ajeetdsouza/zoxide#comparison

This module provides essential smart directory navigation capabilities for development workflows, significantly improving command-line productivity by reducing time spent navigating directory structures within the devgita ecosystem.
