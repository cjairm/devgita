# Ripgrep Module Documentation

## Overview

The Ripgrep module provides installation and command execution management for ripgrep with devgita integration. It follows the standardized devgita app interface while providing ripgrep-specific operations for fast recursive searching with regex pattern matching.

## App Purpose

Ripgrep (rg) is a line-oriented search tool that recursively searches the current directory for a regex pattern. It's significantly faster than grep and ack, respects gitignore rules by default, and provides colorized output. This module ensures ripgrep is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems and provides high-level operations for code searching and pattern matching workflows.

## Lifecycle Summary

1. **Installation**: Install ripgrep package via platform package managers (Homebrew/apt)
2. **Configuration**: Ripgrep uses optional configuration via environment variables and CLI flags (no default configuration applied by devgita)
3. **Execution**: Provide high-level ripgrep operations for searching files and directories

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Ripgrep instance with platform-specific commands         |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install ripgrep                           |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute ripgrep commands  | Runs rg with provided arguments                                       |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
ripgrep := ripgrep.New()
err := ripgrep.Install()
```

- **Purpose**: Standard ripgrep installation
- **Behavior**: Uses `InstallPackage()` to install ripgrep package
- **Use case**: Initial ripgrep installation or explicit reinstall

### ForceInstall()

```go
ripgrep := ripgrep.New()
err := ripgrep.ForceInstall()
```

- **Purpose**: Force ripgrep installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh ripgrep installation or fix corrupted installation

### SoftInstall()

```go
ripgrep := ripgrep.New()
err := ripgrep.SoftInstall()
```

- **Purpose**: Install ripgrep only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing ripgrep installations

### Uninstall()

```go
err := ripgrep.Uninstall()
```

- **Purpose**: Remove ripgrep installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Search tools are typically managed at the system level

### Update()

```go
err := ripgrep.Update()
```

- **Purpose**: Update ripgrep installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Ripgrep updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := ripgrep.ForceConfigure()
err := ripgrep.SoftConfigure()
```

- **Purpose**: Apply ripgrep configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: Ripgrep configuration is handled via CLI flags and environment variables rather than config files

## Execution Methods

### ExecuteCommand()

```go
err := ripgrep.ExecuteCommand("--version")
err := ripgrep.ExecuteCommand("pattern", "path/to/search")
err := ripgrep.ExecuteCommand("-i", "--type", "go", "TODO")
```

- **Purpose**: Execute ripgrep commands with provided arguments
- **Parameters**: Variable arguments passed directly to rg binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### Ripgrep-Specific Operations

The ripgrep CLI provides extensive search capabilities:

#### Basic Searching

```bash
# Search for pattern in current directory
rg pattern

# Search in specific file or directory
rg pattern path/to/search

# Case-insensitive search
rg -i pattern

# Search for whole words
rg -w pattern

# Show line numbers (default)
rg -n pattern
```

#### File Type Filtering

```bash
# Search only Go files
rg --type go pattern

# Search only specific extensions
rg -g "*.md" pattern

# Exclude certain file types
rg --type-not js pattern

# List available file types
rg --type-list
```

#### Output Control

```bash
# Show only matching files
rg -l pattern

# Show count of matches per file
rg -c pattern

# Show context lines
rg -C 3 pattern          # 3 lines before and after
rg -A 3 pattern          # 3 lines after
rg -B 3 pattern          # 3 lines before

# No line numbers
rg -N pattern

# Show column numbers
rg --column pattern
```

#### Advanced Search

```bash
# Regex search
rg "TODO:\s+\w+" .

# Fixed string search (no regex)
rg -F "function()" .

# Match multiline patterns
rg -U "start.*?end" .

# Invert match (show non-matching lines)
rg -v pattern

# Search hidden files and directories
rg --hidden pattern

# Search without respecting gitignore
rg --no-ignore pattern
```

#### Performance Options

```bash
# Parallel search with specific thread count
rg -j 4 pattern

# Search compressed files
rg -z pattern

# Memory-mapped search
rg --mmap pattern

# Follow symbolic links
rg -L pattern
```

#### Replacement

```bash
# Dry-run replacement preview
rg pattern --replace replacement

# Actually replace (requires external tool)
rg pattern -r replacement --passthru | sponge file.txt
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **Search Operations**: `New()` → `SoftInstall()` → `ExecuteCommand()` with search parameters
4. **Version Check**: `New()` → `ExecuteCommand("--version")`

## Constants and Paths

### Relevant Constants

- **Package name**: Referenced via `constants.Ripgrep` (typically "ripgrep")
- **Binary name**: Referenced via `constants.Rg` (typically "rg")
- Used by all installation and execution methods for consistent package reference

### Configuration Approach

- **No traditional config files**: Ripgrep configuration is handled via CLI flags
- **Environment variables**: Supports `RIPGREP_CONFIG_PATH` for config file location
- **Optional config file**: Users can create `~/.ripgreprc` for default flags
- **No default config**: Devgita does not apply default configuration for ripgrep

### Configuration Options

While devgita doesn't apply default configuration, users can customize ripgrep via `~/.ripgreprc`:

```bash
# Example ripgrep configuration file
# One flag per line

# Case-insensitive by default
--smart-case

# Show line numbers
--line-number

# Show column numbers
--column

# Add colors
--colors=line:fg:yellow
--colors=line:style:bold
--colors=path:fg:green
--colors=path:style:bold
--colors=match:fg:black
--colors=match:bg:yellow
--colors=match:style:nobold

# Follow symbolic links
--follow

# Search hidden files
--hidden

# Ignore patterns
--glob=!.git/
--glob=!node_modules/
--glob=!vendor/
--glob=!*.min.js
--glob=!*.lock
```

## Implementation Notes

- **CLI Tool Nature**: Unlike typical applications, ripgrep is a search tool without complex configuration requirements
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since ripgrep uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since ripgrep uses CLI flags and optional config files
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as ripgrep updates should be handled by system package managers

## Usage Examples

### Basic Ripgrep Operations

```go
ripgrep := ripgrep.New()

// Install ripgrep
err := ripgrep.SoftInstall()
if err != nil {
    return err
}

// Search for pattern in current directory
err = ripgrep.ExecuteCommand("TODO", ".")

// Case-insensitive search
err = ripgrep.ExecuteCommand("-i", "fixme")

// Search specific file types
err = ripgrep.ExecuteCommand("--type", "go", "func main")
```

### Advanced Operations

```go
// Check version
err := ripgrep.ExecuteCommand("--version")

// Search with context
err = ripgrep.ExecuteCommand("-C", "3", "pattern", "path/")

// Search hidden files
err = ripgrep.ExecuteCommand("--hidden", "secret")

// Count matches per file
err = ripgrep.ExecuteCommand("-c", "import")

// List matching files only
err = ripgrep.ExecuteCommand("-l", "TODO")
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **No Results Found**: Check if gitignore is hiding files (use `--no-ignore`)
3. **Permission Errors**: Some directories may require elevated permissions
4. **Regex Syntax**: Use `rg --help` to verify regex syntax
5. **Commands Don't Work**: Verify ripgrep is installed and accessible in PATH

### Platform Considerations

- **macOS**: Installed via Homebrew package manager
- **Linux**: Installed via apt package manager
- **Binary Name**: Command is `rg`, not `ripgrep`
- **Performance**: Extremely fast due to Rust implementation and parallel searching

### Prerequisites

Before using ripgrep effectively:
- Understand basic regex patterns
- Familiarize with gitignore behavior (ripgrep respects it by default)
- Learn file type filters for efficient searching
- Consider creating `~/.ripgreprc` for personal defaults

### Performance Tips

- Use `--type` to limit search scope
- Use `--glob` for fine-grained file selection
- Leverage parallel search (enabled by default)
- Use `-l` to quickly find files without showing matches
- Use `--mmap` for very large files

### Comparison with Other Tools

**Ripgrep vs grep:**
- Ripgrep is significantly faster (5-10x on average)
- Respects gitignore by default
- Better defaults (recursive, colored output, line numbers)
- Superior regex engine

**Ripgrep vs ack:**
- Ripgrep is faster (2-5x)
- Better memory usage
- More file type filters
- Active development

**Ripgrep vs ag (The Silver Searcher):**
- Ripgrep is faster (2-3x)
- More accurate regex matching
- Better gitignore handling
- More configuration options

## Key Features

### Speed and Performance
- Written in Rust for maximum performance
- Parallel searching across multiple cores
- Memory-mapped file access
- Optimized regex engine

### Smart Defaults
- Respects `.gitignore` automatically
- Recursive search by default
- Colored output for readability
- Line numbers shown by default

### File Filtering
- 100+ built-in file type definitions
- Custom glob patterns
- Exclude patterns
- Hidden file control

### Output Control
- Flexible output formatting
- Context lines (before/after)
- Match counts
- File-only listing
- JSON output for scripting

## Integration with Devgita

Ripgrep integrates with devgita's terminal category:

- **Installation**: Installed as part of terminal tools setup
- **Configuration**: User-managed via CLI flags or `~/.ripgreprc`
- **Usage**: Available system-wide after installation
- **Updates**: Managed through system package manager
- **Common Use**: Code searching, log analysis, refactoring support

## External References

- **Ripgrep Repository**: https://github.com/BurntSushi/ripgrep
- **User Guide**: https://github.com/BurntSushi/ripgrep/blob/master/GUIDE.md
- **FAQ**: https://github.com/BurntSushi/ripgrep/blob/master/FAQ.md
- **Regex Syntax**: https://docs.rs/regex/latest/regex/#syntax
- **Performance Comparison**: https://blog.burntsushi.net/ripgrep/
