# Xcode Command Line Tools Module Documentation

## Overview

The Xcode Command Line Tools module provides installation and command execution management for macOS development tools with devgita integration. It follows the standardized devgita app interface while providing macOS-specific operations for essential development toolchain installation.

## App Purpose

Xcode Command Line Tools is a macOS-specific package that provides essential development tools including compilers (clang, gcc), build tools (make, git), SDKs, headers, and frameworks needed for software development on macOS. This module ensures Xcode Command Line Tools are properly installed on macOS systems and provides operations for managing the active developer directory.

## Lifecycle Summary

1. **Installation**: Install Xcode Command Line Tools via `xcode-select --install` on macOS
2. **Configuration**: No configuration required - tools are system-level utilities
3. **Execution**: Provide high-level xcode-select operations for developer directory management

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new XcodeCommandLineTools instance with platform-specific commands |
| `Install()`        | Standard installation     | Installs Xcode Command Line Tools if not present                     |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Checks if installed before attempting installation                   |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute xcode-select      | Runs xcode-select with provided arguments                            |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
xcode := xcode.New()
err := xcode.Install()
```

- **Purpose**: Standard Xcode Command Line Tools installation
- **Behavior**: 
  - Checks if already installed via `xcode-select -p`
  - Returns early if already installed
  - Executes `xcode-select --install` if not present
- **Use case**: Initial installation or explicit reinstall
- **Platform**: macOS only

### ForceInstall()

```go
xcode := xcode.New()
err := xcode.ForceInstall()
```

- **Purpose**: Force installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Not applicable - uninstall not supported
- **Note**: Will always return error due to unsupported uninstall

### SoftInstall()

```go
xcode := xcode.New()
err := xcode.SoftInstall()
```

- **Purpose**: Install Xcode Command Line Tools only if not already present
- **Behavior**: 
  - Checks installation status first
  - Returns nil if already installed
  - Calls `Install()` if not installed
- **Use case**: Standard installation that respects existing installations

### Uninstall()

```go
err := xcode.Uninstall()
```

- **Purpose**: Remove Xcode Command Line Tools installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: System-level development tools require manual management

### Update()

```go
err := xcode.Update()
```

- **Purpose**: Update Xcode Command Line Tools installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Updates are handled through macOS Software Update

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := xcode.ForceConfigure()
err := xcode.SoftConfigure()
```

- **Purpose**: Apply configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: Xcode Command Line Tools are system-level utilities that don't require separate configuration files

## Execution Methods

### ExecuteCommand()

```go
err := xcode.ExecuteCommand("--version")
err := xcode.ExecuteCommand("--print-path")
err := xcode.ExecuteCommand("--switch", "/Applications/Xcode.app/Contents/Developer")
err := xcode.ExecuteCommand("--reset")
```

- **Purpose**: Execute xcode-select commands with provided arguments
- **Parameters**: Variable arguments passed directly to xcode-select binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### Xcode-Select Operations

The xcode-select CLI provides developer directory management capabilities:

#### Information Commands

```bash
# Show version
xcode-select --version

# Print active developer directory
xcode-select --print-path

# Show installation status
xcode-select -p
```

#### Management Commands

```bash
# Install Xcode Command Line Tools
xcode-select --install

# Switch active developer directory
xcode-select --switch /path/to/developer/directory

# Reset to default developer directory
xcode-select --reset
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()`
2. **Force Setup**: `New()` → `ForceInstall()` (will fail - uninstall not supported)
3. **Check Installation**: `New()` → `SoftInstall()` (returns immediately if installed)
4. **Developer Directory Management**: `New()` → `ExecuteCommand()` with xcode-select arguments

## Constants and Paths

### Relevant Constants

- **Package name**: Referenced via `constants.XcodeCommandLineTools`
- Used by installation methods for consistent reference
- Command binary: `xcode-select`

### Installation Paths

Xcode Command Line Tools install to one of these locations:
- `/Library/Developer/CommandLineTools` (standalone installation)
- `/Applications/Xcode.app/Contents/Developer` (full Xcode installation)

The active developer directory is determined by `xcode-select -p`.

## Implementation Notes

- **macOS-Only**: This module is specific to macOS and won't function on other platforms
- **System Integration**: Uses `xcode-select` command-line utility for all operations
- **Installation Check**: Uses `xcode-select -p` to verify installation status
- **Path Detection**: Checks for both "xcode.app" and "commandlinetools" in path output
- **ForceInstall Logic**: Calls `Uninstall()` first and returns error since uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` - no configuration needed
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: While the interface is platform-independent, the functionality is macOS-specific
- **Update Method**: Not implemented as updates should be handled by macOS Software Update

## Platform Considerations

### macOS Requirements

- **Operating System**: macOS 10.9 (Mavericks) or later
- **Installation Method**: Interactive prompt via `xcode-select --install`
- **Disk Space**: Approximately 150-300 MB for Command Line Tools
- **Permissions**: May require administrator password for installation

### What Gets Installed

The Xcode Command Line Tools package includes:
- **Compilers**: clang, gcc, g++
- **Build Tools**: make, autoconf, automake, libtool
- **Version Control**: git, svn
- **SDKs**: macOS SDK headers and frameworks
- **Debuggers**: lldb
- **Other Tools**: ar, ld, nm, ranlib, strip

### Verification

After installation, verify with:
```bash
xcode-select -p
# Should output: /Library/Developer/CommandLineTools

gcc --version
# Should output: Apple clang version...

git --version
# Should output: git version...
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure macOS version is 10.9 or later
2. **Permission Denied**: Installation may require administrator password
3. **Already Installed**: Module checks and skips if already present
4. **Commands Not Found**: Verify installation path with `xcode-select -p`
5. **Wrong Developer Directory**: Use `xcode-select --switch` to change active directory

### Installation Methods

There are several ways to install Xcode Command Line Tools:

1. **Via xcode-select** (this module):
   ```bash
   xcode-select --install
   ```

2. **Via Xcode IDE**:
   - Install full Xcode from Mac App Store
   - Xcode includes Command Line Tools

3. **Direct Download**:
   - Download from Apple Developer website
   - Requires Apple ID

### Developer Directory Management

If you have both standalone tools and full Xcode:

```bash
# Check current active directory
xcode-select -p

# Switch to full Xcode
sudo xcode-select --switch /Applications/Xcode.app/Contents/Developer

# Switch to standalone tools
sudo xcode-select --switch /Library/Developer/CommandLineTools

# Reset to default
sudo xcode-select --reset
```

## Integration with Devgita

Xcode Command Line Tools integrate with devgita's terminal category:

- **Installation**: Installed as part of macOS-specific terminal setup
- **Configuration**: No configuration required
- **Usage**: System-level development tools available globally
- **Updates**: Managed through macOS Software Update
- **Dependencies**: Required for building many terminal tools from source

## External References

- **Apple Documentation**: https://developer.apple.com/xcode/features/
- **Command Line Tools Guide**: https://developer.apple.com/library/archive/technotes/tn2339/_index.html
- **xcode-select Manual**: `man xcode-select`
- **Installation Guide**: https://mac.install.guide/commandlinetools/index.html

## Notes and Assumptions

### Required Constant

This module requires the following constant to be defined in `pkg/constants/constants.go`:

```go
const XcodeCommandLineTools = "xcode-select"
```

### Platform Check (Future Enhancement)

Currently, the module assumes it's running on macOS. A future enhancement would add explicit platform checking:

```go
func (x *XcodeCommandLineTools) Install() error {
    if !x.Base.Platform.IsMac() {
        return fmt.Errorf("xcode command line tools are only available on macOS")
    }
    // ... rest of installation logic
}
```

### Usage in Terminal Coordinator

The terminal coordinator should use this module as follows:

```go
// In internal/tooling/terminal/terminal.go
import "github.com/cjairm/devgita/internal/tooling/terminal/core/xcode"

func (t *Terminal) InstallAll() {
    // ... other setup
    
    if t.Base.Platform.IsMac() {
        xcode := xcode.New()
        displayMessage(xcode.SoftInstall(), constants.XcodeCommandLineTools)
    }
    
    // ... continue with other installations
}
```

This module provides essential macOS development toolchain installation capabilities for creating a functional development environment within the devgita ecosystem.
