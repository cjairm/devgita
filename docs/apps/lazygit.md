# LazyGit Module Documentation

## Overview

The LazyGit module provides installation and command execution management for lazygit terminal UI with devgita integration. It follows the standardized devgita app interface while providing lazygit-specific operations for Git repository management, branch operations, commit history, and staging through an interactive terminal interface.

## App Purpose

LazyGit is a simple terminal UI for git commands, written in Go with the gocui library. It provides an interactive interface to manage Git repositories, branches, commits, staging operations, and merge conflicts, making Git workflow management more accessible and efficient directly from the terminal. This module ensures lazygit is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems and provides high-level operations for launching the interactive TUI.

## Lifecycle Summary

1. **Installation**: Install lazygit package via platform package managers (Homebrew/apt)
2. **Configuration**: LazyGit uses optional user-specific configuration files (no default configuration applied by devgita)
3. **Execution**: Provide high-level lazygit operations for launching the interactive TUI and managing Git repositories

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new LazyGit instance with platform-specific commands         |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install lazygit                           |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute lazygit           | Runs lazygit with provided arguments                                 |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
lazygit := lazygit.New()
err := lazygit.Install()
```

- **Purpose**: Standard lazygit installation
- **Behavior**: Uses `InstallPackage()` to install lazygit package
- **Use case**: Initial lazygit installation or explicit reinstall

### ForceInstall()

```go
lazygit := lazygit.New()
err := lazygit.ForceInstall()
```

- **Purpose**: Force lazygit installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh lazygit installation or fix corrupted installation

### SoftInstall()

```go
lazygit := lazygit.New()
err := lazygit.SoftInstall()
```

- **Purpose**: Install lazygit only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing lazygit installations

### Uninstall()

```go
err := lazygit.Uninstall()
```

- **Purpose**: Remove lazygit installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Git tools are typically managed at the system level

### Update()

```go
err := lazygit.Update()
```

- **Purpose**: Update lazygit installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: LazyGit updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := lazygit.ForceConfigure()
err := lazygit.SoftConfigure()
```

- **Purpose**: Apply lazygit configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: LazyGit configuration is optional and user-specific. Configuration files are typically located at `~/.config/lazygit/config.yml` and are managed by users based on their preferences.

## Execution Methods

### ExecuteCommand()

```go
err := lazygit.ExecuteCommand()                    // Launch TUI
err := lazygit.ExecuteCommand("--version")         // Show version
err := lazygit.ExecuteCommand("--config")          // Show config path
err := lazygit.ExecuteCommand("--help")            // Show help
```

- **Purpose**: Execute lazygit commands with provided arguments
- **Parameters**: Variable arguments passed directly to lazygit binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### LazyGit-Specific Operations

The lazygit CLI provides interactive Git repository management capabilities:

#### Launch Interactive TUI

```bash
# Launch lazygit (default behavior)
lazygit

# The TUI provides:
# - Status view with staged/unstaged changes
# - Branch management (create, checkout, delete, merge, rebase)
# - Commit operations (commit, amend, fixup, squash)
# - Log viewer with commit history
# - Stash management
# - Merge conflict resolution
# - File staging and unstaging
# - Diff viewer
# - Remote operations (push, pull, fetch)
# - Submodule management
# - Cherry-picking commits
# - Interactive rebase
```

#### Command-Line Options

```bash
# Show version information
lazygit --version

# Show configuration file path
lazygit --config

# Display help information
lazygit --help

# Use custom config file
lazygit --use-config-file /path/to/config.yml

# Start in specific directory
lazygit --path /path/to/repo

# Debug mode
lazygit --debug

# Filter by path
lazygit --filter path/to/filter
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **Launch TUI**: `New()` → `SoftInstall()` → `ExecuteCommand()`
4. **Version Check**: `New()` → `ExecuteCommand("--version")`

## Constants and Paths

### Relevant Constants

- **Package name**: Referenced via `constants.LazyGit` (typically "lazygit")
- Used by all installation methods for consistent package reference

### Configuration Approach

- **Optional configuration**: LazyGit uses optional configuration files
- **Default location**: `~/.config/lazygit/config.yml` (user-managed)
- **No default config**: Devgita does not apply default configuration for lazygit
- **User customization**: Users can create their own configuration based on preferences

### Configuration Options

While devgita doesn't apply default configuration, users can customize lazygit via `~/.config/lazygit/config.yml`:

```yaml
# Example lazygit configuration
gui:
  theme:
    activeBorderColor:
      - green
      - bold
    inactiveBorderColor:
      - white
  showFileTree: true
  showListFooter: true
  showRandomTip: true
  showBranchCommitHash: false
  showBottomLine: true
  showCommandLog: true
  commandLogSize: 8
  splitDiff: auto

git:
  paging:
    colorArg: always
    pager: delta --dark --paging=never
  commit:
    signOff: false
  merging:
    manualCommit: false
    args: ''
  log:
    order: 'topo-order'
    showGraph: 'when-maximised'
  skipHookPrefix: WIP
  autoFetch: true
  autoRefresh: true
  branchLogCmd: 'git log --graph --color=always --abbrev-commit --decorate --date=relative --pretty=medium {{branchName}} --'

os:
  editCommand: ''
  editCommandTemplate: ''
  openCommand: ''

refresher:
  refreshInterval: 10
  fetchInterval: 60

update:
  method: prompt
  days: 14

confirmOnQuit: false
quitOnTopLevelReturn: false

keybinding:
  universal:
    quit: 'q'
    quit-alt1: '<c-c>'
    return: '<esc>'
    quitWithoutChangingDirectory: 'Q'
    togglePanel: '<tab>'
    prevItem: '<up>'
    nextItem: '<down>'
    prevItem-alt: 'k'
    nextItem-alt: 'j'
```

## Implementation Notes

- **CLI Tool Nature**: LazyGit is an interactive terminal UI tool without complex configuration requirements
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since lazygit uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since lazygit uses optional user-specific configuration
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as lazygit updates should be handled by system package managers

## Usage Examples

### Basic LazyGit Operations

```go
lazygit := lazygit.New()

// Install lazygit
err := lazygit.SoftInstall()
if err != nil {
    return err
}

// Launch interactive TUI
err = lazygit.ExecuteCommand()

// Check version
err = lazygit.ExecuteCommand("--version")
```

### Advanced Usage

```go
// Check configuration file location
err := lazygit.ExecuteCommand("--config")

// Show help
err = lazygit.ExecuteCommand("--help")

// Start in specific directory
err = lazygit.ExecuteCommand("--path", "/path/to/repo")
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **TUI Won't Launch**: Verify Git is installed and accessible
3. **Repository Not Found**: Ensure you're in a Git repository directory
4. **Commands Don't Work**: Verify lazygit is installed and accessible in PATH
5. **Git Not Found**: Install Git via devgita's terminal category

### Platform Considerations

- **macOS**: Installed via Homebrew package manager, requires Git
- **Linux**: Installed via apt package manager, requires Git
- **Git Dependency**: LazyGit requires Git to be installed and configured
- **Terminal Support**: Works best with modern terminal emulators that support colors

### Prerequisites

Before using lazygit, ensure:
- Git is installed and configured
- You're in a Git repository directory (or use `--path` flag)
- Git user name and email are configured
- Terminal supports colors and special characters for optimal TUI experience

### Git Integration

LazyGit integrates with Git via:
- **Git commands**: Executes git commands under the hood
- **Git config**: Respects git configuration settings
- **Git hooks**: Supports pre-commit and commit-msg hooks
- **Git credentials**: Uses configured credential helpers

## Key Features

### Branch Management
- Visual branch tree with relationships
- Create, checkout, delete, rename branches
- Merge and rebase operations
- Branch comparison and diff viewing

### Commit Operations
- Interactive staging and unstaging
- Commit, amend, fixup, squash operations
- Commit message editing
- Commit history navigation

### Merge Conflict Resolution
- Visual conflict markers
- Side-by-side diff view
- Easy conflict resolution
- Automatic staging after resolution

### Stash Management
- Create and apply stashes
- Stash pop and drop operations
- View stash contents
- Named stash support

### Remote Operations
- Push and pull with visual feedback
- Fetch from remotes
- Force push with lease
- Remote branch tracking

## External References

- **LazyGit Repository**: https://github.com/jesseduffield/lazygit
- **Configuration Guide**: https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md
- **Keybindings**: https://github.com/jesseduffield/lazygit/blob/master/docs/keybindings/Keybindings_en.md
- **Custom Commands**: https://github.com/jesseduffield/lazygit/blob/master/docs/Custom_Commands.md
- **Git Documentation**: https://git-scm.com/doc

## Integration with Devgita

LazyGit integrates with devgita's terminal category:

- **Installation**: Installed as part of terminal tools setup
- **Configuration**: User-managed configuration (no default applied)
- **Usage**: Available system-wide after installation
- **Updates**: Managed through system package manager
- **Dependencies**: Requires Git to be installed (also in terminal category)

This module provides essential Git repository management capabilities through an intuitive terminal interface, significantly improving Git workflow efficiency within the devgita ecosystem.
