# OpenCode Module Documentation

## Overview

The OpenCode module provides installation and configuration management for OpenCode AI-powered code editor with devgita integration. It follows the standardized devgita app interface while providing OpenCode-specific operations for editor setup, theme customization, and template-based configuration management.

## App Purpose

OpenCode is an AI-powered code editor designed for modern development workflows. This module ensures OpenCode is properly installed and configured with devgita's optimized settings, including template-based configuration generation and customizable theme support. The module provides flexible theme management allowing users to choose between the default Gruvbox theme or custom themes.

## Lifecycle Summary

1. **Installation**: Install OpenCode package via platform package managers (Homebrew/apt)
2. **Configuration**: Generate configuration from templates with theme customization support
3. **Execution**: Provide high-level OpenCode operations for editor management and command execution

## Exported Functions

| Function           | Purpose                   | Behavior                                                                |
| ------------------ | ------------------------- | ----------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new OpenCode instance with platform-specific commands           |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install OpenCode                             |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()`    |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing                 |
| `ForceConfigure()` | Force configuration       | Generates config from template, overwrites existing                      |
| `SoftConfigure()`  | Conditional configuration | Preserves existing opencode.json if present                             |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                       |
| `ExecuteCommand()` | Execute opencode commands | Runs opencode with provided arguments                                    |
| `Update()`         | Update installation       | **Not implemented** - returns error                                     |

## Installation Methods

### Install()

```go
opencode := opencode.New()
err := opencode.Install()
```

- **Purpose**: Standard OpenCode installation
- **Behavior**: Uses `InstallPackage()` to install OpenCode package
- **Use case**: Initial OpenCode installation or explicit reinstall

### ForceInstall()

```go
opencode := opencode.New()
err := opencode.ForceInstall()
```

- **Purpose**: Force OpenCode installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh OpenCode installation or fix corrupted installation

### SoftInstall()

```go
opencode := opencode.New()
err := opencode.SoftInstall()
```

- **Purpose**: Install OpenCode only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing OpenCode installations

### Uninstall()

```go
err := opencode.Uninstall()
```

- **Purpose**: Remove OpenCode installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Code editors are typically managed at the system level

### Update()

```go
err := opencode.Update()
```

- **Purpose**: Update OpenCode installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: OpenCode updates are typically handled by the system package manager

## Configuration Methods

### ConfigureOptions Structure

```go
type ConfigureOptions struct {
    Theme string  // Theme name (empty = "default", or specify custom theme)
}
```

The `ConfigureOptions` struct allows customization of the configuration generation process:

- **Theme**: Specifies which theme to use
  - Empty string or "default": Uses devgita's Gruvbox theme, copies theme file
  - Custom name (e.g., "my-theme"): Uses custom theme name, skips theme file copy

### Configuration Paths

- **Source Template**: `paths.Paths.App.Configs.OpenCode/opencode.json.tmpl`
- **Generated Config**: `paths.Paths.Config.OpenCode/opencode.json`
- **Source Theme**: `paths.Paths.App.Configs.OpenCode/themes/devgita-gruvbox.json`
- **Copied Theme**: `paths.Paths.Config.OpenCode/themes/devgita-gruvbox.json`
- **Marker file**: `opencode.json` in user's OpenCode config directory

### ForceConfigure()

```go
// Use default theme (Gruvbox)
options := opencode.ConfigureOptions{Theme: ""}
err := opencode.ForceConfigure(options)

// Use custom theme
options := opencode.ConfigureOptions{Theme: "my-custom-theme"}
err := opencode.ForceConfigure(options)
```

- **Purpose**: Apply OpenCode configuration regardless of existing files
- **Behavior**:
  1. Remove existing config directory completely
  2. Create fresh config directory
  3. Load GlobalConfig
  4. Determine theme (options.Theme or "default")
  5. Generate `opencode.json` from template with theme variable
  6. If theme == "default", create themes directory and copy Gruvbox theme file
  7. Register installation in GlobalConfig
  8. Save GlobalConfig
- **Use case**: Reset to devgita defaults, apply config updates, change themes

**Theme Handling:**
- **Default theme** (`Theme: ""` or `Theme: "default"`):
  - Sets `"theme": "default"` in config
  - Creates `~/.config/opencode/themes/` directory
  - Copies `devgita-gruvbox.json` to themes directory
- **Custom theme** (`Theme: "my-theme"`):
  - Sets `"theme": "my-theme"` in config
  - Skips theme file copying (user provides their own)

### SoftConfigure()

```go
options := opencode.ConfigureOptions{Theme: ""}
err := opencode.SoftConfigure(options)
```

- **Purpose**: Apply OpenCode configuration only if not already configured
- **Behavior**: 
  - Checks for `opencode.json` marker file in user config directory
  - If exists, returns nil (preserves user customizations)
  - If not exists, calls `ForceConfigure(options)`
- **Use case**: Initial setup that preserves user customizations

## Execution Methods

### ExecuteCommand()

```go
err := opencode.ExecuteCommand("--version")
err := opencode.ExecuteCommand("--help")
err := opencode.ExecuteCommand("file.go")
```

- **Purpose**: Execute OpenCode commands with provided arguments
- **Parameters**: Variable arguments passed directly to opencode binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### OpenCode-Specific Operations

The OpenCode CLI provides various editor operations:

#### Basic Operations

```bash
# Show version
opencode --version

# Show help
opencode --help

# Open file
opencode file.go

# Open directory
opencode /path/to/project

# Open multiple files
opencode file1.go file2.go
```

#### Editor Features

```bash
# Start with AI assistant
opencode --ai-assist

# Open with specific theme
opencode --theme gruvbox file.go

# Verbose mode
opencode --verbose
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure(ConfigureOptions{})`
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure(ConfigureOptions{})`
3. **Custom Theme Setup**: `New()` → `SoftInstall()` → `ForceConfigure(ConfigureOptions{Theme: "custom"})`
4. **Update Configuration**: `New()` → `SoftInstall()` → `ForceConfigure(ConfigureOptions{})`
5. **OpenCode Operations**: `New()` → `ExecuteCommand()` with specific arguments

## Constants and Paths

### Relevant Constants

- `constants.OpenCode`: Package name ("opencode") for installation and configuration
- Used by all installation methods for consistent package reference

### Configuration Paths

- `paths.Paths.App.Configs.OpenCode`: Source directory for OpenCode configuration templates
- `paths.Paths.Config.OpenCode`: Target directory for user's OpenCode configuration
- Configuration uses template system with variable substitution

## Implementation Notes

- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since OpenCode uninstall is not supported
- **Configuration Strategy**: Uses template-based generation with `GenerateFromTemplate()` instead of simple file copying
- **Theme Management**: Conditionally copies theme file only when using default theme
- **Directory Management**: Removes and recreates config directory in `ForceConfigure()` for clean state
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **GlobalConfig Integration**: Tracks OpenCode as a "package" installation
- **Update Method**: Not implemented as OpenCode updates should be handled by system package managers

## Configuration Structure

The OpenCode configuration (`opencode.json`) includes:

### Template Structure

```json
{
  "version": "1.0.0",
  "theme": "{{ .Theme }}",
  "settings": {
    "fontSize": 14,
    "fontFamily": "JetBrains Mono",
    "aiProvider": "openai",
    "autoSave": true,
    "autoFormat": true,
    "lineNumbers": true,
    "wordWrap": true,
    "tabSize": 2
  },
  "keybindings": {
    "save": "Ctrl+S",
    "quit": "Ctrl+Q",
    "aiAssist": "Ctrl+Space",
    "search": "Ctrl+F",
    "replace": "Ctrl+H"
  },
  "editor": {
    "cursorStyle": "line",
    "cursorBlinking": "smooth",
    "renderWhitespace": "selection",
    "minimap": {
      "enabled": true,
      "side": "right"
    }
  }
}
```

### Template Variables

- `{{ .Theme }}`: Replaced with theme name from `ConfigureOptions.Theme`
  - Default: "default" (when `Theme: ""`)
  - Custom: User-specified value (e.g., "my-theme")

### Gruvbox Theme Structure

```json
{
  "name": "Devgita Gruvbox",
  "type": "dark",
  "colors": {
    "background": "#282828",
    "foreground": "#ebdbb2",
    "cursor": "#fe8019",
    "selection": "#504945",
    "comment": "#928374",
    // ... ANSI colors
  },
  "syntax": {
    "keyword": "#fb4934",
    "string": "#b8bb26",
    "number": "#d3869b",
    "function": "#83a598",
    // ... syntax highlighting
  },
  "ui": {
    "lineHighlight": "#3c3836",
    "statusBar": "#3c3836",
    // ... UI elements
  }
}
```

## Usage Examples

### Default Theme Installation

```go
opencode := opencode.New()

// Install OpenCode
if err := opencode.SoftInstall(); err != nil {
    return err
}

// Configure with default Gruvbox theme
options := opencode.ConfigureOptions{Theme: ""}
if err := opencode.SoftConfigure(options); err != nil {
    return err
}
```

**Result:**
- Config generated with `"theme": "default"`
- Gruvbox theme file copied to `~/.config/opencode/themes/devgita-gruvbox.json`

### Custom Theme Installation

```go
opencode := opencode.New()

// Install OpenCode
if err := opencode.SoftInstall(); err != nil {
    return err
}

// Configure with custom theme
options := opencode.ConfigureOptions{Theme: "solarized-dark"}
if err := opencode.ForceConfigure(options); err != nil {
    return err
}
```

**Result:**
- Config generated with `"theme": "solarized-dark"`
- Theme file copying skipped (user provides their own theme file)

### Command Execution

```go
opencode := opencode.New()

// Check version
if err := opencode.ExecuteCommand("--version"); err != nil {
    return err
}

// Open file
if err := opencode.ExecuteCommand("main.go"); err != nil {
    return err
}
```

## Testing Patterns

### Test Structure

```go
func init() {
    testutil.InitLogger()
}

func TestForceConfigure(t *testing.T) {
    tc := testutil.SetupCompleteTest(t)
    defer tc.Cleanup()
    
    // Setup OpenCode-specific paths
    appConfigDir := filepath.Join(tc.AppDir, "configs", "opencode")
    userConfigDir := filepath.Join(tc.ConfigDir, "opencode")
    
    paths.Paths.App.Configs.OpenCode = appConfigDir
    paths.Paths.Config.OpenCode = userConfigDir
    
    // Create template and theme files
    // ... test logic
    
    testutil.VerifyNoRealCommands(t, tc.MockApp.Base)
}
```

### Test Coverage

- ✅ `TestNew` - Constructor validation
- ✅ `TestInstall` - Package installation
- ✅ `TestSoftInstall` - Conditional installation
- ✅ `TestForceInstall` - Error handling (Uninstall not supported)
- ✅ `TestUninstall` - Returns error
- ✅ `TestUpdate` - Returns error
- ✅ `TestForceConfigure` - Template generation with default theme
- ✅ `TestForceConfigure_CustomTheme` - Template generation with custom theme
- ✅ `TestSoftConfigure` - Marker file preservation
- ✅ `TestExecuteCommand` - Command execution
- ✅ `TestExecuteCommand_Error` - Error handling

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Template Not Found**: Verify `configs/opencode/opencode.json.tmpl` exists in devgita repository
3. **Theme Not Found**: Verify `configs/opencode/themes/devgita-gruvbox.json` exists (only for default theme)
4. **Configuration Not Applied**: Check file permissions in config directory
5. **Commands Don't Work**: Verify OpenCode is installed and accessible in PATH
6. **GlobalConfig Errors**: Check `~/.config/devgita/global_config.yaml` exists and is valid

### Platform Considerations

- **macOS**: Installed via Homebrew package manager
- **Linux**: Installed via apt package manager
- **Configuration Location**: Cross-platform config directory handling via `paths` package
- **Theme Support**: Works on all platforms with proper color terminal support

### Configuration Validation

**After installation, verify configuration:**

```bash
# Check config file exists
ls -la ~/.config/opencode/opencode.json

# Check theme file (only if using default theme)
ls -la ~/.config/opencode/themes/devgita-gruvbox.json

# View generated config
cat ~/.config/opencode/opencode.json

# Check OpenCode version
opencode --version
```

## Key Features

### Template-Based Configuration
- Uses Go templates for flexible config generation
- Single source of truth in `configs/opencode/opencode.json.tmpl`
- Variable substitution for customization
- Clean regeneration ensures consistency

### Theme Customization
- **Default theme**: Devgita Gruvbox (dark theme optimized for development)
- **Custom themes**: User can specify any theme name
- **Conditional copying**: Only copies theme file when using default
- **Theme flexibility**: Users can provide their own theme files

### Clean State Management
- `ForceConfigure()` removes old config before generating new
- Creates fresh directory structure
- Ensures no stale configuration files
- Proper directory permissions (0755)

### GlobalConfig Integration
- Tracks OpenCode installation as "package" type
- Differentiates devgita-installed vs pre-existing
- Enables safe future uninstall operations
- Maintains installation state

## Integration with Devgita

OpenCode integrates with devgita's terminal category:

- **Installation**: Part of terminal tools or standalone installation
- **Configuration**: Template-based with theme support
- **Tracking**: Registered in GlobalConfig as package
- **Updates**: Managed through system package manager
- **Category**: Terminal/Editor tools

## External References

- **Configuration Template**: `configs/opencode/opencode.json.tmpl`
- **Theme File**: `configs/opencode/themes/devgita-gruvbox.json`
- **Testing Guide**: `docs/guides/testing-patterns.md`
- **Project Overview**: `docs/project-overview.md`

## Key Design Decisions

### Why Template-Based Configuration?

1. **Flexibility**: Easy to add new variables and customize
2. **Maintainability**: Single source of truth for default config
3. **Consistency**: Same config generation pattern as other devgita apps
4. **Version Control**: Templates are easier to track than binary configs

### Why Conditional Theme Copying?

1. **Default convenience**: New users get working theme immediately
2. **Custom flexibility**: Advanced users can use their own themes
3. **Efficiency**: Don't copy unnecessary files for custom themes
4. **Clear intent**: Empty theme = default, specified theme = custom

### Why Not Support Uninstall?

1. **System integration**: Editors often have system-wide integrations
2. **User data**: May have user-created files and projects
3. **Package manager**: Better handled by system package manager
4. **Safety**: Prevents accidental removal of important tools

This module provides a robust foundation for OpenCode editor integration within the devgita development environment ecosystem, with flexible theme management and template-based configuration for optimal developer experience.
