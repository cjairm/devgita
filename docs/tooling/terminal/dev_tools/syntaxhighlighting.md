# Syntax Highlighting Module Documentation

## Overview

The Syntax Highlighting module provides installation and configuration management for zsh-syntax-highlighting with devgita integration. It follows the standardized devgita app interface while providing shell enhancement operations for command syntax highlighting and improved terminal readability.

**Note**: This module uses devgita's template-based configuration management system. Configuration changes are tracked in GlobalConfig and shell configuration is regenerated from templates rather than direct file manipulation.

## App Purpose

Zsh-syntax-highlighting is a Fish shell-like syntax highlighting plugin for Zsh that provides real-time syntax highlighting for commands as you type. This module ensures zsh-syntax-highlighting is properly installed and configured with devgita's shell environment setup for enhanced command-line readability and error detection.

## Lifecycle Summary

1. **Installation**: Install zsh-syntax-highlighting package via platform package managers (Homebrew/apt)
2. **Configuration**: Enable syntax highlighting feature in GlobalConfig and regenerate devgita.zsh from template
3. **Execution**: Provides placeholder operations for consistency with standardized app interface

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Syntaxhighlighting instance with platform-specific commands |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install zsh-syntax-highlighting          |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing             |
| `ForceConfigure()` | Force configuration       | Enables feature in GlobalConfig and regenerates shell configuration  |
| `SoftConfigure()`  | Conditional configuration | Checks GlobalConfig; enables only if not already enabled             |
| `Uninstall()`      | Remove installation       | **Fully supported** - Disables feature and regenerates shell config  |
| `ExecuteCommand()` | Execute commands          | **No operation** - returns nil                                      |
| `Update()`         | Update installation       | **Not implemented** - returns error                                 |

## Installation Methods

### Install()

```go
syntaxHighlighting := syntaxhighlighting.New()
err := syntaxHighlighting.Install()
```

- **Purpose**: Standard zsh-syntax-highlighting installation
- **Behavior**: Uses `InstallPackage()` to install zsh-syntax-highlighting package
- **Use case**: Initial zsh-syntax-highlighting installation or explicit reinstall

### ForceInstall()

```go
syntaxHighlighting := syntaxhighlighting.New()
err := syntaxHighlighting.ForceInstall()
```

- **Purpose**: Force zsh-syntax-highlighting installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh syntax highlighting installation or fix corrupted installation

### SoftInstall()

```go
syntaxHighlighting := syntaxhighlighting.New()
err := syntaxHighlighting.SoftInstall()
```

- **Purpose**: Install zsh-syntax-highlighting only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing syntax highlighting installations

### Uninstall()

```go
err := syntaxHighlighting.Uninstall()
```

- **Purpose**: Remove zsh-syntax-highlighting from shell configuration
- **Behavior**: Loads GlobalConfig → Disables feature → Regenerates shell config → Saves
- **Use case**: Remove syntax highlighting from devgita shell without uninstalling the package
- **Note**: This disables the feature in the shell configuration but does not remove the package itself

### Update()

```go
err := syntaxHighlighting.Update()
```

- **Purpose**: Update zsh-syntax-highlighting installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Updates are typically handled by the system package manager

## Configuration Methods

### Configuration Strategy

This module uses **template-based configuration management** via GlobalConfig:

- **Template**: `configs/templates/devgita.zsh.tmpl` contains conditional sections
- **GlobalConfig**: `~/.config/devgita/global_config.yaml` tracks enabled features
- **Generated file**: `~/.config/devgita/devgita.zsh` (regenerated from template)
- **Feature tracking**: `shell.zsh_syntax_highlighting` boolean field in GlobalConfig

### ForceConfigure()

```go
err := syntaxHighlighting.ForceConfigure()
```

- **Purpose**: Enable syntax highlighting in shell configuration
- **Behavior**: 
  1. Loads GlobalConfig from disk
  2. Enables `shell.zsh_syntax_highlighting` feature
  3. Regenerates `devgita.zsh` from template
  4. Saves GlobalConfig back to disk
- **Use case**: Enable syntax highlighting or re-apply configuration

### SoftConfigure()

```go
err := syntaxHighlighting.SoftConfigure()
```

- **Purpose**: Enable syntax highlighting only if not already enabled
- **Behavior**: 
  1. Loads GlobalConfig from disk
  2. Checks if `shell.zsh_syntax_highlighting` is already enabled
  3. If enabled, returns nil (no operation)
  4. If not enabled, calls `ForceConfigure()`
- **Use case**: Initial setup that preserves existing configuration state

## Execution Methods

### ExecuteCommand()

```go
err := syntaxHighlighting.ExecuteCommand("--version")
```

- **Purpose**: Execute syntax highlighting-related commands
- **Behavior**: **No operation** - returns nil (success)
- **Rationale**: Syntax highlighting is a shell plugin without standalone commands, but returns success for interface compliance

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()`
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()`
3. **Update Configuration**: `New()` → `SoftInstall()` → `ForceConfigure()`
4. **Remove from Shell**: `New()` → `Uninstall()`
5. **Shell Integration**: Automatically loaded when shell starts via devgita.zsh

## Constants and Paths

### Relevant Constants

- `constants.Syntaxhighlighting`: Package name ("zsh-syntax-highlighting") used for installation and feature tracking
- Used consistently across all methods for package management and GlobalConfig operations

### Configuration Paths

- `paths.TemplatesAppDir`: Source directory for shell configuration templates
- `paths.AppDir`: Target directory for generated shell configuration
- Template file: `filepath.Join(paths.TemplatesAppDir, "devgita.zsh.tmpl")`
- Generated file: `filepath.Join(paths.AppDir, "devgita.zsh")`
- GlobalConfig file: `~/.config/devgita/global_config.yaml`

## Implementation Notes

- **Shell Plugin Nature**: Unlike typical applications, syntax highlighting is a shell enhancement that doesn't run independently
- **Template-Based Configuration**: Uses GlobalConfig and template regeneration instead of direct file manipulation
- **Load-Modify-Regenerate-Save Pattern**: Each configuration method follows this transaction pattern
- **Fresh GlobalConfig Instances**: Each method creates a new `&config.GlobalConfig{}` and loads from disk to prevent stale data
- **Stateless Configuration**: GlobalConfig represents disk state, not app instance state
- **ForceInstall Logic**: Calls `Uninstall()` first, which now properly disables the feature
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as syntax highlighting updates should be handled by system package managers

## Template Integration

### Template Structure

The syntax highlighting configuration is defined in `configs/templates/devgita.zsh.tmpl`:

```bash
{{if .ZshSyntaxHighlighting}}
# Zsh Syntax Highlighting - Fish-like syntax highlighting
if [[ -f $(brew --prefix)/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh ]]; then
    source $(brew --prefix)/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh
fi
{{end}}
```

### GlobalConfig Tracking

The feature state is tracked in `~/.config/devgita/global_config.yaml`:

```yaml
shell:
  zsh_syntax_highlighting: true  # Enabled
  # ... other shell features
```

### Generated Configuration

When enabled, the generated `devgita.zsh` contains:

```bash
# Zsh Syntax Highlighting - Fish-like syntax highlighting
if [[ -f $(brew --prefix)/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh ]]; then
    source $(brew --prefix)/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh
fi
```

This approach:
- Provides single source of truth (template file)
- Enables clean enable/disable operations
- Prevents configuration conflicts
- Makes tracking and version control easier
- Ensures consistent regeneration

## Deprecated Functions

The module maintains backward compatibility through deprecated functions:

- `MaybeInstall()` → Use `SoftInstall()` instead
- `Setup()` → Use `ForceConfigure()` instead
- `MaybeSetup()` → Use `SoftConfigure()` instead

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Highlighting Not Working**: Verify shell configuration is properly sourced and GlobalConfig has feature enabled
3. **Configuration Not Applied**: Check GlobalConfig file exists and feature is enabled
4. **GlobalConfig Load Errors**: Ensure `~/.config/devgita/global_config.yaml` is valid YAML
5. **Template Not Found**: Verify `configs/templates/devgita.zsh.tmpl` exists in devgita repository

### Shell Integration

- Syntax highlighting requires Zsh shell to function
- Configuration must be sourced in shell initialization (add `source ~/.config/devgita/devgita.zsh` to `.zshrc`)
- Works best with devgita's complete shell setup including autosuggestions
- May conflict with other syntax highlighting plugins
- Should be loaded after other shell enhancements for optimal performance

### Template System Benefits

- **Single Source of Truth**: Template file is the only place shell configuration is defined
- **Trackable**: Git can track template changes easily
- **Predictable**: Regeneration always produces same output for same inputs
- **No Conflicts**: No string manipulation or file appending/removing
- **Clean Uninstall**: Disabling a feature regenerates without it

### Syntax Highlighting Features

- **Command Recognition**: Valid commands are highlighted in green, invalid in red
- **String Highlighting**: Quoted strings are highlighted with distinct colors
- **Path Recognition**: Existing file paths are highlighted differently
- **Bracket Matching**: Matching brackets and parentheses are highlighted
- **Globbing**: Shell glob patterns receive special highlighting
- **Real-time**: Highlighting updates as you type without affecting performance

This module provides essential command-line readability enhancement within the devgita development environment setup, significantly improving error detection and command composition.
