# Bat Module Documentation

## Overview

The Bat module provides syntax highlighting file viewer installation and command execution with devgita integration. It follows the standardized devgita app interface while providing bat-specific operations for file viewing, syntax highlighting, and Git integration capabilities.

## App Purpose

Bat is a cat clone with syntax highlighting and Git integration, written in Rust. It provides automatic paging, line numbers, and syntax highlighting for many programming languages and file formats. This module ensures bat is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems and provides high-level operations for enhanced file viewing in development workflows.

## Lifecycle Summary

1. **Installation**: Install bat package via platform package managers (Homebrew/apt)
2. **Configuration**: Bat uses optional user-specific configuration files (no default configuration applied by devgita)
3. **Execution**: Provide high-level bat operations for file viewing, syntax highlighting, and theme management

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Bat instance with platform-specific commands             |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install bat                               |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute bat commands      | Runs bat with provided arguments                                      |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
bat := bat.New()
err := bat.Install()
```

- **Purpose**: Standard bat installation
- **Behavior**: Uses `InstallPackage()` to install bat package
- **Use case**: Initial bat installation or explicit reinstall

### ForceInstall()

```go
bat := bat.New()
err := bat.ForceInstall()
```

- **Purpose**: Force bat installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh bat installation or fix corrupted installation

### SoftInstall()

```go
bat := bat.New()
err := bat.SoftInstall()
```

- **Purpose**: Install bat only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing bat installations

### Uninstall()

```go
err := bat.Uninstall()
```

- **Purpose**: Remove bat installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: File viewing tools are typically managed at the system level

### Update()

```go
err := bat.Update()
```

- **Purpose**: Update bat installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Bat updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := bat.ForceConfigure()
err := bat.SoftConfigure()
```

- **Purpose**: Apply bat configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: Bat configuration is optional and user-specific. Configuration files can be created at `~/.config/bat/config` and are managed by users based on their preferences.

## Execution Methods

### ExecuteCommand()

```go
err := bat.ExecuteCommand("file.go")                     // View file with syntax highlighting
err := bat.ExecuteCommand("--version")                   // Show version
err := bat.ExecuteCommand("--list-languages")            // Show supported languages
err := bat.ExecuteCommand("--list-themes")               // Show available themes
```

- **Purpose**: Execute bat commands with provided arguments
- **Parameters**: Variable arguments passed directly to bat binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### Bat-Specific Operations

The bat CLI provides extensive file viewing and syntax highlighting capabilities:

#### Basic File Viewing

```bash
# View file with syntax highlighting
bat file.go

# View multiple files
bat file1.go file2.py file3.js

# View file with line numbers (default)
bat file.go

# View without line numbers
bat --style plain file.go

# Show all characters including non-printable
bat -A file.go
```

#### Language and Syntax

```bash
# List all supported languages
bat --list-languages

# Force specific language syntax
bat -l python file.txt
bat --language rust file.txt

# Detect language from file extension
bat file.unknown --language python
```

#### Themes and Styling

```bash
# List all available themes
bat --list-themes

# Use specific theme
bat --theme "Monokai Extended" file.go
bat --theme Dracula file.py

# Preview all themes
bat --list-themes | fzf --preview="bat --theme={} --color=always file.go"

# Configure display style
bat --style full file.go              # Show all UI elements
bat --style header,grid,numbers file.go  # Custom combination
bat --style plain file.go             # Plain output (no decorations)
```

#### Git Integration

```bash
# Show Git modifications
bat file.go  # Automatically shows Git diff markers

# Disable Git integration
bat --no-paging --decorations never file.go

# Works with other Git tools
git show HEAD:file.go | bat -l go
git diff | bat -l diff
```

#### Paging and Output

```bash
# Force paging
bat --paging always file.go

# Disable paging
bat --paging never file.go

# Auto paging (default)
bat file.go

# Use custom pager
bat --pager "less -RF" file.go
```

#### Advanced Options

```bash
# Show only specific line range
bat --line-range 10:20 file.go
bat -r 10:20 file.go

# Wrap long lines
bat --wrap auto file.go
bat --wrap never file.go

# Show file name header
bat --decorations always file.go

# Show tabs as visible characters
bat --tabs 4 file.go
bat --show-all file.go
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **File Viewing**: `New()` → `SoftInstall()` → `ExecuteCommand()` with file path
4. **Version Check**: `New()` → `ExecuteCommand("--version")`

## Constants and Paths

### Relevant Constants

- **Package name**: Referenced via `constants.Bat` (typically "bat")
- Used by all installation methods for consistent package reference

### Configuration Approach

- **Optional configuration**: Bat uses optional configuration files
- **Default location**: `~/.config/bat/config` (user-managed)
- **No default config**: Devgita does not apply default configuration for bat
- **User customization**: Users can create their own configuration based on preferences

### Configuration Options

While devgita doesn't apply default configuration, users can customize bat via `~/.config/bat/config`:

```bash
# Example bat configuration
# Set the theme
--theme="Monokai Extended"

# Show line numbers, Git modifications, and file header
--style="numbers,changes,header"

# Use italic text on the terminal
--italic-text=always

# Add mouse scrolling support
--paging=always

# Use custom pager
--pager="less -FR"

# Map file extensions to languages
--map-syntax "*.conf:INI"
--map-syntax ".ignore:Git Ignore"
```

## Implementation Notes

- **CLI Tool Nature**: Unlike typical applications, bat is a command-line file viewer without complex configuration requirements
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since bat uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since bat uses optional user-specific configuration
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as bat updates should be handled by system package managers

## Usage Examples

### Basic Bat Operations

```go
bat := bat.New()

// Install bat
err := bat.SoftInstall()
if err != nil {
    return err
}

// View file with syntax highlighting
err = bat.ExecuteCommand("file.go")

// Check version
err = bat.ExecuteCommand("--version")
```

### Advanced Usage

```go
// List supported languages
err := bat.ExecuteCommand("--list-languages")

// Use specific theme
err = bat.ExecuteCommand("--theme", "Dracula", "file.py")

// Show line range
err = bat.ExecuteCommand("--line-range", "10:20", "file.go")

// Custom styling
err = bat.ExecuteCommand("--style", "header,grid,numbers", "file.rs")
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Syntax Not Highlighted**: Verify file extension is recognized or use `-l` flag
3. **Theme Not Applied**: Check theme name with `--list-themes`
4. **Paging Issues**: Adjust paging behavior with `--paging` flag
5. **Commands Don't Work**: Verify bat is installed and accessible in PATH

### Platform Considerations

- **macOS**: Installed via Homebrew package manager
- **Linux**: Installed via apt package manager (package may be named `batcat` on some systems)
- **Alias**: On some Linux systems, consider creating alias: `alias bat='batcat'`
- **Terminal Support**: Works best with modern terminal emulators that support colors and Unicode

### Prerequisites

Before using bat, ensure:
- Terminal supports colors and special characters
- For Git integration, bat should be run in a Git repository
- For optimal experience, use a terminal with good Unicode support

### Language Support

Bat supports syntax highlighting for many languages:
- **Programming**: C, C++, C#, Go, Java, JavaScript, TypeScript, Python, Ruby, Rust, Swift, PHP, etc.
- **Markup**: HTML, XML, Markdown, YAML, TOML, JSON
- **Shell**: Bash, Zsh, Fish, PowerShell
- **Config**: Nginx, Apache, Dockerfile, .gitignore, .env
- **Data**: CSV, SQL, GraphQL
- **And many more**: Use `bat --list-languages` for complete list

## Key Features

### Syntax Highlighting
- Automatic language detection based on file extension
- Manual language override with `-l` flag
- Support for 200+ languages and file formats
- Customizable themes

### Git Integration
- Shows Git modifications in sidebar
- Works seamlessly with Git commands
- Visual indicators for added/modified/deleted lines
- Can be disabled if not needed

### User-Friendly Output
- Line numbers by default
- File header with filename
- Grid separators for clarity
- Automatic paging for long files

### Customization
- Multiple color themes
- Configurable UI elements
- Custom pager support
- Line range selection

## External References

- **Bat Repository**: https://github.com/sharkdp/bat
- **Installation Guide**: https://github.com/sharkdp/bat#installation
- **Configuration**: https://github.com/sharkdp/bat#configuration-file
- **Customization**: https://github.com/sharkdp/bat#customization
- **Themes**: https://github.com/sharkdp/bat#adding-new-themes

## Integration with Devgita

Bat integrates with devgita's terminal category:

- **Installation**: Installed as part of terminal tools setup
- **Configuration**: User-managed configuration (no default applied)
- **Usage**: Available system-wide after installation
- **Updates**: Managed through system package manager
- **Replacement**: Designed as a drop-in replacement for `cat` command

This module provides essential file viewing and syntax highlighting capabilities for enhanced development workflows within the devgita ecosystem.
