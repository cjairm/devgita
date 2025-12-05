# GDBM Module Documentation

## Overview

The GDBM module provides installation and library management for GNU dbm (GDBM) with devgita integration. It follows the standardized devgita app interface while providing gdbm-specific operations for enabling persistent key-value database storage, supporting language database modules, and facilitating system-level data persistence.

## App Purpose

GNU dbm (GDBM) is a library of database functions that use extensible hashing and work similar to the standard UNIX dbm. These routines are provided to a programmer needing to create and manipulate a hashed database. GDBM provides a key-value database storage mechanism with extensible hashing, and is required by many system tools, package managers, and language interpreters (Perl, Python, Ruby). This module ensures gdbm is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems as a critical dependency for language database modules and system configuration tools.

## Lifecycle Summary

1. **Installation**: Install gdbm package via platform package managers (Homebrew/apt)
2. **Configuration**: gdbm is a library and typically doesn't require separate configuration files - operations are handled at build time or via language bindings
3. **Execution**: Provides ExecuteCommand() for interface compliance, but gdbm is primarily used as a library dependency

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Gdbm instance with platform-specific commands            |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install gdbm                              |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute gdbm commands     | Runs gdbm with provided arguments (limited practical use)            |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
gdbm := gdbm.New()
err := gdbm.Install()
```

- **Purpose**: Standard gdbm installation
- **Behavior**: Uses `InstallPackage()` to install gdbm package
- **Use case**: Initial gdbm installation or explicit reinstall for language runtime support

### ForceInstall()

```go
gdbm := gdbm.New()
err := gdbm.ForceInstall()
```

- **Purpose**: Force gdbm installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh gdbm installation or fix corrupted library installation

### SoftInstall()

```go
gdbm := gdbm.New()
err := gdbm.SoftInstall()
```

- **Purpose**: Install gdbm only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing gdbm installations

### Uninstall()

```go
err := gdbm.Uninstall()
```

- **Purpose**: Remove gdbm installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: gdbm is a critical system library dependency for many language runtimes and system tools

### Update()

```go
err := gdbm.Update()
```

- **Purpose**: Update gdbm installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: gdbm updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := gdbm.ForceConfigure()
err := gdbm.SoftConfigure()
```

- **Purpose**: Apply gdbm configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: gdbm is a library without traditional config files; configuration is handled at build time via language-specific bindings

## Execution Methods

### ExecuteCommand()

```go
err := gdbm.ExecuteCommand("--version")
err := gdbm.ExecuteCommand("--check", "database.db")
```

- **Purpose**: Execute gdbm commands with provided arguments
- **Parameters**: Variable arguments passed directly to gdbm binary (if available)
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand
- **Note**: gdbm is primarily a library, not a CLI tool. This method is provided for interface compliance but has limited practical use cases.

### GDBM-Specific Operations

GDBM is primarily used as a library dependency rather than a command-line tool:

#### Library Integration via Language Bindings

**Perl Integration**:

```perl
use GDBM_File;
tie %hash, 'GDBM_File', 'database.db', &GDBM_WRCREAT, 0640;
$hash{'key'} = 'value';
untie %hash;
```

**Python Integration**:

```python
import dbm.gnu

# Open/create database
db = dbm.gnu.open('database.db', 'c')
db['key'] = 'value'
db.close()
```

**Ruby Integration**:

```ruby
require 'gdbm'

GDBM.open('database.db') do |db|
  db['key'] = 'value'
end
```

#### Command-Line Tools (if available)

Some systems provide gdbm command-line utilities:

```bash
# Create or modify GDBM database (gdbmtool)
gdbmtool database.db

# Export GDBM database to flat file
gdbm_dump database.db > output.txt

# Import flat file to GDBM database
gdbm_load database.db < input.txt

# Check database integrity
gdbmtool --check database.db
```

#### Build System Integration

**CMake Detection**:

```cmake
find_package(GDBM REQUIRED)
include_directories(${GDBM_INCLUDE_DIRS})
target_link_libraries(myapp ${GDBM_LIBRARIES})
```

**Autotools Detection**:

```bash
# configure.ac
AC_CHECK_LIB([gdbm], [gdbm_open], [], [
  AC_MSG_ERROR([GDBM library not found])
])
```

**pkg-config Integration**:

```bash
# Get compiler flags (if available)
pkg-config --cflags gdbm

# Get linker flags (if available)
pkg-config --libs gdbm
```

#### Language Runtime Dependencies

gdbm is a critical dependency for:

- **Perl**: DBM::Deep, GDBM_File modules
- **Python**: dbm.gnu module (part of standard library)
- **Ruby**: GDBM gem for persistent key-value storage
- **PHP**: dba extension with gdbm handler
- **Scheme/Guile**: gdbm bindings for persistent data
- **System Tools**: Package managers, mail systems, configuration databases

#### Database File Operations

```bash
# Typical GDBM database file locations
~/.cache/*/database.db
/var/lib/*/database.db
/usr/local/var/*/database.db

# GDBM database file characteristics
# - Binary format, not human-readable
# - Uses extensible hashing for performance
# - Single-writer, multiple-reader access
# - Atomic operations for data integrity
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **Language Runtime Dependency**: `New()` → `SoftInstall()` before installing Perl/Python/Ruby
4. **System Tool Dependency**: Install via package manager before system tools that require dbm

## Constants and Paths

### Relevant Constants

- **Package name**: Referenced via `constants.Gdbm` (typically "gdbm")
- Used by all installation methods for consistent package reference

### Configuration Approach

- **No traditional config files**: gdbm is a library without runtime configuration
- **Build-time configuration**: Uses language-specific build systems for integration
- **Language-specific integration**: Each language runtime handles gdbm binding differently
- **Header files**: Installed to system include directories for C/C++ compilation
- **Shared libraries**: Installed to system library directories for dynamic linking

### Library Paths

- **Header files**: `/usr/include/gdbm.h`, `/usr/local/include/gdbm.h`
- **Shared libraries**: `/usr/lib/libgdbm.so`, `/usr/local/lib/libgdbm.dylib`
- **pkg-config data**: `/usr/lib/pkgconfig/gdbm.pc` (if available)
- **Database files**: Created by applications in various locations

## Implementation Notes

- **Library Nature**: Unlike typical applications, gdbm is a library dependency without standalone CLI functionality
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since gdbm uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since gdbm doesn't use config files
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as gdbm updates should be handled by system package managers
- **Critical Dependency**: Removing gdbm can break Perl, Python, Ruby, and system tools

## Usage Examples

### Installation as Dependency

```go
gdbm := gdbm.New()

// Install gdbm before language runtimes
err := gdbm.SoftInstall()
if err != nil {
    return err
}

// Now safe to install Perl, Python, Ruby, etc.
// These runtimes will use the installed gdbm library
```

### Integration with Language Installation

```go
// In language installer coordinator
func InstallPython(ctx context.Context) error {
    // Ensure gdbm is available first
    gdbm := gdbm.New()
    if err := gdbm.SoftInstall(); err != nil {
        return fmt.Errorf("failed to install gdbm dependency: %w", err)
    }

    // Now install Python which may use gdbm for dbm module
    python := python.New()
    return python.SoftInstall()
}
```

### Using GDBM in Applications

```go
// Example: Using gdbm via Python's dbm module
func UsePythonDBM() error {
    python := exec.Command("python3", "-c", `
import dbm.gnu
db = dbm.gnu.open('data.db', 'c')
db['key'] = 'value'
print(db['key'])
db.close()
    `)
    return python.Run()
}
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Python dbm.gnu Fails**: Verify gdbm is installed before Python was built
3. **Perl GDBM_File Missing**: Check gdbm development package is installed
4. **Ruby GDBM Gem Fails**: Install gdbm library before building Ruby native extensions
5. **Database Corruption**: Use gdbmtool --check to verify database integrity

### Platform Considerations

- **macOS**: Installed via Homebrew as `gdbm`, includes headers
- **Linux (Debian/Ubuntu)**: Install both `libgdbm6` (runtime) and `libgdbm-dev` (development)
- **Linux (RHEL/CentOS)**: Install both `gdbm` (runtime) and `gdbm-devel` (development)
- **Dependencies**: Required by many language runtimes - avoid uninstalling

### Development vs Runtime Packages

On Linux systems, gdbm typically requires both packages:

- **Runtime package**: `libgdbm6` or `gdbm` - Shared library for running programs
- **Development package**: `libgdbm-dev` or `gdbm-devel` - Headers for building

```bash
# Debian/Ubuntu - Install both runtime and development
sudo apt-get install libgdbm-dev

# RHEL/CentOS - Install both runtime and development
sudo yum install gdbm-devel

# macOS - Homebrew includes both
brew install gdbm
```

### Verifying Installation

```bash
# Check if gdbm library is installed
ldconfig -p | grep gdbm  # Linux
ls /usr/local/lib/libgdbm*  # macOS

# Find library files
find /usr -name "libgdbm*" 2>/dev/null

# Check for header file
find /usr -name "gdbm.h" 2>/dev/null

# Test with Python
python3 -c "import dbm.gnu; print('gdbm works')"

# Test with Perl
perl -MGDBM_File -e 'print "gdbm works\n"'

# Test with Ruby
ruby -rgdbm -e 'puts "gdbm works"'
```

### Language Runtime Issues

**Python dbm.gnu not working**:

```bash
# Verify gdbm is available
python3 -c "import dbm.gnu; print('gdbm works')"

# Check if Python was built with gdbm support
python3 -c "import dbm; print(dbm.whichdb('test.db'))"

# Rebuild Python with gdbm if necessary
# Install gdbm-dev first, then rebuild Python
```

**Perl GDBM_File not found**:

```bash
# Install gdbm development package first
sudo apt-get install libgdbm-dev  # Debian/Ubuntu
brew install gdbm                 # macOS

# Reinstall Perl GDBM module
cpan GDBM_File
```

**Ruby GDBM gem installation fails**:

```bash
# Ensure gdbm library is available
brew install gdbm  # macOS
sudo apt-get install libgdbm-dev  # Debian/Ubuntu

# Install Ruby GDBM gem
gem install gdbm
```

### Database File Issues

**Database corruption**:

```bash
# Check database integrity
gdbmtool --check database.db

# Export to flat file (if recoverable)
gdbm_dump database.db > backup.txt

# Recreate from flat file
gdbm_load newdatabase.db < backup.txt
```

**Permission errors**:

```bash
# Check database file permissions
ls -l database.db

# Fix ownership
chown user:group database.db

# Fix permissions (typically 0600 or 0640)
chmod 0640 database.db
```

## Integration with Devgita

gdbm integrates with devgita's terminal category as a core dependency:

- **Installation**: Installed as part of core terminal tools setup before language runtimes
- **Configuration**: No configuration files - library is automatically available after installation
- **Usage**: Transparent to users - used internally by Perl, Python, Ruby, and system tools
- **Updates**: Managed through system package manager alongside OS updates
- **Dependencies**: Critical for language runtime installations - installed early in setup process

### Installation Order

Devgita installs gdbm early in the terminal tools setup:

1. **Core libraries**: pkg-config, autoconf, gdbm, libffi, openssl, readline, etc.
2. **Runtime managers**: mise (depends on proper library setup)
3. **Language runtimes**: Perl, Python, Ruby (may depend on gdbm)

This ensures language runtimes have all required dependencies when installed.

## External References

- **GDBM Official Documentation**: https://www.gnu.org.ua/software/gdbm/
- **GDBM Manual**: https://www.gnu.org.ua/software/gdbm/manual.html
- **GDBM on GNU**: https://www.gnu.org/software/gdbm/
- **Python dbm.gnu**: https://docs.python.org/3/library/dbm.html#module-dbm.gnu
- **Perl GDBM_File**: https://perldoc.perl.org/GDBM_File
- **Ruby GDBM**: https://ruby-doc.org/stdlib-3.0.0/libdoc/gdbm/rdoc/GDBM.html

This module provides essential key-value database library support for language runtimes, system tools, and persistent data storage within the devgita ecosystem.
