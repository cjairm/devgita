# GIMP Module Documentation

## Overview

The GIMP module provides installation and management for GIMP (GNU Image Manipulation Program) desktop application with devgita integration. It follows the standardized devgita app interface while providing GIMP-specific operations for graphics editing software installation.

## App Purpose

GIMP is a free and open-source raster graphics editor used for tasks such as photo retouching, image composition, image authoring, and image editing. This module ensures GIMP is properly installed across macOS (Homebrew cask) and Debian/Ubuntu (apt) systems as a desktop application for design and development workflows.

## Lifecycle Summary

1. **Installation**: Install GIMP desktop application via platform package managers (Homebrew cask/apt)
2. **Configuration**: GIMP uses GUI-based configuration (no config file management needed)
3. **Execution**: GIMP is launched as a desktop application, not via CLI commands

## Exported Functions

| Function           | Purpose                   | Behavior                                                    |
| ------------------ | ------------------------- | ----------------------------------------------------------- |
| `New()`            | Factory method            | Creates new GIMP instance with platform-specific commands   |
| `Install()`        | Standard installation     | Uses `InstallDesktopApp()` to install GIMP                  |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallDesktopApp()` to check before installing  |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                            |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                            |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                           |
| `ExecuteCommand()` | Execute GIMP commands     | **Not applicable** - returns nil                            |
| `Update()`         | Update installation       | **Not implemented** - returns error                         |

## Installation Methods

### Install()

```go
gimp := gimp.New()
err := gimp.Install()
```

- **Purpose**: Standard GIMP installation
- **Behavior**: Uses `InstallDesktopApp()` to install GIMP desktop application
- **Use case**: Initial GIMP installation or explicit reinstall

### ForceInstall()

```go
gimp := gimp.New()
err := gimp.ForceInstall()
```

- **Purpose**: Force GIMP installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error since uninstall is not supported), then `Install()`
- **Use case**: Ensure fresh GIMP installation

### SoftInstall()

```go
gimp := gimp.New()
err := gimp.SoftInstall()
```

- **Purpose**: Install GIMP only if not already present
- **Behavior**: Uses `MaybeInstallDesktopApp()` to check before installing
- **Use case**: Standard installation that respects existing GIMP installations

### Uninstall()

```go
err := gimp.Uninstall()
```

- **Purpose**: Remove GIMP installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Desktop application uninstallation requires platform-specific handling and elevated privileges

### Update()

```go
err := gimp.Update()
```

- **Purpose**: Update GIMP installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: GIMP updates are typically handled by the system package manager or GIMP's built-in updater

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := gimp.ForceConfigure()
err := gimp.SoftConfigure()
```

- **Purpose**: Apply GIMP configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: GIMP configuration is managed through the GIMP GUI preferences. Devgita does not apply default configuration for GIMP to allow users to configure based on their specific design needs.

### Configuration Options

While devgita doesn't apply default configuration, users can customize GIMP via:

- **GUI Preferences**: Edit → Preferences in GIMP
- **Configuration files**: `~/.config/GIMP/` (user-managed)
- **Plugins**: User-installable extensions and filters
- **Brushes/Patterns**: Custom assets in GIMP config directory

## Execution Methods

### ExecuteCommand()

```go
err := gimp.ExecuteCommand("--version")
```

- **Purpose**: Execute GIMP commands
- **Behavior**: **Not applicable** - returns nil (success)
- **Rationale**: GIMP is a desktop application without CLI commands typically managed by devgita, but returns success for interface compliance

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` (fails if GIMP installed) → `ForceConfigure()` (no-op)
3. **Desktop Usage**: Launch GIMP from Applications/Programs menu after installation

## Constants and Paths

### Relevant Constants

- **Package name**: Referenced via `constants.Gimp` (expected to be defined as `"gimp"`)
- Used by all installation methods for consistent desktop app reference

### Configuration Approach

- **GUI-based**: Primary configuration through GIMP GUI preferences
- **Config directory**: `~/.config/GIMP/` (Linux) or `~/Library/Application Support/GIMP/` (macOS)
- **No default config**: Devgita does not apply default configuration for GIMP
- **User customization**: Users configure GIMP based on their specific design workflows

## Implementation Notes

- **Desktop App Installation**: Uses `InstallDesktopApp()` instead of `InstallPackage()` for GUI applications
- **ForceInstall Logic**: Calls `Uninstall()` first but will fail since GIMP uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since GIMP uses GUI configuration
- **Error Handling**: All methods return errors that should be checked by callers
- **Cross-Platform**: Works on both macOS and Linux through desktop app installation methods
- **Update Method**: Not implemented as GIMP updates should be handled by system package managers

## Usage Examples

### Basic GIMP Installation

```go
gimp := gimp.New()

// Install GIMP
err := gimp.SoftInstall()
if err != nil {
    return err
}

// No configuration needed - GIMP uses GUI preferences
err = gimp.SoftConfigure()  // No-op, returns nil
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **GIMP Won't Launch**: Check desktop environment compatibility
3. **Permission Errors**: Some systems require additional permissions for graphics applications
4. **Display Issues**: Verify X11/Wayland configuration on Linux

### Platform Considerations

- **macOS**: Installed via Homebrew cask as GIMP.app
- **Linux**: Installed via apt or other package managers
- **Configuration Location**: Platform-specific config directories
- **Dependencies**: May require graphics libraries and display server

### Prerequisites

Before using GIMP:

- **macOS**: macOS 10.12 (Sierra) or newer
- **Linux**: 64-bit Linux distribution with X11 or Wayland
- **Graphics**: OpenGL-capable graphics card recommended

## External References

- **GIMP Official Website**: https://www.gimp.org/
- **GIMP Documentation**: https://docs.gimp.org/
- **GIMP Tutorials**: https://www.gimp.org/tutorials/
- **GIMP Plugin Registry**: https://registry.gimp.org/

## Integration with Devgita

GIMP integrates with devgita's desktop category:

- **Installation**: Installed as part of desktop applications setup
- **Configuration**: User-managed through GIMP GUI
- **Usage**: Available system-wide after installation as desktop application
- **Updates**: Managed through system package manager or GIMP updater
- **Category**: Design tools in desktop category

## Key Features

### Graphics Editing

- Photo retouching and enhancement
- Image composition and authoring
- Layer-based editing
- Advanced selection tools

### File Format Support

- Native XCF format
- PNG, JPEG, GIF, TIFF, BMP
- PSD (Photoshop) import/export
- SVG and other vector formats

### Extensibility

- Plugin system for additional features
- Script-Fu automation
- Custom brushes and patterns
- Third-party extensions

### Development Use Cases

- UI/UX design mockups
- Icon and asset creation
- Image optimization for web
- Screenshot editing and annotation

This module provides essential graphics editing capabilities for design workflows within the devgita development environment ecosystem.
