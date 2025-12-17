# Devgita Module Documentation

## Overview

The Devgita module provides self-installation and configuration management for the devgita development environment manager itself. It follows the standardized devgita app interface while providing devgita-specific operations for repository cloning, global configuration management, and complete environment setup/teardown.

## App Purpose

Devgita is the core meta-application that manages its own installation and configuration. This module handles cloning the devgita repository from GitHub, initializing the global configuration file that tracks installed packages and their states, and managing the complete lifecycle of the devgita toolchain. Unlike other apps, devgita fully supports uninstallation and provides clean removal of all installed components and configuration.

## Lifecycle Summary

1. **Installation**: Clone devgita repository from GitHub to `~/.config/devgita/`
2. **Configuration**: Create and manage global configuration file (`global_config.yaml`) tracking installation state
3. **Uninstallation**: Complete removal of repository, global configuration, and config directories

## Exported Functions

| Function           | Purpose                   | Behavior                                                                |
| ------------------ | ------------------------- | ----------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Devgita instance with Git dependency                        |
| `Install()`        | Standard installation     | Clones devgita repository to app directory                              |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first, then `Install()` for clean reinstall         |
| `SoftInstall()`    | Conditional installation  | Checks if repository exists and has content before installing           |
| `ForceConfigure()` | Force configuration       | Creates or overwrites global config with current paths                  |
| `SoftConfigure()`  | Conditional configuration | Preserves existing global_config.yaml if present                        |
| `Uninstall()`      | Remove installation       | **Fully supported** - removes repository, config file, and directories  |

## Installation Methods

### Install()

```go
devgita := devgita.New()
err := devgita.Install()
```

- **Purpose**: Standard devgita repository installation
- **Behavior**: 
  - Creates app directory if it doesn't exist
  - Removes any existing content in app directory
  - Clones devgita repository from GitHub to `paths.AppDir`
- **Use case**: Initial devgita installation or explicit reinstall
- **Note**: Destructive operation - removes existing app directory content

### ForceInstall()

```go
devgita := devgita.New()
err := devgita.ForceInstall()
```

- **Purpose**: Force devgita installation regardless of existing state
- **Behavior**: 
  - Calls `Uninstall()` first (includes error handling)
  - Returns error if uninstall fails
  - Calls `Install()` to clone fresh repository
- **Use case**: Complete reinstall with clean state, fixing corrupted installation
- **Error handling**: Returns wrapped error if uninstall fails

### SoftInstall()

```go
devgita := devgita.New()
err := devgita.SoftInstall()
```

- **Purpose**: Install devgita only if not already present
- **Behavior**: 
  - Checks if `paths.AppDir` exists and is not empty
  - Logs info message if repository already exists
  - Calls `Install()` if directory doesn't exist or is empty
- **Use case**: Standard installation that preserves existing devgita repository
- **Detection logic**: `files.DirAlreadyExist() && !files.IsDirEmpty()`

### Uninstall()

```go
err := devgita.Uninstall()
```

- **Purpose**: Remove devgita installation and configuration
- **Behavior**: 
  - Removes app directory (repository) if exists and not empty
  - Removes empty app directory if it exists
  - Removes global config file if exists
  - Removes empty config directory if exists
  - Logs info messages for each removal operation
- **Use case**: Complete removal of devgita installation
- **Error handling**: Returns wrapped errors for any failed removals
- **Note**: Unlike other apps, devgita fully supports uninstallation

## Configuration Methods

### Configuration Paths

- **Config directory**: `~/.config/devgita/` (paths.ConfigDir + constants.AppName)
- **Global config file**: `~/.config/devgita/global_config.yaml`
- **App directory**: `~/.config/devgita/` (paths.AppDir - repository location)

### ForceConfigure()

```go
err := devgita.ForceConfigure()
```

- **Purpose**: Create or overwrite global configuration file
- **Behavior**:
  - Creates new GlobalConfig instance
  - Calls `Create()` to initialize config file
  - Loads existing config
  - Sets `AppPath` to current repository location
  - Sets `ConfigPath` to current config directory
  - Saves updated configuration
- **Use case**: Initialize configuration, update paths after repository changes
- **Configuration content**: 
  ```yaml
  app_path: /Users/username/.config/devgita
  config_path: /Users/username/.config/devgita
  installed:
    packages: []
    languages: []
    databases: []
  already_installed:
    packages: []
  ```

### SoftConfigure()

```go
err := devgita.SoftConfigure()
```

- **Purpose**: Create global configuration only if it doesn't exist
- **Behavior**: 
  - Checks for existing `global_config.yaml`
  - Logs info message if config already exists
  - Calls `ForceConfigure()` if config doesn't exist
- **Marker logic**: `files.FileAlreadyExist(globalConfigPath)`
- **Use case**: Initial configuration that preserves user customizations

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()`
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()`
3. **Update Configuration**: `New()` → `SoftInstall()` → `ForceConfigure()`
4. **Complete Removal**: `New()` → `Uninstall()`
5. **Clean Reinstall**: `New()` → `ForceInstall()` → `ForceConfigure()`

## Constants and Paths

### Relevant Constants

- `constants.AppName`: "devgita" - used for directory naming
- `constants.GlobalConfigFile`: "global_config.yaml" - config filename
- `constants.DevgitaRepositoryUrl`: GitHub repository URL for cloning

### Configuration Paths

- `paths.ConfigDir`: User's config directory (`~/.config`)
- `paths.AppDir`: Devgita repository location (`~/.config/devgita/`)
- `configDirPath`: Full config directory path (ConfigDir + AppName)
- `globalConfigPath`: Full path to global_config.yaml

## Implementation Notes

- **Git Dependency**: Uses Git module for repository cloning operations
- **Self-Management**: Unlike other apps, devgita manages its own installation
- **Full Uninstall Support**: Complete removal including repository and configuration
- **Destructive Install**: `Install()` removes existing app directory before cloning
- **Path Management**: Updates global config with current path locations
- **Error Handling**: All methods return wrapped errors with context
- **Directory Cleanup**: Removes empty directories during uninstallation
- **Logging**: Info messages for skipped operations and successful removals

## Global Configuration Structure

The global configuration file tracks devgita's state and installed packages:

### Configuration Schema

```yaml
# Global configuration at ~/.config/devgita/global_config.yaml
app_path: /Users/username/.config/devgita
config_path: /Users/username/.config/devgita

installed:
  packages:
    - neovim
    - alacritty
    - tmux
  languages:
    - node@lts
    - python@latest
    - go@latest
  databases:
    - postgresql
    - redis

already_installed:
  packages:
    - git
    - curl
```

### Configuration Fields

- **app_path**: Absolute path to devgita repository
- **config_path**: Absolute path to devgita configuration directory
- **installed.packages**: Packages installed by devgita (can be uninstalled)
- **installed.languages**: Language runtimes installed via Mise
- **installed.databases**: Database systems installed by devgita
- **already_installed.packages**: Pre-existing packages (won't be uninstalled)

## Installation Flow

### Standard Installation

```go
dg := devgita.New()

// Install repository
if err := dg.SoftInstall(); err != nil {
    return err
}

// Initialize configuration
if err := dg.SoftConfigure(); err != nil {
    return err
}
```

### Force Reinstall

```go
dg := devgita.New()

// Complete clean install
if err := dg.ForceInstall(); err != nil {
    return err
}

// Reset configuration
if err := dg.ForceConfigure(); err != nil {
    return err
}
```

### Complete Removal

```go
dg := devgita.New()

// Remove everything
if err := dg.Uninstall(); err != nil {
    return err
}
```

## Uninstall Behavior

The `Uninstall()` method performs comprehensive cleanup:

### Cleanup Steps

1. **App Directory Removal**
   - Checks if app directory exists
   - Logs info if directory doesn't exist
   - Removes empty directory with `os.Remove()`
   - Removes directory with contents using `os.RemoveAll()`
   - Returns error if removal fails

2. **Global Config Removal**
   - Checks if global_config.yaml exists
   - Logs info and removes file if exists
   - Logs info if file doesn't exist
   - Returns error if removal fails

3. **Config Directory Cleanup**
   - Checks if config directory exists and is empty
   - Removes empty config directory
   - Returns error if removal fails

### Uninstall Safety

- Non-destructive if directories don't exist (logs info instead)
- Handles both empty and non-empty directories
- Removes only devgita-specific directories
- Preserves other config directories in `~/.config/`
- Returns wrapped errors with context for debugging

## Testing Patterns

### Test Coverage

The devgita module includes comprehensive tests:

- **Installation Tests**: `TestSoftInstall_DirectoryDoesNotExist`, `TestSoftInstall_DirectoryExistsWithFiles`
- **Uninstall Tests**: `TestUninstall_DirectoryEmpty`, `TestUninstall_RemovesRepositoryAndConfig`, `TestUninstall_NoRepositoryNoConfig`
- **Configuration Tests**: `TestForceConfigure_CreatesConfig`, `TestForceConfigure_OverwritesExisting`, `TestSoftConfigure_PreservesExistingConfig`
- **Force Install Test**: `TestForceInstall`

### Testing Approach

```go
func init() {
    logger.Init(false)
}

func TestSoftInstall_DirectoryExistsWithFiles(t *testing.T) {
    tempDir := t.TempDir()
    
    // Create test file
    testFile := filepath.Join(tempDir, "test.txt")
    if err := os.WriteFile(testFile, []byte("content"), 0o644); err != nil {
        t.Fatalf("Failed to create test file: %v", err)
    }
    
    // Override paths
    oldAppDir := paths.AppDir
    paths.AppDir = tempDir
    t.Cleanup(func() {
        paths.AppDir = oldAppDir
    })
    
    dg := New()
    err := dg.SoftInstall()
    if err != nil {
        t.Fatalf("SoftInstall() failed: %v", err)
    }
    
    // Verify existing file was preserved
    if _, err := os.Stat(testFile); os.IsNotExist(err) {
        t.Fatal("Expected existing file to be preserved")
    }
}
```

### Testing Best Practices

- Use `t.TempDir()` for isolated test directories
- Override global paths with cleanup via `t.Cleanup()`
- Test both file existence and content preservation
- Verify error handling for all operations
- Test edge cases (empty directories, missing files, etc.)

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure Git is installed and GitHub is accessible
2. **Clone Error**: Check network connectivity and repository URL
3. **Permission Errors**: Verify write access to `~/.config/` directory
4. **Config Not Created**: Check directory permissions and path overrides
5. **Uninstall Fails**: Review error messages for specific file/directory issues

### Platform Considerations

- **macOS/Linux**: Standard installation to `~/.config/devgita/`
- **Git Dependency**: Requires Git to be installed for repository cloning
- **Configuration Location**: Uses XDG Base Directory specification
- **Path Resolution**: All paths are absolute for cross-platform compatibility

### Directory Structure

After successful installation and configuration:

```
~/.config/devgita/
├── configs/                    # Configuration templates
│   ├── alacritty/
│   ├── bash/
│   ├── fastfetch/
│   ├── fonts/
│   ├── neovim/
│   ├── themes/
│   └── tmux/
├── internal/                   # Go source code
├── pkg/                        # Shared utilities
├── cmd/                        # CLI commands
├── go.mod
├── go.sum
├── main.go
└── global_config.yaml          # Global state tracking
```

## Integration with Devgita Ecosystem

### Repository Management

- **Source**: GitHub repository cloned to local system
- **Updates**: Repository can be updated via `git pull` or `ForceInstall()`
- **Version Control**: Uses Git for version tracking and updates

### Configuration Management

- **Global State**: Tracks all installed packages and their states
- **Package Tracking**: Differentiates devgita-installed vs pre-existing packages
- **Safe Uninstall**: Only removes packages explicitly installed by devgita

### CLI Integration

The devgita module is used by the `dg install` command:

```bash
# Standard installation
dg install

# Specific categories
dg install --only terminal
dg install --skip desktop

# Reinstall
dg install --force
```

## Key Features

### Self-Contained Installation
- Clones entire repository with all configuration templates
- No external dependencies beyond Git
- Single command installation process

### State Management
- Tracks installed vs pre-existing packages
- Maintains global configuration file
- Enables safe uninstallation

### Clean Uninstall
- Complete removal of all devgita components
- Removes repository and configuration
- Cleans up empty directories

### Path Flexibility
- Uses standard config directory conventions
- Supports path overrides for testing
- Absolute paths for cross-platform compatibility

## External References

- **Devgita Repository**: https://github.com/cjairm/devgita
- **Configuration Spec**: `internal/config/fromFile.go`
- **Global Config**: `configs/bash/global_config.yaml` (template)
- **Installation Guide**: `docs/project-overview.md`

This module provides the foundational self-installation and configuration management capabilities that enable devgita to bootstrap and manage its entire development environment ecosystem.
