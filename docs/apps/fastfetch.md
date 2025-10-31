# Fastfetch Module Documentation

## Overview

The Fastfetch module provides installation and configuration management for fastfetch with devgita integration. It follows the standardized devgita app interface while providing fastfetch-specific operations for system information display and terminal customization.

## App Purpose

Fastfetch is a neofetch-like tool written in C that fetches system information and displays it prettily in the terminal. This module ensures fastfetch is properly installed and configured with devgita's optimized settings for displaying system information in development environments.

## Lifecycle Summary

1. **Installation**: Install fastfetch package via platform package managers (Homebrew/apt)
2. **Configuration**: Apply devgita's fastfetch configuration templates with developer-focused display settings
3. **Execution**: Provide high-level fastfetch operations for system information display and customization

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Fastfetch instance with platform-specific commands       |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install fastfetch                         |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | Overwrites existing configs with devgita defaults                    |
| `SoftConfigure()`  | Conditional configuration | Preserves existing config.jsonc if present                           |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute fastfetch command | Runs fastfetch with provided arguments                               |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
fastfetch := fastfetch.New()
err := fastfetch.Install()
```

- **Purpose**: Standard fastfetch installation
- **Behavior**: Uses `InstallPackage()` to install fastfetch package
- **Use case**: Initial fastfetch installation or explicit reinstall

### ForceInstall()

```go
fastfetch := fastfetch.New()
err := fastfetch.ForceInstall()
```

- **Purpose**: Force fastfetch installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh fastfetch installation or fix corrupted installation

### SoftInstall()

```go
fastfetch := fastfetch.New()
err := fastfetch.SoftInstall()
```

- **Purpose**: Install fastfetch only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing fastfetch installations

### Uninstall()

```go
err := fastfetch.Uninstall()
```

- **Purpose**: Remove fastfetch installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: System information tools are typically managed at the system level

### Update()

```go
err := fastfetch.Update()
```

- **Purpose**: Update fastfetch installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Fastfetch updates are typically handled by the system package manager

## Configuration Methods

### Configuration Paths

- **Source**: `paths.FastFetchConfigAppDir` (devgita's fastfetch configs)
- **Destination**: `paths.FastFetchConfigLocalDir` (user's config directory)
- **Marker file**: `config.jsonc` in `FastFetchConfigLocalDir`

### ForceConfigure()

```go
err := fastfetch.ForceConfigure()
```

- **Purpose**: Apply fastfetch configuration regardless of existing files
- **Behavior**: Copies all configs from app dir to local dir using `files.CopyDir()`
- **Use case**: Reset to devgita defaults, apply config updates

### SoftConfigure()

```go
err := fastfetch.SoftConfigure()
```

- **Purpose**: Apply fastfetch configuration only if not already configured
- **Behavior**: Checks for `config.jsonc` marker file; if exists, does nothing
- **Marker logic**: `filepath.Join(paths.FastFetchConfigLocalDir, "config.jsonc")`
- **Use case**: Initial setup that preserves user customizations

## Execution Methods

### ExecuteCommand()

```go
err := fastfetch.ExecuteCommand("--version")
err := fastfetch.ExecuteCommand("--config", "/path/to/config.jsonc")
err := fastfetch.ExecuteCommand("--logo", "small")
```

- **Purpose**: Execute fastfetch commands with provided arguments
- **Parameters**: Variable arguments passed directly to fastfetch binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### Fastfetch-Specific Operations

The fastfetch CLI provides extensive system information display options:

#### System Information Display

```bash
# Basic system information
fastfetch

# Use custom configuration
fastfetch --config /path/to/config.jsonc

# Display specific modules only
fastfetch --modules CPU,Memory,Disk

# Use different logo
fastfetch --logo small
fastfetch --logo none
fastfetch --logo /path/to/custom.txt
```

#### Output Formatting

```bash
# JSON output for scripting
fastfetch --format json

# Custom separator
fastfetch --separator " | "

# Disable colors
fastfetch --color-keys never
fastfetch --color-title never

# Set custom title
fastfetch --title "Development Environment"
```

#### Performance Options

```bash
# Fast mode (skip slow modules)
fastfetch --fast

# Disable logo
fastfetch --logo none

# Show only specific information
fastfetch --modules OS,Host,CPU,Memory
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()`
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()`
3. **Update Configuration**: `New()` → `SoftInstall()` → `ForceConfigure()`
4. **Fastfetch Operations**: `New()` → `ExecuteCommand()` with specific fastfetch arguments

## Constants and Paths

### Relevant Constants

- `constants.Fastfetch`: Package name for fastfetch installation
- Used by all installation methods for consistent package reference

### Configuration Paths

- `paths.FastFetchConfigAppDir`: Source directory for devgita's fastfetch configuration templates
- `paths.FastFetchConfigLocalDir`: Target directory for user's fastfetch configuration
- Configuration copying preserves JSON structure and file permissions

## Implementation Notes

- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since fastfetch uninstall is not supported
- **Configuration Strategy**: Uses marker file (`config.jsonc`) to determine if configuration exists
- **Directory Copying**: Uses `files.CopyDir()` for complete configuration directory replication
- **Command Execution**: Uses `BaseCommand.ExecCommand()` with `CommandParams` structure for consistent execution
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as fastfetch updates should be handled by system package managers

## Configuration Structure

The fastfetch configuration (`config.jsonc`) typically includes:

### Basic Configuration

```jsonc
{
  "$schema": "https://github.com/fastfetch-cli/fastfetch/raw/dev/doc/json_schema.json",
  "logo": {
    "source": "auto",
    "padding": {
      "top": 1,
      "left": 4,
    },
  },
  "display": {
    "separator": " ~ ",
  },
}
```

### Module Configuration

```jsonc
{
  "modules": [
    "title",
    "separator",
    "os",
    "host",
    "kernel",
    "uptime",
    "packages",
    "shell",
    "display",
    "de",
    "wm",
    "wmtheme",
    "theme",
    "icons",
    "font",
    "cursor",
    "terminal",
    "terminalfont",
    "cpu",
    "gpu",
    "memory",
    "swap",
    "disk",
    "localip",
    "battery",
    "poweradapter",
    "locale",
    "break",
    "colors",
  ],
}
```

### Custom Display Settings

```jsonc
{
  "display": {
    "separator": " │ ",
    "color": {
      "keys": "blue",
      "title": "yellow",
    },
    "key": {
      "width": 18,
      "paddingLeft": 2,
    },
  },
  "logo": {
    "source": "small",
    "color": {
      "1": "blue",
      "2": "cyan",
    },
  },
}
```

## Deprecated Functions

The module maintains backward compatibility through deprecated functions:

- `MaybeInstall()` → Use `SoftInstall()` instead
- `Setup()` → Use `ForceConfigure()` instead
- `MaybeSetup()` → Use `SoftConfigure()` instead
- `Run()` → Use `ExecuteCommand()` instead

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Configuration Not Applied**: Check file permissions in config directory
3. **Commands Don't Work**: Verify fastfetch is installed and accessible in PATH
4. **JSON Syntax Errors**: Validate config.jsonc syntax and schema compliance
5. **Logo Issues**: Ensure logo files are accessible and properly referenced

### Platform Considerations

- **macOS**: Installed via Homebrew package manager
- **Linux**: Installed via apt package manager
- **Configuration Location**: Cross-platform config directory handling
- **Module Availability**: Some modules may not be available on all platforms

This module provides a robust foundation for fastfetch system information display integration within the devgita development environment setup.

