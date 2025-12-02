# GitHub CLI (gh) Module Documentation

## Overview

The GitHub CLI (gh) module provides installation and command execution management for GitHub's official command-line tool with devgita integration. It follows the standardized devgita app interface while providing gh-specific operations for GitHub interactions, repository management, pull request operations, issue tracking, and API access.

## App Purpose

GitHub CLI (gh) is the official command-line tool for GitHub that brings pull requests, issues, and other GitHub concepts to the terminal. This module ensures gh is properly installed across macOS (Homebrew) and Debian/Ubuntu (apt) systems and provides high-level operations for common GitHub workflows including authentication, repository operations, pull request management, issue tracking, and release operations.

## Lifecycle Summary

1. **Installation**: Install gh package via platform package managers (Homebrew/apt)
2. **Configuration**: GitHub CLI configuration is typically handled via interactive commands (`gh auth login`) rather than config file copying
3. **Execution**: Provide high-level gh operations for GitHub workflows and API access

## Exported Functions

| Function           | Purpose                   | Behavior                                                             |
| ------------------ | ------------------------- | -------------------------------------------------------------------- |
| `New()`            | Factory method            | Creates new GH instance with platform-specific commands              |
| `Install()`        | Standard installation     | Uses `InstallPackage()` to install gh                                |
| `ForceInstall()`   | Force installation        | Calls `Uninstall()` first (returns error if fails), then `Install()` |
| `SoftInstall()`    | Conditional installation  | Uses `MaybeInstallPackage()` to check before installing              |
| `ForceConfigure()` | Force configuration       | **Not applicable** - returns nil                                     |
| `SoftConfigure()`  | Conditional configuration | **Not applicable** - returns nil                                     |
| `Uninstall()`      | Remove installation       | **Not supported** - returns error                                    |
| `ExecuteCommand()` | Execute gh commands       | Runs gh with provided arguments                                       |
| `Update()`         | Update installation       | **Not implemented** - returns error                                  |

## Installation Methods

### Install()

```go
gh := gh.New()
err := gh.Install()
```

- **Purpose**: Standard GitHub CLI installation
- **Behavior**: Uses `InstallPackage()` to install gh package
- **Use case**: Initial gh installation or explicit reinstall

### ForceInstall()

```go
gh := gh.New()
err := gh.ForceInstall()
```

- **Purpose**: Force GitHub CLI installation regardless of existing state
- **Behavior**: Calls `Uninstall()` first (returns error if it fails), then `Install()`
- **Use case**: Ensure fresh gh installation or fix corrupted installation

### SoftInstall()

```go
gh := gh.New()
err := gh.SoftInstall()
```

- **Purpose**: Install GitHub CLI only if not already present
- **Behavior**: Uses `MaybeInstallPackage()` to check before installing
- **Use case**: Standard installation that respects existing gh installations

### Uninstall()

```go
err := gh.Uninstall()
```

- **Purpose**: Remove GitHub CLI installation
- **Behavior**: **Not supported** - returns error
- **Rationale**: GitHub CLI is typically managed at the system level

### Update()

```go
err := gh.Update()
```

- **Purpose**: Update GitHub CLI installation
- **Behavior**: **Not implemented** - returns error
- **Rationale**: GitHub CLI updates are typically handled by the system package manager

## Configuration Methods

### ForceConfigure() & SoftConfigure()

```go
err := gh.ForceConfigure()
err := gh.SoftConfigure()
```

- **Purpose**: Apply GitHub CLI configuration
- **Behavior**: **Not applicable** - both return nil
- **Rationale**: GitHub CLI configuration is handled via interactive commands (`gh auth login`, `gh config set`) rather than config file copying

## Execution Methods

### ExecuteCommand()

```go
err := gh.ExecuteCommand("--version")
err := gh.ExecuteCommand("auth", "login")
err := gh.ExecuteCommand("pr", "list", "--state", "open")
```

- **Purpose**: Execute gh commands with provided arguments
- **Parameters**: Variable arguments passed directly to gh binary
- **Error handling**: Wraps errors with context from BaseCommand.ExecCommand

### GitHub CLI-Specific Operations

The gh CLI provides extensive GitHub interaction capabilities:

#### Authentication

```bash
# Interactive login
gh auth login

# Check authentication status
gh auth status

# Logout
gh auth logout

# Set authentication token
gh auth login --with-token < token.txt
```

#### Repository Operations

```bash
# Clone repository
gh repo clone owner/repo

# Create new repository
gh repo create my-repo --public

# View repository information
gh repo view owner/repo

# Fork repository
gh repo fork owner/repo

# List repositories
gh repo list
gh repo list owner
```

#### Pull Request Management

```bash
# List pull requests
gh pr list
gh pr list --state open
gh pr list --author @me

# Create pull request
gh pr create --title "Title" --body "Description"

# View pull request
gh pr view 123

# Check out pull request
gh pr checkout 123

# Review pull request
gh pr review 123 --approve
gh pr review 123 --comment --body "Looks good"

# Merge pull request
gh pr merge 123 --squash
```

#### Issue Management

```bash
# List issues
gh issue list
gh issue list --state open
gh issue list --assignee @me

# Create issue
gh issue create --title "Bug title" --body "Description"

# View issue
gh issue view 456

# Close issue
gh issue close 456

# Reopen issue
gh issue reopen 456
```

#### Release Operations

```bash
# List releases
gh release list

# Create release
gh release create v1.0.0 --title "Version 1.0.0" --notes "Release notes"

# View release
gh release view v1.0.0

# Download release assets
gh release download v1.0.0
```

#### GitHub API Access

```bash
# Make authenticated API request
gh api /repos/owner/repo

# POST request
gh api /repos/owner/repo/issues --field title="Title"

# GraphQL query
gh api graphql -f query='query { viewer { login } }'
```

#### Workflow Operations

```bash
# List workflows
gh workflow list

# View workflow runs
gh workflow view

# Run workflow
gh workflow run workflow.yml

# View workflow run logs
gh run view 12345
```

## Expected Function Interactions

1. **Standard Setup**: `New()` → `SoftInstall()` → `SoftConfigure()` (no-op)
2. **Force Setup**: `New()` → `ForceInstall()` → `ForceConfigure()` (no-op)
3. **GitHub Operations**: `New()` → `SoftInstall()` → `ExecuteCommand()` with gh arguments
4. **Authentication Flow**: `New()` → `SoftInstall()` → `ExecuteCommand("auth", "login")`

## Constants and Paths

### Relevant Constants

- **Package name**: Referenced via `constants.GH` (typically "gh")
- Used by all installation methods for consistent package reference

### Configuration Approach

- **No traditional config files**: GitHub CLI configuration is managed via interactive commands
- **Authentication**: Handled via `gh auth login` command
- **Configuration**: Managed via `gh config set <key> <value>` commands
- **Config location**: GitHub CLI stores configuration in `~/.config/gh/` (managed by gh itself)

## Implementation Notes

- **CLI Tool Nature**: Unlike typical applications, gh is a command-line interface without traditional config file templates
- **ForceInstall Logic**: Calls `Uninstall()` first and returns the error if it fails since gh uninstall is not supported
- **Configuration Strategy**: Returns nil for both `ForceConfigure()` and `SoftConfigure()` since gh uses interactive configuration
- **Error Handling**: All methods return errors that should be checked by callers
- **Platform Independence**: Uses command interface abstraction for cross-platform compatibility
- **Update Method**: Not implemented as gh updates should be handled by system package managers

## Usage Examples

### Basic GitHub Operations

```go
gh := gh.New()

// Install gh
err := gh.SoftInstall()
if err != nil {
    return err
}

// Check version
err = gh.ExecuteCommand("--version")

// Authenticate
err = gh.ExecuteCommand("auth", "login")

// List pull requests
err = gh.ExecuteCommand("pr", "list")

// Clone repository
err = gh.ExecuteCommand("repo", "clone", "owner/repo")
```

### Advanced Operations

```go
// Create pull request
err := gh.ExecuteCommand("pr", "create", "--title", "Feature", "--body", "Description")

// Check authentication status
err = gh.ExecuteCommand("auth", "status")

// Make API request
err = gh.ExecuteCommand("api", "/repos/owner/repo/issues")

// View workflow runs
err = gh.ExecuteCommand("run", "list")
```

## Troubleshooting

### Common Issues

1. **Installation Fails**: Ensure package manager is available and updated
2. **Authentication Errors**: Run `gh auth login` to authenticate with GitHub
3. **Permission Issues**: Verify GitHub account has required permissions for operations
4. **API Rate Limits**: GitHub API has rate limits; use authenticated requests
5. **Commands Don't Work**: Verify gh is installed and accessible in PATH

### Platform Considerations

- **macOS**: Installed via Homebrew package manager
- **Linux**: Installed via apt package manager (GitHub official repository)
- **Authentication**: Uses OAuth for secure authentication
- **Token Storage**: Tokens stored securely in system keychain

### Authentication Setup

GitHub CLI supports multiple authentication methods:

- **Browser-based OAuth**: Interactive login via web browser (recommended)
- **Personal access token**: Manual token input for automation
- **SSH key authentication**: Git operations use SSH keys

### Best Practices

- **Authenticate after installation**: Run `gh auth login` immediately after install
- **Use verbose output for debugging**: Add `--verbose` flag to gh commands
- **Check authentication status**: Use `gh auth status` to verify authentication
- **Use aliases for common operations**: Configure via `gh alias set`
- **Leverage extensions**: Install gh extensions for additional functionality

### Security Notes

- Tokens are stored securely in system keychain
- Use `gh auth logout` to revoke access when needed
- Avoid exposing authentication tokens in scripts
- Consider using environment variables for CI/CD automation

## Integration with Devgita

GitHub CLI integrates with devgita's terminal category:

- **Installation**: Installed as part of terminal tools setup
- **Configuration**: Interactive authentication via `gh auth login`
- **Usage**: Available system-wide after installation
- **Updates**: Managed through system package manager

## External References

- **GitHub CLI Documentation**: https://cli.github.com/manual/
- **GitHub CLI Repository**: https://github.com/cli/cli
- **Installation Guide**: https://github.com/cli/cli#installation
- **Command Reference**: https://cli.github.com/manual/gh
- **GitHub API**: https://docs.github.com/en/rest

This module provides essential GitHub integration capabilities for development workflows, repository management, and CI/CD automation within the devgita ecosystem.
