# Readline Module Documentation

## Overview

The Readline module provides installation and management for GNU Readline library with devgita integration. It follows the standardized devgita app interface while providing support for this essential line-editing and history library used by many interactive command-line programs.

## App Purpose

GNU Readline is a software library that provides line-editing and history capabilities for interactive programs with a command-line interface. It allows users to edit text as they type, recall previous commands from history, perform text searches, and provides customizable key bindings. Readline is used by many popular command-line programs including bash, python, psql, mysql, and gdb to provide a consistent and powerful editing interface. This module ensures readline is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems, making it available as a system library for other programs to use.

## Lifecycle Summary

1. **Installation**: Install readline package via platform package managers (Homebrew/apt)
2. **Configuration**: readline is a library that doesn't require separate configuration files in devgita - user configuration is handled via `~/.inputrc`
3. **Execution**: Provide interface consistency (readline is primarily a library used by other programs)

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Readline instance with platform-specific commands        |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install readline                          |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute readline          | Interface consistency (limited use for libraries)                    |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
readline := readline.New()
err := readline.Install()
```

- **Purpose**: Standard readline installation
- **Behavior**: Uses `InstallPackage()` to install readline package
- **Use case**: Initial readline installation or explicit reinstall

### ForceInstall()

```go
readline := readline.New()
err := readline.ForceInstall()
```

- **Purpose**: Force readline installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh readline installation or fix corrupted installation

### SoftInstall()

```go
readline := readline.New()
err := readline.SoftInstall()
```

- **Purpose**: Install readline only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing readline installations

### Uninstall()

```go
err := readline.Uninstall()
```

- **Purpose**: Remove readline installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: readline is a critical system library used by many programs and should be managed at the OS level

### Update()

```go
err := readline.Update()
```

- **Purpose**: Update readline installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: readline updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := readline.ForceConfigure()
err := readline.SoftConfigure()
```

- **Purpose**: Apply readline configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: readline is a library that doesn't require separate configuration files in devgita; user-level configuration is handled via `~/.inputrc` which users manage directly

## Execution Methods

### ExecuteCommand()

```go
err := readline.ExecuteCommand("--version")
```

- **Purpose**: Execute readline-related commands
- **Parameters**: Variable arguments passed to readline
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand
- **Note**: Readline is primarily a library without standalone commands; this method exists for interface consistency

### Readline Library Nature

Unlike command-line tools, readline is a **library** (libreadline) that provides line-editing functionality to other programs. It doesn't have standalone command-line tools but is used by many interactive programs.

#### Programs That Use Readline

```bash
# Bash shell (primary user of readline)
bash

# Python interactive interpreter
python
python3

# PostgreSQL interactive terminal
psql

# MySQL command-line client
mysql

# GNU Debugger
gdb

# Arbitrary precision calculator
bc

# R statistical computing environment
R

# Lua interactive interpreter
lua
```

#### Readline Features

Readline provides these capabilities to programs that use it:

**Line Editing**:
- Emacs-style editing (default)
- Vi-style editing mode
- Move cursor with arrow keys
- Cut, copy, and paste text
- Word-based navigation
- Line kill and yank operations

**History Management**:
- Command history recall with up/down arrows
- Reverse incremental search (Ctrl+R)
- History expansion with `!` syntax
- Persistent history across sessions
- History file management

**Completion**:
- Tab completion for commands and paths
- Programmable completion
- Multiple completion matches
- Case-insensitive completion options

**Key Bindings**:
- Customizable key mappings
- Macros for common operations
- Mode-specific bindings
- User-defined key sequences

## User Configuration via .inputrc

While devgita doesn't manage readline configuration, users can customize readline behavior through `~/.inputrc`:

### Basic .inputrc Configuration

```bash
# Enable vi mode
set editing-mode vi

# Set bell style (none, visible, audible)
set bell-style none

# Show completion matches immediately
set show-all-if-ambiguous on

# Case-insensitive completion
set completion-ignore-case on

# Show completion mode indicator
set show-mode-in-prompt on

# Color file completion by type
set colored-stats on

# Show all completion matches with one tab
set menu-complete-display-prefix on
```

### Custom Key Bindings

```bash
# Make Tab cycle through completions
TAB: menu-complete

# Make Shift+Tab cycle backwards
"\e[Z": menu-complete-backward

# Ctrl+Left/Right to move by words
"\e[1;5D": backward-word
"\e[1;5C": forward-word

# Alt+Backspace to delete word
"\e\d": backward-kill-word

# Ctrl+Delete to delete word forward
"\e[3;5~": kill-word
```

### History Configuration

```bash
# Increase history size
$if Bash
    set history-size 10000
$endif

# Ignore duplicate commands in history
$if Bash
    set history-control ignoredups
$endif
```

### Mode-Specific Settings

```bash
# Settings for vi mode
$if mode=vi
    set keymap vi-command
    "gg": beginning-of-history
    "G": end-of-history
    set keymap vi-insert
    "jk": vi-movement-mode
$endif

# Settings for emacs mode
$if mode=emacs
    "\C-p": history-search-backward
    "\C-n": history-search-forward
$endif
```

### Application-Specific Settings

```bash
# Settings specific to Python
$if Python
    set show-all-if-ambiguous on
    TAB: complete
$endif

# Settings specific to MySQL
$if mysql
    set completion-ignore-case on
$endif

# Settings specific to Bash
$if Bash
    set editing-mode emacs
    set completion-query-items 100
$endif
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **Library Installation**: `New()` → `SoftInstall()` (makes readline available to other programs)

## Constants and Paths

### Relevant Constants

- **Package name**: `"readline"` used directly for installation
- Referenced via `constants.Readline` in the codebase
- Used by all installation methods for consistent package reference

### Configuration Approach

- **No devgita configuration**: Readline doesn't require configuration files managed by devgita
- **User configuration**: Users can create `~/.inputrc` for personal readline customization
- **System configuration**: System-level configuration at `/etc/inputrc` (managed by OS)
- **Application-specific**: Some programs have their own readline configuration

### Configuration File Locations

```bash
# User-level configuration (highest priority)
~/.inputrc

# System-level configuration
/etc/inputrc

# Application-specific history files
~/.bash_history      # Bash history
~/.python_history    # Python history
~/.psql_history      # PostgreSQL history
~/.mysql_history     # MySQL history
```

## Implementation Notes

- **Library Nature**: Readline is a shared library (libreadline.so on Linux, libreadline.dylib on macOS) that other programs link against
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since readline uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since readline configuration is user-managed
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as readline updates should be handled by system package managers
- **Critical Dependency**: Many command-line programs depend on readline; removing it could break system functionality

## Usage Examples

### Basic Readline Installation

```go
readline := readline.New()

// Install readline library
err := readline.SoftInstall()
if err != nil {
    return err
}

// Now programs that use readline can be installed
// For example: bash, python, psql, mysql, etc.
```

### Integration with Other Tools

```go
// Install readline first
readline := readline.New()
err := readline.SoftInstall()
if err != nil {
    return err
}

// Then install programs that depend on readline
// These programs will automatically use the readline library
// for line-editing and history capabilities
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Programs Don't Use Readline**: Verify programs are compiled with readline support
3. **Key Bindings Don't Work**: Check `~/.inputrc` syntax and permissions
4. **History Not Saved**: Verify history file permissions and disk space
5. **Completion Issues**: Check readline version and application-specific settings

### Platform Considerations

- **macOS**: Installed via Homebrew, system readline also available but may be outdated
- **Linux**: Installed via apt, usually pre-installed as many programs depend on it
- **Library Versions**: Different readline versions may have different features
- **Compatibility**: Some programs require specific readline versions

### Readline Versions

```bash
# Check readline version (via bash)
bash --version | grep readline

# Check installed readline package
# On Debian/Ubuntu:
dpkg -l | grep readline

# On macOS:
brew list readline

# Check readline library files
ls -l /usr/lib/libreadline* (Linux)
ls -l /usr/local/lib/libreadline* (macOS)
```

### Debugging Readline Configuration

```bash
# Test .inputrc syntax
bash -c 'bind -f ~/.inputrc'

# Show current readline key bindings
bind -p

# Show readline variables
bind -v

# Show readline functions
bind -l

# Reload .inputrc
bind -f ~/.inputrc
```

### Common Readline Commands

These keyboard shortcuts work in programs using readline:

```bash
# Navigation
Ctrl+A          # Move to beginning of line
Ctrl+E          # Move to end of line
Alt+F           # Move forward one word
Alt+B           # Move backward one word

# Editing
Ctrl+K          # Kill (cut) to end of line
Ctrl+U          # Kill to beginning of line
Ctrl+Y          # Yank (paste) killed text
Alt+D           # Delete word forward
Alt+Backspace   # Delete word backward

# History
Ctrl+R          # Reverse search history
Ctrl+P          # Previous command
Ctrl+N          # Next command
!!              # Repeat last command
!$              # Last argument of previous command

# Completion
Tab             # Complete command/filename
Alt+?           # List possible completions
Alt+*           # Insert all completions

# Control
Ctrl+L          # Clear screen
Ctrl+C          # Interrupt current command
Ctrl+D          # Exit shell (EOF)
```

## Integration with Devgita

Readline integrates with devgita's terminal category:

- **Installation**: Installed as part of terminal tools setup (system libraries)
- **Configuration**: User-managed via `~/.inputrc` (not managed by devgita)
- **Usage**: Provides line-editing capabilities to other programs
- **Updates**: Managed through system package manager
- **Dependencies**: Required by bash, python, and many other interactive programs

## External References

- **GNU Readline Homepage**: https://tiswww.case.edu/php/chet/readline/rltop.html
- **Readline Manual**: https://www.gnu.org/software/bash/manual/html_node/Command-Line-Editing.html
- **Readline Source**: https://git.savannah.gnu.org/cgit/readline.git
- **Readline Documentation**: https://tiswww.cwru.edu/php/chet/readline/readline.html
- **Inputrc Guide**: https://www.gnu.org/software/bash/manual/html_node/Readline-Init-File.html
- **Readline Programming**: https://web.mit.edu/gnu/doc/html/rlman_2.html

## Library Linking

Programs can link against readline in different ways:

### Compile-Time Linking

```bash
# C program compilation with readline
gcc -o myprogram myprogram.c -lreadline

# Check if program uses readline
ldd myprogram | grep readline  # Linux
otool -L myprogram | grep readline  # macOS
```

### Runtime Dependencies

```bash
# Check readline dependencies
apt-cache depends libreadline8  # Debian/Ubuntu
brew deps readline  # macOS
```

## Related Libraries

- **Editline (libedit)**: BSD alternative to readline
- **Linenoise**: Minimal readline alternative
- **GNU History**: Often bundled with readline
- **Ncurses**: Terminal handling library often used with readline

This module ensures the GNU Readline library is properly installed as a system dependency, providing essential line-editing and history capabilities for interactive command-line programs within the devgita development environment.
