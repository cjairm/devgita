# Command Interface Contract

**Feature**: 002-debian-package-fixes  
**Date**: 2026-04-09  
**Context**: This document defines the Go interface contract for platform-specific command execution.

---

## BaseCommandExecutor Interface

**Location**: `internal/commands/base.go`  
**Purpose**: Abstracts platform-specific command execution for testability and cross-platform support

### Interface Definition

```go
// BaseCommandExecutor defines the contract for platform-specific command execution
type BaseCommandExecutor interface {
    ExecCommand(cmd CommandParams) (stdout string, stderr string, err error)
    IsPackageInstalled(packageName string) (bool, error)
    // Other existing methods unchanged...
}
```

**Implementations**:
- `MacOSCommand` - Homebrew-based implementation (**UNCHANGED**)
- `DebianCommand` - apt/PPA/GitHub binary implementation (**EXTENDED**)
- `MockBaseCommand` - Test mock implementation

### Contract Guarantees

1. **Platform Isolation**: macOS and Debian implementations MUST NOT share code
2. **Interface Stability**: Method signatures MUST NOT change (backward compatibility)
3. **Error Handling**: All errors MUST be wrapped with context
4. **Idempotency**: All operations MUST be safe to re-run

---

## Command Interface

**Location**: `internal/commands/base.go`  
**Purpose**: High-level package management operations

### Interface Definition

```go
type Command interface {
    InstallPackage(packageName string) error
    MaybeInstallPackage(packageName string) error
    IsPackageInstalled(packageName string) (bool, error)
    InstallDesktopApp(packageName string) error
    MaybeInstallDesktopApp(packageName string) error
    // Other existing methods...
}
```

**Implementations**:
- `MacOSCommand` (**UNCHANGED**)
- `DebianCommand` (**EXTENDED**)
- `MockCommand` (test mock)

### Platform-Specific Behavior

#### macOS (Unchanged)

```go
// internal/commands/macos.go
func (m *MacOSCommand) InstallPackage(packageName string) error {
    // Direct Homebrew installation
    cmd := CommandParams{
        Command: "brew",
        Args:    []string{"install", packageName},
    }
    _, _, err := m.ExecCommand(cmd)
    return err
}
```

**Guarantees**:
- Uses Homebrew for all packages
- No package name translation
- No retry logic (Homebrew handles this)
- **Code must not be modified for this feature**

#### Debian/Ubuntu (Extended)

```go
// internal/commands/debian.go
func (d *DebianCommand) InstallPackage(packageName string) error {
    // 1. Select installation strategy based on package
    strategy := d.getInstallationStrategy(packageName)
    
    // 2. Execute strategy
    return strategy.Install(packageName)
}
```

**New Behavior**:
- Uses strategy pattern for different installation methods
- Package name translation for libraries
- Retry logic for downloads (3 attempts)
- Component-level failure recovery

---

## InstallationStrategy Interface (NEW)

**Location**: `internal/commands/debian_strategies.go`  
**Purpose**: Defines contract for different installation methods on Debian/Ubuntu

### Interface Definition

```go
// InstallationStrategy defines the contract for package installation methods
type InstallationStrategy interface {
    Install(packageName string) error
    IsInstalled(packageName string) (bool, error)
}
```

### Implementations

#### AptStrategy

```go
type AptStrategy struct {
    cmd *DebianCommand
}

// Install installs via apt with package name translation
func (s *AptStrategy) Install(packageName string) error {
    debianName := constants.GetDebianPackageName(packageName)
    return s.cmd.installWithApt(debianName)
}

// IsInstalled checks if package is installed via dpkg
func (s *AptStrategy) IsInstalled(packageName string) (bool, error) {
    debianName := constants.GetDebianPackageName(packageName)
    return s.cmd.IsPackageInstalled(debianName)
}
```

**Use Cases**: curl, git, libgdbm-dev, libffi-dev

#### PPAStrategy

```go
type PPAStrategy struct {
    cmd       *DebianCommand
    ppaConfig apt.PPAConfig
}

// Install adds PPA first, then installs package via apt
func (s *PPAStrategy) Install(packageName string) error {
    manager := apt.NewPPAManager()
    if err := manager.AddPPA(s.ppaConfig); err != nil {
        return fmt.Errorf("failed to add PPA: %w", err)
    }
    return s.cmd.installWithApt(packageName)
}

// IsInstalled checks if package is installed via dpkg
func (s *PPAStrategy) IsInstalled(packageName string) (bool, error) {
    return s.cmd.IsPackageInstalled(packageName)
}
```

**Use Cases**: mise, fastfetch

#### GitHubBinaryStrategy

```go
type GitHubBinaryStrategy struct {
    cmd         *DebianCommand
    downloadURL string
    installPath string
    binaryName  string
}

// Install downloads binary from GitHub with retry, extracts, installs
func (s *GitHubBinaryStrategy) Install(packageName string) error {
    ctx := context.Background()
    config := downloader.DefaultRetryConfig()
    
    // Download with retry (3 attempts: 1s, 2s, 4s backoff)
    tmpFile := "/tmp/" + packageName + ".tar.gz"
    if err := downloader.DownloadFileWithRetry(ctx, s.downloadURL, tmpFile, config); err != nil {
        return fmt.Errorf("download failed: %w", err)
    }
    defer os.Remove(tmpFile)
    
    // Extract and install binary
    // ... implementation details
    return nil
}

// IsInstalled checks if binary exists in PATH
func (s *GitHubBinaryStrategy) IsInstalled(packageName string) (bool, error) {
    _, err := exec.LookPath(s.binaryName)
    return err == nil, nil
}
```

**Use Cases**: neovim, lazygit, lazydocker

#### GitCloneStrategy

```go
type GitCloneStrategy struct {
    cmd         *DebianCommand
    repoURL     string
    installPath string
}

// Install clones Git repository to specified path
func (s *GitCloneStrategy) Install(packageName string) error {
    cmd := exec.Command("git", "clone", s.repoURL, s.installPath)
    return cmd.Run()
}

// IsInstalled checks if directory exists
func (s *GitCloneStrategy) IsInstalled(packageName string) (bool, error) {
    _, err := os.Stat(s.installPath)
    return err == nil, nil
}
```

**Use Cases**: powerlevel10k

---

## Strategy Selection Contract

```go
// internal/commands/debian.go

// getInstallationStrategy returns the appropriate strategy for a package
func (d *DebianCommand) getInstallationStrategy(packageName string) InstallationStrategy {
    switch packageName {
    case constants.Mise:
        return &PPAStrategy{
            cmd: d,
            ppaConfig: apt.PPAConfig{
                Name:         "mise",
                KeyURL:       "https://mise.jdx.dev/gpg-key.pub",
                RepoURL:      "https://mise.jdx.dev/deb",
                Distribution: "stable",
                Component:    "main",
            },
        }
    
    case constants.Neovim:
        return &GitHubBinaryStrategy{
            cmd:         d,
            downloadURL: "https://github.com/neovim/neovim/releases/download/stable/nvim-linux-x86_64.tar.gz",
            binaryName:  "nvim",
            installPath: "/usr/local/bin/nvim",
        }
    
    case constants.LazyGit:
        return &GitHubBinaryStrategy{
            cmd:         d,
            downloadURL: "https://github.com/jesseduffield/lazygit/releases/latest/download/lazygit_Linux_x86_64.tar.gz",
            binaryName:  "lazygit",
            installPath: "/usr/local/bin/lazygit",
        }
    
    // ... more cases
    
    default:
        // Default to apt strategy with package name translation
        return &AptStrategy{cmd: d}
    }
}
```

**Contract Guarantees**:
1. **Always returns a strategy**: Never returns nil
2. **Default fallback**: Unknown packages use AptStrategy
3. **No side effects**: Selection is pure function (no state changes)
4. **Thread-safe**: Can be called concurrently

---

## Error Handling Contract

### Error Types

```go
// Retryable errors (handled by retry logic)
- net.Error with Timeout() == true
- net.Error with Temporary() == true
- HTTP 429 (Too Many Requests)
- HTTP 502, 503, 504 (Server errors)

// Non-retryable errors (fail immediately)
- HTTP 404 (Not Found)
- HTTP 401, 403 (Auth errors)
- Invalid URL
- File system errors
- Package not found in apt
```

### Error Wrapping

```go
// All errors MUST be wrapped with context
if err := strategy.Install(pkg); err != nil {
    return fmt.Errorf("failed to install %s: %w", pkg, err)
}
```

### Error Propagation

```go
Component Installation
    ↓ (error)
Strategy.Install()
    ↓ (wrapped error)
DebianCommand.InstallPackage()
    ↓ (logged + tracked)
Coordinator (terminal/languages/databases)
    ↓ (continue with next package)
Installation Summary
```

---

## Testing Contract

### Mock Interface

```go
// internal/commands/mock.go

type MockBaseCommand struct {
    ExecCommandResult      func(CommandParams) (string, string, error)
    IsPackageInstalledFunc func(string) (bool, error)
    CallCount              int
    LastCommand            CommandParams
}

func (m *MockBaseCommand) ExecCommand(cmd CommandParams) (string, string, error) {
    m.CallCount++
    m.LastCommand = cmd
    if m.ExecCommandResult != nil {
        return m.ExecCommandResult(cmd)
    }
    return "", "", nil
}

func (m *MockBaseCommand) IsPackageInstalled(pkg string) (bool, error) {
    if m.IsPackageInstalledFunc != nil {
        return m.IsPackageInstalledFunc(pkg)
    }
    return false, nil
}
```

### Testing Pattern

```go
func TestDebianInstallPackage(t *testing.T) {
    mockBase := commands.NewMockBaseCommand()
    mockBase.ExecCommandResult = func(cmd CommandParams) (string, string, error) {
        // Verify correct command
        if cmd.Command != "apt" {
            t.Errorf("Expected apt, got %s", cmd.Command)
        }
        return "", "", nil
    }
    
    debian := &DebianCommand{Base: mockBase}
    err := debian.InstallPackage("curl")
    
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    if mockBase.CallCount != 1 {
        t.Errorf("Expected 1 call, got %d", mockBase.CallCount)
    }
}
```

**Testing Guarantees**:
1. **No real commands**: Tests MUST NOT execute real system commands
2. **Isolated**: Each test uses fresh mock instances
3. **Verifiable**: All mock calls are tracked and verifiable

---

## Backward Compatibility Guarantees

### Interface Stability

✅ **No changes to existing interface methods**
```go
// BaseCommandExecutor - NO changes
ExecCommand(cmd CommandParams) (stdout string, stderr string, err error)
IsPackageInstalled(packageName string) (bool, error)

// Command - NO changes
InstallPackage(packageName string) error
MaybeInstallPackage(packageName string) error
```

### macOS Implementation

✅ **macOS code completely untouched**
- `internal/commands/macos.go` - NO modifications
- All macOS tests pass unchanged
- Homebrew workflow identical

### Factory Pattern

✅ **Platform detection unchanged**
```go
// internal/commands/factory.go (existing)
func NewCommand() Command {
    if runtime.GOOS == "darwin" {
        return NewMacOSCommand()
    }
    return NewDebianCommand()
}
```

---

## Performance Contract

| Operation | Max Time | Notes |
|-----------|----------|-------|
| Strategy selection | < 1ms | Switch statement, O(1) |
| Package name lookup | < 1ms | Map lookup, O(1) |
| Single download attempt | 5 minutes | HTTP client timeout |
| Total retry window | 7 seconds | 3 retries: 1s + 2s + 4s |
| PPA addition | < 30 seconds | Network dependent |
| Apt install | Variable | Package size dependent |

---

## Concurrency Contract

**Single-threaded execution**:
- No concurrent package installations within a category
- No locking required for GlobalConfig (single writer)
- Retry logic is sequential (one attempt at a time)

**Future consideration**:
- Parallel installation across categories (terminal + languages simultaneously)
- Would require GlobalConfig locking mechanism

---

This interface contract ensures type-safe, platform-independent package management while maintaining complete isolation between macOS and Debian/Ubuntu implementations.
