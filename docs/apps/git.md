# Git Module Documentation

## Overview

The Git module provides version control system installation and configuration management with devgita integration. It follows the standardized devgita app interface while providing Git-specific operations for repository management, branch operations, and working directory tasks.

## App Purpose

Git is the fundamental distributed version control system that tracks changes in source code during software development. This module ensures Git is properly installed and configured across macOS and Debian/Ubuntu systems with sensible defaults and devgita integration.

## Lifecycle Summary

1. **Installation**: Install Git package via platform package managers (Homebrew/apt)
2. **Configuration**: Apply devgita's Git configuration templates with user-friendly defaults
3. **Execution**: Provide high-level Git operations for common development workflows

## Exported Functions

| Function           | Purpose                   | Behavior                                                 |
| ------------------ | ------------------------- | -------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Git instance with platform-specific commands |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install Git                   |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (ignored), then `Install()`    |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing  |
| `ForceConfigure()` | Force configuration       | Overwrites existing configs with devgita defaults        |
| `SoftConfigure()`  | Conditional configuration | Preserves existing `.gitconfig` if present               |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                        |
| `ExecuteCommand()` | Execute git commands      | Runs arbitrary git commands with error handling          |
| `Update()`         | Update installation       | **Not implemented** - returns error                      |

## Installation Methods

### Install()

```go
git := git.New()
err := git.Install()
```

- **Purpose**: Standard Git installation
- **Behavior**: Uses `InstallPackage()` to install Git package
- **Use case**: Initial Git installation or explicit reinstall

### ForceInstall()

```go
git := git.New()
err := git.ForceInstall()
```

- **Purpose**: Force Git installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (ignored for Git), then `Install()`
- **Use case**: Ensure fresh Git installation or fix corrupted installation

### SoftInstall()

```go
git := git.New()
err := git.SoftInstall()
```

- **Purpose**: Install Git only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing Git installations

### Uninstall()

```go
err := git.Uninstall()
```

- **Purpose**: Remove Git installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Git is a fundamental tool that shouldn't be uninstalled via devgita

### Update()

```go
err := git.Update()
```

- **Purpose**: Update Git installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Git updates are typically handled by the system package manager

## Configuration Methods

### Configuration Paths

- **Source**: `paths.GitConfigAppDir` (devgita's git configs)
- **Destination**: `paths.GitConfigLocalDir` (user's config directory)
- **Marker file**: `.gitconfig` in `GitConfigLocalDir`

### ForceConfigure()

```go
err := git.ForceConfigure()
```

- **Purpose**: Apply Git configuration regardless of existing files
- **Behavior**: Copies all configs from app dir to local dir, overwriting existing
- **Use case**: Reset to devgita defaults, apply config updates

### SoftConfigure()

```go
err := git.SoftConfigure()
```

- **Purpose**: Apply Git configuration only if not already configured
- **Behavior**: Checks for `.gitconfig` marker file; if exists, does nothing
- **Marker logic**: `filepath.Join(paths.GitConfigLocalDir, ".gitconfig")`
- **Use case**: Initial setup that preserves user customizations

## Execution Methods

### ExecuteCommand()

```go
err := git.ExecuteCommand("status")
err := git.ExecuteCommand("add", ".")
err := git.ExecuteCommand("commit", "-m", "message")
```

- **Purpose**: Execute arbitrary git commands
- **Parameters**: Variable arguments passed directly to git
- **Error handling**: Wraps errors with context

### Git-Specific Operations

#### Repository Operations

```go
// Clone repository
err := git.Clone("https://github.com/user/repo.git", "/path/to/destination")

// Fetch from origin
err := git.FetchOrigin()

// Pull changes
err := git.Pull("main")        // Pull from origin/main
err := git.Pull("")           // Pull from current tracking branch
```

#### Branch Management

```go
// Switch to existing branch
err := git.SwitchBranch("feature-branch")

// Delete branch (soft)
err := git.DeleteBranch("feature-branch", false)

// Force delete branch
err := git.DeleteBranch("feature-branch", true)
```

#### Working Directory Operations

```go
// Apply stashed changes
err := git.Pop("branch-name")  // Note: branch parameter not used in implementation

// Deep clean ignored files
err := git.DeepClean("url", "path")  // Note: parameters not used in implementation

// Restore files from branch
err := git.Restore("main", "src/")   // Restore src/ from main
err := git.Restore("", "src/")       // Restore src/ from main (default)
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()`
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()`
3. **Update Configuration**: `New()` → `SoftInstall()` → `ForceConfigure()`
4. **Git Operations**: `New()` → `ExecuteCommand()` / specific methods

## Constants and Paths

### Relevant Constants

- `constants.Git`: Package name for Git installation
- Used by all installation methods for consistent package reference

### Configuration Paths

- `paths.GitConfigAppDir`: Source directory for devgita's Git configuration templates
- `paths.GitConfigLocalDir`: Target directory for user's Git configuration
- Configuration copying preserves directory structure and file permissions

## Implementation Notes

- **ForceInstall Logic**: Calls `Uninstall()` first but ignores the error since Git uninstall is not supported
- **Configuration Strategy**: Uses marker file (`.gitconfig`) to determine if configuration exists
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Unused Parameters**: Some Git-specific methods have legacy parameters that aren't used in implementation
- **Update Method**: Not implemented as Git updates should be handled by system package managers
