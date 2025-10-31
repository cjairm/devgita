# Neovim Module Documentation

## Overview

The Neovim module provides installation and configuration management for Neovim editor with devgita integration. It follows the standardized devgita app interface while providing Neovim-specific operations for editor setup, LSP configuration, version validation, and development environment customization.

## App Purpose

Neovim is a hyperextensible Vim-based text editor that provides powerful features for code editing, including Language Server Protocol (LSP) support, syntax highlighting, and extensive plugin ecosystem. This module ensures Neovim is properly installed and configured with devgita's optimized settings for development workflows, including kickstart configurations and custom editor setups.

## Lifecycle Summary

1. **Installation**: Install Neovim package via platform package managers (Homebrew/apt)
2. **Configuration**: Apply devgita's Neovim configuration templates with developer-focused settings including LSP and kickstart setup
3. **Execution**: Provide high-level Neovim operations for editor management, version validation, and command execution

## Exported Functions

| Function           | Purpose                   | Behavior                                                                |
| ------------------ | ------------------------- | ----------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Neovim instance with platform-specific commands             |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install Neovim                               |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()`    |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing                 |
| `ForceConfigure()` | Force configuration       | Validates version and overwrites existing configs with devgita defaults |
| `SoftConfigure()`  | Conditional configuration | Preserves existing init.lua if present                                  |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                       |
| `ExecuteCommand()` | Execute neovim commands   | Runs nvim with provided arguments                                       |
| `Update()`         | Update installation       | **Not implemented** - returns error                                     |
| `CheckVersion()`   | Version validation        | Validates Neovim version meets minimum requirements                     |

## Installation Methods

### Install()

```go
neovim := neovim.New()
err := neovim.Install()
```

- **Purpose**: Standard Neovim installation
- **Behavior**: Uses `InstallPackage()` to install Neovim package
- **Use case**: Initial Neovim installation or explicit reinstall

### ForceInstall()

```go
neovim := neovim.New()
err := neovim.ForceInstall()
```

- **Purpose**: Force Neovim installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh Neovim installation or fix corrupted installation

### SoftInstall()

```go
neovim := neovim.New()
err := neovim.SoftInstall()
```

- **Purpose**: Install Neovim only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing Neovim installations

### Uninstall()

```go
err := neovim.Uninstall()
```

- **Purpose**: Remove Neovim installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Text editors are typically managed at the system level

### Update()

```go
err := neovim.Update()
```

- **Purpose**: Update Neovim installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Neovim updates are typically handled by the system package manager

## Configuration Methods

### Configuration Paths

- **Source**: `paths.NeovimConfigAppDir` (devgita's neovim configs)
- **Destination**: `paths.NvimConfigLocalDir` (user's config directory)
- **Marker file**: `init.lua` in `NvimConfigLocalDir`

### ForceConfigure()

```go
err := neovim.ForceConfigure()
```

- **Purpose**: Apply Neovim configuration regardless of existing files
- **Behavior**:
  - Validates Neovim version meets minimum requirements via `CheckVersion()`
  - Copies all configs from app dir to local dir, overwriting existing
- **Use case**: Reset to devgita defaults, apply config updates

### SoftConfigure()

```go
err := neovim.SoftConfigure()
```

- **Purpose**: Apply Neovim configuration only if not already configured
- **Behavior**: Checks for `init.lua` marker file; if exists, does nothing
- **Marker logic**: `filepath.Join(paths.NvimConfigLocalDir, "init.lua")`
- **Use case**: Initial setup that preserves user customizations

## Execution Methods

### ExecuteCommand()

```go
err := neovim.ExecuteCommand("--version")
err := neovim.ExecuteCommand("--config", "/path/to/config.lua")
err := neovim.ExecuteCommand("+checkhealth")
```

- **Purpose**: Execute neovim commands with provided arguments
- **Parameters**: Variable arguments passed directly to nvim binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### Neovim-Specific Operations

The Neovim CLI provides extensive configuration and execution options:

#### Editor Operations

```bash
# Launch with specific configuration
nvim --config /path/to/init.lua

# Check plugin health
nvim +checkhealth

# Execute Lua commands
nvim +"lua print('Hello World')"

# Run headless for scripting
nvim --headless -c "echo 'test'" -c "quit"
```

#### Configuration Management

```bash
# Check version
nvim --version

# List runtime paths
nvim --cmd "echo &runtimepath" -c "quit"

# Profile startup time
nvim --startuptime startup.log

# Clean mode (no plugins)
nvim --clean
```

### CheckVersion()

```go
err := neovim.CheckVersion()
```

- **Purpose**: Validate installed Neovim version meets minimum requirements
- **Behavior**: Executes `nvim --version` and compares against `constants.NeovimVersion`
- **Version logic**: Uses `isVersionEqualOrHigher()` for semantic version comparison
- **Use case**: Ensure compatibility before applying configurations

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()`
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()`
3. **Update Configuration**: `New()` → `SoftInstall()` → `ForceConfigure()`
4. **Neovim Operations**: `New()` → `ExecuteCommand()` with specific nvim arguments
5. **Version Validation**: `New()` → `CheckVersion()` → `ForceConfigure()`

## Constants and Paths

### Relevant Constants

- `constants.Neovim`: Package name ("neovim") for installation
- `constants.Nvim`: Binary name ("nvim") for execution
- `constants.NeovimVersion`: Minimum required version (e.g., "0.11.1")
- Used by all installation methods for consistent package reference

### Configuration Paths

- `paths.NeovimConfigAppDir`: Source directory for devgita's Neovim configuration templates
- `paths.NvimConfigLocalDir`: Target directory for user's Neovim configuration
- Configuration copying preserves Lua structure and file permissions

## Implementation Notes

- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since Neovim uninstall is not supported
- **Configuration Strategy**: Uses marker file (`init.lua`) to determine if configuration exists
- **Version Validation**: `ForceConfigure()` always validates version before applying configuration
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Legacy Compatibility**: Maintains deprecated function aliases to prevent breaking existing code
- **Update Method**: Not implemented as Neovim updates should be handled by system package managers

## Configuration Structure

The Neovim configuration (`init.lua`) typically includes:

### Basic Configuration

```lua
-- Set leader key
vim.g.mapleader = ' '
vim.g.maplocalleader = ' '

-- Basic options
vim.opt.number = true
vim.opt.relativenumber = true
vim.opt.mouse = 'a'
vim.opt.breakindent = true
vim.opt.undofile = true
```

### LSP Configuration

```lua
-- LSP setup with kickstart template
require('kickstart.plugins.lspconfig')

-- Autocompletion
require('kickstart.plugins.autocompletion')

-- Syntax highlighting
require('kickstart.plugins.treesitter')
```

### Plugin Management

```lua
-- Lazy.nvim plugin manager
local lazypath = vim.fn.stdpath 'data' .. '/lazy/lazy.nvim'
if not vim.loop.fs_stat(lazypath) then
  vim.fn.system {
    'git',
    'clone',
    '--filter=blob:none',
    'https://github.com/folke/lazy.nvim.git',
    '--branch=stable',
    lazypath,
  }
end
vim.opt.rtp:prepend(lazypath)
```

## Version Validation

The `CheckVersion()` function implements semantic version comparison:

### Version Parsing

```go
func isVersionEqualOrHigher(currentVersion, requiredVersion string) bool {
    currentParts := strings.Split(currentVersion, ".")
    requiredParts := strings.Split(requiredVersion, ".")

    for i, requiredPartStr := range requiredParts {
        if i >= len(currentParts) {
            return false // Current version has fewer parts
        }
        // Compare numeric parts...
    }
    return true
}
```

### Version Format

- Expected format: "MAJOR.MINOR.PATCH" (e.g., "0.11.1")
- Comparison: Each component compared numerically
- Validation: Ensures current >= required version

## Deprecated Functions

The module maintains backward compatibility through deprecated functions:

- `MaybeInstall()` → Use `SoftInstall()` instead
- `Setup()` → Use `ForceConfigure()` instead
- `MaybeSetup()` → Use `SoftConfigure()` instead
- `Run()` → Use `ExecuteCommand()` instead

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Configuration Not Applied**: Check file permissions in config directory
3. **Commands Don't Work**: Verify Neovim is installed and accessible in PATH
4. **Version Too Old**: Update Neovim to meet minimum version requirements
5. **LSP Issues**: Check `:checkhealth` for plugin and LSP configuration problems

### Platform Considerations

- **macOS**: Installed via Homebrew package manager
- **Linux**: Installed via apt package manager
- **Configuration Location**: Cross-platform config directory handling
- **Plugin Dependencies**: May require additional tools for LSP and syntax highlighting

### Configuration References

- **Kickstart**: https://github.com/nvim-lua/kickstart.nvim
- **Personal Config**: https://github.com/cjairm/devenv/blob/main/nvim/init.lua
- **Releases**: https://github.com/neovim/neovim/releases
- **Color Schemes**: https://linovox.com/the-best-color-schemes-for-neovim-nvim/

This module provides essential text editor capabilities with modern LSP support for creating a powerful development environment within the devgita ecosystem.

