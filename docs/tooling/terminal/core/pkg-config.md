# PkgConfig Module Documentation

## Overview

The PkgConfig module provides installation and command execution management for pkg-config helper tool with devgita integration. It follows the standardized devgita app interface while providing pkg-config-specific operations for compiler flag management, library detection, and build system integration.

## App Purpose

pkg-config is a helper tool used when compiling applications and libraries. It helps you insert the correct compiler options on the command line so an application can use `gcc -o test test.c $(pkg-config --libs --cflags glib-2.0)` for instance, rather than hard-coding values on where to find glib (or other libraries). This module ensures pkg-config is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems and provides high-level operations for querying package metadata and build flags.

## Lifecycle Summary

1. **Installation**: Install pkg-config package via platform package managers (Homebrew/apt)
2. **Configuration**: pkg-config typically doesn't require separate configuration files - operations are handled via PKG_CONFIG_PATH environment variable
3. **Execution**: Provide high-level pkg-config operations for querying package information and build flags

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new PkgConfig instance with platform-specific commands       |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install pkg-config                        |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute pkg-config        | Runs pkg-config with provided arguments                              |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
pkgConfig := pkgconfig.New()
err := pkgConfig.Install()
```

- **Purpose**: Standard pkg-config installation
- **Behavior**: Uses `InstallPackage()` to install pkg-config package
- **Use case**: Initial pkg-config installation or explicit reinstall

### ForceInstall()

```go
pkgConfig := pkgconfig.New()
err := pkgConfig.ForceInstall()
```

- **Purpose**: Force pkg-config installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh pkg-config installation or fix corrupted installation

### SoftInstall()

```go
pkgConfig := pkgconfig.New()
err := pkgConfig.SoftInstall()
```

- **Purpose**: Install pkg-config only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing pkg-config installations

### Uninstall()

```go
err := pkgConfig.Uninstall()
```

- **Purpose**: Remove pkg-config installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: pkg-config is a system dependency managed at the OS level

### Update()

```go
err := pkgConfig.Update()
```

- **Purpose**: Update pkg-config installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: pkg-config updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := pkgConfig.ForceConfigure()
err := pkgConfig.SoftConfigure()
```

- **Purpose**: Apply pkg-config configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: pkg-config doesn't use traditional config files; configuration is handled via PKG_CONFIG_PATH environment variable

## Execution Methods

### ExecuteCommand()

```go
err := pkgConfig.ExecuteCommand("--version")
err := pkgConfig.ExecuteCommand("--modversion", "openssl")
err := pkgConfig.ExecuteCommand("--cflags", "--libs", "glib-2.0")
```

- **Purpose**: Execute pkg-config commands with provided arguments
- **Parameters**: Variable arguments passed directly to pkg-config binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### PkgConfig-Specific Operations

The pkg-config CLI provides extensive package metadata and build flag querying capabilities:

#### Version Information

```bash
# Show pkg-config version
pkg-config --version

# Show package version
pkg-config --modversion openssl
pkg-config --modversion glib-2.0
```

#### Compiler and Linker Flags

```bash
# Get compiler flags for a package
pkg-config --cflags openssl
pkg-config --cflags glib-2.0

# Get linker flags for a package
pkg-config --libs openssl
pkg-config --libs glib-2.0

# Get both compiler and linker flags
pkg-config --cflags --libs openssl
```

#### Package Queries

```bash
# Check if package exists
pkg-config --exists openssl
pkg-config --exists glib-2.0

# List all available packages
pkg-config --list-all

# Show package description
pkg-config --print-errors openssl

# Get variable from package
pkg-config --variable=prefix openssl
pkg-config --variable=libdir glib-2.0
```

#### Advanced Options

```bash
# Show errors if package not found
pkg-config --print-errors missing-package

# Get static library flags
pkg-config --static --libs openssl

# Check minimum version requirement
pkg-config --atleast-version=1.1.0 openssl

# Get package dependencies
pkg-config --print-requires openssl
pkg-config --print-requires-private openssl
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **Build System Integration**: `New()` → `SoftInstall()` → `ExecuteCommand()` with pkg-config queries
4. **Version Checks**: `New()` → `ExecuteCommand("--version")`

## Constants and Paths

### Relevant Constants

- **Package name**: Referenced via `constants.PkgConfig` (typically "pkg-config")
- Used by all installation methods for consistent package reference

### Configuration Approach

- **No traditional config files**: pkg-config operations are configured via command-line arguments and environment variables
- **Runtime configuration**: Parameters passed directly to `ExecuteCommand()`
- **Environment variables**:
  - `PKG_CONFIG_PATH`: Search path for .pc files
  - `PKG_CONFIG_LIBDIR`: Override default search directory
  - `PKG_CONFIG_SYSROOT_DIR`: Modify -I and -L paths for cross-compilation

## Implementation Notes

- **Build System Helper Nature**: Unlike typical applications, pkg-config is a helper tool for build systems without persistent configuration
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since pkg-config uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since pkg-config doesn't use config files
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as pkg-config updates should be handled by system package managers

## Usage Examples

### Basic Package Queries

```go
pkgConfig := pkgconfig.New()

// Install pkg-config
err := pkgConfig.SoftInstall()
if err != nil {
    return err
}

// Check version
err = pkgConfig.ExecuteCommand("--version")

// Query package version
err = pkgConfig.ExecuteCommand("--modversion", "openssl")

// Get compiler flags
err = pkgConfig.ExecuteCommand("--cflags", "glib-2.0")

// Get linker flags
err = pkgConfig.ExecuteCommand("--libs", "openssl")
```

### Advanced Operations

```go
// Get both compiler and linker flags
err := pkgConfig.ExecuteCommand("--cflags", "--libs", "glib-2.0")

// Check if package exists
err = pkgConfig.ExecuteCommand("--exists", "openssl")

// List all packages
err = pkgConfig.ExecuteCommand("--list-all")

// Get package variable
err = pkgConfig.ExecuteCommand("--variable=prefix", "openssl")

// Check minimum version
err = pkgConfig.ExecuteCommand("--atleast-version=1.1.0", "openssl")
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Package Not Found**: Verify package development files are installed (e.g., libssl-dev)
3. **PKG_CONFIG_PATH Issues**: Set correct search path for .pc files
4. **Cross-compilation**: Use PKG_CONFIG_SYSROOT_DIR for proper path resolution
5. **Static Linking**: Use --static flag for static library dependencies

### Platform Considerations

- **macOS**: Installed via Homebrew package manager, .pc files typically in /usr/local/lib/pkgconfig
- **Linux**: Installed via apt package manager, .pc files typically in /usr/lib/pkgconfig
- **Development Packages**: Requires -dev or -devel packages to be installed for libraries
- **Search Paths**: pkg-config searches multiple standard paths for .pc files

### Best Practices

- **Check Package Existence**: Use --exists before querying package details
- **Combine Flags**: Use --cflags --libs together for complete build flags
- **Version Checks**: Verify minimum required versions with --atleast-version
- **Build System Integration**: Use $(pkg-config --cflags --libs package) in Makefiles
- **Environment Variables**: Set PKG_CONFIG_PATH for custom .pc file locations

## Integration with Devgita

pkg-config integrates with devgita's terminal category:

- **Installation**: Installed as part of core terminal tools setup
- **Configuration**: No configuration files - uses environment variables
- **Usage**: Available system-wide after installation for build systems
- **Updates**: Managed through system package manager
- **Dependencies**: Required by many development tools and libraries

## External References

- **pkg-config Guide**: https://www.freedesktop.org/wiki/Software/pkg-config/
- **pkg-config Manual**: https://linux.die.net/man/1/pkg-config
- **Best Practices**: https://people.freedesktop.org/~dbn/pkg-config-guide.html
- **Cross-compilation**: https://autotools.io/pkgconfig/cross-compiling.html

This module provides essential build system support for compiling applications and libraries with proper compiler and linker flags within the devgita ecosystem.
