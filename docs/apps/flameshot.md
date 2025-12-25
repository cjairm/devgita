# Flameshot Module Documentation

## Overview

The Flameshot module provides installation and management for Flameshot screenshot tool with devgita integration. It follows the standardized devgita app interface for desktop applications.

## App Purpose

Flameshot is a powerful yet simple to use screenshot software with built-in annotation tools. This module ensures Flameshot is properly installed across macOS (Homebrew cask) and Debian/Ubuntu (apt) systems as a desktop application for capturing and editing screenshots in development workflows.

## Lifecycle Summary

1. **Installation**: Install Flameshot desktop application via platform package managers (Homebrew cask/apt)
2. **Configuration**: Flameshot uses GUI-based configuration (no config file management needed)
3. **Execution**: Flameshot is launched as a desktop application, not via CLI commands

## Exported Functions

| Function           | Purpose                   | Behavior                                                       |
| ------------------ | ------------------------- | -------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Flameshot instance with platform-specific commands |
| `Install()`        | Standard installation     | Uses `InstallDesktopApp()` to install Flameshot                |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error), then `Install()`    |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallDesktopApp()` to check before installing     |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                               |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                               |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                              |
| `ExecuteCommand()` | Execute Flameshot commands| **Not applicable** - returns nil                               |
| `Update()`         | Update installation       | **Not implemented** - returns error                            |

## Installation Methods

### Install()

```go
flameshot := flameshot.New()
err := flameshot.Install()
```

- **Purpose**: Standard Flameshot installation
- **Behavior**: Uses `InstallDesktopApp()` to install Flameshot desktop application
- **Use case**: Initial Flameshot installation or explicit reinstall

### ForceInstall()

```go
flameshot := flameshot.New()
err := flameshot.ForceInstall()
```

- **Purpose**: Force Flameshot installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error since uninstall is not supported), then `Install()`
- **Use case**: Ensure fresh Flameshot installation

### SoftInstall()

```go
flameshot := flameshot.New()
err := flameshot.SoftInstall()
```

- **Purpose**: Install Flameshot only if not already present
- **Behavior**: Uses `MaybeInstallDesktopApp()` to check before installing
- **Use case**: Standard installation that respects existing Flameshot installations

### Uninstall()

```go
err := flameshot.Uninstall()
```

- **Purpose**: Remove Flameshot installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Desktop application uninstallation requires platform-specific handling and elevated privileges

### Update()

```go
err := flameshot.Update()
```

- **Purpose**: Update Flameshot installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Flameshot updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := flameshot.ForceConfigure()
err := flameshot.SoftConfigure()
```

- **Purpose**: Apply Flameshot configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: Flameshot configuration is managed through the GUI preferences. Devgita does not apply default configuration for Flameshot to allow users to configure based on their specific screenshot workflow needs.

### Configuration Options

While devgita doesn't apply default configuration, users can customize Flameshot via:

- **GUI Preferences**: Configuration menu in Flameshot application
- **Configuration file**: `~/.config/flameshot/flameshot.ini` (user-managed)
- **Keyboard shortcuts**: System-level screenshot hotkeys
- **Upload integrations**: Image hosting service configurations

## Execution Methods

### ExecuteCommand()

```go
err := flameshot.ExecuteCommand("--version")
```

- **Purpose**: Execute Flameshot commands
- **Behavior**: **Not applicable** - returns nil (success)
- **Rationale**: Flameshot is a desktop application without CLI commands typically managed by devgita, but returns success for interface compliance

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` (fails if Flameshot installed) → `ForceConfigure()` (no-op)
3. **Desktop Usage**: Launch Flameshot from Applications/Programs menu or system tray after installation

## Constants and Paths

### Relevant Constants

- **Package name**: Defined as `constants.Flameshot` (typically `"flameshot"`)
- Used by all installation methods for consistent desktop app reference

### Configuration Approach

- **GUI-based**: Primary configuration through Flameshot GUI preferences
- **Configuration file**: `~/.config/flameshot/flameshot.ini` (cross-platform)
- **No default config**: Devgita does not apply default configuration for Flameshot
- **User customization**: Users configure Flameshot based on their specific screenshot workflow preferences

## Implementation Notes

- **Desktop App Installation**: Uses `InstallDesktopApp()` instead of `InstallPackage()` for GUI applications
- **ForceInstall Logic**: Calls `Uninstall()` first but will fail since Flameshot uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since Flameshot uses GUI configuration
- **Error Handling**: All methods return errors that should be checked by callers
- **Cross-Platform**: Works on both macOS and Linux through desktop app installation methods
- **Update Method**: Not implemented as Flameshot updates should be handled by system package managers

## Usage Examples

### Basic Flameshot Installation

```go
flameshot := flameshot.New()

// Install Flameshot
err := flameshot.SoftInstall()
if err != nil {
    return err
}

// No configuration needed - Flameshot uses GUI preferences
err = flameshot.SoftConfigure()  // No-op, returns nil
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Flameshot Won't Launch**: Check desktop environment compatibility
3. **Permission Errors**: Some systems require additional permissions for screenshot access
4. **System Tray Missing**: Ensure system tray is available in your desktop environment

### Platform Considerations

- **macOS**: Installed via Homebrew cask as Flameshot.app, requires screen recording permissions
- **Linux**: Installed via apt or other package managers, works with most desktop environments
- **Configuration Location**: `~/.config/flameshot/flameshot.ini`
- **Dependencies**: May require Qt libraries and X11/Wayland display server

### Prerequisites

Before using Flameshot:

- **macOS**: macOS 10.13 (High Sierra) or newer, screen recording permissions required
- **Linux**: 64-bit Linux distribution with X11 or Wayland
- **Desktop Environment**: Compatible with GNOME, KDE, XFCE, and others
- **Display**: Functional display server (X11/Wayland)

## External References

- **Flameshot Official Website**: https://flameshot.org/
- **Flameshot GitHub**: https://github.com/flameshot-org/flameshot
- **Documentation**: https://flameshot.org/docs/
- **Key Bindings**: https://flameshot.org/docs/guide/key-bindings/

## Integration with Devgita

Flameshot integrates with devgita's desktop category:

- **Installation**: Installed as part of desktop applications setup
- **Configuration**: User-managed through GUI preferences
- **Usage**: Available system-wide after installation as desktop application
- **Updates**: Managed through system package manager
- **Category**: Screenshot utility in desktop category

## Key Features

### Screenshot Capabilities

- Fullscreen screenshots
- Region selection screenshots
- Window screenshots
- Delayed screenshots with timer

### Annotation Tools

- Freehand drawing
- Lines, arrows, rectangles, circles
- Text annotations with custom fonts
- Marker and highlighter tools
- Blur and pixelate tools

### Workflow Features

- Copy to clipboard
- Save to file with custom naming
- Upload to image hosting services
- Pin screenshots to screen
- Undo/redo support

### Customization

- Configurable keyboard shortcuts
- Custom save paths and filename patterns
- Interface color customization
- Upload integration (Imgur, custom servers)

### Development Use Cases

- Bug reporting with annotations
- Documentation screenshots
- UI/UX design captures
- Code snippet sharing
- Tutorial creation

This module provides essential screenshot and annotation capabilities for development and productivity workflows within the devgita development environment ecosystem.
