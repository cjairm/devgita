# Powerlevel10k Module Documentation

## Overview

The Powerlevel10k module provides installation and configuration management for Powerlevel10k Zsh theme with devgita integration. It follows the standardized devgita app interface while providing shell enhancement operations for advanced prompt customization, git integration, and development environment visualization.

## App Purpose

Powerlevel10k is a fast, customizable, and highly optimized Zsh theme that provides a visually rich command-line prompt with git status, system information, and development context indicators. This module ensures Powerlevel10k is properly installed and configured with devgita's shell environment setup for enhanced developer productivity and visual feedback.

## Lifecycle Summary

1. **Installation**: Install powerlevel10k package via platform package managers (Homebrew/apt)
2. **Configuration**: Enable powerlevel10k feature in GlobalConfig and regenerate shell configuration via templates
3. **Execution**: Provide high-level Powerlevel10k operations for prompt configuration and theme management

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new PowerLevel10k instance with platform-specific commands   |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install powerlevel10k                     |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | Enables feature in GlobalConfig and regenerates shell configuration  |
| `SoftConfigure()`  | Conditional configuration | Checks GlobalConfig; enables only if not already enabled             |
| `Uninstall()`      | Remove installation       | **Fully supported** - Disables feature and regenerates shell config  |
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

- **Purpose**: Remove Powerlevel10k from shell configuration
- **Behavior**: Loads GlobalConfig → Disables feature → Regenerates shell config → Saves
- **Use case**: Remove Powerlevel10k from devgita shell without uninstalling the package
- **Note**: This disables the feature in the shell configuration but does not remove the package itself

### Update()

```go
err := p10k.Update()
```

- **Purpose**: Update Powerlevel10k installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Updates are typically handled by the system package manager

## Configuration Methods

### Configuration Strategy

This module uses **template-based configuration management** via GlobalConfig:

- **Template**: `configs/templates/devgita.zsh.tmpl` contains conditional sections
- **GlobalConfig**: `~/.config/devgita/global_config.yaml` tracks enabled features
- **Generated file**: `~/.config/devgita/devgita.zsh` (regenerated from template)
- **Feature tracking**: `shell.powerlevel10k` boolean field in GlobalConfig

### ForceConfigure()

```go
err := p10k.ForceConfigure()
```

- **Purpose**: Enable Powerlevel10k in shell configuration
- **Behavior**: 
  1. Loads GlobalConfig from disk
  2. Enables `shell.powerlevel10k` feature
  3. Regenerates `devgita.zsh` from template
  4. Saves GlobalConfig back to disk
- **Use case**: Enable Powerlevel10k or re-apply configuration

### SoftConfigure()

```go
err := p10k.SoftConfigure()
```

- **Purpose**: Enable Powerlevel10k only if not already enabled
- **Behavior**: 
  1. Loads GlobalConfig from disk
  2. Checks if `shell.powerlevel10k` is already enabled
  3. If enabled, returns nil (no operation)
  4. If not enabled, calls `ForceConfigure()`
- **Use case**: Initial setup that preserves existing configuration state

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
4. **Remove from Shell**: `New()` → `Uninstall()`
5. **Interactive Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` → `Reconfigure()`
6. **P10k Operations**: `New()` → `ExecuteCommand()` with specific p10k arguments
7. **Shell Integration**: Automatically loaded when shell starts via devgita.zsh

## Constants and Paths

### Relevant Constants

- `constants.Powerlevel10k`: Package name ("powerlevel10k") used for installation
- Used by all installation methods for consistent package reference

### Configuration Paths

- `paths.TemplatesAppDir`: Source directory for shell configuration templates
- `paths.AppDir`: Target directory for generated shell configuration
- Template file: `filepath.Join(paths.TemplatesAppDir, "devgita.zsh.tmpl")`
- Generated file: `filepath.Join(paths.AppDir, "devgita.zsh")`
- GlobalConfig file: `~/.config/devgita/global_config.yaml`

## Implementation Notes

- **Shell Theme Nature**: Unlike typical applications, Powerlevel10k is a shell enhancement that modifies prompt appearance and behavior
- **Template-Based Configuration**: Uses GlobalConfig and template regeneration instead of direct file manipulation
- **Load-Modify-Regenerate-Save Pattern**: Each configuration method follows this transaction pattern
- **Fresh GlobalConfig Instances**: Each method creates a new `&config.GlobalConfig{}` and loads from disk to prevent stale data
- **Stateless Configuration**: GlobalConfig represents disk state, not app instance state
- **ForceInstall Logic**: Calls `Uninstall()` first, which now properly disables the feature
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as Powerlevel10k updates should be handled by system package managers

## Template Integration

### Template Structure

The Powerlevel10k configuration is defined in `configs/templates/devgita.zsh.tmpl`:

```bash
{{if .Powerlevel10k}}
# Powerlevel10k - Fast, customizable Zsh theme
if [[ -f $(brew --prefix)/share/powerlevel10k/powerlevel10k.zsh-theme ]]; then
    source $(brew --prefix)/share/powerlevel10k/powerlevel10k.zsh-theme
fi
{{end}}
```

### GlobalConfig Tracking

The feature state is tracked in `~/.config/devgita/global_config.yaml`:

```yaml
shell:
  powerlevel10k: true  # Enabled
  # ... other shell features
```

### Generated Configuration

When enabled, the generated `devgita.zsh` contains:

```bash
# Powerlevel10k - Fast, customizable Zsh theme
if [[ -f $(brew --prefix)/share/powerlevel10k/powerlevel10k.zsh-theme ]]; then
    source $(brew --prefix)/share/powerlevel10k/powerlevel10k.zsh-theme
fi
```

This approach:
- Provides single source of truth (template file)
- Enables clean enable/disable operations
- Prevents configuration conflicts
- Makes tracking and version control easier
- Ensures consistent regeneration

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
2. **Theme Not Loading**: Verify shell configuration is properly sourced and GlobalConfig has feature enabled
3. **Configuration Not Applied**: Check GlobalConfig file exists and feature is enabled
4. **GlobalConfig Load Errors**: Ensure `~/.config/devgita/global_config.yaml` is valid YAML
5. **Template Not Found**: Verify `configs/templates/devgita.zsh.tmpl` exists in devgita repository
6. **Slow Performance**: Run `p10k configure` to optimize settings

### Shell Integration

- Powerlevel10k requires Zsh shell to function
- Configuration must be sourced in shell initialization (add `source ~/.config/devgita/devgita.zsh` to `.zshrc`)
- Works best with devgita's complete shell setup
- May conflict with other prompt themes
- Should be loaded after other shell enhancements for optimal performance

### Template System Benefits

- **Single Source of Truth**: Template file is the only place shell configuration is defined
- **Trackable**: Git can track template changes easily
- **Predictable**: Regeneration always produces same output for same inputs
- **No Conflicts**: No string manipulation or file appending/removing
- **Clean Uninstall**: Disabling a feature regenerates without it

### Performance Optimization

- Use instant prompt for faster shell startup
- Disable unused segments to improve performance
- Configure appropriate prompt refresh intervals
- Utilize Powerlevel10k's built-in caching mechanisms

This module provides essential prompt enhancement and visual feedback within the devgita development environment setup, significantly improving command-line productivity and development context awareness.