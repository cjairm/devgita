# Ncurses Module Documentation

## Overview

The Ncurses module provides installation and management for the ncurses terminal UI library with devgita integration. It follows the standardized devgita app interface while providing ncurses-specific operations for system library installation and management.

## App Purpose

Ncurses (new curses) is a programming library that provides an API which allows programmers to write text-based user interfaces in a terminal-independent manner. It is a toolkit for developing GUI-like application software that runs under a terminal emulator. This module ensures ncurses development libraries are properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems for building terminal-based applications with advanced UI capabilities.

## Lifecycle Summary

1. **Installation**: Install ncurses package via platform package managers (Homebrew/apt)
2. **Configuration**: Ncurses typically doesn't require separate configuration files - it's a system library used by other applications
3. **Execution**: Provides placeholder operations for interface compliance (ncurses is a library, not a CLI tool)

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Ncurses instance with platform-specific commands         |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install ncurses                           |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute ncurses commands  | **Library only** - no CLI commands available                         |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
ncurses := ncurses.New()
err := ncurses.Install()
```

- **Purpose**: Standard ncurses installation
- **Behavior**: Uses `InstallPackage()` to install ncurses package
- **Platform differences**:
  - macOS: Installs `ncurses` via Homebrew
  - Debian/Ubuntu: Installs `libncurses-dev` or `libncurses5-dev` via apt
- **Use case**: Initial ncurses installation or explicit reinstall

### ForceInstall()

```go
ncurses := ncurses.New()
err := ncurses.ForceInstall()
```

- **Purpose**: Force ncurses installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh ncurses installation or fix corrupted installation
- **Note**: Will fail since `Uninstall()` is not supported

### SoftInstall()

```go
ncurses := ncurses.New()
err := ncurses.SoftInstall()
```

- **Purpose**: Install ncurses only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing ncurses installations

### Uninstall()

```go
err := ncurses.Uninstall()
```

- **Purpose**: Remove ncurses installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: System libraries are typically managed at the system level and should not be uninstalled via devgita

### Update()

```go
err := ncurses.Update()
```

- **Purpose**: Update ncurses installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Ncurses updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := ncurses.ForceConfigure()
err := ncurses.SoftConfigure()
```

- **Purpose**: Apply ncurses configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: Ncurses is a system library without separate configuration files. It's configured at compile-time by applications that use it.

## Execution Methods

### ExecuteCommand()

```go
err := ncurses.ExecuteCommand(args...)
```

- **Purpose**: Execute ncurses commands with provided arguments
- **Behavior**: Attempts to execute commands but ncurses has no direct CLI interface
- **Note**: Ncurses is a library, not a standalone CLI tool. This method exists for interface compliance but may not have practical use.
- **Parameters**: Variable arguments that would be passed to a hypothetical ncurses command
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **Terminal Category**: Installed as part of core system libraries in the terminal tooling setup

## Constants and Paths

### Relevant Constants

- **Package name**: Must be defined in `pkg/constants/constants.go`
  ```go
  Ncurses = "ncurses"
  ```
- Used by all installation methods for consistent package reference

### Platform-Specific Package Names

- **macOS**: `ncurses` (via Homebrew)
- **Debian/Ubuntu**: `libncurses-dev` or `libncurses5-dev` (via apt)
- **Note**: Platform abstraction handled by the command factory pattern

### Configuration Approach

- **No traditional config files**: Ncurses doesn't use runtime configuration files
- **Compile-time configuration**: Applications that use ncurses configure terminal behavior at compile-time
- **System library**: Installed to standard system library locations by package managers
- **Terminfo database**: Ncurses uses terminfo for terminal capability descriptions

## Implementation Notes

- **System Library Nature**: Unlike typical CLI applications, ncurses is a system library without direct user interaction
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since ncurses uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since ncurses doesn't use config files
- **ExecuteCommand**: Included for interface compliance but has limited practical use for a library
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as ncurses updates should be handled by system package managers

## Usage in Development

### Purpose as System Library

Ncurses is typically installed as a dependency for:
- Building terminal-based user interfaces
- Text editors and viewers (vim, less, htop)
- Terminal multiplexers (tmux, screen)
- System monitoring tools (top, htop, btop)
- Interactive command-line applications
- Terminal-based games and utilities

### Developer Impact

While ncurses itself has no direct CLI commands, it's essential infrastructure for:
- Compiling text-based applications
- Building terminal UI applications from source
- Creating interactive CLI tools
- Developing terminal-based applications with advanced UI features

### Terminal UI Capabilities

Ncurses provides functionality for:
- **Window Management**: Creating and managing multiple windows
- **Color Support**: Handling terminal colors and color pairs
- **Cursor Control**: Moving cursor and managing cursor visibility
- **Keyboard Input**: Reading and processing keyboard input
- **Mouse Support**: Handling mouse events in terminal
- **Text Attributes**: Bold, underline, reverse video, etc.
- **Special Characters**: Box drawing and special symbols

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Missing Development Headers**: On Linux, ensure `libncurses-dev` is installed (not just `libncurses`)
3. **Build Failures**: Many terminal applications fail to build if ncurses is not present
4. **Version Conflicts**: Use system package manager to handle version dependencies
5. **Terminfo Issues**: Ensure terminfo database is properly installed with ncurses

### Platform Considerations

- **macOS**: Installed via Homebrew package manager, system ncurses also available
- **Linux**: Installed via apt package manager, requires dev package for headers
- **System Integration**: Installed to standard system library paths
- **Wide Character Support**: Many systems use ncursesw (wide character version)
- **Compatibility**: Different ncurses versions may have API differences

### Common Build Errors

```bash
# Application build failing
# Error: "ncurses.h: No such file or directory"
# Solution: Install libncurses-dev (Linux) or ncurses (macOS)

# Tmux or vim compilation failing
# Error: "library ncurses not found"
# Solution: Install ncurses development libraries

# Undefined reference to ncurses functions
# Error: "undefined reference to 'initscr'"
# Solution: Link with -lncurses flag during compilation
```

## References

- **Ncurses Homepage**: https://invisible-island.net/ncurses/
- **GNU Ncurses**: https://www.gnu.org/software/ncurses/
- **Ncurses Programming Guide**: https://tldp.org/HOWTO/NCURSES-Programming-HOWTO/
- **Ncurses Manual**: https://invisible-island.net/ncurses/man/ncurses.3x.html
- **Terminfo Database**: https://invisible-island.net/ncurses/man/terminfo.5.html

## Integration with Devgita

Ncurses integrates with devgita's terminal core category:

- **Installation**: Installed as part of core system libraries setup
- **Configuration**: No configuration needed (system library)
- **Usage**: Provides terminal UI functionality to terminal applications
- **Updates**: Managed through system package manager
- **Dependencies**: Required by vim, tmux, htop, and many other terminal tools

### Installation Order

Ncurses should be installed before:
- Vim or Neovim
- Tmux
- Htop, btop, and other system monitoring tools
- Any terminal-based applications

### Related Libraries

Often installed alongside:
- **readline**: For interactive line editing
- **libevent**: For event-driven programming
- **terminfo**: Terminal capability database

## Version Information

### Current Versions

- **Ncurses 6.x**: Current stable series with wide character support
- **Ncurses 5.x**: Older version still in use on some systems
- **ncursesw**: Wide character version (recommended for modern applications)

### Version Checking

```bash
# Check if ncurses is installed (macOS)
brew list ncurses

# Check if ncurses is installed (Linux)
dpkg -l | grep ncurses

# Find ncurses version (pkg-config)
pkg-config --modversion ncurses

# Check for ncursesw (wide character version)
pkg-config --modversion ncursesw
```

## Programming with Ncurses

### Basic Usage Example (C)

```c
#include <ncurses.h>

int main() {
    initscr();              // Initialize ncurses
    printw("Hello World!"); // Print to virtual screen
    refresh();              // Update physical screen
    getch();                // Wait for key press
    endwin();               // Clean up and restore terminal
    return 0;
}

// Compile: gcc -o hello hello.c -lncurses
```

### Common Ncurses Functions

- **initscr()**: Initialize ncurses mode
- **endwin()**: Restore terminal to normal mode
- **printw()**: Print formatted output
- **mvprintw()**: Move cursor and print
- **getch()**: Get character from keyboard
- **refresh()**: Update screen with virtual screen contents
- **clear()**: Clear the screen
- **newwin()**: Create new window
- **wrefresh()**: Refresh specific window

### Language Bindings

Ncurses has bindings for many languages:
- **Python**: curses module (built-in)
- **Ruby**: curses gem
- **Perl**: Curses module
- **Node.js**: blessed, blessed-contrib
- **Rust**: pancurses, ncurses crates

## Terminal Applications Built with Ncurses

### Popular Applications

- **Text Editors**: vim, nano, emacs (terminal mode)
- **Terminal Multiplexers**: tmux, screen
- **System Monitors**: htop, btop, top
- **File Managers**: ranger, midnight commander (mc)
- **Music Players**: cmus, ncmpcpp
- **IRC Clients**: irssi, weechat
- **Email Clients**: mutt, alpine
- **Games**: nethack, angband, cataclysm-dda

## Advanced Features

### Color Support

```c
start_color();
init_pair(1, COLOR_RED, COLOR_BLACK);
attron(COLOR_PAIR(1));
printw("Red text on black background");
attroff(COLOR_PAIR(1));
```

### Window Management

```c
WINDOW *win = newwin(height, width, starty, startx);
box(win, 0, 0);
wrefresh(win);
```

### Mouse Support

```c
mousemask(ALL_MOUSE_EVENTS, NULL);
MEVENT event;
if (getmouse(&event) == OK) {
    // Handle mouse event
}
```

This module provides essential terminal UI library support for building and running text-based applications within the devgita ecosystem, ensuring that terminal applications requiring advanced UI capabilities can be compiled and executed successfully.
