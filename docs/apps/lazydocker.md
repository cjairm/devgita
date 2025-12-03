# Lazydocker Module Documentation

## Overview

The Lazydocker module provides installation and command execution management for lazydocker terminal UI with devgita integration. It follows the standardized devgita app interface while providing lazydocker-specific operations for Docker container and image management through an interactive terminal interface.

## App Purpose

Lazydocker is a simple terminal UI for both docker and docker-compose, written in Go with the gocui library. It provides an interactive interface to manage Docker containers, images, volumes, and networks, making Docker management more accessible and efficient directly from the terminal. This module ensures lazydocker is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems and provides high-level operations for launching the interactive TUI.

## Lifecycle Summary

1. **Installation**: Install lazydocker package via platform package managers (Homebrew/apt)
2. **Configuration**: Lazydocker uses optional user-specific configuration files (no default configuration applied by devgita)
3. **Execution**: Provide high-level lazydocker operations for launching the interactive TUI and managing Docker resources

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Lazydocker instance with platform-specific commands      |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install lazydocker                        |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute lazydocker        | Runs lazydocker with provided arguments                              |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
lazydocker := lazydocker.New()
err := lazydocker.Install()
```

- **Purpose**: Standard lazydocker installation
- **Behavior**: Uses `InstallPackage()` to install lazydocker package
- **Use case**: Initial lazydocker installation or explicit reinstall

### ForceInstall()

```go
lazydocker := lazydocker.New()
err := lazydocker.ForceInstall()
```

- **Purpose**: Force lazydocker installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh lazydocker installation or fix corrupted installation

### SoftInstall()

```go
lazydocker := lazydocker.New()
err := lazydocker.SoftInstall()
```

- **Purpose**: Install lazydocker only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing lazydocker installations

### Uninstall()

```go
err := lazydocker.Uninstall()
```

- **Purpose**: Remove lazydocker installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Docker tools are typically managed at the system level

### Update()

```go
err := lazydocker.Update()
```

- **Purpose**: Update lazydocker installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Lazydocker updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := lazydocker.ForceConfigure()
err := lazydocker.SoftConfigure()
```

- **Purpose**: Apply lazydocker configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: Lazydocker configuration is optional and user-specific. Configuration files are typically located at `~/.config/lazydocker/config.yml` and are managed by users based on their preferences.

## Execution Methods

### ExecuteCommand()

```go
err := lazydocker.ExecuteCommand()                    // Launch TUI
err := lazydocker.ExecuteCommand("--version")         // Show version
err := lazydocker.ExecuteCommand("--config")          // Show config path
err := lazydocker.ExecuteCommand("--help")            // Show help
```

- **Purpose**: Execute lazydocker commands with provided arguments
- **Parameters**: Variable arguments passed directly to lazydocker binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### Lazydocker-Specific Operations

The lazydocker CLI provides interactive Docker management capabilities:

#### Launch Interactive TUI

```bash
# Launch lazydocker (default behavior)
lazydocker

# The TUI provides:
# - Container management (start, stop, restart, remove)
# - Image management (remove, prune)
# - Volume management
# - Network management
# - Log viewing
# - Stats monitoring
# - Docker Compose support
```

#### Command-Line Options

```bash
# Show version information
lazydocker --version

# Show configuration file path
lazydocker --config

# Display help information
lazydocker --help

# Use custom config file
lazydocker --config-file /path/to/config.yml

# Debug mode
lazydocker --debug
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **Launch TUI**: `New()` → `SoftInstall()` → `ExecuteCommand()`
4. **Version Check**: `New()` → `ExecuteCommand("--version")`

## Constants and Paths

### Relevant Constants

- **Package name**: Referenced via `constants.Lazydocker` (typically "lazydocker")
- Used by all installation methods for consistent package reference

### Configuration Approach

- **Optional configuration**: Lazydocker uses optional configuration files
- **Default location**: `~/.config/lazydocker/config.yml` (user-managed)
- **No default config**: Devgita does not apply default configuration for lazydocker
- **User customization**: Users can create their own configuration based on preferences

### Configuration Options

While devgita doesn't apply default configuration, users can customize lazydocker via `~/.config/lazydocker/config.yml`:

```yaml
# Example lazydocker configuration
gui:
  theme:
    activeBorderColor:
      - green
      - bold
    inactiveBorderColor:
      - white
  returnImmediately: false
  wrapMainPanel: true

stats:
  graphs:
    - caption: CPU (%)
      statPath: DerivedStats.CPUPercentage
      color: cyan
    - caption: Memory (%)
      statPath: DerivedStats.MemoryPercentage
      color: green

logs:
  timestamps: false
  since: '60m'
  tail: '200'
```

## Implementation Notes

- **CLI Tool Nature**: Lazydocker is an interactive terminal UI tool without complex configuration requirements
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since lazydocker uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since lazydocker uses optional user-specific configuration
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as lazydocker updates should be handled by system package managers

## Usage Examples

### Basic Lazydocker Operations

```go
lazydocker := lazydocker.New()

// Install lazydocker
err := lazydocker.SoftInstall()
if err != nil {
    return err
}

// Launch interactive TUI
err = lazydocker.ExecuteCommand()

// Check version
err = lazydocker.ExecuteCommand("--version")
```

### Advanced Usage

```go
// Check configuration file location
err := lazydocker.ExecuteCommand("--config")

// Show help
err = lazydocker.ExecuteCommand("--help")
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **TUI Won't Launch**: Verify Docker is installed and running
3. **Permission Issues**: Ensure user has permissions to access Docker socket
4. **Commands Don't Work**: Verify lazydocker is installed and accessible in PATH
5. **Docker Not Found**: Install Docker Desktop (macOS) or Docker Engine (Linux)

### Platform Considerations

- **macOS**: Installed via Homebrew package manager, requires Docker Desktop
- **Linux**: Installed via apt package manager, requires Docker Engine
- **Docker Dependency**: Lazydocker requires Docker to be installed and running
- **Terminal Support**: Works best with modern terminal emulators that support colors

### Prerequisites

Before using lazydocker, ensure:
- Docker is installed and running
- User has permissions to access Docker (may need to add user to docker group on Linux)
- Terminal supports colors and special characters for optimal TUI experience

### Docker Integration

Lazydocker connects to Docker via:
- **Docker socket**: Unix socket on Linux/macOS
- **Docker daemon**: Must be running before launching lazydocker
- **Docker context**: Respects current Docker context configuration

## External References

- **Lazydocker Repository**: https://github.com/jesseduffield/lazydocker
- **Configuration Guide**: https://github.com/jesseduffield/lazydocker/blob/master/docs/Config.md
- **Keybindings**: https://github.com/jesseduffield/lazydocker/blob/master/docs/keybindings/Keybindings_en.md
- **Docker Documentation**: https://docs.docker.com/

## Integration with Devgita

Lazydocker integrates with devgita's terminal category:

- **Installation**: Installed as part of terminal tools setup
- **Configuration**: User-managed configuration (no default applied)
- **Usage**: Available system-wide after installation
- **Updates**: Managed through system package manager
- **Dependencies**: Requires Docker to be installed separately

This module provides essential Docker container and image management capabilities through an intuitive terminal interface, significantly improving Docker workflow efficiency within the devgita ecosystem.
