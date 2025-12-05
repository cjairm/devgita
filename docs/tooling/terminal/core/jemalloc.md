# Jemalloc Module Documentation

## Overview

The Jemalloc module provides installation and library management for jemalloc (memory allocator) with devgita integration. It follows the standardized devgita app interface while providing jemalloc-specific operations for enabling efficient memory management, supporting high-performance applications, and reducing memory fragmentation in production systems.

## App Purpose

jemalloc is a general-purpose malloc(3) implementation that emphasizes fragmentation avoidance and scalable concurrency support. It provides many introspection, memory management, and tuning features beyond the standard allocator. jemalloc is widely used in production systems, particularly by high-performance databases (Redis, MariaDB, Cassandra) and applications requiring efficient memory management. This module ensures jemalloc is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems as a critical performance dependency for database systems and memory-intensive applications.

## Lifecycle Summary

1. **Installation**: Install jemalloc package via platform package managers (Homebrew/apt)
2. **Configuration**: jemalloc is a library and typically doesn't require separate configuration files - tuning is handled via environment variables or compile-time options
3. **Execution**: Provides ExecuteCommand() for interface compliance, but jemalloc is primarily used as a library replacement for system malloc

## Exported Functions

| Function           | Purpose                   | Behavior                                                                |
| ------------------ | ------------------------- | ----------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Jemalloc instance with platform-specific commands           |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install jemalloc                             |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()`    |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing                 |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                        |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                        |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                       |
| `ExecuteCommand()` | Execute jemalloc commands | Runs jemalloc utilities with provided arguments (limited practical use) |
| `Update()`         | Update installation       | **Not implemented** - returns error                                     |

## Installation Methods

### Install()

```go
jemalloc := jemalloc.New()
err := jemalloc.Install()
```

- **Purpose**: Standard jemalloc installation
- **Behavior**: Uses `InstallPackage()` to install jemalloc package
- **Use case**: Initial jemalloc installation or explicit reinstall for database performance

### ForceInstall()

```go
jemalloc := jemalloc.New()
err := jemalloc.ForceInstall()
```

- **Purpose**: Force jemalloc installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh jemalloc installation or fix corrupted library installation

### SoftInstall()

```go
jemalloc := jemalloc.New()
err := jemalloc.SoftInstall()
```

- **Purpose**: Install jemalloc only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing jemalloc installations

### Uninstall()

```go
err := jemalloc.Uninstall()
```

- **Purpose**: Remove jemalloc installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: jemalloc is a critical performance library for database systems and applications

### Update()

```go
err := jemalloc.Update()
```

- **Purpose**: Update jemalloc installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: jemalloc updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := jemalloc.ForceConfigure()
err := jemalloc.SoftConfigure()
```

- **Purpose**: Apply jemalloc configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: jemalloc is a library without traditional config files; tuning is handled via environment variables (MALLOC_CONF) or compile-time options

## Execution Methods

### ExecuteCommand()

```go
err := jemalloc.ExecuteCommand("--version")
err := jemalloc.ExecuteCommand("--config", "stats_print:true")
```

- **Purpose**: Execute jemalloc commands with provided arguments
- **Parameters**: Variable arguments passed directly to jemalloc utilities (if available)
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand
- **Note**: jemalloc is primarily a library, not a CLI tool. Utilities like jemalloc-config may be available on some systems.

### Jemalloc-Specific Operations

jemalloc is primarily used as a memory allocator library rather than a command-line tool:

#### Using jemalloc as Allocator

**Via LD_PRELOAD (Linux)**:

```bash
# Use jemalloc for a specific application
LD_PRELOAD=/usr/lib/x86_64-linux-gnu/libjemalloc.so.2 redis-server

# Set permanently via environment
export LD_PRELOAD=/usr/lib/libjemalloc.so.2
./myapp
```

**Via DYLD_INSERT_LIBRARIES (macOS)**:

```bash
# Use jemalloc for a specific application
DYLD_INSERT_LIBRARIES=/usr/local/lib/libjemalloc.dylib redis-server
```

**Direct Linking**:

```bash
# Compile with jemalloc
gcc -o myapp myapp.c -ljemalloc

# CMake integration
target_link_libraries(myapp jemalloc)
```

#### Environment Variable Configuration

jemalloc is tuned via the MALLOC_CONF environment variable:

```bash
# Enable statistics printing on exit
export MALLOC_CONF="stats_print:true"

# Set background thread for asynchronous operations
export MALLOC_CONF="background_thread:true"

# Configure dirty page decay time (milliseconds)
export MALLOC_CONF="dirty_decay_ms:5000,muzzy_decay_ms:10000"

# Enable memory profiling
export MALLOC_CONF="prof:true,prof_leak:true,prof_final:true"

# Set arena count (default: number of CPUs)
export MALLOC_CONF="narenas:4"

# Combine multiple options
export MALLOC_CONF="stats_print:true,background_thread:true,dirty_decay_ms:5000"
```

#### jemalloc-config Utility

Some installations provide jemalloc-config for build integration:

```bash
# Show jemalloc configuration
jemalloc-config --help

# Get compiler flags
jemalloc-config --cflags

# Get linker flags
jemalloc-config --libs

# Get jemalloc version
jemalloc-config --version

# Get library directory
jemalloc-config --libdir
```

#### Memory Profiling

jemalloc provides powerful memory profiling capabilities:

```bash
# Enable profiling
export MALLOC_CONF="prof:true,prof_prefix:jeprof.out"

# Run application
./myapp

# Analyze profile with jeprof tool
jeprof --show_bytes ./myapp jeprof.out.*

# Generate heap profile graph
jeprof --pdf ./myapp jeprof.out.* > profile.pdf

# Show top memory consumers
jeprof --text ./myapp jeprof.out.*
```

#### Database Integration Examples

**Redis with jemalloc**:

```bash
# Redis typically uses jemalloc by default
# Verify with redis-server --version
redis-server --version | grep jemalloc

# Force jemalloc via LD_PRELOAD if needed
LD_PRELOAD=/usr/lib/libjemalloc.so.2 redis-server /etc/redis/redis.conf
```

**MariaDB with jemalloc**:

```bash
# MariaDB can use jemalloc via LD_PRELOAD
# Add to systemd service or init script
Environment="LD_PRELOAD=/usr/lib/x86_64-linux-gnu/libjemalloc.so.2"

# Or set in my.cnf
[mysqld]
malloc-lib=/usr/lib/x86_64-linux-gnu/libjemalloc.so.2
```

**Cassandra with jemalloc**:

```bash
# Set in cassandra-env.sh
export LD_PRELOAD=/usr/lib/libjemalloc.so.2

# Verify jemalloc is loaded
lsof -p $(pgrep -f cassandra) | grep jemalloc
```

#### Runtime Statistics

Access jemalloc statistics at runtime:

```c
// In C/C++ applications
#include <jemalloc/jemalloc.h>

// Print statistics
malloc_stats_print(NULL, NULL, NULL);

// Get specific stats
size_t allocated, active, metadata;
size_t sz = sizeof(size_t);
mallctl("stats.allocated", &allocated, &sz, NULL, 0);
mallctl("stats.active", &active, &sz, NULL, 0);
mallctl("stats.metadata", &metadata, &sz, NULL, 0);
```

#### Build System Integration

**CMake Detection**:

```cmake
find_package(PkgConfig)
pkg_check_modules(JEMALLOC jemalloc)

if(JEMALLOC_FOUND)
    include_directories(${JEMALLOC_INCLUDE_DIRS})
    link_directories(${JEMALLOC_LIBRARY_DIRS})
    target_link_libraries(myapp ${JEMALLOC_LIBRARIES})
endif()
```

**Makefile Integration**:

```makefile
JEMALLOC_CFLAGS := $(shell jemalloc-config --cflags)
JEMALLOC_LIBS := $(shell jemalloc-config --libs)

CFLAGS += $(JEMALLOC_CFLAGS)
LDFLAGS += $(JEMALLOC_LIBS)
```

**Autotools Integration**:

```bash
# configure.ac
AC_CHECK_LIB([jemalloc], [malloc], [], [
  AC_MSG_ERROR([jemalloc library not found])
])
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **Database Dependency**: `New()` → `SoftInstall()` before installing Redis/MariaDB/Cassandra
4. **Performance Optimization**: Install via package manager, then use via LD_PRELOAD or direct linking

## Constants and Paths

### Relevant Constants

- **Package name**: Referenced via `constants.Jemalloc` (typically "jemalloc")
- Used by all installation methods for consistent package reference

### Configuration Approach

- **No traditional config files**: jemalloc is tuned via environment variables
- **MALLOC_CONF**: Primary configuration mechanism via environment variable
- **Compile-time options**: Advanced tuning via build configuration
- **Runtime API**: mallctl() interface for programmatic configuration
- **Library paths**: Installed to system library directories

### Library and Header Paths

- **Shared libraries**:
  - Linux: `/usr/lib/x86_64-linux-gnu/libjemalloc.so.2`
  - macOS: `/usr/local/lib/libjemalloc.dylib`
- **Header files**: `/usr/include/jemalloc/jemalloc.h`
- **pkg-config data**: `/usr/lib/pkgconfig/jemalloc.pc`
- **Utilities**: `jemalloc-config`, `jeprof` (profiling tool)

## Implementation Notes

- **Library Nature**: Unlike typical applications, jemalloc is a memory allocator library without standalone CLI functionality
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since jemalloc uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since jemalloc uses environment variables
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as jemalloc updates should be handled by system package managers
- **Critical Dependency**: Required by many high-performance databases and applications

## Usage Examples

### Installation as Dependency

```go
jemalloc := jemalloc.New()

// Install jemalloc before database systems
err := jemalloc.SoftInstall()
if err != nil {
    return err
}

// Now safe to install Redis, MariaDB, etc.
// These databases can use jemalloc for better performance
```

### Integration with Database Installation

```go
// In database installer coordinator
func InstallRedis(ctx context.Context) error {
    // Ensure jemalloc is available first
    jemalloc := jemalloc.New()
    if err := jemalloc.SoftInstall(); err != nil {
        return fmt.Errorf("failed to install jemalloc dependency: %w", err)
    }

    // Now install Redis which benefits from jemalloc
    redis := redis.New()
    return redis.SoftInstall()
}
```

### Using jemalloc in Applications

```go
// Example: Running Redis with jemalloc via LD_PRELOAD
func StartRedisWithJemalloc() error {
    cmd := exec.Command("redis-server", "/etc/redis/redis.conf")
    cmd.Env = append(os.Environ(),
        "LD_PRELOAD=/usr/lib/x86_64-linux-gnu/libjemalloc.so.2")
    return cmd.Start()
}
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **LD_PRELOAD Not Working**: Verify library path is correct for your system
3. **Application Crashes**: Check jemalloc compatibility with application version
4. **Performance Issues**: Tune MALLOC_CONF settings for your workload
5. **Profiling Not Working**: Ensure profiling was enabled at compile time

### Platform Considerations

- **macOS**: Installed via Homebrew as `jemalloc`, includes headers and tools
- **Linux (Debian/Ubuntu)**: Install both `libjemalloc2` (runtime) and `libjemalloc-dev` (development)
- **Linux (RHEL/CentOS)**: Install both `jemalloc` (runtime) and `jemalloc-devel` (development)
- **Dependencies**: Required by high-performance databases - install before database systems

### Development vs Runtime Packages

On Linux systems, jemalloc typically requires both packages:

- **Runtime package**: `libjemalloc2` or `jemalloc` - Shared library for applications
- **Development package**: `libjemalloc-dev` or `jemalloc-devel` - Headers and build tools

```bash
# Debian/Ubuntu - Install both runtime and development
sudo apt-get install libjemalloc-dev

# RHEL/CentOS - Install both runtime and development
sudo yum install jemalloc-devel

# macOS - Homebrew includes both
brew install jemalloc
```

### Verifying Installation

```bash
# Check if jemalloc library is installed
ldconfig -p | grep jemalloc  # Linux
ls /usr/local/lib/libjemalloc*  # macOS

# Find library files
find /usr -name "libjemalloc*" 2>/dev/null

# Check for header file
find /usr -name "jemalloc.h" 2>/dev/null

# Test with LD_PRELOAD
LD_PRELOAD=/usr/lib/libjemalloc.so.2 ls -la

# Check if application uses jemalloc
lsof -p $(pgrep redis-server) | grep jemalloc

# Verify jemalloc-config
jemalloc-config --version
```

### Database Integration Issues

**Redis not using jemalloc**:

```bash
# Check Redis build configuration
redis-server --version | grep jemalloc

# Force jemalloc via LD_PRELOAD
LD_PRELOAD=/usr/lib/libjemalloc.so.2 redis-server

# Add to systemd service
# /etc/systemd/system/redis.service
[Service]
Environment="LD_PRELOAD=/usr/lib/x86_64-linux-gnu/libjemalloc.so.2"
```

**MariaDB/MySQL not using jemalloc**:

```bash
# Check if mysqld uses jemalloc
lsof -p $(pgrep mysqld) | grep jemalloc

# Set in systemd service
# /etc/systemd/system/mariadb.service.d/override.conf
[Service]
Environment="LD_PRELOAD=/usr/lib/libjemalloc.so.2"

# Or use malloc-lib in my.cnf
[mysqld]
malloc-lib=/usr/lib/x86_64-linux-gnu/libjemalloc.so.2
```

**Application crashes with jemalloc**:

```bash
# Check for version compatibility
ldd /path/to/app | grep jemalloc

# Try different jemalloc configuration
export MALLOC_CONF="background_thread:false"

# Check for conflicting allocators
# Remove other allocator preloads
unset LD_PRELOAD

# Test without jemalloc first
./myapp
```

### Performance Tuning

**High fragmentation**:

```bash
# Reduce decay times for more aggressive memory reclamation
export MALLOC_CONF="dirty_decay_ms:1000,muzzy_decay_ms:1000"
```

**High memory usage**:

```bash
# Enable background threads for asynchronous operations
export MALLOC_CONF="background_thread:true"

# Reduce arena count
export MALLOC_CONF="narenas:2"
```

**Memory leaks**:

```bash
# Enable leak detection profiling
export MALLOC_CONF="prof:true,prof_leak:true,prof_final:true,lg_prof_sample:0"

# Run application
./myapp

# Analyze with jeprof
jeprof --text ./myapp jeprof.out.*
```

## Integration with Devgita

jemalloc integrates with devgita's terminal category as a core dependency:

- **Installation**: Installed as part of core terminal tools setup before database systems
- **Configuration**: No configuration files - tuning via MALLOC_CONF environment variable
- **Usage**: Transparent to users - databases and applications use it automatically
- **Updates**: Managed through system package manager alongside OS updates
- **Dependencies**: Critical for database performance - installed early in setup process

### Installation Order

Devgita installs jemalloc early in the terminal tools setup:

1. **Core libraries**: pkg-config, autoconf, libffi, jemalloc, openssl, etc.
2. **Database systems**: Redis, MariaDB, Cassandra (benefit from jemalloc)
3. **Runtime managers**: mise (general development tools)

This ensures database systems have optimal memory management when installed.

## External References

- **Jemalloc Official Site**: http://jemalloc.net/
- **Jemalloc GitHub**: https://github.com/jemalloc/jemalloc
- **Jemalloc Documentation**: https://jemalloc.net/jemalloc.3.html
- **Redis and jemalloc**: https://redis.io/docs/reference/optimization/memory-optimization/
- **MariaDB and jemalloc**: https://mariadb.com/kb/en/using-jemalloc/
- **Memory Profiling Guide**: https://github.com/jemalloc/jemalloc/wiki/Use-Case:-Heap-Profiling

This module provides essential high-performance memory allocation support for database systems and memory-intensive applications within the devgita ecosystem.
