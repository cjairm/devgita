# tldr Module Documentation

## Overview

The tldr module provides simplified command documentation tool installation and command execution with devgita integration. It follows the standardized devgita app interface while providing tldr-specific operations for displaying practical command examples and quick reference documentation.

## App Purpose

tldr (Too Long; Didn't Read) is a command-line utility that provides simplified and community-driven man pages with practical examples. Instead of comprehensive technical documentation, tldr focuses on the most common use cases with clear, concise examples. This module ensures tldr is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems and provides high-level operations for accessing practical command documentation.

## Lifecycle Summary

1. **Installation**: Install tldr package via platform package managers (Homebrew/apt)
2. **Configuration**: tldr typically doesn't require separate configuration files - operations are handled via command-line arguments or environment variables
3. **Execution**: Provide high-level tldr operations for displaying simplified command documentation with practical examples

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Tldr instance with platform-specific commands            |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install tldr                              |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute tldr commands     | Runs tldr with provided arguments                                    |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
tldr := tldr.New()
err := tldr.Install()
```

- **Purpose**: Standard tldr installation
- **Behavior**: Uses `InstallPackage()` to install tldr package
- **Use case**: Initial tldr installation or explicit reinstall

### ForceInstall()

```go
tldr := tldr.New()
err := tldr.ForceInstall()
```

- **Purpose**: Force tldr installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh tldr installation or fix corrupted installation

### SoftInstall()

```go
tldr := tldr.New()
err := tldr.SoftInstall()
```

- **Purpose**: Install tldr only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing tldr installations

### Uninstall()

```go
err := tldr.Uninstall()
```

- **Purpose**: Remove tldr installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Documentation tools are typically managed at the system level

### Update()

```go
err := tldr.Update()
```

- **Purpose**: Update tldr installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: tldr updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := tldr.ForceConfigure()
err := tldr.SoftConfigure()
```

- **Purpose**: Apply tldr configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: tldr doesn't use traditional config files; configuration is handled via environment variables or command-line arguments

## Execution Methods

### ExecuteCommand()

```go
err := tldr.ExecuteCommand("--version")
err := tldr.ExecuteCommand("git")
err := tldr.ExecuteCommand("--platform", "linux", "tar")
```

- **Purpose**: Execute tldr commands with provided arguments
- **Parameters**: Variable arguments passed directly to tldr binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### tldr-Specific Operations

The tldr CLI provides simplified command documentation capabilities:

#### Basic Usage

```bash
# Display tldr page for a command
tldr git

# Display page for specific command
tldr tar

# Show all available commands
tldr --list

# Update local cache of tldr pages
tldr --update
```

#### Platform-Specific Documentation

```bash
# Show examples for specific platform
tldr --platform linux tar
tldr --platform osx brew
tldr --platform windows dir

# Common platform values: linux, osx, windows, sunos
```

#### Language Support

```bash
# Display pages in specific language
tldr --language es git
tldr --language pt-BR docker

# Available languages include: en, es, pt-BR, fr, de, it, ja, ko, zh, and more
```

#### Advanced Options

```bash
# Render a specific markdown file
tldr --render /path/to/page.md

# Show version information
tldr --version

# Get help
tldr --help

# Clear local cache
tldr --clear-cache

# Show debug information
tldr --debug
```

#### Search and Navigation

```bash
# Search for commands
tldr --search "compress"

# Display random command
tldr --random

# Show recent updates
tldr --recent
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **Documentation Lookup**: `New()` → `SoftInstall()` → `ExecuteCommand()` with command name
4. **Update Cache**: `New()` → `ExecuteCommand("--update")`

## Constants and Paths

### Relevant Constants

- **Package name**: Referenced via `constants.Tldr` (value: "tldr")
- Used by all installation methods for consistent package reference

### Configuration Approach

- **No traditional config files**: tldr operations are configured via command-line arguments
- **Runtime configuration**: Parameters passed directly to `ExecuteCommand()`
- **Environment variables**: tldr respects several environment variables for customization:
  - `TLDR_CACHE_DIR`: Custom cache directory (default: `~/.cache/tldr` on Linux, `~/Library/Caches/tldr` on macOS)
  - `TLDR_PAGES_SOURCE_LOCATION`: Custom source for tldr pages
  - `TLDR_COLOR_BLANK`: Color for blank lines
  - `TLDR_COLOR_NAME`: Color for command names
  - `TLDR_COLOR_DESCRIPTION`: Color for descriptions
  - `TLDR_COLOR_EXAMPLE`: Color for example commands
  - `TLDR_COLOR_COMMAND`: Color for command syntax
  - `TLDR_COLOR_PARAMETER`: Color for parameters

## Implementation Notes

- **Documentation Tool Nature**: Unlike typical applications, tldr is a documentation lookup utility without persistent configuration
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since tldr uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since tldr doesn't use config files
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as tldr updates should be handled by system package managers
- **Cache Management**: tldr maintains a local cache of pages that can be updated with `--update` flag

## Usage Examples

### Basic Documentation Lookup

```go
tldr := tldr.New()

// Install tldr
err := tldr.SoftInstall()
if err != nil {
    return err
}

// Look up git documentation
err = tldr.ExecuteCommand("git")

// Look up tar documentation
err = tldr.ExecuteCommand("tar")

// Update local cache
err = tldr.ExecuteCommand("--update")
```

### Advanced Operations

```go
// Check tldr version
err := tldr.ExecuteCommand("--version")

// List all available commands
err = tldr.ExecuteCommand("--list")

// Show platform-specific documentation
err = tldr.ExecuteCommand("--platform", "linux", "sed")

// Show documentation in Spanish
err = tldr.ExecuteCommand("--language", "es", "docker")

// Search for commands
err = tldr.ExecuteCommand("--search", "compress")
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Pages Not Found**: Run `tldr --update` to download/update the local cache
3. **Outdated Information**: Update cache with `--update` flag
4. **Language Not Available**: Check available languages with `--help`
5. **Cache Issues**: Clear cache with `--clear-cache` and re-update

### Platform Considerations

- **macOS**: Installed via Homebrew package manager
- **Linux**: Installed via apt package manager, may also be available via pip/npm
- **Cache Location**: Platform-specific default cache directories
- **Network Required**: Initial setup and updates require internet connectivity

### Performance Tips

- **Update cache periodically**: Keep documentation current with `--update`
- **Use offline**: Once cached, tldr works offline
- **Clear old cache**: Periodically clear cache with `--clear-cache` to free space
- **Platform-specific pages**: Use `--platform` for OS-specific examples

### Best Practices

- **Update cache after installation**: Run `tldr --update` after first install
- **Check language availability**: Use `--list` to see available commands
- **Use specific platforms**: Specify `--platform` for cross-platform commands
- **Combine with man pages**: Use tldr for quick reference, man for comprehensive docs
- **Search functionality**: Use `--search` to find related commands

### Comparison with Traditional Man Pages

| Feature                | tldr                  | man pages                 |
| ---------------------- | --------------------- | ------------------------- |
| Focus                  | Practical examples    | Comprehensive docs        |
| Length                 | Short, concise        | Detailed, lengthy         |
| Learning Curve         | Beginner-friendly     | Requires experience       |
| Examples               | Many practical cases  | Few or technical          |
| Community-driven       | Yes                   | Official documentation    |
| Update Frequency       | Regular contributions | Slower, formal process    |
| Offline Support        | Yes (after cache)     | Yes                       |
| Platform-specific      | Explicit support      | Varies                    |

## Common Use Cases

### Quick Reference

```bash
# Quickly remember git clone syntax
tldr git-clone

# Check docker run options
tldr docker

# Review tar compression options
tldr tar
```

### Learning New Commands

```bash
# Learn about a new tool
tldr kubectl
tldr ansible
tldr terraform

# Explore related commands
tldr --search "network"
```

### Platform-Specific Help

```bash
# Get macOS-specific examples
tldr --platform osx launchctl

# Get Linux-specific examples
tldr --platform linux systemctl
```

## External References

- **tldr Repository**: https://github.com/tldr-pages/tldr
- **Client Specification**: https://github.com/tldr-pages/tldr/blob/main/CLIENT-SPECIFICATION.md
- **Contributing**: https://github.com/tldr-pages/tldr/blob/main/CONTRIBUTING.md
- **Web Interface**: https://tldr.sh/

## Integration with Devgita

tldr integrates with devgita's terminal category:

- **Installation**: Installed as part of terminal tools setup
- **Configuration**: Command-line based (no config files required)
- **Usage**: Available system-wide after installation
- **Updates**: Managed through system package manager
- **Complements**: Works alongside man pages for comprehensive documentation access

This module provides essential quick-reference documentation capabilities for efficient command lookup and learning within the devgita development environment, making it easier to remember command syntax and discover usage patterns.
