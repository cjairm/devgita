# FontConfig Module Documentation

## Overview

The FontConfig module provides installation and command execution management for the fontconfig library with devgita integration. It follows the standardized devgita app interface while providing fontconfig-specific operations for font utilities and cache management.

### Current Implementation Status

**Implemented Features**:

- ✅ Package installation (`Install()`, `SoftInstall()`)
- ✅ Command execution for fontconfig utilities (`ExecuteCommand()`)
- ✅ Command validation (fc-cache, fc-list, fc-match, fc-pattern)

**Not Yet Implemented**:

- ⏸️ Configuration management (`ForceConfigure()`, `SoftConfigure()`)
- ⏸️ Font rendering settings templates
- ⏸️ Automatic font cache rebuilding after configuration

**Not Supported**:

- ❌ Uninstallation (system library - managed by OS)
- ❌ Updates (handled by system package manager)
- ❌ Deprecated wrapper functions (minimal API by design)

## App Purpose

Fontconfig is a library for configuring and customizing font access on Linux/Unix systems. It is used by many applications for font rendering and management, providing:

- Font discovery and matching
- Font pattern configuration
- Font cache management
- Antialiasing and hinting configuration
- Font substitution rules

This module ensures fontconfig is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems and provides configuration management for optimal font rendering in development environments.

## Lifecycle Summary

1. **Installation**: Install fontconfig package via platform package managers (Homebrew/apt)
2. **Configuration**: Currently not implemented - future support for devgita's fontconfig configuration templates
3. **Execution**: Provide fontconfig command execution for font cache management and font utilities

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new FontConfig instance with platform-specific commands      |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install fontconfig                        |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not implemented** - returns error                                  |
| `SoftConfigure()`  | Conditional configuration | **Not implemented** - returns error                                  |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute fontconfig utils  | Runs fontconfig utilities (fc-cache, fc-list, fc-match, fc-pattern)  |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
fc := fontconfig.New()
err := fc.Install()
```

- **Purpose**: Standard fontconfig installation
- **Behavior**: Uses `InstallPackage()` to install fontconfig package
- **Use case**: Initial fontconfig installation or explicit reinstall

### ForceInstall()

```go
fc := fontconfig.New()
err := fc.ForceInstall()
```

- **Purpose**: Force fontconfig installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh fontconfig installation or fix corrupted installation

### SoftInstall()

```go
fc := fontconfig.New()
err := fc.SoftInstall()
```

- **Purpose**: Install fontconfig only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing fontconfig installations

### Uninstall()

```go
err := fc.Uninstall()
```

- **Purpose**: Remove fontconfig installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Fontconfig is a system library used by many applications; uninstalling can break system functionality

### Update()

```go
err := fc.Update()
```

- **Purpose**: Update fontconfig installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Fontconfig updates are typically handled by the system package manager

## Configuration Methods

### Configuration Status

**Note**: Configuration methods are currently **not implemented**. The following sections describe the planned functionality for future implementation.

### Configuration Paths (Planned)

- **Source**: `paths.FontConfigConfigAppDir` (devgita's fontconfig configs)
- **Destination**: `paths.FontConfigConfigLocalDir` (user's config directory, typically `~/.config/fontconfig/`)
- **Marker file**: `fonts.conf` in `FontConfigConfigLocalDir`

### ForceConfigure()

```go
err := fc.ForceConfigure()
```

- **Purpose**: Apply fontconfig configuration regardless of existing files
- **Current Status**: **Not implemented** - returns error
- **Planned Behavior**:
  - Copy all configs from app dir to local dir, overwriting existing
  - Update font cache via `fc-cache -fv`
- **Future Use Case**: Reset to devgita defaults, apply config updates

### SoftConfigure()

```go
err := fc.SoftConfigure()
```

- **Purpose**: Apply fontconfig configuration only if not already configured
- **Current Status**: **Not implemented** - returns error
- **Planned Behavior**: Check for `fonts.conf` marker file; if exists, preserve user config
- **Future Marker Logic**: `filepath.Join(paths.FontConfigConfigLocalDir, "fonts.conf")`
- **Future Use Case**: Initial setup that preserves user customizations

## Execution Methods

### ExecuteCommand()

```go
err := fc.ExecuteCommand("fc-cache", "-fv")
err := fc.ExecuteCommand("fc-list")
err := fc.ExecuteCommand("fc-match", "monospace")
err := fc.ExecuteCommand("fc-pattern", "sans-serif")
```

- **Purpose**: Execute fontconfig utilities with provided arguments
- **Parameters**:
  - `fontConfigCmd string`: Command name (must be one of: fc-cache, fc-list, fc-match, fc-pattern)
  - `args ...string`: Variable arguments passed to the command
- **Supported Commands**:
  - `fc-cache`: Build font information cache
  - `fc-list`: List available fonts
  - `fc-match`: Match available fonts
  - `fc-pattern`: Parse and validate patterns
- **Validation**: Returns error if command name is empty or unsupported
- **Error Handling**: Wraps errors with context from BaseCommand.ExecCommand

### Fontconfig-Specific Operations

The fontconfig utilities provide font management capabilities:

#### Font Cache Management

```bash
# Rebuild font cache
fc-cache -fv

# Clean cache
fc-cache -r

# Update cache for specific directory
fc-cache /usr/share/fonts
```

#### Font Listing

```bash
# List all fonts
fc-list

# List fonts with details
fc-list : family style file

# List monospace fonts
fc-list :mono

# Search for specific font
fc-list | grep "JetBrains"
```

#### Font Matching

```bash
# Match font pattern
fc-match monospace

# Detailed match information
fc-match -v "DejaVu Sans Mono"

# Match with specific properties
fc-match "monospace:weight=bold:slant=italic"
```

#### Pattern Validation

```bash
# Parse and validate pattern
fc-pattern "monospace:size=12"

# Show pattern details
fc-pattern -V "sans-serif"
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()`
2. **Force Setup**: `New()` → `ForceInstall()` (will fail - Uninstall not supported)
3. **Fontconfig Operations**: `New()` → `SoftInstall()` → `ExecuteCommand()` with specific fontconfig utilities
4. **Cache Rebuild**: `New()` → `SoftInstall()` → `ExecuteCommand("fc-cache", "-fv")`

**Note**: Configuration methods (`ForceConfigure()`, `SoftConfigure()`) are not yet implemented.

## Constants and Paths

### Relevant Constants

- `constants.FontConfig`: Package name ("fontconfig") for installation
- Used by all installation methods for consistent package reference

### Configuration Paths

- `paths.FontConfigConfigAppDir`: Source directory for devgita's fontconfig configuration templates
- `paths.FontConfigConfigLocalDir`: Target directory for user's fontconfig configuration (typically `~/.config/fontconfig/`)
- Configuration copying preserves XML structure and file permissions

## Implementation Notes

- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error since fontconfig uninstall is not supported
- **Configuration Strategy**: Not yet implemented - planned to use marker file (`fonts.conf`) to determine if configuration exists
- **Command Validation**: `ExecuteCommand()` validates command names - only fc-cache, fc-list, fc-match, and fc-pattern are supported
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **System Library**: Fontconfig is a core system library; managed carefully to avoid breaking dependencies
- **Update Method**: Not implemented - fontconfig updates should be handled by system package managers
- **No Deprecated Functions**: This module does not include deprecated wrapper functions (MaybeInstall, Setup, MaybeSetup, Run)

## Configuration Structure (Planned)

**Note**: Configuration is not yet implemented. The following describes the planned configuration structure for future implementation.

The fontconfig configuration (`fonts.conf`) will typically include:

### Basic Configuration

```xml
<?xml version="1.0"?>
<!DOCTYPE fontconfig SYSTEM "fonts.dtd">
<fontconfig>
  <!-- Font directories -->
  <dir>/usr/share/fonts</dir>
  <dir>/usr/local/share/fonts</dir>
  <dir>~/.fonts</dir>

  <!-- Cache directory -->
  <cachedir>~/.cache/fontconfig</cachedir>
</fontconfig>
```

### Rendering Settings

```xml
<fontconfig>
  <!-- Enable antialiasing -->
  <match target="font">
    <edit name="antialias" mode="assign">
      <bool>true</bool>
    </edit>
  </match>

  <!-- Hinting settings -->
  <match target="font">
    <edit name="hinting" mode="assign">
      <bool>true</bool>
    </edit>
    <edit name="hintstyle" mode="assign">
      <const>hintfull</const>
    </edit>
  </match>

  <!-- RGB subpixel rendering -->
  <match target="font">
    <edit name="rgba" mode="assign">
      <const>rgb</const>
    </edit>
  </match>
</fontconfig>
```

### Font Substitution

```xml
<fontconfig>
  <!-- Prefer specific fonts for monospace -->
  <alias>
    <family>monospace</family>
    <prefer>
      <family>JetBrains Mono</family>
      <family>DejaVu Sans Mono</family>
      <family>Liberation Mono</family>
    </prefer>
  </alias>

  <!-- Prefer specific fonts for sans-serif -->
  <alias>
    <family>sans-serif</family>
    <prefer>
      <family>Noto Sans</family>
      <family>DejaVu Sans</family>
    </prefer>
  </alias>
</fontconfig>
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Configuration Not Applied**: Configuration is not yet implemented - use `ExecuteCommand()` for manual font management
3. **Font Cache Issues**: Manually rebuild with `fc-cache -fv` using `ExecuteCommand("fc-cache", "-fv")`
4. **Fonts Not Appearing**: Check font directories with `fc-list` command
5. **Rendering Issues**: Configure manually via `~/.config/fontconfig/fonts.conf`
6. **Unsupported Command Error**: Only fc-cache, fc-list, fc-match, and fc-pattern are supported
7. **Empty Command Error**: Command name must be provided to `ExecuteCommand()`

### Platform Considerations

- **macOS**: Installed via Homebrew; native font system often preferred
- **Linux**: Installed via apt; critical for X11/Wayland font rendering
- **Configuration Location**: `~/.config/fontconfig/` (user) or `/etc/fonts/` (system)
- **Font Directories**: Platform-specific default font locations

### Font Cache Management

Font cache rebuilding is necessary after:

- Installing new fonts
- Modifying font configuration
- Adding/removing font directories
- System upgrades

Rebuild command:

```bash
fc-cache -fv
```

## Testing Patterns

### Test Structure

```go
func init() {
    testutil.InitLogger()
}

func TestForceConfigure(t *testing.T) {
    cleanup := testutil.SetupIsolatedPaths(t)
    defer cleanup()

    appDir, configDir, _, _ := testutil.SetupTestDirs(t)
    // ... test logic

    testutil.VerifyNoRealCommands(t, mockApp.Base)
}
```

### Test Coverage

- ✅ `TestNew` - Constructor validation
- ✅ `TestInstall` - Package installation
- ✅ `TestSoftInstall` - Conditional installation
- ⏭️ `TestForceInstall` - Skipped (Uninstall not supported)
- ✅ `TestForceConfigure` - Verifies "not implemented" error
- ✅ `TestSoftConfigure` - Verifies "not implemented" error (both scenarios)
- ⏭️ `TestUninstall` - Skipped (not supported for system library)
- ✅ `TestExecuteCommand` - Command execution and validation
  - fc-cache command
  - fc-list command
  - fc-match command
  - fc-pattern command
  - Unsupported command validation
  - Empty command validation
  - Error handling
- ✅ `TestUpdate` - Verifies "not implemented" error

## Integration with Devgita

Fontconfig integrates with devgita's terminal category as a core system library:

- **Installation**: Part of terminal core tools installation or standalone
- **Configuration**: Not yet implemented - planned for future releases
- **Tracking**: Registered in GlobalConfig as package when installed
- **Updates**: Managed through system package manager
- **Category**: Terminal/Core system library
- **Command Execution**: Provides direct access to fontconfig utilities (fc-cache, fc-list, fc-match, fc-pattern)

## External References

- **Fontconfig Documentation**: https://www.freedesktop.org/wiki/Software/fontconfig/
- **Configuration Guide**: https://www.freedesktop.org/software/fontconfig/fontconfig-user.html
- **Font Utilities**: https://linux.die.net/man/1/fc-cache
- **Testing Patterns**: `docs/guides/testing-patterns.md`
- **Project Overview**: `docs/project-overview.md`

## Key Features

### Font Discovery

- Automatic font directory scanning
- Font cache for fast lookups
- Pattern-based font matching
- Font family and style resolution

### Rendering Configuration

- Antialiasing control
- Subpixel rendering (RGB/BGR/VRGB/VBGR)
- Hinting styles (none/slight/medium/full)
- LCD filter configuration

### Font Substitution

- Family-based substitution
- Style-based matching
- Language-specific font selection
- Fallback font chains

### Performance

- Font information caching
- Fast pattern matching
- Lazy font loading
- Memory-efficient font database

## Use Cases in Development

### Terminal Emulators

- Configure monospace fonts for terminals (Alacritty, tmux, etc.)
- Ensure consistent font rendering across different terminals
- Optimize for code readability

### Code Editors

- Set preferred programming fonts (JetBrains Mono, Fira Code, etc.)
- Enable ligatures for code
- Configure font fallbacks for special characters

### System Integration

- Consistent font rendering across all applications
- High-quality antialiasing for readability
- Proper emoji and symbol rendering

This module provides essential font configuration and rendering capabilities for creating a visually consistent development environment within the devgita ecosystem.
