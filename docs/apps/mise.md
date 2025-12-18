# Mise Module Documentation

## Overview

The Mise module provides runtime environment management with devgita integration. It follows the standardized devgita app interface while providing mise-specific operations for managing programming language runtimes including Node.js, Python, Go, Ruby, and many others via the mise.jdx.dev tool.

**Note**: This module uses devgita's template-based configuration management system for shell integration. Configuration changes are tracked in GlobalConfig and shell configuration is regenerated from templates rather than direct file manipulation.

## App Purpose

Mise (formerly rtx) is a polyglot runtime manager that replaces tools like nvm, pyenv, rbenv, etc. This module ensures mise is properly installed, provides shell integration via template-based GlobalConfig management, and offers high-level operations for setting global runtime versions, managing language environments, and handling development tool chains across different programming languages.

## Lifecycle Summary

1. **Installation**: Install mise package via platform package managers (Homebrew/apt)
2. **Configuration**: Enable mise shell integration in GlobalConfig and regenerate shell configuration via templates
3. **Execution**: Provide high-level mise operations for runtime management and version control

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Mise instance with platform-specific commands            |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install mise                              |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | Enables shell integration and regenerates shell configuration        |
| `SoftConfigure()`  | Conditional configuration | Checks GlobalConfig; enables only if not already enabled             |
| `Uninstall()`      | Remove installation       | **Fully supported** - Disables shell integration and regenerates config |
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

- **Purpose**: Remove mise shell integration
- **Behavior**: Loads GlobalConfig → Disables feature → Regenerates shell config → Persists updated state
- **Use case**: Remove mise shell integration from devgita without uninstalling the package
- **Note**: This disables the shell integration feature but does not remove the package itself

### Update()

```go
err := mise.Update()
```

- **Purpose**: Update mise installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Mise updates are typically handled by the system package manager

## Configuration Methods

### Configuration Strategy

This module uses **template-based configuration management** via GlobalConfig:

- **Template**: `configs/templates/devgita.zsh.tmpl` contains conditional sections
- **GlobalConfig**: `~/.config/devgita/global_config.yaml` tracks enabled features
- **Generated file**: `~/.config/devgita/devgita.zsh` (regenerated from template)
- **Feature tracking**: `shell.mise` boolean field in GlobalConfig

### ForceConfigure()

```go
err := mise.ForceConfigure()
```

- **Purpose**: Enable mise shell integration in shell configuration
- **Behavior**: 
  1. Loads GlobalConfig from disk
  2. Enables `shell.mise` feature
  3. Regenerates `devgita.zsh` from template
  4. Saves GlobalConfig back to disk
- **Use case**: Enable mise shell integration or re-apply configuration

### SoftConfigure()

```go
err := mise.SoftConfigure()
```

- **Purpose**: Enable mise shell integration only if not already enabled
- **Behavior**: 
  1. Loads GlobalConfig from disk
  2. Checks if `shell.mise` is already enabled
  3. If enabled, returns nil (no operation)
  4. If not enabled, calls `ForceConfigure()`
- **Use case**: Initial setup that preserves existing configuration state

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
- **Behavior**: Executes `mise use --global <language>@<version>`
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

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()`
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()`
3. **Runtime Management**: `New()` → `SoftInstall()` → `SoftConfigure()` → `UseGlobal("node", "20")`
4. **Mise Operations**: `New()` → `ExecuteCommand()` with specific mise arguments
5. **Remove Integration**: `New()` → `Uninstall()`

## Constants and Paths

### Relevant Constants

- `constants.Mise`: Package name ("mise") used for installation and feature tracking
- Used consistently across all methods for package management and GlobalConfig operations

### Configuration Paths

- `paths.TemplatesAppDir`: Source directory for shell configuration templates
- `paths.AppDir`: Target directory for generated shell configuration
- Template file: `filepath.Join(paths.TemplatesAppDir, "devgita.zsh.tmpl")`
- Generated file: `filepath.Join(paths.AppDir, "devgita.zsh")`
- GlobalConfig file: `~/.config/devgita/global_config.yaml`

### Runtime Configuration

- **Global versions**: Stored in `~/.config/mise/config.toml` (managed by mise itself)
- **Local versions**: Stored in `.tool-versions` or `.mise.toml` files (project-specific)
- **Shell integration**: Managed via devgita's GlobalConfig and template system

## Implementation Notes

- **Shell Integration Nature**: Mise requires shell activation for automatic version switching
- **Template-Based Configuration**: Uses GlobalConfig and template regeneration for shell integration
- **Load-Modify-Regenerate-Save Pattern**: Each configuration method follows this transaction pattern
- **Fresh GlobalConfig Instances**: Each method creates a new `&config.GlobalConfig{}` and loads from disk to prevent stale data
- **Stateless Configuration**: GlobalConfig represents disk state, not app instance state
- **ForceInstall Logic**: Calls `Uninstall()` first, which now properly disables the shell integration feature
- **Runtime Manager Nature**: Unlike typical applications, mise manages other tools rather than being a standalone application
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Validation**: `UseGlobal()` validates that both language and version parameters are provided
- **Update Method**: Not implemented as mise updates should be handled by system package managers

## Template Integration

### Template Structure

The mise shell integration is defined in `configs/templates/devgita.zsh.tmpl`:

```bash
{{if .Mise}}
# Mise - Polyglot runtime manager
if command -v mise &> /dev/null; then
  eval "$(mise activate zsh)"
fi
{{end}}
```

### GlobalConfig Tracking

The feature state is tracked in `~/.config/devgita/global_config.yaml`:

```yaml
shell:
  mise: true  # Enabled
  # ... other shell features
```

### Generated Configuration

When enabled, the generated `devgita.zsh` contains:

```bash
# Mise - Polyglot runtime manager
if command -v mise &> /dev/null; then
  eval "$(mise activate zsh)"
fi
```

This approach:
- Provides single source of truth (template file)
- Enables clean enable/disable operations
- Prevents configuration conflicts
- Makes tracking and version control easier
- Ensures consistent regeneration

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

// Enable shell integration
err = mise.SoftConfigure()
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
- `Setup()` → Use `ForceConfigure()` instead
- `MaybeSetup()` → Use `SoftConfigure()` instead
- `Run()` → Use `ExecuteCommand()` instead

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Runtime Installation Fails**: Check internet connectivity and tool availability
3. **Commands Don't Work**: Verify mise is installed and accessible in PATH
4. **Version Conflicts**: Use `mise current` to check active versions
5. **Shell Integration Not Working**: Verify GlobalConfig has feature enabled and shell configuration is sourced
6. **GlobalConfig Load Errors**: Ensure `~/.config/devgita/global_config.yaml` is valid YAML
7. **Template Not Found**: Verify `configs/templates/devgita.zsh.tmpl` exists in devgita repository

### Platform Considerations

- **macOS**: Installed via Homebrew package manager
- **Linux**: Installed via apt package manager
- **Shell Integration**: Requires shell activation for automatic version switching
- **PATH Management**: Mise manages runtime paths automatically

### Shell Integration

- Shell integration requires Zsh shell to function
- Configuration must be sourced in shell initialization (add `source ~/.config/devgita/devgita.zsh` to `.zshrc`)
- Works best with devgita's complete shell setup
- May conflict with other runtime managers if not properly configured

### Template System Benefits

- **Single Source of Truth**: Template file is the only place shell configuration is defined
- **Trackable**: Git can track template changes easily
- **Predictable**: Regeneration always produces same output for same inputs
- **No Conflicts**: No string manipulation or file appending/removing
- **Clean Uninstall**: Disabling a feature regenerates without it

This module provides essential runtime management capabilities with template-based shell integration for maintaining consistent development environments across different programming languages and projects.
