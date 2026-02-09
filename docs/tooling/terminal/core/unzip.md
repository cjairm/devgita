# Unzip Module Documentation

## Overview

The Unzip module provides archive extraction tool installation and command execution with devgita integration. It follows the standardized devgita app interface while providing unzip-specific operations for extracting files from ZIP archives.

## App Purpose

Unzip is a command-line utility for listing, testing, and extracting files from ZIP archives. This module ensures unzip is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems and provides high-level operations for common archive extraction tasks.

## Lifecycle Summary

1. **Installation**: Install unzip package via platform package managers (Homebrew/apt)
2. **Configuration**: Unzip typically doesn't require separate configuration files - operations are handled via command-line arguments
3. **Execution**: Provide high-level unzip operations for archive extraction and management

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Unzip instance with platform-specific commands           |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install unzip                             |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute unzip commands    | Runs unzip with provided arguments                                   |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
unzip := unzip.New()
err := unzip.Install()
```

- **Purpose**: Standard unzip installation
- **Behavior**: Uses `InstallPackage()` to install unzip package
- **Use case**: Initial unzip installation or explicit reinstall

### ForceInstall()

```go
unzip := unzip.New()
err := unzip.ForceInstall()
```

- **Purpose**: Force unzip installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh unzip installation or fix corrupted installation

### SoftInstall()

```go
unzip := unzip.New()
err := unzip.SoftInstall()
```

- **Purpose**: Install unzip only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing unzip installations

### Uninstall()

```go
err := unzip.Uninstall()
```

- **Purpose**: Remove unzip installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: Archive extraction tools are typically managed at the system level

### Update()

```go
err := unzip.Update()
```

- **Purpose**: Update unzip installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Unzip updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := unzip.ForceConfigure()
err := unzip.SoftConfigure()
```

- **Purpose**: Apply unzip configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: Unzip doesn't use traditional config files; operation parameters are passed via command-line arguments

## Execution Methods

### ExecuteCommand()

```go
err := unzip.ExecuteCommand("archive.zip")
err := unzip.ExecuteCommand("-l", "archive.zip")
err := unzip.ExecuteCommand("-d", "/target/path", "archive.zip")
```

- **Purpose**: Execute unzip commands with provided arguments
- **Parameters**: Variable arguments passed directly to unzip binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### Unzip-Specific Operations

The unzip CLI provides extensive archive extraction and management capabilities:

#### Basic Extraction

```bash
# Extract to current directory
unzip archive.zip

# Extract to specific directory
unzip -d /target/path archive.zip

# Extract specific file
unzip archive.zip file.txt

# Extract with pattern matching
unzip archive.zip "*.txt"
```

#### Archive Information

```bash
# List archive contents
unzip -l archive.zip

# Verbose listing
unzip -v archive.zip

# Test archive integrity
unzip -t archive.zip

# Show archive comment
unzip -z archive.zip
```

#### Extraction Options

```bash
# Quiet mode (suppress output)
unzip -q archive.zip

# Overwrite existing files without prompting
unzip -o archive.zip

# Never overwrite existing files
unzip -n archive.zip

# Update existing files only
unzip -u archive.zip

# Freshen existing files (no new files)
unzip -f archive.zip
```

#### Case Handling

```bash
# Convert filenames to lowercase
unzip -L archive.zip

# Extract with case-insensitive matching
unzip -C archive.zip "*.TXT"
```

#### Password Protection

```bash
# Extract password-protected archive
unzip -P password archive.zip

# Prompt for password
unzip archive.zip
# (will prompt if password is required)
```

#### Multiple Archives

```bash
# Extract multiple archives
unzip "*.zip"

# Extract multiple archives to separate directories
for file in *.zip; do
    unzip "$file" -d "${file%.zip}"
done
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **Archive Operations**: `New()` → `SoftInstall()` → `ExecuteCommand()` with unzip arguments
4. **Extraction Operations**: `New()` → `ExecuteCommand()` with extraction parameters

## Constants and Paths

### Relevant Constants

- **Package name**: `"unzip"` used directly for installation
- Used by all installation methods for consistent package reference

### Configuration Approach

- **No traditional config files**: Unzip operations are configured via command-line arguments
- **Runtime configuration**: Parameters passed directly to `ExecuteCommand()`
- **Environment variables**: Unzip respects standard environment variables like `UNZIP` for default options

## Implementation Notes

- **Archive Extraction Nature**: Unlike typical applications, unzip is a command-line extraction utility without persistent configuration
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since unzip uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since unzip doesn't use config files
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as unzip updates should be handled by system package managers

## Usage Examples

### Basic Archive Operations

```go
unzip := unzip.New()

// Install unzip
err := unzip.SoftInstall()
if err != nil {
    return err
}

// Extract archive to current directory
err = unzip.ExecuteCommand("archive.zip")

// List archive contents
err = unzip.ExecuteCommand("-l", "archive.zip")

// Extract to specific directory
err = unzip.ExecuteCommand("-d", "/tmp/extracted", "archive.zip")
```

### Advanced Operations

```go
// Test archive integrity
err := unzip.ExecuteCommand("-t", "archive.zip")

// Extract with overwrite
err = unzip.ExecuteCommand("-o", "archive.zip")

// Quiet extraction
err = unzip.ExecuteCommand("-q", "archive.zip")

// Extract specific files
err = unzip.ExecuteCommand("archive.zip", "file1.txt", "file2.txt")

// Extract with pattern
err = unzip.ExecuteCommand("archive.zip", "*.txt")
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Password Errors**: Use `-P password` flag or allow interactive password prompt
3. **Overwrite Conflicts**: Use `-o` (overwrite), `-n` (never overwrite), or `-u` (update) flags
4. **Permission Issues**: Check write permissions in target directory
5. **Corrupted Archives**: Use `-t` flag to test archive integrity before extraction

### Platform Considerations

- **macOS**: Installed via Homebrew package manager, system unzip also available
- **Linux**: Installed via apt package manager, usually pre-installed
- **ZIP Format Support**: Supports ZIP archives created by PKZIP, Info-ZIP, and WinZip
- **Encoding Issues**: May have issues with non-ASCII filenames depending on system locale

### Exit Codes

Unzip uses standard exit codes to indicate success or failure:

- **0**: Normal; no errors or warnings detected
- **1**: One or more warning errors were encountered, but processing completed successfully
- **2**: A generic error in the zipfile format was detected
- **3**: A severe error in the zipfile format was detected
- **4**: Unable to allocate memory
- **5**: Unable to allocate memory or unable to obtain a tty to read the decryption password(s)
- **6**: Unable to allocate memory during decompression to disk
- **7**: Unable to allocate memory during in-memory decompression
- **8**: Unused
- **9**: The specified zipfiles were not found
- **10**: Invalid options were specified on the command line
- **11**: No matching files were found

### Best Practices

- **Test before extraction**: Use `-t` flag to verify archive integrity
- **Use quiet mode for scripts**: `-q` flag suppresses output for automated operations
- **Specify target directory**: Use `-d` flag to control extraction location
- **Handle overwrites explicitly**: Choose `-o`, `-n`, `-u`, or `-f` based on requirements
- **Verify permissions**: Ensure write access to target directory before extraction
- **Check available space**: Verify sufficient disk space for extracted files
