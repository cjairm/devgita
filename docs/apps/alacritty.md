# Alacritty Module Documentation

## Overview

The Alacritty module provides installation and configuration management for Alacritty terminal emulator with devgita integration. It follows the standardized devgita app interface while providing Alacritty-specific operations for terminal customization, configuration management, and cross-platform desktop application deployment.

## App Purpose

Alacritty is a fast, cross-platform terminal emulator written in Rust that uses GPU acceleration for rendering. This module ensures Alacritty is properly installed and configured with devgita's optimized settings for development workflows, including custom fonts, themes, and configuration templates.

## Lifecycle Summary

1. **Installation**: Install Alacritty desktop application via platform package managers (Homebrew cask/apt)
2. **Configuration**: Apply devgita's Alacritty configuration templates with developer-focused settings
3. **Execution**: Provide high-level Alacritty operations for terminal management and customization

## Exported Functions

| Function           | Purpose                   | Behavior                                                       |
| ------------------ | ------------------------- | -------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Alacritty instance with platform-specific commands |
| `Install()`        | Standard installation     | Uses `InstallDesktopApp()` to install Alacritty                |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (ignored), then `Install()`          |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallDesktopApp()` to check before installing     |
| `ForceConfigure()` | Force configuration       | Overwrites existing configs with devgita defaults              |
| `SoftConfigure()`  | Conditional configuration | Preserves existing alacritty.toml if present                   |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                              |
| `ExecuteCommand()` | Execute alacritty command | Runs alacritty with provided arguments                         |
| `Update()`         | Update installation       | **Not implemented** - returns error                            |

## Installation Methods

### Install()

```go
alacritty := alacritty.New()
err := alacritty.Install()
```

- **Purpose**: Standard Alacritty installation
- **Behavior**: Uses `InstallDesktopApp()` to install Alacritty desktop application
- **Use case**: Initial Alacritty installation or explicit reinstall

### ForceInstall()

```go
alacritty := alacritty.New()
err := alacritty.ForceInstall()
```

- **Purpose**: Force Alacritty installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (ignored for Alacritty), then `Install()`
- **Use case**: Ensure fresh Alacritty installation or fix corrupted installation

### SoftInstall()

```go
alacritty := alacritty.New()
err := alacritty.SoftInstall()
```

- **Purpose**: Install Alacritty only if not already present
- **Behavior**: Uses `MaybeInstallDesktopApp()` to check before installing
- **Use case**: Standard installation that respects existing Alacritty installations

### Uninstall()

```go
err := alacritty.Uninstall()
```

- **Purpose**: Remove Alacritty installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Desktop app uninstallation requires careful handling and may be managed separately

### Update()

```go
err := alacritty.Update()
```

- **Purpose**: Update Alacritty installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Alacritty updates are typically handled by the system package manager

## Configuration Methods

### Configuration Paths

- **Source**: `paths.AlacrittyConfigAppDir` (devgita's alacritty configs)
- **Destination**: `paths.AlacrittyConfigLocalDir` (user's config directory)
- **Fonts**: `paths.AlacrittyFontsAppDir/default` → `AlacrittyConfigLocalDir`
- **Themes**: `paths.AlacrittyThemesAppDir/default` → `AlacrittyConfigLocalDir`
- **Marker file**: `alacritty.toml` in `AlacrittyConfigLocalDir`

### ForceConfigure()

```go
err := alacritty.ForceConfigure()
```

- **Purpose**: Apply Alacritty configuration regardless of existing files
- **Behavior**: 
  - Copies app configs from template to local dir, overwriting existing
  - Copies default font configuration
  - Copies default theme configuration  
  - Updates config files with current home directory paths
- **Use case**: Reset to devgita defaults, apply config updates

### SoftConfigure()

```go
err := alacritty.SoftConfigure()
```

- **Purpose**: Apply Alacritty configuration only if not already configured
- **Behavior**: Checks for `alacritty.toml` marker file; if exists, does nothing
- **Marker logic**: `filepath.Join(paths.AlacrittyConfigLocalDir, "alacritty.toml")`
- **Use case**: Initial setup that preserves user customizations

## Execution Methods

### ExecuteCommand()

```go
err := alacritty.ExecuteCommand("--version")
err := alacritty.ExecuteCommand("--config-file", "/path/to/config.toml")
err := alacritty.ExecuteCommand("--title", "Development Terminal")
```

- **Purpose**: Execute alacritty commands with provided arguments
- **Parameters**: Variable arguments passed directly to alacritty binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### Alacritty-Specific Operations

The Alacritty CLI provides extensive configuration and execution options:

#### Terminal Management

```bash
# Launch with specific configuration
alacritty --config-file ~/.config/alacritty/alacritty.toml

# Launch with custom title
alacritty --title "Development Environment"

# Launch in specific directory
alacritty --working-directory /path/to/project

# Launch with command
alacritty --command tmux
```

#### Configuration Options

```bash
# Override window dimensions
alacritty --dimensions 120x40

# Set window position
alacritty --position 100,50

# Enable/disable decorations
alacritty --decorations none
alacritty --decorations full

# Set log level
alacritty --log-level debug
alacritty --log-level error
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()`
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()`
3. **Update Configuration**: `New()` → `SoftInstall()` → `ForceConfigure()`
4. **Alacritty Operations**: `New()` → `ExecuteCommand()` with specific alacritty arguments

## Constants and Paths

### Relevant Constants

- `constants.Alacritty`: Package name for Alacritty installation
- Used by all installation methods for consistent desktop app reference

### Configuration Paths

- `paths.AlacrittyConfigAppDir`: Source directory for devgita's Alacritty configuration templates
- `paths.AlacrittyConfigLocalDir`: Target directory for user's Alacritty configuration
- `paths.AlacrittyFontsAppDir`: Source directory for font configuration templates
- `paths.AlacrittyThemesAppDir`: Source directory for theme configuration templates
- `paths.ConfigDir`: Used for placeholder replacement in configuration files

## Implementation Notes

- **Desktop App Installation**: Uses `InstallDesktopApp()` instead of `InstallPackage()` for GUI applications
- **ForceInstall Logic**: Calls `Uninstall()` first but ignores the error since Alacritty uninstall is not supported
- **Configuration Strategy**: Uses marker file (`alacritty.toml`) to determine if configuration exists
- **Multi-Component Setup**: ForceConfigure handles app config, fonts, and themes in sequence
- **Path Substitution**: Updates config files to replace placeholders with actual system paths
- **Error Handling**: All methods return errors that should be checked by callers
- **Cross-Platform**: Works on both macOS and Linux through desktop app installation methods
- **Legacy Compatibility**: Maintains deprecated function aliases to prevent breaking existing code
- **Update Method**: Not implemented as Alacritty updates should be handled by system package managers

## Configuration Structure

The Alacritty configuration system includes multiple components:

### Main Configuration (`alacritty.toml`)

```toml
# Window settings
[window]
opacity = 0.9
option_as_alt = "both"
decorations = "full"

# Font configuration
[font]
size = 14.0

[font.normal]
family = "JetBrains Mono"
style = "Regular"

# Colors and themes
[colors.primary]
background = "#1e1e1e"
foreground = "#ffffff"

# Key bindings
[[keyboard.bindings]]
key = "N"
mods = "Command"
action = "SpawnNewInstance"
```

### Font Configuration (`font.toml`)

```toml
size = 14
family = "JetBrains Mono"
style = "Regular"

[bold]
family = "JetBrains Mono"
style = "Bold"

[italic]
family = "JetBrains Mono"
style = "Italic"
```

### Theme Configuration (`theme.toml`)

```toml
# Color scheme definitions
[colors.primary]
background = "#1e1e1e"
foreground = "#ffffff"

[colors.normal]
black = "#000000"
red = "#ff5555"
green = "#50fa7b"
yellow = "#f1fa8c"
blue = "#bd93f9"
magenta = "#ff79c6"
cyan = "#8be9fd"
white = "#bfbfbf"
```

## Deprecated Functions

The module maintains backward compatibility through deprecated functions:

- `MaybeInstall()` → Use `SoftInstall()` instead
- `SetupApp()` → Use `ForceConfigure()` instead  
- `SetupFont()` → Use `ForceConfigure()` instead
- `SetupTheme()` → Use `ForceConfigure()` instead
- `MaybeSetupApp()` → Use `SoftConfigure()` instead
- `MaybeSetupFont()` → Use `SoftConfigure()` instead
- `MaybeSetupTheme()` → Use `SoftConfigure()` instead
- `UpdateConfigFilesWithCurrentHomeDir()` → Use `ForceConfigure()` instead

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Configuration Not Applied**: Check file permissions in config directory
3. **Commands Don't Work**: Verify Alacritty is installed and accessible in PATH
4. **Font Issues**: Ensure font cache is updated after configuration changes
5. **Theme Not Loading**: Check TOML syntax in theme configuration files

### Platform Considerations

- **macOS**: Installed as `.app` bundle via Homebrew cask
- **Linux**: Installed as desktop application via apt
- **Configuration Location**: Cross-platform config directory handling
- **GPU Acceleration**: May require graphics driver compatibility

This module provides a robust foundation for Alacritty terminal emulator integration within the devgita development environment setup.