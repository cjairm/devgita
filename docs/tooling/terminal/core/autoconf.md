# Autoconf Module Documentation

## Overview

The Autoconf module provides installation and command execution management for GNU Autoconf with devgita integration. It follows the standardized devgita app interface while providing autoconf-specific operations for automatic configure script generation, build system configuration, and cross-platform software portability.

## App Purpose

GNU Autoconf is an extensible package of M4 macros that produce shell scripts to automatically configure software source code packages to adapt to many kinds of POSIX-like systems. The configuration scripts produced by Autoconf are independent of Autoconf when they are run, so their users do not need to have Autoconf installed. This module ensures autoconf is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems and provides high-level operations for generating configure scripts and managing build configurations.

## Lifecycle Summary

1. **Installation**: Install autoconf package via platform package managers (Homebrew/apt)
2. **Configuration**: autoconf typically doesn't require separate configuration files - operations are handled via configure.ac in project directories
3. **Execution**: Provide high-level autoconf operations for generating configure scripts and managing build systems

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Autoconf instance with platform-specific commands        |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install autoconf                          |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute autoconf          | Runs autoconf with provided arguments                                |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
autoconf := autoconf.New()
err := autoconf.Install()
```

- **Purpose**: Standard autoconf installation
- **Behavior**: Uses `InstallPackage()` to install autoconf package
- **Use case**: Initial autoconf installation or explicit reinstall

### ForceInstall()

```go
autoconf := autoconf.New()
err := autoconf.ForceInstall()
```

- **Purpose**: Force autoconf installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh autoconf installation or fix corrupted installation

### SoftInstall()

```go
autoconf := autoconf.New()
err := autoconf.SoftInstall()
```

- **Purpose**: Install autoconf only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing autoconf installations

### Uninstall()

```go
err := autoconf.Uninstall()
```

- **Purpose**: Remove autoconf installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: autoconf is a system dependency managed at the OS level

### Update()

```go
err := autoconf.Update()
```

- **Purpose**: Update autoconf installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: autoconf updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := autoconf.ForceConfigure()
err := autoconf.SoftConfigure()
```

- **Purpose**: Apply autoconf configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: autoconf doesn't use traditional config files; configuration is handled via configure.ac files in project directories

## Execution Methods

### ExecuteCommand()

```go
err := autoconf.ExecuteCommand("--version")
err := autoconf.ExecuteCommand("configure.ac")
err := autoconf.ExecuteCommand() // Generate configure from configure.ac
```

- **Purpose**: Execute autoconf commands with provided arguments
- **Parameters**: Variable arguments passed directly to autoconf binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### Autoconf-Specific Operations

The autoconf CLI provides extensive build system configuration capabilities:

#### Version and Help

```bash
# Show autoconf version
autoconf --version

# Show help information
autoconf --help

# Show verbose output
autoconf --verbose
```

#### Generate Configure Scripts

```bash
# Generate configure script from configure.ac
autoconf

# Generate from specific input file
autoconf configure.ac

# Generate with output file
autoconf -o configure configure.ac

# Force regeneration
autoconf --force
```

#### Related Tools

```bash
# Update generated configuration files
autoreconf

# Install missing auxiliary files
autoreconf -i
autoreconf --install

# Force verbose install mode
autoreconf -fvi
autoreconf --force --verbose --install

# Create template header for configure
autoheader

# Trace macro calls
autom4te --trace
```

#### Debug and Trace Options

```bash
# Enable debug mode
autoconf --debug

# Trace macro expansion
autoconf --trace=MACRO

# Show warnings
autoconf --warnings=all

# Prepend directory to search path
autoconf --prepend-include=DIR
```

#### Advanced Usage

```bash
# Use specific M4 program
autoconf --melt

# Freeze M4 files
autoconf --freeze

# Include additional M4 files
autoconf --include=DIR
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **Build System Setup**: `New()` → `SoftInstall()` → `ExecuteCommand()` to generate configure
4. **Version Checks**: `New()` → `ExecuteCommand("--version")`

## Constants and Paths

### Relevant Constants

- **Package name**: Referenced via `constants.Autoconf` (typically "autoconf")
- Used by all installation methods for consistent package reference

### Configuration Approach

- **No traditional config files**: autoconf operations are configured via configure.ac in project directories
- **Runtime configuration**: Parameters passed directly to `ExecuteCommand()`
- **Project-specific**: Each software project has its own configure.ac defining build configuration
- **M4 macros**: Configuration uses M4 macro language for build system flexibility

## Implementation Notes

- **Build System Tool Nature**: Unlike typical applications, autoconf is a build system tool without persistent configuration
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since autoconf uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since autoconf doesn't use config files
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as autoconf updates should be handled by system package managers

## Usage Examples

### Basic Configure Generation

```go
autoconf := autoconf.New()

// Install autoconf
err := autoconf.SoftInstall()
if err != nil {
    return err
}

// Check version
err = autoconf.ExecuteCommand("--version")

// Generate configure script from configure.ac
err = autoconf.ExecuteCommand()

// Generate with specific input
err = autoconf.ExecuteCommand("configure.ac")
```

### Advanced Operations

```go
// Force regeneration
err := autoconf.ExecuteCommand("--force")

// Verbose output
err = autoconf.ExecuteCommand("--verbose")

// Custom output file
err = autoconf.ExecuteCommand("-o", "configure", "configure.ac")

// Debug mode
err = autoconf.ExecuteCommand("--debug")

// Trace macros
err = autoconf.ExecuteCommand("--trace=AC_INIT")
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **configure.ac Not Found**: Run autoconf in directory containing configure.ac file
3. **M4 Errors**: Ensure M4 macro processor is installed
4. **Missing Macros**: Install autoconf-archive or automake for additional macros
5. **Permission Issues**: Verify write permissions for generated configure script

### Platform Considerations

- **macOS**: Installed via Homebrew package manager, includes automake suite
- **Linux**: Installed via apt package manager, often pre-installed on development systems
- **Dependencies**: Requires M4 macro processor (usually installed as dependency)
- **Related Tools**: Often used with automake, libtool for complete build system

### Best Practices

- **Version Control**: Commit configure.ac but not generated configure script
- **Regeneration**: Run autoconf when configure.ac changes
- **Testing**: Test generated configure script on target platforms
- **Documentation**: Document required autoconf version in project README
- **Macros**: Use standard macros from autoconf-archive when possible
- **Portability**: Test generated scripts on multiple POSIX systems

### Build System Workflow

1. **Create configure.ac**: Define build configuration using M4 macros
2. **Run autoconf**: Generate configure script from configure.ac
3. **Run configure**: Execute generated script to create Makefile
4. **Run make**: Build software using generated Makefile

```bash
# Typical autotools workflow
autoconf              # Generate configure from configure.ac
./configure           # Generate Makefile from Makefile.in
make                  # Build the software
make install          # Install the software
```

## Integration with Devgita

autoconf integrates with devgita's terminal category:

- **Installation**: Installed as part of core terminal tools setup
- **Configuration**: No configuration files - uses project-specific configure.ac
- **Usage**: Available system-wide after installation for build systems
- **Updates**: Managed through system package manager
- **Dependencies**: Part of GNU build system toolchain (autoconf, automake, libtool)

## External References

- **GNU Autoconf Manual**: https://www.gnu.org/software/autoconf/manual/
- **Autoconf Documentation**: https://www.gnu.org/savannah-checkouts/gnu/autoconf/manual/autoconf.html
- **Autoconf Macro Archive**: https://www.gnu.org/software/autoconf-archive/
- **Autotools Tutorial**: https://www.gnu.org/software/automake/manual/html_node/Autotools-Introduction.html
- **M4 Manual**: https://www.gnu.org/software/m4/manual/

This module provides essential build system configuration support for generating portable configure scripts within the devgita ecosystem.
