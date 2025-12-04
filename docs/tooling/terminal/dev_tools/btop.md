
# Btop Module Documentation

## Overview

The Btop module provides installation and command execution management for btop++ resource monitor with devgita integration. It follows the standardized devgita app interface while providing btop-specific operations for system resource monitoring, process management, and performance visualization.

## App Purpose

Btop++ is a modern, feature-rich resource monitor written in C++ that displays usage and stats for processor, memory, disks, network, and processes. It provides a visually appealing, customizable interface with real-time monitoring capabilities, mouse support, and detailed process information. This module ensures btop is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems and provides high-level operations for launching and managing the resource monitor.

## Lifecycle Summary

1. **Installation**: Install btop package via platform package managers (Homebrew/apt)
2. **Configuration**: Btop uses optional user-specific configuration files (no default configuration applied by devgita)
3. **Execution**: Provide high-level btop operations for launching the resource monitor and managing display settings

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Btop instance with platform-specific commands            |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install btop                              |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute btop commands     | Runs btop with provided arguments                                    |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
btop := btop.New()
err := btop.Install()
```

- **Purpose**: Standard btop installation
- **Behavior**: Uses `InstallPackage()` to install btop package
- **Use case**: Initial btop installation or explicit reinstall

### ForceInstall()

```go
btop := btop.New()
err := btop.ForceInstall()
```

- **Purpose**: Force btop installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh btop installation or fix corrupted installation

### SoftInstall()

```go
btop := btop.New()
err := btop.SoftInstall()
```

- **Purpose**: Install btop only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing btop installations

### Uninstall()

```go
err := btop.Uninstall()
```

- **Purpose**: Remove btop installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: System monitoring tools are typically managed at the system level

### Update()

```go
err := btop.Update()
```

- **Purpose**: Update btop installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Btop updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := btop.ForceConfigure()
err := btop.SoftConfigure()
```

- **Purpose**: Apply btop configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: Btop configuration is optional and user-specific. Configuration files are typically located at `~/.config/btop/btop.conf` and are managed by users based on their preferences.

## Execution Methods

### ExecuteCommand()

```go
err := btop.ExecuteCommand()                    // Launch TUI
err := btop.ExecuteCommand("--version")         // Show version
err := btop.ExecuteCommand("--update", "2000")  // Set update interval to 2000ms
err := btop.ExecuteCommand("--help")            // Show help
```

- **Purpose**: Execute btop commands with provided arguments
- **Parameters**: Variable arguments passed directly to btop binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### Btop-Specific Operations

The btop CLI provides extensive resource monitoring and display options:

#### Launch Interactive TUI

```bash
# Launch btop (default behavior)
btop

# The TUI provides:
# - Real-time CPU, memory, disk, network monitoring
# - Detailed process information with sorting and filtering
# - Process management (kill, terminate, nice adjustment)
# - Mouse support for navigation and interaction
# - Customizable color themes and layouts
# - Historical graphs for resource usage
# - Disk I/O statistics
# - Network bandwidth monitoring
```

#### Command-Line Options

```bash
# Show version information
btop --version

# Set update interval (milliseconds)
btop --update 2000

# Display help information
btop --help

# Use preset theme
btop --preset 0  # Default theme
btop --preset 1  # TTY theme

# Set low color mode
btop --low-color

# Start in TTY mode
btop --tty_on

# UTF-8 force on/off
btop --utf-force

# Set update interval (deprecated, use --update)
btop -u 2000
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **Launch TUI**: `New()` → `SoftInstall()` → `ExecuteCommand()`
4. **Version Check**: `New()` → `ExecuteCommand("--version")`

## Constants and Paths

### Relevant Constants

- **Package name**: Should be defined in `pkg/constants/constants.go` as `Btop` (typically "btop")
- Used by all installation methods for consistent package reference

### Configuration Approach

- **Optional configuration**: Btop uses optional configuration files
- **Default location**: `~/.config/btop/btop.conf` (user-managed)
- **No default config**: Devgita does not apply default configuration for btop
- **User customization**: Users can create their own configuration based on preferences

### Configuration Options

While devgita doesn't apply default configuration, users can customize btop via `~/.config/btop/btop.conf`:

```ini
# Color theme
color_theme = "Default"

# Update interval in milliseconds
update_ms = 2000

# Show CPU temperature
show_cpu_temp = True

# Show disk I/O
show_disks = True

# Show network
show_net = True

# Process sorting
proc_sorting = "cpu lazy"

# Process tree view
proc_tree = False

# Show process memory in bytes
proc_mem_bytes = True

# Check for updates on start
update_check = True

# Log level
log_level = "WARNING"

# Graph symbol
graph_symbol = "braille"

# Rounded corners
rounded_corners = True

# True color support
truecolor = True

# Force TTY mode
force_tty = False

# Vim keys enabled
vim_keys = False
```

## Implementation Notes

- **CLI Tool Nature**: Btop is an interactive terminal UI tool without complex configuration requirements
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since btop uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since btop uses optional user-specific configuration
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as btop updates should be handled by system package managers

## Usage Examples

### Basic Btop Operations

```go
btop := btop.New()

// Install btop
err := btop.SoftInstall()
if err != nil {
    return err
}

// Launch interactive TUI
err = btop.ExecuteCommand()

// Check version
err = btop.ExecuteCommand("--version")
```

### Advanced Usage

```go
// Set custom update interval (2 seconds)
err := btop.ExecuteCommand("--update", "2000")

// Show help
err = btop.ExecuteCommand("--help")

// Launch with low color mode
err = btop.ExecuteCommand("--low-color")
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **TUI Won't Launch**: Verify terminal supports required features (UTF-8, true color)
3. **Display Issues**: Try different color themes or use `--low-color` flag
4. **Commands Don't Work**: Verify btop is installed and accessible in PATH
5. **Performance Issues**: Increase update interval with `--update` flag

### Platform Considerations

- **macOS**: Installed via Homebrew package manager
- **Linux**: Installed via apt package manager
- **Terminal Support**: Works best with modern terminal emulators that support true color and UTF-8
- **System Requirements**: Requires C++20 compatible compiler for building from source

### Prerequisites

Before using btop, ensure:
- Terminal supports UTF-8 encoding
- Terminal supports true color (24-bit color) for best experience
- Sufficient terminal size (recommended: 80x25 minimum)

### Display Requirements

Btop requires certain terminal capabilities for optimal display:
- **UTF-8 support**: For graphs and special characters
- **True color support**: For full color themes
- **Mouse support**: Optional, but enhances interactivity
- **Minimum terminal size**: 80 columns × 25 rows

## Key Features

### CPU Monitoring
- Per-core CPU usage graphs
- CPU temperature monitoring (when available)
- CPU frequency information
- Load average display

### Memory Monitoring
- RAM usage with detailed breakdown
- Swap usage tracking
- Memory graphs with historical data
- Available/used memory display

### Process Management
- Detailed process list with sorting options
- Process tree view
- Process filtering and searching
- Signal sending (kill, terminate, etc.)
- Process priority adjustment

### Disk I/O
- Read/write speeds per disk
- I/O operations per second
- Disk usage statistics
- Mount point information

### Network Monitoring
- Upload/download speeds
- Network interface selection
- Bandwidth graphs
- Total data transferred

### Customization
- Multiple color themes
- Configurable update intervals
- Layout customization
- Vim-style keybindings support
- Mouse support toggle

## External References

- **Btop++ Repository**: https://github.com/aristocratos/btop
- **Configuration Guide**: https://github.com/aristocratos/btop#configurability
- **Keyboard Shortcuts**: https://github.com/aristocratos/btop#keys
- **Themes**: https://github.com/aristocratos/btop#themes

## Integration with Devgita

Btop integrates with devgita's terminal category:

- **Installation**: Installed as part of terminal tools setup
- **Configuration**: User-managed configuration (no default applied)
- **Usage**: Available system-wide after installation
- **Updates**: Managed through system package manager

This module provides essential system resource monitoring capabilities through an intuitive terminal interface, significantly improving system performance visibility and process management within the devgita ecosystem.
