# Syntax Highlighting Module Documentation

## Overview

The Syntax Highlighting module provides installation and configuration management for zsh-syntax-highlighting with devgita integration. It follows the standardized devgita app interface while providing shell enhancement operations for command syntax highlighting and improved terminal readability.

## App Purpose

Zsh-syntax-highlighting is a Fish shell-like syntax highlighting plugin for Zsh that provides real-time syntax highlighting for commands as you type. This module ensures zsh-syntax-highlighting is properly installed and configured with devgita's shell environment setup for enhanced command-line readability and error detection.

## Lifecycle Summary

1. **Installation**: Install zsh-syntax-highlighting package via platform package managers (Homebrew/apt)
2. **Configuration**: Apply devgita's shell configuration by adding syntax highlighting source line to devgita.zsh
3. **Execution**: Provides placeholder operations for consistency with standardized app interface

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new SyntaxHighlighting instance with platform-specific commands |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install zsh-syntax-highlighting          |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing             |
| `ForceConfigure()` | Force configuration       | Adds syntax highlighting source line to devgita.zsh                 |
| `SoftConfigure()`  | Conditional configuration | Preserves existing configuration if already present                  |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                   |
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

- **Purpose**: Remove zsh-syntax-highlighting installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Shell enhancement tools are typically managed at the system level

### Update()

```go
err := syntaxHighlighting.Update()
```

- **Purpose**: Update zsh-syntax-highlighting installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Updates are typically handled by the system package manager

## Configuration Methods

### Configuration Paths

- **Target file**: `paths.AppDir/devgita.zsh` (devgita's shell initialization file)
- **Configuration line**: `source $(brew --prefix)/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh`
- **Detection string**: `zsh-syntax-highlighting.zsh` in the devgita.zsh file

### ForceConfigure()

```go
err := syntaxHighlighting.ForceConfigure()
```

- **Purpose**: Apply syntax highlighting configuration regardless of existing state
- **Behavior**: Adds the syntax highlighting source line to devgita.zsh
- **Use case**: Ensure syntax highlighting is properly sourced in shell configuration

### SoftConfigure()

```go
err := syntaxHighlighting.SoftConfigure()
```

- **Purpose**: Apply syntax highlighting configuration only if not already configured
- **Behavior**: Checks for existing configuration; if found, does nothing
- **Detection logic**: Searches for `zsh-syntax-highlighting.zsh` in `devgita.zsh`
- **Use case**: Initial setup that preserves existing shell customizations

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
4. **Shell Integration**: Automatically loaded when shell starts via devgita.zsh

## Constants and Paths

### Relevant Constants

- Package name: `"zsh-syntax-highlighting"` used for installation
- Used by all installation methods for consistent package reference

### Configuration Paths

- `paths.AppDir`: Directory containing devgita's shell configuration files
- Target file: `filepath.Join(paths.AppDir, "devgita.zsh")`
- Configuration integrates with devgita's overall shell setup

## Implementation Notes

- **Shell Plugin Nature**: Unlike typical applications, syntax highlighting is a shell enhancement that doesn't run independently
- **ForceInstall Logic**: Calls `Uninstall()` first but returns the error since syntax highlighting uninstall is not supported
- **Configuration Strategy**: Uses content detection to determine if configuration exists
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Shell Integration**: Configuration is applied to devgita's shell initialization rather than standalone config files
- **Update Method**: Not implemented as syntax highlighting updates should be handled by system package managers

## Configuration Structure

The syntax highlighting configuration adds a single source line to devgita.zsh:

```bash
source $(brew --prefix)/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh
```

This line:

- Sources the syntax highlighting plugin from the Homebrew installation path
- Integrates with Zsh to provide real-time command syntax highlighting
- Automatically loads when the shell starts
- Works alongside other devgita shell enhancements like autosuggestions

## Deprecated Functions

The module maintains backward compatibility through deprecated functions:

- `MaybeInstall()` → Use `SoftInstall()` instead
- `Setup()` → Use `ForceConfigure()` instead
- `MaybeSetup()` → Use `SoftConfigure()` instead

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Highlighting Not Working**: Verify shell configuration is properly sourced
3. **Configuration Not Applied**: Check file permissions and paths
4. **Duplicate Configuration**: Use `SoftConfigure()` to avoid duplicates
5. **Performance Issues**: Syntax highlighting may slow down very long commands

### Shell Integration

- Syntax highlighting requires Zsh shell to function
- Configuration must be sourced in shell initialization
- Works best with devgita's complete shell setup including autosuggestions
- May conflict with other syntax highlighting plugins
- Should be loaded after other shell enhancements for optimal performance

### Syntax Highlighting Features

- **Command Recognition**: Valid commands are highlighted in green, invalid in red
- **String Highlighting**: Quoted strings are highlighted with distinct colors
- **Path Recognition**: Existing file paths are highlighted differently
- **Bracket Matching**: Matching brackets and parentheses are highlighted
- **Globbing**: Shell glob patterns receive special highlighting
- **Real-time**: Highlighting updates as you type without affecting performance

This module provides essential command-line readability enhancement within the devgita development environment setup, significantly improving error detection and command composition.