# Languages Module Documentation

## Overview

The Languages module coordinates programming language installation and management with devgita integration. It follows the standardized tooling coordinator pattern while providing language-specific operations for runtime management via Mise or native package managers.

## Module Purpose

This coordinator module manages the installation of programming languages across development environments. It provides:

- **Interactive language selection** via TUI with installed language detection
- **Dual installation strategies**: Mise runtime manager (Node, Go, Python, Rust) or native package managers (PHP)
- **GlobalConfig tracking** of installed languages for proper state management
- **Conflict prevention** by filtering already-installed languages from selection

## Lifecycle Summary

1. **Selection**: Present TUI for language selection, filtering out already-installed languages
2. **Installation**: Install selected languages via Mise or native package manager based on configuration
3. **Tracking**: Record installations in GlobalConfig for state management and future reference

## Architecture

### Coordinator Pattern

The Languages module acts as a **coordinator** (not an app), orchestrating language installations:

```
DevLanguages (coordinator)
├── LanguageConfig (configuration)
├── Mise (runtime manager app)
└── Native Package Manager (via Command interface)
```

Unlike individual apps in `internal/apps/`, coordinators:
- Manage multiple installation strategies
- Provide interactive selection UIs
- Coordinate between different apps and tools
- Track state in GlobalConfig

### Comparison with Terminal Coordinator

Similar to `internal/tooling/terminal/terminal.go`:

| Aspect | Terminal | Languages |
|--------|----------|-----------|
| Pattern | Coordinator | Coordinator |
| Selection | N/A (installs all) | Interactive TUI |
| Apps | Multiple (neovim, tmux, etc.) | Multiple (mise, native) |
| Strategy | Structured configs array | LanguageConfig array |
| Tracking | Terminal tools | Dev languages |

## Exported Functions

| Function | Purpose | Behavior |
|----------|---------|----------|
| `New()` | Factory method | Creates DevLanguages instance with command executors |
| `GetSelectionOptions()` | Get TUI options | Returns TUI menu options including control items and language names |
| `ChooseLanguages(ctx)` | Interactive selection | Presents TUI, filters installed languages, returns context |
| `InstallChosen(ctx)` | Install selected | Installs languages from context via Mise/native |

## Language Configuration

### LanguageConfig Structure

```go
type LanguageConfig struct {
    DisplayName string  // Human-readable name (e.g., "Node")
    Name        string  // Package/runtime name (e.g., "node")
    Version     string  // Version spec (e.g., "lts", "latest")
    UseMise     bool    // true = mise, false = native package manager
}
```

### Available Languages

| Language | Name | Version | Installation Method |
|----------|------|---------|---------------------|
| Node | node | lts | Mise |
| Go | go | latest | Mise |
| Python | python | latest | Mise |
| PHP | php | (none) | Native package manager |
| Rust | rust | latest | Mise |

### Language Specifications

Languages are tracked in GlobalConfig with these specifications:

- **Mise languages**: `name@version` (e.g., `node@lts`, `go@latest`)
- **Native languages**: `name` only (e.g., `php`)

## Installation Flow

### 1. Interactive Selection

```go
l := languages.New()
ctx, err := l.ChooseLanguages(ctx)
```

**Process**:
1. Load GlobalConfig to identify installed languages
2. Filter installed languages from available options
3. Display warning if languages are already installed
4. Present TUI multi-select with filtered options
5. Store selections in context

**Example Output**:
```
Already installed languages (skipped from selection): Node, Python
? Select programming languages to install
  ▸ All
    None
    Done
    Go
    PHP
    Rust
```

### 2. Install Selected Languages

```go
l.InstallChosen(ctx)
```

**Process** for each selected language:

#### Mise Installation (Node, Go, Python, Rust)

```
1. Ensure Mise is installed (SoftInstall)
2. Configure Mise shell integration (SoftConfigure)
3. Install language globally via Mise.UseGlobal(name, version)
4. Track in GlobalConfig as "name@version"
```

#### Native Installation (PHP)

```
1. Install via package manager (MaybeInstallPackage)
2. Track in GlobalConfig as "name"
```

### 3. State Tracking

```go
gc.AddToInstalled(languageSpec, "dev_language")
gc.Save()
```

Tracked in `~/.config/devgita/global_config.yaml`:

```yaml
installed:
  dev_languages:
    - node@lts
    - python@latest
    - php
```

## Implementation Details

### Installed Language Detection

```go
func (dl *DevLanguages) getInstalledLanguages(gc *config.GlobalConfig) []string
```

**Logic**:
1. Iterate through language configurations
2. Generate language spec (mise: `name@version`, native: `name`)
3. Check in both `installed.dev_languages` and `already_installed.dev_languages`
4. Return display names of installed languages

**Purpose**: Filter already-installed languages from TUI selection to prevent conflicts

### Installation Strategy Selection

```go
if langCfg.UseMise {
    err = dl.installWithMise(langCfg)
} else {
    err = dl.installNative(langCfg)
}
```

**Mise Installation**:
- Ensures Mise is installed and configured
- Calls `mise use --global name@version`
- Managed by Mise runtime manager (see `docs/apps/mise.md`)

**Native Installation**:
- Uses platform package manager (Homebrew/apt)
- Calls `MaybeInstallPackage(name)`
- Simpler approach for languages without version management needs

### Error Handling

```go
if err != nil {
    utils.PrintError("Error: Unable to install %s: %v", name, err)
    logger.L().Errorw("Language installation failed", "language", name, "error", err)
    return // Non-fatal: continue with other languages
}
```

**Philosophy**:
- Language installation failures are **non-fatal**
- Print error to user, log details, continue with remaining languages
- GlobalConfig tracking failures are **non-fatal** but logged as warnings

## Usage Examples

### Standard Installation Flow

```go
// In cmd/install.go
func installLanguages(ctx context.Context, onlySet, skipSet map[string]bool) {
    if shouldInstall("languages", onlySet, skipSet) {
        l := languages.New()
        ctx, err := l.ChooseLanguages(ctx)
        utils.MaybeExitWithError(err)
        l.InstallChosen(ctx)
    }
}
```

### Adding a New Language

To add a new language (e.g., Ruby):

1. **Add to configurations**:
```go
func (dl *DevLanguages) getLanguageConfigs() []LanguageConfig {
    return []LanguageConfig{
        // ... existing languages ...
        {DisplayName: "Ruby", Name: "ruby", Version: "latest", UseMise: true},
    }
}
```

2. **Add to available list**:
```go
func (dl *DevLanguages) GetSelectionOptions() []string {
    return []string{
        "All", "None", "Done",
        "Node", "Go", "PHP", "Python", "Rust", "Ruby",
    }
}
```

That's it! The coordinator handles the rest automatically.

## GlobalConfig Integration

### Schema

```yaml
installed:
  dev_languages:
    - node@lts       # Mise-managed
    - go@latest      # Mise-managed
    - python@latest  # Mise-managed
    - php            # Native package manager
    - rust@latest    # Mise-managed

already_installed:
  dev_languages:
    - node@16.0.0    # Pre-existing, tracked for awareness
```

### Methods Used

- `gc.AddToInstalled(languageSpec, "dev_language")` - Track devgita installation
- `gc.IsInstalledByDevgita(languageSpec, "dev_language")` - Check if installed
- `gc.IsAlreadyInstalled(languageSpec, "dev_language")` - Check if pre-existing

## Testing Patterns

### Test Structure

```go
func TestInstallChosen_WithSelections(t *testing.T) {
    mockApp := testutil.NewMockApp()
    mockApp.Base.SetExecCommandResult("", "", nil)
    
    dl := &DevLanguages{Cmd: mockApp.Cmd, Base: mockApp.Base}
    
    ctx := context.Background()
    config := config.ContextConfig{SelectedLanguages: []string{"PHP"}}
    ctx = config.WithConfig(ctx, config)
    
    dl.InstallChosen(ctx)
    
    testutil.VerifyNoRealCommands(t, mockApp.Base)
}
```

### Test Coverage

- ✅ Configuration structure validation
- ✅ Language spec generation (mise vs native)
- ✅ Installed language detection
- ✅ Selection filtering logic
- ✅ Installation orchestration
- ✅ GlobalConfig tracking

## Comparison with Individual Apps

| Aspect | Individual App (e.g., Neovim) | Languages Coordinator |
|--------|-------------------------------|----------------------|
| Location | `internal/apps/neovim/` | `internal/tooling/languages/` |
| Interface | Implements standard app interface | Custom coordinator interface |
| Installation | Single package | Multiple languages |
| Configuration | Config files (e.g., init.lua) | GlobalConfig tracking only |
| Selection | N/A | Interactive TUI |
| Structure | Singular focus | Multi-strategy orchestration |

## Integration with Other Modules

### Mise App

Languages coordinator **depends on** Mise app:

```go
m := mise.New()
m.SoftInstall()
m.SoftConfigure()
m.UseGlobal(name, version)
```

See `docs/apps/mise.md` for Mise-specific documentation.

### Command Interface

Uses `Command` and `BaseCommandExecutor` interfaces:

```go
type DevLanguages struct {
    Cmd  cmd.Command              // Package management
    Base cmd.BaseCommandExecutor  // Command execution
}
```

### Context Configuration

Leverages context for passing selections:

```go
// Store selections
initialConfig := config.ContextConfig{SelectedLanguages: selections}
ctx = config.WithConfig(ctx, initialConfig)

// Retrieve selections
selections, ok := config.GetConfig(ctx)
```

## Platform Considerations

### macOS (Homebrew)

- **Mise**: Installed via `brew install mise`
- **Native PHP**: Installed via `brew install php`
- **Runtime management**: Mise handles version switching

### Linux (Debian/Ubuntu)

- **Mise**: Installed via apt (mise repository)
- **Native PHP**: Installed via `apt install php`
- **Runtime management**: Mise handles version switching

## Future Enhancements

### Potential Additions

1. **Version selection**: Allow users to specify versions during TUI selection
2. **Uninstall support**: Remove languages and clean up GlobalConfig
3. **Update support**: Update language versions via Mise
4. **Alternative strategies**: Support for asdf, nvm, etc.
5. **Virtual environments**: Python venv, Node.js nvm integration
6. **Language-specific tooling**: Install related tools (npm, pip, etc.)

### Planned Languages

- Ruby (via Mise)
- Java (via Mise)
- Elixir (via Mise)
- Zig (via Mise)

## Troubleshooting

### Common Issues

1. **Mise not installed**: Ensure Mise installation succeeds before language installation
2. **Shell integration missing**: Run `mise activate zsh` or restart shell
3. **Version conflicts**: Check existing language installations with `mise list`
4. **GlobalConfig errors**: Non-fatal but logged; check `~/.config/devgita/global_config.yaml`

### Debugging

```bash
# Check installed languages
mise list

# Verify global versions
mise current

# Check devgita global config
cat ~/.config/devgita/global_config.yaml

# Run with verbose logging
dg install --only languages --verbose
```

## Key Design Decisions

### Why Coordinator Pattern?

1. **Complexity**: Managing multiple languages with different strategies
2. **Selection**: Interactive TUI requires coordination logic
3. **Dual strategies**: Mise vs native requires orchestration
4. **State management**: GlobalConfig tracking spans multiple installations

### Why Mise for Most Languages?

1. **Version management**: Easy switching between versions
2. **Per-project control**: `.mise.toml` for project-specific versions
3. **Cross-platform**: Works on macOS and Linux
4. **Active development**: Modern alternative to nvm, rbenv, etc.

### Why Native for PHP?

1. **System integration**: PHP often needs system-level configuration
2. **Extensions**: Native packages include common extensions
3. **Version management**: Less critical for PHP in most use cases
4. **Simplicity**: Avoids additional complexity for single version

## External References

- **Mise Documentation**: https://mise.jdx.dev/
- **Mise App Docs**: `docs/apps/mise.md`
- **Terminal Coordinator**: `internal/tooling/terminal/terminal.go`
- **GlobalConfig**: `docs/project-overview.md#configuration-management`
- **Testing Patterns**: `docs/guides/testing-patterns.md`

This coordinator provides a robust, extensible foundation for managing programming language installations within the devgita development environment ecosystem.
