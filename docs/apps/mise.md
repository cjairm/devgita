# Mise Module Documentation

## Overview

The Mise module provides runtime environment management with devgita integration. It follows the standardized devgita app interface while providing mise-specific operations for managing programming language runtimes including Node.js, Python, Go, Ruby, and many others via the mise.jdx.dev tool.

## App Purpose

Mise (formerly rtx) is a polyglot runtime manager that replaces tools like nvm, pyenv, rbenv, etc. This module ensures mise is properly installed and provides high-level operations for setting global runtime versions, managing language environments, and handling development tool chains across different programming languages.

## Lifecycle Summary

1. **Installation**: Install mise package via platform package managers (Homebrew/apt)
2. **Configuration**: No traditional config files - mise manages runtimes via CLI commands
3. **Execution**: Provide high-level mise operations for runtime management and version control

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Mise instance with platform-specific commands            |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install mise                              |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns error                                   |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns error                                   |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute mise command      | Runs mise with provided arguments                                     |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |
| `UseGlobal()`      | Set global runtime        | Sets global version for specified language                           |

## Installation Methods

### Install()

```go
mise := mise.New()
err := mise.Install()
```

- **Purpose**: Standard mise installation
- **Behavior**: Uses `InstallPackage()` to install mise package
- **Use case**: Initial mise installation or explicit reinstall

### ForceInstall()

```go
mise := mise.New()
err := mise.ForceInstall()
```

- **Purpose**: Force mise installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh mise installation or fix corrupted installation

### SoftInstall()

```go
mise := mise.New()
err := mise.SoftInstall()
```

- **Purpose**: Install mise only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing mise installations

### Uninstall()

```go
err := mise.Uninstall()
```

- **Purpose**: Remove mise installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Runtime managers are typically managed at the system level

### Update()

```go
err := mise.Update()
```

- **Purpose**: Update mise installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Mise updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := mise.ForceConfigure()
err := mise.SoftConfigure()
```

- **Purpose**: Apply mise configuration
- **Behavior**: **Not applicable** - both return errors
- **Rationale**: Mise doesn't use traditional config files; runtime management is handled via CLI commands and global settings

## Execution Methods

### ExecuteCommand()

```go
err := mise.ExecuteCommand("--version")
err := mise.ExecuteCommand("install", "node@20")
err := mise.ExecuteCommand("global", "python@3.11")
```

- **Purpose**: Execute mise commands with provided arguments
- **Parameters**: Variable arguments passed directly to mise binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### Mise-Specific Operations

#### UseGlobal()

```go
err := mise.UseGlobal("node", "20")
err := mise.UseGlobal("python", "3.11.0")
err := mise.UseGlobal("go", "latest")
```

- **Purpose**: Set global runtime version for specified language
- **Parameters**: 
  - `language string` - Programming language name (e.g., "node", "python", "go")
  - `version string` - Version to set globally (e.g., "20", "3.11.0", "latest")
- **Behavior**: Executes `mise global <language>@<version>`
- **Validation**: Returns error if language or version parameters are empty
- **Use case**: Configure default runtime versions for development environment

#### Runtime Management Commands

The mise CLI provides extensive runtime management capabilities:

##### Language Installation

```bash
# Install specific versions
mise install node@20
mise install python@3.11.0
mise install go@latest

# Install from .tool-versions file
mise install

# List available versions
mise list-all node
mise list-all python
```

##### Version Management

```bash
# Set global versions
mise global node@20
mise global python@3.11
mise global go@latest

# Set local versions (project-specific)
mise local node@18
mise local python@3.10

# Show current versions
mise current
mise current node
```

##### Environment Operations

```bash
# Show installed versions
mise list
mise list node

# Update all tools
mise update

# Remove unused versions
mise prune
mise prune node

# Show tool information
mise info node
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (returns error)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (returns error)
3. **Runtime Management**: `New()` → `SoftInstall()` → `UseGlobal("node", "20")`
4. **Mise Operations**: `New()` → `ExecuteCommand()` with specific mise arguments

## Constants and Paths

### Relevant Constants

- Package name: `"mise"` used directly for installation
- Used by all installation methods for consistent package reference

### Configuration Approach

- **No traditional config files**: Mise manages runtimes via CLI commands and global settings
- **Global versions**: Stored in `~/.config/mise/config.toml` (managed by mise itself)
- **Local versions**: Stored in `.tool-versions` or `.mise.toml` files (project-specific)
- **Devgita integration**: Uses `UseGlobal()` to set development environment defaults

## Implementation Notes

- **Runtime Manager Nature**: Unlike typical applications, mise manages other tools rather than being a standalone application
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since mise uninstall is not supported
- **Configuration Strategy**: Returns errors for `ForceConfigure()` and `SoftConfigure()` since mise doesn't use traditional config files
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Validation**: `UseGlobal()` validates that both language and version parameters are provided
- **Update Method**: Not implemented as mise updates should be handled by system package managers

## Supported Languages

Mise supports a wide variety of programming languages and tools:

### Core Languages

- **Node.js**: `node`, version management via mise
- **Python**: `python`, multiple versions supported
- **Go**: `golang`, latest and specific versions
- **Ruby**: `ruby`, rbenv compatibility
- **Rust**: `rust`, cargo integration
- **Java**: `java`, multiple JDK versions
- **PHP**: `php`, composer integration

### Additional Tools

- **Terraform**: Infrastructure as code
- **Kubectl**: Kubernetes CLI
- **Helm**: Kubernetes package manager
- **And many more**: 500+ tools supported

## Usage Examples

### Setting Up Development Environment

```go
mise := mise.New()

// Install mise
err := mise.SoftInstall()
if err != nil {
    return err
}

// Set global runtimes for development
err = mise.UseGlobal("node", "20")
if err != nil {
    return err
}

err = mise.UseGlobal("python", "3.11")
if err != nil {
    return err
}

err = mise.UseGlobal("go", "latest")
if err != nil {
    return err
}
```

### Direct Command Execution

```go
// Install specific runtime
err := mise.ExecuteCommand("install", "node@18")

// List installed versions
err = mise.ExecuteCommand("list")

// Update all tools
err = mise.ExecuteCommand("update")
```

## Deprecated Functions

The module maintains backward compatibility through deprecated functions:

- `MaybeInstall()` → Use `SoftInstall()` instead
- `Setup()` → Use `ForceConfigure()` instead (note: will return error)
- `MaybeSetup()` → Use `SoftConfigure()` instead (note: will return error)
- `Run()` → Use `ExecuteCommand()` instead

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Runtime Installation Fails**: Check internet connectivity and tool availability
3. **Commands Don't Work**: Verify mise is installed and accessible in PATH
4. **Version Conflicts**: Use `mise current` to check active versions
5. **Shell Integration**: Ensure mise is properly configured in shell profile

### Platform Considerations

- **macOS**: Installed via Homebrew package manager
- **Linux**: Installed via apt package manager
- **Shell Integration**: Requires shell hook for automatic version switching
- **PATH Management**: Mise manages runtime paths automatically

### Configuration Files

While mise doesn't use traditional config templates like other devgita apps, it does create and manage:

- **Global config**: `~/.config/mise/config.toml`
- **Local config**: `.tool-versions` or `.mise.toml` in project directories
- **Shell integration**: Automatic PATH management and version switching

This module provides essential runtime management capabilities for maintaining consistent development environments across different programming languages and projects.