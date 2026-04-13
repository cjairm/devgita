# Developer Quickstart: Debian/Ubuntu Package Installation

**Feature**: 002-debian-package-fixes  
**Date**: 2026-04-09  
**Audience**: Developers implementing or extending Debian/Ubuntu package installation

---

## Prerequisites

- Go 1.21+ installed
- Debian 12+ (Bookworm) or Ubuntu 24+ VM for testing (optional but recommended)
- Git configured
- Familiarity with devgita project structure

---

## Quick Setup

### 1. Clone and Build

```bash
# Clone repository (if not already)
cd /Users/jair.mendez/Documents/projects/devgita

# Checkout feature branch
git checkout 002-debian-package-fixes

# Build
go build -o devgita main.go

# Verify build
./devgita --version
```

### 2. Run Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./pkg/constants/
go test ./pkg/downloader/
go test ./pkg/apt/
go test ./internal/commands/

# Run with verbose output
go test -v ./internal/commands/

# Run with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 3. Test on Debian/Ubuntu (Optional)

```bash
# Create Debian 12 VM (using multipass, virtualbox, or docker)
multipass launch debian -n devgita-test

# SSH into VM
multipass shell devgita-test

# Copy devgita binary
multipass transfer ./devgita devgita-test:/home/ubuntu/

# Install in VM
./devgita install --only terminal
```

---

## Project Structure

### Key Directories

```
devgita/
├── pkg/
│   ├── constants/
│   │   ├── constants.go           # Existing constants
│   │   └── package_mappings.go    # NEW: Library name mappings
│   ├── downloader/
│   │   └── retry.go               # NEW: Retry logic with backoff
│   ├── apt/
│   │   └── ppa.go                 # NEW: PPA management
│   └── logger/                    # Existing
│
├── internal/
│   ├── commands/
│   │   ├── base.go                # Existing: BaseCommandExecutor interface
│   │   ├── macos.go               # Existing: DO NOT MODIFY
│   │   ├── debian.go              # MODIFY: Add strategy selection
│   │   ├── debian_strategies.go   # NEW: Strategy implementations
│   │   ├── factory.go             # Existing: Platform detection
│   │   └── mock.go                # Existing: Test mocks
│   └── apps/
│       ├── neovim/
│       │   ├── neovim.go          # Existing
│       │   └── neovim_debian.go   # NEW: Debian-specific install
│       └── [other apps]/
│
└── specs/002-debian-package-fixes/
    ├── plan.md                    # This feature's implementation plan
    ├── research.md                # Research findings
    ├── data-model.md              # Data structures
    ├── quickstart.md              # This file
    └── contracts/                 # Interface contracts
```

---

## Common Tasks

### Adding a New Package Mapping

**File**: `pkg/constants/package_mappings.go`

```go
// 1. Add constant to pkg/constants/constants.go (if not exists)
const (
    // ... existing
    Vips = "vips"  // NEW
)

// 2. Add mapping to package_mappings.go
var PackageMappings = map[string]PackageMapping{
    // ... existing mappings
    Vips: {
        MacOS:  "vips",
        Debian: "libvips",
    },
}

// 3. Add test
func TestGetDebianPackageName_Vips(t *testing.T) {
    result := GetDebianPackageName(Vips)
    expected := "libvips"
    if result != expected {
        t.Errorf("Expected %s, got %s", expected, result)
    }
}
```

### Adding a New Installation Strategy

**Scenario**: Install `fastfetch` via PPA

**File**: `internal/commands/debian.go`

```go
// Modify getInstallationStrategy() method
func (d *DebianCommand) getInstallationStrategy(packageName string) InstallationStrategy {
    switch packageName {
    // ... existing cases
    
    case constants.Fastfetch:  // NEW
        return &PPAStrategy{
            cmd: d,
            ppaConfig: apt.PPAConfig{
                Name:         "fastfetch",
                KeyURL:       "https://ppa.launchpadcontent.net/zhangsongcui3371/fastfetch/ubuntu/pool/main/f/fastfetch/",
                RepoURL:      "ppa:zhangsongcui3371/fastfetch",
                Distribution: "$(lsb_release -sc)",
                Component:    "main",
            },
        }
    
    default:
        return &AptStrategy{cmd: d}
    }
}
```

**Test**:

```go
func TestFastfetchPPAStrategy(t *testing.T) {
    mockCmd := NewMockDebianCommand()
    debian := &DebianCommand{Base: mockCmd}
    
    strategy := debian.getInstallationStrategy(constants.Fastfetch)
    
    // Verify it's a PPA strategy
    if _, ok := strategy.(*PPAStrategy); !ok {
        t.Errorf("Expected PPAStrategy, got %T", strategy)
    }
}
```

### Adding a GitHub Binary Installation

**Scenario**: Install `lazydocker` from GitHub releases

**File**: `internal/apps/lazydocker/lazydocker_debian.go` (NEW)

```go
package lazydocker

import (
    "context"
    "fmt"
    "os"
    "os/exec"
    
    "github.com/cjairm/devgita/pkg/downloader"
)

func (l *LazyDocker) InstallDebianFromGitHub() error {
    const (
        downloadURL = "https://github.com/jesseduffield/lazydocker/releases/latest/download/lazydocker_Linux_x86_64.tar.gz"
        tmpFile     = "/tmp/lazydocker.tar.gz"
        installPath = "/usr/local/bin/lazydocker"
    )
    
    // 1. Download with retry
    ctx := context.Background()
    config := downloader.DefaultRetryConfig()
    if err := downloader.DownloadFileWithRetry(ctx, downloadURL, tmpFile, config); err != nil {
        return fmt.Errorf("download failed: %w", err)
    }
    defer os.Remove(tmpFile)
    
    // 2. Extract tar.gz
    cmd := exec.Command("tar", "-xzf", tmpFile, "-C", "/tmp")
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("extraction failed: %w", err)
    }
    
    // 3. Install binary
    cmd = exec.Command("sudo", "install", "-m", "755", "/tmp/lazydocker", installPath)
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("installation failed: %w", err)
    }
    
    // 4. Cleanup
    os.Remove("/tmp/lazydocker")
    
    return nil
}
```

**Test**:

```go
func TestInstallDebianFromGitHub(t *testing.T) {
    // Use mock downloader
    mockDownloader := &MockDownloader{
        DownloadFunc: func(ctx context.Context, url, dest string, config RetryConfig) error {
            // Create fake file
            os.WriteFile(dest, []byte("fake content"), 0644)
            return nil
        },
    }
    
    // Test installation logic
    // ... test implementation
}
```

---

## Testing Patterns

### Unit Testing with Mocks

```go
func TestDebianInstallPackage_WithMapping(t *testing.T) {
    // Setup
    mockBase := commands.NewMockBaseCommand()
    mockBase.SetExecCommandResult("", "", nil)  // Success
    
    debian := &DebianCommand{Base: mockBase}
    
    // Execute
    err := debian.InstallPackage(constants.Gdbm)
    
    // Verify
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    lastCall := mockBase.GetLastExecCommandCall()
    
    // Should install "libgdbm-dev" not "gdbm"
    if !contains(lastCall.Args, "libgdbm-dev") {
        t.Errorf("Expected libgdbm-dev, got args: %v", lastCall.Args)
    }
}
```

### Testing Retry Logic

```go
func TestDownloadWithRetry_Success(t *testing.T) {
    attempt := 0
    mockHTTP := &MockHTTPClient{
        DoFunc: func(req *http.Request) (*http.Response, error) {
            attempt++
            if attempt < 2 {
                return nil, fmt.Errorf("timeout")  // Fail first
            }
            return &http.Response{
                StatusCode: 200,
                Body:       io.NopCloser(bytes.NewReader([]byte("content"))),
            }, nil  // Succeed second
        },
    }
    
    config := RetryConfig{MaxRetries: 3, InitialWait: 10 * time.Millisecond}
    err := DownloadFileWithRetry(context.Background(), "http://test.com/file", "/tmp/test", config)
    
    if err != nil {
        t.Fatalf("Expected success after retry, got %v", err)
    }
    
    if attempt != 2 {
        t.Errorf("Expected 2 attempts, got %d", attempt)
    }
}
```

### Integration Testing (Debian VM)

```bash
# Create test script
cat > test_debian.sh <<'EOF'
#!/bin/bash
set -e

echo "Testing terminal installation..."
./devgita install --only terminal

echo "Verifying installations..."
command -v curl >/dev/null 2>&1 || { echo "curl not found"; exit 1; }
command -v git >/dev/null 2>&1 || { echo "git not found"; exit 1; }
command -v nvim >/dev/null 2>&1 || { echo "nvim not found"; exit 1; }

echo "All tests passed!"
EOF

chmod +x test_debian.sh

# Run in VM
multipass transfer test_debian.sh devgita-test:/home/ubuntu/
multipass exec devgita-test -- ./test_debian.sh
```

---

## Debugging

### Enable Verbose Logging

```go
// In main.go or test
logger.Init(true)  // Enable verbose logging
```

### Check Installation State

```bash
# View global config
cat ~/.config/devgita/global_config.yaml

# Check installed packages
dpkg -l | grep libgdbm-dev
dpkg -l | grep mise

# Check binaries in PATH
command -v nvim
which lazygit
```

### Test Retry Logic Manually

```go
func main() {
    ctx := context.Background()
    config := downloader.RetryConfig{
        MaxRetries:  3,
        InitialWait: 1 * time.Second,
        Multiplier:  2.0,
        Jitter:      0.2,
    }
    
    url := "https://github.com/neovim/neovim/releases/download/stable/nvim-linux-x86_64.tar.gz"
    dest := "/tmp/nvim.tar.gz"
    
    err := downloader.DownloadFileWithRetry(ctx, url, dest, config)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("Download successful!")
}
```

---

## Reference Materials

### Omakub Patterns

```bash
# Browse omakub installation scripts for patterns
ls /Users/jair.mendez/Documents/projects/devgita/omakub/install/

# Example: Mise installation
cat /Users/jair.mendez/Documents/projects/devgita/omakub/install/terminal/mise.sh

# Example: Neovim installation
cat /Users/jair.mendez/Documents/projects/devgita/omakub/install/app-neovim.sh
```

### Documentation

- **Specification**: `specs/002-debian-package-fixes/spec.md`
- **Research**: `specs/002-debian-package-fixes/research.md`
- **Data Model**: `specs/002-debian-package-fixes/data-model.md`
- **Contracts**: `specs/002-debian-package-fixes/contracts/`
- **Testing Guide**: `docs/guides/testing-patterns.md`
- **Constitution**: `.specify/memory/constitution.md`

---

## Common Issues

### Issue: Package mapping not working

**Symptom**: apt tries to install "gdbm" instead of "libgdbm-dev"

**Fix**:
1. Verify mapping exists in `pkg/constants/package_mappings.go`
2. Check constant name matches exactly
3. Ensure `GetDebianPackageName()` is called in InstallPackage

```go
// Wrong
d.installWithApt(packageName)

// Correct
debianName := constants.GetDebianPackageName(packageName)
d.installWithApt(debianName)
```

### Issue: Retry not happening

**Symptom**: Download fails immediately without retrying

**Fix**:
1. Check error is marked as retryable
2. Verify RetryConfig is properly initialized
3. Ensure DownloadFileWithRetry is used (not downloadFile directly)

```go
// Wrong - no retry
downloadFile(ctx, url, dest)

// Correct - with retry
config := downloader.DefaultRetryConfig()
downloader.DownloadFileWithRetry(ctx, url, dest, config)
```

### Issue: Tests executing real commands

**Symptom**: Tests fail on macOS or without sudo

**Fix**:
1. Use MockBaseCommand instead of real command
2. Verify no real ExecCommand calls
3. Use testutil.VerifyNoRealCommands

```go
// Wrong
func TestInstall(t *testing.T) {
    debian := &DebianCommand{Base: commands.NewBaseCommand()}  // REAL
    debian.InstallPackage("curl")  // Executes real apt command!
}

// Correct
func TestInstall(t *testing.T) {
    mockBase := commands.NewMockBaseCommand()  // MOCK
    debian := &DebianCommand{Base: mockBase}
    debian.InstallPackage("curl")  // No real command executed
}
```

---

## Next Steps

1. **Implement**: Start with foundation (package mappings, retry logic, PPA management)
2. **Test**: Write unit tests for each component
3. **Integrate**: Wire up strategies in DebianCommand
4. **Validate**: Test on Debian/Ubuntu VM
5. **Audit**: Run Constitution compliance check

---

## Getting Help

- **Code questions**: Review `docs/` directory and existing code patterns
- **Testing questions**: See `docs/guides/testing-patterns.md`
- **Architecture questions**: Review `docs/project-overview.md`
- **Constitution questions**: See `.specify/memory/constitution.md`

---

This quickstart guide provides everything needed to start implementing or extending Debian/Ubuntu package installation support in devgita.
