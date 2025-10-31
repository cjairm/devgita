# Autosuggestions Module Documentation

## Overview

The Autosuggestions module provides installation and configuration management for zsh-autosuggestions with devgita integration. It follows the standardized devgita app interface while providing shell enhancement operations for command suggestions and productivity improvements.

## App Purpose

Zsh-autosuggestions is a Fish shell-like autosuggestion plugin for Zsh that suggests commands as you type based on history and completions. This module ensures zsh-autosuggestions is properly installed and configured with devgita's shell environment setup for enhanced command-line productivity.

## Lifecycle Summary

1. **Installation**: Install zsh-autosuggestions package via platform package managers (Homebrew/apt)
2. **Configuration**: Apply devgita's shell configuration by adding autosuggestions source line to devgita.zsh
3. **Execution**: Provides placeholder operations for consistency with standardized app interface

## Exported Functions

| Function           | Purpose                   | Behavior                                                         |
| ------------------ | ------------------------- | ---------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Autosuggestions instance with platform-specific commands |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install zsh-autosuggestions          |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (ignored), then `Install()`           |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing         |
| `ForceConfigure()` | Force configuration       | Adds autosuggestions source line to devgita.zsh                 |
| `SoftConfigure()`  | Conditional configuration | Preserves existing configuration if already present              |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                               |
| `ExecuteCommand()` | Execute commands          | **Not implemented** - returns error                             |
| `Update()`         | Update installation       | **Not implemented** - returns error                             |

## Installation Methods

### Install()

```go
autosuggestions := autosuggestions.New()
err := autosuggestions.Install()
```

- **Purpose**: Standard zsh-autosuggestions installation
- **Behavior**: Uses `InstallPackage()` to install zsh-autosuggestions package
- **Use case**: Initial zsh-autosuggestions installation or explicit reinstall

### ForceInstall()

```go
autosuggestions := autosuggestions.New()
err := autosuggestions.ForceInstall()
```

- **Purpose**: Force zsh-autosuggestions installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (ignored for autosuggestions), then `Install()`
- **Use case**: Ensure fresh autosuggestions installation or fix corrupted installation

### SoftInstall()

```go
autosuggestions := autosuggestions.New()
err := autosuggestions.SoftInstall()
```

- **Purpose**: Install zsh-autosuggestions only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing autosuggestions installations

### Uninstall()

```go
err := autosuggestions.Uninstall()
```

- **Purpose**: Remove zsh-autosuggestions installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Shell enhancement tools are typically managed at the system level

### Update()

```go
err := autosuggestions.Update()
```

- **Purpose**: Update zsh-autosuggestions installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Updates are typically handled by the system package manager

## Configuration Methods

### Configuration Paths

- **Target file**: `paths.AppDir/devgita.zsh` (devgita's shell initialization file)
- **Configuration line**: `source $(brew --prefix)/share/zsh-autosuggestions/zsh-autosuggestions.zsh`
- **Detection string**: `zsh-autosuggestions.zsh` in the devgita.zsh file

### ForceConfigure()

```go
err := autosuggestions.ForceConfigure()
```

- **Purpose**: Apply autosuggestions configuration regardless of existing state
- **Behavior**: Adds the autosuggestions source line to devgita.zsh
- **Use case**: Ensure autosuggestions is properly sourced in shell configuration

### SoftConfigure()

```go
err := autosuggestions.SoftConfigure()
```

- **Purpose**: Apply autosuggestions configuration only if not already configured
- **Behavior**: Checks for existing configuration; if found, does nothing
- **Detection logic**: Searches for `zsh-autosuggestions.zsh` in `devgita.zsh`
- **Use case**: Initial setup that preserves existing shell customizations

## Execution Methods

### ExecuteCommand()

```go
err := autosuggestions.ExecuteCommand("--version")
```

- **Purpose**: Execute autosuggestions-related commands
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Autosuggestions is a shell plugin without standalone commands

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()`
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()`
3. **Update Configuration**: `New()` → `SoftInstall()` → `ForceConfigure()`
4. **Shell Integration**: Automatically loaded when shell starts via devgita.zsh

## Constants and Paths

### Relevant Constants

- Package name: `"zsh-autosuggestions"` used for installation
- Used by all installation methods for consistent package reference

### Configuration Paths

- `paths.AppDir`: Directory containing devgita's shell configuration files
- Target file: `filepath.Join(paths.AppDir, "devgita.zsh")`
- Configuration integrates with devgita's overall shell setup

## Implementation Notes

- **Shell Plugin Nature**: Unlike typical applications, autosuggestions is a shell enhancement that doesn't run independently
- **ForceInstall Logic**: Calls `Uninstall()` first but ignores the error since autosuggestions uninstall is not supported
- **Configuration Strategy**: Uses content detection to determine if configuration exists
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Shell Integration**: Configuration is applied to devgita's shell initialization rather than standalone config files
- **Update Method**: Not implemented as autosuggestions updates should be handled by system package managers

## Configuration Structure

The autosuggestions configuration adds a single source line to devgita.zsh:

```bash
source $(brew --prefix)/share/zsh-autosuggestions/zsh-autosuggestions.zsh
```

This line:
- Sources the autosuggestions plugin from the Homebrew installation path
- Integrates with Zsh to provide command suggestions
- Automatically loads when the shell starts
- Works alongside other devgita shell enhancements

## Deprecated Functions

The module maintains backward compatibility through deprecated functions:

- `MaybeInstall()` → Use `SoftInstall()` instead
- `Setup()` → Use `ForceConfigure()` instead  
- `MaybeSetup()` → Use `SoftConfigure()` instead

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Suggestions Don't Work**: Verify shell configuration is properly sourced
3. **Configuration Not Applied**: Check file permissions and paths
4. **Duplicate Configuration**: Use `SoftConfigure()` to avoid duplicates

### Shell Integration

- Autosuggestions requires Zsh shell to function
- Configuration must be sourced in shell initialization
- Works best with devgita's complete shell setup including syntax highlighting
- May conflict with other autosuggestion plugins

This module provides essential command-line productivity enhancement within the devgita development environment setup.