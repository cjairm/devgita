# Ulauncher Module Documentation

## Overview

The Ulauncher module provides installation and management for Ulauncher application launcher with devgita integration. It follows the standardized devgita app interface for desktop applications.

## App Purpose

Ulauncher is a fast application launcher for Linux with customizable shortcuts, extensions, and fuzzy search. This module ensures Ulauncher is properly installed on Debian/Ubuntu systems via apt as a desktop application for productivity workflows. It serves as the Debian/Ubuntu equivalent to Raycast on macOS.

## Lifecycle Summary

1. **Installation**: Install Ulauncher desktop application via apt package manager
2. **Configuration**: Ulauncher uses GUI-based configuration (no config file management needed)
3. **Execution**: Ulauncher is launched as a desktop application, not via CLI commands

## Exported Functions

| Function           | Purpose                   | Behavior                                                       |
| ------------------ | ------------------------- | -------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Ulauncher instance with platform-specific commands |
| `Install()`        | Standard installation     | Uses `InstallDesktopApp()` to install Ulauncher                |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error), then `Install()`    |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallDesktopApp()` to check before installing     |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                               |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                               |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                              |
| `ExecuteCommand()` | Execute Ulauncher commands| **Not applicable** - returns nil                               |
| `Update()`         | Update installation       | **Not implemented** - returns error                            |

## Installation Methods

### Install()

```go
ulauncher := ulauncher.New()
err := ulauncher.Install()
```

- **Purpose**: Standard Ulauncher installation
- **Behavior**: Uses `InstallDesktopApp()` to install Ulauncher desktop application
- **Use case**: Initial Ulauncher installation or explicit reinstall

### SoftInstall()

```go
ulauncher := ulauncher.New()
err := ulauncher.SoftInstall()
```

- **Purpose**: Install Ulauncher only if not already present
- **Behavior**: Uses `MaybeInstallDesktopApp()` to check before installing
- **Use case**: Standard installation that respects existing Ulauncher installations

### ForceInstall()

```go
ulauncher := ulauncher.New()
err := ulauncher.ForceInstall()
```

- **Purpose**: Force Ulauncher installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error since uninstall is not supported), then `Install()`
- **Use case**: Ensure fresh Ulauncher installation

### Uninstall()

```go
err := ulauncher.Uninstall()
```

- **Purpose**: Remove Ulauncher installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Desktop application uninstallation requires platform-specific handling and elevated privileges

### Update()

```go
err := ulauncher.Update()
```

- **Purpose**: Update Ulauncher installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Ulauncher updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := ulauncher.ForceConfigure()
err := ulauncher.SoftConfigure()
```

- **Purpose**: Apply Ulauncher configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: Ulauncher configuration is managed through the application GUI preferences. Devgita does not apply default configuration for Ulauncher to allow users to configure based on their specific workflow needs.

### Configuration Options

While devgita doesn't apply default configuration, users can customize Ulauncher via:

- **GUI Preferences**: Settings within Ulauncher application
- **Configuration directory**: `~/.config/ulauncher/` (user-managed)
- **Extensions**: User-installable Ulauncher extensions
- **Hotkeys**: Custom keyboard shortcuts (default: Ctrl+Space)
- **Themes**: Custom themes and appearance settings

## Execution Methods

### ExecuteCommand()

```go
err := ulauncher.ExecuteCommand("--version")
```

- **Purpose**: Execute Ulauncher commands
- **Behavior**: **Not applicable** - returns nil (success)
- **Rationale**: Ulauncher is a desktop application without CLI commands typically managed by devgita, but returns success for interface compliance

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` (fails if Ulauncher installed) → `ForceConfigure()` (no-op)
3. **Desktop Usage**: Launch Ulauncher via configured hotkey (default: Ctrl+Space) after installation

## Constants and Paths

### Relevant Constants

- **Package name**: Defined as `constants.Ulauncher` (typically `"ulauncher"`)
- Used by all installation methods for consistent desktop app reference

### Configuration Approach

- **GUI-based**: Primary configuration through Ulauncher application preferences
- **Config directory**: `~/.config/ulauncher/` (Linux)
- **No default config**: Devgita does not apply default configuration for Ulauncher
- **User customization**: Users configure Ulauncher based on their specific productivity workflows

## Implementation Notes

- **Desktop App Installation**: Uses `InstallDesktopApp()` instead of `InstallPackage()` for GUI applications
- **Linux-only**: Ulauncher is exclusively available for Linux (Debian/Ubuntu equivalent to macOS Raycast)
- **ForceInstall Logic**: Calls `Uninstall()` first but will fail since Ulauncher uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since Ulauncher uses GUI configuration
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Specificity**: Only works on Debian/Ubuntu via apt package manager
- **Update Method**: Not implemented as Ulauncher updates should be handled by apt

## Usage Examples

### Basic Ulauncher Installation

```go
ulauncher := ulauncher.New()

// Install Ulauncher
err := ulauncher.SoftInstall()
if err != nil {
    return err
}

// No configuration needed - Ulauncher uses GUI preferences
err = ulauncher.SoftConfigure()  // No-op, returns nil
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure apt is available and updated on Debian/Ubuntu
2. **Ulauncher Won't Launch**: Check desktop environment compatibility
3. **Permission Errors**: Some systems require additional permissions
4. **Hotkey Conflicts**: Configure different hotkey in Ulauncher preferences if Ctrl+Space is taken

### Platform Considerations

- **Debian/Ubuntu Only**: Ulauncher is exclusively available for Linux
- **Installation**: Installed via apt package manager
- **Configuration Location**: `~/.config/ulauncher/`
- **Dependencies**: Requires GTK+ 3 and Python 3

### Prerequisites

Before using Ulauncher:

- **Linux**: Debian 12+ (Bookworm) or Ubuntu 24+ with desktop environment
- **Desktop Environment**: Compatible with GNOME, KDE, XFCE, and others
- **Display**: Functional display server (X11/Wayland)

## External References

- **Ulauncher Official Website**: https://ulauncher.io/
- **Ulauncher Documentation**: https://docs.ulauncher.io/
- **Ulauncher GitHub**: https://github.com/Ulauncher/Ulauncher
- **Extension Marketplace**: https://ext.ulauncher.io/

## Integration with Devgita

Ulauncher integrates with devgita's desktop category:

- **Installation**: Installed as part of desktop applications setup (Debian/Ubuntu only)
- **Configuration**: User-managed through Ulauncher GUI
- **Usage**: Available system-wide after installation via configured hotkey
- **Updates**: Managed through apt package manager
- **Category**: Application launcher in desktop category
- **Platform Equivalent**: Debian/Ubuntu alternative to macOS Raycast

## Key Features

### Application Launching

- Fuzzy search for applications
- Custom keywords and shortcuts
- Recent applications history
- Application categories

### Extensions System

- Python-based extension framework
- Community extension marketplace
- Custom extension development
- API integrations (Google, GitHub, etc.)

### Customization

- Custom hotkeys (default: Ctrl+Space)
- Theme support with custom colors
- Custom search providers
- Shortcut aliases

### Productivity Features

- Calculator with unit conversions
- File search and navigation
- Web search shortcuts
- Command execution
- Clipboard management

### Development Use Cases

- Quick access to development tools
- Repository search and navigation
- Documentation lookup
- Code snippet management
- Terminal command shortcuts
- Custom development workflows

## Key Differences from Raycast

| Feature | Raycast (macOS) | Ulauncher (Linux) |
|---------|-----------------|-------------------|
| Extensions | Built-in store | Community marketplace |
| AI Features | Yes | No |
| Customization | Extensive | Theme-based |
| Clipboard History | Yes | Via extensions |
| Window Management | Yes | Via extensions |
| Script Commands | Yes | Custom shortcuts |

This module provides essential application launcher capabilities for Debian/Ubuntu development workflows within the devgita development environment ecosystem.
