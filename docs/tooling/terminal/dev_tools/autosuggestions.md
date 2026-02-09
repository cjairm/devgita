# Autosuggestions Module Documentation

## Overview

The Autosuggestions module provides installation and configuration management for zsh-autosuggestions with devgita integration. It follows the standardized devgita app interface while providing shell enhancement operations for command suggestions and productivity improvements.

**Note**: This module uses devgita's template-based configuration management system. Configuration changes are tracked in GlobalConfig and shell configuration is regenerated from templates rather than direct file manipulation.

## App Purpose

Zsh-autosuggestions is a Fish shell-like autosuggestion plugin for Zsh that suggests commands as you type based on history and completions. This module ensures zsh-autosuggestions is properly installed and configured with devgita's shell environment setup for enhanced command-line productivity.

## Lifecycle Summary

1. **Installation**: Install zsh-autosuggestions package via platform package managers (Homebrew/apt)
2. **Configuration**: Enable autosuggestions feature in GlobalConfig and regenerate devgita.zsh from template
3. **Execution**: Provides placeholder operations for consistency with standardized app interface

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Autosuggestions instance with platform-specific commands |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install zsh-autosuggestions               |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | Enables feature in GlobalConfig and regenerates shell configuration  |
| `SoftConfigure()`  | Conditional configuration | Checks GlobalConfig; enables only if not already enabled             |
| `Uninstall()`      | Remove installation       | **Fully supported** - Disables feature and regenerates shell config  |
| `ExecuteCommand()` | Execute commands          | **No operation** - returns nil                                       |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

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
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
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

- **Purpose**: Remove zsh-autosuggestions from shell configuration
- **Behavior**: Loads GlobalConfig → Disables feature → Regenerates shell config → Saves
- **Use case**: Remove autosuggestions from devgita shell without uninstalling the package
- **Note**: This disables the feature in the shell configuration but does not remove the package itself

### Update()

```go
err := autosuggestions.Update()
```

- **Purpose**: Update zsh-autosuggestions installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Updates are typically handled by the system package manager

## Configuration Methods

### Configuration Strategy

This module uses **template-based configuration management** via GlobalConfig:

- **Template**: `configs/templates/devgita.zsh.tmpl` contains conditional sections
- **GlobalConfig**: `~/.config/devgita/global_config.yaml` tracks enabled features
- **Generated file**: `~/.config/devgita/devgita.zsh` (regenerated from template)
- **Feature tracking**: `shell.zsh_autosuggestions` boolean field in GlobalConfig

### ForceConfigure()

```go
err := autosuggestions.ForceConfigure()
```

- **Purpose**: Enable autosuggestions in shell configuration
- **Behavior**: 
  1. Loads GlobalConfig from disk
  2. Enables `shell.zsh_autosuggestions` feature
  3. Regenerates `devgita.zsh` from template
  4. Saves GlobalConfig back to disk
- **Use case**: Enable autosuggestions or re-apply configuration

### SoftConfigure()

```go
err := autosuggestions.SoftConfigure()
```

- **Purpose**: Enable autosuggestions only if not already enabled
- **Behavior**: 
  1. Loads GlobalConfig from disk
  2. Checks if `shell.zsh_autosuggestions` is already enabled
  3. If enabled, returns nil (no operation)
  4. If not enabled, calls `ForceConfigure()`
- **Use case**: Initial setup that preserves existing configuration state

## Execution Methods

### ExecuteCommand()

```go
err := autosuggestions.ExecuteCommand("--version")
```

- **Purpose**: Execute autosuggestions-related commands
- **Behavior**: **No operation** - returns nil (success)
- **Rationale**: Autosuggestions is a shell plugin without standalone commands, but returns success for interface compliance

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()`
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()`
3. **Update Configuration**: `New()` → `SoftInstall()` → `ForceConfigure()`
4. **Remove from Shell**: `New()` → `Uninstall()`
5. **Shell Integration**: Automatically loaded when shell starts via devgita.zsh

## Constants and Paths

### Relevant Constants

- `constants.ZshAutosuggestions`: Package name used for installation and feature tracking
- Used consistently across all methods for package management and GlobalConfig operations

### Configuration Paths

- `paths.TemplatesAppDir`: Source directory for shell configuration templates
- `paths.AppDir`: Target directory for generated shell configuration
- Template file: `filepath.Join(paths.TemplatesAppDir, "devgita.zsh.tmpl")`
- Generated file: `filepath.Join(paths.AppDir, "devgita.zsh")`
- GlobalConfig file: `~/.config/devgita/global_config.yaml`

## Implementation Notes

- **Shell Plugin Nature**: Unlike typical applications, autosuggestions is a shell enhancement that doesn't run independently
- **Template-Based Configuration**: Uses GlobalConfig and template regeneration instead of direct file manipulation
- **Load-Modify-Regenerate-Save Pattern**: Each configuration method follows this transaction pattern
- **Fresh GlobalConfig Instances**: Each method creates a new `&config.GlobalConfig{}` and loads from disk to prevent stale data
- **Stateless Configuration**: GlobalConfig represents disk state, not app instance state
- **ForceInstall Logic**: Calls `Uninstall()` first, which now properly disables the feature
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as autosuggestions updates should be handled by system package managers

## Template Integration

### Template Structure

The autosuggestions configuration is defined in `configs/templates/devgita.zsh.tmpl`:

```bash
{{if .ZshAutosuggestions}}
# Zsh Autosuggestions - Fish-like autosuggestions
if [[ -f $(brew --prefix)/share/zsh-autosuggestions/zsh-autosuggestions.zsh ]]; then
    source $(brew --prefix)/share/zsh-autosuggestions/zsh-autosuggestions.zsh
fi
{{end}}
```

### GlobalConfig Tracking

The feature state is tracked in `~/.config/devgita/global_config.yaml`:

```yaml
shell:
  zsh_autosuggestions: true  # Enabled
  # ... other shell features
```

### Generated Configuration

When enabled, the generated `devgita.zsh` contains:

```bash
# Zsh Autosuggestions - Fish-like autosuggestions
if [[ -f $(brew --prefix)/share/zsh-autosuggestions/zsh-autosuggestions.zsh ]]; then
    source $(brew --prefix)/share/zsh-autosuggestions/zsh-autosuggestions.zsh
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
2. **Suggestions Don't Work**: Verify shell configuration is properly sourced and GlobalConfig has feature enabled
3. **Configuration Not Applied**: Check GlobalConfig file exists and feature is enabled
4. **GlobalConfig Load Errors**: Ensure `~/.config/devgita/global_config.yaml` is valid YAML
5. **Template Not Found**: Verify `configs/templates/devgita.zsh.tmpl` exists in devgita repository

### Shell Integration

- Autosuggestions requires Zsh shell to function
- Configuration must be sourced in shell initialization (add `source ~/.config/devgita/devgita.zsh` to `.zshrc`)
- Works best with devgita's complete shell setup including syntax highlighting
- May conflict with other autosuggestion plugins

### Template System Benefits

- **Single Source of Truth**: Template file is the only place shell configuration is defined
- **Trackable**: Git can track template changes easily
- **Predictable**: Regeneration always produces same output for same inputs
- **No Conflicts**: No string manipulation or file appending/removing
- **Clean Uninstall**: Disabling a feature regenerates without it

This module provides essential command-line productivity enhancement within the devgita development environment setup.

