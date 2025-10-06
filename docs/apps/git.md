# Git Module Documentation

## Overview

The Git module provides version control system installation and configuration management with devgita integration. It follows the standardized devgita app interface while providing Git-specific operations.

## Installation Methods

### ForceInstall()

```go
git := git.New()
err := git.ForceInstall()
```

- **Purpose**: Installs Git package regardless of existing installation
- **Behavior**: Uses `InstallPackage()` to force installation
- **Use case**: When you need to ensure Git is installed or updated

### SoftInstall()

```go
git := git.New()
err := git.SoftInstall()
```

- **Purpose**: Installs Git only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing Git installations

### Uninstall()

```go
err := git.Uninstall()
```

- **Purpose**: Remove Git installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Git is a fundamental tool that shouldn't be uninstalled via devgita

## Configuration Methods

### Paths

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

### Git-Specific Methods

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

## Agent Guidelines

### Installation Strategy

1. **Use SoftInstall()** for standard setup to respect existing Git installations
2. **Use ForceInstall()** only when explicitly updating or fixing Git installation
3. **Never call Uninstall()** - it's not supported and will return an error

### Configuration Strategy

1. **Use SoftConfigure()** for initial setup to preserve user customizations
2. **Use ForceConfigure()** when resetting to devgita defaults or applying updates
3. **Check marker file**: Configuration is considered applied if `.gitconfig` exists in `GitConfigLocalDir`

### Method Selection

- **Standard workflow**: `SoftInstall()` → `SoftConfigure()`
- **Reset workflow**: `ForceInstall()` → `ForceConfigure()`
- **Update workflow**: `SoftInstall()` → `ForceConfigure()`

### Error Handling

- All methods return errors that should be checked
- Git operations may fail if repository is in invalid state
- Some methods have unused parameters (legacy from interface) - ignore them

### Common Patterns

```go
// Standard setup
git := git.New()
if err := git.SoftInstall(); err != nil {
    return err
}
if err := git.SoftConfigure(); err != nil {
    return err
}

// Quick git operations
if err := git.ExecuteCommand("status"); err != nil {
    return err
}

// Repository management
if err := git.Clone(repoURL, destPath); err != nil {
    return err
}
if err := git.SwitchBranch("main"); err != nil {
    return err
}
```

## Implementation Notes

- **Non-standard methods**: Some methods like `ForceInstall()`/`SoftInstall()` use different naming than the standard `Install()`/`MaybeInstall()` pattern
- **Unused parameters**: Some methods accept parameters that aren't used in the implementation (e.g., `Pop()`, `DeepClean()`)
- **Git-specific operations**: The module provides high-level Git operations beyond basic installation
- **Configuration dependency**: Git configuration depends on devgita's config templates in `configs/git/`

