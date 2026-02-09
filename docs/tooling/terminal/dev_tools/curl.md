# Curl Module Documentation

## Overview

The Curl module provides HTTP client tool installation and command execution with devgita integration. It follows the standardized devgita app interface while providing curl-specific operations for data transfer, HTTP requests, and URL-based operations.

## App Purpose

Curl is a command-line tool for transferring data with URLs, supporting various protocols including HTTP, HTTPS, FTP, FTPS, SCP, SFTP, TFTP, LDAP, and more. This module ensures curl is properly installed across macOS and Debian/Ubuntu systems and provides high-level operations for common HTTP client tasks and data transfer operations.

## Lifecycle Summary

1. **Installation**: Install curl package via platform package managers (Homebrew/apt)
2. **Configuration**: Curl typically doesn't require separate configuration files - operations are handled via command-line arguments
3. **Execution**: Provide high-level curl operations for HTTP requests, file downloads, and data transfer

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new Curl instance with platform-specific commands            |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install curl                              |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute curl commands     | Runs curl with provided arguments                                     |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
curl := curl.New()
err := curl.Install()
```

- **Purpose**: Standard curl installation
- **Behavior**: Uses `InstallPackage()` to install curl package
- **Use case**: Initial curl installation or explicit reinstall

### ForceInstall()

```go
curl := curl.New()
err := curl.ForceInstall()
```

- **Purpose**: Force curl installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh curl installation or fix corrupted installation

### SoftInstall()

```go
curl := curl.New()
err := curl.SoftInstall()
```

- **Purpose**: Install curl only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing curl installations

### Uninstall()

```go
err := curl.Uninstall()
```

- **Purpose**: Remove curl installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: HTTP client tools are typically managed at the system level

### Update()

```go
err := curl.Update()
```

- **Purpose**: Update curl installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: Curl updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := curl.ForceConfigure()
err := curl.SoftConfigure()
```

- **Purpose**: Apply curl configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: Curl doesn't use traditional config files; operation parameters are passed via command-line arguments

## Execution Methods

### ExecuteCommand()

```go
err := curl.ExecuteCommand("--version")
err := curl.ExecuteCommand("-o", "file.txt", "https://example.com/file.txt")
err := curl.ExecuteCommand("-X", "POST", "-d", "data", "https://api.example.com")
```

- **Purpose**: Execute curl commands with provided arguments
- **Parameters**: Variable arguments passed directly to curl binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### Curl-Specific Operations

The curl CLI provides extensive HTTP client and data transfer capabilities:

#### Basic Requests

```bash
# Get webpage content
curl https://example.com

# Download file
curl -o filename.zip https://example.com/file.zip
curl -O https://example.com/file.zip  # Use remote filename

# Follow redirects
curl -L https://example.com

# Silent mode (no progress bar)
curl -s https://example.com
```

#### HTTP Methods

```bash
# POST request with data
curl -X POST -d "param1=value1&param2=value2" https://api.example.com

# POST with JSON data
curl -X POST -H "Content-Type: application/json" -d '{"key":"value"}' https://api.example.com

# PUT request
curl -X PUT -d "data" https://api.example.com/resource

# DELETE request
curl -X DELETE https://api.example.com/resource
```

#### Headers and Authentication

```bash
# Add custom headers
curl -H "Authorization: Bearer token" https://api.example.com
curl -H "User-Agent: MyApp/1.0" https://example.com

# Basic authentication
curl -u username:password https://example.com

# Include response headers
curl -i https://example.com

# Show only headers
curl -I https://example.com
```

#### Advanced Options

```bash
# Set timeout
curl --connect-timeout 30 https://example.com

# Follow redirects with limit
curl -L --max-redirs 5 https://example.com

# Use proxy
curl --proxy http://proxy:8080 https://example.com

# Ignore SSL errors (development only)
curl -k https://self-signed.example.com

# Save cookies
curl -c cookies.txt https://example.com

# Use cookies
curl -b cookies.txt https://example.com
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **HTTP Operations**: `New()` → `SoftInstall()` → `ExecuteCommand()` with curl arguments
4. **Download Operations**: `New()` → `ExecuteCommand()` with download parameters

## Constants and Paths

### Relevant Constants

- **Package name**: `"curl"` used directly for installation
- Used by all installation methods for consistent package reference

### Configuration Approach

- **No traditional config files**: Curl operations are configured via command-line arguments
- **Runtime configuration**: Parameters passed directly to `ExecuteCommand()`
- **Environment variables**: Curl respects standard environment variables like `HTTP_PROXY`, `HTTPS_PROXY`

## Implementation Notes

- **HTTP Client Nature**: Unlike typical applications, curl is a command-line HTTP client without persistent configuration
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since curl uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since curl doesn't use config files
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as curl updates should be handled by system package managers

## Usage Examples

### Basic HTTP Operations

```go
curl := curl.New()

// Install curl
err := curl.SoftInstall()
if err != nil {
    return err
}

// Get webpage content
err = curl.ExecuteCommand("https://example.com")

// Download file
err = curl.ExecuteCommand("-o", "file.txt", "https://example.com/file.txt")

// POST request with JSON
err = curl.ExecuteCommand("-X", "POST", "-H", "Content-Type: application/json", "-d", `{"key":"value"}`, "https://api.example.com")
```

### Advanced Operations

```go
// Check curl version
err := curl.ExecuteCommand("--version")

// Download with progress bar
err = curl.ExecuteCommand("--progress-bar", "-o", "large-file.zip", "https://example.com/large-file.zip")

// API call with authentication
err = curl.ExecuteCommand("-H", "Authorization: Bearer token", "https://api.example.com/data")
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **SSL Certificate Errors**: Use `-k` flag for self-signed certificates (development only)
3. **Connection Timeouts**: Adjust timeout with `--connect-timeout` and `--max-time`
4. **Proxy Issues**: Configure proxy settings via `--proxy` or environment variables
5. **Large Downloads**: Use `--continue-at -` to resume interrupted downloads

### Platform Considerations

- **macOS**: Installed via Homebrew package manager, system curl also available
- **Linux**: Installed via apt package manager, usually pre-installed
- **SSL Support**: Requires proper SSL/TLS certificate validation
- **Protocol Support**: HTTP, HTTPS, FTP, FTPS, SCP, SFTP, TFTP, LDAP, LDAPS

### Security Notes

- Always validate SSL certificates in production (`-k` flag should only be used for development)
- Be careful with authentication tokens in command history
- Use environment variables or files for sensitive data
- Consider using `--netrc` for credential management

This module provides essential HTTP client capabilities for development workflows, API testing, and data transfer operations within the devgita ecosystem.