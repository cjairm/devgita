# Powerlevel10k Module Documentation

## Overview

The Powerlevel10k module provides installation and configuration management for Powerlevel10k Zsh theme with devgita integration. It follows the standardized devgita app interface while providing shell enhancement operations for advanced prompt customization, git integration, and development environment visualization.

## App Purpose

Powerlevel10k is a fast, customizable, and highly optimized Zsh theme that provides a visually rich command-line prompt with git status, system information, and development context indicators. This module ensures Powerlevel10k is properly installed and configured with devgita's shell environment setup for enhanced developer productivity and visual feedback.

## Lifecycle Summary

1. **Installation**: Install powerlevel10k package via platform package managers (Homebrew/apt)
2. **Configuration**: Apply devgita's shell configuration by adding Powerlevel10k source line to devgita.zsh
3. **Execution**: Provide high-level Powerlevel10k operations for prompt configuration and theme management

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new PowerLevel10k instance with platform-specific commands   |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install powerlevel10k                     |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | Adds powerlevel10k source line to devgita.zsh                        |
| `SoftConfigure()`  | Conditional configuration | Preserves existing configuration if already present                  |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute p10k commands     | Runs p10k with provided arguments                                    |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |
| `Reconfigure()`    | Reconfigure theme         | Runs p10k configure command for interactive setup                    |

## Installation Methods

### Install()

```go
p10k := powerlevel10k.New()
err := p10k.Install()
```

- **Purpose**: Standard Powerlevel10k installation
- **Behavior**: Uses `InstallPackage()` to install powerlevel10k package
- **Use case**: Initial Powerlevel10k installation or explicit reinstall

### ForceInstall()

```go
p10k := powerlevel10k.New()
err := p10k.ForceInstall()
```

- **Purpose**: Force Powerlevel10k installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh Powerlevel10k installation or fix corrupted installation

### SoftInstall()

```go
p10k := powerlevel10k.New()
err := p10k.SoftInstall()
```

- **Purpose**: Install Powerlevel10k only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing Powerlevel10k installations

### Uninstall()

```go
err := p10k.Uninstall()
```

- **Purpose**: Remove Powerlevel10k installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Shell theme tools are typically managed at the system level

### Update()

```go
err := p10k.Update()
```

- **Purpose**: Update Powerlevel10k installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Updates are typically handled by the system package manager

## Configuration Methods

### Configuration Paths

- **Target file**: `paths.AppDir/devgita.zsh` (devgita's shell initialization file)
- **Configuration line**: `source $(brew --prefix)/share/powerlevel10k/powerlevel10k.zsh-theme`
- **Detection string**: `powerlevel10k.zsh-theme` in the devgita.zsh file

### ForceConfigure()

```go
err := p10k.ForceConfigure()
```

- **Purpose**: Apply Powerlevel10k configuration regardless of existing state
- **Behavior**: Adds the Powerlevel10k source line to devgita.zsh
- **Use case**: Ensure Powerlevel10k is properly sourced in shell configuration

### SoftConfigure()

```go
err := p10k.SoftConfigure()
```

- **Purpose**: Apply Powerlevel10k configuration only if not already configured
- **Behavior**: Checks for existing configuration; if found, does nothing
- **Detection logic**: Searches for `powerlevel10k.zsh-theme` in `devgita.zsh`
- **Use case**: Initial setup that preserves existing shell customizations

## Execution Methods

### ExecuteCommand()

```go
err := p10k.ExecuteCommand("configure")
err := p10k.ExecuteCommand("reload")
err := p10k.ExecuteCommand("segment", "show", "time")
```

- **Purpose**: Execute p10k commands with provided arguments
- **Parameters**: Variable arguments passed directly to p10k binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### Powerlevel10k-Specific Operations

#### Reconfigure()

```go
err := p10k.Reconfigure()
```

- **Purpose**: Launch interactive Powerlevel10k configuration wizard
- **Behavior**: Executes `p10k configure` command
- **Use case**: Initial theme setup or reconfiguring visual preferences

#### Theme Management Commands

The p10k CLI provides extensive theme management capabilities:

##### Configuration Management

```bash
# Interactive configuration wizard
p10k configure

# Reload configuration
p10k reload

# Show configuration status
p10k status

# Finalize instant prompt
p10k finalize
```

##### Segment Control

```bash
# Show specific segments
p10k segment show dir
p10k segment show git
p10k segment show time

# Hide specific segments
p10k segment hide dir
p10k segment hide git
p10k segment hide time

# List all segments
p10k segment list
```

##### Display Options

```bash
# Show instant prompt log
p10k instant-prompt verbose

# Debug mode
p10k debug

# Profile performance
p10k benchmark
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()`
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()`
3. **Update Configuration**: `New()` → `SoftInstall()` → `ForceConfigure()`
4. **Interactive Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` → `Reconfigure()`
5. **P10k Operations**: `New()` → `ExecuteCommand()` with specific p10k arguments

## Constants and Paths

### Relevant Constants

- `constants.Powerlevel10k`: Package name ("powerlevel10k") used for installation
- Used by all installation methods for consistent package reference

### Configuration Paths

- `paths.AppDir`: Directory containing devgita's shell configuration files
- Target file: `filepath.Join(paths.AppDir, "devgita.zsh")`
- Configuration integrates with devgita's overall shell setup

## Implementation Notes

- **Shell Theme Nature**: Unlike typical applications, Powerlevel10k is a shell enhancement that modifies prompt appearance and behavior
- **ForceInstall Logic**: Calls `Uninstall()` first but returns the error since Powerlevel10k uninstall is not supported
- **Configuration Strategy**: Uses content detection to determine if configuration exists
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Shell Integration**: Configuration is applied to devgita's shell initialization rather than standalone config files
- **Update Method**: Not implemented as Powerlevel10k updates should be handled by system package managers

## Configuration Structure

The Powerlevel10k configuration adds a single source line to devgita.zsh:

```bash
source $(brew --prefix)/share/powerlevel10k/powerlevel10k.zsh-theme
```

This line:

- Sources the Powerlevel10k theme from the Homebrew installation path
- Integrates with Zsh to provide enhanced prompt functionality
- Automatically loads when the shell starts
- Works alongside other devgita shell enhancements

## Powerlevel10k Features

### Visual Elements

- **Git Integration**: Branch name, status indicators, ahead/behind counts
- **System Information**: User, hostname, current directory, exit codes
- **Development Context**: Python virtual environments, Node.js versions, etc.
- **Performance**: Fast rendering with instant prompt capabilities
- **Customization**: Extensive configuration options via interactive wizard

### Prompt Styles

- **Rainbow**: Colorful, information-rich display
- **Lean**: Minimal, clean appearance
- **Classic**: Traditional prompt styling
- **Pure**: Minimalist design inspired by Pure theme

### Advanced Features

- **Instant Prompt**: Fast shell startup with cached prompt
- **Transient Prompt**: Simplified previous command display
- **Right Prompt**: Additional information on the right side
- **Multiline Support**: Complex layouts across multiple lines

## Deprecated Functions

The module maintains backward compatibility through deprecated functions:

- `MaybeInstall()` → Use `SoftInstall()` instead
- `Setup()` → Use `ForceConfigure()` instead
- `MaybeSetup()` → Use `SoftConfigure()` instead
- `Run()` → Use `ExecuteCommand()` instead

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Theme Not Loading**: Verify shell configuration is properly sourced
3. **Configuration Not Applied**: Check file permissions and paths
4. **Duplicate Configuration**: Use `SoftConfigure()` to avoid duplicates
5. **Slow Performance**: Run `p10k configure` to optimize settings

### Shell Integration

- Powerlevel10k requires Zsh shell to function
- Configuration must be sourced in shell initialization
- Works best with devgita's complete shell setup
- May conflict with other prompt themes

### Performance Optimization

- Use instant prompt for faster shell startup
- Disable unused segments to improve performance
- Configure appropriate prompt refresh intervals
- Utilize Powerlevel10k's built-in caching mechanisms

This module provides essential prompt enhancement and visual feedback within the devgita development environment setup, significantly improving command-line productivity and development context awareness.