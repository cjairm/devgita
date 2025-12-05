# Libffi Module Documentation

## Overview

The Libffi module provides installation and library management for libffi (Foreign Function Interface) with devgita integration. It follows the standardized devgita app interface while providing libffi-specific operations for enabling runtime foreign function calls, supporting dynamic language interpreters, and facilitating cross-language integration.

## App Purpose

Libffi is a portable, high-level programming interface library that provides a portable interface to various calling conventions. This allows a programmer to call any function specified by a call interface description at runtime. FFI stands for Foreign Function Interface - the popular name for the interface that allows code written in one language to call code written in another language. This module ensures libffi is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems as a critical dependency for dynamic language interpreters (Python, Ruby, Node.js native modules) and runtime code generation systems.

## Lifecycle Summary

1. **Installation**: Install libffi package via platform package managers (Homebrew/apt)
2. **Configuration**: libffi is a library and typically doesn't require separate configuration files - operations are handled at build time or via pkg-config
3. **Execution**: Provides ExecuteCommand() for interface compliance, but libffi is primarily used as a library dependency

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Libffi instance with platform-specific commands          |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install libffi                            |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute libffi commands   | Runs libffi with provided arguments (limited practical use)          |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
libffi := libffi.New()
err := libffi.Install()
```

- **Purpose**: Standard libffi installation
- **Behavior**: Uses `InstallPackage()` to install libffi package
- **Use case**: Initial libffi installation or explicit reinstall for language runtime support

### ForceInstall()

```go
libffi := libffi.New()
err := libffi.ForceInstall()
```

- **Purpose**: Force libffi installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh libffi installation or fix corrupted library installation

### SoftInstall()

```go
libffi := libffi.New()
err := libffi.SoftInstall()
```

- **Purpose**: Install libffi only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing libffi installations

### Uninstall()

```go
err := libffi.Uninstall()
```

- **Purpose**: Remove libffi installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: libffi is a critical system library dependency for many language runtimes

### Update()

```go
err := libffi.Update()
```

- **Purpose**: Update libffi installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: libffi updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := libffi.ForceConfigure()
err := libffi.SoftConfigure()
```

- **Purpose**: Apply libffi configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: libffi is a library without traditional config files; configuration is handled at build time via pkg-config or language-specific build systems

## Execution Methods

### ExecuteCommand()

```go
err := libffi.ExecuteCommand("--version")
err := libffi.ExecuteCommand("--cflags")
```

- **Purpose**: Execute libffi commands with provided arguments
- **Parameters**: Variable arguments passed directly to libffi binary (if available)
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand
- **Note**: libffi is primarily a library, not a CLI tool. This method is provided for interface compliance but has limited practical use cases.

### Libffi-Specific Operations

Libffi is primarily used as a library dependency rather than a command-line tool:

#### Library Integration via pkg-config

```bash
# Get compiler flags for libffi
pkg-config --cflags libffi

# Get linker flags for libffi
pkg-config --libs libffi

# Get libffi version
pkg-config --modversion libffi

# Check if libffi is installed
pkg-config --exists libffi && echo "installed"
```

#### Build System Integration

```bash
# CMake detection
find_package(PkgConfig REQUIRED)
pkg_check_modules(LIBFFI REQUIRED libffi)

# Autotools detection
PKG_CHECK_MODULES([LIBFFI], [libffi >= 3.0])

# Makefile integration
CFLAGS += $(shell pkg-config --cflags libffi)
LDFLAGS += $(shell pkg-config --libs libffi)
```

#### Language Runtime Dependencies

libffi is a critical dependency for:

- **Python**: ctypes, cffi modules
- **Ruby**: FFI gem for native extensions
- **Node.js**: Native addons and N-API
- **GObject Introspection**: Dynamic language bindings
- **LLVM/Clang**: JIT compilation support
- **Java JNA**: Java Native Access library
- **Haskell**: Foreign function interface
- **Common Lisp**: CFFI library

#### Header and Library Locations

```bash
# Typical header location (macOS/Linux)
/usr/include/ffi.h
/usr/include/x86_64-linux-gnu/ffi.h
/usr/local/include/ffi.h

# Typical library location
/usr/lib/libffi.so
/usr/lib/x86_64-linux-gnu/libffi.so
/usr/local/lib/libffi.dylib
/usr/local/lib/libffi.a
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **Language Runtime Dependency**: `New()` → `SoftInstall()` before installing Python/Ruby/etc.
4. **Build System Integration**: Install via package manager, then use pkg-config for build flags

## Constants and Paths

### Relevant Constants

- **Package name**: Referenced via `constants.Libffi` (typically "libffi")
- Used by all installation methods for consistent package reference

### Configuration Approach

- **No traditional config files**: libffi is a library without runtime configuration
- **Build-time configuration**: Uses pkg-config for compiler/linker flags
- **Language-specific integration**: Each language runtime handles libffi integration differently
- **Header files**: Installed to system include directories for C/C++ compilation
- **Shared libraries**: Installed to system library directories for dynamic linking

### Library Paths

- **pkg-config data**: `/usr/lib/pkgconfig/libffi.pc` or `/usr/local/lib/pkgconfig/libffi.pc`
- **Header files**: System include directories (discovered via pkg-config)
- **Shared libraries**: System library directories (discovered via pkg-config)

## Implementation Notes

- **Library Nature**: Unlike typical applications, libffi is a library dependency without standalone CLI functionality
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since libffi uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since libffi doesn't use config files
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as libffi updates should be handled by system package managers
- **Critical Dependency**: Removing libffi can break Python, Ruby, and other language runtimes

## Usage Examples

### Installation as Dependency

```go
libffi := libffi.New()

// Install libffi before language runtimes
err := libffi.SoftInstall()
if err != nil {
    return err
}

// Now safe to install Python, Ruby, etc.
// These runtimes will use the installed libffi
```

### Integration with Language Installation

```go
// In language installer coordinator
func InstallPython(ctx context.Context) error {
    // Ensure libffi is available first
    libffi := libffi.New()
    if err := libffi.SoftInstall(); err != nil {
        return fmt.Errorf("failed to install libffi dependency: %w", err)
    }
    
    // Now install Python which depends on libffi
    python := python.New()
    return python.SoftInstall()
}
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Python ctypes Fails**: Verify libffi is installed and headers are available
3. **Ruby FFI Gem Fails**: Check libffi development package is installed
4. **pkg-config Not Found**: Install pkg-config package for build system integration
5. **Header Files Missing**: Install libffi-dev (Debian/Ubuntu) or libffi-devel (RHEL/CentOS)

### Platform Considerations

- **macOS**: Installed via Homebrew as `libffi`, includes headers and pkg-config
- **Linux (Debian/Ubuntu)**: Install both `libffi` (runtime) and `libffi-dev` (development)
- **Linux (RHEL/CentOS)**: Install both `libffi` (runtime) and `libffi-devel` (development)
- **Dependencies**: Required by many language runtimes - avoid uninstalling

### Development vs Runtime Packages

On Linux systems, libffi typically requires both packages:

- **Runtime package**: `libffi` or `libffi8` - Shared library for running programs
- **Development package**: `libffi-dev` or `libffi-devel` - Headers and pkg-config for building

```bash
# Debian/Ubuntu - Install both runtime and development
sudo apt-get install libffi-dev

# RHEL/CentOS - Install both runtime and development
sudo yum install libffi-devel

# macOS - Homebrew includes both
brew install libffi
```

### Verifying Installation

```bash
# Check if libffi is installed
pkg-config --exists libffi && echo "libffi installed"

# Show libffi version
pkg-config --modversion libffi

# Show compiler flags
pkg-config --cflags libffi

# Show linker flags
pkg-config --libs libffi

# Find library files
find /usr -name "libffi*" 2>/dev/null

# Check for header file
find /usr -name "ffi.h" 2>/dev/null
```

### Language Runtime Issues

**Python ctypes not working**:
```bash
# Verify libffi is available
python -c "import ctypes; print('ctypes works')"

# Check Python was built with libffi
python -c "import sysconfig; print(sysconfig.get_config_var('LIBFFI_INCLUDEDIR'))"
```

**Ruby FFI gem installation fails**:
```bash
# Install libffi development package first
sudo apt-get install libffi-dev  # Debian/Ubuntu
brew install libffi              # macOS

# Then install Ruby FFI gem
gem install ffi
```

**Node.js native addon build fails**:
```bash
# Ensure libffi headers are available
pkg-config --cflags libffi

# Rebuild native addons
npm rebuild
```

## Integration with Devgita

libffi integrates with devgita's terminal category as a core dependency:

- **Installation**: Installed as part of core terminal tools setup before language runtimes
- **Configuration**: No configuration files - library is automatically available after installation
- **Usage**: Transparent to users - used internally by Python, Ruby, Node.js, and other runtimes
- **Updates**: Managed through system package manager alongside OS updates
- **Dependencies**: Critical for language runtime installations - installed early in setup process

### Installation Order

Devgita installs libffi early in the terminal tools setup:

1. **Core libraries**: pkg-config, autoconf, libffi, openssl, readline, etc.
2. **Runtime managers**: mise (depends on proper library setup)
3. **Language runtimes**: Python, Ruby, Node.js (depend on libffi)

This ensures language runtimes have all required dependencies when installed.

## External References

- **Libffi Official Site**: https://sourceware.org/libffi/
- **Libffi GitHub**: https://github.com/libffi/libffi
- **Libffi Documentation**: https://www.sourceware.org/libffi/
- **Python ctypes**: https://docs.python.org/3/library/ctypes.html
- **Ruby FFI**: https://github.com/ffi/ffi
- **pkg-config Guide**: https://people.freedesktop.org/~dbn/pkg-config-guide.html

This module provides essential foreign function interface library support for dynamic language runtimes and cross-language integration within the devgita ecosystem.
