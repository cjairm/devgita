# Brave Browser Module Documentation

## Overview

The Brave module provides installation and management for Brave browser desktop application with devgita integration. It follows the standardized devgita app interface for desktop applications.

## App Purpose

Brave is a privacy-focused web browser built on Chromium that blocks ads and trackers by default. This module ensures Brave browser is properly installed across macOS (Homebrew cask) and Debian/Ubuntu (apt) systems as a desktop application for web browsing and development workflows.

## Lifecycle Summary

1. **Installation**: Install Brave desktop application via platform package managers (Homebrew cask/apt)
2. **Configuration**: Brave uses GUI-based configuration (no config file management needed)
3. **Execution**: Brave is launched as a desktop application, not via CLI commands

## Exported Functions

| Function           | Purpose                   | Behavior                                                    |
| ------------------ | ------------------------- | ----------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Brave instance with platform-specific commands  |
| `Install()`        | Standard installation     | Uses `InstallDesktopApp()` to install Brave                 |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallDesktopApp()` to check before installing  |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                            |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                            |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                           |
| `ExecuteCommand()` | Execute Brave commands    | **Not applicable** - returns nil                            |
| `Update()`         | Update installation       | **Not implemented** - returns error                         |

## Installation Methods

### Install()

```go
brave := brave.New()
err := brave.Install()
```

- **Purpose**: Standard Brave browser installation
- **Behavior**: Uses `InstallDesktopApp()` to install Brave desktop application
- **Use case**: Initial Brave installation or explicit reinstall

### ForceInstall()

```go
brave := brave.New()
err := brave.ForceInstall()
```

- **Purpose**: Force Brave installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error since uninstall is not supported), then `Install()`
- **Use case**: Ensure fresh Brave installation

### SoftInstall()

```go
brave := brave.New()
err := brave.SoftInstall()
```

- **Purpose**: Install Brave only if not already present
- **Behavior**: Uses `MaybeInstallDesktopApp()` to check before installing
- **Use case**: Standard installation that respects existing Brave installations

### Uninstall()

```go
err := brave.Uninstall()
```

- **Purpose**: Remove Brave installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Desktop application uninstallation requires platform-specific handling and elevated privileges

### Update()

```go
err := brave.Update()
```

- **Purpose**: Update Brave installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Brave updates are typically handled by the system package manager or Brave's built-in updater

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := brave.ForceConfigure()
err := brave.SoftConfigure()
```

- **Purpose**: Apply Brave configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: Brave configuration is managed through the browser GUI preferences. Devgita does not apply default configuration for Brave to allow users to configure based on their specific browsing needs.

### Configuration Options

While devgita doesn't apply default configuration, users can customize Brave via:

- **GUI Preferences**: Settings menu in Brave browser
- **Profile settings**: `~/Library/Application Support/BraveSoftware/Brave-Browser/` (macOS)
- **Extensions**: User-installable browser extensions
- **Sync settings**: Brave Sync for cross-device configuration

## Execution Methods

### ExecuteCommand()

```go
err := brave.ExecuteCommand("--version")
```

- **Purpose**: Execute Brave commands
- **Behavior**: **Not applicable** - returns nil (success)
- **Rationale**: Brave is a desktop application without CLI commands typically managed by devgita, but returns success for interface compliance

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` (fails if Brave installed) → `ForceConfigure()` (no-op)
3. **Desktop Usage**: Launch Brave from Applications/Programs menu after installation

## Constants and Paths

### Relevant Constants

- **Package name**: Expected to be defined as `constants.Brave` (typically `"brave-browser"` or `"brave"`)
- Used by all installation methods for consistent desktop app reference

### Configuration Approach

- **GUI-based**: Primary configuration through Brave browser GUI preferences
- **Profile directory**: 
  - macOS: `~/Library/Application Support/BraveSoftware/Brave-Browser/`
  - Linux: `~/.config/BraveSoftware/Brave-Browser/`
- **No default config**: Devgita does not apply default configuration for Brave
- **User customization**: Users configure Brave based on their specific browsing preferences

## Implementation Notes

- **Desktop App Installation**: Uses `InstallDesktopApp()` instead of `InstallPackage()` for GUI applications
- **ForceInstall Logic**: Calls `Uninstall()` first but will fail since Brave uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since Brave uses GUI configuration
- **Error Handling**: All methods return errors that should be checked by callers
- **Cross-Platform**: Works on both macOS and Linux through desktop app installation methods
- **Update Method**: Not implemented as Brave updates should be handled by system package managers or browser updater

## Usage Examples

### Basic Brave Installation

```go
brave := brave.New()

// Install Brave
err := brave.SoftInstall()
if err != nil {
    return err
}

// No configuration needed - Brave uses GUI preferences
err = brave.SoftConfigure()  // No-op, returns nil
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Brave Won't Launch**: Check desktop environment compatibility
3. **Permission Errors**: Some systems require additional permissions for browser installation
4. **Missing Dependencies**: Verify system has required graphics and system libraries

### Platform Considerations

- **macOS**: Installed via Homebrew cask as Brave-Browser.app
- **Linux**: Installed via apt or other package managers
- **Configuration Location**: Platform-specific profile directories
- **Dependencies**: May require graphics libraries and display server

### Prerequisites

Before using Brave:

- **macOS**: macOS 10.13 (High Sierra) or newer
- **Linux**: 64-bit Linux distribution with X11 or Wayland
- **Graphics**: Modern graphics card with driver support
- **Display**: Functional display server (X11/Wayland)

## External References

- **Brave Official Website**: https://brave.com/
- **Brave Browser Documentation**: https://support.brave.com/
- **Brave GitHub**: https://github.com/brave/brave-browser
- **Privacy Features**: https://brave.com/privacy-features/

## Integration with Devgita

Brave integrates with devgita's desktop category:

- **Installation**: Installed as part of desktop applications setup
- **Configuration**: User-managed through browser GUI
- **Usage**: Available system-wide after installation as desktop application
- **Updates**: Managed through system package manager or Brave's built-in updater
- **Category**: Web browser in desktop category

## Key Features

### Privacy & Security

- Built-in ad and tracker blocking
- HTTPS Everywhere integration
- Fingerprinting protection
- Cookie control and management

### Performance

- Faster page loading (blocked ads/trackers)
- Lower memory usage vs other browsers
- Built-in performance optimizations
- Resource-efficient design

### Development Features

- Chromium DevTools integration
- Web3 support (crypto wallets)
- Tor private browsing mode
- Developer-friendly extensions

### Brave Rewards

- Opt-in privacy-respecting ads
- BAT (Basic Attention Token) rewards
- Creator support and tipping
- Privacy-preserving analytics

This module provides essential privacy-focused web browsing capabilities for development and productivity workflows within the devgita development environment ecosystem.
