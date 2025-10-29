# Aerospace Module Documentation

## Overview

The Aerospace module provides window manager installation and configuration management with devgita integration. It follows the standardized devgita app interface while providing Aerospace-specific operations for window management, workspace navigation, and macOS-specific desktop organization.

## App Purpose

Aerospace is a modern tiling window manager for macOS that provides automatic window arrangement, workspace management, and keyboard-driven navigation. This module ensures Aerospace is properly installed and configured with devgita's optimized settings for development workflows.

## Lifecycle Summary

1. **Installation**: Install Aerospace desktop application via Homebrew cask
2. **Configuration**: Apply devgita's Aerospace configuration templates with developer-focused layouts
3. **Execution**: Provide high-level Aerospace operations for window and workspace management

## Exported Functions

| Function           | Purpose                   | Behavior                                                       |
| ------------------ | ------------------------- | -------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Aerospace instance with platform-specific commands |
| `Install()`        | Standard installation     | Uses `InstallDesktopApp()` to install Aerospace                |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (ignored), then `Install()`          |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallDesktopApp()` to check before installing     |
| `ForceConfigure()` | Force configuration       | Overwrites existing configs with devgita defaults              |
| `SoftConfigure()`  | Conditional configuration | Preserves existing aerospace.toml if present                   |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                              |
| `ExecuteCommand()` | NO COMMANDS               | Place holder to keep compatibility with other apps             |
| `Update()`         | Update installation       | **Not implemented** - returns error                            |

## Installation Methods

### Install()

```go
aerospace := aerospace.New()
err := aerospace.Install()
```

- **Purpose**: Standard Aerospace installation
- **Behavior**: Uses `InstallDesktopApp()` to install Aerospace cask
- **Use case**: Initial Aerospace installation or explicit reinstall

### ForceInstall()

```go
aerospace := aerospace.New()
err := aerospace.ForceInstall()
```

- **Purpose**: Force Aerospace installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (ignored for Aerospace), then `Install()`
- **Use case**: Ensure fresh Aerospace installation or fix corrupted installation

### SoftInstall()

```go
aerospace := aerospace.New()
err := aerospace.SoftInstall()
```

- **Purpose**: Install Aerospace only if not already present
- **Behavior**: Uses `MaybeInstallDesktopApp()` to check before installing
- **Use case**: Standard installation that respects existing Aerospace installations

### Uninstall()

```go
err := aerospace.Uninstall()
```

- **Purpose**: Remove Aerospace installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Desktop app uninstallation requires careful handling and may be managed separately

### Update()

```go
err := aerospace.Update()
```

- **Purpose**: Update Aerospace installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Aerospace updates are typically handled by Homebrew cask upgrade

## Configuration Methods

### Configuration Paths

- **Source**: `paths.AerospaceConfigAppDir` (devgita's aerospace configs)
- **Destination**: `paths.AerospaceConfigLocalDir` (user's config directory)
- **Marker file**: `aerospace.toml` in `AerospaceConfigLocalDir`

### ForceConfigure()

```go
err := aerospace.ForceConfigure()
```

- **Purpose**: Apply Aerospace configuration regardless of existing files
- **Behavior**: Copies all configs from app dir to local dir, overwriting existing
- **Use case**: Reset to devgita defaults, apply config updates

### SoftConfigure()

```go
err := aerospace.SoftConfigure()
```

- **Purpose**: Apply Aerospace configuration only if not already configured
- **Behavior**: Checks for `aerospace.toml` marker file; if exists, does nothing
- **Marker logic**: `filepath.Join(paths.AerospaceConfigLocalDir, "aerospace.toml")`
- **Use case**: Initial setup that preserves user customizations

## Execution Methods

### ExecuteCommand()

- **Purpose**: Placeholder

### Aerospace-Specific Operations

The Aerospace CLI provides extensive window and workspace management capabilities:

#### Workspace Management

```bash
# List all workspaces
aerospace list-workspaces

# Switch to workspace
aerospace workspace 1
aerospace workspace next
aerospace workspace prev

# Move windows between workspaces
aerospace move-node-to-workspace 2
aerospace move-node-to-workspace next
aerospace move-node-to-workspace prev
```

#### Window Operations

```bash
# Focus windows
aerospace focus left
aerospace focus right
aerospace focus up
aerospace focus down

# Move windows
aerospace move left
aerospace move right
aerospace move up
aerospace move down

# Resize windows
aerospace resize smart +50
aerospace resize smart -50
aerospace resize width +100
aerospace resize height -100
```

#### Layout Management

```bash
# Change layout mode
aerospace layout tiling
aerospace layout floating
aerospace layout accordion
aerospace layout tabs

# Toggle layouts
aerospace layout toggle tiling floating
aerospace layout toggle accordion tabs

# Split operations
aerospace split horizontal
aerospace split vertical
aerospace split opposite
```

#### Monitor Operations

```bash
# List monitors
aerospace list-monitors

# Focus monitors
aerospace focus-monitor left
aerospace focus-monitor right
aerospace focus-monitor up
aerospace focus-monitor down

# Move workspaces between monitors
aerospace move-workspace-to-monitor next
aerospace move-workspace-to-monitor prev
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()`
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()`
3. **Update Configuration**: `New()` → `SoftInstall()` → `ForceConfigure()`
4. **Aerospace Operations**: `New()` → `ExecuteCommand()` with specific aerospace commands

## Constants and Paths

### Relevant Constants

- `constants.Aerospace`: Package name for Aerospace installation
- Used by all installation methods for consistent desktop app reference

### Configuration Paths

- `paths.AerospaceConfigAppDir`: Source directory for devgita's Aerospace configuration templates
- `paths.AerospaceConfigLocalDir`: Target directory for user's Aerospace configuration
- Configuration copying preserves TOML structure and file permissions

## Implementation Notes

- **Desktop App Installation**: Uses `InstallDesktopApp()` instead of `InstallPackage()` for macOS applications
- **ForceInstall Logic**: Calls `Uninstall()` first but ignores the error since Aerospace uninstall is not supported
- **Configuration Strategy**: Uses marker file (`aerospace.toml`) to determine if configuration exists
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Specificity**: Aerospace is macOS-only; installation will fail on other platforms
- **Legacy Compatibility**: Maintains deprecated function aliases to prevent breaking existing code
- **Update Method**: Not implemented as Aerospace updates should be handled by Homebrew cask system

## Configuration Structure

The Aerospace configuration (`aerospace.toml`) typically includes:

### Workspace Configuration

```toml
# Workspace definitions
[workspace-to-monitor-force-assignment]
1 = 'main'
2 = 'secondary'

# Default workspace per monitor
[workspace.default]
monitor = 'main'
```

### Layout Settings

```toml
# Default layout mode
default-layout = 'tiling'

# Layout-specific settings
[layout.tiling]
split-ratio = 0.5
orientation = 'auto'

[layout.floating]
border-width = 2
border-color = '#007ACC'
```

### Key Bindings

```toml
# Window focus
[key-mapping.focus]
left = 'alt-h'
right = 'alt-l'
up = 'alt-k'
down = 'alt-j'

# Workspace switching
[key-mapping.workspace]
'1' = 'alt-1'
'2' = 'alt-2'
'next' = 'alt-tab'
'prev' = 'alt-shift-tab'
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure Homebrew is installed and updated
2. **Configuration Not Applied**: Check file permissions in config directory
3. **Commands Don't Work**: Verify Aerospace is running and accessible in PATH
4. **Layout Issues**: Restart Aerospace after configuration changes

### macOS Permissions

Aerospace requires several macOS permissions:

- **Accessibility**: For window management
- **Screen Recording**: For certain operations
- **Automation**: For controlling other applications

Users may need to manually grant these permissions in System Preferences → Security & Privacy.

This module provides a robust foundation for Aerospace window manager integration within the devgita development environment setup.

