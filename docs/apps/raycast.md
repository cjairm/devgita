# Raycast Module Documentation

## Overview

The Raycast module provides installation and management for Raycast productivity launcher with devgita integration. It follows the standardized devgita app interface for desktop applications.

## App Purpose

Raycast is a blazingly fast, extendable launcher for macOS that replaces Spotlight. It provides quick access to applications, files, clipboard history, snippets, and extensible commands. This module ensures Raycast is properly installed on macOS systems via Homebrew cask as a desktop application for productivity and development workflows.

## Lifecycle Summary

1. **Installation**: Install Raycast desktop application via Homebrew cask (macOS only)
2. **Configuration**: Raycast uses GUI-based configuration (no config file management needed)
3. **Execution**: Raycast is launched as a desktop application, not via CLI commands

## Exported Functions

| Function           | Purpose                   | Behavior                                                       |
| ------------------ | ------------------------- | -------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Raycast instance with platform-specific commands   |
| `Install()`        | Standard installation     | Uses `InstallDesktopApp()` to install Raycast                  |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error), then `Install()`    |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallDesktopApp()` to check before installing     |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                               |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                               |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                              |
| `ExecuteCommand()` | Execute Raycast commands  | **Not applicable** - returns nil                               |
| `Update()`         | Update installation       | **Not implemented** - returns error                            |

## Installation Methods

### Install()

```go
raycast := raycast.New()
err := raycast.Install()
```

- **Purpose**: Standard Raycast installation
- **Behavior**: Uses `InstallDesktopApp()` to install Raycast desktop application
- **Use case**: Initial Raycast installation or explicit reinstall

### ForceInstall()

```go
raycast := raycast.New()
err := raycast.ForceInstall()
```

- **Purpose**: Force Raycast installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error since uninstall is not supported), then `Install()`
- **Use case**: Ensure fresh Raycast installation

### SoftInstall()

```go
raycast := raycast.New()
err := raycast.SoftInstall()
```

- **Purpose**: Install Raycast only if not already present
- **Behavior**: Uses `MaybeInstallDesktopApp()` to check before installing
- **Use case**: Standard installation that respects existing Raycast installations

### Uninstall()

```go
err := raycast.Uninstall()
```

- **Purpose**: Remove Raycast installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Desktop application uninstallation requires platform-specific handling and elevated privileges

### Update()

```go
err := raycast.Update()
```

- **Purpose**: Update Raycast installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Raycast updates are typically handled by the system package manager or Raycast's built-in updater

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := raycast.ForceConfigure()
err := raycast.SoftConfigure()
```

- **Purpose**: Apply Raycast configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: Raycast configuration is managed through the application GUI preferences. Devgita does not apply default configuration for Raycast to allow users to configure based on their specific workflow needs.

### Configuration Options

While devgita doesn't apply default configuration, users can customize Raycast via:

- **GUI Preferences**: Settings within Raycast application
- **Configuration directory**: `~/Library/Application Support/com.raycast.macos/` (macOS)
- **Extensions**: User-installable Raycast extensions
- **Hotkeys**: Custom keyboard shortcuts
- **Snippets**: Text expansion snippets

## Execution Methods

### ExecuteCommand()

```go
err := raycast.ExecuteCommand("--version")
```

- **Purpose**: Execute Raycast commands
- **Behavior**: **Not applicable** - returns nil (success)
- **Rationale**: Raycast is a desktop application without CLI commands typically managed by devgita, but returns success for interface compliance

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` (fails if Raycast installed) → `ForceConfigure()` (no-op)
3. **Desktop Usage**: Launch Raycast via configured hotkey (default: Cmd+Space) after installation

## Constants and Paths

### Relevant Constants

- **Package name**: Defined as `constants.Raycast` (typically `"raycast"`)
- Used by all installation methods for consistent desktop app reference

### Configuration Approach

- **GUI-based**: Primary configuration through Raycast application preferences
- **Config directory**: `~/Library/Application Support/com.raycast.macos/` (macOS)
- **No default config**: Devgita does not apply default configuration for Raycast
- **User customization**: Users configure Raycast based on their specific productivity workflows

## Implementation Notes

- **Desktop App Installation**: Uses `InstallDesktopApp()` instead of `InstallPackage()` for GUI applications
- **macOS-only**: Raycast is exclusively available for macOS
- **ForceInstall Logic**: Calls `Uninstall()` first but will fail since Raycast uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since Raycast uses GUI configuration
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Specificity**: Only works on macOS via Homebrew cask
- **Update Method**: Not implemented as Raycast updates should be handled by Homebrew or Raycast's built-in updater

## Usage Examples

### Basic Raycast Installation

```go
raycast := raycast.New()

// Install Raycast
err := raycast.SoftInstall()
if err != nil {
    return err
}

// No configuration needed - Raycast uses GUI preferences
err = raycast.SoftConfigure()  // No-op, returns nil
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure Homebrew is installed and updated on macOS
2. **Raycast Won't Launch**: Check macOS version compatibility
3. **Permission Errors**: Grant accessibility permissions in System Preferences
4. **Hotkey Conflicts**: Configure different hotkey in Raycast preferences if Cmd+Space is taken

### Platform Considerations

- **macOS Only**: Raycast is exclusively available for macOS
- **Installation**: Installed via Homebrew cask as Raycast.app
- **Configuration Location**: `~/Library/Application Support/com.raycast.macos/`
- **Permissions**: Requires accessibility permissions in macOS System Preferences

### Prerequisites

Before using Raycast:

- **macOS**: macOS 11 (Big Sur) or newer
- **Homebrew**: Homebrew package manager installed
- **Permissions**: Accessibility permissions granted in System Preferences → Security & Privacy

## External References

- **Raycast Official Website**: https://www.raycast.com/
- **Raycast Documentation**: https://developers.raycast.com/
- **Raycast GitHub**: https://github.com/raycast
- **Raycast Store**: https://www.raycast.com/store
- **API Documentation**: https://developers.raycast.com/api-reference

## Integration with Devgita

Raycast integrates with devgita's desktop category:

- **Installation**: Installed as part of desktop applications setup (macOS only)
- **Configuration**: User-managed through Raycast GUI
- **Usage**: Available system-wide after installation via configured hotkey
- **Updates**: Managed through Homebrew or Raycast's built-in updater
- **Category**: Productivity launcher in desktop category

## Key Features

### Productivity Tools

- Quick application launcher
- File search and navigation
- Clipboard history manager
- Snippet expansion
- Window management

### Development Features

- Script commands execution
- API integrations
- Custom extensions
- Quicklinks for frequent tasks
- Calculator and unit conversions

### Extensions

- GitHub integration
- Jira and Linear
- Slack and communication tools
- Development tools (npm, Docker, etc.)
- Custom scripts and workflows

### Workflow Automation

- Custom hotkeys
- Alias creation
- Script execution
- API calls
- Text transformations

### System Integration

- System commands
- Application control
- File operations
- Clipboard management
- Notification center

## Key Differences from Spotlight

| Feature | Spotlight | Raycast |
|---------|-----------|---------|
| Extensions | Limited | Extensive store |
| Customization | Minimal | Highly customizable |
| Clipboard History | No | Yes |
| Snippets | No | Yes |
| Window Management | No | Yes |
| Script Commands | No | Yes |
| API Integrations | Limited | Extensive |

## Development Use Cases

- Quick access to development tools
- Repository search and navigation
- PR and issue management
- Documentation lookup
- Code snippet management
- Terminal command execution
- System command shortcuts
- Clipboard history for code snippets
- Custom development workflows

This module provides essential productivity launcher capabilities for macOS development workflows within the devgita development environment ecosystem.
