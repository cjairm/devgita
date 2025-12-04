# fd-find Module Documentation

## Overview

The fd-find module provides a fast and user-friendly file search tool installation and command execution with devgita integration. It follows the standardized devgita app interface while providing fd-specific operations for finding files and directories with intuitive syntax and powerful filtering capabilities.

## App Purpose

fd-find (commonly called 'fd') is a simple, fast, and user-friendly alternative to the traditional Unix 'find' command. It provides intuitive syntax, respects .gitignore by default, uses regular expressions for pattern matching, and offers colorized output. This module ensures fd-find is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems and provides high-level operations for file search and discovery tasks.

## Lifecycle Summary

1. **Installation**: Install fd-find package via platform package managers (Homebrew/apt)
2. **Configuration**: fd-find typically doesn't require separate configuration files - operations are handled via command-line arguments or environment variables
3. **Execution**: Provide high-level fd operations for file and directory search with pattern matching

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new FdFind instance with platform-specific commands          |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install fd-find                           |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute fd commands       | Runs fd-find with provided arguments                                 |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
fdFind := fdFind.New()
err := fdFind.Install()
```

- **Purpose**: Standard fd-find installation
- **Behavior**: Uses `InstallPackage()` to install fd-find package
- **Use case**: Initial fd-find installation or explicit reinstall

### ForceInstall()

```go
fdFind := fdFind.New()
err := fdFind.ForceInstall()
```

- **Purpose**: Force fd-find installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh fd-find installation or fix corrupted installation

### SoftInstall()

```go
fdFind := fdFind.New()
err := fdFind.SoftInstall()
```

- **Purpose**: Install fd-find only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing fd-find installations

### Uninstall()

```go
err := fdFind.Uninstall()
```

- **Purpose**: Remove fd-find installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: File search tools are typically managed at the system level

### Update()

```go
err := fdFind.Update()
```

- **Purpose**: Update fd-find installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: fd-find updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := fdFind.ForceConfigure()
err := fdFind.SoftConfigure()
```

- **Purpose**: Apply fd-find configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: fd-find doesn't use traditional config files; operation parameters are passed via command-line arguments or environment variables (FD_OPTS)

## Execution Methods

### ExecuteCommand()

```go
err := fdFind.ExecuteCommand("--version")
err := fdFind.ExecuteCommand("pattern")
err := fdFind.ExecuteCommand("-e", "go", "main")
```

- **Purpose**: Execute fd-find commands with provided arguments
- **Parameters**: Variable arguments passed directly to fd binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### fd-find-Specific Operations

The fd CLI provides extensive file search and filtering capabilities:

#### Basic Search

```bash
# Find files by name pattern
fd pattern

# Find files in specific directory
fd pattern /path/to/search

# Search for exact filename
fd -g exact-name.txt

# Search with regex
fd '^x.*rc$'
```

#### Type Filtering

```bash
# Search only files
fd -t f pattern

# Search only directories
fd -t d pattern

# Search symlinks
fd -t l pattern

# Search executables
fd -t x pattern
```

#### Extension Filtering

```bash
# Find files with specific extension
fd -e go

# Multiple extensions
fd -e go -e rs pattern

# Exclude extension
fd --exclude '*.tmp'
```

#### Hidden Files and Ignored Files

```bash
# Include hidden files
fd -H pattern

# Include .gitignore files
fd -I pattern

# Include both hidden and ignored
fd -HI pattern

# Show full path
fd -p pattern
```

#### Execution on Results

```bash
# Execute command on each result
fd -x echo

# Execute with placeholders
fd -x rm {}

# Execute command with multiple args
fd -e jpg -x convert {} {.}.png

# Execute in parallel
fd -x -j 4 command
```

#### Search Depth

```bash
# Limit search depth
fd -d 3 pattern

# Max depth
fd --max-depth 2 pattern

# Min depth (exclude top-level)
fd --min-depth 1 pattern
```

#### Output Control

```bash
# Colorized output (default)
fd pattern

# No color
fd --color never pattern

# Show absolute paths
fd -a pattern

# Show size
fd --size +1m pattern

# Show changed time
fd --changed-within 1day
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **File Search Operations**: `New()` → `SoftInstall()` → `ExecuteCommand()` with fd arguments
4. **Filtered Search**: `New()` → `ExecuteCommand()` with type/extension filters

## Constants and Paths

### Relevant Constants

- **Package name**: Referenced via `constants.FdFind` (value: "fd")
- **macOS**: Package is "fd" via Homebrew
- **Linux**: Package is "fd-find" via apt (handled via alias in SoftInstall)
- Used by all installation methods for consistent package reference

### Configuration Approach

- **No traditional config files**: fd-find operations are configured via command-line arguments
- **Runtime configuration**: Parameters passed directly to `ExecuteCommand()`
- **Environment variables**: fd respects `FD_OPTS` environment variable for default options
- **Optional .fdignore**: Users can create `~/.fdignore` for custom ignore patterns (similar to .gitignore)

## Implementation Notes

- **File Search Tool Nature**: Unlike typical applications, fd-find is a command-line search utility without persistent configuration
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since fd-find uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since fd-find doesn't use config files
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as fd-find updates should be handled by system package managers
- **Binary Name**: The binary is often installed as `fd` but the package name is `fd-find` to avoid conflicts

## Usage Examples

### Basic File Search Operations

```go
fdFind := fdFind.New()

// Install fd-find
err := fdFind.SoftInstall()
if err != nil {
    return err
}

// Find files by pattern
err = fdFind.ExecuteCommand("main")

// Find Go files
err = fdFind.ExecuteCommand("-e", "go")

// Find files in specific directory
err = fdFind.ExecuteCommand("pattern", "/path/to/search")
```

### Advanced Operations

```go
// Check fd-find version
err := fdFind.ExecuteCommand("--version")

// Find hidden files
err = fdFind.ExecuteCommand("-H", "config")

// Find directories only
err = fdFind.ExecuteCommand("-t", "d", "src")

// Execute command on results
err = fdFind.ExecuteCommand("-x", "echo", "test")

// Search with depth limit
err = fdFind.ExecuteCommand("-d", "3", "README")
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Binary Not Found**: On some systems, the binary might be named `fd` instead of `fd-find`
3. **Slow Performance**: Use depth limiting (`-d`) or exclude directories (`--exclude`)
4. **Too Many Results**: Use more specific patterns or type filters (`-t`)
5. **Missing Files**: Check if files are ignored by .gitignore (use `-I` flag)

### Platform Considerations

- **macOS**: Installed via Homebrew package manager as `fd`
- **Linux**: Installed via apt package manager as `fd-find`, binary might be `fd` or `fdfind`
- **Package Name**: On macOS the package is "fd", on Linux it's "fd-find" (automatically handled by devgita)
- **Binary Command**: The devgita module uses "fd" as the command name for cross-platform compatibility

### Performance Tips

- **Limit search depth**: Use `-d` flag to avoid searching entire filesystem
- **Use type filters**: `-t f` or `-t d` significantly speeds up searches
- **Exclude directories**: Use `--exclude` to skip unnecessary directories
- **Parallel execution**: Use `-j` flag for parallel command execution on results
- **Pattern specificity**: More specific patterns reduce search space

### Best Practices

- **Start with specific patterns**: Begin searches with narrow criteria
- **Use extension filters**: `-e` flag is faster than pattern matching
- **Leverage .gitignore**: Default behavior respects .gitignore, reducing search space
- **Combine with other tools**: Pipe results to `xargs`, `grep`, or use `-x` for execution
- **Check version compatibility**: Run `--version` to ensure supported features

### Comparison with Traditional find

| Feature                | fd-find        | traditional find              |
| ---------------------- | -------------- | ----------------------------- |
| Speed                  | Fast (Rust)    | Slower                        |
| Syntax                 | Simple         | Complex                       |
| .gitignore support     | Default        | Not available                 |
| Color output           | Default        | Requires configuration        |
| Regex support          | Built-in       | Limited                       |
| Hidden files           | Excluded (opt) | Included by default           |
| Performance            | Parallel       | Sequential                    |

## External References

- **fd Repository**: https://github.com/sharkdp/fd
- **User Guide**: https://github.com/sharkdp/fd#how-to-use
- **Installation**: https://github.com/sharkdp/fd#installation
- **Benchmarks**: https://github.com/sharkdp/fd#benchmark

## Integration with Devgita

fd-find integrates with devgita's terminal category:

- **Installation**: Installed as part of terminal tools setup
- **Configuration**: Command-line based (no config files required)
- **Usage**: Available system-wide after installation
- **Updates**: Managed through system package manager
- **Complements**: Works alongside other search tools like ripgrep (rg) and fzf

This module provides essential file search capabilities for efficient file discovery and navigation within the devgita development environment, significantly improving productivity compared to traditional find commands.
