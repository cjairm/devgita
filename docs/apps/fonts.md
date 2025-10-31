# Fonts Module Documentation

## Overview

The Fonts module provides font installation and management with devgita integration. It follows the standardized devgita app interface while providing font-specific operations for installing developer-focused fonts, managing font collections, and handling font-related desktop applications.

## App Purpose

Fonts in devgita refer to developer-oriented typefaces, particularly Nerd Fonts that include programming ligatures and icon glyphs. This module ensures proper fonts are installed across macOS and Debian/Ubuntu systems for terminal emulators, code editors, and development environments.

## Lifecycle Summary

1. **Installation**: Install font packages via platform package managers (Homebrew/apt)
2. **Configuration**: Fonts don't require separate configuration files - installation handles setup
3. **Execution**: Provide high-level font operations for collection management and individual font installation

## Exported Functions

| Function                 | Purpose                   | Behavior                                                             |
| ------------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`                  | Factory method            | Creates new Fonts instance with platform-specific commands           |
| `Install(fontName)`      | Standard installation     | Installs specified font by name                                      |
| `ForceInstall(fontName)` | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall(fontName)`  | Conditional installation  | Uses `MaybeInstallFont()` to check before installing font            |
| `ForceConfigure()`       | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`        | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`            | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()`       | Execute font commands     | **Not applicable** - returns nil                                     |
| `Update()`               | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
fonts := fonts.New()
err := fonts.Install("font-jetbrains-mono-nerd-font")
```

- **Purpose**: Standard font installation with specified font name
- **Parameters**: `fontName string` - Name of the font package to install
- **Behavior**: Uses `InstallDesktopApp()` to install the specified font
- **Use case**: Install specific font for development environment

### ForceInstall()

```go
fonts := fonts.New()
err := fonts.ForceInstall("font-hack-nerd-font")
```

- **Purpose**: Force font installation regardless of existing state
- **Parameters**: `fontName string` - Name of the font package to install
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh font installation or fix corrupted installation

### SoftInstall()

```go
fonts := fonts.New()
err := fonts.SoftInstall("font-meslo-lg-nerd-font")
```

- **Purpose**: Install specified font only if not already present
- **Parameters**: `fontName string` - Name of the font package to install
- **Behavior**: Uses `MaybeInstallFont()` to check before installing the specified font
- **Use case**: Standard installation that respects existing font installations

### Uninstall()

```go
err := fonts.Uninstall()
```

- **Purpose**: Remove font installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Font uninstallation requires careful system-level handling

### Update()

```go
err := fonts.Update()
```

- **Purpose**: Update font installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Font updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := fonts.ForceConfigure()
err := fonts.SoftConfigure()
```

- **Purpose**: Apply font configuration
- **Behavior**: **Not applicable** - always returns nil
- **Rationale**: Fonts don't have separate configuration files; installation handles setup

## Execution Methods

### ExecuteCommand()

```go
err := fonts.ExecuteCommand("--version")
```

- **Purpose**: Execute font-related commands
- **Behavior**: **Not applicable** - returns nil (success)
- **Rationale**: Fonts are desktop applications without CLI commands, but returns success for interface compliance

## Font-Specific Operations

### Available()

```go
fonts := fonts.New()
availableFonts := fonts.Available()
```

- **Purpose**: Get list of available fonts for installation
- **Returns**: String slice of font package names
- **Available fonts**:
  - `font-hack-nerd-font`
  - `font-meslo-lg-nerd-font`
  - `font-caskaydia-mono-nerd-font`
  - `font-fira-mono`
  - `font-jetbrains-mono-nerd-font`

### SoftInstallAll()

```go
err := fonts.SoftInstallAll()
```

- **Purpose**: Install all available fonts from the collection
- **Behavior**: Iterates through `Available()` list and calls `SoftInstall()` for each
- **Error handling**: Returns last error encountered, but continues installation
- **Use case**: Set up complete development font collection

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall("font-name")` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall("font-name")` → `ForceConfigure()` (no-op)
3. **Collection Setup**: `New()` → `SoftInstallAll()`
4. **Individual Font**: `New()` → `Install("font-name")`
5. **Safe Individual**: `New()` → `SoftInstall("font-name")`

## Constants and Paths

### Default Font

- **No default**: All methods require explicit `fontName` parameter
- Recommended starter: `"font-jetbrains-mono-nerd-font"`
- JetBrains Mono chosen for excellent programming ligatures and readability

### Font Collection

All fonts in the `Available()` list are Nerd Fonts or developer-focused typefaces:

- **Hack Nerd Font**: Clean, highly legible programming font
- **Meslo LG Nerd Font**: Customized version of Apple's Menlo
- **Caskaydia Mono Nerd Font**: Microsoft's Cascadia Code with Nerd Font patches
- **Fira Mono**: Mozilla's monospaced font designed for coding
- **JetBrains Mono Nerd Font**: JetBrains' font optimized for developers

## Implementation Notes

- **Desktop App Installation**: Uses `InstallDesktopApp()` and `MaybeInstallFont()` for platform-appropriate installation
- **No Configuration**: Fonts don't require separate config files - installation handles system integration
- **Error Propagation**: `SoftInstallAll()` continues installing even if individual fonts fail
- **ForceInstall Logic**: Calls `Uninstall()` first (returns error if it fails) since font uninstall is not supported
- **Platform Independence**: Works on both macOS (Homebrew) and Linux (apt) through desktop app installation
- **Font Caching**: System handles font cache updates automatically after installation
- **Required Parameters**: All installation methods require explicit `fontName` parameter
- **Legacy Compatibility**: Maintains deprecated function aliases for backward compatibility

## Font Installation Process

### macOS (Homebrew)

```bash
# Individual font installation
brew install --cask font-jetbrains-mono-nerd-font

# Multiple fonts via tap
brew tap homebrew/cask-fonts
brew install --cask font-hack-nerd-font font-meslo-lg-nerd-font
```

### Linux (apt)

```bash
# Font installation via package manager
sudo apt install fonts-hack-ttf fonts-firacode

# Manual installation to user font directory
mkdir -p ~/.local/share/fonts
cp *.ttf ~/.local/share/fonts/
fc-cache -fv
```

## Deprecated Functions

The module maintains backward compatibility through deprecated functions:

- `MaybeInstall(fontName string)` → Use `SoftInstall(fontName)` instead
- `MaybeInstallAll()` → Use `SoftInstallAll()` instead

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Font Not Appearing**: Restart applications or refresh font cache with `fc-cache -fv` (Linux)
3. **Permission Issues**: Verify write access to system font directories
4. **Homebrew Tap Missing**: Run `brew tap homebrew/cask-fonts` on macOS

### Platform Considerations

- **macOS**: Fonts installed via Homebrew cask to `/Library/Fonts/`
- **Linux**: Fonts installed via apt to system directories or `~/.local/share/fonts/`
- **Font Rendering**: Different platforms may render fonts slightly differently
- **Application Support**: Not all applications automatically detect newly installed fonts

### Font Verification

```bash
# Check if font is installed (macOS)
ls /Library/Fonts/ | grep -i jetbrains

# Check if font is installed (Linux)
fc-list | grep -i "JetBrains Mono"

# Test font in terminal
echo "Font test: → ← ↑ ↓ ★ ♠ ♥ ♦ ♣"
```

This module provides essential font management capabilities for creating a consistent, visually appealing development environment across different platforms and applications.

